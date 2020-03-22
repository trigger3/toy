package util

import "strconv"

func IsNumeric(str string) bool {
	_, err := strconv.ParseInt(str, 10, 64)
	if err == nil {
		return true
	}

	return false
}
