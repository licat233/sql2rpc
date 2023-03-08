package jwtBlacklist

import (
	"context"

	"github.com/licat233/sql2rpc/examples/api_single_file/api/internal/svc"
	"github.com/licat233/sql2rpc/examples/api_single_file/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddJwtBlacklistLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddJwtBlacklistLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddJwtBlacklistLogic {
	return &AddJwtBlacklistLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddJwtBlacklistLogic) AddJwtBlacklist(req *types.AddJwtBlacklistReq) (resp *types.BaseResp, err error) {
	// todo: add your logic here and delete this line

	return
}
