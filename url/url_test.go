// antlabs, guonaihong 2023
// apache 2.0
package url

import (
	"fmt"
	"go/format"
	"testing"

	"github.com/antlabs/tostruct/option"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	need  string
	query string
}

func TestUrl2struct(t *testing.T) {
	for _, tc := range []testCase{
		{
			need: "type test struct{\n" +
				"Bool    bool    `form:\"bool\"`\n" +
				"Float64 float64 `form:\"float64\"`\n" +
				"Int     int     `form:\"int\"`\n" +
				"String  string  `form:\"string\"`\n" +
				"}",

			query: "http://127.0.0.1:8080?int=1&float64=1.1&bool=true&string=hello",
		},
		{
			need: "type test struct{\n" +
				"HotwordId    int    `form:\"hotword_id\"`\n" +
				"Token  string  `form:\"token\"`\n" +
				"}",

			query: "/a/v1/b/c/d?token=a&hotword_id=1",
		},
	} {
		res, err := Marshal(tc.query, option.WithStructName("test"))
		fmt.Println(string(res))

		assert.NoError(t, err)
		b, err := format.Source([]byte(tc.need))
		assert.NoError(t, err)
		assert.Equal(t, string(b), string(res))

	}
}
