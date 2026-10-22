<script setup lang="ts">
import { computed } from 'vue';

defineOptions({ name: 'ProjectSyncStatusTag' });

type SyncStatus = Api.Central.NodeSyncRunStatus | 'unknown' | string;

interface Props {
  status?: SyncStatus | null;
}

const props = withDefaults(defineProps<Props>(), { status: 'unknown' });

const statusMap: Record<string, { label: string; tone: string }> = {
  queued: { label: '排队中', tone: 'is-info' },
  running: { label: '同步中', tone: 'is-warn' },
  success: { label: '同步成功', tone: 'is-ok' },
  failed: { label: '同步失败', tone: 'is-danger' },
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

.v2-tag.is-info {
  border-color: rgba(90, 169, 240, 0.32);
  background: rgba(90, 169, 240, 0.15);
  color: #8cc4f7;
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
