package shorturl

import (
	"context"
	"shorturl/internal/model/short"

	"shorturl/internal/svc"
	"shorturl/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLogic struct {
	logx.Logger
	ctx context.Context
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
	}
}

func (l *UpdateLogic) Update(req *types.UpdateRequest) (resp *types.Response, err error) {
	nameDao := req.DataType
	urlDao := "url"
	if req.ExprieInMinutes < 0 {
		nameDao += "_permanent"
		urlDao += "_permanent"
	}

	if req.Status == 3 {
		err = switchShortDao(l.ctx, nameDao).DeleteByCode(req.Code)
		if err != nil {
			return
		}
		if req.ShortCode != "" {
			_ = switchShortDao(l.ctx, urlDao).DeleteByCode(req.ShortCode)
		}
		return
	}

	exprie, err := getExpire(req.ExprieInMinutes, req.ExprieTime)
	if err != nil {
		return
	}
	s := short.Short{
		Code:     req.Code,
		Data:     req.Data,
		Status:   req.Status,
		CreateBy: req.CreateBy,
		Exprie:   exprie,
	}

	if req.ShortCode != "" {
		err = switchShortDao(l.ctx, urlDao).UpdateByCodeStatus(req.ShortCode, int64(req.Status))
		if err != nil {
			return nil, err
		}
	}

	err = switchShortDao(l.ctx, nameDao).UpdateByCode(req.Code, s)
	if err != nil {
		return
	}
	return &types.Response{}, nil
}
