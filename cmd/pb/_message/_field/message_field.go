/*
 * @Author: licat
 * @Date: 2023-02-03 19:26:33
 * @LastEditors: licat
 * @LastEditTime: 2023-02-08 11:28:08
 * @Description: licat233@gmail.com
 */
package _field

import (
	"fmt"

	"github.com/licat233/sql2rpc/tools"
)

// MessageField represents the field of a _message.
type MessageField struct {
	Typ     string
	Name    string
	Tag     int
	Comment string
}

// NewMessageField creates a new _message field.
func NewMessageField(typ, name string, tag int, comment string) *MessageField {
	return &MessageField{
		Typ:     typ,
		Name:    name,
		Tag:     tag,
		Comment: comment,
	}
}

// String returns a string representation of a _message field.
func (f MessageField) String() string {
	// matched, _ := regexp.MatchString("^repeated\\s[A-Z].*", f.Typ)
	return fmt.Sprintf("  %s %s = %d; //%s\n", f.Typ, tools.ToSnake(f.Name), f.Tag, f.Comment)
}
