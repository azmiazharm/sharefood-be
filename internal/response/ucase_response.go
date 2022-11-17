package response

import (
	"context"
	"sharefood/internal/appctx"
	"sharefood/internal/consts"
	"sharefood/internal/entity"

	// "sharefood/internal/requestlog"
	"sharefood/pkg/logger"
	"sharefood/pkg/tracer"

	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type response struct {
	endpointName string
	request      *http.Request
	// clientRequestLog requestlog.ClientRequestLogInterface
}

func NewResponse(endpointName string, request *http.Request) UcaseResponseInterface {
	return &response{
		endpointName: endpointName,
		request:      request,
		// clientRequestLog: clientRequestLog,
	}
}

func (resp *response) Failed(ctx context.Context, transactionID *uuid.UUID, err error) *appctx.Response {
	// errorEvent := consts.ErrorEvent("ucase_failed_response")
	loggerFields := logger.NewFields(logger.EventName(resp.endpointName))
	ctx = tracer.SpanStart(ctx, "ucase_failed_response")
	defer tracer.SpanFinish(ctx)

	var statusCode int
	switch causer := errors.Cause(err).(type) {
	case consts.Errors:
		statusCode = causer[len(causer)-1].StatusCode
	}

	responseStatus := "error"
	if statusCode >= 500 {
		responseStatus = "failure"
	}

	output := parseResponseStatus(responseStatus, resp.endpointName)

	output = output.WithCode(statusCode).WithError(err)

	if transactionID != nil {
		meta := entity.Meta{
			TransactionID: transactionID,
		}
		output = output.WithMeta(meta)
		// errLog := resp.clientRequestLog.Save(ctx, resp.request, *transactionID, output)
		// if errLog != nil {
		// 	err := errorEvent.WithMessage(consts.CreateClientRequestLogErrorMessage).WrapError(errLog)
		// 	tracer.SpanError(ctx, err)
		// 	return resp.Failed(ctx, nil, err)
		// }
	}

	logger.ErrorWithContext(ctx, err, loggerFields...)

	return output
}

func (resp *response) Success(ctx context.Context, statusCode int, transactionID *uuid.UUID, responseData interface{}) *appctx.Response {
	// errorEvent := consts.ErrorEvent("ucase_success_response")
	loggerFields := logger.NewFields(logger.EventName(resp.endpointName))
	ctx = tracer.SpanStart(ctx, "ucase_success_response")
	defer tracer.SpanFinish(ctx)

	responseStatus := strings.ToLower(consts.StatusSuccess)
	output := parseResponseStatus(responseStatus, resp.endpointName)

	output = output.WithCode(statusCode).WithData(responseData)

	if transactionID != nil {
		meta := entity.Meta{
			TransactionID: transactionID,
		}
		output = output.WithMeta(meta)
		// errLog := resp.clientRequestLog.Save(ctx, resp.request, *transactionID, output)
		// if errLog != nil {
		// 	err := errorEvent.WithMessage(consts.CreateClientRequestLogErrorMessage).WrapError(errLog)
		// 	tracer.SpanError(ctx, err)
		// 	return resp.Failed(ctx, nil, err)
		// }
	}

	logger.InfoWithContext(ctx, output.Message, loggerFields...)

	return output
}

func (resp *response) SuccessWithMetadata(ctx context.Context, statusCode int, metadata *entity.Metadata, responseData interface{}) *appctx.Response {
	// errorEvent := consts.ErrorEvent("ucase_success_response")
	loggerFields := logger.NewFields(logger.EventName(resp.endpointName))
	ctx = tracer.SpanStart(ctx, "ucase_success_response")
	defer tracer.SpanFinish(ctx)

	responseStatus := strings.ToLower(consts.StatusSuccess)
	output := parseResponseStatus(responseStatus, resp.endpointName)
	meta := entity.Metadata{
		TransactionID: metadata.TransactionID,
		PerPage:       metadata.PerPage,
		Page:          metadata.Page,
		Total:         metadata.Total,
		OrderBy:       metadata.OrderBy,
		OrderType:     metadata.OrderType,
	}
	output = output.WithCode(statusCode).WithMeta(meta).WithData(responseData)

	// if metadata.TransactionID != nil {
	// 	errLog := resp.clientRequestLog.Save(ctx, resp.request, *metadata.TransactionID, output)
	// 	if errLog != nil {
	// 		err := errorEvent.WithMessage(consts.CreateClientRequestLogErrorMessage).WrapError(errLog)
	// 		tracer.SpanError(ctx, err)
	// 		return resp.Failed(ctx, nil, err)
	// 	}
	// }

	logger.InfoWithContext(ctx, output.Message, loggerFields...)

	return output
}

func parseResponseStatus(responseStatus string, endpointName string) *appctx.Response {
	var responseEntity, responseState, responseMessage string
	splittedString := strings.Split(endpointName, "_")
	for index, str := range splittedString {
		caser := cases.Title(language.English)
		switch index {
		case 0:
			responseEntity += str
			responseState += str
			responseMessage += caser.String(str)
		case len(splittedString) - 1:
			responseEntity += caser.String(str)
			responseState += caser.String(str) + caser.String(responseStatus)
			responseMessage += " " + caser.String(str) + " " + caser.String(responseStatus)
		default:
			responseEntity += caser.String(str)
			responseState += caser.String(str)
			responseMessage += " " + caser.String(str)
		}
	}

	responseStatus = strings.ToUpper(responseStatus)
	output := appctx.NewResponse().
		WithStatus(responseStatus).
		WithEntity(responseEntity).
		WithState(responseState).
		WithMessage(responseMessage)

	return output
}
