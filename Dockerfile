# 构建可执行文件
FROM golang:1.23-alpine AS builder

# 设置工作目录
WORKDIR /app

# 拷贝代码
COPY . .

# 下载依赖
ENV http_proxy=""
ENV https_proxy=""
ENV HTTP_PROXY=""
ENV HTTPS_PROXY=""
ENV GOPROXY=https://goproxy.io,direct
RUN go mod download

# 构建可执行文件，区分环境
RUN go build -o ./bin/godemo ./cmd/main.go

# 用 ENTRYPOINT 保证参数可被 docker run/docker-compose 传递
ENTRYPOINT ["./bin/godemo"]

# 默认参数，可被 docker run/docker-compose 覆盖
CMD ["-c", "./config/debug.toml"]