<script setup lang="ts">
import { computed, h, onMounted, reactive, ref, watch } from 'vue';
import { NButton, NSpace } from 'naive-ui';
import type { DataTableColumns } from 'naive-ui';
import { useRouter } from 'vue-router';
import { createNode, deleteNode, fetchNodes, requestNodeSync, updateNode } from '@/service/api';
import ModulePage from '@/components/project/module-page.vue';
import DataFreshness from '@/components/project/data-freshness.vue';
import NodeStatusTag from '@/components/project/node-status-tag.vue';
import NodeTypeTag from '@/components/project/node-type-tag.vue';

defineOptions({ name: 'NodeManagement' });

const loading = ref(false);
const syncingNodeId = ref('');
const togglingNodeId = ref('');
const deletingNodeId = ref('');
const errorMessage = ref('');
const rows = ref<Api.Central.NodeSummary[]>([]);
const total = ref(0);
const dataAt = ref('');
const selectedType = ref('all');
const onboardingVisible = ref(false);
const onboardingSaving = ref(false);
const onboardingError = ref('');
const onboardingResult = ref<Api.Central.NodeRegistrationResult | null>(null);
const nodeForm = reactive<Api.Central.NodeCreatePayload>({
  nodeKey: '',
  name: '',
  type: 'relay',
  hostname: '',
  managementUrl: '',
  region: '',
  provider: '',
  panelBasePath: '/'
});
const exitIpInputs = ref(['']);
const xpanelUrl = ref('');
const router = useRouter();
const filters = reactive({
  page: 1,
  page_size: 50,
  keyword: '',
  status: 'all',
  node_type: undefined as Api.Central.NodeType | undefined
});

const statusOptions = [
  { label: '全部状态', value: 'all' },
  { label: '在线', value: 'online' },
  { label: '异常', value: 'degraded' },
  { label: '离线', value: 'offline' },
  { label: '已停用', value: 'disabled' }
];
const typeOptions = [
  { label: '全部类型', value: 'all' },
  { label: '线路机', value: 'relay' },
  { label: '落地机', value: 'landing' }
];

function typeLabel(type: Api.Central.NodeType) {
  return typeOptions.find(option => option.value === type)?.label || '未知类型';
}

function formatDate(value?: string) {
  if (!value) return '--';
  const date = new Date(value);
  return Number.isNaN(date.getTime()) ? '时间未知' : date.toLocaleString('zh-CN', { hour12: false });
}

function formatMoney(value?: number, currency = 'CNY') {
  if (value === undefined || value === null) return '--';
  return new Intl.NumberFormat('zh-CN', { style: 'currency', currency }).format(value);
}

function parseManagementAddress(value: string) {
  try {
    const parsed = new URL(value.trim());
    if (
      !/^https?:$/i.test(parsed.protocol) ||
      !parsed.hostname ||
      parsed.username ||
      parsed.password ||
      parsed.search ||
      parsed.hash
    ) {
      return null;
    }
    const path = parsed.pathname ? `/${parsed.pathname.replace(/^\/+|\/+$/g, '')}` : '/';
    return { url: parsed.origin, basePath: path === '/' ? '/' : path };
  } catch {
    return null;
  }
}

watch(
  () => nodeForm.managementUrl,
  value => {
    const parsed = parseManagementAddress(value || '');
    xpanelUrl.value = parsed?.url || '';
    nodeForm.panelBasePath = parsed?.basePath || '/';
  }
);

function centralApiUrl() {
  const configuredCentralUrl = String(import.meta.env.VITE_SERVICE_BASE_URL || '/api').trim();
  return /^https?:\/\//i.test(configuredCentralUrl)
    ? configuredCentralUrl.replace(/\/$/, '')
    : new URL(configuredCentralUrl || '/api', window.location.origin).toString().replace(/\/$/, '');
}

function shellQuote(value: string) {
  return `'${value.replace(/'/g, `'\\''`)}'`;
}

const agentInstallCommand = computed(() => {
  if (!onboardingResult.value?.installerToken) return '';
  const installerUrl = String(
    import.meta.env.VITE_AGENT_INSTALL_URL ||
      'https://raw.githubusercontent.com/Arm0ne/Alien-Panel/main/agent/deploy/install-online.sh'
  ).trim();
  return [
    'curl -fsSL',
    shellQuote(installerUrl),
    '| sudo bash -s --',
    '--central-url',
    shellQuote(centralApiUrl()),
    '--install-token',
    shellQuote(onboardingResult.value.installerToken),
    '--xpanel-url',
    shellQuote(xpanelUrl.value.trim()),
    '--xpanel-base-path',
    shellQuote(nodeForm.panelBasePath?.trim() || '/')
  ].join(' ');
});

const agentConfig = computed(() => {
  if (!onboardingResult.value) return '';
  const centralUrl = centralApiUrl();
  return [
    `node_key: ${onboardingResult.value.nodeKey}`,
    `node_name: ${onboardingResult.value.name}`,
    `node_type: ${onboardingResult.value.type}`,
    `central_url: ${centralUrl}`,
    `central_token: ${onboardingResult.value.token}`,
    `xpanel_url: ${xpanelUrl.value.trim() || 'http://127.0.0.1:2053'}`,
    `xpanel_base_path: ${nodeForm.panelBasePath || '/'}`,
    'xpanel_username: admin',
    'xpanel_password: 请在节点本机填写',
    'sync_interval: 60s',
    'http_timeout: 15s'
  ].join('\n');
});

function resetNodeForm() {
  nodeForm.nodeKey = '';
  nodeForm.name = '';
  nodeForm.type = 'relay';
  nodeForm.hostname = '';
  nodeForm.managementUrl = '';
  nodeForm.region = '';
  nodeForm.provider = '';
  nodeForm.panelBasePath = '/';
  xpanelUrl.value = '';
  exitIpInputs.value = [''];
  onboardingError.value = '';
  onboardingResult.value = null;
}

function addExitIpInput() {
  if (exitIpInputs.value.length >= 100) {
    onboardingError.value = '最多可登记 100 个出口 IP';
    return;
  }
  exitIpInputs.value.push('');
}

function removeExitIpInput(index: number) {
  if (exitIpInputs.value.length === 1) {
    exitIpInputs.value[0] = '';
    return;
  }
  exitIpInputs.value.splice(index, 1);
}

function openOnboarding() {
  resetNodeForm();
  onboardingVisible.value = true;
}

function closeOnboarding() {
  if (!onboardingSaving.value) onboardingVisible.value = false;
}

async function copyText(value: string, successMessage: string) {
  try {
    await navigator.clipboard.writeText(value);
    window.$message?.success(successMessage);
  } catch {
    window.$message?.error('复制失败，请手动选择并复制');
  }
}

async function saveNode() {
  const managementUrl = nodeForm.managementUrl?.trim() || '';
  const parsedManagement = parseManagementAddress(managementUrl);
  if (!parsedManagement) {
    onboardingError.value = '请填写有效的管理地址，例如 https://server.example.com:18086/Alien/';
    return;
  }
  const exitIps = exitIpInputs.value.map(value => value.trim()).filter(Boolean);
  const uniqueExitIps = [...new Set(exitIps)];
  if (uniqueExitIps.length !== exitIps.length) {
    onboardingError.value = '出口 IP 列表中存在重复地址，请检查后重试';
    return;
  }
  if (uniqueExitIps.length > 100) {
    onboardingError.value = '最多可登记 100 个出口 IP，请拆分后重试';
    return;
  }
  const payload: Api.Central.NodeCreatePayload = {
    nodeKey: nodeForm.nodeKey?.trim() || '',
    name: nodeForm.name.trim(),
    type: nodeForm.type,
    hostname: nodeForm.hostname?.trim(),
    managementUrl,
    exitIps: uniqueExitIps,
    region: nodeForm.region?.trim(),
    provider: nodeForm.provider?.trim(),
    panelBasePath: parsedManagement.basePath
  };
  if (!payload.name) {
    onboardingError.value = '请填写节点名称';
    return;
  }
  onboardingError.value = '';
  onboardingSaving.value = true;
  const { data, error } = await createNode(payload);
  onboardingSaving.value = false;
  if (error || !data) {
    onboardingError.value = '节点创建失败，请检查节点信息后重试';
    return;
  }
  onboardingResult.value = data;
  window.$message?.success('节点已创建，请立即保存 Agent Token');
  await loadNodes();
}

async function toggleNode(row: Api.Central.NodeSummary) {
  if (togglingNodeId.value) return;
  const enabled = !row.enabled;
  const action = enabled ? '启用' : '停用';
  const apply = async () => {
    togglingNodeId.value = row.id;
    const { error } = await updateNode(row.id, { enabled });
    togglingNodeId.value = '';
    if (error) {
      window.$message?.error(`${action}节点失败，请稍后重试`);
      return;
    }
    window.$message?.success(`节点已${action}`);
    await loadNodes();
  };
  if (!enabled && window.$dialog) {
    window.$dialog.warning({
      title: '停用节点',
      content: `停用“${row.name}”后，Agent 心跳和同步请求会被拒绝。确认继续吗？`,
      positiveText: '停用',
      negativeText: '取消',
      onPositiveClick: apply
    });
    return;
  }
  await apply();
}

async function removeNode(id: string) {
  if (deletingNodeId.value) return;
  deletingNodeId.value = id;
  const { error } = await deleteNode(id);
  deletingNodeId.value = '';
  if (error) {
    // The shared request handler already displays the backend's actionable
    // conflict message (including the exact association counts). Only show a
    // fallback here for transport errors that do not have a response body.
    if (!error.response?.data?.msg) {
      window.$message?.error('节点删除失败，请稍后重试');
    }
    return;
  }
  window.$message?.success('节点及其关联配置已删除');
  await loadNodes();
}

function confirmDelete(row: Api.Central.NodeSummary) {
  const remove = () => void removeNode(row.id);
  if (!window.$dialog) {
    remove();
    return;
  }
  window.$dialog.warning({
    title: '删除节点',
    content: `删除“${row.name}”后，节点、Agent Token、同步记录、Inbound、出口 IP 和用户路径都会永久删除，关联用户本身不会删除。确定继续吗？`,
    positiveText: '永久删除',
    negativeText: '取消',
    onPositiveClick: remove
  });
}

const columns: DataTableColumns<Api.Central.NodeSummary> = [
  {
    title: '节点',
    key: 'name',
    minWidth: 260,
    render: row =>
      h('div', [
        h('div', { class: 'font-medium' }, row.name),
        row.managementUrl
          ? h(
              'a',
              {
                class: 'text-12px text-primary hover:underline break-all',
                href: row.managementUrl,
                target: '_blank',
                rel: 'noopener noreferrer',
                title: row.managementUrl
              },
              row.managementUrl
            )
          : h('div', { class: 'text-12px text-gray-500' }, '未设置管理地址')
      ])
  },
  {
    title: '类型',
    key: 'type',
    width: 100,
    render: row => h(NodeTypeTag, { type: row.type })
  },
  {
    title: '健康状态',
    key: 'status',
    width: 110,
    render: row => h(NodeStatusTag, { status: row.status })
  },
  { title: 'X-Panel', key: 'xpanelVersion', minWidth: 110, render: row => row.xpanelVersion || '--' },
  { title: 'Xray', key: 'xrayVersion', minWidth: 110, render: row => row.xrayVersion || '--' },
  { title: '出口 IP', key: 'exitIpCount', width: 90 },
  { title: '最近心跳', key: 'lastSeenAt', minWidth: 170, render: row => formatDate(row.lastSeenAt) },
  {
    title: '最近同步',
    key: 'lastSyncAt',
    minWidth: 190,
    render: row =>
      h('div', { class: 'flex flex-col items-start gap-2px' }, [
        h(DataFreshness, { dataAt: row.lastSyncAt, compact: true }),
        h('span', { class: 'text-12px text-gray-500' }, formatDate(row.lastSyncAt))
      ])
  },
  { title: '月成本', key: 'monthlyCost', minWidth: 110, render: row => formatMoney(row.monthlyCost, row.currency) },
  {
    title: '操作',
    key: 'actions',
    width: 330,
    fixed: 'right',
    render: row =>
      h(
        NSpace,
        { size: 4 },
        {
          default: () => [
            h(
              NButton,
              { size: 'small', type: 'primary', text: true, onClick: () => void openDetail(row.id) },
              { default: () => '详情' }
            ),
            h(
              NButton,
              {
                size: 'small',
                type: 'info',
                text: true,
                loading: syncingNodeId.value === row.id,
                disabled: !row.enabled || (Boolean(syncingNodeId.value) && syncingNodeId.value !== row.id),
                onClick: () => void syncNode(row.id)
              },
              { default: () => '立即同步' }
            ),
            h(
              NButton,
              {
                size: 'small',
                type: row.enabled ? 'warning' : 'success',
                text: true,
                loading: togglingNodeId.value === row.id,
                disabled: Boolean(togglingNodeId.value) && togglingNodeId.value !== row.id,
                onClick: () => void toggleNode(row)
              },
              { default: () => (row.enabled ? '停用' : '启用') }
            ),
            h(
              NButton,
              {
                size: 'small',
                type: 'error',
                text: true,
                loading: deletingNodeId.value === row.id,
                disabled: Boolean(deletingNodeId.value) && deletingNodeId.value !== row.id,
                onClick: () => confirmDelete(row)
              },
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
    loadNodes();
  },
  onUpdatePageSize: (pageSize: number) => {
    filters.page_size = pageSize;
    filters.page = 1;
    loadNodes();
  }
}));

async function loadNodes() {
  loading.value = true;
  errorMessage.value = '';
  const { data, error } = await fetchNodes({
    page: filters.page,
    page_size: filters.page_size,
    keyword: filters.keyword || undefined,
    status: filters.status === 'all' ? undefined : filters.status,
    node_type: filters.node_type
  });
  if (error || !data) {
    errorMessage.value = '中央后端暂不可用，无法读取节点数据';
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

function openDetail(id: string) {
  router.push({ name: 'nodes-detail', params: { id } });
}

async function syncNode(id: string) {
  if (syncingNodeId.value) return;
  syncingNodeId.value = id;
  const { data, error } = await requestNodeSync(id);
  syncingNodeId.value = '';
  if (error || !data) {
    window.$message?.error('同步请求发送失败，请检查节点状态后重试');
    return;
  }
  window.$message?.success('已提交立即同步请求，Agent 将在下一次周期同步中执行');
  await loadNodes();
}

function submitFilters() {
  filters.page = 1;
  filters.node_type = selectedType.value === 'all' ? undefined : (selectedType.value as Api.Central.NodeType);
  loadNodes();
}

function resetFilters() {
  filters.keyword = '';
  filters.status = 'all';
  selectedType.value = 'all';
  submitFilters();
}

onMounted(loadNodes);
</script>

<template>
  <div class="nodes-page">
    <ModulePage
      title="节点管理"
      description="查看线路机和落地机的健康状态、版本、同步时间与节点成本。"
      :loading="loading"
      :error="errorMessage"
      :empty="rows.length === 0"
      empty-description="暂无节点；点击“接入节点”创建中央记录，然后按模板配置 Agent"
      :data-at="dataAt"
      @refresh="loadNodes"
    >
      <template #actions>
        <NButton size="small" type="primary" @click="openOnboarding">
          <template #icon><icon-mdi-plus /></template>
          接入节点
        </NButton>
      </template>
      <template #toolbar>
        <div class="border-b border-gray-200 p-16px dark:border-gray-700">
          <NSpace wrap>
            <NInput
              v-model:value="filters.keyword"
              clearable
              class="w-220px"
              placeholder="搜索节点名称或地址"
              @keyup.enter="submitFilters"
            />
            <NSelect v-model:value="selectedType" :options="typeOptions" class="w-140px" />
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
        :scroll-x="1180"
      />
    </ModulePage>

    <NModal
      v-model:show="onboardingVisible"
      preset="card"
      :title="onboardingResult ? '节点接入完成' : '接入节点'"
      class="w-720px max-w-95vw"
      :mask-closable="false"
    >
      <NAlert v-if="onboardingError" type="error" :show-icon="true" class="mb-16px">{{ onboardingError }}</NAlert>
      <template v-if="onboardingResult">
        <NAlert type="success" :show-icon="true" class="mb-16px">
          节点已创建。请在目标 Ubuntu/Debian 服务器执行下面的一行命令，脚本会自动安装并启动 Agent。
        </NAlert>
        <NDescriptions bordered size="small" :column="1" class="mb-16px">
          <NDescriptionsItem label="节点 ID">{{ onboardingResult.nodeId }}</NDescriptionsItem>
          <NDescriptionsItem label="节点 Key">{{ onboardingResult.nodeKey }}</NDescriptionsItem>
          <NDescriptionsItem label="节点类型">{{ typeLabel(onboardingResult.type) }}</NDescriptionsItem>
          <NDescriptionsItem label="已登记出口 IP">{{ onboardingResult.exitIpCount }}</NDescriptionsItem>
        </NDescriptions>
        <NForm label-placement="top">
          <NFormItem label="Ubuntu / Debian 一键部署命令">
            <NSpace vertical :size="8" class="w-full">
              <NInput :value="agentInstallCommand" readonly type="textarea" :autosize="{ minRows: 2, maxRows: 4 }" />
              <NButton
                size="small"
                type="primary"
                class="self-start"
                :disabled="!agentInstallCommand"
                @click="copyText(agentInstallCommand, '一键部署命令已复制')"
              >
                复制部署命令
              </NButton>
            </NSpace>
          </NFormItem>
          <NFormItem label="节点 Token（手动部署备用）">
            <NSpace vertical :size="8" class="w-full">
              <NInput :value="onboardingResult.token" readonly type="textarea" :autosize="{ minRows: 2, maxRows: 4 }" />
              <NButton
                size="small"
                type="primary"
                class="self-start"
                @click="copyText(onboardingResult.token, '节点 Token 已复制')"
              >
                复制 Token
              </NButton>
            </NSpace>
          </NFormItem>
          <NFormItem label="Agent 配置模板（手动部署备用）">
            <NSpace vertical :size="8" class="w-full">
              <NInput :value="agentConfig" readonly type="textarea" :autosize="{ minRows: 10, maxRows: 18 }" />
              <NButton
                size="small"
                type="primary"
                class="self-start"
                @click="copyText(agentConfig, 'Agent 配置已复制')"
              >
                复制配置
              </NButton>
            </NSpace>
          </NFormItem>
        </NForm>
      </template>
      <template v-else>
        <NAlert type="info" :show-icon="true" class="mb-16px">
          选择节点类型并填写管理地址后，创建完成即可复制 Ubuntu/Debian 一键部署命令。节点 Key 由中央自动生成。
        </NAlert>
        <NForm label-placement="top">
          <div class="grid gap-x-12px md:grid-cols-2">
            <NFormItem label="节点名称" required>
              <NInput v-model:value="nodeForm.name" maxlength="120" placeholder="例如：东京线路机 01" />
            </NFormItem>
            <NFormItem label="节点类型" required>
              <NSelect
                v-model:value="nodeForm.type"
                :options="[
                  { label: '线路机', value: 'relay' },
                  { label: '落地机', value: 'landing' }
                ]"
              />
            </NFormItem>
            <NFormItem label="主机名">
              <NInput v-model:value="nodeForm.hostname" maxlength="255" placeholder="可选，例如 relay.example.com" />
            </NFormItem>
            <NFormItem label="管理地址" required>
              <NInput
                v-model:value="nodeForm.managementUrl"
                maxlength="2048"
                placeholder="例如 https://server.example.com:18086/Alien/"
              />
            </NFormItem>
            <NFormItem label="区域">
              <NInput v-model:value="nodeForm.region" maxlength="120" placeholder="可选，例如 东京" />
            </NFormItem>
            <NFormItem label="服务商">
              <NInput v-model:value="nodeForm.provider" maxlength="200" placeholder="可选，例如 Example Cloud" />
            </NFormItem>
            <NFormItem label="X-Panel 地址（自动解析）">
              <NInput :value="xpanelUrl || '--'" readonly disabled />
            </NFormItem>
            <NFormItem label="X-Panel API 路径（自动解析）">
              <NInput :value="nodeForm.panelBasePath || '--'" readonly disabled />
            </NFormItem>
            <NFormItem label="节点出口 IP（可选，可添加多个）" class="md:col-span-2">
              <NSpace vertical :size="8" class="w-full">
                <NSpace v-for="(value, index) in exitIpInputs" :key="index" align="center" :size="8" class="w-full">
                  <NInput
                    v-model:value="exitIpInputs[index]"
                    class="flex-1"
                    maxlength="45"
                    placeholder="例如 38.117.67.14 或 2001:db8::10"
                  />
                  <NButton
                    v-if="index === exitIpInputs.length - 1"
                    quaternary
                    circle
                    type="primary"
                    aria-label="增加出口 IP"
                    @click="addExitIpInput"
                  >
                    <template #icon><icon-mdi-plus /></template>
                  </NButton>
                  <NButton
                    v-else
                    quaternary
                    circle
                    type="error"
                    aria-label="删除出口 IP"
                    @click="removeExitIpInput(index)"
                  >
                    <template #icon><icon-mdi-minus /></template>
                  </NButton>
                </NSpace>
              </NSpace>
            </NFormItem>
          </div>
        </NForm>
      </template>
      <template #footer>
        <NSpace justify="end">
          <NButton :disabled="onboardingSaving" @click="closeOnboarding">关闭</NButton>
          <NButton v-if="!onboardingResult" type="primary" :loading="onboardingSaving" @click="saveNode">
            创建并生成 Token
          </NButton>
        </NSpace>
      </template>
    </NModal>
  </div>
</template>
