<script setup lang="ts">
import { computed, h, onMounted, reactive, ref } from 'vue';
import { useRoute } from 'vue-router';
import { NButton } from 'naive-ui';
import type { DataTableColumns } from 'naive-ui';
import {
  assignUserPath,
  confirmUserRenewal,
  clearUserPath,
  fetchUserDetail,
  fetchUserPathAssets,
  fetchUsers,
  rejectUserRenewal,
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
const rows = ref<Api.Central.UserSummary[]>([]);
const total = ref(0);
const dataAt = ref('');
const drawerVisible = ref(false);
const detailLoading = ref(false);
const detailError = ref('');
const saving = ref(false);
const detail = ref<Api.Central.UserDetail | null>(null);
const selectedUserID = ref('');
const editForm = reactive({
  displayName: '',
  billingCycle: 'monthly' as Api.Central.BillingCycle,
  billingAmount: 0 as number | null,
  notes: ''
});
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
const pathExitIPID = ref<string | null>(null);
const pathNotes = ref('');
const route = useRoute();
let pathAssetsRequestID = 0;

const filters = reactive({ page: 1, page_size: 50, keyword: '', status: 'all' });

const statusOptions = [
  { label: '全部状态', value: 'all' },
  { label: '有效', value: 'active' },
  { label: '即将到期', value: 'expiring' },
  { label: '已到期', value: 'expired' },
  { label: '已停用', value: 'disabled' }
];

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

function copyDetailToForm(value: Api.Central.UserDetail) {
  editForm.displayName = value.displayName;
  editForm.billingCycle = value.billingCycle || 'monthly';
  editForm.billingAmount = value.billingAmount ?? value.monthlyFee;
  editForm.notes = value.notes || '';
}

function syncPathForm(value: Api.Central.UserDetail) {
  const path = value.path;
  pathMode.value = path?.mode || 'relay';
  pathLandingNodeID.value = path?.landingNodeId || null;
  pathLandingInboundID.value = path?.landingInboundId || null;
  pathExitIPID.value = path?.exitIpId || null;
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
  pathExitIPID.value = null;
  pathError.value = '';
}

function handlePathLandingNodeChange(value: string | null) {
  pathLandingNodeID.value = value;
  pathLandingInboundID.value = null;
  pathExitIPID.value = null;
  pathError.value = '';
  pathLandingInboundsLoading.value = false;
  pathLandingInbounds.value = value
    ? pathAssets.value?.landingNodes.find(node => node.id === value)?.inbounds || []
    : [];
}

async function savePathAssignment() {
  if (!detail.value) return;
  if (!pathExitIPID.value) {
    pathError.value = '请选择出口 IP；线路机直出也需要明确指定一个出口 IP。';
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
    exitIpId: pathExitIPID.value,
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
  void loadUsers();
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
  void loadUsers();
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
  loadUsers();
}

function cycleLabel(cycle: Api.Central.BillingCycle) {
  return cycle === 'annual' ? '年付' : '月付';
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
  void loadUsers();
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
  { title: '线路机', key: 'nodeName', minWidth: 140, render: row => row.nodeName || '未关联' },
  { title: '落地机', key: 'landingNodeName', minWidth: 140, render: row => row.landingNodeName || '线路机直出' },
  {
    title: '出口 IP',
    key: 'exitIpAddress',
    minWidth: 220,
    render: row =>
      h('div', [
        h('div', { class: 'font-medium' }, row.exitIpAddress || '未配置'),
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
  { title: '最近活动', key: 'lastActivityAt', minWidth: 170, render: row => formatDate(row.lastActivityAt) },
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

const pagination = computed(() => ({
  page: filters.page,
  pageSize: filters.page_size,
  itemCount: total.value,
  showSizePicker: true,
  pageSizes: [20, 50, 100],
  onChange: (page: number) => {
    filters.page = page;
    loadUsers();
  },
  onUpdatePageSize: (pageSize: number) => {
    filters.page_size = pageSize;
    filters.page = 1;
    loadUsers();
  }
}));

async function loadUsers() {
  loading.value = true;
  errorMessage.value = '';
  const { data, error } = await fetchUsers({
    page: filters.page,
    page_size: filters.page_size,
    keyword: filters.keyword || undefined,
    status: filters.status === 'all' ? undefined : filters.status
  });
  if (error || !data) {
    errorMessage.value = '中央后端暂不可用，无法读取用户数据';
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

function submitFilters() {
  filters.page = 1;
  loadUsers();
}

function resetFilters() {
  filters.keyword = '';
  filters.status = 'all';
  submitFilters();
}

onMounted(() => {
  void loadUsers();
  const userId = typeof route.query.userId === 'string' ? route.query.userId : '';
  if (userId) void openDetail(userId);
});
</script>

<template>
  <div class="users-page">
    <ModulePage
      title="用户管理"
      description="一个 Inbound 对应一个业务用户，Client / Email 仅作为设备凭证。"
      :loading="loading"
      :error="errorMessage"
      :empty="rows.length === 0"
      empty-description="暂无用户同步数据"
      :data-at="dataAt"
      @refresh="loadUsers"
    >
      <template #actions>
        <NButton v-if="errorMessage" size="small" type="warning" secondary @click="loadUsers">重试</NButton>
      </template>
      <template #toolbar>
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
            <NButton type="primary" @click="submitFilters">
              <template #icon><icon-mdi-magnify /></template>
              查询
            </NButton>
            <NButton @click="resetFilters">重置</NButton>
          </NSpace>
        </div>
      </template>
      <NDataTable
        :columns="columns"
        :data="rows"
        :pagination="pagination"
        :bordered="false"
        :single-line="false"
        size="small"
        :scroll-x="1120"
      />
    </ModulePage>

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
              <div
                v-for="record in detail.billingRecords"
                :key="record.id"
                class="flex flex-wrap items-center justify-between gap-8px border-b border-gray-200 py-7px last:border-b-0 dark:border-gray-700"
              >
                <span class="text-12px">
                  {{ cycleLabel(record.billingCycle) }} · {{ formatMoney(record.amount, record.currency) }}
                </span>
                <span class="text-12px text-gray-500">
                  服务 {{ formatDate(record.serviceFrom) }} → {{ formatDate(record.serviceTo) }} · 收款
                  {{ formatDate(record.paidAt) }}
                </span>
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
                    v-model:value="pathExitIPID"
                    :options="pathExitIpOptions"
                    :loading="pathAssetsLoading"
                    :disabled="pathAssetsLoading || !detail.node.id"
                    filterable
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
                      !pathExitIPID
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
                    出口：{{ detail.path.exitIpAddress || '--' }} · 归属：{{
                      detail.path.exitIpOwnerNodeName || '独立 S5'
                    }}
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
                        {{ history.exitIpAddress || '--' }}
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
