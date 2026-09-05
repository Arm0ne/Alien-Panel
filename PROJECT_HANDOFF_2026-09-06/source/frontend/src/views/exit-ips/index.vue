<script setup lang="ts">
import { computed, h, onMounted, reactive, ref, watch } from 'vue';
import { NButton, NSpace, NTag } from 'naive-ui';
import type { DataTableColumns } from 'naive-ui';
import { createExitIp, deleteExitIp, fetchExitIpDetail, fetchExitIps, fetchNodes, updateExitIp } from '@/service/api';
import ModulePage from '@/components/project/module-page.vue';

defineOptions({ name: 'ExitIpManagement' });

const loading = ref(false);
const errorMessage = ref('');
const rows = ref<Api.Central.ExitIpSummary[]>([]);
const total = ref(0);
const dataAt = ref('');
const filters = reactive({ page: 1, page_size: 50, keyword: '', status: 'all' });
const nodes = ref<Api.Central.NodeSummary[]>([]);
const nodesLoading = ref(false);
const nodeError = ref('');
const modalVisible = ref(false);
const modalLoading = ref(false);
const saving = ref(false);
const deletingExitIpId = ref('');
const formError = ref('');
const editingExitIpId = ref<string | null>(null);
const exitIpForm = reactive({
  address: '',
  sourceType: 'node' as 'node' | 's5',
  ownerNodeId: '',
  family: 4 as 4 | 6,
  provider: '',
  monthlyCost: 0 as number | null,
  enabled: true,
  validFrom: null as string | null,
  validTo: null as string | null,
  notes: ''
});

const statusOptions = [
  { label: '全部状态', value: 'all' },
  { label: '启用', value: 'active' },
  { label: '停用', value: 'disabled' }
];
const familyOptions = [
  { label: 'IPv4', value: 4 },
  { label: 'IPv6', value: 6 }
];
const sourceTypeOptions = [
  { label: '节点公网出口', value: 'node' },
  { label: '独立 S5 出口', value: 's5' }
];

function formatMoney(value?: number, currency = 'CNY') {
  if (value === undefined || value === null) return '--';
  return new Intl.NumberFormat('zh-CN', { style: 'currency', currency }).format(value);
}

function formatDate(value?: string | null) {
  if (!value) return '--';
  const date = new Date(value);
  return Number.isNaN(date.getTime()) ? value : date.toLocaleString('zh-CN', { hour12: false });
}

function familyLabel(family: number) {
  return family === 6 ? 'IPv6' : 'IPv4';
}

const ownerNodeOptions = computed(() =>
  nodes.value
    .filter(node => node.type === 'landing' || node.type === 'relay')
    .map(node => ({ label: `${node.name} · ${node.type === 'relay' ? '线路机' : '落地机'}（${node.host || '地址未知'}）`, value: node.id }))
);
const modalTitle = computed(() => (editingExitIpId.value ? '编辑出口 IP' : '新建出口 IP'));

watch(() => exitIpForm.sourceType, value => {
  if (value === 's5') exitIpForm.ownerNodeId = '';
});

function resetExitIpForm() {
  editingExitIpId.value = null;
  exitIpForm.address = '';
  exitIpForm.sourceType = 'node';
  exitIpForm.ownerNodeId = '';
  exitIpForm.family = 4;
  exitIpForm.provider = '';
  exitIpForm.monthlyCost = 0;
  exitIpForm.enabled = true;
  exitIpForm.validFrom = null;
  exitIpForm.validTo = null;
  exitIpForm.notes = '';
  formError.value = '';
}

function copyExitIpToForm(detail: Api.Central.ExitIpDetail) {
  exitIpForm.address = detail.address;
  exitIpForm.sourceType = detail.sourceType || 'node';
  exitIpForm.ownerNodeId = detail.ownerNodeId || detail.landingNodeId || '';
  exitIpForm.family = detail.family;
  exitIpForm.provider = detail.provider || '';
  exitIpForm.monthlyCost = detail.monthlyCost;
  exitIpForm.enabled = detail.enabled;
  exitIpForm.validFrom = detail.validFrom ? detail.validFrom.slice(0, 10) : null;
  exitIpForm.validTo = detail.validTo ? detail.validTo.slice(0, 10) : null;
  exitIpForm.notes = detail.notes || '';
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
  resetExitIpForm();
  modalVisible.value = true;
}

async function openEdit(id: string) {
  resetExitIpForm();
  editingExitIpId.value = id;
  modalVisible.value = true;
  modalLoading.value = true;
  const { data, error } = await fetchExitIpDetail(id);
  modalLoading.value = false;
  if (error || !data) {
    formError.value = '无法读取出口 IP 详情，请刷新后重试。';
    return;
  }
  copyExitIpToForm(data);
}

function closeModal() {
  if (!saving.value) modalVisible.value = false;
}

async function saveExitIp() {
  const address = exitIpForm.address.trim();
  if (!address) {
    formError.value = '请填写出口 IP 地址';
    return;
  }
  if (exitIpForm.sourceType === 'node' && !exitIpForm.ownerNodeId) {
    formError.value = '请选择所属节点';
    return;
  }
  if (exitIpForm.monthlyCost === null || !Number.isFinite(exitIpForm.monthlyCost) || exitIpForm.monthlyCost < 0) {
    formError.value = '月成本必须是不小于 0 的数字';
    return;
  }
  if (exitIpForm.validFrom && exitIpForm.validTo && exitIpForm.validTo < exitIpForm.validFrom) {
    formError.value = '失效日期不能早于生效日期';
    return;
  }
  formError.value = '';
  saving.value = true;
  const payload: Api.Central.ExitIpUpsertPayload = {
    address,
    sourceType: exitIpForm.sourceType,
    ownerNodeId: exitIpForm.sourceType === 'node' ? exitIpForm.ownerNodeId : null,
    family: exitIpForm.family,
    provider: exitIpForm.provider.trim(),
    monthlyCost: exitIpForm.monthlyCost,
    currency: 'CNY',
    enabled: exitIpForm.enabled,
    validFrom: exitIpForm.validFrom || '',
    validTo: exitIpForm.validTo || '',
    notes: exitIpForm.notes.trim()
  };
  const result = editingExitIpId.value ? await updateExitIp(editingExitIpId.value, payload) : await createExitIp(payload);
  saving.value = false;
  if (result.error || !result.data) {
    formError.value = '出口 IP 保存失败，请检查地址、来源和日期后重试。';
    return;
  }
  modalVisible.value = false;
  window.$message?.success(editingExitIpId.value ? '出口 IP 已更新' : '出口 IP 已创建');
  await loadExitIps();
}

async function removeExitIp(id: string) {
  if (deletingExitIpId.value) return;
  deletingExitIpId.value = id;
  const { error } = await deleteExitIp(id);
  deletingExitIpId.value = '';
  if (error) return;
  window.$message?.success('出口 IP 已删除');
  await loadExitIps();
}

function confirmDelete(row: Api.Central.ExitIpSummary) {
  const remove = () => void removeExitIp(row.id);
  if (!window.$dialog) {
    remove();
    return;
  }
  window.$dialog.warning({
    title: '删除出口 IP',
    content: `确定删除“${row.address}”吗？已绑定线路的出口 IP 需要先解除绑定或停用。`,
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: remove
  });
}

const columns: DataTableColumns<Api.Central.ExitIpSummary> = [
  { title: '出口 IP', key: 'address', minWidth: 170, render: row => h('span', { class: 'font-medium' }, row.address) },
  {
    title: '出口来源',
    key: 'ownerNodeName',
    minWidth: 180,
    render: row => row.sourceType === 's5' ? '独立 S5' : `${row.ownerNodeName || '--'}（${row.ownerNodeType === 'relay' ? '线路机' : '落地机'}）`
  },
  { title: '协议族', key: 'family', width: 90, render: row => familyLabel(row.family || 4) },
  { title: '服务商', key: 'provider', minWidth: 140, render: row => row.provider || '--' },
  {
    title: '状态',
    key: 'status',
    width: 90,
    render: row =>
      h(
        NTag,
        { size: 'small', type: row.status === 'active' ? 'success' : 'default' },
        { default: () => (row.status === 'active' ? '启用' : row.status === 'disabled' ? '停用' : '未知') }
      )
  },
  { title: '月成本', key: 'monthlyCost', minWidth: 110, render: row => formatMoney(row.monthlyCost, row.currency) },
  { title: '配置归属用户数', key: 'allocatedUserCount', width: 150 },
  { title: '最近检查', key: 'checkedAt', minWidth: 170, render: row => formatDate(row.checkedAt) },
  {
    title: '操作',
    key: 'actions',
    width: 150,
    fixed: 'right',
    render: row =>
      h(
        NSpace,
        { size: 4 },
        {
          default: () => [
            h(NButton, { size: 'small', type: 'primary', text: true, onClick: () => void openEdit(row.id) }, { default: () => '编辑' }),
            h(
              NButton,
              { size: 'small', type: 'error', text: true, loading: deletingExitIpId.value === row.id, disabled: Boolean(deletingExitIpId.value) && deletingExitIpId.value !== row.id, onClick: () => confirmDelete(row) },
              { default: () => '删除' }
            )
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
    loadExitIps();
  },
  onUpdatePageSize: (pageSize: number) => {
    filters.page_size = pageSize;
    filters.page = 1;
    loadExitIps();
  }
}));

async function loadExitIps() {
  loading.value = true;
  errorMessage.value = '';
  const { data, error } = await fetchExitIps({
    page: filters.page,
    page_size: filters.page_size,
    keyword: filters.keyword || undefined,
    status: filters.status === 'all' ? undefined : filters.status
  });
  if (error || !data) {
    errorMessage.value = '中央后端暂不可用，无法读取出口 IP 数据';
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
  loadExitIps();
}

function resetFilters() {
  filters.keyword = '';
  filters.status = 'all';
  submitFilters();
}

onMounted(() => {
  void Promise.all([loadExitIps(), loadNodes()]);
});
</script>

<template>
  <div class="exit-ip-page">
    <ModulePage
      title="出口 IP"
      description="维护线路机/落地机公网出口与独立 S5 资产、成本和线路配置归属关系。"
      :loading="loading"
      :error="errorMessage"
      :empty="rows.length === 0"
      empty-description="暂无出口 IP 数据"
      :data-at="dataAt"
      @refresh="loadExitIps"
    >
      <template #actions>
        <NButton size="small" type="primary" @click="openCreate">
          <template #icon><icon-mdi-plus /></template>
          新建出口 IP
        </NButton>
      </template>
      <template #toolbar>
        <div class="border-b border-gray-200 p-16px dark:border-gray-700">
          <NSpace wrap>
              <NInput
              v-model:value="filters.keyword"
              clearable
              class="w-240px"
              placeholder="搜索 IP、节点、S5 或服务商"
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
        :scroll-x="1240"
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
            <NFormItem label="出口 IP 地址" required>
              <NInput v-model:value="exitIpForm.address" maxlength="45" placeholder="例如：203.0.113.10 或 2001:db8::10" />
            </NFormItem>
            <NFormItem label="协议族" required>
              <NSelect v-model:value="exitIpForm.family" :options="familyOptions" />
            </NFormItem>
            <NFormItem label="出口来源" required>
              <NSelect v-model:value="exitIpForm.sourceType" :options="sourceTypeOptions" />
            </NFormItem>
            <NFormItem v-if="exitIpForm.sourceType === 'node'" label="所属节点" required>
              <NSelect v-model:value="exitIpForm.ownerNodeId" :options="ownerNodeOptions" :loading="nodesLoading" filterable clearable placeholder="选择线路机或落地机" />
            </NFormItem>
            <NFormItem v-else label="S5 说明">
              <NText depth="3">独立购买的 S5 只登记出口地址和成本；账号凭据不在中央面板保存。</NText>
            </NFormItem>
            <NFormItem label="状态">
              <NSwitch v-model:value="exitIpForm.enabled"><template #checked>启用</template><template #unchecked>停用</template></NSwitch>
            </NFormItem>
            <NFormItem label="服务商">
              <NInput v-model:value="exitIpForm.provider" maxlength="200" placeholder="可选" />
            </NFormItem>
            <NFormItem label="月成本（CNY）">
              <NInputNumber v-model:value="exitIpForm.monthlyCost" :min="0" :max="100000000" :precision="2" class="w-full" />
            </NFormItem>
            <NFormItem label="生效日期">
              <NDatePicker v-model:formatted-value="exitIpForm.validFrom" type="date" value-format="yyyy-MM-dd" clearable class="w-full" />
            </NFormItem>
            <NFormItem label="失效日期">
              <NDatePicker v-model:formatted-value="exitIpForm.validTo" type="date" value-format="yyyy-MM-dd" clearable class="w-full" />
            </NFormItem>
          </div>
          <NFormItem label="备注">
            <NInput v-model:value="exitIpForm.notes" type="textarea" maxlength="2000" show-count :autosize="{ minRows: 3, maxRows: 6 }" placeholder="可选" />
          </NFormItem>
        </NForm>
      </NSpin>
      <template #footer>
        <NSpace justify="end">
          <NButton :disabled="saving" @click="closeModal">取消</NButton>
          <NButton type="primary" :loading="saving" :disabled="modalLoading" @click="saveExitIp">保存</NButton>
        </NSpace>
      </template>
    </NModal>
  </div>
</template>
