/*
 * @Author: licat
 * @Date: 2023-02-07 14:56:28
 * @LastEditors: licat
 * @LastEditTime: 2023-02-17 22:33:51
 * @Description: licat233@gmail.com
 */

package _api

import (
	"fmt"

	"github.com/licat233/sql2rpc/cmd/common"
)

type Api struct {
	Method  string
	Path    string
	Handler string
	Req     string
	Resp    string
	Comment string
}

func NewApi(method, path, handler, req, resp, comment string) *Api {
	if path == "" {
		path = "/"
	}
	return &Api{
		Method:  method,
		Path:    path,
		Handler: handler,
		Req:     req,
		Resp:    resp,
		Comment: comment,
	}
}

func (r *Api) String() string {
	doc := fmt.Sprintf("%s@doc \"%s\"", common.Indent, r.Comment)
	handler := fmt.Sprintf("%s@handler %s", common.Indent, r.Handler)
	api := fmt.Sprintf("%s%s %s(%s) returns (%s)", common.Indent, r.Method, r.Path, r.Req, r.Resp)
	return fmt.Sprintf("\n%s\n%s\n%s\n", doc, handler, api)
}
