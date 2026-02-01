#!/bin/bash
# HttpCat 安装脚本
# 用法: sudo ./install.sh [选项]

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
LOG_DIR="/var/log/httpcat"

# 显示帮助
show_help() {
    echo "HttpCat 安装脚本"
    echo ""
    echo "用法: $0 [选项]"
    echo ""
    echo "选项:"
    echo "  -h, --help          显示帮助信息"
    echo "  -p, --port PORT     指定服务端口 (默认: 8888)"
    echo "  --prefix DIR        自定义安装目录前缀 (默认: /usr/local)"
    echo "  --no-service        不安装 systemd 服务"
    echo ""
    echo "示例:"
    echo "  sudo ./install.sh              # 默认安装"
    echo "  sudo ./install.sh -p 9000      # 指定端口"
    echo "  sudo ./install.sh --no-service # 不安装服务"
    echo ""
}

# 默认参数
PORT=""
INSTALL_SERVICE=true
PREFIX="/usr/local"

# 解析参数
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
        --prefix)
            PREFIX="$2"
            INSTALL_DIR="$PREFIX/bin"
            shift 2
            ;;
        --no-service)
            INSTALL_SERVICE=false
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

# 检测系统类型
detect_system() {
    if [ -f /etc/os-release ]; then
        . /etc/os-release
        OS_NAME=$NAME
    elif [ "$(uname)" = "Darwin" ]; then
        OS_NAME="macOS"
    else
        OS_NAME="Unknown"
    fi
    echo -e "${BLUE}检测到系统: $OS_NAME${NC}"
}

# 检测 systemd
has_systemd() {
    if command -v systemctl &> /dev/null && [ -d /run/systemd/system ]; then
        return 0
    fi
    return 1
}

# 安装函数
install_httpcat() {
    echo ""
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}       开始安装 HttpCat${NC}"
    echo -e "${BLUE}========================================${NC}"
    echo ""
    
    # 检测系统
    detect_system
    
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
    echo -e "${BLUE}[1/5] 创建目录结构...${NC}"
    mkdir -p "$INSTALL_DIR"
    mkdir -p "$CONFIG_DIR"
    mkdir -p "$LOG_DIR"
    mkdir -p "$DATA_DIR/static"
    mkdir -p "$DATA_DIR/upload"
    mkdir -p "$DATA_DIR/download"
    mkdir -p "$DATA_DIR/data"
    echo -e "${GREEN}✓ 目录创建完成${NC}"
    
    # 复制可执行文件
    echo -e "${BLUE}[2/5] 安装可执行文件...${NC}"
    cp -f ./httpcat "$INSTALL_DIR/httpcat"
    chmod +x "$INSTALL_DIR/httpcat"
    echo -e "${GREEN}✓ 可执行文件已安装到 $INSTALL_DIR/httpcat${NC}"
    
    # 复制配置文件
    echo -e "${BLUE}[3/5] 安装配置文件...${NC}"
    if [ -f "$CONFIG_DIR/svr.yml" ]; then
        echo -e "${YELLOW}! 配置文件已存在，保留现有配置${NC}"
        cp -f ./conf/svr.yml "$CONFIG_DIR/svr.yml.new"
        echo -e "${YELLOW}  新配置已保存为 $CONFIG_DIR/svr.yml.new${NC}"
    else
        cp -f ./conf/svr.yml "$CONFIG_DIR/svr.yml"
        
        # 更新配置文件中的路径
        if [ "$(uname)" = "Darwin" ]; then
            # macOS
            sed -i '' "s|path: ./log/httpcat.log|path: $LOG_DIR/httpcat.log|g" "$CONFIG_DIR/svr.yml"
            sed -i '' "s|sqlite_db_path: \"./data/httpcat_sqlite.db\"|sqlite_db_path: \"$DATA_DIR/data/httpcat_sqlite.db\"|g" "$CONFIG_DIR/svr.yml"
        else
            # Linux
            sed -i "s|path: ./log/httpcat.log|path: $LOG_DIR/httpcat.log|g" "$CONFIG_DIR/svr.yml"
            sed -i "s|sqlite_db_path: \"./data/httpcat_sqlite.db\"|sqlite_db_path: \"$DATA_DIR/data/httpcat_sqlite.db\"|g" "$CONFIG_DIR/svr.yml"
        fi
        
        echo -e "${GREEN}✓ 配置文件已安装到 $CONFIG_DIR/svr.yml${NC}"
    fi
    
    # 复制静态资源
    echo -e "${BLUE}[4/5] 安装静态资源...${NC}"
    cp -rf ./static/* "$DATA_DIR/static/"
    echo -e "${GREEN}✓ 静态资源已安装到 $DATA_DIR/static/${NC}"
    
    # 修改端口（如果指定）
    if [ -n "$PORT" ]; then
        echo -e "${BLUE}设置端口为: $PORT${NC}"
        if [ "$(uname)" = "Darwin" ]; then
            sed -i '' "s/port: 8888/port: $PORT/" "$CONFIG_DIR/svr.yml"
        else
            sed -i "s/port: 8888/port: $PORT/" "$CONFIG_DIR/svr.yml"
        fi
    fi
    
    # 安装 systemd 服务
    echo -e "${BLUE}[5/5] 配置系统服务...${NC}"
    if [ "$INSTALL_SERVICE" = true ] && has_systemd; then
        cat > /etc/systemd/system/httpcat.service << EOF
[Unit]
Description=HttpCat File Server
Documentation=https://github.com/puge/httpcat
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=$DATA_DIR
ExecStart=$INSTALL_DIR/httpcat -C $CONFIG_DIR/svr.yml --static=$DATA_DIR/static/ --upload=$DATA_DIR/upload/ --download=$DATA_DIR/download/
Restart=on-failure
RestartSec=5
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
EOF
        
        # 重新加载 systemd
        systemctl daemon-reload
        
        # 启用服务
        systemctl enable httpcat
        
        echo -e "${GREEN}✓ systemd 服务已安装并启用${NC}"
    elif [ "$INSTALL_SERVICE" = true ]; then
        echo -e "${YELLOW}! 系统不支持 systemd，跳过服务安装${NC}"
        echo -e "${YELLOW}  您可以手动启动: $INSTALL_DIR/httpcat -C $CONFIG_DIR/svr.yml${NC}"
    else
        echo -e "${YELLOW}! 已跳过服务安装 (--no-service)${NC}"
    fi
    
    # 显示安装结果
    echo ""
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}       HttpCat 安装成功！${NC}"
    echo -e "${GREEN}========================================${NC}"
    echo ""
    echo -e "${BLUE}安装路径:${NC}"
    echo "  可执行文件: $INSTALL_DIR/httpcat"
    echo "  配置文件:   $CONFIG_DIR/svr.yml"
    echo "  日志目录:   $LOG_DIR/"
    echo "  数据目录:   $DATA_DIR/"
    echo "    ├── static/    静态资源"
    echo "    ├── upload/    上传文件"
    echo "    ├── download/  下载文件"
    echo "    └── data/      数据库文件"
    echo ""
    
    if [ "$INSTALL_SERVICE" = true ] && has_systemd; then
        echo -e "${BLUE}服务管理:${NC}"
        echo "  启动服务: sudo systemctl start httpcat"
        echo "  停止服务: sudo systemctl stop httpcat"
        echo "  重启服务: sudo systemctl restart httpcat"
        echo "  查看状态: sudo systemctl status httpcat"
        echo "  查看日志: sudo journalctl -u httpcat -f"
        echo ""
    fi
    
    echo -e "${BLUE}访问地址:${NC}"
    local display_port=${PORT:-8888}
    echo "  http://localhost:$display_port"
    echo ""
    echo -e "${BLUE}默认账号:${NC}"
    echo "  用户名: admin"
    echo "  密码:   admin"
    echo ""
    echo -e "${YELLOW}安全提示: 首次登录后请立即修改默认密码！${NC}"
    echo ""
    
    if [ "$INSTALL_SERVICE" = true ] && has_systemd; then
        echo -e "${BLUE}立即启动服务？${NC}"
        read -p "是否立即启动 HttpCat? (y/n) [y]: " start_now
        start_now=${start_now:-y}
        if [ "$start_now" = "y" ] || [ "$start_now" = "Y" ]; then
            systemctl start httpcat
            echo -e "${GREEN}✓ 服务已启动${NC}"
            sleep 1
            systemctl status httpcat --no-pager || true
        fi
    fi
}

# 主逻辑
check_root
install_httpcat
