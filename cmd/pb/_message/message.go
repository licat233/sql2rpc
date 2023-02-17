/*
 * @Author: licat
 * @Date: 2023-02-03 19:25:08
 * @LastEditors: licat
 * @LastEditTime: 2023-02-08 13:17:36
 * @Description: licat233@gmail.com
 */
package _message

import (
	"bytes"
	"fmt"

	"github.com/licat233/sql2rpc/cmd/pb/_conf"
	"github.com/licat233/sql2rpc/cmd/pb/_message/_field"

	"github.com/licat233/sql2rpc/tools"
)

// Message represents a protocol buffer _message.
type Message struct {
	Name    string
	Comment string
	Fields  _field.MessageFieldCollection
}

func NewMessage(name, comment string, fields _field.MessageFieldCollection) *Message {
	return &Message{
		Name:    name,
		Comment: comment,
		Fields:  fields,
	}
}

func (ms *Message) Copy() *Message {
	return &Message{
		Name:    ms.Name,
		Comment: ms.Comment,
		Fields:  ms.Fields.Copy(),
	}
}

func (ms *Message) handlerMessageFields(needIgnoreFields []string, more ...string) {
	all := append(needIgnoreFields, more...)
	curFields := []*_field.MessageField{}
	var filedTag int
	for _, field := range ms.Fields {
		if tools.HasInSlice(all, field.Name) {
			continue
		}
		filedTag++
		field.Tag = filedTag
		curFields = append(curFields, field)
	}
	ms.Fields = curFields
}

// GenRpcDefaultMessage gen default _message
func (ms *Message) GenRpcDefaultMessage(buf *bytes.Buffer) {
	m := ms.Copy()
	m.handlerMessageFields(_conf.IgnoreColumns)
	buf.WriteString(fmt.Sprintf("%s\n", m))
	m = nil
}

// GenRpcAddReqRespMessage gen req add resp _message
func (ms *Message) GenRpcAddReqRespMessage(buf *bytes.Buffer) {
	//req
	req := ms.Copy()
	req.Name = "Add" + tools.ToCamel(ms.Name) + "Req"
	req.Comment = "添加" + ms.Comment + "请求"
	req.handlerMessageFields(_conf.IgnoreColumns, _conf.MoreIgnoreColumns...)
	buf.WriteString(fmt.Sprintf("%s\n", req))
	req = nil

	//resp
	resp := ms.Copy()
	resp.Name = "Add" + tools.ToCamel(ms.Name) + "Resp"
	resp.Comment = "添加" + ms.Comment + "响应"
	resp.Fields = []*_field.MessageField{
		_field.NewMessageField(tools.ToCamel(ms.Name), ms.Name, 1, ms.Comment+"信息"),
	}
	buf.WriteString(fmt.Sprintf("%s\n", resp))
	resp = nil
}

// GenRpcPutReqRespMessage gen req add resp _message
func (ms *Message) GenRpcPutReqRespMessage(buf *bytes.Buffer) {
	//req
	req := ms.Copy()
	req.Name = "Put" + tools.ToCamel(ms.Name) + "Req"
	req.Comment = "更新" + ms.Comment + "请求"
	req.handlerMessageFields(_conf.IgnoreColumns)
	buf.WriteString(fmt.Sprintf("%s\n", req))
	req = nil

	//resp
	resp := ms.Copy()
	resp.Name = "Put" + tools.ToCamel(ms.Name) + "Resp"
	resp.Comment = "更新" + ms.Comment + "响应"
	resp.Fields = []*_field.MessageField{}
	buf.WriteString(fmt.Sprintf("%s\n", resp))
	resp = nil
}

// GenRpcDelReqRespMessage gen req add resp _message
func (ms *Message) GenRpcDelReqRespMessage(buf *bytes.Buffer) {
	//req
	req := ms.Copy()
	req.Name = "Del" + tools.ToCamel(ms.Name) + "Req"
	req.Comment = "删除" + ms.Comment + "请求"
	req.Fields = []*_field.MessageField{
		_field.NewMessageField("int64", "id", 1, ms.Comment+" ID"),
	}
	buf.WriteString(fmt.Sprintf("%s\n", req))
	req = nil

	//resp
	resp := ms.Copy()
	resp.Name = "Del" + tools.ToCamel(ms.Name) + "Resp"
	resp.Comment = "删除" + ms.Comment + "响应"
	resp.Fields = []*_field.MessageField{}
	buf.WriteString(fmt.Sprintf("%s\n", resp))
	resp = nil
}

// GenRpcGetReqRespMessage gen req add resp _message
func (ms *Message) GenRpcGetReqRespMessage(buf *bytes.Buffer) {
	//req
	req := ms.Copy()
	req.Name = "Get" + tools.ToCamel(ms.Name) + "Req"
	req.Comment = "获取" + ms.Comment + "请求"
	req.Fields = []*_field.MessageField{
		_field.NewMessageField("int64", "id", 1, ms.Comment+" ID"),
	}
	buf.WriteString(fmt.Sprintf("%s\n", req))
	req = nil

	//resp
	resp := ms.Copy()
	resp.Name = "Get" + tools.ToCamel(ms.Name) + "Resp"
	resp.Comment = "获取" + ms.Comment + "响应"
	resp.Fields = []*_field.MessageField{
		_field.NewMessageField(tools.ToCamel(ms.Name), ms.Name, 1, ms.Comment+" 信息"),
	}
	buf.WriteString(fmt.Sprintf("%s\n", resp))
	resp = nil
}

// GenRpcGetListReqRespMessage gen req add resp _message
func (ms *Message) GenRpcGetListReqRespMessage(buf *bytes.Buffer) {
	//req
	req := ms.Copy()
	req.Name = "Get" + tools.ToCamel(ms.Name) + "ListReq"
	req.Comment = "获取" + ms.Comment + "列表请求"
	req.Fields = []*_field.MessageField{
		_field.NewMessageField("ListReq", "list_req", 1, "列表页码参数"),
		_field.NewMessageField(tools.ToCamel(ms.Name), ms.Name, 2, ms.Comment+"参数"),
	}
	buf.WriteString(fmt.Sprintf("%s\n", req))
	req = nil

	//resp
	resp := ms.Copy()
	resp.Name = "Get" + tools.ToCamel(ms.Name) + "ListResp"
	resp.Comment = "获取" + ms.Comment + "列表响应"
	resp.Fields = []*_field.MessageField{
		_field.NewMessageField("repeated "+tools.ToCamel(ms.Name), tools.PluralizedName(ms.Name), 1, ms.Comment+"列表"),
		_field.NewMessageField("int64", "total", 2, ms.Comment+"总数量"),
	}
	buf.WriteString(fmt.Sprintf("%s\n", resp))
	resp = nil
}

// GenRpcGetEnumsReqRespMessage gen req add resp _message
func (ms *Message) GenRpcGetEnumsReqRespMessage(buf *bytes.Buffer) {
	//req
	req := ms.Copy()
	req.Name = "Get" + tools.ToCamel(ms.Name) + "EnumsReq"
	req.Comment = "获取" + ms.Comment + "列表请求"
	req.Fields = []*_field.MessageField{
		_field.NewMessageField("int64", "parent_id", 1, "父级ID"),
	}
	buf.WriteString(fmt.Sprintf("%s\n", req))
	req = nil
}

// String returns a string representation of a Message.
func (ms *Message) String() string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("\n//%s\nmessage %s {\n", ms.Comment, tools.ToCamel(ms.Name)))
	buf.WriteString(fmt.Sprint(ms.Fields))
	buf.WriteString("}\n")

	return buf.String()
}

// AppendField appends a _message field to a _message. If the tag of the _message field is in use, an error will be returned.
func (ms *Message) AppendField(mf *_field.MessageField) error {
	for _, f := range ms.Fields {
		if f.Tag == mf.Tag {
			return fmt.Errorf("tag `%d` is already in use by field `%s`", mf.Tag, f.Name)
		}
	}

	ms.Fields = append(ms.Fields, mf)
	return nil
}
