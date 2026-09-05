<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue';
import { useRouterPush } from '@/hooks/common/router';
import { fetchEventSummary, fetchEvents } from '@/service/api';
import { sessionStg } from '@/utils/storage';

defineOptions({ name: 'EventBell' });

const { routerPushByKey } = useRouterPush();
const pendingCount = ref(0);
const unreadCount = ref(0);
const pendingEvents = ref<Api.Central.EventSummary[]>([]);
const open = ref(false);
let pollTimer: ReturnType<typeof setInterval> | undefined;
const knownIds = new Set<string>();
const alertedIds = new Set(sessionStg.get('eventAlertedIds') || []);

function formatDate(value: string) {
  const date = new Date(value);
  return Number.isNaN(date.getTime()) ? value : date.toLocaleString('zh-CN', { hour12: false });
}

async function readPendingEvents() {
  const { data } = await fetchEvents({ status: 'pending', page: 1, page_size: 5 });
  if (data) {
    pendingEvents.value = data.items;
    data.items.forEach(item => knownIds.add(item.id));
  }
}

function alertPendingEvents(items: Api.Central.EventSummary[]) {
  const newItems = items.filter(item => !knownIds.has(item.id) && !alertedIds.has(item.id));
  const shouldAlert = newItems.length > 0;
  items.forEach(item => knownIds.add(item.id));
  if (!shouldAlert) return;
  const count = newItems.length;
  newItems.forEach(item => alertedIds.add(item.id));
  sessionStg.set('eventAlertedIds', [...alertedIds].slice(-100));
  window.$dialog?.warning({
    title: '有待处理事件',
    content: `检测到 ${count} 条需要处理的事件${count > 0 ? '，请及时确认' : ''}`,
    positiveText: '立即处理',
    negativeText: '稍后处理',
    maskClosable: false,
    onPositiveClick: () => routerPushByKey('events')
  });
}

async function refresh() {
  const { data } = await fetchEventSummary();
  if (data) {
    pendingCount.value = data.pendingCount;
    unreadCount.value = data.unreadCount;
  }
  const { data: events } = await fetchEvents({ status: 'pending', page: 1, page_size: 5 });
  if (events) {
    pendingEvents.value = events.items;
    alertPendingEvents(events.items);
  }
}

async function toggle() {
  open.value = !open.value;
  if (open.value) await readPendingEvents();
}

function openEvents() {
  open.value = false;
  routerPushByKey('events');
}

function handleEventChanged() {
  void refresh();
}

onMounted(() => {
  window.addEventListener('xpanel-events-updated', handleEventChanged);
  refresh();
  pollTimer = setInterval(refresh, 30_000);
});

onUnmounted(() => {
  window.removeEventListener('xpanel-events-updated', handleEventChanged);
  if (pollTimer) clearInterval(pollTimer);
});
</script>

<template>
  <NPopover v-model:show="open" trigger="manual" placement="bottom-end" :style="{ width: '360px' }">
    <template #trigger>
      <NBadge :value="pendingCount" :max="99" :show-zero="false" :dot="pendingCount > 0" type="error" processing>
        <NButton quaternary circle aria-label="事件中心" @click="toggle">
          <template #icon><icon-mdi-bell-outline /></template>
        </NButton>
      </NBadge>
    </template>
    <div class="p-4px">
      <div class="mb-10px flex items-center justify-between">
        <span class="text-16px font-600">待处理事件</span>
        <span class="text-12px text-gray-500">未读 {{ unreadCount }}</span>
      </div>
      <NEmpty v-if="pendingEvents.length === 0" description="暂无待处理事件" />
      <NList v-else hoverable>
        <NListItem v-for="event in pendingEvents" :key="event.id" class="cursor-pointer" @click="openEvents">
          <div class="min-w-0">
            <div class="truncate text-14px font-500">{{ event.title || event.type }}</div>
            <div class="truncate text-12px text-gray-500">{{ event.nodeName || '系统' }} · {{ formatDate(event.occurredAt) }}</div>
          </div>
        </NListItem>
      </NList>
      <NButton block class="mt-10px" @click="openEvents">查看全部事件</NButton>
    </div>
  </NPopover>
</template>
