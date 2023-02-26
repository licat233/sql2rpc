package model

import (
	"bytes"
	"fmt"
	"os"
	"sort"

	"github.com/licat233/sql2rpc/cmd"
	"github.com/licat233/sql2rpc/cmd/common"
	"github.com/licat233/sql2rpc/cmd/model/internal/conf"
	"github.com/licat233/sql2rpc/cmd/model/internal/table"
	"github.com/licat233/sql2rpc/config"
	"github.com/licat233/sql2rpc/db"
	"github.com/licat233/sql2rpc/tools"
)

type ModelCore struct {
	FilePath string
	Tables   table.TableCollection
}

var _ cmd.Core = (*ModelCore)(nil)

func New() *ModelCore {
	name := tools.ToLowerCamel(config.C.DBSchema.GetString())
	filename := fmt.Sprintf("./%sModel_extend.go", name)
	return &ModelCore{
		FilePath: filename,
		Tables:   []*table.Table{},
	}
}

func (m *ModelCore) Name() string {
	return config.ModelCoreName
}

func (m *ModelCore) Allow() bool {
	return config.C.Model.GetBool()
}

func (m *ModelCore) init() error {
	conf.InitConfig()
	dbs, err := common.DbSchema(db.Conn)
	if nil != err {
		return err
	}
	cols, err := common.DbColumns(db.Conn, dbs, config.C.DBTable.GetString())
	if nil != err {
		return err
	}
	m.tablesFromColumns(cols, conf.IgnoreTables, conf.IgnoreColumns)
	sort.Sort(m.Tables)
	return nil
}

func (m *ModelCore) Run() (err error) {
	m.init()
	err = m.UpdateMultipleFile()
	return
}

func (m *ModelCore) UpdateMultipleFile() error {
	for _, table := range m.Tables {
		if err := table.UpdateGoModelFile(); err != nil {
			return err
		}
	}
	return nil
}

func (m *ModelCore) GenSingleFile() error {
	filename := m.FilePath
	m.FilePath = tools.SetFileType(filename, ".go")
	has, err := tools.PathExists(filename)
	if err != nil {
		return err
	}
	_, f, err := tools.RTCFile(filename)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		e := f.Close()
		if e != nil {
			err = e
		}
	}(f)
	dir, err := tools.GetCurrentDirectoryName()
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	buf.WriteString(fmt.Sprintf("package %s\n\n", dir))
	buf.WriteString(fmt.Sprint(m.Tables))
	// write
	_, err = f.WriteString(buf.String())
	if has {
		fmt.Println(config.UpdatedFileMsg, filename)
	} else {
		fmt.Println(config.CreatedFileMsg, filename)
	}
	return nil
}

func (m *ModelCore) tablesFromColumns(cols []*common.Column, ignoreTables, ignoreColumns []string) {
	ignoreMap := map[string]bool{}
	ignoreColumnMap := map[string]bool{}
	for _, ig := range ignoreTables {
		ignoreMap[ig] = true
	}
	for _, ic := range ignoreColumns {
		ignoreColumnMap[ic] = true
	}

	tableMap := map[string]*table.Table{}
	for _, c := range cols {
		if _, ok := ignoreMap[c.TableName]; ok {
			continue
		}
		if _, ok := ignoreColumnMap[c.ColumnName]; ok {
			continue
		}

		tableName := c.TableName
		msg, ok := tableMap[tableName]
		if !ok {
			tableMap[tableName] = table.New(tableName, []*common.Column{c})
		} else {
			msg.Columns = append(msg.Columns, c)
		}
	}
	for _, table := range tableMap {
		m.Tables = append(m.Tables, table)
	}
}
