package short

import (
	"github.com/zeromicro/go-zero/core/jsonx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/gorm"
)

//CREATE TABLE `urls` (
//`id` int(11) NOT NULL AUTO_INCREMENT,
//`code` varchar(255) NOT NULL COMMENT 'code 短码',
//`data` varchar(512) NOT NULL COMMENT '数据',
//`exprie` datetime DEFAULT NULL COMMENT '过期时间',
//`status` int(2) NOT NULL,
//`create_at` datetime DEFAULT NULL,
//`update_at` datetime DEFAULT NULL ON UPDATE current_timestamp(),
//`delete_at` datetime DEFAULT NULL,
//PRIMARY KEY (`id`)
//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4

type Url struct {
	Short
}

func (Url) TableName() string {
	return "urls"
}

type UrlDao struct {
	*shortDao
}

func NewUrlDao(db *gorm.DB, cache *redis.Redis) *UrlDao {
	return &UrlDao{shortDao: newShortDao(db, Url{}, cache, Url{}.TableName())}
}

func (d *UrlDao) Create(url Short) (err error) {
	json, err := jsonx.MarshalToString(url)
	if err != nil {
		return err
	}
	u := Url{}
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
