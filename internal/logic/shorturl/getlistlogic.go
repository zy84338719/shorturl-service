package shorturl

import (
	"context"
	"shorturl/internal/svc"
	"shorturl/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetListLogic struct {
	logx.Logger
	ctx context.Context
}

func NewGetListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetListLogic {
	return &GetListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
	}
}

func (l *GetListLogic) GetList(req *types.GetListRequest) (resp *types.Response, err error) {
	nameDao := req.DataType
	if req.Permanent {
		nameDao += "_permanent"
	}
	url, count, err := switchShortDao(l.ctx, nameDao).GetList(req.Page, req.PageSize, req.ShortCode, req.CreateBy)
	if err != nil {
		return &types.Response{Code: 1, Msg: err.Error()}, nil
	}
	return &types.Response{Data: map[string]interface{}{
		"list":  url,
		"total": count,
	}}, nil
}
