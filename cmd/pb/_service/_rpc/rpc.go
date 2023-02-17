/*
 * @Author: licat
 * @Date: 2023-02-06 21:49:04
 * @LastEditors: licat
 * @LastEditTime: 2023-02-07 23:47:37
 * @Description: licat233@gmail.com
 */

package _rpc

import (
	"fmt"

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
	return fmt.Sprintf("  //%s\n  rpc %s(%s) returns (%s);\n", r.Comment, tools.ToCamel(r.Name), r.Req, r.Resp)
}
