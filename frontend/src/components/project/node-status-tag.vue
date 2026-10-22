<script setup lang="ts">
import { computed } from 'vue';

defineOptions({ name: 'ProjectNodeStatusTag' });

interface Props {
  status?: Api.Central.NodeStatus | string | null;
}

const props = withDefaults(defineProps<Props>(), { status: 'unknown' });

const statusMap: Record<string, { label: string; tone: string }> = {
  online: { label: '在线', tone: 'is-ok' },
  degraded: { label: '异常', tone: 'is-warn' },
  offline: { label: '离线', tone: 'is-danger' },
  disabled: { label: '已停用', tone: 'is-neutral' },
  unknown: { label: '状态未知', tone: 'is-neutral' }
};

const current = computed(() => statusMap[props.status || 'unknown'] || statusMap.unknown);
</script>

<template>
  <span class="v2-tag" :class="current.tone"><i class="v2-dot"></i>{{ current.label }}</span>
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

.v2-dot {
  width: 5px;
  height: 5px;
  flex: none;
  border-radius: 50%;
  background: currentColor;
}

.v2-tag.is-ok {
  border-color: rgba(46, 212, 122, 0.32);
  background: rgba(46, 212, 122, 0.15);
  color: #5ce49b;
}

.v2-tag.is-warn {
  border-color: rgba(245, 185, 61, 0.32);
  background: rgba(245, 185, 61, 0.15);
  color: #f7ca6b;
}

.v2-tag.is-danger {
  border-color: rgba(243, 101, 111, 0.32);
  background: rgba(243, 101, 111, 0.15);
  color: #ff9ba2;
}

.v2-tag.is-neutral {
  border-color: #2a2e35;
  background: #1d2026;
  color: #bcc2cc;
}
</style>
