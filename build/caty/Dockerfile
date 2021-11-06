FROM golang:1.17.3-alpine as builder
ARG Project=caty
ARG URL=https://gitee.com/crochee/${Project}.git
ENV GOSU_VERSION=1.14
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
    GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o ${Project} -tags jsoniter ./cmd/${Project} &&\
    go install github.com/tianon/gosu@${GOSU_VERSION} &&\
    cp ${GOPATH}/bin/gosu ./ &&\
    gosu --version &&\
    gosu nobody true

FROM alpine:latest as runner
ARG Project=/opt/cloud/caty
WORKDIR ${Project}
RUN mkdir -p ${Project}/conf ${Project}/log
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder ${Project}/gosu /usr/local/bin/
COPY --from=builder ${Project}/caty /usr/local/bin/
COPY --from=builder ${Project}/build/caty/entrypoint.sh /usr/local/bin/
COPY --from=builder ${Project}/conf ./conf

# 赋予执行权限
RUN chmod +x /usr/local/bin/caty /usr/local/bin/entrypoint.sh /usr/local/bin/gosu
# add our user and group first to make sure their IDs get assigned consistently, regardless of whatever dependencies get added
RUN groupadd -g 1000 -r cloud && useradd -r -g cloud dev -u 1000
# 将工作目录加入用户组
RUN chgrp -R  cloud .
# 配置文件目录和文件0440,只有读权限
RUN chmod -R ug=r,o-rwx ${Project}/conf
# 日志文件夹0744,其他为0644
RUN chmod -R u=rw,g=r,o=r ${Project}/log &&\
    chmod u+x ${Project}/log

EXPOSE 8120
STOPSIGNAL 2

ENTRYPOINT ["entrypoint.sh"]
