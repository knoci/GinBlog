FROM golang:alpine AS builder

# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# 移动到工作目录：/build
WORKDIR /build

# 将代码复制到容器中
COPY . .

# 下载依赖信息
RUN go mod download

# 将我们的代码编译成二进制可执行文件 bubble
RUN go build -o ginblog .

###################
# 接下来创建一个小镜像
###################
FROM debian:stretch-slim

# 从builder镜像中把脚本拷贝到当前目录
COPY ./wait-for.sh /

# 从builder镜像中把静态文件拷贝到当前目录
COPY ./templates /templates
COPY ./static /static

# 从builder镜像中把配置文件拷贝到当前目录
COPY ./config /config

# 从builder镜像中把/dist/app 拷贝到当前目录
COPY --from=builder /build/ginblog /

EXPOSE 8808

RUN echo "" > /etc/apt/sources.list; \
    echo "deb http://mirrors.aliyun.com/debian buster main" >> /etc/apt/sources.list ; \
    echo "deb http://mirrors.aliyun.com/debian-security buster/updates main" >> /etc/apt/sources.list ; \
    echo "deb http://mirrors.aliyun.com/debian buster-updates main" >> /etc/apt/sources.list ; \
    set -eux; \
	apt-get update; \
	apt-get install -y \
		--no-install-recommends \
		netcat; \
    chmod 755 wait-for.sh

## 需要运行的命令
#ENTRYPOINT ["/ginblog", "config/config.yaml"]
