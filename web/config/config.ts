// https://umijs.org/config/
import { defineConfig } from 'umi';
import { join } from 'path';

import defaultSettings from './defaultSettings';
import proxy from './proxy';
import routes from './routes';

const { REACT_APP_ENV } = process.env;

// // 在这里添加打印语句
// // 打印当前环境变量的值
// if (process.env.NODE_ENV === 'production') {
//   console.log('当前环境：生产环境');
// } else {
//   console.log('当前环境：开发环境');
// }


export default defineConfig({
  // 前端配置了 publicPath 为 /static/，那么前端应用程序在生产环境下应该通过 /static/ 路径来访问静态资源
  publicPath: process.env.NODE_ENV === 'production' ? '/static/' : '/',
  hash: true,
  antd: {},
  dva: {
    hmr: true,
  },
  // for Ant Design Charts https://pro.ant.design/zh-CN/docs/graph
  scripts: [
    //全部注释掉，不使用cdn源，直接pacakge.json中引入
    // 'https://unpkg.com/react@17/umd/react.production.min.js',
    // 'https://unpkg.com/react-dom@17/umd/react-dom.production.min.js',
    // 'https://unpkg.com/@ant-design/charts@1.0.5/dist/charts.min.js',
    //使用 组织架构图、流程图、资金流向图、缩进树图 才需要使用
    //'https://unpkg.com/@ant-design/charts@1.0.5/dist/charts_g6.min.js',
  ],
  // externals 是 webpack 中的一个配置项，它允许你将一些模块标记为外部依赖，即不会被打包到最终的输出文件中。在这个配置项中，你可以将某些模块指定为外部依赖，并且指定他们在全局变量中的名称，这样在你的代码中使用这些模块时，webpack 就会从全局变量中引用它们，而不是将它们打包进输出文件中。
  externals: {
    // react: 'React',
    // 'react-dom': 'ReactDOM',
    // "@ant-design/charts": "charts"
   },
  layout: {
    // https://umijs.org/zh-CN/plugins/plugin-layout
    locale: true,
    siderWidth: 208,
    ...defaultSettings,
  },
  // https://umijs.org/zh-CN/plugins/plugin-locale
  locale: {
    // default zh-CN
    default: 'zh-CN',
    antd: true,
    // default true, when it is true, will use `navigator.language` overwrite default
    baseNavigator: true,
  },
  dynamicImport: {
    loading: '@ant-design/pro-layout/es/PageLoading',
  },
  targets: {
    ie: 11,
  },
  // umi routes: https://umijs.org/docs/routing
  routes,
  access: {},
  // Theme for antd: https://ant.design/docs/react/customize-theme-cn
  theme: {
    // 如果不想要 configProvide 动态设置主题需要把这个设置为 default
    // 只有设置为 variable， 才能使用 configProvide 动态设置主色调
    // https://ant.design/docs/react/customize-theme-variable-cn
    'root-entry-name': 'variable',
  },
  // esbuild is father build tools
  // https://umijs.org/plugins/plugin-esbuild
  esbuild: {},
  title: false,
  ignoreMomentLocale: true,
  proxy: proxy[REACT_APP_ENV || 'dev'],
  manifest: {
    basePath: '/',
  },
  // Fast Refresh 热更新
  fastRefresh: {},
  openAPI: [
    {
      requestLibPath: "import { request } from 'umi'",
      // 或者使用在线的版本
      // schemaPath: "https://gw.alipayobjects.com/os/antfincdn/M%24jrzTTYJN/oneapi.json"
      schemaPath: join(__dirname, 'oneapi.json'),
      mock: false,
    },
    {
      requestLibPath: "import { request } from 'umi'",
      schemaPath: 'https://gw.alipayobjects.com/os/antfincdn/CA1dOm%2631B/openapi.json',
      projectName: 'swagger',
    },
  ],
  nodeModulesTransform: { type: 'none' },
  mfsu: {},
  webpack5: {},
  exportStatic: {},
});
