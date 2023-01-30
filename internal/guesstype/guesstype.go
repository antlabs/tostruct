// antlabs, guonaihong 2023
// apache 2.0
package guesstype

import (
	"strconv"
)

//使用能力检测, 查看string类型的是哪种类型

func IsInt(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func IsFloat(s string) bool {
	if IsInt(s) {
		return false
	}

	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func IsBool(s string) bool {
	_, err := strconv.ParseBool(s)
	return err == nil
}

func TypeOf(s string) string {
	if IsFloat(s) {
		return "float64"
	}

	if IsInt(s) {
		return "int"
	}

	if IsBool(s) {
		return "bool"
	}

	return "string"
}

func ToAny(s string) any {
	if IsFloat(s) {
		f, _ := strconv.ParseFloat(s, 0)
		return f
	}

	if IsInt(s) {
		i, _ := strconv.Atoi(s)
		return i
	}

	if IsBool(s) {
		b, _ := strconv.ParseBool(s)
		return b
	}

	return s
}
