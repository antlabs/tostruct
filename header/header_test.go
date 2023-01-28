// antlabs, guonaihong 2023
// apache 2.0
package header

import (
	"bytes"
	"fmt"
	"go/format"
	"net/http"
	"os"
	"testing"

	"github.com/antlabs/tostruct/option"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	need string
	head []string
}

func TestHead2struct(t *testing.T) {
	for _, tc := range []testCase{
		{
			need: "type test struct{\n" +
				"Bool    bool    `header:\"bool\"`\n" +
				"Float64 float64 `header:\"float64\"`\n" +
				"Int     int     `header:\"int\"`\n" +
				"String  string  `header:\"string\"`\n" +
				"}",
			head: []string{"int", "1", "float64", "1.1", "bool", "true", "string", "hello"},
		},
	} {
		h := make(http.Header)
		for i := 0; i < len(tc.head); i += 2 {
			h.Set(tc.head[i], tc.head[i+1])
		}

		res, err := Marshal(h, option.WithStructName("test"))

		assert.NoError(t, err)
		b, err := format.Source([]byte(tc.need))
		fmt.Println(string(b))
		assert.NoError(t, err)
		assert.Equal(t, string(b), string(res))

	}
}

func TestHead2struct2(t *testing.T) {
	for _, tc := range [][]string{
		{"Content-Type", "application/json", "Accept", "application/json"},
	} {
		h := make(http.Header)
		for i := 0; i < len(tc); i += 2 {
			h.Set(tc[i], tc[i+1])
		}

		res, err := Marshal(h, option.WithStructName("test"), option.WithTagNameFromKey())
		assert.NoError(t, err)

		all, err := os.ReadFile("../testdata/testheader.0.txt")
		assert.NoError(t, err)
		assert.Equal(t, string(bytes.TrimSpace(all)), string(res))

	}
}

func TestHead2struct3(t *testing.T) {
	for _, h := range []http.Header{
		{
			"Content-Type": []string{"application/json"},
			"Accept":       []string{"application/json"},
		},
	} {

		res, err := Marshal(h, option.WithStructName("test"), option.WithTagNameFromKey())
		assert.NoError(t, err)

		all, err := os.ReadFile("../testdata/testheader.0.txt")
		assert.NoError(t, err)
		assert.Equal(t, string(bytes.TrimSpace(all)), string(res))

	}
}
