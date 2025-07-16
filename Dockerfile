# 构建阶段
FROM golang:1.24.4-alpine AS builder

WORKDIR /app

# 复制依赖文件并下载
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码并构建
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server cmd/server/main.go

# 运行阶段
FROM alpine:latest

# 安装必要的运行时依赖
RUN apk --no-cache add ca-certificates netcat-openbsd

WORKDIR /app

# 复制二进制文件
COPY --from=builder /app/server .

# 复制迁移文件
COPY --from=builder /app/migrations ./migrations

# 复制健康检查脚本
COPY --from=builder /app/scripts ./scripts
RUN chmod +x ./scripts/health_check.sh

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD ./scripts/health_check.sh

# 启动应用程序
CMD ["./server"]
