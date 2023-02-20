#!/bin/bash
###
# @Author: licat
# @Date: 2023-01-14 15:53:16
 # @LastEditors: licat
 # @LastEditTime: 2023-02-20 10:40:39
# @Description: licat233@gmail.com
###

#进入monitor mode
set -m

current_path=$(
    cd $(dirname $0)
    pwd
)

cd $current_path

./exemples/update.sh

comment=$1

if [ -z "$comment" ]; then
    comment="optimize"
fi

git tag -a "v1.3.0" -m "$comment"
git add .
git commit -m $comment
git push -u origin main
