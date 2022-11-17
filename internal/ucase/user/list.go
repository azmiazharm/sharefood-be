package user

import (
	"sharefood/internal/appctx"
	"sharefood/internal/consts"
	"sharefood/internal/repositories"
	"sharefood/internal/ucase/contract"
	"sharefood/pkg/logger"
)

type userList struct {
	userRepository repositories.User
}

func NewUserList(userRepository repositories.User) contract.UseCase {
	return &userList{
		userRepository: userRepository,
	}
}

// Serve implements contract.UseCase
func (u *userList) Serve(data *appctx.Data) appctx.Response {
	users, err := u.userRepository.List(data.Request.Context())
	if err != nil {
		logger.Error("Error get list of users")
		logger.Error(err)
		return *appctx.NewResponse().WithMessage("Ini messagenya").WithCode(consts.CodeInternalServerError)
	}

	return *appctx.NewResponse().WithCode(consts.CodeSuccess).WithData(users)
}
