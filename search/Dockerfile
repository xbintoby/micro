# 基础镜像
FROM golang:latest

# 容器环境变量设置，会覆盖默认的变量值
ENV GOPROXY=https://goproxy.cn,direct

ENV GO111MODULE="on"

# 作者
MAINTAINER jam3 xbinwar3@gmail.com

RUN mkdir -p /go/src/config

# 工作区
WORKDIR /go/src/

# 复制仓库源文件到容器里
COPY ./config/* /go/src/config/
COPY ./static/css/* /go/src/static/css/
COPY ./static/js/* /go/src/static/js/
COPY ./templates/* /go/src/templates/
COPY ./main /go/src

# 编译可执行二进制文件

# 容器向外提供服务的暴露端口
#EXPOSE  81
#EXPOSE  8881

# 启动服务
ENTRYPOINT ["./main"]
