package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"shorturl/internal/config"
)

type ServiceContext struct {
	Config config.Config
	Cache  *redis.Redis
	DB     *gorm.DB
	Short  config.Short
}

func NewServiceContext(c config.Config, cc config.Conf) *ServiceContext {
	mysqlConfig := mysql.Config{
		DSN:                       cc.Mysql.DataSource, // DSN data source name
		DefaultStringSize:         191,                 // string 类型字段的默认长度
		SkipInitializeWithVersion: false,               // 根据版本自动配置

	}
	db, err := gorm.Open(mysql.New(mysqlConfig))
	if err != nil {
		panic(err)
	}
	cache := redis.MustNewRedis(cc.Redis)

	return &ServiceContext{
		Config: c,
		Cache:  cache,
		DB:     db,
		Short:  cc.Short,
	}
}
