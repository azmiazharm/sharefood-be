// Package example
package example

import (
	"sharefood/internal/appctx"
	"sharefood/internal/consts"
	"sharefood/internal/repositories"
	"sharefood/internal/ucase/contract"

	"sharefood/pkg/logger"
)

type exampleList struct {
	repo repositories.Example
}

func NewExampleList(repo repositories.Example) contract.UseCase {
	return &exampleList{repo: repo}
}

// Serve partner list data
func (u *exampleList) Serve(data *appctx.Data) appctx.Response {

	p, err := u.repo.Find(data.Request.Context())

	if err != nil {
		logger.Error(logger.MessageFormat("[example-list] %v", err))

		return *appctx.NewResponse().WithCode(consts.CodeInternalServerError)
	}

	return *appctx.NewResponse().WithCode(consts.CodeSuccess).WithData(p)
}
