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

type myFoodDelete struct {
	foodRepositories repositories.Food
}

func NewMyFoodDelete(foodRepositories repositories.Food) contract.UseCase {
	return &myFoodDelete{
		foodRepositories: foodRepositories,
	}
}

// Serve implements contract.UseCase
func (u *myFoodDelete) Serve(data *appctx.Data) appctx.Response {
	request := data.Request
	response := response.NewResponse("delete_my_food", request)
	errorEvent := consts.ErrorEvent("delete_my_food")
	ctx := tracer.SpanStart(request.Context(), "delete_my_food")
	defer tracer.SpanFinish(ctx)

	transactionID := uuid.New()

	params := mux.Vars(data.Request)
	rawID := params["id"]

	idUser := data.Request.Header.Get("idUser")
	uuidUser, errFood := uuid.Parse(idUser)
	if errFood != nil {
		logger.Error(logger.MessageFormat("[food-delete] parsing id error: %v", errFood))
		err := errorEvent.WithMessage(consts.CreateFoodErrorMessage).WithCode(consts.CodeUnprocessableEntity).WrapError(errFood)
		return *response.Failed(ctx, &transactionID, err)
	}

	idFood, errFood := uuid.Parse(rawID)
	if errFood != nil {
		logger.Error(logger.MessageFormat("[food-delete] parsing id error: %v", errFood))

		err := errorEvent.WithMessage(consts.IdNotValidMessage).WithCode(consts.CodeBadRequest).WrapError(errFood)
		return *response.Failed(ctx, &transactionID, err)
	}

	food, errFood := u.foodRepositories.GetDetailByID(ctx, idFood)
	if errFood != nil {
		logger.Error(logger.MessageFormat("[food-delete] Error get food: %v", errFood))
		logger.Error(errFood)

		err := errorEvent.WithMessage(consts.FoodNotFoundMessage).WithCode(consts.CodeNotFound).WrapError(errFood)
		return *response.Failed(ctx, &transactionID, err)
	}

	//check permission
	if uuidUser != food.IDUser {
		logger.Error(logger.MessageFormat("[food-delete] id user not match"))
		err := errorEvent.WithMessage(consts.IdNotValidMessage).WithCode(consts.CodeForbidden).WrapError(consts.Error(consts.StatusForbidden))
		return *response.Failed(ctx, &transactionID, err)
	}

	//delete food
	err := u.foodRepositories.DeleteByID(ctx, idFood)
	if err != nil {
		logger.Error(logger.MessageFormat("[food-delete] delete food error: %v", err))
		err := errorEvent.WithMessage(consts.CreateFoodErrorMessage).WithCode(consts.CodeUnprocessableEntity).WrapError(err)
		return *response.Failed(ctx, &transactionID, err)
	}

	return *response.Success(ctx, consts.CodeSuccess, &transactionID, nil)
}
