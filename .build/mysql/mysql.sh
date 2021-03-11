#!/bin/bash
set -e

if [ ! -z $(docker images -q mysql:latest) ]; then
  docker pull mysql
fi

docker run -it -p 3307:3306 --name mysql-test --restart=always -v /opt/data/mysql/conf:/etc/mysql -v /opt/data/mysql/data:/var/lib/mysql -v /opt/data/mysql/mysql-files:/var/lib/mysql-files/ -e MYSQL_ROOT_PASSWORD=1234567 mysql
#登录mysql
mysql -u root -p
ALTER USER 'root'@'localhost' IDENTIFIED BY '1234567'
#查看用户信息
select user,host,authentication_string from mysql.user;
#设置权限（为root分配权限，以便可以远程连接）
grant all PRIVILEGES on *.* to root@'%' WITH GRANT OPTION;
#添加远程登录用户
CREATE USER 'crochee'@'%' IDENTIFIED WITH mysql_native_password BY '1234567'
GRANT ALL PRIVILEGES ON *.* TO 'crochee'@'%'


