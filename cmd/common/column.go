/*
 * @Author: licat
 * @Date: 2023-02-05 18:01:41
 * @LastEditors: licat
 * @LastEditTime: 2023-02-09 23:57:37
 * @Description: licat233@gmail.com
 */
package common

import (
	"database/sql"
	"log"
	"regexp"
	"strings"
)

// Column represents a database column.
type Column struct {
	Style                  string
	TableName              string
	TableComment           string
	ColumnName             string
	IsNullable             string
	DataType               string
	CharacterMaximumLength sql.NullInt64
	NumericPrecision       sql.NullInt64
	NumericScale           sql.NullInt64
	ColumnType             string
	ColumnComment          string
}

func DbSchema(db *sql.DB) (string, error) {
	var schema string

	err := db.QueryRow("SELECT SCHEMA()").Scan(&schema)

	return schema, err
}

func DbColumns(db *sql.DB, schema, table string) ([]*Column, error) {

	tableArr := strings.Split(strings.Trim(table, ","), ",")

	q := "SELECT c.TABLE_NAME, c.COLUMN_NAME, c.IS_NULLABLE, c.DATA_TYPE, " +
		"c.CHARACTER_MAXIMUM_LENGTH, c.NUMERIC_PRECISION, c.NUMERIC_SCALE, c.COLUMN_TYPE ,c.COLUMN_COMMENT,t.TABLE_COMMENT " +
		"FROM INFORMATION_SCHEMA.COLUMNS as c  LEFT JOIN  INFORMATION_SCHEMA.TABLES as t  on c.TABLE_NAME = t.TABLE_NAME and  c.TABLE_SCHEMA = t.TABLE_SCHEMA" +
		" WHERE c.TABLE_SCHEMA = ?"

	if table != "" && table != "*" {
		q += " AND c.TABLE_NAME IN('" + strings.TrimRight(strings.Join(tableArr, "' ,'"), ",") + "')"
	}

	q += " ORDER BY c.TABLE_NAME, c.ORDINAL_POSITION"

	rows, err := db.Query(q, schema)
	defer func() {
		rows.Close()
	}()
	if nil != err {
		return nil, err
	}

	cols := []*Column{}

	for rows.Next() {
		cs := &Column{}
		err := rows.Scan(&cs.TableName, &cs.ColumnName, &cs.IsNullable, &cs.DataType,
			&cs.CharacterMaximumLength, &cs.NumericPrecision, &cs.NumericScale, &cs.ColumnType, &cs.ColumnComment, &cs.TableComment)
		if err != nil {
			log.Fatal(err)
		}
		typName := cs.ColumnType
		re := regexp.MustCompile(`\(\d*\)`)
		typName = re.ReplaceAllString(typName, "")
		cs.ColumnType = typName
		cols = append(cols, cs)
	}
	if err := rows.Err(); nil != err {
		return nil, err
	}

	return cols, nil
}
