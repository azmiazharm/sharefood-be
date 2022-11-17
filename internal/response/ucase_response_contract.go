package response

import (
	"context"
	"sharefood/internal/appctx"
	"sharefood/internal/entity"

	"github.com/google/uuid"
)

type UcaseResponseInterface interface {
	Failed(ctx context.Context, transactionID *uuid.UUID, err error) *appctx.Response
	Success(ctx context.Context, statusCode int, transactionID *uuid.UUID, responseData interface{}) *appctx.Response
	SuccessWithMetadata(ctx context.Context, statusCode int, metadata *entity.Metadata, responseData interface{}) *appctx.Response
}
