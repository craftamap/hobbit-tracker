package httputil

import (
	"net/http"
	"regexp"
	"strconv"
)

func GetUint64PathValue(req *http.Request, name string) (uint64, bool) {
	value := req.PathValue(name)
	if value == "" {
		return 0, false
	}
	intValue, err := strconv.ParseUint(value, 10, 64)
	return intValue, err == nil
}


func GetAlphanumericMinusPathValue(req *http.Request, name string) (string, bool) {
	value := req.PathValue(name)
	if value == "" {
		return "", false
	}
	pattern := regexp.MustCompile(`^[a-zA-Z0-9\-]+$`)
	matches := pattern.MatchString(value)
	if !matches {
		return "", false
	}
	return value, true
}
