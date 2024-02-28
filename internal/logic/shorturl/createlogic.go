package shorturl

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"shorturl/internal/model/short"
	"shorturl/internal/svc"
	"shorturl/internal/types"
)

type CreateLogic struct {
	logx.Logger
	ctx context.Context
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {

	return &CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
	}
}

func (l *CreateLogic) Create(req *types.CreateRequest) (resp *types.Response, err error) {
	// 成码规则
	// 共计7位
	// 前6位位hash计算 第七位为分库分表预留扩展字段
	exprie, err := getExpire(req.ExprieInMinutes, req.ExprieTime)
	if err != nil {
		return
	}
	permanent := req.ExprieInMinutes < 0
	nameDao := req.DataType
	if permanent {
		nameDao += "_permanent"
	}

	s := short.Short{
		Exprie:   exprie,
		Data:     req.Data,
		Status:   1,
		CreateBy: req.CreateBy,
	}
	dao := switchShortDao(l.ctx, nameDao)
	for i := 0; i < 10; i++ {
		if i > 0 {
			url, err := dao.GetShortUrl(s.Code)
			if err != nil {
				return &types.Response{Code: 1, Msg: err.Error()}, nil
			}
			if l.compare(&s, url) {
				return &types.Response{Code: 1, Msg: "数据重复", Data: url.Code}, nil
			}
		}

		if !permanent {
			s.Code = getTempShortCode(req.DataType, req.Data, i)
		} else {
			s.Code = getShortCode(req.DataType, req.Data, i)
		}

		err = dao.Create(s)
		if err != nil {
			continue
		}
		break
	}
	return &types.Response{Data: s.Code}, nil
}

// Compare this snippet from drama_stick/shorturl/internal/model/short/short.go:
func (l *CreateLogic) compare(url1 *short.Short, url2 *short.Short) bool {
	if url1.Data != url2.Data {
		return false
	}
	return true
}
