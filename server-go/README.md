# HttpCat Backend (Go)

HttpCat 的后端服务，使用 Go 语言和 Gin 框架构建。

## 🚀 快速开始

### 环境要求

- Go 1.19+

### 编译

```bash
cd server-go
go build -o httpcat ./cmd/httpcat.go
```

### 运行

```bash
# 使用配置文件
./httpcat -C ./internal/conf/svr.yml

# 指定静态资源目录
./httpcat -C ./internal/conf/svr.yml --static=../static/

# 完整参数
./httpcat \
  --port=8888 \
  --static=../static/ \
  --upload=./upload/ \
  --download=./upload/ \
  -C ./internal/conf/svr.yml
```

### 命令行参数

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `-C` | 配置文件路径 | `./conf/svr.yml` |
| `--port` | 监听端口 | `8888` |
| `--static` | 静态资源目录 | `./website/static/` |
| `--upload` | 上传文件目录 | `./upload/` |
| `--download` | 下载文件目录 | `./upload/` |
| `-v` | 显示版本信息 | - |
| `-h` | 显示帮助信息 | - |

## 📁 目录结构

```
server-go/
├── cmd/
│   └── httpcat.go          # 应用入口
│
├── internal/               # 内部包（不对外暴露）
│   ├── common/             # 公共模块
│   │   ├── db.go           # 数据库操作
│   │   ├── defs.go         # 常量定义
│   │   ├── init.go         # 初始化
│   │   ├── page.go         # 分页工具
│   │   ├── response.go     # 响应封装
│   │   ├── user.go         # 用户管理
│   │   ├── userconfig/     # 用户配置
│   │   ├── utils/          # 工具函数
│   │   └── ylog/           # 日志模块
│   │
│   ├── conf/               # 配置文件
│   │   └── svr.yml
│   │
│   ├── handler/            # HTTP 处理器
│   │   └── v1/             # API v1
│   │       ├── conf.go     # 配置接口
│   │       ├── file.go     # 文件操作
│   │       ├── image_manage.go
│   │       ├── statistics.go
│   │       └── user.go     # 用户接口
│   │
│   ├── mcp/                # MCP 服务器
│   │   ├── server.go       # MCP 实现
│   │   └── auth_example.go
│   │
│   ├── midware/            # 中间件
│   │   ├── akskAuth.go     # AK/SK 认证
│   │   ├── tokenAuth.go    # Token 认证
│   │   └── metrics.go      # 指标收集
│   │
│   ├── models/             # 数据模型
│   ├── p2p/                # P2P 功能
│   ├── server/             # 服务器核心
│   │   ├── svr.go          # 服务启动
│   │   └── router.go       # 路由配置
│   │
│   └── storage/            # 存储层
│
├── data/                   # 数据目录（SQLite 等）
├── log/                    # 日志目录
├── go.mod
└── go.sum
```

## ⚙️ 配置文件

配置文件位于 `internal/conf/svr.yml`：

```yaml
# 服务器配置
port: 8888

# 文件目录
base_dir: "./"                    # 文件根目录（默认项目工作目录，生产环境建议改为 /data/httpcat_data/ 等绝对路径）
upload_dir: "website/upload/"     # 上传子目录（相对于 base_dir）
download_dir: "website/download/" # 下载子目录（相对于 base_dir）
static_dir: "./static"

# 认证配置
app_key: "httpcat"
app_secret: "httpcat_app_secret"
enable_upload_token: true

# 数据库配置
enable_sqlite: true
sqlite_db_path: "./data/sqlite.db"

# 通知配置（企业微信）
persistent_notify_url: ""

# P2P 配置（默认关闭）
enable_p2p: false

# MCP 配置
enable_mcp: true
```

## 🔌 API 接口

### 公开接口

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/file/upload` | 上传文件 |
| GET | `/api/v1/file/download` | 下载文件 |
| GET | `/api/v1/file/listFiles` | 列出文件 |
| GET | `/api/v1/conf/getVersion` | 获取版本 |

### 需要认证的接口

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/user/login/account` | 用户登录 |
| POST | `/api/v1/user/login/outLogin` | 用户登出 |
| GET | `/api/v1/user/currentUser` | 获取当前用户 |
| POST | `/api/v1/user/changePasswd` | 修改密码 |
| GET | `/api/v1/statistics/*` | 统计接口 |

### MCP 接口

| 路径 | 说明 |
|------|------|
| `/mcp/sse` | MCP SSE 连接端点 |

## 🔐 认证方式

### 1. JWT Token (用户认证)

用于 Web 界面登录后的 API 调用：

```bash
curl -H "Authorization: Bearer <jwt_token>" \
  http://localhost:8888/api/v1/user/currentUser
```

### 2. Upload Token (文件上传认证)

基于 AK/SK 生成的上传凭证：

```bash
curl -F "f1=@/path/to/file" \
  -H "UploadToken: httpcat:xxx:xxx" \
  http://localhost:8888/api/v1/file/upload
```

## 🧪 开发

### 运行测试

```bash
go test ./...
```

### 代码检查

```bash
go vet ./...
golangci-lint run
```

### 热重载开发

推荐使用 [air](https://github.com/cosmtrek/air)：

```bash
# 安装 air
go install github.com/cosmtrek/air@latest

# 启动热重载
air
```

## 📝 日志

日志文件位于 `log/` 目录，支持以下级别：

- DEBUG
- INFO
- WARN
- ERROR

配置日志级别：

```yaml
log_level: "info"
log_path: "./log/"
```

## 🐛 常见问题

### 1. 端口被占用

```bash
lsof -i :8888
kill -9 <PID>
```

### 2. 权限问题

确保上传/下载目录有写入权限：

```bash
chmod 755 ./upload/
```

### 3. SQLite 锁定

如果遇到数据库锁定问题，检查是否有多个进程访问同一数据库文件。

---

更多信息请参考 [项目主 README](../README.md)
