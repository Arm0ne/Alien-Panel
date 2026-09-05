<script setup lang="ts">
import { computed, h, onMounted, reactive, ref, watch } from 'vue';
import { NButton, NInputNumber, NSelect, NSpace, NSwitch, NTag } from 'naive-ui';
import type { DataTableColumns } from 'naive-ui';
import {
  bindRouteExitIp,
  createRoute,
  deleteRoute,
  fetchExitIps,
  fetchNodes,
  fetchRouteDetail,
  fetchRouteExitIpBindings,
  fetchRoutes,
  unbindRouteExitIp,
  updateRoute,
  updateRouteExitIp
} from '@/service/api';
import ModulePage from '@/components/project/module-page.vue';

defineOptions({ name: 'RouteManagement' });

const loading = ref(false);
const errorMessage = ref('');
const rows = ref<Api.Central.RouteSummary[]>([]);
const total = ref(0);
const dataAt = ref('');
const filters = reactive({ page: 1, page_size: 50, keyword: '', status: 'all' });
const nodes = ref<Api.Central.NodeSummary[]>([]);
const nodesLoading = ref(false);
const nodeError = ref('');
const modalVisible = ref(false);
const modalLoading = ref(false);
const saving = ref(false);
const deletingRouteId = ref('');
const formError = ref('');
const editingRouteId = ref<string | null>(null);
const routeForm = reactive({
  name: '',
  relayNodeId: '',
  landingNodeId: '',
  relayOutboundTag: '',
  landingInboundId: '',
  landingInboundTag: '',
  enabled: true,
  validFrom: null as string | null,
  validTo: null as string | null,
  notes: ''
});
const bindingVisible = ref(false);
const bindingRoute = ref<Api.Central.RouteSummary | null>(null);
const bindingRows = ref<Api.Central.RouteExitIpBinding[]>([]);
const bindingExitIps = ref<Api.Central.ExitIpSummary[]>([]);
const bindingLoading = ref(false);
const bindingSaving = ref(false);
const bindingError = ref('');
const bindingDeletingId = ref('');
const bindingUpdatingId = ref('');
const bindingForm = reactive({ exitIpId: '', scope: 'landing' as Api.Central.RouteExitIpScope, allocationWeight: 1 as number | null });
const bindingDrafts = reactive<Record<string, { scope: Api.Central.RouteExitIpScope; allocationWeight: number | null; enabled: boolean }>>({});
const statusOptions = [
  { label: '全部状态', value: 'all' },
  { label: '启用', value: 'active' },
  { label: '停用', value: 'disabled' }
];

function statusLabel(status: Api.Central.RouteSummary['status']) {
  return statusOptions.find(option => option.value === status)?.label || '未知';
}

const relayNodeOptions = computed(() =>
  nodes.value.filter(node => node.type === 'relay').map(node => ({ label: `${node.name}（${node.host || '地址未知'}）`, value: node.id }))
);
const landingNodeOptions = computed(() =>
  nodes.value.filter(node => node.type === 'landing').map(node => ({ label: `${node.name}（${node.host || '地址未知'}）`, value: node.id }))
);
const bindingScopeOptions = [
  { label: '落地机出口', value: 'landing' },
  { label: '线路机直出', value: 'relay' },
  { label: '外部 S5', value: 'external' }
];
const modalTitle = computed(() => (editingRouteId.value ? '编辑线路关系' : '新建线路关系'));
const bindingExitIpOptions = computed(() => {
  const route = bindingRoute.value;
  const boundIDs = new Set(bindingRows.value.map(item => item.exitIpId));
  const scope = bindingForm.scope;
  return bindingExitIps.value
    .filter(item => {
      if (!route) return false;
      if (scope === 'external') return item.sourceType === 's5';
      if (item.sourceType !== 'node') return false;
      return scope === 'relay' ? item.ownerNodeId === route.relayNodeId : item.ownerNodeId === route.landingNodeId;
    })
    .filter(item => !boundIDs.has(item.id))
    .map(item => ({ label: `${item.address} · ${item.ownerNodeName || '独立 S5'} · ${item.provider || '未填写服务商'}`, value: item.id }));
});
const bindingModalTitle = computed(() => (bindingRoute.value ? `出口 IP 配置 · ${bindingRoute.value.name}` : '出口 IP 配置'));

watch(() => bindingForm.scope, () => {
  bindingForm.exitIpId = '';
});

function clearBindingDrafts() {
  Object.keys(bindingDrafts).forEach(id => Reflect.deleteProperty(bindingDrafts, id));
}

function syncBindingDrafts() {
  clearBindingDrafts();
  bindingRows.value.forEach(item => {
    bindingDrafts[item.id] = { scope: item.scope, allocationWeight: item.allocationWeight, enabled: item.enabled };
  });
}

async function loadBindingData() {
  if (!bindingRoute.value) return;
  bindingLoading.value = true;
  bindingError.value = '';
  const [bindingResult, exitIpResult] = await Promise.all([
    fetchRouteExitIpBindings(bindingRoute.value.id),
    fetchExitIps({ page: 1, page_size: 200, status: 'active' })
  ]);
  if (bindingResult.error || !bindingResult.data || exitIpResult.error || !exitIpResult.data) {
    bindingError.value = '无法读取线路出口 IP 配置，请刷新后重试。';
    bindingRows.value = [];
    bindingExitIps.value = [];
  } else {
    bindingRows.value = bindingResult.data;
    bindingExitIps.value = exitIpResult.data.items;
    syncBindingDrafts();
  }
  bindingLoading.value = false;
}

async function openBindings(row: Api.Central.RouteSummary) {
  bindingRoute.value = row;
  bindingVisible.value = true;
  bindingForm.exitIpId = '';
  bindingForm.scope = 'landing';
  bindingForm.allocationWeight = 1;
  await loadBindingData();
}

async function saveBinding() {
  if (!bindingRoute.value || !bindingForm.exitIpId) {
    bindingError.value = '请选择要绑定的出口 IP';
    return;
  }
  if (bindingForm.allocationWeight === null || !Number.isFinite(bindingForm.allocationWeight) || bindingForm.allocationWeight <= 0) {
    bindingError.value = '分配权重必须大于 0';
    return;
  }
  bindingSaving.value = true;
  bindingError.value = '';
  const { error } = await bindRouteExitIp(bindingRoute.value.id, {
    exitIpId: bindingForm.exitIpId,
    scope: bindingForm.scope,
    allocationWeight: bindingForm.allocationWeight,
    enabled: true
  });
  bindingSaving.value = false;
  if (error) return;
  bindingForm.exitIpId = '';
  bindingForm.scope = 'landing';
  bindingForm.allocationWeight = 1;
  window.$message?.success('出口 IP 已绑定');
  await loadBindingData();
  await loadRoutes();
}

async function saveBindingRow(item: Api.Central.RouteExitIpBinding) {
  if (!bindingRoute.value) return;
  const draft = bindingDrafts[item.id];
  if (!draft || draft.allocationWeight === null || !Number.isFinite(draft.allocationWeight) || draft.allocationWeight <= 0) {
    bindingError.value = '分配权重必须大于 0';
    return;
  }
  bindingUpdatingId.value = item.id;
  bindingError.value = '';
  const { error } = await updateRouteExitIp(bindingRoute.value.id, item.exitIpId, {
    scope: draft.scope,
    allocationWeight: draft.allocationWeight,
    enabled: draft.enabled
  });
  bindingUpdatingId.value = '';
  if (error) return;
  window.$message?.success('绑定配置已更新');
  await loadBindingData();
  await loadRoutes();
}

async function removeBinding(item: Api.Central.RouteExitIpBinding) {
  if (!bindingRoute.value || bindingDeletingId.value) return;
  bindingDeletingId.value = item.id;
  const { error } = await unbindRouteExitIp(bindingRoute.value.id, item.exitIpId);
  bindingDeletingId.value = '';
  if (error) return;
  window.$message?.success('出口 IP 已解绑');
  await loadBindingData();
  await loadRoutes();
}

function confirmRemoveBinding(item: Api.Central.RouteExitIpBinding) {
  const remove = () => void removeBinding(item);
  if (!window.$dialog) {
    remove();
    return;
  }
  window.$dialog.warning({
    title: '解除出口 IP 绑定',
    content: `确定解除 ${item.address} 与线路的绑定吗？`,
    positiveText: '解除绑定',
    negativeText: '取消',
    onPositiveClick: remove
  });
}

function resetRouteForm() {
  editingRouteId.value = null;
  routeForm.name = '';
  routeForm.relayNodeId = '';
  routeForm.landingNodeId = '';
  routeForm.relayOutboundTag = '';
  routeForm.landingInboundId = '';
  routeForm.landingInboundTag = '';
  routeForm.enabled = true;
  routeForm.validFrom = null;
  routeForm.validTo = null;
  routeForm.notes = '';
  formError.value = '';
}

function copyRouteToForm(route: Api.Central.RouteDetail) {
  routeForm.name = route.name;
  routeForm.relayNodeId = route.relayNodeId;
  routeForm.landingNodeId = route.landingNodeId;
  routeForm.relayOutboundTag = route.relayOutboundTag || '';
  routeForm.landingInboundId = route.landingInboundId || '';
  routeForm.landingInboundTag = route.landingInboundTag || '';
  routeForm.enabled = route.enabled;
  routeForm.validFrom = route.validFrom ? route.validFrom.slice(0, 10) : null;
  routeForm.validTo = route.validTo ? route.validTo.slice(0, 10) : null;
  routeForm.notes = route.notes || '';
}

async function loadNodes() {
  nodesLoading.value = true;
  nodeError.value = '';
  const { data, error } = await fetchNodes({ page: 1, page_size: 200 });
  if (error || !data) {
    nodeError.value = '无法读取节点列表，请确认中央后端已启动后重试。';
    nodes.value = [];
  } else {
    nodes.value = data.items;
  }
  nodesLoading.value = false;
}

function openCreate() {
  resetRouteForm();
  modalVisible.value = true;
}

async function openEdit(id: string) {
  resetRouteForm();
  editingRouteId.value = id;
  modalVisible.value = true;
  modalLoading.value = true;
  const { data, error } = await fetchRouteDetail(id);
  modalLoading.value = false;
  if (error || !data) {
    formError.value = '无法读取线路详情，请刷新后重试。';
    return;
  }
  copyRouteToForm(data);
}

function closeModal() {
  if (!saving.value) modalVisible.value = false;
}

async function saveRoute() {
  const name = routeForm.name.trim();
  if (!name) {
    formError.value = '请填写线路名称';
    return;
  }
  if (!routeForm.relayNodeId || !routeForm.landingNodeId) {
    formError.value = '请选择线路机和落地机';
    return;
  }
  if (routeForm.relayNodeId === routeForm.landingNodeId) {
    formError.value = '线路机和落地机不能相同';
    return;
  }
  if (routeForm.validFrom && routeForm.validTo && routeForm.validTo < routeForm.validFrom) {
    formError.value = '失效日期不能早于生效日期';
    return;
  }
  formError.value = '';
  saving.value = true;
  const payload: Api.Central.RouteUpsertPayload = {
    name,
    relayNodeId: routeForm.relayNodeId,
    landingNodeId: routeForm.landingNodeId,
    // Send empty strings for cleared optional fields so PATCH can clear an
    // existing value (the backend treats them as NULL when persisting).
    relayOutboundTag: routeForm.relayOutboundTag.trim(),
    landingInboundId: routeForm.landingInboundId.trim(),
    landingInboundTag: routeForm.landingInboundTag.trim(),
    enabled: routeForm.enabled,
    validFrom: routeForm.validFrom || '',
    validTo: routeForm.validTo || '',
    notes: routeForm.notes.trim()
  };
  const result = editingRouteId.value ? await updateRoute(editingRouteId.value, payload) : await createRoute(payload);
  saving.value = false;
  if (result.error || !result.data) {
    formError.value = '线路保存失败，请检查节点关系和日期后重试。';
    return;
  }
  modalVisible.value = false;
  window.$message?.success(editingRouteId.value ? '线路已更新' : '线路已创建');
  await loadRoutes();
}

async function removeRoute(id: string) {
  if (deletingRouteId.value) return;
  deletingRouteId.value = id;
  const { error } = await deleteRoute(id);
  deletingRouteId.value = '';
  if (error) return;
  window.$message?.success('线路已删除');
  await loadRoutes();
}

function confirmDelete(row: Api.Central.RouteSummary) {
  const remove = () => void removeRoute(row.id);
  if (!window.$dialog) {
    remove();
    return;
  }
  window.$dialog.warning({
    title: '删除线路关系',
    content: `确定删除“${row.name}”吗？已绑定用户或出口 IP 的线路需要先停用。`,
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: remove
  });
}

const columns: DataTableColumns<Api.Central.RouteSummary> = [
  { title: '线路', key: 'name', minWidth: 160, render: row => h('span', { class: 'font-medium' }, row.name) },
  { title: '线路机', key: 'relayNodeName', minWidth: 140 },
  { title: '落地机', key: 'landingNodeName', minWidth: 140 },
  { title: '落地 SS Inbound', key: 'landingInboundTag', minWidth: 170, render: row => row.landingInboundTag || '--' },
  { title: '出口 IP 数', key: 'exitIpCount', width: 110 },
  { title: '配置归属用户数', key: 'allocatedUserCount', width: 150 },
  {
    title: '状态',
    key: 'status',
    width: 90,
    render: row =>
      h(
        NTag,
        { size: 'small', type: row.status === 'active' ? 'success' : 'default' },
        { default: () => statusLabel(row.status) }
      )
  },
  {
    title: '操作',
    key: 'actions',
    width: 220,
    fixed: 'right',
    render: row =>
      h(
        NSpace,
        { size: 4 },
        {
          default: () => [
            h(NButton, { size: 'small', type: 'primary', text: true, onClick: () => void openEdit(row.id) }, { default: () => '编辑' }),
            h(NButton, { size: 'small', type: 'info', text: true, onClick: () => void openBindings(row) }, { default: () => '出口 IP' }),
            h(
              NButton,
              { size: 'small', type: 'error', text: true, loading: deletingRouteId.value === row.id, disabled: Boolean(deletingRouteId.value) && deletingRouteId.value !== row.id, onClick: () => confirmDelete(row) },
              { default: () => '删除' }
            )
          ]
        }
      )
  }
];

const bindingColumns: DataTableColumns<Api.Central.RouteExitIpBinding> = [
  { title: '出口 IP', key: 'address', minWidth: 170, render: row => h('span', { class: 'font-medium' }, row.address) },
  { title: '协议族', key: 'family', width: 90, render: row => (row.family === 6 ? 'IPv6' : 'IPv4') },
  {
    title: '出口位置',
    key: 'scope',
    width: 150,
    render: row => {
      const draft = bindingDrafts[row.id];
      return h(NSelect, {
        value: draft?.scope || row.scope,
        options: bindingScopeOptions,
        size: 'small',
        onUpdateValue: (value: Api.Central.RouteExitIpScope) => {
          if (draft) draft.scope = value;
        }
      });
    }
  },
  { title: '来源', key: 'ownerNodeName', minWidth: 150, render: row => row.sourceType === 's5' ? '独立 S5' : row.ownerNodeName || '--' },
  {
    title: '分配权重',
    key: 'allocationWeight',
    width: 140,
    render: row => {
      const draft = bindingDrafts[row.id];
      return h(NInputNumber, {
        value: draft?.allocationWeight ?? row.allocationWeight,
        min: 0.01,
        max: 1000000,
        precision: 2,
        size: 'small',
        onUpdateValue: (value: number | null) => {
          if (draft) draft.allocationWeight = value;
        }
      });
    }
  },
  {
    title: '状态',
    key: 'enabled',
    width: 90,
    render: row => {
      const draft = bindingDrafts[row.id];
      return h(NSwitch, {
        value: draft?.enabled ?? row.enabled,
        size: 'small',
        onUpdateValue: (value: boolean) => {
          if (draft) draft.enabled = value;
        }
      });
    }
  },
  {
    title: '操作',
    key: 'actions',
    width: 170,
    fixed: 'right',
    render: row =>
      h(
        NSpace,
        { size: 4 },
        {
          default: () => [
            h(NButton, { size: 'small', type: 'primary', text: true, loading: bindingUpdatingId.value === row.id, onClick: () => void saveBindingRow(row) }, { default: () => '保存' }),
            h(NButton, { size: 'small', type: 'error', text: true, loading: bindingDeletingId.value === row.id, disabled: Boolean(bindingDeletingId.value) && bindingDeletingId.value !== row.id, onClick: () => confirmRemoveBinding(row) }, { default: () => '解绑' })
          ]
        }
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
    loadRoutes();
  },
  onUpdatePageSize: (pageSize: number) => {
    filters.page_size = pageSize;
    filters.page = 1;
    loadRoutes();
  }
}));

async function loadRoutes() {
  loading.value = true;
  errorMessage.value = '';
  const { data, error } = await fetchRoutes({
    page: filters.page,
    page_size: filters.page_size,
    keyword: filters.keyword || undefined,
    status: filters.status === 'all' ? undefined : filters.status
  });
  if (error || !data) {
    errorMessage.value = '中央后端暂不可用，无法读取线路数据';
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
  loadRoutes();
}

function resetFilters() {
  filters.keyword = '';
  filters.status = 'all';
  submitFilters();
}

onMounted(() => {
  void Promise.all([loadRoutes(), loadNodes()]);
});
</script>

<template>
  <div class="routes-page">
    <ModulePage
      title="线路关系"
      description="维护线路机、落地机、SS Inbound，以及线路机直出、落地机和外部 S5 出口关系。"
      :loading="loading"
      :error="errorMessage"
      :empty="rows.length === 0"
      empty-description="暂无线路数据；创建线路后，可在对应行的“出口 IP”操作中进行绑定"
      :data-at="dataAt"
      @refresh="loadRoutes"
    >
      <template #actions>
        <NButton size="small" type="primary" @click="openCreate">
          <template #icon><icon-mdi-plus /></template>
          新建线路
        </NButton>
      </template>
      <template #toolbar>
        <div class="border-b border-gray-200 p-16px dark:border-gray-700">
          <NSpace wrap>
            <NInput
              v-model:value="filters.keyword"
              clearable
              class="w-240px"
              placeholder="搜索线路、节点或 Inbound"
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
        :scroll-x="980"
      />
    </ModulePage>

    <NModal v-model:show="modalVisible" preset="card" :title="modalTitle" class="w-720px max-w-95vw" :mask-closable="false">
      <NSpin :show="modalLoading">
        <NAlert v-if="nodeError" type="warning" :show-icon="true" class="mb-16px">
          {{ nodeError }}
          <NButton size="small" class="ml-8px" @click="loadNodes">重试</NButton>
        </NAlert>
        <NAlert v-if="formError" type="error" :show-icon="true" class="mb-16px">{{ formError }}</NAlert>
        <NForm label-placement="top">
          <div class="grid gap-x-12px md:grid-cols-2">
            <NFormItem label="线路名称" required>
              <NInput v-model:value="routeForm.name" maxlength="120" show-count placeholder="例如：东京中转线路" />
            </NFormItem>
            <NFormItem label="状态">
              <NSwitch v-model:value="routeForm.enabled"><template #checked>启用</template><template #unchecked>停用</template></NSwitch>
            </NFormItem>
            <NFormItem label="线路机" required>
              <NSelect v-model:value="routeForm.relayNodeId" :options="relayNodeOptions" :loading="nodesLoading" filterable clearable placeholder="选择线路机" />
            </NFormItem>
            <NFormItem label="落地机" required>
              <NSelect v-model:value="routeForm.landingNodeId" :options="landingNodeOptions" :loading="nodesLoading" filterable clearable placeholder="选择落地机" />
            </NFormItem>
            <NFormItem label="线路机 Outbound Tag">
              <NInput v-model:value="routeForm.relayOutboundTag" maxlength="200" placeholder="可选" />
            </NFormItem>
            <NFormItem label="落地机 SS Inbound Tag">
              <NInput v-model:value="routeForm.landingInboundTag" maxlength="200" placeholder="可选" />
            </NFormItem>
            <NFormItem label="落地机 SS Inbound ID">
              <NInput v-model:value="routeForm.landingInboundId" maxlength="120" placeholder="可选，需属于落地机" />
            </NFormItem>
            <div />
            <NFormItem label="生效日期">
              <NDatePicker v-model:formatted-value="routeForm.validFrom" type="date" value-format="yyyy-MM-dd" clearable class="w-full" />
            </NFormItem>
            <NFormItem label="失效日期">
              <NDatePicker v-model:formatted-value="routeForm.validTo" type="date" value-format="yyyy-MM-dd" clearable class="w-full" />
            </NFormItem>
          </div>
          <NFormItem label="备注">
            <NInput v-model:value="routeForm.notes" type="textarea" maxlength="2000" show-count :autosize="{ minRows: 3, maxRows: 6 }" placeholder="可选" />
          </NFormItem>
        </NForm>
      </NSpin>
      <template #footer>
        <NSpace justify="end">
          <NButton :disabled="saving" @click="closeModal">取消</NButton>
          <NButton type="primary" :loading="saving" :disabled="modalLoading" @click="saveRoute">保存</NButton>
        </NSpace>
      </template>
    </NModal>

    <NModal v-model:show="bindingVisible" preset="card" :title="bindingModalTitle" class="w-900px max-w-95vw" :mask-closable="false">
      <NSpin :show="bindingLoading">
        <NAlert v-if="bindingError" type="error" :show-icon="true" class="mb-16px">{{ bindingError }}</NAlert>
        <NAlert type="info" :show-icon="true" class="mb-16px">
          线路机直出表示用户直接使用线路机公网 IP；落地机出口经过落地机；外部 S5 不属于任何节点。
        </NAlert>
        <NCard title="新增出口 IP 绑定" size="small" :segmented="{ content: true }">
          <NSpace align="center" wrap>
            <NSelect v-model:value="bindingForm.scope" :options="bindingScopeOptions" :disabled="bindingSaving" class="w-150px" />
            <NSelect v-model:value="bindingForm.exitIpId" :options="bindingExitIpOptions" :disabled="bindingSaving" :loading="bindingLoading" filterable clearable class="w-360px" placeholder="选择符合出口位置的资产" />
            <NInputNumber v-model:value="bindingForm.allocationWeight" :min="0.01" :max="1000000" :precision="2" :disabled="bindingSaving" class="w-140px" placeholder="分配权重" />
            <NButton type="primary" :loading="bindingSaving" :disabled="bindingLoading" @click="saveBinding">绑定</NButton>
          </NSpace>
        </NCard>
        <NDivider />
        <NEmpty v-if="bindingRows.length === 0" description="该线路尚未绑定出口 IP" size="small" />
        <NDataTable v-else :columns="bindingColumns" :data="bindingRows" :bordered="false" :single-line="false" size="small" :scroll-x="980" />
      </NSpin>
      <template #footer>
        <NSpace justify="end">
          <NButton @click="bindingVisible = false">关闭</NButton>
        </NSpace>
      </template>
    </NModal>
  </div>
</template>
