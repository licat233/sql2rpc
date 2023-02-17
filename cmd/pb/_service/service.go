/*
 * @Author: licat
 * @Date: 2023-02-03 19:34:01
 * @LastEditors: licat
 * @LastEditTime: 2023-02-17 22:51:39
 * @Description: licat233@gmail.com
 */
package _service

import (
	"bytes"
	"fmt"

	"github.com/licat233/sql2rpc/cmd/common"
	"github.com/licat233/sql2rpc/cmd/pb/_service/_rpc"

	"github.com/licat233/sql2rpc/tools"
)

type Service struct {
	Name    string
	Comment string
	Rpcs    _rpc.RpcCollection
}

func NewService(name, comment string) *Service {
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
	// if config.C.PbMultiple.GetBool() {
	// 	buf.WriteString(fmt.Sprintf("\n// %s base service (%s) ", s.Name, s.Comment))
	// 	buf.WriteString(fmt.Sprintf("\nservice %sBase {\n\n", tools.ToCamel(s.Name)))
	// 	buf.WriteString(fmt.Sprint(s.Rpcs))
	// 	buf.WriteString("\n}\n")
	// } else {
	buf.WriteString(fmt.Sprintf("\n\n%s// %s base service (%s)  \n", common.Indent, s.Name, s.Comment))
	buf.WriteString(fmt.Sprint(s.Rpcs))
	// }
	return buf.String()
}

func (s *Service) initBaseServiceRpcs() {
	name := tools.ToCamel(s.Name)
	s.Rpcs = []*_rpc.Rpc{
		_rpc.NewRpc("BaseAdd"+name, "Add"+name+"Req", "Add"+name+"Resp", "添加"+s.Comment),
		_rpc.NewRpc("BasePut"+name, "Put"+name+"Req", "Put"+name+"Resp", "更新"+s.Comment),
		_rpc.NewRpc("BaseGet"+name, "Get"+name+"Req", "Get"+name+"Resp", "获取"+s.Comment),
		_rpc.NewRpc("BaseDel"+name, "Del"+name+"Req", "Del"+name+"Resp", "删除"+s.Comment),
		_rpc.NewRpc("BaseGet"+name+"List", "Get"+name+"ListReq", "Get"+name+"ListResp", "获取"+s.Comment+"列表"),
		_rpc.NewRpc("BaseGet"+name+"Enums", "Get"+name+"EnumsReq", "Enums", "获取"+s.Comment+"枚举列表"),
	}
}
