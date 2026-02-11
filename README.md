[English](docs/README-en.md) | 简体中文

# 🐱 HttpCat

> 轻量级、高效的 HTTP 文件传输服务，配备现代化 Web 界面和 AI 集成。

HttpCat 是一个基于 HTTP 的文件传输服务，旨在提供简单、高效、稳定的文件上传和下载功能。无论是临时分享还是批量传输文件，HttpCat 都将是你的优秀助手。

## ✨ 功能特点

- 🚀 **简单高效** - 易于部署，无需外部依赖
- 🎨 **现代化界面** - 基于 React 的美观管理界面
- 🤖 **MCP 支持** - AI 助手（Claude、Cursor、CodeBuddy）可直接管理你的文件
- 🐳 **Docker 就绪** - 一键 Docker 部署
- 🔐 **安全可靠** - 基于 Token 的上传认证
- 🔑 **Open API** - AK/SK 签名认证，脚本/CI/AI 可直接调用所有 API
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
│   ├── install.sh          # Linux/macOS 安装脚本
│   ├── uninstall.sh        # 卸载脚本
│   └── translations.sh     # i18n 翻译脚本
│
├── docs/                   # 📚 文档目录
│   ├── README-en.md        # English README
│   ├── BUILD.md            # 编译构建指南
│   ├── INSTALL.md          # 安装部署指南
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
./scripts/build.sh -a -f

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
httpcat_version="v0.2.1"
tar -zxvf httpcat_${httpcat_version}_linux-amd64.tar.gz
cd httpcat_${httpcat_version}_linux-amd64

# 安装（交互式）
sudo ./install.sh

# 或指定端口安装
sudo ./install.sh -p 9000

# 管理服务
sudo systemctl start httpcat
sudo systemctl status httpcat
```

### 安装后目录结构

使用 `install.sh` 安装后，文件按照 Linux FHS 标准分布：

```
/usr/local/bin/
└── httpcat                         # 可执行文件

/etc/httpcat/
└── svr.yml                         # 配置文件

/var/log/httpcat/
└── httpcat.log                     # 日志文件

/var/lib/httpcat/
├── static/                         # Web 界面静态资源
├── upload/                         # 上传文件存储
├── download/                       # 下载文件缓存
└── data/
    └── httpcat_sqlite.db           # SQLite 数据库
```

### 服务管理

```bash
# 启动/停止/重启
sudo systemctl start httpcat
sudo systemctl stop httpcat
sudo systemctl restart httpcat

# 查看状态和日志
sudo systemctl status httpcat
sudo journalctl -u httpcat -f
```

### 卸载

```bash
# 标准卸载（保留配置和数据）
sudo ./uninstall.sh

# 完全卸载（删除所有配置和数据）
sudo ./uninstall.sh --purge

# 完全卸载但保留用户上传的文件
sudo ./uninstall.sh --purge --keep-data
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

📖 详细 MCP 使用指南请查看 [docs/MCP_USAGE.md](docs/MCP_USAGE.md)

## 🤝 OpenClaw + httpcat 联动部署

HttpCat 可以与 [OpenClaw](https://clawd.org.cn/)（系统级 AI Agent）配合使用，通过 MCP 协议让 AI 助手在企业微信、QQ、钉钉、飞书等 IM 中直接管理文件。

### 架构概览

```
用户 ──→ IM（企微/QQ/钉钉/飞书）──→ OpenClaw（AI Agent）──MCP──→ httpcat（文件服务）
```

📖 完整的 OpenClaw + httpcat 联动部署教程请查看 [docs/OPENCLAW_HTTPCAT_GUIDE.md](docs/OPENCLAW_HTTPCAT_GUIDE.md)

## 📡 API 接口

### 认证方式

HttpCat 支持两种 API 认证方式：

| 方式 | 适合场景 | 请求头 |
|------|----------|--------|
| JWT Token | Web 前端登录 | `Authorization: Bearer <token>` |
| AK/SK 签名 | 脚本/CI/AI/Open API | `AccessKey` + `Signature` + `TimeStamp` |

### Open API（AK/SK 签名认证）

启用后，所有 `/api/v1/*` 接口均可通过 AK/SK 签名方式调用，无需先登录获取 JWT Token。

#### 1. 启用配置

在 `svr.yml` 中开启：

```yaml
server:
  http:
    auth:
      open_api_enable: true    # 启用 Open API
      aksk:
        your_access_key: your_secret_key   # AK/SK 密钥对
```

#### 2. 签名算法

```
Signature = HMAC-SHA256(
  "{Method}\n{Path}\n{Query}\n{AccessKey}\n{TimeStamp}\n{BodySHA256}",
  SecretKey
)
```

> **分隔符说明**：各字段之间使用**真正的换行符**（`\n` = `0x0a`）分隔，与 AWS Signature V4、腾讯云 TC3-HMAC-SHA256 等行业标准一致。

| 字段 | 说明 |
|------|------|
| Method | HTTP 方法（GET/POST/DELETE 等） |
| Path | 请求路径（末尾不含 `/`） |
| Query | URL Query 参数（`c.Request.URL.RawQuery`） |
| AccessKey | 你的 Access Key |
| TimeStamp | Unix 时间戳（秒），±60 秒有效 |
| BodySHA256 | 请求体的 SHA256 Hex 值（无 Body 时为空字节的 SHA256：`e3b0c44298fc1c14...`） |

**安全机制**：
- 服务端使用**恒定时间比较**（`hmac.Equal`）校验签名，防止时序攻击
- 时间戳窗口 ±60 秒，防止重放攻击

#### 3. 请求头

| Header | 必填 | 说明 |
|--------|------|------|
| `AccessKey` | 是 | Access Key |
| `Signature` | 是 | HMAC-SHA256 签名 |
| `TimeStamp` | 是 | Unix 时间戳（秒） |

#### 4. 支持的接口

启用 Open API 后，以下所有 `/api/v1/*` 接口均可通过 AK/SK 签名调用：

**文件操作**

| 接口 | 方法 | 说明 |
|------|------|------|
| `/api/v1/file/listFiles?dir=/` | GET | 获取目录文件列表 |
| `/api/v1/file/getFileInfo` | GET | 获取文件详细信息 |
| `/api/v1/file/getDirConf` | GET | 获取上传/下载目录配置 |
| `/api/v1/file/uploadHistoryLogs` | GET | 获取上传历史记录 |
| `/api/v1/file/uploadHistoryLogs` | DELETE | 删除上传历史记录 |
| `/api/v1/file/upload` | POST | 上传文件（白名单接口，需 UploadToken 头，详见[上传流程](#7-aksk-场景下上传文件完整流程)） |
| `/api/v1/file/download` | GET | 下载文件（白名单，无需认证） |

**图片管理**

| 接口 | 方法 | 说明 |
|------|------|------|
| `/api/v1/imageManage/upload` | POST | 上传图片 |
| `/api/v1/imageManage/rename` | POST | 图片改名 |
| `/api/v1/imageManage/delete` | DELETE | 删除图片 |
| `/api/v1/imageManage/clear` | DELETE | 清空所有图片 |
| `/api/v1/imageManage/download` | GET | 下载图片 |
| `/api/v1/imageManage/listThumbImages` | GET | 分页获取缩略图列表 |

**统计与监控**

| 接口 | 方法 | 说明 |
|------|------|------|
| `/api/v1/user/dataOverview` | GET | 数据概览（文件数、磁盘使用等） |
| `/api/v1/user/getUploadAvailableSpace` | GET | 可用上传空间 |
| `/api/v1/statistics/getUploadStatistics` | GET | 上传统计 |
| `/api/v1/statistics/getDownloadStatistics` | GET | 下载统计 |

**系统配置**

| 接口 | 方法 | 说明 |
|------|------|------|
| `/api/v1/conf/getVersion` | GET | 获取版本信息 |
| `/api/v1/conf/getConf` | GET | 获取系统配置 |

**用户与 Token 管理**

| 接口 | 方法 | 说明 |
|------|------|------|
| `/api/v1/user/currentUser` | GET | 当前用户信息 |
| `/api/v1/user/changePasswd` | POST | 修改密码 |
| `/api/v1/user/uploadTokenLists` | GET | 上传 Token 列表 |
| `/api/v1/user/createUploadToken` | POST | 创建上传 Token |
| `/api/v1/user/checkUploadToken` | POST | 校验上传 Token |
| `/api/v1/user/removeUploadToken` | DELETE | 删除上传 Token |
| `/api/v1/user/generateAppSecret` | GET | 生成 AppSecret |

#### 5. 使用场景

| 场景 | 说明 |
|------|------|
| **CI/CD 流水线** | 构建完成后自动上传产物、查询文件列表、清理旧版本 |
| **运维脚本** | 批量上传/下载文件、磁盘空间监控、统计数据采集 |
| **AI Agent** | MCP 不可用时，通过 AK/SK + curl 直接操作 httpcat |
| **第三方系统集成** | 其他系统对接 httpcat，无需模拟登录获取 JWT |
| **定时任务/Cron** | 定期拉取统计数据、自动清理过期文件 |
| **命令行工具** | 封装签名逻辑，打造专用 CLI 工具 |

> **AK/SK vs JWT**：AK/SK 是无状态认证——每次请求自包含签名信息，不需要先登录。特别适合非交互式的自动化场景。JWT 需要先登录获取 Token 且有过期时间，适合 Web 前端交互。

#### 6. 调用示例

##### GET 请求示例（列出文件）

**Shell 脚本：**

```bash
#!/bin/bash
AK="your_access_key"
SK="your_secret_key"
HOST="http://localhost:8888"

METHOD="GET"
PATH_URL="/api/v1/file/listFiles"
QUERY="dir=/"
TIMESTAMP=$(date +%s)

# 无 Body 时也计算空字节的 SHA256（与 AWS Signature V4 行业标准一致）
BODY_HASH=$(printf '' | openssl dgst -sha256 -hex | awk '{print $NF}')

# 用真正的换行符 \n 拼接签名字符串
SIGN_STR=$(printf "%s\n%s\n%s\n%s\n%s\n%s" "${METHOD}" "${PATH_URL}" "${QUERY}" "${AK}" "${TIMESTAMP}" "${BODY_HASH}")
SIGNATURE=$(printf '%s' "${SIGN_STR}" | openssl dgst -sha256 -hmac "${SK}" -hex | awk '{print $NF}')

# 发起请求
curl -s "${HOST}${PATH_URL}?${QUERY}" \
  -H "AccessKey: ${AK}" \
  -H "Signature: ${SIGNATURE}" \
  -H "TimeStamp: ${TIMESTAMP}"
```

**Python 脚本：**

```python
import hmac, hashlib, time, requests

AK = "your_access_key"
SK = "your_secret_key"
HOST = "http://localhost:8888"

method = "GET"
path = "/api/v1/file/listFiles"
query = "dir=/"
timestamp = str(int(time.time()))

# 无 Body 时也计算空字节的 SHA256（与 AWS Signature V4 行业标准一致）
body_hash = hashlib.sha256(b"").hexdigest()

# 用真正的换行符拼接
sign_str = f"{method}\n{path}\n{query}\n{AK}\n{timestamp}\n{body_hash}"
signature = hmac.new(SK.encode(), sign_str.encode(), hashlib.sha256).hexdigest()

resp = requests.get(f"{HOST}{path}?{query}", headers={
    "AccessKey": AK,
    "Signature": signature,
    "TimeStamp": timestamp,
})
print(resp.json())
```

##### POST 请求示例（带 JSON Body）

POST 请求需要对**实际的请求体**计算 SHA256，而非空字节：

```bash
#!/bin/bash
AK="your_access_key"
SK="your_secret_key"
HOST="http://localhost:8888"

METHOD="POST"
PATH_URL="/api/v1/user/createUploadToken"
QUERY=""
TIMESTAMP=$(date +%s)
BODY='{"appkey":"httpcat","appsecret":"httpcat_app_secret"}'

# POST 请求：对实际 Body 计算 SHA256
BODY_HASH=$(printf '%s' "${BODY}" | openssl dgst -sha256 -hex | awk '{print $NF}')

SIGN_STR=$(printf "%s\n%s\n%s\n%s\n%s\n%s" "${METHOD}" "${PATH_URL}" "${QUERY}" "${AK}" "${TIMESTAMP}" "${BODY_HASH}")
SIGNATURE=$(printf '%s' "${SIGN_STR}" | openssl dgst -sha256 -hmac "${SK}" -hex | awk '{print $NF}')

curl -s "${HOST}${PATH_URL}" -X POST \
  -H "Content-Type: application/json" \
  -H "AccessKey: ${AK}" \
  -H "Signature: ${SIGNATURE}" \
  -H "TimeStamp: ${TIMESTAMP}" \
  -d "${BODY}"
```

#### 7. AK/SK 场景下上传文件（完整流程）

上传文件需要 **两步操作**，因为 `/api/v1/file/upload` 接口使用独立的 **UploadToken** 机制认证（不走 AK/SK 签名），需要先通过 AK/SK 生成 UploadToken，再用该 Token 上传：

```
┌─────────────┐     AK/SK 签名      ┌──────────┐    UploadToken     ┌──────────┐
│  你的脚本    │ ──────────────────→ │ 生成Token │ ────────────────→ │ 上传文件  │
│             │  createUploadToken  │          │   file/upload     │          │
└─────────────┘                     └──────────┘                   └──────────┘
     Step 1: POST /api/v1/user/createUploadToken (需要 AK/SK 签名)
     Step 2: POST /api/v1/file/upload (需要 UploadToken 请求头)
```

##### 完整 Shell 脚本

```bash
#!/bin/bash
# AK/SK 场景下上传文件的完整流程
AK="your_access_key"
SK="your_secret_key"
HOST="http://localhost:8888"
FILE_PATH="/path/to/your/file.txt"

# ── Step 1: 通过 AK/SK 签名生成 UploadToken ──
METHOD="POST"
PATH_URL="/api/v1/user/createUploadToken"
TIMESTAMP=$(date +%s)
BODY='{"appkey":"httpcat","appsecret":"httpcat_app_secret"}'

BODY_HASH=$(printf '%s' "${BODY}" | openssl dgst -sha256 -hex | awk '{print $NF}')
SIGN_STR=$(printf "%s\n%s\n%s\n%s\n%s\n%s" "${METHOD}" "${PATH_URL}" "" "${AK}" "${TIMESTAMP}" "${BODY_HASH}")
SIGNATURE=$(printf '%s' "${SIGN_STR}" | openssl dgst -sha256 -hmac "${SK}" -hex | awk '{print $NF}')

UPLOAD_TOKEN=$(curl -s "${HOST}${PATH_URL}" -X POST \
  -H "Content-Type: application/json" \
  -H "AccessKey: ${AK}" \
  -H "Signature: ${SIGNATURE}" \
  -H "TimeStamp: ${TIMESTAMP}" \
  -d "${BODY}" | python3 -c "import sys,json; print(json.load(sys.stdin)['data'])")

echo "UploadToken: ${UPLOAD_TOKEN}"

# ── Step 2: 用 UploadToken 上传文件 ──
curl -s "${HOST}/api/v1/file/upload" -X POST \
  -H "UploadToken: ${UPLOAD_TOKEN}" \
  -F "f1=@${FILE_PATH}"
```

##### 完整 Python 脚本

```python
import hmac, hashlib, time, json, requests

AK = "your_access_key"
SK = "your_secret_key"
HOST = "http://localhost:8888"
FILE_PATH = "/path/to/your/file.txt"
# appkey/appsecret 是 httpcat 的应用凭证（svr.yml 中配置的 app_key/app_secret）
APP_KEY = "httpcat"
APP_SECRET = "httpcat_app_secret"

def sign_request(method, path, query, ak, sk, body=b""):
    """计算 AK/SK 签名"""
    timestamp = str(int(time.time()))
    body_hash = hashlib.sha256(body).hexdigest()
    sign_str = f"{method}\n{path}\n{query}\n{ak}\n{timestamp}\n{body_hash}"
    signature = hmac.new(sk.encode(), sign_str.encode(), hashlib.sha256).hexdigest()
    return {"AccessKey": ak, "Signature": signature, "TimeStamp": timestamp}

# Step 1: 生成 UploadToken
body = json.dumps({"appkey": APP_KEY, "appsecret": APP_SECRET}).encode()
headers = sign_request("POST", "/api/v1/user/createUploadToken", "", AK, SK, body)
headers["Content-Type"] = "application/json"

resp = requests.post(f"{HOST}/api/v1/user/createUploadToken", headers=headers, data=body)
upload_token = resp.json()["data"]
print(f"UploadToken: {upload_token}")

# Step 2: 上传文件
with open(FILE_PATH, "rb") as f:
    resp = requests.post(
        f"{HOST}/api/v1/file/upload",
        headers={"UploadToken": upload_token},
        files={"f1": f},
    )
print(resp.json())
```

> **说明**：
> - `appkey` / `appsecret` 是 `svr.yml` 中配置的应用凭证（`app_key` / `app_secret`），用于生成 UploadToken
> - `AccessKey` / `SecretKey` 是 `svr.yml` 中 `aksk` 字段配置的签名密钥对，用于 Open API 认证
> - UploadToken 生成后**不会过期**（当前默认 `deadline=0`），可以缓存复用，无需每次上传都重新生成
> - 上传接口使用 `multipart/form-data` 格式，文件字段名为 `f1`

### 下载文件

```bash
# 下载文件（白名单接口，无需认证）
wget -O filename.jpg http://localhost:8888/api/v1/file/download?filename=filename.jpg
```

### 列出文件

```bash
# 需要 JWT 或 AK/SK 认证
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

# Open API 配置（AK/SK 签名认证）
open_api_enable: false          # 是否启用
aksk:
  your_access_key: your_secret_key

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
sudo find /var/lib/httpcat -name "*.db"
sudo rm /var/lib/httpcat/data/httpcat_sqlite.db
sudo systemctl restart httpcat
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

- **Go 1.23+** - 后端编译
- **Node.js 20+** - 前端构建（推荐 v20.x LTS）
- **npm** - 包管理器（随 Node.js 安装）

> 💡 **提示**: 使用 nvm 的用户可在 `web/` 目录运行 `nvm use` 自动切换到项目指定版本

### 构建命令

```bash
# 交互式构建
./scripts/build.sh

# 构建所有平台（含前端）
./scripts/build.sh -a -f

# 构建指定平台
./scripts/build.sh -p linux_amd64 -f

# 使用 Docker 构建（完整 CGO 支持）
./scripts/build.sh -d -f

# 显示帮助
./scripts/build.sh -h
```

## 📝 许可证

本软件仅供个人使用，禁止用于商业目的。

- 禁止用于商业目的
- 必须保留版权和许可声明
- 本软件按 "原样" 提供，不承担任何保证

## 🌟 参与贡献

欢迎关注我们的 GitHub 项目！⭐ 点亮 star 了解我们的实时动态。

欢迎提出 issue 或提交 pull request。Good luck! 🍀
