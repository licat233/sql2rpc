/*
 * @Author: licat
 * @Date: 2023-02-08 16:56:31
 * @LastEditors: licat
 * @LastEditTime: 2023-02-16 15:30:51
 * @Description: licat233@gmail.com
 */

package _item

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

type Field interface {
	Name() string
	Value() interface{}
	Desc() string
	DefaultValue() interface{}
	Required() bool
	TypeName() string

	// GetString 获取具体的值
	GetString() string
	// GetBool 获取具体的值
	GetBool() bool
	// GetInt 获取具体的值
	GetInt() int

	// String 实现fmt.Sprintf接口，用于转化为yaml内容
	String() string

	Copy() *ConfigItem
	Set(interface{}) *ConfigItem
	Init(interface{}) *ConfigItem
	FlagString() (name string, value string, usage string)
	FlagBool() (name string, value bool, usage string)
	FlagInt() (name string, value int, usage string)
}

type ConfigItem struct {
	name         string      //字段名
	value        interface{} //字段值
	desc         string      //描述说明
	required     bool        //必须？
	defaultValue interface{} //默认值
	typename     string      //值类型
}

var _ Field = (*ConfigItem)(nil)

func NewConfigItem(name string, defaultValue interface{}, desc string, required bool) *ConfigItem {
	t := reflect.TypeOf(defaultValue).Kind().String()
	return &ConfigItem{
		name:         name,
		value:        defaultValue,
		desc:         desc,
		required:     required,
		defaultValue: defaultValue,
		typename:     t,
	}
}

func (c *ConfigItem) Name() string {
	return c.name
}

func (c *ConfigItem) Value() interface{} {
	return c.value
}

func (c *ConfigItem) Desc() string {
	if c.required {
		return fmt.Sprintf("%s (required)", c.desc)
	}
	return c.desc
}

func (c *ConfigItem) DefaultValue() interface{} {
	return c.defaultValue
}

func (c *ConfigItem) Required() bool {
	return c.required
}

func (c *ConfigItem) TypeName() string {
	return c.typename
}

// GetString 获取value字符串
func (c *ConfigItem) GetString() string {
	return fmt.Sprintf("%+v", c.value)
}

func (c *ConfigItem) GetBool() bool {
	if b, ok := c.value.(bool); ok {
		return b
	}
	if i, ok := c.value.(int); ok {
		return i > 0
	}
	if s, ok := c.value.(string); ok {
		s = strings.TrimSpace(s)
		return len(s) != 0
	}
	return false
}

func (c *ConfigItem) GetInt() int {
	if i, ok := c.value.(int); ok {
		return i
	}
	if b, ok := c.value.(bool); ok {
		if b {
			return 1
		}
		return 0
	}
	if s, ok := c.value.(string); ok {
		s = strings.TrimSpace(s)
		i, err := strconv.Atoi(s)
		if err != nil {
			return 0
		}
		return i
	}
	return 0
}

// String 用于转化为yaml内容
func (c *ConfigItem) String() string {
	desc := strings.TrimSpace(c.desc)
	if desc == "" {
		desc = c.name
	}
	var value = c.value
	if val, ok := c.value.(string); ok {
		val = strings.TrimSpace(val)
		value = fmt.Sprintf("\"%s\"", val)
	}
	return fmt.Sprintf("\n# %s\n%s: %+v", desc, c.name, value)
}

func (c *ConfigItem) Copy() *ConfigItem {
	var res = new(ConfigItem)
	*res = *c
	return res
}

// Set 设置value值，inValue 类型：string | bool | int
func (c *ConfigItem) Set(inValue interface{}) (ci *ConfigItem) {
	ci = c
	//检查输入的值类型是否正确
	inType := reflect.TypeOf(inValue).Kind().String()
	if inType != c.typename {
		//如果不正确，则修改失败
		panic(fmt.Errorf("type missmatch, have: %s, want: %s", inType, c.typename))
		return
	}
	// 配置数据来源优先级：文件->终端->默认值
	switch inType {
	case "string":
		if v, ok := inValue.(string); ok {
			c.value = v
			return
		}
	case "bool":
		if v, ok := inValue.(bool); ok {
			c.value = v
			return
		}
	case "int":
		if v, ok := inValue.(int); ok && v > 0 {
			c.value = v
			return
		}
	default:
		panic(fmt.Errorf("unsupported data type: %s", inType))
		return
	}
	return
}

// Init 初始化值，inValue 类型：string | bool | int
func (c *ConfigItem) Init(inValue interface{}) (ci *ConfigItem) {
	ci = c
	//检查输入的值类型是否正确
	inType := reflect.TypeOf(inValue).Kind().String()
	if inType != c.typename {
		//如果不正确，则修改失败
		panic(fmt.Errorf("type missmatch, have: %s, want: %s", inType, c.typename))
		return
	}

	//配置数据来源优先级：文件->终端->默认值
	switch inType {
	case "string":
		if fileValue := viper.GetString(c.name); fileValue != "" {
			c.value = fileValue
			return
		}
		if v, ok := inValue.(string); ok {
			c.value = v
			return
		}
	case "bool":
		if viper.GetBool(c.name) {
			c.value = true
			return
		}
		if v, ok := inValue.(bool); ok {
			c.value = v
			return
		}
	case "int":
		if fileValue := viper.GetInt(c.name); fileValue > 0 {
			c.value = fileValue
			return
		}
		if v, ok := inValue.(int); ok && v > 0 {
			c.value = v
			return
		}
	default:
		panic(fmt.Errorf("unsupported data type: %s", inType))
	}
	return
}

func (c *ConfigItem) FlagString() (name string, value string, usage string) {
	return c.Name(), c.GetString(), c.Desc()
}

func (c *ConfigItem) FlagBool() (name string, value bool, usage string) {
	return c.Name(), c.GetBool(), c.Desc()
}

func (c *ConfigItem) FlagInt() (name string, value int, usage string) {
	return c.Name(), c.GetInt(), c.Desc()
}
