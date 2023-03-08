package adminer

import (
	"net/http"

	"github.com/licat233/sql2rpc/examples/api_single_file/api/internal/logic/adminer"
	"github.com/licat233/sql2rpc/examples/api_single_file/api/internal/svc"
	"github.com/licat233/sql2rpc/examples/api_single_file/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetAdminerHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetAdminerReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := adminer.NewGetAdminerLogic(r.Context(), svcCtx)
		resp, err := l.GetAdminer(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
