package qiniu

// MustNotEmpty 字符串为空时返回的信息
func MustNotEmpty(paramName string) string {
	return paramName + " should not be empty"
}

// MustAboveZero 整型小于等于0时候返回的信息
func MustAboveZero(paramName string) string {
	return paramName + " must be greater than 0"
}
