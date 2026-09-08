<script setup lang="ts">
import DataFreshness from './data-freshness.vue';

defineOptions({ name: 'ProjectModulePage' });

interface Props {
  title: string;
  description?: string;
  loading?: boolean;
  error?: string;
  empty?: boolean;
  emptyDescription?: string;
  dataAt?: string | null;
}

withDefaults(defineProps<Props>(), {
  description: '',
  loading: false,
  error: '',
  empty: true,
  emptyDescription: '暂无数据',
  dataAt: null
});

const emit = defineEmits<{
  (event: 'refresh'): void;
}>();
</script>

<template>
  <NSpace vertical :size="16">
    <NCard :bordered="false" size="small">
      <div class="flex flex-wrap items-center justify-between gap-16px">
        <div>
          <h2 class="m-0 text-20px font-600">{{ title }}</h2>
          <p v-if="description" class="mb-0 mt-8px text-14px text-gray-500">{{ description }}</p>
        </div>
        <NSpace align="center" :size="8">
          <DataFreshness :data-at="dataAt" />
          <NButton size="small" :loading="loading" @click="emit('refresh')">
            <template #icon><icon-mdi-refresh /></template>
            刷新
          </NButton>
          <slot name="actions" />
        </NSpace>
      </div>
    </NCard>
    <slot name="toolbar" />
    <NCard :bordered="false" size="small" content-class="!p-0">
      <NAlert v-if="error" type="warning" :show-icon="false" class="m-16px">
        {{ error }}
      </NAlert>
      <div v-if="loading && !error" class="min-h-240px flex-center">
        <NSpin size="medium" />
      </div>
      <NEmpty v-else-if="empty && !error" class="min-h-240px flex-center" :description="emptyDescription" />
      <slot v-else />
    </NCard>
  </NSpace>
</template>
