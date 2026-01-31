# HttpCat Dockerfile
# 多阶段构建：前端构建 + Go 编译 + 运行

# ============ 第一阶段：前端构建（可选） ============
FROM node:20-alpine AS frontend-builder

WORKDIR /app/web

# 复制前端文件
COPY web/package*.json ./

# 安装前端依赖
RUN npm ci --registry=https://registry.npmmirror.com

# 复制前端源码
COPY web/ ./

# 构建前端
RUN npm run build

# ============ 第二阶段：Go 编译 ============
FROM golang:1.23-alpine AS go-builder

# 安装编译依赖（SQLite 需要 CGO，linux-headers 解决 off64_t 问题）
RUN apk add --no-cache gcc musl-dev sqlite-dev linux-headers

# 设置 Go 代理（加速中国大陆下载）
ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /app

# 复制 Go 依赖文件
COPY server-go/go.mod server-go/go.sum ./

RUN go mod download

# 复制 Go 源代码
COPY server-go/ ./

# 编译（启用 CGO 支持 SQLite，CGO_CFLAGS 修复 musl 兼容性）
ARG VERSION=v0.2.0
ARG BUILD_TIME
ARG COMMIT_ID
RUN CGO_ENABLED=1 GOOS=linux CGO_CFLAGS="-D_LARGEFILE64_SOURCE" go build \
    -ldflags "-s -w -X httpcat/internal/common.Version=${VERSION} -X httpcat/internal/common.Build=${BUILD_TIME} -X httpcat/internal/common.Commit=${COMMIT_ID}" \
    -o httpcat ./cmd/httpcat.go

# ============ 第三阶段：运行 ============
FROM alpine:3.19

# 安装运行时依赖
RUN apk add --no-cache ca-certificates tzdata sqlite-libs

# 设置时区
ENV TZ=Asia/Shanghai

WORKDIR /app

# 从 Go 编译阶段复制二进制文件
COPY --from=go-builder /app/httpcat /app/httpcat

# 复制配置文件
COPY server-go/internal/conf /app/conf

# 从前端构建阶段复制静态资源（如果有）
COPY --from=frontend-builder /app/web/dist /app/website/static

# 如果没有前端构建，使用默认的 static 目录
# COPY static /app/website/static

# 创建数据目录
RUN mkdir -p /app/data /app/website/upload /app/website/download /app/log

# 暴露端口
EXPOSE 8888

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget -q --spider http://localhost:8888/api/v1/conf/getVersion || exit 1

# 启动命令
ENTRYPOINT ["/app/httpcat"]
CMD ["--port=8888", "-C", "/app/conf/svr.yml"]
