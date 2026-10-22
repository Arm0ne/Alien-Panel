<script lang="ts" setup>
import { onMounted, onUnmounted } from 'vue';
import { useBoolean } from '@sa/hooks';
import SearchModal from './components/search-modal.vue';

defineOptions({ name: 'GlobalSearch' });

const { bool: show, setTrue: open, toggle } = useBoolean();

function handleShortcut(event: KeyboardEvent) {
  if ((event.ctrlKey || event.metaKey) && event.key.toLowerCase() === 'k') {
    event.preventDefault();
    open();
  }
}

onMounted(() => {
  window.addEventListener('keydown', handleShortcut);
});

onUnmounted(() => {
  window.removeEventListener('keydown', handleShortcut);
});
</script>

<template>
  <button class="search-pill" type="button" aria-label="全局搜索" @click="toggle">
    <SvgIcon icon="mdi:magnify" class="search-pill-ico" />
    <span class="search-pill-text">搜索页面 / 用户 / 节点</span>
    <kbd class="search-pill-kbd">Ctrl K</kbd>
  </button>
  <SearchModal v-model:show="show" />
</template>

<style scoped></style>
