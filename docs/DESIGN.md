# Polygame - 模拟预测市场交易平台

## 项目概述

Polygame 是一个纯娱乐性质的模拟预测市场交易平台，用户使用虚拟货币对各类事件进行预测交易。平台不涉及真实货币交易，所有交易均使用平台内的虚拟积分。

## 核心特性

### 1. 纯娱乐模拟
- 用户使用虚拟积分（Virtual Credits）进行交易
- 新用户注册赠送初始虚拟积分
- 不支持充值或提现真实货币
- 明确标注为"模拟交易游戏平台"

### 2. 预测市场机制
- 类似 Polymarket 的预测市场交易机制
- 用户可以买卖事件结果的"份额"
- 价格反映市场对事件发生的概率预测
- 事件结算后，正确预测的用户获得虚拟积分奖励

### 3. 支持的市场类型
- 体育赛事（足球、篮球、网球等）
- 电竞比赛（LOL、DOTA2、CS:GO 等）
- 娱乐事件（电影票房、颁奖典礼、综艺节目）
- 科技事件（产品发布、公司业绩）
- 其他趣味预测

### 4. 实时交易
- 支持赛事进行中的实时交易（模拟"滚球"）
- 实时更新赔率和价格
- WebSocket 推送实时数据

## 技术架构

### 后端系统
- **语言**：Golang 1.21+
- **框架**：Gin Web Framework
- **数据库**：PostgreSQL 15+
- **缓存**：Redis
- **实时通信**：WebSocket
- **API 设计**：RESTful API + WebSocket

### 网页前台
- **框架**：Vue 3 + Vite
- **UI 库**：TailwindCSS
- **状态管理**：Pinia
- **路由**：Vue Router
- **HTTP 客户端**：Axios
- **响应式设计**：支持 PC、平板、手机自适应

### 管理后台
- **框架**：Vue 3 + TypeScript + Composition API
- **UI 库**：PrimeVue
- **图表库**：ECharts
- **功能模块**：
  - 赛事管理（创建、编辑、结算市场）
  - 用户管理（查看用户信息、虚拟积分余额）
  - 数据统计（交易量、活跃用户、市场热度）
  - 系统配置（平台参数、公告管理）

### 移动应用
- **框架**：Flutter 3.x
- **状态管理**：Provider / Riverpod
- **平台支持**：iOS + Android
- **功能**：与网页端功能一致

## 数据库设计

### 核心表结构

#### users（用户表）
- id, username, email, password_hash
- virtual_balance（虚拟积分余额）
- created_at, updated_at

#### markets（市场表）
- id, title, description, category
- start_time, end_time, resolution_time
- status（pending, active, closed, resolved）
- created_at, updated_at

#### market_outcomes（市场结果选项表）
- id, market_id, outcome_name
- current_price（当前价格，0-1 之间）
- total_shares（总份额数）

#### orders（订单表）
- id, user_id, market_id, outcome_id
- order_type（buy/sell）, shares, price
- status（pending, filled, cancelled）
- created_at

#### positions（持仓表）
- id, user_id, market_id, outcome_id
- shares（持有份额数）
- avg_price（平均成本价）

#### transactions（交易记录表）
- id, user_id, type（deposit, trade, settlement）
- amount, balance_after
- created_at

## API 设计

### 用户相关
- POST /api/auth/register - 注册
- POST /api/auth/login - 登录
- GET /api/user/profile - 获取用户信息
- GET /api/user/balance - 获取虚拟积分余额
- GET /api/user/positions - 获取持仓

### 市场相关
- GET /api/markets - 获取市场列表
- GET /api/markets/:id - 获取市场详情
- GET /api/markets/trending - 获取热门市场

### 交易相关
- POST /api/orders - 创建订单
- GET /api/orders - 获取订单列表
- DELETE /api/orders/:id - 取消订单

### 管理后台
- POST /api/admin/markets - 创建市场
- PUT /api/admin/markets/:id - 更新市场
- POST /api/admin/markets/:id/resolve - 结算市场
- GET /api/admin/statistics - 获取统计数据

### WebSocket
- /ws/market/:id - 订阅市场实时数据
- /ws/orderbook/:id - 订阅订单簿实时更新

## 项目结构

```
polygame/
├── backend/                 # 后端服务
│   ├── cmd/
│   │   └── server/         # 主程序入口
│   ├── internal/
│   │   ├── api/            # API 处理器
│   │   ├── service/        # 业务逻辑
│   │   ├── repository/     # 数据访问层
│   │   ├── model/          # 数据模型
│   │   ├── middleware/     # 中间件
│   │   └── websocket/      # WebSocket 处理
│   ├── config/             # 配置文件
│   ├── migrations/         # 数据库迁移
│   └── go.mod
│
├── frontend/               # 网页前台
│   ├── src/
│   │   ├── components/     # Vue 组件
│   │   ├── views/          # 页面视图
│   │   ├── stores/         # Pinia 状态管理
│   │   ├── router/         # 路由配置
│   │   ├── api/            # API 封装
│   │   └── assets/         # 静态资源
│   ├── package.json
│   └── vite.config.js
│
├── admin/                  # 管理后台
│   ├── src/
│   │   ├── components/
│   │   ├── views/
│   │   ├── stores/
│   │   ├── router/
│   │   └── api/
│   ├── package.json
│   └── vite.config.ts
│
├── mobile/                 # Flutter 移动应用
│   ├── lib/
│   │   ├── models/
│   │   ├── services/
│   │   ├── providers/
│   │   ├── screens/
│   │   └── widgets/
│   ├── pubspec.yaml
│   └── README.md
│
├── docs/                   # 文档
│   ├── API.md
│   ├── DEPLOYMENT.md
│   └── USER_GUIDE.md
│
└── README.md
```

## 部署方案

### 开发环境
- Docker Compose 一键启动所有服务
- 包含 PostgreSQL、Redis、后端服务

### 生产环境
- 后端：容器化部署（Docker/Kubernetes）
- 前端：静态文件部署（Nginx/CDN）
- 数据库：PostgreSQL 主从复制
- 缓存：Redis 集群

## 合规声明

平台将在显著位置声明：
- "本平台为纯娱乐性质的模拟交易游戏"
- "所有交易使用虚拟积分，不涉及真实货币"
- "虚拟积分无法兑换为真实货币或任何实物"
- "仅供娱乐和学习预测市场机制使用"

## 开发计划

1. 创建 GitHub 仓库
2. 搭建后端基础架构和数据库
3. 实现核心 API 和业务逻辑
4. 开发网页前台界面
5. 开发管理后台
6. 开发 Flutter 移动应用
7. 编写完整文档
8. 测试和优化
