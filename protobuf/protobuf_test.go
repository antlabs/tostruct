package protobuf

import (
	"os"
	"testing"

	"github.com/antlabs/tostruct/option"
	"github.com/stretchr/testify/assert"
)

// obj
func Test_Protobuf_OBJ(t *testing.T) {
	all, err := Marshal(map[string]any{"int": 0, "float64": 3.3, "string": "hello"}, option.WithStructName("reqProtobuf"))
	assert.NoError(t, err)
	os.Stdout.Write(all)
}

// obj with array
func Test_Protobuf_OBJWithArray(t *testing.T) {
	all, err := Marshal(map[string]any{"int": 0, "float64": 3.3, "string": "hello", "zarray": []any{"aa", "bb", "cc"}}, option.WithStructName("reqProtobuf"))
	assert.NoError(t, err)
	os.Stdout.Write(all)
}

// obj
func Test_Protobuf_OBJWithArray2(t *testing.T) {
	all, err := Marshal(map[string]any{"int": 0, "float64": 3.3, "string": "hello", "zarray": []any{
		map[string]any{
			"hello": "world",
		},
	}}, option.WithStructName("reqProtobuf"))
	assert.NoError(t, err)
	os.Stdout.Write(all)
}
