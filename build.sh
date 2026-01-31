#!/bin/bash

set -e  # 遇到错误立即退出

# ============ 颜色定义 ============
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# ============ 版本信息 ============
HTTPCAT_VERSION=v0.2.0
HTTPCAT_BUILD=$(date "+%Y%m%d%H%M")
COMMIT_ID=$(git rev-parse HEAD 2>/dev/null || echo "unknown")

# ============ 检测操作系统 ============
OS_TYPE=$(uname -s)
echo -e "${GREEN}检测到操作系统: $OS_TYPE${NC}"

# ============ 环境检查 ============
if ! type go >/dev/null 2>&1; then
    echo -e "${RED}错误: Go 未安装${NC}"
    exit 1
fi

echo -e "${GREEN}Go 版本: $(go version)${NC}"

# ============ 清理并创建目录 ============
rm -rf release
mkdir -p release

echo "执行 go mod tidy..."
go mod tidy

# ============ 检查前端文件 ============
if [ ! -f "dist.zip" ]; then
    echo -e "${YELLOW}警告: dist.zip 文件不存在${NC}"
    echo "请先执行 npm run build 命令打包前端，或者跳过前端打包继续构建"
    read -p "是否继续构建？(y/n): " continue_build
    if [ "$continue_build" != "y" ]; then
        exit 1
    fi
else
    cp -rf dist.zip release/
fi

# ============ 构建函数 ============
build_binary() {
    local target_os=$1
    local target_arch=$2
    local output_name=$3
    local cc_compiler=${4:-""}
    
    echo -e "${GREEN}构建 $target_os $target_arch 版本...${NC}"
    
    local build_cmd="GOOS=$target_os GOARCH=$target_arch"
    
    # 如果需要 CGO（SQLite 依赖）
    if [ -n "$cc_compiler" ]; then
        build_cmd="CC=$cc_compiler CGO_ENABLED=1 $build_cmd"
    else
        # x86_64 本机编译，启用 CGO
        if [ "$target_arch" = "amd64" ] && [ "$target_os" = "linux" ]; then
            build_cmd="CGO_ENABLED=1 $build_cmd"
        fi
    fi
    
    eval $build_cmd go build \
        -ldflags \"-s -w -X gin_web_demo/server/common.Version=$HTTPCAT_VERSION -X gin_web_demo/server/common.Build=$HTTPCAT_BUILD -X gin_web_demo/server/common.Commit=$COMMIT_ID\" \
        -o ./release/$output_name ./cmd/httpcat.go
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ $output_name 构建成功${NC}"
    else
        echo -e "${RED}✗ $output_name 构建失败${NC}"
        return 1
    fi
}

# ============ 构建 Linux x86_64 版本 ============
# 注意：如果在 macOS 上交叉编译 Linux 且需要 CGO，需要安装交叉编译工具链
# 可以使用 Docker 来构建，或者禁用 CGO（但会失去 SQLite 支持）

if [ "$OS_TYPE" = "Darwin" ]; then
    echo -e "${YELLOW}macOS 交叉编译 Linux 需要禁用 CGO 或使用 Docker${NC}"
    echo "建议使用 Docker 构建: docker build -t httpcat:latest ."
    
    # macOS 上禁用 CGO 构建（将无法使用 SQLite）
    # 如果需要 SQLite 支持，请使用 Docker 构建
    echo -e "${YELLOW}使用禁用 CGO 模式构建（无 SQLite 支持）...${NC}"
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
        -ldflags "-s -w -X gin_web_demo/server/common.Version=$HTTPCAT_VERSION -X gin_web_demo/server/common.Build=$HTTPCAT_BUILD -X gin_web_demo/server/common.Commit=$COMMIT_ID" \
        -o ./release/httpcat-linux-x86 ./cmd/httpcat.go && \
        echo -e "${GREEN}✓ httpcat-linux-x86 构建成功（无 SQLite）${NC}"
    
    # ARM64 版本
    CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build \
        -ldflags "-s -w -X gin_web_demo/server/common.Version=$HTTPCAT_VERSION -X gin_web_demo/server/common.Build=$HTTPCAT_BUILD -X gin_web_demo/server/common.Commit=$COMMIT_ID" \
        -o ./release/httpcat-linux-aarch64 ./cmd/httpcat.go && \
        echo -e "${GREEN}✓ httpcat-linux-aarch64 构建成功（无 SQLite）${NC}"
        
    # macOS 本机版本（启用 CGO，支持 SQLite）
    echo -e "${GREEN}构建 macOS 本机版本（支持 SQLite）...${NC}"
    CGO_ENABLED=1 go build \
        -ldflags "-s -w -X gin_web_demo/server/common.Version=$HTTPCAT_VERSION -X gin_web_demo/server/common.Build=$HTTPCAT_BUILD -X gin_web_demo/server/common.Commit=$COMMIT_ID" \
        -o ./release/httpcat-darwin ./cmd/httpcat.go && \
        echo -e "${GREEN}✓ httpcat-darwin 构建成功${NC}"
else
    # Linux 环境
    echo "构建 Linux x86_64 版本..."
    CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build \
        -ldflags "-s -w -X gin_web_demo/server/common.Version=$HTTPCAT_VERSION -X gin_web_demo/server/common.Build=$HTTPCAT_BUILD -X gin_web_demo/server/common.Commit=$COMMIT_ID" \
        -o ./release/httpcat-linux-x86 ./cmd/httpcat.go && \
        echo -e "${GREEN}✓ httpcat-linux-x86 构建成功${NC}"

    # Linux ARM64 版本（需要交叉编译工具链）
    if command -v aarch64-linux-gnu-gcc &> /dev/null; then
        echo "构建 Linux ARM64 版本..."
        CC=aarch64-linux-gnu-gcc CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build \
            -ldflags "-s -w -X gin_web_demo/server/common.Version=$HTTPCAT_VERSION -X gin_web_demo/server/common.Build=$HTTPCAT_BUILD -X gin_web_demo/server/common.Commit=$COMMIT_ID" \
            -o ./release/httpcat-linux-aarch64 ./cmd/httpcat.go && \
            echo -e "${GREEN}✓ httpcat-linux-aarch64 构建成功${NC}"
    else
        echo -e "${YELLOW}跳过 ARM64 构建：未找到 aarch64-linux-gnu-gcc${NC}"
        echo "安装方法: sudo apt-get install gcc-aarch64-linux-gnu"
    fi
fi

# ============ 更新 README 版本号 ============
echo "更新 README 版本号..."
if [ "$OS_TYPE" = "Darwin" ]; then
    # macOS sed 需要 -i '' 
    sed -i '' "s/httpcat_version=\".*\"/httpcat_version=\"$HTTPCAT_VERSION\"/g" README.md 2>/dev/null || true
    sed -i '' "s/httpcat_version=\".*\"/httpcat_version=\"$HTTPCAT_VERSION\"/g" translations/README-cn.md 2>/dev/null || true
else
    # Linux sed
    sed -i "s/httpcat_version=\".*\"/httpcat_version=\"$HTTPCAT_VERSION\"/g" README.md
    sed -i "s/httpcat_version=\".*\"/httpcat_version=\"$HTTPCAT_VERSION\"/g" translations/README-cn.md
fi

# ============ 打包发布文件 ============
echo "打包发布文件..."
cp -r server/conf release/
cp -r static release/
cp -r httpcat.service release/
cp -r README.md release/
mkdir -p release/translations
cp -rf translations/* release/translations/
cp -r install.sh release/
cp -r uninstall.sh release/
chmod +x release/install.sh
chmod +x release/uninstall.sh

# 创建发布压缩包
tar zcvf release/httpcat_$HTTPCAT_VERSION.tar.gz -C release .

# ============ 构建完成 ============
echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}构建完成！版本: $HTTPCAT_VERSION${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo "构建产物:"
ls -lh release/

echo ""
echo -e "${YELLOW}提示：${NC}"
echo "1. 如需带 SQLite 的 Linux 版本，请使用 Docker 构建:"
echo "   docker build -t httpcat:latest ."
echo ""
echo "2. 发布包位置: release/httpcat_$HTTPCAT_VERSION.tar.gz"
