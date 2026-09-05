<script setup lang="ts">
import { computed, h, onMounted, reactive, ref, watch } from 'vue';
import { NButton, NSpace, NTag } from 'naive-ui';
import type { DataTableColumns } from 'naive-ui';
import { useRoute, useRouter } from 'vue-router';
import {
  createNodeCost,
  fetchNodeDetail,
  issueNodeInstallToken,
  requestNodeSync,
  updateNode,
  updateNodeCost
} from '@/service/api';
import ModulePage from '@/components/project/module-page.vue';
import DataFreshness from '@/components/project/data-freshness.vue';
import NodeStatusTag from '@/components/project/node-status-tag.vue';
import NodeTypeTag from '@/components/project/node-type-tag.vue';
import SyncStatusTag from '@/components/project/sync-status-tag.vue';
import TrafficValue from '@/components/project/traffic-value.vue';

defineOptions({ name: 'NodeDetail' });

const route = useRoute();
const router = useRouter();
const nodeId = computed(() => String(route.params.id || ''));
const loading = ref(false);
const syncing = ref(false);
const errorMessage = ref('');
const detail = ref<Api.Central.NodeDetail | null>(null);
const costModalVisible = ref(false);
const costSaving = ref(false);
const costError = ref('');
const installerModalVisible = ref(false);
const installerLoading = ref(false);
const installerError = ref('');
const installerResult = ref<Api.Central.NodeInstallTokenResult | null>(null);
const installerXpanelUrl = ref('');
const installerBasePath = ref('/');
const managementModalVisible = ref(false);
const managementSaving = ref(false);
const managementError = ref('');
const managementUrl = ref('');
const editingCostId = ref<string | null>(null);
const costForm = reactive<Api.Central.NodeCostCreatePayload>({
  category: '',
  monthlyAmount: 0,
  currency: 'CNY',
  effectiveFrom: new Date().toISOString().slice(0, 10),
  effectiveTo: null,
  notes: ''
});

function formatDate(value?: string | null) {
  if (!value) return '--';
  const date = new Date(value);
  return Number.isNaN(date.getTime()) ? '时间未知' : date.toLocaleString('zh-CN', { hour12: false });
}

function formatBytes(value: number) {
  if (!Number.isFinite(value)) return '--';
  if (value < 1024) return `${value} B`;
  const units = ['KB', 'MB', 'GB', 'TB', 'PB'];
  let size = value;
  let index = -1;
  do {
    size /= 1024;
    index += 1;
  } while (size >= 1024 && index < units.length - 1);
  return `${size.toFixed(size >= 10 ? 1 : 2)} ${units[index]}`;
}

function formatMoney(value: number, currency = 'CNY') {
  if (!Number.isFinite(value)) return '--';
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

function isValidXPanelAddress(value: string) {
  try {
    const parsed = new URL(value.trim());
    return /^https?:$/i.test(parsed.protocol) && Boolean(parsed.hostname) && !parsed.username && !parsed.password && !parsed.search && !parsed.hash;
  } catch {
    return false;
  }
}

function formatPercent(value?: number | null) {
  if (value === undefined || value === null || !Number.isFinite(value)) return '--';
  return `${value.toFixed(1)}%`;
}

function formatResource(used?: number | null, total?: number | null) {
  if (used === undefined || used === null || total === undefined || total === null || total <= 0) return '--';
  return `${formatBytes(used)} / ${formatBytes(total)}`;
}

const nodeExitIpTitle = computed(() => (detail.value?.type === 'relay' ? '线路机直出 IP' : '落地机出口 IP'));
const hasManagementUrl = computed(() => Boolean(detail.value?.managementUrl));

function centralApiUrl() {
  const configuredCentralUrl = String(import.meta.env.VITE_SERVICE_BASE_URL || '/api').trim();
  return /^https?:\/\//i.test(configuredCentralUrl)
    ? configuredCentralUrl.replace(/\/$/, '')
    : new URL(configuredCentralUrl || '/api', window.location.origin).toString().replace(/\/$/, '');
}

function shellQuote(value: string) {
  return `'${value.replace(/'/g, `'\\''`)}'`;
}

async function copyText(value: string, successMessage: string) {
  try {
    await navigator.clipboard.writeText(value);
    window.$message?.success(successMessage);
  } catch {
    window.$message?.error('复制失败，请手动选择并复制');
  }
}

const installerCommand = computed(() => {
  if (!installerResult.value?.installerToken) return '';
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
    shellQuote(installerResult.value.installerToken),
    '--xpanel-url',
    shellQuote(installerXpanelUrl.value.trim()),
    '--xpanel-base-path',
    shellQuote(installerBasePath.value.trim() || '/')
  ].join(' ');
});

function inboundStatus(value: Api.Central.NodeInboundStatus) {
  if (value === 'active') return { label: '启用', type: 'success' as const };
  if (value === 'disabled') return { label: '停用', type: 'default' as const };
  return { label: '已归档', type: 'warning' as const };
}

function eventType(value: Api.Central.EventSeverity) {
  if (value === 'error') return 'error' as const;
  if (value === 'warning') return 'warning' as const;
  return 'info' as const;
}

const inboundColumns: DataTableColumns<Api.Central.NodeInboundDetail> = [
  {
    title: 'Tag / Remark',
    key: 'tag',
    minWidth: 220,
    render: row =>
      h('div', [
        h('div', { class: 'font-medium' }, row.tag || row.remoteId),
        h('div', { class: 'text-12px text-gray-500' }, row.remark || '--')
      ])
  },
  { title: '协议 / 端口', key: 'protocol', width: 130, render: row => `${row.protocol || '--'} / ${row.port || '--'}` },
  { title: '类型', key: 'kind', width: 120 },
  { title: '设备数', key: 'clientCount', width: 90 },
  { title: '累计流量', key: 'allTime', minWidth: 120, render: row => h(TrafficValue, { value: row.allTime }) },
  { title: '到期时间', key: 'expiresAt', minWidth: 170, render: row => formatDate(row.expiresAt) },
  {
    title: '状态',
    key: 'status',
    width: 100,
    render: row => {
      const value = inboundStatus(row.status);
      return h(NTag, { size: 'small', type: value.type }, { default: () => value.label });
    }
  },
  { title: '最近同步', key: 'lastSeenAt', minWidth: 170, render: row => formatDate(row.lastSeenAt) }
];

const exitIpColumns: DataTableColumns<Api.Central.NodeExitIpDetail> = [
  { title: '地址', key: 'address', minWidth: 180 },
  { title: '提供商', key: 'provider', minWidth: 150, render: row => row.provider || '--' },
  { title: '协议族', key: 'family', width: 90, render: row => `IPv${row.family}` },
  {
    title: '状态',
    key: 'status',
    width: 90,
    render: row =>
      h(
        NTag,
        { size: 'small', type: row.status === 'active' ? 'success' : 'default' },
        { default: () => (row.status === 'active' ? '启用' : '停用') }
      )
  },
  { title: '月成本', key: 'monthlyCost', minWidth: 120, render: row => formatMoney(row.monthlyCost, row.currency) },
  {
    title: '有效期',
    key: 'validFrom',
    minWidth: 190,
    render: row => `${formatDate(row.validFrom)} ~ ${formatDate(row.validTo)}`
  }
];

const costColumns: DataTableColumns<Api.Central.NodeCostRecord> = [
  { title: '成本类别', key: 'category', minWidth: 150 },
  { title: '月金额', key: 'monthlyAmount', minWidth: 130, render: row => formatMoney(row.monthlyAmount, row.currency) },
  {
    title: '生效区间',
    key: 'effectiveFrom',
    minWidth: 210,
    render: row => `${row.effectiveFrom} ~ ${row.effectiveTo || '长期有效'}`
  },
  { title: '备注', key: 'notes', minWidth: 180, render: row => row.notes || '--' },
  {
    title: '操作',
    key: 'actions',
    width: 80,
    render: row =>
      h(
        NButton,
        { size: 'small', type: 'primary', text: true, onClick: () => openEditCost(row) },
        { default: () => '编辑' }
      )
  }
];

const costModalTitle = computed(() => (editingCostId.value ? '编辑节点成本' : '新增节点成本'));

function resetCostForm() {
  editingCostId.value = null;
  costForm.category = '';
  costForm.monthlyAmount = 0;
  costForm.currency = 'CNY';
  costForm.effectiveFrom = new Date().toISOString().slice(0, 10);
  costForm.effectiveTo = null;
  costForm.notes = '';
  costError.value = '';
}

function openCreateCost() {
  resetCostForm();
  costModalVisible.value = true;
}

function openEditCost(record: Api.Central.NodeCostRecord) {
  editingCostId.value = record.id;
  costForm.category = record.category;
  costForm.monthlyAmount = record.monthlyAmount;
  costForm.currency = record.currency;
  costForm.effectiveFrom = record.effectiveFrom;
  costForm.effectiveTo = record.effectiveTo;
  costForm.notes = record.notes || '';
  costError.value = '';
  costModalVisible.value = true;
}

function closeCostModal() {
  if (!costSaving.value) costModalVisible.value = false;
}

async function saveCost() {
  const category = costForm.category.trim();
  if (!category) {
    costError.value = '请填写成本类别';
    return;
  }
  if (!Number.isFinite(costForm.monthlyAmount) || costForm.monthlyAmount < 0) {
    costError.value = '月金额必须是不小于 0 的数字';
    return;
  }
  if (!costForm.effectiveFrom) {
    costError.value = '请选择生效日期';
    return;
  }
  if (costForm.effectiveTo && costForm.effectiveTo < costForm.effectiveFrom) {
    costError.value = '失效日期不能早于生效日期';
    return;
  }
  if (!detail.value) return;
  costError.value = '';
  costSaving.value = true;
  const payload: Api.Central.NodeCostCreatePayload = {
    category,
    monthlyAmount: costForm.monthlyAmount,
    currency: 'CNY',
    effectiveFrom: costForm.effectiveFrom,
    effectiveTo: costForm.effectiveTo || null,
    notes: costForm.notes?.trim() || null
  };
  const result = editingCostId.value
    ? await updateNodeCost(detail.value.id, editingCostId.value, payload)
    : await createNodeCost(detail.value.id, payload);
  costSaving.value = false;
  if (result.error || !result.data) {
    costError.value = result.error
      ? '成本保存失败；如果修改了生效日期，请新增一条成本记录。'
      : '成本保存失败，请检查表单后重试。';
    return;
  }
  costModalVisible.value = false;
  window.$message?.success(editingCostId.value ? '节点成本已更新' : '节点成本已新增');
  await loadNode();
}

async function loadNode() {
  if (!nodeId.value) {
    errorMessage.value = '节点地址无效';
    detail.value = null;
    return;
  }
  loading.value = true;
  errorMessage.value = '';
  const { data, error } = await fetchNodeDetail(nodeId.value);
  if (error || !data) {
    detail.value = null;
    errorMessage.value = '无法读取节点详情，请确认中央后端已启动后重试。';
  } else {
    detail.value = data;
  }
  loading.value = false;
}

async function syncNode() {
  if (!detail.value || syncing.value) return;
  syncing.value = true;
  const { data, error } = await requestNodeSync(detail.value.id);
  syncing.value = false;
  if (error || !data) {
    window.$message?.error('同步请求发送失败，请检查节点状态后重试');
    return;
  }
  window.$message?.success('已提交立即同步请求，Agent 将在下一次周期同步中执行');
  await loadNode();
}

function openInstaller() {
  if (!detail.value) return;
  installerError.value = '';
  installerResult.value = null;
  const parsed = parseManagementAddress(detail.value.managementUrl || '');
  installerXpanelUrl.value = parsed?.url || '';
  installerBasePath.value = parsed?.basePath || detail.value.panelBasePath || '/';
  installerModalVisible.value = true;
}

function closeInstaller() {
  if (!installerLoading.value) installerModalVisible.value = false;
}

async function generateInstaller() {
  if (!detail.value || installerLoading.value) return;
  const xpanelUrl = installerXpanelUrl.value.trim();
  const basePath = installerBasePath.value.trim() || '/';
  if (!isValidXPanelAddress(xpanelUrl)) {
    installerError.value = 'X-Panel 地址必须包含 http:// 或 https://，例如 http://127.0.0.1:2053';
    return;
  }
  if (!basePath.startsWith('/') || basePath.includes('://') || /\s/.test(basePath)) {
    installerError.value = 'X-Panel API 路径必须以 / 开头，例如 /Alien';
    return;
  }
  installerError.value = '';
  installerResult.value = null;
  installerLoading.value = true;
  const { data, error } = await issueNodeInstallToken(detail.value.id);
  installerLoading.value = false;
  if (error || !data) {
    installerError.value = '安装 Token 生成失败，请确认节点仍存在且中央后端可用。';
    return;
  }
  installerXpanelUrl.value = xpanelUrl;
  installerBasePath.value = basePath;
  installerResult.value = data;
}

function openManagementEditor() {
  if (!detail.value) return;
  managementError.value = '';
  managementUrl.value = detail.value.managementUrl || '';
  managementModalVisible.value = true;
}

function closeManagementEditor() {
  if (!managementSaving.value) managementModalVisible.value = false;
}

async function saveManagementEditor() {
  const value = managementUrl.value.trim();
  const parsed = parseManagementAddress(value);
  if (!parsed) {
    managementError.value = '请填写有效的管理地址，例如 https://server.example.com:18086/Alien/';
    return;
  }
  if (!detail.value || managementSaving.value) return;
  managementError.value = '';
  managementSaving.value = true;
  const { error } = await updateNode(detail.value.id, {
    managementUrl: value,
    panelBasePath: parsed.basePath
  });
  managementSaving.value = false;
  if (error) {
    managementError.value = '管理地址保存失败，请检查地址后重试。';
    return;
  }
  managementModalVisible.value = false;
  window.$message?.success('管理地址已更新');
  await loadNode();
}

function goBack() {
  router.push({ name: 'nodes' });
}

watch(nodeId, () => void loadNode());
onMounted(() => void loadNode());
</script>

<template>
  <div class="node-detail-page">
    <ModulePage
      title="节点详情"
      description="查看节点元数据、Inbound、节点公网出口 IP、节点成本、同步运行和最近状态事件。"
      :loading="loading"
      :error="errorMessage"
      :empty="!detail"
      empty-description="暂无节点详情"
      :data-at="detail?.dataAt || ''"
      @refresh="loadNode"
    >
      <template #actions>
        <NButton size="small" @click="goBack">返回节点列表</NButton>
        <NButton
          v-if="detail"
          size="small"
          type="primary"
          :loading="syncing"
          :disabled="!detail.enabled"
          @click="syncNode"
        >
          <template #icon><icon-mdi-sync /></template>
          立即同步
        </NButton>
        <NButton v-if="detail" size="small" :disabled="!detail.enabled" @click="openInstaller">
          重新生成部署命令
        </NButton>
      </template>

      <template v-if="detail">
        <div class="grid gap-16px lg:grid-cols-4">
          <NCard size="small" class="lg:col-span-3">
            <div class="flex flex-wrap items-start justify-between gap-12px">
              <div>
                <div class="text-20px font-600">{{ detail.name }}</div>
                <a
                  v-if="detail.managementUrl"
                  class="mt-4px inline-block break-all text-13px text-primary hover:underline"
                  :href="detail.managementUrl"
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  {{ detail.managementUrl }}
                </a>
                <div v-else class="mt-4px text-13px text-gray-500">未设置管理地址</div>
              </div>
              <NSpace>
                <NodeStatusTag :status="detail.status" />
                <NodeTypeTag :type="detail.type" />
                <NButton size="small" secondary @click="openManagementEditor">编辑管理地址</NButton>
              </NSpace>
            </div>
            <NDescriptions class="mt-16px" label-placement="left" :column="2" bordered size="small">
              <NDescriptionsItem label="节点 Key">{{ detail.nodeKey }}</NDescriptionsItem>
              <NDescriptionsItem label="启用状态">{{ detail.enabled ? '启用' : '停用' }}</NDescriptionsItem>
              <NDescriptionsItem label="主机名">{{ detail.hostname || '--' }}</NDescriptionsItem>
              <NDescriptionsItem label="管理地址">
                <a
                  v-if="detail.managementUrl"
                  class="break-all text-primary hover:underline"
                  :href="detail.managementUrl"
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  {{ detail.managementUrl }}
                </a>
                <span v-else>--</span>
              </NDescriptionsItem>
              <NDescriptionsItem label="区域">{{ detail.region || '--' }}</NDescriptionsItem>
              <NDescriptionsItem label="提供商">{{ detail.provider || '--' }}</NDescriptionsItem>
              <NDescriptionsItem label="Agent 版本">{{ detail.agentVersion || '未上报' }}</NDescriptionsItem>
              <NDescriptionsItem label="X-Panel / Xray">
                {{ detail.xpanelVersion || '未检测到' }} / {{ detail.xrayVersion || '未检测到' }}
              </NDescriptionsItem>
              <NDescriptionsItem label="CPU 使用率">{{ formatPercent(detail.cpuUsage) }}</NDescriptionsItem>
              <NDescriptionsItem label="内存使用">
                {{ formatResource(detail.memoryUsed, detail.memoryTotal) }}
              </NDescriptionsItem>
              <NDescriptionsItem label="磁盘使用">
                {{ formatResource(detail.diskUsed, detail.diskTotal) }}
              </NDescriptionsItem>
              <NDescriptionsItem label="最近心跳">{{ formatDate(detail.lastSeenAt) }}</NDescriptionsItem>
              <NDescriptionsItem label="最近成功同步">
                <NSpace align="center" :size="8">
                  <span>{{ formatDate(detail.lastSyncAt) }}</span>
                  <DataFreshness :data-at="detail.lastSyncAt" compact />
                </NSpace>
              </NDescriptionsItem>
            </NDescriptions>
          </NCard>

          <NCard size="small" title="当前汇总">
            <div class="space-y-14px">
              <div class="flex items-center justify-between">
                <span class="text-gray-500">业务用户</span>
                <span class="text-18px font-600">{{ detail.userCount }}</span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-gray-500">Inbound</span>
                <span class="text-18px font-600">{{ detail.inbounds.length }}</span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-gray-500">出口 IP</span>
                <span class="text-18px font-600">{{ detail.exitIps.length }}</span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-gray-500">累计流量</span>
                <span class="font-600"><TrafficValue :value="detail.trafficBytes" /></span>
              </div>
            </div>
          </NCard>
        </div>

        <NAlert v-if="detail.status === 'offline'" type="error" :show-icon="true" class="mt-16px">
          节点已离线，立即同步请求会在 Agent 恢复连接后的下一次周期同步中执行。
        </NAlert>

        <NCard title="节点成本" size="small" class="mt-16px">
          <template #header-extra>
            <NButton size="small" type="primary" @click="openCreateCost">
              <template #icon><icon-mdi-plus /></template>
              新增成本
            </NButton>
          </template>
          <NDataTable
            v-if="detail.costs.length"
            :data="detail.costs"
            :bordered="false"
            :single-line="false"
            size="small"
            :scroll-x="760"
            :columns="costColumns"
          />
          <NEmpty v-else description="暂无节点成本记录" size="small" />
        </NCard>

        <NCard title="Inbound 列表" size="small" class="mt-16px">
          <NDataTable
            v-if="detail.inbounds.length"
            :data="detail.inbounds"
            :bordered="false"
            :single-line="false"
            size="small"
            :scroll-x="1080"
            :columns="inboundColumns"
          />
          <NEmpty v-else description="暂无 Inbound 快照" size="small" />
        </NCard>

        <NCard :title="nodeExitIpTitle" size="small" class="mt-16px">
          <NDataTable
            v-if="detail.exitIps.length"
            :data="detail.exitIps"
            :bordered="false"
            :single-line="false"
            size="small"
            :columns="exitIpColumns"
          />
          <NEmpty v-else description="暂无出口 IP 配置" size="small" />
        </NCard>

        <div class="grid gap-16px lg:grid-cols-2 mt-16px">
          <NCard title="最近同步记录" size="small">
            <div v-if="detail.syncRuns.slice(0, 5).length" class="space-y-10px">
              <div
                v-for="run in detail.syncRuns.slice(0, 5)"
                :key="run.id"
                class="rounded-6px bg-gray-50 p-10px dark:bg-dark"
              >
                <div class="flex items-center justify-between gap-8px">
                  <span class="text-12px text-gray-500">{{ formatDate(run.startedAt) }}</span>
                  <SyncStatusTag :status="run.status" />
                </div>
                <div class="mt-4px text-12px">
                  Inbound {{ run.inboundCount }} · Client {{ run.clientCount }} · {{ run.syncId }}
                </div>
                <div v-if="run.errorMessage" class="mt-4px text-12px text-red-500">{{ run.errorMessage }}</div>
              </div>
            </div>
            <NEmpty v-else description="暂无同步记录" size="small" />
          </NCard>

          <NCard title="最近状态事件" size="small">
            <div v-if="detail.statusHistory.length" class="space-y-10px">
              <div
                v-for="event in detail.statusHistory"
                :key="event.id"
                class="rounded-6px bg-gray-50 p-10px dark:bg-dark"
              >
                <div class="flex items-center justify-between gap-8px">
                  <span class="font-500">{{ event.title || event.type }}</span>
                  <NTag size="small" :type="eventType(event.severity)">{{ event.severity }}</NTag>
                </div>
                <div class="mt-4px text-12px text-gray-500">{{ formatDate(event.occurredAt) }}</div>
                <div class="mt-4px text-13px">{{ event.message }}</div>
              </div>
            </div>
            <NEmpty v-else description="暂无状态事件" size="small" />
          </NCard>
        </div>
      </template>
    </ModulePage>

    <NModal
      v-model:show="installerModalVisible"
      preset="card"
      title="生成 Agent 一键部署命令"
      class="w-720px max-w-95vw"
      :mask-closable="false"
    >
      <NAlert v-if="installerError" type="error" :show-icon="true" class="mb-16px">{{ installerError }}</NAlert>
      <template v-if="installerResult">
        <NAlert type="success" :show-icon="true" class="mb-16px">
          安装 Token 已生成，15 分钟内有效且只能使用一次。请立即复制命令并在目标 Ubuntu/Debian 服务器执行。
        </NAlert>
        <NForm label-placement="top">
          <NFormItem label="一键部署命令">
            <NSpace vertical :size="8" class="w-full">
              <NInput :value="installerCommand" readonly type="textarea" :autosize="{ minRows: 2, maxRows: 4 }" />
              <NButton
                size="small"
                type="primary"
                class="self-start"
                @click="copyText(installerCommand, '一键部署命令已复制')"
              >
                复制部署命令
              </NButton>
            </NSpace>
          </NFormItem>
        </NForm>
      </template>
      <template v-else>
        <NAlert type="info" :show-icon="true" class="mb-16px">
          已设置管理地址的节点会自动带入 X-Panel 地址和 API 路径；历史节点未设置时可在这里补充。
        </NAlert>
        <NForm label-placement="top">
          <NFormItem label="X-Panel 地址" required>
            <NInput
              v-model:value="installerXpanelUrl"
              :readonly="hasManagementUrl"
              placeholder="例如 http://127.0.0.1:2053"
            />
          </NFormItem>
          <NFormItem label="X-Panel API 路径">
            <NInput v-model:value="installerBasePath" :readonly="hasManagementUrl" placeholder="默认 /；例如 /Alien" />
          </NFormItem>
        </NForm>
      </template>
      <template #footer>
        <NSpace justify="end">
          <NButton :disabled="installerLoading" @click="closeInstaller">关闭</NButton>
          <NButton v-if="!installerResult" type="primary" :loading="installerLoading" @click="generateInstaller">
            生成命令
          </NButton>
        </NSpace>
      </template>
    </NModal>

    <NModal
      v-model:show="managementModalVisible"
      preset="card"
      title="编辑管理地址"
      class="w-620px max-w-95vw"
      :mask-closable="false"
    >
      <NAlert v-if="managementError" type="error" :show-icon="true" class="mb-16px">{{ managementError }}</NAlert>
      <NForm label-placement="top">
        <NFormItem label="管理地址" required>
          <NInput
            v-model:value="managementUrl"
            maxlength="2048"
            placeholder="例如 https://server.example.com:18086/Alien/"
          />
        </NFormItem>
        <NFormItem label="解析结果">
          <div class="text-13px text-gray-500">
            X-Panel 地址：{{ parseManagementAddress(managementUrl)?.url || '--' }}；API 路径：{{
              parseManagementAddress(managementUrl)?.basePath || '--'
            }}
          </div>
        </NFormItem>
      </NForm>
      <template #footer>
        <NSpace justify="end">
          <NButton :disabled="managementSaving" @click="closeManagementEditor">取消</NButton>
          <NButton type="primary" :loading="managementSaving" @click="saveManagementEditor">保存</NButton>
        </NSpace>
      </template>
    </NModal>

    <NModal
      v-model:show="costModalVisible"
      preset="card"
      :title="costModalTitle"
      class="w-620px max-w-95vw"
      :mask-closable="false"
    >
      <NAlert v-if="costError" type="error" :show-icon="true" class="mb-16px">{{ costError }}</NAlert>
      <NAlert v-if="editingCostId" type="info" :show-icon="true" class="mb-16px">
        生效日期用于保留历史口径，编辑已有记录时不可改日期；如需更换日期，请新增一条成本记录。
      </NAlert>
      <NForm label-placement="top">
        <NFormItem label="成本类别" required>
          <NInput v-model:value="costForm.category" maxlength="120" placeholder="例如：服务器、带宽、托管" />
        </NFormItem>
        <NFormItem label="月金额（CNY）" required>
          <NInputNumber
            v-model:value="costForm.monthlyAmount"
            :min="0"
            :max="100000000"
            :precision="2"
            class="w-full"
          />
        </NFormItem>
        <div class="grid gap-x-12px md:grid-cols-2">
          <NFormItem label="生效日期" required>
            <NDatePicker
              v-model:formatted-value="costForm.effectiveFrom"
              type="date"
              value-format="yyyy-MM-dd"
              :disabled="Boolean(editingCostId)"
              :clearable="false"
              class="w-full"
            />
          </NFormItem>
          <NFormItem label="失效日期">
            <NDatePicker
              v-model:formatted-value="costForm.effectiveTo"
              type="date"
              value-format="yyyy-MM-dd"
              :disabled="Boolean(editingCostId)"
              clearable
              class="w-full"
            />
          </NFormItem>
        </div>
        <NFormItem label="备注">
          <NInput
            v-model:value="costForm.notes"
            type="textarea"
            maxlength="2000"
            show-count
            :autosize="{ minRows: 3, maxRows: 6 }"
            placeholder="可选"
          />
        </NFormItem>
      </NForm>
      <template #footer>
        <NSpace justify="end">
          <NButton :disabled="costSaving" @click="closeCostModal">取消</NButton>
          <NButton type="primary" :loading="costSaving" @click="saveCost">保存</NButton>
        </NSpace>
      </template>
    </NModal>
  </div>
</template>
