package short

import (
	"context"
	"github.com/yi-nology/sdk/conf"
	"github.com/zeromicro/go-zero/core/jsonx"
	"github.com/zeromicro/go-zero/core/logx"
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

func NewUrlDao(ctx context.Context) *UrlDao {
	return &UrlDao{shortDao: newShortDao(ctx, conf.MysqlClient, Url{}, conf.RedisClient, Url{}.TableName(), logx.WithContext(ctx))}
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
