/*
 * @Author: licat
 * @Date: 2023-02-03 19:51:18
 * @LastEditors: licat
 * @LastEditTime: 2023-02-22 17:34:14
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

var listReqFields = _field.StructFieldCollection{
	_field.New("PageSize", "int64", "json", "pageSize", "optional,default=20", "页面容量，默认20，可选"),
	_field.New("Page", "int64", "json", "page", "optional,default=1", "当前页码，默认1，可选"),
	_field.New("Current", "int64", "json", "current", "optional,default=1", "当前页码，默认1，用于对接umijs，可选"),
	_field.New("Keyword", "string", "json", "keyword", "optional", "关键词，可选"),
}

var baseStructCollection StructCollection = []*Struct{
	New("Enum", "json", "枚举", _field.StructFieldCollection{
		_field.New("Label", "interface{}", "json", "label", "", "名"),
		_field.New("Value", "interface{}", "json", "value", "", "值"),
	}),
	New("Enums", "json", "枚举列表", _field.StructFieldCollection{
		_field.New("List", "[]Enum", "json", "list", "", "枚举列表数据"),
	}),
	New("Option", "json", "选项", _field.StructFieldCollection{
		_field.New("Title", "string", "json", "title", "", "名"),
		_field.New("Value", "int64", "json", "value", "", "值"),
	}),
	New("Options", "json", "选项列表", _field.StructFieldCollection{
		_field.New("List", "[]Option", "json", "list", "", "选项列表数据"),
	}),
	New("TreeOption", "json", "树形选项", _field.StructFieldCollection{
		_field.New("Title", "string", "json", "title", "", "名"),
		_field.New("Value", "int64", "json", "value", "", "值"),
		_field.New("Children", "[]TreeOption", "json", "children", "optional", "子集"),
	}),
	New("TreeOptions", "json", "树形选项列表", _field.StructFieldCollection{
		_field.New("List", "[]TreeOption", "json", "list", "", "树形选项列表数据"),
	}),
	New("JwtToken", "json", "jwt token", _field.StructFieldCollection{
		_field.New("AccessToken", "string", "json", "accessToken", "", "token"),
		_field.New("AccessExpire", "int64", "json", "accessExpire", "", "expire"),
		_field.New("RefreshAfter", "int64", "json", "refreshAfter", "", "refresh at"),
	}),
	New("ListReq", "form", "列表数据请求", listReqFields.Copy().PutTagType("form")),
	New("ByIdReq", "form", "通过ID请求", _field.StructFieldCollection{
		_field.New("Id", "int64", "form", "id", "", "主键"),
	}),
	New("NilReq", "json", "空请求", nil),
	New("NilResp", "json", "空响应", nil),
	New("Resp", "json", "空响应", _field.StructFieldCollection{
		_field.New("Body", "interface{}", "json", "body", "", "响应数据"),
	}),
	New("CaptchaResp", "json", "验证码响应", _field.StructFieldCollection{
		_field.New("CaptchaId", "string", "json", "captchaId", "", "captcha id"),
		_field.New("ExpiresAt", "int64", "json", "expiresAt", "", "expires time"),
	}),
	New("BaseResp", "json", "规范响应体", _field.StructFieldCollection{
		_field.New("Status", "bool", "json", "status", "", "响应状态"),
		_field.New("Success", "bool", "json", "success", "", "响应状态，用于对接umijs"),
		_field.New("Message", "string", "json", "message", "optional,omitempty", "给予的提示信息"),
		_field.New("Data", "interface{}", "json", "data", "optional,omitempty", "【选填】响应的业务数据"),
		_field.New("Total", "int64", "json", "total", "optional,omitempty", "【选填】数据总个数"),
		_field.New("PageSize", "int64", "json", "pageSize", "optional,omitempty", "【选填】单页数量"),
		_field.New("Current", "int64", "json", "current", "optional,omitempty", "【选填】当前页码，用于对接umijs"),
		_field.New("Page", "int64", "json", "page", "optional,omitempty", "【选填】当前页码"),
		_field.New("TotalPage", "int64", "json", "totalPage", "optional,omitempty", "【选填】自增项，总共有多少页，根据前端的pageSize来计算"),
		_field.New("ErrorCode", "int64", "json", "errorCode", "optional,omitempty", "【选填】错误类型代码：400错误请求，401未授权，500服务器内部错误，200成功"),
		_field.New("ErrorMessage", "string", "json", "errorMessage", "optional,omitempty", "【选填】向用户显示消息"),
		_field.New("TraceMessage", "string", "json", "traceMessage", "optional,omitempty", "【选填】调试错误信息，请勿在生产环境下使用，可有可无"),
		_field.New("ShowType", "int64", "json", "showType", "optional,omitempty", "【选填】错误显示类型：0.不提示错误;1.警告信息提示；2.错误信息提示；4.通知提示；9.页面跳转"),
		_field.New("TraceId", "string", "json", "traceId", "optional,omitempty", "【选填】方便后端故障排除：唯一的请求ID"),
		_field.New("Host", "string", "json", "host", "optional,omitempty", "【选填】方便后端故障排除：当前访问服务器的主机"),
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
	err := common.InsertCustomContent(buf, config.CustomStructStartMark, config.CustomStructEndMark, _conf.FileContent, "")
	if err != nil {
		log.Fatal(err)
	}
}
