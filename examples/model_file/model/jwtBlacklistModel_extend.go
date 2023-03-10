// Code generated by sql2rpc. DO NOT EDIT.

package model

import (
	"context"
	"strings"

	"github.com/Masterminds/squirrel"
)

type (
	jwtBlacklist_model interface {
		FindList(ctx context.Context, pageSize, page int64, keyword string, jwtBlacklist *JwtBlacklist) (resp []*JwtBlacklist, total int64, err error)
		TableName() string
	}
)

func (m *defaultJwtBlacklistModel) FindList(ctx context.Context, pageSize, page int64, keyword string, jwtBlacklist *JwtBlacklist) (resp []*JwtBlacklist, total int64, err error) {
	sq := squirrel.Select(jwtBlacklistRows).From(m.table)
	if jwtBlacklist != nil {
		if jwtBlacklist.Id > 0 {
			sq = sq.Where("id = ?", jwtBlacklist.Id)
		}
		if jwtBlacklist.AdminerId > 0 {
			sq = sq.Where("adminer_id = ?", jwtBlacklist.AdminerId)
		}
		if jwtBlacklist.Uuid != "" {
			sq = sq.Where("uuid = ?", jwtBlacklist.Uuid)
		}
		if jwtBlacklist.Token != "" {
			sq = sq.Where("token = ?", jwtBlacklist.Token)
		}
		if jwtBlacklist.Platform != "" {
			sq = sq.Where("platform = ?", jwtBlacklist.Platform)
		}
		if jwtBlacklist.Ip != "" {
			sq = sq.Where("ip = ?", jwtBlacklist.Ip)
		}
	}
	if pageSize > 0 && page > 0 {
		sqCount := sq.RemoveLimit().RemoveOffset()
		sq = sq.Limit(uint64(pageSize)).Offset(uint64((page - 1) * pageSize))
		queryCount, agrsCount, e := sqCount.ToSql()
		if e != nil {
			err = e
			return
		}
		queryCount = strings.ReplaceAll(queryCount, jwtBlacklistRows, "COUNT(*)")
		if err = m.conn.QueryRowCtx(ctx, &total, queryCount, agrsCount...); err != nil {
			return
		}
	}
	query, agrs, err := sq.ToSql()
	if err != nil {
		return
	}
	resp = make([]*JwtBlacklist, 0)
	if err = m.conn.QueryRowsCtx(ctx, &resp, query, agrs...); err != nil {
		return
	}
	return
}

func (m *defaultJwtBlacklistModel) TableName() string {
	return m.table
}
