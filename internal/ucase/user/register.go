package user

import (
	"sharefood/internal/appctx"
	"sharefood/internal/consts"
	"sharefood/internal/entity"
	"sharefood/internal/repositories"
	"sharefood/internal/ucase"
	"sharefood/internal/ucase/contract"
	"sharefood/pkg/logger"
	"sharefood/pkg/util"

	"github.com/google/uuid"
	"github.com/thedevsaddam/govalidator"
	"golang.org/x/crypto/bcrypt"
)

type userRegister struct {
	userRepository repositories.User
}

func NewUserRegister(userRepository repositories.User) contract.UseCase {
	return &userRegister{
		userRepository: userRepository,
	}
}

// Serve implements contract.UseCase
func (u *userRegister) Serve(data *appctx.Data) appctx.Response {
	payload := entity.User{}

	err := data.Cast(&payload)
	if err != nil {
		logger.Error(logger.MessageFormat("[user-create] parsing body request error: %v", err))
		return *appctx.NewResponse().WithStatus(consts.StatusFailed).WithEntity("registerUser").WithState("registerUserFailed").WithCode(consts.CodeBadRequest).WithError(err.Error())
	}

	fl := []logger.Field{
		logger.Any("payload", payload),
	}

	rules := govalidator.MapData{
		"name":         []string{"required", "between:3,50"},
		"email":        []string{"required", "min:4", "max:50", "email"},
		"password":     []string{"required", "min:4"},
		"phone_number": []string{"required", "digits_between:6,14"},
	}

	opts := govalidator.Options{
		Data:  &payload,
		Rules: rules,
		// RequiredDefault: true,
	}

	v := govalidator.New(opts)
	ev := v.ValidateStruct()

	if len(ev) != 0 {
		logger.Warn(
			logger.MessageFormat("[user-create] validate request param err: %s", util.DumpToString(ev)),
			fl...)

		err := map[string]interface{}{"validationError": ev}
		return *appctx.NewResponse().WithStatus(consts.StatusFailed).WithEntity("registerUser").WithState("registerUserFailed").WithCode(consts.CodeBadRequest).WithError(err)
	}

	// Check if email registered
	isRegistered := u.userRepository.IsRegistered(data.Request.Context(), payload.Email)
	// fmt.Println(isRegistered, "kalo true berarti udah keregister")
	if isRegistered {
		return *appctx.NewResponse().WithStatus(consts.StatusFailed).WithEntity("registerUser").WithState("registerUserFailed").WithCode(consts.CodeUnprocessableEntity).WithMessage("Failed Register User").WithError("Email Already Registered")
	}

	// Hashing the password with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return *appctx.NewResponse().WithStatus(consts.StatusFailed).WithEntity("registerUser").WithState("registerUserFailed").WithCode(consts.CodeBadRequest).WithError(err.Error())
	}

	//create hashedpass
	payload.Password = string(hashedPassword)

	// create uuid
	payload.ID = uuid.New()
	secret := []byte(data.Config.App.JWTSecret)
	token, errToken := ucase.GenerateJWT(payload, secret)
	if errToken != nil {
		logger.Error(logger.MessageFormat("[user-create] %v", err))
		return *appctx.NewResponse().WithStatus(consts.StatusFailed).WithEntity("registerUser").WithState("registerUserFailed").WithCode(consts.CodeInternalServerError).WithError(errToken.Error())
	}

	// create the account
	err = u.userRepository.Create(data.Request.Context(), &payload)
	if err != nil {
		logger.Error(logger.MessageFormat("[user-create] %v", err))
		return *appctx.NewResponse().WithStatus(consts.StatusFailed).WithEntity("registerUser").WithState("registerUserFailed").WithCode(consts.CodeInternalServerError).WithError(err.Error())
	}

	return *appctx.NewResponse().WithStatus(consts.StatusSuccess).WithEntity("registerUser").WithState("registerUserSuccess").WithCode(consts.CodeCreated).WithData(token)
}
