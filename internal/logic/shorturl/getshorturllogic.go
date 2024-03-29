package shorturl

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"shorturl/internal/model/short"
	"shorturl/internal/svc"
	"shorturl/internal/types"
)

type GetShorturlLogic struct {
	logx.Logger
	ctx    context.Context
	domain string
}

func NewGetShorturlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetShorturlLogic {
	return &GetShorturlLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		domain: svcCtx.Short.Domain,
	}
}

func (l *GetShorturlLogic) GetShorturl(req *types.GetShortRequest) (s string, err error) {
	url := &short.Short{}
	nameDao := "url"
	if judgePermanentShort(req.Code) {
		nameDao += "_permanent"
	}
	url, err = switchShortDao(l.ctx, nameDao).GetShortUrl(req.Code)
	if err != nil {
		return l.domain + "/errors.html", nil
	}
	err = checkExprieTime(url.Exprie, url.Status)
	if err != nil {
		return l.domain + "/errors.html", nil
	}
	addCount(l.ctx, req.Code, nameDao)
	return url.Data, nil
}
