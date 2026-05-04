# ExchangeApp

货币汇率查询与财经资讯平台，采用前后端分离架构，支持云原生部署。

## 功能

- 用户注册 / 登录（JWT 鉴权）
- 多币种实时汇率查询与兑换计算
- 财经文章浏览与点赞
- 路由守卫与权限控制

## 技术栈

| 层级 | 技术 |
|------|------|
| **前端** | Vue 3 + TypeScript + Vite + Element Plus + Pinia + Vue Router |
| **后端** | Go 1.24 + Gin + GORM + go-redis v9 |
| **数据库** | MySQL 8.0 + Redis 7 |
| **可观测性** | OpenTelemetry → Jaeger / Prometheus + Grafana / Zap 结构化日志 |
| **部署** | Docker + Kubernetes (Deployment / HPA / PDB / NetworkPolicy) |
| **CI/CD** | GitHub Actions + ArgoCD GitOps |

## 项目结构

```
exchangeapp/
├── Exchangeapp_backend/            # Go 后端
│   ├── cmd/server/main.go          # 入口：依赖注入 + 优雅关闭
│   ├── internal/
│   │   ├── handler/                # HTTP 处理层（参数绑定 + 响应）
│   │   ├── service/                # 业务逻辑层
│   │   ├── repository/             # 数据访问层（interface + GORM 实现）
│   │   ├── model/                  # 领域模型
│   │   └── middleware/             # 认证 + 链路追踪 + 指标采集
│   ├── pkg/
│   │   ├── auth/                   # JWT 签发 / 验证
│   │   ├── cache/                  # Cache interface + Redis 实现
│   │   ├── config/                 # 环境变量优先 + config.yml fallback
│   │   ├── database/               # GORM 连接管理
│   │   ├── logger/                 # Zap 结构化日志（自动注入 TraceID）
│   │   ├── metrics/                # Prometheus 指标（Counter / Histogram / Gauge）
│   │   └── tracing/                # OpenTelemetry 初始化
│   ├── migrations/                 # 数据库迁移脚本
│   ├── deploy/k8s/                 # Kubernetes 生产清单
│   ├── deploy/observability/       # Grafana Dashboard JSON
│   ├── scripts/loadtest.js         # K6 阶梯式压测脚本
│   ├── docs/                       # SRE 性能剖析报告
│   ├── Dockerfile
│   └── docker-compose.yml          # 后端本地开发环境
│
├── Exchangeapp_frontend/           # Vue 3 前端
│   ├── src/
│   │   ├── components/             # Login / Register 组件
│   │   ├── views/                  # 首页 / 汇率 / 文章 / 文章详情
│   │   ├── router/                 # 路由 + 鉴权守卫
│   │   ├── store/                  # Pinia 状态管理
│   │   ├── types/                  # TypeScript 类型定义
│   │   ├── axios.ts                # HTTP 客户端（统一错误处理）
│   │   └── main.ts                 # 应用入口
│   ├── Dockerfile                  # 多阶段构建 → Nginx
│   ├── nginx.conf                  # SPA 回退 + API 反向代理
│   └── .env / .env.development     # 环境变量
│
├── docker-compose.yml              # 全栈一键启动
└── .github/workflows/deploy.yml    # CI/CD 流水线
```

## 快速开始

### Docker Compose 一键启动（推荐）

```bash
git clone https://github.com/Motionists/exchangeApp.git
cd exchangeApp

# 启动全部服务（前端 + 后端 + MySQL + Redis）
docker-compose up --build

# 首次运行需要执行数据库迁移
docker-compose exec backend ./server --migrate
```

访问：
- 前端：http://localhost
- 后端 API：http://localhost:3000
- 健康检查：http://localhost:3000/healthz

### 本地开发

**后端：**

```bash
cd Exchangeapp_backend

# 设置环境变量
export JWT_SECRET=dev-secret
export DB_DSN="root:your_password@tcp(127.0.0.1:3306)/exchangeapp?charset=utf8mb4&parseTime=True&loc=Local"
export REDIS_ADDR=localhost:6379

# 安装依赖 + 迁移 + 启动
go mod tidy
go run ./cmd/server/ --migrate   # 首次建表
go run ./cmd/server/             # 启动服务（默认 :3000）
```

**前端：**

```bash
cd Exchangeapp_frontend

npm install
npm run dev      # 开发服务器 http://localhost:5173
npm run build    # 生产构建
```

前端开发模式下，Vite 会自动将 `/api` 请求代理到 `http://localhost:3000`。

## API 端点

### 公开接口

| 方法 | 路径 | 说明 |
|------|------|------|
| `POST` | `/api/auth/register` | 用户注册 |
| `POST` | `/api/auth/login` | 用户登录 |
| `GET` | `/api/exchangeRates` | 获取全部汇率 |
| `GET` | `/healthz` | 健康检查 |
| `GET` | `/metrics` | Prometheus 指标 |

### 需要认证（Authorization: Bearer \<token\>）

| 方法 | 路径 | 说明 |
|------|------|------|
| `POST` | `/api/exchangeRates` | 创建汇率 |
| `POST` | `/api/articles` | 创建文章 |
| `GET` | `/api/articles` | 获取文章列表 |
| `GET` | `/api/articles/:id` | 获取文章详情 |
| `POST` | `/api/articles/:id/like` | 点赞文章 |
| `GET` | `/api/articles/:id/like` | 获取点赞数 |

## 架构设计

### 后端分层

```
HTTP 请求 → Middleware(认证/追踪/指标) → Handler(参数绑定)
  → Service(业务逻辑) → Repository(数据访问) → MySQL / Redis
```

- **Handler** 只做 HTTP 参数绑定和响应，不包含业务逻辑
- **Service** 纯业务逻辑，通过 interface 依赖 Repository 和 Cache，可独立单元测试
- **Repository** 定义 interface，替换数据库只需新写一个实现

### 可观测性

```
                  ┌─ Jaeger (链路追踪)
请求 → OTel SDK ──┤
                  ├─ Prometheus (指标) → Grafana Dashboard
                  │
                  └─ Zap JSON 日志 (含 TraceID) → ELK / Loki
```

- **链路追踪**：OpenTelemetry 自动传播 SpanContext，跨 handler → service → repository 全链路
- **指标**：请求总数 / 延迟分布 / 并发数 / P99（供 HPA 自定义伸缩）
- **日志**：结构化 JSON 输出，自动关联 TraceID

### Kubernetes 部署

| 组件 | 说明 |
|------|------|
| Deployment | 3 副本，跨可用区 topologySpreadConstraints，非 root 运行 |
| PDB | minAvailable: 2，确保滚动更新零宕机 |
| HPA v2 | CPU 70% + Memory 80% + P99 延迟三维伸缩，3-20 副本 |
| NetworkPolicy | Ingress 仅放行 nginx + prometheus，Egress 仅放行 DB/Redis/OTel |
| ArgoCD | 自动同步 + selfHeal + prune |

```bash
# 部署到 K8s
kubectl apply -f Exchangeapp_backend/deploy/k8s/base/
```

### CI/CD 流水线

```
push to main → GitHub Actions:
  1. golangci-lint + go test -race
  2. Docker multi-arch build → GHCR (tag: sha-abc1234)
  3. 更新 gitops 仓库 → ArgoCD 自动同步到集群
```

## 环境变量

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `JWT_SECRET` | JWT 签名密钥（**必填**） | — |
| `DB_DSN` | MySQL 连接串（**必填**） | — |
| `REDIS_ADDR` | Redis 地址 | `localhost:6379` |
| `APP_PORT` | 服务监听端口 | `:3000` |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | OTel Collector 地址 | — |

前端环境变量（`.env` 文件）：

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `VITE_API_BASE_URL` | API 基础路径 | `/api` |

## 许可证

MIT
