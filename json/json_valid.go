// antlabs, guonaihong 2023
// apache 2.0
package json

import (
	"bytes"
	"encoding/json"
	"errors"
)

var ErrValid = errors.New("valid json")

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

func ValidAndUnmarshal(b []byte) (o map[string]any, a []any, err error) {
	if !Valid(b) {
		return nil, nil, ErrValid
	}

	if err := json.Unmarshal(b, &o); err == nil /*没有错误说明是json 对象字符串*/ {
		return o, nil, nil
	}

	// 可能是array对象
	err = json.Unmarshal(b, &a)
	return nil, a, err
}
