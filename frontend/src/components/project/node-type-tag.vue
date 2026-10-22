<script setup lang="ts">
import { computed } from 'vue';

defineOptions({ name: 'ProjectNodeTypeTag' });

interface Props {
  type?: Api.Central.NodeType | string | null;
}

const props = withDefaults(defineProps<Props>(), { type: 'unknown' });

const typeMap: Record<string, { label: string; tone: string }> = {
  relay: { label: '线路机', tone: 'is-brand' },
  landing: { label: '落地机', tone: 'is-neutral' },
  unknown: { label: '未知类型', tone: 'is-neutral' }
};

const current = computed(() => typeMap[props.type || 'unknown'] || typeMap.unknown);
</script>

<template>
  <span class="v2-tag" :class="current.tone">{{ current.label }}</span>
</template>

<style scoped>
.v2-tag {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  height: 21px;
  padding: 0 7px;
  border: 1px solid transparent;
  border-radius: 6px;
  font-size: 11.5px;
  font-weight: 550;
  line-height: 1;
  white-space: nowrap;
}

.v2-tag.is-brand {
  border-color: rgba(76, 111, 255, 0.34);
  background: rgba(76, 111, 255, 0.16);
  color: #a8b8ff;
}

.v2-tag.is-neutral {
  border-color: #2a2e35;
  background: #1d2026;
  color: #bcc2cc;
}
</style>
