package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ FanModel = (*customFanModel)(nil)

type (
	// FanModel is an interface to be customized, add more methods here,
	// and implement the added methods in customFanModel.
	FanModel interface {
		fanModel
		withSession(session sqlx.Session) FanModel
	}

	customFanModel struct {
		*defaultFanModel
	}
)

// NewFanModel returns a model for the database table.
func NewFanModel(conn sqlx.SqlConn) FanModel {
	return &customFanModel{
		defaultFanModel: newFanModel(conn),
	}
}

func (m *customFanModel) withSession(session sqlx.Session) FanModel {
	return NewFanModel(sqlx.NewSqlConnFromSession(session))
}
