<script setup lang="ts">
import { h } from 'vue';
import type { DataTableColumns } from 'naive-ui';
import EnabledStatusTag from './enabled-status-tag.vue';
import TrafficValue from './traffic-value.vue';

defineOptions({ name: 'ProjectDeviceTable' });

interface Props {
  devices: Api.Central.UserClientDetail[];
  emptyDescription?: string;
}

withDefaults(defineProps<Props>(), { emptyDescription: '暂无设备凭证' });

function formatDate(value?: string | null) {
  if (!value) return '--';
  const date = new Date(value);
  return Number.isNaN(date.getTime()) ? '时间未知' : date.toLocaleString('zh-CN', { hour12: false });
}

const columns: DataTableColumns<Api.Central.UserClientDetail> = [
  {
    title: '设备 / Email',
    key: 'email',
    minWidth: 180,
    render: row =>
      h('div', [
        h('div', { class: 'font-medium' }, row.email || '未填写 Email'),
        h('div', { class: 'text-12px text-gray-500' }, row.remoteId)
      ])
  },
  { title: '状态', key: 'enabled', width: 90, render: row => h(EnabledStatusTag, { enabled: row.enabled }) },
  { title: '累计流量', key: 'allTime', minWidth: 120, render: row => h(TrafficValue, { value: row.allTime }) },
  { title: '最近在线', key: 'lastOnlineAt', minWidth: 170, render: row => formatDate(row.lastOnlineAt) }
];
</script>

<template>
  <NDataTable
    v-if="devices.length"
    :columns="columns"
    :data="devices"
    :bordered="false"
    :single-line="false"
    size="small"
    :scroll-x="560"
  />
  <NEmpty v-else :description="emptyDescription" size="small" />
</template>
