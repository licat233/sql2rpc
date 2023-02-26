package all_pb_api

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"strings"
)

type (
	adminermodel interface {
		FindList(ctx context.Context, pageSize, page int64, keyword string, adminer *Adminer) (resp []*Adminer, total int64, err error)
	}
	extendAdminerModel struct {
		*defaultAdminerModel
	}
)

func (m *extendAdminerModel) FindList(ctx context.Context, pageSize, page int64, keyword string, adminer *Adminer) (resp []*Adminer, total int64, err error) {
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
	resp = make([]*Adminer, 0)
	if err = m.conn.QueryRowsCtx(ctx, &resp, query, agrs...); err != nil {
		return
	}
	return
}
