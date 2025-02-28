FROM --platform=$TARGETPLATFORM golang:1.23.6 AS builder

# 设置工作目录
WORKDIR /build

# 复制源代码
COPY . .

# 基于目标平台设置 GOARCH
ARG TARGETPLATFORM
RUN if [ "$TARGETPLATFORM" = "linux/amd64" ]; then \
      GOARCH=amd64; \
    elif [ "$TARGETPLATFORM" = "linux/arm64" ]; then \
      GOARCH=arm64; \
    else \
      echo "Unsupported platform: $TARGETPLATFORM"; \
      exit 1; \
    fi; \
    CGO_ENABLED=0 GOOS=linux GOARCH=$GOARCH go build -o mcp-server-gitlab

# 使用多阶段构建，最终镜像
FROM --platform=$TARGETPLATFORM debian:latest

# 设置工作目录
WORKDIR /app

# 安装 ca-certificates 包
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# 从 builder 阶段复制编译好的二进制文件
COPY --from=builder /build/mcp-server-gitlab /app/
COPY .env /app/

# 确保二进制文件有执行权限
RUN chmod +x /app/mcp-server-gitlab

# 设置容器启动命令
CMD ["/app/mcp-server-gitlab", ".env"]