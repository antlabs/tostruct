package json

import (
	"testing"

	"github.com/antlabs/tostruct/option"
	"github.com/stretchr/testify/assert"
)

func Test_GetVal(t *testing.T) {
	obj := `{"a":"b"}`
	getValue := map[string]string{
		".a": "",
	}

	_, err := Marshal([]byte(obj), option.WithStructName("reqName"), option.WithTagName("json"), option.WithGetValue(getValue), option.WithOutputFmtBefore(nil))
	assert.NoError(t, err)
	assert.Equal(t, getValue[".a"], "b")
}
