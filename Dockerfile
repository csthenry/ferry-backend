FROM golang:1.18 AS build

WORKDIR /opt/workflow/ferry
COPY . .
ARG GOPROXY="https://goproxy.cn"
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ferry .

FROM alpine AS prod

MAINTAINER lanyulei

RUN echo -e "http://mirrors.aliyun.com/alpine/v3.11/main\nhttp://mirrors.aliyun.com/alpine/v3.11/community" > /etc/apk/repositories \
    && apk add -U tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime 

WORKDIR /opt/workflow/ferry

RUN apk add python3
RUN apk add py3-pip
RUN pip3 install requests

#获取golang:1.18中编译的go二进制文件
COPY --from=build /opt/workflow/ferry/ferry /opt/workflow/ferry/

COPY config/ /opt/workflow/ferry/default_config/
COPY template/ /opt/workflow/ferry/template/
COPY static/template/ /opt/workflow/ferry/static/template/
COPY docker/entrypoint.sh /opt/workflow/ferry/
RUN mkdir -p logs static/uploadfile static/scripts static/template

RUN chmod 755 /opt/workflow/ferry/entrypoint.sh
RUN chmod 755 /opt/workflow/ferry/ferry

EXPOSE 8001
VOLUME [ "/opt/workflow/ferry/config" ]
ENTRYPOINT [ "/opt/workflow/ferry/entrypoint.sh" ]
