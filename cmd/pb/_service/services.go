/*
 * @Author: licat
 * @Date: 2023-02-06 20:32:11
 * @LastEditors: licat
 * @LastEditTime: 2023-02-17 22:56:39
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
	svrName := tools.ToCamel(config.C.ServiceName.GetString())
	buf := new(bytes.Buffer)
	buf.WriteString(fmt.Sprintf("\n%s", config.ServiceStartMark))
	if config.C.PbMultiple.GetBool() {
		buf.WriteString("\nservice Base" + svrName + " {")
		for _, s := range sc {
			buf.WriteString(fmt.Sprint(s))
		}
		buf.WriteString("\n}\n")
		sc.insertCustomContent(buf, "")
	} else {
		buf.WriteString("\n// " + svrName + " service")
		buf.WriteString("\nservice " + svrName + " {")
		for _, s := range sc {
			buf.WriteString(fmt.Sprint(s))
		}
		sc.insertCustomContent(buf, common.Indent)
		buf.WriteString("\n}\n")
	}
	return buf.String()
}

func (sc ServiceCollection) insertCustomContent(buf *bytes.Buffer, indent string) {
	err := common.InsertCustomContent(buf, config.CustomServiceStartMark, config.CustomServiceEndMark, _conf.FileContent, indent)
	if err != nil {
		log.Fatal(err)
	}
}
