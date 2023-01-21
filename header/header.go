// antlabs, guonaihong 2022
// apache 2.0
package header

import (
	"net/http"

	"github.com/antlabs/tostruct/internal/map2struct"
	"github.com/antlabs/tostruct/option"
)

func Marshal(h http.Header, opt ...option.OptionFunc) (structBytes []byte, err error) {

	var opts option.Option
	opts.Tag = "header"
	for _, o := range opt {
		o(&opts)
	}

	return map2struct.MapGenStruct(map[string][]string(h), opts)
}
