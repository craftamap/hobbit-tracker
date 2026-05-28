package httputil

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func GetUint64PathValue(req *http.Request, name string) (uint64, bool) {
	value := req.PathValue(name)
	if value == "" {
		return 0, false
	}
	intValue, err := strconv.ParseUint(value, 10, 64)
	return intValue, err == nil
}


func GetUuidPathValue(req *http.Request, name string) (string, bool) {
	value := req.PathValue(name)
	if value == "" {
		return "", false
	}
	pattern := regexp.MustCompile("^[a-zA-Z0-9\\-]+$")
	matches := pattern.MatchString(value)
	if !matches {
		return "", false
	}
	return value, true
}

func IsContentType(contentTypeHeader string, contentType string) bool {
	if i := strings.IndexRune(contentTypeHeader, ';'); i != -1 {
		contentTypeHeader = contentTypeHeader[0:i]
	}
	return contentTypeHeader == contentType
}
