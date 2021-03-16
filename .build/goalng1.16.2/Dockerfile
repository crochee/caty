FROM centos:latest
ENV GO_VERSION "1.16.2"
#下载go
RUN wget https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz &&\
    tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz &&\
    rm -rf go${GO_VERSION}.linux-amd64.tar.gz
#环境变量设置
ENV GOROOT=/usr/local/go
ENV GOPATH=/root/go
ENV PATH=${GOROOT}/bin/:${GOPATH}/bin/:${PATH}
#代理设置
RUN go env -w GOPROXY=https://goproxy.io,direct
RUN go env -w GO111MODULE=on

RUN mkdir -p ${GOPATH}/bin