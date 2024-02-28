package shorturl

import (
	"context"
	"shorturl/internal/svc"
	"shorturl/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDataLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDataLogic {
	return &GetDataLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDataLogic) GetData(req *types.GetShortRequest) (resp *types.Response, err error) {

	nameDao := "data"
	if judgePermanentShort(req.Code) {
		nameDao += "_permanent"
	}
	url, err := switchShortDao(l.ctx, nameDao).GetShortUrl(req.Code)
	if err != nil {
		return &types.Response{Code: 1, Msg: err.Error()}, nil
	}
	err = checkExprieTime(url.Exprie, url.Status)
	if err != nil {
		return &types.Response{Code: 1, Msg: err.Error()}, nil
	}
	addCount(l.ctx, req.Code, nameDao)
	return &types.Response{Data: url.Data}, nil
}
