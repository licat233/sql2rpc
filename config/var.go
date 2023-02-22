/*
 * @Author: licat
 * @Date: 2023-02-07 13:40:54
 * @LastEditors: licat
 * @LastEditTime: 2023-02-22 12:13:56
 * @Description: licat233@gmail.com
 */

package config

const (
	// CurrentVersion 当前项目版本
	CurrentVersion = "v1.3.6"

	// ProjectName 当前项目名称
	ProjectName = "sql2rpc"

	// ProjectURL 当前项目地址
	ProjectURL = "https://github.com/licat233/" + ProjectName
	// ProjectInfoURL 当前项目的信息接口
	ProjectInfoURL = "https://api.github.com/repos/licat233/" + ProjectName + "/releases/latest"

	// DefaultFileName 配置文件名称
	DefaultFileName = "sql2rpcConfig.yaml"

	Syntax = "proto3"

	CamelCase      = "SqlRpc"
	LowerCamelCase = "sqlRpc"
	SnakeCase      = "sql_rpc"

	ApiCoreName = "api" //api服务名
	PbCoreName  = "pb"  //pb服务名

	UpdatedFileMsg = "已更新文件"
	CreatedFileMsg = "已创建文件"
)

var (
	C    *Config = nil
	Info *info   = &info{}

	IgnoreTables  = []string{}
	IgnoreColumns = []string{
		"version",
		"create_time",
		"created_time",
		"create_at",
		"created_at",
		"update_time",
		"updated_time",
		"update_at",
		"updated_at",
		"delete_time",
		"deleted_time",
		"delete_at",
		"deleted_at",
		"del_state",
		"is_deleted",
		"is_delete",
	}

	InfoStartMark, InfoEndMark       = GetMark("Info")
	ImportStartMark, ImportEndMark   = GetMark("Import")
	StructStartMark, StructEndMark   = GetMark("Struct")
	EnumStartMark, EnumEndMark       = GetMark("Enum")
	MessageStartMark, MessageEndMark = GetMark("Message")
	ServiceStartMark, ServiceEndMark = GetMark("Service")

	CustomImportStartMark, CustomImportEndMark   = GetCustomMark("import")
	CustomStructStartMark, CustomStructEndMark   = GetCustomMark("struct")
	CustomEnumStartMark, CustomEnumEndMark       = GetCustomMark("enum")
	CustomMessageStartMark, CustomMessageEndMark = GetCustomMark("message")
	CustomServiceStartMark, CustomServiceEndMark = GetCustomMark("service")
)
