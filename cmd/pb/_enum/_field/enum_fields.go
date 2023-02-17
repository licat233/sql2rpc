/*
 * @Author: licat
 * @Date: 2023-02-06 23:14:15
 * @LastEditors: licat
 * @LastEditTime: 2023-02-06 23:14:20
 * @Description: licat233@gmail.com
 */
package _field

import (
	"bytes"
	"fmt"
)

type EnumFieldCollection []*EnumField

func (efc EnumFieldCollection) String() string {
	buf := new(bytes.Buffer)
	for _, f := range efc {
		buf.WriteString(fmt.Sprint(f))
	}
	return buf.String()
}

func (efc EnumFieldCollection) Copy() EnumFieldCollection {
	var data = make(EnumFieldCollection, len(efc))
	for k, v := range efc {
		cp := *v
		data[k] = &cp
	}
	return data
}
