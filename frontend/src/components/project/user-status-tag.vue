<script setup lang="ts">
import { computed } from 'vue';

defineOptions({ name: 'ProjectUserStatusTag' });

interface Props {
  status?: Api.Central.UserStatus | string | null;
  expiresAt?: string | null;
}

const props = withDefaults(defineProps<Props>(), { status: 'unknown', expiresAt: null });

const statusValue = computed(() => {
  if (props.expiresAt) {
    const timestamp = new Date(props.expiresAt).getTime();
    if (!Number.isFinite(timestamp)) return 'unknown-expiry';
  }
  return props.status || 'unknown';
});

const statusMap: Record<string, { label: string; type: 'default' | 'warning' | 'success' | 'error' }> = {
  active: { label: '有效', type: 'success' },
  expiring: { label: '即将到期', type: 'warning' },
  expired: { label: '已到期', type: 'error' },
  disabled: { label: '已停用', type: 'default' },
  unknown: { label: '状态未知', type: 'default' },
  'unknown-expiry': { label: '到期时间未知', type: 'warning' }
};
</script>

<template>
  <NTag size="small" :type="(statusMap[statusValue] || statusMap.unknown).type">
    {{ (statusMap[statusValue] || statusMap.unknown).label }}
  </NTag>
</template>
