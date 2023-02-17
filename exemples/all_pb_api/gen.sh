#!/bin/bash
###
# @Author: licat
# @Date: 2023-01-11 15:40:07
 # @LastEditors: licat
 # @LastEditTime: 2023-02-17 09:49:04
# @Description: licat233@gmail.com
###

#进入monitor mode
set -m

current_path=$(
    cd $(dirname $0) || exit
    pwd
)

cd $current_path

if [ ! -f "../sql2rpc" ]; then
    ../build.sh
fi

rm -f "./sql2rpcConfig.yaml"

if [ ! -f "./sql2rpcConfig.yaml" ]; then
    cp ../sql2rpcConfig.example.yaml sql2rpcConfig.yaml
fi

# 会根据yaml配置文件来生成服务的配置文件
echo " - 将优先根据sql2rpcConfig.yaml配置文件内容来生成服务"
../sql2rpc

exit  # 如果需要生成gozero框架的服务代码，请注释这行

if [ $? -ne 0 ]; then
    exit 1
fi

if [ -d api ]; then
    rm -rf api/*
fi

if [ -d rpc ]; then
    rm -rf rpc/*
fi

# Please install the goctl first
# command: go install github.com/zeromicro/go-zero/tools/goctl@latest
# github: https://github.com/zeromicro/go-zero
goctl api go -api admin.api -dir ./api -style goZero
if [ $? -ne 0 ]; then
    exit 1
fi

goctl rpc protoc admin.proto --go_out=./rpc --go-grpc_out=./rpc --zrpc_out=./rpc
if [ $? -ne 0 ]; then
    exit 1
fi
cd ../
go mod tidy