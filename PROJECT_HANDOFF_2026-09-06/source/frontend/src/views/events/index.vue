<script setup lang="ts">
import { computed, h, onMounted, reactive, ref } from 'vue';
import { useRouter } from 'vue-router';
import { NButton, NTag, NTabs, NTabPane } from 'naive-ui';
import type { DataTableColumns } from 'naive-ui';
import {
  confirmUserRenewal,
  fetchEvents,
  markAllEventsRead,
  markEventRead,
  rejectUserRenewal,
  requestNodeSync,
  resolveEvent
} from '@/service/api';
import ModulePage from '@/components/project/module-page.vue';

defineOptions({ name: 'EventManagement' });

type EventTab = 'pending' | 'alerts' | 'all';

const loading = ref(false);
const actionLoading = ref('');
const router = useRouter();
const errorMessage = ref('');
const rows = ref<Api.Central.EventSummary[]>([]);
const total = ref(0);
const dataAt = ref('');
const filters = reactive({ page: 1, page_size: 20, keyword: '', severity: 'all', tab: 'pending' as EventTab });

const severityOptions = [
  { label: '全部级别', value: 'all' },
  { label: '信息', value: 'info' },
  { label: '警告', value: 'warning' },
  { label: '错误', value: 'error' }
];

function severityType(severity: Api.Central.EventSeverity): 'info' | 'warning' | 'error' {
  if (severity === 'warning') return 'warning';
  if (severity === 'error') return 'error';
  return 'info';
}

function severityLabel(severity: Api.Central.EventSeverity) {
  return severityOptions.find(option => option.value === severity)?.label || '未知';
}

function formatDate(value: string) {
  const date = new Date(value);
  return Number.isNaN(date.getTime()) ? value : date.toLocaleString('zh-CN', { hour12: false });
}

function payloadValue(row: Api.Central.EventSummary, key: string) {
  return row.payload && typeof row.payload[key] === 'string' ? row.payload[key] : '';
}

function isRenewal(row: Api.Central.EventSummary) {
  return row.type === 'renewal_candidate_detected' && row.status !== 'resolved' && row.status !== 'dismissed';
}

function eventStatusLabel(row: Api.Central.EventSummary) {
  if (row.status === 'resolved' || row.status === 'dismissed') return '已处理';
  if (row.requiresAction) return '待处理';
  return row.acknowledged ? '已读' : '未读';
}

function notifyEventChanged() {
  window.dispatchEvent(new Event('xpanel-events-updated'));
}

function openUser(row: Api.Central.EventSummary) {
  const userId = payloadValue(row, 'userId');
  if (userId) router.push({ name: 'users', query: { userId } });
}

async function loadEvents() {
  loading.value = true;
  errorMessage.value = '';
  const status = filters.tab === 'pending' ? 'pending' : filters.tab === 'alerts' ? 'alerts' : undefined;
  const { data, error } = await fetchEvents({
    page: filters.page,
    page_size: filters.page_size,
    keyword: filters.keyword || undefined,
    severity: filters.severity === 'all' ? undefined : (filters.severity as Api.Central.EventSeverity),
    status,
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
  loading.value = false;
}

function selectTab(value: string) {
  filters.tab = value as EventTab;
  filters.page = 1;
  loadEvents();
}

function submitFilters() {
  filters.page = 1;
  loadEvents();
}

function resetFilters() {
  filters.keyword = '';
  filters.severity = 'all';
  submitFilters();
}

async function readEvent(row: Api.Central.EventSummary) {
  if (row.acknowledged) return;
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

async function confirmRenewal(row: Api.Central.EventSummary) {
  const userId = payloadValue(row, 'userId');
  const candidateId = payloadValue(row, 'candidateId') || row.resourceId || '';
  if (!userId || !candidateId) return;
  actionLoading.value = row.id;
  const billingCycle = payloadValue(row, 'billingCycle');
  const suggestedAmount = Number(row.payload?.suggestedAmount || 0);
  const { error } = await confirmUserRenewal(userId, candidateId, {
    billingCycle: billingCycle === 'annual' ? 'annual' : 'monthly',
    amount: suggestedAmount > 0 ? suggestedAmount : undefined
  });
  actionLoading.value = '';
  if (error) {
    window.$message?.error('续费确认失败，请稍后重试');
    return;
  }
  window.$message?.success('续费已确认并计入财务');
  notifyEventChanged();
  loadEvents();
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
  loadEvents();
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
  loadEvents();
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
    row.acknowledged = true;
  });
  notifyEventChanged();
}

const columns: DataTableColumns<Api.Central.EventSummary> = [
  {
    title: '级别',
    key: 'severity',
    width: 80,
    render: row => h(NTag, { size: 'small', type: severityType(row.severity) }, { default: () => severityLabel(row.severity) })
  },
  { title: '事件', key: 'title', minWidth: 190, render: row => row.title || row.type },
  { title: '节点', key: 'nodeName', minWidth: 130, render: row => row.nodeName || '--' },
  { title: '内容', key: 'message', minWidth: 300, ellipsis: { tooltip: true } },
  { title: '发生时间', key: 'occurredAt', minWidth: 165, render: row => formatDate(row.occurredAt) },
  {
    title: '状态',
    key: 'status',
    width: 90,
    render: row => h(NTag, { size: 'small', type: row.status === 'resolved' || row.status === 'dismissed' ? 'success' : row.requiresAction ? 'warning' : 'default' }, { default: () => eventStatusLabel(row) })
  },
  {
    title: '操作',
    key: 'actions',
    width: 230,
    render: row => {
      const buttons = [];
      if (isRenewal(row)) {
        buttons.push(
          h(NButton, { size: 'small', quaternary: true, onClick: () => openUser(row) }, { default: () => '查看用户' }),
          h(NButton, { size: 'small', type: 'primary', loading: actionLoading.value === row.id, onClick: () => confirmRenewal(row) }, { default: () => '确认续费' }),
          h(NButton, { size: 'small', loading: actionLoading.value === row.id, onClick: () => rejectRenewal(row) }, { default: () => '非收费变更' })
        );
      } else if (row.status !== 'resolved' && row.actionType === 'retry_sync') {
        buttons.push(
          h(NButton, { size: 'small', type: 'primary', loading: actionLoading.value === row.id, onClick: () => retrySync(row) }, { default: () => '重新同步' }),
          h(NButton, { size: 'small', quaternary: true, loading: actionLoading.value === row.id, onClick: () => resolve(row) }, { default: () => '关闭' })
        );
      } else if (row.status !== 'resolved' && row.actionType) {
        buttons.push(h(NButton, { size: 'small', loading: actionLoading.value === row.id, onClick: () => resolve(row) }, { default: () => '标记已处理' }));
      }
      if (!row.acknowledged) {
        buttons.push(h(NButton, { size: 'small', quaternary: true, loading: actionLoading.value === `read:${row.id}`, onClick: () => readEvent(row) }, { default: () => '已读' }));
      }
      return h('div', { class: 'flex flex-wrap gap-6px' }, buttons);
    }
  }
];

const pagination = computed(() => ({
  page: filters.page,
  pageSize: filters.page_size,
  itemCount: total.value,
  showSizePicker: true,
  pageSizes: [20, 50, 100],
  onChange: (page: number) => {
    filters.page = page;
    loadEvents();
  },
  onUpdatePageSize: (pageSize: number) => {
    filters.page_size = pageSize;
    filters.page = 1;
    loadEvents();
  }
}));

onMounted(loadEvents);
</script>

<template>
  <ModulePage title="事件中心" :loading="loading" :error="errorMessage" :empty="rows.length === 0" empty-description="暂无事件" :data-at="dataAt" @refresh="loadEvents">
    <template #actions>
      <NButton size="small" :loading="actionLoading === 'all'" @click="readAll">全部标为已读</NButton>
    </template>
    <template #toolbar>
      <div class="border-b border-gray-200 p-16px dark:border-gray-700">
        <NTabs :value="filters.tab" type="line" @update:value="selectTab">
          <NTabPane name="pending" tab="待处理" />
          <NTabPane name="alerts" tab="告警" />
          <NTabPane name="all" tab="全部" />
        </NTabs>
        <NSpace wrap class="mt-12px">
          <NInput v-model:value="filters.keyword" clearable class="w-240px" placeholder="搜索事件、节点或内容" @keyup.enter="submitFilters" />
          <NSelect v-model:value="filters.severity" :options="severityOptions" class="w-140px" />
          <NButton type="primary" @click="submitFilters"><template #icon><icon-mdi-magnify /></template>查询</NButton>
          <NButton @click="resetFilters">重置</NButton>
        </NSpace>
      </div>
    </template>
    <NDataTable :columns="columns" :data="rows" :pagination="pagination" :bordered="false" :single-line="false" size="small" :scroll-x="1180" />
  </ModulePage>
</template>
