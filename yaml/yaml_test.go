// antlabs, guonaihong 2023
// apache 2.0
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

// 测试文件以 data编号命名
// 生成后的结构体文件以 need编码.struct
const (
	dataFormat  = "../testdata/data%d.yaml"
	needFormat1 = "../testdata/need%d.yaml.0.struct"
	needFormat2 = "../testdata/need%d.yaml.1.struct"
)

func Test_Gen_Obj_YAML2(t *testing.T) {
	for i := 0; i < 1; i++ {
		data, err := os.ReadFile(fmt.Sprintf(dataFormat, i))
		assert.NoError(t, err)

		need1, err := os.ReadFile(fmt.Sprintf(needFormat1, i))
		assert.NoError(t, err)

		need2, err := os.ReadFile(fmt.Sprintf(needFormat2, i))
		assert.NoError(t, err)

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
	}
}
