<script setup lang="ts">
import { computed } from 'vue';

defineOptions({ name: 'ProjectTrafficValue' });

interface Props {
  value?: number | null;
  emptyLabel?: string;
}

const props = withDefaults(defineProps<Props>(), { value: null, emptyLabel: '--' });

const validValue = computed(() => {
  if (props.value === null || props.value === undefined || !Number.isFinite(props.value) || props.value < 0) return null;
  return props.value;
});

const displayValue = computed(() => {
  const value = validValue.value;
  if (value === null) return props.emptyLabel;
  if (value < 1024) return `${value} B`;

  const units = ['KB', 'MB', 'GB', 'TB', 'PB'];
  let size = value;
  let index = -1;
  do {
    size /= 1024;
    index += 1;
  } while (size >= 1024 && index < units.length - 1);
  return `${size.toFixed(size >= 10 ? 1 : 2)} ${units[index]}`;
});

const exactValue = computed(() => (validValue.value === null ? '' : `${validValue.value.toLocaleString('zh-CN')} B`));
</script>

<template>
  <span :title="exactValue">{{ displayValue }}</span>
</template>
