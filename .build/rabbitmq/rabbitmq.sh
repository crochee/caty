#!/bin/bash
set -e

image=docker.io/library/rabbitmq:latest
if [ -n "$(docker images -q ${image})" ]; then
  docker pull ${image}
fi

# 构建容器集群
RABBITMQ_NAME_0=rabbitmq_node_0
RABBITMQ_NAME_1=rabbitmq_node_1
RABBITMQ_NAME_2=rabbitmq_node_2
# node_0
docker run -itd \
--hostname ${RABBITMQ_NAME_0}_host \
--name ${RABBITMQ_NAME_0} \
--log-opt max-size=10m \
--log-opt max-file=3 \
-e RABBITMQ_DEFAULT_USER=admin \
-e RABBITMQ_DEFAULT_PASS='1234567' \
-e RABBITMQ_ERLANG_COOKIE='crochee secret cookie here' \
-p 15672:15672 \
-p 5672:5672 \
--restart=always \
${image}

# node_1
docker run -itd \
--hostname ${RABBITMQ_NAME_1}_host \
--name ${RABBITMQ_NAME_1} \
--log-opt max-size=10m \
--log-opt max-file=3 \
-e RABBITMQ_DEFAULT_USER=admin \
-e RABBITMQ_DEFAULT_PASS='1234567' \
-e RABBITMQ_ERLANG_COOKIE='crochee secret cookie here' \
-p 5673:5672 \
--restart=always \
--link ${RABBITMQ_NAME_0}:${RABBITMQ_NAME_0}_host \
${image}

# node_2
docker run -itd \
--hostname ${RABBITMQ_NAME_2}_host \
--name ${RABBITMQ_NAME_2} \
--log-opt max-size=10m \
--log-opt max-file=3 \
-e RABBITMQ_DEFAULT_USER=admin \
-e RABBITMQ_DEFAULT_PASS='1234567' \
-e RABBITMQ_ERLANG_COOKIE='crochee secret cookie here' \
-p 5674:5672 \
--restart=always \
--link ${RABBITMQ_NAME_0}:${RABBITMQ_NAME_0}_host \
--link ${RABBITMQ_NAME_1}:${RABBITMQ_NAME_1}_host \
${image}
#多个容器之间使用“--link”连接，此属性不能少；
id=$(docker ps | grep ${RABBITMQ_NAME_0} | awk '{print $1}')
echo "${id}"
#开启控制台 并加入集群

# node_0
#docker exec -it ${RABBITMQ_NAME_0} bash
#rabbitmqctl stop_app
#rabbitmqctl reset
#rabbitmq-plugins enable rabbitmq_management 开启管理ui
#rabbitmqctl start_app
#exit

# node_1
#docker exec -it ${RABBITMQ_NAME_1} bash
#rabbitmqctl stop_app
#rabbitmqctl reset
#rabbitmqctl join_cluster --ram rabbit@${RABBITMQ_NAME_0}_host #参数“--ram”表示设置为内存节点，忽略次参数默认为磁盘节点
#rabbitmqctl start_app
#exit

# node_2
#docker exec -it ${RABBITMQ_NAME_2} bash
#rabbitmqctl stop_app
#rabbitmqctl reset
#rabbitmqctl join_cluster --ram rabbit@${RABBITMQ_NAME_0}_host #参数“--ram”表示设置为内存节点，忽略次参数默认为磁盘节点
#rabbitmqctl start_app
#exit

#账户名：admin
#密码：1234567

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
