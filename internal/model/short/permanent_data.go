package short

import (
	"context"
	"github.com/yi-nology/sdk/conf"
	"github.com/zeromicro/go-zero/core/jsonx"
	"github.com/zeromicro/go-zero/core/logx"
)

type PermanentData struct {
	Short
}

func (PermanentData) TableName() string {
	return "permanent_data"
}

type PermanentDataDao struct {
	*shortDao
}

func NewPermanentDataDao(ctx context.Context) *PermanentDataDao {
	return &PermanentDataDao{shortDao: newShortDao(ctx, conf.MysqlClient, PermanentData{}, conf.RedisClient, PermanentData{}.TableName(), logx.WithContext(ctx))}
}

func (d *PermanentDataDao) Create(url Short) (err error) {
	json, err := jsonx.MarshalToString(url)
	if err != nil {
		return err
	}
	u := PermanentData{}
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
