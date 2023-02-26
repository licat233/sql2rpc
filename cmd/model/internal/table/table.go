package table

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/licat233/sql2rpc/cmd/common"
	"github.com/licat233/sql2rpc/cmd/model/internal/conf"
	"github.com/licat233/sql2rpc/config"
	"github.com/licat233/sql2rpc/tools"
)

type Table struct {
	Name    string
	Columns []*common.Column
}

func (t *Table) String() string {
	var buf = new(bytes.Buffer)
	buf.WriteString(t.FindList())
	return buf.String()
}

func (t *Table) FindList() string {
	structName := tools.ToCamel(t.Name)
	basisName := tools.ToLowerCamel(t.Name)
	var buf = new(bytes.Buffer)
	funcString := fmt.Sprintf("\nfunc (m *custom%sModel) FindList(ctx context.Context, pageSize, page int64, keyword string, %s *%s) (resp []*%s, total int64, err error) {", structName, basisName, structName, structName)
	buf.WriteString(funcString)
	baseSq := fmt.Sprintf("\n\tsq := squirrel.Select(%sRows).From(m.table)", basisName)
	buf.WriteString(baseSq)
	buf.WriteString(fmt.Sprintf("\n\tif %s != nil {", basisName))
	t.thanString(buf)
	buf.WriteString("\n\t}")
	buf.WriteString("\n\tif pageSize > 0 && page > 0 {")
	buf.WriteString("\n\t\tsqCount := sq.RemoveLimit().RemoveOffset()")
	buf.WriteString("\n\t\tsq = sq.Offset(uint64((page - 1) * pageSize)).Limit(uint64(pageSize))")
	buf.WriteString("\n\t\tqueryCount, agrsCount, e := sqCount.ToSql()")
	buf.WriteString("\n\t\tif e != nil {\n\t\t\terr = e\n\t\t\treturn\n\t\t}")
	buf.WriteString(fmt.Sprintf("\n\t\tqueryCount = strings.ReplaceAll(queryCount, %sRows, \"COUNT(*)\")", basisName))
	buf.WriteString("\n\t\tif err = m.conn.QueryRowCtx(ctx, &total, queryCount, agrsCount...); err != nil {\n\t\t\treturn\n\t\t}")
	buf.WriteString("\n\t}")
	buf.WriteString("\n\tquery, agrs, err := sq.ToSql()\n\tif err != nil {\n\t\treturn\n\t}\n\treturn")
	buf.WriteString(fmt.Sprintf("\n\tresp = make([]*%s, 0)", structName))
	buf.WriteString("\n\tif err = m.conn.QueryRowsCtx(ctx, &resp, query, agrs...); err != nil {\n\t\treturn\n\t}")
	buf.WriteString("\n\treturn")
	buf.WriteString("\n}\n")
	return buf.String()
}

func (t *Table) thanString(buf *bytes.Buffer) {
	basisName := tools.ToLowerCamel(t.Name)
	for _, c := range t.Columns {
		var than string
		//判断是字符串，还是数字
		tName := convTypeName(c.ColumnType)
		if tName == "int64" || tName == "int" || tName == "float64" || tName == "float32" {
			than = ">= 0"
			if isIdColumn(c.ColumnName) {
				than = "> 0"
			}
		} else if tName == "string" {
			than = "!= \"\""
		} else {
			continue
		}
		cName := tools.ToCamel(c.ColumnName)
		cV := fmt.Sprintf("%s.%s", basisName, cName)
		buf.WriteString(fmt.Sprintf("\n\t\tif %s %s {", cV, than))
		buf.WriteString(fmt.Sprintf("\n\t\t\tsq = sq.Where(\"%s = ?\", %s)", c.ColumnName, cV))
		buf.WriteString("\n\t\t}")
	}
}

func isIdColumn(name string) bool {
	name = tools.ToSnake(name)
	names := strings.Split(name, "_")
	n := len(names)
	if n == 0 {
		return false
	}
	return names[n-1] == "id"
}

func convTypeName(columnType string) string {
	typ := strings.ToLower(columnType)
	switch typ {
	case "char", "varchar", "text", "longtext", "mediumtext", "tinytext":
		return "string"
	case "blob", "mediumblob", "longblob", "varbinary", "binary":
		return "bytes"
	case "date", "time", "datetime", "timestamp":
		return "time"
	case "bool", "bit":
		return "bool"
	case "tinyint", "smallint", "int", "mediumint", "bigint":
		return "int64"
	case "float", "decimal", "double":
		return "float64"
	case "json":
		return "string"
	default:
		return ""
	}
}

func (t *Table) UpdateGoModelFile() error {
	filename := fmt.Sprintf("%sModel.go", tools.ToLowerCamel(t.Name))
	has, err := tools.PathExists(filename)
	if err != nil {
		return err
	}
	if !has {
		return fmt.Errorf("%sModel.go文件不存在，请先使用goctl工具创建", tools.ToLowerCamel(t.Name))
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
	conf.FileContent = fileContent
	content := t.String()
	fileContent, err = common.UpdateMarkContent(config.BaseFuncsStartMark, config.BaseFuncsEndMark, fileContent, content)
	if err != nil {
		return err
	}
	_, err = f.WriteString(fileContent)
	if has {
		fmt.Println(config.UpdatedFileMsg, filename)
	} else {
		fmt.Println(config.CreatedFileMsg, filename)
	}
	return err
}
