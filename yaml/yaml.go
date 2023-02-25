// antlabs, guonaihong 2023
// apache 2.0
package yaml

import (
	"github.com/antlabs/tostruct/json"
	"github.com/antlabs/tostruct/option"
	"gopkg.in/yaml.v3"
)

func Unmarshal(bytes []byte) (o map[string]any, a []any, err error) {

	err = yaml.Unmarshal(bytes, &o)
	if err != nil {
		if err = yaml.Unmarshal(bytes, &a); err != nil {
			return nil, nil, err
		}
	}
	return
}

func Marshal(bytes []byte, opt ...option.OptionFunc) (b []byte, err error) {

	o, a, err := Unmarshal(bytes)
	if err != nil {
		return nil, err
	}

	opt = append(opt, option.WithTagName("yaml"))
	if o != nil {
		return json.Marshal(o, opt...)
	}

	return json.Marshal(a, opt...)
}
