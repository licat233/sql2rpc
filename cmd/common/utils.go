/*
 * @Author: licat
 * @Date: 2023-02-10 00:51:06
 * @LastEditors: licat
 * @LastEditTime: 2023-02-20 11:56:40
 * @Description: licat233@gmail.com
 */
package common

import (
	"github.com/licat233/sql2rpc/config"
	"github.com/licat233/sql2rpc/tools"
)

// ConvertStringStyle 转化字符串风格，默认 snake 风格
func ConvertStringStyle(style, name string) string {
	switch style {
	case config.CamelCase:
		return tools.ToCamel(name)
	case config.LowerCamelCase:
		return tools.ToLowerCamel(name)
	case config.SnakeCase:
		return tools.ToSnake(name)
	default:
		return tools.ToSnake(name)
	}
}
