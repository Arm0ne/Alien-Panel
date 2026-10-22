<script setup lang="ts">
import { computed, h, onMounted, reactive, ref, watch } from 'vue';
import { NDropdown, NInput } from 'naive-ui';
import type { DropdownOption } from 'naive-ui';
import { useRoute, useRouter } from 'vue-router';
import { deleteNode, fetchNodeDetail, issueNodeInstallToken, requestNodeSync, updateNode } from '@/service/api';
import DataFreshness from '@/components/project/data-freshness.vue';
import NodeStatusTag from '@/components/project/node-status-tag.vue';
import NodeTypeTag from '@/components/project/node-type-tag.vue';
import SyncStatusTag from '@/components/project/sync-status-tag.vue';
import TrafficValue from '@/components/project/traffic-value.vue';
import SvgIcon from '@/components/custom/svg-icon.vue';

defineOptions({ name: 'NodeDetail' });

const route = useRoute();
const router = useRouter();
const nodeId = computed(() => String(route.params.id || ''));

const loading = ref(false);
const errorMessage = ref('');
const detail = ref<Api.Central.NodeDetail | null>(null);

const editVisible = ref(false);
const editSaving = ref(false);
const editError = ref('');
const editForm = reactive({ name: '', hostname: '', region: '', provider: '', managementUrl: '' });

const syncVisible = ref(false);
const syncing = ref(false);
const syncError = ref('');

const statusVisible = ref(false);
const statusSaving = ref(false);
const statusError = ref('');

const deleteVisible = ref(false);
const deleteSaving = ref(false);
const deleteError = ref('');
const deleteConfirmName = ref('');

const installerVisible = ref(false);
const installerLoading = ref(false);
const installerError = ref('');
const installerResult = ref<Api.Central.NodeInstallTokenResult | null>(null);
const installerXpanelUrl = ref('');
const installerBasePath = ref('/');
const installerCopied = ref(false);

const syncPage = ref(1);
const syncPageSize = 5;
const expandedErrors = ref<string[]>([]);

function formatDate(value?: string | null) {
  if (!value) return '--';
  const date = new Date(value);
  return Number.isNaN(date.getTime()) ? '时间未知' : date.toLocaleString('zh-CN', { hour12: false });
}

function formatDateOnly(value?: string | null) {
  if (!value) return '--';
  const date = new Date(value);
  return Number.isNaN(date.getTime()) ? value : date.toLocaleDateString('zh-CN');
}

function relativeTime(value?: string | null) {
  if (!value) return '时间未知';
  const timestamp = new Date(value).getTime();
  if (!Number.isFinite(timestamp)) return '时间未知';
  const seconds = Math.max(0, Math.round((Date.now() - timestamp) / 1000));
  if (seconds < 60) return `${seconds} 秒前`;
  const minutes = Math.round(seconds / 60);
  if (minutes < 60) return `${minutes} 分钟前`;
  const hours = Math.round(minutes / 60);
  if (hours < 24) return `${hours} 小时前`;
  return `${Math.round(hours / 24)} 天前`;
}

function formatMoney(value?: number, currency = 'CNY') {
  if (value === undefined || value === null || !Number.isFinite(value)) return '--';
  return new Intl.NumberFormat('zh-CN', { style: 'currency', currency }).format(value);
}

function formatDuration(run: Api.Central.NodeSyncRun) {
  if (!run.finishedAt) return run.status === 'running' ? '运行中' : '--';
  const started = new Date(run.startedAt).getTime();
  const finished = new Date(run.finishedAt).getTime();
  if (!Number.isFinite(started) || !Number.isFinite(finished) || finished < started) return '--';
  return `${((finished - started) / 1000).toFixed(1)} s`;
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
    return (
      /^https?:$/i.test(parsed.protocol) &&
      Boolean(parsed.hostname) &&
      !parsed.username &&
      !parsed.password &&
      !parsed.search &&
      !parsed.hash
    );
  } catch {
    return false;
  }
}

function centralApiUrl() {
  const configuredCentralUrl = String(import.meta.env.VITE_SERVICE_BASE_URL || '/api').trim();
  return /^https?:\/\//i.test(configuredCentralUrl)
    ? configuredCentralUrl.replace(/\/$/, '')
    : new URL(configuredCentralUrl || '/api', window.location.origin).toString().replace(/\/$/, '');
}

function shellQuote(value: string) {
  return `'${value.replace(/'/g, `'\\''`)}'`;
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

const hasManagementUrl = computed(() => Boolean(detail.value?.managementUrl));
const exitIpTitle = computed(() => (detail.value?.type === 'relay' ? '线路机直出 IP' : '落地机出口 IP'));

const exitFamilyStats = computed(() => {
  const list = detail.value?.exitIps || [];
  const v4 = list.filter(item => item.family === 4).length;
  const v6 = list.filter(item => item.family === 6).length;
  const enabled = list.filter(item => item.status === 'active').length;
  return { v4, v6, enabled, total: list.length };
});

const monthlyCostValue = computed(() => {
  const node = detail.value;
  if (!node) return undefined;
  if (typeof node.monthlyCost === 'number' && Number.isFinite(node.monthlyCost)) return node.monthlyCost;
  const today = new Date().toISOString().slice(0, 10);
  return (node.costs || [])
    .filter(cost => cost.effectiveFrom <= today && (!cost.effectiveTo || cost.effectiveTo >= today))
    .reduce((sum, cost) => sum + cost.monthlyAmount, 0);
});

const monthlyCostLabel = computed(() => formatMoney(monthlyCostValue.value, detail.value?.currency));

const syncPageCount = computed(() => Math.max(1, Math.ceil((detail.value?.syncRuns.length || 0) / syncPageSize)));

const pagedSyncRuns = computed(() => {
  const list = detail.value?.syncRuns || [];
  const start = (syncPage.value - 1) * syncPageSize;
  return list.slice(start, start + syncPageSize);
});

watch(
  () => detail.value?.syncRuns.length,
  () => {
    if (syncPage.value > syncPageCount.value) syncPage.value = syncPageCount.value;
  }
);

function inboundPurposeLabel(row: Api.Central.NodeInboundDetail) {
  const purpose = row.purpose || (row.kind === 'user' ? 'business' : 'infrastructure');
  return purpose === 'business' ? '业务用户' : '基础设施';
}

function inboundPurposeTone(row: Api.Central.NodeInboundDetail) {
  const purpose = row.purpose || (row.kind === 'user' ? 'business' : 'infrastructure');
  return purpose === 'business' ? 'is-brand' : 'is-neutral';
}

function inboundStatusLabel(status: Api.Central.NodeInboundStatus) {
  if (status === 'active') return '启用';
  if (status === 'disabled') return '停用';
  return '已归档';
}

function inboundStatusTone(status: Api.Central.NodeInboundStatus) {
  if (status === 'active') return 'is-ok';
  if (status === 'disabled') return 'is-neutral';
  return 'is-warn';
}

function severityLabel(severity: Api.Central.EventSeverity) {
  if (severity === 'error') return '错误';
  if (severity === 'warning') return '警告';
  return '提示';
}

function severityTone(severity: Api.Central.EventSeverity) {
  if (severity === 'error') return 'is-danger';
  if (severity === 'warning') return 'is-warn';
  return 'is-info';
}

function isErrorExpanded(id: string) {
  return expandedErrors.value.includes(id);
}

function toggleError(id: string) {
  if (isErrorExpanded(id)) {
    expandedErrors.value = expandedErrors.value.filter(item => item !== id);
  } else {
    expandedErrors.value = [...expandedErrors.value, id];
  }
}

const menuOptions = computed<DropdownOption[]>(() => [
  { label: '编辑节点信息', key: 'edit', icon: () => h(SvgIcon, { icon: 'mdi:pencil-outline' }) },
  { label: '复制节点 Key', key: 'copy-key', icon: () => h(SvgIcon, { icon: 'mdi:content-copy' }) },
  { label: '重新生成部署命令', key: 'installer', icon: () => h(SvgIcon, { icon: 'mdi:console-line' }) },
  { type: 'divider', key: 'd1' },
  {
    label: detail.value?.enabled ? '停用节点' : '启用节点',
    key: 'toggle',
    icon: () => h(SvgIcon, { icon: detail.value?.enabled ? 'mdi:power' : 'mdi:check-circle-outline' })
  },
  { label: '删除节点 · 不可恢复', key: 'delete', icon: () => h(SvgIcon, { icon: 'mdi:delete-outline' }) }
]);

function handleMenu(key: string) {
  if (key === 'edit') openEdit();
  else if (key === 'copy-key') copyNodeKey();
  else if (key === 'installer') openInstaller();
  else if (key === 'toggle') openStatusModal();
  else if (key === 'delete') openDelete();
}

async function copyText(value: string, successMessage: string) {
  try {
    await navigator.clipboard.writeText(value);
    window.$message?.success(successMessage);
  } catch {
    window.$message?.error('复制失败，请手动选择并复制');
  }
}

function copyNodeKey() {
  if (detail.value?.nodeKey) void copyText(detail.value.nodeKey, '节点 Key 已复制');
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
    if (syncPage.value > syncPageCount.value) syncPage.value = syncPageCount.value;
  }
  loading.value = false;
}

function openEdit() {
  if (!detail.value) return;
  editError.value = '';
  editForm.name = detail.value.name;
  editForm.hostname = detail.value.hostname || '';
  editForm.region = detail.value.region || '';
  editForm.provider = detail.value.provider || '';
  editForm.managementUrl = detail.value.managementUrl || '';
  editVisible.value = true;
}

function closeEdit() {
  if (!editSaving.value) editVisible.value = false;
}

function openSyncModal() {
  if (!detail.value || !detail.value.enabled) return;
  syncError.value = '';
  syncVisible.value = true;
}

function closeSyncModal() {
  if (!syncing.value) syncVisible.value = false;
}

function openStatusModal() {
  if (!detail.value) return;
  statusError.value = '';
  statusVisible.value = true;
}

function closeStatusModal() {
  if (!statusSaving.value) statusVisible.value = false;
}

function openDelete() {
  if (!detail.value) return;
  deleteError.value = '';
  deleteConfirmName.value = '';
  deleteVisible.value = true;
}

function closeDelete() {
  if (!deleteSaving.value) deleteVisible.value = false;
}

function openInstaller() {
  if (!detail.value) return;
  installerError.value = '';
  installerResult.value = null;
  installerCopied.value = false;
  const parsed = parseManagementAddress(detail.value.managementUrl || '');
  installerXpanelUrl.value = parsed?.url || '';
  installerBasePath.value = parsed?.basePath || detail.value.panelBasePath || '/';
  installerVisible.value = true;
}

function closeInstaller() {
  if (!installerLoading.value) installerVisible.value = false;
}

async function saveEdit() {
  if (!detail.value || editSaving.value) return;
  const name = editForm.name.trim();
  if (!name) {
    editError.value = '请填写节点名称';
    return;
  }
  if (name.length > 120) {
    editError.value = '节点名称不能超过 120 个字符';
    return;
  }
  const managementUrl = editForm.managementUrl.trim();
  const parsed = parseManagementAddress(managementUrl);
  if (!parsed) {
    editError.value = '请填写有效的管理地址，例如 https://server.example.com:18086/Alien/';
    return;
  }
  editError.value = '';
  editSaving.value = true;
  const { error } = await updateNode(detail.value.id, {
    name,
    hostname: editForm.hostname.trim() || undefined,
    region: editForm.region.trim() || undefined,
    provider: editForm.provider.trim() || undefined,
    managementUrl,
    panelBasePath: parsed.basePath
  });
  editSaving.value = false;
  if (error) {
    editError.value = '节点信息保存失败，请检查填写内容后重试。';
    return;
  }
  editVisible.value = false;
  window.$message?.success('节点信息已更新');
  await loadNode();
}

async function confirmSync() {
  if (!detail.value || syncing.value) return;
  syncing.value = true;
  syncError.value = '';
  const { data, error } = await requestNodeSync(detail.value.id);
  syncing.value = false;
  if (error || !data) {
    syncError.value = '同步请求发送失败，请检查节点状态后重试';
    return;
  }
  syncVisible.value = false;
  window.$message?.success('已提交立即同步请求，Agent 将在下一次周期同步中执行');
  await loadNode();
}

async function confirmStatusChange() {
  if (!detail.value || statusSaving.value) return;
  const enabled = !detail.value.enabled;
  statusSaving.value = true;
  statusError.value = '';
  const { error } = await updateNode(detail.value.id, { enabled });
  statusSaving.value = false;
  if (error) {
    statusError.value = enabled ? '启用节点失败，请稍后重试' : '停用节点失败，请稍后重试';
    return;
  }
  statusVisible.value = false;
  window.$message?.success(enabled ? '节点已启用' : '节点已停用');
  await loadNode();
}

async function confirmDelete() {
  if (!detail.value || deleteSaving.value) return;
  if (deleteConfirmName.value.trim() !== detail.value.name) return;
  deleteSaving.value = true;
  deleteError.value = '';
  const { error } = await deleteNode(detail.value.id);
  deleteSaving.value = false;
  if (error) {
    deleteError.value = error.response?.data?.msg || '节点删除失败，请稍后重试';
    return;
  }
  deleteVisible.value = false;
  window.$message?.success('节点及其关联配置已删除');
  router.push({ name: 'nodes' });
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
  installerCopied.value = false;
}

async function copyInstallerCommand() {
  await copyText(installerCommand.value, '一键部署命令已复制');
  installerCopied.value = true;
}

function goBack() {
  router.push({ name: 'nodes' });
}

function goUsers() {
  if (!detail.value) return;
  router.push({ name: 'users', query: { nodeId: detail.value.id } });
}

function goFinance() {
  router.push({ name: 'finance' });
}

watch(nodeId, () => void loadNode());
onMounted(() => void loadNode());
</script>

<template>
  <div class="d-page">
    <!-- 页头 -->
    <header class="d-head">
      <div class="d-head-main">
        <div class="d-crumb">
          <button type="button" @click="goBack">业务</button>
          <span class="d-crumb-sep">/</span>
          <button type="button" class="d-crumb-back" @click="goBack">
            <SvgIcon icon="mdi:chevron-left" class="d-crumb-ico" />节点管理
          </button>
          <span class="d-crumb-sep">/</span>
          <span>节点详情</span>
        </div>

        <div class="d-title-row">
          <h1 class="d-title">{{ detail?.name || '节点详情' }}</h1>
          <template v-if="detail">
            <NodeTypeTag :type="detail.type" />
            <NodeStatusTag :status="detail.status" />
            <SyncStatusTag :status="detail.syncStatus" />
          </template>
        </div>

        <div v-if="detail" class="d-meta-row">
          <a
            v-if="detail.managementUrl"
            class="d-mana-link"
            :href="detail.managementUrl"
            target="_blank"
            rel="noopener noreferrer"
          >
            <SvgIcon icon="mdi:open-in-new" class="d-mana-ico" />
            <span class="d-mono">{{ detail.managementUrl }}</span>
          </a>
          <span v-else class="d-dim">未设置管理地址</span>
          <span class="d-dim">·</span>
          <span class="d-chip">{{ detail.region || '区域未知' }}</span>
          <span class="d-chip">{{ detail.provider || '服务商未知' }}</span>
          <span class="d-chip">
            <SvgIcon icon="mdi:account-outline" class="d-chip-ico" />
            {{ detail.type === 'relay' ? `承载 ${detail.userCount} 个业务用户` : '落地机不承载业务用户' }}
          </span>
        </div>

        <p class="d-desc">
          查看节点元数据、Inbound、节点公网出口 IP、同步运行和最近状态事件；本页只维护中央节点记录，
          <b class="d-strong">不写入 Xray</b>，所有 Inbound / Client / 流量数值来自最近一次 Agent 同步。
        </p>
      </div>

      <div v-if="detail" class="d-head-actions">
        <DataFreshness :data-at="detail.dataAt" compact />
        <button class="d-btn" type="button" @click="openEdit">
          <SvgIcon icon="mdi:pencil-outline" class="d-btn-ico" />编辑
        </button>
        <button class="d-btn is-primary" type="button" :disabled="!detail.enabled" @click="openSyncModal">
          <SvgIcon icon="mdi:refresh" class="d-btn-ico" />立即同步
        </button>
        <NDropdown trigger="click" :options="menuOptions" @select="key => handleMenu(String(key))">
          <button class="d-icon-btn" type="button" aria-label="更多操作">
            <SvgIcon icon="mdi:dots-horizontal" />
          </button>
        </NDropdown>
      </div>
    </header>

    <div v-if="errorMessage" class="d-alert is-danger">{{ errorMessage }}</div>

    <div v-if="loading && !detail" class="d-loading">正在读取节点详情…</div>
    <div v-else-if="!detail" class="d-empty">
      <div class="d-empty-title">暂无节点详情</div>
      <p class="d-empty-desc">请确认节点仍存在，或返回节点列表重新选择。</p>
      <button class="d-btn" type="button" @click="goBack">
        <SvgIcon icon="mdi:chevron-left" class="d-btn-ico" />返回节点列表
      </button>
    </div>

    <template v-else>
      <!-- KPI -->
      <section class="d-kpis">
        <div class="d-kpi">
          <div class="d-kpi-top">
            <span class="d-kpi-label">当前用户 <span class="d-dim">· 业务用户</span></span>
            <span class="d-kpi-icon is-blue"><SvgIcon icon="mdi:account-group-outline" /></span>
          </div>
          <div class="d-kpi-value">
            {{ detail.userCount }}<span class="d-kpi-unit">户</span>
          </div>
          <div class="d-kpi-foot">
            <template v-if="detail.type === 'relay'">收费 {{ detail.paidUserCount }} · 免费 {{ detail.freeUserCount }}</template>
            <template v-else>落地机不承载业务用户</template>
            <span class="d-dim">·</span>
            <span>仅统计线路机 Inbound</span>
          </div>
        </div>

        <div class="d-kpi">
          <div class="d-kpi-top">
            <span class="d-kpi-label">启用 Client</span>
            <span class="d-kpi-icon is-green"><SvgIcon icon="mdi:account-outline" /></span>
          </div>
          <div class="d-kpi-value">{{ detail.clientCount }}<span class="d-kpi-unit">个</span></div>
          <div class="d-kpi-foot">
            <span>跨 {{ detail.inbounds.length }} 个 Inbound</span>
            <span class="d-dim">·</span>
            <span>仅统计启用凭证</span>
          </div>
        </div>

        <div class="d-kpi">
          <div class="d-kpi-top">
            <span class="d-kpi-label">出口 IP</span>
            <span class="d-kpi-icon is-cyan"><SvgIcon icon="mdi:globe" /></span>
          </div>
          <div class="d-kpi-value">{{ exitFamilyStats.total }}<span class="d-kpi-unit">个</span></div>
          <div class="d-kpi-foot">
            <span>IPv4 {{ exitFamilyStats.v4 }} · IPv6 {{ exitFamilyStats.v6 }}</span>
            <span class="d-dim">·</span>
            <span>启用 {{ exitFamilyStats.enabled }}</span>
          </div>
        </div>

        <div class="d-kpi">
          <div class="d-kpi-top">
            <span class="d-kpi-label">月成本</span>
            <span class="d-kpi-icon is-amber"><SvgIcon icon="mdi:wallet-outline" /></span>
          </div>
          <div class="d-kpi-value">{{ monthlyCostLabel }}</div>
          <div class="d-kpi-foot">
            <span>{{ detail.costs.length }} 条成本记录</span>
            <span class="d-dim">·</span>
            <button class="d-link-btn" type="button" @click="goFinance">查看成本明细</button>
          </div>
        </div>
      </section>

      <!-- 节点元数据 -->
      <section class="d-card">
        <div class="d-card-head">
          <h3>节点元数据</h3>
          <span class="d-badge is-neutral">只读</span>
          <div class="d-right">
            <span class="d-dim">最近成功同步 {{ formatDate(detail.lastSyncAt) }}</span>
            <span class="d-chip">
              <SvgIcon icon="mdi:clock-outline" class="d-chip-ico" />{{ relativeTime(detail.lastSyncAt) }}
            </span>
          </div>
        </div>
        <div class="d-card-body">
          <div class="d-meta-grid">
            <div class="d-field d-meta-wide">
              <span class="d-label">节点 Key</span>
              <div class="d-meta-key">
                <span class="d-mono">{{ detail.nodeKey }}</span>
                <button class="d-icon-btn" type="button" title="复制节点 Key" @click="copyNodeKey">
                  <SvgIcon icon="mdi:content-copy" />
                </button>
              </div>
            </div>
            <div class="d-field">
              <span class="d-label">主机名</span>
              <b class="d-mono">{{ detail.hostname || '--' }}</b>
            </div>
            <div class="d-field">
              <span class="d-label">区域</span>
              <b>{{ detail.region || '--' }}</b>
            </div>
            <div class="d-field">
              <span class="d-label">提供商</span>
              <b>{{ detail.provider || '--' }}</b>
            </div>
            <div class="d-field">
              <span class="d-label">Agent 版本</span>
              <b class="d-mono">{{ detail.agentVersion || '未上报' }}</b>
            </div>
            <div class="d-field">
              <span class="d-label">Xray 版本</span>
              <b class="d-mono">{{ detail.xrayVersion || '未检测到' }}</b>
            </div>
            <div class="d-field">
              <span class="d-label">X-Panel 同步</span>
              <div class="d-inline">
                <SyncStatusTag :status="detail.syncStatus" />
                <span v-if="detail.lastSyncError" class="d-field-error">{{ detail.lastSyncError }}</span>
              </div>
            </div>
            <div class="d-field">
              <span class="d-label">启用状态</span>
              <div>
                <span class="d-badge" :class="detail.enabled ? 'is-ok' : 'is-neutral'">
                  {{ detail.enabled ? '启用' : '已停用' }}
                </span>
              </div>
            </div>
            <div class="d-field">
              <span class="d-label">最近成功同步</span>
              <div class="d-inline">
                <b class="d-num">{{ formatDate(detail.lastSyncAt) }}</b>
                <DataFreshness :data-at="detail.lastSyncAt" compact />
              </div>
            </div>
            <div class="d-field">
              <span class="d-label">月成本</span>
              <div class="d-inline">
                <b class="d-num">{{ monthlyCostLabel }}</b>
                <button class="d-link-btn" type="button" @click="goFinance">在成本与收入中管理 →</button>
              </div>
            </div>
            <div class="d-field d-meta-mgmt">
              <span class="d-label">管理地址</span>
              <a
                v-if="detail.managementUrl"
                class="d-mono d-mana-link"
                :href="detail.managementUrl"
                target="_blank"
                rel="noopener noreferrer"
              >
                {{ detail.managementUrl }}
              </a>
              <span v-else class="d-dim">未设置管理地址</span>
            </div>
          </div>
        </div>
      </section>

      <!-- 承载用户与 Inbound -->
      <section class="d-card">
        <div class="d-card-head">
          <h3>承载用户与 Inbound</h3>
          <span class="d-badge is-neutral">{{ detail.inbounds.length }} 个 Inbound · {{ detail.clientCount }} 个启用 Client</span>
          <div class="d-right">
            <button class="d-btn d-btn-xs is-ghost" type="button" @click="goUsers">
              查看该节点 {{ detail.userCount }} 个业务用户
              <SvgIcon icon="mdi:chevron-right" class="d-btn-ico" />
            </button>
          </div>
        </div>
        <div class="d-table-wrap">
          <table class="d-table">
            <thead>
              <tr>
                <th>Tag / Remark</th>
                <th>协议 / 端口</th>
                <th>类型</th>
                <th class="is-num">设备数</th>
                <th class="is-num">累计流量</th>
                <th>到期时间</th>
                <th>状态</th>
                <th>最近同步</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="inbound in detail.inbounds" :key="inbound.id || inbound.remoteId">
                <td>
                  <b class="d-strong">{{ inbound.tag || inbound.remoteId }}</b>
                  <span class="d-sub">{{ inbound.remark || '--' }}</span>
                </td>
                <td class="d-mono">{{ inbound.protocol || '--' }} / {{ inbound.port || '--' }}</td>
                <td><span class="d-badge" :class="inboundPurposeTone(inbound)">{{ inboundPurposeLabel(inbound) }}</span></td>
                <td class="is-num">{{ inbound.clientCount }}</td>
                <td class="is-num"><TrafficValue :value="inbound.allTime" /></td>
                <td class="d-num">{{ formatDate(inbound.expiresAt) }}</td>
                <td>
                  <span class="d-badge" :class="inboundStatusTone(inbound.status)">{{ inboundStatusLabel(inbound.status) }}</span>
                </td>
                <td class="d-num d-dim">{{ formatDate(inbound.lastSeenAt) }}</td>
              </tr>
              <tr v-if="!detail.inbounds.length">
                <td colspan="8" class="d-table-empty">暂无 Inbound 快照</td>
              </tr>
            </tbody>
          </table>
        </div>
        <div class="d-card-body d-card-note">
          只有线路机的业务用户 Inbound 会产生业务用户；落地机 Inbound 恒为基础设施，不承载业务用户与业务流量。
          「累计流量」为 X-Panel 上报的累计总量，不是相对流量基线的增量。
        </div>
      </section>

      <!-- 出口 IP -->
      <section class="d-card">
        <div class="d-card-head">
          <h3>{{ exitIpTitle }}</h3>
          <span class="d-badge is-neutral">{{ detail.exitIps.length }} 个出口</span>
          <div class="d-right">
            <span class="d-dim">出口 IP 月成本不计入节点月成本</span>
          </div>
        </div>
        <div class="d-table-wrap">
          <table class="d-table">
            <thead>
              <tr>
                <th>地址</th>
                <th>区域</th>
                <th>提供商</th>
                <th>协议族</th>
                <th>状态</th>
                <th class="is-num">月成本</th>
                <th>有效期</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="ip in detail.exitIps" :key="ip.id">
                <td class="d-strong d-mono">{{ ip.address }}</td>
                <td>{{ ip.region || '--' }}</td>
                <td>{{ ip.provider || '--' }}</td>
                <td>
                  <span class="d-badge" :class="ip.family === 6 ? 'is-violet' : 'is-info'">IPv{{ ip.family }}</span>
                </td>
                <td>
                  <span class="d-badge" :class="ip.status === 'active' ? 'is-ok' : 'is-neutral'">
                    {{ ip.status === 'active' ? '启用' : '停用' }}
                  </span>
                </td>
                <td class="is-num">{{ formatMoney(ip.monthlyCost, ip.currency) }}</td>
                <td class="d-num">{{ formatDateOnly(ip.validFrom) }} → {{ formatDateOnly(ip.validTo) }}</td>
              </tr>
              <tr v-if="!detail.exitIps.length">
                <td colspan="7" class="d-table-empty">暂无出口 IP 配置</td>
              </tr>
            </tbody>
          </table>
        </div>
        <div class="d-card-body d-card-note">
          <div class="d-callout is-info">
            <SvgIcon icon="mdi:source-branch" class="d-callout-ico" />
            <div>
              <b>路径关联与删除保护</b>
              出口 IP 被用户路径引用时不可删除：校验会同时检查路径上的主出口（<span class="d-mono">exit_ip_id</span>）与有序多出口关联表（<span class="d-mono">user_path_exit_ips</span>），放在第二跳也不能绕过。
              出口 IP 月成本属于资产成本，不计入上面的节点月成本。
            </div>
          </div>
        </div>
      </section>

      <!-- 同步记录 + 事件 -->
      <div class="d-two-col">
        <section class="d-card">
          <div class="d-card-head">
            <h3>最近同步记录</h3>
            <span class="d-badge is-neutral">共 {{ detail.syncRuns.length }} 条 · 按开始时间倒序</span>
            <div class="d-right">
              <button class="d-btn d-btn-xs is-ghost" type="button" :disabled="syncing" @click="loadNode">
                <SvgIcon icon="mdi:refresh" class="d-btn-ico" />刷新
              </button>
              <button class="d-btn d-btn-xs is-ghost" type="button" :disabled="!detail.enabled" @click="openSyncModal">
                立即同步
              </button>
            </div>
          </div>
          <div class="d-table-wrap">
            <table class="d-table">
              <thead>
                <tr>
                  <th>开始时间</th>
                  <th>结果</th>
                  <th>数据量</th>
                  <th class="is-num">耗时</th>
                  <th>同步 ID</th>
                  <th></th>
                </tr>
              </thead>
              <tbody>
                <template v-for="run in pagedSyncRuns" :key="run.id">
                  <tr>
                    <td class="d-num">{{ formatDate(run.startedAt) }}</td>
                    <td><SyncStatusTag :status="run.status" /></td>
                    <td class="d-dim">
                      {{ run.status === 'success' ? `Inbound ${run.inboundCount} · Client ${run.clientCount}` : '--' }}
                    </td>
                    <td class="is-num">{{ formatDuration(run) }}</td>
                    <td class="d-mono d-dim">{{ run.syncId }}</td>
                    <td class="is-num">
                      <button v-if="run.errorMessage" class="d-link-btn" type="button" @click="toggleError(run.id)">
                        {{ isErrorExpanded(run.id) ? '收起' : '错误详情' }}
                      </button>
                    </td>
                  </tr>
                  <tr v-if="run.errorMessage && isErrorExpanded(run.id)">
                    <td colspan="6" class="d-error-cell">
                      <pre class="d-error-text">{{ run.errorMessage }}</pre>
                    </td>
                  </tr>
                </template>
                <tr v-if="!detail.syncRuns.length">
                  <td colspan="6" class="d-table-empty">暂无同步记录</td>
                </tr>
              </tbody>
            </table>
          </div>
          <div class="d-table-foot">
            <span class="d-dim">第 {{ syncPage }} / {{ syncPageCount }} 页 · 每页 {{ syncPageSize }} 条</span>
            <div class="d-foot-actions">
              <button class="d-btn d-btn-xs" type="button" :disabled="syncPage <= 1" @click="syncPage -= 1">上一页</button>
              <button class="d-btn d-btn-xs" type="button" :disabled="syncPage >= syncPageCount" @click="syncPage += 1">下一页</button>
            </div>
          </div>
          <div class="d-card-body d-card-note">
            Agent 每 60 秒同步一次；失败不重试，等下一个周期。同步失败只标记采集异常，节点仍可在线。
          </div>
        </section>

        <section class="d-card">
          <div class="d-card-head">
            <h3>最近状态事件</h3>
            <span class="d-badge is-neutral">最近 30 条</span>
            <div class="d-right">
              <button class="d-btn d-btn-xs is-ghost" type="button" @click="router.push({ name: 'events' })">
                事件中心
                <SvgIcon icon="mdi:chevron-right" class="d-btn-ico" />
              </button>
            </div>
          </div>
          <div class="d-event-list">
            <div v-for="event in detail.statusHistory" :key="event.id" class="d-event-row">
              <div class="d-event-main">
                <b>{{ event.title || event.type }}</b>
                <span>{{ formatDate(event.occurredAt) }} · {{ event.message }}</span>
              </div>
              <span class="d-badge" :class="severityTone(event.severity)">{{ severityLabel(event.severity) }}</span>
            </div>
            <div v-if="!detail.statusHistory.length" class="d-table-empty">暂无状态事件</div>
          </div>
        </section>
      </div>

      <div class="d-footnote">
        同步记录保留最近 20 条、状态事件保留最近 30 条；「累计流量」为 X-Panel 上报的累计总量；
        节点的 CPU / 内存 / 磁盘等资源指标由独立探针监控，不在本页展示。本页只维护中央节点记录，不会写入 Xray 或 X-Panel。
      </div>
    </template>

    <!-- ===================== 编辑节点信息 ===================== -->
    <Teleport to="body">
      <div v-if="editVisible" class="d-mask" @click.self="closeEdit">
        <div class="d-modal" style="width: 620px">
          <div class="d-modal-head">
            <div class="d-modal-head-main">
              <h3>编辑节点信息</h3>
              <p class="d-modal-sub">维护中央节点记录，不影响 X-Panel 本机配置</p>
            </div>
            <button class="d-icon-btn" type="button" :disabled="editSaving" @click="closeEdit">
              <SvgIcon icon="mdi:close" />
            </button>
          </div>
          <div class="d-modal-body">
            <div v-if="editError" class="d-alert is-danger mb-12">{{ editError }}</div>
            <div class="d-callout is-info mb-12">
              <SvgIcon icon="mdi:link-variant" class="d-callout-ico" />
              <div>已设置管理地址的节点会自动带入 X-Panel 地址和 API 路径；历史节点未设置时可在这里补充。</div>
            </div>
            <div class="d-form-grid">
              <div class="d-field">
                <span class="d-label">节点名称 <span class="d-req">*</span></span>
                <NInput v-model:value="editForm.name" maxlength="120" placeholder="例如：hk-01" />
                <span class="d-help">最多 120 个字符</span>
              </div>
              <div class="d-field">
                <span class="d-label">主机名</span>
                <NInput v-model:value="editForm.hostname" maxlength="255" placeholder="hk01.example.com" />
              </div>
              <div class="d-field">
                <span class="d-label">区域</span>
                <NInput v-model:value="editForm.region" maxlength="120" placeholder="例如：香港" />
              </div>
              <div class="d-field">
                <span class="d-label">提供商</span>
                <NInput v-model:value="editForm.provider" maxlength="200" placeholder="例如：Example Cloud" />
              </div>
              <div class="d-field d-field-wide">
                <span class="d-label">管理地址 <span class="d-req">*</span></span>
                <NInput v-model:value="editForm.managementUrl" maxlength="2048" placeholder="例如 https://server.example.com:18086/Alien/" />
                <span class="d-help">保存后会同时写入 managementUrl 与 panelBasePath</span>
              </div>
              <div class="d-field d-field-wide">
                <span class="d-label">解析结果</span>
                <span class="d-mono d-dim">X-Panel 地址：{{ parseManagementAddress(editForm.managementUrl)?.url || '--' }}；API 路径：{{ parseManagementAddress(editForm.managementUrl)?.basePath || '--' }}</span>
              </div>
            </div>
          </div>
          <div class="d-modal-foot">
            <span class="d-dim">保存只更新中央记录，不会写入 X-Panel</span>
            <div class="d-foot-actions">
              <button class="d-btn" type="button" :disabled="editSaving" @click="closeEdit">取消</button>
              <button class="d-btn is-primary" type="button" :disabled="editSaving" @click="saveEdit">
                {{ editSaving ? '保存中…' : '保存' }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- ===================== 立即同步确认 ===================== -->
    <Teleport to="body">
      <div v-if="syncVisible" class="d-mask" @click.self="closeSyncModal">
        <div class="d-modal" style="width: 520px">
          <div class="d-modal-head">
            <div class="d-modal-head-main">
              <h3>立即同步</h3>
              <p class="d-modal-sub">请求会排入 Agent 的下一次周期同步</p>
            </div>
            <button class="d-icon-btn" type="button" :disabled="syncing" @click="closeSyncModal">
              <SvgIcon icon="mdi:close" />
            </button>
          </div>
          <div class="d-modal-body">
            <div v-if="syncError" class="d-alert is-danger mb-12">{{ syncError }}</div>
            <div class="d-callout is-info mb-12">
              <SvgIcon icon="mdi:refresh" class="d-callout-ico" />
              <div>立即同步只把请求排入 Agent 的下一次周期同步，<b>不会直接写 Xray</b>；Agent 离线时会在其恢复后的第一个周期执行。</div>
            </div>
            <div class="d-detail-grid">
              <div class="d-field">
                <span class="d-label">节点</span>
                <b>{{ detail?.name }}<span class="d-dim">（{{ detail?.type === 'relay' ? '线路机' : '落地机' }}）</span></b>
              </div>
              <div class="d-field">
                <span class="d-label">当前同步状态</span>
                <div><SyncStatusTag v-if="detail" :status="detail.syncStatus" /></div>
              </div>
              <div class="d-field">
                <span class="d-label">最近成功同步</span>
                <b class="d-num">{{ formatDate(detail?.lastSyncAt) }}</b>
              </div>
              <div class="d-field">
                <span class="d-label">上一次耗时</span>
                <b class="d-num">{{ detail?.syncRuns.length ? formatDuration(detail.syncRuns[0]) : '--' }}</b>
              </div>
              <div class="d-field">
                <span class="d-label">同步周期</span>
                <b class="d-num">60 秒</b>
              </div>
              <div class="d-field">
                <span class="d-label">失败重试</span>
                <b>等下一周期（不重试）</b>
              </div>
            </div>
          </div>
          <div class="d-modal-foot">
            <span class="d-dim">本次请求预计 60 秒内执行</span>
            <div class="d-foot-actions">
              <button class="d-btn" type="button" :disabled="syncing" @click="closeSyncModal">取消</button>
              <button class="d-btn is-primary" type="button" :disabled="syncing" @click="confirmSync">
                <SvgIcon icon="mdi:refresh" class="d-btn-ico" />{{ syncing ? '提交中…' : '立即同步' }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- ===================== 停用 / 启用确认 ===================== -->
    <Teleport to="body">
      <div v-if="statusVisible" class="d-mask" @click.self="closeStatusModal">
        <div class="d-modal" style="width: 520px">
          <div class="d-modal-head">
            <div class="d-modal-head-main">
              <h3>{{ detail?.enabled ? '停用节点' : '启用节点' }}</h3>
              <p class="d-modal-sub">
                {{ detail?.enabled ? '停用后不再接受该节点的心跳与同步' : '启用后恢复接受该节点的心跳与同步' }}
              </p>
            </div>
            <span class="d-badge is-warn">需二次确认</span>
            <button class="d-icon-btn" type="button" :disabled="statusSaving" @click="closeStatusModal">
              <SvgIcon icon="mdi:close" />
            </button>
          </div>
          <div class="d-modal-body">
            <div v-if="statusError" class="d-alert is-danger mb-12">{{ statusError }}</div>
            <div class="d-callout" :class="detail?.enabled ? 'is-warn' : 'is-info'">
              <SvgIcon icon="mdi:alert-outline" class="d-callout-ico" />
              <div v-if="detail?.enabled">
                停用后该节点不再接受 Agent 心跳与同步；它承载的
                <b>{{ detail.userCount }} 个业务用户</b> 与 <b>{{ detail.clientCount }} 个启用 Client</b>
                的中央快照停留在最后一次同步结果，不再更新，也不会自动停用任何 X-Panel 配置。
              </div>
              <div v-else>启用后该节点恢复接受 Agent 心跳与同步，中央快照会在下一次成功同步后更新。</div>
            </div>
            <div class="d-detail-grid mt-12">
              <div class="d-field">
                <span class="d-label">节点</span>
                <b>{{ detail?.name }}<span class="d-dim">（{{ detail?.type === 'relay' ? '线路机' : '落地机' }}）</span></b>
              </div>
              <div class="d-field">
                <span class="d-label">当前状态</span>
                <div><NodeStatusTag v-if="detail" :status="detail.status" /></div>
              </div>
              <div class="d-field">
                <span class="d-label">承载业务用户</span>
                <b class="d-num">{{ detail?.userCount || 0 }} 户</b>
              </div>
              <div class="d-field">
                <span class="d-label">启用 Client</span>
                <b class="d-num">{{ detail?.clientCount || 0 }} 个</b>
              </div>
            </div>
            <div class="d-help mt-12">
              {{ detail?.enabled ? '停用可在「节点管理」列表中重新启用；停用不会删除任何中央记录。' : '启用不会修改 X-Panel 配置。' }}
            </div>
          </div>
          <div class="d-modal-foot">
            <span class="d-dim">历史同步记录与事件会保留</span>
            <div class="d-foot-actions">
              <button class="d-btn" type="button" :disabled="statusSaving" @click="closeStatusModal">取消</button>
              <button
                class="d-btn"
                :class="detail?.enabled ? 'is-danger' : 'is-success'"
                type="button"
                :disabled="statusSaving"
                @click="confirmStatusChange"
              >
                <SvgIcon icon="mdi:alert-outline" class="d-btn-ico" />
                {{ statusSaving ? '提交中…' : detail?.enabled ? '停用节点' : '启用节点' }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- ===================== 删除节点 ===================== -->
    <Teleport to="body">
      <div v-if="deleteVisible" class="d-mask" @click.self="closeDelete">
        <div class="d-modal" style="width: 560px">
          <div class="d-modal-head">
            <div class="d-modal-head-main">
              <h3>删除节点</h3>
              <p class="d-modal-sub">不可恢复 · 需要输入节点名确认</p>
            </div>
            <span class="d-badge is-danger">危险操作</span>
            <button class="d-icon-btn" type="button" :disabled="deleteSaving" @click="closeDelete">
              <SvgIcon icon="mdi:close" />
            </button>
          </div>
          <div class="d-modal-body">
            <div v-if="deleteError" class="d-alert is-danger mb-12">{{ deleteError }}</div>
            <div class="d-callout is-danger mb-12">
              <SvgIcon icon="mdi:alert-outline" class="d-callout-ico" />
              <div>
                <b>删除后无法恢复</b>
                将移除中央节点记录，并永久删除该节点的 Inbound、Client、流量快照、出口 IP 与成本记录；
                关联业务用户会被软删，账务记录与审计日志永久保留。
              </div>
            </div>
            <div class="d-field">
              <span class="d-label">
                输入节点名 <span class="d-mono">{{ detail?.name }}</span> 以确认 <span class="d-req">*</span>
              </span>
              <NInput v-model:value="deleteConfirmName" :placeholder="`输入 ${detail?.name || ''} 以确认删除`" />
              <span class="d-help">节点名不匹配时无法提交</span>
            </div>
          </div>
          <div class="d-modal-foot">
            <span class="d-dim">删除后该节点的出口 IP 会解除绑定</span>
            <div class="d-foot-actions">
              <button class="d-btn" type="button" :disabled="deleteSaving" @click="closeDelete">取消</button>
              <button
                class="d-btn is-danger"
                type="button"
                :disabled="deleteSaving || deleteConfirmName.trim() !== (detail?.name || '')"
                @click="confirmDelete"
              >
                <SvgIcon icon="mdi:delete-outline" class="d-btn-ico" />{{ deleteSaving ? '删除中…' : '永久删除' }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- ===================== 重新生成部署命令 ===================== -->
    <Teleport to="body">
      <div v-if="installerVisible" class="d-mask" @click.self="closeInstaller">
        <div class="d-modal" style="width: 700px">
          <div class="d-modal-head">
            <div class="d-modal-head-main">
              <h3>重新生成部署命令</h3>
              <p class="d-modal-sub">一次性 installerToken · 不会重新显示长期 token</p>
            </div>
            <button class="d-icon-btn" type="button" :disabled="installerLoading" @click="closeInstaller">
              <SvgIcon icon="mdi:close" />
            </button>
          </div>
          <div class="d-modal-body">
            <div v-if="installerError" class="d-alert is-danger mb-12">{{ installerError }}</div>
            <template v-if="installerResult">
              <div class="d-callout is-info mb-12">
                <SvgIcon icon="mdi:check-circle-outline" class="d-callout-ico" />
                <div>
                  <b>安装 Token 已生成</b>
                  15 分钟内有效且只能使用一次。请立即复制命令并在目标 Ubuntu / Debian 服务器执行。
                </div>
              </div>
              <div class="d-code-title">一键部署命令</div>
              <pre class="d-code">{{ installerCommand }}</pre>
              <div class="d-code-actions">
                <button class="d-btn d-btn-xs is-soft" type="button" @click="copyInstallerCommand">
                  <SvgIcon icon="mdi:content-copy" class="d-btn-ico" />{{ installerCopied ? '已复制' : '复制部署命令' }}
                </button>
                <span class="d-dim">安装完成后 Agent 自动启动并完成首次同步</span>
              </div>
              <div class="d-callout is-warn mt-12">
                <SvgIcon icon="mdi:alert-outline" class="d-callout-ico" />
                <div>
                  <b>不会重新显示长期 token</b>
                  该弹窗只返回一次性 installerToken；节点的长期 <span class="d-mono">central_token</span>
                  仅在创建时显示一次，之后无法找回。
                </div>
              </div>
            </template>
            <template v-else>
              <div class="d-callout is-info mb-12">
                <SvgIcon icon="mdi:information-outline" class="d-callout-ico" />
                <div>已设置管理地址的节点会自动带入 X-Panel 地址和 API 路径；历史节点未设置时可在这里补充。</div>
              </div>
              <div class="d-form-grid">
                <div class="d-field d-field-wide">
                  <span class="d-label">X-Panel 地址 <span class="d-req">*</span></span>
                  <NInput v-model:value="installerXpanelUrl" :readonly="hasManagementUrl" placeholder="例如 http://127.0.0.1:2053" />
                </div>
                <div class="d-field d-field-wide">
                  <span class="d-label">X-Panel API 路径</span>
                  <NInput v-model:value="installerBasePath" :readonly="hasManagementUrl" placeholder="默认 /；例如 /Alien" />
                </div>
              </div>
            </template>
          </div>
          <div class="d-modal-foot">
            <span class="d-dim">安装 Token 只能使用一次</span>
            <div class="d-foot-actions">
              <button class="d-btn" type="button" :disabled="installerLoading" @click="closeInstaller">关闭</button>
              <button v-if="!installerResult" class="d-btn is-primary" type="button" :disabled="installerLoading" @click="generateInstaller">
                {{ installerLoading ? '生成中…' : '生成命令' }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<style scoped>
/* ============================================================
   节点详情 · UI v2
   颜色 / 字号 / 间距对齐 design-preview/v2/nodes-detail.html 与 kit-v2
   ============================================================ */
.d-page {
  --d-card: #141519;
  --d-card-hover: #191b20;
  --d-line: #23262c;
  --d-line-soft: #1c1f24;
  --d-text-1: #f2f3f5;
  --d-text-2: #bcc2cc;
  --d-text-3: #9aa1ad;
  --d-text-4: #868d99;
  --d-brand: #4c6fff;
  --d-ok: #2ed47a;
  --d-warn: #f5b93d;
  --d-danger: #f3656f;
  --d-info: #5aa9f0;
  --d-violet: #a78bfa;
  --d-radius: 14px;

  display: flex;
  flex-direction: column;
  gap: 15px;
  color: var(--d-text-1);
  font-size: 13px;
  min-width: 0;
}

.d-dim { color: var(--d-text-4); }
.d-strong { color: var(--d-text-1); }
.d-mono { font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace; font-size: 12px; }
.d-num { font-variant-numeric: tabular-nums; }

/* ---------- 页头 ---------- */
.d-head { display: flex; align-items: flex-end; gap: 18px; flex-wrap: wrap; }
.d-head-main { min-width: 0; flex: 1; }
.d-crumb { display: flex; align-items: center; gap: 6px; margin-bottom: 7px; font-size: 12px; color: var(--d-text-4); }
.d-crumb button { padding: 0; border: 0; background: transparent; color: inherit; font-size: inherit; cursor: pointer; }
.d-crumb button:hover { color: #a8b8ff; }
.d-crumb-back { display: inline-flex; align-items: center; gap: 3px; }
.d-crumb-ico { width: 11px; height: 11px; }
.d-crumb-sep { color: #545b64; }
.d-title-row { display: flex; align-items: center; gap: 10px; flex-wrap: wrap; }
.d-title { margin: 0; font-size: 22px; font-weight: 620; line-height: 1.25; color: #f6f7f9; }
.d-meta-row { display: flex; align-items: center; gap: 8px; flex-wrap: wrap; margin-top: 8px; font-size: 12px; }
.d-mana-link { display: inline-flex; align-items: center; gap: 5px; color: #a8b8ff; text-decoration: none; }
.d-mana-link:hover { text-decoration: underline; }
.d-mana-ico { width: 13px; height: 13px; flex: none; }
.d-chip { display: inline-flex; align-items: center; gap: 5px; height: 24px; padding: 0 8px; border: 1px solid #2a2e35; border-radius: 8px; background: #14161a; color: #9aa1ab; font-size: 11.5px; }
.d-chip-ico { width: 12px; height: 12px; color: #6f7681; }
.d-desc { max-width: 860px; margin: 8px 0 0; font-size: 12px; line-height: 1.7; color: var(--d-text-3); }
.d-head-actions { display: flex; align-items: center; gap: 8px; flex: none; padding-bottom: 2px; }

/* ---------- 按钮 / 输入 ---------- */
.d-btn {
  display: inline-flex; align-items: center; justify-content: center; gap: 6px;
  height: 34px; padding: 0 13px;
  border: 1px solid #2a2e35; border-radius: 9px;
  background: #1b1e24; color: #e6e8ec; font-size: 12.5px; white-space: nowrap; cursor: pointer;
  transition: background 0.15s ease, border-color 0.15s ease;
}
.d-btn:hover { background: #232730; border-color: #363b44; }
.d-btn:disabled { opacity: 0.45; cursor: not-allowed; }
.d-btn-ico { width: 15px; height: 15px; flex: none; }
.d-btn.is-primary { border-color: var(--d-brand); background: var(--d-brand); color: #fff; }
.d-btn.is-primary:hover { background: #5d7cff; border-color: #5d7cff; }
.d-btn.is-soft { border-color: rgba(76, 111, 255, 0.28); background: rgba(76, 111, 255, 0.13); color: #a8b8ff; }
.d-btn.is-ghost { border-color: transparent; background: transparent; color: #a8b8ff; }
.d-btn.is-ghost:hover { background: rgba(76, 111, 255, 0.1); border-color: transparent; }
.d-btn.is-danger { border-color: rgba(243, 101, 111, 0.3); background: rgba(243, 101, 111, 0.13); color: #ff9ba2; }
.d-btn.is-success { border-color: rgba(46, 212, 122, 0.32); background: rgba(46, 212, 122, 0.14); color: #5ce49b; }
.d-btn-xs { height: 28px; padding: 0 10px; font-size: 12px; border-radius: 7px; }
.d-icon-btn {
  display: grid; place-items: center; width: 28px; height: 28px;
  border: 1px solid transparent; border-radius: 7px; background: transparent; color: #9aa1ab; cursor: pointer;
}
.d-icon-btn:hover { background: #1a1d22; color: #e6e8ec; }
.d-icon-btn:disabled { opacity: 0.4; cursor: not-allowed; }
.d-link-btn { padding: 0; border: 0; background: transparent; color: #8fa8ff; font-size: 12px; cursor: pointer; white-space: nowrap; }
.d-link-btn:hover { text-decoration: underline; }

/* ---------- KPI ---------- */
.d-kpis { display: grid; grid-template-columns: repeat(4, minmax(0, 1fr)); gap: 12px; }
.d-kpi { padding: 15px 16px 14px; border: 1px solid var(--d-line); border-radius: var(--d-radius); background: var(--d-card); min-width: 0; }
.d-kpi-top { display: flex; align-items: center; gap: 8px; margin-bottom: 10px; }
.d-kpi-label { font-size: 12px; color: var(--d-text-3); white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.d-kpi-icon { display: grid; place-items: center; width: 24px; height: 24px; margin-left: auto; flex: none; border-radius: 8px; border: 1px solid #2a2e35; background: #1d2026; color: #a8adb8; }
.d-kpi-icon :deep(.svg-icon) { width: 13px; height: 13px; }
.d-kpi-icon.is-blue { border-color: rgba(76, 111, 255, 0.36); background: rgba(76, 111, 255, 0.16); color: #8fa8ff; }
.d-kpi-icon.is-green { border-color: rgba(46, 212, 122, 0.34); background: rgba(46, 212, 122, 0.14); color: #5ce49b; }
.d-kpi-icon.is-cyan { border-color: rgba(52, 211, 224, 0.34); background: rgba(52, 211, 224, 0.14); color: #6fe0e9; }
.d-kpi-icon.is-amber { border-color: rgba(245, 185, 61, 0.34); background: rgba(245, 185, 61, 0.14); color: #f7ca6b; }
.d-kpi-value { display: flex; align-items: baseline; gap: 6px; font-size: 25px; font-weight: 620; line-height: 1.15; letter-spacing: -0.015em; font-variant-numeric: tabular-nums; color: #fff; }
.d-kpi-unit { font-size: 12px; font-weight: 500; color: var(--d-text-3); }
.d-kpi-foot { display: flex; align-items: center; gap: 6px; flex-wrap: wrap; margin-top: 9px; font-size: 12px; color: var(--d-text-4); }

/* ---------- 卡片 ---------- */
.d-card { position: relative; border: 1px solid var(--d-line); border-radius: var(--d-radius); background: var(--d-card); min-width: 0; }
.d-card-head { display: flex; align-items: center; gap: 10px; padding: 14px 18px; border-bottom: 1px solid var(--d-line-soft); }
.d-card-head h3 { margin: 0; font-size: 13.5px; font-weight: 600; color: var(--d-text-1); }
.d-right { display: flex; align-items: center; gap: 8px; margin-left: auto; flex: none; font-size: 12px; }
.d-card-body { padding: 16px 18px; }
.d-card-note { padding-top: 12px; border-top: 1px solid var(--d-line-soft); font-size: 11.5px; line-height: 1.7; color: var(--d-text-4); }

/* ---------- 元数据 ---------- */
.d-meta-grid { display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 14px 22px; }
.d-meta-wide { grid-column: 1 / -1; }
.d-meta-mgmt { grid-column: span 2; }
.d-field { display: flex; flex-direction: column; gap: 4px; min-width: 0; }
.d-field-wide { grid-column: 1 / -1; }
.d-label { font-size: 11.5px; color: var(--d-text-4); }
.d-field b { font-size: 13px; font-weight: 500; color: var(--d-text-1); }
.d-field-error { font-size: 11.5px; color: #ff9ba2; word-break: break-word; }
.d-inline { display: flex; align-items: center; gap: 8px; flex-wrap: wrap; }
.d-meta-key { display: flex; align-items: flex-start; gap: 8px; }
.d-meta-key .d-mono { line-height: 1.6; word-break: break-all; color: var(--d-text-2); }

/* ---------- 表格 ---------- */
.d-table-wrap { overflow-x: auto; }
.d-table { width: 100%; border-collapse: separate; border-spacing: 0; font-size: 12.5px; }
.d-table thead th { position: sticky; top: 0; z-index: 1; padding: 9px 12px; border-bottom: 1px solid var(--d-line); background: var(--d-card); color: var(--d-text-3); font-size: 11px; font-weight: 600; letter-spacing: 0.05em; text-align: left; white-space: nowrap; }
.d-table thead th.is-num { text-align: right; }
.d-table tbody td { padding: 8px 12px; border-bottom: 1px solid var(--d-line-soft); color: var(--d-text-2); vertical-align: middle; white-space: nowrap; }
.d-table tbody tr:last-child td { border-bottom: none; }
.d-table tbody tr:hover td { background: var(--d-card-hover); }
.d-table td.is-num { text-align: right; font-variant-numeric: tabular-nums; color: var(--d-text-1); }
.d-sub { display: block; margin-top: 1px; font-size: 11.5px; color: var(--d-text-4); white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.d-table-empty { padding: 28px 12px !important; text-align: center; color: var(--d-text-4); }
.d-error-cell { padding: 0 12px 10px !important; }
.d-error-text { margin: 0; padding: 8px 10px; border: 1px solid #23262c; border-radius: 8px; background: #0d0e11; color: #ffb3b9; font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace; font-size: 11.5px; line-height: 1.6; white-space: pre-wrap; word-break: break-word; }
.d-table-foot { display: flex; align-items: center; justify-content: space-between; gap: 10px; padding: 10px 18px; border-top: 1px solid var(--d-line-soft); font-size: 12px; }
.d-foot-actions { display: flex; align-items: center; gap: 8px; margin-left: auto; }

/* ---------- 两列 ---------- */
.d-two-col { display: grid; grid-template-columns: minmax(0, 1.5fr) minmax(0, 1fr); gap: 12px; align-items: start; }

/* ---------- 事件列表 ---------- */
.d-event-list { max-height: 560px; overflow-y: auto; }
.d-event-row { display: flex; align-items: flex-start; gap: 10px; padding: 11px 18px; border-bottom: 1px solid var(--d-line-soft); }
.d-event-row:last-child { border-bottom: none; }
.d-event-main { flex: 1; min-width: 0; }
.d-event-main b { display: block; font-size: 12.5px; font-weight: 550; color: var(--d-text-1); }
.d-event-main span { display: block; margin-top: 2px; font-size: 11.5px; line-height: 1.6; color: var(--d-text-4); word-break: break-word; }

/* ---------- 提示条 ---------- */
.d-callout { display: flex; gap: 10px; padding: 12px 13px; border: 1px solid var(--d-line); border-radius: 11px; background: rgba(255, 255, 255, 0.022); font-size: 12px; line-height: 1.7; color: var(--d-text-2); }
.d-callout-ico { width: 16px; height: 16px; flex: none; margin-top: 2px; color: var(--d-text-3); }
.d-callout b { display: block; margin-bottom: 1px; font-size: 12.5px; font-weight: 600; color: var(--d-text-1); }
.d-callout.is-info { border-color: rgba(90, 169, 240, 0.24); background: rgba(90, 169, 240, 0.08); }
.d-callout.is-info .d-callout-ico { color: #8cc4f7; }
.d-callout.is-warn { border-color: rgba(245, 185, 61, 0.28); background: rgba(245, 185, 61, 0.08); }
.d-callout.is-warn .d-callout-ico { color: #f7ca6b; }
.d-callout.is-danger { border-color: rgba(243, 101, 111, 0.28); background: rgba(243, 101, 111, 0.08); }
.d-callout.is-danger .d-callout-ico { color: #ff9ba2; }

/* ---------- 徽标 ---------- */
.d-badge { display: inline-flex; align-items: center; gap: 5px; height: 21px; padding: 0 7px; border: 1px solid transparent; border-radius: 6px; font-size: 11.5px; font-weight: 550; line-height: 1; white-space: nowrap; }
.d-badge.is-neutral { border-color: #2a2e35; background: #1d2026; color: #bcc2cc; }
.d-badge.is-brand { border-color: rgba(76, 111, 255, 0.34); background: rgba(76, 111, 255, 0.16); color: #a8b8ff; }
.d-badge.is-ok { border-color: rgba(46, 212, 122, 0.32); background: rgba(46, 212, 122, 0.15); color: #5ce49b; }
.d-badge.is-warn { border-color: rgba(245, 185, 61, 0.32); background: rgba(245, 185, 61, 0.15); color: #f7ca6b; }
.d-badge.is-danger { border-color: rgba(243, 101, 111, 0.32); background: rgba(243, 101, 111, 0.15); color: #ff9ba2; }
.d-badge.is-info { border-color: rgba(90, 169, 240, 0.32); background: rgba(90, 169, 240, 0.15); color: #8cc4f7; }
.d-badge.is-violet { border-color: rgba(167, 139, 250, 0.34); background: rgba(167, 139, 250, 0.16); color: #c0a9ff; }

/* ---------- 弹窗 ---------- */
.d-mask {
  position: fixed;
  inset: 0;
  z-index: 2000;
  display: grid;
  place-items: center;
  padding: 24px;
  background: rgba(0, 0, 0, 0.66);
}
.d-modal { display: flex; flex-direction: column; max-height: min(86vh, 880px); border: 1px solid var(--d-line); border-radius: 16px; background: var(--d-card); color: var(--d-text-1); text-align: left; box-shadow: 0 30px 80px -24px rgba(0, 0, 0, 0.9); }
.d-modal-head { display: flex; align-items: flex-start; gap: 10px; padding: 16px 20px 14px; border-bottom: 1px solid var(--d-line-soft); }
.d-modal-head-main { flex: 1; min-width: 0; }
.d-modal-head h3 { margin: 0; font-size: 15px; font-weight: 600; }
.d-modal-sub { margin: 3px 0 0; font-size: 12px; color: var(--d-text-4); }
.d-modal-body { flex: 1; min-height: 0; overflow-y: auto; padding: 16px 20px 20px; }
.d-modal-foot { display: flex; align-items: center; gap: 10px; padding: 13px 20px; border-top: 1px solid var(--d-line-soft); font-size: 12px; }
.d-form-grid { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 12px; }
.d-detail-grid { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 12px 20px; }
.d-help { font-size: 11.5px; line-height: 1.6; color: var(--d-text-4); }
.d-req { color: var(--d-danger); }
.d-code-title { margin-bottom: 7px; font-size: 12px; color: var(--d-text-2); }
.d-code { margin: 0; padding: 11px 13px; border: 1px solid #23262c; border-radius: 9px; background: #0d0e11; color: #b8bec8; font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace; font-size: 11.5px; line-height: 1.7; white-space: pre-wrap; word-break: break-word; }
.d-code-actions { display: flex; align-items: center; gap: 8px; margin-top: 8px; font-size: 11.5px; }

/* ---------- 加载 / 空态 / 错误 ---------- */
.d-loading, .d-empty { display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 10px; min-height: 220px; padding: 24px; border: 1px solid var(--d-line); border-radius: var(--d-radius); background: var(--d-card); color: var(--d-text-4); font-size: 12.5px; }
.d-empty-title { font-size: 14px; color: var(--d-text-1); }
.d-empty-desc { margin: 0; color: var(--d-text-4); }
.d-alert { padding: 10px 13px; border: 1px solid rgba(243, 101, 111, 0.32); border-radius: 10px; background: rgba(243, 101, 111, 0.1); color: #ffb3b9; font-size: 12.5px; }
.d-footnote { font-size: 12px; line-height: 1.7; color: var(--d-text-4); }

.mt-12 { margin-top: 12px; }
.mb-12 { margin-bottom: 12px; }

@media (max-width: 1280px) {
  .d-meta-grid { grid-template-columns: repeat(2, minmax(0, 1fr)); }
  .d-meta-mgmt { grid-column: 1 / -1; }
  .d-two-col { grid-template-columns: minmax(0, 1fr); }
}
@media (max-width: 1100px) {
  .d-kpis { grid-template-columns: repeat(2, minmax(0, 1fr)); }
  .d-form-grid, .d-detail-grid { grid-template-columns: minmax(0, 1fr); }
}
</style>
