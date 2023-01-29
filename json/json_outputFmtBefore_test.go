package json

import (
	"bytes"
	"testing"

	"github.com/antlabs/tostruct/option"
	"github.com/stretchr/testify/assert"
)

func Test_OutputFmtBefore(t *testing.T) {

	obj := `{"a":"b"}`
	var out bytes.Buffer
	all, err := Marshal([]byte(obj), option.WithStructName("reqName"), option.WithTagName("json"), option.WithOutputFmtBefore(&out))
	assert.NoError(t, err)
	assert.NotEqual(t, out.Len(), 0)
	assert.NotEqual(t, 0, len(all))
}
