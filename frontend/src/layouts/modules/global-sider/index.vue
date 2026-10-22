<script setup lang="ts">
import { computed } from 'vue';
import { GLOBAL_SIDER_MENU_ID } from '@/constants/app';
import { useAppStore } from '@/store/modules/app';
import { useThemeStore } from '@/store/modules/theme';

defineOptions({
  name: 'GlobalSider'
});

const appStore = useAppStore();
const themeStore = useThemeStore();

const showLogo = computed(() => themeStore.layout.mode === 'vertical');
const menuWrapperClass = computed(() => (showLogo.value ? 'flex-1-hidden' : 'h-full'));
</script>

<template>
  <div class="v2-sider size-full flex-col-stretch">
    <div v-if="showLogo" class="v2-sider-head">
      <SystemLogo class="v2-sider-logo" />
      <div v-show="!appStore.siderCollapse" class="v2-sider-brand">
        <strong>Alien Panel</strong>
        <span>中央运营控制台</span>
      </div>
    </div>
    <div :id="GLOBAL_SIDER_MENU_ID" class="v2-sider-menu" :class="menuWrapperClass"></div>
  </div>
</template>

<style scoped></style>
