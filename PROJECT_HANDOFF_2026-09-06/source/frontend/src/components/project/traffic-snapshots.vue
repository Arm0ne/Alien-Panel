<script setup lang="ts">
import TrafficValue from './traffic-value.vue';

defineOptions({ name: 'ProjectTrafficSnapshots' });

interface Snapshot {
  collectedAt: string;
  allTime: number;
  resetDetected: boolean;
}

interface Props {
  snapshots: Snapshot[];
  emptyDescription?: string;
}

withDefaults(defineProps<Props>(), { emptyDescription: '暂无流量快照' });

function formatDate(value?: string | null) {
  if (!value) return '--';
  const date = new Date(value);
  return Number.isNaN(date.getTime()) ? '时间未知' : date.toLocaleString('zh-CN', { hour12: false });
}
</script>

<template>
  <NEmpty v-if="snapshots.length === 0" :description="emptyDescription" size="small" />
  <div v-else class="space-y-6px">
    <div v-for="snapshot in snapshots" :key="snapshot.collectedAt" class="flex items-center justify-between gap-8px text-12px">
      <span class="text-gray-500">{{ formatDate(snapshot.collectedAt) }}</span>
      <TrafficValue :value="snapshot.allTime" />
      <NTag v-if="snapshot.resetDetected" size="small" type="warning">基线重置</NTag>
    </div>
  </div>
</template>
