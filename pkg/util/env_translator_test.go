// Package util
package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvironmentTransform(t *testing.T) {
	t.Parallel()
	testCase := map[string]string{
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

	for k, v := range testCase {
		assert.Equal(t, v, EnvironmentTransform(k))
	}
}
