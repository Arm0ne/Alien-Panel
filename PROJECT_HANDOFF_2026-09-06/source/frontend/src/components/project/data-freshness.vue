<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue';

defineOptions({ name: 'ProjectDataFreshness' });

interface Props {
  /** Timestamp of the newest snapshot used by the API response. */
  dataAt?: string | null;
  /** Mark data as delayed after this many minutes. */
  delayedAfterMinutes?: number;
  /** Mark data as expired after this many minutes. */
  expiredAfterMinutes?: number;
  /** Render only the state tag, without the timestamp detail. */
  compact?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  dataAt: null,
  delayedAfterMinutes: 10,
  expiredAfterMinutes: 30,
  compact: false
});

const now = ref(Date.now());
let timer: ReturnType<typeof setInterval> | undefined;

const parsedAt = computed(() => {
  if (!props.dataAt) return null;
  const timestamp = new Date(props.dataAt).getTime();
  return Number.isFinite(timestamp) ? timestamp : null;
});

const state = computed(() => {
  if (parsedAt.value === null) {
    return { label: '数据时间未知', type: 'warning' as const, detail: '中央接口没有提供有效的同步时间' };
  }

  const age = now.value - parsedAt.value;
  if (age < -5 * 60 * 1000) {
    return { label: '数据时间未知', type: 'warning' as const, detail: '同步时间晚于当前时间，可能存在时钟偏差' };
  }

  const delayedAfter = Math.max(props.delayedAfterMinutes, 1) * 60 * 1000;
  const expiredAfter = Math.max(props.expiredAfterMinutes, props.delayedAfterMinutes) * 60 * 1000;
  if (age > expiredAfter) {
    return { label: '数据过期', type: 'error' as const, detail: `已超过 ${Math.round(expiredAfter / 60000)} 分钟未同步` };
  }
  if (age > delayedAfter) {
    return { label: '数据延迟', type: 'warning' as const, detail: `已超过 ${Math.round(delayedAfter / 60000)} 分钟未同步` };
  }
  return { label: '数据新鲜', type: 'success' as const, detail: '最近一次同步在可接受范围内' };
});

const formattedAt = computed(() => {
  if (parsedAt.value === null) return '';
  return new Date(parsedAt.value).toLocaleString('zh-CN', { hour12: false });
});

onMounted(() => {
  if (!props.compact) {
    timer = setInterval(() => {
      now.value = Date.now();
    }, 60_000);
  }
});

onBeforeUnmount(() => {
  if (timer) clearInterval(timer);
});
</script>

<template>
  <NTooltip>
    <template #trigger>
      <NTag size="small" :type="state.type">
        {{ state.label }}<span v-if="!compact && formattedAt" class="ml-4px">· {{ formattedAt }}</span>
      </NTag>
    </template>
    <div>{{ state.detail }}</div>
    <div v-if="formattedAt" class="mt-2px text-12px opacity-80">数据时间：{{ formattedAt }}</div>
  </NTooltip>
</template>
