package short

import (
	"github.com/zeromicro/go-zero/core/jsonx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/gorm"
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

func NewPermanentDataDao(db *gorm.DB, cache *redis.Redis) *PermanentDataDao {
	return &PermanentDataDao{shortDao: newShortDao(db, PermanentData{}, cache, PermanentData{}.TableName())}
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
