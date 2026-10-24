<script setup lang="ts">
import { computed } from 'vue';
import type { Component } from 'vue';
import { loginModuleRecord } from '@/constants/app';
import brandBg from '@/assets/imgs/login-brand-bg.webp';
import PwdLogin from './modules/pwd-login.vue';

interface Props {
  /** The login module */
  module?: UnionKey.LoginModule;
}

const props = defineProps<Props>();

interface LoginModule {
  label: App.I18n.I18nKey;
  component: Component;
}

const moduleMap: Record<UnionKey.LoginModule, LoginModule> = {
  'pwd-login': { label: loginModuleRecord['pwd-login'], component: PwdLogin }
};

const activeModule = computed(() => moduleMap[props.module || 'pwd-login']);
</script>

<template>
  <div class="login-page">
    <!-- 背景：模糊夜色 + 两团光晕 -->
    <div class="bg-layer" aria-hidden="true">
      <div class="bg-blob blob-a"></div>
      <div class="bg-blob blob-b"></div>
      <div class="bg-vignette"></div>
    </div>

    <!-- 居中卡片 -->
    <main class="stage">
      <section class="login-card" aria-label="登录">
        <!-- 左：品牌视觉 -->
        <div class="brand-panel" :style="{ backgroundImage: `url(${brandBg})` }">
          <div class="brand-content">
            <div class="brand-mark">
              <svg viewBox="0 0 28 28" fill="none" aria-hidden="true">
                <path d="M14 3.4 L25.4 24.6 H19.4 L14 13.6 L8.6 24.6 H2.6 Z" fill="currentColor" />
              </svg>
              <span class="brand-word">Alien Panel</span>
            </div>

            <div class="brand-copy">
              <h2>掌控每一台节点</h2>
              <p>统一的运营控制台，节点、用户与账单清晰可控。</p>
            </div>
          </div>
        </div>

        <!-- 右：登录表单 -->
        <div class="form-panel">
          <div class="form-inner">
            <header class="form-head">
              <span class="form-mark" aria-hidden="true">
                <svg viewBox="0 0 28 28" fill="none">
                  <path d="M14 3.4 L25.4 24.6 H19.4 L14 13.6 L8.6 24.6 H2.6 Z" fill="currentColor" />
                </svg>
              </span>
              <p class="form-sub">中央管理面板</p>
            </header>

            <Transition name="fade-slide" mode="out-in" appear>
              <component :is="activeModule.component" />
            </Transition>

            <p class="form-foot">Alien Panel · 中央运营控制台</p>
          </div>
        </div>
      </section>
    </main>
  </div>
</template>

<style scoped>
.login-page {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  min-height: 100vh;
  min-height: 100dvh;
  overflow: hidden;
  background: #05050d;
  color: #f2f2f5;
  font-family: Inter, -apple-system, BlinkMacSystemFont, 'Segoe UI', 'PingFang SC', 'Hiragino Sans GB',
    'Microsoft YaHei', sans-serif;
}

/* ---------- 背景 ---------- */

.bg-layer {
  position: fixed;
  inset: 0;
  z-index: 0;
  overflow: hidden;
  background:
    radial-gradient(92% 86% at 10% 8%, rgba(42, 39, 145, 0.34) 0%, transparent 64%),
    radial-gradient(86% 92% at 92% 88%, rgba(10, 92, 119, 0.3) 0%, transparent 67%),
    radial-gradient(74% 100% at 52% 50%, rgba(24, 16, 70, 0.24) 0%, transparent 72%), #05050d;
}

.bg-blob {
  position: absolute;
  border-radius: 50%;
  filter: blur(110px);
  will-change: transform;
}

.blob-a {
  width: 66vw;
  height: 66vw;
  left: -18vw;
  top: -8vh;
  background: radial-gradient(circle at 46% 44%, #4a4bd2 0%, #1b216e 48%, transparent 74%);
  opacity: 0.28;
  animation: drift 28s ease-in-out infinite alternate;
}

.blob-b {
  width: 54vw;
  height: 54vw;
  right: -14vw;
  bottom: -12vh;
  background: radial-gradient(circle at 55% 50%, #1dc2d0 0%, #104878 50%, transparent 74%);
  opacity: 0.22;
  animation: drift 36s ease-in-out infinite alternate-reverse;
}

@keyframes drift {
  from {
    transform: translate3d(0, 0, 0) scale(1);
  }
  to {
    transform: translate3d(2.5vw, -1.6vh, 0) scale(1.07);
  }
}

.bg-vignette {
  position: absolute;
  inset: 0;
  background: radial-gradient(72% 74% at 50% 50%, transparent 30%, rgba(2, 3, 8, 0.74) 100%);
  pointer-events: none;
}

/* ---------- 舞台与卡片 ---------- */

.stage {
  position: relative;
  z-index: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  min-height: 100dvh;
  padding: 32px;
}

.login-card {
  display: grid;
  grid-template-columns: 1fr 1fr;
  width: min(1000px, 92vw);
  height: min(700px, 88vh);
  border: 1px solid rgba(164, 178, 255, 0.11);
  border-radius: 18px;
  background: #0d0e19;
  box-shadow:
    0 60px 130px -34px rgba(0, 0, 0, 0.96),
    0 0 0 1px rgba(255, 255, 255, 0.015) inset;
  overflow: hidden;
  animation: cardIn 0.85s cubic-bezier(0.16, 1, 0.3, 1) both;
}

@keyframes cardIn {
  from {
    opacity: 0;
    transform: translateY(16px) scale(0.988);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

/* ---------- 左侧品牌视觉 ---------- */

.brand-panel {
  position: relative;
  overflow: hidden;
  isolation: isolate;
  background-color: #04060d;
  background-size: cover;
  background-position: 25% center;
  background-repeat: no-repeat;
}

/* 上弱下强的暗角，保证 logo 与文案对比度 */
.brand-panel::after {
  content: '';
  position: absolute;
  inset: 0;
  z-index: 1;
  pointer-events: none;
  background: linear-gradient(
    180deg,
    rgba(2, 4, 11, 0.14) 0%,
    rgba(2, 4, 11, 0.02) 38%,
    rgba(2, 4, 11, 0.28) 66%,
    rgba(2, 4, 11, 0.62) 100%
  );
}

.brand-content {
  position: relative;
  z-index: 4;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  height: 100%;
  padding: 38px 42px 52px;
}

.brand-mark {
  display: flex;
  align-items: center;
  gap: 10px;
}

.brand-mark svg {
  width: 23px;
  height: 23px;
  color: #ffffff;
  flex: none;
}

.brand-word {
  font-size: 13px;
  font-weight: 600;
  letter-spacing: 0.12em;
  text-transform: uppercase;
  color: #ffffff;
  white-space: nowrap;
}

.brand-copy h2 {
  margin: 0 0 13px;
  font-size: 26px;
  line-height: 1.38;
  font-weight: 650;
  letter-spacing: 0;
  color: #ffffff;
  text-wrap: balance;
  text-shadow: 0 2px 24px rgba(0, 0, 0, 0.55);
}

.brand-copy p {
  margin: 0;
  max-width: 300px;
  font-size: 13px;
  line-height: 1.75;
  letter-spacing: 0.01em;
  color: rgba(214, 228, 245, 0.76);
  text-shadow: 0 1px 14px rgba(0, 0, 0, 0.5);
}

/* ---------- 右侧表单 ---------- */

.form-panel {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 52px 56px;
  background:
    radial-gradient(140% 68% at 28% -10%, rgba(94, 87, 220, 0.12) 0%, transparent 62%),
    linear-gradient(180deg, #17182a 0%, #131525 58%, #0f111d 100%);
}

.form-inner {
  width: 100%;
  max-width: 348px;
}

.form-head {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 14px;
  margin-bottom: 38px;
}

.form-mark {
  display: block;
  color: #ffffff;
}

.form-mark svg {
  display: block;
  width: 48px;
  height: 48px;
}

.form-sub {
  margin: 0;
  font-size: 12.5px;
  line-height: 1.5;
  font-weight: 500;
  letter-spacing: 0.16em;
  /* letter-spacing 会在末尾多出一份空隙，左侧补等量缩进让视觉真正居中 */
  text-indent: 0.16em;
  color: #8b90a6;
}

.form-foot {
  margin: 26px 0 0;
  font-size: 11px;
  letter-spacing: 0.05em;
  color: #50556f;
  text-align: center;
}

/* ---------- 过渡动画 ---------- */

.fade-slide-enter-active,
.fade-slide-leave-active {
  transition: opacity 0.22s ease, transform 0.22s ease;
}

.fade-slide-enter-from {
  opacity: 0;
  transform: translateY(8px);
}

.fade-slide-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}

/* ---------- 响应式 ---------- */

/* ============ 移动端：全屏版式，不再使用浮动卡片 ============ */
@media (max-width: 900px) {
  .login-page {
    overflow-y: auto;
  }

  .stage {
    width: 100%;
    padding: 0;
    align-items: stretch;
    min-height: 100dvh;
  }

  .login-card {
    display: flex;
    flex-direction: column;
    width: 100%;
    max-width: none;
    min-height: 100dvh;
    height: auto;
    border: 0;
    border-radius: 0;
    box-shadow: none;
    background: transparent;
    animation: none;
  }

  .login-card > * {
    min-width: 0;
  }

  /* 品牌图：顶部主视觉，占屏幕约 1/3 */
  .brand-panel {
    flex: none;
    width: 100%;
    height: 36vh;
    min-height: 210px;
    max-height: 320px;
    background-position: center 35%;
  }

  .brand-panel::after {
    background: linear-gradient(
      180deg,
      rgba(2, 4, 11, 0.16) 0%,
      rgba(2, 4, 11, 0.04) 38%,
      rgba(18, 19, 34, 0.55) 78%,
      #121322 100%
    );
  }

  .brand-content {
    padding: 20px 24px 22px;
    justify-content: space-between;
  }

  .brand-mark {
    gap: 9px;
  }

  .brand-mark svg {
    width: 21px;
    height: 21px;
  }

  .brand-word {
    font-size: 12px;
    letter-spacing: 0.11em;
  }

  .brand-copy h2 {
    margin: 0;
    font-size: 21px;
    line-height: 1.32;
  }

  .brand-copy p {
    display: none;
  }

  /* 表单区：填满剩余高度并垂直居中 */
  .form-panel {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 30px 24px calc(24px + env(safe-area-inset-bottom, 0px));
    border-radius: 0;
    background: linear-gradient(180deg, #121322 0%, #0d0f1b 100%);
  }

  .form-inner {
    width: 100%;
    max-width: 400px;
    margin: 0 auto;
  }

  .form-head {
    gap: 12px;
    margin-bottom: 28px;
  }

  .form-mark svg {
    width: 42px;
    height: 42px;
  }

  .form-sub {
    font-size: 12.5px;
  }

  .form-foot {
    margin-top: 24px;
  }
}

/* 矮屏（横屏手机 / 小屏）压缩主视觉，避免表单被挤出视口 */
@media (max-width: 900px) and (max-height: 700px) {
  .brand-panel {
    height: 26vh;
    min-height: 150px;
    max-height: 200px;
  }

  .brand-copy h2 {
    font-size: 18px;
  }

  .form-panel {
    padding: 24px 22px calc(22px + env(safe-area-inset-bottom, 0px));
  }

  .form-head {
    margin-bottom: 22px;
  }

  .form-mark svg {
    width: 36px;
    height: 36px;
  }
}

@media (max-width: 380px) {
  .brand-panel {
    height: 32vh;
    min-height: 180px;
  }

  .form-panel {
    padding: 28px 18px calc(22px + env(safe-area-inset-bottom, 0px));
  }
}

@media (prefers-reduced-motion: reduce) {
  .login-card,
  .bg-blob {
    animation: none;
  }

  .fade-slide-enter-active,
  .fade-slide-leave-active {
    transition: none;
  }
}
</style>
