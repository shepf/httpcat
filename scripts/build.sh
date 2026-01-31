#!/bin/bash

set -e  # 遇到错误立即退出

# macOS 上禁用 tar 的扩展属性（避免 Linux 解压警告）
export COPYFILE_DISABLE=1

# ============ 颜色定义 ============
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# ============ 项目目录 ============
# 脚本位于 scripts/ 目录，需要回到上级目录
SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd)
PROJECT_ROOT=$(cd "$SCRIPT_DIR/.." && pwd)
SERVER_DIR="$PROJECT_ROOT/server-go"
WEB_DIR="$PROJECT_ROOT/web"
STATIC_DIR="$PROJECT_ROOT/static"
RELEASE_DIR="$PROJECT_ROOT/release"
SCRIPTS_DIR="$PROJECT_ROOT/scripts"

# ============ 版本信息 ============
HTTPCAT_VERSION="${HTTPCAT_VERSION:-v0.2.0}"
HTTPCAT_BUILD=$(date "+%Y%m%d%H%M")
COMMIT_ID=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# ============ 目标平台 ============
# 格式: OS_ARCH
PLATFORMS=(
    "linux_amd64"
    "linux_arm64"
    "darwin_amd64"
    "darwin_arm64"
    "windows_amd64"
)

# ============ 帮助信息 ============
show_help() {
    echo -e "${GREEN}HttpCat 构建脚本${NC}"
    echo ""
    echo "用法: $0 [选项]"
    echo ""
    echo "选项:"
    echo "  -h, --help          显示帮助信息"
    echo "  -v, --version VER   指定版本号 (默认: $HTTPCAT_VERSION)"
    echo "  -p, --platform PLAT 只构建指定平台 (linux_amd64, linux_arm64, darwin_amd64, darwin_arm64, windows_amd64)"
    echo "  -a, --all           构建所有平台 (默认)"
    echo "  -f, --frontend      构建前端 (默认询问)"
    echo "  -s, --skip-frontend 跳过前端构建"
    echo "  -d, --docker        使用 Docker 构建 (支持完整 CGO)"
    echo "  --clean             只清理构建目录"
    echo ""
    echo "示例:"
    echo "  $0                      # 交互式构建"
    echo "  $0 -a -f                # 构建所有平台，包含前端"
    echo "  $0 -p linux_amd64       # 只构建 Linux x86_64"
    echo "  $0 -d                   # 使用 Docker 构建所有 Linux 版本"
    echo "  $0 -v v1.0.0 -a -f      # 指定版本号构建所有平台"
}

# ============ 检测操作系统 ============
detect_os() {
    OS_TYPE=$(uname -s)
    ARCH_TYPE=$(uname -m)
    
    echo -e "${GREEN}检测到操作系统: $OS_TYPE ($ARCH_TYPE)${NC}"
    
    case "$OS_TYPE" in
        Darwin) HOST_OS="darwin" ;;
        Linux)  HOST_OS="linux" ;;
        MINGW*|MSYS*|CYGWIN*) HOST_OS="windows" ;;
        *) HOST_OS="unknown" ;;
    esac
    
    case "$ARCH_TYPE" in
        x86_64|amd64) HOST_ARCH="amd64" ;;
        arm64|aarch64) HOST_ARCH="arm64" ;;
        *) HOST_ARCH="unknown" ;;
    esac
}

# ============ 环境检查 ============
check_environment() {
    echo -e "${BLUE}检查构建环境...${NC}"
    
    if ! type go >/dev/null 2>&1; then
        echo -e "${RED}错误: Go 未安装${NC}"
        exit 1
    fi
    echo -e "${GREEN}✓ Go 版本: $(go version)${NC}"
    
    if type node >/dev/null 2>&1; then
        echo -e "${GREEN}✓ Node 版本: $(node -v)${NC}"
    else
        echo -e "${YELLOW}⚠ Node.js 未安装，将跳过前端构建${NC}"
    fi
    
    if type docker >/dev/null 2>&1; then
        echo -e "${GREEN}✓ Docker 可用${NC}"
        DOCKER_AVAILABLE=true
    else
        echo -e "${YELLOW}⚠ Docker 未安装${NC}"
        DOCKER_AVAILABLE=false
    fi
}

# ============ 清理目录 ============
clean_release() {
    echo -e "${BLUE}清理构建目录...${NC}"
    rm -rf "$RELEASE_DIR"
    mkdir -p "$RELEASE_DIR"
    echo -e "${GREEN}✓ 清理完成${NC}"
}

# ============ 构建前端 ============
build_frontend() {
    if [ ! -d "$WEB_DIR" ] || [ ! -f "$WEB_DIR/package.json" ]; then
        echo -e "${YELLOW}警告: 未检测到前端项目${NC}"
        return 1
    fi
    
    echo -e "${BLUE}开始构建前端...${NC}"
    cd "$WEB_DIR"
    
    # 检查 node_modules
    if [ ! -d "node_modules" ]; then
        echo "安装前端依赖..."
        npm install --registry=https://registry.npmmirror.com
    fi
    
    # 设置 OpenSSL 兼容性选项（Node.js 17+ 需要）
    export NODE_OPTIONS=--openssl-legacy-provider
    
    # 构建前端
    echo "构建前端..."
    npm run build
    
    # 复制构建产物到 static 目录
    if [ -d "dist" ]; then
        echo "复制前端构建产物到 static 目录..."
        rm -rf "$STATIC_DIR"/*
        cp -r dist/* "$STATIC_DIR/"
        echo -e "${GREEN}✓ 前端构建完成${NC}"
    else
        echo -e "${RED}✗ 前端构建失败：未生成 dist 目录${NC}"
        return 1
    fi
    
    cd "$PROJECT_ROOT"
}

# ============ Go 依赖处理 ============
prepare_go_deps() {
    echo -e "${BLUE}准备 Go 依赖...${NC}"
    cd "$SERVER_DIR"
    go mod tidy
    go mod download
    cd "$PROJECT_ROOT"
    echo -e "${GREEN}✓ Go 依赖准备完成${NC}"
}

# ============ 构建单个平台 ============
build_platform() {
    local platform=$1
    local os="${platform%_*}"
    local arch="${platform#*_}"
    
    local output_name="httpcat-${os}-${arch}"
    [ "$os" = "windows" ] && output_name="${output_name}.exe"
    
    echo -e "${BLUE}构建 $os/$arch...${NC}"
    
    cd "$SERVER_DIR"
    
    # 判断是否可以启用 CGO
    local cgo_enabled=0
    local cc_compiler=""
    
    # 本机编译可以启用 CGO
    if [ "$os" = "$HOST_OS" ] && [ "$arch" = "$HOST_ARCH" ]; then
        cgo_enabled=1
        echo -e "  ${GREEN}本机编译，启用 CGO${NC}"
    # Linux 环境交叉编译其他 Linux 架构
    elif [ "$HOST_OS" = "linux" ] && [ "$os" = "linux" ]; then
        if [ "$arch" = "arm64" ] && command -v aarch64-linux-gnu-gcc &> /dev/null; then
            cgo_enabled=1
            cc_compiler="aarch64-linux-gnu-gcc"
            echo -e "  ${GREEN}使用交叉编译器: $cc_compiler${NC}"
        elif [ "$arch" = "amd64" ]; then
            cgo_enabled=1
        fi
    fi
    
    # 构建命令
    local build_env="CGO_ENABLED=$cgo_enabled GOOS=$os GOARCH=$arch"
    [ -n "$cc_compiler" ] && build_env="CC=$cc_compiler $build_env"
    
    # 执行构建
    eval $build_env go build \
        -ldflags \"-s -w -X httpcat/internal/common.Version=$HTTPCAT_VERSION -X httpcat/internal/common.Build=$HTTPCAT_BUILD -X httpcat/internal/common.Commit=$COMMIT_ID\" \
        -o \"$RELEASE_DIR/$output_name\" ./cmd/httpcat.go
    
    if [ $? -eq 0 ]; then
        local size=$(ls -lh "$RELEASE_DIR/$output_name" | awk '{print $5}')
        if [ $cgo_enabled -eq 1 ]; then
            echo -e "${GREEN}✓ $output_name ($size) [CGO 启用，支持 SQLite]${NC}"
        else
            echo -e "${YELLOW}✓ $output_name ($size) [CGO 禁用，无 SQLite]${NC}"
        fi
        return 0
    else
        echo -e "${RED}✗ $output_name 构建失败${NC}"
        return 1
    fi
    
    cd "$PROJECT_ROOT"
}

# ============ Docker 构建 (完整 CGO 支持) ============
build_with_docker() {
    if [ "$DOCKER_AVAILABLE" != "true" ]; then
        echo -e "${RED}错误: Docker 不可用${NC}"
        exit 1
    fi
    
    echo -e "${BLUE}使用 Docker 构建 Linux 版本...${NC}"
    
    # 默认使用国内镜像，可通过环境变量覆盖
    local GO_IMAGE="${GO_BASE_IMAGE:-m.daocloud.io/docker.io/golang:1.23-alpine}"
    local ALPINE_IMAGE="${ALPINE_BASE_IMAGE:-m.daocloud.io/docker.io/alpine:3.19}"
    
    echo -e "  使用镜像: ${GO_IMAGE}, ${ALPINE_IMAGE}"
    echo -e "  提示: 可通过 GO_BASE_IMAGE 和 ALPINE_BASE_IMAGE 环境变量自定义镜像"
    
    # 创建临时 Dockerfile 用于多平台构建
    cat > "$PROJECT_ROOT/.build.Dockerfile" << DOCKERFILE
FROM ${GO_IMAGE} AS builder

# 使用国内 Alpine 镜像源
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# 安装编译依赖
RUN apk add --no-cache gcc musl-dev sqlite-dev linux-headers

# Go 模块代理（国内加速）
ENV GOPROXY=https://goproxy.cn,direct
ENV GOSUMDB=sum.golang.google.cn

WORKDIR /app
COPY server-go/go.mod server-go/go.sum ./
RUN go mod download
COPY server-go/ ./

ARG VERSION
ARG BUILD_TIME
ARG COMMIT_ID
ARG TARGETOS
ARG TARGETARCH

RUN CGO_ENABLED=1 GOOS=\${TARGETOS} GOARCH=\${TARGETARCH} \
    CGO_CFLAGS="-D_LARGEFILE64_SOURCE" \
    CGO_LDFLAGS="-static" \
    go build -ldflags "-s -w -extldflags '-static' -X httpcat/internal/common.Version=\${VERSION} -X httpcat/internal/common.Build=\${BUILD_TIME} -X httpcat/internal/common.Commit=\${COMMIT_ID}" \
    -o /httpcat ./cmd/httpcat.go

FROM ${ALPINE_IMAGE}
COPY --from=builder /httpcat /httpcat
CMD ["cat", "/httpcat"]
DOCKERFILE
    
    # 构建 Linux amd64
    echo -e "${BLUE}构建 Linux amd64 (Docker)...${NC}"
    docker build --platform linux/amd64 \
        --build-arg VERSION="$HTTPCAT_VERSION" \
        --build-arg BUILD_TIME="$HTTPCAT_BUILD" \
        --build-arg COMMIT_ID="$COMMIT_ID" \
        --cache-from httpcat-build:amd64-cache \
        -f "$PROJECT_ROOT/.build.Dockerfile" \
        -t httpcat-build:amd64 "$PROJECT_ROOT"
    
    # 提取二进制文件
    docker run --rm httpcat-build:amd64 > "$RELEASE_DIR/httpcat-linux-amd64"
    chmod +x "$RELEASE_DIR/httpcat-linux-amd64"
    echo -e "${GREEN}✓ httpcat-linux-amd64 [CGO 启用，支持 SQLite]${NC}"
    
    # 构建 Linux arm64
    echo -e "${BLUE}构建 Linux arm64 (Docker)...${NC}"
    docker build --platform linux/arm64 \
        --build-arg VERSION="$HTTPCAT_VERSION" \
        --build-arg BUILD_TIME="$HTTPCAT_BUILD" \
        --build-arg COMMIT_ID="$COMMIT_ID" \
        --cache-from httpcat-build:arm64-cache \
        -f "$PROJECT_ROOT/.build.Dockerfile" \
        -t httpcat-build:arm64 "$PROJECT_ROOT"
    
    docker run --rm httpcat-build:arm64 > "$RELEASE_DIR/httpcat-linux-arm64"
    chmod +x "$RELEASE_DIR/httpcat-linux-arm64"
    echo -e "${GREEN}✓ httpcat-linux-arm64 [CGO 启用，支持 SQLite]${NC}"
    
    # 清理临时文件
    rm -f "$PROJECT_ROOT/.build.Dockerfile"
    docker rmi httpcat-build:amd64 httpcat-build:arm64 2>/dev/null || true
}

# ============ 打包发布文件 ============
package_release() {
    echo -e "${BLUE}打包发布文件...${NC}"
    
    cd "$PROJECT_ROOT"
    
    # 复制配置和资源文件
    cp -r "$SERVER_DIR/internal/conf" "$RELEASE_DIR/"
    cp -r "$PROJECT_ROOT/static" "$RELEASE_DIR/"
    cp "$PROJECT_ROOT/httpcat.service" "$RELEASE_DIR/"
    # 使用用户安装指南作为发布包 README（而非开发者 README）
    cp "$PROJECT_ROOT/docs/INSTALL.md" "$RELEASE_DIR/README.md"
    cp "$SCRIPTS_DIR/install.sh" "$RELEASE_DIR/"
    cp "$SCRIPTS_DIR/uninstall.sh" "$RELEASE_DIR/"
    
    chmod +x "$RELEASE_DIR/install.sh"
    chmod +x "$RELEASE_DIR/uninstall.sh"
    
    # 为每个平台创建独立压缩包
    cd "$RELEASE_DIR"
    
    for binary in httpcat-*; do
        [ ! -f "$binary" ] && continue
        
        local platform_name="${binary#httpcat-}"
        platform_name="${platform_name%.exe}"
        local package_name="httpcat_${HTTPCAT_VERSION}_${platform_name}"
        
        echo -e "  打包 $package_name..."
        
        mkdir -p "$package_name"
        cp "$binary" "$package_name/httpcat${binary##*httpcat-*}"
        [ "${binary##*.}" = "exe" ] && mv "$package_name/httpcat${binary##*httpcat-*}" "$package_name/httpcat.exe"
        [ "${binary##*.}" != "exe" ] && mv "$package_name/httpcat${binary##*httpcat-*}" "$package_name/httpcat"
        
        cp -r conf "$package_name/"
        cp -r static "$package_name/"
        cp README.md "$package_name/"
        
        # Linux/macOS 版本添加安装脚本
        if [[ "$platform_name" != windows* ]]; then
            cp httpcat.service "$package_name/"
            cp install.sh "$package_name/"
            cp uninstall.sh "$package_name/"
            chmod +x "$package_name/httpcat"
            chmod +x "$package_name/install.sh"
            chmod +x "$package_name/uninstall.sh"
        fi
        
        # 压缩（排除 macOS 扩展属性）
        if [[ "$platform_name" == windows* ]]; then
            zip -rq "${package_name}.zip" "$package_name"
            echo -e "${GREEN}  ✓ ${package_name}.zip${NC}"
        else
            tar -czf "${package_name}.tar.gz" "$package_name"
            echo -e "${GREEN}  ✓ ${package_name}.tar.gz${NC}"
        fi
        
        rm -rf "$package_name"
    done
    
    # 清理临时文件
    rm -rf conf static httpcat.service README.md install.sh uninstall.sh
    
    cd "$PROJECT_ROOT"
}

# ============ 更新版本号 ============
update_version() {
    echo -e "${BLUE}更新版本号到 $HTTPCAT_VERSION...${NC}"
    
    if [ "$OS_TYPE" = "Darwin" ]; then
        sed -i '' "s/httpcat_version=\".*\"/httpcat_version=\"$HTTPCAT_VERSION\"/g" README.md 2>/dev/null || true
        sed -i '' "s/httpcat_version=\".*\"/httpcat_version=\"$HTTPCAT_VERSION\"/g" translations/README-cn.md 2>/dev/null || true
    else
        sed -i "s/httpcat_version=\".*\"/httpcat_version=\"$HTTPCAT_VERSION\"/g" README.md 2>/dev/null || true
        sed -i "s/httpcat_version=\".*\"/httpcat_version=\"$HTTPCAT_VERSION\"/g" translations/README-cn.md 2>/dev/null || true
    fi
}

# ============ 显示构建结果 ============
show_result() {
    echo ""
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}       构建完成！版本: $HTTPCAT_VERSION${NC}"
    echo -e "${GREEN}========================================${NC}"
    echo ""
    echo -e "${BLUE}构建产物:${NC}"
    ls -lh "$RELEASE_DIR"/*.tar.gz "$RELEASE_DIR"/*.zip 2>/dev/null || ls -lh "$RELEASE_DIR"/httpcat-*
    echo ""
    echo -e "${YELLOW}使用说明:${NC}"
    echo "  1. 解压对应平台的压缩包"
    echo "  2. Linux/macOS: ./install.sh 或直接运行 ./httpcat"
    echo "  3. Windows: 直接运行 httpcat.exe"
    echo "  4. 默认端口: 8888, 默认账号: admin/admin"
}

# ============ 主流程 ============
main() {
    # 解析参数
    BUILD_FRONTEND="ask"
    BUILD_PLATFORMS=()
    USE_DOCKER=false
    CLEAN_ONLY=false
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                exit 0
                ;;
            -v|--version)
                HTTPCAT_VERSION="$2"
                shift 2
                ;;
            -p|--platform)
                BUILD_PLATFORMS+=("$2")
                shift 2
                ;;
            -a|--all)
                BUILD_PLATFORMS=("${PLATFORMS[@]}")
                shift
                ;;
            -f|--frontend)
                BUILD_FRONTEND="yes"
                shift
                ;;
            -s|--skip-frontend)
                BUILD_FRONTEND="no"
                shift
                ;;
            -d|--docker)
                USE_DOCKER=true
                shift
                ;;
            --clean)
                CLEAN_ONLY=true
                shift
                ;;
            *)
                echo -e "${RED}未知参数: $1${NC}"
                show_help
                exit 1
                ;;
        esac
    done
    
    # 检测系统
    detect_os
    
    # 只清理
    if [ "$CLEAN_ONLY" = true ]; then
        clean_release
        exit 0
    fi
    
    # 环境检查
    check_environment
    
    # 清理并创建目录
    clean_release
    
    # 前端构建
    if [ "$BUILD_FRONTEND" = "ask" ]; then
        read -p "是否构建前端？(y/n) [y]: " answer
        answer=${answer:-y}
        [ "$answer" = "y" ] && build_frontend
    elif [ "$BUILD_FRONTEND" = "yes" ]; then
        build_frontend
    fi
    
    # Go 依赖
    prepare_go_deps
    
    # 如果没有指定平台，使用默认的三个主要平台
    if [ ${#BUILD_PLATFORMS[@]} -eq 0 ]; then
        BUILD_PLATFORMS=("linux_amd64" "linux_arm64" "darwin_arm64")
        
        # 交互式选择
        echo ""
        echo -e "${BLUE}选择构建平台:${NC}"
        echo "  1) Linux x86_64 + Linux ARM64 + macOS ARM64 (推荐)"
        echo "  2) 所有平台 (Linux/macOS/Windows, amd64/arm64)"
        echo "  3) 仅 Linux (amd64 + arm64)"
        echo "  4) 仅当前平台 (${HOST_OS}_${HOST_ARCH})"
        read -p "请选择 [1]: " choice
        choice=${choice:-1}
        
        case $choice in
            1) BUILD_PLATFORMS=("linux_amd64" "linux_arm64" "darwin_arm64") ;;
            2) BUILD_PLATFORMS=("${PLATFORMS[@]}") ;;
            3) BUILD_PLATFORMS=("linux_amd64" "linux_arm64") ;;
            4) BUILD_PLATFORMS=("${HOST_OS}_${HOST_ARCH}") ;;
        esac
    fi
    
    echo ""
    echo -e "${BLUE}将构建以下平台:${NC}"
    printf '  - %s\n' "${BUILD_PLATFORMS[@]}"
    echo ""
    
    # 检查是否需要 Docker（macOS 交叉编译 Linux 时）
    local need_docker_for_linux=false
    if [ "$HOST_OS" = "darwin" ]; then
        for platform in "${BUILD_PLATFORMS[@]}"; do
            if [[ "$platform" == linux_* ]]; then
                need_docker_for_linux=true
                break
            fi
        done
    fi
    
    # macOS 上编译 Linux 版本时，自动建议使用 Docker
    if [ "$need_docker_for_linux" = true ] && [ "$USE_DOCKER" != true ]; then
        echo -e "${YELLOW}⚠ 检测到在 macOS 上交叉编译 Linux 版本${NC}"
        echo -e "${YELLOW}  交叉编译的 Linux 版本将不支持 SQLite（无法使用用户登录等功能）${NC}"
        if [ "$DOCKER_AVAILABLE" = true ]; then
            echo -e "${YELLOW}  建议使用 Docker 构建以获得完整功能${NC}"
            # 非交互式环境默认使用 Docker
            if [ -t 0 ]; then
                read -p "是否使用 Docker 构建 Linux 版本？(y/n) [y]: " use_docker_answer < /dev/tty
                use_docker_answer=${use_docker_answer:-y}
                if [ "$use_docker_answer" = "y" ] || [ "$use_docker_answer" = "Y" ]; then
                    USE_DOCKER=true
                fi
            else
                echo -e "${YELLOW}  非交互式环境，自动使用 Docker 构建${NC}"
                USE_DOCKER=true
            fi
        else
            echo -e "${YELLOW}  提示：安装 Docker 后可使用 -d 参数获得完整 SQLite 支持${NC}"
        fi
        echo ""
    fi
    
    # Docker 构建 Linux 版本
    if [ "$USE_DOCKER" = true ]; then
        build_with_docker
        
        # 构建非 Linux 平台（本机）
        for platform in "${BUILD_PLATFORMS[@]}"; do
            [[ "$platform" == linux_* ]] && continue
            build_platform "$platform"
        done
    else
        # 普通构建（Linux 本机编译或明确不使用 Docker）
        for platform in "${BUILD_PLATFORMS[@]}"; do
            build_platform "$platform"
        done
    fi
    
    # 更新版本号
    update_version
    
    # 打包
    package_release
    
    # 显示结果
    show_result
}

# 执行主流程
main "$@"
