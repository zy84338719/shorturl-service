package short

import (
	"github.com/zeromicro/go-zero/core/jsonx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/gorm"
)

//CREATE TABLE `data` (
//`id` int(11) NOT NULL AUTO_INCREMENT,
//`code` varchar(255) NOT NULL COMMENT 'code 短码',
//`data` varchar(4096) NOT NULL COMMENT '数据',
//`exprie` datetime DEFAULT NULL COMMENT '过期时间',
//`status` int(2) NOT NULL,
//`create_at` datetime DEFAULT NULL,
//`update_at` datetime DEFAULT NULL ON UPDATE current_timestamp(),
//`delete_at` datetime DEFAULT NULL,
//PRIMARY KEY (`id`)
//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4

// 我想将上面的表 转换为 golang 带有gorm标签的结构体
// 1. 读取表结构

type Data struct {
	Short
}

func (Data) TableName() string {
	return "data"
}

type DataDao struct {
	*shortDao
}

func NewDataDao(db *gorm.DB, cache *redis.Redis) *DataDao {
	return &DataDao{shortDao: newShortDao(db, Data{}, cache, Data{}.TableName())}
}
func (d *DataDao) Create(url Short) (err error) {
	json, err := jsonx.MarshalToString(url)
	if err != nil {
		return err
	}
	u := Data{}
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
