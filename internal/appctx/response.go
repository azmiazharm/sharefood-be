// Package appctx
package appctx

import (
	"encoding/json"
	"sync"

	"sharefood/internal/consts"
	"sharefood/pkg/msg"
)

var rsp *Response
var oneRsp sync.Once

// Response presentation contract object
type Response struct {
	Code    int         `json:"-"`
	Status  string      `json:"status,omitempty"`
	Entity  string      `json:"entity,omitempty"`
	State   string      `json:"state,omitempty"`
	Message interface{} `json:"message,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	lang    string      `json:"-"`
	msgKey  string
}

// MetaData represent meta data response for multi data
type MetaData struct {
	Page       uint64 `json:"page"`
	Limit      uint64 `json:"limit"`
	TotalPage  uint64 `json:"total_page"`
	TotalCount uint64 `json:"total_count"`
}

// GetMessage method to transform response name var to message detail
func (r *Response) GetMessage() string {
	return msg.Get(r.msgKey, r.lang).Text()
}

// Generate setter message
func (r *Response) Generate() *Response {
	if r.lang == "" {
		r.lang = consts.LangDefault
	}
	msg := msg.Get(r.msgKey, r.lang)
	if r.Message == nil {
		r.Message = msg.Text()
	}

	if r.Code == 0 {
		r.Code = msg.Status()
	}

	return r
}

// WithCode setter response var name
func (r *Response) WithCode(c int) *Response {
	r.Code = c
	return r
}

// With setter status response
func (r *Response) WithStatus(s string) *Response {
	r.Status = s
	return r
}

// With setter entity response
func (r *Response) WithEntity(s string) *Response {
	r.Entity = s
	return r
}

// With setter state response
func (r *Response) WithState(s string) *Response {
	r.State = s
	return r
}

// WithData setter data response
func (r *Response) WithData(v interface{}) *Response {
	r.Data = v
	return r
}

// WithError setter error messages
func (r *Response) WithError(v interface{}) *Response {
	r.Errors = v
	return r
}

func (r *Response) WithMsgKey(v string) *Response {
	r.msgKey = v
	return r
}

// WithMeta setter meta data response
func (r *Response) WithMeta(v interface{}) *Response {
	r.Meta = v
	return r
}

// WithLang setter language response
func (r *Response) WithLang(v string) *Response {
	r.lang = v
	return r
}

// WithMessage setter custom message response
func (r *Response) WithMessage(v interface{}) *Response {
	if v != nil {
		r.Message = v
	}

	return r
}

func (r *Response) Byte() []byte {
	if r.Code == 0 || r.Message == nil {
		r.Generate()
	}

	b, _ := json.Marshal(r)
	return b
}

// NewResponse initialize response
func NewResponse() *Response {
	oneRsp.Do(func() {
		rsp = &Response{}
	})

	// clone response
	x := *rsp

	return &x
}
