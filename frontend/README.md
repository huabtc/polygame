# Polygame Frontend

Vue 3 + TailwindCSS 响应式前端应用

## 技术栈

- **框架**: Vue 3 (Composition API)
- **构建工具**: Vite
- **UI**: TailwindCSS
- **状态管理**: Pinia
- **路由**: Vue Router
- **HTTP 客户端**: Axios

## 功能特性

- ✅ 响应式设计（PC/平板/手机自适应）
- ✅ 用户注册/登录
- ✅ 市场浏览和搜索
- ✅ 实时交易
- ✅ 持仓管理
- ✅ 订单历史

## 快速开始

### 安装依赖

```bash
npm install
```

### 开发模式

```bash
npm run dev
```

访问 http://localhost:3000

### 生产构建

```bash
npm run build
```

构建产物在 `dist/` 目录

### 预览生产构建

```bash
npm run preview
```

## 项目结构

```
src/
├── api/              # API 封装
├── assets/           # 静态资源
├── components/       # Vue 组件
│   ├── Navbar.vue
│   ├── Footer.vue
│   └── MarketCard.vue
├── stores/           # Pinia 状态管理
│   ├── auth.js
│   ├── market.js
│   └── trading.js
├── views/            # 页面视图
│   ├── Home.vue
│   ├── Login.vue
│   ├── Register.vue
│   ├── Markets.vue
│   ├── MarketDetail.vue
│   ├── Portfolio.vue
│   └── Profile.vue
├── router/           # 路由配置
├── App.vue           # 根组件
└── main.js           # 入口文件
```

## 环境变量

创建 `.env` 文件：

```env
VITE_API_URL=http://localhost:8080/api/v1
```

## 响应式断点

TailwindCSS 默认断点：

- `sm`: 640px
- `md`: 768px
- `lg`: 1024px
- `xl`: 1280px
- `2xl`: 1536px

## 页面路由

- `/` - 首页
- `/login` - 登录
- `/register` - 注册
- `/markets` - 市场列表
- `/markets/:id` - 市场详情
- `/portfolio` - 投资组合
- `/profile` - 个人资料

## API 集成

所有 API 请求通过 `src/api/axios.js` 统一处理：

- 自动添加 JWT token
- 统一错误处理
- 401 自动跳转登录

## 许可证

MIT License
