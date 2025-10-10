# Exchange App

一个基于 Go 和 Vue.js 的货币汇率应用，提供用户注册登录、汇率查询、新闻文章浏览和点赞功能。

## 项目结构

- `Exchangeapp_backend/` - Go 后端服务
- `Exchangeapp_frontend/` - Vue.js 前端应用

## 功能特性

- 用户注册和登录
- 货币汇率查询
- 新闻文章浏览
- 文章点赞功能
- JWT 身份验证

## 技术栈

### 后端
- Go 1.24.4
- Gin Web 框架
- GORM ORM
- MySQL 数据库
- Redis 缓存
- JWT 认证
- Viper 配置管理

### 前端
- Vue 3
- TypeScript
- Vite 构建工具
- Element Plus UI 组件库
- Vant UI 组件库
- Pinia 状态管理
- Vue Router 路由管理
- Axios HTTP 客户端

## 安装和运行

### 后端

1. 进入后端目录：
   ```bash
   cd Exchangeapp_backend
   ```

2. 下载依赖：
   ```bash
   go mod tidy
   ```

3. 配置数据库：
   - 确保 MySQL 运行在 127.0.0.1:3306
   - 创建数据库 `bogger`
   - 更新 `config/config.yml` 中的数据库连接信息

4. 运行后端服务：
   ```bash
   go run main.go
   ```

   服务将在 http://localhost:3000 启动

### 前端

1. 进入前端目录：
   ```bash
   cd Exchangeapp_frontend
   ```

2. 安装依赖：
   ```bash
   npm install
   ```

3. 启动开发服务器：
   ```bash
   npm run dev
   ```

   前端将在 http://localhost:5173 启动

## 使用说明

1. 访问前端应用
2. 注册新用户或登录
3. 浏览汇率信息
4. 查看新闻文章
5. 对文章进行点赞

## API 文档

后端 API 端点：

- `POST /auth/register` - 用户注册
- `POST /auth/login` - 用户登录
- `GET /exchange-rates` - 获取汇率
- `GET /articles` - 获取文章列表
- `POST /articles/:id/like` - 点赞文章


