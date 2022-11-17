package request

import (
	"fmt"
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

type requestFoodList struct {
	requestRepository repositories.Request
}

func NewRequestFoodList(requestRepository repositories.Request) contract.UseCase {
	return &requestFoodList{
		requestRepository: requestRepository,
	}
}

// Serve implements contract.UseCase
func (u *requestFoodList) Serve(data *appctx.Data) appctx.Response {
	request := data.Request
	response := response.NewResponse("list_requests_food", request)
	errorEvent := consts.ErrorEvent("list_requests_food")
	ctx := tracer.SpanStart(request.Context(), "list_requests_food")
	defer tracer.SpanFinish(ctx)

	transactionID := uuid.New()

	// idUser := data.Request.Header.Get("idUser")
	// uuidUser, errFood := uuid.Parse(idUser)
	// if errFood != nil {
	// 	logger.Error(logger.MessageFormat("[request-get] parsing id error: %v", errFood))
	// 	err := errorEvent.WithMessage(consts.CreateFoodErrorMessage).WithCode(consts.CodeUnprocessableEntity).WrapError(errFood)
	// 	return *response.Failed(ctx, &transactionID, err)
	// }

	params := mux.Vars(data.Request)
	rawID := params["id"]
	idFood, errFood := uuid.Parse(rawID)
	if errFood != nil {
		logger.Error(logger.MessageFormat("[list_requests_food parsing id error: %v", errFood))

		err := errorEvent.WithMessage(consts.IdNotValidMessage).WithCode(consts.CodeBadRequest).WrapError(errFood)
		return *response.Failed(ctx, &transactionID, err)
	}

	requests, err := u.requestRepository.ListbyFood(ctx, idFood)
	if err != nil {
		logger.Error("Error get list of request")
		logger.Error(err)
		return *response.Failed(ctx, &transactionID, err)
	}
	fmt.Println(requests)

	return *response.Success(ctx, consts.CodeSuccess, &transactionID, requests)
}
