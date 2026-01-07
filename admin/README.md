# Polygame Admin Dashboard

Vue 3 + TypeScript + PrimeVue 管理后台

## 技术栈

- **框架**: Vue 3 (Composition API + TypeScript)
- **UI 库**: PrimeVue
- **图表**: Chart.js + vue-chartjs
- **状态管理**: Pinia
- **路由**: Vue Router

## 功能模块

### 赛事管理
- 创建市场
- 编辑市场信息
- 结算市场
- 市场状态管理

### 用户管理
- 用户列表
- 用户详情
- 虚拟积分查看

### 数据统计
- 总交易量
- 活跃用户数
- 市场热度
- 实时数据图表

### 系统配置
- 平台参数设置
- 公告管理

## 快速开始

### 安装依赖

```bash
npm install
```

### 开发模式

```bash
npm run dev
```

访问 http://localhost:3001

### 生产构建

```bash
npm run build
```

## 主要页面

- `/dashboard` - 仪表盘
- `/markets` - 市场管理
- `/users` - 用户管理
- `/statistics` - 数据统计
- `/settings` - 系统设置

## 权限要求

管理后台需要管理员权限（`is_admin: true`）才能访问。

## 许可证

MIT License
