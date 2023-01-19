package fromjson

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

type FromJSON struct {
	obj    interface{}
	Indent int
	buf    bytes.Buffer
}

const (
	structStart     = "type %s struct {\n"
	structEnd       = "}"
	startArrayStart = "%s []"
	startMap        = "%s struct {\n"
	emptyMap        = "%s struct {" +
		"} `json:\"%s\"`" +
		"}"
	keyName = "%s %s `json:\"%s\"`"
	endMap  = "} `json:\"%s\"`"
)

// TODO:
// 分开结构体
// 分开结构体，同名情况
func New(jsonBytes []byte, structName string) (f *FromJSON, err error) {
	var o map[string]interface{}
	jsonBytes = bytes.TrimSpace(jsonBytes)
	if b := valid(jsonBytes); !b {
		return nil, fmt.Errorf("tostruct:Not qualified json")
	}

	var a []interface{}

	rv := &FromJSON{}

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

func (f *FromJSON) Marshal() (b []byte, err error) {
	f.marshalValue("", f.obj, false, 0)
	f.buf.WriteString(structEnd)
	return f.buf.Bytes(), nil
}

func (f *FromJSON) marshalMap(key string, m map[string]interface{}, depth int) {

	buf := &f.buf
	remaining := len(m)

	if remaining == 0 {
		buf.WriteString(fmt.Sprintf(emptyMap, key, key))
		return
	}

	keys := make([]string, 0)
	for key := range m {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	if len(key) > 0 {
		buf.WriteString(fmt.Sprintf(startMap, key))
	}

	for _, key := range keys {

		f.writeIndent(buf, depth+1)

		f.marshalValue(key, m[key], false, depth+1)

		f.writeObjSep(buf)
	}

	f.writeIndent(buf, depth)
	if len(key) > 0 {
		buf.WriteString(fmt.Sprintf(endMap, key))
	}
}

func (f *FromJSON) marshalArray(key string, a []interface{}, depth int) {
	buf := &f.buf
	if len(a) == 0 {
		buf.WriteString(fmt.Sprintf("%s interface{} `json:\"json:%s\"`", key, key))
		return
	}

	f.marshalValue(key, a[0], true, depth)
}

func (f *FromJSON) marshalValue(key string, obj interface{}, fromArray bool, depth int) {
	buf := &f.buf
	typePrefix := ""
	if fromArray {
		typePrefix = "[]"
	}

	switch v := obj.(type) {
	case map[string]interface{}:
		f.marshalMap(key, v, depth)
	case []interface{}:
		f.marshalArray(key, v, depth)
	case string:
		buf.WriteString(fmt.Sprintf("%s %sstring `json:\"%s\"`", key, typePrefix, key))
	case float64:
		// int
		if float64(int(v)) == v {
			buf.WriteString(fmt.Sprintf("%s %sint `json:\"%s\"`", key, typePrefix, key))
			return
		}

		// float64
		buf.WriteString(fmt.Sprintf("%s %sfloat64 `json:\"%s\"`", key, typePrefix, key))
	case bool:
		buf.WriteString(fmt.Sprintf("%s %sbool `json:\"%s\"`", key, typePrefix, key))
	case nil:
		buf.WriteString(fmt.Sprintf("%s interface{} `json:\"%s\"`", key, key))
	}
}

func (f *FromJSON) writeIndent(buf *bytes.Buffer, depth int) {
	buf.WriteString(strings.Repeat(" ", f.Indent*depth))
}

func (f *FromJSON) writeObjSep(buf *bytes.Buffer) {
	if f.Indent != 0 {
		buf.WriteByte('\n')
	} else {
		buf.WriteByte(' ')
	}
}
