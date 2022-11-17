package request

import (
	"sharefood/internal/appctx"
	"sharefood/internal/consts"
	"sharefood/internal/entity"
	"sharefood/internal/repositories"
	"sharefood/internal/response"
	"sharefood/internal/ucase/contract"
	"sharefood/pkg/logger"
	"sharefood/pkg/tracer"

	"github.com/google/uuid"
)

type requestAction struct {
	requestRepository repositories.Request
	foodRepository    repositories.Food
}

func NewRequestAction(requestRepository repositories.Request, foodRepository repositories.Food) contract.UseCase {
	return &requestAction{
		requestRepository: requestRepository,
		foodRepository:    foodRepository,
	}
}

// Serve implements contract.UseCase
func (u *requestAction) Serve(data *appctx.Data) appctx.Response {
	request := data.Request
	response := response.NewResponse("request_action", request)
	errorEvent := consts.ErrorEvent("request_action")
	ctx := tracer.SpanStart(request.Context(), "request_action")
	defer tracer.SpanFinish(ctx)

	transactionID := uuid.New()

	payload := entity.RequestAction{}

	err := data.Cast(&payload)
	if err != nil {
		logger.Error(logger.MessageFormat("[request-action] parsing body request error: %v", err))
		err := errorEvent.WithMessage(consts.CreateRequestErrorMessage).WithCode(consts.CodeUnprocessableEntity).WrapError(err)
		return *response.Failed(ctx, &transactionID, err)
	}

	idUser := data.Request.Header.Get("idUser")
	uuidUser, errFood := uuid.Parse(idUser)
	if errFood != nil {
		logger.Error(logger.MessageFormat("[request-action] parsing id error: %v", errFood))
		err := errorEvent.WithMessage(consts.CreateRequestErrorMessage).WithCode(consts.CodeUnprocessableEntity).WrapError(err)
		return *response.Failed(ctx, &transactionID, err)
	}

	// get requestnya dulu - join ke food
	reqFood, err := u.requestRepository.GetRequestFoodByIDRequest(ctx, payload.ID)
	if err != nil {
		logger.Error(logger.MessageFormat("[request-action] %v", err))
		err := errorEvent.WithMessage(consts.ActionRequestErrorMessage).WithCode(consts.CodeUnprocessableEntity).WrapError(err)
		return *response.Failed(ctx, &transactionID, err)
	}

	// kalo food-id usernya sama kaya user header berarti boleh action (cuma user yang punya yang boleh action)
	if reqFood.IDUserFood != uuidUser {
		logger.Error(logger.MessageFormat("[request-action] id user not match"))
		err := errorEvent.WithMessage(consts.IdNotValidMessage).WithCode(consts.CodeForbidden).WrapError(consts.Error(consts.StatusForbidden))
		return *response.Failed(ctx, &transactionID, err)
	}

	// misal reject -> update statusnya aja jadi 2
	if payload.Action == "reject" {
		err := u.requestRepository.RejectRequest(ctx, payload.ID)
		if err != nil {
			logger.Error(logger.MessageFormat("[request-action] %v", err))
			err := errorEvent.WithMessage(consts.ActionRequestNotValid).WithCode(consts.CodeUnprocessableEntity).WrapError(err)
			return *response.Failed(ctx, &transactionID, err)
		}
		return *response.Success(ctx, consts.CodeSuccess, &transactionID, nil)

	} else if payload.Action == "accept" {
		// misal accept -> kalo request quantitynya lebih banyak daripada stok di food -> cancel
		if reqFood.Stock < reqFood.Quantity {
			logger.Error(logger.MessageFormat("[request-action] not enough stok"))

			err := errorEvent.WithMessage(consts.NotEnoughQuantity).WithCode(consts.CodeUnprocessableEntity).WrapError(consts.Error(consts.NotEnoughQuantity))
			return *response.Failed(ctx, &transactionID, err)
		}

		err := u.requestRepository.AcceptRequest(ctx, payload.ID)
		if err != nil {
			logger.Error(logger.MessageFormat("[request-action] %v", err))
			err := errorEvent.WithMessage(consts.ActionRequestNotValid).WithCode(consts.CodeUnprocessableEntity).WrapError(err)
			return *response.Failed(ctx, &transactionID, err)
		}
		return *response.Success(ctx, consts.CodeSuccess, &transactionID, reqFood)

	} else {
		logger.Error(logger.MessageFormat("[request-action] Unknown error"))
		err := errorEvent.WithMessage(consts.RespUnknown).WithCode(consts.CodeUnknownError).WrapError(err)
		return *response.Failed(ctx, &transactionID, err)
	}
}
