package middleware

import (
	"fmt"
	"net/http"
	"sharefood/internal/appctx"
	"sharefood/internal/consts"
	"sharefood/internal/response"
	"sharefood/pkg/logger"
	"sharefood/pkg/tracer"
	"strings"

	"github.com/dgrijalva/jwt-go/v4"
)

func ValidateBearerToken(w http.ResponseWriter, r *http.Request, conf *appctx.Config) error {
	errorEvent := consts.ErrorEvent("validate_bearer_token_middleware")
	response := response.NewResponse("validate_bearer_token_middleware", r)
	ctx := tracer.SpanStart(r.Context(), "validate_bearer_token_middleware")
	defer tracer.SpanFinish(ctx)

	authToken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	if len(authToken) == 0 {
		err := errorEvent.WithCode(consts.CodeBadRequest).WrapError(consts.Error(consts.StatusUnauthorized))
		logger.Error(logger.MessageFormat("[user-login] parsing body request error:"))
		return NewError(*response.Failed(ctx, nil, err))
	}

	secret := []byte(conf.App.JWTSecret)
	token, errToken := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("parsing error")
		}
		return secret, nil
	})
	if errToken != nil {
		err := errorEvent.WithCode(consts.CodeBadRequest).WrapError(errToken)
		tracer.SpanError(ctx, err)
		return NewError(*response.Failed(ctx, nil, err))
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		err := errorEvent.WithCode(consts.CodeBadRequest).WrapError(consts.Error(consts.TokenNotValid))
		tracer.SpanError(ctx, err)
		return NewError(*response.Failed(ctx, nil, err))
	}

	idUser := claims["id_user"].(string)
	r.Header.Set("idUser", idUser)

	// role := claims.(jwt.MapClaims)["role"].(string)
	// fmt.Println(id_user)
	// r.Header.Set("role", role)

	return nil
}
