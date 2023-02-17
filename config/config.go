/*
 * @Author: licat
 * @Date: 2023-02-06 14:29:01
 * @LastEditors: licat
 * @LastEditTime: 2023-02-17 17:19:51
 * @Description: licat233@gmail.com
 */
package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/fatih/structs"
	"github.com/licat233/sql2rpc/config/_item"
	"github.com/licat233/sql2rpc/tools"
)

type info struct {
	Version string
}

type Config struct {
	//数据库配置
	DBType     _item.Field
	DBHost     _item.Field
	DBPort     _item.Field
	DBUser     _item.Field
	DBPassword _item.Field
	DBSchema   _item.Field
	DBTable    _item.Field

	//ingore配置
	IgnoreTableStr  _item.Field
	IgnoreColumnStr _item.Field

	//通用配置
	ServiceName _item.Field
	Filename    _item.Field

	//pb配置
	Pb          _item.Field
	PbPackage   _item.Field
	PbGoPackage _item.Field
	PbMultiple  _item.Field

	// api配置
	Api           _item.Field
	Style         _item.Field
	ApiJwt        _item.Field
	ApiMiddleware _item.Field
	ApiPrefix     _item.Field
	ApiMultiple   _item.Field
}

// 获取一个默认的Config
func NewDefaultConfig() *Config {
	return &Config{
		DBType:          _item.NewConfigItem("db_type", "mysql", "the database type", false),
		DBHost:          _item.NewConfigItem("db_host", "localhost", "the database host", false),
		DBPort:          _item.NewConfigItem("db_port", 3306, "the database port", false),
		DBUser:          _item.NewConfigItem("db_user", "root", "the database user", false),
		DBPassword:      _item.NewConfigItem("db_password", "", "the database password", false),
		DBSchema:        _item.NewConfigItem("db_schema", "", "the database schema", true),
		DBTable:         _item.NewConfigItem("db_table", "*", "the database table, split multiple tables with \",\"", false),
		IgnoreTableStr:  _item.NewConfigItem("ignore_table", "", "the table to ignore, split multiple value by \",\"", false),
		IgnoreColumnStr: _item.NewConfigItem("ignore_column", "", "the column to ignore, split multiple value by \",\"", false),
		ServiceName:     _item.NewConfigItem("service_name", "", "the service name, defaults to the database schema", false),
		Filename:        _item.NewConfigItem("filename", "", "the generated file name, defaults to the service name", false),
		Pb:              _item.NewConfigItem("pb", false, "generate .proto files", false),
		PbPackage:       _item.NewConfigItem("pb_package", "", "the protocol buffer package, defaults to the service name", false),
		PbGoPackage:     _item.NewConfigItem("pb_gopackage", "", "the protocol buffer go_package, defaults to the service name", false),
		PbMultiple:      _item.NewConfigItem("pb_multiple", false, "the generated in multiple rpc service mode", false),
		Api:             _item.NewConfigItem("api", false, "generate .api files", false),
		Style:           _item.NewConfigItem("api_style", SnakeCase, fmt.Sprintf("the struct json naming format: %s | %s | %s ", SnakeCase, CamelCase, LowerCamelCase), false),
		ApiJwt:          _item.NewConfigItem("api_jwt", "", "the api service jwt, example: Auth", false),
		ApiMiddleware:   _item.NewConfigItem("api_middleware", "", "the api service middleware,  split multiple value by \",\", example: AuthMiddleware", false),
		ApiPrefix:       _item.NewConfigItem("api_prefix", "", "the api service route prefix, example: api", false),
		ApiMultiple:     _item.NewConfigItem("api_multiple", false, "Generate multiple api files according to table", false),
	}
}

/**
 * @description: Assignment 合并yamlConfig、cmdConfig、defaultConfig，优先级由左到右
 * @param {*Config} cmdConfig
 * @return {*Config} config
 */
func (c *Config) Assignment(cmdConfig *Config) *Config {
	if err := InitViper(); err != nil && err != os.ErrNotExist {
		log.Fatal(err)
	}
	return &Config{
		DBType:          c.DBType.Init(cmdConfig.DBType.Value()),
		DBHost:          c.DBHost.Init(cmdConfig.DBHost.Value()),
		DBPort:          c.DBPort.Init(cmdConfig.DBPort.Value()),
		DBUser:          c.DBUser.Init(cmdConfig.DBUser.Value()),
		DBPassword:      c.DBPassword.Init(cmdConfig.DBPassword.Value()),
		DBSchema:        c.DBSchema.Init(cmdConfig.DBSchema.Value()),
		DBTable:         c.DBTable.Init(cmdConfig.DBTable.Value()),
		IgnoreTableStr:  c.IgnoreTableStr.Init(cmdConfig.IgnoreTableStr.Value()),
		IgnoreColumnStr: c.IgnoreColumnStr.Init(cmdConfig.IgnoreColumnStr.Value()),
		ServiceName:     c.ServiceName.Init(cmdConfig.ServiceName.Value()),
		Filename:        c.Filename.Init(cmdConfig.Filename.Value()),
		Pb:              c.Pb.Init(cmdConfig.Pb.Value()),
		PbPackage:       c.PbPackage.Init(cmdConfig.PbPackage.Value()),
		PbGoPackage:     c.PbGoPackage.Init(cmdConfig.PbGoPackage.Value()),
		PbMultiple:      c.PbMultiple.Init(cmdConfig.PbMultiple.Value()),
		Api:             c.Api.Init(cmdConfig.Api.Value()),
		Style:           c.Style.Init(cmdConfig.Style.Value()),
		ApiJwt:          c.ApiJwt.Init(cmdConfig.ApiJwt.Value()),
		ApiMiddleware:   c.ApiMiddleware.Init(cmdConfig.ApiMiddleware.Value()),
		ApiPrefix:       c.ApiPrefix.Init(cmdConfig.ApiPrefix.Value()),
		ApiMultiple:     c.ApiMultiple.Init(cmdConfig.ApiMultiple.Value()),
	}
}

// Validate 验证配置是否正确
func (c *Config) Validate() error {
	if c.DBSchema.GetString() == "" {
		alert := "\n - please set the database schema \n"
		alert += " - example command: \n"
		alert += fmt.Sprintf("   %s -db_schema admin \n", ProjectName)
		alert += " - or use the config file，the command: \n"
		alert += fmt.Sprintf("   %s -init\n", ProjectName)
		alert += fmt.Sprintf(" Run the \"%s -h\" command for help\n", ProjectName)
		return errors.New(alert)
	}
	return nil
}

// Initialize 初始化配置，自动完善部分配置
func (c *Config) Initialize() {
	if c.ServiceName.GetString() == "" {
		c.ServiceName.Set(tools.ToCamel(c.DBSchema.GetString()))
	}
	if c.Filename.GetString() == "" {
		c.Filename.Set(tools.ToLowerCamel(c.ServiceName.GetString()))
	}
	if c.PbPackage.GetString() == "" {
		c.PbPackage.Set(tools.ToSnake(c.ServiceName.GetString()) + "_proto")
	}
	if c.PbGoPackage.GetString() == "" {
		c.PbGoPackage.Set("./" + tools.ToSnake(c.ServiceName.GetString()) + "_pb")
	}

	IgnoreTables = append(IgnoreTables, strings.Split(strings.Trim(C.IgnoreTableStr.GetString(), ","), ",")...)
	IgnoreColumns = append(IgnoreColumns, strings.Split(strings.Trim(C.IgnoreColumnStr.GetString(), ","), ",")...)
}

func (c *Config) String1() string {
	buf := new(bytes.Buffer)
	m := structs.Map(c)
	for name, value := range m {
		if _, ok := value.(string); ok {
			buf.WriteString(fmt.Sprintf("%s: \"%+v\"\n", name, value))
			continue
		}
		buf.WriteString(fmt.Sprintf("%s: %+v\n", name, value))
	}
	return buf.String()
}

// 这里用于生产yaml配置文件内容
func (c *Config) String() string {
	buf := new(bytes.Buffer)
	buf.WriteString(fmt.Sprintf("# Generated by %s: %s\n", ProjectName, ProjectURL))
	buf.WriteString(fmt.Sprintf("# Version: %s\n", CurrentVersion))
	cV := reflect.ValueOf(c).Elem()
	for i, n := 0, cV.NumField(); i < n; i++ {
		buf.WriteString(fmt.Sprint(cV.Field(i)))
	}
	return buf.String()
}

func (c *Config) ToJson() string {
	b, err := json.Marshal(*c)
	if err != nil {
		return fmt.Sprintf("%+v", *c)
	}
	buf := new(bytes.Buffer)
	err = json.Indent(buf, b, "", "    ")
	if err != nil {
		return fmt.Sprintf("%+v", *c)
	}
	return buf.String()
}
