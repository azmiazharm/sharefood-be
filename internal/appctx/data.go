// Package appctx
package appctx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/gorilla/schema"

	"sharefood/internal/consts"
)

var jsonContent = regexp.MustCompile(`(?i)application/((\w|\.|-)+\+)?json(-seq)?`)

// Data context for http use case
type Data struct {
	Request     *http.Request
	Config      *Config
	ServiceType string
	BytesValue  []byte
}

// ConsumerData context for use case message processor
type ConsumerData struct {
	Body        []byte
	Key         []byte
	Topic       string
	Partition   int32
	TimeStamp   time.Time
	Offset      int64
	ServiceType string
	Lang        string
	Commit      func()
}

// Cast casts data based on servcice type
// args:
//
//	target: object target holder
//
// returns:
//
//	error operation
func (d *Data) Cast(target interface{}) error {

	v := reflect.ValueOf(target)
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("target %T cannot addressable, must pointer target", target)
	}

	if d.ServiceType == "" {
		return fmt.Errorf("empty service type")
	}
	switch d.ServiceType {
	case consts.ServiceTypeHTTP:
		return d.httpCast(target)
	case consts.ServiceTypeConsumer:
		return d.mqCast(target)
	default:
		return nil
	}
}

func (d *Data) httpCast(target interface{}) error {
	if d.Request == nil {
		return fmt.Errorf("unable to cast http data, null request")
	}

	// httpCast transform request payload data
	// GET -> params-query-string
	// POST -> json-body
	if err := d.grabMethod(target); err != nil {
		return err
	}
	return nil
}

func (d *Data) mqCast(target interface{}) error {
	return json.NewDecoder(bytes.NewReader(d.BytesValue)).Decode(target)
}

// Transform query-string into json struct
func (d *Data) transform(target interface{}, src map[string][]string) error {
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	decoder.SetAliasTag("url")
	if err := decoder.Decode(target, src); err != nil {
		return fmt.Errorf("unable to decode query string:%s", err.Error())
	}
	return nil
}

// Grab request method
// Take a destination source of struct
func (d *Data) grabMethod(target interface{}) error {
	if d.Request.Method != http.MethodPost && d.Request.Method != http.MethodPut && d.Request.Method != http.MethodPatch {
		return d.transform(target, d.Request.URL.Query())
	}

	cType := d.Request.Header.Get("Content-Type")
	if strings.Contains(cType, "form-urlencoded") {
		if err := d.Request.ParseForm(); err != nil {
			return err
		}

		return d.transform(target, d.Request.PostForm)
	}

	if jsonContent.MatchString(cType) {
		return d.decodeJSON(d.Request.Body, target)
	}

	return fmt.Errorf("unsupported decode content-type=%s", cType)
}

func (d *Data) decodeJSON(body io.ReadCloser, dst interface{}) error {
	if body == nil {
		return nil
	}
	err := json.NewDecoder(body).Decode(dst)
	if err != nil {
		return fmt.Errorf("unable decode request body, err:%s", err.Error())
	}

	return nil
}
