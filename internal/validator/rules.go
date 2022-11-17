// Package validator
package validator

import (
	"regexp"

	"github.com/go-ozzo/ozzo-validation/v4"

)

// ValidateAlphaNumericDash for reference transaction id
func ValidateAlphaNumericDash(v string) bool {
	pattern := `^[0-9a-zA-Z\-]+$`

	rgx, err := regexp.Compile(pattern)

	if err != nil {
		return false
	}

	return rgx.MatchString(v)
}

func AlphaNumericDash() validation.StringRule {
	return validation.NewStringRuleWithError(
		ValidateAlphaNumericDash,
		validation.NewError("validation_is_alphanumeric", "must contain alpha, digits and dash only"))
}

