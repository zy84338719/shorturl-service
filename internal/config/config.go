package config

import (
	"github.com/yi-nology/sdk/conf"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	conf.Config
	Short Short `yaml:"short" json:"short"`
}

type Short struct {
	Domain string `yaml:"domain" json:"domain"`
}
