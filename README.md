[English](docs/README-en.md) | 简体中文

# 🐱 HttpCat

> 轻量级、高效的 HTTP 文件传输服务，配备现代化 Web 界面和 AI 集成。

HttpCat 是一个基于 HTTP 的文件传输服务，旨在提供简单、高效、稳定的文件上传和下载功能。无论是临时分享还是批量传输文件，HttpCat 都将是你的优秀助手。

## ✨ 功能特点

### 📦 文件传输核心

- 🚀 **大文件上传** - 分片上传 + 断点续传 + 秒传，支持 **100GB** 级超大文件（v0.7.0+）
- ⚡ **高效下载** - HTTP Range 支持，`wget -c` / `curl -C -` 断点续传、浏览器拖拽视频进度（v0.7.0+）
- 📤 **单次上传** - 小文件走一次性 HTTP，简单脚本直接 `curl -F` 即可
- 🗂️ **子目录支持** - 完整的目录树浏览、创建、批量操作（v0.5.0+）
- 🧩 **多文件 ZIP 下载** - 勾选多个文件，服务端即时打包为 ZIP（v0.6.0+）
- 👁️ **文件在线预览** - 文本/代码/Markdown/图片/视频/音频/PDF 直接在浏览器查看（v0.6.0+）
- 🖼️ **图片专用管理** - 独立的图片画廊视图 + 自动缩略图

### 🔗 分享与协作

- 🔗 **文件分享链接** - 一键生成分享链接，支持：
  - ⏳ 有效期（自定义到期时间）
  - 🎯 下载次数限制（达到上限自动失效）
  - 🔑 提取码保护（4 位数字密码）
  - 👁️ 匿名访问开关（可配置分享是否需要登录）
- 📊 **分享管理** - 统一查看/删除所有已创建的分享，实时统计访问次数

### 🤖 AI 与自动化

- 🤖 **MCP 协议支持** - 15 个 MCP Tools，AI 助手（Claude Desktop / Cursor / CodeBuddy / OpenClaw）可直接管理文件
- 🧠 **AI Skill 集成** - 提供符合 [Agent Skills 规范](https://agentskills.io/) 的 Skill 包，一键软链接到 AI IDE
- 🔑 **Open API** - AK/SK 签名认证（AWS V4 风格），脚本/CI/自动化系统无需登录即可调用全部 API
- 🎫 **UploadToken 机制** - 独立的上传凭证，可配置策略（大小限制/有效期/回调通知）

### 🔐 安全与审计

- 🔐 **bcrypt 密码加密** - 行业标准的密码哈希，自动升级旧版 SHA1 哈希
- 🛡️ **登录限流防爆破** - 5 分钟内失败 5 次自动锁定 15 分钟（v0.7.0+）
- 🚧 **路径穿越防护** - 所有文件路径经 `ResolvePathWithinBase` 校验，杜绝越权访问
- 📜 **操作日志审计** - 完整记录上传/下载/删除/分享/登录等操作，支持筛选和统计（v0.6.0+）
- 🔒 **HTTPS/TLS** - 内置 SSL 支持，配置证书即可启用
- 🎭 **JWT + AK/SK 双认证** - Web 端 JWT 会话、API 端 AK/SK 签名，按场景选择

### 📊 监控与管理

- 📊 **数据总览** - 文件总数、磁盘使用、上传/下载趋势图
- 📈 **统计报表** - 上传/下载历史、文件类型分布、Top 文件排行
- 💾 **磁盘监控** - 实时查看可用空间，告警阈值可配
- 📝 **操作日志** - 谁、在什么时间、对哪个文件、做了什么，一目了然

### ⚙️ 部署与运维

- 🐳 **Docker 就绪** - 一键 `docker-compose up`，支持多架构镜像
- 🛠️ **一键安装脚本** - `install.sh` 自动创建 systemd 服务、部署到 `/var/lib/httpcat`
- 🌐 **多平台二进制** - Linux amd64/arm64、macOS arm64、Windows x64 四平台发布
- 📦 **零外部依赖** - 单个二进制 + 内嵌 SQLite，无需 MySQL/Redis/Nginx
- ⚙️ **Web 配置管理** - 浏览器内直接改 `svr.yml`，支持热更新 + 一键重启

### 🎨 用户体验

- 🎨 **现代化界面** - 基于 Ant Design Pro 5.x 的响应式 Web 控制台
- 🌍 **国际化** - 中英文双语切换
- 📱 **移动端友好** - 响应式布局，手机浏览器可正常使用
- 🏠 **快捷上传页** - 首页拖拽即传，无需进入文件管理

## 📁 项目结构

```
httpcat/
├── server-go/                  # 🔧 Go 后端
│   ├── cmd/
│   │   └── httpcat.go          # 应用入口
│   ├── internal/
│   │   ├── common/             # 公共常量、工具函数、初始化
│   │   ├── conf/               # 配置文件（svr.yml）
│   │   ├── handler/            # HTTP 处理器（v1 版本）
│   │   ├── mcp/                # MCP Server 实现（SSE 传输）
│   │   ├── midware/            # 中间件（JWT/AK-SK 认证、指标采集）
│   │   ├── models/             # 数据模型
│   │   ├── p2p/                # P2P 节点发现（实验性，未用于文件传输）
│   │   ├── server/             # 路由注册、核心业务逻辑
│   │   └── storage/            # SQLite 存储层
│   ├── go.mod
│   └── go.sum
│
├── web/                        # 🎨 React 前端（Ant Design Pro + UmiJS）
│   ├── src/
│   │   ├── components/         # 通用组件（Footer、Header 等）
│   │   ├── pages/              # 页面
│   │   │   ├── Welcome.tsx     # 首页（快捷上传）
│   │   │   ├── FileManage/     # 文件管理（列表、图片管理）
│   │   │   ├── ShareManage/    # 分享管理
│   │   │   ├── SharePage/      # 分享访问页面（匿名可访问）
│   │   │   ├── sysInfo/        # 系统信息
│   │   │   ├── SysConfig/      # 系统配置管理
│   │   │   ├── uploadTokenManage/  # Token 管理
│   │   │   └── user/           # 登录、修改密码
│   │   ├── services/           # API 接口定义
│   │   └── locales/            # 国际化（中文/英文）
│   ├── config/                 # UmiJS 路由与构建配置
│   ├── mock/                   # Mock 数据（仅开发环境）
│   └── package.json
│
├── scripts/                    # 🛠️ 运维脚本
│   ├── build.sh                # 多平台交叉编译（支持 Docker 构建）
│   ├── install.sh              # Linux/macOS 一键安装
│   ├── uninstall.sh            # 卸载脚本
│   ├── httpcat-api.sh          # AK/SK 签名调用示例
│   ├── search-and-upload-image.sh  # 图片搜索上传脚本
│   └── translations.sh         # i18n 翻译脚本
│
├── docs/                       # 📚 文档
│   ├── README-en.md            # English README
│   ├── BUILD.md                # 编译构建指南
│   ├── INSTALL.md              # 安装部署指南
│   ├── MCP_USAGE.md            # MCP 集成指南
│   ├── ReleaseNote.md          # 版本发布记录
│   └── ...                     # 其他设计文档
│
├── static/                     # 📦 前端构建产物（不纳入 git，由 npm run build 生成）
├── release/                    # 📤 构建输出（发行包）
├── website/                    # 📂 运行时数据目录
│   ├── upload/                 # 上传文件存储
│   └── download/               # 下载文件存储
├── data/                       # 💾 运行时数据
│   ├── httpcat_sqlite.db       # SQLite 数据库
│   └── chunks/                 # 分片上传临时目录（v0.7.0+，会话过期自动清理）
│
├── Dockerfile                  # Docker 镜像配置
├── docker-compose.yml          # Docker Compose 编排
├── httpcat.service             # systemd 服务文件
└── LANGUAGES                   # 支持的语言列表
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

> ⚠️ **安全提示**: 默认账号仅供初次登录使用，请**立即通过「修改密码」页面修改**。为防止暴力破解，登录接口内置了限流保护：同一 IP 5 分钟内失败 ≥ 5 次将锁定 15 分钟（v0.7.0+）。

## 🎉 生产环境安装

### 快速安装

```bash
# 下载并解压
httpcat_version="v0.7.0"
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
    ├── httpcat_sqlite.db           # SQLite 数据库
    └── chunks/                     # 分片上传临时目录（v0.7.0+）
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

> ⚠️ **v0.4.0 起，MCP 开启时必须配置 `auth_token`**，在 `svr.yml` 中设置：
>
> ```yaml
> server:
>   mcp:
>     enable: true
>     auth_token: "你的安全密码"
> ```

在你的 MCP 客户端配置（Claude Desktop、Cursor、CodeBuddy 等）中添加：

```json
{
  "mcpServers": {
    "httpcat": {
      "type": "sse",
      "url": "http://your-server:8888/mcp/sse",
      "headers": {
        "Authorization": "Bearer 你的安全密码"
      }
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
| `upload_image` | 通过 MCP 上传图片 |
| `create_folder` | 创建文件夹（v0.5.0+） |
| `rename_file` | 重命名文件/文件夹（v0.5.0+） |
| `batch_delete_files` | 批量删除文件（v0.5.0+） |
| `request_delete_file` | 请求删除文件（第一步） |
| `confirm_delete_file` | 确认删除文件（第二步） |
| `get_disk_usage` | 获取磁盘使用情况 |
| `get_statistics` | 获取上传/下载统计 |
| `get_upload_history` | 查询上传历史记录 |
| `get_download_history` | 查询下载历史记录（v0.5.0+） |
| `get_file_overview` | 文件总览统计（v0.5.0+） |
| `verify_file_md5` | 验证文件 MD5 完整性 |

📖 详细 MCP 使用指南请查看 [docs/MCP_USAGE.md](docs/MCP_USAGE.md)

## 🧠 AI Skill（Agent Skills 规范）

HttpCat 提供了符合 [Agent Skills 规范](https://agentskills.io/) 的 Skill 包，可安装到 Claude Code / CodeBuddy / Cursor 等 AI IDE 中，让 AI 助手通过自然语言管理你的文件服务器。

```bash
# 安装到 Claude Code
ln -s /path/to/httpcat/httpcat-skill ~/.claude/skills/httpcat

# 安装到 CodeBuddy
ln -s /path/to/httpcat/httpcat-skill .codebuddy/skills/httpcat

# 安装到 Cursor
ln -s /path/to/httpcat/httpcat-skill .cursor/skills/httpcat
```

安装后即可在 AI 对话中说 "列出 httpcat 上的文件"、"上传文件到服务器"、"查看磁盘使用情况" 等。

📖 详细说明请查看 [httpcat-skill/README.md](httpcat-skill/README.md)

## 🤝 OpenClaw + httpcat 联动部署

HttpCat 可以与 [OpenClaw](https://clawd.org.cn/)（系统级 AI Agent）配合使用，通过 MCP 协议让 AI 助手在企业微信、QQ、钉钉、飞书等 IM 中直接管理文件。

### 架构概览

```
用户 ──→ IM（企微/QQ/钉钉/飞书）──→ OpenClaw（AI Agent）──MCP──→ httpcat（文件服务）
```

📖 完整的 OpenClaw + httpcat 联动部署教程请查看 [docs/OPENCLAW_HTTPCAT_GUIDE.md](docs/OPENCLAW_HTTPCAT_GUIDE.md)

## 📦 大文件上传（v0.7.0+）

HttpCat v0.7.0 新增**分片上传 + 断点续传**能力，解决以下痛点：

- ❌ 1GB+ 大文件单次上传容易因网络抖动失败
- ❌ 上传到 99% 断网，只能从头再来
- ❌ 浏览器关闭就丢失进度
- ❌ 同样的文件重复上传浪费带宽

### 🎯 前端使用：全自动

**浏览器里什么都不用改**：前端会根据文件大小自动选择：

| 文件大小 | 前端策略 | 实际接口 |
|---------|---------|---------|
| **< 10 MB** | 单次上传 | `POST /api/v1/file/upload` |
| **≥ 10 MB** | 分片上传（5MB/片，3 并发） | `POST /upload/init` → `chunk` → `complete` |

打开文件管理页 → 拖拽文件 → 自动走对应模式，进度条实时更新。

### 🔧 脚本/CI 使用：两种模式按需选

**小文件 / 简单场景**（不需要断点续传，接口同 v0.6.0）：

```bash
# 直接单次上传（适合 < 100MB、网络稳定的场景）
curl -X POST http://localhost:8888/api/v1/file/upload \
  -H "UploadToken: httpcat:xxx:xxx" \
  -F "f1=@file.zip" \
  -F "dir=backup/2026"
```

**大文件 / 弱网 / 需要断点续传**：

```bash
#!/bin/bash
# 完整的分片上传脚本示例
HOST="http://localhost:8888"
TOKEN="httpcat:vO9Mt5UtCXWVEaYumi4LxXFImh4=:e30="
FILE="my_big_file.zip"
SIZE=$(stat -c%s "$FILE" 2>/dev/null || stat -f%z "$FILE")
MD5=$(md5sum "$FILE" 2>/dev/null | cut -d' ' -f1 || md5 -q "$FILE")
CHUNK_SIZE=$((5 * 1024 * 1024))  # 5MB
TOTAL=$(( (SIZE + CHUNK_SIZE - 1) / CHUNK_SIZE ))

# Step 1: 初始化会话
RESP=$(curl -s -X POST "$HOST/api/v1/file/upload/init" \
  -H "Content-Type: application/json" -H "UploadToken: $TOKEN" \
  -d "{\"fileName\":\"$FILE\",\"fileSize\":$SIZE,\"chunkSize\":$CHUNK_SIZE,\"fileMD5\":\"$MD5\",\"dir\":\"uploads\",\"overwrite\":true}")
UPLOAD_ID=$(echo "$RESP" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['uploadId'])")
INSTANT=$(echo "$RESP" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['instant'])")

# 命中秒传直接完成
if [ "$INSTANT" = "True" ]; then
  echo "✅ 秒传成功（文件已存在）"
  exit 0
fi

# Step 2: 上传每个分片
for i in $(seq 0 $((TOTAL-1))); do
  dd if="$FILE" of=/tmp/chunk_$i bs=$CHUNK_SIZE skip=$i count=1 2>/dev/null
  curl -s -X POST "$HOST/api/v1/file/upload/chunk" \
    -H "UploadToken: $TOKEN" \
    -F "uploadId=$UPLOAD_ID" \
    -F "chunkIndex=$i" \
    -F "chunk=@/tmp/chunk_$i" > /dev/null
  echo "  chunk $i / $TOTAL"
  rm -f /tmp/chunk_$i
done

# Step 3: 合并分片
curl -s -X POST "$HOST/api/v1/file/upload/complete" \
  -H "Content-Type: application/json" -H "UploadToken: $TOKEN" \
  -d "{\"uploadId\":\"$UPLOAD_ID\"}"
```

**Python 版本（含断点续传）**：

```python
import hashlib, requests, os, json

HOST = "http://localhost:8888"
TOKEN = "httpcat:xxx:xxx"
FILE = "my_big_file.zip"
CHUNK = 5 * 1024 * 1024

size = os.path.getsize(FILE)
total = (size + CHUNK - 1) // CHUNK
md5 = hashlib.md5(open(FILE, "rb").read()).hexdigest()

# 初始化会话
r = requests.post(f"{HOST}/api/v1/file/upload/init",
    headers={"UploadToken": TOKEN, "Content-Type": "application/json"},
    data=json.dumps({
        "fileName": os.path.basename(FILE),
        "fileSize": size, "chunkSize": CHUNK,
        "fileMD5": md5, "dir": "uploads", "overwrite": True
    })).json()

if r["data"]["instant"]:
    print("✅ 秒传命中"); exit()

uid = r["data"]["uploadId"]

# 【断点续传】查询已上传分片，跳过它们
status = requests.get(f"{HOST}/api/v1/file/upload/status",
    headers={"UploadToken": TOKEN}, params={"uploadId": uid}).json()
uploaded = set(status["data"]["uploadedIdx"])

# 上传缺失分片
with open(FILE, "rb") as f:
    for i in range(total):
        if i in uploaded:
            continue
        f.seek(i * CHUNK)
        data = f.read(CHUNK)
        requests.post(f"{HOST}/api/v1/file/upload/chunk",
            headers={"UploadToken": TOKEN},
            data={"uploadId": uid, "chunkIndex": i},
            files={"chunk": data})
        print(f"  chunk {i+1}/{total}")

# 合并
r = requests.post(f"{HOST}/api/v1/file/upload/complete",
    headers={"UploadToken": TOKEN, "Content-Type": "application/json"},
    data=json.dumps({"uploadId": uid})).json()
print(f"✅ 完成：{r['data']}")
```

### 🎁 附加能力

#### 秒传（Instant Upload）

若文件已存在于服务器且 MD5 相同，`init` 接口会直接返回 `instant: true`，服务端通过**硬链接**瞬间生成新文件，无需重传：

```json
{
  "data": {
    "uploadId": "instant-abc123...",
    "instant": true,
    "totalChunks": 0,
    "uploadedIdx": []
  }
}
```

#### 断点续传（Resume）

任何时候调用 `GET /upload/status?uploadId=xxx` 可获得：
- `uploadedIdx`：已上传分片索引
- `missingIdx`：缺失分片索引

客户端只需重传 `missingIdx` 中的分片。**即使服务端重启**，会话状态也能从 SQLite 完整恢复。

#### 下载断点续传

`GET /api/v1/file/download` 自动支持 HTTP Range：

```bash
# wget 断点续传
wget -c "http://localhost:8888/api/v1/file/download?filename=big.zip"

# curl 断点续传
curl -C - -o big.zip "http://localhost:8888/api/v1/file/download?filename=big.zip"

# 4 段并行加速下载（自行合并）
for i in 0 1 2 3; do
  START=$((i * 25000000)); END=$((START + 24999999))
  curl -H "Range: bytes=$START-$END" -o part$i "..." &
done; wait
cat part{0,1,2,3} > big.zip
```

### 📊 参数约束

| 参数 | 默认 | 最小 | 最大 | 说明 |
|------|------|------|------|------|
| `chunkSize` | 5 MB | 64 KB | 100 MB | 单分片大小 |
| `fileSize` | - | 1 字节 | **100 GB** | 单文件总大小 |
| 会话有效期 | 24 小时 | - | - | 过期后分片自动清理 |
| 前端分片阈值 | 10 MB | - | - | 可在 `FileList/index.tsx` 的 `CHUNK_THRESHOLD` 调整 |
| 前端并发数 | 3 | 1 | - | 可在 `chunkedUpload()` 的 `concurrent` 参数调整 |

### 🛡️ 安全机制

- ✅ 路径穿越防护：所有文件路径都经过 `ResolvePathWithinBase` 校验
- ✅ 分片大小校验：除最后一片外必须精确等于 `chunkSize`
- ✅ 可选 `chunkMD5` 字段：服务端校验单片完整性
- ✅ 整体 MD5 校验：`complete` 时若客户端声明了 `fileMD5`，服务端合并后必须一致
- ✅ 失败不留残渣：合并过程中任何错误都会清理临时文件

> 📖 **想深入了解底层实现？** 请阅读 [**分片上传、断点续传、秒传原理详解**](docs/CHUNK_UPLOAD_PRINCIPLE.md)，包含完整的工程设计、安全机制、性能对比和业界方案对比。

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
| `/api/v1/file/delete` | POST | 批量删除文件/文件夹（v0.5.0+） |
| `/api/v1/file/mkdir` | POST | 创建文件夹（v0.5.0+） |
| `/api/v1/file/rename` | POST | 重命名文件/文件夹（v0.5.0+） |
| `/api/v1/file/upload/init` | POST | 初始化分片上传会话（v0.7.0+） |
| `/api/v1/file/upload/status` | GET | 查询分片上传状态（断点续传用，v0.7.0+） |
| `/api/v1/file/upload/chunk` | POST | 上传单个分片（v0.7.0+） |
| `/api/v1/file/upload/complete` | POST | 合并分片为最终文件（v0.7.0+） |
| `/api/v1/file/upload/abort` | POST | 中止分片上传会话（v0.7.0+） |

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
| `/api/v1/statistics/getFileOverview` | GET | 文件总览统计（v0.5.0+） |
| `/api/v1/statistics/downloadHistoryLogs` | GET | 下载历史日志（v0.5.0+） |

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

配置文件：`conf/svr.yml`（Docker 镜像内）或 `/etc/httpcat/svr.yml`（install.sh 安装）

```yaml
server:
  http:
    port: 8888
    # 大文件上传需要较长的读写超时（单位：秒）
    read_timeout: 1800
    write_timeout: 1800
    auth:
      open_api_enable: false       # 是否启用 Open API（AK/SK 签名认证）
      aksk:                        # AK/SK 凭证对
        your_access_key: your_secret_key
    file:
      upload_enable: true
      enable_upload_token: true    # 是否开启 UploadToken 验证
      app_key: "httpcat"           # 生成 UploadToken 的 app_key
      app_secret: "httpcat_app_secret"
      upload_policy:
        deadline: 7200             # UploadToken 有效期（秒）
        fsizeLimit: 0              # 单文件大小上限（字节，0=不限）
      download_dir: "website/upload/"
      enable_sqlite: true
      sqlite_db_path: "./data/httpcat_sqlite.db"

  mcp:
    enable: true                   # 启用 MCP Server（AI 助手接入）
    auth_token: "替换为安全密码"    # MCP Bearer Token

  share:
    enable: true                   # 启用文件分享
    anonymous_access: true         # 分享链接是否允许匿名访问
```

> 💡 **v0.7.0 分片上传参数**：默认 5MB/片、24 小时会话有效期、100GB 单文件上限，均为代码默认值，无需配置。前端分片阈值（10MB）可在 `web/src/pages/FileManage/FileList/index.tsx` 的 `CHUNK_THRESHOLD` 调整。

## 🍀 常见问题

### 忘记密码？

**方案 A：只重置 admin 密码（保留所有数据，推荐）**

需要 Python 环境和 bcrypt 库：

```bash
# 安装依赖（只需一次）
pip3 install bcrypt

# 重置 admin 密码为 admin123（可改）
sudo systemctl stop httpcat
sudo python3 -c "
import bcrypt, sqlite3, time
pwd = b'admin123'  # 改成你想要的新密码
h = bcrypt.hashpw(pwd, bcrypt.gensalt(10)).decode()
c = sqlite3.connect('/var/lib/httpcat/data/httpcat_sqlite.db')
c.execute(\"UPDATE users SET password=?, salt='', password_update_time=? WHERE username='admin'\", (h, int(time.time())))
c.commit(); c.close()
print('✅ 密码已重置')
"
sudo systemctl start httpcat
```

**方案 B：重置数据库（会清空所有数据）**

```bash
sudo systemctl stop httpcat
sudo rm /var/lib/httpcat/data/httpcat_sqlite.db
sudo systemctl start httpcat
# 重启后会自动创建默认管理员账号 admin / admin
```

### 登录提示 "too many failed attempts, please try again later"？

这是 v0.7.0 新增的**登录限流防爆破**机制：同一 IP 5 分钟内失败 ≥ 5 次会锁定 15 分钟。等待返回的 `lockedRemainingSeconds` 秒即可，或**重启服务**立即清除内存中的限流状态：

```bash
sudo systemctl restart httpcat
```

### 分片上传失败，会残留临时文件吗？

不会。有三重保障：
1. 未完成的分片写入时使用 `.part` 临时后缀，失败不会污染 bitmap
2. 会话默认 **24 小时**过期，后台每 30 分钟扫描清理 `data/chunks/{uploadId}/`
3. 也可调用 `POST /api/v1/file/upload/abort` 主动中止并清理

查看当前残留：`ls /var/lib/httpcat/data/chunks/`

### 大文件上传后，浏览器没有进度条？

确认文件大小 **≥ 10 MB**：前端代码里 `CHUNK_THRESHOLD = 10MB`，小于此值会走旧版单次上传（没有精细进度）。可在 `web/src/pages/FileManage/FileList/index.tsx` 修改阈值。

### Node.js 版本问题？

本项目使用 UmiJS 3.x，需要在 Node.js 17+ 环境启用旧版 OpenSSL provider：

```bash
NODE_OPTIONS=--openssl-legacy-provider npm run start:dev
```

> 推荐使用 **Node.js v20.x LTS**（构建脚本已处理兼容）。使用 nvm 的用户可在 `web/` 目录运行 `nvm use` 自动切换到项目指定版本。

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

## 📅 版本演进

| 版本 | 主题 | 核心能力 |
|------|------|---------|
| **v0.7.0** | 📦 大文件 & 安全 | 分片上传 + 断点续传 + 秒传 + 下载 Range + 登录限流 |
| v0.6.0 | 🔍 审计 & 体验 | 操作日志、文件在线预览（文本/图片/视频/PDF）、多文件 ZIP 下载 |
| v0.5.0 | 📂 深度文件管理 | 子目录、文件总览、批量操作、15 个 MCP Tools |
| v0.4.0 | 🔒 安全 & 分享 | bcrypt 密码加密、文件分享（有效期/下载次数/提取码） |
| v0.3.0 | 🌐 Web 自治 | 浏览器内系统配置管理 |
| v0.2.x | 🤖 AI 集成 | MCP 协议、Docker 镜像、AK/SK 签名认证 |
| v0.1.x | 🎯 基础 | 上传下载、SQLite、Web 界面 |

> 📖 详细变更请查看 [CHANGELOG.md](CHANGELOG.md) 和 [docs/ReleaseNote.md](docs/ReleaseNote.md)

## 📝 许可证

本软件仅供个人使用，禁止用于商业目的。

- 禁止用于商业目的
- 必须保留版权和许可声明
- 本软件按 "原样" 提供，不承担任何保证

## 🌟 参与贡献

欢迎关注我们的 GitHub 项目！⭐ 点亮 star 了解我们的实时动态。

欢迎提出 issue 或提交 pull request。Good luck! 🍀
