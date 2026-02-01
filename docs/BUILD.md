# HttpCat 编译构建指南

## 快速构建

```bash
# 给脚本添加执行权限
chmod +x scripts/*.sh

# 交互式构建（推荐）
./scripts/build.sh

# 一键构建所有平台 + 前端
./scripts/build.sh -a -f
```

## 构建选项

| 参数 | 说明 |
|------|------|
| `-h, --help` | 显示帮助信息 |
| `-v, --version VER` | 指定版本号（默认 v0.2.0） |
| `-p, --platform PLAT` | 只构建指定平台 |
| `-a, --all` | 构建所有平台 |
| `-f, --frontend` | 构建前端 |
| `-s, --skip-frontend` | 跳过前端构建 |
| `-d, --docker` | 使用 Docker 构建（完整 CGO 支持） |
| `--clean` | 只清理构建目录 |

## 支持的平台

| 平台 | 标识 | 说明 |
|------|------|------|
| Linux x86_64 | `linux_amd64` | 主流服务器 |
| Linux ARM64 | `linux_arm64` | ARM 服务器、树莓派等 |
| macOS Intel | `darwin_amd64` | Intel Mac |
| macOS Apple Silicon | `darwin_arm64` | M1/M2/M3 Mac |
| Windows x64 | `windows_amd64` | Windows 系统 |

## 构建示例

### 1. 交互式构建（推荐新手）

```bash
./scripts/build.sh
```

脚本会引导你选择：
- 是否构建前端
- 选择目标平台

### 2. 构建三个主要平台

```bash
# Linux x86_64 + Linux ARM64 + macOS ARM64
./scripts/build.sh -f
```

### 3. 构建所有平台

```bash
./scripts/build.sh -a -f
```

### 4. 只构建特定平台

```bash
# 只构建 Linux x86_64
./scripts/build.sh -p linux_amd64 -f

# 只构建 Linux 两个架构
./scripts/build.sh -p linux_amd64 -p linux_arm64 -f
```

### 5. 使用 Docker 构建（推荐 macOS 用户）

在 macOS 上交叉编译 Linux 版本时，由于 CGO 限制无法启用 SQLite 支持。
使用 Docker 构建可以获得完整的 CGO 支持：

```bash
# 使用 Docker 构建 Linux 版本（完整 SQLite 支持）
./scripts/build.sh -d -f
```

**Docker 镜像加速**（默认已启用）：

构建脚本默认使用国内镜像源加速：
- Go 基础镜像：`m.daocloud.io/docker.io/golang:1.23-alpine`
- Alpine 基础镜像：`m.daocloud.io/docker.io/alpine:3.19`
- Go 模块代理：`https://goproxy.cn`
- Alpine APK 源：`mirrors.aliyun.com`

**静态链接**：Docker 构建使用静态链接，不依赖系统库，兼容所有 Linux 发行版（Ubuntu、CentOS、Debian 等）

**直接使用**（无需任何配置）：

```bash
./scripts/build.sh -d -f
```

**使用其他镜像源**（可选）：

```bash
# 使用官方镜像（不使用国内加速）
GO_BASE_IMAGE=golang:1.23-alpine \
ALPINE_BASE_IMAGE=alpine:3.19 \
./scripts/build.sh -d -f

# 使用阿里云镜像
GO_BASE_IMAGE=registry.cn-hangzhou.aliyuncs.com/acs-sample/golang:1.23-alpine \
ALPINE_BASE_IMAGE=registry.cn-hangzhou.aliyuncs.com/acs-sample/alpine:3.19 \
./scripts/build.sh -d -f
```

**配置 Docker 镜像加速器**（全局配置，可选）：

```bash
# 编辑 Docker 配置文件
# Linux: /etc/docker/daemon.json
# macOS: ~/.docker/daemon.json

{
  "registry-mirrors": [
    "https://docker.mirrors.ustc.edu.cn",
    "https://hub-mirror.c.163.com",
    "https://mirror.baidubce.com"
  ]
}

# 重启 Docker 服务
# Linux
sudo systemctl restart docker

# macOS: 通过 Docker Desktop 界面重启
```

### 6. 指定版本号

```bash
# 指定版本号为 v1.0.0
./scripts/build.sh -v v1.0.0 -a -f

# 或通过环境变量
HTTPCAT_VERSION=v1.0.0 ./scripts/build.sh -a -f
```

## 构建产物

构建完成后，`release/` 目录结构：

```
release/
├── httpcat_v0.2.0_linux-amd64.tar.gz    # Linux x86_64 安装包
├── httpcat_v0.2.0_linux-arm64.tar.gz    # Linux ARM64 安装包
├── httpcat_v0.2.0_darwin-arm64.tar.gz   # macOS ARM64 安装包
├── httpcat_v0.2.0_darwin-amd64.tar.gz   # macOS Intel 安装包
└── httpcat_v0.2.0_windows-amd64.zip     # Windows 安装包
```

每个安装包包含：
- `httpcat` - 可执行文件
- `conf/` - 配置文件目录
- `static/` - 前端静态资源
- `README.md` - 说明文档
- `install.sh` - 安装脚本（Linux/macOS）
- `uninstall.sh` - 卸载脚本（Linux/macOS）
- `httpcat.service` - systemd 服务文件（Linux）

## CGO 与 SQLite 说明

HttpCat 使用 SQLite 作为数据库，需要 CGO 支持。**SQLite 是核心功能依赖**，无 SQLite 将导致用户登录等功能无法使用。

### 不同系统的构建策略

| 构建环境 | 目标平台 | CGO | SQLite | 推荐方案 |
|---------|---------|-----|--------|---------|
| **Linux** | Linux (本机) | ✅ | ✅ | 直接编译 |
| **Linux** | Linux (其他架构) | ✅ | ✅ | 使用交叉编译器 |
| **macOS** | macOS (本机) | ✅ | ✅ | 直接编译 |
| **macOS** | Linux | ❌ | ❌ | **使用 Docker** |

### Linux 环境构建（推荐）

在 Linux 服务器上直接编译，天然支持 CGO：

```bash
# Linux 上直接编译，自动启用 CGO
./scripts/build.sh -p linux_amd64 -f

# 编译多个 Linux 架构（需要交叉编译器）
./scripts/build.sh -p linux_amd64 -p linux_arm64 -f
```

### macOS 环境构建

在 macOS 上交叉编译 Linux 版本时，**脚本会自动检测并建议使用 Docker**：

```bash
# 方式 1：交互式构建（推荐）
# 脚本会自动询问是否使用 Docker
./scripts/build.sh -f

# 方式 2：明确使用 Docker
./scripts/build.sh -d -f

# 方式 3：只构建 macOS 版本（不需要 Docker）
./scripts/build.sh -p darwin_arm64 -f
```

**注意**：如果在 macOS 上不使用 Docker 构建 Linux 版本，编译产物将**无法登录**（无 SQLite 支持）。

## 环境要求

### 必需

- **Go 1.23+** - 后端编译
- Git

### 可选

- **Node.js 20+** - 构建前端（推荐 v20.x LTS）
- Docker - 跨平台构建（macOS 上编译 Linux 版本时推荐）
- `aarch64-linux-gnu-gcc` - Linux 上交叉编译 ARM64

> 💡 **提示**: 使用 nvm 的用户可在 `web/` 目录运行 `nvm use` 自动切换版本

### 安装交叉编译工具链（Linux）

```bash
# Ubuntu/Debian
sudo apt-get install gcc-aarch64-linux-gnu

# CentOS/RHEL
sudo yum install gcc-aarch64-linux-gnu
```

## 手动编译

如果不想使用构建脚本，可以手动编译：

```bash
# 进入后端目录
cd server-go

# 下载依赖
go mod tidy
go mod download

# 编译当前平台
CGO_ENABLED=1 go build -o httpcat ./cmd/httpcat.go

# 交叉编译（禁用 CGO）
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o httpcat-linux-amd64 ./cmd/httpcat.go
```

## 使用编译产物

```bash
# 解压安装包
tar -xzf httpcat_v0.2.0_linux-amd64.tar.gz
cd httpcat_v0.2.0_linux-amd64

# 方式 1: 使用安装脚本（推荐）
sudo ./install.sh

# 方式 2: 直接运行
./httpcat --port=8888 -C conf/svr.yml

# 方式 3: 后台运行
nohup ./httpcat --port=8888 -C conf/svr.yml > httpcat.log 2>&1 &
```

## 常见问题

### Q: macOS 上编译报错 `cgo: C compiler "gcc" not found`

安装 Xcode 命令行工具：

```bash
xcode-select --install
```

### Q: Node.js 17+ 构建前端报错

设置环境变量：

```bash
export NODE_OPTIONS=--openssl-legacy-provider
npm run build
```

### Q: 如何验证版本信息？

```bash
./httpcat -v
# 输出: httpcat version v0.2.0 (build: 202401311200, commit: abc1234)
```

### Q: 如何只清理构建目录？

```bash
./scripts/build.sh --clean
```

## 提交代码前

检查 Git 用户配置：

```bash
# 查看全局配置
git config --global user.name
git config --global user.email

# 查看当前仓库配置
git config user.name
git config user.email
```
