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

type foodList struct {
	foodRepositories repositories.Food
}

func NewFoodList(foodRepositories repositories.Food) contract.UseCase {
	return &foodList{
		foodRepositories: foodRepositories,
	}
}

// Serve implements contract.UseCase
func (u *foodList) Serve(data *appctx.Data) appctx.Response {
	request := data.Request
	response := response.NewResponse("get_foods", request)
	errorEvent := consts.ErrorEvent("get_foods")
	ctx := tracer.SpanStart(request.Context(), "get_foods")
	defer tracer.SpanFinish(ctx)

	transactionID := uuid.New()

	foods, errList := u.foodRepositories.List(data.Request.Context())
	if errList != nil {
		logger.Error("Error get list of foods")
		logger.Error(errList)

		err := errorEvent.WithMessage(consts.FindVendorEnvironmentByClientIDErrorMessage).WithCode(consts.CodeInternalServerError).WrapError(errList)
		return *response.Failed(ctx, &transactionID, err)

		// return *appctx.NewResponse().WithMessage("Failed Get all food").WithCode(consts.CodeInternalServerError).WithEntity("getAllSharedFood").WithState("getAllSharedFoodFailed")
	}

	// return *response.SuccessWithMetadata(ctx, consts.CodeSuccess, &transactionID, foods)
	return *response.Success(ctx, consts.CodeSuccess, &transactionID, foods)
	// return *appctx.NewResponse().WithEntity("getAllSharedFood").WithState("getAllSharedFoodSuccess").WithCode(consts.CodeSuccess).WithData(foods).WithStatus(consts.StatusSuccess).WithMessage("")
}
