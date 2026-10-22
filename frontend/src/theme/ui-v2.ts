import type { GlobalThemeOverrides } from 'naive-ui';

/**
 * Alien Panel UI v2 设计 token → Naive UI 主题覆盖
 *
 * 颜色、圆角、字号、间距均与 design-preview/assets/kit-v2.css 对齐。
 * 设计原则：
 * - 深色实心面板，不做半透明/毛玻璃
 * - 蓝 #4C6FFF 主操作，绿/琥珀/紫/青/红做语义色
 * - 次级文字最低 12px，正文 13px，标题 14px
 */
export const uiV2ThemeOverrides: GlobalThemeOverrides = {
  common: {
    primaryColor: '#4c6fff',
    primaryColorHover: '#5d7cff',
    primaryColorPressed: '#3f5ee8',
    primaryColorSuppl: '#5d7cff',
    infoColor: '#4c6fff',
    infoColorHover: '#5d7cff',
    infoColorPressed: '#3f5ee8',
    successColor: '#2ed47a',
    successColorHover: '#45dc8c',
    successColorPressed: '#1fb865',
    warningColor: '#f5b93d',
    warningColorHover: '#f7c65b',
    warningColorPressed: '#d99f22',
    errorColor: '#f3656f',
    errorColorHover: '#f77d86',
    errorColorPressed: '#d94f59',

    bodyColor: '#0a0b0d',
    cardColor: '#141519',
    modalColor: '#141519',
    popoverColor: '#191b20',
    tableColor: '#141519',
    actionColor: '#191b20',
    inputColor: '#0d0e11',
    borderColor: '#23262c',
    dividerColor: '#1c1f24',

    textColorBase: '#f2f3f5',
    textColor1: '#f2f3f5',
    textColor2: '#bcc2cc',
    textColor3: '#9aa1ad',
    textColorDisabled: '#6b7280',
    placeholderColor: '#8f96a1',

    borderRadius: '10px',
    borderRadiusSmall: '8px',

    fontSize: '13px',
    fontSizeMini: '12px',
    fontSizeTiny: '12px',
    fontSizeSmall: '12.5px',
    fontSizeMedium: '13px',
    fontSizeLarge: '14px',
    fontSizeHuge: '16px',

    lineHeight: '1.65',
    heightSmall: '29px',
    heightMedium: '34px',
    heightLarge: '38px',

    boxShadow1: '0 2px 6px rgba(0, 0, 0, 0.3)',
    boxShadow2: '0 6px 18px rgba(0, 0, 0, 0.45)',
    boxShadow3: '0 24px 60px -24px rgba(0, 0, 0, 0.85)'
  },

  Card: {
    color: '#141519',
    borderColor: '#23262c',
    borderRadius: '14px',
    paddingMedium: '18px',
    paddingSmall: '14px',
    titleFontSizeMedium: '14px',
    titleFontWeight: '600'
  },

  Button: {
    color: '#1b1e24',
    colorHover: '#232730',
    colorPressed: '#232730',
    colorFocus: '#232730',
    textColor: '#e6e8ec',
    textColorHover: '#ffffff',
    textColorPressed: '#ffffff',
    textColorFocus: '#ffffff',
    border: '1px solid #2a2e35',
    borderHover: '1px solid #363b44',
    borderPressed: '1px solid #363b44',
    borderFocus: '1px solid #363b44',
    borderRadiusTiny: '7px',
    borderRadiusSmall: '8px',
    borderRadiusMedium: '9px',
    borderRadiusLarge: '10px',
    fontWeight: '500',
    fontSizeMedium: '13px',
    fontSizeSmall: '12.5px',
    heightMedium: '34px',
    heightSmall: '29px',
    heightLarge: '38px'
  },

  Input: {
    color: '#0d0e11',
    colorFocus: '#0d0e11',
    colorDisabled: '#101114',
    textColor: '#f2f3f5',
    textColorDisabled: '#6b7280',
    placeholderColor: '#8f96a1',
    border: '1px solid #262a31',
    borderHover: '1px solid #32373f',
    borderFocus: '1px solid #4c6fff',
    borderDisabled: '1px solid #23262c',
    boxShadowFocus: '0 0 0 3px rgba(76, 111, 255, 0.18)',
    borderRadius: '9px',
    heightMedium: '34px',
    fontSizeMedium: '13px',
    paddingMedium: '0 10px'
  },

  InternalSelection: {
    color: '#0d0e11',
    colorActive: '#0d0e11',
    textColor: '#f2f3f5',
    placeholderColor: '#8f96a1',
    border: '1px solid #262a31',
    borderHover: '1px solid #32373f',
    borderActive: '1px solid #4c6fff',
    boxShadowActive: '0 0 0 3px rgba(76, 111, 255, 0.18)',
    borderRadius: '9px',
    heightMedium: '34px',
    fontSizeMedium: '13px'
  },

  InternalSelectMenu: {
    color: '#191b20',
    optionColorHover: 'rgba(163, 177, 255, 0.08)',
    optionTextColor: '#bcc2cc',
    optionTextColorActive: '#ffffff',
    borderRadius: '10px'
  },

  DataTable: {
    thColor: '#141519',
    thColorHover: '#191b20',
    thTextColor: '#a2a9b4',
    thFontWeight: '600',
    tdColor: '#141519',
    tdColorHover: '#191c21',
    tdColorStriped: '#17181d',
    tdTextColor: '#bcc2cc',
    borderColor: '#1c1f24',
    borderRadius: '14px',
    thPaddingMedium: '9px 14px',
    tdPaddingMedium: '10px 14px',
    fontSizeMedium: '13px',
    thFontSizeMedium: '11.5px'
  },

  Drawer: {
    color: '#141519',
    headerBorderBottom: '1px solid #1c1f24',
    footerBorderTop: '1px solid #1c1f24',
    titleFontSize: '16px',
    titleFontWeight: '600'
  },

  Modal: {
    color: '#141519',
    borderRadius: '16px'
  },

  Dialog: {
    color: '#141519',
    borderRadius: '14px',
    padding: '18px 20px'
  },

  Tag: {
    color: 'rgba(163, 177, 255, 0.1)',
    textColor: '#bcc2cc',
    border: '1px solid #2a2e35',
    borderRadius: '6px',
    fontSizeSmall: '11.5px'
  },

  Menu: {
    itemColorActive: 'rgba(76, 111, 255, 0.18)',
    itemColorActiveHover: 'rgba(76, 111, 255, 0.24)',
    itemColorHover: 'rgba(163, 177, 255, 0.06)',
    itemTextColor: '#b0b6c1',
    itemTextColorActive: '#ffffff',
    itemTextColorHover: '#ffffff',
    itemTextColorChildActive: '#ffffff',
    itemIconColor: '#6f7681',
    itemIconColorActive: '#a8b8ff',
    itemIconColorHover: '#b8bec7',
    itemIconColorChildActive: '#a8b8ff',
    borderRadius: '9px',
    itemHeight: '38px',
    itemFontSize: '13.5px'
  },

  Layout: {
    color: '#0a0b0d',
    siderColor: '#0f1013',
    headerColor: '#0f1013',
    tabColor: '#0f1013'
  },

  Dropdown: {
    color: '#191b20',
    optionColorHover: 'rgba(163, 177, 255, 0.08)',
    optionTextColor: '#bcc2cc',
    optionTextColorHover: '#ffffff',
    borderRadius: '10px',
    optionHeightMedium: '30px',
    optionFontSizeMedium: '13px'
  },

  Popover: {
    color: '#191b20',
    borderRadius: '10px',
    textColor: '#bcc2cc'
  },

  Tooltip: {
    color: '#23252b',
    textColor: '#e6e8ec',
    borderRadius: '8px'
  },

  Progress: {
    railColor: '#24272e',
    fillColor: '#4c6fff'
  },

  Switch: {
    railColor: '#2a2e35',
    railColorActive: '#4c6fff',
    buttonColor: '#ffffff'
  },

  Pagination: {
    itemColor: 'transparent',
    itemColorHover: 'rgba(163, 177, 255, 0.08)',
    itemColorActive: 'rgba(76, 111, 255, 0.18)',
    itemColorActiveHover: 'rgba(76, 111, 255, 0.24)',
    itemTextColor: '#9aa1ad',
    itemTextColorHover: '#f2f3f5',
    itemTextColorActive: '#ffffff',
    itemBorder: '1px solid transparent',
    itemBorderActive: '1px solid rgba(76, 111, 255, 0.45)',
    borderRadius: '7px',
    itemSize: '28px'
  },

  Empty: {
    textColor: '#9aa1ad',
    iconColor: '#4e5560'
  },

  Spin: {
    color: '#4c6fff',
    textColor: '#9aa1ad'
  },

  Message: {
    color: '#191b20',
    textColor: '#e6e8ec',
    borderRadius: '10px',
    boxShadow: '0 12px 30px -12px rgba(0, 0, 0, 0.8)'
  },

  Notification: {
    color: '#191b20',
    textColor: '#e6e8ec',
    borderRadius: '12px',
    boxShadow: '0 20px 50px -20px rgba(0, 0, 0, 0.85)'
  },

  Form: {
    labelTextColor: '#b0b6c1',
    labelFontSizeTopMedium: '12px',
    labelFontSizeLeftMedium: '12.5px',
    feedbackHeightMedium: '18px'
  }
};
