// Package bootstrap
package bootstrap

import (
	"context"
	"crypto/tls"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"

	"sharefood/internal/appctx"
	"sharefood/pkg/cache"
	"sharefood/pkg/logger"
)

const (
	redisInitializeNil         = `redis cannot connect, please check your config or network`
	redisPingError             = `redis cannot connect, error: %v`
	logFieldHost               = "host"
	logFieldDB                 = "db"
	logFieldMaxRedirect        = "max_redirect"
	logFieldReadOnly           = "read_only"
	logFieldTlsEnable          = "tls_enable"
	logFieldInsecureSKipVerify = "insecure_skip_verify"
	logFieldRouteRandomly      = "route_randomly"
	logFieldRouteByLatency     = "route_by_latency"
)

// RegistryRedisNative initiate redis session
func RegistryRedisNative(conf *appctx.Config) redis.Cmdable {

	lf := []logger.Field{
		logger.Any(logFieldHost, conf.Redis.Hosts),
		logger.Any(logFieldDB, conf.Redis.DB),
		logger.Any(logFieldReadOnly, conf.Redis.ReadOnly),
		logger.Any(logFieldRouteByLatency, conf.Redis.RouteByLatency),
		logger.Any(logFieldRouteRandomly, conf.Redis.RouteRandomly),
		logger.Any(logFieldMaxRedirect, conf.Redis.MaxRedirect),
		logger.Any(logFieldTlsEnable, conf.Redis.TLSEnable),
		logger.Any(logFieldInsecureSKipVerify, conf.Redis.InsecureSkipVerify),
	}

	if conf.Redis.ClusterMode {
		return registryRedisCluster(conf, lf)
	}

	return registryRedisUniversal(conf, lf)
}

// registryRedisUniversal initiate redis session
func registryRedisUniversal(conf *appctx.Config, lf []logger.Field) redis.Cmdable {
	cfg := redis.UniversalOptions{
		Addrs:          strings.Split(conf.Redis.Hosts, ","),
		ReadTimeout:    time.Duration(conf.Redis.ReadTimeoutSecond) * time.Second,
		WriteTimeout:   time.Duration(conf.Redis.WriteTimeoutSecond) * time.Second,
		DB:             conf.Redis.DB,
		PoolSize:       conf.Redis.PoolSize,
		PoolTimeout:    time.Duration(conf.Redis.PoolTimeoutSecond) * time.Second,
		MinIdleConns:   conf.Redis.MinIdleConn,
		IdleTimeout:    time.Duration(conf.Redis.IdleTimeoutSecond) * time.Second,
		RouteByLatency: conf.Redis.RouteByLatency,
		Password:       conf.Redis.Password,
		ReadOnly:       conf.Redis.ReadOnly,
		RouteRandomly:  conf.Redis.RouteRandomly,
		MaxRedirects:   conf.Redis.MaxRedirect,
	}

	if conf.Redis.TLSEnable {
		cfg.TLSConfig = &tls.Config{
			InsecureSkipVerify: conf.Redis.InsecureSkipVerify,
		}
	}

	r := redis.NewUniversalClient(&cfg)
	r.AddHook(cache.NewRedisHook(cfg.Addrs, 0))

	if r == nil {
		logger.Fatal(redisInitializeNil, lf...)
	}

	c := r.Ping(context.Background())

	if err := c.Err(); err != nil {
		logger.Fatal(fmt.Sprintf(redisPingError, err), lf...)
	}

	return r
}

// registryRedisCluster initiate redis session
func registryRedisCluster(conf *appctx.Config, lf []logger.Field) redis.Cmdable {

	cfg := &redis.ClusterOptions{
		Addrs:          strings.Split(conf.Redis.Hosts, ","),
		ReadTimeout:    time.Duration(conf.Redis.ReadTimeoutSecond) * time.Second,
		WriteTimeout:   time.Duration(conf.Redis.WriteTimeoutSecond) * time.Second,
		PoolSize:       conf.Redis.PoolSize,
		PoolTimeout:    time.Duration(conf.Redis.PoolTimeoutSecond) * time.Second,
		MinIdleConns:   conf.Redis.MinIdleConn,
		IdleTimeout:    time.Duration(conf.Redis.IdleTimeoutSecond) * time.Second,
		RouteByLatency: conf.Redis.RouteByLatency,
		Password:       conf.Redis.Password,
		ReadOnly:       conf.Redis.ReadOnly,
		RouteRandomly:  conf.Redis.RouteRandomly,
		MaxRedirects:   conf.Redis.MaxRedirect,
	}

	if conf.Redis.TLSEnable {
		cfg.TLSConfig = &tls.Config{
			InsecureSkipVerify: conf.Redis.InsecureSkipVerify,
		}
	}

	r := redis.NewClusterClient(cfg)
	r.AddHook(cache.NewRedisHook(cfg.Addrs, 0))

	if r == nil {
		logger.Fatal(redisInitializeNil, lf...)
	}

	c := r.Ping(context.Background())

	if err := c.Err(); err != nil {
		logger.Fatal(fmt.Sprintf(redisPingError, err), lf...)
	}

	return r
}
