#!/bin/bash
###
# @Author: licat
# @Date: 2023-01-11 15:40:07
 # @LastEditors: licat
 # @LastEditTime: 2023-02-09 17:14:22
# @Description: licat233@gmail.com
###

#进入monitor mode
set -m

current_path=$(
    cd $(dirname $0)
    pwd
)

cd ${current_path}/../ || exit

if [ -f examples/sql2rpc ]; then
    rm -f examples/sql2rpc
fi

go mod tidy
go mod download
go build -o examples/sql2rpc main.go
