/*
 * @Author: licat
 * @Date: 2023-02-03 19:26:33
 * @LastEditors: licat
 * @LastEditTime: 2023-02-17 22:36:33
 * @Description: licat233@gmail.com
 */

package _field

import (
	"fmt"
	"strings"

	"github.com/licat233/sql2rpc/cmd/common"
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

// NewStructField creates a new type field.
func NewStructField(name, typ, tagType, tagName, tagOpt, comment string) *StructField {
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
	opt := strings.TrimSpace(f.TagOpt)
	if opt != "" {
		if opt[:1] != "," {
			opt = "," + opt
		}
	}
	tName := common.ConvertStringStyle(f.TagName)
	//name 必须camel风格
	return fmt.Sprintf("%s%s %s `%s:\"%s%s\"` //%s\n", common.Indent, tools.ToCamel(f.Name), f.Typ, f.TagType, tName, opt, f.Comment)
}
