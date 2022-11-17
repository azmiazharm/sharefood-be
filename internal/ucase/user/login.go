package user

import (
	"sharefood/internal/appctx"
	"sharefood/internal/consts"
	"sharefood/internal/entity"
	"sharefood/internal/repositories"
	"sharefood/internal/ucase"
	"sharefood/internal/ucase/contract"
	"sharefood/pkg/logger"

	"golang.org/x/crypto/bcrypt"
)

type userLogin struct {
	userRepository repositories.User
}

func NewUserLogin(userRepository repositories.User) contract.UseCase {
	return &userLogin{
		userRepository: userRepository,
	}
}

// Serve implements contract.UseCase
func (u *userLogin) Serve(data *appctx.Data) appctx.Response {
	payload := entity.UserLogin{}

	err := data.Cast(&payload)
	if err != nil {
		logger.Error(logger.MessageFormat("[user-login] parsing body request error: %v", err))
		return *appctx.NewResponse().WithCode(consts.CodeBadRequest).WithError(err.Error()).WithStatus(consts.StatusFailed).WithEntity("login").WithState("loginFailed")
	}

	// Check is account available
	userAccount, err := u.userRepository.GetByEmail(data.Request.Context(), payload.Email)
	if err != nil {
		return *appctx.NewResponse().WithCode(consts.CodeNotFound).WithMessage("Failed Login User").WithError(err.Error()).WithStatus(consts.StatusFailed).WithEntity("login").WithState("loginFailed")
	}

	pass := []byte(payload.Password)
	err = bcrypt.CompareHashAndPassword([]byte(userAccount.Password), pass)
	if err != nil {
		return *appctx.NewResponse().WithCode(consts.CodeAuthenticationFailure).WithMessage("Failed Login User").WithError(err.Error()).WithStatus(consts.StatusFailed).WithEntity("login").WithState("loginFailed")
	}

	secret := []byte(data.Config.App.JWTSecret)
	token, err := ucase.GenerateJWT(userAccount, secret)
	if err != nil {
		return *appctx.NewResponse().WithCode(consts.CodeAuthenticationFailure).WithMessage("Failed Login User").WithError(err.Error()).WithStatus(consts.StatusFailed).WithEntity("login").WithState("loginFailed")
	}

	return *appctx.NewResponse().WithCode(consts.CodeSuccess).WithData(token).WithMessage("Success Login User").WithStatus(consts.StatusSuccess).WithEntity("login").WithState("loginSuccess")
}

// func generateJWT(user entity.User, secret []byte) (entity.TokenResponse, error) {
// 	expiredAt := time.Now().Add(time.Hour * 2).Local()

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, entity.TokenClaims{
// 		ID: user.ID,
// 		StandardClaims: jwt.StandardClaims{
// 			ExpiresAt: &jwt.Time{
// 				Time: expiredAt,
// 			},
// 		},
// 	})

// 	signedToken, errToken := token.SignedString(secret)
// 	if errToken != nil {
// 		return entity.TokenResponse{}, errToken
// 	}

// 	tokenResponse := entity.TokenResponse{
// 		Type:      "bearer",
// 		Token:     signedToken,
// 		ExpiredAt: expiredAt,
// 	}

// 	return tokenResponse, nil
// }
