# HttpCat Web Frontend

HttpCat 的现代化 Web 管理界面，基于 [Ant Design Pro](https://pro.ant.design) 构建。

## 🚀 快速开始

### 环境要求

- Node.js 16+ (推荐 v16.18.0)
- npm 或 yarn

### 安装依赖

```bash
# 使用 npm（推荐使用国内镜像）
npm install --registry=https://registry.npmmirror.com

# 或使用 yarn
yarn
```

### 开发模式

```bash
# 启动开发服务器（连接本地后端 8888 端口）
npm run start:dev

# 如果使用 Node.js 17+，需要添加 OpenSSL 兼容参数
NODE_OPTIONS=--openssl-legacy-provider npm run start:dev
```

开发服务器运行在 http://localhost:8000，API 请求会自动代理到 http://127.0.0.1:8888

### 生产构建

```bash
npm run build
```

构建产物输出到 `dist/` 目录，需要复制到 `../static/` 供后端服务。

## 📁 目录结构

```
web/
├── config/                 # UmiJS 配置
│   ├── config.ts           # 主配置文件
│   ├── routes.ts           # 路由配置
│   ├── proxy.ts            # 开发代理配置
│   └── defaultSettings.ts  # 默认主题设置
│
├── src/
│   ├── components/         # 公共组件
│   ├── locales/            # 国际化文件
│   ├── pages/              # 页面组件
│   │   ├── user/           # 用户相关（登录等）
│   │   ├── Welcome/        # 首页
│   │   └── ...
│   ├── services/           # API 服务
│   └── app.tsx             # 应用入口
│
├── mock/                   # Mock 数据（仅开发环境）
├── public/                 # 静态资源
└── package.json
```

## 🔧 可用脚本

| 命令 | 说明 |
|------|------|
| `npm run start:dev` | 启动开发服务器（禁用 mock，代理到后端） |
| `npm run start` | 启动开发服务器（启用 mock） |
| `npm run build` | 生产环境构建 |
| `npm run lint` | 代码检查 |
| `npm run lint:fix` | 自动修复代码问题 |
| `npm test` | 运行测试 |

## ⚙️ 配置说明

### 代理配置 (config/proxy.ts)

开发环境下，API 请求会代理到后端服务：

```typescript
export default {
  dev: {
    '/api/': {
      target: 'http://127.0.0.1:8888',
      changeOrigin: true,
    },
  },
};
```

### 路由配置 (config/routes.ts)

所有页面路由在此配置，支持权限控制和嵌套路由。

### 环境变量

| 变量 | 说明 |
|------|------|
| `REACT_APP_ENV` | 环境标识 (dev/test/pre/prod) |
| `MOCK` | 是否启用 mock (`none` 禁用) |
| `UMI_ENV` | UmiJS 环境配置 |

## 🎨 技术栈

- **框架**: React 18 + UmiJS 3
- **UI 组件**: Ant Design 4 + Ant Design Pro Components
- **状态管理**: DVA (基于 Redux)
- **国际化**: UmiJS i18n
- **图表**: @ant-design/charts
- **HTTP 客户端**: umi-request

## 🐛 常见问题

### 1. OpenSSL 错误

**错误**: `Error: error:0308010C:digital envelope routines::unsupported`

**解决方案**: Node.js 17+ 需要使用 legacy OpenSSL provider：

```bash
NODE_OPTIONS=--openssl-legacy-provider npm run start:dev
```

或者降级到 Node.js 16.x。

### 2. Husky 安装失败

由于 `.git` 目录在父级目录，Husky 可能无法正确安装。可以跳过：

```bash
npm install --ignore-scripts
```

### 3. 依赖安装问题

清除缓存后重试：

```bash
rm -rf node_modules package-lock.json
npm cache clean --force
npm install --registry=https://registry.npmmirror.com
```

### 4. 登录提示 "错误的用户名和密码"

确保：
1. 后端服务正在运行 (`http://localhost:8888`)
2. 使用 `start:dev` 命令启动（禁用 mock）
3. 默认账号: `admin` / `admin`

## 🔗 相关资源

- [Ant Design Pro 文档](https://pro.ant.design/docs/getting-started)
- [UmiJS 文档](https://umijs.org/)
- [Ant Design 组件](https://4x.ant.design/components/overview-cn/)
- [图标库](https://www.iconfont.cn/)

## 📝 开发约定

### 代码风格

- 使用 ESLint + Prettier 保持代码一致性
- 组件使用 TypeScript 编写
- 样式使用 Less

### 提交规范

提交前会自动运行 lint-staged 检查代码质量。

### API 调用

所有 API 调用统一在 `src/services/` 目录定义，使用 `umi-request` 发起请求。

---

更多信息请参考 [项目主 README](../README.md)
