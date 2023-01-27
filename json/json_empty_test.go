package json

import (
	"testing"

	"github.com/antlabs/tostruct/option"
	"github.com/stretchr/testify/assert"
)

func TestEmpty(t *testing.T) {
	str := []byte("")
	all, err := Marshal(str, option.WithStructName("empty"), option.WithTagName("json"))
	assert.Error(t, err)
	assert.Equal(t, all, []byte(nil))
}
