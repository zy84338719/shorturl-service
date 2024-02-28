package shorturl

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"shorturl/internal/logic/shorturl"
	"shorturl/internal/svc"
	"shorturl/internal/types"
)

func CreateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		if req.DataType == "url" && !req.ProhibitEffective && checkUrlEffective(req.Data) {
			httpx.OkJsonCtx(r.Context(), w, types.Response{Code: 1, Msg: "url禁止访问"})
		}

		l := shorturl.NewCreateLogic(r.Context(), svcCtx)
		resp, err := l.Create(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
