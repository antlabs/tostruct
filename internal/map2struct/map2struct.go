// antlabs, guonaihong 2022
// apache 2.0
package map2struct

import (
	"bytes"
	"fmt"
	"go/format"
	"sort"
	"strings"

	"github.com/antlabs/tostruct/internal/guesstype"
	"github.com/antlabs/tostruct/option"
)

func MapGenStruct(m map[string][]string, opt option.Option) (res []byte, err error) {
	var out bytes.Buffer

	structName := opt.StructName
	tagName := opt.Tag
	fmt.Fprintf(&out, "type %s struct {", structName)
	var ks []string
	for k := range m {
		ks = append(ks, k)
	}

	// 排序
	sort.StringSlice(ks).Sort()

	for _, k := range ks {
		v := m[k]
		tagStr := fmt.Sprintf("`%s:%q`", tagName, k)
		k = strings.Title(k)
		if len(v) == 0 {
			fmt.Fprintf(&out, "%s string %s\n", k, tagStr)
			continue
		}

		if len(v) > 1 {
			fmt.Fprintf(&out, "%s []%s %s\n", k, guesstype.TypeOf(v[0]), tagStr)
			continue
		}

		fmt.Fprintf(&out, "%s %s %s\n", k, guesstype.TypeOf(v[0]), tagStr)
	}

	fmt.Fprint(&out, "}")

	src, err := format.Source(out.Bytes())
	if err != nil {
		return nil, err
	}

	return src, nil
}
