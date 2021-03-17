FROM golang1.16.2:latest as builder
ARG URL="https://crochee:lcf13149825.lmy@github.com/crochee/obs.git"
ARG COMPONENT="obs"
WORKDIR /opt/cloud
RUN git clone ${URL} &&\
    cd ${COMPONENT} &&\
    go mod tidy &&\
    go build -o ${COMPONENT} -tags jsoniter .

FROM centos:latest
WORKDIR /opt/cloud/obs
RUN mkdir -p /opt/cloud/obs/docs &&\
    mkdir -p /opt/cloud/obs/conf
COPY --from=builder /opt/cloud/obs .
COPY --from=builder /opt/cloud/obs/docs ./docs
COPY --from=builder /opt/cloud/obs/conf ./conf
COPY --from=builder /opt/cloud/obs/.build/obs/entrypoint.sh .
RUN rm -rf ./docs/docs.go && chmod +x ./obs && chmod +x ./entrypoint.sh

EXPOSE 8150
ENTRYPOINT ["./entrypoint.sh"]
CMD ["./obs"]
