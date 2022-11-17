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
	"github.com/gorilla/mux"
)

type requestCreate struct {
	requestRepository repositories.Request
	foodRepository    repositories.Food
}

func NewRequestFoodCreate(requestRepository repositories.Request, foodRepository repositories.Food) contract.UseCase {
	return &requestCreate{
		requestRepository: requestRepository,
		foodRepository:    foodRepository,
	}
}

// Serve implements contract.UseCase
func (u *requestCreate) Serve(data *appctx.Data) appctx.Response {
	request := data.Request
	response := response.NewResponse("create_request", request)
	errorEvent := consts.ErrorEvent("create_request")
	ctx := tracer.SpanStart(request.Context(), "create_request")
	defer tracer.SpanFinish(ctx)

	transactionID := uuid.New()

	payload := entity.Request{}

	err := data.Cast(&payload)
	if err != nil {
		logger.Error(logger.MessageFormat("[request-create] parsing body request error: %v", err))
		err := errorEvent.WithMessage(consts.CreateRequestErrorMessage).WithCode(consts.CodeUnprocessableEntity).WrapError(err)
		return *response.Failed(ctx, &transactionID, err)
	}

	idUser := data.Request.Header.Get("idUser")
	uuidUser, errFood := uuid.Parse(idUser)
	if errFood != nil {
		logger.Error(logger.MessageFormat("[request-create] parsing id error: %v", errFood))
		err := errorEvent.WithMessage(consts.CreateRequestErrorMessage).WithCode(consts.CodeUnprocessableEntity).WrapError(err)
		return *response.Failed(ctx, &transactionID, err)
	}

	params := mux.Vars(data.Request)
	rawID := params["id"]
	uuidFood, errFood := uuid.Parse(rawID)
	if errFood != nil {
		logger.Error(logger.MessageFormat("[request-create] parsing id error: %v", errFood))

		err := errorEvent.WithMessage(consts.IdNotValidMessage).WithCode(consts.CodeBadRequest).WrapError(errFood)
		return *response.Failed(ctx, &transactionID, err)
	}

	// get current food to check quantity left
	food, errFood := u.foodRepository.GetDetailByID(ctx, uuidFood)
	if errFood != nil {
		logger.Error(logger.MessageFormat("[request-create] food not found: %v", errFood))

		err := errorEvent.WithMessage(consts.FoodNotFoundMessage).WithCode(consts.CodeBadRequest).WrapError(errFood)
		return *response.Failed(ctx, &transactionID, err)
	}

	// check if quantity avail requested is greeater than requested
	if food.Quantity <= payload.Quantity {
		logger.Error(logger.MessageFormat("[request-create] too many quantity requested"))

		err := errorEvent.WithMessage(consts.NotEnoughQuantity).WithCode(consts.CodeUnprocessableEntity).WrapError(consts.Error(consts.NotEnoughQuantity))
		return *response.Failed(ctx, &transactionID, err)
	}

	// create uuid for food
	payload.ID = uuid.New()

	// pass user id to payload
	payload.IDUser = uuidUser

	// pass food id to payload
	payload.IDFood = uuidFood

	// create food to db
	err = u.requestRepository.Create(data.Request.Context(), &payload)
	if err != nil {
		logger.Error(logger.MessageFormat("[user-create] %v", err))
		err := errorEvent.WithMessage(consts.CreateFoodErrorMessage).WithCode(consts.CodeUnprocessableEntity).WrapError(err)
		return *response.Failed(ctx, &transactionID, err)
	}

	return *response.Success(ctx, consts.CodeCreated, &transactionID, nil)
}
