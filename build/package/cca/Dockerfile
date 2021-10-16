FROM centos-go:latest as builder
ARG Project=cca
ARG URL=https://gitee.com/crochee/${Project}.git
WORKDIR /opt/cloud
RUN git clone ${URL} &&\
    cd ${Project} &&\
    go mod tidy &&\
    go build -ldflags="-s -w" -o cca -tags jsoniter cmd/${Project}

FROM centos:latest
ARG Project=/opt/cloud/cca
WORKDIR ${Project}
RUN mkdir -p ${Project}/conf
COPY --from=builder /opt/cloud/cca/cca .
COPY --from=builder /opt/cloud/cca/conf ./conf
COPY --from=builder /opt/cloud/cca/build/package/cca/entrypoint.sh .
RUN chmod +x ./cca && chmod +x ./entrypoint.sh

EXPOSE 8120
STOPSIGNAL 2

ENTRYPOINT ["./entrypoint.sh"]
CMD ["./cca"]
