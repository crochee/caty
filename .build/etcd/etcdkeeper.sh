#!/bin/bash
set -e

#docker pull evildecay/etcdkeeper

docker run -itd -p 8080:8080 --restart=always --name etcdkeeper docker.io/evildecay/etcdkeeper:latest