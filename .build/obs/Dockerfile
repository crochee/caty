FROM centos:latest

WORKDIR /opt/cloud/obs
RUN mkdir -p ./docs &&\
    mkdir -p ./conf
COPY ./obs .
COPY ./conf ./conf
COPY ./docs ./docs
RUN chmod +x ./obs

EXPOSE 8150
ENTRYPOINT ["entrypoint.sh"]
CMD ["obs"]