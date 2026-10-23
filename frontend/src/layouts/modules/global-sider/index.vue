<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue';
import { useRoute } from 'vue-router';
import { fetchEventSummary, fetchNodes } from '@/service/api';
import { useAppStore } from '@/store/modules/app';

defineOptions({
  name: 'GlobalSider'
});

const route = useRoute();
const appStore = useAppStore();

interface NavItem {
  key: string;
  label: string;
  icon: string;
  path: string;
  badge?: number;
}

const pendingEventCount = ref(0);
const nodesOnline = ref(0);
const nodesTotal = ref(0);
const lastSyncAt = ref<string | null>(null);

const groups = computed<Array<{ title: string; items: NavItem[] }>>(() => [
  {
    title: '运营',
    items: [
      { key: 'dashboard', label: '总览', icon: 'mdi:monitor-dashboard', path: '/dashboard' },
      {
        key: 'events',
        label: '事件中心',
        icon: 'mdi:bell-alert-outline',
        path: '/events',
        badge: pendingEventCount.value
      }
    ]
  },
  {
    title: '业务',
    items: [
      { key: 'users', label: '用户管理', icon: 'mdi:account-group-outline', path: '/users' },
      { key: 'nodes', label: '节点管理', icon: 'mdi:server-network', path: '/nodes' },
      { key: 'exit-ips', label: '出口 IP', icon: 'mdi:earth', path: '/exit-ips' }
    ]
  },
  {
    title: '财务',
    items: [{ key: 'finance', label: '成本与收入', icon: 'mdi:finance', path: '/finance' }]
  },
  {
    title: '系统',
    items: [{ key: 'system_backup', label: '数据备份', icon: 'mdi:database-outline', path: '/system/backup' }]
  }
]);

const lastSyncText = computed(() => {
  if (!lastSyncAt.value) return '--';
  const timestamp = new Date(lastSyncAt.value).getTime();
  if (!Number.isFinite(timestamp)) return '--';
  const minutes = Math.max(0, Math.round((Date.now() - timestamp) / 60_000));
  if (minutes < 1) return '刚刚';
  if (minutes < 60) return `${minutes} 分钟前`;
  const hours = Math.round(minutes / 60);
  if (hours < 24) return `${hours} 小时前`;
  return new Date(lastSyncAt.value).toLocaleString('zh-CN', { hour12: false });
});

function isActive(path: string) {
  if (path === '/dashboard') return route.path === '/' || route.path.startsWith('/dashboard');
  if (path === '/nodes') return route.path.startsWith('/nodes');
  return route.path === path || route.path.startsWith(`${path}/`);
}

async function loadSiderStatus() {
  const [{ data: eventData }, { data: nodeData }] = await Promise.all([
    fetchEventSummary(),
    fetchNodes({ page: 1, page_size: 100 })
  ]);

  if (eventData) {
    pendingEventCount.value = eventData.pendingCount;
  }

  if (nodeData) {
    nodesOnline.value = nodeData.stats?.online ?? 0;
    nodesTotal.value = nodeData.stats?.total ?? 0;
    const times = (nodeData.items || [])
      .map(item => item.lastSyncAt)
      .filter((value): value is string => Boolean(value))
      .map(value => new Date(value).getTime())
      .filter(value => Number.isFinite(value));
    lastSyncAt.value = times.length ? new Date(Math.max(...times)).toISOString() : null;
  }
}

let pollTimer: ReturnType<typeof setInterval> | undefined;

onMounted(() => {
  loadSiderStatus();
  window.addEventListener('xpanel-events-updated', loadSiderStatus);
  pollTimer = setInterval(loadSiderStatus, 60_000);
});

onUnmounted(() => {
  window.removeEventListener('xpanel-events-updated', loadSiderStatus);
  if (pollTimer) clearInterval(pollTimer);
});
</script>

<template>
  <div class="v2-sider size-full flex-col-stretch">
    <div class="v2-sider-head">
      <div class="v2-brand-mark">
        <svg viewBox="0 0 28 28" fill="none" aria-hidden="true">
          <path d="M14 3.4 25.4 24.6H19.4L14 13.6 8.6 24.6H2.6Z" fill="currentColor" />
        </svg>
      </div>
      <div v-show="!appStore.siderCollapse" class="v2-sider-brand">
        <strong>Alien Panel</strong>
        <span>中央运营控制台 · v2.0.7</span>
      </div>
    </div>

    <nav class="v2-sider-nav">
      <div v-for="group in groups" :key="group.title" class="v2-nav-group">
        <div v-show="!appStore.siderCollapse" class="v2-nav-group-title">{{ group.title }}</div>
        <RouterLink
          v-for="item in group.items"
          :key="item.key"
          :to="item.path"
          class="v2-nav-item"
          :class="{ 'is-active': isActive(item.path), 'is-collapsed': appStore.siderCollapse }"
          :title="item.label"
        >
          <SvgIcon :icon="item.icon" class="v2-nav-icon" />
          <span v-show="!appStore.siderCollapse" class="v2-nav-label">{{ item.label }}</span>
          <i
            v-if="item.badge && !appStore.siderCollapse"
            class="v2-nav-badge"
            :class="{ 'is-quiet': item.key === 'exit-ips' }"
          >
            {{ item.badge }}
          </i>
        </RouterLink>
      </div>
    </nav>

    <div v-show="!appStore.siderCollapse" class="v2-sider-status">
      <div class="v2-sider-status-top">
        <span>系统状态</span>
        <i class="v2-status-dot" :class="{ 'is-off': nodesTotal > 0 && nodesOnline < nodesTotal }"></i>
      </div>
      <div class="v2-sider-status-row">
        <span>在线节点</span>
        <b>{{ nodesOnline }} / {{ nodesTotal }}</b>
      </div>
      <div class="v2-sider-status-row">
        <span>最近同步</span>
        <b>{{ lastSyncText }}</b>
      </div>
      <div class="v2-sider-status-row">
        <span>待处理事件</span>
        <b :class="{ 'is-danger': pendingEventCount > 0 }">{{ pendingEventCount }}</b>
      </div>
    </div>
  </div>
</template>

<style scoped></style>
