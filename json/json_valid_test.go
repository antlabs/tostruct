// antlabs, guonaihong 2023
// apache 2.0
package json

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testJSONValid struct {
	data string
	need bool
}

type testJSONValid2 struct {
	data []byte
	need bool
}

func Test_JSON_Valid(t *testing.T) {
	for _, tc := range []testJSONValid{
		{data: "[}", need: false},
		{data: "{]", need: false},
		{data: "{}", need: true},
		{data: "[]", need: true},
		{data: "", need: false},
	} {
		assert.Equal(t, Valid(tc.data), tc.need)
	}

}

func Test_JSON_Valid2(t *testing.T) {
	for i, tc := range []testJSONValid2{
		{data: []byte("[}"), need: false},
		{data: []byte("{]"), need: false},
		{data: []byte("{}"), need: true},
		{data: []byte("[]"), need: true},
		{data: []byte(""), need: false},
	} {
		assert.Equal(t, Valid(tc.data), tc.need, fmt.Sprintf("index:%d", i))
	}

}
