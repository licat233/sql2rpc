/*
 * @Author: licat
 * @Date: 2023-02-06 23:23:16
 * @LastEditors: licat
 * @LastEditTime: 2023-02-07 21:38:51
 * @Description: licat233@gmail.com
 */

package _field

import (
	"bytes"
	"fmt"
)

type StructFieldCollection []*StructField

func (sfc StructFieldCollection) String() string {
	buf := new(bytes.Buffer)
	for _, f := range sfc {
		buf.WriteString(fmt.Sprint(f))
	}
	return buf.String()
}

func (sfc StructFieldCollection) Copy() StructFieldCollection {
	var data = make(StructFieldCollection, len(sfc))
	for k, v := range sfc {
		cp := *v
		data[k] = &cp
	}
	return data
}

func (sfc StructFieldCollection) PutTagType(tagType string) StructFieldCollection {
	for k := range sfc {
		sfc[k].TagType = tagType
	}
	return sfc
}
