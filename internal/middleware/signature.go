// Package middleware
package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"sharefood/internal/appctx"
	"sharefood/internal/consts"
	"sharefood/pkg/hash"
	"sharefood/pkg/logger"
)

// ValidateSignature header signature
func ValidateSignature(r *http.Request, conf *appctx.Config) int {

	var err error
	key := conf.App.ApiKey

	signature := r.Header.Get("Signature")

	save := r.Body

	save, r.Body, err = drainBody(r.Body)

	nameField := "request"

	if err != nil {
		return consts.CodeBadRequest
	}

	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		logger.Warn(fmt.Sprintf("[middleware] cannot read request body , error %s", err.Error()), logger.Any(nameField, string(b)))
		return consts.CodeBadRequest
	}

	buff := &bytes.Buffer{}
	err = json.Compact(buff, b)
	if err != nil {
		logger.Error(fmt.Sprintf("[middleware] cannot compat bytes , error %s", err.Error()), logger.Any(nameField, string(b)))
		return consts.CodeBadRequest
	}

	if !hash.HmacComparator(buff.String(), signature, hash.SHA256(key)) {
		logger.Warn(fmt.Sprintf("[middleware] invalid signature %s", signature), logger.Any(nameField, string(b)))
		return consts.CodeBadRequest
	}

	r.Body = save

	return consts.CodeSuccess
}

func drainBody(b io.ReadCloser) (r1, r2 io.ReadCloser, err error) {
	if b == http.NoBody {
		// No copying needed. Preserve the magic sentinel meaning of NoBody.
		return http.NoBody, http.NoBody, nil
	}

	var buf bytes.Buffer
	if _, err = buf.ReadFrom(b); err != nil {
		return nil, b, err
	}
	if err = b.Close(); err != nil {
		return nil, b, err
	}
	return ioutil.NopCloser(&buf), ioutil.NopCloser(bytes.NewReader(buf.Bytes())), nil
}
