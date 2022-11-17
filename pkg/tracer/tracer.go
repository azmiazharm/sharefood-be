// Package tracer
// @author Daud Valentino
package tracer

import (
	"context"
	"fmt"

	"github.com/opentracing/opentracing-go"
	otext "github.com/opentracing/opentracing-go/ext"
	otlog "github.com/opentracing/opentracing-go/log"

	"sharefood/pkg/util"
)

type Option struct {
	TagKey string
	Value  string
}

func WithOptions(key, value string) Option {
	return Option{
		TagKey: key,
		Value:  value,
	}
}

func WithResourceNameOptions(value string) Option {
	return Option{
		TagKey: "resource.name",
		Value:  value,
	}
}

// SpanStart starts a new query span from ctx, then returns a new context with the new span.
func SpanStart(ctx context.Context, eventName string) context.Context {
	_, ctx = opentracing.StartSpanFromContext(ctx, eventName)
	return ctx
}

// SpanFinish finishes the span associated with ctx.
func SpanFinish(ctx context.Context) {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span.Finish()
	}
}

// SpanError adds an error to the span associated with ctx.
func SpanError(ctx context.Context, err error) {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		otext.Error.Set(span, true)
		span.LogFields(otlog.String("event", "error"), otlog.Error(err))
	}
}

func SpanStartWithOption(ctx context.Context, eventName string, opts ...Option) context.Context {

	var spOptions []opentracing.StartSpanOption

	for x := 0; x < len(opts); x++ {
		if util.InArray(opts[x].TagKey, []string{
			"span.type",
			"service.name",
			"resource.name",
		}) {
			spOptions = append(spOptions, opentracing.Tag{Key: opts[x].TagKey, Value: opts[x].Value})
		}
	}

	sp, ctx := opentracing.StartSpanFromContext(ctx, eventName, spOptions...)
	for i := 0; i < len(opts); i++ {
		sp.SetTag(opts[i].TagKey, opts[i].Value)
	}

	return ctx
}

func DBSpanStartWithOption(ctx context.Context, dbName, eventName string, opts ...Option) context.Context {

	svcName := fmt.Sprintf("%s.%s", "mysql", dbName)
	opts = append(opts,
		WithOptions("db.type", "sql"),
		//WithOptions("span.type", "db"),
		WithOptions("span.kind", "clients"),
		WithOptions("peer.service", svcName),
		WithOptions("service.name", svcName),
	)

	return SpanStartWithOption(ctx, eventName, opts...)
}
