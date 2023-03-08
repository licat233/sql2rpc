#!/bin/bash
###
# @Author: licat
# @Date: 2023-01-11 15:40:07
 # @LastEditors: licat
 # @LastEditTime: 2023-02-18 11:42:31
# @Description: licat233@gmail.com
###

#进入monitor mode
set -m

current_path=$(
    cd "$(dirname "$0")" || exit
    pwd
)

cd "${current_path}" || exit

chmod +x ./*glob*.sh

rm -f ./*glob*.api
rm -f ./*glob*.proto

./build.sh
./model_file/gen.sh
./all_pb_api/gen.sh
./api_multiple_file/gen.sh
./api_single_file/gen.sh
./pb_multiple_group/gen.sh
./pb_single_group/gen.sh
