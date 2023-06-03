package short

import (
	"github.com/zeromicro/go-zero/core/jsonx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/gorm"
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

func NewPermanentUrlDao(db *gorm.DB, cache *redis.Redis) *PermanentUrlDao {
	return &PermanentUrlDao{shortDao: newShortDao(db, Url{}, cache, Url{}.TableName())}
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
