package food

import (
	"fmt"
	"sharefood/internal/appctx"
	"sharefood/internal/consts"
	"sharefood/internal/entity"
	"sharefood/internal/repositories"
	"sharefood/internal/response"
	"sharefood/internal/ucase/contract"
	"sharefood/pkg/logger"
	"sharefood/pkg/tracer"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type myFoodUpdate struct {
	foodRepositories repositories.Food
}

func NewMyFoodUpdate(foodRepositories repositories.Food) contract.UseCase {
	return &myFoodUpdate{
		foodRepositories: foodRepositories,
	}
}

// Serve implements contract.UseCase
func (u *myFoodUpdate) Serve(data *appctx.Data) appctx.Response {
	request := data.Request
	response := response.NewResponse("update_my_foods", request)
	errorEvent := consts.ErrorEvent("update_my_foods")
	ctx := tracer.SpanStart(request.Context(), "update_my_foods")
	defer tracer.SpanFinish(ctx)

	transactionID := uuid.New()

	// get id food from param
	params := mux.Vars(data.Request)
	rawID := params["id"]
	if rawID == "" {
		err := errorEvent.WithCode(consts.CodeUnprocessableEntity).WrapError(consts.Error(consts.IDRequired))
		tracer.SpanError(ctx, err)
		return *response.Failed(ctx, nil, err)
	}

	// parse id_food from req
	uuidFood, errFood := uuid.Parse(rawID)
	if errFood != nil {
		logger.Error(logger.MessageFormat("[food-get] parsing id error: %v", errFood))
		err := errorEvent.WithMessage(consts.UpdateFoodErrorMessage).WithCode(consts.CodeUnprocessableEntity).WrapError(errFood)
		return *response.Failed(ctx, &transactionID, err)
	}

	// parse id_user from header
	idUser := data.Request.Header.Get("idUser")
	uuidUser, errFood := uuid.Parse(idUser)
	if errFood != nil {
		logger.Error(logger.MessageFormat("[food-get] parsing id error: %v", errFood))
		err := errorEvent.WithMessage(consts.UpdateFoodErrorMessage).WithCode(consts.CodeUnprocessableEntity).WrapError(errFood)
		return *response.Failed(ctx, &transactionID, err)
	}

	// cast to entity
	payload := entity.Food{}
	errCast := data.Cast(&payload)
	if errCast != nil {
		logger.Error(logger.MessageFormat("[food-update] parsing body request error: %v", errCast))
		err := errorEvent.WithMessage(consts.UpdateFoodErrorMessage).WithCode(consts.CodeUnprocessableEntity).WrapError(errCast)
		return *response.Failed(ctx, &transactionID, err)
	}

	// check if eligible to update
	oldFood, errFood := u.foodRepositories.GetDetailByID(data.Request.Context(), uuidFood)
	if errFood != nil {
		logger.Error(logger.MessageFormat("[food-get] Error get food: %v", errFood))
		logger.Error(errFood)

		err := errorEvent.WithMessage(consts.FoodNotFoundMessage).WithCode(consts.CodeNotFound).WrapError(errFood)
		return *response.Failed(ctx, &transactionID, err)
	}

	if uuidUser != oldFood.IDUser {
		logger.Error(logger.MessageFormat("[food-get] id user not match"))
		err := errorEvent.WithMessage(consts.IdNotValidMessage).WithCode(consts.CodeForbidden).WrapError(consts.Error(consts.StatusForbidden))
		return *response.Failed(ctx, &transactionID, err)
	}

	fmt.Println(payload)

	// do update with payload
	payload.ID = uuidFood
	errUpdateFood := u.foodRepositories.Update(ctx, &payload)
	if errUpdateFood != nil {
		err := errorEvent.WithMessage(consts.UpdateFoodErrorMessage).WrapError(errUpdateFood)
		tracer.SpanError(ctx, err)
		return *response.Failed(ctx, nil, err)
	}

	// return last data of food
	outputFood, errFood := u.foodRepositories.GetDetailByID(data.Request.Context(), uuidFood)
	if errFood != nil {
		logger.Error(logger.MessageFormat("[food-get] Error get food: %v", errFood))
		logger.Error(errFood)

		err := errorEvent.WithMessage(consts.FoodNotFoundMessage).WithCode(consts.CodeNotFound).WrapError(errFood)
		return *response.Failed(ctx, &transactionID, err)
	}

	return *response.Success(ctx, consts.CodeSuccess, &transactionID, outputFood)
}
