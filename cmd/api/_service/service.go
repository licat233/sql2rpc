/*
 * @Author: licat
 * @Date: 2023-02-03 19:34:01
 * @LastEditors: licat
 * @LastEditTime: 2023-02-18 18:23:21
 * @Description: licat233@gmail.com
 */

package _service

import (
	"bytes"
	"fmt"

	"github.com/licat233/sql2rpc/cmd/api/_service/_api"
	"github.com/licat233/sql2rpc/cmd/api/_service/_server"

	"github.com/licat233/sql2rpc/config"
	"github.com/licat233/sql2rpc/tools"
)

type Service struct {
	Name    string
	Comment string
	Apis    _api.ApiCollection
	Server  *_server.Server
}

func New(name, comment string) *Service {
	s := &Service{
		Name:    name,
		Comment: comment,
		Apis:    _api.ApiCollection{},
		Server:  _server.NewServer(name, config.C.ApiJwt.GetString(), name, config.C.ApiMiddleware.GetString(), config.C.ApiPrefix.GetString()),
	}
	s.initBaseApiServiceItems()
	return s
}

func (s *Service) String() string {
	buf := new(bytes.Buffer)
	if s.Comment != "" {
		buf.WriteString(fmt.Sprintf("\n// %s \n", s.Comment))
	}
	buf.WriteString(fmt.Sprint(s.Server))
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

func GenarateDefaultCustomService() string {
	defSrv := New("", "")
	defSrv.Apis = nil
	buf := new(bytes.Buffer)
	buf.WriteString(fmt.Sprintf("\n%s\n", config.CustomServiceStartMark))
	buf.WriteString(fmt.Sprint(defSrv))
	buf.WriteString(fmt.Sprintf("\n%s\n", config.CustomServiceEndMark))

	return buf.String()
}
