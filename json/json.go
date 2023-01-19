// antlabs, guonaihong 2023
// apache 2.0
package json

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/gobeam/stringy"
)

type JSON struct {
	obj        interface{}  // json/yaml 解成map[string]interface{} 或者[]interface{}
	Indent     int          // 控制输出缩进
	buf        bytes.Buffer // 存放内联结构体的数据
	inline     bool         // 是否内联
	tag        string       // json, yaml, 默认json
	structName string       //最外层结构体名
}

const (
	structStart     = "type %s struct {\n"
	structEnd       = "}"
	startArrayStart = "%s []"
	startInlineMap  = "%s struct {\n"
	endInlineMap    = "} `json:\"%s\"`"
	emptyMap        = "%s struct {" +
		"} `json:\"%s\"`" +
		"}"
	keyName = "%s %s `json:\"%s\"`"
)

// TODO:
// 分开结构体
// 分开结构体，同名情况
func New(jsonBytes []byte, structName string, tag string) (f *JSON, err error) {
	var o map[string]interface{}
	jsonBytes = bytes.TrimSpace(jsonBytes)
	if b := Valid(jsonBytes); !b {
		return nil, fmt.Errorf("tostruct:Not qualified json")
	}

	var a []interface{}

	rv := &JSON{}

	rv.buf.WriteString(fmt.Sprintf(structStart, structName))
	if jsonBytes[0] == '{' {
		json.Unmarshal(jsonBytes, &o)
		rv.obj = o
	} else if jsonBytes[0] == '[' {
		json.Unmarshal(jsonBytes, &a)
		rv.obj = a
	}

	rv.Indent = 4
	return rv, nil
}

func (f *JSON) Marshal() (b []byte, err error) {
	f.marshalValue("", f.obj, false, 0)
	f.buf.WriteString(structEnd)
	return f.buf.Bytes(), nil
}

func (f *JSON) marshalMap(key string, m map[string]interface{}, depth int) {

	buf := &f.buf
	remaining := len(m)

	fieldName, tagName := getFieldAndTagName(key)
	if remaining == 0 {
		buf.WriteString(fmt.Sprintf(emptyMap, fieldName, tagName))
		return
	}

	keys := make([]string, 0)
	for key := range m {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	if len(key) > 0 {
		buf.WriteString(fmt.Sprintf(startInlineMap, fieldName))
	}

	for _, key := range keys {

		f.writeIndent(buf, depth+1)

		f.marshalValue(key, m[key], false, depth+1)

		f.writeObjSep(buf)
	}

	f.writeIndent(buf, depth)
	if len(key) > 0 {
		buf.WriteString(fmt.Sprintf(endInlineMap, tagName))
	}
}

func (f *JSON) marshalArray(key string, a []interface{}, depth int) {
	buf := &f.buf
	if len(a) == 0 {
		buf.WriteString(fmt.Sprintf("%s interface{} `json:\"json:%s\"`", key, key))
		return
	}

	f.marshalValue(key, a[0], true, depth)
}

func getFieldAndTagName(key string) (string, string) {
	str := stringy.New(key)
	fieldName := str.CamelCase("?", "")
	tagName := str.SnakeCase("?", "").ToLower()
	return fieldName, tagName
}

func (f *JSON) marshalValue(key string, obj interface{}, fromArray bool, depth int) {
	buf := &f.buf
	typePrefix := ""
	if fromArray {
		typePrefix = "[]"
	}

	fieldName, tagName := getFieldAndTagName(key)

	switch v := obj.(type) {
	case map[string]interface{}:
		f.marshalMap(key, v, depth)
	case []interface{}:
		f.marshalArray(key, v, depth)
	case string:
		buf.WriteString(fmt.Sprintf("%s %sstring `json:\"%s\"`", fieldName, typePrefix, tagName))
	case float64:
		// int
		if float64(int(v)) == v {
			buf.WriteString(fmt.Sprintf("%s %sint `json:\"%s\"`", fieldName, typePrefix, tagName))
			return
		}

		// float64
		buf.WriteString(fmt.Sprintf("%s %sfloat64 `json:\"%s\"`", fieldName, typePrefix, tagName))
	case bool:
		buf.WriteString(fmt.Sprintf("%s %sbool `json:\"%s\"`", fieldName, typePrefix, tagName))
	case nil:
		buf.WriteString(fmt.Sprintf("%s interface{} `json:\"%s\"`", fieldName, tagName))
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
