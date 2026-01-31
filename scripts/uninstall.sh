#!/bin/bash
# HttpCat 卸载脚本

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# 安装路径
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="/etc/httpcat"
DATA_DIR="/var/lib/httpcat"

echo -e "${BLUE}正在卸载 HttpCat...${NC}"

# 检查是否为 root
if [ "$EUID" -ne 0 ]; then
    echo -e "${RED}请使用 sudo 运行此脚本${NC}"
    exit 1
fi

# 停止服务
if systemctl is-active --quiet httpcat 2>/dev/null; then
    echo -e "${BLUE}停止服务...${NC}"
    systemctl stop httpcat
    echo -e "${GREEN}✓ 服务已停止${NC}"
fi

# 禁用服务
if systemctl is-enabled --quiet httpcat 2>/dev/null; then
    echo -e "${BLUE}禁用服务...${NC}"
    systemctl disable httpcat
    echo -e "${GREEN}✓ 服务已禁用${NC}"
fi

# 删除文件
echo -e "${BLUE}删除文件...${NC}"
rm -f "$INSTALL_DIR/httpcat"
rm -f /etc/systemd/system/httpcat.service
systemctl daemon-reload 2>/dev/null || true

echo ""
echo -e "${YELLOW}! 保留以下目录（如需完全删除请手动执行）:${NC}"
echo "  - $CONFIG_DIR (配置文件)"
echo "  - $DATA_DIR (数据文件)"
echo ""
echo -e "${GREEN}✓ HttpCat 卸载完成${NC}"
