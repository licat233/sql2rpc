/*
 * @Author: licat
 * @Date: 2023-02-03 19:51:18
 * @LastEditors: licat
 * @LastEditTime: 2023-02-17 09:42:20
 * @Description: licat233@gmail.com
 */

package _struct

import (
	"bytes"
	"fmt"
	"log"

	"github.com/licat233/sql2rpc/cmd/api/_conf"
	"github.com/licat233/sql2rpc/cmd/api/_struct/_field"

	"github.com/licat233/sql2rpc/cmd/common"
	"github.com/licat233/sql2rpc/config"
)

// StructCollection represents a sortable collection of messages.
type StructCollection []*Struct

var baseStructCollection StructCollection = []*Struct{
	NewStruct("Enum", "json", "枚举", _field.StructFieldCollection{
		_field.NewStructField("Label", "interface{}", "json", "label", "", "名"),
		_field.NewStructField("Value", "interface{}", "json", "label", "", "值"),
	}),
	NewStruct("Enums", "json", "枚举列表", _field.StructFieldCollection{
		_field.NewStructField("List", "[]Enum", "json", "list", "", "枚举列表数据"),
	}),
	NewStruct("Option", "json", "选项", _field.StructFieldCollection{
		_field.NewStructField("Title", "string", "json", "title", "", "名"),
		_field.NewStructField("Value", "int64", "json", "value", "", "值"),
	}),
	NewStruct("Options", "json", "选项列表", _field.StructFieldCollection{
		_field.NewStructField("List", "[]Option", "json", "list", "", "选项列表数据"),
	}),
	NewStruct("TreeOption", "json", "树形选项", _field.StructFieldCollection{
		_field.NewStructField("Title", "string", "json", "title", "", "名"),
		_field.NewStructField("Value", "int64", "json", "value", "", "值"),
		_field.NewStructField("Children", "[]TreeOption", "json", "children", "optional", "子集"),
	}),
	NewStruct("TreeOptions", "json", "树形选项列表", _field.StructFieldCollection{
		_field.NewStructField("List", "[]TreeOption", "json", "list", "", "树形选项列表数据"),
	}),
	NewStruct("ListReq", "form", "列表数据请求", listReqFields.Copy().PutTagType("form")),
	NewStruct("ByIdReq", "form", "通过ID请求", _field.StructFieldCollection{
		_field.NewStructField("Id", "int64", "form", "id", "", "主键"),
	}),
	NewStruct("NilReq", "json", "空请求", nil),
	NewStruct("NilResp", "json", "空响应", nil),
	NewStruct("BaseResp", "json", "规范响应体", _field.StructFieldCollection{
		_field.NewStructField("Status", "bool", "json", "success", "", "响应状态"),
		_field.NewStructField("Message", "string", "json", "_message", "optional", "给予的提示信息"),
		_field.NewStructField("Data", "interface{}", "json", "data", "optional", "【选填】响应的业务数据"),
		_field.NewStructField("Total", "int64", "json", "total", "optional", "【选填】数据总个数"),
		_field.NewStructField("PageSize", "int64", "json", "pageSize", "optional", "【选填】单页数量"),
		_field.NewStructField("Page", "int64", "json", "current", "optional", "【选填】当前页码，current与antd前端对接"),
		_field.NewStructField("TotalPage", "int64", "json", "totalPage", "optional", "【选填】自增项，总共有多少页，根据前端的pageSize来计算"),
		_field.NewStructField("ErrorCode", "int64", "json", "errorCode", "optional", "【选填】错误类型代码：400错误请求，401未授权，500服务器内部错误，200成功"),
		_field.NewStructField("ErrorMessage", "string", "json", "errorMessage", "optional", "【选填】向用户显示消息"),
		_field.NewStructField("TraceMessage", "string", "json", "traceMessage", "optional", "【选填】调试错误信息，请勿在生产环境下使用，可有可无"),
		_field.NewStructField("ShowStruct", "int64", "json", "showStruct", "optional", "【选填】错误显示类型：0.不提示错误;1.警告信息提示；2.错误信息提示；4.通知提示；9.页面跳转"),
		_field.NewStructField("TraceId", "string", "json", "traceId", "optional", "【选填】方便后端故障排除：唯一的请求ID"),
		_field.NewStructField("Host", "string", "json", "host", "optional", "【选填】方便后端故障排除：当前访问服务器的主机"),
	}),
}

func (sc StructCollection) Len() int {
	return len(sc)
}

func (sc StructCollection) Less(i, j int) bool {
	return sc[i].Name < sc[j].Name
}

func (sc StructCollection) Swap(i, j int) {
	sc[i], sc[j] = sc[j], sc[i]
}

func (sc StructCollection) Base() string {
	buf := new(bytes.Buffer)
	buf.WriteString("\n//--------- base struct ---------\n")
	for _, m := range baseStructCollection {
		buf.WriteString(fmt.Sprint(m))
	}
	buf.WriteString("\n")
	return buf.String()
}

func (sc StructCollection) String() string {
	buf := new(bytes.Buffer)
	buf.WriteString(fmt.Sprintf("\n\n%s\n", config.StructStartMark))

	if _conf.CurrentIsCoreFile {
		buf.WriteString(sc.Base())
	}

	for _, s := range sc {
		buf.WriteString("\n//---------" + s.Comment + "---------\n")
		s.GenApiDefaultStruct(buf)
		s.GenApiAddReqRespStruct(buf)
		s.GenApiPutReqRespStruct(buf)
		s.GenApiDelReqRespStruct(buf)
		s.GenApiGetReqRespStruct(buf)
		s.GenApiGetListReqRespStruct(buf)
		s.GenApiGetEnumsReqRespStruct(buf)
	}
	sc.insertCustomContent(buf)
	return buf.String()
}

func (sc StructCollection) insertCustomContent(buf *bytes.Buffer) {
	err := common.InsertCustomContent(buf, config.CustomStructStartMark, config.CustomStructEndMark, _conf.FileContent, false, false)
	if err != nil {
		log.Fatal(err)
	}
}
