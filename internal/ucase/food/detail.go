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

type foodGet struct {
	foodRepositories repositories.Food
}

func NewFoodGet(foodRepositories repositories.Food) contract.UseCase {
	return &foodGet{
		foodRepositories: foodRepositories,
	}
}

// Serve implements contract.UseCase
func (u *foodGet) Serve(data *appctx.Data) appctx.Response {
	request := data.Request
	response := response.NewResponse("get_detail_shared_food", request)
	errorEvent := consts.ErrorEvent("get_detail_shared_food")
	ctx := tracer.SpanStart(request.Context(), "get_detail_shared_food")
	defer tracer.SpanFinish(ctx)

	transactionID := uuid.New()

	params := mux.Vars(data.Request)
	rawID := params["id"]

	idFood, errFood := uuid.Parse(rawID)
	if errFood != nil {
		logger.Error(logger.MessageFormat("[food-get] parsing id error: %v", errFood))

		err := errorEvent.WithMessage(consts.IdNotValidMessage).WithCode(consts.CodeBadRequest).WrapError(errFood)
		return *response.Failed(ctx, &transactionID, err)
		// return *appctx.NewResponse().WithCode(consts.CodeBadRequest).WithError("Invalid Food ID").WithMessage("Get Food Failed").WithEntity("getDetailSharedFood").WithState("getDetailSharedFoodFailed").WithStatus(consts.StatusFailed)
	}

	food, errFood := u.foodRepositories.GetDetailByID(data.Request.Context(), idFood)
	if errFood != nil {
		logger.Error(logger.MessageFormat("[food-get] Error get food: %v", errFood))
		logger.Error(errFood)

		err := errorEvent.WithMessage(consts.FoodNotFoundMessage).WithCode(consts.CodeNotFound).WrapError(errFood)
		return *response.Failed(ctx, &transactionID, err)
		// return *appctx.NewResponse().WithCode(consts.CodeNotFound).WithError("Invalid Food ID").WithMessage("Get Food Failed").WithEntity("getDetailSharedFood").WithState("getDetailSharedFoodFailed").WithStatus(consts.StatusFailed)
	}

	return *response.Success(ctx, consts.CodeSuccess, &transactionID, food)
	// return *appctx.NewResponse().WithCode(consts.CodeSuccess).WithData(food).WithMessage("Get Food Success").WithEntity("getDetailSharedFood").WithState("getDetailSharedFoodSuccess").WithStatus(consts.StatusSuccess)
}
