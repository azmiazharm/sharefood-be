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
)

type myFoodList struct {
	foodRepositories repositories.Food
}

func NewMyFoodList(foodRepositories repositories.Food) contract.UseCase {
	return &myFoodList{
		foodRepositories: foodRepositories,
	}
}

// Serve implements contract.UseCase
func (u *myFoodList) Serve(data *appctx.Data) appctx.Response {
	request := data.Request
	response := response.NewResponse("list_my_foods", request)
	errorEvent := consts.ErrorEvent("list_my_foods")
	ctx := tracer.SpanStart(request.Context(), "list_my_foods")
	defer tracer.SpanFinish(ctx)

	transactionID := uuid.New()

	idUser := data.Request.Header.Get("idUser")
	uuidUser, errFood := uuid.Parse(idUser)
	if errFood != nil {
		logger.Error(logger.MessageFormat("[food-get] parsing id error: %v", errFood))
		err := errorEvent.WithMessage(consts.CreateFoodErrorMessage).WithCode(consts.CodeUnprocessableEntity).WrapError(errFood)
		return *response.Failed(ctx, &transactionID, err)
	}

	foods, errList := u.foodRepositories.ListMy(data.Request.Context(), uuidUser)
	if errList != nil {
		logger.Error("Error get list of foods")
		logger.Error(errList)

		err := errorEvent.WithMessage(consts.FindVendorEnvironmentByClientIDErrorMessage).WithCode(consts.CodeInternalServerError).WrapError(errList)
		return *response.Failed(ctx, &transactionID, err)
	}

	return *response.Success(ctx, consts.CodeSuccess, &transactionID, foods)
}
