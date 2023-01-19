// antlabs, guonaihong 2022
// apache 2.0
package fromhead

import (
	"go/format"
	"testing"

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
				"Bool    bool    `header:\"Bool\"`\n" +
				"Float64 float64 `header:\"Float64\"`\n" +
				"Int     int     `header:\"Int\"`\n" +
				"String  string  `header:\"String\"`\n" +
				"}",
			head: []string{"int", "1", "float64", "1.1", "bool", "true", "string", "hello"},
		},
	} {
		gs := New("test")
		for i := 0; i < len(tc.head); i += 2 {
			gs.AppendHeader([]byte(tc.head[i]), []byte(tc.head[i+1]))
		}

		res, err := gs.Gen()

		assert.NoError(t, err)
		b, err := format.Source([]byte(tc.need))
		assert.NoError(t, err)
		assert.Equal(t, string(b), res)

	}
}
