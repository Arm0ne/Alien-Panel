<script setup lang="ts">
import { computed } from 'vue';

defineOptions({ name: 'ProjectEnabledStatusTag' });

interface Props {
  enabled?: boolean | null;
}

const props = withDefaults(defineProps<Props>(), { enabled: null });
const current = computed(() => {
  if (props.enabled === true) return { label: '启用', type: 'success' as const };
  if (props.enabled === false) return { label: '停用', type: 'default' as const };
  return { label: '状态未知', type: 'default' as const };
});
</script>

<template>
  <NTag size="small" :type="current.type">{{ current.label }}</NTag>
</template>
