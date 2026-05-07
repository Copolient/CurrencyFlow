# ExchangeApp

货币汇率查询与财经资讯平台，采用前后端分离的云原生架构。

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
| **测试** | Go testing + Mock 层 / Vitest + @testing-library/vue + @vue/test-utils |
| **可观测性** | OpenTelemetry → Jaeger / Prometheus + Grafana / Zap 结构化日志 |
| **代码规范** | golangci-lint (13 linters) / error wrapping / 输入校验 |
| **部署** | Docker + Kubernetes (Deployment / HPA / PDB / NetworkPolicy) |
| **CI/CD** | GitHub Actions + ArgoCD GitOps |

## 项目结构

```
exchangeapp/
├── Exchangeapp_backend/
│   ├── cmd/server/main.go              # 入口：依赖注入 + 优雅关闭
│   ├── internal/
│   │   ├── handler/                    # HTTP 处理层（参数绑定 + 响应）
│   │   ├── service/                    # 业务逻辑层 + 单元测试
│   │   ├── repository/                 # 数据访问层（interface + GORM 实现）
│   │   ├── mock/                       # Mock 实现（Cache / UserRepo / ArticleRepo / ExchangeRateRepo）
│   │   ├── model/                      # 领域模型（含索引定义 + 校验标签）
│   │   └── middleware/
│   │       ├── auth.go                 # JWT 认证
│   │       ├── observability.go        # OTel 链路追踪 + Prometheus 指标 + 结构化日志
│   │       ├── ratelimit.go            # IP 级限流（100 req/s）
│   │       └── security.go             # 安全响应头（CSP / X-Frame-Options / HSTS）
│   ├── pkg/
│   │   ├── auth/                       # JWT 签发 / 验证
│   │   ├── cache/                      # Cache interface + Redis 实现
│   │   ├── config/                     # 环境变量优先 + config.yml fallback
│   │   ├── database/                   # GORM 连接管理
│   │   ├── logger/                     # Zap JSON 日志（自动注入 TraceID）
│   │   ├── metrics/                    # Prometheus 指标（Counter / Histogram / Gauge）
│   │   └── tracing/                    # OpenTelemetry 初始化
│   ├── migrations/                     # 数据库迁移脚本
│   ├── deploy/k8s/                     # K8s 生产清单（Deployment / PDB / HPA / NetworkPolicy）
│   ├── deploy/observability/           # Grafana Dashboard JSON
│   ├── scripts/loadtest.js             # K6 阶梯式压测脚本
│   ├── docs/sre-performance-analysis.md # SRE 性能剖析报告
│   ├── .golangci.yml                   # Lint 配置
│   ├── Makefile                        # 工程操作入口
│   └── Dockerfile
│
├── Exchangeapp_frontend/
│   ├── src/
│   │   ├── components/
│   │   │   ├── Login.vue               # 登录（表单验证 + loading）
│   │   │   ├── Register.vue            # 注册（密码确认 + 校验）
│   │   │   └── ErrorBoundary.vue       # 全局错误捕获，防白屏
│   │   ├── views/                      # 首页 / 汇率 / 文章 / 文章详情
│   │   ├── router/                     # 路由 + 鉴权守卫
│   │   ├── store/                      # Pinia 状态管理
│   │   ├── types/                      # TypeScript 类型定义
│   │   ├── __tests__/                  # Vitest 测试（store / router / component）
│   │   ├── axios.ts                    # HTTP 客户端（统一错误处理 + 429 限流）
│   │   └── main.ts
│   ├── vitest.config.ts                # 测试配置
│   ├── Dockerfile                      # 多阶段构建 → Nginx
│   ├── nginx.conf                      # SPA 回退 + API 反向代理
│   └── .env / .env.development
│
├── docker-compose.yml                  # 全栈一键启动
└── .github/workflows/deploy.yml        # CI/CD 流水线
```

## 快速开始

### Docker Compose（推荐）

```bash
git clone https://github.com/Motionists/exchangeApp.git
cd exchangeApp

docker-compose up --build
docker-compose exec backend ./server --migrate   # 首次建表
```

访问：前端 http://localhost / 后端 http://localhost:3000 / 健康检查 http://localhost:3000/healthz

### 本地开发

**后端：**

```bash
cd Exchangeapp_backend
export JWT_SECRET=dev-secret
export DB_DSN="root:password@tcp(127.0.0.1:3306)/exchangeapp?charset=utf8mb4&parseTime=True&loc=Local"
export REDIS_ADDR=localhost:6379

make migrate   # 首次建表
make run       # 启动服务 :3000
make test      # 运行测试
make lint      # 代码检查
```

**前端：**

```bash
cd Exchangeapp_frontend
npm install
npm run dev    # http://localhost:5173（自动代理 /api → localhost:3000）
npm test       # 运行 Vitest 测试
```

## API 端点（v1）

> 向后兼容：`/api/v1/*` 和 `/api/*` 均可访问

### 公开接口

| 方法 | 路径 | 说明 |
|------|------|------|
| `POST` | `/api/v1/auth/register` | 用户注册（用户名 2-32 字符，密码 >= 6） |
| `POST` | `/api/v1/auth/login` | 用户登录 |
| `GET` | `/api/v1/exchangeRates` | 获取全部汇率 |
| `GET` | `/healthz` | 健康检查（DB 连接状态） |
| `GET` | `/metrics` | Prometheus 指标 |

### 需要认证（Authorization: Bearer \<token\>）

| 方法 | 路径 | 说明 |
|------|------|------|
| `POST` | `/api/v1/exchangeRates` | 创建汇率（币种 3 字符，汇率 > 0） |
| `POST` | `/api/v1/articles` | 创建文章 |
| `GET` | `/api/v1/articles` | 获取文章列表（带 Redis 缓存） |
| `GET` | `/api/v1/articles/:id` | 获取文章详情 |
| `POST` | `/api/v1/articles/:id/like` | 点赞文章（Redis INCR） |
| `GET` | `/api/v1/articles/:id/like` | 获取点赞数 |

## 架构设计

### 后端分层

```
HTTP 请求 → Middleware(认证 / 追踪 / 指标 / 限流 / 安全头)
  → Handler(参数绑定 + 校验) → Service(业务逻辑) → Repository(数据访问) → MySQL / Redis
```

- **Handler** 只做 HTTP 参数绑定和响应，不包含业务逻辑
- **Service** 纯业务逻辑，通过 interface 依赖 Repository 和 Cache，可独立单元测试
- **Repository** 定义 interface + Mock 实现，替换数据库只需新写一个实现

### 测试策略

```
┌─────────────────────────────────────────────────┐
│  Service 层测试（20 个用例）                      │
│  - 依赖注入 MockRepo + MockCache                 │
│  - 覆盖正常路径 / 缓存命中 / DB 错误 / 未找到    │
├─────────────────────────────────────────────────┤
│  前端测试（13 个用例）                            │
│  - Store: 登录 / 注册 / 登出 / loading 状态      │
│  - Router: 守卫拦截 / 认证放行 / 404 重定向       │
│  - Component: 表单渲染 / 字段存在 / 按钮文本      │
└─────────────────────────────────────────────────┘
```

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

### 安全设计

| 层级 | 措施 |
|------|------|
| 传输 | JWT Bearer Token，72h 过期 |
| 限流 | IP 级滑动窗口，100 req/s，自动清理过期 limiter |
| 响应头 | X-Frame-Options / CSP / X-Content-Type-Options / HSTS / Permissions-Policy |
| 输入校验 | binding 标签：min/max/len/gt，币种 3 字符，汇率 > 0 |
| 容器 | 非 root 运行，只读 FS，drop ALL capabilities |

### Kubernetes 部署

| 组件 | 说明 |
|------|------|
| Deployment | 3 副本，跨可用区 topologySpreadConstraints，非 root 运行 |
| PDB | minAvailable: 2，确保滚动更新零宕机 |
| HPA v2 | CPU 70% + Memory 80% + P99 延迟三维伸缩，3-20 副本 |
| NetworkPolicy | Ingress 仅放行 nginx + prometheus，Egress 仅放行 DB/Redis/OTel |
| ArgoCD | 自动同步 + selfHeal + prune |

```bash
kubectl apply -f Exchangeapp_backend/deploy/k8s/base/
```

### CI/CD 流水线

```
push to main → GitHub Actions:
  1. golangci-lint + go test -race
  2. Docker multi-arch build → GHCR (tag: sha-abc1234 / v1.2.3)
  3. 更新 gitops 仓库 → ArgoCD 自动同步到集群
```

## Makefile 命令

| 命令 | 说明 |
|------|------|
| `make build` | 编译后端二进制 |
| `make test` | 运行测试（-race -coverprofile） |
| `make test-cover` | 生成测试覆盖率报告 |
| `make lint` | golangci-lint 检查 |
| `make run` | 启动服务 |
| `make migrate` | 执行数据库迁移 |
| `make docker-build` | 构建 Docker 镜像 |
| `make docker-up` | docker-compose 启动 |
| `make clean` | 清理构建产物 |

## 环境变量

### 后端

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `JWT_SECRET` | JWT 签名密钥（**必填**） | — |
| `DB_DSN` | MySQL 连接串（**必填**） | — |
| `REDIS_ADDR` | Redis 地址 | `localhost:6379` |
| `APP_PORT` | 服务监听端口 | `:3000` |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | OTel Collector 地址 | — |

### 前端

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `VITE_API_BASE_URL` | API 基础路径 | `/api` |

