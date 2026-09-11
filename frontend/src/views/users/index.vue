<script setup lang="ts">
import { computed, h, onMounted, reactive, ref } from 'vue';
import { useRoute } from 'vue-router';
import { NButton } from 'naive-ui';
import type { DataTableColumns } from 'naive-ui';
import {
  assignUserPath,
  cancelUserBillingRecord,
  createUserBillingRecord,
  confirmUserRenewal,
  clearUserPath,
  fetchUserDetail,
  fetchUserGroups,
  fetchUserPathAssets,
  fetchUsers,
  importUserBillingRecords,
  rejectUserRenewal,
  verifyUserBillingRecord,
  updateUser
} from '@/service/api';
import ModulePage from '@/components/project/module-page.vue';
import UserStatusTag from '@/components/project/user-status-tag.vue';
import DeviceTable from '@/components/project/device-table.vue';
import TrafficTrendChart from '@/components/project/traffic-trend-chart.vue';
import TrafficValue from '@/components/project/traffic-value.vue';

defineOptions({ name: 'UserManagement' });

const loading = ref(false);
const errorMessage = ref('');
const groups = ref<Api.Central.UserGroupSummary[]>([]);
const groupTotal = ref(0);
const stats = ref<Api.Central.UserListStats | null>(null);
const dataAt = ref('');
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
  serviceFrom: null as string | null,
  serviceTo: null as string | null,
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
const pathError = ref('');
const pathMode = ref<'relay' | 'landing' | 'external'>('relay');
const pathLandingNodeID = ref<string | null>(null);
const pathLandingInboundID = ref<string | null>(null);
const pathExitIPIDs = ref<string[]>([]);
const pathNotes = ref('');
const route = useRoute();
let pathAssetsRequestID = 0;

const filters = reactive({ page: 1, page_size: 20, keyword: '', status: 'all', billingType: 'all', nodeID: '' });

const statusOptions = [
  { label: '全部状态', value: 'all' },
  { label: '有效', value: 'active' },
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
  return assets.map(item => ({
    label: `${item.address} · ${item.sourceType === 's5' ? '独立 S5' : item.ownerNodeName || '节点出口'}${item.provider ? ` · ${item.provider}` : ''}`,
    value: item.id
  }));
});

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
  void loadGroups();
}

function resetOrderForm() {
  orderForm.billingCycle = detail.value?.billingCycle || 'monthly';
  orderForm.amount = detail.value?.billingAmount ?? detail.value?.monthlyFee ?? 0;
  const records = detail.value?.billingRecords || [];
  const latest = records.reduce<string | null>((value, record) => (value && value > record.serviceTo ? value : record.serviceTo), null);
  orderForm.serviceFrom = latest ? toDateInputValue(latest) : null;
  orderForm.serviceTo = detail.value?.expiresAt ? toDateInputValue(detail.value.expiresAt) : null;
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
  const serviceFrom = dayToRFC3339(orderForm.serviceFrom);
  const serviceTo = dayToRFC3339(orderForm.serviceTo);
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
  void loadGroups();
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
  void loadGroups();
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
    void loadGroups();
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
    void loadGroups();
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
  void loadGroups();
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
  void loadGroups();
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

async function openDetail(id: string) {
  selectedUserID.value = id;
  pathAssetsRequestID += 1;
  drawerVisible.value = true;
  detailLoading.value = true;
  detailError.value = '';
  detail.value = null;
  resetInitialOrderForm();
  pathAssets.value = null;
  pathLandingInbounds.value = [];
  pathAssetsLoading.value = false;
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
  window.$message?.success('中央业务信息已保存');
  void loadGroups();
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

async function confirmRenewal(candidate: Api.Central.UserRenewalCandidate) {
  if (!detail.value) return;
  const { data, error } = await confirmUserRenewal(detail.value.id, candidate.id, {
    billingCycle: candidate.billingCycle,
    amount: candidate.suggestedAmount
  });
  if (error || !data) {
    window.$message?.error('续费确认失败，请稍后重试');
    return;
  }
  detail.value = data;
  copyDetailToForm(data);
  window.$message?.success('续费已确认并计入财务');
  void loadGroups();
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
}

const columns: DataTableColumns<Api.Central.UserSummary> = [
  {
    title: '用户 / Inbound',
    key: 'name',
    minWidth: 220,
    render: row =>
      h('div', [
        h('div', { class: 'font-medium' }, row.name || '未命名用户'),
        h('div', { class: 'text-12px text-gray-500' }, row.inboundTag)
      ])
  },
  {
    title: '计费类型',
    key: 'billingType',
    width: 100,
    render: row => (row.billingType === 'free' ? h('span', { class: 'text-orange-500' }, '免费') : '收费')
  },
  { title: '落地机', key: 'landingNodeName', minWidth: 140, render: row => row.landingNodeName || '线路机直出' },
  {
    title: '出口 IP',
    key: 'exitIpAddress',
    minWidth: 220,
    render: row =>
      h('div', [
        h('div', { class: 'font-medium whitespace-normal break-all' }, row.exitIpAddress || '未配置'),
        row.exitIpOwnerNodeName
          ? h('div', { class: 'text-12px text-gray-500' }, `归属：${row.exitIpOwnerNodeName}`)
          : null
      ])
  },
  { title: '路径状态', key: 'pathMode', width: 110, render: row => pathModeLabel(row.pathMode) },
  {
    title: '状态',
    key: 'status',
    width: 100,
    render: row => h(UserStatusTag, { status: row.status, expiresAt: row.expiresAt })
  },
  { title: '到期时间', key: 'expiresAt', minWidth: 170, render: row => formatExpiry(row.expiresAt, row.status) },
  { title: '设备数', key: 'clientCount', width: 90 },
  {
    title: '累计流量',
    key: 'trafficBytes',
    minWidth: 120,
    render: row => h(TrafficValue, { value: row.trafficBytes })
  },
  { title: '最近在线', key: 'lastActivityAt', minWidth: 170, render: row => formatDate(row.lastActivityAt) },
  {
    title: '操作',
    key: 'actions',
    width: 90,
    fixed: 'right',
    render: row =>
      h(
        NButton,
        { size: 'small', type: 'primary', text: true, onClick: () => void openDetail(row.id) },
        { default: () => '查看' }
      )
  }
];

function ensureGroupUsers(nodeID: string) {
  if (!groupUsers[nodeID]) {
    groupUsers[nodeID] = { rows: [], total: 0, page: 1, pageSize: 20, loading: false, error: '' };
  }
  return groupUsers[nodeID];
}

function groupPagination(nodeID: string) {
  const state = ensureGroupUsers(nodeID);
  return {
    page: state.page,
    pageSize: state.pageSize,
    itemCount: state.total,
    showSizePicker: true,
    pageSizes: [20, 50, 100],
    onChange: (page: number) => {
      state.page = page;
      void loadGroupUsers(nodeID);
    },
    onUpdatePageSize: (pageSize: number) => {
      state.pageSize = pageSize;
      state.page = 1;
      void loadGroupUsers(nodeID);
    }
  };
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

function syncStatusType(status: Api.Central.UserGroupSummary['syncStatus']): 'success' | 'warning' | 'error' | 'default' {
  if (status === 'success') return 'success';
  if (status === 'failed') return 'error';
  if (status === 'unknown') return 'default';
  return 'warning';
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
  const previouslyExpanded = [...expandedGroups.value];
  const { data, error } = await fetchUserGroups({
    page: filters.page,
    page_size: filters.page_size,
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
    if (groups.value.length === 1) {
      expandedGroups.value = new Set([groups.value[0].nodeId]);
      void loadGroupUsers(groups.value[0].nodeId);
    } else {
      for (const nodeID of previouslyExpanded) {
        if (available.has(nodeID)) void loadGroupUsers(nodeID);
      }
    }
  }
  loading.value = false;
}

function submitFilters() {
  filters.page = 1;
  expandedGroups.value = new Set();
  for (const state of Object.values(groupUsers)) {
    state.page = 1;
    state.rows = [];
    state.total = 0;
    state.error = '';
  }
  loadGroups();
}

function resetFilters() {
  filters.keyword = '';
  filters.status = 'all';
  filters.billingType = 'all';
  filters.nodeID = '';
  submitFilters();
}

onMounted(() => {
  void loadGroups();
  const userId = typeof route.query.userId === 'string' ? route.query.userId : '';
  if (userId) void openDetail(userId);
});
</script>

<template>
  <div class="users-page">
    <ModulePage
      title="用户管理"
      description="按线路机查看业务用户；一个 Inbound 对应一个业务用户，Client / Email 仅作为设备凭证。"
      :loading="loading"
      :error="errorMessage"
      :empty="groups.length === 0"
      empty-description="暂无用户同步数据"
      :data-at="dataAt"
      @refresh="loadGroups"
    >
      <template #actions>
        <NButton size="small" secondary @click="importModalVisible = true">
          <template #icon><icon-mdi-upload /></template>
          导入历史账单
        </NButton>
        <NButton v-if="errorMessage" size="small" type="warning" secondary @click="loadGroups">重试</NButton>
      </template>
      <template #toolbar>
        <div>
          <div v-if="stats" class="users-kpis p-16px pb-0">
            <NCard :bordered="false" size="small">
              <NStatistic label="有效用户" :value="stats.active" />
              <div class="mt-4px text-12px text-gray-400">与总览页口径一致，按 Client 到期状态同步</div>
            </NCard>
            <NCard :bordered="false" size="small">
              <NStatistic label="付费用户" :value="stats.paid" />
              <div class="mt-4px text-12px text-gray-400">有效用户中的收费用户</div>
            </NCard>
            <NCard :bordered="false" size="small">
              <NStatistic label="免费用户" :value="stats.free" />
              <div class="mt-4px text-12px text-gray-400">有效用户中的免费用户</div>
            </NCard>
          </div>
          <div class="border-b border-gray-200 p-16px dark:border-gray-700">
            <NSpace wrap>
              <NInput
                v-model:value="filters.keyword"
                clearable
                class="w-240px"
                placeholder="搜索用户、Inbound 或线路机"
                @keyup.enter="submitFilters"
              />
              <NSelect v-model:value="filters.status" :options="statusOptions" class="w-140px" />
              <NSelect v-model:value="filters.billingType" :options="billingTypeOptions" class="w-140px" />
              <NSelect v-model:value="filters.nodeID" :options="nodeOptions" class="w-180px" filterable />
              <NButton type="primary" @click="submitFilters">
                <template #icon><icon-mdi-magnify /></template>
                查询
              </NButton>
              <NButton @click="resetFilters">重置</NButton>
              <NButton secondary :disabled="groups.length === 0" @click="expandAllGroups">全部展开</NButton>
              <NButton secondary :disabled="groups.length === 0" @click="collapseAllGroups">全部收起</NButton>
            </NSpace>
          </div>
        </div>
      </template>
      <div v-if="groups.length" class="space-y-12px p-16px pt-0">
        <div
          v-for="group in groups"
          :key="group.nodeId"
          class="user-group-card overflow-hidden rounded-6px border border-gray-200 bg-white dark:border-gray-700 dark:bg-dark"
        >
          <div
            class="flex items-start justify-between gap-16px border-b border-gray-200 px-16px py-14px dark:border-gray-700"
          >
            <div class="min-w-0 cursor-pointer" @click="toggleGroup(group)">
              <div class="flex flex-wrap items-center gap-8px">
                <span class="text-16px font-600">{{ group.nodeName }}</span>
                <NTag size="small" :type="nodeStatusType(group.status)">{{ nodeStatusLabel(group.status) }}</NTag>
                <NTag size="small" :type="syncStatusType(group.syncStatus)">{{ syncStatusLabel(group.syncStatus) }}</NTag>
              </div>
              <div class="mt-4px text-12px text-gray-500">
                {{ group.lastSyncAt ? `最近同步：${formatDate(group.lastSyncAt)}` : '尚未完成成功同步' }}
              </div>
            </div>
            <NButton text type="primary" @click="toggleGroup(group)">
              {{ isGroupExpanded(group.nodeId) ? '收起用户' : '查看用户' }}
            </NButton>
          </div>

          <div class="px-16px py-14px">
            <div class="flex flex-wrap items-center gap-x-16px gap-y-8px text-13px">
              <span>用户 <strong>{{ group.stats.total }}</strong></span>
              <span class="text-green-600">有效 <strong>{{ group.stats.active }}</strong></span>
              <span class="text-blue-600">付费 <strong>{{ group.stats.paid }}</strong></span>
              <span class="text-orange-500">免费 <strong>{{ group.stats.free }}</strong></span>
              <span class="text-yellow-600">即将到期 <strong>{{ group.stats.expiring }}</strong></span>
              <span class="text-gray-500">流量 <TrafficValue :value="group.stats.trafficBytes" /></span>
            </div>

            <div v-if="isGroupExpanded(group.nodeId)" class="mt-12px">
              <NSpin :show="groupUsers[group.nodeId]?.loading || false">
                <NAlert v-if="groupUsers[group.nodeId]?.error" type="warning" :show-icon="false" class="mb-12px">
                  {{ groupUsers[group.nodeId]?.error }}
                  <NButton size="small" class="ml-8px" @click="loadGroupUsers(group.nodeId)">重试</NButton>
                </NAlert>
                <NDataTable
                  v-else-if="groupUsers[group.nodeId]"
                  :columns="columns"
                  :data="groupUsers[group.nodeId].rows"
                  :pagination="groupPagination(group.nodeId)"
                  :bordered="false"
                  :single-line="false"
                  size="small"
                  :scroll-x="1020"
                />
              </NSpin>
            </div>
          </div>
        </div>
        <div v-if="groupTotal > filters.page_size" class="flex justify-end pt-4px">
          <NPagination
            v-model:page="filters.page"
            :page-count="Math.ceil(groupTotal / filters.page_size)"
            @update:page="loadGroups"
          />
        </div>
      </div>
    </ModulePage>

    <NModal v-model:show="importModalVisible" preset="card" title="导入历史账单" class="w-900px max-w-92vw">
      <NAlert type="info" :show-icon="true" class="mb-12px">
        每行一笔订单。表头支持：userId、nodeId、remoteInboundId、billingCycle、amount、serviceFrom、serviceTo、paidAt、orderType、verified、notes。
        userId 与 nodeId + remoteInboundId 二选一；日期支持 YYYY-MM-DD。verified 填 false 的记录会进入待核实，不计入财务。
      </NAlert>
      <NInput
        v-model:value="importText"
        type="textarea"
        :autosize="{ minRows: 8, maxRows: 18 }"
        placeholder="userId,billingCycle,amount,serviceFrom,serviceTo,paidAt,orderType,verified,notes\n用户ID,monthly,100,2026-03-01,2026-04-01,2026-03-01,initial,true,首购"
      />
      <div v-if="importPreview" class="mt-12px">
        <NAlert :type="importPreview.errors.length ? 'error' : 'success'" :show-icon="true">
          {{ importPreview.errors.length ? `预览发现 ${importPreview.errors.length} 个错误，请修正后重新预览` : `预览通过：${importPreview.records.length} 笔，其中 ${importPreview.unverified} 笔待核实` }}
        </NAlert>
        <div v-if="importPreview.errors.length" class="mt-8px max-h-160px overflow-auto text-12px text-red-500">
          <div v-for="item in importPreview.errors" :key="`${item.row}-${item.message}`">第 {{ item.row }} 行：{{ item.message }}</div>
        </div>
      </div>
      <template #footer>
        <NSpace justify="end">
          <NButton @click="importModalVisible = false">关闭</NButton>
          <NButton :loading="importLoading" @click="previewBillingImport">预览校验</NButton>
          <NButton type="primary" :loading="importLoading" :disabled="!importPreview || importPreview.errors.length > 0" @click="submitBillingImport">
            确认导入
          </NButton>
        </NSpace>
      </template>
    </NModal>

    <NModal v-model:show="orderModalVisible" preset="card" title="新增订单" class="w-680px max-w-92vw">
      <NForm label-placement="top">
        <div class="grid gap-x-12px md:grid-cols-2">
          <NFormItem label="订单类型" required>
            <NSelect
              v-model:value="orderForm.orderType"
              :options="[
                { label: '首购', value: 'initial' },
                { label: '续费', value: 'renewal' },
                { label: '逾期恢复', value: 'recovery' }
              ]"
            />
          </NFormItem>
          <NFormItem label="收费周期" required>
            <NSelect
              v-model:value="orderForm.billingCycle"
              :options="[
                { label: '月付', value: 'monthly' },
                { label: '年付', value: 'annual' }
              ]"
            />
          </NFormItem>
          <NFormItem label="订单金额（CNY）" required>
            <NInputNumber v-model:value="orderForm.amount" :min="0" :max="100000000" :precision="2" class="w-full" />
          </NFormItem>
          <NFormItem label="实际收款日期" required>
            <NDatePicker v-model:formatted-value="orderForm.paidAt" type="date" value-format="yyyy-MM-dd" class="w-full" />
          </NFormItem>
          <NFormItem label="服务开始日期" required>
            <NDatePicker v-model:formatted-value="orderForm.serviceFrom" type="date" value-format="yyyy-MM-dd" class="w-full" />
          </NFormItem>
          <NFormItem label="服务结束日期" required>
            <NDatePicker v-model:formatted-value="orderForm.serviceTo" type="date" value-format="yyyy-MM-dd" class="w-full" />
          </NFormItem>
        </div>
        <NFormItem label="订单备注"><NInput v-model:value="orderForm.notes" maxlength="2000" /></NFormItem>
      </NForm>
      <template #footer>
        <NSpace justify="end">
          <NButton @click="orderModalVisible = false">取消</NButton>
          <NButton type="primary" :loading="orderSaving" @click="createAdditionalOrder">保存订单</NButton>
        </NSpace>
      </template>
    </NModal>

    <NModal v-model:show="verifyModalVisible" preset="card" title="核验历史订单" class="w-520px max-w-92vw">
      <NAlert type="warning" :show-icon="true" class="mb-12px">核验后该订单会计入实收、服务期收入和续费统计。</NAlert>
      <NForm label-placement="top">
        <NFormItem label="实际收款日期" required>
          <NDatePicker v-model:formatted-value="verifyPaidAt" type="date" value-format="yyyy-MM-dd" class="w-full" />
        </NFormItem>
        <NFormItem label="核验备注"><NInput v-model:value="verifyNotes" type="textarea" /></NFormItem>
      </NForm>
      <template #footer>
        <NSpace justify="end">
          <NButton @click="verifyModalVisible = false">取消</NButton>
          <NButton type="primary" :loading="verifySaving" @click="submitVerifyRecord">确认核验</NButton>
        </NSpace>
      </template>
    </NModal>

    <NDrawer v-model:show="drawerVisible" :width="680" placement="right">
      <NDrawerContent title="用户详情" closable>
        <NSpin :show="detailLoading">
          <NAlert v-if="detailError" type="error" class="mb-16px" :show-icon="true">
            {{ detailError }}
            <NButton class="ml-8px" size="small" @click="openDetail(selectedUserID)">重试</NButton>
          </NAlert>
          <template v-if="detail">
            <div class="mb-16px flex items-center justify-between gap-12px">
              <div>
                <div class="text-18px font-600">{{ detail.displayName }}</div>
                <div class="mt-4px text-12px text-gray-500">一个 Reality Inbound 对应一个业务用户</div>
              </div>
              <UserStatusTag :status="detail.status" :expires-at="detail.expiresAt" />
            </div>

            <NAlert type="info" :show-icon="true" class="mb-16px">
              下列 Inbound、Client、流量和到期时间均来自 X-Panel 同步，只读展示；保存只会更新中央面板的业务信息。
            </NAlert>

            <NCard title="中央业务信息" size="small" :segmented="{ content: true }">
              <NForm label-placement="top">
                <div class="grid gap-x-12px md:grid-cols-2">
                  <NFormItem label="业务用户名称" required>
                    <NInput v-model:value="editForm.displayName" maxlength="120" show-count />
                  </NFormItem>
                  <NFormItem label="收费周期">
                    <NSelect
                      v-model:value="editForm.billingCycle"
                      :options="[
                        { label: '月付', value: 'monthly' },
                        { label: '年付', value: 'annual' }
                      ]"
                    />
                  </NFormItem>
                  <NFormItem label="计费类型">
                    <NSelect
                      v-model:value="editForm.billingType"
                      :options="[
                        { label: '收费用户', value: 'paid' },
                        { label: '免费用户', value: 'free' }
                      ]"
                    />
                  </NFormItem>
                  <NFormItem :label="editForm.billingCycle === 'annual' ? '年费（CNY）' : '月费（CNY）'">
                    <NInputNumber
                      v-model:value="editForm.billingAmount"
                      :min="0"
                      :max="100000000"
                      :precision="2"
                      class="w-full"
                    />
                    <div v-if="editForm.billingCycle === 'annual'" class="mt-4px text-12px text-gray-500">
                      月均 {{ formatMoney((editForm.billingAmount || 0) / 12) }}
                    </div>
                  </NFormItem>
                </div>
                <NFormItem v-if="editForm.billingType === 'free'" label="免费原因">
                  <NInput v-model:value="editForm.freeReason" maxlength="200" placeholder="例如：朋友体验、内部测试" />
                </NFormItem>
                <NFormItem label="备注">
                  <NInput
                    v-model:value="editForm.notes"
                    type="textarea"
                    maxlength="2000"
                    show-count
                    :autosize="{ minRows: 3, maxRows: 6 }"
                  />
                </NFormItem>
                <NButton type="primary" :loading="saving" @click="saveBusinessFields">保存业务信息</NButton>
              </NForm>
            </NCard>

            <NCard v-if="detail.billingRecords.length === 0" title="首笔订单" size="small" class="mt-16px" :segmented="{ content: true }">
              <NAlert type="info" :show-icon="true" class="mb-12px">
                首笔订单对应 Agent 同步过来的线路机 Inbound
                业务用户。服务开始日期请按真实业务手动填写，服务结束日期直接采用 Agent 同步的 Client 到期日。
              </NAlert>
              <NForm label-placement="top">
                <div class="grid gap-x-12px md:grid-cols-2">
                  <NFormItem label="订单金额（CNY）" required>
                    <NInputNumber
                      v-model:value="initialOrderForm.amount"
                      :min="0"
                      :max="100000000"
                      :precision="2"
                      class="w-full"
                    />
                  </NFormItem>
                  <NFormItem label="实际收款日期" required>
                    <NDatePicker
                      v-model:formatted-value="initialOrderForm.paidAt"
                      type="date"
                      value-format="yyyy-MM-dd"
                      class="w-full"
                    />
                  </NFormItem>
                  <NFormItem label="服务开始日期" required>
                    <NDatePicker
                      v-model:formatted-value="initialOrderForm.serviceFrom"
                      type="date"
                      value-format="yyyy-MM-dd"
                      placeholder="请选择服务开始日期"
                      clearable
                      class="w-full"
                    />
                  </NFormItem>
                  <NFormItem label="服务结束日期（Agent 同步）" required>
                    <NDatePicker
                      v-model:formatted-value="initialOrderForm.serviceTo"
                      type="date"
                      value-format="yyyy-MM-dd"
                      placeholder="Agent 尚未同步到期日"
                      disabled
                      class="w-full"
                    />
                  </NFormItem>
                </div>
                <NFormItem label="订单备注">
                  <NInput v-model:value="initialOrderForm.notes" maxlength="2000" placeholder="可选" />
                </NFormItem>
                <NButton type="primary" :loading="initialOrderSaving" @click="createInitialOrder">确认首笔订单</NButton>
              </NForm>
            </NCard>

            <NCard
              v-if="detail.renewalCandidates.some(candidate => candidate.status === 'pending')"
              title="待确认续费"
              size="small"
              class="mt-16px"
              :segmented="{ content: true }"
            >
              <div
                v-for="candidate in detail.renewalCandidates.filter(item => item.status === 'pending')"
                :key="candidate.id"
                class="flex flex-wrap items-center justify-between gap-10px border-b border-gray-200 py-8px last:border-b-0 dark:border-gray-700"
              >
                <div>
                  <div class="text-13px font-500">
                    {{ cycleLabel(candidate.billingCycle) }} · {{ formatMoney(candidate.suggestedAmount) }}
                  </div>
                  <div class="mt-3px text-12px text-gray-500">
                    {{ formatDate(candidate.oldExpiryAt) }} → {{ formatDate(candidate.newExpiryAt) }}
                  </div>
                </div>
                <NSpace size="small">
                  <NButton size="small" type="primary" @click="confirmRenewal(candidate)">确认续费</NButton>
                  <NButton size="small" secondary @click="rejectRenewal(candidate)">非收费变更</NButton>
                </NSpace>
              </div>
            </NCard>

            <NCard v-if="detail.billingRecords.length > 0" title="收费记录" size="small" class="mt-16px">
              <template #header-extra>
                <NButton size="small" type="primary" secondary @click="openOrderModal">新增订单</NButton>
              </template>
              <div
                v-for="record in detail.billingRecords"
                :key="record.id"
                class="flex flex-wrap items-center justify-between gap-8px border-b border-gray-200 py-7px last:border-b-0 dark:border-gray-700"
              >
                <div>
                  <div class="text-12px">
                    {{ orderTypeLabel(record.orderType) }} · {{ cycleLabel(record.billingCycle) }} ·
                    {{ formatMoney(record.amount, record.currency) }}
                    <NTag v-if="record.verificationStatus === 'unverified'" size="small" type="warning" class="ml-6px">待核实</NTag>
                    <NTag v-else-if="record.origin === 'historical_import'" size="small" type="info" class="ml-6px">历史导入</NTag>
                    <NTag v-if="record.status === 'cancelled'" size="small" type="error" class="ml-6px">已取消</NTag>
                  </div>
                  <div class="mt-3px text-12px text-gray-500">
                    服务 {{ formatDate(record.serviceFrom) }} → {{ formatDate(record.serviceTo) }} · 收款
                    {{ formatDate(record.paidAt) }}
                  </div>
                </div>
                <NSpace v-if="record.source === 'manual' && record.status !== 'cancelled'" size="small">
                  <NButton v-if="record.verificationStatus === 'unverified'" size="tiny" type="primary" @click="openVerifyModal(record)">核验</NButton>
                  <NButton size="tiny" secondary type="error" @click="cancelRecord(record)">取消</NButton>
                </NSpace>
              </div>
            </NCard>

            <NCard title="X-Panel Inbound 快照（只读）" size="small" class="mt-16px">
              <NDescriptions label-placement="left" :column="1" bordered size="small">
                <NDescriptionsItem label="线路机">
                  {{ detail.node.name || '--'
                  }}{{ detail.node.type ? `（${detail.node.type === 'relay' ? '线路机' : '落地机'}）` : '' }}
                </NDescriptionsItem>
                <NDescriptionsItem label="Inbound">
                  {{ detail.inbound.tag || detail.inbound.remoteId || '--' }}
                </NDescriptionsItem>
                <NDescriptionsItem label="协议 / 端口">
                  {{ detail.inbound.protocol || '--' }} / {{ detail.inbound.port || '--' }}
                </NDescriptionsItem>
                <NDescriptionsItem label="启用状态">{{ detail.inbound.enabled ? '启用' : '停用' }}</NDescriptionsItem>
                <NDescriptionsItem label="到期时间">
                  {{ formatExpiry(detail.expiresAt, detail.status) }}
                </NDescriptionsItem>
                <NDescriptionsItem label="累计流量"><TrafficValue :value="detail.inbound.allTime" /></NDescriptionsItem>
                <NDescriptionsItem label="上次同步">{{ formatDate(detail.inbound.lastSeenAt) }}</NDescriptionsItem>
              </NDescriptions>
            </NCard>

            <NCard title="设备 Client / Email（只读）" size="small" class="mt-16px">
              <DeviceTable :devices="detail.clients" />
            </NCard>

            <NCard title="流量趋势（只读）" size="small" class="mt-16px">
              <TrafficTrendChart :user-id="detail.id" />
            </NCard>

            <NCard title="用户路径" size="small" class="mt-16px" :segmented="{ content: true }">
              <NAlert v-if="pathAssetsError" type="warning" :show-icon="true" class="mb-12px">
                {{ pathAssetsError }}
                <NButton size="small" class="ml-8px" @click="() => loadPathAssets()">重试</NButton>
              </NAlert>
              <NForm label-placement="top">
                <NFormItem label="线路机（自动）">
                  <NInput :value="detail.node.name || '未关联线路机'" readonly />
                </NFormItem>
                <NFormItem label="出口路径">
                  <NSelect
                    :value="pathMode"
                    :options="pathModeOptions"
                    :disabled="pathAssetsLoading"
                    @update:value="handlePathModeChange"
                  />
                </NFormItem>
                <NFormItem v-if="pathMode === 'landing'" label="落地机（执行节点）" required>
                  <NSelect
                    :value="pathLandingNodeID"
                    :options="pathLandingNodeOptions"
                    :loading="pathAssetsLoading"
                    :disabled="pathAssetsLoading"
                    filterable
                    placeholder="选择启用的落地机"
                    @update:value="handlePathLandingNodeChange"
                  />
                  <div
                    v-if="pathAssets && pathLandingNodeOptions.length === 0"
                    class="mt-4px text-12px text-orange-500"
                  >
                    暂无启用的落地机，请先在节点管理中添加并启用。
                  </div>
                </NFormItem>
                <NFormItem v-if="pathMode === 'landing'" label="落地 Inbound（基础设施入口）" required>
                  <NSelect
                    v-model:value="pathLandingInboundID"
                    :options="pathInboundOptions"
                    :loading="pathLandingInboundsLoading"
                    :disabled="pathAssetsLoading || !pathLandingNodeID"
                    filterable
                    placeholder="选择该落地机已同步的 Inbound"
                  />
                  <div
                    v-if="pathLandingNodeID && pathInboundOptions.length === 0"
                    class="mt-4px text-12px text-orange-500"
                  >
                    {{
                      pathAssets?.landingNodes.find(node => node.id === pathLandingNodeID)?.inboundState === 'pending'
                        ? '该落地机尚未完成首次 Agent 同步。'
                        : '该落地机暂无启用的 Inbound。'
                    }}
                  </div>
                </NFormItem>
                <NFormItem label="固定出口 IP">
                  <NSelect
                    v-model:value="pathExitIPIDs"
                    :options="pathExitIpOptions"
                    :loading="pathAssetsLoading"
                    :disabled="pathAssetsLoading || !detail.node.id"
                    multiple
                    clearable
                    filterable
                    max-tag-count="responsive"
                    placeholder="选择当前路径允许的出口 IP"
                  />
                </NFormItem>
                <NFormItem label="路径备注">
                  <NInput
                    v-model:value="pathNotes"
                    type="textarea"
                    maxlength="500"
                    show-count
                    :autosize="{ minRows: 2, maxRows: 4 }"
                    placeholder="可选"
                  />
                </NFormItem>
                <NAlert v-if="pathError" type="error" :show-icon="true" class="mb-12px">{{ pathError }}</NAlert>
                <NSpace>
                  <NButton
                    type="primary"
                    :loading="pathSaving"
                    :disabled="
                      pathAssetsLoading ||
                        (pathMode === 'landing' && (!pathLandingNodeID || !pathLandingInboundID)) ||
                        !pathExitIPIDs.length
                    "
                    @click="savePathAssignment"
                  >
                    保存用户路径
                  </NButton>
                  <NButton v-if="detail.path" secondary :loading="pathSaving" @click="confirmRemovePathAssignment">
                    解除路径
                  </NButton>
                </NSpace>
              </NForm>

              <div class="mt-16px border-t border-gray-200 pt-12px">
                <div class="mb-8px text-13px font-500">当前生效路径</div>
                <NEmpty v-if="!detail.path" description="尚未配置用户路径" size="small" />
                <div v-else class="rounded-6px bg-gray-50 p-10px dark:bg-dark">
                  <div class="flex items-center justify-between gap-8px">
                    <span class="font-500">
                      {{ detail.path.relayNodeName || detail.node.name || '--' }} →
                      {{ detail.path.landingNodeName || '线路机直出' }}
                    </span>
                    <NTag size="small" :type="detail.path.valid ? 'success' : 'warning'">
                      {{ detail.path.valid ? '有效' : '需检查' }}
                    </NTag>
                  </div>
                  <div class="mt-4px text-12px text-gray-500">模式：{{ pathModeLabel(detail.path.mode) }}</div>
                  <div class="mt-4px text-12px text-gray-500">
                    出口：{{ pathExitIPText(detail.path) }} · 归属：{{ detail.path.exitIpOwnerNodeName || '独立 S5' }}
                  </div>
                  <div
                    v-if="detail.path.landingInboundTag || detail.path.landingInboundId"
                    class="mt-4px text-12px text-gray-500"
                  >
                    落地 Inbound：{{ detail.path.landingInboundTag || detail.path.landingInboundId }}
                  </div>
                  <div v-if="detail.path.notes" class="mt-4px text-12px text-gray-500">
                    备注：{{ detail.path.notes }}
                  </div>
                </div>
              </div>

              <div v-if="detail.pathHistory.length > 0" class="mt-16px border-t border-gray-200 pt-12px">
                <div class="mb-8px text-13px font-500">路径历史</div>
                <div class="space-y-8px">
                  <div
                    v-for="history in detail.pathHistory"
                    :key="history.id"
                    class="rounded-6px border border-gray-200 p-8px dark:border-gray-700"
                  >
                    <div class="flex items-center justify-between gap-8px text-12px">
                      <span>
                        {{ history.relayNodeName || '--' }} → {{ history.landingNodeName || '线路机直出' }} ·
                        {{ pathExitIPText(history) }}
                      </span>
                      <span class="text-gray-500">
                        {{ formatDate(history.activeFrom)
                        }}{{ history.activeTo ? ` 至 ${formatDate(history.activeTo)}` : ' 至今' }}
                      </span>
                    </div>
                  </div>
                </div>
              </div>
            </NCard>
          </template>
        </NSpin>
      </NDrawerContent>
    </NDrawer>
  </div>
</template>

<style scoped>
.users-kpis {
  display: grid;
  grid-template-columns: minmax(0, 1fr);
  gap: 12px;
}

@media (min-width: 768px) {
  .users-kpis {
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }
}
</style>
