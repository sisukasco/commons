package http_utils

import (
	"net/http"
	"strconv"
)

func GetQueryValue(r *http.Request, key string, def int32) int32 {
	val := def
	sval := r.URL.Query().Get(key)
	if len(sval) > 0 {
		l, err := strconv.Atoi(sval)
		if err != nil {
			val = def
		}
		val = int32(l)
	}
	return val
}
