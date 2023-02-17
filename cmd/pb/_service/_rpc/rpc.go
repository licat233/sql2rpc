/*
 * @Author: licat
 * @Date: 2023-02-06 21:49:04
 * @LastEditors: licat
 * @LastEditTime: 2023-02-17 22:48:48
 * @Description: licat233@gmail.com
 */

package _rpc

import (
	"fmt"

	"github.com/licat233/sql2rpc/cmd/common"
	"github.com/licat233/sql2rpc/tools"
)

type Rpc struct {
	Name    string
	Req     string
	Resp    string
	Comment string
}

func NewRpc(name, req, resp, comment string) *Rpc {
	return &Rpc{
		Name:    name,
		Req:     req,
		Resp:    resp,
		Comment: comment,
	}
}

func (r *Rpc) String() string {
	comment := fmt.Sprintf("\n%s//%s", common.Indent, r.Comment)
	rpcContent := fmt.Sprintf("\n%srpc %s(%s) returns (%s);", common.Indent, tools.ToCamel(r.Name), r.Req, r.Resp)
	return comment + rpcContent
}
