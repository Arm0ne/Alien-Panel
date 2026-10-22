<script setup lang="ts">
import { computed, h, onMounted, reactive, ref } from 'vue';
import { NDropdown, NInput, NPagination, NSelect } from 'naive-ui';
import { useRouter } from 'vue-router';
import { createNode, deleteNode, fetchNodes, requestNodeSync, updateNode } from '@/service/api';
import DataFreshness from '@/components/project/data-freshness.vue';
import SvgIcon from '@/components/custom/svg-icon.vue';
import NodeTypeTag from '@/components/project/node-type-tag.vue';

defineOptions({ name: 'NodeManagement' });

const router = useRouter();

const loading = ref(false);
const syncingNodeId = ref('');
const togglingNodeId = ref('');
const deletingNodeId = ref('');
const errorMessage = ref('');
const rows = ref<Api.Central.NodeSummary[]>([]);
const total = ref(0);
const stats = ref<Api.Central.NodeListStats | null>(null);
const dataAt = ref('');
const selectedType = ref('all');
const onboardingVisible = ref(false);
const onboardingSaving = ref(false);
const onboardingError = ref('');
const onboardingResult = ref<Api.Central.NodeRegistrationResult | null>(null);
const wizardCopied = ref(false);

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

const typeSegments = [
  { label: '全部类型', value: 'all' },
  { label: '线路机', value: 'relay' },
  { label: '落地机', value: 'landing' }
];

const addNodeTypeOptions = [
  { label: '线路机', value: 'relay' },
  { label: '落地机', value: 'landing' }
];

function formatDate(value?: string | null) {
  if (!value) return '--';
  const date = new Date(value);
  return Number.isNaN(date.getTime()) ? '时间未知' : date.toLocaleString('zh-CN', { hour12: false });
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

function statusDotClass(row: Api.Central.NodeSummary) {
  if (!row.enabled) return 'is-neutral';
  if (row.status === 'online') return 'is-ok';
  if (row.status === 'degraded') return 'is-warn';
  if (row.status === 'offline') return 'is-danger';
  return 'is-neutral';
}

function syncLabel(row: Api.Central.NodeSummary) {
  if (row.syncStatus === 'success') return '同步成功';
  if (row.syncStatus === 'failed') return '同步失败';
  return '状态未知';
}

function syncDotClass(row: Api.Central.NodeSummary) {
  if (row.syncStatus === 'success') return 'is-ok';
  if (row.syncStatus === 'failed') return 'is-danger';
  return 'is-neutral';
}

const nodeListHint = computed(() => {
  const count = rows.value.length;
  return `共 ${total.value} 个节点${count ? ` · 当前显示 ${count} 个` : ''} · 刷新于 ${formatDate(dataAt.value || new Date().toISOString())}`;
});

const pageCount = computed(() => Math.max(1, Math.ceil(total.value / filters.page_size)));

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
  return [
    `node_key: ${onboardingResult.value.nodeKey}`,
    `node_name: ${onboardingResult.value.name}`,
    `node_type: ${onboardingResult.value.type}`,
    `central_url: ${centralApiUrl()}`,
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
  wizardCopied.value = false;
}

function syncAutoFields() {
  const parsed = parseManagementAddress(nodeForm.managementUrl || '');
  xpanelUrl.value = parsed?.url || '';
  nodeForm.panelBasePath = parsed?.basePath || '/';
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

async function copyWizardValue(value: string, successMessage: string) {
  await copyText(value, successMessage);
  if (onboardingResult.value) wizardCopied.value = true;
}

async function saveNode() {
  syncAutoFields();
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
  if (!nodeForm.name.trim()) {
    onboardingError.value = '请填写节点名称';
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
  onboardingError.value = '';
  onboardingSaving.value = true;
  const { data, error } = await createNode(payload);
  onboardingSaving.value = false;
  if (error || !data) {
    onboardingError.value = '节点创建失败，请检查节点信息后重试';
    return;
  }
  onboardingResult.value = data;
  wizardCopied.value = false;
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
    stats.value = null;
    dataAt.value = '';
  } else {
    rows.value = data.items;
    total.value = data.total;
    stats.value = data.stats;
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
  void loadNodes();
}

function selectType(value: string) {
  selectedType.value = value;
  submitFilters();
}

function resetFilters() {
  filters.keyword = '';
  filters.status = 'all';
  selectedType.value = 'all';
  submitFilters();
}

function onPageChange(page: number) {
  filters.page = page;
  void loadNodes();
}

function nodeMenuOptions(row: Api.Central.NodeSummary) {
  return [
    {
      label: '查看详情',
      key: 'detail',
      icon: () => h(SvgIcon, { icon: 'mdi:chevron-right' })
    },
    {
      label: '立即同步',
      key: 'sync',
      disabled: !row.enabled || Boolean(syncingNodeId.value),
      icon: () => h(SvgIcon, { icon: 'mdi:refresh' })
    },
    {
      label: '复制管理地址',
      key: 'copy',
      disabled: !row.managementUrl,
      icon: () => h(SvgIcon, { icon: 'mdi:content-copy' })
    },
    { type: 'divider', key: 'd1' },
    {
      label: row.enabled ? '停用节点' : '启用节点',
      key: 'toggle',
      disabled: Boolean(togglingNodeId.value),
      icon: () => h(SvgIcon, { icon: row.enabled ? 'mdi:power' : 'mdi:check-circle-outline' })
    },
    {
      label: '删除节点',
      key: 'delete',
      disabled: Boolean(deletingNodeId.value),
      icon: () => h(SvgIcon, { icon: 'mdi:delete-outline' })
    }
  ];
}

function handleNodeAction(key: string, row: Api.Central.NodeSummary) {
  if (key === 'detail') {
    openDetail(row.id);
    return;
  }
  if (key === 'sync') {
    void syncNode(row.id);
    return;
  }
  if (key === 'copy') {
    if (row.managementUrl) void copyText(row.managementUrl, '管理地址已复制');
    return;
  }
  if (key === 'toggle') {
    void toggleNode(row);
    return;
  }
  if (key === 'delete') {
    confirmDelete(row);
  }
}

onMounted(() => {
  void loadNodes();
});
</script>

<template>
  <div class="n-page">
    <!-- 页头 -->
    <header class="n-head">
      <div class="n-head-main">
        <div class="n-eyebrow">
          业务 · 节点管理 · 数据更新于 {{ relativeTime(dataAt) }}
        </div>
        <h1 class="n-title">节点管理</h1>
        <p class="n-desc">
          只读同步各节点本机 X-Panel，不写入 Xray。点击卡片里的「管理地址」可直接打开该节点的 X-Panel 登录页，方便一键跳转管理。
        </p>
      </div>
      <div class="n-head-actions">
        <DataFreshness :data-at="dataAt" compact />
        <button class="n-btn" type="button" :disabled="loading" @click="loadNodes">
          <SvgIcon icon="mdi:refresh" class="n-btn-ico" />刷新
        </button>
        <button class="n-btn is-primary" type="button" @click="openOnboarding">
          <SvgIcon icon="mdi:plus" class="n-btn-ico" />接入新节点
        </button>
      </div>
    </header>

    <div v-if="errorMessage" class="n-alert is-danger">{{ errorMessage }}</div>

    <!-- KPI -->
    <section class="n-kpis">
      <div class="n-kpi">
        <div class="n-kpi-top">
          <span class="n-kpi-label">节点总数</span>
          <span class="n-kpi-icon is-blue"><SvgIcon icon="mdi:server-network" /></span>
        </div>
        <div class="n-kpi-value">{{ stats ? stats.total : '--' }}<span class="n-kpi-unit">台</span></div>
        <div class="n-kpi-foot">已接入且未删除</div>
      </div>
      <div class="n-kpi">
        <div class="n-kpi-top">
          <span class="n-kpi-label">在线节点</span>
          <span class="n-kpi-icon is-green"><SvgIcon icon="mdi:activity" /></span>
        </div>
        <div class="n-kpi-value">
          {{ stats ? stats.online : '--' }}<span class="n-kpi-unit">/ {{ stats ? stats.total : '--' }}</span>
        </div>
        <div class="n-kpi-foot">启用且健康状态正常</div>
      </div>
      <div class="n-kpi">
        <div class="n-kpi-top">
          <span class="n-kpi-label">线路机</span>
          <span class="n-kpi-icon is-violet"><SvgIcon icon="mdi:source-branch" /></span>
        </div>
        <div class="n-kpi-value">{{ stats ? stats.relay : '--' }}<span class="n-kpi-unit">台</span></div>
        <div class="n-kpi-foot">承载业务用户 Inbound</div>
      </div>
      <div class="n-kpi">
        <div class="n-kpi-top">
          <span class="n-kpi-label">落地机</span>
          <span class="n-kpi-icon is-amber"><SvgIcon icon="mdi:globe" /></span>
        </div>
        <div class="n-kpi-value">{{ stats ? stats.landing : '--' }}<span class="n-kpi-unit">台</span></div>
        <div class="n-kpi-foot">用于落地出口和路径</div>
      </div>
    </section>

    <!-- 工具条 -->
    <section class="n-toolbar">
      <label class="n-input n-search">
        <SvgIcon icon="mdi:magnify" class="n-input-ico" />
        <input v-model="filters.keyword" placeholder="搜索节点名称或管理地址" @keyup.enter="submitFilters" />
      </label>
      <div class="n-seg">
        <button
          v-for="segment in typeSegments"
          :key="segment.value"
          type="button"
          :class="{ 'is-active': selectedType === segment.value }"
          @click="selectType(segment.value)"
        >
          {{ segment.label }}
        </button>
      </div>
      <NSelect
        v-model:value="filters.status"
        :options="statusOptions"
        size="small"
        style="width: 140px"
        @update:value="submitFilters"
      />
      <span class="n-spacer"></span>
      <span class="n-dim n-hint">{{ nodeListHint }}</span>
    </section>

    <!-- 节点卡片 -->
    <div v-if="loading && rows.length === 0" class="n-loading">正在读取节点…</div>
    <div v-else-if="!rows.length" class="n-empty">
      <div class="n-empty-title">暂无节点</div>
      <p class="n-empty-desc">点击「接入新节点」创建中央记录，然后按模板配置 Agent。</p>
      <button class="n-btn is-primary" type="button" @click="openOnboarding">
        <SvgIcon icon="mdi:plus" class="n-btn-ico" />接入新节点
      </button>
    </div>
    <div v-else class="n-grid">
      <article
        v-for="row in rows"
        :key="row.id"
        class="n-card n-node-card"
        :class="{ 'is-offline': row.status === 'offline', 'is-disabled': !row.enabled }"
      >
        <div class="n-card-head">
          <h3>
            <i class="n-dot" :class="statusDotClass(row)"></i>
            {{ row.name }}
            <NodeTypeTag :type="row.type" />
            <span v-if="!row.enabled" class="n-badge is-neutral">已停用</span>
            <span v-else-if="row.status === 'offline'" class="n-badge is-danger">离线</span>
            <span v-else-if="row.status === 'degraded'" class="n-badge is-warn">异常</span>
          </h3>
          <div class="n-right">
            <span class="n-dim n-region">{{ row.region || '—' }}</span>
            <NDropdown trigger="click" :options="nodeMenuOptions(row)" @select="key => handleNodeAction(String(key), row)">
              <button class="n-icon-btn" type="button" aria-label="更多操作">
                <SvgIcon icon="mdi:dots-horizontal" />
              </button>
            </NDropdown>
          </div>
        </div>

        <a
          v-if="row.managementUrl"
          class="n-mana"
          :href="row.managementUrl"
          target="_blank"
          rel="noopener noreferrer"
          :title="row.managementUrl"
        >
          <SvgIcon icon="mdi:link-variant" class="n-mana-ico" />
          <span class="n-mana-url">{{ row.managementUrl }}</span>
          <SvgIcon icon="mdi:open-in-new" class="n-mana-ico" />
        </a>
        <div v-else class="n-mana is-empty">
          <SvgIcon icon="mdi:link-off" class="n-mana-ico" />
          <span class="n-mana-url">未设置管理地址</span>
        </div>

        <div class="n-stats">
          <div>
            <span>业务用户</span>
            <b v-if="row.type === 'relay'">{{ row.userCount }}<small>户</small></b>
            <b v-else class="n-dim">—</b>
            <i v-if="row.type === 'relay'" class="n-mini">收费 {{ row.paidUserCount }} / 免费 {{ row.freeUserCount }}</i>
            <i v-else class="n-mini">落地机不承载业务用户</i>
          </div>
          <div>
            <span>启用 Client</span>
            <b>{{ row.clientCount }}<small>台</small></b>
          </div>
          <div>
            <span>出口 IP</span>
            <b>{{ row.exitIpCount }}<small>个</small></b>
          </div>
          <div>
            <span>月成本</span>
            <b class="n-money">{{ formatMoney(row.monthlyCost, row.currency) }}</b>
          </div>
        </div>

        <div class="n-card-foot">
          <div class="n-foot-meta">
            <b>
              <i class="n-dot" :class="syncDotClass(row)"></i>
              {{ syncLabel(row) }}<template v-if="row.lastSyncAt"> · {{ relativeTime(row.lastSyncAt) }}</template>
            </b>
            <span>
              Xray {{ row.xrayVersion || '未检测到' }} · 最近成功同步 {{ formatDate(row.lastSyncAt) }}
            </span>
          </div>
          <div class="n-foot-actions">
            <button v-if="row.status === 'offline'" class="n-btn n-btn-xs" type="button" @click="openDetail(row.id)">排查</button>
            <button class="n-btn n-btn-xs is-soft" type="button" @click="openDetail(row.id)">详情</button>
            <button
              class="n-btn n-btn-xs"
              type="button"
              :disabled="!row.enabled || syncingNodeId === row.id"
              @click="syncNode(row.id)"
            >
              {{ syncingNodeId === row.id ? '同步中…' : '立即同步' }}
            </button>
          </div>
        </div>
      </article>

      <!-- 接入引导卡 -->
      <div class="n-card n-add-card" @click="openOnboarding">
        <div class="n-add-inner">
          <div class="n-kpi-icon is-blue"><SvgIcon icon="mdi:plus" /></div>
          <b>接入新节点</b>
          <span class="n-dim">填写管理地址后生成<br />一键部署命令与 Token</span>
          <button class="n-btn n-btn-xs is-primary" type="button" @click.stop="openOnboarding">开始接入</button>
        </div>
      </div>
    </div>

    <!-- 分页 -->
    <div v-if="total > filters.page_size" class="n-pager">
      <span class="n-dim">第 {{ filters.page }} / {{ pageCount }} 页 · 共 {{ total }} 个节点</span>
      <NPagination
        :page="filters.page"
        :page-count="pageCount"
        :page-slot="7"
        size="small"
        @update:page="onPageChange"
      />
    </div>

    <!-- ===================== 接入向导 ===================== -->
    <Teleport to="body">
      <div v-if="onboardingVisible" class="n-mask" @click.self="closeOnboarding">
        <div class="n-modal" style="width: 680px">
        <div class="n-modal-head">
          <div class="n-modal-head-main">
            <h3>{{ onboardingResult ? '节点接入完成' : '接入新节点' }}</h3>
            <p class="n-modal-sub">
              {{ onboardingResult ? '创建中央节点记录并交付 Agent 部署三件套' : '创建中央节点记录后生成一键部署命令与 Token' }}
            </p>
          </div>
          <span class="n-badge" :class="onboardingResult ? 'is-brand' : 'is-neutral'">
            {{ onboardingResult ? '第 3 / 3 步' : '第 1 / 3 步' }}
          </span>
          <button class="n-icon-btn" type="button" :disabled="onboardingSaving" @click="closeOnboarding">
            <SvgIcon icon="mdi:close" />
          </button>
        </div>

        <div class="n-modal-body">
          <div v-if="onboardingError" class="n-alert is-danger mb-12">{{ onboardingError }}</div>

          <div class="n-steps">
            <div class="n-step" :class="onboardingResult ? 'is-done' : 'is-active'">
              <i>{{ onboardingResult ? '✓' : '1' }}</i>
              <div>
                <b>填写管理地址</b>
                <span>{{ onboardingResult ? '已解析并保存管理地址' : '解析 X-Panel 地址与 API 路径' }}</span>
              </div>
            </div>
            <div class="n-step" :class="onboardingResult ? 'is-done' : ''">
              <i>{{ onboardingResult ? '✓' : '2' }}</i>
              <div>
                <b>创建节点</b>
                <span>{{ onboardingResult ? '已生成 node_key 与长期 token' : '生成中央记录与一次性凭证' }}</span>
              </div>
            </div>
            <div class="n-step" :class="onboardingResult ? 'is-active' : ''">
              <i>3</i>
              <div>
                <b>部署 Agent</b>
                <span>{{ onboardingResult ? '复制命令与 token 到目标服务器' : '复制命令、token 与配置模板' }}</span>
              </div>
            </div>
          </div>

          <!-- 表单态 -->
          <template v-if="!onboardingResult">
            <div class="n-callout is-info">
              <SvgIcon icon="mdi:information-outline" class="n-callout-ico" />
              <div>
                选择节点类型并填写管理地址后，创建完成即可复制 Ubuntu / Debian 一键部署命令。节点 Key 由中央自动生成。
              </div>
            </div>

            <div class="n-form-grid mt-16">
              <div class="n-field">
                <span class="n-label">节点名称 <span class="n-req">*</span></span>
                <NInput v-model:value="nodeForm.name" maxlength="120" placeholder="例如：东京线路机 01" />
              </div>
              <div class="n-field">
                <span class="n-label">节点类型 <span class="n-req">*</span></span>
                <NSelect v-model:value="nodeForm.type" :options="addNodeTypeOptions" />
              </div>
              <div class="n-field">
                <span class="n-label">主机名</span>
                <NInput v-model:value="nodeForm.hostname" maxlength="255" placeholder="可选，例如 relay.example.com" />
              </div>
              <div class="n-field">
                <span class="n-label">区域</span>
                <NInput v-model:value="nodeForm.region" maxlength="120" placeholder="可选，例如 东京" />
              </div>
              <div class="n-field">
                <span class="n-label">服务商</span>
                <NInput v-model:value="nodeForm.provider" maxlength="200" placeholder="可选，例如 Example Cloud" />
              </div>
              <div class="n-field">
                <span class="n-label">X-Panel 地址（自动解析）</span>
                <NInput :value="xpanelUrl || '--'" readonly disabled />
              </div>
              <div class="n-field n-field-wide">
                <span class="n-label">管理地址 <span class="n-req">*</span></span>
                <NInput
                  v-model:value="nodeForm.managementUrl"
                  maxlength="2048"
                  placeholder="例如 https://server.example.com:18086/Alien/"
                  @blur="syncAutoFields"
                />
                <span class="n-help">保存后会同时写入 managementUrl 与 panelBasePath</span>
              </div>
              <div class="n-field n-field-wide">
                <span class="n-label">节点出口 IP（可选，可添加多个）</span>
                <div class="n-exit-ip-list">
                  <div v-for="(_, index) in exitIpInputs" :key="index" class="n-exit-ip-row">
                    <NInput
                      v-model:value="exitIpInputs[index]"
                      maxlength="45"
                      placeholder="例如 38.117.67.14 或 2001:db8::10"
                    />
                    <button
                      v-if="index === exitIpInputs.length - 1"
                      class="n-icon-btn"
                      type="button"
                      aria-label="增加出口 IP"
                      @click="addExitIpInput"
                    >
                      <SvgIcon icon="mdi:plus-circle-outline" />
                    </button>
                    <button
                      v-else
                      class="n-icon-btn"
                      type="button"
                      aria-label="删除出口 IP"
                      @click="removeExitIpInput(index)"
                    >
                      <SvgIcon icon="mdi:minus-circle-outline" />
                    </button>
                  </div>
                </div>
                <span class="n-help">这些地址会登记为该节点名下的出口 IP，最多 100 个</span>
              </div>
            </div>
          </template>

          <!-- 完成态 -->
          <template v-else>
            <div class="n-callout is-info">
              <SvgIcon icon="mdi:check-circle-outline" class="n-callout-ico" />
              <div>
                <b>节点 {{ onboardingResult.name }} 已创建</b>
                {{ onboardingResult.type === 'relay' ? '线路机' : '落地机' }} ·
                {{ nodeForm.region || '未设置区域' }} · 已登记 {{ onboardingResult.exitIpCount }} 个出口 IP
              </div>
            </div>

            <div class="n-code-title">Ubuntu / Debian 一键部署命令（推荐）</div>
            <pre class="n-code">{{ agentInstallCommand }}</pre>
            <div class="n-code-actions">
              <button class="n-btn n-btn-xs is-soft" type="button" @click="copyWizardValue(agentInstallCommand, '一键部署命令已复制')">
                <SvgIcon icon="mdi:content-copy" class="n-btn-ico" />复制命令
              </button>
              <span class="n-dim">脚本会自动安装 Agent、写入 /etc/xpanel-agent/agent.yaml 并启动 systemd 服务</span>
            </div>

            <div class="n-code-title mt-16">Install token（一次性 · 15 分钟内有效）</div>
            <pre class="n-code">{{ onboardingResult.installerToken }}</pre>
            <div class="n-code-actions">
              <button
                class="n-btn n-btn-xs is-soft"
                type="button"
                @click="copyWizardValue(onboardingResult.installerToken, 'install token 已复制')"
              >
                <SvgIcon icon="mdi:content-copy" class="n-btn-ico" />复制 token
              </button>
              <span class="n-dim">用于 bootstrap 换取 Agent 正式凭证，只能使用一次</span>
            </div>

            <div class="n-callout is-warn mt-12">
              <SvgIcon icon="mdi:alert-outline" class="n-callout-ico" />
              <div>
                <b>token 只显示一次，关闭后无法找回</b>
                install token 与节点长期 <span class="n-mono">central_token</span> 都只在本次弹窗内显示。关闭弹窗即无法再次查看；
                丢失后只能删除节点重建，或在节点详情「重新生成部署命令」（不会重新显示长期 token）。请先复制并确认已保存，再关闭本窗口。
                <div class="n-badge-row mt-8">
                  <span class="n-badge is-warn">15 分钟内有效</span>
                  <span class="n-badge is-warn">仅可使用一次</span>
                  <span class="n-badge is-warn">关闭后不可恢复</span>
                </div>
              </div>
            </div>

            <div class="n-code-title mt-16">Agent 配置模板（手动部署备用）</div>
            <pre class="n-code">{{ agentConfig }}</pre>
            <div class="n-code-actions">
              <button class="n-btn n-btn-xs is-soft" type="button" @click="copyWizardValue(agentConfig, 'Agent 配置已复制')">
                <SvgIcon icon="mdi:content-copy" class="n-btn-ico" />复制配置
              </button>
              <span class="n-dim">凭据永不进命令行；install 脚本会在节点本机交互询问 X-Panel 账号密码</span>
            </div>
          </template>
        </div>

        <div class="n-modal-foot">
          <span class="n-dim">
            {{ onboardingResult ? (wizardCopied ? '已复制，可安全关闭' : '请先复制命令或 token 再关闭') : '完成后 Agent 将在 60 秒内完成首次同步' }}
          </span>
          <div class="n-foot-actions">
            <button class="n-btn" type="button" :disabled="onboardingSaving" @click="closeOnboarding">关闭</button>
            <button
              v-if="!onboardingResult"
              class="n-btn is-primary"
              type="button"
              :disabled="onboardingSaving"
              @click="saveNode"
            >
              {{ onboardingSaving ? '创建中…' : '创建并生成 Token' }}
            </button>
            <button
              v-else
              class="n-btn is-primary"
              type="button"
              :disabled="!wizardCopied"
              :title="wizardCopied ? '' : '请先复制部署命令或 token'"
              @click="closeOnboarding"
            >
              <SvgIcon icon="mdi:check" class="n-btn-ico" />我已保存，关闭
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
   节点管理 · UI v2
   颜色 / 字号 / 间距对齐 design-preview/v2/nodes.html 与 kit-v2
   ============================================================ */
.n-page {
  --n-card: #141519;
  --n-card-hover: #191b20;
  --n-line: #23262c;
  --n-line-soft: #1c1f24;
  --n-text-1: #f2f3f5;
  --n-text-2: #bcc2cc;
  --n-text-3: #9aa1ad;
  --n-text-4: #868d99;
  --n-brand: #4c6fff;
  --n-brand-soft: rgba(76, 111, 255, 0.14);
  --n-ok: #2ed47a;
  --n-warn: #f5b93d;
  --n-danger: #f3656f;
  --n-info: #5aa9f0;
  --n-violet: #a78bfa;
  --n-radius: 14px;

  display: flex;
  flex-direction: column;
  gap: 15px;
  color: var(--n-text-1);
  font-size: 13px;
  min-width: 0;
}

.n-dim { color: var(--n-text-4); }
.n-mono { font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace; }
.n-money { font-variant-numeric: tabular-nums; }

/* ---------- 页头 ---------- */
.n-head { display: flex; align-items: flex-end; gap: 18px; flex-wrap: wrap; }
.n-head-main { min-width: 0; flex: 1; }
.n-eyebrow { margin-bottom: 6px; font-size: 12px; letter-spacing: 0.1em; color: #7f8792; }
.n-title { margin: 0; font-size: 22px; font-weight: 620; line-height: 1.25; color: #f6f7f9; }
.n-desc { max-width: 840px; margin: 7px 0 0; font-size: 12px; line-height: 1.7; color: var(--n-text-3); }
.n-head-actions { display: flex; align-items: center; gap: 8px; flex: none; padding-bottom: 2px; }

/* ---------- 按钮 / 输入 ---------- */
.n-btn {
  display: inline-flex; align-items: center; justify-content: center; gap: 6px;
  height: 34px; padding: 0 13px;
  border: 1px solid #2a2e35; border-radius: 9px;
  background: #1b1e24; color: #e6e8ec;
  font-size: 12.5px; white-space: nowrap; cursor: pointer;
  transition: background 0.15s ease, border-color 0.15s ease;
}
.n-btn:hover { background: #232730; border-color: #363b44; }
.n-btn:disabled { opacity: 0.45; cursor: not-allowed; }
.n-btn-ico { width: 15px; height: 15px; flex: none; }
.n-btn.is-primary { border-color: var(--n-brand); background: var(--n-brand); color: #fff; }
.n-btn.is-primary:hover { background: #5d7cff; border-color: #5d7cff; }
.n-btn.is-soft { border-color: rgba(76, 111, 255, 0.28); background: rgba(76, 111, 255, 0.13); color: #a8b8ff; }
.n-btn-xs { height: 28px; padding: 0 10px; font-size: 12px; border-radius: 7px; }

.n-icon-btn {
  display: grid; place-items: center; width: 28px; height: 28px;
  border: 1px solid transparent; border-radius: 7px; background: transparent;
  color: #9aa1ab; cursor: pointer;
}
.n-icon-btn:hover { background: #1a1d22; color: #e6e8ec; }
.n-icon-btn:disabled { opacity: 0.4; cursor: not-allowed; }

.n-input {
  display: flex; align-items: center; gap: 8px; height: 34px; padding: 0 10px;
  border: 1px solid #262a31; border-radius: 9px; background: #0e1013;
}
.n-input:focus-within { border-color: rgba(76, 111, 255, 0.7); }
.n-input input { min-width: 0; flex: 1; border: 0; outline: none; background: transparent; color: var(--n-text-1); font-size: 12.5px; font-family: inherit; }
.n-input input::placeholder { color: var(--n-text-4); }
.n-input-ico { width: 15px; height: 15px; flex: none; color: var(--n-text-4); }
.n-search { width: 260px; flex: none; }

/* ---------- KPI ---------- */
.n-kpis { display: grid; grid-template-columns: repeat(4, minmax(0, 1fr)); gap: 12px; }
.n-kpi { padding: 15px 16px 14px; border: 1px solid var(--n-line); border-radius: var(--n-radius); background: var(--n-card); min-width: 0; }
.n-kpi-top { display: flex; align-items: center; gap: 8px; margin-bottom: 10px; }
.n-kpi-label { font-size: 12px; color: var(--n-text-3); white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.n-kpi-icon { display: grid; place-items: center; width: 24px; height: 24px; margin-left: auto; flex: none; border-radius: 8px; border: 1px solid #2a2e35; background: #1d2026; color: #a8adb8; }
.n-kpi-icon :deep(.svg-icon) { width: 13px; height: 13px; }
.n-kpi-icon.is-blue { border-color: rgba(76, 111, 255, 0.36); background: rgba(76, 111, 255, 0.16); color: #8fa8ff; }
.n-kpi-icon.is-green { border-color: rgba(46, 212, 122, 0.34); background: rgba(46, 212, 122, 0.14); color: #5ce49b; }
.n-kpi-icon.is-violet { border-color: rgba(167, 139, 250, 0.36); background: rgba(167, 139, 250, 0.16); color: #c0a9ff; }
.n-kpi-icon.is-amber { border-color: rgba(245, 185, 61, 0.34); background: rgba(245, 185, 61, 0.14); color: #f7ca6b; }
.n-kpi-value { display: flex; align-items: baseline; gap: 6px; font-size: 25px; font-weight: 620; line-height: 1.15; letter-spacing: -0.015em; font-variant-numeric: tabular-nums; color: #fff; }
.n-kpi-unit { font-size: 12px; font-weight: 500; color: var(--n-text-3); }
.n-kpi-foot { margin-top: 9px; font-size: 12px; color: var(--n-text-4); }

/* ---------- 工具条 ---------- */
.n-toolbar { display: flex; align-items: center; gap: 8px; flex-wrap: wrap; padding: 12px 14px; border: 1px solid var(--n-line); border-radius: var(--n-radius); background: var(--n-card); }
.n-spacer { flex: 1; }
.n-hint { font-size: 12px; }
.n-seg { display: inline-flex; gap: 2px; padding: 3px; border: 1px solid #24272e; border-radius: 9px; background: #14161a; }
.n-seg button { height: 26px; padding: 0 11px; border: 0; border-radius: 7px; background: transparent; color: var(--n-text-3); font-size: 12px; cursor: pointer; }
.n-seg button:hover { color: var(--n-text-1); }
.n-seg button.is-active { background: #2b3038; color: #fff; }

/* ---------- 卡片网格 ---------- */
.n-grid { display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 12px; }
.n-card { position: relative; border: 1px solid var(--n-line); border-radius: var(--n-radius); background: var(--n-card); min-width: 0; }
.n-card-head { display: flex; align-items: center; gap: 10px; padding: 14px 16px; }
.n-card-head h3 { display: flex; align-items: center; gap: 8px; margin: 0; min-width: 0; font-size: 13.5px; font-weight: 600; color: var(--n-text-1); }
.n-right { display: flex; align-items: center; gap: 6px; margin-left: auto; flex: none; }
.n-region { font-size: 12px; }
.n-dot { width: 7px; height: 7px; flex: none; border-radius: 50%; background: #545b64; }
.n-dot.is-ok { background: var(--n-ok); box-shadow: 0 0 7px rgba(46, 212, 122, 0.5); }
.n-dot.is-warn { background: var(--n-warn); box-shadow: 0 0 7px rgba(245, 185, 61, 0.45); }
.n-dot.is-danger { background: var(--n-danger); box-shadow: 0 0 7px rgba(243, 101, 111, 0.45); }
.n-dot.is-neutral { background: #545b64; }

.n-node-card.is-offline { border-color: rgba(243, 101, 111, 0.42); }
.n-node-card.is-disabled { opacity: 0.82; }

.n-mana { display: flex; align-items: center; gap: 7px; margin: 0 16px; padding: 9px 10px; border: 1px solid #24272e; border-radius: 9px; background: #101114; font-size: 12.5px; color: #aab0ba; overflow: hidden; transition: border-color 0.16s ease, background 0.16s ease; }
.n-mana:hover { border-color: rgba(76, 111, 255, 0.42); background: rgba(76, 111, 255, 0.08); }
.n-mana.is-empty { color: var(--n-text-4); }
.n-mana-url { flex: 1; min-width: 0; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace; font-size: 12px; color: #a8b8ff; }
.n-mana.is-empty .n-mana-url { color: var(--n-text-4); font-family: inherit; }
.n-mana-ico { width: 13px; height: 13px; color: #98a0ab; flex: none; }
.n-mana:hover .n-mana-ico { color: #8fa8ff; }

.n-stats { display: grid; grid-template-columns: repeat(4, 1fr); margin: 12px 16px 0; padding: 11px 0; border-top: 1px solid var(--n-line-soft); border-bottom: 1px solid var(--n-line-soft); }
.n-stats > div { min-width: 0; padding: 0 10px; border-left: 1px solid var(--n-line-soft); }
.n-stats > div:first-child { padding-left: 0; border-left: 0; }
.n-stats span { display: block; margin-bottom: 3px; font-size: 11.5px; color: #98a0ab; white-space: nowrap; }
.n-stats b { display: block; font-size: 14px; font-weight: 600; color: var(--n-text-1); font-variant-numeric: tabular-nums; white-space: nowrap; }
.n-stats b small { margin-left: 2px; font-size: 11.5px; font-weight: 400; color: #98a0ab; }
.n-stats i.n-mini { display: block; margin-top: 2px; font-size: 11px; font-style: normal; color: #98a0ab; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }

.n-card-foot { display: flex; align-items: center; gap: 8px; padding: 12px 16px 14px; }
.n-foot-meta { display: flex; flex-direction: column; gap: 2px; min-width: 0; font-size: 11.5px; color: #98a0ab; line-height: 1.5; }
.n-foot-meta b { display: flex; align-items: center; gap: 6px; color: #bcc2cc; font-size: 11.5px; font-weight: 500; }
.n-foot-meta span { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.n-foot-actions { display: flex; align-items: center; gap: 6px; margin-left: auto; flex: none; }

.n-add-card { display: grid; place-items: center; min-height: 196px; border-style: dashed; background: rgba(255, 255, 255, 0.012); cursor: pointer; }
.n-add-card:hover { border-color: rgba(76, 111, 255, 0.4); }
.n-add-inner { display: flex; flex-direction: column; align-items: center; gap: 7px; padding: 18px; text-align: center; }
.n-add-inner .n-kpi-icon { margin: 0 0 3px; }
.n-add-inner b { font-size: 13px; }
.n-add-inner span { font-size: 12px; line-height: 1.7; }

/* ---------- 加载 / 空 / 错误 / 分页 ---------- */
.n-loading, .n-empty { display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 10px; min-height: 220px; padding: 24px; border: 1px solid var(--n-line); border-radius: var(--n-radius); background: var(--n-card); color: var(--n-text-4); font-size: 12.5px; }
.n-empty-title { font-size: 14px; color: var(--n-text-1); }
.n-empty-desc { margin: 0; color: var(--n-text-4); }
.n-alert { padding: 10px 13px; border: 1px solid rgba(243, 101, 111, 0.32); border-radius: 10px; background: rgba(243, 101, 111, 0.1); color: #ffb3b9; font-size: 12.5px; }
.n-pager { display: flex; align-items: center; justify-content: space-between; gap: 10px; font-size: 12px; }

/* ---------- 徽标 ---------- */
.n-badge { display: inline-flex; align-items: center; gap: 5px; height: 21px; padding: 0 7px; border: 1px solid transparent; border-radius: 6px; font-size: 11.5px; font-weight: 550; white-space: nowrap; }
.n-badge.is-neutral { border-color: #2a2e35; background: #1d2026; color: #bcc2cc; }
.n-badge.is-brand { border-color: rgba(76, 111, 255, 0.34); background: rgba(76, 111, 255, 0.16); color: #a8b8ff; }
.n-badge.is-danger { border-color: rgba(243, 101, 111, 0.32); background: rgba(243, 101, 111, 0.15); color: #ff9ba2; }
.n-badge.is-warn { border-color: rgba(245, 185, 61, 0.32); background: rgba(245, 185, 61, 0.15); color: #f7ca6b; }
.n-badge-row { display: flex; gap: 6px; flex-wrap: wrap; }

/* ---------- 弹窗 ---------- */
.n-mask {
  position: fixed;
  inset: 0;
  z-index: 2000;
  display: grid;
  place-items: center;
  padding: 24px;
  background: rgba(0, 0, 0, 0.66);
}
.n-modal { display: flex; flex-direction: column; max-height: min(86vh, 880px); border: 1px solid var(--n-line); border-radius: 16px; background: var(--n-card); color: var(--n-text-1); text-align: left; box-shadow: 0 30px 80px -24px rgba(0, 0, 0, 0.9); }
.n-modal-head { display: flex; align-items: flex-start; gap: 10px; padding: 16px 20px 14px; border-bottom: 1px solid var(--n-line-soft); }
.n-modal-head-main { flex: 1; min-width: 0; }
.n-modal-head h3 { margin: 0; font-size: 15px; font-weight: 600; }
.n-modal-sub { margin: 3px 0 0; font-size: 12px; color: var(--n-text-4); }
.n-modal-body { flex: 1; min-height: 0; overflow-y: auto; padding: 16px 20px 20px; }
.n-modal-foot { display: flex; align-items: center; gap: 10px; padding: 13px 20px; border-top: 1px solid var(--n-line-soft); font-size: 12px; }
.n-modal-foot .n-foot-actions { margin-left: auto; }

/* ---------- 步骤 ---------- */
.n-steps { display: flex; align-items: stretch; gap: 10px; margin-bottom: 16px; }
.n-step { flex: 1; display: flex; align-items: center; gap: 9px; padding: 10px 12px; border: 1px solid #24272e; border-radius: 10px; background: #101114; min-width: 0; }
.n-step i { display: grid; place-items: center; width: 20px; height: 20px; flex: none; border-radius: 50%; border: 1px solid #2a2e35; background: #1d2026; font-size: 11.5px; font-style: normal; font-weight: 650; color: #9aa1ad; }
.n-step b { display: block; font-size: 12.5px; font-weight: 550; color: #e6e8ec; }
.n-step span { display: block; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; font-size: 11.5px; color: #98a0ab; }
.n-step.is-done i { border-color: rgba(46, 212, 122, 0.34); background: rgba(46, 212, 122, 0.16); color: #5ce49b; }
.n-step.is-active { border-color: rgba(76, 111, 255, 0.42); background: rgba(76, 111, 255, 0.08); }
.n-step.is-active i { border-color: rgba(76, 111, 255, 0.45); background: rgba(76, 111, 255, 0.2); color: #a8b8ff; }

/* ---------- 提示条 ---------- */
.n-callout { display: flex; gap: 10px; padding: 12px 13px; border: 1px solid var(--n-line); border-radius: 11px; background: rgba(255, 255, 255, 0.022); font-size: 12px; line-height: 1.7; color: var(--n-text-2); }
.n-callout-ico { width: 16px; height: 16px; flex: none; margin-top: 2px; color: var(--n-text-3); }
.n-callout b { display: block; margin-bottom: 1px; font-size: 12.5px; font-weight: 600; color: var(--n-text-1); }
.n-callout.is-info { border-color: rgba(90, 169, 240, 0.24); background: rgba(90, 169, 240, 0.08); }
.n-callout.is-info .n-callout-ico { color: #8cc4f7; }
.n-callout.is-warn { border-color: rgba(245, 185, 61, 0.28); background: rgba(245, 185, 61, 0.08); }
.n-callout.is-warn .n-callout-ico { color: #f7ca6b; }
.n-callout.is-danger { border-color: rgba(243, 101, 111, 0.28); background: rgba(243, 101, 111, 0.08); }

/* ---------- 表单 ---------- */
.n-form-grid { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 12px; }
.n-field { display: flex; flex-direction: column; gap: 5px; min-width: 0; }
.n-field-wide { grid-column: 1 / -1; }
.n-label { font-size: 12px; color: var(--n-text-2); }
.n-req { color: var(--n-danger); }
.n-help { font-size: 11px; line-height: 1.5; color: var(--n-text-4); }
.n-exit-ip-list { display: flex; flex-direction: column; gap: 8px; }
.n-exit-ip-row { display: flex; align-items: center; gap: 8px; }
.n-exit-ip-row :deep(.n-input) { flex: 1; }
.n-code-title { margin-bottom: 7px; font-size: 12px; color: var(--n-text-2); }
.n-code { margin: 0; padding: 11px 13px; border: 1px solid #23262c; border-radius: 9px; background: #0d0e11; color: #b8bec8; font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace; font-size: 11.5px; line-height: 1.7; white-space: pre-wrap; word-break: break-word; }
.n-code-actions { display: flex; align-items: center; gap: 8px; margin-top: 8px; font-size: 11.5px; }

.mt-8 { margin-top: 8px; }
.mt-12 { margin-top: 12px; }
.mt-16 { margin-top: 16px; }
.mb-12 { margin-bottom: 12px; }

@media (max-width: 1280px) {
  .n-grid { grid-template-columns: repeat(2, minmax(0, 1fr)); }
}
@media (max-width: 1100px) {
  .n-kpis { grid-template-columns: repeat(2, minmax(0, 1fr)); }
}
@media (max-width: 820px) {
  .n-grid { grid-template-columns: minmax(0, 1fr); }
  .n-form-grid { grid-template-columns: minmax(0, 1fr); }
  .n-search { width: 100%; }
}
</style>
