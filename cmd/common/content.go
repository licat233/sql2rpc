/*
 * @Author: licat
 * @Date: 2023-02-05 17:14:15
 * @LastEditors: licat
 * @LastEditTime: 2023-02-17 12:56:13
 * @Description: licat233@gmail.com
 */
package common

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

func PickMarkContents(startMark, endMark, oldContent string) ([]string, error) {
	content := strings.TrimSpace(oldContent)
	if content == "" {
		return []string{}, nil
	}

	startMark = regexp.QuoteMeta(startMark)
	endMark = regexp.QuoteMeta(endMark)

	expr := fmt.Sprintf("%s[\n]*((?s).*?)[\n]*%s", startMark, endMark)

	reg, err := regexp.Compile(expr)
	if err != nil {
		return nil, err
	}
	listArr := reg.FindAllStringSubmatch(content, -1)
	strArr := []string{}
	for _, list := range listArr {
		if len(list) != 2 {
			continue
		}
		target := strings.TrimSpace(list[len(list)-1])
		if target == "" {
			continue
		}
		target = strings.Trim(target, "")
		target = strings.Trim(target, "\n")
		target = fmt.Sprintf("\n%s\n", target)
		strArr = append(strArr, target)
	}
	return strArr, nil
}

func InsertCustomContent(buf *bytes.Buffer, startMark, endMark, oldContent string, isService, multiple bool) error {
	prefix := ""
	if isService && !multiple {
		prefix = "  "
	}

	buf.WriteString(fmt.Sprintf("\n\n%s// The content in this block will not be updated\n%s// 此区块内的内容不会被更新", prefix, prefix))

	list, err := PickMarkContents(startMark, endMark, oldContent)
	if err != nil {
		return err
	}

	customContent := strings.Join(list, "\n")
	customContent = strings.Trim(customContent, "\n")

	customContent = fmt.Sprintf("\n%s%s\n", prefix, customContent)
	startMark = fmt.Sprintf("\n%s%s\n", prefix, startMark)
	endMark = fmt.Sprintf("\n%s%s\n", prefix, endMark)

	buf.WriteString(startMark)
	buf.WriteString(customContent)
	buf.WriteString(endMark)
	return nil
}

func PickInfoContent(content string) string {
	re := regexp.MustCompile(`(?s)info\s*\((.*?)\n\)`)
	match := re.FindStringSubmatch(content)

	if len(match) == 2 {
		res := match[1]
		str := strings.TrimSpace(res)
		str = strings.Trim(str, "\n")
		if str == "" {
			return ""
		}
		return res
	} else {
		return ""
	}
}
