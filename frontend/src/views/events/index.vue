<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';
import { useRouter } from 'vue-router';
import { NPagination, NSelect } from 'naive-ui';
import {
  confirmUserRenewal,
  fetchEventStats,
  fetchEventSummary,
  fetchEvents,
  markAllEventsRead,
  markEventRead,
  resetUserFromReplacementEvent,
  rejectUserRenewal,
  requestNodeSync,
  resolveEvent
} from '@/service/api';
import DataFreshness from '@/components/project/data-freshness.vue';
import SvgIcon from '@/components/custom/svg-icon.vue';

defineOptions({ name: 'EventManagement' });

type EventTab = 'pending' | 'all';
type TimeRange = 'all' | '1d' | '7d' | '30d';

const loading = ref(false);
const actionLoading = ref('');
const modalSaving = ref(false);
const router = useRouter();
const errorMessage = ref('');
const rows = ref<Api.Central.EventSummary[]>([]);
const total = ref(0);
const dataAt = ref('');
const statsReady = ref(false);
const statsCapped = ref(false);
const allCount = ref<number | null>(null);
const pendingCount = ref<number | null>(null);
const renewalTarget = ref<Api.Central.EventSummary | null>(null);
const replacementTarget = ref<Api.Central.EventSummary | null>(null);

const filters = reactive({
  page: 1,
  page_size: 20,
  keyword: '',
  severity: 'all',
  category: 'all',
  timeRange: 'all' as TimeRange,
  tab: 'pending' as EventTab
});

const stats = reactive({
  pending: 0,
  attention: 0,
  attentionError: 0,
  attentionReplacement: 0,
  resolvedToday: 0,
  autoRecoveredToday: 0,
  resolved7d: 0,
  avgHours: null as number | null,
  medianHours: null as number | null
});

const distribution = ref<Array<{ label: string; count: number; tone: string; width: string }>>([]);
const distributionTotal = ref(0);

const severitySegments = [
  { label: '全部级别', value: 'all' },
  { label: '信息', value: 'info' },
  { label: '警告', value: 'warning' },
  { label: '错误', value: 'error' }
];

const categoryOptions = [
  { label: '全部类型', value: 'all' },
  { label: '业务事件', value: 'business' },
  { label: '节点事件', value: 'node' },
  { label: '同步事件', value: 'sync' },
  { label: '系统事件', value: 'system' }
];

const timeRangeOptions = [
  { label: '全部时间', value: 'all' },
  { label: '近 24 小时', value: '1d' },
  { label: '近 7 天', value: '7d' },
  { label: '近 30 天', value: '30d' }
];

const pageSizeOptions = [
  { label: '每页 20 条', value: 20 },
  { label: '每页 50 条', value: 50 },
  { label: '每页 100 条', value: 100 }
];

const typeLabels: Record<string, string> = {
  renewal_candidate_detected: '续费候选',
  renewal_non_billable: '非收费变更',
  client_set_replacement_detected: '客户更换',
  new_user_profile_required: '新增用户',
  client_expiry_mismatch: '到期不一致',
  sync_failed: '同步失败',
  xpanel_sync_failed: '采集异常',
  xpanel_sync_recovered: '采集恢复',
  node_recovered: '节点恢复',
  node_offline: '节点离线',
  traffic_reset: '流量回退',
  inbound_missing: 'Inbound 缺失',
  inbound_archived: 'Inbound 归档',
  demo_seed: '演示事件'
};

const severityLabels: Record<string, string> = {
  info: '信息',
  warning: '警告',
  error: '错误'
};

const typeGroups = [
  {
    label: '续费候选',
    tone: 'is-warn',
    types: ['renewal_candidate_detected', 'renewal_non_billable']
  },
  {
    label: '同步异常',
    tone: 'is-danger',
    types: ['sync_failed', 'xpanel_sync_failed']
  },
  {
    label: '节点事件',
    tone: '',
    types: ['node_offline', 'node_recovered', 'traffic_reset', 'inbound_missing', 'inbound_archived']
  },
  {
    label: '业务资料',
    tone: 'is-violet',
    types: ['new_user_profile_required', 'client_expiry_mismatch']
  },
  {
    label: '客户更换',
    tone: 'is-cyan',
    types: ['client_set_replacement_detected']
  },
  {
    label: '其它',
    tone: 'is-ok',
    types: [] as string[]
  }
];

function eventTypeLabel(type: string) {
  return typeLabels[type] || type;
}

function payloadValue(row: Api.Central.EventSummary, key: string) {
  return row.payload && typeof row.payload[key] === 'string' ? (row.payload[key] as string) : '';
}

function payloadNumber(row: Api.Central.EventSummary, key: string) {
  const value = row.payload?.[key];
  const parsed = typeof value === 'number' ? value : Number(value);
  return Number.isFinite(parsed) && parsed > 0 ? parsed : 0;
}

function targetPayloadValue(target: Api.Central.EventSummary | null, key: string) {
  return target ? payloadValue(target, key) : '';
}

function formatDateTime(value?: string | null) {
  if (!value) return '--';
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return value;
  const pad = (input: number) => String(input).padStart(2, '0');
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())} ${pad(date.getHours())}:${pad(date.getMinutes())}`;
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

function formatMoney(value: number, currency = 'CNY') {
  return new Intl.NumberFormat('zh-CN', { style: 'currency', currency }).format(value);
}

function severityLabel(severity: string) {
  return severityLabels[severity] || '未知';
}

function isRenewal(row: Api.Central.EventSummary) {
  return row.type === 'renewal_candidate_detected' && row.status !== 'resolved' && row.status !== 'dismissed';
}

function isClientReplacement(row: Api.Central.EventSummary) {
  return row.type === 'client_set_replacement_detected' && row.status !== 'resolved' && row.status !== 'dismissed';
}

function isNewUserProfile(row: Api.Central.EventSummary) {
  return row.type === 'new_user_profile_required' && row.status !== 'resolved' && row.status !== 'dismissed';
}

function isPendingEvent(row: Api.Central.EventSummary) {
  return row.status !== 'resolved' && row.status !== 'dismissed' && (row.requiresAction || Boolean(row.actionType?.trim()));
}

function canMarkRead(row: Api.Central.EventSummary) {
  return row.status !== 'resolved' && row.status !== 'dismissed' && !isPendingEvent(row);
}

function eventStatusLabel(row: Api.Central.EventSummary) {
  if (row.status === 'resolved' || row.status === 'dismissed') return '已处理';
  if (isPendingEvent(row)) return '待处理';
  return row.acknowledged ? '已读' : '未读';
}

function statusTone(row: Api.Central.EventSummary) {
  if (row.status === 'resolved' || row.status === 'dismissed') return 'is-ok';
  if (isPendingEvent(row)) return 'is-warn';
  return 'is-neutral';
}

function severityDotClass(row: Api.Central.EventSummary) {
  if (row.severity === 'error') return 'is-danger';
  if (row.severity === 'warning') return 'is-warn';
  return 'is-info';
}

function typeBadgeClass(row: Api.Central.EventSummary) {
  switch (row.type) {
    case 'renewal_candidate_detected':
    case 'renewal_non_billable':
    case 'xpanel_sync_failed':
    case 'client_expiry_mismatch':
    case 'traffic_reset':
      return 'is-warn';
    case 'new_user_profile_required':
      return 'is-brand';
    case 'client_set_replacement_detected':
      return 'is-violet';
    case 'sync_failed':
    case 'node_offline':
      return 'is-danger';
    case 'xpanel_sync_recovered':
    case 'node_recovered':
      return 'is-ok';
    case 'inbound_missing':
    case 'inbound_archived':
      return 'is-neutral';
    default:
      return row.severity === 'error' ? 'is-danger' : row.severity === 'warning' ? 'is-warn' : 'is-neutral';
  }
}

function hasUserId(row: Api.Central.EventSummary) {
  return Boolean(payloadValue(row, 'userId'));
}

function userLabel(row: Api.Central.EventSummary) {
  return payloadValue(row, 'userName') || payloadValue(row, 'userId') || '—';
}

function hasRowActions(row: Api.Central.EventSummary) {
  return (
    isRenewal(row) ||
    isNewUserProfile(row) ||
    isClientReplacement(row) ||
    (row.status !== 'resolved' && Boolean(row.actionType)) ||
    (!row.acknowledged && canMarkRead(row)) ||
    (!isRenewal(row) && !isNewUserProfile(row) && !isClientReplacement(row) && hasUserId(row))
  );
}

function notifyEventChanged() {
  window.dispatchEvent(new Event('xpanel-events-updated'));
}

function openUser(row: Api.Central.EventSummary) {
  const userId = payloadValue(row, 'userId');
  if (userId) router.push({ name: 'users', query: { userId } });
}

async function completeNewUserProfile(row: Api.Central.EventSummary) {
  openUser(row);
}


function isSameLocalDay(value: string | null | undefined) {
  if (!value) return false;
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return false;
  const now = new Date();
  return (
    date.getFullYear() === now.getFullYear() && date.getMonth() === now.getMonth() && date.getDate() === now.getDate()
  );
}

function timeRangeBounds(range: TimeRange): { from?: string; to?: string } {
  if (range === 'all') return {};
  const days = range === '1d' ? 1 : range === '7d' ? 7 : 30;
  return { from: new Date(Date.now() - days * 24 * 60 * 60 * 1000).toISOString() };
}

async function loadEvents() {
  loading.value = true;
  errorMessage.value = '';
  const status = filters.tab === 'pending' ? 'pending' : undefined;
  const range = timeRangeBounds(filters.timeRange);
  const { data, error } = await fetchEvents({
    keyword: filters.keyword || undefined,
    severity: filters.severity === 'all' ? undefined : (filters.severity as Api.Central.EventSeverity),
    category: filters.category === 'all' ? undefined : (filters.category as Api.Central.EventCategory),
    status,
    from: range.from,
    to: range.to,
    page: filters.page,
    page_size: filters.page_size
  });
  if (error || !data) {
    errorMessage.value = '中央后端暂不可用，无法读取事件';
    rows.value = [];
    total.value = 0;
    dataAt.value = '';
  } else {
    rows.value = data.items;
    total.value = data.total;
    dataAt.value = data.dataAt || '';
  }
  // 处理完最后一条后当前页可能已越界，回退到最后一页避免空白表格
  if (!rows.value.length && total.value > 0 && filters.page > 1) {
    filters.page = Math.min(filters.page, Math.max(1, Math.ceil(total.value / filters.page_size)));
    loading.value = false;
    await loadEvents();
    return;
  }
  loading.value = false;
}

async function loadStats() {
  const { data, error } = await fetchEventStats({ window: '7d' });
  if (error || !data) {
    return;
  }

  allCount.value = data.totalCount;
  pendingCount.value = data.pendingCount;
  statsCapped.value = false;

  const typeCount = (type: string) => data.byType.find(item => item.type === type)?.count ?? 0;
  const severityCount = (severity: string) => data.bySeverity.find(item => item.severity === severity)?.count ?? 0;

  stats.pending = data.pendingCount;
  stats.attention = data.attentionCount;
  stats.attentionError = severityCount('error');
  stats.attentionReplacement = typeCount('client_set_replacement_detected');
  stats.resolvedToday = data.processedToday;
  stats.autoRecoveredToday = data.autoRecoveredToday;
  stats.resolved7d = data.resolvedInWindow;
  stats.avgHours = data.avgHandleMinutes > 0 ? data.avgHandleMinutes / 60 : null;
  stats.medianHours = data.medianHandleMinutes > 0 ? data.medianHandleMinutes / 60 : null;

  const totalInWindow = data.byType.reduce((sum, item) => sum + item.count, 0);
  distributionTotal.value = totalInWindow;
  const grouped = typeGroups.map(group => ({
    label: group.label,
    tone: group.tone,
    count: group.types.length
      ? group.types.reduce((sum, type) => sum + typeCount(type), 0)
      : data.byType
          .filter(item => !typeGroups.some(other => other.types.includes(item.type)))
          .reduce((sum, item) => sum + item.count, 0)
  }));
  distribution.value = grouped
    .filter(group => group.count > 0)
    .sort((left, right) => right.count - left.count)
    .map(group => ({
      ...group,
      width: `${Math.max(4, Math.round((group.count / Math.max(1, totalInWindow)) * 100))}%`
    }));
  statsReady.value = true;
}

async function reloadAll() {
  await Promise.all([loadEvents(), loadStats()]);
}

function selectTab(value: EventTab) {
  if (filters.tab === value) return;
  filters.tab = value;
  filters.page = 1;
  void loadEvents();
}

function submitFilters() {
  filters.page = 1;
  void loadEvents();
}

function resetFilters() {
  filters.keyword = '';
  filters.severity = 'all';
  filters.category = 'all';
  filters.timeRange = 'all';
  submitFilters();
}

function selectSeverity(value: string) {
  filters.severity = value;
  submitFilters();
}

function changeTimeRange(value: TimeRange) {
  filters.timeRange = value;
  submitFilters();
}

function changePageSize(value: number) {
  filters.page_size = value;
  filters.page = 1;
  void loadEvents();
}

function changePage(page: number) {
  filters.page = page;
  void loadEvents();
}

async function readEvent(row: Api.Central.EventSummary) {
  if (row.acknowledged || !canMarkRead(row)) return;
  actionLoading.value = `read:${row.id}`;
  const { error } = await markEventRead(row.id);
  actionLoading.value = '';
  if (error) {
    window.$message?.error('标记已读失败');
    return;
  }
  row.acknowledged = true;
  notifyEventChanged();
}

function openRenewalModal(row: Api.Central.EventSummary) {
  renewalTarget.value = row;
}

function closeRenewalModal() {
  if (modalSaving.value) return;
  renewalTarget.value = null;
}

async function confirmRenewal(row: Api.Central.EventSummary) {
  const userId = payloadValue(row, 'userId');
  const candidateId = payloadValue(row, 'candidateId') || row.resourceId || '';
  if (!userId || !candidateId) return;
  modalSaving.value = true;
  actionLoading.value = row.id;
  const billingCycle = payloadValue(row, 'billingCycle');
  const suggestedAmount = Number(row.payload?.suggestedAmount || 0);
  const { error } = await confirmUserRenewal(userId, candidateId, {
    billingCycle: billingCycle === 'annual' ? 'annual' : 'monthly',
    amount: suggestedAmount > 0 ? suggestedAmount : undefined
  });
  actionLoading.value = '';
  modalSaving.value = false;
  if (error) {
    window.$message?.error('续费确认失败，请稍后重试');
    return;
  }
  renewalTarget.value = null;
  window.$message?.success('续费已确认并计入财务');
  notifyEventChanged();
  void reloadAll();
}

async function rejectRenewal(row: Api.Central.EventSummary) {
  const userId = payloadValue(row, 'userId');
  const candidateId = payloadValue(row, 'candidateId') || row.resourceId || '';
  if (!userId || !candidateId) return;
  actionLoading.value = row.id;
  const { error } = await rejectUserRenewal(userId, candidateId);
  actionLoading.value = '';
  if (error) {
    window.$message?.error('标记非收费变更失败，请稍后重试');
    return;
  }
  window.$message?.success('已标记为非收费变更');
  notifyEventChanged();
  void reloadAll();
}

function openReplacementModal(row: Api.Central.EventSummary) {
  replacementTarget.value = row;
}

function closeReplacementModal() {
  if (modalSaving.value) return;
  replacementTarget.value = null;
}

async function confirmClientReplacement(row: Api.Central.EventSummary) {
  modalSaving.value = true;
  actionLoading.value = row.id;
  const { error } = await resetUserFromReplacementEvent(row.id);
  actionLoading.value = '';
  modalSaving.value = false;
  if (error) {
    window.$message?.error('客户更换确认失败，可能是 Inbound 已发生变化');
    return;
  }
  replacementTarget.value = null;
  window.$message?.success('已关闭旧用户并创建新用户，旧账务记录已保留');
  notifyEventChanged();
  void reloadAll();
}

async function resolve(row: Api.Central.EventSummary) {
  actionLoading.value = row.id;
  const { error } = await resolveEvent(row.id);
  actionLoading.value = '';
  if (error) {
    window.$message?.error('关闭事件失败');
    return;
  }
  notifyEventChanged();
  void reloadAll();
}

async function retrySync(row: Api.Central.EventSummary) {
  if (!row.nodeId) return;
  actionLoading.value = row.id;
  const { error } = await requestNodeSync(row.nodeId);
  actionLoading.value = '';
  if (error) {
    window.$message?.error('同步请求发送失败');
    return;
  }
  window.$message?.success('同步请求已提交');
}

async function readAll() {
  actionLoading.value = 'all';
  const { error } = await markAllEventsRead();
  actionLoading.value = '';
  if (error) {
    window.$message?.error('批量标记已读失败');
    return;
  }
  rows.value.forEach(row => {
    if (canMarkRead(row)) row.acknowledged = true;
  });
  notifyEventChanged();
  void loadStats();
}

const pageCount = computed(() => Math.max(1, Math.ceil(total.value / filters.page_size)));
const rangeStart = computed(() => (total.value === 0 ? 0 : (filters.page - 1) * filters.page_size + 1));
const rangeEnd = computed(() => Math.min(filters.page * filters.page_size, total.value));
const pendingTabCount = computed(() => pendingCount.value ?? 0);
const allTabCount = computed(() => allCount.value ?? 0);
const renewAmount = computed(() => {
  if (!renewalTarget.value) return 0;
  const amount = Number(renewalTarget.value.payload?.suggestedAmount || 0);
  return Number.isFinite(amount) ? amount : 0;
});
const renewCycleLabel = computed(() =>
  targetPayloadValue(renewalTarget.value, 'billingCycle') === 'annual' ? '年付' : '月付'
);
const renewOldExpiry = computed(() => formatDateTime(targetPayloadValue(renewalTarget.value, 'oldExpiryAt')));
const renewNewExpiry = computed(() => formatDateTime(targetPayloadValue(renewalTarget.value, 'newExpiryAt')));
const replacementClientText = computed(() => {
  if (!replacementTarget.value) return '—';
  const count = payloadNumber(replacementTarget.value, 'clientCount');
  return count ? `${count} 个（已全部更换）` : '已全部更换';
});
const replacementInbound = computed(() => {
  if (!replacementTarget.value) return '—';
  const row = replacementTarget.value;
  const inboundTag = payloadValue(row, 'inboundTag');
  const nodeName = row.nodeName || '—';
  return inboundTag ? `${nodeName} · ${inboundTag}` : nodeName;
});

onMounted(() => {
  void reloadAll();
});
</script>

<template>
  <div class="e-page">
    <!-- ===== 页头 ===== -->
    <header class="e-head">
      <div class="e-head-main">
        <div class="e-eyebrow">
          运营 · 事件中心
          <span class="e-dim">·</span>
          <span>数据更新于 {{ relativeTime(dataAt) }}</span>
        </div>
        <h1 class="e-title">事件中心</h1>
        <p class="e-desc">
          事件由 Agent 提议，确认后才入账：续费与到期变更需要运营确认，确认后生成不可修改的收费记录。事件只反映差异与异常，<b
            class="e-strong"
            >不直接写入 X-Panel / Xray</b
          >。
        </p>
      </div>
      <div class="e-head-actions">
        <DataFreshness :data-at="dataAt" compact />
        <button class="e-btn" type="button" :disabled="loading" @click="reloadAll">
          <SvgIcon icon="mdi:refresh" class="e-btn-ico" />刷新
        </button>
        <button class="e-btn is-soft" type="button" :disabled="actionLoading === 'all'" @click="readAll">
          <SvgIcon icon="mdi:check-all" class="e-btn-ico" />全部标为已读
        </button>
      </div>
    </header>

    <!-- ===== KPI ===== -->
    <section class="e-kpis">
      <div class="e-kpi">
        <div class="e-kpi-top">
          <span class="e-kpi-label">待处理事件</span>
          <span class="e-kpi-icon is-amber"><SvgIcon icon="mdi:bell-outline" /></span>
        </div>
        <div class="e-kpi-value">
          {{ statsReady ? stats.pending : '--' }}<span class="e-kpi-unit">条</span>
        </div>
        <div class="e-kpi-foot">需要运营动作 · 未处理</div>
      </div>
      <div class="e-kpi">
        <div class="e-kpi-top">
          <span class="e-kpi-label">需关注</span>
          <span class="e-kpi-icon is-red"><SvgIcon icon="mdi:alert-outline" /></span>
        </div>
        <div class="e-kpi-value">
          {{ statsReady ? stats.attention : '--' }}<span class="e-kpi-unit">条</span>
        </div>
        <div class="e-kpi-foot">错误 {{ stats.attentionError }} · 客户更换 {{ stats.attentionReplacement }}</div>
      </div>
      <div class="e-kpi">
        <div class="e-kpi-top">
          <span class="e-kpi-label">今日已处理</span>
          <span class="e-kpi-icon is-green"><SvgIcon icon="mdi:check-circle-outline" /></span>
        </div>
        <div class="e-kpi-value">
          {{ statsReady ? stats.resolvedToday : '--' }}<span class="e-kpi-unit">条</span>
        </div>
        <div class="e-kpi-foot">含自动恢复 {{ stats.autoRecoveredToday }} 条</div>
      </div>
      <div class="e-kpi">
        <div class="e-kpi-top">
          <span class="e-kpi-label">平均处理时长</span>
          <span class="e-kpi-icon is-blue"><SvgIcon icon="mdi:clock-outline" /></span>
        </div>
        <div class="e-kpi-value">
          <template v-if="stats.avgHours !== null">
            {{ stats.avgHours.toFixed(1) }}<span class="e-kpi-unit">h</span>
          </template>
          <template v-else>--</template>
        </div>
        <div class="e-kpi-foot">
          近 7 天 · 中位数 {{ stats.medianHours === null ? '--' : `${stats.medianHours.toFixed(1)}h` }} · 已处理
          {{ stats.resolved7d }} 条
        </div>
      </div>
    </section>

    <!-- ===== 事件表 ===== -->
    <section class="e-card">
      <div class="e-tabs">
        <button type="button" :class="{ 'is-active': filters.tab === 'pending' }" @click="selectTab('pending')">
          待处理
          <span class="e-tab-badge is-warn">{{ pendingTabCount }}</span>
        </button>
        <button type="button" :class="{ 'is-active': filters.tab === 'all' }" @click="selectTab('all')">
          全部事件
          <span class="e-tab-badge is-neutral">{{ allTabCount }}</span>
        </button>
      </div>

      <div class="e-toolbar">
        <label class="e-input e-search">
          <SvgIcon icon="mdi:magnify" class="e-input-ico" />
          <input v-model="filters.keyword" placeholder="搜索事件、节点或内容" @keyup.enter="submitFilters" />
        </label>
        <NSelect
          v-model:value="filters.category"
          :options="categoryOptions"
          size="small"
          style="width: 132px"
          @update:value="submitFilters"
        />
        <div class="e-seg">
          <button
            v-for="segment in severitySegments"
            :key="segment.value"
            type="button"
            :class="{ 'is-active': filters.severity === segment.value }"
            @click="selectSeverity(segment.value)"
          >
            {{ segment.label }}
          </button>
        </div>
        <NSelect
          :value="filters.timeRange"
          :options="timeRangeOptions"
          size="small"
          style="width: 118px"
          @update:value="value => changeTimeRange(value as TimeRange)"
        />
        <span class="e-spacer"></span>
        <button class="e-btn e-btn-sm is-primary" type="button" @click="submitFilters">查询</button>
        <button class="e-btn e-btn-sm is-ghost" type="button" @click="resetFilters">重置</button>
      </div>

      <div v-if="errorMessage" class="e-alert">{{ errorMessage }}</div>

      <template v-else>
        <div v-if="loading && rows.length === 0" class="e-loading">正在读取事件…</div>
        <div v-else-if="!rows.length" class="e-empty">
          <div class="e-empty-title">暂无事件</div>
          <p class="e-empty-desc">Agent 上报差异后，事件会出现在这里等待运营确认。</p>
        </div>
        <div v-else class="e-table-wrap">
        <table class="e-table">
          <thead>
            <tr>
              <th>严重度 / 事件</th>
              <th>类型</th>
              <th>节点</th>
              <th>用户</th>
              <th>发生时间</th>
              <th>状态</th>
              <th class="e-th-actions"></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="row in rows" :key="row.id">
              <td>
                <div class="e-event-cell">
                  <i class="e-dot" :class="severityDotClass(row)" :title="severityLabel(row.severity)"></i>
                  <div class="e-event-text">
                    <div class="e-event-title">{{ row.title || eventTypeLabel(row.type) }}</div>
                    <div class="e-sub truncate" :title="row.message">{{ row.message }}</div>
                  </div>
                </div>
              </td>
              <td>
                <span class="e-badge" :class="typeBadgeClass(row)">{{ eventTypeLabel(row.type) }}</span>
              </td>
              <td class="e-mono">{{ row.nodeName || '—' }}</td>
              <td class="e-user">{{ userLabel(row) }}</td>
              <td class="e-num">{{ formatDateTime(row.occurredAt) }}</td>
              <td>
                <span class="e-badge" :class="statusTone(row)"><i class="e-badge-dot"></i>{{ eventStatusLabel(row) }}</span>
              </td>
              <td>
                <div class="e-actions">
                  <template v-if="isRenewal(row)">
                    <button class="e-btn e-btn-xs is-ghost" type="button" @click="openUser(row)">查看用户</button>
                    <button
                      class="e-btn e-btn-xs is-primary"
                      type="button"
                      :disabled="actionLoading === row.id"
                      @click="openRenewalModal(row)"
                    >
                      确认续费
                    </button>
                    <button
                      class="e-btn e-btn-xs"
                      type="button"
                      :disabled="actionLoading === row.id"
                      @click="rejectRenewal(row)"
                    >
                      非收费变更
                    </button>
                  </template>
                  <template v-else-if="isNewUserProfile(row)">
                    <button class="e-btn e-btn-xs is-primary" type="button" @click="completeNewUserProfile(row)">
                      完善资料
                    </button>
                  </template>
                  <template v-else-if="isClientReplacement(row)">
                    <button
                      class="e-btn e-btn-xs is-danger"
                      type="button"
                      :disabled="actionLoading === row.id"
                      @click="openReplacementModal(row)"
                    >
                      确认更换客户
                    </button>
                    <button
                      class="e-btn e-btn-xs is-ghost"
                      type="button"
                      :disabled="actionLoading === row.id"
                      @click="resolve(row)"
                    >
                      取消
                    </button>
                  </template>
                  <template v-else-if="row.status !== 'resolved' && row.actionType === 'retry_sync'">
                    <button
                      class="e-btn e-btn-xs is-primary"
                      type="button"
                      :disabled="actionLoading === row.id"
                      @click="retrySync(row)"
                    >
                      重新同步
                    </button>
                    <button
                      class="e-btn e-btn-xs is-ghost"
                      type="button"
                      :disabled="actionLoading === row.id"
                      @click="resolve(row)"
                    >
                      关闭
                    </button>
                  </template>
                  <template v-else-if="row.status !== 'resolved' && row.actionType">
                    <button
                      class="e-btn e-btn-xs"
                      type="button"
                      :disabled="actionLoading === row.id"
                      @click="resolve(row)"
                    >
                      标记已处理
                    </button>
                  </template>
                  <button
                    v-if="!isRenewal(row) && !isNewUserProfile(row) && !isClientReplacement(row) && hasUserId(row)"
                    class="e-btn e-btn-xs is-ghost"
                    type="button"
                    @click="openUser(row)"
                  >
                    查看用户
                  </button>
                  <button
                    v-if="!row.acknowledged && canMarkRead(row)"
                    class="e-btn e-btn-xs is-soft"
                    type="button"
                    :disabled="actionLoading === `read:${row.id}`"
                    @click="readEvent(row)"
                  >
                    标记已读
                  </button>
                  <span v-if="!hasRowActions(row)" class="e-dim">—</span>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
        </div>
      </template>

      <div class="e-table-foot">
        <span class="e-dim e-foot-hint">
          显示 {{ rangeStart }}–{{ rangeEnd }} 条 · 共 {{ total }} 条 · 每页 {{ filters.page_size }} 条
          <template v-if="filters.timeRange !== 'all'"> · 时间范围在最近 200 条内筛选</template>
        </span>
        <div class="e-pager">
          <NSelect
            :value="filters.page_size"
            :options="pageSizeOptions"
            size="small"
            style="width: 116px"
            @update:value="value => changePageSize(value as number)"
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

    <!-- ===== 分布 + 处理口径 ===== -->
    <div class="e-split">
      <section class="e-card">
        <div class="e-card-head">
          <h3>事件类型分布</h3>
          <span class="e-badge is-neutral">近 7 天 · {{ distributionTotal }} 条</span>
          <div class="e-right">
            <span class="e-dim e-card-hint">
              按事件类型统计<template v-if="statsCapped"> · 最近 200 条</template>
            </span>
          </div>
        </div>
        <div class="e-card-body">
          <div v-if="!distribution.length" class="e-dist-empty">近 7 天暂无事件</div>
          <div
            v-for="(item, index) in distribution"
            :key="item.label"
            class="e-dist-row"
            :class="{ 'has-line': index > 0 }"
          >
            <span class="e-dim">{{ item.label }}</span>
            <div class="e-progress"><span :class="item.tone" :style="{ width: item.width }"></span></div>
            <b class="e-num">{{ item.count }}</b>
          </div>
        </div>
      </section>

      <section class="e-card">
        <div class="e-card-head">
          <h3>处理口径</h3>
          <span class="e-badge is-brand"><i class="e-badge-dot"></i>Agent 只提议不记账</span>
        </div>
        <div class="e-card-body is-flush">
          <div class="e-list-row">
            <span class="e-list-ico"><SvgIcon icon="mdi:currency-cny" /></span>
            <div class="e-list-main">
              <b>续费确认后生成不可修改的收费记录</b>
              <span>金额与服务期此后不可改；确认动作才计入服务期收入与现金收款</span>
            </div>
          </div>
          <div class="e-list-row">
            <span class="e-list-ico"><SvgIcon icon="mdi:account-switch-outline" /></span>
            <div class="e-list-main">
              <b>客户更换会关闭旧用户路径与流量</b>
              <span>仅保留账务记录，并为当前 Inbound 创建新用户；此操作不可恢复</span>
            </div>
          </div>
          <div class="e-list-row">
            <span class="e-list-ico"><SvgIcon icon="mdi:chart-timeline-variant" /></span>
            <div class="e-list-main">
              <b>免费用户与逾期恢复不计入续费率</b>
              <span>续费率分母只取本期到期的收费用户；逾期恢复单独统计</span>
            </div>
          </div>
        </div>
        <div class="e-card-note">已处理事件仍保留在「全部事件」中用于审计；标记已读不会关闭待处理事件。</div>
      </section>
    </div>

    <!-- ===================== 确认续费 ===================== -->
    <Teleport to="body">
      <div v-if="renewalTarget" class="e-mask" @click.self="closeRenewalModal">
        <div class="e-modal" style="width: 560px">
          <div class="e-modal-head">
            <div class="e-modal-head-main">
              <h3>确认续费</h3>
              <p class="e-modal-sub">确认后生成不可修改的收费记录</p>
            </div>
            <span class="e-badge is-warn"><i class="e-badge-dot"></i>计入财务</span>
            <button class="e-icon-btn" type="button" :disabled="modalSaving" @click="closeRenewalModal">
              <SvgIcon icon="mdi:close" />
            </button>
          </div>

          <div class="e-modal-body">
            <div class="e-callout is-warn">
              <SvgIcon icon="mdi:alert-outline" class="e-callout-ico" />
              <div>
                <b>金额与服务期确认后不可修改</b>
                请核对 Agent 提议的续费金额与计费周期；确认动作才会计入服务期收入与现金收款。
              </div>
            </div>

            <div class="e-amount">
              <span class="e-amount-label">建议续费金额</span>
              <strong v-if="renewAmount > 0">{{ formatMoney(renewAmount) }}</strong>
              <strong v-else class="e-dim">未提供建议金额</strong>
              <span class="e-amount-cycle">{{ renewCycleLabel }}</span>
            </div>

            <div class="e-grid-2 mt-16">
              <div class="e-field">
                <span class="e-field-label">用户</span>
                <div class="e-strong">{{ renewalTarget ? userLabel(renewalTarget) : '—' }}</div>
              </div>
              <div class="e-field">
                <span class="e-field-label">节点</span>
                <div class="e-strong e-mono">{{ renewalTarget?.nodeName || '—' }}</div>
              </div>
              <div class="e-field">
                <span class="e-field-label">原到期时间</span>
                <div class="e-strong e-num">{{ renewOldExpiry }}</div>
              </div>
              <div class="e-field">
                <span class="e-field-label">续费后到期时间</span>
                <div class="e-strong e-num">{{ renewNewExpiry }}</div>
              </div>
            </div>

            <div class="e-callout is-info mt-12">
              <SvgIcon icon="mdi:information-outline" class="e-callout-ico" />
              <div>
                <b>也可以标记为非收费变更</b>
                若本次到期变化不涉及收费，请关闭本弹窗并使用「非收费变更」，该事件不会生成收费记录。
              </div>
            </div>
          </div>

          <div class="e-modal-foot">
            <span class="e-dim">确认后事件关闭，收费记录不可撤销</span>
            <div class="e-foot-actions">
              <button class="e-btn" type="button" :disabled="modalSaving" @click="closeRenewalModal">取消</button>
              <button
                class="e-btn is-primary"
                type="button"
                :disabled="modalSaving || !renewalTarget"
                @click="renewalTarget && confirmRenewal(renewalTarget)"
              >
                {{ modalSaving ? '确认中…' : '确认续费并计入财务' }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- ===================== 确认更换客户 ===================== -->
    <Teleport to="body">
      <div v-if="replacementTarget" class="e-mask" @click.self="closeReplacementModal">
        <div class="e-modal" style="width: 560px">
          <div class="e-modal-head">
            <div class="e-modal-head-main">
              <h3>确认更换客户</h3>
              <p class="e-modal-sub">关闭旧用户并创建新用户，旧账务记录保留</p>
            </div>
            <span class="e-badge is-danger"><i class="e-badge-dot"></i>不可恢复</span>
            <button class="e-icon-btn" type="button" :disabled="modalSaving" @click="closeReplacementModal">
              <SvgIcon icon="mdi:close" />
            </button>
          </div>

          <div class="e-modal-body">
            <div class="e-callout is-danger">
              <SvgIcon icon="mdi:alert-outline" class="e-callout-ico" />
              <div>
                <b>此操作不可恢复</b>
                确认后将关闭旧用户的路径和流量，仅保留账务记录，并为当前 Inbound
                创建新用户。请确认 Client 更换是真实发生的业务变更。
              </div>
            </div>

            <div class="e-grid-2 mt-16">
              <div class="e-field">
                <span class="e-field-label">旧用户</span>
                <div class="e-strong">{{ replacementTarget ? userLabel(replacementTarget) : '—' }}</div>
              </div>
              <div class="e-field">
                <span class="e-field-label">线路机 / Inbound</span>
                <div class="e-strong e-mono">{{ replacementInbound }}</div>
              </div>
              <div class="e-field">
                <span class="e-field-label">当前 Client</span>
                <div class="e-strong e-num">{{ replacementClientText }}</div>
              </div>
              <div class="e-field">
                <span class="e-field-label">事件发生</span>
                <div class="e-strong e-num">{{ formatDateTime(replacementTarget?.occurredAt) }}</div>
              </div>
            </div>

            <div class="e-callout is-info mt-12">
              <SvgIcon icon="mdi:information-outline" class="e-callout-ico" />
              <div>
                <b>保留什么</b>
                旧用户的收费订单、历史账单与审计记录全部保留；流量快照、设备凭证与路径配置不再关联旧用户。
              </div>
            </div>
          </div>

          <div class="e-modal-foot">
            <span class="e-dim">确认后事件关闭，并生成一条新用户资料待完善事件</span>
            <div class="e-foot-actions">
              <button class="e-btn" type="button" :disabled="modalSaving" @click="closeReplacementModal">取消</button>
              <button
                class="e-btn is-danger"
                type="button"
                :disabled="modalSaving || !replacementTarget"
                @click="replacementTarget && confirmClientReplacement(replacementTarget)"
              >
                <SvgIcon icon="mdi:alert-outline" class="e-btn-ico" />{{ modalSaving ? '确认中…' : '确认更换' }}
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
   事件中心 · UI v2
   颜色 / 字号 / 间距对齐 design-preview/v2/events.html 与 kit-v2
   ============================================================ */
.e-page {
  --e-card: #141519;
  --e-card-hover: #191b20;
  --e-line: #23262c;
  --e-line-soft: #1c1f24;
  --e-text-1: #f2f3f5;
  --e-text-2: #bcc2cc;
  --e-text-3: #9aa1ad;
  --e-text-4: #868d99;
  --e-brand: #4c6fff;
  --e-ok: #2ed47a;
  --e-warn: #f5b93d;
  --e-danger: #f3656f;
  --e-info: #5aa9f0;
  --e-violet: #a78bfa;
  --e-cyan: #34d3e0;
  --e-radius: 14px;

  display: flex;
  flex-direction: column;
  gap: 15px;
  min-width: 0;
  color: var(--e-text-1);
  font-size: 13px;
}

.e-dim { color: var(--e-text-4); }
.e-strong { color: var(--e-text-1); font-weight: 550; }
.e-mono { font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace; font-size: 12px; }
.e-num { font-variant-numeric: tabular-nums; }
.truncate { min-width: 0; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }

/* ---------- 页头 ---------- */
.e-head { display: flex; align-items: flex-end; gap: 18px; flex-wrap: wrap; }
.e-head-main { min-width: 0; flex: 1; }
.e-eyebrow { display: flex; align-items: center; gap: 6px; margin-bottom: 6px; font-size: 12px; letter-spacing: 0.1em; color: #7f8792; }
.e-title { margin: 0; font-size: 23px; font-weight: 620; line-height: 1.25; color: #f6f7f9; }
.e-desc { max-width: 840px; margin: 7px 0 0; font-size: 12px; line-height: 1.72; color: var(--e-text-3); }
.e-head-actions { display: flex; align-items: center; gap: 8px; flex: none; padding-bottom: 2px; }

/* ---------- 按钮 / 输入 ---------- */
.e-btn {
  display: inline-flex; align-items: center; justify-content: center; gap: 6px;
  height: 34px; padding: 0 13px;
  border: 1px solid #2a2e35; border-radius: 9px;
  background: #1b1e24; color: #e6e8ec;
  font-size: 13px; font-family: inherit; white-space: nowrap; cursor: pointer;
  transition: background 0.15s ease, border-color 0.15s ease;
}
.e-btn:hover { background: #232730; border-color: #363b44; }
.e-btn:disabled { opacity: 0.45; cursor: not-allowed; }
.e-btn-ico { width: 15px; height: 15px; flex: none; }
.e-btn.is-primary { border-color: var(--e-brand); background: var(--e-brand); color: #fff; }
.e-btn.is-primary:hover { background: #5d7cff; border-color: #5d7cff; }
.e-btn.is-soft { border-color: rgba(76, 111, 255, 0.28); background: rgba(76, 111, 255, 0.13); color: #a8b8ff; }
.e-btn.is-soft:hover { background: rgba(76, 111, 255, 0.2); }
.e-btn.is-ghost { border-color: transparent; background: transparent; color: var(--e-text-3); }
.e-btn.is-ghost:hover { background: #1a1d22; color: #e6e8ec; }
.e-btn.is-danger { border-color: rgba(243, 101, 111, 0.3); background: rgba(243, 101, 111, 0.13); color: #ff9ba2; }
.e-btn.is-danger:hover { background: rgba(243, 101, 111, 0.2); }
.e-btn-sm { height: 30px; padding: 0 11px; font-size: 12px; }
.e-btn-xs { height: 28px; padding: 0 10px; font-size: 12px; border-radius: 7px; }

.e-icon-btn {
  display: grid; place-items: center; width: 28px; height: 28px; flex: none;
  border: 1px solid transparent; border-radius: 7px; background: transparent;
  color: var(--e-text-3); cursor: pointer;
}
.e-icon-btn:hover { background: #1a1d22; color: #e6e8ec; }
.e-icon-btn:disabled { opacity: 0.4; cursor: not-allowed; }

.e-input {
  display: flex; align-items: center; gap: 8px; height: 34px; padding: 0 10px;
  border: 1px solid #262a31; border-radius: 9px; background: #0e1013;
}
.e-input:focus-within { border-color: rgba(76, 111, 255, 0.7); }
.e-input input { min-width: 0; flex: 1; border: 0; outline: none; background: transparent; color: var(--e-text-1); font-size: 13px; font-family: inherit; }
.e-input input::placeholder { color: var(--e-text-4); }
.e-input-ico { width: 15px; height: 15px; flex: none; color: var(--e-text-4); }
.e-search { width: 250px; flex: none; }

/* ---------- KPI ---------- */
.e-kpis { display: grid; grid-template-columns: repeat(4, minmax(0, 1fr)); gap: 12px; }
.e-kpi { padding: 15px 16px 14px; border: 1px solid var(--e-line); border-radius: var(--e-radius); background: var(--e-card); min-width: 0; }
.e-kpi-top { display: flex; align-items: center; gap: 8px; margin-bottom: 10px; }
.e-kpi-label { font-size: 12px; color: #b0b6c1; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.e-kpi-icon { display: grid; place-items: center; width: 24px; height: 24px; margin-left: auto; flex: none; border-radius: 8px; border: 1px solid #2a2e35; background: #1d2026; color: #a8adb8; }
.e-kpi-icon :deep(svg) { width: 13px; height: 13px; }
.e-kpi-icon.is-blue { border-color: rgba(76, 111, 255, 0.36); background: rgba(76, 111, 255, 0.16); color: #8fa8ff; }
.e-kpi-icon.is-green { border-color: rgba(46, 212, 122, 0.34); background: rgba(46, 212, 122, 0.14); color: #5ce49b; }
.e-kpi-icon.is-amber { border-color: rgba(245, 185, 61, 0.34); background: rgba(245, 185, 61, 0.14); color: #f7ca6b; }
.e-kpi-icon.is-red { border-color: rgba(243, 101, 111, 0.34); background: rgba(243, 101, 111, 0.14); color: #ff9ba2; }
.e-kpi-value { display: flex; align-items: baseline; gap: 6px; font-size: 25px; font-weight: 620; line-height: 1.15; letter-spacing: -0.015em; font-variant-numeric: tabular-nums; color: #fff; }
.e-kpi-unit { font-size: 13px; font-weight: 500; color: var(--e-text-3); }
.e-kpi-foot { margin-top: 9px; font-size: 12px; line-height: 1.55; color: var(--e-text-4); }

/* ---------- 卡片 ---------- */
.e-card { border: 1px solid var(--e-line); border-radius: var(--e-radius); background: var(--e-card); min-width: 0; }
.e-card-head { display: flex; align-items: center; gap: 10px; padding: 15px 18px; border-bottom: 1px solid var(--e-line-soft); }
.e-card-head h3 { margin: 0; font-size: 14px; font-weight: 600; color: var(--e-text-1); }
.e-card-head .e-right { display: flex; align-items: center; gap: 8px; margin-left: auto; }
.e-card-hint { font-size: 12px; }
.e-card-body { padding: 16px 18px; }
.e-card-body.is-flush { padding: 0; }
.e-card-note { padding: 11px 18px; border-top: 1px solid var(--e-line-soft); font-size: 12px; line-height: 1.6; color: var(--e-text-4); }

/* ---------- Tab ---------- */
.e-tabs { display: flex; align-items: center; gap: 22px; padding: 0 18px; border-bottom: 1px solid var(--e-line-soft); }
.e-tabs button {
  position: relative; display: flex; align-items: center; gap: 6px; height: 42px; padding: 0;
  border: 0; background: none; color: var(--e-text-3); font-size: 13px; font-family: inherit; cursor: pointer;
  transition: color 0.14s ease;
}
.e-tabs button:hover { color: var(--e-text-2); }
.e-tabs button.is-active { color: var(--e-text-1); font-weight: 550; }
.e-tabs button.is-active::after {
  content: ''; position: absolute; left: 0; right: 0; bottom: -1px; height: 2px;
  border-radius: 2px 2px 0 0; background: var(--e-brand);
}
.e-tab-badge { display: inline-flex; align-items: center; height: 20px; padding: 0 7px; border: 1px solid transparent; border-radius: 99px; font-size: 12px; font-weight: 550; font-variant-numeric: tabular-nums; }
.e-tab-badge.is-warn { border-color: rgba(245, 185, 61, 0.32); background: rgba(245, 185, 61, 0.15); color: #f7ca6b; }
.e-tab-badge.is-neutral { border-color: #2a2e35; background: #1d2026; color: var(--e-text-2); }

/* ---------- 工具条 ---------- */
.e-toolbar { display: flex; align-items: center; gap: 8px; flex-wrap: wrap; padding: 12px 16px; border-bottom: 1px solid var(--e-line-soft); }
.e-spacer { flex: 1; }
.e-seg { display: inline-flex; gap: 2px; padding: 3px; border: 1px solid #24272e; border-radius: 9px; background: #14161a; }
.e-seg button { height: 26px; padding: 0 11px; border: 0; border-radius: 7px; background: transparent; color: var(--e-text-3); font-size: 12px; font-family: inherit; cursor: pointer; }
.e-seg button:hover { color: var(--e-text-1); }
.e-seg button.is-active { background: #2b3038; color: #fff; }

/* ---------- 表格 ---------- */
.e-table-wrap { overflow-x: auto; }
.e-table { width: 100%; min-width: 1120px; border-collapse: separate; border-spacing: 0; font-size: 13px; }
.e-table thead th {
  padding: 9px 14px; border-bottom: 1px solid var(--e-line); background: var(--e-card);
  font-size: 12px; font-weight: 600; letter-spacing: 0.06em; color: #a2a9b4; text-align: left; white-space: nowrap;
}
.e-table thead th.e-th-actions { width: 240px; }
.e-table tbody td { padding: 11px 14px; border-bottom: 1px solid var(--e-line-soft); color: var(--e-text-2); vertical-align: middle; white-space: nowrap; }
.e-table tbody tr:last-child td { border-bottom: none; }
.e-table tbody tr { transition: background 0.12s ease; }
.e-table tbody tr:hover td { background: var(--e-card-hover); }
.e-table td .e-sub { display: block; margin-top: 2px; max-width: 360px; font-size: 12px; line-height: 1.45; color: var(--e-text-4); }
.e-event-cell { display: flex; align-items: flex-start; gap: 8px; min-width: 300px; }
.e-event-cell .e-dot { margin-top: 6px; }
.e-event-text { min-width: 0; max-width: 360px; }
.e-event-title { color: var(--e-text-1); font-weight: 550; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.e-user { max-width: 190px; overflow: hidden; text-overflow: ellipsis; }

/* ---------- 状态点 / 徽标 ---------- */
.e-dot { width: 7px; height: 7px; flex: none; border-radius: 50%; background: #545b64; }
.e-dot.is-ok { background: var(--e-ok); box-shadow: 0 0 7px rgba(46, 212, 122, 0.5); }
.e-dot.is-warn { background: var(--e-warn); box-shadow: 0 0 7px rgba(245, 185, 61, 0.45); }
.e-dot.is-danger { background: var(--e-danger); box-shadow: 0 0 7px rgba(243, 101, 111, 0.45); }
.e-dot.is-info { background: var(--e-info); box-shadow: 0 0 7px rgba(90, 169, 240, 0.45); }
.e-dot.is-neutral { background: #545b64; }

.e-badge { display: inline-flex; align-items: center; gap: 5px; height: 22px; padding: 0 8px; border: 1px solid transparent; border-radius: 6px; font-size: 12px; font-weight: 550; white-space: nowrap; }
.e-badge .e-badge-dot { width: 6px; height: 6px; flex: none; border-radius: 50%; background: currentColor; }
.e-badge.is-ok { border-color: rgba(46, 212, 122, 0.32); background: rgba(46, 212, 122, 0.15); color: #5ce49b; }
.e-badge.is-warn { border-color: rgba(245, 185, 61, 0.32); background: rgba(245, 185, 61, 0.15); color: #f7ca6b; }
.e-badge.is-danger { border-color: rgba(243, 101, 111, 0.32); background: rgba(243, 101, 111, 0.15); color: #ff9ba2; }
.e-badge.is-info { border-color: rgba(90, 169, 240, 0.32); background: rgba(90, 169, 240, 0.15); color: #8cc4f7; }
.e-badge.is-brand { border-color: rgba(76, 111, 255, 0.34); background: rgba(76, 111, 255, 0.16); color: #9db2ff; }
.e-badge.is-violet { border-color: rgba(167, 139, 250, 0.34); background: rgba(167, 139, 250, 0.16); color: #c0a9ff; }
.e-badge.is-neutral { border-color: #2a2e35; background: #1d2026; color: var(--e-text-3); }

/* ---------- 行操作 ---------- */
.e-actions { display: flex; align-items: center; gap: 4px; flex-wrap: wrap; justify-content: flex-end; }
.e-table thead th:last-child { text-align: right; }

/* ---------- 表尾 / 分页 ---------- */
.e-table-foot { display: flex; align-items: center; gap: 12px; flex-wrap: wrap; padding: 11px 16px; border-top: 1px solid var(--e-line-soft); }
.e-foot-hint { font-size: 12px; }
.e-pager { display: flex; align-items: center; gap: 8px; margin-left: auto; }

/* ---------- 分布 + 处理口径 ---------- */
.e-split { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 12px; }
.e-dist-row { display: grid; grid-template-columns: 78px minmax(0, 1fr) 34px; align-items: center; gap: 10px; padding: 8px 0; font-size: 13px; }
.e-dist-row.has-line { border-top: 1px solid var(--e-line-soft); }
.e-dist-row b { text-align: right; font-weight: 600; color: var(--e-text-1); }
.e-dist-empty { padding: 6px 0; font-size: 12px; color: var(--e-text-4); }
.e-progress { position: relative; height: 5px; border-radius: 99px; background: rgba(163, 177, 255, 0.1); overflow: hidden; }
.e-progress > span { position: absolute; inset: 0 auto 0 0; border-radius: 99px; background: linear-gradient(90deg, #5961c8, #7b83ee); }
.e-progress > span.is-warn { background: linear-gradient(90deg, #b07f22, #e8b04b); }
.e-progress > span.is-danger { background: linear-gradient(90deg, #b8434c, #f0656f); }
.e-progress > span.is-ok { background: linear-gradient(90deg, #26906b, #3ecf9a); }
.e-progress > span.is-cyan { background: linear-gradient(90deg, #16828c, #1dc2d0); }
.e-progress > span.is-violet { background: linear-gradient(90deg, #7c5cf0, #a78bfa); }

.e-list-row { display: flex; align-items: center; gap: 12px; padding: 11px 18px; border-bottom: 1px solid var(--e-line-soft); }
.e-list-row:last-child { border-bottom: none; }
.e-list-ico { display: grid; place-items: center; width: 30px; height: 30px; flex: none; border-radius: 9px; border: 1px solid var(--e-line-soft); background: rgba(163, 177, 255, 0.07); color: var(--e-text-2); }
.e-list-ico :deep(svg) { width: 15px; height: 15px; }
.e-list-main { min-width: 0; flex: 1; }
.e-list-main b { display: block; font-size: 13px; font-weight: 550; color: var(--e-text-1); }
.e-list-main span { display: block; margin-top: 1px; font-size: 12px; line-height: 1.55; color: var(--e-text-4); }

/* ---------- 提示 / 空 / 加载 ---------- */
.e-alert { margin: 14px 16px; padding: 10px 13px; border: 1px solid rgba(243, 101, 111, 0.32); border-radius: 10px; background: rgba(243, 101, 111, 0.1); color: #ffb3b9; font-size: 12px; }
.e-loading, .e-empty { display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 10px; min-height: 220px; padding: 24px; color: var(--e-text-4); font-size: 12px; }
.e-empty-title { font-size: 14px; color: var(--e-text-1); }
.e-empty-desc { margin: 0; color: var(--e-text-4); }

/* ---------- 弹窗 ---------- */
.e-mask { position: fixed; inset: 0; z-index: 2000; display: grid; place-items: center; padding: 24px; background: rgba(0, 0, 0, 0.66); }
.e-modal { display: flex; flex-direction: column; max-height: min(86vh, 880px); border: 1px solid var(--e-line); border-radius: 16px; background: var(--e-card); color: var(--e-text-1); text-align: left; box-shadow: 0 30px 80px -24px rgba(0, 0, 0, 0.9); }
.e-modal-head { display: flex; align-items: flex-start; gap: 10px; padding: 16px 20px 14px; border-bottom: 1px solid var(--e-line-soft); }
.e-modal-head-main { flex: 1; min-width: 0; }
.e-modal-head h3 { margin: 0; font-size: 15px; font-weight: 600; }
.e-modal-sub { margin: 3px 0 0; font-size: 12px; color: var(--e-text-4); }
.e-modal-body { flex: 1; min-height: 0; overflow-y: auto; padding: 16px 20px 20px; }
.e-modal-foot { display: flex; align-items: center; gap: 10px; padding: 13px 20px; border-top: 1px solid var(--e-line-soft); font-size: 12px; }
.e-foot-actions { display: flex; align-items: center; gap: 8px; margin-left: auto; }

.e-callout { display: flex; gap: 10px; padding: 12px 13px; border: 1px solid var(--e-line); border-radius: 11px; background: rgba(255, 255, 255, 0.022); font-size: 12px; line-height: 1.7; color: var(--e-text-2); }
.e-callout-ico { width: 16px; height: 16px; flex: none; margin-top: 2px; color: var(--e-text-3); }
.e-callout b { display: block; margin-bottom: 1px; font-size: 12.5px; font-weight: 600; color: var(--e-text-1); }
.e-callout.is-info { border-color: rgba(90, 169, 240, 0.24); background: rgba(90, 169, 240, 0.08); }
.e-callout.is-info .e-callout-ico { color: #8cc4f7; }
.e-callout.is-warn { border-color: rgba(245, 185, 61, 0.28); background: rgba(245, 185, 61, 0.08); }
.e-callout.is-warn .e-callout-ico { color: #f7ca6b; }
.e-callout.is-danger { border-color: rgba(243, 101, 111, 0.28); background: rgba(243, 101, 111, 0.08); }
.e-callout.is-danger .e-callout-ico { color: #ff9ba2; }

.e-amount { display: flex; align-items: baseline; gap: 10px; margin-top: 14px; padding: 13px 15px; border: 1px solid var(--e-line); border-radius: 11px; background: #101114; }
.e-amount-label { font-size: 12px; color: var(--e-text-3); }
.e-amount strong { font-size: 22px; font-weight: 620; letter-spacing: -0.01em; font-variant-numeric: tabular-nums; color: #fff; }
.e-amount-cycle { margin-left: auto; padding: 2px 8px; border: 1px solid #2a2e35; border-radius: 6px; background: #1d2026; font-size: 12px; color: var(--e-text-2); }

.e-grid-2 { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 12px; }
.e-field { display: flex; flex-direction: column; gap: 5px; min-width: 0; }
.e-field-label { font-size: 12px; color: #b0b6c1; }

.mt-12 { margin-top: 12px; }
.mt-16 { margin-top: 16px; }

/* ---------- 响应式 ---------- */
@media (max-width: 1280px) {
  .e-kpis { grid-template-columns: repeat(2, minmax(0, 1fr)); }
}
@media (max-width: 1100px) {
  .e-split { grid-template-columns: minmax(0, 1fr); }
}
@media (max-width: 820px) {
  .e-search { width: 100%; }
  .e-grid-2 { grid-template-columns: minmax(0, 1fr); }
}
</style>
