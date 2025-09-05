# 多阶段构建 Dockerfile
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装必要的工具
RUN apk add --no-cache git

# 复制 go mod 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建所有服务
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o user-service ./services/user/user/user.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o course-service ./services/course/course/course.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ai-service ./services/ai_assistant/ai_assistant/ai_assistant.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o exam-service ./services/exam/exam/exam.go

# 最终阶段
FROM alpine:latest

# 安装 ca-certificates 和基础工具
RUN apk --no-cache add ca-certificates tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    apk del tzdata

# 创建非 root 用户
RUN addgroup -g 1000 appgroup && \
    adduser -D -u 1000 -G appgroup appuser

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/user-service .
COPY --from=builder /app/course-service .
COPY --from=builder /app/ai-service .
COPY --from=builder /app/exam-service .

# 复制配置文件
COPY services/user/user/etc/user-api.yaml ./etc/
COPY services/course/course/etc/course-api.yaml ./etc/
COPY services/ai_assistant/ai_assistant/etc/ai_assistant-api.yaml ./etc/
COPY services/exam/exam/etc/exam-api.yaml ./etc/

# 创建必要的目录
RUN mkdir -p logs && \
    chown -R appuser:appgroup /app

# 切换到非 root 用户
USER appuser

# 暴露端口
EXPOSE 8881 8882 8883 8884

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8881/api/v1/user/info || exit 1

# 启动脚本
CMD ["sh", "-c", "./user-service -f etc/user-api.yaml & ./course-service -f etc/course-api.yaml & ./ai-service -f etc/ai_assistant-api.yaml & ./exam-service -f etc/exam-api.yaml & wait"]