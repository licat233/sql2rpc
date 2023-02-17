/*
 * @Author: licat
 * @Date: 2023-02-09 09:42:29
 * @LastEditors: licat
 * @LastEditTime: 2023-02-09 09:46:02
 * @Description: licat233@gmail.com
 */

package _item

import (
	"bytes"
	"fmt"
)

type FieldCollection []*Field

func (fc FieldCollection) String() string {
	buf := new(bytes.Buffer)
	for _, v := range fc {
		buf.WriteString(fmt.Sprint(v))
	}
	return buf.String()
}
