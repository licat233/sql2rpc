/*
 * @Author: licat
 * @Date: 2023-02-03 19:51:18
 * @LastEditors: licat
 * @LastEditTime: 2023-02-20 17:11:24
 * @Description: licat233@gmail.com
 */
package _message

import (
	"bytes"
	"fmt"
	"log"

	"github.com/licat233/sql2rpc/cmd/pb/_conf"
	"github.com/licat233/sql2rpc/cmd/pb/_message/_field"

	"github.com/licat233/sql2rpc/cmd/common"
	"github.com/licat233/sql2rpc/config"
)

// MessageCollection represents a sortable collection of messages.
type MessageCollection []*Message

var baseMessageCollection MessageCollection = []*Message{
	New("Enum", "枚举", _field.MessageFieldCollection{
		_field.New("string", "label", 1, "标签"),
		_field.New("int64", "value", 2, "值"),
	}),
	New("Enums", "枚举列表", _field.MessageFieldCollection{
		_field.New("repeated Enum", "enums", 1, "枚举列表数据"),
	}),
	New("Option", "选项", _field.MessageFieldCollection{
		_field.New("string", "title", 1, "标题"),
		_field.New("int64", "value", 2, "值"),
	}),
	New("Options", "选项列表", _field.MessageFieldCollection{
		_field.New("repeated Option", "options", 1, "选项列表数据"),
	}),
	New("TreeOption", "树形选项", _field.MessageFieldCollection{
		_field.New("string", "title", 1, "标题"),
		_field.New("int64", "value", 2, "值"),
		_field.New("repeated TreeOption", "children", 3, "子集"),
	}),
	New("TreeOptions", "树形选项列表", _field.MessageFieldCollection{
		_field.New("repeated TreeOption", "tree_options", 1, "树形选项列表数据"),
	}),
	New("StatusResp", "状态响应", _field.MessageFieldCollection{
		_field.New("bool", "status", 1, "状态"),
	}),
	New("ListReq", "列表数据请求", _field.MessageFieldCollection{
		_field.New("int64", "page_size", 1, "页容量"),
		_field.New("int64", "page", 2, "页码"),
		_field.New("string", "keyword", 3, "关键词"),
	}),
	New("ByIdReq", "通过ID请求", _field.MessageFieldCollection{
		_field.New("int64", "id", 1, "主键"),
	}),
	New("NilReq", "空请求", nil),
	New("NilResp", "空响应", nil),
}

func (mc MessageCollection) Len() int {
	return len(mc)
}

func (mc MessageCollection) Less(i, j int) bool {
	return mc[i].Name < mc[j].Name
}

func (mc MessageCollection) Swap(i, j int) {
	mc[i], mc[j] = mc[j], mc[i]
}

func (mc MessageCollection) String() string {
	buf := new(bytes.Buffer)
	buf.WriteString(fmt.Sprintf("\n\n%s\n", config.MessageStartMark))
	buf.WriteString("\n//--------- base message ---------\n")
	for _, m := range baseMessageCollection {
		buf.WriteString(fmt.Sprint(m))
	}
	for _, m := range mc {
		buf.WriteString("\n//--------- " + m.Comment + " ---------\n")
		m.GenRpcDefaultMessage(buf)
		m.GenRpcAddReqRespMessage(buf)
		m.GenRpcPutReqRespMessage(buf)
		m.GenRpcDelReqRespMessage(buf)
		m.GenRpcGetReqRespMessage(buf)
		m.GenRpcGetListReqRespMessage(buf)
		m.GenRpcGetEnumsReqRespMessage(buf)
	}
	mc.insertCustomContent(buf)
	return buf.String()
}

func (mc MessageCollection) insertCustomContent(buf *bytes.Buffer) {
	err := common.InsertCustomContent(buf, config.CustomMessageStartMark, config.CustomMessageEndMark, _conf.FileContent, "")
	if err != nil {
		log.Fatal(err)
	}
}
