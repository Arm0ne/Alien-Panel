<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue';
import { useRoute } from 'vue-router';
import { fetchEventSummary } from '@/service/api';
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

function isActive(path: string) {
  if (path === '/dashboard') return route.path === '/' || route.path.startsWith('/dashboard');
  if (path === '/nodes') return route.path.startsWith('/nodes');
  return route.path === path || route.path.startsWith(`${path}/`);
}

async function loadEventCount() {
  const { data } = await fetchEventSummary();
  if (data) {
    pendingEventCount.value = data.pendingCount;
  }
}

let pollTimer: ReturnType<typeof setInterval> | undefined;

onMounted(() => {
  loadEventCount();
  window.addEventListener('xpanel-events-updated', loadEventCount);
  pollTimer = setInterval(loadEventCount, 60_000);
});

onUnmounted(() => {
  window.removeEventListener('xpanel-events-updated', loadEventCount);
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
        <span>中央运营控制台 · v1.0.0</span>
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
  </div>
</template>

<style scoped></style>
