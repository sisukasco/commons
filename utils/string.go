package utils

import (
	"encoding/base64"
	"encoding/json"
	"strings"
)

//CleanupString trims string and makes all lowercase
func CleanupString(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func ToJSONString(obj interface{}) string {
	bjson, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return ""
	}
	return string(bjson)
}

func MakeToken(obj interface{}) string {
	bLoginInfo, err := json.Marshal(obj)
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(bLoginInfo)
}

func AddUnique(arr []string, str string) []string {
	for _, s := range arr {
		if s == str {
			return arr
		}
	}
	return append(arr, str)
}

func RemoveString(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func SearchString(arr []string, needle string) bool {
	for s := 0; s < len(arr); s++ {
		if arr[s] == needle {
			return true
		}
	}
	return false
}
