# HttpCat 安装使用指南

## 📦 发布包内容

```
httpcat/
├── httpcat              # 可执行文件
├── conf/                # 配置文件目录
│   └── svr.yml          # 主配置文件
├── static/              # Web 界面静态资源
├── install.sh           # 安装脚本 (Linux)
├── uninstall.sh         # 卸载脚本 (Linux)
├── httpcat.service      # systemd 服务文件
└── README.md            # 本文档
```

## 🚀 快速启动

### 方式一：直接运行

```bash
# Linux/macOS
chmod +x httpcat
./httpcat --port=8888 -C conf/svr.yml

# Windows
httpcat.exe --port=8888 -C conf/svr.yml
```

### 方式二：使用安装脚本（推荐 Linux）

```bash
# 安装到系统
sudo ./install.sh

# 启动服务
sudo systemctl start httpcat
sudo systemctl enable httpcat  # 开机自启

# 查看状态
sudo systemctl status httpcat
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

## ⚙️ 配置说明

编辑 `conf/svr.yml` 文件：

```yaml
# 服务端口
port: 8888

# 数据存储目录
data_dir: ./data

# 文件上传目录
upload_dir: ./upload

# 日志配置
log:
  level: info
  file: ./log/httpcat.log
```

## 📁 目录说明

安装后会创建以下目录：

| 目录 | 路径 | 说明 |
|------|------|------|
| 可执行文件 | `/usr/local/bin/httpcat` | 主程序 |
| 配置文件 | `/etc/httpcat/svr.yml` | 主配置 |
| 数据目录 | `/var/lib/httpcat/` | 数据根目录 |
| 静态资源 | `/var/lib/httpcat/static/` | Web 界面 |
| 上传目录 | `/var/lib/httpcat/upload/` | 上传文件存储 |
| 下载目录 | `/var/lib/httpcat/download/` | 下载缓存 |
| 数据库 | `/var/lib/httpcat/data/` | SQLite 数据库 |

## 🔧 命令行参数

```bash
./httpcat -h

选项:
  --port, -p         服务端口 (默认: 8888)
  -C                 配置文件路径 (默认: ./conf/svr.yml)
  --static           静态资源目录 (默认: ./static/)
  --upload           上传目录 (默认: ./website/upload/)
  --download         下载目录 (默认: ./website/download/)
  --p2pport          P2P 监听端口
  -v                 显示版本信息
```

## 🔧 命令行参数

```bash
./httpcat -h

选项:
  --port, -p     服务端口 (默认: 8888)
  -C             配置文件路径
  --static       静态资源目录
  -v             显示版本信息
```

## ❓ 常见问题

### Q: 登录提示"账号或密码错误"？

首先确认默认账号密码：`admin` / `admin`

如果仍然无法登录，可能是使用了**不支持 SQLite 的版本**。检查方法：

```bash
# 查看启动日志
./httpcat --port=8888 -C conf/svr.yml 2>&1 | grep -i "sqlite\|CGO"
```

如果看到 SQLite 相关错误，说明该版本无数据库支持。请使用以下版本：
- **Linux**：在 Linux 服务器上编译的版本
- **macOS**：在 macOS 上使用 Docker 编译的 Linux 版本

### Q: 端口被占用？

```bash
# 查看端口占用
lsof -i :8888  # Linux/macOS
netstat -ano | findstr 8888  # Windows

# 使用其他端口
./httpcat --port=9999 -C conf/svr.yml
```

### Q: 如何修改密码？

登录管理界面后，点击右上角用户头像 → 个人设置 → 修改密码

### Q: 如何卸载？

```bash
# 使用卸载脚本（推荐）
sudo ./uninstall.sh

# 或手动删除
sudo systemctl stop httpcat
sudo systemctl disable httpcat
sudo rm /usr/local/bin/httpcat
sudo rm -f /etc/systemd/system/httpcat.service

# 删除数据（谨慎操作）
sudo rm -rf /etc/httpcat      # 配置文件
sudo rm -rf /var/lib/httpcat  # 数据文件
```

## 📞 获取帮助

- GitHub: https://github.com/xxxx/httpcat
- Issues: https://github.com/xxxx/httpcat/issues
