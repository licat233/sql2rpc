/*
 * @Author: licat
 * @Date: 2023-02-03 19:48:19
 * @LastEditors: licat
 * @LastEditTime: 2023-02-18 12:30:41
 * @Description: licat233@gmail.com
 */

package pb

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/licat233/sql2rpc/cmd"
	"github.com/licat233/sql2rpc/cmd/pb/_conf"
	"github.com/licat233/sql2rpc/cmd/pb/_enum"
	"github.com/licat233/sql2rpc/cmd/pb/_import"
	"github.com/licat233/sql2rpc/cmd/pb/_message"
	"github.com/licat233/sql2rpc/cmd/pb/_message/_field"
	"github.com/licat233/sql2rpc/cmd/pb/_service"

	"github.com/licat233/sql2rpc/cmd/common"
	"github.com/licat233/sql2rpc/config"
	"github.com/licat233/sql2rpc/db"
	"github.com/licat233/sql2rpc/tools"

	"github.com/chuckpreslar/inflect"
)

// PbCore is a representation of a protobuf schema.
type PbCore struct {
	FilePath string
	Imports  _import.ImpCollection
	Messages _message.MessageCollection
	Enums    _enum.EnumCollection
	Services _service.ServiceCollection
}

var _ cmd.Core = (*PbCore)(nil)

func New() *PbCore {
	filename := config.C.Filename.GetString()
	if filename == "" {
		filename = tools.ToLowerCamel(config.C.ServiceName.GetString())
	}
	filename = tools.SetFileType(filename, ".proto")
	return &PbCore{
		FilePath: filename,
		Imports:  _import.ImpCollection{},
		Messages: _message.MessageCollection{},
		Enums:    _enum.EnumCollection{},
		Services: _service.ServiceCollection{},
	}
}

func (s *PbCore) Name() string {
	return "pb"
}

func (s *PbCore) Allow() bool {
	return config.C.Pb.GetBool()
}

func (s *PbCore) Run() (err error) {
	_conf.InitConfig()
	err = s.Gen()
	return
}

func (s *PbCore) Gen() error {
	s.FilePath = tools.SetFileType(s.FilePath, ".proto")
	has, err := tools.PathExists(s.FilePath)
	if err != nil {
		return err
	}
	fileContent, f, err := tools.RTCFile(s.FilePath)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		e := f.Close()
		if e != nil {
			err = e
		}
	}(f)

	//只在创建的时候才会加入默认service板块
	if fileContent == "" && config.C.PbMultiple.GetBool() {
		fileContent = _service.GenarateDefaultCustomService()
	}

	_conf.FileContent = fileContent

	dbs, err := common.DbSchema(db.Conn)
	if nil != err {
		return err
	}

	cols, err := common.DbColumns(db.Conn, dbs, config.C.DBTable.GetString())
	if nil != err {
		return err
	}

	err = s.typesFromColumns(cols, _conf.IgnoreTables, _conf.IgnoreColumns)
	if nil != err {
		return err
	}

	sort.Sort(s.Imports)
	sort.Sort(s.Messages)
	sort.Sort(s.Enums)
	sort.Sort(s.Services)

	// write
	_, err = f.WriteString(fmt.Sprint(s))
	if err != nil {
		return err
	}

	if has {
		fmt.Println(config.UpdatedFileMsg, s.FilePath)
	} else {
		fmt.Println(config.CreatedFileMsg, s.FilePath)
	}

	return nil
}

// AppendImport adds an importto a schema if the specific importdoes not already exist in the schema.
func (s *PbCore) AppendImport(imports string) {
	shouldAdd := true
	for _, si := range s.Imports {
		if si.Filename == imports {
			shouldAdd = false
			break
		}
	}

	if shouldAdd {
		s.Imports = append(s.Imports, _import.New(imports))
	}
}

// String returns a string representation of a Schema.
// 用于实现 fmt.Stringer 接口
func (s *PbCore) String() string {
	buf := new(bytes.Buffer)
	buf.WriteString(config.HeaderContent())
	buf.WriteString(fmt.Sprintf("syntax = \"%s\";\n", config.Syntax))
	buf.WriteString("\n")
	buf.WriteString(fmt.Sprintf("package %s;\n", config.C.PbPackage.GetString()))
	buf.WriteString("\n")
	buf.WriteString(fmt.Sprintf("option go_package =\"%s\";\n", config.C.PbGoPackage.GetString()))

	// import start
	buf.WriteString(fmt.Sprint(s.Imports))
	// import end

	//enum start
	buf.WriteString(fmt.Sprint(s.Enums))
	//enum end

	//message start
	buf.WriteString(fmt.Sprint(s.Messages))
	//message end

	buf.WriteString("\n\n")

	// service start
	buf.WriteString(fmt.Sprint(s.Services))
	// service end

	content := common.FormatContent(buf.String())

	return content
}

// typesFromColumns creates the appropriate schema properties from a collection of column types.
func (s *PbCore) typesFromColumns(cols []*common.Column, ignoreTables, ignoreColumns []string) error {
	messageMap := map[string]*_message.Message{}
	ignoreMap := map[string]bool{}
	ignoreColumnMap := map[string]bool{}
	for _, ig := range ignoreTables {
		ignoreMap[ig] = true
	}
	for _, ic := range ignoreColumns {
		ignoreColumnMap[ic] = true
	}

	for _, c := range cols {
		if _, ok := ignoreMap[c.TableName]; ok {
			continue
		}
		if _, ok := ignoreColumnMap[c.ColumnName]; ok {
			continue
		}

		messageName := c.TableName

		msg, ok := messageMap[messageName]
		if !ok {
			messageMap[messageName] = &_message.Message{
				Name:    messageName,
				Comment: c.TableComment,
				Fields:  []*_field.MessageField{},
			}
			msg = messageMap[messageName]
		}

		err := s.parseColumn(msg, c)
		if nil != err {
			return err
		}
	}

	for _, v := range messageMap {
		s.Messages = append(s.Messages, v)
		s.Services = append(s.Services, _service.New(v.Name, v.Comment))
	}

	return nil
}

// ParseColumn parses a column and inserts the relevant fields in the Message. If an enumerated type is encountered, an Enum will
// be added to the Schema. Returns an error if an incompatible protobuf data type cannot be found for the database column type.
func (s *PbCore) parseColumn(msg *_message.Message, col *common.Column) error {
	typ := strings.ToLower(col.DataType)
	var fieldType string

	switch typ {
	case "char", "varchar", "text", "longtext", "mediumtext", "tinytext":
		fieldType = "string"
	case "_enum", "set":
		// Parse c.ColumnType to get the _enum list
		enumList := regexp.MustCompile(`[_enum|set]\((.+?)\)`).FindStringSubmatch(col.ColumnType)
		enums := strings.FieldsFunc(enumList[1], func(c rune) bool {
			cs := string(c)
			return cs == "," || cs == "'"
		})

		enumName := inflect.Singularize(tools.ToCamel(col.TableName)) + tools.ToCamel(col.ColumnName)
		enum, err := _enum.NewEnumFromStrings(enumName, col.ColumnComment, enums)
		if nil != err {
			return err
		}

		s.Enums = append(s.Enums, enum)

		fieldType = enumName
	case "blob", "mediumblob", "longblob", "varbinary", "binary":
		fieldType = "bytes"
	case "date", "time", "datetime", "timestamp":
		//s.AppendImport("google/protobuf/timestamp.proto")
		fieldType = "int64"
	case "bool", "bit":
		fieldType = "bool"
	case "tinyint", "smallint", "int", "mediumint", "bigint":
		fieldType = "int64"
	case "float", "decimal", "double":
		fieldType = "double"
	case "json":
		fieldType = "string"
	}

	if fieldType == "" {
		return fmt.Errorf("no compatible protobuf type found for `%s`. column: `%s`.`%s`", col.DataType, col.TableName, col.ColumnName)
	}

	field := _field.New(fieldType, col.ColumnName, len(msg.Fields)+1, col.ColumnComment)

	err := msg.AppendField(field)
	if nil != err {
		return err
	}

	return nil
}
