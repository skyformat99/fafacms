#!/bin/bash

mkdir /data
chmod 777 /data
mkdir /data/mydocker
chmod 777 /data/mydocker
mkdir -p /data/mydocker/redis/data
mkdir -p /data/mydocker/redis/conf
mkdir -p /data/mydocker/mysql/data
mkdir -p /data/mydocker/mysql/conf
mkdir -p /data/mydocker/fafacms
mkdir -p /data/mydocker/fafacms/storage
mkdir -p /data/mydocker/fafacms/storage_x
mkdir -p /data/mydocker/fafacms/log
cp my.cnf /data/mydocker/mysql/conf/my.cnf
cp redis.conf /data/mydocker/redis/conf/redis.conf
cp config_mac.json /data/mydocker/fafacms/config.json
docker-compose stop
docker-compose rm -f
docker-compose up -d