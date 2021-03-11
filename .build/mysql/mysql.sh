#!/bin/bash
set -e

if [ ! -z $(docker images -q mysql:latest) ]; then
  docker pull mysql
fi

docker run -it --name mysql-test -p 3306:3306 -e MYSQL_ROOT_PASSWORD=1234567 mysql

#登录mysql
mysql -u root -p
ALTER USER 'root'@'localhost' IDENTIFIED BY '1234567'

#添加远程登录用户
CREATE USER 'crochee'@'%' IDENTIFIED WITH mysql_native_password BY '1234567'
GRANT ALL PRIVILEGES ON *.* TO 'crochee'@'%'
exit
