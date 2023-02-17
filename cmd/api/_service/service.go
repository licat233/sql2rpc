/*
 * @Author: licat
 * @Date: 2023-02-03 19:34:01
 * @LastEditors: licat
 * @LastEditTime: 2023-02-18 00:27:17
 * @Description: licat233@gmail.com
 */

package _service

import (
	"bytes"
	"fmt"

	"github.com/licat233/sql2rpc/cmd/api/_service/_api"
	"github.com/licat233/sql2rpc/cmd/common"

	"github.com/licat233/sql2rpc/config"
	"github.com/licat233/sql2rpc/tools"
)

type ServerConfig struct {
	Name       string
	Jwt        string
	Group      string
	Middleware string
	Prefix     string
}

func NewServerConfig(name, jwt, group, middleware, prefix string) *ServerConfig {
	return &ServerConfig{
		Name:       name,
		Jwt:        jwt,
		Group:      group,
		Middleware: middleware,
		Prefix:     prefix,
	}
}

func (sc *ServerConfig) String() string {
	name := tools.ToLowerCamel(sc.Name)
	var buf bytes.Buffer
	buf.WriteString("@server(\n")
	if sc.Jwt != "" {
		buf.WriteString(fmt.Sprintf("%sjwt: %s\n", common.Indent, sc.Jwt))
	}
	buf.WriteString(fmt.Sprintf("%sgroup: %s\n", common.Indent, name))
	if sc.Middleware != "" {
		buf.WriteString(fmt.Sprintf("%smiddleware: %s\n", common.Indent, sc.Middleware))
	}
	if name == "" {
		buf.WriteString(fmt.Sprintf("%sprefix: %s\n", common.Indent, sc.Prefix))
	} else {
		buf.WriteString(fmt.Sprintf("%sprefix: %s/%s\n", common.Indent, sc.Prefix, name))
	}
	buf.WriteString(")\n")
	return buf.String()
}

type Service struct {
	Name    string
	Comment string
	Apis    _api.ApiCollection
	Config  *ServerConfig
}

func NewService(name, comment string) *Service {
	s := &Service{
		Name:    name,
		Comment: comment,
		Apis:    _api.ApiCollection{},
		Config:  NewServerConfig(name, config.C.ApiJwt.GetString(), name, config.C.ApiMiddleware.GetString(), config.C.ApiPrefix.GetString()),
	}
	s.initBaseApiServiceItems()
	return s
}

func (s *Service) String() string {
	buf := new(bytes.Buffer)
	buf.WriteString(fmt.Sprintf("\n// %s \n", s.Comment))
	buf.WriteString(fmt.Sprint(s.Config))
	buf.WriteString(fmt.Sprintf("service %s {\n", config.C.ServiceName.GetString()))
	buf.WriteString(fmt.Sprint(s.Apis))
	buf.WriteString("\n}\n\n")
	return buf.String()
}

func (s *Service) initBaseApiServiceItems() {
	name := tools.ToCamel(s.Name)
	s.Apis = []*_api.Api{
		_api.NewApi("post", "/", "Add"+name, "Add"+name+"Req", "BaseResp", "添加"+s.Comment+" 基础API"),
		_api.NewApi("put", "/", "Put"+name, "Put"+name+"Req", "BaseResp", "更新"+s.Comment+" 基础API"),
		_api.NewApi("get", "/", "Get"+name, "Get"+name+"Req", "BaseResp", "获取"+s.Comment+" 基础API"),
		_api.NewApi("delete", "/", "Del"+name, "Del"+name+"Req", "BaseResp", "删除"+s.Comment+" 基础API"),
		_api.NewApi("get", "/list", "Get"+name+"List", "Get"+name+"ListReq", "BaseResp", "获取"+s.Comment+"列表"+" 基础API"),
		_api.NewApi("get", "/enums", "Get"+name+"Enums", "Get"+name+"EnumsReq", "BaseResp", "获取"+s.Comment+"枚举列表"+" 基础API"),
	}
}

func GenarateDefaultService() string {
	svcConfig := NewServerConfig("", config.C.ApiJwt.GetString(), "", config.C.ApiMiddleware.GetString(), config.C.ApiPrefix.GetString())
	startMark := fmt.Sprintf("\n%s\n", config.CustomServiceStartMark)
	endMark := fmt.Sprintf("\n%s\n", config.CustomServiceEndMark)
	buf := new(bytes.Buffer)
	buf.WriteString(startMark)
	buf.WriteString(svcConfig.String())
	buf.WriteString(fmt.Sprintf("service %s {\n", config.C.ServiceName.GetString()))
	buf.WriteString("\n}\n")
	buf.WriteString(endMark)
	return buf.String()
}
