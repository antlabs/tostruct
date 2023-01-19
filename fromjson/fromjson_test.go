// guonaihong 2023
// apache 2.0
package fromjson

import (
	"fmt"
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
	j, err := New([]byte(obj), "reqName")
	assert.NoError(t, err)
	all, err := j.Marshal()
	assert.NoError(t, err)
	fmt.Println(string(all))
}
