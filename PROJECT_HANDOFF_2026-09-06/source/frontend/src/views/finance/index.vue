<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import type { DataTableColumns } from 'naive-ui';
import { fetchFinanceSummary } from '@/service/api';
import ModulePage from '@/components/project/module-page.vue';

defineOptions({ name: 'FinanceManagement' });

const loading = ref(false);
const errorMessage = ref('');
const summary = ref<Api.Central.FinanceSummary | null>(null);
const period = ref<string | null>(new Date().toISOString().slice(0, 7));

function formatMoney(value: number | undefined, currency = 'CNY') {
  if (value === undefined || value === null) return '--';
  return new Intl.NumberFormat('zh-CN', { style: 'currency', currency }).format(value);
}

const breakdownColumns: DataTableColumns<{ label: string; amount: number }> = [
  { title: '成本/收入项目', key: 'label', minWidth: 220 },
  { title: '金额', key: 'amount', minWidth: 140, render: row => formatMoney(row.amount, summary.value?.currency) }
];

const breakdown = computed(() => summary.value?.breakdown || []);

async function loadFinance() {
  loading.value = true;
  errorMessage.value = '';
  const { data, error } = await fetchFinanceSummary({ period: period.value || undefined });
  if (error || !data) {
    errorMessage.value = '中央后端暂不可用，无法读取财务数据';
    summary.value = null;
  } else {
    summary.value = data;
  }
  loading.value = false;
}

onMounted(loadFinance);
</script>

<template>
  <ModulePage
    title="成本与收入"
    description="按月份查看用户应计收入、已确认实收、节点成本、出口 IP 成本和预计毛利润。"
    :loading="loading"
    :error="errorMessage"
    :empty="!summary"
    empty-description="暂无财务汇总数据"
    :data-at="summary?.dataAt"
    @refresh="loadFinance"
  >
    <template #toolbar>
      <div class="border-b border-gray-200 p-16px dark:border-gray-700">
        <NSpace align="center" wrap>
          <span class="text-14px text-gray-600 dark:text-gray-300">统计月份</span>
          <NDatePicker
            v-model:formatted-value="period"
            type="month"
            value-format="yyyy-MM"
            class="w-160px"
            placeholder="选择统计月份"
          />
          <NButton type="primary" @click="loadFinance">
            <template #icon><icon-mdi-magnify /></template>
            查询
          </NButton>
        </NSpace>
      </div>
    </template>
    <div class="grid grid-cols-1 gap-16px p-16px sm:grid-cols-2 lg:grid-cols-4">
      <NCard size="small" embedded>
        <NStatistic label="有效用户数" :value="summary?.effectiveUserCount ?? 0" />
      </NCard>
      <NCard size="small" embedded>
        <NStatistic label="本月收入" :value="formatMoney(summary?.monthIncome, summary?.currency)" />
      </NCard>
      <NCard size="small" embedded>
        <NStatistic label="已确认实收" :value="formatMoney(summary?.cashIncome, summary?.currency)" />
      </NCard>
      <NCard size="small" embedded>
        <NStatistic label="本月成本" :value="formatMoney(summary?.monthCost, summary?.currency)" />
      </NCard>
      <NCard size="small" embedded>
        <NStatistic label="预计毛利润" :value="formatMoney(summary?.grossProfit, summary?.currency)" />
      </NCard>
    </div>
    <NDataTable v-if="breakdown.length" :columns="breakdownColumns" :data="breakdown" :bordered="false" size="small" />
    <NEmpty v-else class="pb-24px" description="暂无成本明细" />
  </ModulePage>
</template>
