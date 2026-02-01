English | [简体中文](docs/README-cn.md)

# 🐱 HttpCat

> A lightweight, efficient HTTP file transfer service with modern web interface and AI integration.

HttpCat is designed to provide a simple, efficient, and stable solution for file uploading and downloading. Whether it's for temporary sharing or bulk file transfers, HttpCat will be your excellent assistant.

## ✨ Key Features

- 🚀 **Simple & Efficient** - Easy to deploy, no external dependencies
- 🎨 **Modern Web UI** - Beautiful React-based management interface
- 🤖 **MCP Support** - AI assistants (Claude, Cursor, CodeBuddy) can directly manage your files
- 🐳 **Docker Ready** - One-command deployment with Docker
- 🔐 **Secure** - Token-based authentication for uploads
- 📊 **Statistics** - Track uploads/downloads with detailed history

## 📁 Project Structure

```
httpcat/
├── server-go/              # 🔧 Go Backend
│   ├── cmd/                # Application entry point
│   │   └── httpcat.go
│   ├── internal/           # Internal packages
│   │   ├── common/         # Shared utilities
│   │   ├── handler/        # HTTP handlers
│   │   ├── mcp/            # MCP server implementation
│   │   ├── midware/        # Middleware (auth, metrics)
│   │   ├── models/         # Data models
│   │   ├── p2p/            # P2P functionality
│   │   ├── server/         # Server core
│   │   ├── storage/        # Storage layer
│   │   └── conf/           # Configuration files
│   ├── go.mod
│   └── go.sum
│
├── web/                    # 🎨 React Frontend
│   ├── src/                # Source code
│   ├── config/             # UmiJS configuration
│   ├── mock/               # Mock data (dev only)
│   └── package.json
│
├── scripts/                # 🛠️ Scripts
│   ├── build.sh            # Multi-platform build script
│   ├── install.sh          # Linux installation script
│   ├── uninstall.sh        # Uninstallation script
│   └── translations.sh     # i18n translation script
│
├── docs/                   # 📚 Documentation
│   ├── README-cn.md        # Chinese README
│   ├── BUILD.md            # Build guide
│   ├── ReleaseNote.md      # Release history
│   ├── MCP_USAGE.md        # MCP integration guide
│   └── ...                 # Other design docs
│
├── static/                 # 📦 Frontend build output
├── release/                # 📤 Build artifacts (gitignored)
│
├── Dockerfile              # Docker configuration
├── docker-compose.yml      # Docker Compose setup
└── httpcat.service         # Systemd service file
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

Access the application:
- **Frontend**: http://localhost:8000 (dev) or http://localhost:8888 (prod)
- **Backend API**: http://localhost:8888/api/v1/

### Default Credentials

| Field | Value |
|-------|-------|
| Username | `admin` |
| Password | `admin` |

> ⚠️ **Security**: Change the default password after first login!

## 🎉 Installation (Production)

### Quick Install

```bash
# Download and extract
httpcat_version="v0.2.0"
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

### Installation Directory Structure

After running `install.sh`, files are organized following Linux FHS standard:

```
/usr/local/bin/
└── httpcat                         # Executable

/etc/httpcat/
└── svr.yml                         # Configuration

/var/log/httpcat/
└── httpcat.log                     # Log files

/var/lib/httpcat/
├── static/                         # Web UI assets
├── upload/                         # Uploaded files
├── download/                       # Download cache
└── data/
    └── httpcat_sqlite.db           # SQLite database
```

### Service Management

```bash
# Start/Stop/Restart
sudo systemctl start httpcat
sudo systemctl stop httpcat
sudo systemctl restart httpcat

# View status and logs
sudo systemctl status httpcat
sudo journalctl -u httpcat -f
```

### Uninstall

```bash
# Standard uninstall (keeps config and data)
sudo ./uninstall.sh

# Complete removal (deletes everything)
sudo ./uninstall.sh --purge

# Keep uploaded files only
sudo ./uninstall.sh --purge --keep-data
```

## 🤖 MCP (Model Context Protocol) Support

HttpCat supports MCP, allowing AI assistants to directly manage your file server.

### Quick Setup

Add to your MCP client configuration (Claude Desktop, Cursor, CodeBuddy, etc.):

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

### Available MCP Tools

| Tool | Description |
|------|-------------|
| `list_files` | List files in upload directory |
| `get_file_info` | Get file details (size, MD5, etc.) |
| `upload_file` | Upload file via MCP (requires Token) |
| `get_disk_usage` | Get disk usage information |
| `get_upload_history` | Query upload history |
| `request_delete_file` | Request file deletion (step 1) |
| `confirm_delete_file` | Confirm file deletion (step 2) |
| `get_statistics` | Get upload/download statistics |
| `verify_file_md5` | Verify file MD5 checksum |

📖 For detailed MCP usage guide, see [docs/MCP_USAGE.md](docs/MCP_USAGE.md)

## 📡 API Reference

### Upload File

```bash
curl -v -F "f1=@/path/to/file" \
  -H "UploadToken: your-token" \
  http://localhost:8888/api/v1/file/upload
```

### Download File

```bash
wget -O filename.jpg http://localhost:8888/api/v1/file/download?filename=filename.jpg
```

### List Files

```bash
curl http://localhost:8888/api/v1/file/listFiles?dir=/
```

## ⚙️ Configuration

Configuration file: `svr.yml`

```yaml
# Server settings
port: 8888
upload_dir: "./upload"
download_dir: "./upload"
static_dir: "./static"

# Authentication
app_key: "httpcat"
app_secret: "httpcat_app_secret"
enable_upload_token: true

# Database
enable_sqlite: true
sqlite_db_path: "./data/sqlite.db"

# Notifications
persistent_notify_url: ""  # WeChat webhook URL
```

## 🍀 FAQ

### Forgot Password?

Delete the SQLite database and restart:

```bash
find / -name "*.db" | grep httpcat
rm /path/to/httpcat_sqlite.db
systemctl restart httpcat
```

A new admin user will be created with default credentials.

### Node.js Version Issues?

For Node.js 17+, use the legacy OpenSSL provider:

```bash
NODE_OPTIONS=--openssl-legacy-provider npm run start:dev
```

Recommended: Use Node.js v16.x for best compatibility.

## 🛠️ Development

### Prerequisites

- Go 1.19+
- Node.js 16+ (recommended: v16.18.0)
- npm or yarn

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

## 📝 License

This software is for personal use only and is strictly prohibited for commercial purposes.

- Prohibited for commercial purposes
- Copyright and license statements must be preserved
- This software is provided "as is" without any warranties

## 🌟 Contributing

Welcome to follow our GitHub project! ⭐ Star it to stay updated with our real-time developments.

Feel free to raise issues or submit pull requests. Good luck! 🍀
