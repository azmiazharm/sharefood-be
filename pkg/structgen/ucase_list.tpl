// Package {{.PackageName}}
// Automatic generated
package {{.PackageName}}

import (
	"fmt"

	"sharefood/internal/appctx"
	"sharefood/internal/common"
	"sharefood/internal/consts"
	"sharefood/internal/presentations"
	"sharefood/internal/repositories"
	"sharefood/pkg/logger"
	"sharefood/pkg/tracer"

	ucase "sharefood/internal/ucase/contract"
)

type {{.StructName}}List struct {
	repo repositories.{{.EntityName}}er
}

func New{{.EntityName}}List(repo repositories.{{.EntityName}}er) ucase.UseCase {
	return &{{.StructName}}List{repo: repo}
}

// Serve {{.EntityName}} list data
func (u *{{.StructName}}List) Serve(dctx *appctx.Data) appctx.Response {
	var (
		param presentations.{{.EntityName}}Query
		ctx   = tracer.SpanStart(dctx.Request.Context(), "ucase.{{.FileName}}_list")
		lf    = logger.NewFields(
			logger.EventName("{{.StructName}}List"),
		)
	)
    defer tracer.SpanFinish(ctx)

	err := dctx.Cast(&param)
	if err != nil {
		logger.WarnWithContext(ctx, fmt.Sprintf("error parsing query url: %v", err), lf...)
		return *appctx.NewResponse().WithMsgKey(consts.RespValidationError)
	}

	param.Limit = common.LimitDefaultValue(param.Limit)
	param.Page = common.PageDefaultValue(param.Page)

	p, count, err := u.repo.FindWithCount(ctx, param)
	if err != nil {
	    tracer.SpanError(ctx, err)
		logger.ErrorWithContext(ctx, fmt.Sprintf("error find data to database: %v", err), lf...)
		return *appctx.NewResponse().WithMsgKey(consts.RespError)
	}

	logger.InfoWithContext(ctx, fmt.Sprintf("success fetch {{.TableName}} to database"), lf...)
	return *appctx.NewResponse().
            WithMsgKey(consts.RespSuccess).
            WithData(p).
            WithMeta(appctx.MetaData{
                    Page:       param.Page,
                    Limit:      param.Limit,
                    TotalCount: count,
                    TotalPage:  common.PageCalculate(count, param.Limit),
            })
}