/*
 * @Author: licat
 * @Date: 2023-02-03 19:48:19
 * @LastEditors: licat
 * @LastEditTime: 2023-02-18 00:24:19
 * @Description: licat233@gmail.com
 */

package api

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/licat233/sql2rpc/cmd"
	"github.com/licat233/sql2rpc/cmd/api/_conf"
	"github.com/licat233/sql2rpc/cmd/api/_import"
	"github.com/licat233/sql2rpc/cmd/api/_info"
	"github.com/licat233/sql2rpc/cmd/api/_service"
	"github.com/licat233/sql2rpc/cmd/api/_struct"
	"github.com/licat233/sql2rpc/cmd/api/_struct/_field"

	"github.com/licat233/sql2rpc/cmd/common"
	"github.com/licat233/sql2rpc/config"
	"github.com/licat233/sql2rpc/db"
	"github.com/licat233/sql2rpc/tools"
)

type ApiCore struct {
	Filename string
	Info     *_info.Info
	Imports  _import.ImpCollection
	Structs  _struct.StructCollection
	Services _service.ServiceCollection
}

var _ cmd.Core = (*ApiCore)(nil)

func New() *ApiCore {
	return &ApiCore{
		Filename: config.C.Filename.GetString(),
		Info:     _info.New(),
		Imports:  _import.ImpCollection{},
		Structs:  _struct.StructCollection{},
		Services: _service.ServiceCollection{},
	}
}

func (s *ApiCore) Name() string {
	return "api"
}

func (s *ApiCore) Allow() bool {
	return config.C.Api.GetBool()
}

func (s *ApiCore) Init() error {
	_conf.InitConfig()
	if s.Filename == "" {
		s.Filename = fmt.Sprintf("%s.api", tools.ToLowerCamel(config.C.ServiceName.GetString()))
	} else {
		s.Filename = tools.SetFileType(s.Filename, ".api")
	}
	return nil
}

// Reset reset structs and services
func (s *ApiCore) Reset() {
	s.Structs = _struct.StructCollection{}
	s.Services = _service.ServiceCollection{}
}

func (s *ApiCore) Run() (err error) {
	s.Init()
	if config.C.ApiMultiple.GetBool() {
		err = s.GenerateMultipleFile()
	} else {
		err = s.GenerateSingleFile(config.C.DBTable.GetString(), s.Filename)
	}
	return
}

// MultipleFile 生成多个文件
func (s *ApiCore) GenerateMultipleFile() error {
	mainFilname := s.Filename
	has, err := tools.PathExists(mainFilname)
	if err != nil {
		return err
	}
	fileContent, f, err := tools.RTCFile(mainFilname)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		e := f.Close()
		if e != nil {
			err = e
		}
	}(f)

	tableStr := config.C.DBTable.GetString()
	imports := _import.ImpCollection{}
	var tables []string
	// get tables
	if tableStr != "" && tableStr != "*" {
		tables = strings.Split(strings.Trim(tableStr, ","), ",")
	} else {
		tables, err = db.ShowTables()
		if err != nil {
			return err
		}
	}

	//暂时关闭，不让子文件生成baseStructs，见structs.go 98行
	_conf.CurrentIsCoreFile = false
	ignoreTables := make(map[string]bool)
	for _, ig := range _conf.IgnoreTables {
		ignoreTables[ig] = true
	}

	schema := config.C.DBSchema.GetString()
	for _, table := range tables {
		if _, ok := ignoreTables[table]; ok {
			continue
		}
		filename := fmt.Sprintf("%s.api", tools.ToLowerCamel(fmt.Sprintf("%s_%s", schema, table)))
		if filename == mainFilname {
			oldFileName := mainFilname
			mainFilname = tools.FileRename(mainFilname, fmt.Sprintf("%s_%s", tools.GetFilename(mainFilname), "core"))
			fmt.Println(" - The file name conflicts. Please reset the file name. The file " + oldFileName + ".api rename  to '" + mainFilname + "'")
		}
		imports = append(imports, _import.New(filename))
		err := s.GenerateSingleFile(table, filename)
		if nil != err {
			return err
		}
		s.Reset()
	}

	if fileContent == "" {
		fileContent = _service.GenarateDefaultService()
	}

	_conf.FileContent = fileContent
	_conf.CurrentIsCoreFile = true
	s.Imports = imports

	sort.Sort(s.Imports)

	// write
	_, err = f.WriteString(fmt.Sprint(s))
	if err != nil {
		return err
	}

	if has {
		fmt.Println(config.UpdatedFileMsg, mainFilname)
	} else {
		fmt.Println(config.CreatedFileMsg, mainFilname)
	}
	return nil
}

// SingleFile 生成一个文件
func (s *ApiCore) GenerateSingleFile(table, filename string) error {
	filename = tools.SetFileType(filename, ".api")
	has, err := tools.PathExists(filename)
	if err != nil {
		return err
	}
	fileContent, f, err := tools.RTCFile(filename)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		e := f.Close()
		if e != nil {
			err = e
		}
	}(f)

	if fileContent == "" {
		fileContent = _service.GenarateDefaultService()
	}

	_conf.FileContent = fileContent

	dbs, err := common.DbSchema(db.Conn)
	if nil != err {
		return err
	}

	cols, err := common.DbColumns(db.Conn, dbs, table)
	if nil != err {
		return err
	}

	err = s.typesFromColumns(cols, _conf.IgnoreTables, _conf.IgnoreColumns)
	if nil != err {
		return err
	}

	sort.Sort(s.Imports)
	sort.Sort(s.Structs)
	sort.Sort(s.Services)

	// write
	_, err = f.WriteString(fmt.Sprint(s))
	if has {
		fmt.Println(config.UpdatedFileMsg, filename)
	} else {
		fmt.Println(config.CreatedFileMsg, filename)
	}
	return err
}

// AppendImport adds an importto a schema if the specific importdoes not already exist in the schema.
func (s *ApiCore) AppendImport(imports string) {
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
func (s *ApiCore) String() string {
	buf := new(bytes.Buffer)
	// header start
	buf.WriteString("syntax = \"v1\"")
	// header end

	// info start
	buf.WriteString(fmt.Sprint(s.Info))
	// info end

	// import start
	buf.WriteString(fmt.Sprint(s.Imports))
	// import end

	//struct start
	buf.WriteString(fmt.Sprint(s.Structs))
	//struct end

	buf.WriteString("\n\n")

	// service start
	buf.WriteString(fmt.Sprint(s.Services))
	// service end

	return buf.String()
}

// typesFromColumns creates the appropriate schema properties from a collection of column types.
func (s *ApiCore) typesFromColumns(cols []*common.Column, ignoreTables, ignoreColumns []string) error {
	structMap := map[string]*_struct.Struct{}
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

		structName := c.TableName

		stt, ok := structMap[structName]
		if !ok {
			structMap[structName] = _struct.NewStruct(structName, "json", c.TableComment, nil)
			stt = structMap[structName]
		}

		err := s.parseColumn(stt, c)
		if nil != err {
			return err
		}
	}

	s.Structs = _struct.StructCollection{}
	s.Services = _service.ServiceCollection{}
	for _, v := range structMap {
		s.Structs = append(s.Structs, v)
		svr := _service.NewService(v.Name, v.Comment)
		s.Services = append(s.Services, svr)
	}

	return nil
}

func (s *ApiCore) parseColumn(stt *_struct.Struct, col *common.Column) error {
	typ := strings.ToLower(col.DataType)
	var fieldType string

	switch typ {
	case "char", "varchar", "text", "longtext", "mediumtext", "tinytext":
		fieldType = "string"
	case "blob", "mediumblob", "longblob", "varbinary", "binary":
		fieldType = "bytes"
	case "date", "time", "datetime", "timestamp":
		fieldType = "int64"
	case "bool", "bit":
		fieldType = "bool"
	case "tinyint", "smallint", "int", "mediumint", "bigint":
		fieldType = "int64"
	case "float", "decimal", "double":
		fieldType = "float64"
	case "json":
		fieldType = "string"
	default:
		fieldType = ""
	}

	if fieldType == "" {
		return fmt.Errorf("no compatible golang type found for `%s`. column: `%s`.`%s`", col.DataType, col.TableName, col.ColumnName)
	}

	sf := _field.NewStructField(col.ColumnName, fieldType, "json", col.ColumnName, "", col.ColumnComment)

	stt.Fields = append(stt.Fields, sf)

	return nil
}
