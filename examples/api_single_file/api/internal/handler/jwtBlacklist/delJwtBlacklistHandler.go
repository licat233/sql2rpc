package jwtBlacklist

import (
	"net/http"

	"github.com/licat233/sql2rpc/examples/api_single_file/api/internal/logic/jwtBlacklist"
	"github.com/licat233/sql2rpc/examples/api_single_file/api/internal/svc"
	"github.com/licat233/sql2rpc/examples/api_single_file/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func DelJwtBlacklistHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DelJwtBlacklistReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := jwtBlacklist.NewDelJwtBlacklistLogic(r.Context(), svcCtx)
		resp, err := l.DelJwtBlacklist(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
