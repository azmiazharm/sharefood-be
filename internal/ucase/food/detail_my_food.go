package food

import (
	"sharefood/internal/appctx"
	"sharefood/internal/consts"
	"sharefood/internal/repositories"
	"sharefood/internal/response"
	"sharefood/internal/ucase/contract"
	"sharefood/pkg/logger"
	"sharefood/pkg/tracer"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type myFoodGet struct {
	foodRepositories repositories.Food
}

func NewMyFoodGet(foodRepositories repositories.Food) contract.UseCase {
	return &myFoodGet{
		foodRepositories: foodRepositories,
	}
}

// Serve implements contract.UseCase
func (u *myFoodGet) Serve(data *appctx.Data) appctx.Response {
	request := data.Request
	response := response.NewResponse("get_detail_my_food", request)
	errorEvent := consts.ErrorEvent("get_detail_my_food")
	ctx := tracer.SpanStart(request.Context(), "get_detail_my_food")
	defer tracer.SpanFinish(ctx)

	transactionID := uuid.New()

	params := mux.Vars(data.Request)
	rawID := params["id"]

	idUser := data.Request.Header.Get("idUser")
	uuidUser, errFood := uuid.Parse(idUser)
	if errFood != nil {
		logger.Error(logger.MessageFormat("[food-get] parsing id error: %v", errFood))
		err := errorEvent.WithMessage(consts.CreateFoodErrorMessage).WithCode(consts.CodeUnprocessableEntity).WrapError(errFood)
		return *response.Failed(ctx, &transactionID, err)
	}

	id_food, errFood := uuid.Parse(rawID)
	if errFood != nil {
		logger.Error(logger.MessageFormat("[food-get] parsing id error: %v", errFood))

		err := errorEvent.WithMessage(consts.IdNotValidMessage).WithCode(consts.CodeBadRequest).WrapError(errFood)
		return *response.Failed(ctx, &transactionID, err)
	}

	food, errFood := u.foodRepositories.GetDetailByID(data.Request.Context(), id_food)
	if errFood != nil {
		logger.Error(logger.MessageFormat("[food-get] Error get food: %v", errFood))
		logger.Error(errFood)

		err := errorEvent.WithMessage(consts.FoodNotFoundMessage).WithCode(consts.CodeNotFound).WrapError(errFood)
		return *response.Failed(ctx, &transactionID, err)
	}

	if uuidUser != food.IDUser {
		logger.Error(logger.MessageFormat("[food-get] id user not match"))
		err := errorEvent.WithMessage(consts.IdNotValidMessage).WithCode(consts.CodeForbidden).WrapError(consts.Error(consts.StatusForbidden))
		return *response.Failed(ctx, &transactionID, err)
	}

	return *response.Success(ctx, consts.CodeSuccess, &transactionID, food)
}
