package bootstrap

import (
	"fmt"

	"github.com/opentracing/opentracing-go"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/opentracer"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"sharefood/internal/appctx"
	"sharefood/pkg/logger"
)

// RegistryOpenTracing setup
func RegistryOpenTracing(cfg *appctx.Config) opentracing.Tracer {

	if !cfg.APM.Enable {
		return opentracing.NoopTracer{}
	}

	lf := logger.NewFields(logger.EventName("TracerInitiated"))
	logger.Debug(fmt.Sprint("apm address : ", cfg.APM.Address), lf...)
	tr := opentracer.New(
		tracer.WithAgentAddr(cfg.APM.Address),
		tracer.WithService(cfg.APM.Name),
		tracer.WithGlobalTag("env", cfg.App.Env),
	)

	opentracing.SetGlobalTracer(tr)
	return tr

}
