/*
 * @Author: licat
 * @Date: 2023-02-07 22:35:57
 * @LastEditors: licat
 * @LastEditTime: 2023-02-16 14:48:17
 * @Description: licat233@gmail.com
 */

package _conf

import "github.com/licat233/sql2rpc/config"

var (
	FileContent       string       //文件内容
	baseIgnoreTables  = []string{} //当前服务必须忽略的表
	baseIgnoreColumns = []string{} //每个结构体必须忽略的列

	MoreIgnoreTables  = []string{}             //当前服务可能忽略的表
	MoreIgnoreColumns = []string{"id", "uuid"} //某个结构可能忽略的列

	IgnoreTables  = []string{} //当前解析需要忽略的表
	IgnoreColumns = []string{} //每个结构体需要忽略的列

	CurrentIsCoreFile = true //标记当前是否为核心文件，用于生成唯一内容，避免重复
)

func InitConfig() {
	IgnoreTables = append(config.IgnoreTables, baseIgnoreTables...)
	IgnoreColumns = append(config.IgnoreColumns, baseIgnoreColumns...)
}
