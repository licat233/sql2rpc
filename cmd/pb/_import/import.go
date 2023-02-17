/*
 * @Author: licat
 * @Date: 2023-02-08 12:31:20
 * @LastEditors: licat
 * @LastEditTime: 2023-02-16 14:29:39
 * @Description: licat233@gmail.com
 */

package _import

import (
	"fmt"
	"strings"
)

type Imp struct {
	Filename string
}

func New(name string) *Imp {
	return &Imp{
		Filename: name,
	}
}

func (i *Imp) String() string {
	name := strings.TrimSpace(i.Filename)
	if name == "" {
		return ""
	}
	return fmt.Sprintf("import \"%s\";\n", name)
}
