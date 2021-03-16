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
COPY --from=builder /opt/cloud/obs .
COPY --from=builder /opt/cloud/obs/docs .
COPY --from=builder /opt/cloud/obs/conf .
COPY --from=builder /opt/cloud/obs/.build/obs/entrypoint.sh .
RUN rm -rf ./docs/docs.go&&chmod +x ./obs

EXPOSE 8150
ENTRYPOINT ["entrypoint.sh"]
CMD ["obs"]
