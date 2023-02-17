/*
 * @Author: licat
 * @Date: 2023-02-06 23:12:48
 * @LastEditors: licat
 * @LastEditTime: 2023-02-08 10:17:58
 * @Description: licat233@gmail.com
 */
package _field

import (
	"bytes"
	"fmt"

	"github.com/licat233/sql2rpc/tools"
)

type MessageFieldCollection []*MessageField

func (mfc MessageFieldCollection) String() string {
	buf := new(bytes.Buffer)
	for _, f := range mfc {
		buf.WriteString(fmt.Sprint(f))
	}
	return buf.String()
}

func (mfc MessageFieldCollection) Copy() MessageFieldCollection {
	var data = make(MessageFieldCollection, len(mfc))
	for k, v := range mfc {
		cp := *v
		data[k] = &cp
	}
	return data
}

func ListRespFields(dataType string) MessageFieldCollection {
	name := tools.PluralizedName(dataType)
	return MessageFieldCollection{
		NewMessageField("repeated "+dataType, name, 1, "数据列表"),
		NewMessageField("int64", "total", 2, "总数量"),
	}
}
