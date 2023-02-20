/*
 * @Author: licat
 * @Date: 2023-02-03 19:26:33
 * @LastEditors: licat
 * @LastEditTime: 2023-02-20 12:00:04
 * @Description: licat233@gmail.com
 */

package _field

import (
	"fmt"
	"strings"

	"github.com/licat233/sql2rpc/cmd/common"
	"github.com/licat233/sql2rpc/config"
	"github.com/licat233/sql2rpc/tools"
)

// StructField represents the field of a type.
type StructField struct {
	Name    string
	Typ     string
	TagType string
	TagName string
	TagOpt  string
	Comment string
}

// New creates a new type field.
func New(name, typ, tagType, tagName, tagOpt, comment string) *StructField {
	if tagName == "" {
		tagName = name
	}
	return &StructField{
		Name:    name,
		Typ:     typ,
		TagType: tagType,
		TagName: tagName,
		TagOpt:  tagOpt,
		Comment: comment,
	}
}

// String returns a string representation of a type field.
func (f StructField) String() string {
	if f.TagName == "" {
		f.TagName = f.Name
	}
	tName := common.ConvertStringStyle(config.C.ApiStyle.GetString(), f.TagName)
	optString := handleOptContent(tName, f.TagOpt)
	comment := strings.TrimSpace(f.Comment)
	if comment == "" {
		comment = ""
	} else {
		comment = " // " + comment
	}
	//name 必须camel风格
	return fmt.Sprintf("%s%s %s `%s:\"%s\"` %s\n", common.Indent, tools.ToCamel(f.Name), f.Typ, f.TagType, optString, comment)
}

func handleOptContent(name string, opts ...string) string {
	name = strings.Trim(name, ",")
	str := strings.Join(opts, ",")
	list := strings.Split(str, ",")
	filter := []string{name}
	for _, arg := range list {
		arg = strings.TrimSpace(arg)
		if arg == "" {
			continue
		}
		filter = append(filter, arg)
	}
	return strings.Join(filter, ",")
}
