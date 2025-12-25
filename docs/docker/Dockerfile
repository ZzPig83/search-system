# ---------- 构建阶段 ----------
FROM golang:1.25.2-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装依赖
RUN apk add --no-cache git ca-certificates

# 复制go模块文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源码
COPY . .

# 编译
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o search-system .

# ---------- 运行阶段 ----------
FROM alpine:latest

# 设置工作目录
WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/search-system .

# 复制数据文件（如果需要）
COPY data_source/ ./data_source/

# 暴露端口
EXPOSE 8080

# 启动应用
CMD ["./search-system"]
