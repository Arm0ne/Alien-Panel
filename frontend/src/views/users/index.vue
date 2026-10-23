<script setup lang="ts">
import { computed, nextTick, onMounted, reactive, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import {
  assignUserPath,
  cancelUserBillingRecord,
  createUserBillingRecord,
  confirmUserRenewal,
  clearUserPath,
  fetchEvents,
  fetchUserDetail,
  fetchUserGroups,
  fetchUserPathAssets,
  fetchUsers,
  importUserBillingRecords,
  rejectUserRenewal,
  resetInboundUser,
  verifyUserBillingRecord,
  updateUser
} from '@/service/api';
import DataFreshness from '@/components/project/data-freshness.vue';
import DeviceTable from '@/components/project/device-table.vue';
import TrafficTrendChart from '@/components/project/traffic-trend-chart.vue';
import TrafficValue from '@/components/project/traffic-value.vue';
import type { DropdownOption } from 'naive-ui';
import { billingServiceInterval, billingTimestampISO } from '@/utils/billing';

defineOptions({ name: 'UserManagement' });

const loading = ref(false);
const errorMessage = ref('');
const groups = ref<Api.Central.UserGroupSummary[]>([]);
const groupTotal = ref(0);
const stats = ref<Api.Central.UserListStats | null>(null);
const dataAt = ref('');
const userRows = ref<Api.Central.UserSummary[]>([]);
const userTotal = ref(0);
const userLoading = ref(false);
const userError = ref('');
const totalUsers = ref<number | null>(null);
const pendingRenewalCount = ref<number | null>(null);
const pendingRenewalAmount = ref(0);
const viewMode = ref<'table' | 'grouped'>('table');
const detailTab = ref<'overview' | 'billing' | 'path' | 'device'>('overview');
const businessEditing = ref(false);
const renewModalVisible = ref(false);
const renewCandidate = ref<Api.Central.UserRenewalCandidate | null>(null);
const renewSaving = ref(false);
const renewForm = reactive({
  amount: 0 as number | null,
  paidAt: new Date().toISOString().slice(0, 10)
});
const expandedGroups = ref(new Set<string>());
interface GroupUsersState {
  rows: Api.Central.UserSummary[];
  total: number;
  page: number;
  pageSize: number;
  loading: boolean;
  error: string;
}
const groupUsers = reactive<Record<string, GroupUsersState>>({});
const drawerVisible = ref(false);
const detailLoading = ref(false);
const detailError = ref('');
const saving = ref(false);
const detail = ref<Api.Central.UserDetail | null>(null);
const selectedUserID = ref('');
const editForm = reactive({
  displayName: '',
  billingType: 'paid' as Api.Central.BillingType,
  freeReason: '',
  billingCycle: 'monthly' as Api.Central.BillingCycle,
  billingAmount: 0 as number | null,
  notes: ''
});
const initialOrderForm = reactive({
  amount: null as number | null,
  serviceFrom: null as string | null,
  serviceTo: null as string | null,
  serviceToSyncedAt: null as string | null,
  paidAt: new Date().toISOString().slice(0, 10),
  notes: ''
});
const initialOrderSaving = ref(false);
const orderModalVisible = ref(false);
const orderSaving = ref(false);
const orderForm = reactive({
  billingCycle: 'monthly' as Api.Central.BillingCycle,
  amount: null as number | null,
  serviceFrom: null as number | null,
  serviceTo: null as number | null,
  paidAt: new Date().toISOString().slice(0, 10),
  orderType: 'renewal' as 'initial' | 'renewal' | 'recovery',
  notes: ''
});
const verifyModalVisible = ref(false);
const verifySaving = ref(false);
const verifyRecord = ref<Api.Central.UserBillingRecord | null>(null);
const verifyPaidAt = ref('');
const verifyNotes = ref('');
const importModalVisible = ref(false);
const importText = ref('');
const importPreview = ref<Api.Central.BillingImportResult | null>(null);
const importLoading = ref(false);
const pathAssets = ref<Api.Central.UserPathAssets | null>(null);
const pathLandingInbounds = ref<Api.Central.NodeInboundDetail[]>([]);
const pathAssetsLoading = ref(false);
const pathLandingInboundsLoading = ref(false);
const pathAssetsError = ref('');
const pathSaving = ref(false);
const resettingUser = ref(false);
const pathError = ref('');
const pathMode = ref<'relay' | 'landing' | 'external'>('relay');
const pathLandingNodeID = ref<string | null>(null);
const pathLandingInboundID = ref<string | null>(null);
const pathExitIPIDs = ref<string[]>([]);
const pathNotes = ref('');
const route = useRoute();
const router = useRouter();
let pathAssetsRequestID = 0;

const defaultGroupPageSize = 20;
const filters = reactive({
  page: 1,
  page_size: defaultGroupPageSize,
  keyword: '',
  status: 'all',
  billingType: 'all',
  nodeID: ''
});

const statusOptions = [
  { label: '全部', value: 'all' },
  { label: '启用中', value: 'active' },
  { label: '即将到期', value: 'expiring' },
  { label: '已到期', value: 'expired' },
  { label: '已停用', value: 'disabled' }
];

const billingTypeOptions = [
  { label: '全部类型', value: 'all' },
  { label: '付费用户', value: 'paid' },
  { label: '免费用户', value: 'free' }
];

const nodeOptions = computed(() => [
  { label: '全部线路机', value: '' },
  ...groups.value.map(group => ({ label: group.nodeName, value: group.nodeId }))
]);

const billingCycleOptions = [
  { label: '月付', value: 'monthly' },
  { label: '年付', value: 'annual' }
];

const billingTypeFormOptions = [
  { label: '收费用户', value: 'paid' },
  { label: '免费用户', value: 'free' }
];

const detailTabs = [
  { key: 'overview', label: '概览' },
  { key: 'billing', label: '账务' },
  { key: 'path', label: '路径' },
  { key: 'device', label: '设备与流量' }
] as const;

const moreMenuOptions = computed<DropdownOption[]>(() => [
  { label: '新增订单', key: 'order' },
  { label: '导入历史账单', key: 'import' },
  { type: 'divider', key: 'divider' },
  { label: '更换客户 · 不可恢复', key: 'replace' }
]);

const hasActiveFilters = computed(
  () =>
    Boolean(filters.keyword.trim()) ||
    filters.status !== 'all' ||
    filters.billingType !== 'all' ||
    filters.nodeID !== ''
);

const userPageCount = computed(() => Math.max(1, Math.ceil(userTotal.value / Math.max(filters.page_size, 1))));

const pendingCandidates = computed(() =>
  (detail.value?.renewalCandidates || []).filter(candidate => candidate.status === 'pending')
);

const enabledClientCount = computed(() => (detail.value?.clients || []).filter(client => client.enabled).length);

const lastOnlineAt = computed(() => {
  const clients = detail.value?.clients || [];
  const times = clients
    .map(client => (client.lastOnlineAt ? new Date(client.lastOnlineAt).getTime() : 0))
    .filter(value => Number.isFinite(value) && value > 0);
  return times.length ? new Date(Math.max(...times)).toISOString() : null;
});

const detailInitial = computed(() => {
  const name = detail.value?.displayName || selectedUserID.value || '?';
  return name.trim().slice(0, 1).toUpperCase();
});

const expiringHint = computed(() => {
  if (!stats.value) return '等待中央数据';
  return stats.value.expiring > 0 ? '无宽限期，到期即失效' : '当前没有 7 天内到期用户';
});

const distribution = computed(() => {
  const total = totalUsers.value || 0;
  const active = stats.value?.active || 0;
  const expiring = stats.value?.expiring || 0;
  const stopped = Math.max(0, total - active - expiring);
  const pending = pendingRenewalCount.value || 0;
  const percent = (value: number) => (total > 0 ? Math.min(100, Math.round((value / total) * 100)) : 0);
  return [
    { key: 'active', label: '启用中', value: active, tone: 'is-ok', percent: percent(active) },
    { key: 'expiring', label: '7 天内到期', value: expiring, tone: 'is-warn', percent: percent(expiring) },
    { key: 'stopped', label: '已到期 / 停用', value: stopped, tone: 'is-danger', percent: percent(stopped) },
    { key: 'pending', label: '待确认续费', value: pending, tone: 'is-neutral', percent: percent(pending) }
  ];
});

const pathModeOptions = [
  { label: '线路机直出', value: 'relay' },
  { label: '经落地机', value: 'landing' },
  { label: '独立 S5', value: 'external' }
];

const pathLandingNodeOptions = computed(() =>
  (pathAssets.value?.landingNodes || [])
    .filter(node => node.type === 'landing' && node.enabled)
    .map(node => ({ label: `${node.name} · ${node.host || '地址未知'}`, value: node.id }))
);

const pathInboundOptions = computed(() =>
  pathLandingInbounds.value
    .filter(inbound => inbound.enabled && inbound.status === 'active')
    .map(inbound => ({
      label: `${inbound.tag || inbound.remark || inbound.remoteId} · ${inbound.protocol || '--'}:${inbound.port || '--'}`,
      value: inbound.id
    }))
);

const pathExitIpOptions = computed(() => {
  let assets: Api.Central.ExitIpSummary[] = [];
  if (pathMode.value === 'relay') {
    assets = pathAssets.value?.relayExitIps || [];
  } else if (pathMode.value === 'landing') {
    assets = pathAssets.value?.landingNodes.find(node => node.id === pathLandingNodeID.value)?.exitIps || [];
  } else {
    assets = pathAssets.value?.externalExitIps || [];
  }
  const options = assets.map(item => ({
    label: pathExitIpOptionLabel(item),
    value: item.id,
    disabled: pathExitIpIsUnavailable(item)
  }));
  const currentPath = detail.value?.path;
  const samePathContext =
    currentPath?.mode === pathMode.value &&
    (pathMode.value !== 'landing' || currentPath.landingNodeId === pathLandingNodeID.value);
  if (samePathContext) {
    const optionIDs = new Set(options.map(option => option.value));
    for (const item of currentPath?.exitIps || []) {
      if (optionIDs.has(item.id)) continue;
      options.unshift({
        label: pathExitIpOptionLabel(item),
        value: item.id,
        disabled: pathExitIpIsUnavailable(item)
      });
      optionIDs.add(item.id);
    }
  }
  return options;
});

function exitIpExpired(value?: string | null) {
  if (!value) return false;
  const expiry = value.slice(0, 10);
  const now = new Date();
  const today = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-${String(now.getDate()).padStart(2, '0')}`;
  return expiry < today;
}

function pathExitIpIsUnavailable(item: { enabled?: boolean; validTo?: string | null }) {
  return item.enabled === false || exitIpExpired(item.validTo);
}

function pathExitIpStatusLabel(item: { enabled?: boolean; validTo?: string | null }) {
  if (item.enabled === false) return '已停用';
  if (exitIpExpired(item.validTo)) return '已到期';
  return '';
}

function pathExitIpOptionLabel(item: {
  id: string;
  address: string;
  sourceType: 'node' | 's5';
  ownerNodeName?: string | null;
  provider?: string | null;
  enabled?: boolean;
  validTo?: string | null;
}) {
  const address = item.address || item.id;
  const source = item.sourceType === 's5' ? '独立 S5' : item.ownerNodeName || '节点出口';
  const provider = item.provider ? ` · ${item.provider}` : '';
  const status = pathExitIpStatusLabel(item);
  return `${address} · ${source}${provider}${status ? ` · ${status}` : ''}`;
}

function formatDate(value?: string | null) {
  if (!value) return '--';
  const date = new Date(value);
  return Number.isNaN(date.getTime()) ? '时间未知' : date.toLocaleString('zh-CN', { hour12: false });
}

function formatExpiry(value: string | null | undefined, status: Api.Central.UserStatus) {
  if (!value) return status === 'active' ? '长期有效' : '时间未知';
  return formatDate(value);
}

function pathModeLabel(mode?: Api.Central.UserSummary['pathMode']) {
  if (mode === 'relay') return '线路机直出';
  if (mode === 'landing') return '经落地机';
  if (mode === 'external') return '独立 S5';
  return '未配置';
}

/** 状态文案统一（有效 → 启用中） */
function formatQuota(value?: number | null) {
  if (value === undefined || value === null || !Number.isFinite(value) || value <= 0) return '--';
  const units = ['B', 'KB', 'MB', 'GB', 'TB', 'PB'];
  let amount = value;
  let index = 0;
  while (amount >= 1024 && index < units.length - 1) {
    amount /= 1024;
    index += 1;
  }
  return `${amount.toFixed(amount >= 100 ? 0 : amount >= 10 ? 1 : 2)} ${units[index]}`;
}

function quotaUsageHint(row: Api.Central.UserSummary) {
  const quota = Number(row.quotaBytes || 0);
  if (!quota) return '';
  const used = Number(row.trafficBytes || 0);
  const percent = Math.min(999, Math.round((used / quota) * 100));
  return `已用 ${percent}%`;
}

function billingHint(row: Api.Central.UserSummary) {
  if (row.billingType === 'free') {
    return row.billingAmount ? `免费 · 参考 ¥${Number(row.billingAmount).toFixed(2)}` : '免费用户';
  }
  const amount = Number(row.billingAmount || 0);
  if (!amount) return '金额未设置';
  const cycle = row.billingCycle === 'annual' ? '年' : '月';
  return `¥${amount.toFixed(2)} / ${cycle}`;
}

function statusLabel(status: Api.Central.UserStatus) {
  if (status === 'active') return '启用中';
  if (status === 'expiring') return '即将到期';
  if (status === 'expired') return '已到期';
  if (status === 'disabled') return '已停用';
  return '状态未知';
}

function statusBadgeClass(status: Api.Central.UserStatus) {
  if (status === 'active') return 'is-ok';
  if (status === 'expiring') return 'is-warn';
  if (status === 'expired') return 'is-danger';
  if (status === 'disabled') return 'is-neutral';
  return 'is-violet';
}

function nodeDotClass(status: Api.Central.NodeStatus) {
  if (status === 'online') return 'is-ok';
  if (status === 'degraded' || status === 'unknown') return 'is-warn';
  if (status === 'offline' || status === 'disabled') return 'is-danger';
  return 'is-warn';
}

function nodeBadgeClass(status: Api.Central.NodeStatus) {
  if (status === 'online') return 'is-ok';
  if (status === 'degraded' || status === 'unknown') return 'is-warn';
  if (status === 'offline') return 'is-danger';
  return 'is-neutral';
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

function daysUntil(value?: string | null) {
  if (!value) return 0;
  const target = new Date(value).getTime();
  if (!Number.isFinite(target)) return 0;
  return Math.ceil((target - Date.now()) / 86_400_000);
}

function expiryHint(row: Api.Central.UserSummary) {
  if (!row.expiresAt) return row.status === 'active' ? '不限期' : '时间未知';
  if (row.status === 'expired') return `已逾期 ${Math.abs(daysUntil(row.expiresAt))} 天`;
  if (row.status === 'disabled') return '停用前最后同步';
  if (row.status === 'expiring') return `还剩 ${Math.max(0, daysUntil(row.expiresAt))} 天`;
  return `还剩 ${Math.max(0, daysUntil(row.expiresAt))} 天`;
}

function expiryBadgeClass(value: { expiresAt?: string | null; status: Api.Central.UserStatus }) {
  if (!value.expiresAt) return value.status === 'active' ? 'is-ok' : 'is-neutral';
  const days = daysUntil(value.expiresAt);
  if (days < 0) return 'is-danger';
  if (days <= 7) return 'is-warn';
  return 'is-ok';
}

function expiryText(value: { expiresAt?: string | null; status: Api.Central.UserStatus }) {
  if (!value.expiresAt) return value.status === 'active' ? '长期有效' : '时间未知';
  const days = daysUntil(value.expiresAt);
  if (days < 0) return `已逾期 ${Math.abs(days)} 天`;
  return `还剩 ${days} 天`;
}

function exitIpHint(row: Api.Central.UserSummary) {
  if (row.pathMode === 'landing') return row.landingNodeName ? `经 ${row.landingNodeName}` : '经落地机';
  if (row.pathMode === 'external') return '独立 S5';
  if (row.pathMode === 'relay') return '线路机直出';
  return '尚未配置路径';
}

function groupSummary(group: Api.Central.UserGroupSummary) {
  const statsValue = group.stats;
  const parts = [`启用中 ${statsValue.active}`, `即将到期 ${statsValue.expiring}`];
  if (statsValue.expired) parts.push(`已到期 ${statsValue.expired}`);
  if (statsValue.disabled) parts.push(`已停用 ${statsValue.disabled}`);
  return parts.join(' · ');
}

function groupRows(nodeID: string) {
  return groupUsers[nodeID]?.rows || [];
}

function groupUserTotal(nodeID: string) {
  return groupUsers[nodeID]?.total || 0;
}

function groupPage(nodeID: string) {
  return groupUsers[nodeID]?.page || 1;
}

function groupPageCount(nodeID: string) {
  const state = groupUsers[nodeID];
  if (!state) return 1;
  return Math.max(1, Math.ceil(state.total / Math.max(state.pageSize, 1)));
}

function changeGroupPage(nodeID: string, page: number) {
  void loadGroupUsers(nodeID, page);
}

function pathExitIPText(path?: Api.Central.UserPathDetail | null) {
  if (!path) return '--';
  const addresses = path.exitIps?.map(item => item.address).filter(Boolean) || [];
  return addresses.join('、') || path.exitIpAddress || '--';
}

function toDateInputValue(value?: string | null) {
  if (!value) return null;
  if (/^\d{4}-\d{2}-\d{2}/.test(value)) return value.slice(0, 10);
  const date = new Date(value);
  return Number.isNaN(date.getTime()) ? null : date.toISOString().slice(0, 10);
}

function copyDetailToForm(value: Api.Central.UserDetail) {
  editForm.displayName = value.displayName;
  editForm.billingType = value.billingType || 'paid';
  editForm.freeReason = value.freeReason || '';
  editForm.billingCycle = value.billingCycle || 'monthly';
  editForm.billingAmount = value.billingAmount ?? value.monthlyFee;
  editForm.notes = value.notes || '';
  initialOrderForm.amount = value.billingAmount ?? value.monthlyFee;
  initialOrderForm.serviceTo = toDateInputValue(value.expiresAt);
  initialOrderForm.serviceToSyncedAt = value.expiresAt || null;
}

function resetInitialOrderForm() {
  initialOrderForm.serviceFrom = null;
  initialOrderForm.serviceTo = null;
  initialOrderForm.serviceToSyncedAt = null;
  initialOrderForm.paidAt = new Date().toISOString().slice(0, 10);
  initialOrderForm.notes = '';
}

function dayToRFC3339(value: string | null, endOfDay = false) {
  return value ? `${value}T${endOfDay ? '23:59:59' : '00:00:00'}Z` : '';
}

async function createInitialOrder() {
  if (!detail.value) return;
  if (detail.value.billingRecords.length > 0) {
    window.$message?.warning('该用户已有收费记录，续费请使用续费候选或新增订单流程');
    return;
  }
  if (initialOrderForm.amount === null || !Number.isFinite(initialOrderForm.amount) || initialOrderForm.amount < 0) {
    window.$message?.warning('请填写有效的首笔订单金额');
    return;
  }
  const serviceFrom = dayToRFC3339(initialOrderForm.serviceFrom);
  const serviceTo = initialOrderForm.serviceToSyncedAt || '';
  if (!serviceFrom || !serviceTo || new Date(serviceTo) <= new Date(serviceFrom)) {
    window.$message?.warning('请填写正确的服务开始和结束日期');
    return;
  }
  initialOrderSaving.value = true;
  const { data, error } = await createUserBillingRecord(detail.value.id, {
    billingCycle: editForm.billingCycle,
    amount: initialOrderForm.amount,
    serviceFrom,
    serviceTo,
    paidAt: dayToRFC3339(initialOrderForm.paidAt, true),
    orderType: 'initial',
    notes: initialOrderForm.notes.trim()
  });
  initialOrderSaving.value = false;
  if (error || !data) {
    window.$message?.error('首笔订单保存失败，请检查服务区间');
    return;
  }
  detail.value = data;
  copyDetailToForm(data);
  window.$message?.success('首笔订单已确认并计入实收');
  void reloadAll();
}

function resetOrderForm() {
  orderForm.billingCycle = detail.value?.billingCycle || 'monthly';
  orderForm.amount = detail.value?.billingAmount ?? detail.value?.monthlyFee ?? 0;
  const records = (detail.value?.billingRecords || []).filter(record => record.status !== 'cancelled');
  Object.assign(orderForm, billingServiceInterval(records, detail.value?.expiresAt));
  orderForm.paidAt = new Date().toISOString().slice(0, 10);
  orderForm.orderType = records.length === 0 ? 'initial' : 'renewal';
  orderForm.notes = '';
}

function openOrderModal() {
  resetOrderForm();
  orderModalVisible.value = true;
}

async function createAdditionalOrder() {
  if (!detail.value) return;
  if (orderForm.amount === null || !Number.isFinite(orderForm.amount) || orderForm.amount < 0) {
    window.$message?.warning('请填写有效的订单金额');
    return;
  }
  const serviceFrom = billingTimestampISO(orderForm.serviceFrom);
  const serviceTo = billingTimestampISO(orderForm.serviceTo);
  if (!serviceFrom || !serviceTo || new Date(serviceTo) <= new Date(serviceFrom)) {
    window.$message?.warning('请填写正确的服务区间');
    return;
  }
  orderSaving.value = true;
  const { data, error } = await createUserBillingRecord(detail.value.id, {
    billingCycle: orderForm.billingCycle,
    amount: orderForm.amount,
    serviceFrom,
    serviceTo,
    paidAt: dayToRFC3339(orderForm.paidAt, true),
    orderType: orderForm.orderType,
    notes: orderForm.notes.trim()
  });
  orderSaving.value = false;
  if (error || !data) {
    window.$message?.error('订单保存失败，请检查服务区间是否与已有记录重叠');
    return;
  }
  detail.value = data;
  copyDetailToForm(data);
  orderModalVisible.value = false;
  window.$message?.success('订单已保存');
  void reloadAll();
}

function openVerifyModal(record: Api.Central.UserBillingRecord) {
  verifyRecord.value = record;
  verifyPaidAt.value = toDateInputValue(record.paidAt) || new Date().toISOString().slice(0, 10);
  verifyNotes.value = record.notes || '';
  verifyModalVisible.value = true;
}

async function submitVerifyRecord() {
  if (!detail.value || !verifyRecord.value) return;
  if (!verifyPaidAt.value) {
    window.$message?.warning('请填写实际收款日期');
    return;
  }
  verifySaving.value = true;
  const { data, error } = await verifyUserBillingRecord(detail.value.id, verifyRecord.value.id, {
    paidAt: dayToRFC3339(verifyPaidAt.value, true),
    notes: verifyNotes.value.trim()
  });
  verifySaving.value = false;
  if (error || !data) {
    window.$message?.error('历史订单核验失败');
    return;
  }
  detail.value = data;
  copyDetailToForm(data);
  verifyModalVisible.value = false;
  window.$message?.success('历史订单已核验并计入财务');
  void reloadAll();
}

function cancelRecord(record: Api.Central.UserBillingRecord) {
  if (!detail.value) return;
  const cancel = async () => {
    const { data, error } = await cancelUserBillingRecord(detail.value!.id, record.id, { notes: '管理员取消错误或重复的历史账单' });
    if (error || !data) {
      window.$message?.error('取消订单失败');
      return;
    }
    detail.value = data;
    copyDetailToForm(data);
    window.$message?.success('订单已取消，财务统计已排除');
    void reloadAll();
  };
  if (!window.$dialog) {
    void cancel();
    return;
  }
  window.$dialog.warning({
    title: '取消订单',
    content: '取消后该订单不会再计入收入，但记录仍会保留在历史中。确定继续吗？',
    positiveText: '取消订单',
    negativeText: '返回',
    onPositiveClick: cancel
  });
}

function parseCsvLine(line: string) {
  const values: string[] = [];
  let value = '';
  let quoted = false;
  for (let index = 0; index < line.length; index += 1) {
    const char = line[index];
    if (char === '"') {
      if (quoted && line[index + 1] === '"') {
        value += '"';
        index += 1;
      } else quoted = !quoted;
    } else if (char === ',' && !quoted) {
      values.push(value.trim());
      value = '';
    } else value += char;
  }
  values.push(value.trim());
  return values;
}

function parseBillingCsv(value: string) {
  const lines = value.split(/\r?\n/).map(line => line.trim()).filter(Boolean);
  if (lines.length < 2) throw new Error('CSV 至少需要一行表头和一行数据');
  const headers = parseCsvLine(lines[0]);
  const required = ['billingCycle', 'amount', 'serviceFrom', 'serviceTo'];
  for (const name of required) {
    if (!headers.includes(name)) throw new Error(`CSV 缺少列：${name}`);
  }
  return lines.slice(1).map((line, rowIndex) => {
    const values = parseCsvLine(line);
    const row = Object.fromEntries(headers.map((header, index) => [header, values[index] || ''])) as Record<string, string>;
    if (!row.userId && (!row.nodeId || !row.remoteInboundId)) throw new Error(`第 ${rowIndex + 2} 行缺少 userId 或 nodeId + remoteInboundId`);
    return {
      userId: row.userId || undefined,
      nodeId: row.nodeId || undefined,
      remoteInboundId: row.remoteInboundId || undefined,
      billingCycle: row.billingCycle as Api.Central.BillingCycle,
      amount: Number(row.amount),
      serviceFrom: row.serviceFrom,
      serviceTo: row.serviceTo,
      paidAt: row.paidAt || undefined,
      orderType: (row.orderType || 'renewal') as 'initial' | 'renewal' | 'recovery',
      verified: row.verified === '' || ['true', '1', 'yes', '是'].includes(row.verified.toLowerCase()),
      notes: row.notes || undefined
    };
  });
}

async function previewBillingImport() {
  importLoading.value = true;
  importPreview.value = null;
  try {
    const records = parseBillingCsv(importText.value);
    const { data, error } = await importUserBillingRecords({ records, dryRun: true });
    if (error || !data) window.$message?.error('导入预览失败，请检查 CSV 内容');
    else importPreview.value = data;
  } catch (error) {
    window.$message?.error(error instanceof Error ? error.message : 'CSV 格式错误');
  } finally {
    importLoading.value = false;
  }
}

async function submitBillingImport() {
  if (!importPreview.value || importPreview.value.errors.length > 0) return;
  importLoading.value = true;
  try {
    const records = parseBillingCsv(importText.value);
    const { data, error } = await importUserBillingRecords({ records, dryRun: false });
    if (error || !data) {
      window.$message?.error('历史账单导入失败');
      return;
    }
    importModalVisible.value = false;
    importPreview.value = null;
    importText.value = '';
    window.$message?.success(`已导入 ${data.imported} 笔历史订单${data.unverified ? `，其中 ${data.unverified} 笔待核实` : ''}`);
    void reloadAll();
    if (selectedUserID.value) void openDetail(selectedUserID.value);
  } catch (error) {
    window.$message?.error(error instanceof Error ? error.message : 'CSV 格式错误');
  } finally {
    importLoading.value = false;
  }
}

function syncPathForm(value: Api.Central.UserDetail) {
  const path = value.path;
  pathMode.value = path?.mode || 'relay';
  pathLandingNodeID.value = path?.landingNodeId || null;
  pathLandingInboundID.value = path?.landingInboundId || null;
  pathExitIPIDs.value = path?.exitIpIds?.length ? [...path.exitIpIds] : path?.exitIpId ? [path.exitIpId] : [];
  pathNotes.value = path?.notes || '';
  pathError.value = '';
  pathLandingInbounds.value = pathLandingNodeID.value
    ? pathAssets.value?.landingNodes.find(node => node.id === pathLandingNodeID.value)?.inbounds || []
    : [];
}

async function loadPathAssets(userID = selectedUserID.value) {
  if (!userID) return;
  const requestID = ++pathAssetsRequestID;
  pathAssetsLoading.value = true;
  pathAssetsError.value = '';
  const { data, error } = await fetchUserPathAssets(userID);
  if (requestID !== pathAssetsRequestID) return;
  if (error || !data) {
    pathAssets.value = null;
    pathLandingInbounds.value = [];
    pathAssetsError.value = '无法读取用户路径资源，请确认中央后端已启动并检查线路机主 Inbound。';
  } else {
    pathAssets.value = data;
    if (pathLandingNodeID.value) {
      pathLandingInbounds.value = data.landingNodes.find(node => node.id === pathLandingNodeID.value)?.inbounds || [];
    }
  }
  pathAssetsLoading.value = false;
}

function handlePathModeChange(value: 'relay' | 'landing' | 'external') {
  pathMode.value = value;
  pathLandingNodeID.value = null;
  pathLandingInboundID.value = null;
  pathLandingInbounds.value = [];
  pathExitIPIDs.value = [];
  pathError.value = '';
}

function handlePathLandingNodeChange(value: string | null) {
  pathLandingNodeID.value = value;
  pathLandingInboundID.value = null;
  pathExitIPIDs.value = [];
  pathError.value = '';
  pathLandingInboundsLoading.value = false;
  pathLandingInbounds.value = value
    ? pathAssets.value?.landingNodes.find(node => node.id === value)?.inbounds || []
    : [];
}

async function savePathAssignment() {
  if (!detail.value) return;
  if (pathExitIPIDs.value.length === 0) {
    pathError.value = '请选择至少一个出口 IP；线路机直出也需要明确指定出口 IP。';
    return;
  }
  if (pathMode.value === 'landing' && !pathLandingNodeID.value) {
    pathError.value = '经落地机模式必须选择落地机。';
    return;
  }
  if (pathMode.value === 'landing' && !pathLandingInboundID.value) {
    pathError.value = '经落地机模式必须选择落地 Inbound。';
    return;
  }
  pathSaving.value = true;
  pathError.value = '';
  const { data, error } = await assignUserPath(detail.value.id, {
    landingNodeId: pathMode.value === 'landing' ? pathLandingNodeID.value : null,
    landingInboundId: pathMode.value === 'landing' ? pathLandingInboundID.value : null,
    exitIpIds: pathExitIPIDs.value,
    exitIpId: pathExitIPIDs.value[0],
    notes: pathNotes.value.trim() || null
  });
  pathSaving.value = false;
  if (error || !data) {
    pathError.value = '用户路径保存失败，请检查节点和出口 IP 是否仍然启用。';
    return;
  }
  detail.value = data;
  syncPathForm(data);
  window.$message?.success('用户路径已保存');
  void reloadAll();
}

async function removePathAssignment() {
  if (!detail.value || pathSaving.value) return;
  pathSaving.value = true;
  pathError.value = '';
  const { data, error } = await clearUserPath(detail.value.id);
  pathSaving.value = false;
  if (error || !data) {
    pathError.value = '解除用户路径失败，请稍后重试。';
    return;
  }
  detail.value = data;
  syncPathForm(data);
  window.$message?.success('已解除用户路径');
  void reloadAll();
}

function confirmRemovePathAssignment() {
  const remove = () => void removePathAssignment();
  if (!window.$dialog) {
    remove();
    return;
  }
  window.$dialog.warning({
    title: '解除用户路径',
    content: '解除后用户将不再有当前的中央路径配置，历史记录仍会保留。确定继续吗？',
    positiveText: '解除',
    negativeText: '取消',
    onPositiveClick: remove
  });
}

function confirmReplaceCurrentUser() {
  const inboundID = detail.value?.inbound.id;
  if (!inboundID || resettingUser.value) return;
  const reset = async () => {
    resettingUser.value = true;
    const { data, error } = await resetInboundUser(inboundID);
    resettingUser.value = false;
    if (error || !data) {
      window.$message?.error('更换客户失败，请稍后重试');
      return;
    }
    window.$message?.success('已关闭旧用户并创建新用户，旧账务记录已保留');
    drawerVisible.value = false;
    void reloadAll();
    void openDetail(data.newUserId);
  };
  if (!window.$dialog) {
    void reset();
    return;
  }
  window.$dialog.warning({
    title: '更换当前 Inbound 的客户',
    content: '确认后将关闭当前用户的路径和流量，仅保留账务记录，并创建新用户。此操作不可恢复。',
    positiveText: '确认更换',
    negativeText: '取消',
    maskClosable: false,
    onPositiveClick: reset
  });
}

async function openDetail(id: string) {
  selectedUserID.value = id;
  pathAssetsRequestID += 1;
  drawerVisible.value = true;
  detailLoading.value = true;
  detailError.value = '';
  detail.value = null;
  detailTab.value = 'overview';
  businessEditing.value = false;
  resetInitialOrderForm();
  pathAssets.value = null;
  pathLandingInbounds.value = [];
  pathAssetsLoading.value = false;
  if (route.query.userId !== id) {
    void router.replace({ query: { ...route.query, userId: id } });
  }
  const { data, error } = await fetchUserDetail(id);
  if (error || !data) {
    detailError.value = '无法读取用户详情，请确认中央后端已启动后重试。';
    pathAssetsError.value = '';
  } else {
    detail.value = data;
    copyDetailToForm(data);
    syncPathForm(data);
    void loadPathAssets(id);
  }
  detailLoading.value = false;
}

async function saveBusinessFields() {
  if (!detail.value) return;
  const displayName = editForm.displayName.trim();
  const billingAmount = editForm.billingAmount;
  if (!displayName) {
    window.$message?.warning('请填写业务用户名称');
    return;
  }
  if (billingAmount === null || !Number.isFinite(billingAmount) || billingAmount < 0) {
    window.$message?.warning('收费金额必须是不小于 0 的数字');
    return;
  }
  saving.value = true;
  const { data, error } = await updateUser(detail.value.id, {
    displayName,
    billingType: editForm.billingType,
    freeReason: editForm.billingType === 'free' ? editForm.freeReason.trim() : '',
    billingCycle: editForm.billingCycle,
    billingAmount,
    currency: 'CNY',
    notes: editForm.notes.trim()
  });
  saving.value = false;
  if (error || !data) return;
  detail.value = data;
  copyDetailToForm(data);
  businessEditing.value = false;
  window.$message?.success('中央业务信息已保存');
  window.dispatchEvent(new Event('xpanel-events-updated'));
  void reloadAll();
}

function cycleLabel(cycle: Api.Central.BillingCycle) {
  return cycle === 'annual' ? '年付' : '月付';
}

function orderTypeLabel(orderType: Api.Central.UserBillingRecord['orderType']) {
  return orderType === 'renewal' ? '续费' : orderType === 'recovery' ? '逾期恢复' : '首购';
}

function formatMoney(amount: number, currency = 'CNY') {
  return `${currency} ${Number(amount || 0).toFixed(2)}`;
}

function openRenewConfirm(candidate: Api.Central.UserRenewalCandidate) {
  renewCandidate.value = candidate;
  renewForm.amount = candidate.suggestedAmount;
  renewForm.paidAt = new Date().toISOString().slice(0, 10);
  renewModalVisible.value = true;
}

async function submitRenewConfirm() {
  if (!detail.value || !renewCandidate.value) return;
  const amount = renewForm.amount;
  if (amount === null || !Number.isFinite(amount) || amount < 0) {
    window.$message?.warning('请填写有效的续费金额');
    return;
  }
  if (!renewForm.paidAt) {
    window.$message?.warning('请填写实际收款日期');
    return;
  }
  renewSaving.value = true;
  const { data, error } = await confirmUserRenewal(detail.value.id, renewCandidate.value.id, {
    billingCycle: renewCandidate.value.billingCycle,
    amount,
    paidAt: dayToRFC3339(renewForm.paidAt, true)
  });
  renewSaving.value = false;
  if (error || !data) {
    window.$message?.error('续费确认失败，请稍后重试');
    return;
  }
  detail.value = data;
  copyDetailToForm(data);
  renewModalVisible.value = false;
  renewCandidate.value = null;
  window.$message?.success('续费已确认并计入财务');
  void reloadAll();
}

async function rejectRenewal(candidate: Api.Central.UserRenewalCandidate) {
  if (!detail.value) return;
  const { data, error } = await rejectUserRenewal(detail.value.id, candidate.id);
  if (error || !data) {
    window.$message?.error('忽略续费变更失败，请稍后重试');
    return;
  }
  detail.value = data;
  window.$message?.success('已标记为非收费变更');
  void reloadAll();
}

function ensureGroupUsers(nodeID: string) {
  if (!groupUsers[nodeID]) {
    groupUsers[nodeID] = { rows: [], total: 0, page: 1, pageSize: 20, loading: false, error: '' };
  }
  return groupUsers[nodeID];
}


function isGroupExpanded(nodeID: string) {
  return expandedGroups.value.has(nodeID);
}

function toggleGroup(group: Api.Central.UserGroupSummary) {
  const next = new Set(expandedGroups.value);
  if (next.has(group.nodeId)) {
    next.delete(group.nodeId);
  } else {
    next.add(group.nodeId);
    void loadGroupUsers(group.nodeId);
  }
  expandedGroups.value = next;
}

function expandAllGroups() {
  expandedGroups.value = new Set(groups.value.map(group => group.nodeId));
  void Promise.all(groups.value.map(group => loadGroupUsers(group.nodeId)));
}

function collapseAllGroups() {
  expandedGroups.value = new Set();
}

function nodeStatusLabel(status: Api.Central.NodeStatus) {
  if (status === 'online') return '在线';
  if (status === 'degraded') return '异常';
  if (status === 'offline') return '离线';
  if (status === 'disabled') return '停用';
  return '未知';
}

function nodeStatusType(status: Api.Central.NodeStatus): 'success' | 'warning' | 'error' | 'default' {
  if (status === 'online') return 'success';
  if (status === 'degraded' || status === 'unknown') return 'warning';
  if (status === 'offline' || status === 'disabled') return 'error';
  return 'default';
}

function syncStatusLabel(status: Api.Central.UserGroupSummary['syncStatus']) {
  if (status === 'success') return '同步成功';
  if (status === 'failed') return '同步异常';
  return '待同步';
}

async function loadGroupUsers(nodeID: string, requestedPage?: number) {
  const state = ensureGroupUsers(nodeID);
  if (requestedPage) state.page = requestedPage;
  state.loading = true;
  state.error = '';
  const { data, error } = await fetchUsers({
    page: state.page,
    page_size: state.pageSize,
    keyword: filters.keyword || undefined,
    status: filters.status === 'all' ? undefined : filters.status,
    billing_type: filters.billingType === 'all' ? undefined : (filters.billingType as Api.Central.BillingType),
    node_id: nodeID === '__unassigned__' ? '__unassigned__' : nodeID
  });
  if (error || !data) {
    state.rows = [];
    state.total = 0;
    state.error = '无法读取该线路机下的用户';
  } else {
    state.rows = data.items;
    state.total = data.total;
  }
  state.loading = false;
}

async function loadGroups() {
  loading.value = true;
  errorMessage.value = '';
  const { data, error } = await fetchUserGroups({
    page: 1,
    page_size: 100,
    keyword: filters.keyword || undefined,
    status: filters.status === 'all' ? undefined : filters.status,
    billing_type: filters.billingType === 'all' ? undefined : (filters.billingType as Api.Central.BillingType),
    node_id: filters.nodeID || undefined
  });
  if (error || !data) {
    errorMessage.value = '中央后端暂不可用，无法读取用户数据';
    groups.value = [];
    groupTotal.value = 0;
    stats.value = null;
    dataAt.value = '';
  } else {
    groups.value = data.items;
    groupTotal.value = data.total;
    stats.value = data.stats;
    dataAt.value = data.dataAt || '';
    const available = new Set(groups.value.map(group => group.nodeId));
    expandedGroups.value = new Set([...expandedGroups.value].filter(nodeID => available.has(nodeID)));
    for (const nodeID of Object.keys(groupUsers)) {
      if (!available.has(nodeID)) Reflect.deleteProperty(groupUsers, nodeID);
    }
    if (viewMode.value === 'grouped' && !expandedGroups.value.size && groups.value.length) {
      expandedGroups.value = new Set([groups.value[0].nodeId]);
      void loadGroupUsers(groups.value[0].nodeId);
    }
  }
  loading.value = false;
}

async function loadUsers() {
  userLoading.value = true;
  userError.value = '';
  const { data, error } = await fetchUsers({
    page: filters.page,
    page_size: filters.page_size,
    keyword: filters.keyword || undefined,
    status: filters.status === 'all' ? undefined : filters.status,
    billing_type: filters.billingType === 'all' ? undefined : (filters.billingType as Api.Central.BillingType),
    node_id: filters.nodeID || undefined
  });
  if (error || !data) {
    userRows.value = [];
    userTotal.value = 0;
    userError.value = '无法读取用户列表，请确认中央后端已启动后重试。';
  } else {
    userRows.value = data.items;
    userTotal.value = data.total;
  }
  userLoading.value = false;
}

/**
 * 全局概览：后端用户列表 stats 只有 启用/即将到期/付费/免费，
 * 因此「总用户」用列表 total、「待确认续费」用待处理事件里的续费候选事件统计。
 */
async function loadOverview() {
  const [allUsers, pendingEvents] = await Promise.all([
    fetchUsers({ page: 1, page_size: 1 }),
    fetchEvents({ page: 1, page_size: 200, status: 'pending' })
  ]);
  if (!allUsers.error && allUsers.data) totalUsers.value = allUsers.data.total;
  if (!pendingEvents.error && pendingEvents.data) {
    const renewals = pendingEvents.data.items.filter(item => {
      const type = item.type || '';
      return type === 'renewal_candidate_detected' && item.status !== 'resolved' && item.status !== 'dismissed';
    });
    pendingRenewalCount.value = renewals.length;
    pendingRenewalAmount.value = renewals.reduce((sum, item) => {
      const amount = Number(item.payload?.suggestedAmount ?? 0);
      return sum + (Number.isFinite(amount) ? amount : 0);
    }, 0);
  }
}

async function reloadAll() {
  await Promise.all([loadGroups(), loadUsers(), loadOverview()]);
}

function resetGroupUserState() {
  expandedGroups.value = new Set();
  for (const state of Object.values(groupUsers)) {
    state.page = 1;
    state.pageSize = defaultGroupPageSize;
    state.rows = [];
    state.total = 0;
    state.error = '';
  }
}

function submitFilters() {
  filters.page = 1;
  resetGroupUserState();
  void (async () => {
    await loadGroups();
    await loadUsers();
    if (viewMode.value === 'grouped' && groups.value.length) {
      expandedGroups.value = new Set([groups.value[0].nodeId]);
      void loadGroupUsers(groups.value[0].nodeId);
    }
  })();
}

function setStatusFilter(value: string) {
  filters.status = value;
  submitFilters();
}

function setNodeFilter(value: string) {
  filters.nodeID = value;
  submitFilters();
}

function setViewMode(mode: 'table' | 'grouped') {
  viewMode.value = mode;
  const current = route.query.view === 'grouped' ? 'grouped' : 'table';
  if (current !== mode) {
    const query = { ...route.query };
    if (mode === 'grouped') query.view = 'grouped';
    else delete query.view;
    void router.replace({ query });
  }
  if (mode === 'grouped' && !expandedGroups.value.size && groups.value.length) {
    expandedGroups.value = new Set([groups.value[0].nodeId]);
    void loadGroupUsers(groups.value[0].nodeId);
  }
}

function changeUserPage(page: number) {
  filters.page = page;
  void loadUsers();
}

function resetFilters() {
  filters.keyword = '';
  filters.status = 'all';
  filters.billingType = 'all';
  filters.nodeID = '';
  filters.page_size = defaultGroupPageSize;
  submitFilters();
}

async function openOrderForUser(userID: string) {
  await openDetail(userID);
  if (detail.value) openOrderModal();
}

function setDetailTab(tab: 'overview' | 'billing' | 'path' | 'device') {
  detailTab.value = tab;
  void nextTick(() => {
    const target = document.getElementById(`u-sec-${tab}`);
    target?.scrollIntoView({ behavior: 'smooth', block: 'start' });
  });
}

function scrollToSection(tab: string) {
  setDetailTab(tab as 'overview' | 'billing' | 'path' | 'device');
}

function handleMoreMenu(key: string) {
  if (key === 'order') openOrderModal();
  else if (key === 'import') importModalVisible.value = true;
  else if (key === 'replace') confirmReplaceCurrentUser();
}

function closeDetail() {
  drawerVisible.value = false;
  businessEditing.value = false;
  if (route.query.userId) {
    const query = { ...route.query };
    delete query.userId;
    void router.replace({ query });
  }
}

function exportCurrentRows() {
  if (!userRows.value.length) return;
  const header = ['用户', 'Inbound', '状态', '到期时间', '累计流量(字节)', '线路机', '出口 IP', '计费类型'];
  const lines = userRows.value.map(row => [
    row.name,
    row.inboundTag,
    statusLabel(row.status),
    row.expiresAt || '',
    String(row.trafficBytes ?? ''),
    row.nodeName,
    row.exitIpAddress || '',
    row.billingType === 'free' ? '免费' : '收费'
  ]);
  const escape = (value: string) => `"${String(value ?? '').replace(/"/g, '""')}"`;
  const csv = [header, ...lines].map(cells => cells.map(escape).join(',')).join('\r\n');
  const blob = new Blob([`\ufeff${csv}`], { type: 'text/csv;charset=utf-8;' });
  const url = URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = `users-${new Date().toISOString().slice(0, 10)}.csv`;
  link.click();
  URL.revokeObjectURL(url);
}

onMounted(() => {
  void reloadAll();
  if (route.query.view === 'grouped') setViewMode('grouped');
  const userId = typeof route.query.userId === 'string' ? route.query.userId : '';
  if (userId) void openDetail(userId);
});

</script>

<template>
  <div class="u-page">
    <!-- ===================== 页头 ===================== -->
    <header class="u-head">
      <div class="u-head-main">
        <div class="u-eyebrow">
          <span>业务</span>
          <span class="u-eyebrow-sep">·</span>
          <span>用户管理</span>
          <span class="u-eyebrow-sep">·</span>
          <span>数据更新于 {{ relativeTime(dataAt) }}</span>
        </div>
        <h1 class="u-title">用户管理</h1>
        <p class="u-desc">
          业务用户 = 一台线路机上的一个用户 Inbound；Email 与 Client 只是这台用户下的设备凭证，不参与身份。
          到期时间由启用 Client 派生（取最早有限的到期值），不可手改；收费与路径由中央面板维护。
        </p>
      </div>
      <div class="u-head-actions">
        <DataFreshness :data-at="dataAt" />
        <button class="u-btn" type="button" :disabled="loading || userLoading" @click="reloadAll">
          <SvgIcon icon="mdi:refresh" class="u-btn-ico" />刷新
        </button>
      </div>
    </header>

    <!-- ===================== KPI ===================== -->
    <section class="u-kpis">
      <div class="u-kpi">
        <div class="u-kpi-top">
          <span class="u-kpi-label">总用户 <span class="u-dim">· 业务用户</span></span>
          <span class="u-kpi-icon is-blue"><SvgIcon icon="mdi:account-group-outline" /></span>
        </div>
        <div class="u-kpi-value">{{ totalUsers === null ? '--' : totalUsers }}<span class="u-kpi-unit">户</span></div>
        <div class="u-kpi-foot">{{ groupTotal }} 个线路机分组</div>
      </div>

      <div class="u-kpi">
        <div class="u-kpi-top">
          <span class="u-kpi-label">启用中</span>
          <span class="u-kpi-icon is-green"><SvgIcon icon="mdi:account-check-outline" /></span>
        </div>
        <div class="u-kpi-value">{{ stats ? stats.active : '--' }}<span class="u-kpi-unit">户</span></div>
        <div class="u-kpi-foot">
          付费 {{ stats ? stats.paid : '--' }} · 免费 {{ stats ? stats.free : '--' }}
        </div>
      </div>

      <div class="u-kpi is-clickable" title="查看 7 天内到期用户" @click="setStatusFilter('expiring')">
        <div class="u-kpi-top">
          <span class="u-kpi-label">7 天内到期</span>
          <span class="u-kpi-icon is-amber"><SvgIcon icon="mdi:clock-alert-outline" /></span>
        </div>
        <div class="u-kpi-value">{{ stats ? stats.expiring : '--' }}<span class="u-kpi-unit">户</span></div>
        <div class="u-kpi-foot u-warn">{{ expiringHint }}</div>
      </div>

      <div class="u-kpi">
        <div class="u-kpi-top">
          <span class="u-kpi-label">待确认续费</span>
          <span class="u-kpi-icon is-red"><SvgIcon icon="mdi:credit-card-clock-outline" /></span>
        </div>
        <div class="u-kpi-value">{{ pendingRenewalCount === null ? '--' : pendingRenewalCount }}<span class="u-kpi-unit">笔</span></div>
        <div class="u-kpi-foot">
          <span class="u-badge is-danger">{{ formatMoney(pendingRenewalAmount) }}</span>
          <span class="u-kpi-foot-text">确认后才计入财务</span>
        </div>
      </div>
    </section>

    <!-- ===================== 工具条 ===================== -->
    <section class="u-toolbar">
      <label class="u-input u-search">
        <SvgIcon icon="mdi:magnify" class="u-input-ico" />
        <input v-model="filters.keyword" placeholder="搜索用户、Inbound 或线路机" @keyup.enter="submitFilters" />
      </label>
      <NSelect
        v-model:value="filters.nodeID"
        :options="nodeOptions"
        size="small"
        class="u-select"
        style="width: 150px"
        @update:value="submitFilters"
      />
      <NSelect
        v-model:value="filters.billingType"
        :options="billingTypeOptions"
        size="small"
        class="u-select"
        style="width: 120px"
        @update:value="submitFilters"
      />
      <div class="u-seg">
        <button
          v-for="option in statusOptions"
          :key="option.value"
          type="button"
          :class="{ 'is-active': filters.status === option.value }"
          @click="setStatusFilter(option.value)"
        >
          {{ option.label }}
        </button>
      </div>
      <button class="u-btn u-btn-xs is-soft" type="button" @click="submitFilters">查询</button>
      <button v-if="hasActiveFilters" class="u-btn u-btn-xs" type="button" @click="resetFilters">重置</button>
      <span class="u-spacer"></span>
      <div class="u-seg">
        <button type="button" :class="{ 'is-active': viewMode === 'grouped' }" @click="setViewMode('grouped')">
          分组视图
        </button>
        <button type="button" :class="{ 'is-active': viewMode === 'table' }" @click="setViewMode('table')">
          表格视图
        </button>
      </div>
      <button class="u-btn" type="button" :disabled="!userRows.length" @click="exportCurrentRows">
        <SvgIcon icon="mdi:download" class="u-btn-ico" />导出
      </button>
    </section>

    <!-- ===================== 表格视图 ===================== -->
    <div v-if="viewMode === 'table'" class="u-split">
      <aside class="u-card u-facet">
        <div class="u-card-head">
          <h3>线路机</h3>
          <span class="u-badge is-neutral">{{ groupTotal }}</span>
          <div class="u-right"><span class="u-dim">按用户数</span></div>
        </div>
        <div class="u-facet-body">
          <button
            type="button"
            class="u-facet-item"
            :class="{ 'is-active': filters.nodeID === '' }"
            @click="setNodeFilter('')"
          >
            <i class="u-dot is-brand"></i>
            <span class="u-facet-name">全部线路机</span>
            <span class="u-num u-dim">{{ totalUsers === null ? '--' : totalUsers }}</span>
          </button>
          <button
            v-for="group in groups"
            :key="group.nodeId"
            type="button"
            class="u-facet-item"
            :class="{ 'is-active': filters.nodeID === group.nodeId }"
            @click="setNodeFilter(group.nodeId)"
          >
            <i class="u-dot" :class="nodeDotClass(group.status)"></i>
            <span class="u-facet-name">{{ group.nodeName }}</span>
            <span class="u-num u-dim">{{ group.stats.total }}</span>
          </button>

          <div class="u-hr"></div>
          <div class="u-facet-title">到期分布</div>
          <div v-for="item in distribution" :key="item.key" class="u-dist">
            <div class="u-dist-top">
              <span class="u-dim">{{ item.label }}</span>
              <span class="u-num">{{ item.value }}</span>
            </div>
            <div class="u-bar"><span :class="item.tone" :style="{ width: item.percent + '%' }"></span></div>
          </div>

          <div class="u-hr"></div>
          <button type="button" class="u-btn u-btn-block" @click="importModalVisible = true">
            <SvgIcon icon="mdi:upload" class="u-btn-ico" />导入历史账单
          </button>
          <p class="u-hint">CSV 逐行预览校验，有错禁止提交；verified=false 的记录进入待核实，不计入财务。</p>

          <div class="u-hr"></div>
          <p class="u-hint">落地机不产生业务用户：落地机 Inbound 恒为基础设施，只作为路径出口。</p>
        </div>
      </aside>

      <section class="u-card u-main">
        <div class="u-card-head">
          <h3>用户</h3>
          <span class="u-badge is-neutral">{{ userTotal }} 户</span>
          <div class="u-right"><span class="u-dim">按到期时间升序 · 不限期排最后</span></div>
        </div>

        <NAlert v-if="errorMessage" type="warning" :show-icon="false" class="u-alert">{{ errorMessage }}</NAlert>

        <div v-if="userError" class="u-inline-error">
          {{ userError }}
          <button class="u-btn u-btn-xs" type="button" @click="loadUsers">重试</button>
        </div>

        <div class="u-table-wrap">
          <table class="u-table">
            <thead>
              <tr>
                <th>用户 / Inbound</th>
                <th>状态</th>
                <th>到期时间</th>
                <th>用量 / 配额 <span class="u-star">*</span></th>
                <th>线路机</th>
                <th>出口 IP</th>
                <th class="is-num">计费 <span class="u-star">*</span></th>
                <th></th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="row in userRows" :key="row.id" class="u-row" @click="openDetail(row.id)">
                <td class="u-td-strong">
                  <span class="u-name">{{ row.name || '未命名用户' }}</span>
                  <span class="u-sub">{{ row.inboundTag || '--' }} · {{ row.clientCount }} 台设备</span>
                </td>
                <td>
                  <span class="u-badge" :class="statusBadgeClass(row.status)">
                    <i class="u-badge-dot"></i>{{ statusLabel(row.status) }}
                  </span>
                  <span v-if="row.status === 'expiring'" class="u-sub u-warn">还剩 {{ daysUntil(row.expiresAt) }} 天</span>
                </td>
                <td class="u-num">
                  {{ formatExpiry(row.expiresAt, row.status) }}
                  <span class="u-sub">{{ expiryHint(row) }}</span>
                </td>
                <td>
                  <TrafficValue :value="row.trafficBytes" />
                  <span class="u-dim"> / {{ row.quotaBytes ? formatQuota(row.quotaBytes) : '--' }}</span>
                  <span v-if="row.quotaBytes" class="u-sub">{{ quotaUsageHint(row) }}</span>
                </td>
                <td><span class="u-badge is-neutral">{{ row.nodeName || '未关联' }}</span></td>
                <td>
                  <span class="u-mono">{{ row.exitIpAddress || '未配置' }}</span>
                  <span class="u-sub">{{ exitIpHint(row) }}</span>
                </td>
                <td class="is-num">
                  <span v-if="row.billingType === 'free'" class="u-badge is-neutral">免费</span>
                  <span v-else class="u-strong">收费</span>
                  <span class="u-sub">{{ billingHint(row) }}</span>
                </td>
                <td class="u-td-actions" @click.stop>
                  <button
                    v-if="row.status === 'expiring' || row.status === 'expired'"
                    class="u-btn u-btn-xs is-soft"
                    type="button"
                    @click="openOrderForUser(row.id)"
                  >
                    续费
                  </button>
                  <button class="u-icon-btn" type="button" title="查看详情" @click="openDetail(row.id)">
                    <SvgIcon icon="mdi:chevron-right" />
                  </button>
                </td>
              </tr>
            </tbody>
          </table>

          <div v-if="userLoading" class="u-loading"><NSpin size="small" /></div>
          <div v-if="!userLoading && !userRows.length && !userError" class="u-empty">
            {{ hasActiveFilters ? '没有符合筛选条件的用户' : '暂无用户同步数据' }}
          </div>
        </div>

        <div class="u-table-foot">
          <span class="u-dim">共 {{ userTotal }} 户 · 每页 {{ filters.page_size }} 条 · 第 {{ filters.page }} 页</span>
          <NPagination
            :page="filters.page"
            :page-count="userPageCount"
            :page-slot="5"
            size="small"
            @update:page="changeUserPage"
          />
        </div>
      </section>
    </div>

    <!-- ===================== 分组视图 ===================== -->
    <div v-else class="u-grouped">
      <div class="u-grouped-tools">
        <span class="u-dim">按线路机分组 · 展开后懒加载该组用户</span>
        <span class="u-spacer"></span>
        <button class="u-btn u-btn-xs" type="button" :disabled="!groups.length" @click="expandAllGroups">全部展开</button>
        <button class="u-btn u-btn-xs" type="button" :disabled="!groups.length" @click="collapseAllGroups">全部收起</button>
      </div>

      <section
        v-for="group in groups"
        :key="group.nodeId"
        class="u-card u-group"
        :class="{ 'is-open': isGroupExpanded(group.nodeId) }"
      >
        <button type="button" class="u-group-head" @click="toggleGroup(group)">
          <SvgIcon icon="mdi:chevron-right" class="u-group-chev" />
          <span class="u-group-name">{{ group.nodeName }}</span>
          <span class="u-badge" :class="nodeBadgeClass(group.status)">
            <i class="u-badge-dot"></i>{{ nodeStatusLabel(group.status) }}
          </span>
          <span class="u-badge is-brand">线路机</span>
          <span class="u-badge is-neutral">{{ group.stats.total }} 户</span>
          <span class="u-group-summary u-dim">{{ groupSummary(group) }}</span>
          <span class="u-group-traffic u-dim">累计 <TrafficValue :value="group.stats.trafficBytes" /></span>
          <span class="u-group-sync u-dim">
            {{ group.lastSyncAt ? '最近同步 ' + formatDate(group.lastSyncAt) : '尚未完成成功同步' }}
          </span>
          <span class="u-group-toggle">{{ isGroupExpanded(group.nodeId) ? '收起用户' : '查看用户' }}</span>
        </button>

        <div v-if="isGroupExpanded(group.nodeId)" class="u-group-body">
          <NSpin :show="groupUsers[group.nodeId]?.loading || false">
            <div v-if="groupUsers[group.nodeId]?.error" class="u-inline-error">
              {{ groupUsers[group.nodeId]?.error }}
              <button class="u-btn u-btn-xs" type="button" @click="loadGroupUsers(group.nodeId)">重试</button>
            </div>
            <div v-else class="u-table-wrap">
              <table class="u-table is-compact">
                <thead>
                  <tr>
                    <th>用户 / Inbound</th>
                    <th>状态</th>
                    <th>到期时间</th>
                    <th>用量 / 配额 <span class="u-star">*</span></th>
                    <th>出口 IP</th>
                    <th class="is-num">计费 <span class="u-star">*</span></th>
                    <th></th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="row in groupRows(group.nodeId)" :key="row.id" class="u-row" @click="openDetail(row.id)">
                    <td class="u-td-strong">
                      <span class="u-name">{{ row.name || '未命名用户' }}</span>
                      <span class="u-sub">{{ row.inboundTag || '--' }}</span>
                    </td>
                    <td>
                      <span class="u-badge" :class="statusBadgeClass(row.status)">
                        <i class="u-badge-dot"></i>{{ statusLabel(row.status) }}
                      </span>
                    </td>
                    <td class="u-num">
                      {{ formatExpiry(row.expiresAt, row.status) }}
                      <span class="u-sub">{{ expiryHint(row) }}</span>
                    </td>
                    <td>
                      <TrafficValue :value="row.trafficBytes" />
                      <span class="u-dim"> / --</span>
                    </td>
                    <td>
                      <span class="u-mono">{{ row.exitIpAddress || '未配置' }}</span>
                    </td>
                    <td class="is-num">
                      <span v-if="row.billingType === 'free'" class="u-badge is-neutral">免费</span>
                      <span v-else class="u-strong">收费</span>
                    </td>
                    <td class="u-td-actions" @click.stop>
                      <button class="u-icon-btn" type="button" title="查看详情" @click="openDetail(row.id)">
                        <SvgIcon icon="mdi:chevron-right" />
                      </button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </NSpin>
          <div class="u-table-foot">
            <span class="u-dim">共 {{ groupUserTotal(group.nodeId) }} 户 · 已显示 {{ groupRows(group.nodeId).length }} 户</span>
            <NPagination
              :page="groupPage(group.nodeId)"
              :page-count="groupPageCount(group.nodeId)"
              :page-slot="5"
              size="small"
              @update:page="page => changeGroupPage(group.nodeId, page)"
            />
          </div>
        </div>
      </section>

      <div v-if="!groups.length" class="u-card u-empty-card">
        {{ hasActiveFilters ? '没有符合筛选条件的线路机分组' : '暂无用户同步数据' }}
      </div>
    </div>

    <p class="u-footnote">
      * 「用量 / 配额」的用量为相对流量基线的增量，配额来自 Agent 上报的 X-Panel Client
      <span class="u-mono">totalGB</span>（未上报时显示 <span class="u-mono">--</span>）；「计费」列展示中央维护的月费与周期（<span class="u-mono">billingAmount</span> / <span class="u-mono">billingCycle</span>）。
    </p>

    <!-- ===================== 详情抽屉 ===================== -->
    <NDrawer v-model:show="drawerVisible" :width="760" placement="right" :mask-closable="false">
      <NDrawerContent body-content-style="padding: 0">
        <template #header>
          <div class="u-dhead">
            <span class="u-avatar">{{ detailInitial }}</span>
            <div class="u-dhead-main">
              <div class="u-dhead-title">
                <h2>{{ detail?.displayName || selectedUserID || '用户详情' }}</h2>
                <span v-if="detail" class="u-badge" :class="statusBadgeClass(detail.status)">
                  <i class="u-badge-dot"></i>{{ statusLabel(detail.status) }}
                </span>
                <span v-if="pendingCandidates.length" class="u-badge is-danger">续费待确认</span>
                <span v-if="detail" class="u-badge is-neutral">
                  {{ cycleLabel(detail.billingCycle) }} · {{ formatMoney(detail.billingAmount) }}
                </span>
                <span v-if="detail?.path" class="u-badge is-brand">{{ pathModeLabel(detail.path.mode) }}</span>
              </div>
              <p class="u-dhead-sub">
                {{
                  detail
                    ? `${detail.node.name || '未关联线路机'} · ${detail.inbound.tag || detail.inbound.remoteId || '--'} · 一个 Reality Inbound 对应一个业务用户`
                    : '正在读取用户详情…'
                }}
              </p>
            </div>
            <div class="u-dhead-actions">
              <NDropdown trigger="click" :options="moreMenuOptions" @select="handleMoreMenu">
                <button class="u-icon-btn" type="button" title="更多操作">
                  <SvgIcon icon="mdi:dots-horizontal" />
                </button>
              </NDropdown>
              <button class="u-btn is-primary" type="button" :disabled="!detail" @click="openOrderModal">续费</button>
              <button class="u-icon-btn" type="button" title="关闭" @click="closeDetail"><SvgIcon icon="mdi:close" /></button>
            </div>
          </div>
        </template>

        <div v-if="detailLoading" class="u-drawer-loading"><NSpin size="medium" /></div>

        <template v-else>
          <NAlert v-if="detailError" type="error" :show-icon="true" class="u-drawer-alert">
            {{ detailError }}
            <button class="u-btn u-btn-xs" type="button" @click="openDetail(selectedUserID)">重试</button>
          </NAlert>

          <template v-if="detail">
            <div class="u-dtabs">
              <button
                v-for="tab in detailTabs"
                :key="tab.key"
                type="button"
                :class="{ 'is-active': detailTab === tab.key }"
                @click="scrollToSection(tab.key)"
              >
                {{ tab.label }}
              </button>
            </div>

            <div class="u-dbody">
              <div class="u-stack">
                <!-- 1. 待确认续费 -->
                <div v-if="pendingCandidates.length" class="u-callout is-warn">
                  <SvgIcon icon="mdi:alert-outline" class="u-callout-ico" />
                  <div class="u-callout-main">
                    <div class="u-callout-head">
                      <b>待确认续费</b>
                      <span class="u-badge is-warn">Agent 检测到 Client 到期延长</span>
                      <span class="u-dim u-ml-auto">{{ relativeTime(pendingCandidates[0].detectedAt) }}</span>
                    </div>
                    <div v-for="candidate in pendingCandidates" :key="candidate.id" class="u-cand">
                      <div class="u-cand-info">
                        <b class="u-num">{{ formatDate(candidate.oldExpiryAt) }} → {{ formatDate(candidate.newExpiryAt) }}</b>
                        <span class="u-cand-sub">
                          {{ cycleLabel(candidate.billingCycle) }} · 建议金额 {{ formatMoney(candidate.suggestedAmount) }} ·
                          来源 {{ detail.node.name || '--' }} 的 Client 到期变更
                        </span>
                      </div>
                      <div class="u-cand-actions">
                        <button class="u-btn u-btn-xs is-success" type="button" @click="openRenewConfirm(candidate)">
                          确认续费
                        </button>
                        <button class="u-btn u-btn-xs" type="button" @click="rejectRenewal(candidate)">非收费变更</button>
                      </div>
                    </div>
                    <p class="u-callout-note">
                      确认后会立即生成 1 条 agent_detected 收费记录并计入实收，金额与服务期不可再修改；
                      选择「非收费变更」只关闭候选与事件，不计入财务。代理只提议、不自动记账。
                    </p>
                  </div>
                </div>

                <!-- 2. 业务资料 -->
                <section id="u-sec-overview" class="u-dcard">
                  <div class="u-dcard-head">
                    <h3>业务资料</h3>
                    <span class="u-badge is-brand">中央面板可写</span>
                    <div class="u-right">
                      <button class="u-btn u-btn-xs" type="button" @click="businessEditing = !businessEditing">
                        {{ businessEditing ? '取消编辑' : '编辑' }}
                      </button>
                    </div>
                  </div>
                  <div class="u-dcard-body">
                    <div class="u-grid-2">
                      <label class="u-field">
                        <span class="u-label">业务用户名称 <i class="u-req">*</i></span>
                        <input v-model="editForm.displayName" class="u-input is-plain" maxlength="120" :disabled="!businessEditing" />
                      </label>
                      <div class="u-field">
                        <span class="u-label">收费周期</span>
                        <NSelect v-model:value="editForm.billingCycle" :options="billingCycleOptions" size="small" :disabled="!businessEditing" />
                      </div>
                      <div class="u-field">
                        <span class="u-label">计费类型</span>
                        <NSelect v-model:value="editForm.billingType" :options="billingTypeFormOptions" size="small" :disabled="!businessEditing" />
                      </div>
                      <label class="u-field">
                        <span class="u-label">{{ editForm.billingCycle === 'annual' ? '年费（CNY）' : '月费（CNY）' }}</span>
                        <NInputNumber
                          v-model:value="editForm.billingAmount"
                          :min="0"
                          :max="100000000"
                          :precision="2"
                          size="small"
                          :disabled="!businessEditing"
                        />
                      </label>
                    </div>
                    <label v-if="editForm.billingType === 'free'" class="u-field u-mt-12">
                      <span class="u-label">免费原因</span>
                      <input
                        v-model="editForm.freeReason"
                        class="u-input is-plain"
                        maxlength="200"
                        placeholder="例如：朋友体验、内部测试"
                        :disabled="!businessEditing"
                      />
                    </label>
                    <label class="u-field u-mt-12">
                      <span class="u-label">备注</span>
                      <textarea v-model="editForm.notes" class="u-input is-plain is-area" rows="3" maxlength="2000" :disabled="!businessEditing"></textarea>
                    </label>
                    <div v-if="businessEditing" class="u-form-actions">
                      <button class="u-btn is-primary" type="button" :disabled="saving" @click="saveBusinessFields">保存业务信息</button>
                    </div>
                    <p v-else class="u-hint">点「编辑」后可修改；保存只更新中央面板业务信息，不改 X-Panel。</p>
                  </div>
                </section>

                <!-- 3. 服务与到期 -->
                <section class="u-dcard">
                  <div class="u-dcard-head">
                    <h3>服务与到期</h3>
                    <span class="u-badge is-neutral">由启用 Client 派生</span>
                    <div class="u-right u-dim">不可手改</div>
                  </div>
                  <div class="u-dcard-body">
                    <div class="u-grid-3">
                      <div class="u-field">
                        <span class="u-label">到期时间</span>
                        <b class="u-num">{{ formatExpiry(detail.expiresAt, detail.status) }}</b>
                      </div>
                      <div class="u-field">
                        <span class="u-label">剩余</span>
                        <span class="u-badge" :class="expiryBadgeClass(detail)">{{ expiryText(detail) }}</span>
                      </div>
                      <div class="u-field">
                        <span class="u-label">用户状态</span>
                        <span class="u-badge" :class="statusBadgeClass(detail.status)">{{ statusLabel(detail.status) }}</span>
                      </div>
                      <div class="u-field">
                        <span class="u-label">启用 Client</span>
                        <b class="u-num">{{ enabledClientCount }} 个</b>
                      </div>
                      <div class="u-field">
                        <span class="u-label">X-Panel 累计流量</span>
                        <b><TrafficValue :value="detail.inbound.allTime" /></b>
                      </div>
                      <div class="u-field">
                        <span class="u-label">最近在线</span>
                        <b class="u-num">{{ formatDate(lastOnlineAt) }}</b>
                      </div>
                    </div>
                    <p class="u-hint u-mt-12">
                      到期时间取启用 Client 中最早有限的到期值；不限期用户不会生成待确认续费。无宽限期，到期后连接能力由本机配置决定。
                    </p>
                  </div>
                </section>

                <!-- 4. 收费记录 -->
                <section id="u-sec-billing" class="u-dcard">
                  <div class="u-dcard-head">
                    <h3>收费记录</h3>
                    <span class="u-badge is-neutral">{{ detail.billingRecords.length }} 笔</span>
                    <div class="u-right">
                      <button class="u-btn u-btn-xs is-soft" type="button" @click="openOrderModal">新增订单</button>
                    </div>
                  </div>
                  <div v-if="detail.billingRecords.length" class="u-list">
                    <div v-for="record in detail.billingRecords" :key="record.id" class="u-list-row">
                      <div class="u-list-main">
                        <div class="u-list-title">
                          <span class="u-strong">
                            {{ orderTypeLabel(record.orderType) }} · {{ cycleLabel(record.billingCycle) }} ·
                            {{ formatMoney(record.amount, record.currency) }}
                          </span>
                          <span v-if="record.verificationStatus === 'unverified'" class="u-badge is-warn">待核实</span>
                          <span v-else-if="record.origin === 'historical_import'" class="u-badge is-info">历史导入</span>
                          <span v-if="record.status === 'cancelled'" class="u-badge is-danger">已取消</span>
                        </div>
                        <div class="u-list-sub">
                          服务 {{ formatDate(record.serviceFrom) }} → {{ formatDate(record.serviceTo) }} · 收款
                          {{ formatDate(record.paidAt) }}
                        </div>
                      </div>
                      <div v-if="record.source === 'manual' && record.status !== 'cancelled'" class="u-list-actions">
                        <button
                          v-if="record.verificationStatus === 'unverified'"
                          class="u-btn u-btn-xs is-soft"
                          type="button"
                          @click="openVerifyModal(record)"
                        >
                          核验
                        </button>
                        <button class="u-btn u-btn-xs is-danger" type="button" @click="cancelRecord(record)">取消</button>
                      </div>
                    </div>
                  </div>
                  <p v-else class="u-hint u-px">暂无收费记录；首次收费请使用下方「首笔订单」。</p>
                </section>

                <!-- 5. 首笔订单（无任何收费记录时） -->
                <section v-if="detail.billingRecords.length === 0" class="u-dcard">
                  <div class="u-dcard-head">
                    <h3>首笔订单</h3>
                    <span class="u-badge is-neutral">对应 Agent 同步的 Inbound 业务用户</span>
                  </div>
                  <div class="u-dcard-body">
                    <div class="u-grid-2">
                      <div class="u-field">
                        <span class="u-label">订单金额（CNY）<i class="u-req">*</i></span>
                        <NInputNumber v-model:value="initialOrderForm.amount" :min="0" :max="100000000" :precision="2" size="small" />
                      </div>
                      <div class="u-field">
                        <span class="u-label">实际收款日期 <i class="u-req">*</i></span>
                        <NDatePicker v-model:formatted-value="initialOrderForm.paidAt" type="date" value-format="yyyy-MM-dd" size="small" />
                      </div>
                      <div class="u-field">
                        <span class="u-label">服务开始日期 <i class="u-req">*</i></span>
                        <NDatePicker
                          v-model:formatted-value="initialOrderForm.serviceFrom"
                          type="date"
                          value-format="yyyy-MM-dd"
                          size="small"
                          placeholder="请选择服务开始日期"
                        />
                      </div>
                      <div class="u-field">
                        <span class="u-label">服务结束日期（Agent 同步）</span>
                        <NDatePicker
                          v-model:formatted-value="initialOrderForm.serviceTo"
                          type="date"
                          value-format="yyyy-MM-dd"
                          size="small"
                          disabled
                          placeholder="Agent 尚未同步到期日"
                        />
                      </div>
                    </div>
                    <label class="u-field u-mt-12">
                      <span class="u-label">订单备注</span>
                      <input v-model="initialOrderForm.notes" class="u-input is-plain" maxlength="2000" placeholder="可选" />
                    </label>
                    <div class="u-form-actions">
                      <button class="u-btn is-primary" type="button" :disabled="initialOrderSaving" @click="createInitialOrder">
                        确认首笔订单
                      </button>
                    </div>
                    <p class="u-hint">
                      首笔订单对应 Agent 同步过来的线路机 Inbound 业务用户；服务开始日期请按真实业务填写，服务结束日期采用 Agent 同步的 Client 到期日。
                    </p>
                  </div>
                </section>

                <!-- 6. 用户路径 -->
                <section id="u-sec-path" class="u-dcard">
                  <div class="u-dcard-head">
                    <h3>用户路径</h3>
                    <span v-if="detail.path" class="u-badge" :class="detail.path.valid ? 'is-ok' : 'is-warn'">
                      {{ detail.path.valid ? '有效' : '需检查' }}
                    </span>
                    <div class="u-right">
                      <button class="u-btn u-btn-xs" type="button" :disabled="pathAssetsLoading" @click="loadPathAssets()">
                        重新读取资源
                      </button>
                    </div>
                  </div>
                  <div class="u-dcard-body">
                    <NAlert v-if="pathAssetsError" type="warning" :show-icon="true" class="u-alert">
                      {{ pathAssetsError }}
                    </NAlert>
                    <div class="u-grid-2">
                      <div class="u-field">
                        <span class="u-label">线路机（自动）</span>
                        <b>{{ detail.node.name || '未关联线路机' }}</b>
                      </div>
                      <div class="u-field">
                        <span class="u-label">出口路径</span>
                        <NSelect
                          :value="pathMode"
                          :options="pathModeOptions"
                          size="small"
                          :disabled="pathAssetsLoading"
                          @update:value="handlePathModeChange"
                        />
                      </div>
                      <div v-if="pathMode === 'landing'" class="u-field">
                        <span class="u-label">落地机（执行节点）<i class="u-req">*</i></span>
                        <NSelect
                          :value="pathLandingNodeID"
                          :options="pathLandingNodeOptions"
                          size="small"
                          filterable
                          :loading="pathAssetsLoading"
                          :disabled="pathAssetsLoading"
                          placeholder="选择启用的落地机"
                          @update:value="handlePathLandingNodeChange"
                        />
                      </div>
                      <div v-if="pathMode === 'landing'" class="u-field">
                        <span class="u-label">落地 Inbound（基础设施入口）<i class="u-req">*</i></span>
                        <NSelect
                          v-model:value="pathLandingInboundID"
                          :options="pathInboundOptions"
                          size="small"
                          filterable
                          :loading="pathLandingInboundsLoading"
                          :disabled="pathAssetsLoading || !pathLandingNodeID"
                          placeholder="选择该落地机已同步的 Inbound"
                        />
                      </div>
                      <div class="u-field u-span-2">
                        <span class="u-label">固定出口 IP</span>
                        <NSelect
                          v-model:value="pathExitIPIDs"
                          :options="pathExitIpOptions"
                          size="small"
                          multiple
                          filterable
                          clearable
                          max-tag-count="responsive"
                          :loading="pathAssetsLoading"
                          :disabled="pathAssetsLoading || !detail.node.id"
                          placeholder="选择当前路径允许的出口 IP"
                        />
                      </div>
                    </div>
                    <p v-if="pathMode === 'landing' && pathAssets && !pathLandingNodeOptions.length" class="u-hint u-warn">
                      暂无启用的落地机，请先在节点管理中添加并启用。
                    </p>
                    <p v-if="pathMode === 'landing' && pathLandingNodeID && !pathInboundOptions.length" class="u-hint u-warn">
                      {{
                        pathAssets?.landingNodes.find(node => node.id === pathLandingNodeID)?.inboundState === 'pending'
                          ? '该落地机尚未完成首次 Agent 同步。'
                          : '该落地机暂无启用的 Inbound。'
                      }}
                    </p>
                    <label class="u-field u-mt-12">
                      <span class="u-label">路径备注</span>
                      <input v-model="pathNotes" class="u-input is-plain" maxlength="500" placeholder="可选" />
                    </label>
                    <p v-if="pathError" class="u-hint u-warn">{{ pathError }}</p>
                    <div class="u-form-actions">
                      <button
                        class="u-btn is-primary"
                        type="button"
                        :disabled="
                          pathSaving ||
                          pathAssetsLoading ||
                          (pathMode === 'landing' && (!pathLandingNodeID || !pathLandingInboundID)) ||
                          !pathExitIPIDs.length
                        "
                        @click="savePathAssignment"
                      >
                        保存用户路径
                      </button>
                      <button v-if="detail.path" class="u-btn" type="button" :disabled="pathSaving" @click="confirmRemovePathAssignment">
                        解除路径
                      </button>
                    </div>

                    <div class="u-hr"></div>
                    <div class="u-sub-title">当前生效路径</div>
                    <div v-if="!detail.path" class="u-empty-small">尚未配置用户路径</div>
                    <div v-else class="u-path-card">
                      <div class="u-path-head">
                        <b>{{ detail.path.relayNodeName || detail.node.name || '--' }} → {{ detail.path.landingNodeName || '线路机直出' }}</b>
                        <span class="u-badge" :class="detail.path.valid ? 'is-ok' : 'is-warn'">
                          {{ detail.path.valid ? '有效' : '需检查' }}
                        </span>
                      </div>
                      <div class="u-path-meta">模式：{{ pathModeLabel(detail.path.mode) }}</div>
                      <div class="u-path-meta">
                        出口：{{ pathExitIPText(detail.path) }} · 归属：{{ detail.path.exitIpOwnerNodeName || '独立 S5' }}
                      </div>
                      <div v-if="detail.path.landingInboundTag || detail.path.landingInboundId" class="u-path-meta">
                        落地 Inbound：{{ detail.path.landingInboundTag || detail.path.landingInboundId }}
                      </div>
                      <div v-if="detail.path.notes" class="u-path-meta">备注：{{ detail.path.notes }}</div>
                    </div>

                    <template v-if="detail.pathHistory.length">
                      <div class="u-hr"></div>
                      <div class="u-sub-title">路径历史</div>
                      <div v-for="history in detail.pathHistory" :key="history.id" class="u-history-row">
                        <span>{{ history.relayNodeName || '--' }} → {{ history.landingNodeName || '线路机直出' }} · {{ pathExitIPText(history) }}</span>
                        <span class="u-dim u-num">
                          {{ formatDate(history.activeFrom) }}{{ history.activeTo ? ` 至 ${formatDate(history.activeTo)}` : ' 至今' }}
                        </span>
                      </div>
                    </template>
                  </div>
                </section>

                <!-- 7. 设备与流量 -->
                <section id="u-sec-device" class="u-dcard">
                  <div class="u-dcard-head">
                    <h3>设备 Client / Email</h3>
                    <span class="u-badge is-neutral">只读</span>
                    <div class="u-right u-dim">Email 仅作设备凭证，不参与用户身份</div>
                  </div>
                  <DeviceTable :devices="detail.clients" />
                </section>

                <section class="u-dcard">
                  <div class="u-dcard-head">
                    <h3>流量趋势</h3>
                    <span class="u-badge is-neutral">只读</span>
                  </div>
                  <div class="u-dcard-body is-flush">
                    <TrafficTrendChart :user-id="detail.id" />
                  </div>
                </section>

                <!-- 8. X-Panel 快照 -->
                <section class="u-dcard is-muted">
                  <div class="u-dcard-head">
                    <h3>X-Panel 快照</h3>
                    <span class="u-badge is-neutral">只读 · 同步而来</span>
                    <div class="u-right u-dim">上次同步 {{ formatDate(detail.inbound.lastSeenAt) }}</div>
                  </div>
                  <div class="u-kv">
                    <div class="u-kv-row">
                      <span>线路机</span>
                      <b>{{ detail.node.name || '--' }}{{ detail.node.type ? `（${detail.node.type === 'relay' ? '线路机' : '落地机'}）` : '' }}</b>
                    </div>
                    <div class="u-kv-row">
                      <span>Inbound</span>
                      <b class="u-mono">{{ detail.inbound.tag || detail.inbound.remoteId || '--' }}</b>
                    </div>
                    <div class="u-kv-row">
                      <span>协议 / 端口</span>
                      <b class="u-mono">{{ detail.inbound.protocol || '--' }} / {{ detail.inbound.port || '--' }}</b>
                    </div>
                    <div class="u-kv-row">
                      <span>启用状态</span>
                      <span class="u-badge" :class="detail.inbound.enabled ? 'is-ok' : 'is-neutral'">
                        {{ detail.inbound.enabled ? '启用' : '停用' }}
                      </span>
                    </div>
                    <div class="u-kv-row">
                      <span>X-Panel 累计流量</span>
                      <TrafficValue :value="detail.inbound.allTime" />
                    </div>
                  </div>
                </section>
              </div>
            </div>
          </template>
        </template>

        <template #footer>
          <div class="u-drawer-foot">
            <span class="u-dim">只读字段来自 X-Panel 同步；保存只更新中央面板业务信息。</span>
            <span class="u-spacer"></span>
            <button class="u-btn" type="button" @click="closeDetail">关闭</button>
          </div>
        </template>
      </NDrawerContent>
    </NDrawer>

    <!-- ===================== 弹窗：确认续费 ===================== -->
    <NModal v-model:show="renewModalVisible" preset="card" title="确认续费" class="u-modal" :bordered="false" :style="{ width: '520px' }">
      <div v-if="renewCandidate" class="u-modal-body">
        <div class="u-callout is-danger u-callout-tight">
          <SvgIcon icon="mdi:alert-outline" class="u-callout-ico" />
          <div class="u-callout-main">
            <div class="u-callout-head"><b>计入后不可修改</b></div>
            <p class="u-callout-note">
              确认后立即生成 1 条 agent_detected 收费记录，计入实收、服务期收入与续费统计；金额与服务期此后不可再改。
            </p>
          </div>
        </div>
        <div class="u-grid-2 u-mt-12">
          <div class="u-field"><span class="u-label">用户</span><b>{{ detail?.displayName || selectedUserID }}</b></div>
          <div class="u-field"><span class="u-label">来源</span><b>{{ detail?.node.name || '--' }} Client 到期延长</b></div>
          <div class="u-field"><span class="u-label">服务开始</span><b class="u-num">{{ formatDate(renewCandidate.oldExpiryAt) }}</b></div>
          <div class="u-field"><span class="u-label">服务结束</span><b class="u-num">{{ formatDate(renewCandidate.newExpiryAt) }}</b></div>
          <div class="u-field">
            <span class="u-label">金额（CNY）<i class="u-req">*</i></span>
            <NInputNumber v-model:value="renewForm.amount" :min="0" :max="100000000" :precision="2" size="small" />
          </div>
          <div class="u-field">
            <span class="u-label">实际收款日期 <i class="u-req">*</i></span>
            <NDatePicker v-model:formatted-value="renewForm.paidAt" type="date" value-format="yyyy-MM-dd" size="small" />
          </div>
        </div>
        <p class="u-hint">服务期按半开区间 [开始, 结束) 校验：与已有记录重叠会被拒绝，首尾相接允许；收款日期按本地时区计。</p>
      </div>
      <template #footer>
        <div class="u-modal-foot">
          <span class="u-dim u-mr-auto">确认后若发现错误，只能取消该记录（现金回退）。</span>
          <button class="u-btn" type="button" @click="renewModalVisible = false">取消</button>
          <button class="u-btn is-primary" type="button" :disabled="renewSaving" @click="submitRenewConfirm">
            确认并计入财务
          </button>
        </div>
      </template>
    </NModal>

    <!-- ===================== 弹窗：新增订单 ===================== -->
    <NModal v-model:show="orderModalVisible" preset="card" title="新增订单" class="u-modal" :bordered="false" :style="{ width: '680px' }">
      <div class="u-modal-body">
        <div class="u-grid-2">
          <div class="u-field">
            <span class="u-label">订单类型 <i class="u-req">*</i></span>
            <NSelect
              v-model:value="orderForm.orderType"
              :options="[
                { label: '首购', value: 'initial' },
                { label: '续费', value: 'renewal' },
                { label: '逾期恢复', value: 'recovery' }
              ]"
              size="small"
            />
          </div>
          <div class="u-field">
            <span class="u-label">收费周期 <i class="u-req">*</i></span>
            <NSelect v-model:value="orderForm.billingCycle" :options="billingCycleOptions" size="small" />
          </div>
          <div class="u-field">
            <span class="u-label">订单金额（CNY）<i class="u-req">*</i></span>
            <NInputNumber v-model:value="orderForm.amount" :min="0" :max="100000000" :precision="2" size="small" />
          </div>
          <div class="u-field">
            <span class="u-label">实际收款日期 <i class="u-req">*</i></span>
            <NDatePicker v-model:formatted-value="orderForm.paidAt" type="date" value-format="yyyy-MM-dd" size="small" />
          </div>
          <div class="u-field">
            <span class="u-label">服务开始时间 <i class="u-req">*</i></span>
            <NDatePicker v-model:value="orderForm.serviceFrom" type="datetime" format="yyyy-MM-dd HH:mm:ss" size="small" />
          </div>
          <div class="u-field">
            <span class="u-label">服务结束时间 <i class="u-req">*</i></span>
            <NDatePicker v-model:value="orderForm.serviceTo" type="datetime" format="yyyy-MM-dd HH:mm:ss" size="small" />
          </div>
        </div>
        <label class="u-field u-mt-12">
          <span class="u-label">订单备注</span>
          <input v-model="orderForm.notes" class="u-input is-plain" maxlength="2000" />
        </label>
        <p class="u-hint">服务时间按本地时区显示。开始时间默认接续上一笔未取消订单的结束时间，结束时间默认使用已同步的用户到期时间。</p>
      </div>
      <template #footer>
        <div class="u-modal-foot">
          <button class="u-btn" type="button" @click="orderModalVisible = false">取消</button>
          <button class="u-btn is-primary" type="button" :disabled="orderSaving" @click="createAdditionalOrder">保存订单</button>
        </div>
      </template>
    </NModal>

    <!-- ===================== 弹窗：核验历史订单 ===================== -->
    <NModal v-model:show="verifyModalVisible" preset="card" title="核验历史订单" class="u-modal" :bordered="false" :style="{ width: '520px' }">
      <div class="u-modal-body">
        <div class="u-callout is-warn u-callout-tight">
          <SvgIcon icon="mdi:alert-outline" class="u-callout-ico" />
          <div class="u-callout-main">
            <p class="u-callout-note">核验后该订单会计入实收、服务期收入和续费统计。</p>
          </div>
        </div>
        <div class="u-field u-mt-12">
          <span class="u-label">实际收款日期 <i class="u-req">*</i></span>
          <NDatePicker v-model:formatted-value="verifyPaidAt" type="date" value-format="yyyy-MM-dd" size="small" />
        </div>
        <label class="u-field u-mt-12">
          <span class="u-label">核验备注</span>
          <textarea v-model="verifyNotes" class="u-input is-plain is-area" rows="3" maxlength="2000"></textarea>
        </label>
      </div>
      <template #footer>
        <div class="u-modal-foot">
          <button class="u-btn" type="button" @click="verifyModalVisible = false">取消</button>
          <button class="u-btn is-primary" type="button" :disabled="verifySaving" @click="submitVerifyRecord">确认核验</button>
        </div>
      </template>
    </NModal>

    <!-- ===================== 弹窗：导入历史账单 ===================== -->
    <NModal v-model:show="importModalVisible" preset="card" title="导入历史账单" class="u-modal" :bordered="false" :style="{ width: '900px', maxWidth: '92vw' }">
      <div class="u-modal-body">
        <div class="u-callout is-info u-callout-tight">
          <SvgIcon icon="mdi:information-outline" class="u-callout-ico" />
          <div class="u-callout-main">
            <p class="u-callout-note">
              每行一笔订单。表头支持：userId、nodeId、remoteInboundId、billingCycle、amount、serviceFrom、serviceTo、paidAt、orderType、verified、notes。
              userId 与 nodeId + remoteInboundId 二选一；日期支持 YYYY-MM-DD。verified 填 false 的记录会进入待核实，不计入财务。
            </p>
          </div>
        </div>
        <textarea
          v-model="importText"
          class="u-input is-plain is-area u-mt-12"
          rows="10"
          placeholder="userId,billingCycle,amount,serviceFrom,serviceTo,paidAt,orderType,verified,notes&#10;用户ID,monthly,100,2026-03-01,2026-04-01,2026-03-01,initial,true,首购"
        ></textarea>
        <div v-if="importPreview" class="u-mt-12">
          <NAlert :type="importPreview.errors.length ? 'error' : 'success'" :show-icon="true">
            {{
              importPreview.errors.length
                ? `预览发现 ${importPreview.errors.length} 个错误，请修正后重新预览`
                : `预览通过：${importPreview.records.length} 笔，其中 ${importPreview.unverified} 笔待核实`
            }}
          </NAlert>
          <div v-if="importPreview.errors.length" class="u-import-errors">
            <div v-for="item in importPreview.errors" :key="`${item.row}-${item.message}`">第 {{ item.row }} 行：{{ item.message }}</div>
          </div>
        </div>
      </div>
      <template #footer>
        <div class="u-modal-foot">
          <button class="u-btn" type="button" @click="importModalVisible = false">关闭</button>
          <button class="u-btn" type="button" :disabled="importLoading" @click="previewBillingImport">预览校验</button>
          <button
            class="u-btn is-primary"
            type="button"
            :disabled="importLoading || !importPreview || importPreview.errors.length > 0"
            @click="submitBillingImport"
          >
            确认导入
          </button>
        </div>
      </template>
    </NModal>
  </div>
</template>

<style scoped>
/* ============================================================
   用户管理 · UI v2（对齐 design-preview/v2/users.html 与 kit-v2）
   仅本页使用；颜色取自 ui-v2 / kit-v2 既有 token，不新增色值
   ============================================================ */

.u-page {
  --u-card: #141519;
  --u-card-hover: #191b20;
  --u-line: #23262c;
  --u-line-soft: #1c1f24;
  --u-text-1: #f2f3f5;
  --u-text-2: #bcc2cc;
  --u-text-3: #9aa1ad;
  --u-text-4: #868d99;
  --u-brand: #4c6fff;
  --u-brand-soft: rgba(76, 111, 255, 0.14);
  --u-ok: #2ed47a;
  --u-ok-soft: rgba(46, 212, 122, 0.14);
  --u-warn: #f5b93d;
  --u-warn-soft: rgba(245, 185, 61, 0.14);
  --u-danger: #f3656f;
  --u-danger-soft: rgba(243, 101, 111, 0.14);
  --u-info: #5aa9f0;
  --u-info-soft: rgba(90, 169, 240, 0.14);
  --u-violet: #a78bfa;
  --u-violet-soft: rgba(167, 139, 250, 0.14);
  --u-radius: 14px;

  display: flex;
  flex-direction: column;
  gap: 15px;
  color: var(--u-text-1);
  font-size: 13px;
}

/* ---------- 页头 ---------- */
.u-head {
  display: flex;
  align-items: flex-end;
  gap: 18px;
}

.u-head-main {
  min-width: 0;
}

.u-eyebrow {
  display: flex;
  align-items: center;
  gap: 7px;
  margin-bottom: 6px;
  font-size: 12px;
  color: var(--u-text-4);
}

.u-eyebrow-sep {
  color: #545b64;
}

.u-title {
  margin: 0;
  font-size: 22px;
  font-weight: 620;
  line-height: 1.25;
  color: #f6f7f9;
}

.u-desc {
  max-width: 820px;
  margin: 7px 0 0;
  font-size: 12px;
  line-height: 1.65;
  color: var(--u-text-3);
}

.u-head-actions {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-left: auto;
  flex: none;
  padding-bottom: 2px;
}

/* ---------- 通用按钮 / 输入 ---------- */
.u-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  height: 34px;
  padding: 0 13px;
  border: 1px solid #2a2e35;
  border-radius: 9px;
  background: #1b1e24;
  color: #e6e8ec;
  font-size: 12.5px;
  white-space: nowrap;
  cursor: pointer;
  transition: background 0.15s ease, border-color 0.15s ease;
}

.u-btn:hover {
  background: #232730;
  border-color: #363b44;
}

.u-btn:disabled {
  opacity: 0.45;
  cursor: not-allowed;
}

.u-btn-ico {
  width: 15px;
  height: 15px;
  flex: none;
}

.u-btn.is-primary {
  border-color: var(--u-brand);
  background: var(--u-brand);
  color: #fff;
}

.u-btn.is-primary:hover {
  background: #5d7cff;
  border-color: #5d7cff;
}

.u-btn.is-soft {
  border-color: rgba(76, 111, 255, 0.28);
  background: rgba(76, 111, 255, 0.13);
  color: #a8b8ff;
}

.u-btn.is-success {
  border-color: rgba(46, 212, 122, 0.32);
  background: rgba(46, 212, 122, 0.14);
  color: #5ce49b;
}

.u-btn.is-danger {
  border-color: rgba(243, 101, 111, 0.3);
  background: rgba(243, 101, 111, 0.13);
  color: #ff9ba2;
}

.u-btn-xs {
  height: 28px;
  padding: 0 10px;
  font-size: 12px;
  border-radius: 7px;
}

.u-btn-block {
  width: 100%;
}

.u-icon-btn {
  display: grid;
  place-items: center;
  width: 28px;
  height: 28px;
  border: 1px solid transparent;
  border-radius: 7px;
  background: transparent;
  color: #9aa1ab;
  cursor: pointer;
  transition: background 0.15s ease, color 0.15s ease;
}

.u-icon-btn:hover {
  background: #1a1d22;
  color: #e6e8ec;
}

.u-input {
  display: flex;
  align-items: center;
  gap: 8px;
  height: 34px;
  padding: 0 10px;
  border: 1px solid #262a31;
  border-radius: 9px;
  background: #0e1013;
}

.u-input:focus-within {
  border-color: rgba(76, 111, 255, 0.7);
  box-shadow: 0 0 0 3px rgba(76, 111, 255, 0.18);
}

.u-input input,
.u-input textarea {
  min-width: 0;
  flex: 1;
  border: 0;
  outline: none;
  background: transparent;
  color: var(--u-text-1);
  font-size: 12.5px;
  font-family: inherit;
}

.u-input input::placeholder,
.u-input textarea::placeholder {
  color: var(--u-text-4);
}

.u-input.is-plain {
  width: 100%;
}

.u-input.is-plain input:disabled,
.u-input.is-plain textarea:disabled {
  color: var(--u-text-2);
  -webkit-text-fill-color: var(--u-text-2);
  cursor: default;
}

.u-input.is-area {
  height: auto;
  padding: 8px 10px;
  align-items: flex-start;
}

.u-input.is-area textarea {
  resize: vertical;
  line-height: 1.65;
}

.u-input-ico {
  width: 15px;
  height: 15px;
  flex: none;
  color: var(--u-text-4);
}

.u-search {
  width: 250px;
  flex: none;
}

.u-select :deep(.n-base-selection) {
  border-radius: 9px;
}

/* ---------- 卡片 ---------- */
.u-card {
  position: relative;
  border: 1px solid var(--u-line);
  border-radius: var(--u-radius);
  background: var(--u-card);
  min-width: 0;
}

.u-card-head {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 14px 18px;
  border-bottom: 1px solid var(--u-line-soft);
}

.u-card-head h3 {
  margin: 0;
  font-size: 13.5px;
  font-weight: 600;
  color: var(--u-text-1);
}

.u-right {
  margin-left: auto;
  display: flex;
  align-items: center;
  gap: 8px;
}

/* ---------- KPI ---------- */
.u-kpis {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 12px;
}

.u-kpi {
  padding: 15px 16px 14px;
  border: 1px solid var(--u-line);
  border-radius: var(--u-radius);
  background: var(--u-card);
  transition: border-color 0.18s ease;
}

.u-kpi:hover {
  border-color: #31353d;
}

.u-kpi.is-clickable {
  cursor: pointer;
}

.u-kpi-top {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 10px;
}

.u-kpi-label {
  font-size: 12px;
  color: var(--u-text-3);
}

.u-kpi-icon {
  display: grid;
  place-items: center;
  width: 24px;
  height: 24px;
  margin-left: auto;
  flex: none;
  border-radius: 8px;
  border: 1px solid #2a2e35;
  background: #1d2026;
  color: #a8adb8;
}

.u-kpi-icon :deep(.svg-icon) {
  width: 13px;
  height: 13px;
}

.u-kpi-icon.is-blue {
  border-color: rgba(76, 111, 255, 0.36);
  background: rgba(76, 111, 255, 0.16);
  color: #8fa8ff;
}

.u-kpi-icon.is-green {
  border-color: rgba(46, 212, 122, 0.34);
  background: rgba(46, 212, 122, 0.14);
  color: #5ce49b;
}

.u-kpi-icon.is-amber {
  border-color: rgba(245, 185, 61, 0.34);
  background: rgba(245, 185, 61, 0.14);
  color: #f7ca6b;
}

.u-kpi-icon.is-red {
  border-color: rgba(243, 101, 111, 0.34);
  background: rgba(243, 101, 111, 0.14);
  color: #ff9ba2;
}

.u-kpi-value {
  display: flex;
  align-items: baseline;
  gap: 6px;
  font-size: 25px;
  font-weight: 620;
  line-height: 1.15;
  letter-spacing: -0.015em;
  font-variant-numeric: tabular-nums;
  color: #fff;
}

.u-kpi-unit {
  font-size: 12px;
  font-weight: 500;
  color: var(--u-text-3);
}

.u-kpi-foot {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 9px;
  font-size: 12px;
  color: var(--u-text-4);
}

.u-kpi-foot-text {
  color: var(--u-text-4);
}

/* ---------- 工具条 ---------- */
.u-toolbar {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
  padding: 12px 14px;
  border: 1px solid var(--u-line);
  border-radius: var(--u-radius);
  background: var(--u-card);
}

.u-spacer {
  flex: 1;
}

.u-seg {
  display: inline-flex;
  gap: 2px;
  padding: 3px;
  border: 1px solid #24272e;
  border-radius: 9px;
  background: #14161a;
}

.u-seg button {
  height: 26px;
  padding: 0 11px;
  border: 0;
  border-radius: 7px;
  background: transparent;
  color: var(--u-text-3);
  font-size: 12px;
  cursor: pointer;
  transition: background 0.14s ease, color 0.14s ease;
}

.u-seg button:hover {
  color: var(--u-text-1);
}

.u-seg button.is-active {
  background: #2b3038;
  color: #fff;
}

/* ---------- 主体：分面 + 宽表 ---------- */
.u-split {
  display: flex;
  align-items: flex-start;
  gap: 12px;
}

.u-facet {
  width: 216px;
  flex: none;
}

.u-facet-body {
  padding: 6px 8px 12px;
}

.u-facet-item {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
  height: 33px;
  padding: 0 9px;
  border: 0;
  border-radius: 8px;
  background: transparent;
  color: var(--u-text-2);
  font-size: 12.5px;
  text-align: left;
  cursor: pointer;
  transition: background 0.14s ease, color 0.14s ease;
}

.u-facet-item:hover {
  background: #1a1c21;
  color: var(--u-text-1);
}

.u-facet-item.is-active {
  background: var(--u-brand-soft);
  color: #fff;
  box-shadow: inset 2px 0 0 var(--u-brand);
}

.u-facet-name {
  min-width: 0;
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.u-dot {
  width: 7px;
  height: 7px;
  flex: none;
  border-radius: 50%;
  background: #545b64;
}

.u-dot.is-ok {
  background: var(--u-ok);
  box-shadow: 0 0 7px rgba(46, 212, 122, 0.5);
}

.u-dot.is-warn {
  background: var(--u-warn);
  box-shadow: 0 0 7px rgba(245, 185, 61, 0.45);
}

.u-dot.is-danger {
  background: var(--u-danger);
  box-shadow: 0 0 7px rgba(243, 101, 111, 0.45);
}

.u-dot.is-brand {
  background: var(--u-brand);
  box-shadow: 0 0 7px rgba(76, 111, 255, 0.5);
}

.u-facet-title {
  padding: 2px 9px 8px;
  font-size: 11px;
  font-weight: 600;
  letter-spacing: 0.12em;
  color: #7f8792;
}

.u-dist {
  padding: 0 9px 9px;
}

.u-dist-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 12px;
}

.u-bar {
  position: relative;
  height: 5px;
  margin-top: 5px;
  border-radius: 99px;
  background: #24272e;
  overflow: hidden;
}

.u-bar > span {
  position: absolute;
  inset: 0 auto 0 0;
  border-radius: 99px;
  background: var(--u-brand);
}

.u-bar > span.is-ok {
  background: var(--u-ok);
}

.u-bar > span.is-warn {
  background: var(--u-warn);
}

.u-bar > span.is-danger {
  background: var(--u-danger);
}

.u-bar > span.is-neutral {
  background: #545b64;
}

.u-hr {
  height: 1px;
  margin: 10px 8px;
  background: var(--u-line-soft);
}

.u-hint {
  margin: 8px 0 0;
  padding: 0 9px;
  font-size: 12px;
  line-height: 1.6;
  color: var(--u-text-4);
}

.u-warn {
  color: var(--u-warn);
}

/* ---------- 表格 ---------- */
.u-main {
  flex: 1;
  min-width: 0;
}

.u-table-wrap {
  overflow-x: auto;
}

.u-table {
  width: 100%;
  border-collapse: separate;
  border-spacing: 0;
  font-size: 13px;
}

.u-table thead th {
  position: sticky;
  top: 0;
  z-index: 2;
  padding: 9px 12px;
  border-bottom: 1px solid var(--u-line);
  background: var(--u-card);
  color: var(--u-text-3);
  font-size: 11.5px;
  font-weight: 600;
  letter-spacing: 0.05em;
  text-align: left;
  white-space: nowrap;
}

.u-table thead th.is-num {
  text-align: right;
}

.u-table tbody td {
  height: 46px;
  padding: 7px 12px;
  border-bottom: 1px solid var(--u-line-soft);
  color: var(--u-text-2);
  vertical-align: middle;
  white-space: nowrap;
}

.u-table tbody tr:last-child td {
  border-bottom: none;
}

.u-row {
  cursor: pointer;
  transition: background 0.12s ease;
}

.u-row:hover td {
  background: var(--u-card-hover);
}

.u-td-strong {
  color: var(--u-text-1);
}

.u-name {
  display: block;
  font-weight: 500;
  color: var(--u-text-1);
}

.u-sub {
  display: block;
  margin-top: 1px;
  font-size: 12px;
  color: var(--u-text-4);
}

.u-sub.u-warn {
  color: var(--u-warn);
}

.u-num {
  font-variant-numeric: tabular-nums;
}

.u-table td.is-num {
  text-align: right;
  font-variant-numeric: tabular-nums;
  color: var(--u-text-1);
}

.u-strong {
  color: var(--u-text-1);
  font-weight: 500;
}

.u-mono {
  font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
  font-size: 12px;
  letter-spacing: -0.01em;
}

.u-dim {
  color: var(--u-text-4);
}

.u-star {
  color: var(--u-warn);
}

.u-td-actions {
  width: 1%;
  white-space: nowrap;
}

.u-td-actions > * + * {
  margin-left: 4px;
}

.u-table.is-compact tbody td {
  height: 42px;
}

.u-empty,
.u-loading {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 46px 0;
  color: var(--u-text-4);
  font-size: 12.5px;
}

.u-table-foot {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  padding: 11px 18px;
  border-top: 1px solid var(--u-line-soft);
}

/* ---------- 分段视图切换 ---------- */
.u-grouped {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.u-grouped-tools {
  display: flex;
  align-items: center;
  gap: 8px;
}

.u-group {
  overflow: hidden;
}

.u-group-head {
  display: flex;
  align-items: center;
  gap: 10px;
  width: 100%;
  padding: 13px 18px;
  border: 0;
  border-bottom: 1px solid transparent;
  background: transparent;
  color: var(--u-text-1);
  text-align: left;
  cursor: pointer;
  user-select: none;
}

.u-group.is-open .u-group-head {
  border-bottom-color: var(--u-line-soft);
}

.u-group-chev {
  width: 16px;
  height: 16px;
  flex: none;
  color: var(--u-text-3);
  transition: transform 0.18s ease;
}

.u-group.is-open .u-group-chev {
  transform: rotate(90deg);
  color: #8fa8ff;
}

.u-group-name {
  font-size: 14px;
  font-weight: 600;
}

.u-group-summary {
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: 12px;
}

.u-group-traffic {
  flex: none;
  font-size: 12px;
}

.u-group-sync {
  margin-left: auto;
  font-size: 12px;
  flex: none;
}

.u-group-toggle {
  flex: none;
  font-size: 12px;
  color: #8fa8ff;
}

.u-group-body {
  padding: 0 12px 12px;
}

.u-group-body .u-table thead th {
  background: var(--u-card);
}

.u-group-body .u-table-foot {
  padding: 10px 6px 0;
  border-top: 1px solid var(--u-line-soft);
}

.u-empty-card {
  padding: 40px 0;
  text-align: center;
  font-size: 12.5px;
  color: var(--u-text-4);
}

/* ---------- 徽标 ---------- */
.u-badge {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  height: 21px;
  padding: 0 7px;
  border: 1px solid transparent;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 550;
  white-space: nowrap;
}

.u-badge-dot {
  width: 5px;
  height: 5px;
  flex: none;
  border-radius: 50%;
  background: currentColor;
}

.u-badge.is-ok {
  border-color: rgba(46, 212, 122, 0.32);
  background: rgba(46, 212, 122, 0.15);
  color: #5ce49b;
}

.u-badge.is-warn {
  border-color: rgba(245, 185, 61, 0.32);
  background: rgba(245, 185, 61, 0.15);
  color: #f7ca6b;
}

.u-badge.is-danger {
  border-color: rgba(243, 101, 111, 0.32);
  background: rgba(243, 101, 111, 0.15);
  color: #ff9ba2;
}

.u-badge.is-info {
  border-color: rgba(90, 169, 240, 0.32);
  background: rgba(90, 169, 240, 0.15);
  color: #8cc4f7;
}

.u-badge.is-brand {
  border-color: rgba(76, 111, 255, 0.34);
  background: rgba(76, 111, 255, 0.16);
  color: #9db2ff;
}

.u-badge.is-violet {
  border-color: rgba(167, 139, 250, 0.34);
  background: rgba(167, 139, 250, 0.16);
  color: #c0a9ff;
}

.u-badge.is-neutral {
  border-color: #2a2e35;
  background: #1d2026;
  color: #9aa1ab;
}

/* ---------- 提示条 ---------- */
.u-callout {
  display: flex;
  gap: 10px;
  padding: 13px 15px;
  border: 1px solid #24272e;
  border-radius: 11px;
  background: #15171b;
  line-height: 1.7;
}

.u-callout.is-warn {
  border-color: rgba(245, 185, 61, 0.24);
  background: rgba(245, 185, 61, 0.08);
}

.u-callout.is-danger {
  border-color: rgba(243, 101, 111, 0.24);
  background: rgba(243, 101, 111, 0.08);
}

.u-callout.is-info {
  border-color: rgba(90, 169, 240, 0.22);
  background: rgba(90, 169, 240, 0.08);
}

.u-callout-tight {
  padding: 10px 12px;
}

.u-callout-ico {
  width: 17px;
  height: 17px;
  flex: none;
  margin-top: 2px;
  color: #f7ca6b;
}

.u-callout.is-danger .u-callout-ico {
  color: #ff9ba2;
}

.u-callout.is-info .u-callout-ico {
  color: #8cc4f7;
}

.u-callout-main {
  min-width: 0;
  flex: 1;
}

.u-callout-head {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.u-callout-head b {
  font-size: 13px;
  font-weight: 600;
  color: var(--u-text-1);
}

.u-callout-note {
  margin: 6px 0 0;
  font-size: 12px;
  line-height: 1.65;
  color: var(--u-text-3);
}

/* ---------- 抽屉 ---------- */
.u-dhead {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding-right: 4px;
}

.u-avatar {
  display: grid;
  place-items: center;
  width: 40px;
  height: 40px;
  flex: none;
  border-radius: 50%;
  background: #2a2f38;
  box-shadow: 0 0 0 1px rgba(255, 255, 255, 0.08) inset;
  color: #e8eaee;
  font-size: 14px;
  font-weight: 600;
}

.u-dhead-main {
  min-width: 0;
}

.u-dhead-title {
  display: flex;
  align-items: center;
  gap: 7px;
  flex-wrap: wrap;
}

.u-dhead-title h2 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: var(--u-text-1);
}

.u-dhead-sub {
  margin: 4px 0 0;
  font-size: 12px;
  color: var(--u-text-4);
}

.u-dhead-actions {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-left: auto;
  flex: none;
}

.u-dtabs {
  position: sticky;
  top: 0;
  z-index: 5;
  display: flex;
  gap: 22px;
  padding: 0 20px;
  border-bottom: 1px solid var(--u-line-soft);
  background: #141519;
}

.u-dtabs button {
  position: relative;
  height: 42px;
  padding: 0;
  border: 0;
  background: none;
  color: var(--u-text-3);
  font-size: 12.5px;
  cursor: pointer;
}

.u-dtabs button:hover {
  color: var(--u-text-1);
}

.u-dtabs button.is-active {
  color: #fff;
  font-weight: 550;
}

.u-dtabs button.is-active::after {
  content: '';
  position: absolute;
  left: 0;
  right: 0;
  bottom: -1px;
  height: 2px;
  border-radius: 2px 2px 0 0;
  background: var(--u-brand);
}

.u-dbody {
  padding: 14px 20px 24px;
}

.u-stack {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.u-drawer-loading {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 260px;
}

.u-drawer-alert {
  margin: 14px 20px 0;
}

.u-drawer-foot {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 12px;
}

.u-dcard {
  border: 1px solid var(--u-line);
  border-radius: var(--u-radius);
  background: var(--u-card);
  overflow: hidden;
}

.u-dcard.is-muted {
  opacity: 0.92;
}

.u-dcard-head {
  display: flex;
  align-items: center;
  gap: 9px;
  padding: 13px 15px;
  border-bottom: 1px solid var(--u-line-soft);
}

.u-dcard-head h3 {
  margin: 0;
  font-size: 13px;
  font-weight: 600;
}

.u-dcard-body {
  padding: 14px 15px 15px;
}

.u-dcard-body.is-flush {
  padding: 0;
}

.u-field {
  display: flex;
  flex-direction: column;
  gap: 6px;
  min-width: 0;
}

.u-field > .u-badge {
  align-self: flex-start;
}

.u-label {
  font-size: 12px;
  color: var(--u-text-3);
}

.u-req {
  color: var(--u-danger);
  font-style: normal;
}

.u-grid-2 {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.u-grid-3 {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 14px 12px;
}

.u-span-2 {
  grid-column: span 2;
}

.u-mt-12 {
  margin-top: 12px;
}

.u-ml-auto {
  margin-left: auto;
}

.u-mr-auto {
  margin-right: auto;
}

.u-px {
  padding: 0 15px;
}

.u-form-actions {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 14px;
}

.u-sub-title {
  margin-bottom: 8px;
  font-size: 12.5px;
  font-weight: 550;
  color: var(--u-text-2);
}

.u-empty-small {
  padding: 14px 0;
  font-size: 12px;
  color: var(--u-text-4);
}

/* 待确认续费行 */
.u-cand {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
  margin-top: 8px;
  padding-top: 8px;
  border-top: 1px dashed rgba(245, 185, 61, 0.22);
}

.u-cand-info {
  min-width: 0;
  flex: 1;
}

.u-cand-info b {
  display: block;
  font-size: 12.5px;
  font-weight: 550;
  color: var(--u-text-1);
}

.u-cand-sub {
  display: block;
  margin-top: 2px;
  font-size: 12px;
  color: var(--u-text-4);
}

.u-cand-actions {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: none;
}

/* 收费记录 */
.u-list {
  display: flex;
  flex-direction: column;
}

.u-list-row {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 15px;
  border-bottom: 1px solid var(--u-line-soft);
}

.u-list-row:last-child {
  border-bottom: none;
}

.u-list-main {
  min-width: 0;
  flex: 1;
}

.u-list-title {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-wrap: wrap;
  font-size: 12.5px;
}

.u-list-sub {
  margin-top: 2px;
  font-size: 12px;
  color: var(--u-text-4);
}

.u-list-actions {
  display: flex;
  align-items: center;
  gap: 6px;
  flex: none;
}

/* 路径 */
.u-path-card {
  padding: 10px 12px;
  border: 1px solid var(--u-line-soft);
  border-radius: 10px;
  background: #111317;
}

.u-path-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  font-size: 12.5px;
}

.u-path-meta {
  margin-top: 3px;
  font-size: 12px;
  color: var(--u-text-3);
}

.u-history-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  padding: 8px 10px;
  margin-top: 6px;
  border: 1px solid var(--u-line-soft);
  border-radius: 9px;
  font-size: 12px;
  color: var(--u-text-3);
}

/* X-Panel 快照 */
.u-kv {
  display: flex;
  flex-direction: column;
}

.u-kv-row {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 15px;
  border-bottom: 1px solid var(--u-line-soft);
  font-size: 12.5px;
}

.u-kv-row:last-child {
  border-bottom: none;
}

.u-kv-row > span:first-child {
  width: 128px;
  flex: none;
  font-size: 12px;
  color: var(--u-text-4);
}

.u-kv-row b {
  font-weight: 500;
  color: var(--u-text-1);
}

/* ---------- 弹窗 ---------- */
.u-modal :deep(.n-card-header__main) {
  font-size: 14px;
  font-weight: 600;
}

.u-modal-body {
  font-size: 13px;
}

.u-modal-foot {
  display: flex;
  align-items: center;
  gap: 8px;
  justify-content: flex-end;
}

.u-alert {
  margin: 12px 15px 0;
}

.u-inline-error {
  display: flex;
  align-items: center;
  gap: 10px;
  margin: 12px 15px 0;
  padding: 9px 12px;
  border: 1px solid rgba(245, 185, 61, 0.24);
  border-radius: 9px;
  background: rgba(245, 185, 61, 0.08);
  font-size: 12.5px;
  color: #f7ca6b;
}

.u-import-errors {
  max-height: 160px;
  margin-top: 8px;
  overflow: auto;
  font-size: 12px;
  line-height: 1.7;
  color: #ff9ba2;
}

/* ---------- 脚注 ---------- */
.u-footnote {
  margin: 0 2px;
  font-size: 12px;
  line-height: 1.7;
  color: var(--u-text-4);
}

/* ---------- 响应式 ---------- */
@media (max-width: 1440px) {
  .u-kpis {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 1180px) {
  .u-split {
    flex-direction: column;
  }

  .u-facet {
    width: 100%;
  }

  .u-search {
    width: 100%;
  }
}


@media (max-width: 900px) {
  .u-head {
    flex-wrap: wrap;
    align-items: flex-start;
    gap: 12px;
  }

  .u-head-actions {
    margin-left: 0;
    flex-wrap: wrap;
    padding-bottom: 0;
  }
}
</style>
