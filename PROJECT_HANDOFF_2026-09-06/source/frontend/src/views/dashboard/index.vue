<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { fetchDashboard } from '@/service/api';
import DataFreshness from '@/components/project/data-freshness.vue';
import TrafficValue from '@/components/project/traffic-value.vue';

defineOptions({ name: 'CentralDashboard' });

const loading = ref(false);
const errorMessage = ref('');
const dashboard = ref<Api.Central.DashboardSummary | null>(null);

const nodeCards = computed(() => {
  const nodes = dashboard.value?.nodes;

  return [
    { label: '节点总数', value: nodes?.total ?? null },
    { label: '在线节点', value: nodes?.online ?? null },
    { label: '线路机', value: nodes?.relay ?? null },
    { label: '落地机', value: nodes?.landing ?? null }
  ];
});

const userCards = computed(() => {
  const users = dashboard.value?.users;

  return [
    { label: '有效用户', value: users?.active ?? null, type: 'success' as const },
    { label: '即将到期', value: users?.expiring ?? null, type: 'warning' as const },
    { label: '已到期', value: users?.expired ?? null, type: 'error' as const }
  ];
});

function formatMoney(value: number | undefined, currency = 'CNY') {
  if (value === undefined) return '--';
  return new Intl.NumberFormat('zh-CN', { style: 'currency', currency }).format(value);
}

async function loadDashboard() {
  loading.value = true;
  errorMessage.value = '';

  const { data, error } = await fetchDashboard();

  if (error) {
    errorMessage.value = '中央后端暂不可用，请稍后重试';
    dashboard.value = null;
  } else {
    dashboard.value = data;
  }

  loading.value = false;
}

onMounted(loadDashboard);
</script>

<template>
  <NSpace vertical :size="16">
    <NAlert v-if="errorMessage" type="warning" :show-icon="false">
      <div class="flex items-center justify-between gap-16px">
        <span>{{ errorMessage }}</span>
        <NButton size="small" secondary @click="loadDashboard">重新同步</NButton>
      </div>
    </NAlert>

    <NCard :bordered="false" size="small" title="运营总览">
      <template #header-extra>
        <NSpin v-if="loading" size="small" />
        <NSpace v-else align="center" :size="8">
          <DataFreshness :data-at="dashboard?.dataAt" />
          <span class="text-12px text-gray-500">
            {{ dashboard?.generatedAt ? `查询时间：${dashboard.generatedAt}` : '等待中央数据' }}
          </span>
        </NSpace>
      </template>
      <NGrid cols="1 s:2 m:4" responsive="screen" :x-gap="16" :y-gap="16">
        <NGi v-for="card in nodeCards" :key="card.label">
          <NCard size="small" embedded>
            <NStatistic :label="card.label" :value="card.value ?? '--'" />
          </NCard>
        </NGi>
      </NGrid>
    </NCard>

    <NCard :bordered="false" size="small" title="用户状态">
      <NGrid cols="1 s:3" responsive="screen" :x-gap="16" :y-gap="16">
        <NGi v-for="card in userCards" :key="card.label">
          <NCard size="small" embedded>
            <NStatistic :label="card.label" :value="card.value ?? '--'">
              <template #suffix><NTag :type="card.type" size="small">Inbound</NTag></template>
            </NStatistic>
          </NCard>
        </NGi>
      </NGrid>
    </NCard>

    <NGrid cols="1 s:2" responsive="screen" :x-gap="16" :y-gap="16">
      <NGi>
        <NCard :bordered="false" size="small" title="流量">
          <div class="grid grid-cols-2 gap-16px">
            <NStatistic label="今日流量"><TrafficValue :value="dashboard?.traffic.todayBytes" /></NStatistic>
            <NStatistic label="本月流量"><TrafficValue :value="dashboard?.traffic.monthBytes" /></NStatistic>
          </div>
        </NCard>
      </NGi>
      <NGi>
        <NCard :bordered="false" size="small" title="财务估算">
          <div class="grid grid-cols-3 gap-16px">
            <NStatistic
              label="本月收入"
              :value="formatMoney(dashboard?.finance.monthIncome, dashboard?.finance.currency)"
            />
            <NStatistic
              label="本月成本"
              :value="formatMoney(dashboard?.finance.monthCost, dashboard?.finance.currency)"
            />
            <NStatistic
              label="预计毛利润"
              :value="formatMoney(dashboard?.finance.grossProfit, dashboard?.finance.currency)"
            />
          </div>
        </NCard>
      </NGi>
    </NGrid>
  </NSpace>
</template>
