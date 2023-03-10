/*
 * @Author: licat
 * @Date: 2023-02-05 17:14:15
 * @LastEditors: licat
 * @LastEditTime: 2023-02-20 14:27:33
 * @Description: licat233@gmail.com
 */
package common

import (
	"fmt"
	"testing"

	"github.com/licat233/sql2rpc/config"
)

func TestUpdateMarkContent(t *testing.T) {
	fileContent := `
	package model_single_file

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ AdminerModel = (*customAdminerModel)(nil)

type (
	// AdminerModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAdminerModel.
	AdminerModel interface {
		adminerModel
	}

	customAdminerModel struct {
		*defaultAdminerModel
	}
)

// NewAdminerModel returns a model for the database table.
func NewAdminerModel(conn sqlx.SqlConn) AdminerModel {
	return &customAdminerModel{
		defaultAdminerModel: newAdminerModel(conn),
	}
}


//[base Funcs start]


func (m *customAdminerModel) FindList(ctx context.Context, pageSize, page int64, keyword string, adminer *Adminer) (resp []*Adminer, total int64, err error) {
	hasName := false
	sq := squirrel.Select(adminerRows).From(m.table)
	if adminer != nil {
		if adminer.Id > 0 {
			sq = sq.Where("id = ?", adminer.Id)
		}
		if adminer.Uuid != "" {
			sq = sq.Where("uuid = ?", adminer.Uuid)
		}
		if adminer.Name != "" {
			sq = sq.Where("name = ?", adminer.Name)
			hasName = true
		}
		if adminer.Avatar != "" {
			sq = sq.Where("avatar = ?", adminer.Avatar)
		}
		if adminer.Passport != "" {
			sq = sq.Where("passport = ?", adminer.Passport)
		}
		if adminer.Password != "" {
			sq = sq.Where("password = ?", adminer.Password)
		}
		if adminer.Email != "" {
			sq = sq.Where("email = ?", adminer.Email)
		}
		if adminer.Status >= 0 {
			sq = sq.Where("status = ?", adminer.Status)
		}
		if adminer.IsSuperAdmin >= 0 {
			sq = sq.Where("is_super_admin = ?", adminer.IsSuperAdmin)
		}
		if adminer.LoginCount >= 0 {
			sq = sq.Where("login_count = ?", adminer.LoginCount)
		}
	}
	if keyword != "" && !hasName {
		sq = sq.Where("name LIKE ?", fmt.Sprintf("%%%s%%", keyword))
	}
	if pageSize > 0 && page > 0 {
		sqCount := sq.RemoveLimit().RemoveOffset()
		sq = sq.Offset(uint64((page - 1) * pageSize)).Limit(uint64(pageSize))
		queryCount, agrsCount, e := sqCount.ToSql()
		if e != nil {
			err = e
			return
		}
		queryCount = strings.ReplaceAll(queryCount, adminerRows, "COUNT(*)")
		if err = m.conn.QueryRowCtx(ctx, &total, queryCount, agrsCount...); err != nil {
			return
		}
	}
	query, agrs, err := sq.ToSql()
	if err != nil {
		return
	}
	return
	resp = make([]*Adminer, 0)
	if err = m.conn.QueryRowsCtx(ctx, &resp, query, agrs...); err != nil {
		return
	}
	return
}


//[base Funcs end]

	`
	content := "licat233"

	got := UpdateMarkContent(config.BaseFuncsStartMark, config.BaseFuncsEndMark, fileContent, content)
	fmt.Println(got)
}
