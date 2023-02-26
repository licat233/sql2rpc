package table

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/licat233/sql2rpc/cmd/common"
	"github.com/licat233/sql2rpc/cmd/model/internal/conf"
	"github.com/licat233/sql2rpc/config"
	"github.com/licat233/sql2rpc/tools"
)

type Table struct {
	Name          string
	Columns       []*common.Column
	HasNameColumn bool
	interfaceName string
	strcutName    string
	findListName  string
}

func New(name string, columns []*common.Column) *Table {
	lowerName := tools.ToLowerCamel(name)
	camelName := tools.ToCamel(name)
	findList := fmt.Sprintf("FindList(ctx context.Context, pageSize, page int64, keyword string, %s *%s) (resp []*%s, total int64, err error)", lowerName, camelName, camelName)
	return &Table{
		Name:          name,
		Columns:       columns,
		HasNameColumn: false,
		interfaceName: fmt.Sprintf("%smodel", lowerName),
		strcutName:    fmt.Sprintf("extend%sModel", camelName),
		findListName:  findList,
	}
}

func (t *Table) String() string {
	// lowerName := tools.ToLowerCamel(t.Name)
	camelName := tools.ToCamel(t.Name)
	var buf = new(bytes.Buffer)
	dir, _ := tools.GetCurrentDirectoryName()
	buf.WriteString(fmt.Sprintf("package %s\n\n", dir))
	buf.WriteString("import (")
	buf.WriteString("\n\t\"context\"")
	if t.hasName() {
		buf.WriteString("\n\t\"fmt\"")
	}
	buf.WriteString("\n\t\"strings\"")
	buf.WriteString("\n\t\"github.com/Masterminds/squirrel\"")
	buf.WriteString("\n)\n")

	buf.WriteString("\ntype (")
	buf.WriteString(fmt.Sprintf("\n\t%s interface {", t.interfaceName))

	buf.WriteString(fmt.Sprintf("\n\t\t%s", t.findListName))

	buf.WriteString("\n\t}")
	buf.WriteString(fmt.Sprintf("\n\t%s struct {", t.strcutName))
	buf.WriteString(fmt.Sprintf("\n\t\t*default%sModel", camelName))
	buf.WriteString("\n\t}")
	buf.WriteString("\n)\n")

	buf.WriteString(t.FindList())
	return buf.String()
}

func (t *Table) GenFile() error {
	filename := fmt.Sprintf("%smodel_extend.go", strings.ToLower(t.Name))
	filename = path.Join(config.C.Dir.GetString(), filename)
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
	buf := new(bytes.Buffer)

	content := t.String()

	content, _ = tools.FormatGoContent(content)

	buf.WriteString(content)
	// write
	_, err = f.WriteString(buf.String())
	if has {
		fmt.Println(config.UpdatedFileMsg, filename)
	} else {
		fmt.Println(config.CreatedFileMsg, filename)
	}
	return nil
}

func (t *Table) FindList() string {
	t.HasNameColumn = t.hasName()
	camelName := tools.ToCamel(t.Name)
	lowerName := tools.ToLowerCamel(t.Name)
	var buf = new(bytes.Buffer)
	funcString := fmt.Sprintf("\nfunc (m *%s) %s {", t.strcutName, t.findListName)
	buf.WriteString(funcString)
	if t.HasNameColumn {
		buf.WriteString("\n\thasName := false")
	}
	baseSq := fmt.Sprintf("\n\tsq := squirrel.Select(%sRows).From(m.table)", lowerName)
	buf.WriteString(baseSq)

	t.thanString(buf)

	buf.WriteString("\n\tif pageSize > 0 && page > 0 {")
	buf.WriteString("\n\t\tsqCount := sq.RemoveLimit().RemoveOffset()")
	buf.WriteString("\n\t\tsq = sq.Offset(uint64((page - 1) * pageSize)).Limit(uint64(pageSize))")
	buf.WriteString("\n\t\tqueryCount, agrsCount, e := sqCount.ToSql()")
	buf.WriteString("\n\t\tif e != nil {\n\t\t\terr = e\n\t\t\treturn\n\t\t}")
	buf.WriteString(fmt.Sprintf("\n\t\tqueryCount = strings.ReplaceAll(queryCount, %sRows, \"COUNT(*)\")", lowerName))
	buf.WriteString("\n\t\tif err = m.conn.QueryRowCtx(ctx, &total, queryCount, agrsCount...); err != nil {\n\t\t\treturn\n\t\t}")
	buf.WriteString("\n\t}")
	buf.WriteString("\n\tquery, agrs, err := sq.ToSql()\n\tif err != nil {\n\t\treturn\n\t}")
	buf.WriteString(fmt.Sprintf("\n\tresp = make([]*%s, 0)", camelName))
	buf.WriteString("\n\tif err = m.conn.QueryRowsCtx(ctx, &resp, query, agrs...); err != nil {\n\t\treturn\n\t}")
	buf.WriteString("\n\treturn")
	buf.WriteString("\n}\n")
	return buf.String()
}

func (t *Table) thanString(buf *bytes.Buffer) {
	basisName := tools.ToLowerCamel(t.Name)
	hasName := false
	buf.WriteString(fmt.Sprintf("\n\tif %s != nil {", basisName))
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
			// fmt.Println("unknow column type:", c.ColumnName, "-", tName, "-", c.ColumnType)
			continue
		}
		cName := tools.ToCamel(c.ColumnName)
		cV := fmt.Sprintf("%s.%s", basisName, cName)
		buf.WriteString(fmt.Sprintf("\n\t\tif %s %s {", cV, than))
		buf.WriteString(fmt.Sprintf("\n\t\t\tsq = sq.Where(\"%s = ?\", %s)", c.ColumnName, cV))
		if strings.ToLower(c.ColumnName) == "name" {
			buf.WriteString("\n\t\t\thasName = true")
			hasName = true
		}
		buf.WriteString("\n\t\t}")
	}
	buf.WriteString("\n\t}")
	if hasName {
		buf.WriteString("\n\tif keyword != \"\" && !hasName {")
		buf.WriteString("\n\t\tsq = sq.Where(\"name LIKE ?\", fmt.Sprintf(\"%%%s%%\", keyword))")
		buf.WriteString("\n\t}")
	}
}

func (t *Table) hasName() bool {
	if t.HasNameColumn {
		return true
	}
	for _, c := range t.Columns {
		if strings.ToLower(c.ColumnName) == "name" {
			tName := convTypeName(c.ColumnType)
			if tName == "string" {
				t.HasNameColumn = true
				return true
			}
			return true
		}
	}
	return false
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
	fileContent = common.UpdateMarkContent(config.BaseFuncsStartMark, config.BaseFuncsEndMark, fileContent, content)
	_, err = f.WriteString(fileContent)
	if has {
		fmt.Println(config.UpdatedFileMsg, filename)
	} else {
		fmt.Println(config.CreatedFileMsg, filename)
	}
	return err
}
