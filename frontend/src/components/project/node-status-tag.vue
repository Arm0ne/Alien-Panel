<script setup lang="ts">
import { computed } from 'vue';

defineOptions({ name: 'ProjectNodeStatusTag' });

interface Props {
  status?: Api.Central.NodeStatus | string | null;
}

const props = withDefaults(defineProps<Props>(), { status: 'unknown' });
const statusMap: Record<string, { label: string; type: 'success' | 'warning' | 'error' | 'default' }> = {
  online: { label: '在线', type: 'success' },
  degraded: { label: '异常', type: 'warning' },
  offline: { label: '离线', type: 'error' },
  disabled: { label: '已停用', type: 'default' },
  unknown: { label: '状态未知', type: 'default' }
};
const current = computed(() => statusMap[props.status || 'unknown'] || statusMap.unknown);
</script>

<template>
  <NTag size="small" :type="current.type">{{ current.label }}</NTag>
</template>
