package svc

import (
	"github.com/redis/go-redis/v9"
	"github.com/stumble/dcache"
	"github.com/stumble/wpgx"

	"{{ cookiecutter.go_module_path.strip('/') }}/internal/config"
	"{{ cookiecutter.go_module_path.strip('/') }}/internal/repos/users"
)

type Service struct {
	wpgxPool *wpgx.Pool
	cache    *dcache.DCache
	redisCli redis.UniversalClient
	config   *config.Config

	usersRepo *users.Queries
}

func NewService(
	wpgxPool *wpgx.Pool,
	cache *dcache.DCache,
	redisCli redis.UniversalClient,
	config *config.Config,
) *Service {
	return &Service{
		wpgxPool:  wpgxPool,
		cache:     cache,
		redisCli:  redisCli,
		config:    config,
		usersRepo: users.New(wpgxPool.WConn(), cache),
	}
}
