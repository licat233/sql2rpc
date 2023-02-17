/*
 * @Author: licat
 * @Date: 2023-02-06 00:34:15
 * @LastEditors: licat
 * @LastEditTime: 2023-02-17 09:41:31
 * @Description: licat233@gmail.com
 */

package _service

import (
	"bytes"
	"fmt"
	"log"

	"github.com/licat233/sql2rpc/cmd/api/_conf"

	"github.com/licat233/sql2rpc/cmd/common"
	"github.com/licat233/sql2rpc/config"
)

type ServiceCollection []*Service

func (sc ServiceCollection) Len() int {
	return len(sc)
}

func (sc ServiceCollection) Less(i, j int) bool {
	return sc[i].Name < sc[j].Name
}

func (sc ServiceCollection) Swap(i, j int) {
	sc[i], sc[j] = sc[j], sc[i]
}

func (sc ServiceCollection) String() string {
	buf := new(bytes.Buffer)
	buf.WriteString(fmt.Sprintf("\n%s\n", config.ServiceStartMark))
	for _, s := range sc {
		buf.WriteString(fmt.Sprint(s))
	}
	buf.WriteString("\n")
	sc.insertCustomContent(buf)
	return buf.String()
}

func (sc ServiceCollection) insertCustomContent(buf *bytes.Buffer) {
	err := common.InsertCustomContent(buf, config.CustomServiceStartMark, config.CustomServiceEndMark, _conf.FileContent, false, false)
	if err != nil {
		log.Fatal(err)
	}
}
