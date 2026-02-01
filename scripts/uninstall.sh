#!/bin/bash
# HttpCat 卸载脚本
# 用法: sudo ./uninstall.sh [选项]

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# 安装路径（与 install.sh 保持一致）
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="/etc/httpcat"
DATA_DIR="/var/lib/httpcat"
LOG_DIR="/var/log/httpcat"

# 显示帮助
show_help() {
    echo "HttpCat 卸载脚本"
    echo ""
    echo "用法: $0 [选项]"
    echo ""
    echo "选项:"
    echo "  -h, --help       显示帮助信息"
    echo "  --purge          完全卸载（包括配置和数据文件）"
    echo "  --keep-data      保留上传的文件数据"
    echo "  -y, --yes        跳过确认提示"
    echo ""
    echo "示例:"
    echo "  sudo ./uninstall.sh           # 标准卸载（保留配置和数据）"
    echo "  sudo ./uninstall.sh --purge   # 完全卸载"
    echo "  sudo ./uninstall.sh -y        # 无需确认"
    echo ""
}

# 默认参数
PURGE=false
KEEP_DATA=false
SKIP_CONFIRM=false

# 解析参数
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        --purge)
            PURGE=true
            shift
            ;;
        --keep-data)
            KEEP_DATA=true
            shift
            ;;
        -y|--yes)
            SKIP_CONFIRM=true
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
if [ "$EUID" -ne 0 ]; then
    echo -e "${RED}请使用 sudo 运行此脚本${NC}"
    exit 1
fi

# 检查是否已安装
check_installed() {
    if [ ! -f "$INSTALL_DIR/httpcat" ] && [ ! -d "$CONFIG_DIR" ] && [ ! -d "$DATA_DIR" ] && [ ! -d "$LOG_DIR" ]; then
        echo -e "${YELLOW}HttpCat 似乎未安装在标准位置${NC}"
        echo "  检查路径: $INSTALL_DIR/httpcat"
        echo "  检查路径: $CONFIG_DIR"
        echo "  检查路径: $DATA_DIR"
        echo "  检查路径: $LOG_DIR"
        exit 0
    fi
}

# 确认卸载
confirm_uninstall() {
    if [ "$SKIP_CONFIRM" = true ]; then
        return 0
    fi
    
    echo -e "${YELLOW}========================================${NC}"
    echo -e "${YELLOW}       即将卸载 HttpCat${NC}"
    echo -e "${YELLOW}========================================${NC}"
    echo ""
    
    if [ "$PURGE" = true ]; then
        echo -e "${RED}警告: 这将删除所有配置和数据！${NC}"
        if [ "$KEEP_DATA" = true ]; then
            echo -e "${YELLOW}  (--keep-data: 将保留上传文件)${NC}"
        fi
    else
        echo "将删除:"
        echo "  - 可执行文件: $INSTALL_DIR/httpcat"
        echo "  - systemd 服务配置"
        echo ""
        echo "将保留:"
        echo "  - 配置目录: $CONFIG_DIR"
        echo "  - 数据目录: $DATA_DIR"
    fi
    echo ""
    
    read -p "确认卸载? (y/n) [n]: " confirm
    if [ "$confirm" != "y" ] && [ "$confirm" != "Y" ]; then
        echo "已取消卸载"
        exit 0
    fi
}

# 停止服务
stop_service() {
    echo -e "${BLUE}[1/4] 停止服务...${NC}"
    
    if command -v systemctl &> /dev/null; then
        if systemctl is-active --quiet httpcat 2>/dev/null; then
            systemctl stop httpcat
            echo -e "${GREEN}✓ 服务已停止${NC}"
        else
            echo -e "${YELLOW}  服务未运行${NC}"
        fi
        
        if systemctl is-enabled --quiet httpcat 2>/dev/null; then
            systemctl disable httpcat
            echo -e "${GREEN}✓ 服务已禁用${NC}"
        fi
    else
        # 尝试终止进程
        if pgrep -x "httpcat" > /dev/null; then
            pkill -x "httpcat" || true
            echo -e "${GREEN}✓ httpcat 进程已终止${NC}"
        else
            echo -e "${YELLOW}  未发现运行中的 httpcat 进程${NC}"
        fi
    fi
}

# 删除文件
remove_files() {
    echo -e "${BLUE}[2/4] 删除程序文件...${NC}"
    
    # 删除可执行文件
    if [ -f "$INSTALL_DIR/httpcat" ]; then
        rm -f "$INSTALL_DIR/httpcat"
        echo -e "${GREEN}✓ 已删除 $INSTALL_DIR/httpcat${NC}"
    fi
    
    # 删除 systemd 服务文件
    if [ -f /etc/systemd/system/httpcat.service ]; then
        rm -f /etc/systemd/system/httpcat.service
        systemctl daemon-reload 2>/dev/null || true
        echo -e "${GREEN}✓ 已删除 systemd 服务配置${NC}"
    fi
}

# 清理配置和数据
cleanup_data() {
    echo -e "${BLUE}[3/4] 清理配置和数据...${NC}"
    
    if [ "$PURGE" = true ]; then
        # 删除配置目录
        if [ -d "$CONFIG_DIR" ]; then
            rm -rf "$CONFIG_DIR"
            echo -e "${GREEN}✓ 已删除配置目录: $CONFIG_DIR${NC}"
        fi
        
        # 删除日志目录
        if [ -d "$LOG_DIR" ]; then
            rm -rf "$LOG_DIR"
            echo -e "${GREEN}✓ 已删除日志目录: $LOG_DIR${NC}"
        fi
        
        # 删除数据目录
        if [ -d "$DATA_DIR" ]; then
            if [ "$KEEP_DATA" = true ]; then
                # 保留上传文件，删除其他
                rm -rf "$DATA_DIR/static" 2>/dev/null || true
                rm -rf "$DATA_DIR/data" 2>/dev/null || true
                echo -e "${GREEN}✓ 已清理数据目录（保留上传文件）${NC}"
                echo -e "${YELLOW}  保留: $DATA_DIR/upload${NC}"
                echo -e "${YELLOW}  保留: $DATA_DIR/download${NC}"
            else
                rm -rf "$DATA_DIR"
                echo -e "${GREEN}✓ 已删除数据目录: $DATA_DIR${NC}"
            fi
        fi
    else
        echo -e "${YELLOW}  保留配置目录: $CONFIG_DIR${NC}"
        echo -e "${YELLOW}  保留日志目录: $LOG_DIR${NC}"
        echo -e "${YELLOW}  保留数据目录: $DATA_DIR${NC}"
        echo -e "${YELLOW}  使用 --purge 选项可完全删除${NC}"
    fi
}

# 显示结果
show_result() {
    echo -e "${BLUE}[4/4] 卸载完成${NC}"
    echo ""
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}       HttpCat 卸载完成${NC}"
    echo -e "${GREEN}========================================${NC}"
    echo ""
    
    if [ "$PURGE" != true ]; then
        echo -e "${YELLOW}保留的目录（如需完全删除请执行）:${NC}"
        [ -d "$CONFIG_DIR" ] && echo "  sudo rm -rf $CONFIG_DIR"
        [ -d "$LOG_DIR" ] && echo "  sudo rm -rf $LOG_DIR"
        [ -d "$DATA_DIR" ] && echo "  sudo rm -rf $DATA_DIR"
        echo ""
    elif [ "$KEEP_DATA" = true ] && [ -d "$DATA_DIR" ]; then
        echo -e "${YELLOW}保留的用户数据:${NC}"
        [ -d "$DATA_DIR/upload" ] && echo "  $DATA_DIR/upload"
        [ -d "$DATA_DIR/download" ] && echo "  $DATA_DIR/download"
        echo ""
        echo "手动删除: sudo rm -rf $DATA_DIR"
        echo ""
    fi
    
    echo -e "${GREEN}感谢使用 HttpCat！${NC}"
}

# 主流程
echo ""
echo -e "${BLUE}HttpCat 卸载程序${NC}"
echo ""

check_installed
confirm_uninstall
stop_service
remove_files
cleanup_data
show_result
