#!/bin/bash
set -e

image=docker.io/library/rabbitmq:latest
if [ -n "$(docker images -q ${image})" ]; then
  docker pull ${image}
fi

RABBITMQ_NAME=rabbitmq_node_0
RABBITMQ_DIR=/home/lcf/cloud/rabbitmq
rm -rf ${RABBITMQ_DIR} && mkdir -p ${RABBITMQ_DIR}
docker run -d \
--net host \
--hostname ${RABBITMQ_NAME} \
--name ${RABBITMQ_NAME} \
--log-opt max-size=10m \
--log-opt max-file=3 \
-v ${RABBITMQ_DIR}:/var/lib/rabbitmq:z \
-v ${RABBITMQ_DIR}/hosts:/etc/hosts \
-e RABBITMQ_DEFAULT_USER=admin \
-e RABBITMQ_DEFAULT_PASS='1234567' \
-e RABBITMQ_ERLANG_COOKIE='secret cookie here' \
-p 15672:15672 \
-p 5672:5672 \
--restart=always \
${image}

#开启控制台
id=$(docker ps | grep ${RABBITMQ_NAME} | awk '{print $1}')
echo "${id}"
#docker exec -it ${id} /bin/bash && rabbitmq-plugins enable rabbitmq_management && exit

#docker run -d --hostname my-rabbit --name rabbit-test -p 8080:15672 rabbitmq

#15672：控制台端口号
#5672：应用访问端口号

#   http://localhost:15672
#默认账户名：guest
#密码：guest

#  http://127.0.0.1:15672/#/queues

#4.将节点2,3加入集群
#在rabbit2机器进入容器的命令行
#sudo docker exec -it rabbit2 /bin/bash
#加入集群
#
#rabbitmqctl stop_app
#rabbitmqctl join_cluster rabbit@rabbit1
#rabbitmqctl start_app

#rabbit3执行相同的命令
#
#查询集群状态
#rabbitmactl cluster_status
#
#5.故障节点的处理
#docker exec -it rabbit2 /bin/bash
#rabbitmqctl stop
##在一个正常的节点上移除有问题的节点
#rabbitmqctl  -n rabbit@rabbit1 forget_cluster_node rabbit@rabbit2
