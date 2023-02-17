/*
 * @Author: licat
 * @Date: 2023-02-07 14:56:28
 * @LastEditors: licat
 * @LastEditTime: 2023-02-07 14:57:44
 * @Description: licat233@gmail.com
 */

package _api

import "fmt"

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
	doc := fmt.Sprintf("  @doc \"%s\"\n", r.Comment)
	handler := fmt.Sprintf("  @handler %s\n", r.Handler)
	api := fmt.Sprintf("  %s %s(%s) returns (%s)\n", r.Method, r.Path, r.Req, r.Resp)
	return fmt.Sprintf("%s%s%s\n", doc, handler, api)
}
