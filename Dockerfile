# 机器人路径编辑器 Dockerfile
# 多阶段构建，优化镜像大小

# 阶段1: 构建Go应用
FROM golang:1.21-alpine AS builder

# 安装构建依赖
RUN apk add --no-cache git gcc musl-dev

WORKDIR /app

# 复制依赖文件并下载
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码并构建
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -o robot-path-editor cmd/server/main.go

# 阶段2: 运行时镜像
FROM alpine:latest

# 安装运行时依赖
RUN apk --no-cache add ca-certificates tzdata wget && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    adduser -D -s /bin/sh robot

WORKDIR /app

# 创建必要目录
RUN mkdir -p data logs && chown -R robot:robot /app

# 复制应用文件
COPY --from=builder /app/robot-path-editor .
COPY --chown=robot:robot configs ./configs
COPY --chown=robot:robot web ./web

# 切换到非root用户
USER robot

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# 启动应用
CMD ["./robot-path-editor"]