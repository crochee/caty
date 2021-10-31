FROM centos-go:latest as builder
ARG Project=caty
ARG User=crochee
ARG Password=lcf13149825.???
ARG URL=https://${User}:${Password}@gitee.com/crochee/${Project}.git
WORKDIR /opt/cloud
RUN git clone ${URL} &&\
    cd ${Project} &&\
    go mod tidy &&\
    GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-s -w" -o caty -tags jsoniter ./cmd/${Project}

FROM centos:latest
ARG Project=/opt/cloud/caty
WORKDIR ${Project}
RUN mkdir -p ${Project}/conf
COPY --from=builder ${Project}/caty .
COPY --from=builder ${Project}/conf ./conf
COPY --from=builder ${Project}/build/caty/entrypoint.sh .
RUN chmod +x ./caty && chmod +x ./entrypoint.sh

EXPOSE 8120
STOPSIGNAL 2

ENTRYPOINT ["./entrypoint.sh"]
CMD ["./caty"]
