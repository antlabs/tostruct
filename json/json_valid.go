// antlabs, guonaihong 2023
// apache 2.0
package json

import (
	"bytes"
	"encoding/json"
)

func Valid[T string | []byte](t T) bool {
	var x interface{} = t

	var b []byte
	switch v := x.(type) {
	case string:
		// TODO使用强制类型转换
		b = []byte(v)
	case []byte:
		b = v
	}

	b = bytes.TrimSpace(b)
	if len(b) <= 1 {
		return false
	}

	if !(b[0] == '{' && b[len(b)-1] == '}' || b[0] == '[' && b[len(b)-1] == ']') {
		return false
	}

	return json.Valid(b)
}
