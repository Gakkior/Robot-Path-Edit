# 多阶段构建的Dockerfile
# 阶段1: 构建Go应用
FROM golang:1.21-alpine AS builder

# 安装必要的工具
RUN apk add --no-cache git gcc musl-dev

# 设置工作目录
WORKDIR /app

# 复制go mod文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用 (CGO_ENABLED=1 支持SQLite)
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o robot-path-editor cmd/server/main.go

# 阶段2: 运行时镜像
FROM alpine:latest

# 安装运行时依赖
RUN apk --no-cache add ca-certificates tzdata

# 设置时区
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

# 创建应用用户
RUN adduser -D -s /bin/sh robot

# 设置工作目录
WORKDIR /app

# 创建必要的目录
RUN mkdir -p data logs && \
    chown -R robot:robot /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/robot-path-editor .

# 复制配置文件和静态资源
COPY --chown=robot:robot configs ./configs
COPY --chown=robot:robot web ./web

# 切换到应用用户
USER robot

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# 启动命令
CMD ["./robot-path-editor"]

# 构建命令:
# docker build -t robot-path-editor .
# 
# 运行命令:
# docker run -d \
#   --name robot-path-editor \
#   -p 8080:8080 \
#   -v $(pwd)/data:/app/data \
#   -v $(pwd)/logs:/app/logs \
#   robot-path-editor
#
# 开发模式运行:
# docker run -it --rm \
#   --name robot-path-editor-dev \
#   -p 8080:8080 \
#   -v $(pwd):/app \
#   golang:1.21-alpine \
#   sh -c "cd /app && go run cmd/server/main.go"