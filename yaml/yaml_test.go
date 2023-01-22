package yaml

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/antlabs/tostruct/option"
	"github.com/stretchr/testify/assert"
)

func Test_Gen_Obj_YAML(t *testing.T) {
	obj := `
a: 1
b: 3.14
c: hello
d:
  - aa
  - bb
e:
  - a: 1
    b: 3.14
    c: hello
f:
  first: 1
  second: 3.14

  `
	for k, v := range [][]byte{
		func() []byte {
			all, err := Marshal([]byte(obj), option.WithStructName("reqName"))
			assert.NoError(t, err)
			return all
		}(),
		func() []byte {
			all, err := Marshal([]byte(obj), option.WithStructName("reqName"), option.WithNotInline())
			assert.NoError(t, err)
			return all
		}(),
	} {

		all := v
		fmt.Println(string(all))
		need, err := os.ReadFile(fmt.Sprintf("../testdata/testyaml1.%d.txt", k))
		assert.NoError(t, err)
		assert.Equal(t, string(bytes.TrimSpace(need)), string(bytes.TrimSpace(all)))
	}
}
