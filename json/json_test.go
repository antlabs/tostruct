// antlabs, guonaihong 2023
// apache 2.0
package json

import (
	"bytes"
	"fmt"
	"os"
	"testing"

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
	all, err := Marshal([]byte(obj), WithStructName("reqName"), WithTagName("json"))
	assert.NoError(t, err)
	//fmt.Println(string(all))

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
	all, err := Marshal([]byte(obj), WithStructName("reqName"), WithTagName("json"), WithNotInline())
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
	all, err := Marshal([]byte(obj), WithStructName("reqName"), WithTagName("json"), WithNotInline())
	assert.NoError(t, err)

	fmt.Println(string(all))
	need, err := os.ReadFile("../testdata/test3.txt")
	assert.NoError(t, err)
	assert.Equal(t, string(bytes.TrimSpace(need)), string(all))
}

func Test_Gen_Obj_JSON4(t *testing.T) {

	obj := `
 [{
   "first" : true,
   "second" : 0,
   "third" : 1.1,
   "fourth" : null
}]
  `
	all, err := Marshal([]byte(obj), WithStructName("reqName"), WithTagName("json"), WithNotInline())
	assert.NoError(t, err)

	fmt.Println(string(all))
	need, err := os.ReadFile("../testdata/test4.txt")
	assert.NoError(t, err)
	assert.Equal(t, string(bytes.TrimSpace(need)), string(all))
}
