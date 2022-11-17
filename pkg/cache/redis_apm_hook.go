// Package cache
package cache

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
)

type redisHook struct {
	addrs []string
	db    int
}

func NewRedisHook(addrs []string, db int) redis.Hook {
	return &redisHook{
		addrs: addrs,
		db:    db,
	}
}

func (h *redisHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {

	span, newCtx := opentracing.StartSpanFromContext(ctx, fmt.Sprintf("redis.%s", cmd.Name()))
	ext.DBType.Set(span, "cache")
	ext.DBInstance.Set(span, strconv.Itoa(h.db))
	ext.PeerAddress.Set(span, strings.Join(h.addrs, ", "))
	ext.PeerService.Set(span, "redis")
	ext.SpanKind.Set(span, ext.SpanKindEnum("client"))
	ext.DBStatement.Set(span, strings.ToUpper(cmd.Name()))
	span.SetTag("service.name", "redis")
	return newCtx, nil

}

func (h *redisHook) AfterProcess(ctx context.Context, _ redis.Cmder) error {
	span := opentracing.SpanFromContext(ctx)

	if span != nil {
		// if context is raised an error.
		if ctx.Err() != nil {
			ext.Error.Set(span, true)
			span.LogFields(log.Error(ctx.Err()))
		}
		span.Finish()
	}

	return nil
}

func (h *redisHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	span, newCtx := opentracing.StartSpanFromContext(ctx, "redis:pipeline:cmd")
	ext.DBType.Set(span, "cache")
	ext.DBInstance.Set(span, strconv.Itoa(h.db))
	ext.PeerAddress.Set(span, strings.Join(h.addrs, ", "))
	ext.PeerService.Set(span, "redis")
	ext.SpanKind.Set(span, ext.SpanKindEnum("client"))
	merge := make([]string, len(cmds))
	for i, cmd := range cmds {
		merge[i] = strings.ToUpper(cmd.Name())
	}
	ext.DBStatement.Set(span, strings.Join(merge, " --> "))
	span.SetTag("service.name", "redis")
	return newCtx, nil
}

func (h *redisHook) AfterProcessPipeline(ctx context.Context, _ []redis.Cmder) error {
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		// if context is raised an error.
		if ctx.Err() != nil {
			ext.Error.Set(span, true)
			span.LogFields(log.Error(ctx.Err()))
		}
		span.Finish()
	}
	return nil
}
