package yaml

import (
	"github.com/antlabs/tostruct/json"
	"github.com/antlabs/tostruct/option"
	"gopkg.in/yaml.v3"
)

func Marshal(bytes []byte, opt ...option.OptionFunc) (b []byte, err error) {
	var o map[string]any
	var a []any

	err = yaml.Unmarshal(bytes, &o)
	if err != nil {
		if err = yaml.Unmarshal(bytes, &a); err != nil {
			return nil, err
		}

	}

	if a != nil {
		return json.Marshal(a, opt...)
	}

	opt = append(opt, option.WithTagName("yaml"))
	return json.Marshal(o, opt...)
}
