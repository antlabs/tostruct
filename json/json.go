// antlabs, guonaihong 2023
// apache 2.0
package json

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/format"
	"sort"
	"strings"

	"github.com/antlabs/gstl/mapex"
	"github.com/antlabs/tostruct/internal/tab"
	"github.com/antlabs/tostruct/name"
	"github.com/antlabs/tostruct/option"
)

type JSON struct {
	option.Option
	obj       interface{}              // json/yaml 解成map[string]interface{} 或者[]interface{}
	Indent    int                      // 控制输出缩进
	buf       bytes.Buffer             // 存放内联结构体的数据
	structBuf map[string]*bytes.Buffer // 记录拆分结构体
}

type Type interface {
	[]byte | map[string]any | []any
}

// 原始json
/*
{
   "first" : ["a", "b"],
   "second" : {"b1" : "b1", "b2" : "b2"},
   "third" : [{"b1" : "b1", "b2" : "b2"}]
}
*/

// 生成拆开的结构体
/*
type AutoGenerated struct {
	First  []string `json:"first"`
	Second Second   `json:"second"`
	Third  []Third  `json:"third"`
}
type Second struct {
	B1 string `json:"b1"`
	B2 string `json:"b2"`
}
type Third struct {
	B1 string `json:"b1"`
	B2 string `json:"b2"`
}
*/

// 生成内联结构体
/*
	type AutoGenerated struct {
		First  []string `json:"first"`
		Second struct {
			B1 string `json:"b1"`
			B2 string `json:"b2"`
		} `json:"second"`
		Third []struct {
			B1 string `json:"b1"`
			B2 string `json:"b2"`
		} `json:"third"`
	}
*/

const (
	startStruct      = "type %s struct {\n"
	startArrayStruct = "type %s []struct {\n"
	endStruct        = "}"
	startArrayStart  = "%s []"
	startInlineMap   = "%s %sstruct {\n"     // 内联结构体开始
	endInlineMap     = "} `%s:\"%s\"`"       // 内联结构体结束
	startMap         = "%s %s%s `%s:\"%s\"`" // 拆开结构体开始
	endMap           = "}"                   // 拆开结构体结束
	emptyMap         = "%s struct {" +
		"} `%s:\"%s\"`" +
		"}"
	keyName       = "%s %s `%s:\"%s\"`"
	defStructName = "AutoGenerated"
	nilFmt        = "%s interface{} `%s:\"%s\"`"
	stringFmt     = "%s %sstring `%s:\"%s\"`"
	boolFmt       = "%s %sbool `%s:\"%s\"`"
	float64Fmt    = "%s %sfloat64 `%s:\"%s\"`"
	intFmt        = "%s %sint `%s:\"%s\"`"
	specifytFmt   = "%s %s `%s:\"%s\"`"
)

func Marshal[T Type](t T, opt ...option.OptionFunc) (b []byte, err error) {
	f, err := new(t, opt...)
	if err != nil {
		return nil, err
	}

	return f.marshal()
}

func newDefault() *JSON {
	return &JSON{structBuf: make(map[string]*bytes.Buffer),
		Option: option.Option{Tag: "json", StructName: defStructName, Inline: true}}

}

func new[T Type](t T, opt ...option.OptionFunc) (f *JSON, err error) {
	var tmp any = t
	var rv *JSON

	switch data := tmp.(type) {

	case []byte:
		jsonBytes := bytes.TrimSpace(data)

		if b := Valid(jsonBytes); !b {
			return nil, fmt.Errorf("tostruct.json:Not qualified json")
		}

		var a []interface{}
		var o map[string]interface{}

		rv = newDefault()
		for _, o := range opt {
			o(&rv.Option)
		}

		if jsonBytes[0] == '{' {
			rv.buf.WriteString(fmt.Sprintf(startStruct, rv.StructName))
			json.Unmarshal(jsonBytes, &o)
			rv.obj = o
		} else if jsonBytes[0] == '[' {
			rv.buf.WriteString(fmt.Sprintf(startArrayStruct, rv.StructName))
			json.Unmarshal(jsonBytes, &a)
			rv.obj = a
		}
	case map[string]any:
		rv = newDefault()
		for _, o := range opt {
			o(&rv.Option)
		}
		rv.buf.WriteString(fmt.Sprintf(startStruct, rv.StructName))
		rv.obj = data
	case []any:
		rv = newDefault()
		for _, o := range opt {
			o(&rv.Option)
		}
		rv.buf.WriteString(fmt.Sprintf(startArrayStruct, rv.StructName))
		rv.obj = data
	}

	rv.Indent = 4
	return rv, nil
}

func (f *JSON) marshal() (b []byte, err error) {
	f.marshalValue("", f.obj, false, 0, &f.buf, "")
	f.buf.WriteString(endStruct)
	if !f.Inline {
		keys := mapex.Keys(f.structBuf)
		sort.Strings(keys)

		for i, v := range keys {
			if i == 0 {
				f.buf.WriteByte('\n')
			}
			f.buf.WriteString(f.structBuf[v].String())
		}
	}

	if f.OutputFmtBefore != nil {
		f.OutputFmtBefore.Write(f.buf.Bytes())
	}

	if b, err = format.Source(f.buf.Bytes()); err != nil {
		fmt.Printf("%s\n", f.buf.String())
		return nil, err
	}

	return b, nil
}

func (f *JSON) getStructTypeName(fieldName string) (structTypeName string, buf *bytes.Buffer) {

	structTypeName = fieldName
	for count := 0; ; count++ {

		if _, ok := f.structBuf[structTypeName]; ok {
			// 比较少见的情况， 结构体里面有重名变量
			// 使用fieldName + 数字编号的形式解决重名问题
			structTypeName = fmt.Sprintf("%s%d", fieldName, count)
			continue
		}
		buf = bytes.NewBuffer([]byte{})
		f.structBuf[structTypeName] = buf
		return
	}
}

func appendKeyPath(pathKey string, key string) string {
	return fmt.Sprintf("%s.%s", pathKey, strings.ToLower(key))
}

func (f *JSON) marshalMap(key string, m map[string]interface{}, typePrefix string, depth int, buf *bytes.Buffer, pathKey string) {

	remaining := len(m)

	fieldName, tagName := name.GetFieldAndTagName(key)
	if remaining == 0 {
		buf.WriteString(fmt.Sprintf(emptyMap, fieldName, f.Tag, tagName))
		return
	}

	keys := make([]string, 0)
	for key := range m {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	if len(key) > 0 {
		if f.Inline {
			// 如果是内嵌结构体
			buf.WriteString(fmt.Sprintf(startInlineMap, fieldName, typePrefix))
		} else {
			// 生成struct类型名和子结构体可以保存的子buf
			structTypeName, buf2 := f.getStructTypeName(fieldName)
			// 保存子结构体声明语句， type Third struct {
			buf2.WriteString(fmt.Sprintf("\n"+startStruct, structTypeName))
			// 在父结构体里，写入声明类型 Third  []Third  `json:"third"`
			buf.WriteString(fmt.Sprintf(startMap, fieldName, typePrefix, structTypeName, f.Tag, tagName))
			depth = 0
			buf = buf2
		}
	}

	for _, key := range keys {

		f.writeIndent(buf, depth+1)

		f.marshalValue(key, m[key], false, depth+1, buf, appendKeyPath(pathKey, key))

		f.writeObjSep(buf)
	}

	f.writeIndent(buf, depth)
	if len(key) > 0 {
		if f.Inline {
			buf.WriteString(fmt.Sprintf(endInlineMap, f.Tag, tagName))
		} else {
			buf.WriteString(endStruct + "\n")

		}
	}
}

func (f *JSON) marshalArray(key string, a []interface{}, depth int, buf *bytes.Buffer, keyPath string) {

	f.marshalValue(key, a[0], true, depth, buf, keyPath)
}

func (f *JSON) marshalValue(key string, obj interface{}, fromArray bool, depth int, buf *bytes.Buffer, keyPath string) {
	typePrefix := ""
	if fromArray {
		typePrefix = "[]"
	}

	fieldName, tagName := name.GetFieldAndTagName(key)

	// 重写key的类型
	if f.TypeMap != nil {
		fieldType, ok := f.TypeMap[keyPath]
		if ok {
			buf.WriteString(fmt.Sprintf(specifytFmt, fieldName, fieldType, f.Tag, tagName))
			return
		}
	}

	if f.GetValue != nil {
		_, ok := f.GetValue[keyPath]
		if ok {
			f.GetValue[keyPath] = fmt.Sprintf("%s", obj)
		}

	}

	tmpFieldName := strings.ToUpper(fieldName)
	if tab.InitialismsTab[tmpFieldName] {
		fieldName = tmpFieldName
	}

	switch v := obj.(type) {
	case map[string]interface{}:
		f.marshalMap(key, v, typePrefix, depth, buf, keyPath)
	case []interface{}:
		if len(v) == 0 {
			buf.WriteString(fmt.Sprintf("%s interface{} `json:\"%s\"`", fieldName, key))
			return
		}
		f.marshalArray(key, v, depth, buf, keyPath+"[0]")
	case string:
		buf.WriteString(fmt.Sprintf(stringFmt, fieldName, typePrefix, f.Tag, tagName))
	case float64: //json默认解析的数字是float64类型
		// int
		if float64(int(v)) == v {
			buf.WriteString(fmt.Sprintf(intFmt, fieldName, typePrefix, f.Tag, tagName))
			return
		}

		// float64
		buf.WriteString(fmt.Sprintf(float64Fmt, fieldName, typePrefix, f.Tag, tagName))
	case int: //yaml解析成map[string]any，数值是int类型
		buf.WriteString(fmt.Sprintf(intFmt, fieldName, typePrefix, f.Tag, tagName))
	case bool:
		buf.WriteString(fmt.Sprintf(boolFmt, fieldName, typePrefix, f.Tag, tagName))
	case nil:
		buf.WriteString(fmt.Sprintf(nilFmt, fieldName, f.Tag, tagName))
	}
}

func (f *JSON) writeIndent(buf *bytes.Buffer, depth int) {
	buf.WriteString(strings.Repeat(" ", f.Indent*depth))
}

func (f *JSON) writeObjSep(buf *bytes.Buffer) {
	if f.Indent != 0 {
		buf.WriteByte('\n')
	} else {
		buf.WriteByte(' ')
	}
}
