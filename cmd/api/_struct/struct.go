/*
 * @Author: licat
 * @Date: 2023-02-03 19:25:08
 * @LastEditors: licat
 * @LastEditTime: 2023-02-08 09:36:39
 * @Description: licat233@gmail.com
 */
package _struct

import (
	"bytes"
	"fmt"
	"github.com/licat233/sql2rpc/cmd/api/_conf"
	"github.com/licat233/sql2rpc/cmd/api/_struct/_field"
	"strings"

	"github.com/licat233/sql2rpc/tools"
)

// Struct represents a protocol buffer type.
type Struct struct {
	Name    string
	TagType string
	Comment string
	Fields  _field.StructFieldCollection
}

var listReqFields = _field.StructFieldCollection{
	_field.NewStructField("PageSize", "int64", "json", "pageSize", "optional,default=20", "页面容量，默认20，可选"),
	_field.NewStructField("Page", "int64", "json", "page", "optional,default=20", "页码，默认20，可选"),
	_field.NewStructField("Keyword", "int64", "json", "keyword", "optional", "关键词，可选"),
}

func NewStruct(name, TagType, comment string, fields _field.StructFieldCollection) *Struct {
	return &Struct{
		Name:    name,
		TagType: TagType,
		Comment: comment,
		Fields:  fields,
	}
}

func (s *Struct) Copy() *Struct {
	return &Struct{
		Name:    s.Name,
		TagType: s.TagType,
		Comment: s.Comment,
		Fields:  s.Fields.Copy(),
	}
}

func (s *Struct) handleStructFields(needIgnoreFields []string, more ...string) {
	all := append(needIgnoreFields, more...)
	filterFields := []*_field.StructField{}
	for _, field := range s.Fields {
		if tools.HasInSlice(all, field.Name) {
			continue
		}
		field.TagType = s.TagType
		filterFields = append(filterFields, field)
	}
	s.Fields = filterFields
}

// GenApiDefaultStruct gen default type
func (s *Struct) GenApiDefaultStruct(buf *bytes.Buffer) {
	ss := s.Copy()
	ss.handleStructFields(_conf.IgnoreColumns)
	buf.WriteString(fmt.Sprint(ss))
	ss = nil
}

// GenApiAddReqRespStruct gen req add resp type
func (s *Struct) GenApiAddReqRespStruct(buf *bytes.Buffer) {
	ss := s.Copy()
	//req
	ss.Name = "Add" + tools.ToCamel(s.Name) + "Req"
	ss.Comment = "添加" + s.Comment + "请求"
	ss.handleStructFields(_conf.IgnoreColumns, _conf.MoreIgnoreColumns...)
	buf.WriteString(fmt.Sprint(ss))
	ss = nil
}

// GenApiPutReqRespStruct gen req add resp type
func (s *Struct) GenApiPutReqRespStruct(buf *bytes.Buffer) {
	ss := s.Copy()
	//req
	ss.Name = "Put" + tools.ToCamel(s.Name) + "Req"
	ss.Comment = "更新" + s.Comment + "请求"
	ss.handleStructFields(_conf.IgnoreColumns)
	buf.WriteString(fmt.Sprint(ss))
	ss = nil
}

// GenApiDelReqRespStruct gen req add resp type
func (s *Struct) GenApiDelReqRespStruct(buf *bytes.Buffer) {
	ss := s.Copy()
	//req
	ss.Name = "Del" + tools.ToCamel(s.Name) + "Req"
	ss.Comment = "删除" + s.Comment + "请求"
	ss.Fields = []*_field.StructField{
		_field.NewStructField("Id", "int64", "json", "id", "", s.Comment+" ID"),
	}
	buf.WriteString(fmt.Sprint(ss))
	ss = nil
}

// GenApiGetReqRespStruct gen req add resp type
func (s *Struct) GenApiGetReqRespStruct(buf *bytes.Buffer) {
	ss := s.Copy()
	// req
	ss.Name = "Get" + tools.ToCamel(s.Name) + "Req"
	ss.Comment = "获取" + s.Comment + "请求"
	ss.Fields = []*_field.StructField{
		_field.NewStructField("Id", "int64", "form", "id", "", s.Comment+" ID"),
	}
	buf.WriteString(fmt.Sprint(ss))
	ss = nil
}

// GenApiGetListReqRespStruct gen req add resp type
func (s *Struct) GenApiGetListReqRespStruct(buf *bytes.Buffer) {
	ss := s.Copy()
	//req
	ss.Name = "Get" + tools.ToCamel(s.Name) + "ListReq"
	ss.Comment = "获取" + s.Comment + "列表请求"
	ss.TagType = "form"
	curFields := listReqFields.Copy().PutTagType(ss.TagType)
	for _, field := range ss.Fields {
		if tools.HasInSlice(_conf.IgnoreColumns, field.Name) {
			continue
		}
		field.TagType = ss.TagType
		if !strings.Contains(field.TagOpt, "optional") {
			field.TagOpt += ",optional"
			if field.Typ == "int64" {
				if strings.Contains(field.TagOpt, "Time") {
					field.TagOpt += ",default=0"
				} else {
					field.TagOpt += ",default=-1"
				}
			}
		}
		curFields = append(curFields, field)
	}
	ss.Fields = curFields
	buf.WriteString(fmt.Sprint(ss))
	ss = nil
}

// GenApiGetEnumsReqRespStruct gen req add resp type
func (s *Struct) GenApiGetEnumsReqRespStruct(buf *bytes.Buffer) {
	ss := s.Copy()
	//req
	ss.Name = "Get" + tools.ToCamel(s.Name) + "EnumsReq"
	ss.Comment = "获取" + s.Comment + "枚举请求"
	ss.TagType = "form"
	ss.Fields = []*_field.StructField{
		_field.NewStructField("ParentId", "int64", "form", "parent_id", "optional,default=-1", "父级ID"),
	}
	buf.WriteString(fmt.Sprint(ss))
	ss = nil
}

// String returns a string representation of a Message.
func (s *Struct) String() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("\n//%s\n", s.Comment))
	buf.WriteString(fmt.Sprintf("type %s {\n", tools.ToCamel(s.Name)))
	for _, field := range s.Fields {
		buf.WriteString(fmt.Sprint(field))
	}
	buf.WriteString("}\n")

	return buf.String()
}
