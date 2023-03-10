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
    cd $(dirname $0) || exit
    pwd
)

cd "$current_path" || exit

if [ ! -f "../sql2rpc" ]; then
    ../build.sh
fi

rm -f ./model/*.go

goctl model mysql ddl --src "../sql/admin.sql" -dir="model" --style goZero -cache=false
go mod tidy

../sql2rpc -model -db_schema="admin" -db_table="*" -dir="model"

if ! ../sql2rpc -model -db_schema="admin" -db_table="*" -dir="model";
then
    exit 1
fi

if [ -d api ]; then
    rm -rf api/*
fi

# Please install the goctl first
# command: go install github.com/zeromicro/go-zero/tools/goctl@latest
# github: https://github.com/zeromicro/go-zero


exit # 如果需要生成gozero框架的服务代码，请注释这行

goctl api go -api admin.api -dir ./api -style goZero
if [ $? -ne 0 ]; then
    exit 1
fi
cd ../
go mod tidy
