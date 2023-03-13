package json

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/antlabs/gstl/mapex"
	"github.com/antlabs/tostruct/option"
	"github.com/stretchr/testify/assert"
)

const (
	dataFormat     = "../testdata/json/data%d.json"
	needFormat1    = "../testdata/json/need%d.json.0.struct"
	needFormat2    = "../testdata/json/need%d.json.1.struct"
	needPtrFormat1 = "../testdata/json/need%d.json.0.ptr.struct"
	needPtrFormat2 = "../testdata/json/need%d.json.1.ptr.struct"
)

func Test_Gen_Obj_JSONEx(t *testing.T) {
	for i := 0; i < 1; i++ {
		data, err := os.ReadFile(fmt.Sprintf(dataFormat, i))
		assert.NoError(t, err)

		need1, err := os.ReadFile(fmt.Sprintf(needFormat1, i))
		assert.NoError(t, err)

		need2, err := os.ReadFile(fmt.Sprintf(needFormat2, i))
		assert.NoError(t, err)

		needPtr1, err := os.ReadFile(fmt.Sprintf(needPtrFormat1, i))
		assert.NoError(t, err)

		needPtr2, err := os.ReadFile(fmt.Sprintf(needPtrFormat2, i))
		assert.NoError(t, err)

		var m map[string]any

		err = json.Unmarshal(data, &m)
		assert.NoError(t, err)

		keys := mapex.Keys(m)
		for i, v := range keys {
			keys[i] = "." + v
		}

		var out bytes.Buffer
		got1, err := Marshal(data, option.WithStructName("reqName"), option.WithOutputFmtBefore(&out))
		assert.NoError(t, err, out.String())
		out.Reset()

		got2, err := Marshal(data, option.WithStructName("reqName"), option.WithNotInline(), option.WithOutputFmtBefore(&out))
		assert.NoError(t, err, out.String())

		need1 = bytes.TrimSpace(need1)
		need2 = bytes.TrimSpace(need2)

		got1 = bytes.TrimSpace(got1)
		got2 = bytes.TrimSpace(got2)

		assert.Equal(t, need1, got1)
		assert.Equal(t, need2, got2)

		gotPtr1, err := Marshal(data, option.WithStructName("reqName"), option.WithOutputFmtBefore(&out), option.WithUsePtrType(keys))
		assert.NoError(t, err, out.String())
		out.Reset()

		gotPtr2, err := Marshal(data, option.WithStructName("reqName"), option.WithNotInline(), option.WithOutputFmtBefore(&out), option.WithUsePtrType(keys))
		assert.NoError(t, err, out.String())

		needPtr1 = bytes.TrimSpace(needPtr1)
		needPtr2 = bytes.TrimSpace(needPtr2)

		gotPtr1 = bytes.TrimSpace(gotPtr1)
		gotPtr2 = bytes.TrimSpace(gotPtr2)

		assert.Equal(t, string(needPtr1), string(gotPtr1))
		assert.Equal(t, string(needPtr2), string(gotPtr2))
	}
}
