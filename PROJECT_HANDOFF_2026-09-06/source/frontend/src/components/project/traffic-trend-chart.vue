<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { fetchUserTraffic } from '@/service/api';
import { useEcharts, type ECOption } from '@/hooks/common/echarts';

defineOptions({ name: 'ProjectTrafficTrendChart' });

const props = defineProps<{
  userId: string;
}>();

const range = ref<Api.Central.TrafficRange>('1h');
const trend = ref<Api.Central.UserTrafficTrend | null>(null);
const loading = ref(false);
const errorMessage = ref('');
let requestID = 0;

const rangeOptions: Array<{ label: string; value: Api.Central.TrafficRange }> = [
  { label: '1 小时', value: '1h' },
  { label: '6 小时', value: '6h' },
  { label: '1 天', value: '1d' },
  { label: '7 天', value: '7d' }
];

const points = computed(() => trend.value?.points || []);
const summary = computed(() => trend.value?.summary);
const hasData = computed(() => points.value.length > 0);
const hasReset = computed(() => points.value.some(point => point.resetDetected));
const hasGap = computed(() => points.value.some(point => point.hasGap));

function formatBytes(value?: number | null) {
  if (value === null || value === undefined || !Number.isFinite(value)) return '--';
  if (value < 1024) return value.toFixed(0) + ' B';
  const units = ['KB', 'MB', 'GB', 'TB', 'PB'];
  let amount = value;
  let unit = -1;
  while (amount >= 1024 && unit < units.length - 1) {
    amount /= 1024;
    unit += 1;
  }
  return (amount >= 100 ? amount.toFixed(0) : amount.toFixed(1)) + ' ' + units[unit];
}

function formatRate(value?: number | null) {
  if (value === null || value === undefined || !Number.isFinite(value)) return '--';
  return formatBytes(value) + '/s';
}

function formatDate(value?: string | null) {
  if (!value) return '--';
  const date = new Date(value);
  return Number.isNaN(date.getTime())
    ? '时间未知'
    : date.toLocaleString('zh-CN', {
        hour12: false,
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit'
      });
}

function formatAxisDate(value: string) {
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return value;
  if (range.value === '7d') return date.toLocaleDateString('zh-CN', { month: '2-digit', day: '2-digit' });
  return date.toLocaleTimeString('zh-CN', { hour12: false, hour: '2-digit', minute: '2-digit' });
}

function buildChartOptions(): ECOption {
  const chartPoints = points.value;
  const labels = chartPoints.map(point => formatAxisDate(point.time));
  const valueFormatter = (value: number) => formatBytes(value);
  return {
    animation: false,
    grid: { left: 58, right: 20, top: 38, bottom: 28 },
    legend: {
      top: 0,
      left: 0,
      itemWidth: 10,
      itemHeight: 10,
      textStyle: { fontSize: 12 }
    },
    tooltip: {
      trigger: 'axis',
      formatter: (params: any) => {
        const first = Array.isArray(params) ? params[0] : params;
        const index = typeof first?.dataIndex === 'number' ? first.dataIndex : 0;
        const point = chartPoints[index];
        if (!point) return '';
        const flags = [point.resetDetected ? '基线重置' : '', point.hasGap ? '采样间隔较长' : '']
          .filter(Boolean)
          .join(' · ');
        return (
          '<div style="margin-bottom:4px">' +
          formatDate(point.time) +
          '</div><div>累计总流量：' +
          formatBytes(point.totalBytes) +
          '</div><div>累计下载：' +
          formatBytes(point.downloadBytes) +
          '</div><div>累计上传：' +
          formatBytes(point.uploadBytes) +
          '</div><div style="margin-top:4px;color:#64748b">区间下载速率：' +
          formatRate(point.downloadRate) +
          '</div><div style="color:#64748b">区间上传速率：' +
          formatRate(point.uploadRate) +
          '</div>' +
          (flags ? '<div style="margin-top:4px;color:#f0a020">' + flags + '</div>' : '')
        );
      }
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: labels,
      axisLabel: { hideOverlap: true, fontSize: 11 },
      axisTick: { show: false }
    },
    yAxis: {
      type: 'value',
      axisLabel: { formatter: valueFormatter, fontSize: 11 },
      splitLine: { lineStyle: { type: 'dashed', opacity: 0.35 } }
    },
    series: [
      {
        name: '总流量累计',
        type: 'line',
        smooth: 0.2,
        showSymbol: false,
        connectNulls: false,
        data: chartPoints.map(point => point.totalBytes),
        lineStyle: { width: 2.5, color: '#2563eb' },
        itemStyle: { color: '#2563eb' },
        areaStyle: { color: 'rgba(37, 99, 235, 0.12)' }
      },
      {
        name: '下载累计',
        type: 'line',
        smooth: 0.25,
        showSymbol: false,
        connectNulls: false,
        data: chartPoints.map(point => point.downloadBytes),
        lineStyle: { width: 1.5, color: '#3b82f6', type: 'dashed' },
        itemStyle: { color: '#3b82f6' }
      },
      {
        name: '上传累计',
        type: 'line',
        smooth: 0.25,
        showSymbol: false,
        connectNulls: false,
        data: chartPoints.map(point => point.uploadBytes),
        lineStyle: { width: 1.5, color: '#8b5cf6', type: 'dashed' },
        itemStyle: { color: '#8b5cf6' }
      }
    ]
  };
}

const { domRef, updateOptions } = useEcharts(() => buildChartOptions());

async function loadTrend() {
  const userID = props.userId.trim();
  if (!userID) return;
  const currentRequestID = ++requestID;
  loading.value = true;
  errorMessage.value = '';
  const { data, error } = await fetchUserTraffic(userID, range.value);
  if (currentRequestID !== requestID) return;
  loading.value = false;
  if (error || !data) {
    trend.value = null;
    errorMessage.value = '无法读取该时间范围的流量趋势，请确认中央后端已更新。';
    return;
  }
  trend.value = data;
  void updateOptions(() => buildChartOptions());
}

function retry() {
  void loadTrend();
}

watch(
  () => [props.userId, range.value],
  () => {
    void loadTrend();
  },
  { immediate: true }
);
</script>

<template>
  <div>
    <div class="mb-12px flex flex-wrap items-center justify-between gap-8px">
      <div class="text-12px text-gray-500">
        曲线展示所选范围内的累计增长；数据截至 {{ formatDate(trend?.dataAt) }}
      </div>
      <NButtonGroup size="small">
        <NButton
          v-for="item in rangeOptions"
          :key="item.value"
          :type="range === item.value ? 'primary' : 'default'"
          @click="range = item.value"
        >
          {{ item.label }}
        </NButton>
      </NButtonGroup>
    </div>

    <NSpin :show="loading">
      <NAlert v-if="errorMessage" type="error" :show-icon="true" class="mb-12px">
        {{ errorMessage }}
        <NButton size="small" class="ml-8px" @click="retry">重试</NButton>
      </NAlert>

      <NEmpty v-if="!errorMessage && !hasData" description="该时间范围暂无足够的流量快照" size="small" />
      <div v-else-if="!errorMessage" ref="domRef" class="h-260px w-full" />

      <div v-if="hasReset || hasGap" class="mt-8px space-y-6px">
        <NAlert v-if="hasReset" type="warning" :show-icon="true" size="small">
          检测到流量基线重置，重置点之后按新的基线重新统计。
        </NAlert>
        <NAlert v-if="hasGap" type="info" :show-icon="true" size="small">
          部分采样间隔较长，图表仅展示实际采样计算出的速率，没有用 0 填充缺口。
        </NAlert>
      </div>

      <div v-if="summary && hasData && !errorMessage" class="mt-12px grid grid-cols-2 gap-8px md:grid-cols-3">
        <div class="rounded-6px bg-gray-50 px-10px py-8px dark:bg-dark-200">
          <div class="text-12px text-gray-500">下载增长量</div>
          <div class="mt-2px text-14px font-600">{{ formatBytes(summary.downloadBytes) }}</div>
        </div>
        <div class="rounded-6px bg-gray-50 px-10px py-8px dark:bg-dark-200">
          <div class="text-12px text-gray-500">上传增长量</div>
          <div class="mt-2px text-14px font-600">{{ formatBytes(summary.uploadBytes) }}</div>
        </div>
        <div class="rounded-6px bg-gray-50 px-10px py-8px dark:bg-dark-200">
          <div class="text-12px text-gray-500">总增长量</div>
          <div class="mt-2px text-14px font-600">{{ formatBytes(summary.totalBytes) }}</div>
        </div>
        <div class="rounded-6px bg-gray-50 px-10px py-8px dark:bg-dark-200">
          <div class="text-12px text-gray-500">平均下载速率</div>
          <div class="mt-2px text-14px font-600">{{ formatRate(summary.averageDownloadRate) }}</div>
        </div>
        <div class="rounded-6px bg-gray-50 px-10px py-8px dark:bg-dark-200">
          <div class="text-12px text-gray-500">峰值下载速率</div>
          <div class="mt-2px text-14px font-600">{{ formatRate(summary.peakDownloadRate) }}</div>
        </div>
        <div class="rounded-6px bg-gray-50 px-10px py-8px dark:bg-dark-200">
          <div class="text-12px text-gray-500">平均上传速率</div>
          <div class="mt-2px text-14px font-600">{{ formatRate(summary.averageUploadRate) }}</div>
        </div>
        <div class="rounded-6px bg-gray-50 px-10px py-8px dark:bg-dark-200">
          <div class="text-12px text-gray-500">峰值上传速率</div>
          <div class="mt-2px text-14px font-600">{{ formatRate(summary.peakUploadRate) }}</div>
        </div>
        <div class="rounded-6px bg-gray-50 px-10px py-8px dark:bg-dark-200">
          <div class="text-12px text-gray-500">数据覆盖率</div>
          <div class="mt-2px text-14px font-600">{{ (summary.coverage * 100).toFixed(0) }}%</div>
        </div>
      </div>
    </NSpin>
  </div>
</template>
