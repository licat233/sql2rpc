package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ AdminerModel = (*customAdminerModel)(nil)

type (
	// AdminerModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAdminerModel.
	AdminerModel interface {
		adminermodel
		adminerModel
	}

	customAdminerModel struct {
		*extendAdminerModel
		*defaultAdminerModel
	}
)

// NewAdminerModel returns a model for the database table.
func NewAdminerModel(conn sqlx.SqlConn) AdminerModel {
	m := newAdminerModel(conn)
	return &customAdminerModel{
		newExtendAdminerModelModel(m),
		m,
	}
}
