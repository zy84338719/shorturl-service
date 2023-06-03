package shorturl

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/gorm"
	"shorturl/internal/model/short"

	"shorturl/internal/svc"
	"shorturl/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLogic struct {
	logx.Logger
	ctx   context.Context
	cache *redis.Redis
	db    *gorm.DB
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		cache:  svcCtx.Cache,
		db:     svcCtx.DB,
	}
}

func (l *UpdateLogic) Update(req *types.UpdateRequest) (resp *types.Response, err error) {
	nameDao := req.DataType
	permanent := req.ExprieInMinutes < 0
	if permanent {
		nameDao += "_permanent"
	}

	if req.Status == 3 {
		err = switchShortDao(l.db, l.cache, nameDao).DeleteByCode(req.Code)
		if err != nil {
			return
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

	err = switchShortDao(l.db, l.cache, nameDao).UpdateByCode(req.Code, s)
	if err != nil {
		return
	}
	return
}
