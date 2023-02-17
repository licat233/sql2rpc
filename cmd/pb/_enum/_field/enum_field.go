/*
 * @Author: licat
 * @Date: 2023-02-03 19:32:44
 * @LastEditors: licat
 * @LastEditTime: 2023-02-17 22:29:37
 * @Description: licat233@gmail.com
 */
package _field

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/licat233/sql2rpc/cmd/common"
)

// EnumField represents a field in an enumerated type.
type EnumField struct {
	name string
	tag  int
}

// NewEnumField constructs an EnumField type.
func NewEnumField(name string, tag int) *EnumField {
	name = strings.ToUpper(name)

	re := regexp.MustCompile(`([^\w]+)`)
	name = re.ReplaceAllString(name, "_")

	return &EnumField{name, tag}
}

// String returns a string representation of an Enum.
func (ef EnumField) String() string {
	return fmt.Sprintf("%s%s = %d;\n", common.Indent, ef.name, ef.tag)
}

// Name returns the name of the _enum field.
func (ef EnumField) Name() string {
	return ef.name
}

// Tag returns the identifier tag of the _enum field.
func (ef EnumField) Tag() int {
	return ef.tag
}
