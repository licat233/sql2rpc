/*
 * @Author: licat
 * @Date: 2023-02-07 11:12:47
 * @LastEditors: licat
 * @LastEditTime: 2023-02-16 13:26:14
 * @Description: licat233@gmail.com
 */
package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/licat233/sql2rpc/config"
)

var Conn *sql.DB

func InitConn() (err error) {
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", config.C.DBUser.GetString(), config.C.DBPassword.GetString(), config.C.DBHost.GetString(), config.C.DBPort.GetInt(), config.C.DBSchema.GetString())
	Conn, err = sql.Open(config.C.DBType.GetString(), dataSource)
	if err != nil {
		return
	}
	if err = Conn.Ping(); err != nil {
		return
	}
	return
}

// ShowTables show all tables in database
func ShowTables() ([]string, error) {
	var tables []string
	rows, err := Conn.Query("show tables")
	if err != nil {
		return tables, err
	}
	defer rows.Close()
	for rows.Next() {
		var table string
		err = rows.Scan(&table)
		if err != nil {
			return tables, err
		}
		tables = append(tables, table)
	}
	return tables, nil
}
