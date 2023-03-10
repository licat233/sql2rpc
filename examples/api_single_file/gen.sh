#!/bin/bash
###
# @Author: licat
# @Date: 2023-01-11 15:40:07
# @LastEditors: licat
# @LastEditTime: 2023-02-18 10:16:41
# @Description: licat233@gmail.com
###

#进入monitor mode
set -m

current_path=$(
    cd $(dirname $0)
    pwd
)

cd $current_path

if [ ! -f "../sql2rpc" ]; then
    ../build.sh
fi

rm -f ./*.api

if ! ../sql2rpc -api -db_schema="admin" -db_table="*" -service_name="admin-api" -filename="admin.api" -api_jwt="Auth" -api_middleware="AuthMiddleware" -api_prefix="/v1/api/admin" -api_multiple=false; then
    exit 1
fi
# exit


if [ -d api ]; then
    rm -rf api/*
fi

# Please install the goctl first
# command: go install github.com/zeromicro/go-zero/tools/goctl@latest
# github: https://github.com/zeromicro/go-zero

# exit # 如果需要生成gozero框架的服务代码，请注释这行

# goctl template init --home ./template

if ! goctl api go -api admin.api -dir ./api -style goZero --home ./template; then
    exit 1
fi

cd ../
go mod tidy
