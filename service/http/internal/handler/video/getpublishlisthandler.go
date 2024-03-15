package video

import (
	"net/http"

	"douyin/service/http/internal/logic/video"
	"douyin/service/http/internal/svc"
	"douyin/service/http/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetPublishListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PublishListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := video.NewGetPublishListLogic(r.Context(), svcCtx)
		resp, err := l.GetPublishList(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
