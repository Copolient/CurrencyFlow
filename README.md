# ExchangeApp

货币汇率查询与财经资讯平台，采用前后端分离的云原生架构。

## ✨ 功能一览

| 功能 | 说明 |
|------|------|
| 💱 货币兑换 | 多币种实时汇率查询与兑换计算 |
| 📈 行情走势 | ECharts 交互式图表，支持 1D/1W/1M/3M/1Y 时间范围 |
| 🤖 AI 分析师 | 基于历史数据的智能分析，趋势判断、关键价位 |
| 🔔 汇率预警 | 设置目标汇率，触发时实时推送通知 |
| 💰 货币收藏 | 收藏常用货币对，个人仪表盘展示 |
| 👥 社交社区 | 关注用户、发布交易观点、信息流 |
| 📰 财经资讯 | 文章浏览与点赞 |
| 🔒 用户系统 | JWT 鉴权、个人资料、关注/粉丝 |
| 📱 PWA 支持 | 离线缓存、移动端适配、桌面图标 |
| ⚡ WebSocket | 汇率变动实时推送，数字跳动 |

## 🛠️ 技术栈

| 层级 | 技术 |
|------|------|
| **前端** | Vue 3 + TypeScript + Vite + Element Plus + Pinia + Vue Router + ECharts |
| **后端** | Go 1.25 + Gin + GORM + go-redis v9 + gorilla/websocket |
| **数据库** | MySQL 8.0 + Redis 7 |
| **测试** | Go testing + Mock / Vitest + @testing-library/vue |
| **可观测性** | OpenTelemetry → Jaeger / Prometheus + Grafana / Zap |
| **部署** | Docker + Docker Compose + Kubernetes |
| **CI/CD** | GitHub Actions + ArgoCD GitOps |

## 📁 项目结构

```
exchangeapp/
├── Exchangeapp_backend/
│   ├── cmd/server/main.go
│   ├── internal/
│   │   ├── handler/           # HTTP 处理层
│   │   ├── service/           # 业务逻辑层 + 测试
│   │   ├── repository/        # 数据访问层
│   │   ├── model/             # 领域模型
│   │   ├── mock/              # Mock 实现
│   │   ├── middleware/         # 中间件
│   │   ├── websocket/         # WebSocket Hub
│   │   └── scheduler/         # 定时任务（汇率采集）
│   ├── pkg/                   # 公共库
│   ├── migrations/            # 数据库迁移
│   ├── deploy/                # K8s 部署配置
│   └── Dockerfile
│
├── Exchangeapp_frontend/
│   ├── src/
│   │   ├── views/             # 页面组件
│   │   ├── components/        # 通用组件
│   │   ├── composables/       # Vue Composables
│   │   ├── store/             # Pinia 状态管理
│   │   ├── router/            # 路由 + 鉴权守卫
│   │   ├── __tests__/         # 测试
│   │   └── axios.ts           # HTTP 客户端
│   ├── public/icons/          # PWA 图标
│   ├── vite.config.ts         # Vite + PWA 配置
│   └── Dockerfile
│
├── docker-compose.yml         # 全栈一键启动
└── .github/workflows/         # CI/CD 流水线
```

---

## 🚀 快速开始

### 前置条件

| 工具 | 最低版本 | 说明 |
|------|---------|------|
| Go | 1.21+ | 后端语言 |
| Node.js | 18+ | 前端构建 |
| MySQL | 8.0+ | 主数据库 |
| Redis | 7.0+ | 缓存 + 实时数据 |
| Git | 2.0+ | 版本控制 |

---

### 🍎 macOS

#### 1. 安装依赖

```bash
# Homebrew 安装
brew install go node mysql redis git

# 启动 MySQL 和 Redis 服务
brew services start mysql
brew services start redis

# 验证服务运行
mysql -u root -e "SELECT 1"
redis-cli ping
```

#### 2. 克隆项目

```bash
git clone https://github.com/Copolient/exchangeApp.git
cd exchangeApp
```

#### 3. 后端启动

```bash
cd Exchangeapp_backend

# 设置环境变量
export JWT_SECRET="your-secret-key-here"
export DB_DSN="root:@tcp(127.0.0.1:3306)/exchangeapp?charset=utf8mb4&parseTime=True&loc=Local"
export REDIS_ADDR="localhost:6379"

# 首次运行：创建数据库
mysql -u root -e "CREATE DATABASE IF NOT EXISTS exchangeapp CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 数据库迁移
make migrate

# 启动服务
make run
```

后端运行在 `http://localhost:3000`

#### 4. 前端启动

```bash
cd Exchangeapp_frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

前端运行在 `http://localhost:5173`

---

### 🪟 Windows

#### 1. 安装依赖

**方式一：使用 Scoop（推荐）**

```powershell
# 安装 Scoop（如果没有）
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
Invoke-RestMethod -Uri https://get.scoop.sh | Invoke-Expression

# 安装依赖
scoop install go nodejs mysql redis git
```

**方式二：手动安装**

1. **Go**: 从 https://go.dev/dl/ 下载安装
2. **Node.js**: 从 https://nodejs.org/ 下载 LTS 版本
3. **MySQL**: 从 https://dev.mysql.com/downloads/mysql/ 下载安装
4. **Redis**: 从 https://github.com/microsoftarchive/releases/releases/tag/win-3.2.100 下载，或使用 WSL
5. **Git**: 从 https://git-scm.com/download/win 下载

#### 2. 启动 MySQL 和 Redis

```powershell
# MySQL（以管理员身份运行）
net start mysql

# Redis（以管理员身份运行）
redis-server

# 或者使用 Windows 服务
net start Redis
```

#### 3. 克隆项目

```powershell
git clone https://github.com/Copolient/exchangeApp.git
cd exchangeApp
```

#### 4. 后端启动

```powershell
cd Exchangeapp_backend

# 设置环境变量（PowerShell）
$env:JWT_SECRET = "your-secret-key-here"
$env:DB_DSN = "root:@tcp(127.0.0.1:3306)/exchangeapp?charset=utf8mb4&parseTime=True&loc=Local"
$env:REDIS_ADDR = "localhost:6379"

# 或者使用 CMD
# set JWT_SECRET=your-secret-key-here
# set DB_DSN=root:@tcp(127.0.0.1:3306)/exchangeapp?charset=utf8mb4&parseTime=True&loc=Local
# set REDIS_ADDR=localhost:6379

# 创建数据库
mysql -u root -e "CREATE DATABASE IF NOT EXISTS exchangeapp CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 数据库迁移
go run cmd/server/main.go --migrate

# 启动服务
go run cmd/server/main.go
```

#### 5. 前端启动

```powershell
cd Exchangeapp_frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

---

### 🐧 Linux (Ubuntu/Debian)

#### 1. 安装依赖

```bash
# 更新包管理器
sudo apt update

# 安装 Go
wget https://go.dev/dl/go1.25.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.25.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# 安装 Node.js
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt install -y nodejs

# 安装 MySQL
sudo apt install -y mysql-server
sudo systemctl start mysql
sudo systemctl enable mysql

# 安装 Redis
sudo apt install -y redis-server
sudo systemctl start redis-server
sudo systemctl enable redis-server

# 安装 Git
sudo apt install -y git
```

#### 2. 配置 MySQL

```bash
# 设置 root 密码（如果需要）
sudo mysql_secure_installation

# 创建数据库
sudo mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS exchangeapp CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
```

#### 3. 克隆项目

```bash
git clone https://github.com/Copolient/exchangeApp.git
cd exchangeApp
```

#### 4. 后端启动

```bash
cd Exchangeapp_backend

# 设置环境变量
export JWT_SECRET="your-secret-key-here"
export DB_DSN="root:your-password@tcp(127.0.0.1:3306)/exchangeapp?charset=utf8mb4&parseTime=True&loc=Local"
export REDIS_ADDR="localhost:6379"

# 数据库迁移
make migrate

# 启动服务
make run
```

#### 5. 前端启动

```bash
cd Exchangeapp_frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

---

### 🐳 Docker Compose（推荐，全平台通用）

如果你不想手动安装 MySQL 和 Redis，可以使用 Docker Compose 一键启动：

```bash
# 安装 Docker Desktop
# macOS: https://docs.docker.com/desktop/install/mac-install/
# Windows: https://docs.docker.com/desktop/install/windows-install/
# Linux: https://docs.docker.com/desktop/install/linux-install/

# 克隆项目
git clone https://github.com/Copolient/exchangeApp.git
cd exchangeApp

# 一键启动所有服务
docker-compose up --build

# 首次运行：执行数据库迁移（新终端）
docker-compose exec backend ./server --migrate
```

访问：
- 前端: http://localhost
- 后端: http://localhost:3000
- 健康检查: http://localhost:3000/healthz

---

## 🧪 运行测试

### 后端测试

```bash
cd Exchangeapp_backend

# 运行所有测试
make test

# 运行测试并生成覆盖率报告
make test-cover

# 代码检查
make lint
```

### 前端测试

```bash
cd Exchangeapp_frontend

# 运行测试
npm test

# 监听模式
npm run test:watch

# 覆盖率报告
npm run test:coverage
```

---

## 📡 API 端点

### 公开接口

| 方法 | 路径 | 说明 |
|------|------|------|
| `POST` | `/api/v1/auth/register` | 用户注册 |
| `POST` | `/api/v1/auth/login` | 用户登录 |
| `GET` | `/api/v1/exchangeRates` | 获取全部汇率 |
| `GET` | `/api/v1/rates/history` | 汇率历史（支持 range 参数） |
| `GET` | `/api/v1/rates/latest` | 所有货币对最新汇率 |
| `GET` | `/api/v1/posts` | 社区帖子列表 |
| `GET` | `/api/v1/users/:id` | 用户公开资料 |
| `POST` | `/api/v1/ai/analyze` | AI 汇率分析 |
| `GET` | `/api/v1/ws` | WebSocket 连接 |
| `GET` | `/healthz` | 健康检查 |
| `GET` | `/metrics` | Prometheus 指标 |

### 需要认证（Authorization: Bearer token）

| 方法 | 路径 | 说明 |
|------|------|------|
| `POST` | `/api/v1/exchangeRates` | 创建汇率 |
| `POST` | `/api/v1/articles` | 创建文章 |
| `GET` | `/api/v1/articles` | 文章列表 |
| `GET` | `/api/v1/articles/:id` | 文章详情 |
| `POST` | `/api/v1/articles/:id/like` | 点赞文章 |
| `POST` | `/api/v1/favorites` | 添加收藏 |
| `GET` | `/api/v1/favorites` | 我的收藏 |
| `DELETE` | `/api/v1/favorites` | 取消收藏 |
| `POST` | `/api/v1/alerts` | 创建预警 |
| `GET` | `/api/v1/alerts` | 我的预警 |
| `DELETE` | `/api/v1/alerts/:id` | 删除预警 |
| `GET` | `/api/v1/notifications` | 通知列表 |
| `PUT` | `/api/v1/notifications/:id/read` | 标记已读 |
| `PUT` | `/api/v1/notifications/read-all` | 全部已读 |
| `GET` | `/api/v1/notifications/unread-count` | 未读数量 |
| `POST` | `/api/v1/posts` | 发布帖子 |
| `POST` | `/api/v1/posts/:id/like` | 点赞帖子 |
| `POST` | `/api/v1/users/:id/follow` | 关注用户 |
| `DELETE` | `/api/v1/users/:id/follow` | 取消关注 |
| `PUT` | `/api/v1/users/profile` | 更新资料 |

---

## ⚙️ 环境变量

### 后端

| 变量 | 必填 | 说明 | 默认值 |
|------|------|------|--------|
| `JWT_SECRET` | ✅ | JWT 签名密钥 | — |
| `DB_DSN` | ✅ | MySQL 连接串 | — |
| `REDIS_ADDR` | ❌ | Redis 地址 | `localhost:6379` |
| `APP_PORT` | ❌ | 服务端口 | `:3000` |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | ❌ | OTel Collector | — |

### 前端

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `VITE_API_BASE_URL` | API 基础路径 | `/api` |

---

## 📱 PWA 支持

本项目支持 PWA（Progressive Web App），可以：

- 📲 添加到手机桌面，体验接近原生 App
- 📴 离线访问（静态资源缓存）
- 🔄 自动更新

### 构建 PWA 版本

```bash
cd Exchangeapp_frontend

# 构建生产版本
npm run build

# 预览
npm run preview
```

构建后的文件在 `dist/` 目录，可以部署到任何静态文件服务器。

---

## 🏗️ Makefile 命令

| 命令 | 说明 |
|------|------|
| `make build` | 编译后端二进制 |
| `make test` | 运行测试（-race -coverprofile） |
| `make test-cover` | 生成覆盖率报告 |
| `make lint` | golangci-lint 检查 |
| `make run` | 启动服务 |
| `make migrate` | 执行数据库迁移 |
| `make docker-build` | 构建 Docker 镜像 |
| `make docker-up` | docker-compose 启动 |
| `make clean` | 清理构建产物 |

---

## 🔧 常见问题

### macOS: MySQL 连接失败

```bash
# 检查 MySQL 是否运行
brew services list | grep mysql

# 如果没有运行
brew services start mysql

# 如果使用 socket 连接
export DB_DSN="root:@unix(/tmp/mysql.sock)/exchangeapp?charset=utf8mb4&parseTime=True&loc=Local"
```

### Windows: Go 命令找不到

```powershell
# 检查 Go 是否在 PATH 中
go version

# 如果找不到，手动添加
$env:PATH += ";C:\Go\bin"
```

### Linux: Redis 连接被拒绝

```bash
# 检查 Redis 状态
sudo systemctl status redis-server

# 启动 Redis
sudo systemctl start redis-server

# 检查是否监听正确端口
redis-cli ping
```

### 前端: npm install 失败

```bash
# 清除缓存
npm cache clean --force

# 删除 node_modules 重新安装
rm -rf node_modules package-lock.json
npm install
```

---

## 📄 License

MIT License

---

## 🙏 致谢

- [Vue.js](https://vuejs.org/)
- [Gin](https://github.com/gin-gonic/gin)
- [Element Plus](https://element-plus.org/)
- [ECharts](https://echarts.apache.org/)
- [GORM](https://gorm.io/)
