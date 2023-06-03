package shorturl

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/gorm"
	"shorturl/internal/svc"
	"shorturl/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetListLogic struct {
	logx.Logger
	ctx   context.Context
	db    *gorm.DB
	cache *redis.Redis
}

func NewGetListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetListLogic {
	return &GetListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		db:     svcCtx.DB,
		cache:  svcCtx.Cache,
	}
}

func (l *GetListLogic) GetList(req *types.GetListRequest) (resp *types.Response, err error) {
	nameDao := req.DataType
	if req.Permanent {
		nameDao += "_permanent"
	}
	url, count, err := switchShortDao(l.db, l.cache, nameDao).GetList(req.Page, req.PageSize, req.ShortCode)
	if err != nil {
		return &types.Response{Code: 1, Msg: err.Error()}, nil
	}
	return &types.Response{Data: map[string]interface{}{
		"list":  url,
		"total": count,
	}}, nil
}
