package jwtBlacklist

import (
	"context"

	"github.com/licat233/sql2rpc/examples/api_single_file/api/internal/svc"
	"github.com/licat233/sql2rpc/examples/api_single_file/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelJwtBlacklistLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDelJwtBlacklistLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelJwtBlacklistLogic {
	return &DelJwtBlacklistLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DelJwtBlacklistLogic) DelJwtBlacklist(req *types.DelJwtBlacklistReq) (resp *types.BaseResp, err error) {
	// todo: add your logic here and delete this line

	return
}
