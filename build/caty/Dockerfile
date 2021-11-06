FROM golang:1.17.3-alpine as builder
ARG Project=caty
ARG URL=https://gitee.com/crochee/${Project}.git
WORKDIR /opt/cloud
# 下载git 修改配置
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories &&\
    apk add --no-cache ca-certificates git tzdata
# 设置代理环境变量
RUN go env -w GOPROXY=https://goproxy.io,https://goproxy.cn,direct &&\
    go env -w GO111MODULE=on
# 拉去代码
RUN git clone ${URL}
# 文件dos->unix
RUN sed -e 's/.$//' ./${Project}/build/caty/entrypoint.sh > ./entrypoint.sh &&\
    rm -rf ./${Project}/build/caty/entrypoint.sh &&\
    mv ./entrypoint.sh ./${Project}/build/caty/entrypoint.sh
# 代码编译
RUN cd ${Project} &&\
    go mod tidy &&\
    GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o ${Project} -tags jsoniter ./cmd/${Project}

FROM alpine:latest as runner
ARG Project=/opt/cloud/caty
WORKDIR ${Project}
RUN mkdir -p ${Project}/conf
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder ${Project}/caty /usr/local/bin/
COPY --from=builder ${Project}/conf ./conf
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder ${Project}/build/caty/entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/caty && chmod +x /usr/local/bin/entrypoint.sh

EXPOSE 8120
STOPSIGNAL 2

ENTRYPOINT ["entrypoint.sh"]
CMD ["caty"]
