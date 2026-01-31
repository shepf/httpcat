[English](../README.md) | 简体中文

# 🐱 HttpCat

> 轻量级、高效的 HTTP 文件传输服务，配备现代化 Web 界面和 AI 集成。

HttpCat 是一个基于 HTTP 的文件传输服务，旨在提供简单、高效、稳定的文件上传和下载功能。无论是临时分享还是批量传输文件，HttpCat 都将是你的优秀助手。

## ✨ 功能特点

- 🚀 **简单高效** - 易于部署，无需外部依赖
- 🎨 **现代化界面** - 基于 React 的美观管理界面
- 🤖 **MCP 支持** - AI 助手（Claude、Cursor、CodeBuddy）可直接管理你的文件
- 🐳 **Docker 就绪** - 一键 Docker 部署
- 🔐 **安全可靠** - 基于 Token 的上传认证
- 📊 **统计功能** - 详细的上传下载历史记录

## 📁 项目结构

```
httpcat/
├── server-go/              # 🔧 Go 后端
│   ├── cmd/                # 应用入口
│   │   └── httpcat.go
│   ├── internal/           # 内部包
│   │   ├── common/         # 公共工具
│   │   ├── handler/        # HTTP 处理器
│   │   ├── mcp/            # MCP 服务实现
│   │   ├── midware/        # 中间件（认证、指标）
│   │   ├── models/         # 数据模型
│   │   ├── p2p/            # P2P 功能
│   │   ├── server/         # 服务器核心
│   │   ├── storage/        # 存储层
│   │   └── conf/           # 配置文件
│   ├── go.mod
│   └── go.sum
│
├── web/                    # 🎨 React 前端
│   ├── src/                # 源代码
│   ├── config/             # UmiJS 配置
│   ├── mock/               # Mock 数据（仅开发环境）
│   └── package.json
│
├── scripts/                # 🛠️ 脚本目录
│   ├── build.sh            # 多平台构建脚本
│   ├── install.sh          # Linux 安装脚本
│   ├── uninstall.sh        # 卸载脚本
│   └── translations.sh     # i18n 翻译脚本
│
├── docs/                   # 📚 文档目录
│   ├── README-cn.md        # 中文说明
│   ├── BUILD.md            # 编译构建指南
│   ├── ReleaseNote.md      # 版本发布记录
│   ├── MCP_USAGE.md        # MCP 集成指南
│   └── ...                 # 其他设计文档
│
├── static/                 # 📦 前端构建产物
├── release/                # 📤 构建输出目录（已忽略）
│
├── Dockerfile              # Docker 配置
├── docker-compose.yml      # Docker Compose 配置
└── httpcat.service         # systemd 服务文件
```

## 🚀 快速开始

### 方式一：Docker（推荐）

```bash
# 使用 Docker Compose
docker-compose up -d

# 或直接使用 Docker
docker run -d --name httpcat \
  -p 8888:8888 \
  -v $(pwd)/data:/app/data \
  -v $(pwd)/upload:/app/upload \
  httpcat:latest
```

### 方式二：源码构建

```bash
# 完整构建（后端 + 前端）
./build.sh

# 或分别构建：

# 仅后端
cd server-go && go build -o httpcat ./cmd/httpcat.go

# 仅前端
cd web && npm install && npm run build
```

### 方式三：开发模式

```bash
# 终端 1：启动后端
cd server-go
go build -o httpcat ./cmd/httpcat.go
./httpcat -C ./internal/conf/svr.yml --static=../static/

# 终端 2：启动前端开发服务器
cd web
npm install --registry=https://registry.npmmirror.com
NODE_OPTIONS=--openssl-legacy-provider npm run start:dev
```

访问地址：
- **前端**: http://localhost:8000（开发）或 http://localhost:8888（生产）
- **后端 API**: http://localhost:8888/api/v1/

### 默认账号

| 字段 | 值 |
|------|------|
| 用户名 | `admin` |
| 密码 | `admin` |

> ⚠️ **安全提示**: 首次登录后请立即修改默认密码！

## 🎉 生产环境安装

### 快速安装

```bash
# 下载并解压
httpcat_version="v0.2.0"
mkdir httpcat && cd httpcat
tar -zxvf httpcat_$httpcat_version.tar.gz

# 安装
./install.sh

# 管理服务
systemctl status httpcat
systemctl start httpcat
systemctl stop httpcat

# 查看日志
tail -f /root/log/httpcat.log
```

### 手动安装

```bash
# 1. 创建目录
mkdir -p /home/web/website/upload/
mkdir -p /home/web/website/httpcat_web/
mkdir -p /etc/httpdcat/

# 2. 安装后端
cp httpcat /usr/local/bin/
cp conf/svr.yml /etc/httpdcat/

# 3. 安装前端
unzip httpcat_web.zip -d /home/web/website/httpcat_web/

# 4. 启动服务
httpcat --port=8888 \
  --static=/home/web/website/httpcat_web/ \
  --upload=/home/web/website/upload/ \
  --download=/home/web/website/upload/ \
  -C /etc/httpdcat/svr.yml
```

### 使用 systemd

```bash
# 复制服务文件
cp httpcat.service /usr/lib/systemd/system/

# 重载并启动
sudo systemctl daemon-reload
sudo systemctl enable httpcat
sudo systemctl start httpcat
```

## 🤖 MCP（模型上下文协议）支持

HttpCat 支持 MCP 协议，让 AI 助手可以直接管理你的文件服务器。

### 快速配置

在你的 MCP 客户端配置（Claude Desktop、Cursor、CodeBuddy 等）中添加：

```json
{
  "mcpServers": {
    "httpcat": {
      "type": "sse",
      "url": "http://your-server:8888/mcp/sse"
    }
  }
}
```

### 可用的 MCP 工具

| 工具 | 功能说明 |
|------|----------|
| `list_files` | 列出上传目录中的文件 |
| `get_file_info` | 获取文件详情（大小、MD5 等） |
| `upload_file` | 通过 MCP 上传文件（需要 Token） |
| `get_disk_usage` | 获取磁盘使用情况 |
| `get_upload_history` | 查询上传历史记录 |
| `request_delete_file` | 请求删除文件（第一步） |
| `confirm_delete_file` | 确认删除文件（第二步） |
| `get_statistics` | 获取上传/下载统计 |
| `verify_file_md5` | 验证文件 MD5 完整性 |

📖 详细 MCP 使用指南请查看 [docs/MCP_USAGE.md](../docs/MCP_USAGE.md)

## 📡 API 接口

### 上传文件

```bash
curl -v -F "f1=@/path/to/file" \
  -H "UploadToken: your-token" \
  http://localhost:8888/api/v1/file/upload
```

### 下载文件

```bash
wget -O filename.jpg http://localhost:8888/api/v1/file/download?filename=filename.jpg
```

### 列出文件

```bash
curl http://localhost:8888/api/v1/file/listFiles?dir=/
```

## ⚙️ 配置说明

配置文件：`svr.yml`

```yaml
# 服务器设置
port: 8888
upload_dir: "./upload"
download_dir: "./upload"
static_dir: "./static"

# 认证配置
app_key: "httpcat"
app_secret: "httpcat_app_secret"
enable_upload_token: true

# 数据库配置
enable_sqlite: true
sqlite_db_path: "./data/sqlite.db"

# 通知配置（企业微信 Webhook）
persistent_notify_url: ""
```

## 🍀 常见问题

### 忘记密码？

删除 SQLite 数据库并重启：

```bash
find / -name "*.db" | grep httpcat
rm /path/to/httpcat_sqlite.db
systemctl restart httpcat
```

重启后会自动创建默认管理员账号。

### Node.js 版本问题？

Node.js 17+ 需要使用旧版 OpenSSL provider：

```bash
NODE_OPTIONS=--openssl-legacy-provider npm run start:dev
```

推荐使用 Node.js v16.x 以获得最佳兼容性。

## 🛠️ 开发指南

### 环境要求

- Go 1.19+
- Node.js 16+（推荐 v16.18.0）
- npm 或 yarn

### 构建命令

```bash
# 完整构建（前端 + 后端 + Docker 镜像）
./build.sh

# 仅构建后端
./build.sh --backend-only

# 仅构建前端
./build.sh --frontend-only

# 跳过 Docker 构建
./build.sh --skip-docker
```

## 📝 许可证

本软件仅供个人使用，禁止用于商业目的。

- 禁止用于商业目的
- 必须保留版权和许可声明
- 本软件按 "原样" 提供，不承担任何保证

## 🌟 参与贡献

欢迎关注我们的 GitHub 项目！⭐ 点亮 star 了解我们的实时动态。

欢迎提出 issue 或提交 pull request。Good luck! 🍀
