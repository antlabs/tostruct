// antlabs, guonaihong 2023
// apache 2.0
package json

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/antlabs/tostruct/option"
	"github.com/stretchr/testify/assert"
)

func Test_Gen_Obj_JSON(t *testing.T) {
	obj := `
{
   "first" : ["a", "b"],
   "second" : {"b1" : "b1", "b2" : "b2"},
   "third" : [{"b1" : "b1", "b2" : "b2"}]
}
  `
	all, err := Marshal([]byte(obj), option.WithStructName("reqName"), option.WithTagName("json"))
	assert.NoError(t, err)
	fmt.Println(string(all))

	need, err := os.ReadFile("../testdata/test1.txt")
	//fmt.Println(string(all))

	assert.NoError(t, err)
	assert.Equal(t, string(bytes.TrimSpace(need)), string(all))
}

func Test_Gen_Obj_JSON2(t *testing.T) {
	obj := `
{
   "first" : ["a", "b"],
   "second" : {"b1" : "b1", "b2" : "b2"},
   "third" : [{"b1" : "b1", "b2" : "b2"}]
}
  `
	all, err := Marshal([]byte(obj), option.WithStructName("reqName"), option.WithTagName("json"), option.WithNotInline())
	assert.NoError(t, err)

	fmt.Println(string(all))
	need, err := os.ReadFile("../testdata/test2.txt")
	assert.NoError(t, err)
	assert.Equal(t, string(need), string(all))
}

func Test_Gen_Obj_JSON3(t *testing.T) {

	obj := `
{
   "first" : true,
   "second" : 0,
   "third" : 1.1,
   "fourth" : null
}
  `
	all, err := Marshal([]byte(obj), option.WithStructName("reqName"), option.WithTagName("json"), option.WithNotInline())
	assert.NoError(t, err)

	fmt.Println(string(all))
	need, err := os.ReadFile("../testdata/test3.txt")
	assert.NoError(t, err)
	assert.Equal(t, string(bytes.TrimSpace(need)), string(all))
}

// 数组对象
func Test_Gen_Obj_JSON4(t *testing.T) {

	obj := `
 [{
   "first" : true,
   "second" : 0,
   "third" : 1.1,
   "fourth" : null
}]
  `
	all, err := Marshal([]byte(obj), option.WithStructName("reqName"), option.WithTagName("json"), option.WithNotInline())
	assert.NoError(t, err)

	fmt.Println(string(all))
	need, err := os.ReadFile("../testdata/test4.txt")
	assert.NoError(t, err)
	assert.Equal(t, string(bytes.TrimSpace(need)), string(all))
}

func Test_Gen_Obj_JSON5(t *testing.T) {

	obj := `
  {
  "action": "Deactivate user",
  "entities": [
    {
      "uuid": "4759aa70-XXXX-XXXX-925f-6fa0510823ba",
      "type": "user",
      "created": 1542595573399,
      "modified": 1542597578147,
      "username": "user1",
      "activated": false,
      "nickname": "user"
      }  ],
      "timestamp": 1542602157258,
      "duration": 12}

  `

	for k, v := range [][]byte{
		func() []byte {
			all, err := Marshal([]byte(obj), option.WithStructName("reqName"), option.WithTagName("json"))
			assert.NoError(t, err)
			return all
		}(),
		func() []byte {
			all, err := Marshal([]byte(obj), option.WithStructName("reqName"), option.WithTagName("json"), option.WithNotInline())
			assert.NoError(t, err)
			return all
		}(),
	} {

		all := v
		fmt.Println(string(all))
		need, err := os.ReadFile(fmt.Sprintf("../testdata/test5.%d.txt", k))
		assert.NoError(t, err)
		assert.Equal(t, string(bytes.TrimSpace(need)), string(bytes.TrimSpace(all)))
	}

}

func Test_Gen_Obj_JSON6(t *testing.T) {

	obj := `
  {
  "a" : 1,
  "b" : 3.14,
  "c" : "hello",
  "d" : ["aa", "bb"],
  "e" : [{
  "a" : 1,
  "b" : 3.14,
  "c" : "hello"
}],
  "f" : {
"first" : 1,
"second": 3.14
}
}
  `

	for k, v := range [][]byte{
		func() []byte {
			all, err := Marshal([]byte(obj), option.WithStructName("reqName"), option.WithTagName("json"))
			assert.NoError(t, err)
			return all
		}(),
		func() []byte {
			all, err := Marshal([]byte(obj), option.WithStructName("reqName"), option.WithTagName("json"), option.WithNotInline())
			assert.NoError(t, err)
			return all
		}(),
	} {

		all := v
		fmt.Println(string(all))
		need, err := os.ReadFile(fmt.Sprintf("../testdata/test6.%d.txt", k))
		assert.NoError(t, err)
		assert.Equal(t, string(bytes.TrimSpace(need)), string(bytes.TrimSpace(all)))
	}

}
