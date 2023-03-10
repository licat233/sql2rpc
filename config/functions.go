package config

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/licat233/sql2rpc/tools"
	"github.com/spf13/viper"
)

func SetDefaultValue(param *string, value string) {
	if *param == "" {
		*param = value
	}
}

func FilterStringValue(values ...string) string {
	for _, value := range values {
		if v := strings.TrimSpace(value); v != "" {
			return v
		}
	}
	return ""
}

func FilterBoolValue(values ...bool) bool {
	for _, v := range values {
		if v {
			return v
		}
	}
	return false
}

func FilterIntValue(values ...int) int {
	for _, v := range values {
		if v != 0 {
			return v
		}
	}
	return 0
}

func SetStringValue(target *string, value1 string, value2 string) {
	if strings.TrimSpace(value1) != "" {
		*target = value1
	} else if strings.TrimSpace(value2) != "" {
		*target = value2
	}
}

func SetBoolValue(target *bool, value1 bool, value2 bool) {
	if value1 {
		*target = value1
	} else if value2 {
		*target = value2
	}
}

func SetIntValue(target *int, value1 int, value2 int) {
	if value1 != 0 {
		*target = value1
	} else if value2 != 0 {
		*target = value2
	}
}

func ExistConfigFile() (bool, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return false, err
	}
	// 判断文件是否存在
	configfilename := path.Join(currentDir, DefaultFileName)
	has, err := tools.PathExists(configfilename)
	return has, err
}

func InitViper() error {
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	//判断文件是否存在
	configPath := path.Join(currentDir, DefaultFileName)
	has, err := tools.PathExists(configPath)
	if err != nil {
		return err
	}
	if !has {
		return os.ErrNotExist
	}

	viper.SetConfigType("yaml")
	viper.SetConfigName(DefaultFileName)
	viper.AddConfigPath(currentDir) // path to look for the config file inionally look for config in the working directory
	err = viper.ReadInConfig()      // Find and read the config file
	if err != nil {                 // Handle errors reading the config file
		return fmt.Errorf("fatal error config file: %w", err)
	}
	return nil
}

func CreateConfigFile(configFilename string) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	//判断文件是否存在
	configfilename := path.Join(currentDir, configFilename)
	fd, e := os.ReadFile(configfilename)
	if e != nil && !errors.Is(e, os.ErrNotExist) {
		return e
	}
	if e == nil && len(fd) > 0 {
		fmt.Println("warning: ")
		fmt.Printf(" - %s already exists\n", configfilename)
		fmt.Println(" - 已存在配置文件，无需重复创建")
		os.Exit(1)
	}

	//创建文件
	fo, err := os.OpenFile(configfilename, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0766)
	if err != nil {
		return err
	}

	//写入文件
	_, err = fo.WriteString(NewDefaultConfig().String())
	if err != nil {
		return err
	}
	fmt.Println(" - 配置文件创建成功:", configfilename)
	return nil
}

func ParseIgnoreString(s string) []string {
	return strings.Split(s, ",")
}

func GetMark(name string) (startMark string, endMark string) {
	startMark = "// ------------------------------ " + name + " Start ------------------------------"
	endMark = "// ------------------------------ " + name + " End ------------------------------"
	return
}

func GetCustomMark(name string) (startMark string, endMark string) {
	startMark = "//[custom " + name + " start]"
	endMark = "//[custom " + name + " end]"
	return
}

func GetBaseMark(name string) (startMark string, endMark string) {
	startMark = "//[base " + name + " start]"
	endMark = "//[base " + name + " end]"
	return
}

func AddIgnoreColumns(cls ...string) []string {
	return append(IgnoreColumns, cls...)
}

func HeaderContent() string {
	buf := new(bytes.Buffer)
	buf.WriteString("/**\n")
	buf.WriteString(fmt.Sprintf(" * Generated by %s: %s\n", ProjectName, ProjectURL))
	buf.WriteString(fmt.Sprintf(" * Version: %s\n", CurrentVersion))
	buf.WriteString(" */\n\n")
	return buf.String()
}

func ParseToolName() (snake, camel, lowerCamel string) {
	name := strings.ReplaceAll(ProjectName, "2", "_")
	list := strings.Split(name, "_")
	if len(list) < 2 {
		snake = "snake"
		camel = "camel"
		lowerCamel = "lowerCamel"
		return
	}
	snake = tools.ToSnake(name)
	camel = tools.ToCamel(name)
	lowerCamel = tools.ToLowerCamel(name)
	return
}
