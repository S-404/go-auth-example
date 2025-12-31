package validation

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func validateRegex(pattern string) validator.Func {
	return func(fl validator.FieldLevel) bool {
		value, ok := fl.Field().Interface().(string)
		if ok {
			match, _ := regexp.MatchString(pattern, value)
			return match
		}
		return true
	}
}

func validateRegexPatterns(patterns []string) validator.Func {
	return func(fl validator.FieldLevel) bool {
		value, ok := fl.Field().Interface().(string)
		if ok {
			isPassed := true
			for _, test := range patterns {
				match, _ := regexp.MatchString(test, value)
				if !match {
					isPassed = false
					break
				}
			}
			return isPassed
		}
		return true
	}
}
