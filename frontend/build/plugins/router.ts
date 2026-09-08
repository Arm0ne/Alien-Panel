import type { RouteMeta } from 'vue-router';
import ElegantVueRouter from '@elegant-router/vue/vite';
import type { RouteKey } from '@elegant-router/types';

export function setupElegantRouter() {
  return ElegantVueRouter({
    layouts: {
      base: 'src/layouts/base-layout/index.vue',
      blank: 'src/layouts/blank-layout/index.vue'
    },
    routeNameTransformer(routeName) {
      // Keep the node detail page as a standalone route. Using an underscore
      // here makes Elegant Router treat it as a child of `nodes`, while the
      // list page is also an index route. That combination is not supported
      // by the generated route typings and can crash HMR route regeneration.
      if (routeName === 'nodes_detail') return 'nodes-detail';

      return routeName;
    },
    routePathTransformer(routeName, routePath) {
      const key = routeName as RouteKey;

      if (key === 'login') {
        const modules: UnionKey.LoginModule[] = ['pwd-login'];

        const moduleReg = modules.join('|');

        return `/login/:module(${moduleReg})?`;
      }

      if (routeName === 'nodes-detail') return '/nodes/:id';

      return routePath;
    },
    onRouteMetaGen(routeName) {
      const key = routeName as RouteKey;

      const constantRoutes: RouteKey[] = ['login', '403', '404', '500'];

      const meta: Partial<RouteMeta> = {
        title: key,
        i18nKey: `route.${key}` as App.I18n.I18nKey
      };

      if (routeName === 'nodes-detail') {
        meta.hideInMenu = true;
        meta.activeMenu = 'nodes';
      }

      if (constantRoutes.includes(key)) {
        meta.constant = true;
      }

      return meta;
    }
  });
}
