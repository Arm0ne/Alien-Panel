/** Default theme settings */
export const themeSettings: App.Theme.ThemeSetting = {
  themeScheme: 'dark',
  grayscale: false,
  colourWeakness: false,
  recommendColor: false,
  themeColor: '#4c6fff',
  themeRadius: 10,
  otherColor: {
    info: '#4c6fff',
    success: '#2ed47a',
    warning: '#f5b93d',
    error: '#f3656f'
  },
  isInfoFollowPrimary: true,
  layout: {
    mode: 'vertical',
    scrollMode: 'content'
  },
  page: {
    animate: true,
    animateMode: 'fade-slide'
  },
  header: {
    height: 62,
    breadcrumb: {
      visible: true,
      showIcon: true
    },
    multilingual: {
      visible: false
    },
    globalSearch: {
      visible: true
    }
  },
  tab: {
    visible: true,
    cache: true,
    height: 40,
    mode: 'button',
    closeTabByMiddleClick: false
  },
  fixedHeaderAndTab: true,
  sider: {
    inverted: false,
    width: 246,
    collapsedWidth: 68,
    mixWidth: 90,
    mixCollapsedWidth: 64,
    mixChildMenuWidth: 200,
    autoSelectFirstMenu: false
  },
  footer: {
    visible: false,
    fixed: false,
    height: 48,
    right: true
  },
  watermark: {
    visible: false,
    text: 'Alien-Panel',
    enableUserName: false,
    enableTime: false,
    timeFormat: 'YYYY-MM-DD HH:mm'
  },
  tokens: {
    light: {
      colors: {
        container: 'rgb(255, 255, 255)',
        layout: 'rgb(247, 250, 252)',
        inverted: 'rgb(0, 20, 40)',
        'base-text': 'rgb(31, 31, 31)'
      },
      boxShadow: {
        header: '0 1px 2px rgb(0, 21, 41, 0.08)',
        sider: '2px 0 8px 0 rgb(29, 35, 41, 0.05)',
        tab: '0 1px 2px rgb(0, 21, 41, 0.08)'
      }
    },
    dark: {
      colors: {
        container: 'rgb(20, 21, 25)',
        layout: 'rgb(10, 11, 13)',
        inverted: 'rgb(0, 20, 40)',
        'base-text': 'rgb(242, 243, 245)'
      },
      boxShadow: {
        header: 'none',
        sider: 'none',
        tab: 'none'
      }
    }
  }
};

/**
 * Override theme settings
 *
 * If publish new version, use `overrideThemeSettings` to override certain theme settings
 *
 * UI v2 起强制深色：老用户 localStorage 里可能缓存了浅色设置，
 * 这里在版本升级（BUILD_TIME 变化）时覆盖一次。
 */
export const overrideThemeSettings: Partial<App.Theme.ThemeSetting> = {
  themeScheme: 'dark',
  themeColor: '#4c6fff',
  themeRadius: 10,
  header: {
    multilingual: {
      visible: false
    },
    globalSearch: {
      visible: true
    }
  } as App.Theme.ThemeSetting['header'],
  sider: {
    width: 246,
    collapsedWidth: 68
  } as App.Theme.ThemeSetting['sider'],
  footer: {
    visible: false
  } as App.Theme.ThemeSetting['footer']
};
