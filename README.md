# CurrencyFlow

CurrencyFlow 是一个货币汇率查询与财经资讯平台。后端使用 Go 和 Gin，前端使用 Vue 3 与 TypeScript，提供汇率查询、走势图、预警、AI 分析、社区和资讯等功能。

## 功能

- 汇率查询与兑换计算，支持直接和反向货币对
- 汇率历史走势，时间范围 1D / 1W / 1M / 3M / 1Y
- 基于历史统计或 Anthropic API 的 AI 分析
- 汇率预警，触发后写入通知并推送给已登录用户的 WebSocket
- 货币对收藏
- 社区帖子、点赞和关注
- 财经资讯浏览与点赞
- JWT 鉴权与用户资料
- PWA 支持，离线缓存仅包含公开接口

汇率采集默认每 30 分钟运行一次，并通过 WebSocket 广播给在线客户端。

## 技术栈

| 类别 | 技术 |
| --- | --- |
| 后端 | Go 1.25、Gin、GORM |
| 数据库 | MySQL 8 |
| 缓存 | Redis 7 |
| 前端 | Vue 3、TypeScript、Vite、Pinia、Element Plus、ECharts |
| 实时通信 | gorilla/websocket |
| 可观测性 | OpenTelemetry、Prometheus、Grafana、Zap |
| 部署 | Docker Compose、Kubernetes、ArgoCD |

## 项目结构

```text
CurrencyFlow/
├── backend/
│   ├── cmd/server/         # HTTP 服务入口
│   ├── cmd/seed/           # 演示数据
│   ├── internal/
│   │   ├── handler/        # HTTP handlers
│   │   ├── service/        # 业务逻辑
│   │   ├── repository/     # 数据访问
│   │   ├── model/          # 领域模型
│   │   ├── middleware/     # 鉴权、限流、日志与指标
│   │   ├── websocket/      # WebSocket hub
│   │   └── scheduler/      # 汇率采集
│   ├── pkg/                # config、auth、cache、database、logger 等
│   ├── migrations/         # GORM AutoMigrate
│   └── deploy/             # Kubernetes manifests
├── frontend/
│   ├── src/
│   │   ├── views/
│   │   ├── components/
│   │   ├── composables/
│   │   ├── store/
│   │   ├── router/
│   │   └── axios.ts
│   ├── public/
│   └── vite.config.ts
├── docker-compose.yml
└── .github/workflows/
```

## 本地开发

### 环境要求

- Go 1.25+
- Node.js 18+
- MySQL 8.0+
- Redis 7.0+
- Git

### 后端

```bash
git clone https://github.com/Copolient/CurrencyFlow.git
cd CurrencyFlow/backend

export JWT_SECRET="$(openssl rand -base64 48)"
export DB_DSN="root:@tcp(127.0.0.1:3306)/currencyflow?charset=utf8mb4&parseTime=True&loc=Local"
export REDIS_ADDR="localhost:6379"

mysql -u root -e "CREATE DATABASE IF NOT EXISTS currencyflow CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

make migrate
go run ./cmd/seed
make run
```

服务默认监听 `http://localhost:3000`：

- 健康检查：`GET /healthz`
- Prometheus 指标：`GET /metrics`

### 前端

```bash
cd CurrencyFlow/frontend
npm install
npm run dev
```

前端开发服务器位于 `http://localhost:5173`。Vite 会把 `/api` 代理到 `http://localhost:3000`，并启用 WebSocket 代理。

## Docker Compose

```bash
export JWT_SECRET="$(openssl rand -base64 48)"
docker compose up --build

# 首次启动后执行迁移和演示数据
docker compose exec backend ./currencyflow --migrate
docker compose exec backend ./currencyflow-seed
```

访问：

- 前端：`http://localhost`
- 后端：`http://localhost:3000`
- 健康检查：`http://localhost:3000/healthz`

也可以复制 `.env.example` 为 `.env`，在文件中设置 `JWT_SECRET`。

## 配置

后端环境变量：

| 变量 | 必填 | 默认值 | 说明 |
| --- | --- | --- | --- |
| `JWT_SECRET` | 是 | 无 | JWT 签名密钥，至少 32 个字符。包含 `change-me`、`change_me`、`placeholder` 等占位符的值会被拒绝 |
| `DB_DSN` | 是 | 无 | MySQL 连接字符串 |
| `REDIS_ADDR` | 否 | `localhost:6379` | Redis 地址 |
| `REDIS_PASSWORD` | 否 | 空 | Redis 密码 |
| `APP_PORT` | 否 | `:3000` | HTTP 监听端口 |
| `CORS_ALLOWED_ORIGINS` | 否 | `http://localhost:5173,http://localhost,http://localhost:80` | 允许的浏览器 Origin，逗号分隔 |
| `TRUSTED_PROXIES` | 否 | loopback 与私有网段 | 可信反向代理 CIDR。设为 `none` 可关闭 `X-Forwarded-*` 信任 |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | 否 | 空 | OTLP collector 地址 |
| `OTEL_TRACE_SAMPLE_RATE` | 否 | `0.1` | 链路采样率 |
| `LLM_BASE_URL` | 否 | `https://api.anthropic.com` | LLM API 地址 |
| `LLM_API_KEY` | 否 | 空 | LLM API key；为空时使用本地统计回退 |
| `LLM_MODEL` | 否 | `claude-sonnet-4-20250514` | LLM 模型名 |

前端环境变量：

| 变量 | 默认值 | 说明 |
| --- | --- | --- |
| `VITE_API_BASE_URL` | `/api/v1` | API 基础路径 |

## API

所有业务接口位于 `/api/v1`。需要鉴权的接口使用：

```http
Authorization: Bearer <token>
```

公开接口：

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| POST | `/api/v1/auth/register` | 注册 |
| POST | `/api/v1/auth/login` | 登录 |
| GET | `/api/v1/exchangeRates` | 当前汇率快照 |
| GET | `/api/v1/rates/history` | 汇率历史，支持 `from`、`to`、`range` |
| GET | `/api/v1/rates/latest` | 各货币对最新汇率 |
| GET | `/api/v1/posts` | 社区帖子列表 |
| GET | `/api/v1/users/:id` | 用户公开资料 |
| GET | `/api/v1/ws` | WebSocket 连接 |
| GET | `/healthz` | 健康检查 |
| GET | `/metrics` | Prometheus 指标 |

需要鉴权的接口包括：

- `POST /api/v1/ai/analyze`
- `POST /api/v1/exchangeRates`
- `POST /api/v1/articles`、`GET /api/v1/articles`、`GET /api/v1/articles/:id`
- `POST /api/v1/articles/:id/like`、`GET /api/v1/articles/:id/like`
- `POST /api/v1/favorites`、`GET /api/v1/favorites`、`DELETE /api/v1/favorites`、`GET /api/v1/favorites/check`
- `POST /api/v1/alerts`、`GET /api/v1/alerts`、`DELETE /api/v1/alerts/:id`
- `GET /api/v1/notifications`、`PUT /api/v1/notifications/:id/read`、`PUT /api/v1/notifications/read-all`、`GET /api/v1/notifications/unread-count`
- `POST /api/v1/posts`、`POST /api/v1/posts/:id/like`
- `POST /api/v1/users/:id/follow`、`DELETE /api/v1/users/:id/follow`、`GET /api/v1/users/:id/following`
- `PUT /api/v1/users/profile`

## WebSocket

连接地址：

```text
ws://localhost:3000/api/v1/ws
```

通过前端 Nginx 时使用：

```text
ws://localhost/api/v1/ws
```

可选参数：

- `pair`：订阅单个货币对，例如 `?pair=USD/CNY`。不传则接收所有汇率更新。
- `token`：JWT。传入后会绑定用户身份，用于接收告警通知。

汇率消息：

```json
{
  "type": "rate",
  "fromCurrency": "USD",
  "toCurrency": "CNY",
  "rate": 7.24,
  "timestamp": "2026-01-01T12:00:00Z"
}
```

通知消息：

```json
{
  "type": "notification",
  "id": 1,
  "notificationType": "alert_triggered",
  "title": "汇率预警触发: USD/CNY",
  "content": "当前汇率 7.1985 已低于目标汇率 7.2000",
  "read": false,
  "createdAt": "2026-01-01T12:00:00Z"
}
```

## 测试

后端：

```bash
cd backend

make test          # go test -race
make lint          # gofmt 检查 + go vet
make test-cover
```

前端：

```bash
cd frontend

npm test
npm run test:coverage
npm run lint
npx vue-tsc --noEmit
npm run build
```

## 部署

Kubernetes manifests 位于 `backend/deploy/k8s/base`。部署前需要：

- 创建 `currencyflow-secrets`，提供 `jwt-secret` 和 `db-dsn`。仓库中的 `secret.yaml` 是占位模板，应用会拒绝其中的占位 JWT 密钥。
- 修改 `currencyflow-config` 中的 `cors-allowed-origins` 为实际前端域名。
- 根据 Ingress 或负载均衡网段调整 `trusted-proxies`。

NetworkPolicy 默认允许 DNS、MySQL、Redis、OTel，以及访问外部 HTTPS 服务的 443 端口。汇率采集和 AI 分析依赖外部 HTTPS 出网。

汇率采集使用 Redis 租约。同一时间只有一个副本写入历史数据、快照并触发告警，其他副本仍会广播汇率更新给各自连接的客户端。

CI 与 CD：

- `.github/workflows/ci.yml`：运行后端测试、格式检查、`go vet`、前端类型检查、lint、测试和构建。
- `.github/workflows/deploy.yml`：构建多架构镜像并更新 `currencyflow-gitops` 仓库。需要在 GitHub Secrets 中配置 `GITOPS_TOKEN`。

## 安全

- JWT 密钥必须由部署环境提供，至少 32 个字符，默认值和常见占位符会被拒绝。
- 用户密码使用 bcrypt 哈希。
- 全局接口和登录注册接口分别限流。
- 后端通过 `TRUSTED_PROXIES` 决定是否信任 `X-Forwarded-For`，生产环境应只信任 ingress 或负载均衡网段。
- CORS 与 WebSocket Origin 通过 `CORS_ALLOWED_ORIGINS` 配置。
- Kubernetes 生产环境建议使用 Sealed Secrets 或 External Secrets Operator 管理密钥。

## 目录说明

- `backend/cmd/server`：HTTP 服务入口，负责依赖组装和优雅退出。
- `backend/cmd/seed`：写入演示用户、汇率、文章和社区数据。
- `backend/internal/scheduler`：定时获取汇率，维护历史记录和当前快照。
- `backend/internal/websocket`：管理 WebSocket 客户端、汇率广播和用户通知。
- `frontend/src/composables/useWebSocket.ts`：汇率与通知的 WebSocket 客户端。
- `frontend/src/router/index.ts`：路由和鉴权守卫。

## License

MIT
