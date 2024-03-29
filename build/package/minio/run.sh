#!/bin/bash
set -e

docker run -p 9000:9000 -p 9001:9001 --name minio -d --restart=always -e "MINIO_ACCESS_KEY=root" -e "MINIO_SECRET_KEY=123456" -v /data/minio/data:/data -v /data/minio/config:/root/.minio minio/minio server /data --console-address ":9001"

docker run -p 9000:9000 -p 9001:9001 --name minio -d --restart=always -v /data/:/data -v /data/config:/root/.minio -e "MINIO_ACCESS_KEY=root" -e "MINIO_SECRET_KEY=12345678"  minio/minio server /data/minioone /data/miniotwo /data/miniothree /data/miniofour --console-address ":9001"