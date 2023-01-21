// antlabs, guonaihong 2023
// apache 2.0
package url

import (
	"net/url"

	"github.com/antlabs/tostruct/internal/map2struct"
	"github.com/antlabs/tostruct/option"
)

func Marshal(rawURL string, opt ...option.OptionFunc) (structBytes []byte, err error) {

	var opts option.Option
	opts.Tag = "form"
	for _, o := range opt {
		o(&opts)
	}

	u2, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	return map2struct.MapGenStruct(map[string][]string(u2.Query()), opts)
}
