// antlabs, guonaihong 2023
// apache 2.0
package map2struct

import (
	"bytes"
	"fmt"
	"go/format"
	"sort"

	"github.com/antlabs/tostruct/internal/guesstype"
	"github.com/antlabs/tostruct/name"
	"github.com/antlabs/tostruct/option"
)

func getStructStart(buf *bytes.Buffer, opt option.Option) {
	if opt.IsProtobuf {
		fmt.Fprintf(buf, "message %s {\n", opt.StructName)
		return
	}

	fmt.Fprintf(buf, "type %s struct {", opt.StructName)
}

func getEmptyField(buf *bytes.Buffer, opt option.Option, fieldName string, tagStr string, id int) {

	if opt.IsProtobuf {
		fmt.Fprintf(buf, "    string %s = %d;\n", fieldName, id)
		return
	}

	fmt.Fprintf(buf, "%s string %s\n", fieldName, tagStr)
}

func getSliceField(buf *bytes.Buffer, opt option.Option, fieldName string, v string, tagStr string, id int) {

	if opt.IsProtobuf {
		fmt.Fprintf(buf, "    repeated %s %s = %d;\n", guesstype.TypeOf(v, opt.IsProtobuf), fieldName, id)
		return
	}

	fmt.Fprintf(buf, "%s []%s %s\n", fieldName, guesstype.TypeOf(v, opt.IsProtobuf), tagStr)
}

func getField(buf *bytes.Buffer, opt option.Option, fieldName string, v string, tagStr string, id int) {
	if opt.IsProtobuf {
		fmt.Fprintf(buf, "    %s %s = %d;\n", guesstype.TypeOf(v, opt.IsProtobuf), fieldName, id)
		return
	}

	fmt.Fprintf(buf, "%s %s %s\n", fieldName, guesstype.TypeOf(v, opt.IsProtobuf), tagStr)
}

func MapGenStruct(m map[string][]string, opt option.Option) (res []byte, err error) {
	var out bytes.Buffer

	tag := opt.Tag
	getStructStart(&out, opt)
	var ks []string
	for k := range m {
		ks = append(ks, k)
	}

	// 排序
	sort.StringSlice(ks).Sort()

	id := 1
	for _, k := range ks {
		v := m[k]
		if opt.GetRawValue != nil {
			_, ok := opt.GetRawValue[k]
			if ok {
				opt.GetRawValue[k] = guesstype.ToAny(v[0])

			}
		}

		fieldName, tagName := name.GetFieldAndTagName(k)
		if opt.TagNameFromKey {
			tagName = k
		}

		tagStr := fmt.Sprintf("`%s:%q`", tag, tagName)
		if len(v) == 0 {
			getEmptyField(&out, opt, fieldName, tagStr, id)
			goto next
		}

		if len(v) > 1 {
			getSliceField(&out, opt, fieldName, v[0], tagStr, id)
			goto next
		}

		getField(&out, opt, fieldName, v[0], tagStr, id)
	next:
		id++
	}

	fmt.Fprint(&out, "}")

	if opt.OutputFmtBefore != nil {
		opt.OutputFmtBefore.Write(out.Bytes())
	}

	if opt.IsProtobuf {
		return out.Bytes(), nil
	}

	src, err := format.Source(out.Bytes())
	if err != nil {
		return nil, err
	}

	return src, nil
}
