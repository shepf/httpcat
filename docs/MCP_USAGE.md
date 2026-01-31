# HttpCat MCP 使用指南

> 🎯 **本文档面向小白用户**，帮助你 5 分钟内完成 MCP 配置，让 AI 助手帮你管理文件！

---

## 📖 目录

- [什么是 MCP？](#什么是-mcp)
- [🚀 快速入门（5分钟）](#-快速入门5分钟)
- [🔧 详细配置教程](#-详细配置教程)
- [💬 使用示例](#-使用示例)
- [🔐 上传文件教程](#-上传文件教程)
- [❓ 常见问题](#-常见问题)
- [📚 高级配置](#-高级配置)

---

## 什么是 MCP？

**MCP (Model Context Protocol)** 是 Anthropic 推出的开放协议，让 AI 助手可以直接操作外部工具。

简单理解：**配置 MCP 后，你可以用自然语言让 AI 帮你管理 HttpCat 上的文件**！

| 没有 MCP | 有了 MCP |
|---------|---------|
| 手动登录网页 → 找到文件 → 点击删除 | 对 AI 说"删除 test.pdf" |
| 手动打开浏览器 → 查看文件列表 | 对 AI 说"看看有什么文件" |
| 手动下载 → 检查 MD5 值 | 对 AI 说"验证文件完整性" |

---

## 🚀 快速入门（5分钟）

### 第一步：确认服务已运行

打开浏览器访问你的 HttpCat 地址，能看到页面说明服务正常：

```
http://你的服务器IP:8888
```

### 第二步：测试 MCP 端点

在终端运行（替换成你的地址）：

```bash
curl http://你的服务器IP:8888/mcp/sse
```

看到类似输出说明 MCP 正常工作：
```
event: endpoint
data: /mcp/message?sessionId=xxxx-xxxx-xxxx
```

### 第三步：配置 AI 工具

根据你使用的 AI 工具，选择下面对应的配置方法。

---

## 🔧 详细配置教程

### 方法 A：Claude Desktop 配置

> 适用于：macOS / Windows 上使用 Claude Desktop 客户端

**macOS 用户：**

1. 打开终端，运行：
   ```bash
   open ~/Library/Application\ Support/Claude/claude_desktop_config.json
   ```

2. 在打开的文件中添加以下内容：
   ```json
   {
     "mcpServers": {
       "httpcat": {
         "type": "sse",
         "url": "http://你的服务器IP:8888/mcp/sse"
       }
     }
   }
   ```

3. 保存文件，**重启 Claude Desktop**

**Windows 用户：**

1. 按 `Win + R`，输入：
   ```
   %APPDATA%\Claude\claude_desktop_config.json
   ```

2. 用记事本打开，添加同样的配置

3. 保存文件，**重启 Claude Desktop**

**✅ 验证配置成功：**
重启后，在 Claude Desktop 中输入 `"查看 httpcat 有哪些文件"`，如果 AI 能返回文件列表，说明配置成功！

---

### 方法 B：Cursor 配置

> 适用于：使用 Cursor 编辑器的开发者

1. 打开 Cursor，点击左下角的 **设置图标** ⚙️
2. 选择 **Features** → **MCP**
3. 点击 **Add new MCP server**
4. 填写信息：

   | 字段 | 值 |
   |-----|-----|
   | Name | `httpcat` |
   | Type | `SSE` |
   | URL | `http://你的服务器IP:8888/mcp/sse` |

5. 点击保存

**✅ 验证配置成功：**
在 Cursor 的 AI 对话框中输入 `"用 httpcat 查看文件列表"`

---

### 方法 C：CodeBuddy 配置

> 适用于：使用 CodeBuddy IDE 的用户

1. 点击 CodeBuddy 左侧的 **MCP 图标**
2. 点击 **添加 MCP 服务器**
3. 选择 **SSE 类型**
4. 填入 URL：`http://你的服务器IP:8888/mcp/sse`
5. 点击确认

---

### 方法 D：其他 MCP 客户端

通用配置参数：

| 参数 | 值 |
|-----|-----|
| 类型 | SSE |
| SSE 端点 | `http://服务器IP:端口/mcp/sse` |
| 消息端点 | `http://服务器IP:端口/mcp/message` |

---

## 💬 使用示例

配置成功后，你可以用自然语言让 AI 帮你管理文件！

### 示例 1：查看文件列表

> **你**：看看 httpcat 上有什么文件  
> **AI**：当前上传目录有 5 个文件：
> - report.pdf (2.5 MB)  
> - data.csv (150 KB)  
> - image.png (3.2 MB)  
> - backup.zip (50 MB)  
> - notes.txt (10 KB)

### 示例 2：查看磁盘使用情况

> **你**：查看一下磁盘占用情况  
> **AI**：磁盘使用情况：
> - 文件数量：150 个
> - 总大小：2.3 GB
> - 存储路径：/data/upload

### 示例 3：获取文件详情

> **你**：report.pdf 这个文件多大？什么时候上传的？  
> **AI**：report.pdf 文件信息：
> - 大小：2.5 MB
> - 修改时间：2026-01-30 10:30:00
> - MD5：d41d8cd98f00b204e9800998ecf8427e

### 示例 4：验证文件完整性

> **你**：帮我验证 data.zip 的 MD5 是不是 abc123...  
> **AI**：验证结果：
> - 期望 MD5：abc123...
> - 实际 MD5：abc123...
> - ✅ 匹配成功，文件完整无损！

### 示例 5：删除文件（两步确认）

> **你**：删除 old_backup.zip  
> **AI**：正在请求删除 old_backup.zip (50 MB)...  
> **AI**：确认删除中...  
> **AI**：✅ 文件 old_backup.zip 已成功删除！

### 示例 6：查看上传历史

> **你**：今天上传了哪些文件？  
> **AI**：今天的上传记录：
> 1. document.pdf - 14:30:00 - 2.5 MB
> 2. image.png - 13:15:00 - 3.2 MB
> 3. data.xlsx - 11:00:00 - 850 KB

---

## 🔐 上传文件教程

> ⚠️ **注意**：上传文件需要 Token，这是为了安全考虑！

### 为什么需要 Token？

MCP 协议本身没有登录功能，所以上传文件需要一个"通行证"（Token）来证明你有权限。

### 获取 Token 的 3 种方法

#### 方法 1：Web 界面生成（推荐新手）

1. 浏览器打开 `http://你的服务器IP:8888`
2. 使用管理员账号登录（默认 admin/admin）
3. 点击右上角头像 → **个人设置**
4. 找到 **上传 Token 管理** → 点击 **生成 Token**
5. 复制生成的 Token

#### 方法 2：命令行生成

```bash
# 1. 先登录获取 JWT Token
curl -X POST http://你的服务器IP:8888/api/v1/user/login/account \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin"}'

# 返回：{"data":{"token":"eyJhbGciOiJIUzI1NiIs..."}}

# 2. 用 JWT Token 生成上传 Token
curl -X POST http://你的服务器IP:8888/api/v1/user/createUploadToken \
  -H "Authorization: Bearer 上一步返回的token" \
  -H "Content-Type: application/json" \
  -d '{"appkey":"httpcat","appsecret":"httpcat_app_secret"}'

# 返回：{"data":{"uploadToken":"httpcat:xxxx:yyyy"}}
```

#### 方法 3：使用脚本一键生成

创建脚本 `gen_token.sh`：

```bash
#!/bin/bash
HOST="${1:-http://localhost:8888}"
USER="${2:-admin}"
PASS="${3:-admin}"

JWT=$(curl -s -X POST "$HOST/api/v1/user/login/account" \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"$USER\",\"password\":\"$PASS\"}" | grep -o '"token":"[^"]*' | cut -d'"' -f4)

curl -s -X POST "$HOST/api/v1/user/createUploadToken" \
  -H "Authorization: Bearer $JWT" \
  -H "Content-Type: application/json" \
  -d '{"appkey":"httpcat","appsecret":"httpcat_app_secret"}' | grep -o '"uploadToken":"[^"]*' | cut -d'"' -f4
```

使用方法：
```bash
chmod +x gen_token.sh
./gen_token.sh http://你的服务器IP:8888 admin 你的密码
```

### 使用 Token 上传文件

有了 Token 后，告诉 AI：

> **你**：帮我上传 report.pdf，Token 是 `httpcat:abc123:xyz789`  
> **AI**：好的，正在上传...  
> **AI**：✅ 上传成功！
> - 文件名：report.pdf
> - 大小：2.5 MB
> - 下载链接：http://xxx/api/v1/file/download?filename=report.pdf

---

## ❓ 常见问题

### Q1：AI 说连接不上 MCP Server

**可能原因：**
1. HttpCat 服务没有启动
2. URL 配置错误
3. 防火墙阻止了连接

**解决步骤：**
```bash
# 1. 检查服务是否运行
curl http://你的服务器IP:8888/api/v1/conf/getVersion

# 2. 检查 MCP 端点
curl http://你的服务器IP:8888/mcp/sse

# 3. 检查防火墙（Linux）
sudo ufw allow 8888
```

### Q2：查询上传历史显示"SQLite is not enabled"

**原因**：需要启用 SQLite 才能记录上传历史

**解决**：修改配置文件 `conf/svr.yml`：
```yaml
server:
  http:
    file:
      enable_sqlite: true
      sqlite_db_path: "./data/httpcat_sqlite.db"
```

### Q3：上传文件失败，提示 Token 无效

**可能原因：**
1. Token 已过期
2. Token 格式错误
3. 服务端配置不正确

**解决**：重新生成一个 Token

### Q4：删除文件时报权限错误

**可能原因：**
1. 文件被其他进程占用
2. HttpCat 进程没有写权限

**解决**：
```bash
# 检查文件权限
ls -la /你的上传目录/

# 修改权限（如需要）
chmod -R 755 /你的上传目录/
```

### Q5：如何启用 MCP 认证？

如果你的 HttpCat 部署在公网，建议启用认证：

**服务端配置** (`conf/svr.yml`)：
```yaml
server:
  mcp:
    enable: true
    auth_token: "你的安全密码"  # 设置一个强密码
```

**客户端配置**（以 Claude Desktop 为例）：
```json
{
  "mcpServers": {
    "httpcat": {
      "type": "sse",
      "url": "http://你的服务器IP:8888/mcp/sse",
      "headers": {
        "Authorization": "Bearer 你的安全密码"
      }
    }
  }
}
```

---

## 🛠️ MCP 功能一览

HttpCat MCP 提供 **9 个工具** 和 **3 个资源**：

### Tools（工具）

| 工具名 | 功能 | 说明 |
|-------|------|------|
| `list_files` | 列出文件 | 查看上传目录中的文件列表 |
| `get_file_info` | 获取文件信息 | 查看文件大小、修改时间、MD5 |
| `get_upload_history` | 查询上传历史 | 查看谁在什么时候上传了什么 |
| `get_disk_usage` | 磁盘使用情况 | 查看已用空间和文件数量 |
| `request_delete_file` | 请求删除（第1步） | 获取删除确认 Token |
| `confirm_delete_file` | 确认删除（第2步） | 执行实际删除操作 |
| `upload_file` | 上传文件 | 需要上传 Token |
| `get_statistics` | 获取统计 | 上传下载统计信息 |
| `verify_file_md5` | 验证 MD5 | 检查文件完整性 |

### Resources（资源）

| 资源 URI | 说明 |
|---------|------|
| `filelist://current` | 当前文件列表 |
| `disk://usage` | 磁盘使用情况 |
| `system://info` | 系统信息 |

---

## 📚 高级配置

### 安全设计原则

1. **两步删除确认**：删除文件需要先请求，再确认，防止误删
2. **路径安全**：防止目录遍历攻击，只能操作上传目录内的文件
3. **Token 分离**：上传 Token 需要登录后才能生成，保护上传权限
4. **可选认证**：支持 Bearer Token 认证保护 MCP 端点

### 自定义 MCP 端点路径

如需修改 MCP 端点路径，编辑 `server/router.go`：

```go
// 默认路径
r.GET("/mcp/sse", mcpServer.SSEHandler())
r.POST("/mcp/message", mcpServer.MessageHandler())

// 自定义路径示例
r.GET("/api/v1/mcp/sse", mcpServer.SSEHandler())
r.POST("/api/v1/mcp/message", mcpServer.MessageHandler())
```

### 添加自定义 MCP Tool

在 `server/mcp/server.go` 中添加：

```go
s.AddTool(
    mcp.NewTool("my_custom_tool",
        mcp.WithDescription("我的自定义工具"),
        mcp.WithString("param1", mcp.Required()),
    ),
    handleMyCustomTool,
)

func handleMyCustomTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    args := request.GetArguments()
    param1 := args["param1"].(string)
    return mcp.NewToolResultText("Success: " + param1), nil
}
```

---

## 📎 相关链接

- [MCP 协议官方文档](https://modelcontextprotocol.io/)
- [HttpCat 项目主页](https://github.com/shepf/httpcat)
- [HttpCat 部署文档](./DEPLOYMENT_STATUS.md)

---

## 📝 更新日志

| 版本 | 更新内容 |
|-----|---------|
| v0.3.0 | 新增 MCP 认证、两步删除确认、上传 Token 指南 |
| v0.2.0 | 初始 MCP 支持，9 个 Tools + 3 个 Resources |

---

> 📧 **遇到问题？** 欢迎提交 Issue 或查看 [常见问题](#-常见问题) 部分！
