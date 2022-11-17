package consts

import (
	"net/http"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pkg/errors"
)

type Error string

func (e Error) Error() string {
	return string(e)
}

type WrappedError struct {
	StatusCode int    `json:"-"`
	Event      string `json:"event,omitempty"`
	Err        error  `json:"error,omitempty"`
}

func ErrorEvent(event string) *WrappedError {
	return &WrappedError{
		StatusCode: CodeInternalServerError,
		Event:      event,
		Err:        Error(strings.ToLower(http.StatusText(CodeInternalServerError))),
	}
}

func (e *WrappedError) WithCode(statusCode int) *WrappedError {
	e.StatusCode = statusCode
	return e
}

func (e *WrappedError) WithMessage(message string) *WrappedError {
	e.Err = Error(message)
	return e
}

type Errors []WrappedError

func (e Errors) Error() string {
	return e[0].Err.Error()
}

func (e *WrappedError) WrapError(err error) Errors {
	switch causer := errors.Cause(err).(type) {
	case Errors:
		wrappedError := WrappedError{
			StatusCode: causer[len(causer)-1].StatusCode,
			Event:      e.Event,
			Err:        transformFieldError(e.Err),
		}
		causer = append(causer, wrappedError)
		return causer

	default:
		errs := make(Errors, 0)
		wrappedError := WrappedError{
			StatusCode: e.StatusCode,
			Event:      e.Event,
			Err:        transformFieldError(err),
		}
		errs = append(errs, wrappedError)
		return errs
	}
}

type FieldError struct {
	Field string
	Error error
}

type FieldErrors []FieldError

func (err FieldErrors) Error() string {
	return err[0].Error.Error()
}

func transformFieldError(err error) error {
	switch causer := errors.Cause(err).(type) {
	case validation.Errors:
		fieldErrors := make(FieldErrors, 0)
		splittedRawFieldErrors := strings.Split(causer.Error(), ";")
		for _, validationError := range splittedRawFieldErrors {
			splittedValidationError := strings.Split(validationError, ":")
			fieldError := FieldError{
				Field: strings.TrimSpace(splittedValidationError[0]),
				Error: Error(strings.TrimSpace(splittedValidationError[len(splittedValidationError)-1])),
			}
			fieldErrors = append(fieldErrors, fieldError)
		}
		return fieldErrors
	default:
		return Error(causer.Error())
	}
}

const (
	IdNotValidMessage         = "id not valid"
	FoodNotFoundMessage       = "food not found"
	CreateFoodErrorMessage    = "create food error"
	CreateRequestErrorMessage = "create request error"
	ActionRequestErrorMessage = "action request error"
	UpdateFoodErrorMessage    = "update food error"
	StatusForbidden           = "forbidden"
	NotEnoughQuantity         = "too many quantity requested"
	ActionRequestNotValid     = "action cannot be processed"
	ActionAlreadyDone         = "action already accepted or rejected"

	AssetNotFoundMessage                                = "asset not found"
	ClientIDNotValidMessage                             = "clients id not valid"
	ClientRequestLogNotFoundMessage                     = "clients request log not found"
	VendorContentNotFoundMessage                        = "vendor content not found"
	VendorContentTypeNotFoundMessage                    = "vendor content type not found"
	VendorFeatureNotFoundMessage                        = "vendor feature not found"
	UnprocessableClientResponseDataMessage              = "unprocessable clients response data"
	UnprocessablePrivyRequestDataMessage                = "unprocessable privy request data"
	UnprocessablePrivyResponseDataMessage               = "unprocessable privy response data"
	UnprocessableRequestDataMessage                     = "unprocessable request data"
	UnprocessableIsObjectMessage                        = "unprocessable is_object"
	UnprocessableFileNameMessage                        = "unprocessable file name"
	UnprocessableRequestTransactionID                   = "unprocessable request transaction id"
	ContentTypeRequiredMessage                          = "content type required"
	ClientIDRequiredMessage                             = "clients id required"
	TransactionIDRequiredMessage                        = "transaction id required"
	TimestampRequiredMessage                            = "timestamp required"
	SignatureRequiredMessage                            = "signature required"
	RequestTransactionIDRequiredMessage                 = "request transaction id required"
	IsObjectRequiredMessage                             = "is_object required"
	TemplateIDRequiredMessage                           = "template id required"
	AssetAlreadyExistsMessage                           = "asset already exists"
	CreateClientRequestLogErrorMessage                  = "create clients request log error"
	CreateVendorRequestLogErrorMessage                  = "create vendor request log error"
	CreateAssetErrorMessage                             = "create asset error"
	CreateVendorContentTypeErrorMessage                 = "create vendor content type error"
	CreateVendorContentErrorMessage                     = "create vendor content error"
	TransformClientRequestDataByContentTypeErrorMessage = "transform clients request data by content type error"
	TransformVendorRequestDataByContentTypeErrorMessage = "transform vendor request data by content type error"
	ParseContextValuesErrorMessage                      = "parse context values error"
	ParseClientHeaderParamsErrorMessage                 = "parse clients header params error"
	ParseUserDataErrorMessage                           = "parse user data error"
	FindVendorEnvironmentByClientIDErrorMessage         = "find vendor environment by clients id error"
	FindClientRequestLogByTransactionIDErrorMessage     = "find clients request log by transaction id error"
	FindVendorFeatureByIDErrorMessage                   = "find vendor feature by id error"
	RequestMerchantAPIHTTPErrorMessage                  = "request merchant api http error"
	ParseHeaderParamsErrorMessage                       = "parse header param error"
	ParseTypeIdentifierRequestDataErrorMessage          = "parse type identifier request data error"
	GetMerchantAPIParamsErrorMessage                    = "get merchant api error"
	StatusUnauthorized                                  = "unauthorized"
	TokenNotValid                                       = "token not valid"
	ClientNotValid                                      = "clients not valid"
	ClientHasVendorFeatureAlreadyExistsMessage          = "clients has vendor feature already exists"
	FindClientHasVendorFeatureErrorMessage              = "find clients has vendor feature error"
	ClientHasVendorFeatureNotFoundMessage               = "clients has vendor feature not found"
	CreateClientHasVendorFeatureErrorMessage            = "create client has vendor feature error"
	DeleteClientHasVendorFeatureErrorMessage            = "delete client has vendor feature error"
	CreateClientErrorMessage                            = "create clients error"
	UpdateClientErrorMessage                            = "update client error"
	DeleteClientErrorMessage                            = "delete client error"
	CreateVendorErrorMessage                            = "create vendor error"
	CreateVendorEnvironmentErrorMessage                 = "create vendor environment error"
	CreateVendorFeatureErrorMessage                     = "create vendor feature error"
	DeleteVendorEnvironmentErrorMessage                 = "delete vendor Environment error"
	DeleteVendorFeatureErrorMessage                     = "delete vendor Feature error"
	IDRequired                                          = "ID Required"
	UpdateVendorErrorMessage                            = "update vendor error"
	GetClientsErrorMessage                              = "get clients error"
	GetVendorsErrorMessage                              = "get vendors error"
	GetVendorEnvironmentsErrorMessage                   = "get vendor environments error"
	GetVendorFeaturesErrorMessage                       = "get vendor features error"
	GetClientHasVendorFeaturesErrorMessage              = "get client has vendor features error"
	ClientIDAlreadyExists                               = "client_id already exists"
	ClientNameAlreadyExists                             = "client name already exists"
	VendorNameAlreadyExists                             = "vendor name already exists"
	VendorEnvironmentNameAlreadyExists                  = "vendor environment name already exists"
	CreateClientHasBillingErrorMessage                  = "create client has billing error"
	UpdateClientHasBillingErrorMessage                  = "update client has billing error"
	GetClientHasBillingsErrorMessage                    = "get client has billings error"
	DeleteClientWithVendorFeatureErrorMessage           = "delete client billing with vendor feature error"
	DeleteClientHasBillingErrorMessage                  = "delete client has billing error"
	VendorFeatureExistsErrorMessage                     = "vendor feature type already exists"
	ClientHasBillingExistsErrorMessage                  = "client billing already exists"
	ClientHasVendorFeatureExistsErrorMessage            = "client has vendor feature already exists"
	ClientBillingWithVendorFeatureExistsErrorMessage    = "client billing with vendor feature already exists"
)
