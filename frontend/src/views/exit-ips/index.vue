<script setup lang="ts">
import { computed, h, onMounted, reactive, ref, watch } from 'vue';
import dayjs from 'dayjs';
import { useWindowSize } from '@vueuse/core';
import { useRouter } from 'vue-router';
import {
  NDatePicker,
  NDrawer,
  NDrawerContent,
  NDropdown,
  NInput,
  NInputNumber,
  NPagination,
  NSelect,
  NSpin,
  NSwitch
} from 'naive-ui';
import {
  createExitIp,
  deleteExitIp,
  fetchExitIpDetail,
  fetchExitIps,
  fetchExitIpUsers,
  fetchNodes,
  updateExitIp
} from '@/service/api';
import SvgIcon from '@/components/custom/svg-icon.vue';
import DataFreshness from '@/components/project/data-freshness.vue';

defineOptions({ name: 'ExitIpManagement' });

type Tone = 'ok' | 'warn' | 'danger' | 'neutral' | 'brand';
/** 列表接口已返回 notes；为空时用派生说明兜底。 */
type ExitIpRow = Api.Central.ExitIpSummary & { notes?: string | null };

const loading = ref(false);
const errorMessage = ref('');
const rows = ref<ExitIpRow[]>([]);
const total = ref(0);
const stats = ref<Api.Central.ExitIpListStats | null>(null);
const dataAt = ref('');
const filters = reactive({ page: 1, page_size: 50, keyword: '', status: 'all', country: 'all' });
const nodes = ref<Api.Central.NodeSummary[]>([]);
const nodesLoading = ref(false);
const nodeError = ref('');
const modalVisible = ref(false);
const modalLoading = ref(false);
const saving = ref(false);
const deletingExitIpId = ref('');
const formError = ref('');
const editingExitIpId = ref<string | null>(null);
const exitIpForm = reactive({
  address: '',
  sourceType: 'node' as 'node' | 's5',
  ownerNodeId: '',
  region: '',
  family: 4 as 4 | 6,
  provider: '',
  monthlyCost: 0 as number | null,
  enabled: true,
  validFrom: null as string | null,
  validTo: null as string | null,
  notes: ''
});

const statusSegments = [
  { label: '全部', value: 'all' },
  { label: '正常', value: 'normal' },
  { label: '即将到期', value: 'expiring' },
  { label: '已到期', value: 'expired' },
  { label: '停用', value: 'disabled' }
];
const pageSizeOptions = [20, 50, 100].map(value => ({ label: `${value} 条 / 页`, value }));
const familyOptions = [
  { label: 'IPv4', value: 4 },
  { label: 'IPv6', value: 6 }
];
const sourceTypeOptions = [
  { label: '节点公网出口', value: 'node' },
  { label: '独立 S5 出口', value: 's5' }
];
// 色板全部来自 ui-v2 / kit-v2 既有 token，不新增色值
const countryColorVars = ['--e-brand', '--e-cyan', '--e-violet', '--e-ok', '--e-warn', '--e-info', '--e-danger'];

function formatMoney(value?: number, currency = 'CNY') {
  if (value === undefined || value === null || !Number.isFinite(value)) return '--';
  return new Intl.NumberFormat('zh-CN', { style: 'currency', currency }).format(value);
}

function formatDate(value?: string | null) {
  if (!value) return '--';
  const date = dayjs(value);
  return date.isValid() ? date.format('YYYY-MM-DD HH:mm') : '时间未知';
}

function relativeTime(value?: string | null) {
  if (!value) return '时间未知';
  const timestamp = new Date(value).getTime();
  if (!Number.isFinite(timestamp)) return '时间未知';
  const seconds = Math.max(0, Math.round((Date.now() - timestamp) / 1000));
  if (seconds < 60) return `${seconds} 秒前`;
  const minutes = Math.round(seconds / 60);
  if (minutes < 60) return `${minutes} 分钟前`;
  const hours = Math.round(minutes / 60);
  if (hours < 24) return `${hours} 小时前`;
  return `${Math.round(hours / 24)} 天前`;
}

function familyLabel(family: number) {
  return family === 6 ? 'IPv6' : 'IPv4';
}

function formatDateOnly(value?: string | null) {
  return value ? value.slice(0, 10) : '未设置';
}

function daysUntilExpiry(value?: string | null) {
  if (!value) return null;
  const match = /^(\d{4})-(\d{2})-(\d{2})/.exec(value);
  if (!match) return null;
  const [, year, month, day] = match;
  const expiry = Date.UTC(Number(year), Number(month) - 1, Number(day));
  const now = new Date();
  const today = Date.UTC(now.getFullYear(), now.getMonth(), now.getDate());
  return Math.floor((expiry - today) / 86_400_000);
}

function exitIpDisplayStatus(row: Api.Central.ExitIpSummary): { label: string; tone: Tone } {
  if (row.status === 'disabled') return { label: '停用', tone: 'neutral' };
  if (row.status === 'unknown') return { label: '未知', tone: 'neutral' };
  const days = daysUntilExpiry(row.validTo);
  if (days !== null && days < 0) return { label: '已到期', tone: 'danger' };
  if (days !== null && days <= 7) return { label: '即将到期', tone: 'warn' };
  return { label: '正常', tone: 'ok' };
}

function expiryTone(row: Api.Central.ExitIpSummary) {
  if (row.status === 'disabled') return '';
  const days = daysUntilExpiry(row.validTo);
  if (days === null) return '';
  if (days < 0) return 'is-danger';
  if (days <= 7) return 'is-warn';
  return '';
}

function expiryHint(row: Api.Central.ExitIpSummary) {
  if (!row.validTo) return '长期有效';
  const days = daysUntilExpiry(row.validTo);
  if (days !== null && days < 0) return `已于 ${formatDateOnly(row.validTo)} 失效`;
  return `失效 ${formatDateOnly(row.validTo)}`;
}

function rowNote(row: ExitIpRow) {
  const notes = row.notes?.trim();
  if (notes) return notes;
  if (row.status === 'disabled') return '已停用 · 不再分配';
  if (row.allocatedUserCount === 0) return '空闲 · 可分配';
  return '—';
}

function nodeTypeLabel(type?: 'relay' | 'landing' | null) {
  if (type === 'relay') return '线路机';
  if (type === 'landing') return '落地机';
  return '未知类型';
}

function ownerNodeLabel(row: Api.Central.ExitIpSummary) {
  if (row.sourceType === 's5') return '独立 S5';
  return row.ownerNodeName || '未关联节点';
}

function ownerNodeTone(row: Api.Central.ExitIpSummary): Tone {
  if (row.sourceType === 's5') return 'neutral';
  return row.ownerNodeType === 'relay' ? 'brand' : 'neutral';
}

function ownerNodeSub(row: Api.Central.ExitIpSummary) {
  const region = (row.region || '').trim() || '未设置';
  if (row.sourceType === 's5') return `独立 S5 · ${region}`;
  return `${region} · ${nodeTypeLabel(row.ownerNodeType)}`;
}

function countryOf(row: Api.Central.ExitIpSummary) {
  return (row.region || '').trim() || '未设置';
}

function ipCellTitle(row: Api.Central.ExitIpSummary) {
  return [
    `月成本 ${formatMoney(row.monthlyCost, row.currency)}`,
    `最近检查 ${formatDate(row.checkedAt)}`,
    `协议族 ${familyLabel(row.family || 4)}`
  ].join(' · ');
}

function userStatusLabel(status?: Api.Central.UserStatus | null) {
  if (status === 'active') return '有效';
  if (status === 'expiring') return '即将到期';
  if (status === 'expired') return '已到期';
  if (status === 'disabled') return '已停用';
  return '状态未知';
}

function userStatusTone(status?: Api.Central.UserStatus | null): Tone {
  if (status === 'active') return 'ok';
  if (status === 'expiring') return 'warn';
  if (status === 'expired') return 'danger';
  return 'neutral';
}

function pathModeLabel(mode: Api.Central.ExitIpAllocatedUser['pathMode']) {
  if (mode === 'relay') return '线路机直出';
  if (mode === 'landing') return '经落地机';
  if (mode === 'external') return '独立 S5';
  return '--';
}

const ownerNodeOptions = computed(() =>
  nodes.value
    .filter(node => node.type === 'landing' || node.type === 'relay')
    .map(node => ({
      label: `${node.name} · ${node.type === 'relay' ? '线路机' : '落地机'}（${node.host || '地址未知'}）`,
      value: node.id
    }))
);
const selectedOwnerNode = computed(() => nodes.value.find(node => node.id === exitIpForm.ownerNodeId) || null);
const selectedOwnerRegion = computed(() => selectedOwnerNode.value?.region?.trim() || '');
const modalTitle = computed(() => (editingExitIpId.value ? '编辑出口 IP' : '新建出口 IP'));
const modalSubtitle = computed(() =>
  editingExitIpId.value ? '更新出口资产的归属、成本与有效期' : '登记中央出口资产并关联到线路配置'
);

const countryDistribution = computed(() => stats.value?.countries || []);
const countryTotal = computed(() => countryDistribution.value.reduce((sum, item) => sum + item.count, 0) || 1);
const countryBarSegments = computed(() =>
  countryDistribution.value.map((item, index) => ({
    ...item,
    colorVar: countryColorVars[index % countryColorVars.length],
    percent: (item.count / countryTotal.value) * 100
  }))
);
const countryOptions = computed(() => [
  { label: '全部国家 / 地区', value: 'all' },
  ...countryDistribution.value.map(item => ({ label: `${item.name}（${item.count}）`, value: item.name }))
]);

const ownedNodes = computed(() => nodes.value.filter(node => (node.exitIpCount || 0) > 0));
const ownerNodeCount = computed(() => (nodeError.value ? '--' : String(ownedNodes.value.length)));
const ownerRegionSummary = computed(() => {
  if (nodeError.value) return '节点列表读取失败';
  if (!ownedNodes.value.length) return '暂无归属节点';
  const groups = new Map<string, number>();
  ownedNodes.value.forEach(node => {
    const region = (node.region || '').trim() || '未设置';
    groups.set(region, (groups.get(region) || 0) + 1);
  });
  return Array.from(groups.entries())
    .sort((a, b) => b[1] - a[1])
    .map(([region, count]) => `${region} ${count}`)
    .join(' · ');
});

// 国家/地区筛选已下沉到服务端（region 参数），列表本身就是过滤结果
const displayRows = computed(() => rows.value);

const pageCount = computed(() => Math.max(1, Math.ceil(total.value / filters.page_size)));
const rangeStart = computed(() => (total.value === 0 ? 0 : (filters.page - 1) * filters.page_size + 1));
const rangeEnd = computed(() => Math.min(filters.page * filters.page_size, total.value));
const filterHint = computed(() => {
  const parts = [`第 ${filters.page} / ${pageCount.value} 页`];
  if (filters.country !== 'all') parts.push(`国家/地区：${filters.country}`);
  return parts.join(' · ');
});

const router = useRouter();
const { width: viewportWidth } = useWindowSize();
const allocatedUsersDrawerWidth = computed(() =>
  viewportWidth.value < 900 ? '100vw' : Math.min(1000, Math.round(viewportWidth.value * 0.9))
);
const allocatedUsersVisible = ref(false);
const selectedExitIp = ref<Api.Central.ExitIpSummary | null>(null);
const allocatedUsers = ref<Api.Central.ExitIpAllocatedUser[]>([]);
const allocatedUsersLoading = ref(false);
const allocatedUsersError = ref('');
const allocatedUsersTotal = ref(0);
const allocatedUsersPage = ref(1);
const allocatedUsersPageSize = ref(20);
const allocatedUsersKeyword = ref('');
let allocatedUsersRequestID = 0;
const allocatedUsersPageCount = computed(() =>
  Math.max(1, Math.ceil(allocatedUsersTotal.value / allocatedUsersPageSize.value))
);
const allocatedUsersPageSizeOptions = [10, 20, 50].map(value => ({ label: `${value} 条 / 页`, value }));
const selectedExitIpStatus = computed(() => (selectedExitIp.value ? exitIpDisplayStatus(selectedExitIp.value) : null));
const drawerSubtitle = computed(() => {
  if (!selectedExitIp.value) return '';
  const ip = selectedExitIp.value;
  const owner = ip.sourceType === 's5' ? '独立 S5' : ip.ownerNodeName || '未关联节点';
  const region = (ip.region || '').trim() || '未设置';
  return `${owner} · ${region} · ${allocatedUsersTotal.value} 人（含已到期）`;
});

watch(
  () => exitIpForm.sourceType,
  value => {
    if (value === 's5') {
      exitIpForm.ownerNodeId = '';
    } else {
      exitIpForm.region = '';
    }
  }
);

function resetExitIpForm() {
  editingExitIpId.value = null;
  exitIpForm.address = '';
  exitIpForm.sourceType = 'node';
  exitIpForm.ownerNodeId = '';
  exitIpForm.region = '';
  exitIpForm.family = 4;
  exitIpForm.provider = '';
  exitIpForm.monthlyCost = 0;
  exitIpForm.enabled = true;
  exitIpForm.validFrom = null;
  exitIpForm.validTo = null;
  exitIpForm.notes = '';
  formError.value = '';
}

function copyExitIpToForm(detail: Api.Central.ExitIpDetail) {
  exitIpForm.address = detail.address;
  exitIpForm.sourceType = detail.sourceType || 'node';
  exitIpForm.ownerNodeId = detail.ownerNodeId || detail.landingNodeId || '';
  exitIpForm.region = detail.sourceType === 's5' ? detail.region || '' : '';
  exitIpForm.family = detail.family;
  exitIpForm.provider = detail.provider || '';
  exitIpForm.monthlyCost = detail.monthlyCost;
  exitIpForm.enabled = detail.enabled;
  exitIpForm.validFrom = detail.validFrom ? detail.validFrom.slice(0, 10) : null;
  exitIpForm.validTo = detail.validTo ? detail.validTo.slice(0, 10) : null;
  exitIpForm.notes = detail.notes || '';
}

async function loadNodes() {
  nodesLoading.value = true;
  nodeError.value = '';
  const { data, error } = await fetchNodes({ page: 1, page_size: 200 });
  if (error || !data) {
    nodeError.value = '无法读取节点列表，请确认中央后端已启动后重试。';
    nodes.value = [];
  } else {
    nodes.value = data.items;
  }
  nodesLoading.value = false;
}

function openCreate() {
  resetExitIpForm();
  modalVisible.value = true;
}

async function openEdit(id: string) {
  resetExitIpForm();
  editingExitIpId.value = id;
  modalVisible.value = true;
  modalLoading.value = true;
  const { data, error } = await fetchExitIpDetail(id);
  modalLoading.value = false;
  if (error || !data) {
    formError.value = '无法读取出口 IP 详情，请刷新后重试。';
    return;
  }
  copyExitIpToForm(data);
}

function closeModal() {
  if (!saving.value) modalVisible.value = false;
}

async function saveExitIp() {
  const address = exitIpForm.address.trim();
  if (!address) {
    formError.value = '请填写出口 IP 地址';
    return;
  }
  if (exitIpForm.sourceType === 'node' && !exitIpForm.ownerNodeId) {
    formError.value = '请选择所属节点';
    return;
  }
  if (exitIpForm.sourceType === 's5' && !exitIpForm.region.trim()) {
    formError.value = '请填写 S5 所属国家 / 地区';
    return;
  }
  if (exitIpForm.monthlyCost === null || !Number.isFinite(exitIpForm.monthlyCost) || exitIpForm.monthlyCost < 0) {
    formError.value = '月成本必须是不小于 0 的数字';
    return;
  }
  if (exitIpForm.validFrom && exitIpForm.validTo && exitIpForm.validTo < exitIpForm.validFrom) {
    formError.value = '失效日期不能早于生效日期';
    return;
  }
  formError.value = '';
  saving.value = true;
  const payload: Api.Central.ExitIpUpsertPayload = {
    address,
    sourceType: exitIpForm.sourceType,
    ownerNodeId: exitIpForm.sourceType === 'node' ? exitIpForm.ownerNodeId : null,
    region: exitIpForm.sourceType === 's5' ? exitIpForm.region.trim() : null,
    family: exitIpForm.family,
    provider: exitIpForm.provider.trim(),
    monthlyCost: exitIpForm.monthlyCost,
    currency: 'CNY',
    enabled: exitIpForm.enabled,
    validFrom: exitIpForm.validFrom || '',
    validTo: exitIpForm.validTo || '',
    notes: exitIpForm.notes.trim()
  };
  const result = editingExitIpId.value
    ? await updateExitIp(editingExitIpId.value, payload)
    : await createExitIp(payload);
  saving.value = false;
  if (result.error || !result.data) {
    formError.value = '出口 IP 保存失败，请检查地址、来源和日期后重试。';
    return;
  }
  modalVisible.value = false;
  window.$message?.success(editingExitIpId.value ? '出口 IP 已更新' : '出口 IP 已创建');
  await loadExitIps();
}

async function removeExitIp(id: string) {
  if (deletingExitIpId.value) return;
  deletingExitIpId.value = id;
  const { error } = await deleteExitIp(id);
  deletingExitIpId.value = '';
  if (error) return;
  window.$message?.success('出口 IP 已删除');
  await loadExitIps();
}

function confirmDelete(row: Api.Central.ExitIpSummary) {
  const remove = () => void removeExitIp(row.id);
  if (!window.$dialog) {
    remove();
    return;
  }
  window.$dialog.warning({
    title: '删除出口 IP',
    content: `确定删除“${row.address}”吗？已绑定线路的出口 IP 需要先解除绑定或停用。`,
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: remove
  });
}

function rowMenuOptions(row: Api.Central.ExitIpSummary) {
  return [
    {
      label: '删除出口 IP',
      key: 'delete',
      disabled: Boolean(deletingExitIpId.value),
      icon: () => h(SvgIcon, { icon: 'mdi:delete-outline' })
    }
  ];
}

function handleRowAction(key: string, row: Api.Central.ExitIpSummary) {
  if (key === 'delete') confirmDelete(row);
}

function openAllocatedUsers(row: Api.Central.ExitIpSummary) {
  selectedExitIp.value = row;
  allocatedUsersPage.value = 1;
  allocatedUsersKeyword.value = '';
  allocatedUsersVisible.value = true;
  void loadAllocatedUsers();
}

async function loadAllocatedUsers() {
  const exitIp = selectedExitIp.value;
  if (!exitIp) return;
  const requestID = ++allocatedUsersRequestID;
  allocatedUsersLoading.value = true;
  allocatedUsersError.value = '';
  const { data, error } = await fetchExitIpUsers(exitIp.id, {
    page: allocatedUsersPage.value,
    page_size: allocatedUsersPageSize.value,
    keyword: allocatedUsersKeyword.value.trim() || undefined
  });
  if (requestID !== allocatedUsersRequestID) return;
  allocatedUsersLoading.value = false;
  if (error || !data) {
    allocatedUsersError.value = '无法读取该 IP 的用户列表，请重试。';
    allocatedUsers.value = [];
    allocatedUsersTotal.value = 0;
    return;
  }
  allocatedUsers.value = data.items;
  allocatedUsersTotal.value = data.total;
}

watch(allocatedUsersVisible, visible => {
  if (!visible) allocatedUsersRequestID += 1;
});

function submitAllocatedUserSearch() {
  allocatedUsersPage.value = 1;
  void loadAllocatedUsers();
}

function changeAllocatedUsersPage(page: number) {
  allocatedUsersPage.value = page;
  void loadAllocatedUsers();
}

function changeAllocatedUsersPageSize(pageSize: number) {
  allocatedUsersPageSize.value = pageSize;
  allocatedUsersPage.value = 1;
  void loadAllocatedUsers();
}

function openUserDetail(user: Api.Central.ExitIpAllocatedUser) {
  allocatedUsersVisible.value = false;
  void router.push({ path: '/users', query: { userId: user.id } });
}

async function loadExitIps() {
  loading.value = true;
  errorMessage.value = '';
  const { data, error } = await fetchExitIps({
    page: filters.page,
    page_size: filters.page_size,
    keyword: filters.keyword || undefined,
    status: filters.status === 'all' ? undefined : filters.status,
    region: filters.country === 'all' ? undefined : filters.country
  });
  if (error || !data) {
    errorMessage.value = '中央后端暂不可用，无法读取出口 IP 数据';
    rows.value = [];
    total.value = 0;
    stats.value = null;
    dataAt.value = '';
  } else {
    rows.value = data.items;
    total.value = data.total;
    stats.value = data.stats;
    dataAt.value = data.dataAt || '';
  }
  loading.value = false;
}

function submitFilters() {
  filters.page = 1;
  void loadExitIps();
}

function setStatusFilter(value: string) {
  filters.status = value;
  submitFilters();
}

function setCountryFilter(value: string) {
  filters.country = value;
  filters.page = 1;
  void loadExitIps();
}

function changePageSize(pageSize: number) {
  filters.page_size = pageSize;
  filters.page = 1;
  void loadExitIps();
}

function changePage(page: number) {
  filters.page = page;
  void loadExitIps();
}

function resetFilters() {
  filters.keyword = '';
  filters.status = 'all';
  filters.country = 'all';
  submitFilters();
}

onMounted(() => {
  void Promise.all([loadExitIps(), loadNodes()]);
});
</script>

<template>
  <div class="e-page">
    <!-- ===== 页头 ===== -->
    <header class="e-head">
      <div class="e-head-main">
        <div class="e-eyebrow">业务 · 出口 IP · 数据更新于 {{ relativeTime(dataAt) }}</div>
        <h1 class="e-title">出口 IP</h1>
        <p class="e-desc">
          维护线路机 / 落地机公网出口与独立 S5 资产、成本和线路配置归属关系。出口 IP 只记录归属与成本，<b>不在中央面板保存任何账号凭据</b>。
        </p>
      </div>
      <div class="e-head-actions">
        <DataFreshness :data-at="dataAt" compact />
        <button class="e-btn" type="button" :disabled="loading" @click="loadExitIps">
          <SvgIcon icon="mdi:refresh" class="e-btn-ico" />刷新
        </button>
        <button class="e-btn is-primary" type="button" @click="openCreate">
          <SvgIcon icon="mdi:plus" class="e-btn-ico" />新建出口 IP
        </button>
      </div>
    </header>

    <div v-if="errorMessage" class="e-alert is-danger">
      {{ errorMessage }}
      <button class="e-btn e-btn-xs" type="button" @click="loadExitIps">重试</button>
    </div>

    <!-- ===== KPI ===== -->
    <section class="e-kpis">
      <div class="e-kpi">
        <div class="e-kpi-top">
          <span class="e-kpi-label">出口 IP 总数</span>
          <span class="e-kpi-icon is-blue"><SvgIcon icon="mdi:globe" /></span>
        </div>
        <div class="e-kpi-value">{{ stats ? stats.total : '--' }}<span class="e-kpi-unit">个</span></div>
        <div class="e-kpi-foot">已登记的全部出口资产</div>
      </div>
      <div class="e-kpi">
        <div class="e-kpi-top">
          <span class="e-kpi-label">被用户占用</span>
          <span class="e-kpi-icon is-green"><SvgIcon icon="mdi:account-multiple-outline" /></span>
        </div>
        <div class="e-kpi-value">
          {{ stats ? stats.assigned : '--' }}<span class="e-kpi-unit">/ {{ stats ? stats.total : '--' }}</span>
        </div>
        <div class="e-kpi-foot">已有有效用户路径使用</div>
      </div>
      <div class="e-kpi">
        <div class="e-kpi-top">
          <span class="e-kpi-label">空闲</span>
          <span class="e-kpi-icon is-amber"><SvgIcon icon="mdi:power-plug-off" /></span>
        </div>
        <div class="e-kpi-value">{{ stats ? stats.unassigned : '--' }}<span class="e-kpi-unit">个</span></div>
        <div class="e-kpi-foot">可继续配置使用</div>
      </div>
      <div class="e-kpi">
        <div class="e-kpi-top">
          <span class="e-kpi-label">归属节点数</span>
          <span class="e-kpi-icon is-violet"><SvgIcon icon="mdi:server" /></span>
        </div>
        <div class="e-kpi-value">{{ ownerNodeCount }}<span class="e-kpi-unit">台</span></div>
        <div class="e-kpi-foot">{{ ownerRegionSummary }}</div>
      </div>
    </section>

    <!-- ===== 国家分布 + 筛选 ===== -->
    <section class="e-card">
      <div class="e-card-head">
        <h3>国家 / 地区分布</h3>
        <span class="e-badge is-neutral">全部资产 {{ stats ? stats.total : '--' }} 个</span>
        <div class="e-right">
          <span class="e-chip"><SvgIcon icon="mdi:filter-variant" />按归属节点区域统计</span>
        </div>
      </div>
      <div class="e-card-body">
        <div v-if="countryDistribution.length" class="e-bar" aria-label="国家 / 地区分布堆叠条">
          <span
            v-for="segment in countryBarSegments"
            :key="segment.name"
            :style="{ width: `${segment.percent}%`, background: `var(${segment.colorVar})` }"
            :title="`${segment.name} ${segment.count} 个`"
          ></span>
        </div>
        <div v-else class="e-bar is-empty"></div>
        <div class="e-legend">
          <button
            v-for="(segment, index) in countryBarSegments"
            :key="segment.name"
            type="button"
            class="e-legend-item"
            :class="{ 'is-active': filters.country === segment.name }"
            :title="`筛选 ${segment.name}`"
            @click="setCountryFilter(filters.country === segment.name ? 'all' : segment.name)"
          >
            <i class="e-legend-dot" :style="{ background: `var(${countryColorVars[index % countryColorVars.length]})` }"></i>
            <span>{{ segment.name }}</span>
            <b class="e-num">{{ segment.count }}</b>
          </button>
          <span class="e-legend-hint">按数量降序 · 点击图例筛选</span>
        </div>
      </div>

      <div class="e-toolbar">
        <label class="e-input e-search">
          <SvgIcon icon="mdi:magnify" class="e-input-ico" />
          <input
            v-model="filters.keyword"
            placeholder="搜索 IP、国家/地区、节点或服务商"
            @keyup.enter="submitFilters"
          />
        </label>
        <NSelect
          v-model:value="filters.country"
          :options="countryOptions"
          size="small"
          class="e-country-select"
          :title="'按国家 / 地区筛选（服务端过滤）'"
          @update:value="setCountryFilter"
        />
        <div class="e-seg">
          <button
            v-for="segment in statusSegments"
            :key="segment.value"
            type="button"
            :class="{ 'is-active': filters.status === segment.value }"
            @click="setStatusFilter(segment.value)"
          >
            {{ segment.label }}
          </button>
        </div>
        <span class="e-spacer"></span>
        <button class="e-btn e-btn-xs is-primary" type="button" @click="submitFilters">
          <SvgIcon icon="mdi:magnify" class="e-btn-ico" />查询
        </button>
        <button class="e-btn e-btn-xs" type="button" @click="resetFilters">重置</button>
      </div>
    </section>

    <!-- ===== 台账表 ===== -->
    <section class="e-card">
      <div class="e-card-head">
        <h3>出口 IP 台账</h3>
        <span class="e-badge is-neutral">{{ total }} 条</span>
        <span class="e-card-sub">归属关系与成本口径以中央记录为准</span>
        <div class="e-right">
          <span class="e-dim">{{ filterHint }}</span>
        </div>
      </div>

      <div class="e-table-wrap">
        <table class="e-table">
          <thead>
            <tr>
              <th>出口 IP</th>
              <th>归属节点</th>
              <th>国家 / 地区</th>
              <th>使用用户数</th>
              <th>状态</th>
              <th>备注</th>
              <th class="e-th-actions"></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="row in displayRows" :key="row.id">
              <td>
                <span class="e-ip e-mono" :title="ipCellTitle(row)">{{ row.address }}</span>
                <span class="e-sub">{{ familyLabel(row.family || 4) }} · {{ row.provider || '未登记服务商' }}</span>
              </td>
              <td>
                <span class="e-badge" :class="`is-${ownerNodeTone(row)}`">{{ ownerNodeLabel(row) }}</span>
                <span class="e-sub">{{ ownerNodeSub(row) }}</span>
              </td>
              <td>{{ countryOf(row) }}</td>
              <td>
                <button
                  type="button"
                  class="e-user-btn"
                  :class="{ 'is-zero': row.allocatedUserCount === 0 }"
                  :aria-label="`${row.address} 的配置归属用户：${row.allocatedUserCount} 人`"
                  @click="openAllocatedUsers(row)"
                >
                  <SvgIcon icon="mdi:account-multiple-outline" class="e-user-btn-ico" />
                  {{ row.allocatedUserCount }} 人
                </button>
              </td>
              <td>
                <span class="e-badge" :class="`is-${exitIpDisplayStatus(row).tone}`">
                  <i class="e-dot"></i>{{ exitIpDisplayStatus(row).label }}
                </span>
                <span class="e-sub" :class="expiryTone(row)">{{ expiryHint(row) }}</span>
              </td>
              <td class="e-td-notes">
                <span class="e-notes" :title="rowNote(row)">{{ rowNote(row) }}</span>
              </td>
              <td class="e-td-actions">
                <button class="e-btn e-btn-xs" type="button" @click="openEdit(row.id)">编辑</button>
                <NDropdown
                  trigger="click"
                  :options="rowMenuOptions(row)"
                  @select="key => handleRowAction(String(key), row)"
                >
                  <button class="e-icon-btn" type="button" title="更多操作">
                    <SvgIcon icon="mdi:dots-horizontal" />
                  </button>
                </NDropdown>
              </td>
            </tr>
          </tbody>
        </table>

        <div v-if="loading" class="e-table-state">
          <NSpin size="small" />
          <span>正在读取出口 IP…</span>
        </div>
        <div v-else-if="!displayRows.length" class="e-table-state">
          {{ filters.country !== 'all' ? '当前页没有该国家 / 地区的出口 IP' : '暂无出口 IP 数据' }}
        </div>
      </div>

      <div class="e-table-foot">
        <span class="e-dim">显示 {{ rangeStart }}–{{ rangeEnd }} 条 · 共 {{ total }} 条</span>
        <div class="e-foot-right">
          <NSelect
            v-model:value="filters.page_size"
            :options="pageSizeOptions"
            size="small"
            class="e-pagesize"
            @update:value="changePageSize"
          />
          <NPagination
            :page="filters.page"
            :page-count="pageCount"
            :page-slot="7"
            size="small"
            @update:page="changePage"
          />
        </div>
      </div>
    </section>

    <!-- ===================== 抽屉：使用该 IP 的用户 ===================== -->
    <NDrawer v-model:show="allocatedUsersVisible" :width="allocatedUsersDrawerWidth" placement="right">
      <NDrawerContent body-content-style="padding: 0">
        <template #header>
          <div class="e-dhead">
            <div class="e-dhead-main">
              <h2 class="e-dhead-title">使用 {{ selectedExitIp?.address || '' }} 的用户</h2>
              <p class="e-dhead-sub">{{ drawerSubtitle }}</p>
            </div>
            <div class="e-dhead-actions">
              <span v-if="selectedExitIpStatus" class="e-badge" :class="`is-${selectedExitIpStatus.tone}`">
                <i class="e-dot"></i>IP {{ selectedExitIpStatus.label }}
              </span>
              <button class="e-icon-btn" type="button" title="关闭" @click="allocatedUsersVisible = false">
                <SvgIcon icon="mdi:close" />
              </button>
            </div>
          </div>
        </template>

        <div class="e-drawer-body">
          <div class="e-toolbar e-toolbar-flat">
            <label class="e-input e-search">
              <SvgIcon icon="mdi:magnify" class="e-input-ico" />
              <input
                v-model="allocatedUsersKeyword"
                placeholder="搜索用户、线路机或 Inbound"
                @keyup.enter="submitAllocatedUserSearch"
              />
            </label>
            <button class="e-btn e-btn-xs is-primary" type="button" @click="submitAllocatedUserSearch">
              <SvgIcon icon="mdi:magnify" class="e-btn-ico" />查询
            </button>
            <span class="e-spacer"></span>
            <span class="e-dim">当前归属用户：{{ allocatedUsersTotal }} 人（包含已到期用户）</span>
          </div>

          <div v-if="allocatedUsersError" class="e-alert is-danger">
            {{ allocatedUsersError }}
            <button class="e-btn e-btn-xs" type="button" @click="loadAllocatedUsers">重试</button>
          </div>

          <div class="e-table-wrap e-drawer-table">
            <table class="e-table is-dense">
              <thead>
                <tr>
                  <th>用户</th>
                  <th>线路机 / Inbound</th>
                  <th>路径</th>
                  <th>用户状态</th>
                  <th>用户到期</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="user in allocatedUsers" :key="user.id">
                  <td>
                    <button type="button" class="e-link" @click="openUserDetail(user)">{{ user.name }}</button>
                  </td>
                  <td>
                    <span class="e-badge is-brand">{{ user.nodeName || '--' }}</span>
                    <span class="e-sub">{{ user.inboundTag || 'Inbound 未命名' }}</span>
                  </td>
                  <td>{{ pathModeLabel(user.pathMode) }}</td>
                  <td>
                    <span class="e-badge" :class="`is-${userStatusTone(user.status)}`">
                      <i class="e-dot"></i>{{ userStatusLabel(user.status) }}
                    </span>
                  </td>
                  <td class="e-num">{{ formatDate(user.expiresAt) }}</td>
                </tr>
              </tbody>
            </table>

            <div v-if="allocatedUsersLoading" class="e-table-state">
              <NSpin size="small" />
              <span>正在读取用户…</span>
            </div>
            <div v-else-if="!allocatedUsers.length" class="e-table-state">暂无符合条件的用户</div>
          </div>

          <p class="e-drawer-hint">
            归属关系来自 user_paths + user_path_exit_ips；包含已到期用户。点击用户可跳转用户详情。
          </p>
        </div>

        <template #footer>
          <div class="e-drawer-foot">
            <span class="e-dim">路径来源：user_paths / user_path_exit_ips</span>
            <div class="e-foot-right">
              <NSelect
                v-model:value="allocatedUsersPageSize"
                :options="allocatedUsersPageSizeOptions"
                size="small"
                class="e-pagesize"
                @update:value="changeAllocatedUsersPageSize"
              />
              <NPagination
                :page="allocatedUsersPage"
                :page-count="allocatedUsersPageCount"
                :page-slot="5"
                size="small"
                @update:page="changeAllocatedUsersPage"
              />
              <button class="e-btn e-btn-xs" type="button" @click="allocatedUsersVisible = false">关闭</button>
            </div>
          </div>
        </template>
      </NDrawerContent>
    </NDrawer>

    <!-- ===================== 新建 / 编辑出口 IP ===================== -->
    <Teleport to="body">
      <div v-if="modalVisible" class="e-mask" @click.self="closeModal">
        <div class="e-modal" role="dialog" aria-modal="true" :aria-label="modalTitle">
          <div class="e-modal-head">
            <div class="e-modal-head-main">
              <h3>{{ modalTitle }}</h3>
              <p class="e-modal-sub">{{ modalSubtitle }}</p>
            </div>
            <span class="e-badge is-brand">资产登记</span>
            <button class="e-icon-btn" type="button" title="关闭" :disabled="saving" @click="closeModal">
              <SvgIcon icon="mdi:close" />
            </button>
          </div>

          <div class="e-modal-body">
            <div class="e-callout is-info">
              <SvgIcon icon="mdi:help-circle-outline" class="e-callout-ico" />
              <div>
                <b>只登记地址、来源与成本</b>
                独立购买的 S5 只登记出口地址、国家/地区和成本；账号凭据不在中央面板保存。出口来源切换为「独立 S5」时，需手动填写国家 / 地区。
              </div>
            </div>

            <div v-if="nodeError" class="e-alert is-warn">
              {{ nodeError }}
              <button class="e-btn e-btn-xs" type="button" @click="loadNodes">重试</button>
            </div>
            <div v-if="formError" class="e-alert is-danger">{{ formError }}</div>

            <div v-if="modalLoading" class="e-modal-loading"><NSpin size="medium" /></div>

            <div v-else class="e-form-grid">
              <div class="e-field">
                <span class="e-label">出口 IP 地址 <span class="e-req">*</span></span>
                <NInput
                  v-model:value="exitIpForm.address"
                  maxlength="45"
                  placeholder="例如：203.0.113.10 或 2001:db8::10"
                />
              </div>
              <div class="e-field">
                <span class="e-label">协议族 <span class="e-req">*</span></span>
                <NSelect v-model:value="exitIpForm.family" :options="familyOptions" />
              </div>
              <div class="e-field">
                <span class="e-label">出口来源 <span class="e-req">*</span></span>
                <NSelect v-model:value="exitIpForm.sourceType" :options="sourceTypeOptions" />
              </div>
              <div v-if="exitIpForm.sourceType === 'node'" class="e-field">
                <span class="e-label">所属节点 <span class="e-req">*</span></span>
                <NSelect
                  v-model:value="exitIpForm.ownerNodeId"
                  :options="ownerNodeOptions"
                  :loading="nodesLoading"
                  filterable
                  clearable
                  placeholder="选择线路机或落地机"
                />
              </div>
              <div v-if="exitIpForm.sourceType === 'node'" class="e-field">
                <span class="e-label">国家 / 地区</span>
                <NInput :value="selectedOwnerRegion || '未设置国家 / 地区'" readonly />
              </div>
              <div v-else class="e-field">
                <span class="e-label">国家 / 地区 <span class="e-req">*</span></span>
                <NInput v-model:value="exitIpForm.region" maxlength="120" placeholder="例如：美国、英国、香港" />
              </div>
              <div class="e-field">
                <span class="e-label">状态</span>
                <div class="e-switch-row">
                  <NSwitch v-model:value="exitIpForm.enabled">
                    <template #checked>启用</template>
                    <template #unchecked>停用</template>
                  </NSwitch>
                  <span class="e-dim">{{ exitIpForm.enabled ? '启用中，可被线路配置引用' : '已停用，不再分配' }}</span>
                </div>
              </div>
              <div class="e-field">
                <span class="e-label">服务商</span>
                <NInput v-model:value="exitIpForm.provider" maxlength="200" placeholder="可选，例如 Example Cloud" />
              </div>
              <div class="e-field">
                <span class="e-label">月成本（CNY）</span>
                <NInputNumber
                  v-model:value="exitIpForm.monthlyCost"
                  :min="0"
                  :max="100000000"
                  :precision="2"
                  class="w-full"
                />
              </div>
              <div class="e-field">
                <span class="e-label">生效日期</span>
                <NDatePicker
                  v-model:formatted-value="exitIpForm.validFrom"
                  type="date"
                  value-format="yyyy-MM-dd"
                  clearable
                  class="w-full"
                />
              </div>
              <div class="e-field">
                <span class="e-label">失效日期</span>
                <NDatePicker
                  v-model:formatted-value="exitIpForm.validTo"
                  type="date"
                  value-format="yyyy-MM-dd"
                  clearable
                  class="w-full"
                />
              </div>
              <div class="e-field e-field-wide">
                <span class="e-label">备注</span>
                <NInput
                  v-model:value="exitIpForm.notes"
                  type="textarea"
                  maxlength="2000"
                  show-count
                  :autosize="{ minRows: 2, maxRows: 5 }"
                  placeholder="可选，例如：香港出口 · hk-02 直出"
                />
              </div>
              <div v-if="exitIpForm.sourceType === 's5'" class="e-field e-field-wide">
                <span class="e-help">
                  独立购买的 S5 只登记出口地址、国家/地区和成本；账号凭据不在中央面板保存。
                </span>
              </div>
            </div>
          </div>

          <div class="e-modal-foot">
            <span class="e-dim">失效日期留空表示长期有效；已绑定线路的出口 IP 需先解除绑定或停用后才能删除。</span>
            <div class="e-foot-right">
              <button class="e-btn" type="button" :disabled="saving" @click="closeModal">取消</button>
              <button
                class="e-btn is-primary"
                type="button"
                :disabled="saving || modalLoading"
                @click="saveExitIp"
              >
                <SvgIcon v-if="!saving" icon="mdi:check" class="e-btn-ico" />
                {{ saving ? '保存中…' : '保存' }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<style scoped>
/* ============================================================
   出口 IP · UI v2
   颜色 / 字号 / 间距对齐 design-preview/v2/exit-ips.html 与 kit-v2
   ============================================================ */
.e-page,
.e-mask,
.e-dhead,
.e-drawer-body,
.e-drawer-foot {
  --e-card: #141519;
  --e-card-hover: #191b20;
  --e-line: #23262c;
  --e-line-soft: #1c1f24;
  --e-text-1: #f2f3f5;
  --e-text-2: #bcc2cc;
  --e-text-3: #9aa1ad;
  --e-text-4: #868d99;
  --e-brand: #4c6fff;
  --e-cyan: #34d3e0;
  --e-ok: #2ed47a;
  --e-warn: #f5b93d;
  --e-danger: #f3656f;
  --e-info: #5aa9f0;
  --e-violet: #a78bfa;
  --e-radius: 14px;
}

.e-page {
  display: flex;
  flex-direction: column;
  gap: 15px;
  min-width: 0;
  color: var(--e-text-1);
  font-size: 13px;
}

.e-dim { color: var(--e-text-4); }
.e-mono { font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace; }
.e-num { font-variant-numeric: tabular-nums; }

/* ---------- 页头 ---------- */
.e-head { display: flex; align-items: flex-end; gap: 18px; flex-wrap: wrap; }
.e-head-main { min-width: 0; flex: 1; }
.e-eyebrow { margin-bottom: 6px; font-size: 12px; letter-spacing: 0.1em; color: #7f8792; }
.e-title { margin: 0; font-size: 22px; font-weight: 620; line-height: 1.25; color: #f6f7f9; }
.e-desc { max-width: 900px; margin: 7px 0 0; font-size: 12px; line-height: 1.7; color: var(--e-text-3); }
.e-desc b { color: var(--e-text-2); font-weight: 600; }
.e-head-actions { display: flex; align-items: center; gap: 8px; flex: none; padding-bottom: 2px; }

/* ---------- 按钮 / 输入 ---------- */
.e-btn {
  display: inline-flex; align-items: center; justify-content: center; gap: 6px;
  height: 34px; padding: 0 13px;
  border: 1px solid #2a2e35; border-radius: 9px;
  background: #1b1e24; color: #e6e8ec;
  font-size: 12.5px; white-space: nowrap; cursor: pointer;
  transition: background 0.15s ease, border-color 0.15s ease;
}
.e-btn:hover { background: #232730; border-color: #363b44; }
.e-btn:disabled { opacity: 0.45; cursor: not-allowed; }
.e-btn-ico { width: 15px; height: 15px; flex: none; }
.e-btn.is-primary { border-color: var(--e-brand); background: var(--e-brand); color: #fff; }
.e-btn.is-primary:hover { background: #5d7cff; border-color: #5d7cff; }
.e-btn-xs { height: 28px; padding: 0 10px; font-size: 12px; border-radius: 7px; }

.e-icon-btn {
  display: inline-grid; place-items: center; width: 28px; height: 28px;
  border: 1px solid transparent; border-radius: 7px; background: transparent;
  color: #9aa1ab; cursor: pointer;
}
.e-icon-btn:hover { background: #1a1d22; color: #e6e8ec; }
.e-icon-btn:disabled { opacity: 0.4; cursor: not-allowed; }
.e-icon-btn :deep(.svg-icon) { width: 16px; height: 16px; }

.e-input {
  display: flex; align-items: center; gap: 8px; height: 34px; padding: 0 10px;
  border: 1px solid #262a31; border-radius: 9px; background: #0e1013;
}
.e-input:focus-within { border-color: rgba(76, 111, 255, 0.7); }
.e-input input {
  min-width: 0; flex: 1; border: 0; outline: none; background: transparent;
  color: var(--e-text-1); font-size: 12.5px; font-family: inherit;
}
.e-input input::placeholder { color: var(--e-text-4); }
.e-input-ico { width: 15px; height: 15px; flex: none; color: var(--e-text-4); }
.e-search { width: 280px; flex: none; }
.e-country-select { width: 170px; flex: none; }
.e-pagesize { width: 104px; }

/* ---------- KPI ---------- */
.e-kpis { display: grid; grid-template-columns: repeat(4, minmax(0, 1fr)); gap: 12px; }
.e-kpi { padding: 15px 16px 14px; border: 1px solid var(--e-line); border-radius: var(--e-radius); background: var(--e-card); min-width: 0; }
.e-kpi-top { display: flex; align-items: center; gap: 8px; margin-bottom: 10px; }
.e-kpi-label { font-size: 12px; color: var(--e-text-3); white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.e-kpi-icon { display: grid; place-items: center; width: 24px; height: 24px; margin-left: auto; flex: none; border-radius: 8px; border: 1px solid #2a2e35; background: #1d2026; color: #a8adb8; }
.e-kpi-icon :deep(.svg-icon) { width: 13px; height: 13px; }
.e-kpi-icon.is-blue { border-color: rgba(76, 111, 255, 0.36); background: rgba(76, 111, 255, 0.16); color: #8fa8ff; }
.e-kpi-icon.is-green { border-color: rgba(46, 212, 122, 0.34); background: rgba(46, 212, 122, 0.14); color: #5ce49b; }
.e-kpi-icon.is-violet { border-color: rgba(167, 139, 250, 0.36); background: rgba(167, 139, 250, 0.16); color: #c0a9ff; }
.e-kpi-icon.is-amber { border-color: rgba(245, 185, 61, 0.34); background: rgba(245, 185, 61, 0.14); color: #f7ca6b; }
.e-kpi-value { display: flex; align-items: baseline; gap: 6px; font-size: 25px; font-weight: 620; line-height: 1.15; letter-spacing: -0.015em; font-variant-numeric: tabular-nums; color: #fff; }
.e-kpi-unit { font-size: 12px; font-weight: 500; color: var(--e-text-3); }
.e-kpi-foot { margin-top: 9px; font-size: 12px; line-height: 1.5; color: var(--e-text-4); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }

/* ---------- 卡片 ---------- */
.e-card { border: 1px solid var(--e-line); border-radius: var(--e-radius); background: var(--e-card); min-width: 0; overflow: hidden; }
.e-card-head { display: flex; align-items: center; gap: 9px; flex-wrap: wrap; padding: 14px 16px; border-bottom: 1px solid var(--e-line-soft); }
.e-card-head h3 { margin: 0; font-size: 14px; font-weight: 600; color: var(--e-text-1); }
.e-card-sub { font-size: 12px; color: var(--e-text-4); }
.e-right { display: flex; align-items: center; gap: 8px; margin-left: auto; }
.e-card-body { padding: 14px 16px 15px; }
.e-spacer { flex: 1; }

.e-chip {
  display: inline-flex; align-items: center; gap: 5px; height: 22px; padding: 0 8px;
  border: 1px solid #2a2e35; border-radius: 7px; background: #1a1d22;
  font-size: 11.5px; color: var(--e-text-3);
}
.e-chip :deep(.svg-icon) { width: 12px; height: 12px; }

/* ---------- 国家分布 ---------- */
.e-bar { display: flex; gap: 2px; height: 12px; border-radius: 6px; overflow: hidden; background: #1a1d22; }
.e-bar.is-empty { opacity: 0.4; }
.e-bar > span { min-width: 2px; height: 100%; }
.e-legend { display: flex; flex-wrap: wrap; align-items: center; gap: 8px 20px; margin-top: 12px; }
.e-legend-item {
  display: inline-flex; align-items: center; gap: 6px; height: 22px; padding: 0 2px;
  border: 0; border-radius: 6px; background: transparent; color: var(--e-text-3);
  font-size: 12px; cursor: pointer;
}
.e-legend-item:hover { color: var(--e-text-1); }
.e-legend-item.is-active { color: #fff; }
.e-legend-item.is-active span { text-decoration: underline; text-underline-offset: 3px; }
.e-legend-dot { width: 8px; height: 8px; flex: none; border-radius: 2px; }
.e-legend-item b { margin-left: 1px; font-weight: 600; color: var(--e-text-1); }
.e-legend-hint { margin-left: auto; font-size: 12px; color: var(--e-text-4); }

/* ---------- 工具条 / 分段 ---------- */
.e-toolbar {
  display: flex; align-items: center; gap: 8px; flex-wrap: wrap;
  padding: 12px 16px; border-top: 1px solid var(--e-line-soft);
}
.e-toolbar-flat { padding: 14px 18px; border-top: 0; border-bottom: 1px solid var(--e-line-soft); }
.e-seg { display: inline-flex; gap: 2px; padding: 3px; border: 1px solid #24272e; border-radius: 9px; background: #14161a; }
.e-seg button {
  height: 26px; padding: 0 11px; border: 0; border-radius: 7px;
  background: transparent; color: var(--e-text-3); font-size: 12px; cursor: pointer;
}
.e-seg button:hover { color: var(--e-text-1); }
.e-seg button.is-active { background: #2b3038; color: #fff; }

/* ---------- 表格 ---------- */
.e-table-wrap { position: relative; overflow-x: auto; }
.e-table { width: 100%; min-width: 940px; border-collapse: separate; border-spacing: 0; font-size: 13px; }
.e-table thead th {
  padding: 9px 16px; border-bottom: 1px solid var(--e-line);
  background: var(--e-card); color: var(--e-text-3);
  font-size: 11.5px; font-weight: 600; letter-spacing: 0.05em;
  text-align: left; white-space: nowrap;
}
.e-table tbody td {
  height: 52px; padding: 8px 16px;
  border-bottom: 1px solid var(--e-line-soft);
  color: var(--e-text-2); vertical-align: middle; white-space: nowrap;
}
.e-table tbody tr:last-child td { border-bottom: none; }
.e-table tbody tr:hover td { background: var(--e-card-hover); }
.e-table.is-dense { min-width: 720px; }
.e-table.is-dense tbody td { height: 46px; padding: 7px 14px; }
.e-th-actions { width: 1%; }
.e-td-actions { width: 1%; white-space: nowrap; }
.e-td-actions > * + * { margin-left: 4px; }
.e-td-notes { max-width: 220px; }
.e-ip { display: block; font-size: 13px; font-weight: 550; color: var(--e-text-1); letter-spacing: -0.01em; }
.e-sub { display: block; margin-top: 2px; font-size: 12px; color: var(--e-text-4); }
.e-sub.is-warn { color: var(--e-warn); }
.e-sub.is-danger { color: var(--e-danger); }
.e-notes { display: block; max-width: 220px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; font-size: 12px; color: var(--e-text-4); }

.e-user-btn {
  display: inline-flex; align-items: center; gap: 5px; height: 26px; padding: 0 9px;
  border: 1px solid rgba(76, 111, 255, 0.28); border-radius: 7px;
  background: rgba(76, 111, 255, 0.12); color: #a8b8ff;
  font-size: 12px; font-variant-numeric: tabular-nums; cursor: pointer;
  transition: background 0.15s ease, border-color 0.15s ease;
}
.e-user-btn:hover { background: rgba(76, 111, 255, 0.2); border-color: rgba(76, 111, 255, 0.45); }
.e-user-btn.is-zero { border-color: #2a2e35; background: #1a1d22; color: var(--e-text-4); }
.e-user-btn-ico { width: 13px; height: 13px; flex: none; }

.e-table-state {
  display: flex; align-items: center; justify-content: center; gap: 8px;
  padding: 40px 16px; color: var(--e-text-4); font-size: 12.5px;
}
.e-table-foot {
  display: flex; align-items: center; justify-content: space-between; gap: 10px;
  padding: 11px 16px; border-top: 1px solid var(--e-line-soft); font-size: 12px;
}
.e-foot-right { display: flex; align-items: center; gap: 8px; }

/* ---------- 徽标 / 圆点 ---------- */
.e-badge {
  display: inline-flex; align-items: center; gap: 5px; height: 21px; padding: 0 7px;
  border: 1px solid transparent; border-radius: 6px;
  font-size: 11.5px; font-weight: 550; white-space: nowrap;
}
.e-badge.is-neutral { border-color: #2a2e35; background: #1d2026; color: #bcc2cc; }
.e-badge.is-brand { border-color: rgba(76, 111, 255, 0.34); background: rgba(76, 111, 255, 0.16); color: #a8b8ff; }
.e-badge.is-ok { border-color: rgba(46, 212, 122, 0.32); background: rgba(46, 212, 122, 0.14); color: #5ce49b; }
.e-badge.is-warn { border-color: rgba(245, 185, 61, 0.32); background: rgba(245, 185, 61, 0.15); color: #f7ca6b; }
.e-badge.is-danger { border-color: rgba(243, 101, 111, 0.32); background: rgba(243, 101, 111, 0.15); color: #ff9ba2; }
.e-dot { width: 6px; height: 6px; flex: none; border-radius: 50%; background: currentColor; }

/* ---------- 错误 / 提示 ---------- */
.e-alert {
  display: flex; align-items: center; gap: 10px;
  padding: 10px 13px; border: 1px solid rgba(243, 101, 111, 0.32); border-radius: 10px;
  background: rgba(243, 101, 111, 0.1); color: #ffb3b9; font-size: 12.5px;
}
.e-alert.is-warn { border-color: rgba(245, 185, 61, 0.3); background: rgba(245, 185, 61, 0.1); color: #f7ca6b; }
.e-alert.is-danger { border-color: rgba(243, 101, 111, 0.32); background: rgba(243, 101, 111, 0.1); color: #ffb3b9; }
.e-alert .e-btn { margin-left: auto; }

/* ---------- 抽屉 ---------- */
.e-dhead { display: flex; align-items: flex-start; gap: 12px; width: 100%; }
.e-dhead-main { flex: 1; min-width: 0; }
.e-dhead-title { margin: 0; font-size: 15px; font-weight: 600; color: var(--e-text-1); }
.e-dhead-sub { margin: 3px 0 0; font-size: 12px; color: var(--e-text-4); }
.e-dhead-actions { display: flex; align-items: center; gap: 8px; flex: none; }
.e-drawer-body { display: flex; flex-direction: column; gap: 12px; padding-bottom: 8px; }
.e-drawer-body > .e-alert { margin: 0 18px; }
.e-drawer-table { margin: 0 18px; border: 1px solid var(--e-line-soft); border-radius: 10px; }
.e-drawer-hint { margin: 0 18px; font-size: 12px; line-height: 1.7; color: var(--e-text-4); }
.e-drawer-foot { display: flex; align-items: center; justify-content: space-between; gap: 10px; width: 100%; font-size: 12px; }
.e-link {
  border: 0; padding: 0; background: transparent; color: #8fa8ff;
  font-size: 13px; font-weight: 550; cursor: pointer;
}
.e-link:hover { color: #b7c6ff; text-decoration: underline; text-underline-offset: 3px; }

/* ---------- 弹窗 ---------- */
.e-mask {
  position: fixed; inset: 0; z-index: 2000;
  display: grid; place-items: center; padding: 24px;
  background: rgba(0, 0, 0, 0.66);
}
.e-modal {
  display: flex; flex-direction: column;
  width: 720px; max-width: 95vw; max-height: min(88vh, 900px);
  border: 1px solid var(--e-line); border-radius: 16px; background: var(--e-card);
  color: var(--e-text-1); text-align: left;
  box-shadow: 0 30px 80px -24px rgba(0, 0, 0, 0.9);
}
.e-modal-head { display: flex; align-items: flex-start; gap: 10px; padding: 16px 20px 14px; border-bottom: 1px solid var(--e-line-soft); }
.e-modal-head-main { flex: 1; min-width: 0; }
.e-modal-head h3 { margin: 0; font-size: 15px; font-weight: 600; }
.e-modal-sub { margin: 3px 0 0; font-size: 12px; color: var(--e-text-4); }
.e-modal-body { flex: 1; min-height: 0; overflow-y: auto; padding: 16px 20px 20px; display: flex; flex-direction: column; gap: 12px; }
.e-modal-loading { display: flex; align-items: center; justify-content: center; min-height: 220px; }
.e-modal-foot {
  display: flex; align-items: center; gap: 10px;
  padding: 13px 20px; border-top: 1px solid var(--e-line-soft); font-size: 12px;
}
.e-modal-foot .e-foot-right { margin-left: auto; }

.e-callout {
  display: flex; gap: 10px; padding: 12px 13px;
  border: 1px solid var(--e-line); border-radius: 11px;
  background: rgba(255, 255, 255, 0.022);
  font-size: 12px; line-height: 1.7; color: var(--e-text-2);
}
.e-callout-ico { width: 16px; height: 16px; flex: none; margin-top: 2px; color: var(--e-text-3); }
.e-callout b { display: block; margin-bottom: 1px; font-size: 12.5px; font-weight: 600; color: var(--e-text-1); }
.e-callout.is-info { border-color: rgba(90, 169, 240, 0.24); background: rgba(90, 169, 240, 0.08); }
.e-callout.is-info .e-callout-ico { color: #8cc4f7; }

.e-form-grid { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 12px; }
.e-field { display: flex; flex-direction: column; gap: 5px; min-width: 0; }
.e-field-wide { grid-column: 1 / -1; }
.e-label { font-size: 12px; color: var(--e-text-2); }
.e-req { color: var(--e-danger); }
.e-help { font-size: 12px; line-height: 1.6; color: var(--e-text-4); }
.e-switch-row { display: flex; align-items: center; gap: 9px; min-height: 34px; font-size: 12px; }

/* ---------- 响应式 ---------- */
@media (max-width: 1100px) {
  .e-kpis { grid-template-columns: repeat(2, minmax(0, 1fr)); }
  .e-search { width: 100%; }
}
@media (max-width: 820px) {
  .e-form-grid { grid-template-columns: minmax(0, 1fr); }
  .e-legend-hint { margin-left: 0; }
  .e-modal-foot { flex-wrap: wrap; }
  .e-modal-foot .e-foot-right { margin-left: 0; width: 100%; justify-content: flex-end; }
}
</style>
