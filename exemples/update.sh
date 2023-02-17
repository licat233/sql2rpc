#!/bin/bash
###
# @Author: licat
# @Date: 2023-01-11 15:40:07
 # @LastEditors: licat
 # @LastEditTime: 2023-02-18 00:32:32
# @Description: licat233@gmail.com
###

#进入monitor mode
set -m

current_path=$(
    cd $(dirname $0)
    pwd
)

cd ${current_path}

chmod +x **/*.sh

rm -f **/*.api
rm -f **/*.proto

./build.sh
./all_pb_api/gen.sh
./api_multiple_file/gen.sh
./api_single_file/gen.sh
./pb_multiple_group/gen.sh
./pb_single_group/gen.sh
