package protobuf

import (
	"github.com/antlabs/tostruct/json"
	"github.com/antlabs/tostruct/option"
)

type Type interface {
	map[string]any | []any
}

func Marshal[T Type](t T, opt ...option.OptionFunc) (b []byte, err error) {

	return json.Marshal(t, append(opt, option.WithProtobuf())...)
}
