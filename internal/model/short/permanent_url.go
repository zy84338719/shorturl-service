package short

import (
	"context"
	"github.com/yi-nology/sdk/conf"
	"github.com/zeromicro/go-zero/core/jsonx"
	"github.com/zeromicro/go-zero/core/logx"
)

type PermanentUrl struct {
	Short
}

func (PermanentUrl) TableName() string {
	return "permanent_urls"
}

type PermanentUrlDao struct {
	*shortDao
}

func NewPermanentUrlDao(ctx context.Context) *PermanentUrlDao {
	return &PermanentUrlDao{shortDao: newShortDao(ctx, conf.MysqlClient, PermanentUrl{}, conf.RedisClient, PermanentUrl{}.TableName(), logx.WithContext(ctx))}
}

func (d *PermanentUrlDao) Create(url Short) (err error) {
	json, err := jsonx.MarshalToString(url)
	if err != nil {
		return err
	}
	u := PermanentUrl{}
	err = jsonx.UnmarshalFromString(json, &u)
	if err != nil {
		return err
	}
	err = d.db.Model(d.model).Create(&u).Error
	if err != nil {
		return err
	}
	return
}
