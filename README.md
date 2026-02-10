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

### 腾讯云 Lighthouse 部署

#### 1. 部署 httpcat

```bash
mkdir -p /data/httpcat
docker run -d \
  --name httpcat \
  --restart unless-stopped \
  -p 8888:8888 \
  -v /data/httpcat:/app/uploads \
  dockershe/httpcat:latest
```

#### 2. 部署 OpenClaw

```bash
mkdir -p ~/openclaw-docker/data/.openclaw ~/openclaw-docker/data/clawd
cd ~/openclaw-docker

# 生成 Token
GATEWAY_TOKEN=$(openssl rand -hex 16)
echo "你的 Token: $GATEWAY_TOKEN"

# 创建环境变量
cat > .env << EOF
OPENCLAW_IMAGE=jiulingyun803/openclaw-cn:latest
OPENCLAW_CONFIG_DIR=./data/.openclaw
OPENCLAW_WORKSPACE_DIR=./data/clawd
OPENCLAW_GATEWAY_PORT=18789
OPENCLAW_BRIDGE_PORT=18790
OPENCLAW_GATEWAY_TOKEN=$GATEWAY_TOKEN
EOF

# 创建 docker-compose.yml
cat > docker-compose.yml << 'YAML'
services:
  openclaw-gateway:
    image: ${OPENCLAW_IMAGE:-jiulingyun803/openclaw-cn:latest}
    user: node:node
    environment:
      HOME: /home/node
      TERM: xterm-256color
      OPENCLAW_GATEWAY_TOKEN: ${OPENCLAW_GATEWAY_TOKEN}
    volumes:
      - ${OPENCLAW_CONFIG_DIR:-./data/.openclaw}:/home/node/.openclaw
      - ${OPENCLAW_WORKSPACE_DIR:-./data/clawd}:/home/node/clawd
    ports:
      - "${OPENCLAW_GATEWAY_PORT:-18789}:18789"
      - "${OPENCLAW_BRIDGE_PORT:-18790}:18790"
    init: true
    restart: unless-stopped
    command: ["node", "dist/index.js", "gateway", "--bind", "lan", "--port", "18789", "--allow-unconfigured"]
  openclaw-cli:
    image: ${OPENCLAW_IMAGE:-jiulingyun803/openclaw-cn:latest}
    user: node:node
    environment:
      HOME: /home/node
      TERM: xterm-256color
      BROWSER: echo
    volumes:
      - ${OPENCLAW_CONFIG_DIR:-./data/.openclaw}:/home/node/.openclaw
      - ${OPENCLAW_WORKSPACE_DIR:-./data/clawd}:/home/node/clawd
    stdin_open: true
    tty: true
    init: true
    entrypoint: ["node", "dist/index.js"]
YAML

# 关键：修复目录权限（容器以 node uid=1000 运行）
chown -R 1000:1000 ./data/.openclaw ./data/clawd

# 拉取镜像并启动
docker compose pull
docker compose up -d openclaw-gateway
```

#### 3. 配置 httpcat MCP 连接

```bash
cat > ~/openclaw-docker/data/.openclaw/mcp.json << 'EOF'
{
  "mcpServers": {
    "httpcat": {
      "type": "sse",
      "url": "http://172.17.0.1:8888/mcp/sse"
    }
  }
}
EOF
cd ~/openclaw-docker && docker compose restart openclaw-gateway
```

#### 4. 配置 AI 模型

Gateway 启动后需要配置 AI 模型 API Key（DeepSeek / OpenAI / Claude 等），否则 Web UI 会显示"健康状态 离线"：

- **方式 A（Web UI）**：浏览器打开 `http://公网IP:18789` → 所有设置 → 配置 AI Provider
- **方式 B（SSH）**：`docker compose run --rm openclaw-cli onboard`

#### 5. 放通防火墙端口

| 端口 | 用途 |
|------|------|
| 443 | OpenClaw Web UI（HTTPS，Nginx 反向代理） |
| 8888 | httpcat |
| 18790 | OpenClaw Bridge |

### 常见踩坑

| 问题 | 原因 | 解决 |
|------|------|------|
| Web UI 显示"健康状态 离线" | 未配置 AI 模型 API Key | 在设置中配置 AI Provider 或运行 `onboard` |
| `EACCES: permission denied` | 数据目录属主不对 | `chown -R 1000:1000 ./data/.openclaw ./data/clawd` |
| `Invalid --bind` | `--bind` 不支持 `0.0.0.0` | 改为 `--bind lan` |
| `Missing config` 反复重启 | 未配置时 Gateway 无法启动 | 添加 `--allow-unconfigured` 参数 |
| `disconnected (1006): no reason` | 公网 IP 访问时 WebSocket Origin 检查返回 403 | 使用 Nginx 反向代理（见下方 HTTPS 部署） |
| `control ui requires HTTPS or localhost` | OpenClaw 要求安全上下文（HTTPS 或 localhost） | **必须通过 Nginx + HTTPS 反向代理访问**（见下方） |
| `unauthorized: gateway token missing` | Nginx 代理请求被视为不可信，Token 未传递 | 配置 `gateway.trustedProxies`（见下方） |
| `pairing required` | 远程设备首次连接需要配对认证 | 配置 `dangerouslyDisableDeviceAuth`（见下方） |

### HTTPS 反向代理部署（公网访问必需）

OpenClaw 的 Control UI 要求 **HTTPS 安全上下文**（或 localhost）。通过公网 IP 的 HTTP 直接访问会报错 `disconnected (1008): control ui requires HTTPS or localhost (secure context)`。此外，WebSocket 有 Origin 同源检查（CVE 安全补丁 GHSA-g8p2-7wf7-98mq），非 localhost Origin 会返回 403。

**解决方案**：使用 Nginx + 自签名 HTTPS 证书反向代理：

```bash
# 1. 生成自签名证书（10 年有效期）
mkdir -p /etc/nginx/ssl
openssl req -x509 -nodes -days 3650 -newkey rsa:2048 \
  -keyout /etc/nginx/ssl/openclaw.key \
  -out /etc/nginx/ssl/openclaw.crt \
  -subj '/CN=你的公网IP/O=OpenClaw/C=CN'

# 2. 创建 Nginx 配置
cat > /etc/nginx/sites-available/openclaw << 'NGINX'
server {
    listen 443 ssl;
    server_name _;

    ssl_certificate /etc/nginx/ssl/openclaw.crt;
    ssl_certificate_key /etc/nginx/ssl/openclaw.key;
    ssl_protocols TLSv1.2 TLSv1.3;

    location / {
        proxy_pass http://127.0.0.1:18789;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_read_timeout 86400;
        proxy_send_timeout 86400;
    }
}
NGINX

# 3. 启用配置并启动 Nginx
ln -sf /etc/nginx/sites-available/openclaw /etc/nginx/sites-enabled/openclaw
nginx -t && systemctl start nginx && systemctl enable nginx
```

配置完成后通过 **`https://你的公网IP`** 访问（浏览器提示不安全时点击"高级 → 继续前往"）。

> **端口说明**：还需在 Lighthouse 防火墙中放通 **443** 端口。

### 配置 trustedProxies（解决 token missing）

通过 Nginx 反向代理后，Gateway 会检测到 `X-Forwarded-For` 等代理头。如果代理 IP 不在信任列表中，请求会被视为**不可信**，导致 Token 无法传递，报错 `unauthorized: gateway token missing`。

```bash
cd ~/openclaw-docker

# 添加 Docker 网关 IP 到信任代理列表
docker compose run --rm openclaw-cli config set \
  gateway.trustedProxies '["172.19.0.1","172.17.0.1","127.0.0.1"]'

# 重启生效
docker compose restart openclaw-gateway
```

> **如何确定 Docker 网关 IP**：运行 `docker network inspect openclaw-docker_default | grep Gateway`，通常是 `172.19.0.1` 或 `172.18.0.1`。

### 配置 dangerouslyDisableDeviceAuth（解决 pairing required）

OpenClaw 有**设备配对认证**机制：每个新设备（浏览器）首次连接时需要被批准。本地连接（localhost）会自动批准，但通过 Nginx 代理的远程连接不会——即使 trustedProxies 配置正确，真实客户端 IP 仍然不是本地地址，`isLocalClient = false`，配对不会静默批准。

原理：`isLocalClient = !hasUntrustedProxyHeaders && isLocalGatewayAddress(clientIp)`，代理可信后 clientIp 会解析为真实 IP（如 `111.206.96.148`），不是本地地址。

最简单的解决方案——**禁用设备认证**：

```bash
cd ~/openclaw-docker

# 编辑配置（Python 方式）
python3 << 'PYEOF'
import json
with open('./data/.openclaw/openclaw.json', 'r') as f:
    cfg = json.load(f)
if 'controlUi' not in cfg['gateway']:
    cfg['gateway']['controlUi'] = {}
cfg['gateway']['controlUi']['dangerouslyDisableDeviceAuth'] = True
with open('./data/.openclaw/openclaw.json', 'w') as f:
    json.dump(cfg, f, indent=2)
print('Done! dangerouslyDisableDeviceAuth = true')
PYEOF

# 重启生效
docker compose restart openclaw-gateway
```

> **安全提示**：此选项名含 `dangerously` 前缀，意味着禁用后任何知道 Token 的人都能直接访问 Control UI，无需设备配对。在私有服务器上可以接受，但公网暴露时请确保 Token 足够复杂。

### 配置 AI 模型（dashscope 为例）

OpenClaw 的模型配置涉及**两处**，都在 `openclaw.json` 中：

| 配置项 | 作用 | 示例值 |
|--------|------|--------|
| `models.providers.dashscope.models[]` | 定义可用模型列表（ID、上下文窗口、价格等） | `{"id": "qwen3-max", ...}` |
| `agents.defaults.model.primary` | Agent 实际使用的默认模型 | `"dashscope/qwen3-max"` |
| `agents.defaults.models` | Agent 模型别名映射 | `{"dashscope/qwen3-max": {"alias": "Qwen3 Max"}}` |

> **踩坑**：只改 `models.providers` 不改 `agents.defaults.model.primary`，日志里模型不会变！两处必须同步修改。

#### dashscope 模型选择建议

| 模型 ID | 定位 | 适合 OpenClaw Agent？ |
|---------|------|----------------------|
| `qwen3-max` | 文本生成旗舰，能力最强 | **首选** — 编程+对话效果最好，上下文 262K |
| `qwen-coder-plus-latest` | 代码专用模型 | 适合纯编程场景，上下文 1M |
| `qwen-plus-latest` | 效果/速度/成本均衡 | 性价比之选 |
| `qwen-flash` | 速度快、成本低 | 简单任务省钱用 |
| `qwen3-max-2026-01-23` | 深度思考（reasoning）模型 | **不推荐** — OpenClaw 的 `reasoning:true` 会发送不兼容参数导致空回复 |

> **关键**：使用 dashscope 模型时，`reasoning` 必须设为 `false`。OpenClaw 的 reasoning 模式会发送 Claude 专用的 `thinking`/`budget_tokens` 参数，dashscope 的 OpenAI 兼容接口不支持，会导致 AI 静默返回空内容。

#### 切换模型的完整命令

```bash
cd ~/openclaw-docker

python3 << 'PYEOF'
import json
with open('./data/.openclaw/openclaw.json', 'r') as f:
    cfg = json.load(f)

# 1. 更新模型列表
cfg['models']['providers']['dashscope']['models'][0] = {
    "id": "qwen3-max",
    "name": "Qwen3 Max",
    "reasoning": False,          # 必须 False！
    "input": ["text"],
    "cost": {"input": 0.002, "output": 0.008, "cacheRead": 0, "cacheWrite": 0},
    "contextWindow": 131072,
    "maxTokens": 65536
}

# 2. 同步更新 Agent 默认模型（必须！）
cfg['agents']['defaults']['model']['primary'] = 'dashscope/qwen3-max'
cfg['agents']['defaults']['models'] = {
    'dashscope/qwen3-max': {'alias': 'Qwen3 Max'}
}

with open('./data/.openclaw/openclaw.json', 'w') as f:
    json.dump(cfg, f, indent=2)
print('Done!')
PYEOF

# 重启生效
docker compose restart openclaw-gateway
```

### 配置 SOUL.md（自定义 AI 行为和语言）

OpenClaw 通过工作区中的 **`SOUL.md`** 文件定义 AI 的人格、行为风格和语言偏好。这是 Agent 的"灵魂文件"，每次新对话开始时自动读取。

**文件位置**：容器内 `/home/node/openclaw/SOUL.md`（即 `agents.defaults.workspace` 配置的目录下）

#### 为什么默认回复英文？

OpenClaw 内置的 SOUL.md 是全英文的，没有语言指令，而底层模型（qwen3-max）看到英文系统提示后会默认用英文回复。

#### 解决方案：在 SOUL.md 开头加中文语言规则

```bash
docker exec openclaw-docker-openclaw-gateway-1 python3 -c "
content = '''# SOUL.md - Who You Are

## Language Rule (HIGHEST PRIORITY)

**Always respond in Simplified Chinese.** Unless the user explicitly asks in English or requests English replies, all responses must be in Chinese. Technical terms may remain in English.

## Core Truths

Be genuinely helpful, not performatively helpful. Skip filler words, just help.
Have opinions. Be resourceful before asking. Earn trust through competence.

## Vibe

Concise when needed, thorough when it matters. Not a corporate drone. Just good.

## Continuity

Each session, you wake up fresh. These files are your memory. Read and update them.
'''
with open('/home/node/openclaw/SOUL.md', 'w') as f:
    f.write(content)
print('Done!')
"
```

> **注意**：修改 SOUL.md **不需要重启 Gateway**，开一个新对话即可生效。

#### SOUL.md 的其他工作区文件

| 文件 | 作用 |
|------|------|
| `SOUL.md` | AI 人格和语言偏好（**核心文件**） |
| `IDENTITY.md` | AI 的身份信息 |
| `USER.md` | 用户信息（AI 会记住你是谁） |
| `AGENTS.md` | 子 Agent 配置 |
| `TOOLS.md` | 可用工具说明 |
| `BOOTSTRAP.md` | 首次启动引导 |
| `HEARTBEAT.md` | 心跳/保活配置 |

这些文件都在容器内 `/home/node/openclaw/` 目录下，AI 每次醒来会读取它们作为记忆。

📖 完整的 OpenClaw + httpcat 联动教程请查看 [docs/OPENCLAW_HTTPCAT_GUIDE.md](docs/OPENCLAW_HTTPCAT_GUIDE.md)

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
