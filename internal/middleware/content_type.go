// Package middleware
package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"sharefood/internal/appctx"
	"sharefood/internal/consts"
	"sharefood/pkg/logger"
)

// ValidateContentType header
func ValidateContentType(r *http.Request, conf *appctx.Config) int {

	if ct := strings.ToLower(r.Header.Get(`Content-Type`)); ct != `application/json` {
		logger.Warn(fmt.Sprintf("[middleware] invalid content-type %s", ct))

		return consts.CodeBadRequest
	}

	return consts.CodeSuccess
}
