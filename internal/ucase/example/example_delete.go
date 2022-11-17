// Package example
package example

import (
	"github.com/gorilla/mux"
	"github.com/spf13/cast"

	"sharefood/internal/appctx"
	"sharefood/internal/consts"
	"sharefood/internal/repositories"
	"sharefood/internal/ucase/contract"

	"sharefood/pkg/logger"
)

type exampleDelete struct {
	repo repositories.Example
}

func NewExampleDelete(repo repositories.Example) contract.UseCase {
	return &exampleDelete{repo: repo}
}

// Serve partner list data
func (u *exampleDelete) Serve(data *appctx.Data) appctx.Response {

	id := mux.Vars(data.Request)["id"]

	err := u.repo.Delete(data.Request.Context(), cast.ToUint64(id))

	if err != nil {
		logger.Error(logger.MessageFormat("[example-delete] %v", err))

		return *appctx.NewResponse().WithCode(consts.CodeInternalServerError)
	}

	return *appctx.NewResponse().WithCode(consts.CodeSuccess)
}
