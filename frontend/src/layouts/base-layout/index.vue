<script setup lang="ts">
import { computed } from 'vue';
import { LAYOUT_SCROLL_EL_ID } from '@sa/materials';
import { useAppStore } from '@/store/modules/app';
import { useThemeStore } from '@/store/modules/theme';
import GlobalHeader from '../modules/global-header/index.vue';
import GlobalSider from '../modules/global-sider/index.vue';
import GlobalTab from '../modules/global-tab/index.vue';
import GlobalContent from '../modules/global-content/index.vue';
import GlobalFooter from '../modules/global-footer/index.vue';

defineOptions({
  name: 'BaseLayout'
});

const appStore = useAppStore();
const themeStore = useThemeStore();

const headerProps = computed<App.Global.HeaderProps>(() => ({
  showLogo: false,
  showMenu: false,
  showMenuToggler: true
}));

const siderCollapsed = computed(() => appStore.siderCollapse);

function handleMobileMaskClick() {
  appStore.setSiderCollapse(true);
}
</script>

<template>
  <div class="shell" :class="{ 'is-mobile': appStore.isMobile }">
    <!-- 侧栏面板 -->
    <aside
      v-if="!appStore.isMobile"
      class="shell-sider"
      :class="{ 'is-collapsed': siderCollapsed }"
    >
      <GlobalSider />
    </aside>

    <!-- 移动端侧栏抽屉 -->
    <template v-if="appStore.isMobile">
      <aside class="shell-sider is-mobile-drawer" :class="{ 'is-open': !siderCollapsed }">
        <GlobalSider />
      </aside>
      <div v-show="!siderCollapsed" class="shell-mask" @click="handleMobileMaskClick"></div>
    </template>

    <!-- 主面板：顶栏 + 标签栏 + 内容 -->
    <section class="shell-main">
      <header class="shell-header">
        <GlobalHeader v-bind="headerProps" />
      </header>
      <div
        v-if="themeStore.tab.visible"
        class="shell-tab"
        :style="{ height: `${themeStore.tab.height}px` }"
      >
        <GlobalTab />
      </div>
      <main :id="LAYOUT_SCROLL_EL_ID" class="shell-content">
        <GlobalContent />
      </main>
      <GlobalFooter v-if="themeStore.footer.visible" />
    </section>
  </div>
</template>

<style scoped>
.shell {
  display: flex;
  gap: 12px;
  width: 100%;
  height: 100vh;
  padding: 12px;
  overflow: hidden;
  background: var(--ui2-bg);
}

/* ---------- 侧栏面板 ---------- */
.shell-sider {
  flex: none;
  width: 246px;
  min-height: 0;
  overflow: hidden;
  border: 1px solid var(--ui2-line);
  border-radius: 16px;
  background: var(--ui2-panel);
  transition: width 0.25s ease;
}

.shell-sider.is-collapsed {
  width: 68px;
}

/* ---------- 主面板 ---------- */
.shell-main {
  display: flex;
  flex: 1;
  flex-direction: column;
  min-width: 0;
  min-height: 0;
  overflow: hidden;
  border: 1px solid var(--ui2-line);
  border-radius: 16px;
  background: var(--ui2-panel);
}

.shell-header {
  flex: none;
  height: 62px;
}

.shell-tab {
  flex: none;
  min-height: 0;
  overflow: hidden;
}

.shell-content {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  overflow-x: hidden;
}

/* ---------- 移动端 ---------- */
.shell-mask {
  position: fixed;
  inset: 0;
  z-index: 1000;
  background: rgba(0, 0, 0, 0.55);
}

.shell-sider.is-mobile-drawer {
  position: fixed;
  top: 12px;
  bottom: 12px;
  left: 12px;
  z-index: 1001;
  width: 246px;
  transform: translateX(-110%);
  transition: transform 0.25s ease;
}

.shell-sider.is-mobile-drawer.is-open {
  transform: translateX(0);
}

.shell.is-mobile {
  padding: 8px;
  gap: 0;
}

.shell.is-mobile .shell-main {
  border-radius: 14px;
}
</style>

<style lang="scss">
#__SCROLL_EL_ID__ {
  @include scrollbar();
}
</style>
