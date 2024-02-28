package svc

import (
	"shorturl/internal/config"
)

type ServiceContext struct {
	Config config.Config
	Short  config.Short
}

func NewServiceContext(c config.Config) *ServiceContext {

	return &ServiceContext{
		Config: c,
		Short:  c.Short,
	}
}
