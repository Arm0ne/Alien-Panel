<script setup lang="ts">
defineOptions({ name: 'ProjectSyncStatusTag' });

type SyncStatus = Api.Central.NodeSyncRunStatus | 'unknown' | string;

interface Props {
  status?: SyncStatus | null;
}

const props = withDefaults(defineProps<Props>(), { status: 'unknown' });

const statusMap: Record<string, { label: string; type: 'default' | 'info' | 'warning' | 'success' | 'error' }> = {
  queued: { label: '排队中', type: 'info' },
  running: { label: '同步中', type: 'warning' },
  success: { label: '同步成功', type: 'success' },
  failed: { label: '同步失败', type: 'error' },
  unknown: { label: '状态未知', type: 'default' }
};
</script>

<template>
  <NTag size="small" :type="(statusMap[props.status || 'unknown'] || statusMap.unknown).type">
    {{ (statusMap[props.status || 'unknown'] || statusMap.unknown).label }}
  </NTag>
</template>
