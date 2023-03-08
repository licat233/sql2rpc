package svc

import (
	"github.com/licat233/sql2rpc/examples/api_single_file/api/internal/config"
	"github.com/licat233/sql2rpc/examples/api_single_file/api/internal/middleware"
	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config         config.Config
	AuthMiddleware rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:         c,
		AuthMiddleware: middleware.NewAuthMiddleware().Handle,
	}
}
