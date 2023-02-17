/*
 * @Author: licat
 * @Date: 2023-02-03 19:50:50
 * @LastEditors: licat
 * @LastEditTime: 2023-02-08 13:17:16
 * @Description: licat233@gmail.com
 */
package _enum

import (
	"bytes"
	"fmt"
	"log"

	"github.com/licat233/sql2rpc/cmd/pb/_conf"

	"github.com/licat233/sql2rpc/cmd/common"
	"github.com/licat233/sql2rpc/config"
)

// EnumCollection represents a sortable collection of enums.
type EnumCollection []*Enum

func (ec EnumCollection) Len() int {
	return len(ec)
}

func (ec EnumCollection) Less(i, j int) bool {
	return ec[i].Name < ec[j].Name
}

func (ec EnumCollection) Swap(i, j int) {
	ec[i], ec[j] = ec[j], ec[i]
}

func (ec EnumCollection) String() string {
	var buf = new(bytes.Buffer)
	buf.WriteString(fmt.Sprintf("\n\n%s\n\n", config.EnumStartMark))
	if len(ec) > 0 {
		buf.WriteString("\n//--------- base enum ---------\n")
	}
	for _, e := range ec {
		buf.WriteString(fmt.Sprint(e))
		buf.WriteString("\n")
	}
	ec.insertCustomContent(buf)
	return buf.String()
}

func (ec EnumCollection) insertCustomContent(buf *bytes.Buffer) {
	err := common.InsertCustomContent(buf, config.CustomEnumStartMark, config.CustomEnumEndMark, _conf.FileContent, false, config.C.PbMultiple.GetBool())
	if err != nil {
		log.Fatal(err)
	}
}
