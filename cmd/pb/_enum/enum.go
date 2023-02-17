/*
 * @Author: licat
 * @Date: 2023-02-03 19:31:55
 * @LastEditors: licat
 * @LastEditTime: 2023-02-06 23:17:59
 * @Description: licat233@gmail.com
 */
package _enum

import (
	"bytes"
	"fmt"
	"github.com/licat233/sql2rpc/cmd/pb/_enum/_field"
)

// Enum represents a protocol buffer enumerated type.
type Enum struct {
	Name    string
	Comment string
	Fields  _field.EnumFieldCollection
}

// NewEnumFromStrings creates an _enum from a name and a slice of strings that represent the names of each field.
func NewEnumFromStrings(name, comment string, ss []string) (*Enum, error) {
	enum := &Enum{}
	enum.Name = name
	enum.Comment = comment

	for i, s := range ss {
		err := enum.AppendField(_field.NewEnumField(s, i))
		if nil != err {
			return nil, err
		}
	}

	return enum, nil
}

// String returns a string representation of an Enum.
func (e *Enum) String() string {
	buf := new(bytes.Buffer)

	buf.WriteString(fmt.Sprintf("// %s \n", e.Comment))
	buf.WriteString(fmt.Sprintf("_enum %s {\n", e.Name))
	buf.WriteString(fmt.Sprint(e.Fields))
	buf.WriteString("}\n")

	return buf.String()
}

// AppendField appends an EnumField to an Enum.
func (e *Enum) AppendField(ef *_field.EnumField) error {
	for _, f := range e.Fields {
		if f.Tag() == ef.Tag() {
			return fmt.Errorf("tag `%d` is already in use by field `%s`", ef.Tag(), f.Name())
		}
	}

	e.Fields = append(e.Fields, ef)

	return nil
}
