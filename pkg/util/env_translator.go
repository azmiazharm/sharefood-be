// Package util
package util

import (
	"strings"
)

var envArr = map[string]string{
	"production":  "production",
	"staging":     "staging",
	"development": "development",
	"prod":        "production",
	"stg":         "staging",
	"dev":         "development",
	"prd":         "production",
	"green":       "green",
	"blue":        "blue",
}

// EnvironmentTransform transformer
func EnvironmentTransform(s string) string {
	v, ok := envArr[strings.ToLower(strings.Trim(s, " "))]

	if !ok {
		return ""
	}

	return v
}
