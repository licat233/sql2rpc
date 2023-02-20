/*
 * @Author: licat
 * @Date: 2023-02-06 20:32:11
 * @LastEditors: licat
 * @LastEditTime: 2023-02-20 14:36:43
 * @Description: licat233@gmail.com
 */
package _service

import (
	"bytes"
	"fmt"
	"log"
	"regexp"
	"strings"

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
		// buf.WriteString(sc.getCustomContent())
		sc.insertCustomContent(buf, "")
	} else {
		buf.WriteString("\n// " + svrName + " service")
		buf.WriteString("\nservice " + svrName + " {")
		for _, s := range sc {
			buf.WriteString(fmt.Sprint(s))
		}
		// buf.WriteString(sc.getCustomContent())
		sc.insertCustomContent(buf, common.Indent)
		buf.WriteString("\n}\n")
	}
	return buf.String()
}

// getCustomContent 生成自定义的内容版块
// 注意：使用该功能，可能会让用户在大意间，损失自己定义好的service，所以暂时不推荐使用
func (sc ServiceCollection) GetCustomContent() string {
	var custombuf = new(bytes.Buffer)
	sc.insertCustomContent(custombuf, "")
	customContent := custombuf.String()
	pattern := `(?s)service\s*\w+\s*\{(.*?)}`
	re := regexp.MustCompile(pattern)
	match := re.FindAllStringSubmatch(customContent, -1)
	startMark := config.CustomServiceStartMark
	endMark := config.CustomServiceEndMark
	if len(match) == 0 {
		if config.C.PbMultiple.GetBool() {
			svrName := tools.ToCamel(config.C.ServiceName.GetString())
			buf := new(bytes.Buffer)
			buf.WriteString("\n// " + svrName + " service")
			buf.WriteString("\nservice " + svrName + " {")
			buf.WriteString(common.PickInfoContent(customContent))
			buf.WriteString("\n}\n")
			content := buf.String()
			return common.GenCustomBlock(startMark, endMark, content, "")
		} else {
			return customContent
		}
	} else {
		if config.C.PbMultiple.GetBool() {
			return customContent
		} else {
			rpcs := []string{}
			for _, rpc := range match {
				if len(rpc) == 2 {
					rpcs = append(rpcs, rpc[1])
				}
			}
			content := strings.Join(rpcs, "\n")
			return common.GenCustomBlock(startMark, endMark, content, common.Indent)
		}
	}
}

func (sc ServiceCollection) insertCustomContent(buf *bytes.Buffer, indent string) {
	err := common.InsertCustomContent(buf, config.CustomServiceStartMark, config.CustomServiceEndMark, _conf.FileContent, indent)
	if err != nil {
		log.Fatal(err)
	}
}

// func (sc ServiceCollection) existServiceBlock(content string) bool {
// 	pattern := `(?s)service\s*\w+\s*\{(.*?)}`
// 	matched, err := regexp.MatchString(pattern, content)
// 	if err != nil {
// 		log.Fatalf("error: %s", err)
// 	}
// 	return matched
// }
