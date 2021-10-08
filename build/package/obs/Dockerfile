FROM golang1.16.2:latest as builder
ARG URL=https://crochee:lcf13149825.lmy@github.com/crochee/cca.git
WORKDIR /opt/cloud
RUN git clone ${URL} &&\
    cd cca &&\
    go mod tidy &&\
    go build -ldflags="-s -w" -o cca -tags jsoniter cmd/server &&\
    rm -rf ./docs/docs.go

FROM centos:latest
WORKDIR /opt/cloud/cca
ARG cfg=./conf/config.yml
ENV config=${cfg}
RUN mkdir -p /opt/cloud/cca/docs &&\
    mkdir -p /opt/cloud/cca/conf
COPY --from=builder /opt/cloud/cca/cca .
COPY --from=builder /opt/cloud/cca/docs ./docs
COPY --from=builder /opt/cloud/cca/conf ./conf
COPY --from=builder /opt/cloud/cca/.build/cca/entrypoint.sh .
RUN chmod +x ./cca && chmod +x ./entrypoint.sh
VOLUME /cca

EXPOSE 8150
STOPSIGNAL 2

ENTRYPOINT ["./entrypoint.sh"]
CMD ["./cca"]
