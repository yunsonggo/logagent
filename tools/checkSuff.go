package tools

import "strings"

func CheckSuffix(key string) string {
	if strings.HasSuffix(key, "/") == false {
		key = key + "/"
	}
	return key
}
