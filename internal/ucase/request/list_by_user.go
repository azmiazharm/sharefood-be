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
)

type requestUserList struct {
	requestRepository repositories.Request
}

func NewRequestUserList(requestRepository repositories.Request) contract.UseCase {
	return &requestUserList{
		requestRepository: requestRepository,
	}
}

// Serve implements contract.UseCase
func (u *requestUserList) Serve(data *appctx.Data) appctx.Response {
	request := data.Request
	response := response.NewResponse("list_requests_user", request)
	errorEvent := consts.ErrorEvent("list_requests_user")
	ctx := tracer.SpanStart(request.Context(), "list_requests_user")
	defer tracer.SpanFinish(ctx)

	transactionID := uuid.New()

	idUser := data.Request.Header.Get("idUser")
	fmt.Println(idUser)
	uuidUser, errFood := uuid.Parse(idUser)
	if errFood != nil {
		logger.Error(logger.MessageFormat("[request-get] parsing id error: %v", errFood))
		err := errorEvent.WithMessage(consts.CreateFoodErrorMessage).WithCode(consts.CodeUnprocessableEntity).WrapError(errFood)
		return *response.Failed(ctx, &transactionID, err)
	}

	// params := mux.Vars(data.Request)
	// rawID := params["id"]
	// idFood, errFood := uuid.Parse(rawID)
	// if errFood != nil {
	// 	logger.Error(logger.MessageFormat("[request-get] parsing id error: %v", errFood))

	// 	err := errorEvent.WithMessage(consts.IdNotValidMessage).WithCode(consts.CodeBadRequest).WrapError(errFood)
	// 	return *response.Failed(ctx, &transactionID, err)
	// }

	requests, err := u.requestRepository.ListbyUser(ctx, uuidUser)
	if err != nil {
		logger.Error("Error get list of users")
		logger.Error(err)
		return *response.Failed(ctx, &transactionID, err)
	}
	fmt.Println(requests)

	return *response.Success(ctx, consts.CodeSuccess, &transactionID, requests)
}
