package helper

import (
	"github.com/iancoleman/strcase"
	"regexp"
	"sync"
)

var camelSnakeCaseRE *regexp.Regexp
var camelSnakeCaseREOnce sync.Once

func IsSnakeCase(s string) bool {

	camelSnakeCaseREOnce.Do(func() {
		camelSnakeCaseRE = regexp.MustCompile("^[A-Z][a-z]+(_[A-Z][a-z]+)*$")
	})

	return camelSnakeCaseRE.MatchString(s)
}

func ToSnakeCase(s string) string {
	if !IsSnakeCase(s) {
		return strcase.ToSnake(s)
	}
	return s
}

