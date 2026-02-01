# HttpCat 安装使用指南

## 📦 发布包内容

```
httpcat_vX.X.X_linux-amd64/
├── httpcat              # 可执行文件
├── conf/                # 配置文件目录
│   └── svr.yml          # 主配置文件
├── static/              # Web 界面静态资源
├── install.sh           # 安装脚本 (Linux/macOS)
├── uninstall.sh         # 卸载脚本 (Linux/macOS)
├── httpcat.service      # systemd 服务文件
└── README.md            # 本文档
```

## 🚀 快速启动

### 方式一：使用安装脚本（推荐）

```bash
# 解压
tar -zxvf httpcat_vX.X.X_linux-amd64.tar.gz
cd httpcat_vX.X.X_linux-amd64

# 安装到系统
sudo ./install.sh

# 启动服务
sudo systemctl start httpcat

# 查看状态
sudo systemctl status httpcat
```

### 方式二：直接运行（免安装）

```bash
# Linux/macOS
chmod +x httpcat
./httpcat --port=8888 -C conf/svr.yml

# Windows
httpcat.exe --port=8888 -C conf/svr.yml
```

### 方式三：后台运行

```bash
# Linux/macOS
nohup ./httpcat --port=8888 -C conf/svr.yml > httpcat.log 2>&1 &

# 查看日志
tail -f httpcat.log
```

## 🔐 默认账号

| 项目 | 值 |
|------|------|
| **管理地址** | http://localhost:8888 |
| **用户名** | `admin` |
| **密码** | `admin` |

⚠️ **安全提示**：首次登录后请立即修改默认密码！

## 📁 安装后目录结构

使用 `install.sh` 安装后，文件将按照 Linux FHS 标准分布：

```
/usr/local/bin/
└── httpcat                         # 可执行文件

/etc/httpcat/
└── svr.yml                         # 配置文件

/var/log/httpcat/
└── httpcat.log                     # 日志文件

/var/lib/httpcat/
├── static/                         # Web 界面静态资源
├── upload/                         # 上传文件存储目录
├── download/                       # 下载文件缓存目录
└── data/
    └── httpcat_sqlite.db           # SQLite 数据库
```

### 目录说明

| 目录 | 路径 | 用途 |
|------|------|------|
| **可执行文件** | `/usr/local/bin/httpcat` | 主程序 |
| **配置文件** | `/etc/httpcat/svr.yml` | 服务配置 |
| **日志目录** | `/var/log/httpcat/` | 运行日志 |
| **数据目录** | `/var/lib/httpcat/` | 应用数据根目录 |
| **静态资源** | `/var/lib/httpcat/static/` | Web 管理界面 |
| **上传目录** | `/var/lib/httpcat/upload/` | 用户上传的文件 |
| **下载目录** | `/var/lib/httpcat/download/` | 下载缓存 |
| **数据库** | `/var/lib/httpcat/data/` | SQLite 数据库 |

## ⚙️ 安装脚本选项

```bash
# 查看帮助
./install.sh -h

# 默认安装
sudo ./install.sh

# 指定端口
sudo ./install.sh -p 9000

# 自定义安装前缀
sudo ./install.sh --prefix /opt

# 不安装 systemd 服务
sudo ./install.sh --no-service
```

## 🔧 命令行参数

```bash
./httpcat -h

选项:
  --port, -p         服务端口 (默认: 8888)
  -C                 配置文件路径 (默认: ./conf/svr.yml)
  --static           静态资源目录 (默认: ./static/)
  --upload           上传目录 (默认: ./upload/)
  --download         下载目录 (默认: ./download/)
  --p2pport          P2P 监听端口
  -v                 显示版本信息
```

## 🛠️ 服务管理

```bash
# 启动服务
sudo systemctl start httpcat

# 停止服务
sudo systemctl stop httpcat

# 重启服务
sudo systemctl restart httpcat

# 查看状态
sudo systemctl status httpcat

# 开机自启
sudo systemctl enable httpcat

# 取消开机自启
sudo systemctl disable httpcat

# 查看日志
sudo journalctl -u httpcat -f
# 或
tail -f /var/log/httpcat/httpcat.log
```

## ❌ 卸载

### 使用卸载脚本（推荐）

```bash
# 标准卸载（保留配置和数据）
sudo ./uninstall.sh

# 完全卸载（删除所有配置和数据）
sudo ./uninstall.sh --purge

# 完全卸载但保留用户上传的文件
sudo ./uninstall.sh --purge --keep-data

# 无需确认（用于自动化脚本）
sudo ./uninstall.sh -y
```

### 手动卸载

```bash
# 停止并禁用服务
sudo systemctl stop httpcat
sudo systemctl disable httpcat

# 删除程序和服务文件
sudo rm /usr/local/bin/httpcat
sudo rm /etc/systemd/system/httpcat.service
sudo systemctl daemon-reload

# 删除配置文件（可选）
sudo rm -rf /etc/httpcat

# 删除数据文件（谨慎操作！）
sudo rm -rf /var/lib/httpcat
sudo rm -rf /var/log/httpcat
```

## ❓ 常见问题

### Q: 登录提示"账号或密码错误"？

首先确认默认账号密码：`admin` / `admin`

如果仍然无法登录，可能是使用了**不支持 SQLite 的版本**。检查方法：

```bash
# 查看启动日志
./httpcat --port=8888 -C conf/svr.yml 2>&1 | grep -i "sqlite\|CGO"
```

如果看到 SQLite 相关错误，请下载支持 SQLite 的版本：
- **Linux**：使用 Docker 构建或在 Linux 服务器上编译
- **macOS 交叉编译**：需使用 `./build.sh -d` 启用 Docker 构建

### Q: 端口被占用？

```bash
# 查看端口占用
lsof -i :8888  # Linux/macOS
netstat -ano | findstr 8888  # Windows

# 使用其他端口
./httpcat --port=9999 -C conf/svr.yml

# 或修改安装时的端口
sudo ./install.sh -p 9999
```

### Q: 如何修改密码？

登录管理界面后，点击右上角用户头像 → 个人设置 → 修改密码

### Q: 忘记密码怎么办？

删除数据库文件后重启服务，将自动创建默认账号：

```bash
# 找到数据库文件
sudo find /var/lib/httpcat -name "*.db"

# 删除数据库
sudo rm /var/lib/httpcat/data/httpcat_sqlite.db

# 重启服务
sudo systemctl restart httpcat
```

### Q: 如何备份数据？

```bash
# 备份上传的文件和数据库
sudo tar -czf httpcat-backup-$(date +%Y%m%d).tar.gz \
  /var/lib/httpcat/upload \
  /var/lib/httpcat/data \
  /etc/httpcat/svr.yml
```

### Q: 如何迁移到其他服务器？

```bash
# 在旧服务器上备份
sudo tar -czf httpcat-full-backup.tar.gz \
  /var/lib/httpcat \
  /etc/httpcat

# 在新服务器上恢复
sudo tar -xzf httpcat-full-backup.tar.gz -C /

# 重新安装可执行文件
sudo ./install.sh
```

## 📞 获取帮助

- GitHub: https://github.com/puge/httpcat
- Issues: https://github.com/puge/httpcat/issues
