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
  return items.sort((left, right) => right.totalBytes - left.totalBytes).slice(0, 10);
});

const nodeMaximum = computed(() => Math.max(1, ...nodeRanking.value.map(item => item[metricKey(nodeMetric.value)])));

const nodeMetricOptions: Array<{ value: RankingMetric; label: string }> = [
  { value: 'total', label: '总量' },
  { value: 'download', label: '下载' },
  { value: 'upload', label: '上传' }
];

const nodeRows = computed(() =>
  nodeRanking.value.map(item => ({
    nodeId: item.nodeId,
    nodeName: item.nodeName,
    status: item.status,
    bytes: item[metricKey(nodeMetric.value)]
  }))
);

const expiryItems = computed(() => {
  const expiry = dashboard.value?.expiry;
  if (!expiry) return [];
  return [
    { label: '30 天后 / 长期', value: expiry.normalAfter30d, color: 'brand' },
    { label: '8–30 天到期', value: expiry.expiring30d, color: 'cyan' },
    { label: '7 天内到期', value: expiry.expiring7d, color: 'warn' },
    { label: '已到期', value: expiry.expired, color: 'danger' }
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
  return `${bucket}累计统计 · 全部线路机 · 上传与下载均来自用户 Inbound`;
});

const topNodeShare = computed(() => {
  const total = dashboard.value?.traffic.totalBytes || 0;
  if (!total) return 0;
  return nodeRanking.value.slice(0, 5).reduce((sum, item) => sum + item.totalBytes, 0) / total;
});

const pendingEvents = computed(() => (dashboard.value?.events.items || []).slice(0, 5));

const syncRows = computed(() => {
  const items = [...(dashboard.value?.nodeTrafficRanking || [])];
  return items.sort((left, right) => {
    const leftTime = left.dataAt ? new Date(left.dataAt).getTime() : 0;
    const rightTime = right.dataAt ? new Date(right.dataAt).getTime() : 0;
    return rightTime - leftTime;
  });
});

const offlineNodes = computed(() => {
  const nodes = dashboard.value?.nodes;
  if (!nodes) return 0;
  return Math.max(0, nodes.total - nodes.online);
});

function metricKey(metric: RankingMetric): 'totalBytes' | 'downloadBytes' | 'uploadBytes' {
  if (metric === 'download') return 'downloadBytes';
  if (metric === 'upload') return 'uploadBytes';
  return 'totalBytes';
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

function eventIcon(severity: string) {
  if (severity === 'error') return 'i-carbon:warning-alt';
  if (severity === 'warning') return 'i-carbon:warning';
  return 'i-carbon:information';
}

function buildTrafficOptions(): ECOption {
  const points = trendPoints.value;
  return {
    animation: false,
    grid: { left: 62, right: 18, top: 34, bottom: 26 },
    legend: {
      top: 0,
      left: 0,
      itemWidth: 10,
      itemHeight: 10,
      textStyle: { fontSize: 12, color: '#a2a9b4' }
    },
    tooltip: {
      trigger: 'axis',
      backgroundColor: '#191b20',
      borderColor: '#2e323a',
      textStyle: { color: '#f2f3f5', fontSize: 12 },
      formatter: (params: any) => {
        const first = Array.isArray(params) ? params[0] : params;
        const point = points[first?.dataIndex ?? 0];
        if (!point) return '';
        const notices = [
          point.resetDetected ? '检测到基线重置' : '',
          point.hasGap ? '采样间隔较长' : ''
        ]
          .filter(Boolean)
          .join(' · ');
        return `<div>${formatDate(point.time)}</div><div>下载累计：${formatBytes(
          point.downloadBytes
        )}</div><div>上传累计：${formatBytes(point.uploadBytes)}</div><div>总计：${formatBytes(
          point.totalBytes
        )}</div>${notices ? `<div style="margin-top:4px;color:#f5b93d">${notices}</div>` : ''}`;
      }
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: points.map(point => formatAxisDate(point.time)),
      axisTick: { show: false },
      axisLine: { lineStyle: { color: '#23262c' } },
      axisLabel: { hideOverlap: true, fontSize: 11, color: '#8f96a1' }
    },
    yAxis: {
      type: 'value',
      axisLabel: { formatter: (value: number) => formatBytes(value), fontSize: 11, color: '#8f96a1' },
      splitLine: { lineStyle: { type: 'dashed', opacity: 0.35, color: '#23262c' } }
    },
    series: [
      {
        name: '下载累计',
        type: 'line',
        smooth: 0.2,
        showSymbol: false,
        data: points.map(point => point.downloadBytes),
        lineStyle: { width: 2.5, color: '#4c6fff' },
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
        name: '上传累计',
        type: 'line',
        smooth: 0.2,
        showSymbol: false,
        data: points.map(point => point.uploadBytes),
        lineStyle: { width: 2, color: '#34d3e0' },
        itemStyle: { color: '#34d3e0' }
      }
    ]
  };
}

const { domRef: trendChartRef, updateOptions } = useEcharts(() => buildTrafficOptions());

watch(
  () => dashboard.value?.trafficTrend,
  () => {
    if (!trendChartRef.value) return;
    void updateOptions(() => buildTrafficOptions());
  },
  { deep: true }
);

onMounted(loadDashboard);
</script>

<template>
  <div class="dash">
    <!-- 页头 -->
    <header class="dash-head">
      <div class="dash-head-main">
        <div class="dash-eyebrow">运营 · 运营总览 · 业务日 UTC+8</div>
        <h1 class="dash-title">运营总览</h1>
        <p class="dash-desc">
          以业务用户 Inbound 为统计单位，掌握流量、用户与节点状态。财务口径与明细见「成本与收入」，本页不重复财务指标。
        </p>
      </div>
      <div class="dash-actions">
        <div class="dash-range">
          <button
            v-for="item in rangeOptions"
            :key="item.value"
            class="dash-range-btn"
            :class="{ 'is-active': range === item.value }"
            @click="selectRange(item.value)"
          >
            {{ item.label }}
          </button>
        </div>
        <NDatePicker
          v-if="range === 'custom'"
          :value="customRange"
          type="daterange"
          clearable
          size="small"
          class="w-250px"
          @update:value="updateCustomRange"
        />
        <DataFreshness :data-at="dashboard?.dataAt" compact />
        <button class="dash-btn" :disabled="loading" @click="loadDashboard">
          <SvgIcon icon="mdi:refresh" class="text-15px" />
          刷新
        </button>
      </div>
    </header>

    <NAlert v-if="errorMessage" type="warning" :show-icon="false" class="mb-12px">{{ errorMessage }}</NAlert>

    <NSpin :show="loading">
      <template v-if="dashboard">
        <div class="dash-body">
        <!-- KPI -->
        <section class="dash-kpis">
          <div class="dash-kpi">
            <div class="dash-kpi-top">
              <span class="dash-kpi-label">今日业务流量</span>
              <span class="dash-kpi-icon is-blue"><SvgIcon icon="mdi:chart-line" /></span>
            </div>
            <div class="dash-kpi-value"><TrafficValue :value="dashboard.traffic.todayBytes" /></div>
            <div class="dash-kpi-foot">线路机用户 Inbound</div>
          </div>

          <div class="dash-kpi">
            <div class="dash-kpi-top">
              <span class="dash-kpi-label">本月业务流量</span>
              <span class="dash-kpi-icon is-cyan"><SvgIcon icon="mdi:layers-outline" /></span>
            </div>
            <div class="dash-kpi-value"><TrafficValue :value="dashboard.traffic.monthBytes" /></div>
            <div class="dash-kpi-foot">按最近同步时间统计</div>
          </div>

          <div class="dash-kpi">
            <div class="dash-kpi-top">
              <span class="dash-kpi-label">在线节点</span>
              <span class="dash-kpi-icon is-green"><SvgIcon icon="mdi:server-network" /></span>
            </div>
            <div class="dash-kpi-value">
              {{ dashboard.nodes.online }}<span class="dash-kpi-unit">/ {{ dashboard.nodes.total }}</span>
            </div>
            <div class="dash-kpi-foot">{{ dashboard.nodes.relay }} 个线路机 · {{ dashboard.nodes.landing }} 个落地机</div>
          </div>

          <div class="dash-kpi">
            <div class="dash-kpi-top">
              <span class="dash-kpi-label">有效业务用户</span>
              <span class="dash-kpi-icon is-violet"><SvgIcon icon="mdi:account-group-outline" /></span>
            </div>
            <div class="dash-kpi-value">{{ dashboard.users.active }}<span class="dash-kpi-unit">户</span></div>
            <div class="dash-kpi-foot">按 Client 到期状态同步</div>
          </div>

          <div class="dash-kpi">
            <div class="dash-kpi-top">
              <span class="dash-kpi-label">7 天内到期</span>
              <span class="dash-kpi-icon is-amber"><SvgIcon icon="mdi:clock-outline" /></span>
            </div>
            <div class="dash-kpi-value">{{ dashboard.expiry.expiring7d }}<span class="dash-kpi-unit">户</span></div>
            <div class="dash-kpi-foot"><span class="dash-chip is-warn">需要跟进续费</span></div>
          </div>

          <div class="dash-kpi">
            <div class="dash-kpi-top">
              <span class="dash-kpi-label">待处理事件</span>
              <span class="dash-kpi-icon is-red"><SvgIcon icon="mdi:bell-alert-outline" /></span>
            </div>
            <div class="dash-kpi-value">{{ dashboard.events.pendingCount }}<span class="dash-kpi-unit">条</span></div>
            <div class="dash-kpi-foot"><span class="dash-chip is-danger">需要运营处理</span></div>
          </div>
        </section>

        <!-- 流量趋势 + 节点排行 -->
        <div class="dash-split">
          <section class="dash-card dash-card-chart">
            <div class="dash-card-head">
              <div>
                <h3>业务流量趋势</h3>
                <span>{{ trendSubtitle }}</span>
              </div>
              <div class="dash-card-tools">
                <span class="dash-chip is-ghost">{{ rangeDescription }}</span>
              </div>
            </div>
            <div class="dash-card-body">
              <div class="dash-chart-wrap">
                <div ref="trendChartRef" class="dash-chart"></div>
                <div v-if="!hasTrendData" class="dash-chart-empty">当前范围暂无足够的流量快照</div>
              </div>
              <div class="dash-summary">
                <div>
                  下载
                  <b><TrafficValue :value="dashboard.trafficTrend.summary.downloadBytes" /></b>
                </div>
                <i></i>
                <div>
                  上传
                  <b><TrafficValue :value="dashboard.trafficTrend.summary.uploadBytes" /></b>
                </div>
                <i></i>
                <div>
                  总计
                  <b><TrafficValue :value="dashboard.trafficTrend.summary.totalBytes" /></b>
                </div>
                <span v-if="hasTrendReset || hasTrendGap" class="dash-summary-notice">
                  {{ hasTrendReset ? '检测到流量基线重置，重置点后按新基线统计' : '' }}
                  {{ hasTrendReset && hasTrendGap ? ' · ' : '' }}
                  {{ hasTrendGap ? '部分采样间隔较长，仅展示实际采样' : '' }}
                </span>
              </div>
            </div>
          </section>

          <section class="dash-card">
            <div class="dash-card-head">
              <div>
                <h3>节点流量排行</h3>
                <span>按线路机统计 · 落地机不计入</span>
              </div>
              <div class="dash-seg">
                <button
                  v-for="option in nodeMetricOptions"
                  :key="option.value"
                  :class="{ 'is-active': nodeMetric === option.value }"
                  @click="nodeMetric = option.value"
                >
                  {{ option.label }}
                </button>
              </div>
            </div>
            <div class="dash-card-body dash-rank-list">
              <div v-for="item in nodeRows" :key="item.nodeId" class="dash-rank-row" @click="openNode(item.nodeId)">
                <span class="dash-rank-name">{{ item.nodeName }}</span>
                <NodeStatusTag :status="item.status" />
                <div class="dash-rank-bar">
                  <i :style="{ width: `${Math.max((item.bytes / nodeMaximum) * 100, 2)}%` }"></i>
                </div>
                <span class="dash-rank-val"><TrafficValue :value="item.bytes" /></span>
              </div>
              <div v-if="!nodeRanking.length" class="dash-empty-small">当前范围暂无节点流量</div>
              <div class="dash-rank-foot">
                TOP 5 占业务流量 <b>{{ (topNodeShare * 100).toFixed(1) }}%</b> · 点击节点查看详情
              </div>
            </div>
          </section>
        </div>

        <!-- 业务用户流量排行 -->
        <section class="dash-card">
          <div class="dash-card-head">
            <div>
              <h3>业务用户流量排行</h3>
              <span>本期 · 全部线路机 · 按用户 Inbound 汇总，点击行查看用户详情</span>
            </div>
            <div class="dash-card-tools">
              <span class="dash-chip is-ghost">TOP 10</span>
              <RouterLink to="/users" class="dash-link">全部用户 ›</RouterLink>
            </div>
          </div>
          <div class="dash-table-wrap">
            <table class="dash-table">
              <thead>
                <tr>
                  <th style="width: 42px">#</th>
                  <th>业务用户 / Inbound</th>
                  <th>线路机</th>
                  <th class="is-num">Client</th>
                  <th class="is-num">本期下载</th>
                  <th class="is-num">本期上传</th>
                  <th class="is-num">本期总量</th>
                  <th>状态</th>
                  <th>最近在线</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(item, index) in userRanking" :key="item.userId" @click="openUser(item.userId)">
                  <td>
                    <span class="dash-rank-num" :class="{ 'is-top': index < 3 }">{{ index + 1 }}</span>
                  </td>
                  <td class="is-strong">
                    {{ item.userName }}
                    <span class="dash-sub">{{ item.nodeName }} · {{ item.inboundTag || item.inboundId }}</span>
                  </td>
                  <td><span class="dash-node-tag">{{ item.nodeName }}</span></td>
                  <td class="is-num">{{ item.clientCount }} 台</td>
                  <td class="is-num"><TrafficValue :value="item.downloadBytes" /></td>
                  <td class="is-num"><TrafficValue :value="item.uploadBytes" /></td>
                  <td class="is-num is-strong"><TrafficValue :value="item.totalBytes" /></td>
                  <td><UserStatusTag :status="item.status" /></td>
                  <td class="dash-dim">{{ relativeTime(item.lastActivityAt) }}</td>
                </tr>
              </tbody>
            </table>
            <div v-if="!userRanking.length" class="dash-empty-small">当前范围暂无业务用户流量</div>
          </div>
        </section>

        <!-- 三宫格 -->
        <div class="dash-three">
          <section class="dash-card">
            <div class="dash-card-head">
              <h3>用户到期分布</h3>
              <span>按业务日实时刷新</span>
            </div>
            <div class="dash-card-body">
              <div v-for="item in expiryItems" :key="item.label" class="dash-expiry-row">
                <span class="dash-expiry-label">{{ item.label }}</span>
                <div class="dash-expiry-bar">
                  <i :class="`is-${item.color}`" :style="{ width: `${(item.value / expiryMaximum) * 100}%` }"></i>
                </div>
                <b>{{ item.value }}</b>
              </div>
              <p class="dash-note">已逻辑删除或没有线路机关联的用户不进入运营统计；条宽为相对最大值，非占比。</p>
            </div>
          </section>

          <section class="dash-card">
            <div class="dash-card-head">
              <h3>待处理事件</h3>
              <span class="dash-chip is-danger">{{ dashboard.events.pendingCount }} 条</span>
              <div class="dash-card-tools">
                <RouterLink to="/events" class="dash-link">消息中心 ›</RouterLink>
              </div>
            </div>
            <div class="dash-card-body dash-events">
              <div v-for="event in pendingEvents" :key="event.id" class="dash-event-row">
                <span class="dash-event-ico" :class="`is-${event.severity}`">
                  <SvgIcon :icon="eventIcon(event.severity)" />
                </span>
                <div class="dash-event-main">
                  <b>{{ event.title || event.message }}</b>
                  <span>{{ event.nodeName || '系统' }} · {{ relativeTime(event.occurredAt) }}</span>
                </div>
              </div>
              <div v-if="!pendingEvents.length" class="dash-empty-small">暂无待处理事件</div>
            </div>
          </section>

          <section class="dash-card">
            <div class="dash-card-head">
              <h3>同步状态</h3>
              <span class="dash-chip is-ok">在线 {{ dashboard.nodes.online }} / {{ dashboard.nodes.total }}</span>
            </div>
            <div class="dash-card-body">
              <div v-for="item in syncRows" :key="item.nodeId" class="dash-sync-row">
                <span class="dash-sync-name">{{ item.nodeName }}</span>
                <NodeStatusTag :status="item.status" />
                <span class="dash-sync-time">{{ relativeTime(item.dataAt) }}</span>
              </div>
              <div class="dash-sync-foot">
                <span>数据新鲜度取最近一次成功同步</span>
                <span v-if="offlineNodes" class="dash-chip is-danger">离线 {{ offlineNodes }} 台</span>
              </div>
            </div>
          </section>
        </div>
        </div>
      </template>

      <div v-else-if="!loading" class="dash-empty">
        <NEmpty description="暂无运营数据" />
      </div>
    </NSpin>
  </div>
</template>

<style scoped>
.dash {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.dash-body {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

/* ---------- 页头 ---------- */
.dash-head {
  display: flex;
  align-items: flex-end;
  gap: 18px;
  flex-wrap: wrap;
}

.dash-eyebrow {
  font-size: 11.5px;
  letter-spacing: 0.12em;
  color: #7f8792;
  margin-bottom: 6px;
}

.dash-title {
  margin: 0;
  font-size: 23px;
  font-weight: 600;
  letter-spacing: 0.01em;
  color: #f2f3f5;
}

.dash-desc {
  margin: 7px 0 0;
  max-width: 780px;
  font-size: 13px;
  line-height: 1.72;
  color: #a2a9b4;
}

.dash-actions {
  margin-left: auto;
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
  padding-bottom: 2px;
}

.dash-range {
  display: inline-flex;
  gap: 2px;
  padding: 3px;
  border: 1px solid #2a2e35;
  border-radius: 9px;
  background: #14161a;
}

.dash-range-btn {
  height: 27px;
  padding: 0 11px;
  border: 0;
  border-radius: 7px;
  background: transparent;
  color: #8f96a1;
  font-size: 12.5px;
  cursor: pointer;
  transition: background 0.15s ease, color 0.15s ease;
}

.dash-range-btn:hover {
  color: #f2f3f5;
}

.dash-range-btn.is-active {
  background: #2b3038;
  color: #fff;
}

.dash-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  height: 32px;
  padding: 0 12px;
  border: 1px solid #2a2e35;
  border-radius: 9px;
  background: #1b1e24;
  color: #e6e8ec;
  font-size: 12.5px;
  cursor: pointer;
}

.dash-btn:hover {
  background: #232730;
  border-color: #363b44;
}

/* ---------- KPI ---------- */
.dash-kpis {
  display: grid;
  grid-template-columns: repeat(6, minmax(0, 1fr));
  gap: 12px;
}

.dash-kpi {
  padding: 13px 14px 12px;
  border: 1px solid #23262c;
  border-radius: 14px;
  background: #141519;
  min-width: 0;
}

.dash-kpi-top {
  display: flex;
  align-items: center;
  gap: 8px;
}

.dash-kpi-label {
  font-size: 12.5px;
  line-height: 1.4;
  color: #b0b6c1;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.dash-kpi-icon {
  display: inline-grid;
  place-items: center;
  width: 24px;
  height: 24px;
  margin-left: auto;
  flex: none;
  border-radius: 8px;
  font-size: 13px;
  background: #1d2026;
  border: 1px solid #2a2e35;
  color: #a8adb8;
}

.dash-kpi-icon.is-blue {
  background: rgba(76, 111, 255, 0.16);
  border-color: rgba(76, 111, 255, 0.36);
  color: #8fa8ff;
}

.dash-kpi-icon.is-cyan {
  background: rgba(52, 211, 224, 0.14);
  border-color: rgba(52, 211, 224, 0.34);
  color: #6fe0e9;
}

.dash-kpi-icon.is-green {
  background: rgba(46, 212, 122, 0.14);
  border-color: rgba(46, 212, 122, 0.34);
  color: #5ce49b;
}

.dash-kpi-icon.is-violet {
  background: rgba(167, 139, 250, 0.16);
  border-color: rgba(167, 139, 250, 0.36);
  color: #c0a9ff;
}

.dash-kpi-icon.is-amber {
  background: rgba(245, 185, 61, 0.14);
  border-color: rgba(245, 185, 61, 0.34);
  color: #f7ca6b;
}

.dash-kpi-icon.is-red {
  background: rgba(243, 101, 111, 0.14);
  border-color: rgba(243, 101, 111, 0.34);
  color: #ff9ba2;
}

.dash-kpi-value {
  display: flex;
  align-items: baseline;
  gap: 5px;
  margin-top: 8px;
  font-size: 22px;
  line-height: 1.2;
  font-weight: 620;
  letter-spacing: -0.01em;
  font-variant-numeric: tabular-nums;
  color: #fff;
}

.dash-kpi-unit {
  font-size: 12.5px;
  font-weight: 500;
  color: #8f96a1;
}

.dash-kpi-foot {
  margin-top: 7px;
  font-size: 12px;
  line-height: 1.5;
  color: #98a0ab;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* ---------- 卡片 ---------- */
.dash-card {
  border: 1px solid #23262c;
  border-radius: 14px;
  background: #141519;
  min-width: 0;
}

.dash-card-head {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 16px 18px;
  border-bottom: 1px solid #1c1f24;
}

.dash-card-head h3 {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
  color: #f2f3f5;
  white-space: nowrap;
}

.dash-card-head > div > span,
.dash-card-head > span {
  font-size: 12px;
  color: #98a0ab;
}

.dash-card-head > div:first-child {
  display: flex;
  flex-direction: column;
  gap: 3px;
  min-width: 0;
}

.dash-card-tools {
  margin-left: auto;
  display: flex;
  align-items: center;
  gap: 8px;
}

.dash-card-body {
  padding: 16px 18px 18px;
}

/* ---------- 趋势 ---------- */
.dash-split {
  display: grid;
  grid-template-columns: minmax(0, 1.62fr) minmax(0, 1fr);
  gap: 15px;
}

.dash-card-chart {
  display: flex;
  flex-direction: column;
}

.dash-chart-wrap {
  position: relative;
  flex: 1;
  min-height: 280px;
}

.dash-chart {
  width: 100%;
  height: 100%;
  min-height: 280px;
}

.dash-chart-empty {
  position: absolute;
  inset: 0;
  display: grid;
  place-items: center;
  font-size: 12.5px;
  color: #8f96a1;
}

.dash-summary {
  display: flex;
  align-items: center;
  gap: 22px;
  margin-top: 10px;
  padding-top: 12px;
  border-top: 1px solid #1c1f24;
  font-size: 12px;
  color: #98a0ab;
}

.dash-summary b {
  display: block;
  margin-top: 2px;
  font-size: 14px;
  font-weight: 600;
  color: #f2f3f5;
  font-variant-numeric: tabular-nums;
}

.dash-summary i {
  width: 1px;
  height: 26px;
  background: #1c1f24;
}

.dash-summary-notice {
  margin-left: auto;
  font-size: 11.5px;
  color: #f7ca6b;
}

/* ---------- 排行 ---------- */
.dash-rank-list {
  padding-top: 6px;
}

.dash-rank-row {
  display: flex;
  align-items: center;
  gap: 9px;
  padding: 8px 0;
  cursor: pointer;
  border-bottom: 1px solid #1c1f24;
}

.dash-rank-row:last-of-type {
  border-bottom: 0;
}

.dash-rank-row:hover .dash-rank-name {
  color: #fff;
}

.dash-rank-name {
  width: 76px;
  flex: none;
  font-size: 12.5px;
  font-weight: 550;
  color: #e6e8ec;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.dash-rank-bar {
  flex: 1;
  height: 7px;
  border-radius: 99px;
  background: #24272e;
  overflow: hidden;
}

.dash-rank-bar i {
  display: block;
  height: 100%;
  border-radius: 99px;
  background: linear-gradient(90deg, #3f5ee8, #5b7dff);
}

.dash-rank-val {
  width: 68px;
  flex: none;
  text-align: right;
  font-size: 12px;
  font-weight: 550;
  color: #f2f3f5;
  font-variant-numeric: tabular-nums;
}

.dash-rank-foot {
  margin-top: 10px;
  padding-top: 10px;
  border-top: 1px solid #1c1f24;
  font-size: 12px;
  color: #8f96a1;
}

.dash-rank-foot b {
  color: #f2f3f5;
}

/* ---------- 表格 ---------- */
.dash-table-wrap {
  overflow-x: auto;
}

.dash-table {
  width: 100%;
  border-collapse: separate;
  border-spacing: 0;
  font-size: 13px;
}

.dash-table thead th {
  padding: 10px 14px;
  border-bottom: 1px solid #23262c;
  font-size: 11.5px;
  font-weight: 600;
  letter-spacing: 0.06em;
  text-align: left;
  white-space: nowrap;
  color: #a2a9b4;
}

.dash-table thead th.is-num {
  text-align: right;
}

.dash-table tbody td {
  height: 46px;
  padding: 0 14px;
  border-bottom: 1px solid #1c1f24;
  color: #bcc2cc;
  white-space: nowrap;
}

.dash-table tbody tr {
  cursor: pointer;
}

.dash-table tbody tr:hover td {
  background: #191c21;
}

.dash-table tbody tr:last-child td {
  border-bottom: 0;
}

.dash-table .is-num {
  text-align: right;
  font-variant-numeric: tabular-nums;
}

.dash-table .is-strong {
  color: #f2f3f5;
  font-weight: 500;
}

.dash-sub {
  display: block;
  margin-top: 2px;
  font-size: 11.5px;
  color: #98a0ab;
  font-weight: 400;
}

.dash-dim {
  color: #8f96a1;
}

.dash-node-tag {
  display: inline-flex;
  align-items: center;
  height: 21px;
  padding: 0 7px;
  border: 1px solid rgba(76, 111, 255, 0.34);
  border-radius: 6px;
  background: rgba(76, 111, 255, 0.14);
  font-size: 11.5px;
  color: #9db2ff;
}

.dash-rank-num {
  display: inline-grid;
  place-items: center;
  width: 22px;
  height: 22px;
  border-radius: 7px;
  border: 1px solid #2a2e35;
  background: #1d2026;
  font-size: 11px;
  font-weight: 650;
  font-variant-numeric: tabular-nums;
  color: #8f96a1;
}

.dash-rank-num.is-top {
  border-color: rgba(76, 111, 255, 0.5);
  background: linear-gradient(180deg, #5b7dff, #3f5ee8);
  color: #fff;
}

/* ---------- 三宫格 ---------- */
.dash-three {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 15px;
}

.dash-expiry-row {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 7px 0;
}

.dash-expiry-label {
  width: 96px;
  flex: none;
  font-size: 12px;
  color: #bcc2cc;
}

.dash-expiry-bar {
  flex: 1;
  height: 7px;
  border-radius: 99px;
  background: #24272e;
  overflow: hidden;
}

.dash-expiry-bar i {
  display: block;
  height: 100%;
  border-radius: 99px;
  background: linear-gradient(90deg, #3f5ee8, #5b7dff);
}

.dash-expiry-bar i.is-cyan {
  background: linear-gradient(90deg, #16828c, #34d3e0);
}

.dash-expiry-bar i.is-warn {
  background: linear-gradient(90deg, #b07f22, #f5b93d);
}

.dash-expiry-bar i.is-danger {
  background: linear-gradient(90deg, #b8434c, #f3656f);
}

.dash-expiry-row b {
  width: 34px;
  flex: none;
  text-align: right;
  font-size: 12.5px;
  color: #f2f3f5;
  font-variant-numeric: tabular-nums;
}

.dash-note {
  margin: 10px 0 0;
  font-size: 11.5px;
  line-height: 1.7;
  color: #8f96a1;
}

.dash-events {
  padding-top: 6px;
}

.dash-event-row {
  display: flex;
  gap: 10px;
  padding: 8px 0;
  border-bottom: 1px solid #1c1f24;
}

.dash-event-row:last-of-type {
  border-bottom: 0;
}

.dash-event-ico {
  display: inline-grid;
  place-items: center;
  width: 28px;
  height: 28px;
  flex: none;
  border-radius: 9px;
  font-size: 14px;
  background: #1d2026;
  border: 1px solid #2a2e35;
  color: #a8adb8;
}

.dash-event-ico.is-error {
  background: rgba(243, 101, 111, 0.12);
  border-color: rgba(243, 101, 111, 0.3);
  color: #ff9ba2;
}

.dash-event-ico.is-warning {
  background: rgba(245, 185, 61, 0.12);
  border-color: rgba(245, 185, 61, 0.3);
  color: #f7ca6b;
}

.dash-event-main {
  min-width: 0;
}

.dash-event-main b {
  display: block;
  font-size: 12.5px;
  font-weight: 500;
  color: #f2f3f5;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.dash-event-main span {
  display: block;
  margin-top: 2px;
  font-size: 11.5px;
  color: #98a0ab;
}

.dash-sync-row {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 0;
  border-bottom: 1px solid #1c1f24;
  font-size: 12px;
}

.dash-sync-row:last-of-type {
  border-bottom: 0;
}

.dash-sync-name {
  width: 70px;
  flex: none;
  color: #e6e8ec;
  font-weight: 500;
}

.dash-sync-time {
  margin-left: auto;
  color: #98a0ab;
}

.dash-sync-foot {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 10px;
  padding-top: 10px;
  border-top: 1px solid #1c1f24;
  font-size: 11.5px;
  color: #8f96a1;
}

/* ---------- 通用小件 ---------- */
.dash-chip {
  display: inline-flex;
  align-items: center;
  height: 22px;
  padding: 0 8px;
  border-radius: 6px;
  border: 1px solid #2a2e35;
  background: #1d2026;
  font-size: 11.5px;
  color: #bcc2cc;
  white-space: nowrap;
}

.dash-chip.is-warn {
  background: rgba(245, 185, 61, 0.14);
  border-color: rgba(245, 185, 61, 0.32);
  color: #f7ca6b;
}

.dash-chip.is-danger {
  background: rgba(243, 101, 111, 0.14);
  border-color: rgba(243, 101, 111, 0.32);
  color: #ff9ba2;
}

.dash-chip.is-ok {
  background: rgba(46, 212, 122, 0.14);
  border-color: rgba(46, 212, 122, 0.32);
  color: #5ce49b;
}

.dash-chip.is-ghost {
  background: transparent;
  color: #8f96a1;
}

.dash-link {
  font-size: 12px;
  color: #9db2ff;
  text-decoration: none;
}

.dash-link:hover {
  color: #c9d2ff;
}

.dash-seg {
  margin-left: auto;
  display: inline-flex;
  gap: 2px;
  padding: 3px;
  border: 1px solid #2a2e35;
  border-radius: 8px;
  background: #14161a;
}

.dash-seg button {
  height: 25px;
  padding: 0 9px;
  border: 0;
  border-radius: 6px;
  background: transparent;
  color: #8f96a1;
  font-size: 12px;
  cursor: pointer;
}

.dash-seg button.is-active {
  background: #2b3038;
  color: #fff;
}

.dash-empty-small {
  padding: 22px 0;
  text-align: center;
  font-size: 12.5px;
  color: #8f96a1;
}

.dash-empty {
  padding: 60px 0;
}

@media (max-width: 1440px) {
  .dash-kpis {
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }
}

@media (max-width: 1100px) {
  .dash-kpis {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .dash-split,
  .dash-three {
    grid-template-columns: minmax(0, 1fr);
  }
}
</style>
