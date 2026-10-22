<script setup lang="ts">
import { computed, nextTick, onMounted, reactive, ref, watch } from 'vue';
import {
  confirmUserRenewal,
  createNodeCost,
  fetchEvents,
  fetchFinanceSummary,
  fetchNodes
} from '@/service/api';
import DataFreshness from '@/components/project/data-freshness.vue';
import { useEcharts, type ECOption } from '@/hooks/common/echarts';

defineOptions({ name: 'FinanceManagement' });

type FinanceOrder = NonNullable<Api.Central.FinanceSummary['orders']>[number];
type FinanceCost = NonNullable<Api.Central.FinanceSummary['costDetails']>[number];

interface MonthPoint {
  period: string;
  label: string;
  income: number;
  cost: number;
  renewalRate: number;
  dueUsers: number;
  renewedUsers: number;
}

interface RenewalCandidateRow {
  eventId: string;
  candidateId: string;
  userId: string;
  userName: string;
  oldExpiryAt: string;
  newExpiryAt: string;
  billingCycle: Api.Central.BillingCycle;
  amount: number;
  currency: string;
  detectedAt: string;
}

type TabKey = 'overview' | 'orders' | 'renewal' | 'cost';

const tabs: Array<{ label: string; value: TabKey }> = [
  { label: '总览', value: 'overview' },
  { label: '收款订单', value: 'orders' },
  { label: '续费分析', value: 'renewal' },
  { label: '成本明细', value: 'cost' }
];

const loading = ref(false);
const errorMessage = ref('');
const summary = ref<Api.Central.FinanceSummary | null>(null);
const period = ref<string>(new Date().toISOString().slice(0, 7));
const activeTab = ref<TabKey>('overview');

const trend = ref<MonthPoint[]>([]);
const trendLoading = ref(false);
const trendError = ref('');
const trendCache = new Map<string, MonthPoint>();
let trendRequestId = 0;

const candidates = ref<RenewalCandidateRow[]>([]);
const candidateLoading = ref(false);
const candidateError = ref('');
const confirmingCandidateId = ref('');

const orderFilters = reactive({ source: 'all', type: 'all', keyword: '' });
const orderPage = ref(1);
const orderPageSize = 20;

const costModalVisible = ref(false);
const costSaving = ref(false);
const costError = ref('');
const nodesLoading = ref(false);
const nodeOptions = ref<Array<{ label: string; value: string }>>([]);
const costForm = reactive({
  type: 'node',
  nodeId: null as string | null,
  amount: '',
  effectiveFrom: null as string | null,
  notes: ''
});

const sourceOptions = [
  { label: '全部来源', value: 'all' },
  { label: '手工确认', value: 'manual' },
  { label: 'Agent 候选', value: 'agent_detected' }
];
const orderTypeOptions = [
  { label: '全部类型', value: 'all' },
  { label: '首购', value: 'initial' },
  { label: '续费', value: 'renewal' },
  { label: '逾期恢复', value: 'recovery' }
];
const costTypeOptions = [{ label: '节点成本', value: 'node' }];

const moneyFormatters = new Map<string, Intl.NumberFormat>();

function formatMoney(value?: number | null, currency = 'CNY') {
  if (value === undefined || value === null || !Number.isFinite(value)) return '--';
  let formatter = moneyFormatters.get(currency);
  if (!formatter) {
    formatter = new Intl.NumberFormat('zh-CN', {
      style: 'currency',
      currency,
      minimumFractionDigits: 2,
      maximumFractionDigits: 2
    });
    moneyFormatters.set(currency, formatter);
  }
  return formatter.format(value);
}

function pad2(value: number) {
  return `${value}`.padStart(2, '0');
}

function formatDateTime(value?: string | null) {
  if (!value) return '--';
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return value;
  return `${date.getFullYear()}-${pad2(date.getMonth() + 1)}-${pad2(date.getDate())} ${pad2(date.getHours())}:${pad2(
    date.getMinutes()
  )}`;
}

function formatDay(value?: string | null) {
  if (!value) return '--';
  if (/^\d{4}-\d{2}-\d{2}/.test(value)) return value.slice(0, 10);
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return value;
  return `${date.getFullYear()}-${pad2(date.getMonth() + 1)}-${pad2(date.getDate())}`;
}

function formatPercentValue(value: number) {
  return `${value.toFixed(1)}%`;
}

function formatAxisMoney(value: number) {
  if (Math.abs(value) >= 1000) return `${(value / 1000).toFixed(1)}k`;
  return `${value}`;
}

function orderTypeLabel(value: string) {
  if (value === 'renewal') return '续费';
  if (value === 'recovery') return '逾期恢复';
  return '首购';
}

function orderTypeBadge(value: string) {
  if (value === 'renewal') return 'is-brand';
  if (value === 'recovery') return 'is-warn';
  return 'is-ok';
}

function cycleLabel(value: Api.Central.BillingCycle) {
  return value === 'annual' ? '年付' : '月付';
}

function costTypeLabel(value: string) {
  if (value === 'node') return '节点';
  if (value === 'ip') return '出口 IP';
  return '其他';
}

function categoryLabel(value: string) {
  if (!value) return '--';
  if (value === 'exit_ip') return '出口 IP';
  if (value === 'other') return '其他';
  return value;
}

function sourceLabel(value: string) {
  return value === 'agent_detected' ? 'Agent 候选' : '手工确认';
}

function sourceBadge(value: string) {
  return value === 'agent_detected' ? 'is-info' : 'is-neutral';
}

function payloadText(row: Api.Central.EventSummary, key: string) {
  const value = row.payload?.[key];
  if (typeof value === 'string') return value;
  if (value === undefined || value === null) return '';
  return String(value);
}

function currentMonth() {
  return new Date().toISOString().slice(0, 7);
}

function monthLabel(value: string) {
  const month = Number(value.slice(5, 7));
  return Number.isFinite(month) ? `${month}月` : value;
}

function monthRange(endPeriod: string, count: number) {
  const year = Number(endPeriod.slice(0, 4));
  const month = Number(endPeriod.slice(5, 7));
  const periods: string[] = [];
  for (let index = count - 1; index >= 0; index -= 1) {
    const date = new Date(Date.UTC(year, month - 1 - index, 1));
    periods.push(`${date.getUTCFullYear()}-${pad2(date.getUTCMonth() + 1)}`);
  }
  return periods;
}

const orders = computed<FinanceOrder[]>(() => summary.value?.orders ?? []);
const costs = computed<FinanceCost[]>(() => summary.value?.costDetails ?? []);
const renewal = computed(
  () =>
    summary.value?.renewal || {
      dueUsers: 0,
      renewedUsers: 0,
      notRenewedUsers: 0,
      renewalRate: 0,
      recoveryUsers: 0
    }
);

const nodeCosts = computed(() => costs.value.filter(item => item.type === 'node'));
const otherCosts = computed(() => costs.value.filter(item => item.type === 'other'));

const ordersTotal = computed(() => orders.value.reduce((sum, item) => sum + item.amount, 0));
const nodeCostTotal = computed(() => nodeCosts.value.reduce((sum, item) => sum + item.amount, 0));
const otherCostTotal = computed(() => otherCosts.value.reduce((sum, item) => sum + item.amount, 0));
const ipCostTotal = computed(() =>
  costs.value.filter(item => item.type === 'ip').reduce((sum, item) => sum + item.amount, 0)
);

const grossMargin = computed(() => {
  const income = summary.value?.monthIncome;
  if (income === undefined || income <= 0) return null;
  return (summary.value?.grossProfit || 0) / income;
});

const grossMarginText = computed(() => (grossMargin.value === null ? '--' : formatPercentValue(grossMargin.value * 100)));

const costGroup = computed(() => [
  { key: 'node', label: '节点成本', color: '#4c6fff', amount: nodeCostTotal.value },
  { key: 'ip', label: '流量成本（出口 IP）', color: '#f5b93d', amount: ipCostTotal.value },
  { key: 'other', label: '其它成本', color: '#a78bfa', amount: otherCostTotal.value }
]);

const costTotal = computed(() => summary.value?.monthCost ?? 0);

function costPercent(amount: number) {
  if (!costTotal.value) return 0;
  return Number(((amount / costTotal.value) * 100).toFixed(1));
}

const largestCost = computed(() => {
  if (!costs.value.length) return null;
  return costs.value.reduce((current, item) => (item.amount > current.amount ? item : current), costs.value[0]);
});

const filteredOrders = computed(() => {
  const keyword = orderFilters.keyword.trim().toLowerCase();
  return orders.value.filter(item => {
    if (orderFilters.source !== 'all' && item.source !== orderFilters.source) return false;
    if (orderFilters.type !== 'all' && item.orderType !== orderFilters.type) return false;
    if (!keyword) return true;
    return (
      (item.userName || '').toLowerCase().includes(keyword) ||
      (item.orderNo || '').toLowerCase().includes(keyword) ||
      (item.inboundTag || '').toLowerCase().includes(keyword)
    );
  });
});

const filteredOrdersTotal = computed(() => filteredOrders.value.reduce((sum, item) => sum + item.amount, 0));
const orderPageCount = computed(() => Math.max(1, Math.ceil(filteredOrders.value.length / orderPageSize)));
const pagedOrders = computed(() =>
  filteredOrders.value.slice((orderPage.value - 1) * orderPageSize, orderPage.value * orderPageSize)
);
const pageStart = computed(() =>
  filteredOrders.value.length ? (orderPage.value - 1) * orderPageSize + 1 : 0
);
const pageEnd = computed(() => Math.min(orderPage.value * orderPageSize, filteredOrders.value.length));

function typeTotal(type: string) {
  return filteredOrders.value.filter(item => item.orderType === type).reduce((sum, item) => sum + item.amount, 0);
}

const candidateTotal = computed(() => candidates.value.reduce((sum, item) => sum + item.amount, 0));

const trendRangeLabel = computed(() => {
  if (!trend.value.length) return '近 12 个月';
  return `${trend.value[0].period} → ${trend.value[trend.value.length - 1].period}`;
});

const renewalTrendPoints = computed(() => trend.value.slice(-6));

function expiryDays(oldValue: string, newValue: string) {
  const from = new Date(oldValue).getTime();
  const to = new Date(newValue).getTime();
  if (!Number.isFinite(from) || !Number.isFinite(to)) return '--';
  const days = Math.round((to - from) / 86_400_000);
  return days >= 0 ? `+${days} 天` : `${days} 天`;
}

let financeRequestId = 0;

async function loadFinance() {
  const requestId = ++financeRequestId;
  loading.value = true;
  errorMessage.value = '';
  const { data, error } = await fetchFinanceSummary({ period: period.value || undefined });
  if (requestId !== financeRequestId) return;
  if (error || !data) {
    errorMessage.value = '中央后端暂不可用，无法读取财务数据';
    summary.value = null;
  } else {
    summary.value = data;
  }
  loading.value = false;
}

async function loadTrend(force = false) {
  const requestId = ++trendRequestId;
  const periods = monthRange(period.value || currentMonth(), 12);
  if (force) {
    periods.forEach(item => trendCache.delete(item));
  }
  const missing = periods.filter(item => !trendCache.has(item));
  trendLoading.value = true;
  trendError.value = '';
  if (missing.length) {
    const results = await Promise.all(
      missing.map(async item => {
        const { data, error } = await fetchFinanceSummary({ period: item });
        if (error || !data) return null;
        const point: MonthPoint = {
          period: item,
          label: monthLabel(item),
          income: data.monthIncome ?? 0,
          cost: data.monthCost ?? 0,
          renewalRate: data.renewal?.renewalRate ?? 0,
          dueUsers: data.renewal?.dueUsers ?? 0,
          renewedUsers: data.renewal?.renewedUsers ?? 0
        };
        return point;
      })
    );
    if (requestId !== trendRequestId) return;
    results.forEach(item => {
      if (item) trendCache.set(item.period, item);
    });
    if (results.every(item => item === null)) {
      trendError.value = '近 12 个月财务趋势暂不可用';
    }
  }
  if (requestId !== trendRequestId) return;
  trend.value = periods
    .map(item => trendCache.get(item))
    .filter((item): item is MonthPoint => Boolean(item));
  trendLoading.value = false;
}

async function loadRenewalCandidates() {
  candidateLoading.value = true;
  candidateError.value = '';
  const { data, error } = await fetchEvents({ page: 1, page_size: 100, status: 'pending' });
  if (error || !data) {
    candidateError.value = '待确认续费候选读取失败';
    candidates.value = [];
  } else {
    candidates.value = data.items
      .filter(
        item =>
          item.type === 'renewal_candidate_detected' &&
          item.status !== 'resolved' &&
          item.status !== 'dismissed'
      )
      .map(item => ({
        eventId: item.id,
        candidateId: payloadText(item, 'candidateId') || item.resourceId || '',
        userId: payloadText(item, 'userId'),
        userName: payloadText(item, 'userName') || '未知用户',
        oldExpiryAt: payloadText(item, 'oldExpiryAt'),
        newExpiryAt: payloadText(item, 'newExpiryAt'),
        billingCycle: (payloadText(item, 'billingCycle') === 'annual'
          ? 'annual'
          : 'monthly') as Api.Central.BillingCycle,
        amount: Number(item.payload?.suggestedAmount ?? 0),
        currency: payloadText(item, 'currency') || 'CNY',
        detectedAt: item.occurredAt
      }))
      .filter(item => Boolean(item.userId && item.candidateId));
  }
  candidateLoading.value = false;
}

async function confirmCandidate(row: RenewalCandidateRow) {
  if (confirmingCandidateId.value) return;
  confirmingCandidateId.value = row.eventId;
  const { error } = await confirmUserRenewal(row.userId, row.candidateId, {
    billingCycle: row.billingCycle,
    amount: row.amount > 0 ? row.amount : undefined
  });
  confirmingCandidateId.value = '';
  if (error) {
    window.$message?.error('续费确认失败，请稍后重试');
    return;
  }
  window.$message?.success('续费已确认并计入财务');
  trendCache.clear();
  void loadFinance();
  void loadTrend();
  void loadRenewalCandidates();
}

async function loadNodeOptions() {
  nodesLoading.value = true;
  const { data, error } = await fetchNodes({ page: 1, page_size: 200 });
  if (error || !data) {
    nodeOptions.value = [];
    costError.value = '节点列表读取失败，无法新增成本记录';
  } else {
    nodeOptions.value = data.items.map(item => ({
      label: `${item.name}${item.region ? ` · ${item.region}` : ''}`,
      value: item.id
    }));
  }
  nodesLoading.value = false;
}

function openCostModal() {
  costError.value = '';
  costForm.type = 'node';
  costForm.nodeId = null;
  costForm.amount = '';
  costForm.effectiveFrom = `${period.value || currentMonth()}-01`;
  costForm.notes = '';
  costModalVisible.value = true;
  if (!nodeOptions.value.length) void loadNodeOptions();
}

function closeCostModal() {
  if (costSaving.value) return;
  costModalVisible.value = false;
}

async function submitCost() {
  costError.value = '';
  if (!costForm.nodeId) {
    costError.value = '请选择关联节点';
    return;
  }
  const amount = Number(costForm.amount);
  if (!costForm.amount.trim() || !Number.isFinite(amount) || amount <= 0) {
    costError.value = '金额必须为大于 0 的数字';
    return;
  }
  if (!costForm.effectiveFrom) {
    costError.value = '请选择生效日期';
    return;
  }
  costSaving.value = true;
  const { error } = await createNodeCost(costForm.nodeId, {
    category: '服务器',
    monthlyAmount: Number(amount.toFixed(2)),
    currency: 'CNY',
    effectiveFrom: costForm.effectiveFrom,
    notes: costForm.notes.trim() || null
  });
  costSaving.value = false;
  if (error) {
    costError.value = '新增成本记录失败，请检查参数后重试';
    return;
  }
  window.$message?.success('成本记录已新增，立即计入本月成本');
  costModalVisible.value = false;
  trendCache.clear();
  void loadFinance();
  void loadTrend();
}

function csvCell(value: string | number) {
  return `"${String(value).replace(/"/g, '""')}"`;
}

function exportCsv() {
  let rows: Array<Array<string | number>> = [];
  let filename = `finance-${period.value}-orders.csv`;
  if (activeTab.value === 'cost') {
    rows = [
      ['类型', '成本项目', '分类', '本期金额', '生效日期', '失效日期', '备注'],
      ...costs.value.map(item => [
        costTypeLabel(item.type),
        item.name,
        categoryLabel(item.category),
        item.amount.toFixed(2),
        formatDay(item.effectiveFrom),
        formatDay(item.effectiveTo),
        item.notes || ''
      ])
    ];
    filename = `finance-${period.value}-costs.csv`;
  } else if (activeTab.value === 'renewal') {
    rows = [
      ['用户', '原到期', '新到期', '增加', '金额', '检测时间'],
      ...candidates.value.map(item => [
        item.userName,
        formatDateTime(item.oldExpiryAt),
        formatDateTime(item.newExpiryAt),
        expiryDays(item.oldExpiryAt, item.newExpiryAt),
        item.amount.toFixed(2),
        formatDateTime(item.detectedAt)
      ])
    ];
    filename = `finance-${period.value}-renewals.csv`;
  } else {
    const source = activeTab.value === 'orders' ? filteredOrders.value : orders.value;
    rows = [
      ['用户', 'Inbound', '节点', '订单号', '类型', '周期', '金额', '服务期起', '服务期止', '收款时间', '来源'],
      ...source.map(item => [
        item.userName,
        item.inboundTag || '',
        item.nodeName || '',
        item.orderNo || item.id.slice(0, 12),
        orderTypeLabel(item.orderType),
        cycleLabel(item.billingCycle),
        item.amount.toFixed(2),
        formatDay(item.serviceFrom),
        formatDay(item.serviceTo),
        formatDateTime(item.paidAt),
        sourceLabel(item.source)
      ])
    ];
    filename = `finance-${period.value}-orders.csv`;
  }
  const content = rows.map(row => row.map(cell => csvCell(cell)).join(',')).join('\r\n');
  const blob = new Blob([`\uFEFF${content}`], { type: 'text/csv;charset=utf-8;' });
  const url = URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = filename;
  link.click();
  URL.revokeObjectURL(url);
  window.$message?.success('已导出当前视图 CSV');
}

function onPeriodChange(value: string | null) {
  if (!value) return;
  period.value = value;
  orderPage.value = 1;
  void loadFinance();
  void loadTrend();
}

async function refreshAll() {
  await loadFinance();
  void loadTrend(true);
  void loadRenewalCandidates();
}

function buildRevenueOptions(): ECOption {
  const points = trend.value;
  return {
    animation: false,
    grid: { left: 58, right: 18, top: 18, bottom: 28 },
    tooltip: {
      trigger: 'axis',
      backgroundColor: '#191b20',
      borderColor: '#2e323a',
      textStyle: { color: '#f2f3f5', fontSize: 12 },
      formatter: (params: any) => {
        const list = Array.isArray(params) ? params : [params];
        const index = list[0]?.dataIndex ?? 0;
        const point = points[index];
        if (!point) return '';
        return `<div>${point.period}</div><div>服务期收入：${formatMoney(point.income)}</div><div>当月成本：${formatMoney(
          point.cost
        )}</div><div>毛利：${formatMoney(point.income - point.cost)}</div>`;
      }
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: points.map(point => point.label),
      axisTick: { show: false },
      axisLine: { lineStyle: { color: '#23262c' } },
      axisLabel: { hideOverlap: true, fontSize: 11, color: '#8f96a1' }
    },
    yAxis: {
      type: 'value',
      axisLabel: { formatter: (value: number) => formatAxisMoney(value), fontSize: 11, color: '#8f96a1' },
      splitLine: { lineStyle: { type: 'dashed', opacity: 0.35, color: '#23262c' } }
    },
    series: [
      {
        name: '服务期收入',
        type: 'line',
        smooth: 0.25,
        showSymbol: false,
        data: points.map(point => point.income),
        lineStyle: { width: 2.4, color: '#4c6fff' },
        itemStyle: { color: '#4c6fff' },
        areaStyle: {
          color: {
            type: 'linear',
            x: 0,
            y: 0,
            x2: 0,
            y2: 1,
            colorStops: [
              { offset: 0, color: 'rgba(76, 111, 255, 0.26)' },
              { offset: 1, color: 'rgba(76, 111, 255, 0)' }
            ]
          }
        }
      },
      {
        name: '当月成本',
        type: 'line',
        smooth: 0.25,
        showSymbol: false,
        data: points.map(point => point.cost),
        lineStyle: { width: 2, color: '#f5b93d' },
        itemStyle: { color: '#f5b93d' },
        areaStyle: {
          color: {
            type: 'linear',
            x: 0,
            y: 0,
            x2: 0,
            y2: 1,
            colorStops: [
              { offset: 0, color: 'rgba(245, 185, 61, 0.18)' },
              { offset: 1, color: 'rgba(245, 185, 61, 0)' }
            ]
          }
        }
      }
    ]
  };
}

function buildRenewalOptions(): ECOption {
  const points = renewalTrendPoints.value;
  return {
    animation: false,
    grid: { left: 44, right: 16, top: 16, bottom: 28 },
    tooltip: {
      trigger: 'axis',
      backgroundColor: '#191b20',
      borderColor: '#2e323a',
      textStyle: { color: '#f2f3f5', fontSize: 12 },
      formatter: (params: any) => {
        const list = Array.isArray(params) ? params : [params];
        const index = list[0]?.dataIndex ?? 0;
        const point = points[index];
        if (!point) return '';
        return `<div>${point.period}</div><div>严格续费率：${formatPercentValue(
          point.renewalRate * 100
        )}</div><div>到期 ${point.dueUsers} 户 · 已续费 ${point.renewedUsers} 户</div>`;
      }
    },
    xAxis: {
      type: 'category',
      data: points.map(point => point.label),
      axisTick: { show: false },
      axisLine: { lineStyle: { color: '#23262c' } },
      axisLabel: { fontSize: 11, color: '#8f96a1' }
    },
    yAxis: {
      type: 'value',
      max: 100,
      axisLabel: { formatter: (value: number) => `${value}%`, fontSize: 11, color: '#8f96a1' },
      splitLine: { lineStyle: { type: 'dashed', opacity: 0.35, color: '#23262c' } }
    },
    series: [
      {
        name: '严格续费率',
        type: 'bar',
        data: points.map(point => Number((point.renewalRate * 100).toFixed(1))),
        barWidth: 18,
        itemStyle: { color: '#4c6fff', borderRadius: [5, 5, 0, 0] }
      }
    ]
  };
}

const { domRef: revenueChartRef, updateOptions: updateRevenueChart } = useEcharts(() => buildRevenueOptions());
const { domRef: renewalChartRef, updateOptions: updateRenewalChart } = useEcharts(() => buildRenewalOptions());

watch(
  trend,
  () => {
    void updateRevenueChart(() => buildRevenueOptions());
    void updateRenewalChart(() => buildRenewalOptions());
  },
  { deep: true }
);

watch(activeTab, value => {
  if (value === 'overview') {
    void nextTick(() => updateRevenueChart(() => buildRevenueOptions()));
  }
  if (value === 'renewal') {
    void nextTick(() => updateRenewalChart(() => buildRenewalOptions()));
  }
});

watch(
  [() => orderFilters.source, () => orderFilters.type, () => orderFilters.keyword],
  () => {
    orderPage.value = 1;
  }
);

watch(
  () => filteredOrders.value.length,
  () => {
    if (orderPage.value > orderPageCount.value) orderPage.value = orderPageCount.value;
  }
);

onMounted(() => {
  void loadFinance();
  void loadTrend();
  void loadRenewalCandidates();
});
</script>

<template>
  <div class="fin-page">
    <!-- 页头 -->
    <header class="fin-head">
      <div class="fin-head-main">
        <div class="fin-eyebrow">
          财务 · 成本与收入 <span class="fin-dim">·</span> 财务月份按 UTC 统计
        </div>
        <h1 class="fin-title">成本与收入</h1>
        <p class="fin-desc">
          服务期收入按服务期与月份交集按秒分摊（应计口径），现金收款按 paid_at 归属月份（收付实现口径）；两者不可互相抵扣。
          成本按当期生效记录直接相加，<b>不做用户均摊</b>。
        </p>
      </div>
      <div class="fin-head-actions">
        <DataFreshness :data-at="summary?.dataAt" compact />
        <button class="fin-icon-btn" type="button" title="刷新数据" :disabled="loading" @click="refreshAll">
          <SvgIcon icon="mdi:refresh" />
        </button>
        <NDatePicker
          :formatted-value="period"
          type="month"
          value-format="yyyy-MM"
          size="small"
          style="width: 150px"
          @update:formattedValue="onPeriodChange"
        />
        <button class="fin-btn" type="button" @click="exportCsv">
          <SvgIcon icon="mdi:download-outline" class="fin-btn-ico" />导出
        </button>
        <button class="fin-btn is-primary" type="button" @click="openCostModal">
          <SvgIcon icon="mdi:plus" class="fin-btn-ico" />新增记录
        </button>
      </div>
    </header>

    <div v-if="errorMessage" class="fin-alert is-danger">{{ errorMessage }}</div>

    <NSpin :show="loading">
      <template v-if="summary">
        <div class="fin-body">
          <!-- 口径提示条 -->
          <div class="fin-callout is-info">
            <SvgIcon icon="mdi:information-outline" class="fin-callout-ico" />
            <div class="fin-callout-text">
              <b>服务期收入 = 按服务期与月份交集按秒分摊；现金收款 = paid_at 在本月的已确认记录；两者不可互相抵扣。</b>
              毛利 = 服务期收入 − 本月成本；现金余额 = 现金收款 − 本月成本；续费率分母只取本期到期的收费用户。
            </div>
            <span class="fin-badge is-warn fin-callout-badge">每月 1 日 00:00–08:00 财务月与业务日不同步</span>
          </div>

          <!-- 7 KPI -->
          <section class="fin-kpis">
            <div class="fin-kpi">
              <div class="fin-kpi-top">
                <span class="fin-kpi-label">服务期收入</span>
                <span class="fin-kpi-icon is-blue"><SvgIcon icon="mdi:coins" /></span>
              </div>
              <div class="fin-kpi-value">{{ formatMoney(summary.monthIncome, summary.currency) }}</div>
              <div class="fin-kpi-foot">按秒分摊 · 覆盖 {{ summary.effectiveUserCount }} 户</div>
            </div>

            <div class="fin-kpi">
              <div class="fin-kpi-top">
                <span class="fin-kpi-label">现金收款</span>
                <span class="fin-kpi-icon is-green"><SvgIcon icon="mdi:credit-card-outline" /></span>
              </div>
              <div class="fin-kpi-value">{{ formatMoney(summary.cashIncome, summary.currency) }}</div>
              <div class="fin-kpi-foot">本月 {{ orders.length }} 笔已确认</div>
            </div>

            <div class="fin-kpi">
              <div class="fin-kpi-top">
                <span class="fin-kpi-label">毛利</span>
                <span class="fin-kpi-icon is-violet"><SvgIcon icon="mdi:chart-line" /></span>
              </div>
              <div class="fin-kpi-value">{{ formatMoney(summary.grossProfit, summary.currency) }}</div>
              <div class="fin-kpi-foot">收入−成本 · 毛利率 {{ grossMarginText }}</div>
            </div>

            <div class="fin-kpi">
              <div class="fin-kpi-top">
                <span class="fin-kpi-label">现金余额</span>
                <span class="fin-kpi-icon is-cyan"><SvgIcon icon="mdi:wallet-outline" /></span>
              </div>
              <div class="fin-kpi-value">{{ formatMoney(summary.cashBalance, summary.currency) }}</div>
              <div class="fin-kpi-foot">现金收款−本月成本</div>
            </div>

            <div class="fin-kpi">
              <div class="fin-kpi-top">
                <span class="fin-kpi-label">本月成本</span>
                <span class="fin-kpi-icon is-amber"><SvgIcon icon="mdi:layers-outline" /></span>
              </div>
              <div class="fin-kpi-value">{{ formatMoney(summary.monthCost, summary.currency) }}</div>
              <div
                class="fin-kpi-foot"
                :title="`节点 ${formatMoney(nodeCostTotal, summary.currency)} · 流量 ${formatMoney(
                  ipCostTotal,
                  summary.currency
                )} · 其它 ${formatMoney(otherCostTotal, summary.currency)}`"
              >
                节点 {{ formatMoney(nodeCostTotal, summary.currency) }} · 流量
                {{ formatMoney(ipCostTotal, summary.currency) }} · 其它
                {{ formatMoney(otherCostTotal, summary.currency) }}
              </div>
            </div>

            <div class="fin-kpi">
              <div class="fin-kpi-top">
                <span class="fin-kpi-label">毛利率</span>
                <span class="fin-kpi-icon is-pink"><SvgIcon icon="mdi:percent-outline" /></span>
              </div>
              <div class="fin-kpi-value">{{ grossMarginText }}</div>
              <div class="fin-kpi-foot">毛利 ÷ 服务期收入</div>
            </div>

            <div class="fin-kpi">
              <div class="fin-kpi-top">
                <span class="fin-kpi-label">待确认续费金额</span>
                <span class="fin-kpi-icon is-amber"><SvgIcon icon="mdi:clock-outline" /></span>
              </div>
              <div class="fin-kpi-value">
                {{ candidateLoading ? '--' : formatMoney(candidateTotal, 'CNY') }}
              </div>
              <div class="fin-kpi-foot">{{ candidates.length }} 笔候选 · 未计入</div>
            </div>
          </section>

          <!-- Tab 栏 -->
          <section class="fin-tabs">
            <button
              v-for="tab in tabs"
              :key="tab.value"
              type="button"
              :class="{ 'is-active': activeTab === tab.value }"
              @click="activeTab = tab.value"
            >
              {{ tab.label }}
            </button>
          </section>

          <!-- ===== 总览 ===== -->
          <div v-show="activeTab === 'overview'" class="fin-panel">
            <div class="fin-split">
              <section class="fin-card">
                <div class="fin-card-head">
                  <div class="fin-card-title">
                    <h3>收入与成本</h3>
                    <span>近 12 个月 · 服务期收入 vs 当月成本 · CNY</span>
                  </div>
                  <div class="fin-card-tools">
                    <span class="fin-badge is-neutral">{{ trendRangeLabel }}</span>
                  </div>
                </div>
                <div class="fin-card-body">
                  <div class="fin-legend">
                    <span><i style="background: #4c6fff"></i>服务期收入（按秒分摊）</span>
                    <span><i style="background: #f5b93d"></i>当月成本</span>
                    <span class="fin-legend-note">单位：CNY · 成本不做用户均摊</span>
                  </div>
                  <div class="fin-chart-wrap">
                    <div ref="revenueChartRef" class="fin-chart"></div>
                    <div v-if="trendLoading" class="fin-chart-state">正在读取近 12 个月财务数据…</div>
                    <div v-else-if="!trend.length" class="fin-chart-state">
                      {{ trendError || '暂无趋势数据' }}
                    </div>
                  </div>
                </div>
              </section>

              <section class="fin-card">
                <div class="fin-card-head">
                  <div class="fin-card-title">
                    <h3>成本构成</h3>
                  </div>
                  <span class="fin-badge is-neutral">本月 {{ formatMoney(summary.monthCost, summary.currency) }}</span>
                </div>
                <div class="fin-card-body">
                  <div class="fin-legend">
                    <span class="fin-dim">占本月成本比例</span>
                    <span class="fin-legend-note">共 {{ costGroup.length }} 类 · 不做用户均摊</span>
                  </div>
                  <div class="fin-stack">
                    <span
                      v-for="group in costGroup"
                      :key="group.key"
                      :style="{ width: `${costPercent(group.amount)}%`, background: group.color }"
                    ></span>
                  </div>
                  <div class="fin-cost-rows">
                    <div v-for="group in costGroup" :key="group.key" class="fin-cost-row">
                      <i :style="{ background: group.color }"></i>
                      <span class="is-strong">{{ group.label }}</span>
                      <span class="is-num">{{ formatMoney(group.amount, summary.currency) }}</span>
                      <span class="is-num fin-dim">{{ formatPercentValue(costPercent(group.amount)) }}</span>
                    </div>
                  </div>
                  <div class="fin-hr"></div>
                  <div v-if="largestCost" class="fin-largest">
                    <span class="fin-badge is-warn">最大单项</span>
                    <span class="fin-dim">
                      {{ largestCost.name }} {{ formatMoney(largestCost.amount, largestCost.currency) }} ·
                      {{ formatDay(largestCost.effectiveFrom) }} 起
                    </span>
                  </div>
                  <p class="fin-note">
                    流量成本为出口 IP 月租合计（{{ formatMoney(ipCostTotal, summary.currency) }}），明细在「出口 IP」页维护；
                    其它成本来自 other_costs，目前无 CRUD 接口。
                  </p>
                </div>
              </section>
            </div>

            <section class="fin-card">
              <div class="fin-card-head">
                <div class="fin-card-title">
                  <h3>本月收款订单</h3>
                  <span>按实际收款日期归属 {{ period }}</span>
                </div>
                <span class="fin-badge is-ok">
                  {{ orders.length }} 笔 · {{ formatMoney(ordersTotal, summary.currency) }}
                </span>
                <div class="fin-card-tools">
                  <button class="fin-btn is-ghost fin-btn-xs" type="button" @click="activeTab = 'orders'">
                    查看全部订单<svg class="fin-chevron" viewBox="0 0 24 24" aria-hidden="true">
                      <path d="M9 6l6 6-6 6" fill="none" stroke="currentColor" stroke-width="2" />
                    </svg>
                  </button>
                </div>
              </div>
              <div class="fin-table-wrap">
                <table class="fin-table">
                  <thead>
                    <tr>
                      <th>用户 / Inbound</th>
                      <th>类型</th>
                      <th class="is-num">金额</th>
                      <th>收款时间</th>
                      <th>服务期</th>
                      <th>来源</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="order in orders.slice(0, 8)" :key="order.id">
                      <td class="is-strong">
                        {{ order.userName || '未命名用户' }}
                        <span class="fin-sub">{{ order.nodeName || '--' }} · {{ order.inboundTag || '--' }}</span>
                      </td>
                      <td>
                        <span class="fin-badge" :class="orderTypeBadge(order.orderType)">
                          {{ orderTypeLabel(order.orderType) }}
                        </span>
                      </td>
                      <td class="is-num">{{ formatMoney(order.amount, order.currency) }}</td>
                      <td class="is-num">{{ formatDateTime(order.paidAt) }}</td>
                      <td class="is-num">{{ formatDay(order.serviceFrom) }} → {{ formatDay(order.serviceTo) }}</td>
                      <td>
                        <span class="fin-badge" :class="sourceBadge(order.source)">
                          {{ sourceLabel(order.source) }}
                        </span>
                      </td>
                    </tr>
                  </tbody>
                </table>
                <div v-if="!orders.length" class="fin-empty-small">本月暂无已确认收款订单</div>
              </div>
            </section>
          </div>

          <!-- ===== 收款订单 ===== -->
          <section v-show="activeTab === 'orders'" class="fin-card">
            <div class="fin-card-head">
              <div class="fin-card-title">
                <h3>收款订单台账</h3>
                <span>按 paid_at 归属 {{ period }} · UTC</span>
              </div>
              <span class="fin-badge is-ok">
                {{ filteredOrders.length }} 笔 · {{ formatMoney(filteredOrdersTotal, summary.currency) }}
              </span>
              <div class="fin-card-tools">
                <button class="fin-btn is-ghost fin-btn-xs" type="button" @click="exportCsv">
                  <SvgIcon icon="mdi:download-outline" class="fin-btn-ico" />导出
                </button>
              </div>
            </div>
            <div class="fin-toolbar">
              <NSelect
                v-model:value="orderFilters.source"
                :options="sourceOptions"
                size="small"
                style="width: 132px"
              />
              <NSelect
                v-model:value="orderFilters.type"
                :options="orderTypeOptions"
                size="small"
                style="width: 120px"
              />
              <label class="fin-search">
                <SvgIcon icon="mdi:magnify" class="fin-search-ico" />
                <input v-model="orderFilters.keyword" placeholder="搜索用户 / 订单号" />
              </label>
              <span class="fin-spacer"></span>
              <span class="fin-dim fin-hint">
                合计 {{ formatMoney(filteredOrdersTotal, summary.currency) }} · {{ filteredOrders.length }} 笔 ·
                全部为已确认且已核验
              </span>
            </div>
            <div class="fin-table-wrap">
              <table class="fin-table">
                <thead>
                  <tr>
                    <th>用户 / Inbound</th>
                    <th>订单号</th>
                    <th>类型</th>
                    <th>周期</th>
                    <th class="is-num">金额</th>
                    <th>收款时间</th>
                    <th>服务期</th>
                    <th>来源</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="order in pagedOrders" :key="order.id">
                    <td class="is-strong">
                      {{ order.userName || '未命名用户' }}
                      <span class="fin-sub">{{ order.nodeName || '--' }} · {{ order.inboundTag || '--' }}</span>
                    </td>
                    <td class="fin-mono">{{ order.orderNo || order.id.slice(0, 12) }}</td>
                    <td>
                      <span class="fin-badge" :class="orderTypeBadge(order.orderType)">
                        {{ orderTypeLabel(order.orderType) }}
                      </span>
                    </td>
                    <td class="fin-dim">{{ cycleLabel(order.billingCycle) }}</td>
                    <td class="is-num">{{ formatMoney(order.amount, order.currency) }}</td>
                    <td class="is-num">{{ formatDateTime(order.paidAt) }}</td>
                    <td class="is-num">{{ formatDay(order.serviceFrom) }} → {{ formatDay(order.serviceTo) }}</td>
                    <td>
                      <span class="fin-badge" :class="sourceBadge(order.source)">
                        {{ sourceLabel(order.source) }}
                      </span>
                    </td>
                  </tr>
                </tbody>
                <tfoot v-if="filteredOrders.length">
                  <tr>
                    <td class="is-strong" colspan="4">合计 {{ filteredOrders.length }} 笔</td>
                    <td class="is-num is-strong">{{ formatMoney(filteredOrdersTotal, summary.currency) }}</td>
                    <td colspan="3" class="fin-dim">
                      首购 {{ formatMoney(typeTotal('initial'), summary.currency) }} · 续费
                      {{ formatMoney(typeTotal('renewal'), summary.currency) }} · 逾期恢复
                      {{ formatMoney(typeTotal('recovery'), summary.currency) }}
                    </td>
                  </tr>
                </tfoot>
              </table>
              <div v-if="!filteredOrders.length" class="fin-empty-small">没有符合筛选条件的收款订单</div>
            </div>
            <div class="fin-pager">
              <span class="fin-dim">
                显示 {{ pageStart }}–{{ pageEnd }} 条 · 共 {{ filteredOrders.length }} 条 · 每页 {{ orderPageSize }} 条
              </span>
              <div class="fin-pager-actions">
                <button type="button" :disabled="orderPage <= 1" @click="orderPage -= 1">上一页</button>
                <span>{{ orderPage }} / {{ orderPageCount }}</span>
                <button type="button" :disabled="orderPage >= orderPageCount" @click="orderPage += 1">
                  下一页
                </button>
              </div>
            </div>
          </section>

          <!-- ===== 续费分析 ===== -->
          <div v-show="activeTab === 'renewal'" class="fin-panel">
            <section class="fin-kpis is-3">
              <div class="fin-kpi">
                <div class="fin-kpi-top">
                  <span class="fin-kpi-label">续费候选</span>
                  <span class="fin-kpi-icon is-amber"><SvgIcon icon="mdi:clock-outline" /></span>
                </div>
                <div class="fin-kpi-value">
                  {{ candidates.length }}<span class="fin-kpi-unit">笔</span>
                </div>
                <div class="fin-kpi-foot">Agent 只提议不记账 · 待运营确认</div>
              </div>
              <div class="fin-kpi">
                <div class="fin-kpi-top">
                  <span class="fin-kpi-label">待确认金额</span>
                  <span class="fin-kpi-icon is-amber"><SvgIcon icon="mdi:coins" /></span>
                </div>
                <div class="fin-kpi-value">{{ formatMoney(candidateTotal, 'CNY') }}</div>
                <div class="fin-kpi-foot">确认后计入服务期收入与现金收款</div>
              </div>
              <div class="fin-kpi">
                <div class="fin-kpi-top">
                  <span class="fin-kpi-label">本月续费率</span>
                  <span class="fin-kpi-icon is-green"><SvgIcon icon="mdi:activity" /></span>
                </div>
                <div class="fin-kpi-value">{{ formatPercentValue(renewal.renewalRate * 100) }}</div>
                <div class="fin-kpi-foot">
                  到期前已续费 {{ renewal.renewedUsers }} / 到期 {{ renewal.dueUsers }} 户 · 免费用户不计入
                </div>
              </div>
            </section>

            <div class="fin-split is-renewal">
              <section class="fin-card">
                <div class="fin-card-head">
                  <div class="fin-card-title">
                    <h3>续费率趋势</h3>
                    <span>近 6 个月 · 严格续费率（到期前已确认 renewal 订单）</span>
                  </div>
                </div>
                <div class="fin-card-body">
                  <div class="fin-legend">
                    <span><i style="background: #4c6fff"></i>严格续费率</span>
                    <span class="fin-legend-note">本期 {{ formatPercentValue(renewal.renewalRate * 100) }}</span>
                  </div>
                  <div class="fin-chart-wrap is-sm">
                    <div ref="renewalChartRef" class="fin-chart is-sm"></div>
                    <div v-if="trendLoading" class="fin-chart-state">正在读取续费趋势…</div>
                    <div v-else-if="!renewalTrendPoints.length" class="fin-chart-state">
                      {{ trendError || '暂无续费趋势数据' }}
                    </div>
                  </div>
                </div>
              </section>

              <section class="fin-card">
                <div class="fin-card-head">
                  <div class="fin-card-title">
                    <h3>续费候选</h3>
                  </div>
                  <span class="fin-badge is-warn">{{ candidates.length }} 笔待确认</span>
                  <div class="fin-card-tools">
                    <span class="fin-dim fin-hint">合计 {{ formatMoney(candidateTotal, 'CNY') }}</span>
                  </div>
                </div>
                <div v-if="candidateError" class="fin-alert is-danger fin-alert-inline">{{ candidateError }}</div>
                <div class="fin-table-wrap">
                  <table class="fin-table">
                    <thead>
                      <tr>
                        <th>用户</th>
                        <th>原到期</th>
                        <th>新到期</th>
                        <th class="is-num">增加</th>
                        <th class="is-num">金额</th>
                        <th></th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr v-for="row in candidates" :key="row.eventId">
                        <td class="is-strong">{{ row.userName }}</td>
                        <td class="is-num">{{ formatDateTime(row.oldExpiryAt) }}</td>
                        <td class="is-num">{{ formatDateTime(row.newExpiryAt) }}</td>
                        <td class="is-num">{{ expiryDays(row.oldExpiryAt, row.newExpiryAt) }}</td>
                        <td class="is-num">{{ formatMoney(row.amount, row.currency) }}</td>
                        <td class="is-num">
                          <button
                            class="fin-btn is-success fin-btn-xs"
                            type="button"
                            :disabled="Boolean(confirmingCandidateId)"
                            @click="confirmCandidate(row)"
                          >
                            {{ confirmingCandidateId === row.eventId ? '确认中…' : '确认' }}
                          </button>
                        </td>
                      </tr>
                    </tbody>
                  </table>
                  <div v-if="!candidates.length" class="fin-empty-small">
                    {{ candidateLoading ? '正在读取续费候选…' : '暂无待确认的续费候选' }}
                  </div>
                </div>
                <p class="fin-card-note">
                  续费率分母只取本期到期的收费用户；逾期恢复不计入续费率，单独统计。确认后计入服务期收入与现金收款。
                </p>
              </section>
            </div>
          </div>

          <!-- ===== 成本明细 ===== -->
          <div v-show="activeTab === 'cost'" class="fin-panel">
            <section class="fin-card">
              <div class="fin-card-head">
                <div class="fin-card-title">
                  <h3>节点成本</h3>
                </div>
                <span class="fin-badge is-neutral">
                  {{ nodeCosts.length }} 条生效记录 · {{ formatMoney(nodeCostTotal, summary.currency) }}
                </span>
                <div class="fin-card-tools">
                  <button class="fin-btn is-soft fin-btn-xs" type="button" @click="openCostModal">
                    <SvgIcon icon="mdi:plus" class="fin-btn-ico" />新增节点成本
                  </button>
                </div>
              </div>
              <div class="fin-table-wrap">
                <table class="fin-table">
                  <thead>
                    <tr>
                      <th>节点</th>
                      <th class="is-num">月成本</th>
                      <th>生效日期</th>
                      <th>备注</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="item in nodeCosts" :key="item.id">
                      <td class="is-strong">{{ item.name || '未命名节点' }}</td>
                      <td class="is-num">{{ formatMoney(item.amount, item.currency) }}</td>
                      <td class="is-num">{{ formatDay(item.effectiveFrom) }}</td>
                      <td class="fin-dim">{{ item.notes || '--' }}</td>
                    </tr>
                  </tbody>
                </table>
                <div v-if="!nodeCosts.length" class="fin-empty-small">本期没有生效的节点成本记录</div>
              </div>
              <p class="fin-card-note">
                成本时态不可回写：修改金额请新增一条记录；编辑时不可修改生效日期，变更生效日会被拒绝（409）。
                出口 IP 月租合计 {{ formatMoney(ipCostTotal, summary.currency) }}，由「出口 IP」台账维护，本页不重复维护。
              </p>
            </section>

            <section class="fin-card">
              <div class="fin-card-head">
                <div class="fin-card-title">
                  <h3>其它成本</h3>
                </div>
                <span class="fin-badge is-neutral">
                  {{ otherCosts.length }} 条 · {{ formatMoney(otherCostTotal, summary.currency) }}
                </span>
                <div class="fin-card-tools">
                  <span class="fin-badge is-ghost">目前无 CRUD，仅展示</span>
                </div>
              </div>
              <div class="fin-table-wrap">
                <table class="fin-table">
                  <thead>
                    <tr>
                      <th>项目</th>
                      <th>类别</th>
                      <th class="is-num">月金额</th>
                      <th>生效日期</th>
                      <th>备注</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="item in otherCosts" :key="item.id">
                      <td class="is-strong">{{ item.name || '--' }}</td>
                      <td class="fin-dim">{{ categoryLabel(item.category) }}</td>
                      <td class="is-num">{{ formatMoney(item.amount, item.currency) }}</td>
                      <td class="is-num">{{ formatDay(item.effectiveFrom) }}</td>
                      <td class="fin-dim">{{ item.notes || '--' }}</td>
                    </tr>
                  </tbody>
                </table>
                <div v-if="!otherCosts.length" class="fin-empty-small">本期没有其它成本记录</div>
              </div>
              <p class="fin-card-note">
                其它成本来自 other_costs，仅参与财务汇总；当前版本没有创建 / 编辑接口，只能通过数据维护写入。
              </p>
            </section>
          </div>
        </div>
      </template>
      <div v-else-if="!loading" class="fin-empty">
        <NEmpty description="暂无财务汇总数据" />
      </div>
    </NSpin>

    <!-- ===== 新增成本记录 ===== -->
    <Teleport to="body">
      <div v-if="costModalVisible" class="fin-mask" @click.self="closeCostModal">
        <div class="fin-modal">
          <div class="fin-modal-head">
            <div class="fin-modal-head-main">
              <h3>新增成本记录</h3>
              <p class="fin-modal-sub">写入后计入当月成本与后续月份</p>
            </div>
            <span class="fin-badge is-brand">时态成本</span>
            <button class="fin-icon-btn" type="button" :disabled="costSaving" @click="closeCostModal">
              <SvgIcon icon="mdi:close" />
            </button>
          </div>

          <div class="fin-modal-body">
            <div v-if="costError" class="fin-alert is-danger">{{ costError }}</div>
            <div class="fin-callout is-warn">
              <SvgIcon icon="mdi:alert-outline" class="fin-callout-ico" />
              <div class="fin-callout-text">
                <b>成本时态不可回写</b>
                保存后金额与生效日期不可修改，如需调整请新增一条记录；把生效日期改到过去不会回写历史月份。
              </div>
            </div>

            <div class="fin-form-grid">
              <div class="fin-field">
                <span class="fin-label">成本类型 <span class="fin-req">*</span></span>
                <NSelect v-model:value="costForm.type" :options="costTypeOptions" disabled />
              </div>
              <div class="fin-field">
                <span class="fin-label">关联节点 <span class="fin-req">*</span></span>
                <NSelect
                  v-model:value="costForm.nodeId"
                  :options="nodeOptions"
                  :loading="nodesLoading"
                  filterable
                  placeholder="选择节点"
                />
              </div>
              <div class="fin-field">
                <span class="fin-label">金额（CNY） <span class="fin-req">*</span></span>
                <NInput v-model:value="costForm.amount" placeholder="例如 520.00" />
              </div>
              <div class="fin-field">
                <span class="fin-label">生效日期 <span class="fin-req">*</span></span>
                <NDatePicker
                  v-model:formatted-value="costForm.effectiveFrom"
                  type="date"
                  value-format="yyyy-MM-dd"
                  clearable
                  style="width: 100%"
                />
              </div>
            </div>

            <div class="fin-field fin-field-wide">
              <span class="fin-label">备注</span>
              <NInput v-model:value="costForm.notes" maxlength="2000" placeholder="可选，例如：含 1Gbps 带宽" />
            </div>

            <p class="fin-note">
              失效日期留空表示长期有效；金额按「当月生效记录直接相加」计入成本，不按天拆分、不做用户均摊。
            </p>
          </div>

          <div class="fin-modal-foot">
            <span class="fin-dim">保存后立即参与本月成本与毛利计算</span>
            <div class="fin-foot-actions">
              <button class="fin-btn" type="button" :disabled="costSaving" @click="closeCostModal">取消</button>
              <button class="fin-btn is-primary" type="button" :disabled="costSaving" @click="submitCost">
                {{ costSaving ? '保存中…' : '保存' }}
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
   成本与收入 · UI v2
   颜色 / 字号 / 间距对齐 design-preview/v2/finance.html 与 kit-v2
   ============================================================ */
.fin-page {
  --fin-card: #141519;
  --fin-card-hover: #191b20;
  --fin-line: #23262c;
  --fin-line-soft: #1c1f24;
  --fin-text-1: #f2f3f5;
  --fin-text-2: #bcc2cc;
  --fin-text-3: #9aa1ad;
  --fin-text-4: #868d99;
  --fin-brand: #4c6fff;
  --fin-brand-soft: rgba(76, 111, 255, 0.14);
  --fin-ok: #2ed47a;
  --fin-warn: #f5b93d;
  --fin-danger: #f3656f;
  --fin-info: #5aa9f0;
  --fin-violet: #a78bfa;
  --fin-cyan: #34d3e0;
  --fin-pink: #f472b6;
  --fin-radius: 14px;

  display: flex;
  flex-direction: column;
  gap: 15px;
  min-width: 0;
  color: var(--fin-text-1);
  font-size: 13px;
}

.fin-body {
  display: flex;
  flex-direction: column;
  gap: 15px;
  min-width: 0;
}

.fin-dim {
  color: var(--fin-text-4);
}

.fin-mono {
  font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
}

/* ---------- 页头 ---------- */
.fin-head {
  display: flex;
  align-items: flex-end;
  gap: 18px;
  flex-wrap: wrap;
}

.fin-head-main {
  flex: 1;
  min-width: 0;
}

.fin-eyebrow {
  margin-bottom: 6px;
  font-size: 12px;
  letter-spacing: 0.1em;
  color: #7f8792;
}

.fin-title {
  margin: 0;
  font-size: 23px;
  font-weight: 620;
  line-height: 1.25;
  color: #f6f7f9;
}

.fin-desc {
  max-width: 860px;
  margin: 7px 0 0;
  font-size: 13px;
  line-height: 1.72;
  color: var(--fin-text-3);
}

.fin-desc b {
  font-weight: 600;
  color: var(--fin-text-2);
}

.fin-head-actions {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: none;
  flex-wrap: wrap;
  padding-bottom: 2px;
}

/* ---------- 按钮 / 输入 ---------- */
.fin-btn {
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

.fin-btn:hover {
  background: #232730;
  border-color: #363b44;
}

.fin-btn:disabled {
  opacity: 0.45;
  cursor: not-allowed;
}

.fin-btn-ico {
  width: 15px;
  height: 15px;
  flex: none;
}

.fin-btn.is-primary {
  border-color: var(--fin-brand);
  background: var(--fin-brand);
  color: #fff;
}

.fin-btn.is-primary:hover {
  background: #5d7cff;
  border-color: #5d7cff;
}

.fin-btn.is-soft {
  border-color: rgba(76, 111, 255, 0.28);
  background: rgba(76, 111, 255, 0.13);
  color: #a8b8ff;
}

.fin-btn.is-ghost {
  border-color: transparent;
  background: transparent;
  color: var(--fin-text-3);
}

.fin-btn.is-ghost:hover {
  background: #1a1d22;
  color: var(--fin-text-1);
}

.fin-btn.is-success {
  border-color: rgba(46, 212, 122, 0.34);
  background: rgba(46, 212, 122, 0.14);
  color: #5ce49b;
}

.fin-btn.is-success:hover {
  background: rgba(46, 212, 122, 0.22);
}

.fin-btn-xs {
  height: 28px;
  padding: 0 10px;
  font-size: 12px;
  border-radius: 7px;
}

.fin-chevron {
  width: 13px;
  height: 13px;
  flex: none;
}

.fin-icon-btn {
  display: grid;
  place-items: center;
  width: 30px;
  height: 30px;
  border: 1px solid transparent;
  border-radius: 8px;
  background: transparent;
  color: #9aa1ab;
  cursor: pointer;
}

.fin-icon-btn:hover {
  background: #1a1d22;
  color: #e6e8ec;
}

.fin-icon-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.fin-search {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 220px;
  height: 32px;
  padding: 0 10px;
  border: 1px solid #262a31;
  border-radius: 9px;
  background: #0e1013;
}

.fin-search:focus-within {
  border-color: rgba(76, 111, 255, 0.7);
}

.fin-search input {
  flex: 1;
  min-width: 0;
  border: 0;
  outline: none;
  background: transparent;
  color: var(--fin-text-1);
  font-size: 12.5px;
  font-family: inherit;
}

.fin-search input::placeholder {
  color: var(--fin-text-4);
}

.fin-search-ico {
  width: 15px;
  height: 15px;
  flex: none;
  color: var(--fin-text-4);
}

/* ---------- 提示条 / 徽标 ---------- */
.fin-callout {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 11px 14px;
  border: 1px solid #2a2e35;
  border-radius: 11px;
  background: #111317;
  font-size: 12.5px;
  line-height: 1.65;
  color: var(--fin-text-2);
}

.fin-callout.is-info {
  border-color: rgba(90, 169, 240, 0.28);
  background: rgba(90, 169, 240, 0.09);
}

.fin-callout.is-warn {
  border-color: rgba(245, 185, 61, 0.3);
  background: rgba(245, 185, 61, 0.09);
}

.fin-callout-ico {
  width: 15px;
  height: 15px;
  flex: none;
  margin-top: 2px;
  color: var(--fin-info);
}

.fin-callout.is-warn .fin-callout-ico {
  color: var(--fin-warn);
}

.fin-callout-text {
  flex: 1;
  min-width: 0;
}

.fin-callout-text b {
  display: block;
  margin-bottom: 1px;
  font-size: 12.5px;
  font-weight: 600;
  color: var(--fin-text-1);
}

.fin-callout-badge {
  flex: none;
  align-self: center;
}

.fin-badge {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  height: 22px;
  padding: 0 9px;
  border: 1px solid #2a2e35;
  border-radius: 7px;
  background: #1b1e24;
  color: var(--fin-text-3);
  font-size: 11.5px;
  font-weight: 550;
  line-height: 1;
  white-space: nowrap;
}

.fin-badge.is-neutral {
  border-color: #2a2e35;
  background: #1b1e24;
  color: var(--fin-text-2);
}

.fin-badge.is-ok {
  border-color: rgba(46, 212, 122, 0.34);
  background: rgba(46, 212, 122, 0.14);
  color: #5ce49b;
}

.fin-badge.is-warn {
  border-color: rgba(245, 185, 61, 0.34);
  background: rgba(245, 185, 61, 0.14);
  color: #f7ca6b;
}

.fin-badge.is-info {
  border-color: rgba(90, 169, 240, 0.34);
  background: rgba(90, 169, 240, 0.14);
  color: #8cc4f7;
}

.fin-badge.is-brand {
  border-color: rgba(76, 111, 255, 0.38);
  background: rgba(76, 111, 255, 0.14);
  color: #a8b8ff;
}

.fin-badge.is-violet {
  border-color: rgba(167, 139, 250, 0.36);
  background: rgba(167, 139, 250, 0.14);
  color: #c0a9ff;
}

.fin-badge.is-ghost {
  border-color: transparent;
  background: transparent;
  color: var(--fin-text-4);
}

.fin-alert {
  padding: 10px 13px;
  border: 1px solid #2a2e35;
  border-radius: 10px;
  font-size: 12.5px;
  line-height: 1.6;
}

.fin-alert.is-danger {
  border-color: rgba(243, 101, 111, 0.34);
  background: rgba(243, 101, 111, 0.1);
  color: #ff9ba2;
}

.fin-alert-inline {
  margin: 10px 14px 0;
}

/* ---------- KPI ---------- */
.fin-kpis {
  display: grid;
  grid-template-columns: repeat(7, minmax(0, 1fr));
  gap: 12px;
}

.fin-kpis.is-3 {
  grid-template-columns: repeat(3, minmax(0, 1fr));
}

.fin-kpi {
  min-width: 0;
  padding: 12px 13px 11px;
  border: 1px solid var(--fin-line);
  border-radius: var(--fin-radius);
  background: var(--fin-card);
}

.fin-kpi-top {
  display: flex;
  align-items: center;
  gap: 8px;
}

.fin-kpi-label {
  min-width: 0;
  overflow: hidden;
  font-size: 12px;
  color: var(--fin-text-3);
  white-space: nowrap;
  text-overflow: ellipsis;
}

.fin-kpi-icon {
  display: grid;
  place-items: center;
  width: 22px;
  height: 22px;
  margin-left: auto;
  flex: none;
  border: 1px solid #2a2e35;
  border-radius: 7px;
  background: #1d2026;
  color: #a8adb8;
}

.fin-kpi-icon :deep(.svg-icon) {
  width: 12px;
  height: 12px;
}

.fin-kpi-icon.is-blue {
  border-color: rgba(76, 111, 255, 0.36);
  background: rgba(76, 111, 255, 0.16);
  color: #8fa8ff;
}

.fin-kpi-icon.is-green {
  border-color: rgba(46, 212, 122, 0.34);
  background: rgba(46, 212, 122, 0.14);
  color: #5ce49b;
}

.fin-kpi-icon.is-violet {
  border-color: rgba(167, 139, 250, 0.36);
  background: rgba(167, 139, 250, 0.16);
  color: #c0a9ff;
}

.fin-kpi-icon.is-cyan {
  border-color: rgba(52, 211, 224, 0.34);
  background: rgba(52, 211, 224, 0.14);
  color: #6fe0e9;
}

.fin-kpi-icon.is-amber {
  border-color: rgba(245, 185, 61, 0.34);
  background: rgba(245, 185, 61, 0.14);
  color: #f7ca6b;
}

.fin-kpi-icon.is-pink {
  border-color: rgba(244, 114, 182, 0.34);
  background: rgba(244, 114, 182, 0.14);
  color: #f9a8d4;
}

.fin-kpi-value {
  display: flex;
  align-items: baseline;
  gap: 5px;
  margin-top: 8px;
  overflow: hidden;
  font-size: 20px;
  font-weight: 620;
  line-height: 1.2;
  letter-spacing: -0.01em;
  color: #fff;
  font-variant-numeric: tabular-nums;
  white-space: nowrap;
  text-overflow: ellipsis;
}

.fin-kpi-unit {
  font-size: 12px;
  font-weight: 500;
  color: var(--fin-text-3);
}

.fin-kpi-foot {
  margin-top: 6px;
  overflow: hidden;
  font-size: 11.5px;
  line-height: 1.5;
  color: var(--fin-text-4);
  white-space: nowrap;
  text-overflow: ellipsis;
}

/* ---------- Tab 栏 ---------- */
.fin-tabs {
  display: flex;
  align-items: center;
  gap: 2px;
  padding: 5px;
  border: 1px solid var(--fin-line);
  border-radius: var(--fin-radius);
  background: var(--fin-card);
  overflow-x: auto;
}

.fin-tabs button {
  height: 32px;
  padding: 0 15px;
  border: 0;
  border-radius: 8px;
  background: transparent;
  color: var(--fin-text-3);
  font-size: 12.5px;
  font-weight: 550;
  white-space: nowrap;
  cursor: pointer;
  transition: background 0.15s ease, color 0.15s ease;
}

.fin-tabs button:hover {
  color: var(--fin-text-1);
}

.fin-tabs button.is-active {
  background: rgba(76, 111, 255, 0.16);
  color: #fff;
  box-shadow: inset 0 -2px 0 var(--fin-brand);
}

/* ---------- 面板 / 卡片 ---------- */
.fin-panel {
  display: flex;
  flex-direction: column;
  gap: 15px;
  min-width: 0;
}

.fin-split {
  display: grid;
  grid-template-columns: minmax(0, 1.62fr) minmax(0, 1fr);
  gap: 15px;
}

.fin-split.is-renewal {
  grid-template-columns: minmax(0, 0.92fr) minmax(0, 1.62fr);
}

.fin-card {
  min-width: 0;
  border: 1px solid var(--fin-line);
  border-radius: var(--fin-radius);
  background: var(--fin-card);
}

.fin-card-head {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 15px 18px;
  border-bottom: 1px solid var(--fin-line-soft);
}

.fin-card-title {
  display: flex;
  flex-direction: column;
  gap: 3px;
  min-width: 0;
}

.fin-card-title h3 {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
  color: var(--fin-text-1);
  white-space: nowrap;
}

.fin-card-title span {
  overflow: hidden;
  font-size: 12px;
  color: var(--fin-text-4);
  white-space: nowrap;
  text-overflow: ellipsis;
}

.fin-card-tools {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-left: auto;
  flex: none;
}

.fin-card-body {
  padding: 15px 18px 18px;
}

.fin-card-note {
  margin: 0;
  padding: 11px 18px;
  border-top: 1px solid var(--fin-line-soft);
  font-size: 12px;
  line-height: 1.65;
  color: var(--fin-text-4);
}

.fin-hint {
  font-size: 12px;
}

/* ---------- 图表 ---------- */
.fin-legend {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 8px;
  font-size: 12px;
  color: var(--fin-text-3);
}

.fin-legend span {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.fin-legend i {
  width: 8px;
  height: 8px;
  border-radius: 2px;
}

.fin-legend-note {
  margin-left: auto;
  color: var(--fin-text-4);
}

.fin-chart-wrap {
  position: relative;
  min-height: 248px;
}

.fin-chart-wrap.is-sm {
  min-height: 178px;
}

.fin-chart {
  width: 100%;
  height: 248px;
}

.fin-chart.is-sm {
  height: 178px;
}

.fin-chart-state {
  position: absolute;
  inset: 0;
  display: grid;
  place-items: center;
  font-size: 12.5px;
  color: var(--fin-text-4);
  pointer-events: none;
}

/* ---------- 成本构成 ---------- */
.fin-stack {
  display: flex;
  gap: 2px;
  height: 14px;
  border-radius: 7px;
  overflow: hidden;
  background: #1c1f24;
}

.fin-stack span {
  display: block;
  height: 100%;
  min-width: 0;
  transition: width 0.2s ease;
}

.fin-cost-rows {
  margin-top: 12px;
}

.fin-cost-row {
  display: grid;
  grid-template-columns: 9px minmax(0, 1fr) auto 54px;
  align-items: center;
  gap: 9px;
  padding: 9px 0;
  border-bottom: 1px solid var(--fin-line-soft);
  font-size: 12.5px;
}

.fin-cost-row:last-child {
  padding-bottom: 2px;
  border-bottom: 0;
}

.fin-cost-row i {
  width: 8px;
  height: 8px;
  border-radius: 2px;
}

.fin-cost-row .is-strong {
  min-width: 0;
  overflow: hidden;
  color: var(--fin-text-2);
  white-space: nowrap;
  text-overflow: ellipsis;
}

.fin-cost-row .is-num {
  color: var(--fin-text-2);
  font-variant-numeric: tabular-nums;
  text-align: right;
}

.fin-cost-row .is-num.fin-dim {
  color: var(--fin-text-4);
}

.fin-hr {
  height: 1px;
  margin: 13px 0;
  background: var(--fin-line-soft);
}

.fin-largest {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12.5px;
}

.fin-note {
  margin: 12px 0 0;
  font-size: 12px;
  line-height: 1.65;
  color: var(--fin-text-4);
}

/* ---------- 表格 ---------- */
.fin-table-wrap {
  overflow-x: auto;
}

.fin-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 13px;
}

.fin-table th {
  padding: 10px 14px;
  border-bottom: 1px solid var(--fin-line-soft);
  color: var(--fin-text-3);
  font-size: 11.5px;
  font-weight: 600;
  letter-spacing: 0.04em;
  text-align: left;
  white-space: nowrap;
}

.fin-table td {
  padding: 10px 14px;
  border-bottom: 1px solid var(--fin-line-soft);
  color: var(--fin-text-2);
  vertical-align: middle;
}

.fin-table tbody tr:hover td {
  background: var(--fin-card-hover);
}

.fin-table tbody tr:last-child td {
  border-bottom: 0;
}

.fin-table tfoot td {
  padding: 11px 14px;
  border-top: 1px solid var(--fin-line);
  border-bottom: 0;
  background: #111317;
}

.fin-table .is-num {
  font-variant-numeric: tabular-nums;
  text-align: right;
  white-space: nowrap;
}

.fin-table .is-strong {
  color: var(--fin-text-1);
  font-weight: 550;
}

.fin-sub {
  display: block;
  margin-top: 2px;
  overflow: hidden;
  font-size: 12px;
  font-weight: 400;
  color: var(--fin-text-4);
  white-space: nowrap;
  text-overflow: ellipsis;
}

.fin-empty-small {
  display: grid;
  place-items: center;
  min-height: 96px;
  padding: 20px;
  font-size: 12.5px;
  color: var(--fin-text-4);
}

.fin-toolbar {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
  padding: 12px 14px;
  border-bottom: 1px solid var(--fin-line-soft);
}

.fin-spacer {
  flex: 1;
}

.fin-pager {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 11px 16px;
  border-top: 1px solid var(--fin-line-soft);
  font-size: 12px;
}

.fin-pager-actions {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-left: auto;
  color: var(--fin-text-3);
}

.fin-pager-actions button {
  height: 26px;
  padding: 0 10px;
  border: 1px solid #2a2e35;
  border-radius: 7px;
  background: #1b1e24;
  color: var(--fin-text-2);
  font-size: 12px;
  cursor: pointer;
}

.fin-pager-actions button:hover:not(:disabled) {
  border-color: #363b44;
  background: #232730;
}

.fin-pager-actions button:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

/* ---------- 空态 ---------- */
.fin-empty {
  display: grid;
  place-items: center;
  min-height: 240px;
}

/* ---------- 弹窗 ---------- */
.fin-mask {
  position: fixed;
  inset: 0;
  z-index: 2000;
  display: grid;
  place-items: center;
  padding: 20px;
  background: rgba(0, 0, 0, 0.62);
  backdrop-filter: blur(2px);
}

.fin-modal {
  display: flex;
  flex-direction: column;
  width: 560px;
  max-width: 100%;
  max-height: min(86vh, 720px);
  border: 1px solid #2a2e35;
  border-radius: 16px;
  background: var(--fin-card);
  color: var(--fin-text-1);
  text-align: left;
  box-shadow: 0 30px 80px -24px rgba(0, 0, 0, 0.9);
}

.fin-modal-head {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 16px 20px 14px;
  border-bottom: 1px solid var(--fin-line-soft);
}

.fin-modal-head-main {
  flex: 1;
  min-width: 0;
}

.fin-modal-head h3 {
  margin: 0;
  font-size: 15px;
  font-weight: 600;
}

.fin-modal-sub {
  margin: 3px 0 0;
  font-size: 12px;
  color: var(--fin-text-4);
}

.fin-modal-body {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  padding: 16px 20px 20px;
}

.fin-modal-body .fin-callout {
  margin-bottom: 16px;
}

.fin-form-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.fin-field {
  display: flex;
  flex-direction: column;
  gap: 6px;
  min-width: 0;
}

.fin-field-wide {
  margin-top: 12px;
}

.fin-label {
  font-size: 12px;
  color: var(--fin-text-3);
}

.fin-req {
  color: var(--fin-danger);
}

.fin-modal-foot {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 13px 20px;
  border-top: 1px solid var(--fin-line-soft);
  font-size: 12px;
}

.fin-foot-actions {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-left: auto;
}

/* ---------- 响应式 ---------- */
@media (max-width: 1440px) {
  .fin-kpis {
    grid-template-columns: repeat(4, minmax(0, 1fr));
  }
}

@media (max-width: 1180px) {
  .fin-split,
  .fin-split.is-renewal {
    grid-template-columns: minmax(0, 1fr);
  }
}

@media (max-width: 1100px) {
  .fin-kpis {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .fin-kpis.is-3 {
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }
}

@media (max-width: 720px) {
  .fin-kpis.is-3,
  .fin-form-grid {
    grid-template-columns: minmax(0, 1fr);
  }

  .fin-head-actions {
    width: 100%;
  }
}
</style>
