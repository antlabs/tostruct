package protobuf

import (
	"github.com/antlabs/tostruct/json"
	"github.com/antlabs/tostruct/option"
	"github.com/antlabs/tostruct/yaml"
)

type Type interface {
	map[string]any | []any | []byte
}

func Marshal[T Type](t T, opt ...option.OptionFunc) (b []byte, err error) {
	i := any(t)
	switch v := i.(type) {
	case []byte:
		o, a, err := json.ValidAndUnmarshal(v)
		if err != nil {
			o, a, err = yaml.Unmarshal(v)
			if err != nil {
				return nil, err
			}
		}
		if o != nil {
			return json.Marshal(o, append(opt, option.WithProtobuf())...)
		}

		if a != nil {
			return json.Marshal(a, append(opt, option.WithProtobuf())...)
		}
	}

	return json.Marshal(t, append(opt, option.WithProtobuf())...)
}
