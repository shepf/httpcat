#!/bin/bash
# HttpCat 安装脚本

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
WEB_ROOT="$DATA_DIR/website"

# 显示帮助
show_help() {
    echo "用法: $0 [选项]"
    echo ""
    echo "选项:"
    echo "  -h, --help          显示帮助信息"
    echo "  -p, --port PORT     指定服务端口 (默认: 8888)"
    echo "  --uninstall         卸载 HttpCat"
    echo ""
}

# 解析参数
PORT=""
UNINSTALL=false
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        -p|--port)
            PORT="$2"
            shift 2
            ;;
        --uninstall)
            UNINSTALL=true
            shift
            ;;
        *)
            echo -e "${RED}未知选项: $1${NC}"
            show_help
            exit 1
            ;;
    esac
done

# 检查是否为 root
check_root() {
    if [ "$EUID" -ne 0 ]; then
        echo -e "${RED}请使用 sudo 运行此脚本${NC}"
        exit 1
    fi
}

# 卸载函数
uninstall() {
    echo -e "${BLUE}正在卸载 HttpCat...${NC}"
    
    # 停止服务
    if systemctl is-active --quiet httpcat 2>/dev/null; then
        systemctl stop httpcat
        echo -e "${GREEN}✓ 服务已停止${NC}"
    fi
    
    # 禁用服务
    if systemctl is-enabled --quiet httpcat 2>/dev/null; then
        systemctl disable httpcat
        echo -e "${GREEN}✓ 服务已禁用${NC}"
    fi
    
    # 删除文件
    rm -f "$INSTALL_DIR/httpcat"
    rm -f /etc/systemd/system/httpcat.service
    systemctl daemon-reload 2>/dev/null || true
    
    echo ""
    echo -e "${YELLOW}! 保留以下目录（如需完全删除请手动执行）:${NC}"
    echo "  - $CONFIG_DIR (配置文件)"
    echo "  - $DATA_DIR (数据文件)"
    echo ""
    echo -e "${GREEN}✓ 卸载完成${NC}"
}

# 安装函数
install_httpcat() {
    echo -e "${BLUE}开始安装 HttpCat...${NC}"
    
    # 检查可执行文件
    if [ ! -f "./httpcat" ]; then
        echo -e "${RED}错误: 未找到 httpcat 可执行文件${NC}"
        echo "请确保在解压后的目录中运行此脚本"
        exit 1
    fi
    
    # 检查配置文件
    if [ ! -f "./conf/svr.yml" ]; then
        echo -e "${RED}错误: 未找到 conf/svr.yml 配置文件${NC}"
        exit 1
    fi
    
    # 检查静态资源
    if [ ! -d "./static" ]; then
        echo -e "${RED}错误: 未找到 static 静态资源目录${NC}"
        exit 1
    fi
    
    # 创建目录
    echo -e "${BLUE}创建目录...${NC}"
    mkdir -p "$CONFIG_DIR"
    mkdir -p "$CONFIG_DIR/conf"
    mkdir -p "$WEB_ROOT"
    mkdir -p "$WEB_ROOT/upload"
    mkdir -p "$WEB_ROOT/download"
    mkdir -p "$WEB_ROOT/static"
    
    # 复制可执行文件
    echo -e "${BLUE}安装可执行文件...${NC}"
    cp -f ./httpcat "$INSTALL_DIR/httpcat"
    chmod +x "$INSTALL_DIR/httpcat"
    
    # 复制配置文件
    echo -e "${BLUE}安装配置文件...${NC}"
    if [ -f "$CONFIG_DIR/svr.yml" ]; then
        echo -e "${YELLOW}! 配置文件已存在，保留现有配置${NC}"
        cp -f ./conf/svr.yml "$CONFIG_DIR/svr.yml.example"
    else
        cp -f ./conf/svr.yml "$CONFIG_DIR/svr.yml"
    fi
    
    # 复制静态资源
    echo -e "${BLUE}安装静态资源...${NC}"
    cp -rf ./static/* "$WEB_ROOT/static/"
    
    # 修改端口（如果指定）
    if [ -n "$PORT" ]; then
        echo -e "${BLUE}设置端口为: $PORT${NC}"
        sed -i "s/port: 8888/port: $PORT/" "$CONFIG_DIR/svr.yml"
    fi
    
    # 安装 systemd 服务
    echo -e "${BLUE}安装 systemd 服务...${NC}"
    cat > /etc/systemd/system/httpcat.service << EOF
[Unit]
Description=HttpCat File Server
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=$DATA_DIR
ExecStart=$INSTALL_DIR/httpcat -C $CONFIG_DIR/svr.yml --static=$WEB_ROOT/static/ --upload=$WEB_ROOT/upload/ --download=$WEB_ROOT/download/
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF
    
    # 重新加载 systemd
    systemctl daemon-reload
    
    # 启用服务
    systemctl enable httpcat
    
    echo ""
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}       HttpCat 安装成功！${NC}"
    echo -e "${GREEN}========================================${NC}"
    echo ""
    echo -e "${BLUE}目录结构:${NC}"
    echo "  可执行文件: $INSTALL_DIR/httpcat"
    echo "  配置文件:   $CONFIG_DIR/svr.yml"
    echo "  数据目录:   $DATA_DIR"
    echo "  上传目录:   $WEB_ROOT/upload"
    echo "  静态资源:   $WEB_ROOT/static"
    echo ""
    echo -e "${BLUE}使用说明:${NC}"
    echo "  启动服务: sudo systemctl start httpcat"
    echo "  停止服务: sudo systemctl stop httpcat"
    echo "  查看状态: sudo systemctl status httpcat"
    echo "  查看日志: sudo journalctl -u httpcat -f"
    echo ""
    echo -e "${BLUE}访问地址:${NC}"
    if [ -n "$PORT" ]; then
        echo "  http://localhost:$PORT"
    else
        echo "  http://localhost:8888"
    fi
    echo ""
    echo -e "${BLUE}默认账号:${NC}"
    echo "  用户名: admin"
    echo "  密码: admin"
    echo ""
    echo -e "${YELLOW}注意: 首次登录后请立即修改默认密码！${NC}"
}

# 主逻辑
if [ "$UNINSTALL" = true ]; then
    check_root
    uninstall
else
    check_root
    install_httpcat
fi
