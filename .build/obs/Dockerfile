FROM golang1.16.2:latest as builder
ARG URL="https://crochee:lcf13149825.lmy@github.com/crochee/obs.git"
WORKDIR /opt/cloud
RUN git clone ${URL} &&\
    cd obs &&\
    go mod tidy &&\
    go build -o obs -tags jsoniter .

FROM centos:latest
WORKDIR /opt/cloud/obs
ARG cfg=./conf/config.yml
ENV config=${cfg}
RUN mkdir -p /opt/cloud/obs/docs &&\
    mkdir -p /opt/cloud/obs/conf
COPY --from=builder /opt/cloud/obs/obs .
COPY --from=builder /opt/cloud/obs/docs ./docs
COPY --from=builder /opt/cloud/obs/conf ./conf
COPY --from=builder /opt/cloud/obs/.build/obs/entrypoint.sh .
RUN rm -rf ./docs/docs.go && chmod +x ./obs && chmod +x ./entrypoint.sh

EXPOSE 8150
ENTRYPOINT ["./entrypoint.sh"]
CMD ["./obs"]
