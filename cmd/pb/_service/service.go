/*
 * @Author: licat
 * @Date: 2023-02-03 19:34:01
 * @LastEditors: licat
 * @LastEditTime: 2023-02-18 12:38:20
 * @Description: licat233@gmail.com
 */
package _service

import (
	"bytes"
	"fmt"

	"github.com/licat233/sql2rpc/cmd/common"
	"github.com/licat233/sql2rpc/cmd/pb/_service/_rpc"
	"github.com/licat233/sql2rpc/config"

	"github.com/licat233/sql2rpc/tools"
)

type Service struct {
	Name    string
	Comment string
	Rpcs    _rpc.RpcCollection
}

func New(name, comment string) *Service {
	s := &Service{
		Name:    name,
		Comment: comment,
		Rpcs:    nil,
	}
	s.initBaseServiceRpcs()
	return s
}

func (s *Service) String() string {
	buf := new(bytes.Buffer)
	buf.WriteString(fmt.Sprintf("\n\n%s// %s rpc  \n", common.Indent, s.Name))
	buf.WriteString(fmt.Sprint(s.Rpcs))
	buf.WriteString("\n")
	return buf.String()
}

func (s *Service) initBaseServiceRpcs() {
	name := tools.ToCamel(s.Name)
	s.Rpcs = []*_rpc.Rpc{
		_rpc.New("BaseAdd"+name, "Add"+name+"Req", "Add"+name+"Resp", "添加"+s.Comment),
		_rpc.New("BasePut"+name, "Put"+name+"Req", "Put"+name+"Resp", "更新"+s.Comment),
		_rpc.New("BaseGet"+name, "Get"+name+"Req", "Get"+name+"Resp", "获取"+s.Comment),
		_rpc.New("BaseDel"+name, "Del"+name+"Req", "Del"+name+"Resp", "删除"+s.Comment),
		_rpc.New("BaseGet"+name+"List", "Get"+name+"ListReq", "Get"+name+"ListResp", "获取"+s.Comment+"列表"),
		_rpc.New("BaseGet"+name+"Enums", "Get"+name+"EnumsReq", "Enums", "获取"+s.Comment+"枚举列表"),
	}
}

func GenarateDefaultCustomService() string {
	svrName := tools.ToCamel(config.C.ServiceName.GetString())
	var buf = new(bytes.Buffer)
	buf.WriteString(fmt.Sprintf("\n%s\n", config.CustomServiceStartMark))
	buf.WriteString("\n// " + svrName + " service")
	buf.WriteString("\nservice " + svrName + " {")
	buf.WriteString("\n}\n")
	buf.WriteString(fmt.Sprintf("\n%s\n", config.CustomServiceEndMark))
	return buf.String()
}
