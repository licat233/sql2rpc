#!/bin/bash
###
# @Author: licat
# @Date: 2023-01-11 15:40:07
 # @LastEditors: licat
 # @LastEditTime: 2023-02-16 15:16:28
# @Description: licat233@gmail.com
###

#进入monitor mode
set -m

current_path=$(
    cd $(dirname $0)
    pwd
)

cd $current_path

../build.sh

if [ ! -f "../sql2rpc" ]; then
    ../build.sh
fi

../sql2rpc -pb -db_schema="admin" -db_table="*" -service_name="Admin" -filename="admin.proto" -pb_package="admin_proto" -pb_gopackage="./admin_pb" -pb_multiple=true

exit # 如果需要生成gozero框架的服务代码，请注释这行

if [ $? -ne 0 ]; then
    exit 1
fi

if [ -d rpc ]; then
    rm -rf rpc/*
fi

# Please install the goctl first
# command: go install github.com/zeromicro/go-zero/tools/goctl@latest
# github: https://github.com/zeromicro/go-zero
goctl rpc protoc "admin.proto" --go_out="./rpc" --go-grpc_out="./rpc" --zrpc_out="./rpc"

cd ../
go mod tidy
