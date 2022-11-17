package middleware

import (
	"fmt"
	"sharefood/internal/appctx"
)

type Error struct {
	Response appctx.Response
	err      error
}

func (e Error) Error() string {
	var errString string
	if e.err != nil {
		errString = fmt.Sprintf("%s:", e.err.Error())
	}

	return fmt.Sprintf("%s%s", errString, string(e.Response.Byte()))
}

func NewError(response appctx.Response, opts ...ErrorOpt) Error {
	errOpt := ErrorOpts{error: nil}

	for _, opt := range opts {
		opt(&errOpt)
	}

	return Error{Response: response, err: errOpt.error}
}

type ErrorOpts struct {
	error error
}

type ErrorOpt func(opts *ErrorOpts)

func WithError(err error) ErrorOpt {
	return func(opts *ErrorOpts) {
		opts.error = err
	}
}
