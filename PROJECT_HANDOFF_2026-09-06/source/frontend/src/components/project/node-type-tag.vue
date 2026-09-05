<script setup lang="ts">
import { computed } from 'vue';

defineOptions({ name: 'ProjectNodeTypeTag' });

interface Props {
  type?: Api.Central.NodeType | string | null;
}

const props = withDefaults(defineProps<Props>(), { type: 'unknown' });

const typeMap: Record<string, { label: string; type: 'info' | 'success' | 'default' }> = {
  relay: { label: '线路机', type: 'info' },
  landing: { label: '落地机', type: 'success' },
  unknown: { label: '未知类型', type: 'default' }
};

const current = computed(() => typeMap[props.type || 'unknown'] || typeMap.unknown);
</script>

<template>
  <NTag size="small" :type="current.type">{{ current.label }}</NTag>
</template>
