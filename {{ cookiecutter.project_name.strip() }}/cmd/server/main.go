package main

import (
	"context"

	ggapp "github.com/ggsrc/gglib/app"
	gggrpc "github.com/ggsrc/gglib/grpc"
	ggresource "github.com/ggsrc/gglib/resource"
	ggcache "github.com/ggsrc/gglib/resource/cache"
	ggwpgx "github.com/ggsrc/gglib/resource/wpgx"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
	"github.com/stumble/wpgx"
	"google.golang.org/grpc/reflection"

	pb "{{ cookiecutter.go_module_path.strip('/') }}/api/go/{{ cookiecutter.project_name.strip() }}/v1"
	"{{ cookiecutter.go_module_path.strip('/') }}/internal/config"
	"{{ cookiecutter.go_module_path.strip('/') }}/internal/svc"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		log.Fatal().Err(err).Msg("{{ cookiecutter.project_name.strip() }} stopped")
	}
}

// run follows gglib/app's lifecycle: New answers the health probes (alive,
// not ready) from the start, Init waits for Postgres and Redis however long
// they take, and Run serves until SIGTERM, then drains and stops everything.
func run(ctx context.Context) error {
	var c config.Config
	envconfig.MustProcess("{{ cookiecutter.project_name.strip().replace("-", "_") }}", &c)

	a, err := ggapp.New(ctx,
		ggapp.WithServerName(c.AppName),
		ggapp.WithDebug(c.Debug),
	)
	if err != nil {
		return err
	}

	// Register dependencies before the resources that use them: resources
	// stop in reverse order.
	wpgx := newWPGX()
	cache := newCache(c.AppName)
	resourceManager := ggresource.NewResourceManager([]ggresource.Resource{wpgx, cache})
	if err := a.Init(ctx, resourceManager); err != nil {
		return err
	}

	wpgxCache := cache.GetDCache()
	// set max memory cache ttl to offload traffic from redis.
	if err := wpgxCache.SetMemCacheMaxTTLSeconds(c.MemCacheMaxTTLSeconds); err != nil {
		return err
	}

	{{ cookiecutter.project_name.strip().replace("-", "_") }} := svc.NewService(wpgx.GetPool(), wpgxCache, cache.GetRedisClient(), &c)

	grpcServer := newGRPCServer()
	reflection.Register(grpcServer.Server())
	pb.RegisterGServiceServer(grpcServer.Server(), {{ cookiecutter.project_name.strip().replace("-", "_") }})

	return a.Run(ctx, grpcServer, resourceManager)
}

func newWPGX() *ggwpgx.WPGX {
	var wpgxConfig wpgx.Config
	envconfig.MustProcess("postgres", &wpgxConfig)
	return ggwpgx.NewWPGXWithOptions(
		ggwpgx.WithConfig(&wpgxConfig),
	)
}

func newCache(appName string) *ggcache.Cache {
	var redisConfig ggcache.RedisConfig
	envconfig.MustProcess("redis", &redisConfig)
	var dcacheConfig ggcache.DCacheConfig
	envconfig.MustProcess("dcache", &dcacheConfig)
	return ggcache.NewCacheWithOptions(
		ggcache.WithAppName(appName),
		ggcache.WithRedisConfig(&redisConfig),
		ggcache.WithDCacheConfig(&dcacheConfig),
	)
}

func newGRPCServer() *gggrpc.Server {
	var grpcConfig gggrpc.ServerConfig
	envconfig.MustProcess("grpc", &grpcConfig)
	return gggrpc.NewServerWithOptions(
		gggrpc.WithServerVerbose(grpcConfig.Verbose),
		gggrpc.WithServerDebug(grpcConfig.Debug),
		gggrpc.WithServerPort(grpcConfig.Port),
	)
}
