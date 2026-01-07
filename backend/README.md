# Polygame Backend

Golang + PostgreSQL 后端服务

## 技术栈

- **语言**: Go 1.21+
- **框架**: Gin Web Framework
- **数据库**: PostgreSQL 15+
- **ORM**: GORM
- **认证**: JWT
- **缓存**: Redis

## 项目结构

```
backend/
├── cmd/
│   └── server/          # 主程序入口
├── internal/
│   ├── api/             # HTTP 处理器
│   ├── service/         # 业务逻辑层
│   ├── repository/      # 数据访问层
│   ├── model/           # 数据模型
│   ├── middleware/      # 中间件
│   └── websocket/       # WebSocket 处理
├── config/              # 配置管理
└── migrations/          # 数据库迁移
```

## 快速开始

### 1. 安装依赖

```bash
go mod download
```

### 2. 配置环境变量

复制 `.env.example` 为 `.env` 并修改配置：

```bash
cp .env.example .env
```

### 3. 启动 PostgreSQL

```bash
# 使用 Docker
docker run --name polygame-postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=polygame \
  -p 5432:5432 \
  -d postgres:15
```

### 4. 启动服务

```bash
go run cmd/server/main.go
```

服务将在 `http://localhost:8080` 启动

## API 文档

### 认证相关

#### 注册
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "username": "testuser",
  "email": "test@example.com",
  "password": "password123"
}
```

#### 登录
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "testuser",
  "password": "password123"
}
```

### 用户相关

#### 获取个人信息
```http
GET /api/v1/user/profile
Authorization: Bearer <token>
```

#### 获取余额
```http
GET /api/v1/user/balance
Authorization: Bearer <token>
```

### 市场相关

#### 获取市场列表
```http
GET /api/v1/markets?category=sports&status=active&page=1&page_size=20
Authorization: Bearer <token>
```

#### 获取市场详情
```http
GET /api/v1/markets/:id
Authorization: Bearer <token>
```

#### 获取热门市场
```http
GET /api/v1/markets/trending?limit=10
Authorization: Bearer <token>
```

#### 搜索市场
```http
GET /api/v1/markets/search?q=football
Authorization: Bearer <token>
```

### 交易相关

#### 下单
```http
POST /api/v1/trading/orders
Authorization: Bearer <token>
Content-Type: application/json

{
  "market_id": 1,
  "outcome_id": 1,
  "order_type": "buy",
  "shares": 10,
  "price": 0.6
}
```

#### 获取订单列表
```http
GET /api/v1/trading/orders?page=1&page_size=20
Authorization: Bearer <token>
```

#### 获取持仓
```http
GET /api/v1/trading/positions
Authorization: Bearer <token>
```

#### 取消订单
```http
DELETE /api/v1/trading/orders/:id
Authorization: Bearer <token>
```

### 管理员相关

#### 创建市场
```http
POST /api/v1/admin/markets
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "title": "NBA Finals 2026 Winner",
  "description": "Who will win the NBA Finals 2026?",
  "category": "sports",
  "image_url": "https://example.com/image.jpg",
  "outcomes": ["Lakers", "Celtics", "Warriors", "Other"]
}
```

#### 结算市场
```http
POST /api/v1/admin/markets/:id/resolve
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "winning_outcome_id": 1
}
```

## 数据库模型

### User (用户)
- ID, Username, Email, PasswordHash
- VirtualBalance (虚拟积分余额)
- IsAdmin (是否管理员)

### Market (市场)
- ID, Title, Description, Category
- Status (pending, active, closed, resolved)
- StartTime, EndTime, ResolutionTime

### Outcome (结果选项)
- ID, MarketID, OutcomeName
- CurrentPrice (当前价格 0-1)
- TotalShares, TotalVolume

### Order (订单)
- ID, UserID, MarketID, OutcomeID
- OrderType (buy, sell)
- Shares, Price, TotalCost
- Status (pending, filled, cancelled)

### Position (持仓)
- ID, UserID, MarketID, OutcomeID
- Shares (持有份额)
- AvgPrice (平均成本价)

### Transaction (交易记录)
- ID, UserID, Type
- Amount, BalanceAfter
- OrderID, MarketID

## 开发

### 运行测试
```bash
go test ./...
```

### 构建
```bash
go build -o polygame-server cmd/server/main.go
```

### 生产部署
```bash
# 设置生产模式
export GIN_MODE=release

# 运行服务
./polygame-server
```

## 许可证

MIT License
