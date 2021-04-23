#!/bin/bash
set -e

docker run -itd --name redis-test -p 6379:6379 redis
docker run -p 6379:6379 --name redis -d redis:latest --requirepass "123456"

docker run --name my_redis -p 6379:6379 -v /root/docker/redis/data:/data -v /root/docker/redis/conf/redis.conf:/etc/redis/redis.conf -d redis redis-server /etc/redis/redis.conf

#auth 你的密码