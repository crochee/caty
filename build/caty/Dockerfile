FROM golang:1.17.3-alpine as builder
ARG Project=caty
ARG URL=https://gitee.com/crochee/${Project}.git
ARG GOSU_VERSION=1.14
WORKDIR /opt/cloud
# 下载git 修改配置
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories &&\
    apk add --no-cache git tzdata dos2unix
# 设置代理环境变量
RUN go env -w GOPROXY=https://goproxy.io,https://goproxy.cn,direct &&\
    go env -w GO111MODULE=on
# 拉去代码
RUN git clone ${URL}
# 文件dos->unix
RUN dos2unix  ./${Project}/build/caty/entrypoint.sh
# 代码编译
RUN cd ${Project} &&\
    go mod tidy &&\
    GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o caty -tags jsoniter ./cmd/${Project} &&\
    go install github.com/tianon/gosu@${GOSU_VERSION}
# 整理项目需要拷贝的资源
RUN mv ./${Project}/build/caty/entrypoint.sh . &&\
    mv ./${Project}/conf . &&\
    mkdir ./out &&\
    cp ./${Project}/caty ./out/ &&\
    cp ${GOPATH}/bin/gosu .

FROM alpine:latest as runner
ARG WorkDir=/opt/cloud/caty
WORKDIR ${WorkDir}
# add our user and group first to make sure their IDs get assigned consistently, regardless of whatever dependencies get added
RUN addgroup -g 10000 cloud && adduser -g cloud dev -u 5000 -D -H
# 预创建文件夹
RUN mkdir -p ${WorkDir}/conf ${WorkDir}/log
# 将工作目录加入用户组
RUN chown -R  dev:cloud ${WorkDir}
# 日志文件夹0744
RUN chmod u=rwx,g=r,o=r ${WorkDir}/log

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /opt/cloud/gosu /usr/local/bin/
COPY --from=builder /opt/cloud/out/caty /usr/local/bin/
COPY --from=builder /opt/cloud/entrypoint.sh /usr/local/bin/
COPY --from=builder --chown=dev:cloud /opt/cloud/conf ${WorkDir}/conf
# 赋予执行权限
RUN chmod +x /usr/local/bin/caty /usr/local/bin/entrypoint.sh /usr/local/bin/gosu
# 配置文件目录和文件0440,只有读权限
RUN chmod -R ug=r,o-rwx ${WorkDir}/conf

EXPOSE 8120
STOPSIGNAL 2

ENTRYPOINT ["entrypoint.sh"]
CMD ["caty"]