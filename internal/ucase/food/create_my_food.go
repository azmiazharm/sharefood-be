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
)

type foodCreate struct {
	foodRepository repositories.Food
}

func NewFoodCreate(foodRepository repositories.Food) contract.UseCase {
	return &foodCreate{
		foodRepository: foodRepository,
	}
}

// Serve implements contract.UseCase
func (u *foodCreate) Serve(data *appctx.Data) appctx.Response {
	request := data.Request
	response := response.NewResponse("create_food", request)
	errorEvent := consts.ErrorEvent("create_food")
	ctx := tracer.SpanStart(request.Context(), "create_food")
	defer tracer.SpanFinish(ctx)

	transactionID := uuid.New()

	payload := entity.Food{}

	err := data.Cast(&payload)
	if err != nil {
		logger.Error(logger.MessageFormat("[food-create] parsing body request error: %v", err))
		err := errorEvent.WithMessage(consts.CreateFoodErrorMessage).WithCode(consts.CodeUnprocessableEntity).WrapError(err)
		return *response.Failed(ctx, &transactionID, err)
	}

	idUser := data.Request.Header.Get("idUser")
	uuidUser, errFood := uuid.Parse(idUser)
	if errFood != nil {
		logger.Error(logger.MessageFormat("[food-get] parsing id error: %v", errFood))
		err := errorEvent.WithMessage(consts.CreateFoodErrorMessage).WithCode(consts.CodeUnprocessableEntity).WrapError(err)
		return *response.Failed(ctx, &transactionID, err)
	}

	// create uuid for food
	payload.ID = uuid.New()

	// pass user id to payload
	payload.IDUser = uuidUser

	fmt.Println(payload)

	// create food to db
	err = u.foodRepository.Create(data.Request.Context(), &payload)
	if err != nil {
		logger.Error(logger.MessageFormat("[user-create] %v", err))
		err := errorEvent.WithMessage(consts.CreateFoodErrorMessage).WithCode(consts.CodeUnprocessableEntity).WrapError(err)
		return *response.Failed(ctx, &transactionID, err)
	}

	return *response.Success(ctx, consts.CodeCreated, &transactionID, nil)
}
