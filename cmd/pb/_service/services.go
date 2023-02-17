/*
 * @Author: licat
 * @Date: 2023-02-06 20:32:11
 * @LastEditors: licat
 * @LastEditTime: 2023-02-16 15:31:28
 * @Description: licat233@gmail.com
 */
package _service

import (
	"bytes"
	"fmt"
	"log"

	"github.com/licat233/sql2rpc/cmd/pb/_conf"

	"github.com/licat233/sql2rpc/cmd/common"
	"github.com/licat233/sql2rpc/config"
	"github.com/licat233/sql2rpc/tools"
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
	if config.C.PbMultiple.GetBool() {
		for _, s := range sc {
			buf.WriteString(fmt.Sprint(s))
		}
		sc.insertCustomContent(buf)
	} else {
		svrName := tools.ToCamel(config.C.ServiceName.GetString())
		buf.WriteString("\n// " + svrName + " _service\n")
		buf.WriteString("service " + svrName + " { \n")
		for _, s := range sc {
			buf.WriteString(fmt.Sprint(s))
		}
		sc.insertCustomContent(buf)
		buf.WriteString("\n}\n")
	}
	return buf.String()
}

func (sc ServiceCollection) insertCustomContent(buf *bytes.Buffer) {
	err := common.InsertCustomContent(buf, config.CustomServiceStartMark, config.CustomServiceEndMark, _conf.FileContent, true, config.C.PbMultiple.GetBool())
	if err != nil {
		log.Fatal(err)
	}
}
