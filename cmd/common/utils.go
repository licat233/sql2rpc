package common

import (
	"fmt"
	"strings"

	"github.com/licat233/sql2rpc/config"
	"github.com/licat233/sql2rpc/tools"
)

// ConvertStringStyle 转化字符串风格，默认 snake 风格
func ConvertStringStyle(s string) string {
	switch config.C.Style.GetString() {
	case config.CamelCase:
		return tools.ToCamel(s)
	case config.LowerCamelCase:
		return tools.ToLowerCamel(s)
	case config.SnakeCase:
		return tools.ToSnake(s)
	default:
		return tools.ToSnake(s)
	}
}

// FormatContent 格式化内容
func FormatContent(str string, indent string) string {
	s := strings.Split(str, "\n")
	for i := 0; i < len(s); i++ {
		s[i] = fmt.Sprintf("%s%s", indent, strings.TrimSpace(s[i]))
	}
	return strings.Join(s, "\n")
}
