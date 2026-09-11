<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue';
import { useRouter } from 'vue-router';
import { fetchDashboard } from '@/service/api';
import DataFreshness from '@/components/project/data-freshness.vue';
import NodeStatusTag from '@/components/project/node-status-tag.vue';
import TrafficValue from '@/components/project/traffic-value.vue';
import UserStatusTag from '@/components/project/user-status-tag.vue';
import { useEcharts, type ECOption } from '@/hooks/common/echarts';

defineOptions({ name: 'CentralDashboard' });

type RankingMetric = 'total' | 'download' | 'upload';

const router = useRouter();
const loading = ref(false);
const errorMessage = ref('');
const dashboard = ref<Api.Central.DashboardSummary | null>(null);
const range = ref<Api.Central.DashboardRange>('today');
const customRange = ref<[number, number] | null>(null);
const nodeMetric = ref<RankingMetric>('total');
const userMetric = ref<RankingMetric>('total');

const rangeOptions: Array<{ label: string; value: Api.Central.DashboardRange }> = [
  { label: '今日', value: 'today' },
  { label: '近 7 天', value: '7d' },
  { label: '近 30 天', value: '30d' },
  { label: '自定义', value: 'custom' }
];

const nodeRanking = computed(() => {
  const items = [...(dashboard.value?.nodeTrafficRanking || [])];
  const key = metricKey(nodeMetric.value);
  return items.sort((left, right) => right[key] - left[key]);
});

const userRanking = computed(() => {
  const items = [...(dashboard.value?.userTrafficRanking || [])];
  const key = metricKey(userMetric.value);
  return items.sort((left, right) => right[key] - left[key]);
});

const nodeMaximum = computed(() => Math.max(1, ...nodeRanking.value.map(item => item[metricKey(nodeMetric.value)])));

const expiryItems = computed(() => {
  const expiry = dashboard.value?.expiry;
  if (!expiry) return [];
  return [
    { label: '30 天后 / 长期', value: expiry.normalAfter30d, color: 'primary' },
    { label: '8–30 天到期', value: expiry.expiring30d, color: 'info' },
    { label: '7 天内到期', value: expiry.expiring7d, color: 'warning' },
    { label: '已到期', value: expiry.expired, color: 'error' }
  ];
});

const expiryMaximum = computed(() => Math.max(1, ...expiryItems.value.map(item => item.value)));
const trendPoints = computed(() => dashboard.value?.trafficTrend.points || []);
const hasTrendData = computed(() => trendPoints.value.length > 0);
const hasTrendReset = computed(() => trendPoints.value.some(point => point.resetDetected));
const hasTrendGap = computed(() => trendPoints.value.some(point => point.hasGap));
const rangeDescription = computed(() => {
  const selected = dashboard.value?.range;
  if (!selected) return '等待中央数据';
  return `${formatDate(selected.from)} 至 ${formatDate(selected.to)}`;
});
const trendSubtitle = computed(() => {
  const bucket = dashboard.value?.trafficTrend.bucket === '1h' ? '按小时' : '按天';
  return `${bucket}累计统计，上传与下载均来自线路机用户 Inbound`;
});
const topNodeShare = computed(() => {
  const total = dashboard.value?.traffic.totalBytes || 0;
  if (!total) return 0;
  return nodeRanking.value.slice(0, 5).reduce((sum, item) => sum + item.totalBytes, 0) / total;
});
const grossMargin = computed(() => {
  const finance = dashboard.value?.finance;
  if (!finance || finance.monthIncome <= 0) return null;
  return finance.grossProfit / finance.monthIncome;
});

function metricKey(metric: RankingMetric): 'totalBytes' | 'downloadBytes' | 'uploadBytes' {
  if (metric === 'download') return 'downloadBytes';
  if (metric === 'upload') return 'uploadBytes';
  return 'totalBytes';
}

function metricLabel(metric: RankingMetric) {
  return metric === 'download' ? '下载' : metric === 'upload' ? '上传' : '总量';
}

function formatBytes(value?: number | null) {
  if (value === undefined || value === null || !Number.isFinite(value)) return '--';
  if (value < 1024) return `${Math.round(value)} B`;
  const units = ['KB', 'MB', 'GB', 'TB', 'PB'];
  let amount = value;
  let index = -1;
  do {
    amount /= 1024;
    index += 1;
  } while (amount >= 1024 && index < units.length - 1);
  return `${amount.toFixed(amount >= 10 ? 1 : 2)} ${units[index]}`;
}

function formatMoney(value?: number, currency = 'CNY') {
  if (value === undefined || value === null) return '--';
  return new Intl.NumberFormat('zh-CN', { style: 'currency', currency, maximumFractionDigits: 2 }).format(value);
}

function formatDate(value?: string | null) {
  if (!value) return '--';
  const date = new Date(value);
  return Number.isNaN(date.getTime()) ? '时间未知' : date.toLocaleString('zh-CN', { hour12: false });
}

function formatAxisDate(value: string) {
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return value;
  if (dashboard.value?.trafficTrend.bucket === '1d') {
    return date.toLocaleDateString('zh-CN', { month: '2-digit', day: '2-digit' });
  }
  return date.toLocaleTimeString('zh-CN', { hour12: false, hour: '2-digit', minute: '2-digit' });
}

function relativeTime(value?: string | null) {
  if (!value) return '--';
  const timestamp = new Date(value).getTime();
  if (!Number.isFinite(timestamp)) return '时间未知';
  const minutes = Math.max(0, Math.round((Date.now() - timestamp) / 60_000));
  if (minutes < 1) return '刚刚';
  if (minutes < 60) return `${minutes} 分钟前`;
  const hours = Math.round(minutes / 60);
  if (hours < 24) return `${hours} 小时前`;
  return formatDate(value);
}

function userInitial(name: string) {
  return name.trim().slice(0, 1).toUpperCase() || '用';
}

function dateValueToDay(value: number) {
  const date = new Date(value);
  const month = `${date.getMonth() + 1}`.padStart(2, '0');
  const day = `${date.getDate()}`.padStart(2, '0');
  return `${date.getFullYear()}-${month}-${day}`;
}

function loadParameters(): Api.Central.DashboardQuery | null {
  if (range.value !== 'custom') return { range: range.value };
  if (!customRange.value) return null;
  return {
    range: 'custom',
    from: dateValueToDay(customRange.value[0]),
    to: dateValueToDay(customRange.value[1])
  };
}

async function loadDashboard() {
  const params = loadParameters();
  if (!params) return;
  loading.value = true;
  errorMessage.value = '';
  const { data, error } = await fetchDashboard(params);
  if (error || !data) {
    errorMessage.value = '中央后端暂不可用，无法读取运营数据。';
    dashboard.value = null;
  } else {
    dashboard.value = data;
  }
  loading.value = false;
}

function selectRange(value: Api.Central.DashboardRange) {
  range.value = value;
  if (value !== 'custom') void loadDashboard();
}

function updateCustomRange(value: number | [number, number] | null) {
  if (!Array.isArray(value) || value.length !== 2) {
    customRange.value = null;
    return;
  }
  customRange.value = [value[0], value[1]];
  void loadDashboard();
}

function openNode(nodeID: string) {
  if (nodeID) router.push({ name: 'nodes-detail', params: { id: nodeID } });
}

function openUser(userID: string) {
  if (userID) router.push({ name: 'users', query: { userId: userID } });
}

function buildTrafficOptions(): ECOption {
  const points = trendPoints.value;
  return {
    animation: false,
    grid: { left: 58, right: 18, top: 36, bottom: 28 },
    legend: { top: 0, left: 0, itemWidth: 10, itemHeight: 10, textStyle: { fontSize: 12 } },
    tooltip: {
      trigger: 'axis',
      formatter: (params: any) => {
        const first = Array.isArray(params) ? params[0] : params;
        const point = points[first?.dataIndex ?? 0];
        if (!point) return '';
        const notices = [point.resetDetected ? '检测到基线重置' : '', point.hasGap ? '采样间隔较长' : '']
          .filter(Boolean)
          .join(' · ');
        return `<div>${formatDate(point.time)}</div><div>下载累计：${formatBytes(point.downloadBytes)}</div><div>上传累计：${formatBytes(point.uploadBytes)}</div><div>总计：${formatBytes(point.totalBytes)}</div>${notices ? `<div style="margin-top:4px;color:#d9822b">${notices}</div>` : ''}`;
      }
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: points.map(point => formatAxisDate(point.time)),
      axisTick: { show: false },
      axisLabel: { hideOverlap: true, fontSize: 11 }
    },
    yAxis: {
      type: 'value',
      axisLabel: { formatter: (value: number) => formatBytes(value), fontSize: 11 },
      splitLine: { lineStyle: { type: 'dashed', opacity: 0.35 } }
    },
    series: [
      {
        name: '下载累计',
        type: 'line',
        smooth: 0.2,
        showSymbol: false,
        data: points.map(point => point.downloadBytes),
        lineStyle: { width: 2.5, color: '#646cff' },
        itemStyle: { color: '#646cff' },
        areaStyle: { color: 'rgba(100, 108, 255, 0.12)' }
      },
      {
        name: '上传累计',
        type: 'line',
        smooth: 0.2,
        showSymbol: false,
        data: points.map(point => point.uploadBytes),
        lineStyle: { width: 2, color: '#0f9f9a' },
        itemStyle: { color: '#0f9f9a' }
      }
    ]
  };
}

const { domRef: trendChartRef, updateOptions } = useEcharts(() => buildTrafficOptions());

watch(
  () => dashboard.value?.trafficTrend,
  () => {
    void updateOptions(() => buildTrafficOptions());
  },
  { deep: true }
);

onMounted(loadDashboard);
</script>

<template>
  <NSpace vertical :size="16">
    <NCard :bordered="false" size="small">
      <div class="flex flex-wrap items-start justify-between gap-16px">
        <div>
          <h2 class="m-0 text-20px font-600">运营总览</h2>
          <p class="mb-0 mt-8px text-14px text-gray-500">
            以业务用户 Inbound 为统计单位，快速掌握流量、用户和收入状态。
          </p>
        </div>
        <div class="flex flex-wrap items-center gap-8px">
          <NButtonGroup size="small">
            <NButton
              v-for="item in rangeOptions"
              :key="item.value"
              :type="range === item.value ? 'primary' : 'default'"
              @click="selectRange(item.value)"
            >
              {{ item.label }}
            </NButton>
          </NButtonGroup>
          <NDatePicker
            v-if="range === 'custom'"
            :value="customRange"
            type="daterange"
            clearable
            class="w-250px"
            @update:value="updateCustomRange"
          />
          <DataFreshness :data-at="dashboard?.dataAt" compact />
          <NButton size="small" :loading="loading" @click="loadDashboard">
            <template #icon><icon-mdi-refresh /></template>
            刷新数据
          </NButton>
        </div>
      </div>
    </NCard>

    <NAlert v-if="errorMessage" type="warning" :show-icon="false">{{ errorMessage }}</NAlert>

    <NSpin :show="loading">
      <template v-if="dashboard">
        <div class="operations-kpis">
          <NCard :bordered="false" size="small">
            <NStatistic label="今日业务流量"><TrafficValue :value="dashboard.traffic.todayBytes" /></NStatistic>
            <div class="mt-4px text-12px text-gray-400">线路机用户 Inbound</div>
          </NCard>
          <NCard :bordered="false" size="small">
            <NStatistic label="本月业务流量"><TrafficValue :value="dashboard.traffic.monthBytes" /></NStatistic>
            <div class="mt-4px text-12px text-gray-400">截至 {{ formatDate(dashboard.dataAt) }}</div>
          </NCard>
          <NCard :bordered="false" size="small">
            <NStatistic label="在线节点" :value="dashboard.nodes.online">
              <template #suffix>/ {{ dashboard.nodes.total }}</template>
            </NStatistic>
            <div class="mt-4px text-12px text-gray-400">
              {{ dashboard.nodes.relay }} 个线路机 · {{ dashboard.nodes.landing }} 个落地机
            </div>
          </NCard>
          <NCard :bordered="false" size="small">
            <NStatistic label="有效业务用户" :value="dashboard.users.active" />
            <div class="mt-4px text-12px text-gray-400">按 Client 到期状态同步</div>
          </NCard>
          <NCard :bordered="false" size="small">
            <NStatistic label="7 天内到期" :value="dashboard.expiry.expiring7d" />
            <div class="mt-4px text-12px text-warning">需要跟进续费</div>
          </NCard>
          <NCard :bordered="false" size="small">
            <NStatistic label="待处理事件" :value="dashboard.events.pendingCount" />
            <div class="mt-4px text-12px text-error">需要运营人员关注</div>
          </NCard>
          <NCard :bordered="false" size="small">
            <NStatistic
              label="本月服务期收入"
              :value="formatMoney(dashboard.finance.monthIncome, dashboard.finance.currency)"
            />
            <div class="mt-4px text-12px text-gray-400">服务期折算收入</div>
          </NCard>
          <NCard :bordered="false" size="small">
            <NStatistic
              label="服务期毛利"
              :value="formatMoney(dashboard.finance.grossProfit, dashboard.finance.currency)"
            />
            <div class="mt-4px text-12px text-gray-400">
              {{ grossMargin === null ? '收入为零，暂无毛利率' : `毛利率 ${(grossMargin * 100).toFixed(1)}%` }}
            </div>
          </NCard>
        </div>

        <div class="operations-main-grid mt-16px">
          <NCard :bordered="false" size="small">
            <template #header>
              <div class="font-600">业务流量趋势</div>
              <div class="mt-2px text-12px text-gray-400">{{ trendSubtitle }}</div>
            </template>
            <template #header-extra>
              <span class="text-12px text-gray-400">{{ rangeDescription }}</span>
            </template>
            <NEmpty
              v-if="!hasTrendData"
              description="当前范围暂无足够的流量快照"
              size="small"
              class="h-260px flex-center"
            />
            <div v-else ref="trendChartRef" class="h-260px w-full" />
            <div
              v-if="hasTrendData"
              class="mt-12px flex flex-wrap items-center gap-x-18px gap-y-6px text-12px text-gray-500"
            >
              <span>
                下载
                <strong class="ml-4px text-gray-800 dark:text-gray-100">
                  <TrafficValue :value="dashboard.trafficTrend.summary.downloadBytes" />
                </strong>
              </span>
              <span>
                上传
                <strong class="ml-4px text-gray-800 dark:text-gray-100">
                  <TrafficValue :value="dashboard.trafficTrend.summary.uploadBytes" />
                </strong>
              </span>
              <span class="ml-auto">
                总计
                <strong class="ml-4px text-gray-800 dark:text-gray-100">
                  <TrafficValue :value="dashboard.trafficTrend.summary.totalBytes" />
                </strong>
              </span>
            </div>
            <NSpace v-if="hasTrendReset || hasTrendGap" vertical :size="8" class="mt-12px">
              <NAlert v-if="hasTrendReset" type="warning" size="small">
                检测到流量基线重置，重置点后已按新基线继续统计。
              </NAlert>
              <NAlert v-if="hasTrendGap" type="info" size="small">
                部分采样间隔较长，图表只展示实际采样得到的流量。
              </NAlert>
            </NSpace>
          </NCard>

          <NCard :bordered="false" size="small">
            <template #header>
              <div class="font-600">节点流量排行</div>
              <div class="mt-2px text-12px text-gray-400">线路机用户 Inbound 的业务流量</div>
            </template>
            <template #header-extra>
              <NButtonGroup size="tiny">
                <NButton
                  v-for="item in ['total', 'download', 'upload'] as RankingMetric[]"
                  :key="item"
                  :type="nodeMetric === item ? 'primary' : 'default'"
                  @click="nodeMetric = item"
                >
                  {{ metricLabel(item) }}
                </NButton>
              </NButtonGroup>
            </template>
            <NEmpty
              v-if="!nodeRanking.length"
              description="当前范围暂无节点流量"
              size="small"
              class="h-280px flex-center"
            />
            <div v-else class="space-y-14px">
              <button
                v-for="item in nodeRanking"
                :key="item.nodeId"
                type="button"
                class="operations-rank-row"
                @click="openNode(item.nodeId)"
              >
                <div class="min-w-0 text-left">
                  <div class="truncate text-13px font-600">{{ item.nodeName }}</div>
                  <div class="mt-2px flex items-center gap-4px"><NodeStatusTag :status="item.status" /></div>
                </div>
                <div class="h-8px overflow-hidden rounded bg-gray-100 dark:bg-gray-700">
                  <div
                    class="h-full rounded bg-primary transition-all"
                    :style="{ width: `${(item[metricKey(nodeMetric)] / nodeMaximum) * 100}%` }"
                  />
                </div>
                <span class="text-right text-12px text-gray-500">
                  <TrafficValue :value="item[metricKey(nodeMetric)]" />
                </span>
              </button>
            </div>
            <div
              v-if="nodeRanking.length"
              class="mt-16px flex justify-between border-t border-gray-100 pt-12px text-12px text-gray-400 dark:border-gray-700"
            >
              <span>Top {{ Math.min(5, nodeRanking.length) }} 占业务流量 {{ (topNodeShare * 100).toFixed(1) }}%</span>
              <span>点击节点查看详情</span>
            </div>
          </NCard>
        </div>

        <NCard :bordered="false" size="small" class="mt-16px">
          <template #header>
            <div class="font-600">业务用户流量排行</div>
            <div class="mt-2px text-12px text-gray-400">
              每个线路机用户 Inbound 独立统计，不按 Email 合并，不重复累加 Client。
            </div>
          </template>
          <template #header-extra>
            <NButtonGroup size="tiny">
              <NButton
                v-for="item in ['total', 'download', 'upload'] as RankingMetric[]"
                :key="item"
                :type="userMetric === item ? 'primary' : 'default'"
                @click="userMetric = item"
              >
                {{ metricLabel(item) }}
              </NButton>
            </NButtonGroup>
          </template>
          <NEmpty
            v-if="!userRanking.length"
            description="当前范围暂无业务用户流量"
            size="small"
            class="h-180px flex-center"
          />
          <div v-else class="overflow-x-auto">
            <table class="operations-table min-w-760px w-full">
              <thead>
                <tr>
                  <th>业务用户 / Inbound</th>
                  <th>线路机</th>
                  <th>Client</th>
                  <th>本期下载</th>
                  <th>本期上传</th>
                  <th>本期总量</th>
                  <th>状态</th>
                  <th>最近在线</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="item in userRanking" :key="item.inboundId" @click="openUser(item.userId)">
                  <td>
                    <div class="flex items-center gap-8px">
                      <div
                        class="flex h-28px w-28px flex-none items-center justify-center rounded bg-primary/10 text-12px font-600 text-primary"
                      >
                        {{ userInitial(item.userName) }}
                      </div>
                      <div class="min-w-0">
                        <div class="truncate font-600 text-gray-800 dark:text-gray-100">
                          {{ item.userName || item.inboundTag }}
                        </div>
                        <div class="mt-2px truncate text-11px text-gray-400">
                          {{ item.inboundTag || item.inboundId }}
                        </div>
                      </div>
                    </div>
                  </td>
                  <td>{{ item.nodeName }}</td>
                  <td>{{ item.clientCount }} 台</td>
                  <td><TrafficValue :value="item.downloadBytes" /></td>
                  <td><TrafficValue :value="item.uploadBytes" /></td>
                  <td class="font-600 text-gray-800 dark:text-gray-100"><TrafficValue :value="item.totalBytes" /></td>
                  <td><UserStatusTag :status="item.status" :expires-at="item.expiresAt" /></td>
                  <td>{{ relativeTime(item.lastActivityAt) }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </NCard>

        <div class="operations-bottom-grid mt-16px">
          <NCard :bordered="false" size="small">
            <template #header>
              <div class="font-600">用户到期分布</div>
              <div class="mt-2px text-12px text-gray-400">按 Client 到期时间判断业务状态</div>
            </template>
            <div class="space-y-14px">
              <div
                v-for="item in expiryItems"
                :key="item.label"
                class="grid grid-cols-[100px_minmax(0,1fr)_42px] items-center gap-10px text-12px"
              >
                <span class="text-gray-500">{{ item.label }}</span>
                <div class="h-7px overflow-hidden rounded bg-gray-100 dark:bg-gray-700">
                  <div
                    class="h-full rounded transition-all"
                    :class="
                      item.color === 'warning'
                        ? 'bg-warning'
                        : item.color === 'error'
                          ? 'bg-error'
                          : item.color === 'info'
                            ? 'bg-info'
                            : 'bg-primary'
                    "
                    :style="{ width: `${(item.value / expiryMaximum) * 100}%` }"
                  />
                </div>
                <span class="text-right font-600">{{ item.value }}</span>
              </div>
            </div>
            <div class="mt-16px border-t border-gray-100 pt-12px text-11px text-gray-400 dark:border-gray-700">
              已逻辑删除或没有线路机关联的用户不进入运营统计。
            </div>
          </NCard>

          <NCard :bordered="false" size="small">
            <template #header>
              <div class="font-600">待处理事件</div>
              <div class="mt-2px text-12px text-gray-400">需要运营人员关注的变化</div>
            </template>
            <template #header-extra>
              <NButton text type="primary" size="small" @click="router.push({ name: 'events' })">消息中心</NButton>
            </template>
            <NEmpty
              v-if="!dashboard.events.items.length"
              description="暂无待处理事件"
              size="small"
              class="h-168px flex-center"
            />
            <div v-else class="divide-y divide-gray-100 dark:divide-gray-700">
              <button
                v-for="item in dashboard.events.items"
                :key="item.id"
                type="button"
                class="operations-event-row"
                @click="router.push({ name: 'events' })"
              >
                <span
                  class="mt-5px h-7px w-7px flex-none rounded-full"
                  :class="
                    item.severity === 'error' ? 'bg-error' : item.severity === 'warning' ? 'bg-warning' : 'bg-primary'
                  "
                />
                <span class="min-w-0 text-left">
                  <span class="block truncate text-12px text-gray-700 dark:text-gray-200">
                    {{ item.title || item.message }}
                  </span>
                  <span class="mt-2px block truncate text-11px text-gray-400">
                    {{ item.category }}{{ item.nodeName ? ` · ${item.nodeName}` : '' }}
                  </span>
                </span>
                <span class="text-11px text-gray-400">{{ relativeTime(item.occurredAt) }}</span>
              </button>
            </div>
          </NCard>

          <NCard :bordered="false" size="small">
            <template #header>
              <div class="font-600">本月财务</div>
              <div class="mt-2px text-12px text-gray-400">已确认收款保留在财务记录中</div>
            </template>
            <template #header-extra>
              <NButton text type="primary" size="small" @click="router.push({ name: 'finance' })">财务明细</NButton>
            </template>
            <div class="space-y-13px text-13px">
              <div class="flex items-center justify-between gap-12px">
                <span class="text-gray-500">服务期折算收入</span>
                <strong>{{ formatMoney(dashboard.finance.monthIncome, dashboard.finance.currency) }}</strong>
              </div>
              <div class="flex items-center justify-between gap-12px">
                <span class="text-gray-500">已确认实收</span>
                <strong>{{ formatMoney(dashboard.finance.cashIncome, dashboard.finance.currency) }}</strong>
              </div>
              <div class="flex items-center justify-between gap-12px">
                <span class="text-gray-500">本月成本</span>
                <strong>{{ formatMoney(dashboard.finance.monthCost, dashboard.finance.currency) }}</strong>
              </div>
              <div class="border-t border-gray-100 pt-13px dark:border-gray-700">
                <div class="flex items-center justify-between gap-12px">
                  <span class="text-gray-500">服务期毛利</span>
                  <strong class="text-success">
                    {{ formatMoney(dashboard.finance.grossProfit, dashboard.finance.currency) }}
                  </strong>
                </div>
              </div>
            </div>
            <div class="mt-16px border-t border-gray-100 pt-12px text-11px text-gray-400 dark:border-gray-700">
              财务数据独立于用户清理逻辑。
            </div>
          </NCard>
        </div>
      </template>
      <NEmpty v-else-if="!loading && !errorMessage" description="暂无运营数据" class="min-h-280px flex-center" />
    </NSpin>
  </NSpace>
</template>

<style scoped>
.operations-kpis {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}
.operations-main-grid {
  display: grid;
  grid-template-columns: minmax(0, 1fr);
  gap: 16px;
}
.operations-bottom-grid {
  display: grid;
  grid-template-columns: minmax(0, 1fr);
  gap: 16px;
}
.operations-rank-row {
  display: grid;
  width: 100%;
  grid-template-columns: minmax(108px, 1.25fr) minmax(120px, 2fr) auto;
  align-items: center;
  gap: 12px;
  border: 0;
  background: transparent;
  padding: 0;
  color: inherit;
  cursor: pointer;
}
.operations-rank-row:hover .font-600 {
  color: var(--n-primary-color);
}
.operations-table {
  border-collapse: collapse;
}
.operations-table th,
.operations-table td {
  border-bottom: 1px solid rgb(243 244 246);
  padding: 11px 12px;
  text-align: left;
  white-space: nowrap;
  font-size: 12px;
}
.operations-table th {
  color: rgb(156 163 175);
  font-weight: 600;
  background: rgb(249 250 251);
}
.operations-table tbody tr {
  cursor: pointer;
  transition: background-color 0.15s ease;
}
.operations-table tbody tr:hover {
  background: rgb(249 250 255);
}
.operations-event-row {
  display: grid;
  width: 100%;
  grid-template-columns: 8px minmax(0, 1fr) auto;
  align-items: start;
  gap: 9px;
  border: 0;
  background: transparent;
  padding: 9px 0;
  color: inherit;
  cursor: pointer;
}
.operations-event-row:hover .text-gray-700 {
  color: var(--n-primary-color);
}
.dark .operations-table th {
  background: rgb(31 41 55);
}
.dark .operations-table th,
.dark .operations-table td {
  border-color: rgb(55 65 81);
}
.dark .operations-table tbody tr:hover {
  background: rgb(31 41 55);
}
@media (min-width: 768px) {
  .operations-kpis {
    grid-template-columns: repeat(4, minmax(0, 1fr));
  }
  .operations-bottom-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
  .operations-bottom-grid > :last-child {
    grid-column: 1 / -1;
  }
}
@media (min-width: 1100px) {
  .operations-main-grid {
    grid-template-columns: minmax(0, 1.7fr) minmax(360px, 1fr);
  }
  .operations-bottom-grid {
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }
  .operations-bottom-grid > :last-child {
    grid-column: auto;
  }
}
@media (min-width: 1320px) {
  .operations-kpis {
    grid-template-columns: repeat(8, minmax(0, 1fr));
  }
}
@media (max-width: 640px) {
  .operations-rank-row {
    grid-template-columns: minmax(92px, 1fr) minmax(68px, 1fr) auto;
    gap: 8px;
  }
}
</style>
