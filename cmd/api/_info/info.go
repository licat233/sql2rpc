/*
 * @Author: licat
 * @Date: 2023-02-17 12:15:12
 * @LastEditors: licat
 * @LastEditTime: 2023-02-17 13:01:50
 * @Description: licat233@gmail.com
 */
package _info

import (
	"bytes"
	"fmt"
	"time"

	"github.com/licat233/sql2rpc/cmd/api/_conf"
	"github.com/licat233/sql2rpc/cmd/common"
	"github.com/licat233/sql2rpc/config"
	"github.com/licat233/sql2rpc/tools"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

type Info struct {
	*orderedmap.OrderedMap[string, any]
}

func New() *Info {
	defaultInfo := orderedmap.New[string, any]()
	defaultInfo.Set("author", tools.GetCurrentUserName())
	defaultInfo.Set("date", time.Now().Format("2006-01-02 15:04:05"))
	defaultInfo.Set("desc", config.C.ServiceName.GetString()+"API service. "+"Generated by sql2rpc: "+config.ProjectURL)
	return &Info{
		OrderedMap: defaultInfo,
	}
}

func (info *Info) String() string {
	if !_conf.CurrentIsCoreFile {
		return ""
	}
	buf := new(bytes.Buffer)
	buf.WriteString(fmt.Sprintf("\n%s", config.InfoStartMark))
	buf.WriteString("\ninfo(")
	if res := common.PickInfoContent(_conf.FileContent); res != "" {
		buf.WriteString(res)
	} else {
		for pair := info.Oldest(); pair != nil; pair = pair.Next() {
			buf.WriteString(fmt.Sprintf("\n\t%s: \"%v\"", pair.Key, pair.Value))
		}
	}
	buf.WriteString("\n)\n")
	return buf.String()
}