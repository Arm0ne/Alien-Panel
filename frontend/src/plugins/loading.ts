import { getColorPalette, getRgb } from '@sa/color';
import { DARK_CLASS } from '@/constants/app';
import { localStg } from '@/utils/storage';
import { toggleHtmlClass } from '@/utils/common';

/**
 * 启动加载页：在 Vue 挂载前注入 #app。
 * 视觉与 UI v2 一致：深色画布 + 品牌蓝 + 细进度条，不再使用模板默认的紫色渐变 Logo 与跳动方块。
 */
export function setupLoading() {
  const themeColor = localStg.get('themeColor') || '#646cff';
  const darkMode = localStg.get('darkMode') || false;
  const palette = getColorPalette(themeColor);
  const { r, g, b } = getRgb(themeColor);

  // 保留主题色变量，避免首屏期间 UnoCSS 的 primary 色阶缺失
  const cssVars = [`--primary-color: ${r} ${g} ${b}`]
    .concat(Array.from(palette.entries()).map(([key, value]) => `--logo-color-${key}: ${value}`))
    .join(';');

  if (darkMode) {
    toggleHtmlClass(DARK_CLASS).add();
  }

  const mark = `<svg viewBox="0 0 28 28" fill="none" aria-hidden="true">
      <path d="M14 3.4 25.4 24.6H19.4L14 13.6 8.6 24.6H2.6Z" fill="currentColor" />
    </svg>`;

  const loading = `
<div class="xp-boot" style="${cssVars}">
  <div class="xp-boot-mark">${mark}</div>
  <div class="xp-boot-text">
    <strong>ALIEN PANEL</strong>
    <span>中央运营控制台</span>
  </div>
  <div class="xp-boot-track"><span class="xp-boot-bar"></span></div>
</div>
<style>
  .xp-boot {
    position: fixed;
    inset: 0;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 16px;
    background: #0a0b0d;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'PingFang SC', 'Microsoft YaHei', sans-serif;
    -webkit-font-smoothing: antialiased;
  }

  .xp-boot-mark {
    display: grid;
    place-items: center;
    width: 44px;
    height: 44px;
    border-radius: 13px;
    background: linear-gradient(140deg, #6a72e0 0%, #4c54b8 55%, #3f4bb0 100%);
    box-shadow: 0 10px 24px -12px rgba(76, 111, 255, 0.9), inset 0 0 0 1px rgba(255, 255, 255, 0.12);
    color: #ffffff;
  }

  .xp-boot-mark svg {
    width: 25px;
    height: 25px;
  }

  .xp-boot-text {
    text-align: center;
    line-height: 1.4;
  }

  .xp-boot-text strong {
    display: block;
    font-size: 13px;
    font-weight: 650;
    letter-spacing: 0.16em;
    text-transform: uppercase;
    color: #f2f3f5;
  }

  .xp-boot-text span {
    display: block;
    margin-top: 3px;
    font-size: 11.5px;
    letter-spacing: 0.04em;
    color: #6b727c;
  }

  .xp-boot-track {
    position: relative;
    width: 132px;
    height: 2px;
    overflow: hidden;
    border-radius: 2px;
    background: #1b1e23;
  }

  .xp-boot-bar {
    position: absolute;
    top: 0;
    left: 0;
    width: 44px;
    height: 100%;
    border-radius: 2px;
    background: linear-gradient(90deg, rgba(76, 111, 255, 0.2), #4c6fff);
    animation: xp-boot-slide 1.15s ease-in-out infinite;
  }

  @keyframes xp-boot-slide {
    0% {
      transform: translateX(-46px);
    }
    100% {
      transform: translateX(136px);
    }
  }

  @media (prefers-reduced-motion: reduce) {
    .xp-boot-bar {
      width: 100%;
      opacity: 0.45;
      animation: none;
    }
  }
</style>`;

  const app = document.getElementById('app');

  if (app) {
    app.innerHTML = loading;
  }
}
