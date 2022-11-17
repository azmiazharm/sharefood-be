// Package {{.PackageName}}
// Automatic generated
package {{.PackageName}}

import (
	"fmt"

	"github.com/gookit/validate"

	"sharefood/internal/appctx"
	"sharefood/internal/consts"
	"sharefood/internal/presentations"
	"sharefood/internal/repositories"
	"sharefood/pkg/logger"
	"sharefood/pkg/tracer"

	ucase "sharefood/internal/ucase/contract"
)

type {{.StructName}} struct {
	repo repositories.{{.RepoContractName}}
}

// New{{.EntityName}} new instance
func New{{.EntityName}}(repo repositories.{{.RepoContractName}}) ucase.UseCase {
	return &{{.StructName}}{repo: repo}
}

// Serve store {{.StructName}} data
func (u *{{.StructName}}) Serve(dctx *appctx.Data) appctx.Response {
	var (
		param presentations.{{.EntityName}}Param
		ctx   = tracer.SpanStart(dctx.Request.Context(), "ucase.create")
		lf    = logger.NewFields(
			logger.EventName("{{.EntityName}}"),
		)
	)

	defer tracer.SpanFinish(ctx)

	v := validate.Request(dctx.Request)

	// example rule
	// v.AddRule("price", "min", 1)
	// v.StringRule("product_id", `required|minLen:7`)
	// v.StringRule("client_id", `required|minLen:7`)
	// v.StringRule("status", `required|minLen:3`)

    // validate not ok
	if !v.Validate() {
		logger.WarnWithContext(ctx, fmt.Sprintf("validation error"), lf...)
		return *appctx.NewResponse().WithMsgKey(consts.RespValidationError).WithError(v.Errors)
	}

    // cast data into struct
	err := v.BindSafeData(&param)
	if err != nil {
	    tracer.SpanError(ctx, err)
		logger.WarnWithContext(ctx, fmt.Sprintf("error parsing request param: %v", err), lf...)
		return *appctx.NewResponse().WithMsgKey(consts.RespValidationError)
	}

	af, err := u.repo.Store(ctx, param)
	if err != nil {
	    tracer.SpanError(ctx, err)
		logger.WarnWithContext(ctx, fmt.Sprintf("store data to database error: %v", err), lf...)
		return *appctx.NewResponse().WithMsgKey(consts.RespError)
	}

	lf.Append(logger.Any("affected_rows", af))

	logger.InfoWithContext(ctx, fmt.Sprintf("success store data to database"), lf...)
	return *appctx.NewResponse().
		WithMsgKey(consts.RespSuccess)
}
