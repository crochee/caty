#!/bin/bash
set -e
# 官网地址 https://github.com/etcd-io/etcd/releases
image=
if [ -n "$(docker images -q quay.io/coreos/etcd:v3.2.32)" ]; then
  docker pull quay.io/coreos/etcd:v3.2.32
fi

TOKEN=etcd-cluster-k8s-micro-token
CLUSTER_STATE=new
NAME_1=etcd_node0
HOST_1=0.0.0.0
CLUSTER=${NAME_1}=http://${HOST_1}:2380

rm -rf /home/lcf/cloud/etcd/tmp/etcd-data.tmp && mkdir -p /home/lcf/cloud/etcd/tmp/etcd-data.tmp && \
docker run \
-itd \
-p 2379:2379 \
-p 2380:2380 \
--mount type=bind,source=/home/lcf/cloud/etcd/tmp/etcd-data.tmp,destination=/etcd-data \
--name etcd-gcr-v3.2.32 \
quay.io/coreos/etcd:v3.2.32 \
/usr/local/bin/etcd \
--name ${NAME_1} \
--data-dir /etcd-data \
--listen-client-urls http://${HOST_1}:2379 \
--advertise-client-urls http://${HOST_1}:2379 \
--listen-peer-urls http://${HOST_1}:2380 \
--initial-advertise-peer-urls http://${HOST_1}:2380 \
--initial-cluster ${CLUSTER} \
--initial-cluster-token ${TOKEN} \
--initial-cluster-state ${CLUSTER_STATE}

#docker run \
#  -p 2379:2379 \
#  -p 2380:2380 \
#  --mount type=bind,source=/tmp/etcd-data.tmp,destination=/etcd-data \
#  --name etcd_node1 \
#  quay.io/coreos/etcd:v3.4.14 \
#  /usr/local/bin/etcd \
#  --name etcd_node1 \
#  --data-dir /etcd-data \
#  --listen-client-urls http://0.0.0.0:2379 \
#  --advertise-client-urls http://172.24.149.46:2379 \
#  --listen-peer-urls http://0.0.0.0:2380 \
#  --initial-advertise-peer-urls http://172.24.149.46:2380 \
#  --initial-cluster-token etcd-cluster-k8s \
#  --initial-cluster-state new \
#  --initial-cluster etcd_node0=http://172.24.152.17:2380,etcd_node1=http://172.24.149.46:2380,etcd_node2=http://172.24.155.15:2380 \
#  --log-level info \
#  --logger zap \
#  --log-outputs stderr

#--name：节点名称，默认为 default。
#--data-dir：服务运行数据保存的路径，默认为${name}.etcd。
#--snapshot-count：指定有多少事务（transaction）被提交时，触发截取快照保存到磁盘。
#--heartbeat-interval：leader 多久发送一次心跳到 followers。默认值是 100ms。
#--eletion-timeout：重新投票的超时时间，如果 follow 在该时间间隔没有收到心跳包，会触发重新投票，默认为 1000 ms。
#--listen-peer-urls：和同伴通信的地址，比如http://ip:2380，如果有多个，使用逗号分隔。需要所有节点都能够访问，所以不要使用 localhost！
#--listen-client-urls：对外提供服务的地址：比如http://ip:2379,http://127.0.0.1:2379，客户端会连接到这里和 etcd 交互。
#--advertise-client-urls：对外公告的该节点客户端监听地址，这个值会告诉集群中其他节点。
#--initial-advertise-peer-urls：该节点同伴监听地址，这个值会告诉集群中其他节点。
#--initial-cluster：集群中所有节点的信息，格式为node1=http://ip1:2380,node2=http://ip2:2380,…，注意：这里的 node1 是节点的 --name 指定的名字；后面的 ip1:2380 是 --initial-advertise-peer-urls 指定的值。
#--initial-cluster-state：新建集群的时候，这个值为 new；假如已经存在的集群，这个值为 existing。
#--initial-cluster-token：创建集群的 token，这个值每个集群保持唯一。这样的话，如果你要重新创建集群，即使配置和之前一样，也会再次生成新的集群和节点 uuid；否则会导致多个集群之间的冲突，造成未知的错误
