English | [简体中文](../README.md)

# 🐱 HttpCat

> A lightweight, efficient HTTP file transfer service with modern web interface and AI integration.

HttpCat is designed to provide a simple, efficient, and stable solution for file uploading and downloading. Whether it's for temporary sharing or bulk file transfers, HttpCat will be your excellent assistant.

## ✨ Key Features

### 📦 Core File Transfer

- 🚀 **Large File Upload** - Chunked upload + resume + instant-upload, supports files up to **100GB** (v0.7.0+)
- ⚡ **Efficient Download** - HTTP Range support, `wget -c` / `curl -C -` resume, video seek in browser (v0.7.0+)
- 📤 **Simple Upload** - Small files use single-request upload, `curl -F` works out of the box
- 🗂️ **Subdirectory Support** - Full directory tree: browse, create, batch operations (v0.5.0+)
- 🧩 **Multi-file ZIP Download** - Select multiple files, server packages them into ZIP on-the-fly (v0.6.0+)
- 👁️ **Online File Preview** - Text/code/Markdown/image/video/audio/PDF viewable directly in browser (v0.6.0+)
- 🖼️ **Dedicated Image Management** - Standalone gallery view with auto-generated thumbnails

### 🔗 Sharing & Collaboration

- 🔗 **File Share Links** - One-click share links, supporting:
  - ⏳ Expiration time (custom)
  - 🎯 Download count limits (auto-expire when reached)
  - 🔑 Access code protection (4-digit password)
  - 👁️ Anonymous access toggle (configurable)
- 📊 **Share Management** - Unified view/delete of all shares with real-time access stats

### 🤖 AI & Automation

- 🤖 **MCP Protocol Support** - 15 MCP Tools let AI assistants (Claude Desktop / Cursor / CodeBuddy / OpenClaw) manage files directly
- 🧠 **AI Skill Integration** - Skill package following [Agent Skills spec](https://agentskills.io/), symlink to your AI IDE
- 🔑 **Open API** - AK/SK signature auth (AWS V4 style), scripts/CI/automation can call any API without login
- 🎫 **UploadToken Mechanism** - Standalone upload credentials with configurable policies (size/expiry/callback)

### 🔐 Security & Audit

- 🔐 **bcrypt Password Hashing** - Industry-standard password hash, auto-upgrades legacy SHA1 hashes
- 🛡️ **Login Rate-Limit Protection** - 5 failures in 5 min → 15 min lockout (v0.7.0+)
- 🚧 **Path Traversal Protection** - All paths validated by `ResolvePathWithinBase`, prevents unauthorized access
- 📜 **Operation Audit Log** - Full audit trail for upload/download/delete/share/login operations (v0.6.0+)
- 🔒 **HTTPS/TLS** - Built-in SSL support, just provide certs to enable
- 🎭 **JWT + AK/SK Dual Auth** - JWT sessions for web UI, AK/SK signatures for API calls

### 📊 Monitoring & Management

- 📊 **Data Overview** - Total files, disk usage, upload/download trends
- 📈 **Statistics** - Upload/download history, file type distribution, top files
- 💾 **Disk Monitoring** - Real-time free space, configurable alert thresholds
- 📝 **Operation Log** - Who, when, which file, and what operation — at a glance

### ⚙️ Deployment & Operations

- 🐳 **Docker Ready** - One-command `docker-compose up`, multi-arch images
- 🛠️ **One-click Install Script** - `install.sh` auto-creates systemd service, deploys to `/var/lib/httpcat`
- 🌐 **Multi-platform Binaries** - Linux amd64/arm64, macOS arm64, Windows x64 — four platforms
- 📦 **Zero External Dependencies** - Single binary + embedded SQLite, no MySQL/Redis/Nginx needed
- ⚙️ **Web Config Management** - Edit `svr.yml` in browser, hot-reload + one-click restart

### 🎨 User Experience

- 🎨 **Modern UI** - Responsive web console based on Ant Design Pro 5.x
- 🌍 **Internationalization** - Chinese/English bilingual switcher
- 📱 **Mobile-friendly** - Responsive layout, works on mobile browsers
- 🏠 **Quick Upload Homepage** - Drag-and-drop on homepage without entering file management

## 📁 Project Structure

```
httpcat/
├── server-go/                  # 🔧 Go Backend
│   ├── cmd/
│   │   └── httpcat.go          # Application entry
│   ├── internal/
│   │   ├── common/             # Shared utils, constants, initialization
│   │   ├── conf/               # Config (svr.yml)
│   │   ├── handler/            # HTTP handlers (v1)
│   │   ├── mcp/                # MCP Server (SSE transport)
│   │   ├── midware/            # Middleware (JWT/AK-SK auth, metrics, rate-limit)
│   │   ├── models/             # Data models
│   │   ├── p2p/                # P2P node discovery (experimental, not used for file transfer yet)
│   │   ├── server/             # Route registration, core logic
│   │   └── storage/            # SQLite storage layer
│   ├── go.mod
│   └── go.sum
│
├── web/                        # 🎨 React Frontend (Ant Design Pro + UmiJS)
│   ├── src/
│   │   ├── components/         # Shared components
│   │   ├── pages/              # Pages
│   │   │   ├── Welcome.tsx     # Homepage (quick upload)
│   │   │   ├── FileManage/     # File management (list, images)
│   │   │   ├── ShareManage/    # Share management
│   │   │   ├── SharePage/      # Share access page (anonymous)
│   │   │   ├── OperationLog/   # Operation audit log
│   │   │   ├── sysInfo/        # System info
│   │   │   ├── SysConfig/      # System configuration
│   │   │   ├── uploadTokenManage/  # Token management
│   │   │   └── user/           # Login, change password
│   │   ├── services/           # API definitions
│   │   └── locales/            # i18n (Chinese/English)
│   ├── config/                 # UmiJS routes & build config
│   ├── mock/                   # Mock data (dev only)
│   └── package.json
│
├── scripts/                    # 🛠️ Ops scripts
│   ├── build.sh                # Multi-platform cross-compile (Docker supported)
│   ├── install.sh              # Linux/macOS one-click install
│   ├── uninstall.sh            # Uninstall script
│   ├── httpcat-api.sh          # AK/SK signature example
│   ├── test-v070.sh            # v0.7.0 integration test suite
│   └── translations.sh         # i18n script
│
├── docs/                       # 📚 Documentation
│   ├── README-en.md            # English README
│   ├── BUILD.md                # Build guide
│   ├── INSTALL.md              # Install guide
│   ├── MCP_USAGE.md            # MCP integration guide
│   ├── CHUNK_UPLOAD_PRINCIPLE.md  # Chunked upload principle (v0.7.0)
│   ├── ReleaseNote.md          # Release history
│   └── ...                     # Other design docs
│
├── static/                     # 📦 Frontend build output (not in git)
├── release/                    # 📤 Build artifacts
├── website/                    # 📂 Runtime data
│   ├── upload/                 # Uploaded files
│   └── download/               # Download files
├── data/                       # 💾 Runtime data
│   ├── httpcat_sqlite.db       # SQLite database
│   └── chunks/                 # Chunked upload temp dir (v0.7.0+, auto-cleaned)
│
├── Dockerfile                  # Docker image config
├── docker-compose.yml          # Docker Compose orchestration
├── httpcat.service             # systemd service file
└── LANGUAGES                   # Supported languages
```

## 🚀 Quick Start

### Option 1: Docker (Recommended)

```bash
# Using Docker Compose
docker-compose up -d

# Or using Docker directly
docker run -d --name httpcat \
  -p 8888:8888 \
  -v $(pwd)/data:/app/data \
  -v $(pwd)/upload:/app/upload \
  httpcat:latest
```

### Option 2: Build from Source

```bash
# Build everything (backend + frontend)
./scripts/build.sh -a -f

# Or build separately:

# Backend only
cd server-go && go build -o httpcat ./cmd/httpcat.go

# Frontend only
cd web && npm install && npm run build
```

### Option 3: Development Mode

```bash
# Terminal 1: Start backend
cd server-go
go build -o httpcat ./cmd/httpcat.go
./httpcat -C ./internal/conf/svr.yml --static=../static/

# Terminal 2: Start frontend dev server
cd web
npm install --registry=https://registry.npmmirror.com
NODE_OPTIONS=--openssl-legacy-provider npm run start:dev
```

Access URLs:
- **Frontend**: http://localhost:8000 (dev) or http://localhost:8888 (prod)
- **Backend API**: http://localhost:8888/api/v1/

### Default Account

| Field | Value |
|-------|-------|
| Username | `admin` |
| Password | `admin` |

> ⚠️ **Security Note**: The default account is only for initial login — please change the password immediately via the "Change Password" page. To prevent brute-force attacks, the login endpoint has built-in rate-limiting: 5 failures in 5 min from the same IP → 15 min lockout (v0.7.0+).

## 🎉 Production Installation

### Quick Install

```bash
# Download and extract
httpcat_version="v0.7.0"
tar -zxvf httpcat_${httpcat_version}_linux-amd64.tar.gz
cd httpcat_${httpcat_version}_linux-amd64

# Install (interactive)
sudo ./install.sh

# Or install with custom port
sudo ./install.sh -p 9000

# Manage service
sudo systemctl start httpcat
sudo systemctl status httpcat
```

### Post-install Directory Layout

After `install.sh`, files follow Linux FHS conventions:

```
/usr/local/bin/
└── httpcat                         # Executable

/etc/httpcat/
└── svr.yml                         # Config file

/var/log/httpcat/
└── httpcat.log                     # Log file

/var/lib/httpcat/
├── static/                         # Web UI static assets
├── upload/                         # Uploaded files
├── download/                       # Download cache
└── data/
    ├── httpcat_sqlite.db           # SQLite database
    └── chunks/                     # Chunked upload temp dir (v0.7.0+)
```

### Service Management

```bash
# Start/stop/restart
sudo systemctl start httpcat
sudo systemctl stop httpcat
sudo systemctl restart httpcat

# Status and logs
sudo systemctl status httpcat
sudo journalctl -u httpcat -f
```

### Uninstall

```bash
# Standard uninstall (keeps config and data)
sudo ./uninstall.sh

# Full uninstall (removes all config and data)
sudo ./uninstall.sh --purge

# Full uninstall but keep user-uploaded files
sudo ./uninstall.sh --purge --keep-data
```

## 🤖 MCP (Model Context Protocol) Support

HttpCat supports the MCP protocol, allowing AI assistants to manage your file server directly.

### Quick Configuration

> ⚠️ **Since v0.4.0, MCP requires `auth_token` when enabled** — set it in `svr.yml`:
>
> ```yaml
> server:
>   mcp:
>     enable: true
>     auth_token: "your_secure_token"
> ```

In your MCP client config (Claude Desktop, Cursor, CodeBuddy, etc.):

```json
{
  "mcpServers": {
    "httpcat": {
      "type": "sse",
      "url": "http://your-server:8888/mcp/sse",
      "headers": {
        "Authorization": "Bearer your_secure_token"
      }
    }
  }
}
```

### Available MCP Tools

| Tool | Description |
|------|-------------|
| `list_files` | List files in upload directory |
| `get_file_info` | Get file details (size, MD5, etc.) |
| `upload_file` | Upload file via MCP (requires Token) |
| `upload_image` | Upload image via MCP |
| `create_folder` | Create folder (v0.5.0+) |
| `rename_file` | Rename file/folder (v0.5.0+) |
| `batch_delete_files` | Batch delete files (v0.5.0+) |
| `request_delete_file` | Request file deletion (step 1) |
| `confirm_delete_file` | Confirm file deletion (step 2) |
| `get_disk_usage` | Get disk usage |
| `get_statistics` | Get upload/download statistics |
| `get_upload_history` | Query upload history |
| `get_download_history` | Query download history (v0.5.0+) |
| `get_file_overview` | File overview statistics (v0.5.0+) |
| `verify_file_md5` | Verify file MD5 integrity |

📖 For detailed MCP guide, see [docs/MCP_USAGE.md](MCP_USAGE.md)

## 🧠 AI Skill (Agent Skills Spec)

HttpCat provides a Skill package following the [Agent Skills specification](https://agentskills.io/), installable in Claude Code / CodeBuddy / Cursor and other AI IDEs for natural language file management.

```bash
# Install in Claude Code
ln -s /path/to/httpcat/httpcat-skill ~/.claude/skills/httpcat

# Install in CodeBuddy
ln -s /path/to/httpcat/httpcat-skill .codebuddy/skills/httpcat

# Install in Cursor
ln -s /path/to/httpcat/httpcat-skill .cursor/skills/httpcat
```

After installation, you can say things like "list files on httpcat", "upload file to server", "check disk usage" in AI conversations.

📖 For details, see [httpcat-skill/README.md](../httpcat-skill/README.md)

## 🤝 OpenClaw + httpcat Integration

HttpCat can be combined with [OpenClaw](https://clawd.org.cn/) (system-level AI Agent) to let AI assistants manage files directly in WeCom / QQ / DingTalk / Lark via the MCP protocol.

📖 For the complete OpenClaw + httpcat deployment guide, see [docs/OPENCLAW_HTTPCAT_GUIDE.md](OPENCLAW_HTTPCAT_GUIDE.md)

## 📦 Large File Upload (v0.7.0+)

HttpCat v0.7.0 adds **chunked upload + resumable upload**, solving these pain points:

- ❌ 1GB+ files often fail on unstable networks during single-request upload
- ❌ Upload to 99% and disconnect — start over from scratch
- ❌ Closing the browser loses all progress
- ❌ Re-uploading the same file wastes bandwidth

### 🎯 Frontend Usage: Fully Automatic

**No configuration needed in browser**. The frontend auto-selects based on file size:

| File size | Frontend strategy | Actual endpoint |
|-----------|------------------|-----------------|
| **< 10 MB** | Single upload | `POST /api/v1/file/upload` |
| **≥ 10 MB** | Chunked upload (5MB/chunk, 3 concurrent) | `POST /upload/init` → `chunk` → `complete` |

Open File Management page → drag files → the right strategy is used automatically, with a live progress bar.

### 🔧 Script/CI Usage: Choose the Right Mode

**Small files / simple scenarios** (backward compatible with v0.6.0):

```bash
# Direct single upload (best for < 100MB on stable networks)
curl -X POST http://localhost:8888/api/v1/file/upload \
  -H "UploadToken: httpcat:xxx:xxx" \
  -F "f1=@file.zip" \
  -F "dir=backup/2026"
```

**Large files / weak network / resumable upload needed**:

```bash
#!/bin/bash
# Complete chunked upload example
HOST="http://localhost:8888"
TOKEN="httpcat:vO9Mt5UtCXWVEaYumi4LxXFImh4=:e30="
FILE="my_big_file.zip"
SIZE=$(stat -c%s "$FILE" 2>/dev/null || stat -f%z "$FILE")
MD5=$(md5sum "$FILE" 2>/dev/null | cut -d' ' -f1 || md5 -q "$FILE")
CHUNK_SIZE=$((5 * 1024 * 1024))  # 5MB
TOTAL=$(( (SIZE + CHUNK_SIZE - 1) / CHUNK_SIZE ))

# Step 1: Initialize session
RESP=$(curl -s -X POST "$HOST/api/v1/file/upload/init" \
  -H "Content-Type: application/json" -H "UploadToken: $TOKEN" \
  -d "{\"fileName\":\"$FILE\",\"fileSize\":$SIZE,\"chunkSize\":$CHUNK_SIZE,\"fileMD5\":\"$MD5\",\"dir\":\"uploads\",\"overwrite\":true}")
UPLOAD_ID=$(echo "$RESP" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['uploadId'])")
INSTANT=$(echo "$RESP" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['instant'])")

# If instant-upload hit, we're done
if [ "$INSTANT" = "True" ]; then
  echo "✅ Instant upload hit (file already exists)"
  exit 0
fi

# Step 2: Upload each chunk
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

# Step 3: Merge chunks
curl -s -X POST "$HOST/api/v1/file/upload/complete" \
  -H "Content-Type: application/json" -H "UploadToken: $TOKEN" \
  -d "{\"uploadId\":\"$UPLOAD_ID\"}"
```

**Python version (with resume support)**:

```python
import hashlib, requests, os, json

HOST = "http://localhost:8888"
TOKEN = "httpcat:xxx:xxx"
FILE = "my_big_file.zip"
CHUNK = 5 * 1024 * 1024

size = os.path.getsize(FILE)
total = (size + CHUNK - 1) // CHUNK
md5 = hashlib.md5(open(FILE, "rb").read()).hexdigest()

# Initialize session
r = requests.post(f"{HOST}/api/v1/file/upload/init",
    headers={"UploadToken": TOKEN, "Content-Type": "application/json"},
    data=json.dumps({
        "fileName": os.path.basename(FILE),
        "fileSize": size, "chunkSize": CHUNK,
        "fileMD5": md5, "dir": "uploads", "overwrite": True
    })).json()

if r["data"]["instant"]:
    print("✅ Instant upload hit"); exit()

uid = r["data"]["uploadId"]

# [Resume] Query uploaded chunks, skip them
status = requests.get(f"{HOST}/api/v1/file/upload/status",
    headers={"UploadToken": TOKEN}, params={"uploadId": uid}).json()
uploaded = set(status["data"]["uploadedIdx"])

# Upload missing chunks
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

# Merge
r = requests.post(f"{HOST}/api/v1/file/upload/complete",
    headers={"UploadToken": TOKEN, "Content-Type": "application/json"},
    data=json.dumps({"uploadId": uid})).json()
print(f"✅ Complete: {r['data']}")
```

### 🎁 Additional Capabilities

#### Instant Upload

If a file with the same MD5 already exists on the server, the `init` endpoint returns `instant: true`. The server uses a **hard link** to create the new file instantly — no actual upload needed:

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

#### Resume Upload

Anytime, calling `GET /upload/status?uploadId=xxx` returns:
- `uploadedIdx`: uploaded chunk indices
- `missingIdx`: missing chunk indices

The client only needs to re-upload chunks in `missingIdx`. **Even if the server restarts**, session state is fully restored from SQLite.

#### Resumable Download

`GET /api/v1/file/download` supports HTTP Range automatically:

```bash
# wget resume
wget -c "http://localhost:8888/api/v1/file/download?filename=big.zip"

# curl resume
curl -C - -o big.zip "http://localhost:8888/api/v1/file/download?filename=big.zip"

# 4-segment parallel download (merge yourself)
for i in 0 1 2 3; do
  START=$((i * 25000000)); END=$((START + 24999999))
  curl -H "Range: bytes=$START-$END" -o part$i "..." &
done; wait
cat part{0,1,2,3} > big.zip
```

### 📊 Parameter Limits

| Parameter | Default | Min | Max | Notes |
|-----------|---------|-----|-----|-------|
| `chunkSize` | 5 MB | 64 KB | 100 MB | Chunk size |
| `fileSize` | - | 1 byte | **100 GB** | Total file size |
| Session TTL | 24 hours | - | - | Chunks auto-cleaned after expiry |
| Frontend threshold | 10 MB | - | - | Adjustable in `FileList/index.tsx` via `CHUNK_THRESHOLD` |
| Frontend concurrency | 3 | 1 | - | Adjustable via `concurrent` param in `chunkedUpload()` |

### 🛡️ Security

- ✅ Path traversal protection: all paths validated by `ResolvePathWithinBase`
- ✅ Chunk size validation: all chunks except the last must be exactly `chunkSize`
- ✅ Optional `chunkMD5`: server verifies individual chunk integrity
- ✅ Final MD5 check: if client provided `fileMD5`, server must match after merge
- ✅ No orphan files: any error during merge cleans up temp files

> 📖 **Want to dive deeper?** See [**Chunked Upload, Resume & Instant Upload Principles**](CHUNK_UPLOAD_PRINCIPLE.md) for complete engineering design, security mechanisms, performance comparisons, and industry references.

## 📡 API Reference

### Authentication

HttpCat supports two authentication methods:

| Method | Use Case | Header |
|--------|----------|--------|
| JWT Token | Web frontend login | `Authorization: Bearer <token>` |
| AK/SK Signature | Scripts/CI/AI/Open API | `AccessKey` + `Signature` + `TimeStamp` |

### Open API (AK/SK Signature Auth)

When enabled, all `/api/v1/*` endpoints can be called via AK/SK signature without JWT login.

Enable in `svr.yml`:

```yaml
server:
  http:
    auth:
      open_api_enable: true
      aksk:
        your_access_key: your_secret_key
```

Signature algorithm:

```
Signature = HMAC-SHA256(
  "{Method}\n{Path}\n{Query}\n{AccessKey}\n{TimeStamp}\n{BodySHA256}",
  SecretKey
)
```

### Main Endpoints

**File Management**

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/v1/file/upload` | POST | Single-request upload |
| `/api/v1/file/download` | GET | Download file (supports HTTP Range, v0.7.0+) |
| `/api/v1/file/listFiles` | GET | List files (supports subdirectories) |
| `/api/v1/file/getFileInfo` | GET | Get file info |
| `/api/v1/file/preview` | GET | Online file preview (v0.6.0+) |
| `/api/v1/file/previewInfo` | GET | Get preview metadata (v0.6.0+) |
| `/api/v1/file/downloadZip` | POST | Batch download as ZIP (v0.6.0+) |
| `/api/v1/file/delete` | POST | Batch delete files/folders (v0.5.0+) |
| `/api/v1/file/mkdir` | POST | Create folder (v0.5.0+) |
| `/api/v1/file/rename` | POST | Rename file/folder (v0.5.0+) |
| `/api/v1/file/upload/init` | POST | Initialize chunked upload session (v0.7.0+) |
| `/api/v1/file/upload/status` | GET | Query chunked upload status (for resume, v0.7.0+) |
| `/api/v1/file/upload/chunk` | POST | Upload a single chunk (v0.7.0+) |
| `/api/v1/file/upload/complete` | POST | Merge chunks into final file (v0.7.0+) |
| `/api/v1/file/upload/abort` | POST | Abort chunked upload session (v0.7.0+) |

**Sharing**

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/v1/share` | POST | Create share link (v0.4.0+) |
| `/api/v1/share/list` | GET | List all shares |
| `/api/v1/share/:code` | DELETE | Delete a share |
| `/api/v1/share/stats` | GET | Share statistics |
| `/s/:code` | GET | Share landing page (anonymous) |
| `/s/:code/download` | GET | Download via share code |

**Operation Log (v0.6.0+)**

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/v1/oplog/list` | GET | List operation logs (with filters) |
| `/api/v1/oplog/stats` | GET | Operation statistics |

**Statistics & Overview**

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/v1/user/dataOverview` | GET | Data overview (file count, disk usage) |
| `/api/v1/statistics/getUploadStatistics` | GET | Upload statistics |
| `/api/v1/statistics/getDownloadStatistics` | GET | Download statistics |
| `/api/v1/statistics/getFileOverview` | GET | File overview (v0.5.0+) |

**User & Token**

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/v1/user/login/account` | POST | Login (rate-limited, v0.7.0+) |
| `/api/v1/user/changePasswd` | POST | Change password |
| `/api/v1/user/createUploadToken` | POST | Create upload token |
| `/api/v1/user/uploadTokenLists` | GET | List upload tokens |

📖 For detailed API examples (Shell & Python), see the [Chinese README](../README.md#-api-接口).

### Download File

```bash
# Download (whitelist endpoint, no auth required when upload_token_check is off)
wget -O filename.jpg "http://localhost:8888/api/v1/file/download?filename=filename.jpg"

# Resume download (v0.7.0+)
wget -c "http://localhost:8888/api/v1/file/download?filename=big.zip"
curl -C - -o big.zip "http://localhost:8888/api/v1/file/download?filename=big.zip"
```

### List Files

```bash
# Requires JWT or AK/SK auth
curl http://localhost:8888/api/v1/file/listFiles?dir=/
```

## ⚙️ Configuration

Config file: `conf/svr.yml` (inside Docker) or `/etc/httpcat/svr.yml` (via install.sh).

```yaml
server:
  http:
    port: 8888
    # Large file upload needs longer read/write timeout (seconds)
    read_timeout: 1800
    write_timeout: 1800
    auth:
      open_api_enable: false       # Enable Open API (AK/SK signature auth)
      aksk:                        # AK/SK credential pairs
        your_access_key: your_secret_key
    file:
      upload_enable: true
      enable_upload_token: true    # Enable UploadToken verification
      app_key: "httpcat"           # app_key for generating UploadToken
      app_secret: "httpcat_app_secret"
      upload_policy:
        deadline: 7200             # UploadToken TTL (seconds)
        fsizeLimit: 0              # Max file size (bytes, 0 = unlimited)
      download_dir: "website/upload/"
      enable_sqlite: true
      sqlite_db_path: "./data/httpcat_sqlite.db"

  mcp:
    enable: true                   # Enable MCP Server (AI agent access)
    auth_token: "replace-with-secure-token"

  share:
    enable: true                   # Enable file sharing
    anonymous_access: true         # Allow anonymous access to share links
```

> 💡 **v0.7.0 chunked upload parameters**: 5MB/chunk, 24-hour session TTL, 100GB max file size — all code defaults, no config needed. Frontend threshold (10MB) is adjustable in `web/src/pages/FileManage/FileList/index.tsx` via `CHUNK_THRESHOLD`.

## 🍀 FAQ

### Forgot Password?

**Option A: Reset only admin password (keep all data, recommended)**

Requires Python with bcrypt:

```bash
# Install dependency (one-time)
pip3 install bcrypt

# Reset admin password to admin123 (customize as needed)
sudo systemctl stop httpcat
sudo python3 -c "
import bcrypt, sqlite3, time
pwd = b'admin123'  # your new password
h = bcrypt.hashpw(pwd, bcrypt.gensalt(10)).decode()
c = sqlite3.connect('/var/lib/httpcat/data/httpcat_sqlite.db')
c.execute(\"UPDATE users SET password=?, salt='', password_update_time=? WHERE username='admin'\", (h, int(time.time())))
c.commit(); c.close()
print('✅ Password reset')
"
sudo systemctl start httpcat
```

**Option B: Reset database (wipes all data)**

```bash
sudo systemctl stop httpcat
sudo rm /var/lib/httpcat/data/httpcat_sqlite.db
sudo systemctl start httpcat
# After restart, a default admin / admin account is recreated
```

### Login returns "too many failed attempts, please try again later"?

This is the v0.7.0 **login rate-limit** feature: 5 failures in 5 min from the same IP → 15 min lockout. Wait for the `lockedRemainingSeconds` in the response, or **restart the service** to clear the in-memory lockout state:

```bash
sudo systemctl restart httpcat
```

### Will failed chunked uploads leave temp files?

No. Three safeguards:

1. Incomplete chunks use `.part` suffix during write — failures don't pollute the bitmap
2. Sessions expire after **24 hours**; a background job scans every 30 min to clean up `data/chunks/{uploadId}/`
3. You can also call `POST /api/v1/file/upload/abort` to cancel and clean up proactively

Check current residuals: `ls /var/lib/httpcat/data/chunks/`

### Why is there no progress bar after uploading a large file?

Make sure the file is **≥ 10 MB**: the frontend has `CHUNK_THRESHOLD = 10MB`. Smaller files use the legacy single-request upload (no fine-grained progress). You can adjust the threshold in `web/src/pages/FileManage/FileList/index.tsx`.

### Node.js Version Issues?

This project uses UmiJS 3.x and requires the legacy OpenSSL provider on Node.js 17+:

```bash
NODE_OPTIONS=--openssl-legacy-provider npm run start:dev
```

> Recommended: **Node.js v20.x LTS** (build scripts handle compatibility). `nvm` users can run `nvm use` in the `web/` directory to auto-switch.

## 🛠️ Development

### Prerequisites

- **Go 1.23+** - Backend compilation
- **Node.js 20+** - Frontend build (recommended: v20.x LTS)
- **npm** - Package manager (comes with Node.js)

> 💡 **Tip**: nvm users can run `nvm use` in the `web/` directory to auto-switch to the project-specified version.

### Build Commands

```bash
# Interactive build
./scripts/build.sh

# Build all platforms with frontend
./scripts/build.sh -a -f

# Build specific platform
./scripts/build.sh -p linux_amd64 -f

# Build with Docker (full CGO support for Linux)
./scripts/build.sh -d -f

# Show help
./scripts/build.sh -h
```

## 📅 Version Evolution

| Version | Theme | Highlights |
|---------|-------|-----------|
| **v0.7.0** | 📦 Large files & security | Chunked upload + resume + instant-upload + Range download + login rate-limit |
| v0.6.0 | 🔍 Audit & UX | Operation log, online file preview (text/image/video/PDF), multi-file ZIP download |
| v0.5.0 | 📂 Deep file management | Subdirectories, file overview, batch ops, 15 MCP Tools |
| v0.4.0 | 🔒 Security & sharing | bcrypt password hashing, file sharing (TTL/count/access-code) |
| v0.3.0 | 🌐 Web autonomy | In-browser system config management |
| v0.2.x | 🤖 AI integration | MCP protocol, Docker image, AK/SK signature auth |
| v0.1.x | 🎯 Foundation | Upload/download, SQLite, web UI |

> 📖 For detailed changes, see [CHANGELOG.md](../CHANGELOG.md) and [docs/ReleaseNote.md](ReleaseNote.md)

## 📝 License

This software is for personal use only and is strictly prohibited for commercial purposes.

- Prohibited for commercial purposes
- Copyright and license statements must be preserved
- This software is provided "as is" without any warranties

## 🌟 Contributing

Welcome to follow our GitHub project! ⭐ Star it to stay updated with our real-time developments.

Welcome to submit issues or pull requests. Good luck! 🍀
