package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
}

type Conf struct {
	Mysql Mysql           `yaml:"mysql" json:"mysql"`
	Redis redis.RedisConf `yaml:"redis" json:"redis"`
	Short Short           `yaml:"short" json:"short"`
}

type Mysql struct {
	DataSource string `yaml:"dataSource" json:"dataSource"`
}

type Short struct {
	Domain string `yaml:"domain" json:"domain"`
}
