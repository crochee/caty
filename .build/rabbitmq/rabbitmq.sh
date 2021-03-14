#!/bin/bash
set -e

docker run -d --hostname my-rabbit --name rabbit-test -p 8080:15672 rabbitmq

docker run -d --hostname my-rabbit --name rabbit -e RABBITMQ_DEFAULT_USER=admin -e RABBITMQ_DEFAULT_PASS=admin -p 15672:15672 -p 5672:5672 -p 25672:25672 -p 61613:61613 -p 1883:1883 rabbitmq:management

#15672：控制台端口号
#5672：应用访问端口号

#   http://localhost:15672
#默认账户名：guest
#密码：guest

#  http://127.0.0.1:15672/#/queues