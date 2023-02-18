/*
 * @Author: licat
 * @Date: 2023-02-18 10:01:53
 * @LastEditors: licat
 * @LastEditTime: 2023-02-18 20:11:51
 * @Description: licat233@gmail.com
 */

package _server

import (
	"bytes"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/licat233/sql2rpc/cmd/common"
	"github.com/licat233/sql2rpc/tools"
)

type Server struct {
	Name       string
	Jwt        string
	Group      string
	Middleware string
	Prefix     string
}

func NewServer(name, jwt, group, middleware, prefix string) *Server {
	return &Server{
		Name:       name,
		Jwt:        jwt,
		Group:      group,
		Middleware: middleware,
		Prefix:     prefix,
	}
}

func (sc *Server) String() string {
	name := tools.ToLowerCamel(sc.Name)
	prefixValue, err := url.JoinPath(sc.Prefix, name)
	if err != nil {
		log.Fatalln("api url path error:", err)
	}
	var buf bytes.Buffer
	buf.WriteString("@server(")
	if sc.Jwt != "" {
		buf.WriteString(fmt.Sprintf("\n%sjwt: %s", common.Indent, sc.Jwt))
	}
	if sc.Group != "" {
		buf.WriteString(fmt.Sprintf("\n%sgroup: %s", common.Indent, name))
	}
	if sc.Middleware != "" {
		buf.WriteString(fmt.Sprintf("\n%smiddleware: %s", common.Indent, toCamelHandler(sc.Middleware)))
	}

	if prefixValue != "" {
		buf.WriteString(fmt.Sprintf("\n%sprefix: %s", common.Indent, prefixValue))
	}
	buf.WriteString("\n)\n")
	return buf.String()
}

func toCamelHandler(value string) string {
	if value == "" {
		return ""
	}
	list := strings.Split(value, ",")
	for i, v := range list {
		list[i] = tools.ToCamel(v)
	}
	return strings.Join(list, ",")
}
