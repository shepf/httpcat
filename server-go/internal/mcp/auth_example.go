package mcp

/*
MCP 认证配置指南

HttpCat MCP Server 支持可选的 Bearer Token 认证。

=== 配置方式 ===

1. 编辑 conf/svr.yml：

server:
  mcp:
    enable: true
    auth_token: "your-secure-mcp-token-2024"  # 设置认证 Token

2. 客户端配置示例：

Claude Desktop (~/Library/Application Support/Claude/claude_desktop_config.json):
{
  "mcpServers": {
    "httpcat": {
      "type": "sse",
      "url": "http://localhost:8888/mcp/sse",
      "headers": {
        "Authorization": "Bearer your-secure-mcp-token-2024"
      }
    }
  }
}

Cursor:
- Name: httpcat
- Type: SSE
- URL: http://localhost:8888/mcp/sse
- Headers: Authorization: Bearer your-secure-mcp-token-2024

=== 安全建议 ===

1. 生产环境强烈建议启用认证
2. Token 应该是随机生成的强密码
3. 定期轮换 Token
4. 不要在代码仓库中提交真实 Token

=== Token 生成示例 ===

# macOS/Linux
openssl rand -hex 32

# 输出示例: a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0u1v2w3x4y5z6

=== 不启用认证 ===

如果 auth_token 为空或未设置，MCP Server 将允许任何连接。
仅建议在开发环境或内网使用。

server:
  mcp:
    enable: true
    auth_token: ""  # 留空则不验证
*/
