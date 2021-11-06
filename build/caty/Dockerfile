FROM golang:1.17.1-alpine as builder
ARG Project=caty
ARG URL=https://gitee.com/crochee/${Project}.git
WORKDIR /opt/cloud
RUN git clone ${URL} &&\
    cd ${Project} &&\
    go mod tidy &&\
    GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o caty -tags jsoniter ./cmd/${Project}

FROM busybox:latest as runner
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
