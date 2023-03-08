package jwtBlacklist

import (
	"context"

	"github.com/licat233/sql2rpc/examples/api_single_file/api/internal/svc"
	"github.com/licat233/sql2rpc/examples/api_single_file/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetJwtBlacklistListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetJwtBlacklistListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetJwtBlacklistListLogic {
	return &GetJwtBlacklistListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetJwtBlacklistListLogic) GetJwtBlacklistList(req *types.GetJwtBlacklistListReq) (resp *types.BaseResp, err error) {
	// todo: add your logic here and delete this line

	return
}
