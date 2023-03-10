/*
 * @Author: licat
 * @Date: 2023-02-08 12:36:15
 * @LastEditors: licat
 * @LastEditTime: 2023-02-17 22:40:21
 * @Description: licat233@gmail.com
 */

package _import

import (
	"bytes"
	"fmt"
	"log"

	"github.com/licat233/sql2rpc/cmd/common"
	"github.com/licat233/sql2rpc/config"

	"github.com/licat233/sql2rpc/cmd/pb/_conf"
)

type ImpCollection []*Imp

func (ic ImpCollection) Len() int {
	return len(ic)
}

func (ic ImpCollection) Less(i, j int) bool {
	return ic[i].Filename < ic[j].Filename
}

func (ic ImpCollection) Swap(i, j int) {
	ic[i], ic[j] = ic[j], ic[i]
}

func (ic ImpCollection) Search(x string) int {
	for k, v := range ic {
		if v.Filename == x {
			return k
		}
	}
	return -1
}

func (ic ImpCollection) String() string {
	var buf = new(bytes.Buffer)
	buf.WriteString(fmt.Sprintf("\n\n%s\n\n", config.ImportStartMark))
	if len(ic) > 0 {
		buf.WriteString("\n//--------- base import ---------\n")
	}
	for _, i := range ic {
		buf.WriteString(fmt.Sprint(i))
	}
	ic.insertCustomContent(buf)
	return buf.String()
}

func (ic ImpCollection) insertCustomContent(buf *bytes.Buffer) {
	err := common.InsertCustomContent(buf, config.CustomImportStartMark, config.CustomImportEndMark, _conf.FileContent, "")
	if err != nil {
		log.Fatal(err)
	}
}
