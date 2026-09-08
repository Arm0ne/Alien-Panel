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
const activeTab = ref('overview');

function formatMoney(value: number | undefined, currency = 'CNY') {
  if (value === undefined || value === null) return '--';
  return new Intl.NumberFormat('zh-CN', { style: 'currency', currency, maximumFractionDigits: 2 }).format(value);
}

function formatDate(value?: string | null) {
  if (!value) return '--';
  const date = new Date(value);
  return Number.isNaN(date.getTime()) ? value : date.toLocaleDateString('zh-CN');
}

function orderTypeLabel(value: string) {
  return value === 'renewal' ? '续费' : value === 'recovery' ? '逾期恢复' : '首购';
}
function cycleLabel(value: Api.Central.BillingCycle) {
  return value === 'annual' ? '年付' : '月付';
}
function costTypeLabel(value: string) {
  return value === 'node' ? '节点' : value === 'ip' ? '出口 IP' : '其他';
}

const orders = computed(() => summary.value?.orders || []);
const costs = computed(() => summary.value?.costDetails || []);
const paidBreakdown = computed(() => summary.value?.cashBreakdown || []);
const renewal = computed(
  () => summary.value?.renewal || { dueUsers: 0, renewedUsers: 0, notRenewedUsers: 0, renewalRate: 0, recoveryUsers: 0 }
);
const totalBreakdown = computed(() => paidBreakdown.value.reduce((sum, item) => sum + item.amount, 0));

const orderColumns: DataTableColumns<(typeof orders.value)[number]> = [
  { title: '订单号', key: 'orderNo', minWidth: 155, render: row => row.orderNo || row.id.slice(0, 12) },
  { title: '业务用户', key: 'userName', minWidth: 130 },
  { title: '类型', key: 'orderType', width: 90, render: row => orderTypeLabel(row.orderType) },
  { title: '周期', key: 'billingCycle', width: 80, render: row => cycleLabel(row.billingCycle) },
  { title: '金额', key: 'amount', width: 120, render: row => formatMoney(row.amount, row.currency) },
  {
    title: '服务区间',
    key: 'serviceFrom',
    minWidth: 210,
    render: row => `${formatDate(row.serviceFrom)} → ${formatDate(row.serviceTo)}`
  },
  { title: '收款时间', key: 'paidAt', minWidth: 150, render: row => formatDate(row.paidAt) },
  {
    title: '来源',
    key: 'source',
    width: 110,
    render: row => (row.source === 'agent_detected' ? 'Agent 候选' : '手工确认')
  }
];

const costColumns: DataTableColumns<(typeof costs.value)[number]> = [
  { title: '类型', key: 'type', width: 90, render: row => costTypeLabel(row.type) },
  { title: '成本项目', key: 'name', minWidth: 160 },
  { title: '分类', key: 'category', minWidth: 120 },
  { title: '本期金额', key: 'amount', width: 130, render: row => formatMoney(row.amount, row.currency) },
  {
    title: '有效期',
    key: 'effectiveFrom',
    minWidth: 190,
    render: row => `${formatDate(row.effectiveFrom)} → ${formatDate(row.effectiveTo)}`
  },
  { title: '备注', key: 'notes', minWidth: 160, render: row => row.notes || '--' }
];

async function loadFinance() {
  loading.value = true;
  errorMessage.value = '';
  const { data, error } = await fetchFinanceSummary({ period: period.value || undefined });
  if (error || !data) {
    errorMessage.value = '中央后端暂不可用，无法读取财务数据';
    summary.value = null;
  } else summary.value = data;
  loading.value = false;
}

onMounted(loadFinance);
</script>

<template>
  <ModulePage
    title="成本与收入"
    description="以实际收款为主，核对订单、续费和真实成本。"
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
          <NDatePicker v-model:formatted-value="period" type="month" value-format="yyyy-MM" class="w-160px" />
          <NButton type="primary" :loading="loading" @click="loadFinance">
            <template #icon><icon-mdi-refresh /></template>
            刷新数据
          </NButton>
        </NSpace>
      </div>
    </template>

    <div v-if="summary" class="finance-content p-16px">
      <div class="grid grid-cols-1 gap-12px sm:grid-cols-2 xl:grid-cols-5">
        <NCard size="small" embedded>
          <NStatistic label="本期实收" :value="formatMoney(summary.cashIncome, summary.currency)" />
          <div class="mt-4px text-12px text-gray-500">{{ orders.length }} 笔已确认订单</div>
        </NCard>
        <NCard size="small" embedded>
          <NStatistic label="累计实收" :value="formatMoney(summary.cumulativeCashIncome, summary.currency)" />
          <div class="mt-4px text-12px text-gray-500">从首笔确认订单累计</div>
        </NCard>
        <NCard size="small" embedded>
          <NStatistic label="本期成本" :value="formatMoney(summary.monthCost, summary.currency)" />
          <div class="mt-4px text-12px text-gray-500">节点、IP 与其他真实成本</div>
        </NCard>
        <NCard size="small" embedded>
          <NStatistic label="本期现金结余" :value="formatMoney(summary.cashBalance, summary.currency)" />
          <div class="mt-4px text-12px text-gray-500">本期实收 - 本期成本</div>
        </NCard>
        <NCard size="small" embedded>
          <NStatistic label="本期到期收费用户" :value="renewal.dueUsers" />
          <div class="mt-4px text-12px text-gray-500">免费用户不计入续费率</div>
        </NCard>
      </div>

      <NTabs v-model:value="activeTab" type="line" animated class="mt-16px">
        <NTabPane name="overview" tab="财务总览">
          <div class="grid grid-cols-1 gap-16px xl:grid-cols-2">
            <NCard title="收款构成" size="small" :segmented="{ content: true }">
              <div class="breakdown-layout">
                <div
                  class="breakdown-donut"
                  :style="{
                    background: totalBreakdown
                      ? `conic-gradient(#646cff 0 ${((paidBreakdown.find(item => item.orderType === 'renewal')?.amount || 0) / totalBreakdown) * 360}deg, #18a058 0 ${(((paidBreakdown.find(item => item.orderType === 'renewal')?.amount || 0) + (paidBreakdown.find(item => item.orderType === 'initial')?.amount || 0)) / totalBreakdown) * 360}deg, #f0a020 0 360deg)`
                      : undefined
                  }"
                >
                  <div class="breakdown-donut-hole">
                    <strong>{{ formatMoney(summary.cashIncome, summary.currency) }}</strong>
                    <span>本期实收</span>
                  </div>
                </div>
                <div class="breakdown-legend">
                  <div v-for="item in paidBreakdown" :key="item.orderType" class="breakdown-line">
                    <span>
                      <i :class="`breakdown-dot breakdown-${item.orderType}`"></i>
                      {{ orderTypeLabel(item.orderType) }}
                    </span>
                    <strong>
                      {{ formatMoney(item.amount, summary.currency) }} ·
                      {{ totalBreakdown ? `${((item.amount / totalBreakdown) * 100).toFixed(0)}%` : '0%' }}
                    </strong>
                  </div>
                  <div class="mt-8px text-12px text-gray-500">
                    服务期折算收入
                    {{ formatMoney(summary.monthIncome, summary.currency) }} 仅用于经营分析，不代表待收款。
                  </div>
                </div>
              </div>
            </NCard>
            <NCard title="用户与续费" size="small" :segmented="{ content: true }">
              <div class="grid grid-cols-2 gap-12px md:grid-cols-4">
                <NStatistic label="收费用户" :value="summary.paidUserCount" />
                <NStatistic label="免费用户" :value="summary.freeUserCount" />
                <NStatistic label="严格续费率" :value="`${(renewal.renewalRate * 100).toFixed(1)}%`" />
                <NStatistic label="逾期恢复" :value="renewal.recoveryUsers" />
              </div>
              <NAlert class="mt-16px" type="info" :show-icon="true">
                到期即停止连接；到期前确认的下一笔订单计入续费率，逾期付款单独记为逾期恢复。
              </NAlert>
            </NCard>
          </div>
          <NCard class="mt-16px" title="本期收款订单" size="small" :segmented="{ content: true }">
            <NDataTable
              :columns="orderColumns"
              :data="orders.slice(0, 8)"
              :bordered="false"
              :scroll-x="1050"
              size="small"
            />
            <div class="mt-12px text-right">
              <NButton text type="primary" @click="activeTab = 'orders'">查看全部订单</NButton>
            </div>
          </NCard>
        </NTabPane>
        <NTabPane name="orders" tab="收款订单">
          <NCard title="收款订单台账" size="small" :segmented="{ content: true }">
            <template #header-extra>
              <span class="text-12px text-gray-500">按实际收款日期归属 {{ period }}</span>
            </template>
            <NDataTable :columns="orderColumns" :data="orders" :bordered="false" :scroll-x="1050" size="small" />
          </NCard>
        </NTabPane>
        <NTabPane name="renewals" tab="续费分析">
          <div class="grid grid-cols-1 gap-12px sm:grid-cols-2 xl:grid-cols-5">
            <NCard size="small" embedded><NStatistic label="本期到期收费用户" :value="renewal.dueUsers" /></NCard>
            <NCard size="small" embedded><NStatistic label="到期前已续费" :value="renewal.renewedUsers" /></NCard>
            <NCard size="small" embedded>
              <NStatistic label="严格续费率" :value="`${(renewal.renewalRate * 100).toFixed(1)}%`" />
            </NCard>
            <NCard size="small" embedded><NStatistic label="到期未续费" :value="renewal.notRenewedUsers" /></NCard>
            <NCard size="small" embedded><NStatistic label="逾期恢复" :value="renewal.recoveryUsers" /></NCard>
          </div>
          <NCard class="mt-16px" title="续费口径" size="small">
            <p class="m-0 text-13px text-gray-600 dark:text-gray-300">
              续费率分母是本期到期的收费用户，分子是到期前已经确认的 renewal
              订单。免费用户、逾期恢复和没有付款的到期用户不会被混入收费续费率。
            </p>
          </NCard>
        </NTabPane>
        <NTabPane name="costs" tab="成本明细">
          <NCard title="真实成本台账" size="small" :segmented="{ content: true }">
            <template #header-extra><span class="text-12px text-gray-500">不做用户均摊</span></template>
            <NDataTable :columns="costColumns" :data="costs" :bordered="false" :scroll-x="850" size="small" />
          </NCard>
        </NTabPane>
      </NTabs>
    </div>
  </ModulePage>
</template>

<style scoped>
.breakdown-layout {
  display: grid;
  grid-template-columns: minmax(148px, 180px) minmax(0, 1fr);
  align-items: center;
  gap: 24px;
}
.breakdown-donut {
  display: grid;
  width: 156px;
  height: 156px;
  margin: 0 auto;
  place-items: center;
  border-radius: 50%;
}
.breakdown-donut-hole {
  display: grid;
  width: 106px;
  height: 106px;
  place-content: center;
  gap: 4px;
  border-radius: 50%;
  background: var(--n-color);
  text-align: center;
}
.breakdown-donut-hole strong {
  max-width: 100px;
  overflow: hidden;
  font-size: 14px;
  line-height: 20px;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.breakdown-donut-hole span {
  color: var(--n-text-color-3);
  font-size: 12px;
}
.breakdown-legend {
  min-width: 0;
}
.breakdown-line {
  display: flex;
  min-width: 0;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 8px 0;
  border-bottom: 1px solid var(--n-divider-color);
  font-size: 13px;
}
.breakdown-line > span,
.breakdown-line > strong {
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.breakdown-dot {
  display: inline-block;
  width: 8px;
  height: 8px;
  margin-right: 7px;
  border-radius: 50%;
}
.breakdown-renewal {
  background: #646cff;
}
.breakdown-initial {
  background: #18a058;
}
.breakdown-recovery {
  background: #f0a020;
}
@media (max-width: 560px) {
  .breakdown-layout {
    grid-template-columns: 1fr;
    gap: 12px;
  }
  .breakdown-donut {
    width: 136px;
    height: 136px;
  }
  .breakdown-donut-hole {
    width: 94px;
    height: 94px;
  }
}
</style>
