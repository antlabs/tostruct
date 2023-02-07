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
	obj       any                      // json/yaml 解成map[string]any 或者[]any
	Indent    int                      // 控制输出缩进
	buf       bytes.Buffer             // 存放内联结构体的数据
	count     map[string]int           //记录message里面的最大值
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
	// protobuf message开始
	messageStart               = "message %s {\n"      // ok
	messageStartArrayStruct    = "message %s {\n"      //
	messageEndStruct           = "}"                   // TODO
	messageStartInlineMapAfter = "%s%s %s = %d;\n"     // repeated? int32 id = 4;
	messageStartInlineMap      = "message %s {\n"      // 内联message体开始
	messageStartMap            = "%s %s%s `%s:\"%s\"`" // 拆开结构体开始, TODO
	messageEndMap              = "}\n"                 // 拆开结构体结束, TODO
	messageEmptyMap            = "%s struct {" +
		"} `%s:\"%s\"`" +
		"}"
	//messageNilFmt        = "%s any `%s:\"%s\"`" //TODO protobuf 暂时忽略
	messageStringFmt   = "%sstring %s = %d;"  //repeated? int32 id = 2;
	messageBoolFmt     = "%sbool %s = %d;"    //repeated? int32 id = 2;
	messageFloat64Fmt  = "%sfloat64 %s = %d;" //repeated? int32 id = 2;
	messageIntFmt      = "%sint64 %s = %d;"   //repeated? int32 id = 2;
	messageSpecifytFmt = "%s %s = %d;"        // ok

	// json --
	startStruct      = "type %s struct {\n"
	startArrayStruct = "type %s []struct {\n"
	endStruct        = "}"
	startInlineMap   = "%s %sstruct {\n"     // 内联结构体开始
	endInlineMap     = "} `%s:\"%s\"`"       // 内联结构体结束
	startMap         = "%s %s%s `%s:\"%s\"`" // 拆开结构体开始
	emptyMap         = "%s struct {" +
		"} `%s:\"%s\"`" +
		"}"
	defStructName = "AutoGenerated"
	nilFmt        = "%s interface{} `%s:\"%s\"`"
	stringFmt     = "%s %sstring `%s:\"%s\"`"
	boolFmt       = "%s %sbool `%s:\"%s\"`"
	float64Fmt    = "%s %sfloat64 `%s:\"%s\"`"
	intFmt        = "%s %sint `%s:\"%s\"`"
	specifytFmt   = "%s %s `%s:\"%s\"`"
)

func (j *JSON) getSpecifytFmt(buf *bytes.Buffer, fieldName string, fieldType string, tagName string, id int) {

	if j.IsProtobuf {
		buf.WriteString(fmt.Sprintf(messageSpecifytFmt, fieldType, fieldName, id))
		return
	}

	buf.WriteString(fmt.Sprintf(specifytFmt, fieldName, fieldType, j.Tag, tagName))
}

// 生成结构体或者message的开头
func (j *JSON) getStructStart(buf *bytes.Buffer) {
	format := startStruct
	if j.IsProtobuf {
		format = messageStart
	}
	buf.WriteString(fmt.Sprintf(format, j.StructName))
	return
}

// array 对象外层
func (j *JSON) getArrayStart(buf *bytes.Buffer) {
	if j.IsProtobuf {
		buf.WriteString(fmt.Sprintf(messageStartArrayStruct, j.StructName))
		return
	}

	buf.WriteString(fmt.Sprintf(startArrayStruct, j.StructName))
}

func (j *JSON) getInt(buf *bytes.Buffer, fieldName string, typePrefix string, tagName string, id int) {
	if j.IsProtobuf {

		buf.WriteString(fmt.Sprintf(messageIntFmt, typePrefix, fieldName, id))
		return
	}

	buf.WriteString(fmt.Sprintf(intFmt, fieldName, typePrefix, j.Tag, tagName))
	return
}

func (j *JSON) getFloat64(buf *bytes.Buffer, fieldName string, typePrefix string, tagName string, id int) {
	if j.IsProtobuf {

		buf.WriteString(fmt.Sprintf(messageFloat64Fmt, typePrefix, fieldName, id))
		return
	}

	buf.WriteString(fmt.Sprintf(float64Fmt, fieldName, typePrefix, j.Tag, tagName))
	return
}

func (j *JSON) getBool(buf *bytes.Buffer, fieldName string, typePrefix string, tagName string, id int) {
	if j.IsProtobuf {

		buf.WriteString(fmt.Sprintf(messageBoolFmt, typePrefix, fieldName, id))
		return
	}

	buf.WriteString(fmt.Sprintf(boolFmt, fieldName, typePrefix, j.Tag, tagName))
	return
}

func (j *JSON) getString(buf *bytes.Buffer, fieldName string, typePrefix string, tagName string, id int) {
	if j.IsProtobuf {

		buf.WriteString(fmt.Sprintf(messageStringFmt, typePrefix, fieldName, id))
		return
	}

	buf.WriteString(fmt.Sprintf(stringFmt, fieldName, typePrefix, j.Tag, tagName))
	return
}

func (j *JSON) getStartInlineMap(buf *bytes.Buffer, fieldName string, typePrefix string, depth int) {
	if j.IsProtobuf {

		// 这里是message 开始和上个普通成员缩进保持一致
		// int32 x = 3;
		// message xx {
		// }
		//j.writeIndent(buf, depth)
		buf.WriteString(fmt.Sprintf(messageStartInlineMap, fieldName))
		return
	}

	buf.WriteString(fmt.Sprintf(startInlineMap, fieldName, typePrefix))
}

func (j *JSON) getEndInlineMap(buf *bytes.Buffer, fieldName string, typePrefix string, tagName string, depth int, id int) {
	if j.IsProtobuf {
		buf.WriteString(fmt.Sprintf(messageEndMap))

		j.writeIndent(buf, depth)
		buf.WriteString(fmt.Sprintf(messageStartInlineMapAfter, typePrefix, fieldName, tagName, id))
		return
	}

	buf.WriteString(fmt.Sprintf(endInlineMap, j.Tag, tagName))
}

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

		var a []any
		var o map[string]any

		rv = newDefault()
		for _, o := range opt {
			o(&rv.Option)
		}

		if jsonBytes[0] == '{' {
			rv.getStructStart(&rv.buf)
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
		rv.getStructStart(&rv.buf)
		rv.obj = data
	case []any:
		rv = newDefault()
		for _, o := range opt {
			o(&rv.Option)
		}
		rv.getArrayStart(&rv.buf)
		rv.obj = data
	}

	rv.Indent = 4
	return rv, nil
}

func (f *JSON) marshal() (b []byte, err error) {
	key := ""
	_, isArr := f.obj.([]any)
	fromArray := f.IsProtobuf && isArr
	depth := 0
	if fromArray {
		// protobuf不支持像json一样的顶层是[]数组的对象，如果要把类似的结构转成protobuf.
		// 先包装一个message，其中的成员是data
		key = "Data"
		//外层已经是message，里面的成员都要缩进
		// message req {
		//     message Data {
		//       int32 x = 1;
		//     }
		//     repeated Data data = 1;
		// }
		depth = 1
		// 临时补丁, 格式化protubf
		f.writeIndent(&f.buf, depth)
	}
	id := 1
	f.marshalValue(key, f.obj, fromArray, depth, &f.buf, "", &id)
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

	if f.IsProtobuf {
		return f.buf.Bytes(), nil
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

func (f *JSON) marshalMap(key string, m map[string]any,
	typePrefix string,
	depth int,
	buf *bytes.Buffer,
	pathKey string,
	id *int,
) {

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
			//f.writeIndent(buf, depth)
			f.getStartInlineMap(buf, fieldName, typePrefix, depth)
			// 如果是内嵌结构体
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

		f.marshalValue(key, m[key], false, depth+1, buf, appendKeyPath(pathKey, key), id)
		(*id)++

		f.writeObjSep(buf)
	}

	f.writeIndent(buf, depth)
	if len(key) > 0 {
		if f.Inline {
			f.getEndInlineMap(buf, fieldName, typePrefix, tagName, depth, *id)
		} else {
			buf.WriteString(endStruct + "\n")

		}
	}
}

func (f *JSON) marshalArray(key string, a []any, depth int, buf *bytes.Buffer, keyPath string, id *int) {

	f.marshalValue(key, a[0], true, depth, buf, keyPath, id)
}

func (f *JSON) marshalValue(key string, obj any, fromArray bool, depth int, buf *bytes.Buffer, keyPath string, id *int) {
	typePrefix := ""
	if fromArray {
		typePrefix = "[]"
		if f.IsProtobuf {
			typePrefix = "repeated " //这里的空格是为了格式化效果的
		}
	}

	fieldName, tagName := name.GetFieldAndTagName(key)

	// 重写key的类型
	if f.TypeMap != nil {
		fieldType, ok := f.TypeMap[keyPath]
		if ok {
			f.getSpecifytFmt(buf, fieldName, fieldType, tagName, *id)
			return
		}
	}

	if f.GetValue != nil {
		_, ok := f.GetValue[keyPath]
		if ok {
			f.GetValue[keyPath] = fmt.Sprintf("%s", obj)
		}
	}

	if f.GetRawValue != nil {
		_, ok := f.GetRawValue[keyPath]
		if ok {
			f.GetRawValue[keyPath] = obj
		}
	}

	tmpFieldName := strings.ToUpper(fieldName)
	if tab.InitialismsTab[tmpFieldName] {
		fieldName = tmpFieldName
	}

	switch v := obj.(type) {
	case map[string]any:
		f.marshalMap(key, v, typePrefix, depth, buf, keyPath, id)
	case []any:
		if len(v) == 0 {
			if !f.IsProtobuf {
				buf.WriteString(fmt.Sprintf("%s interface{} `json:\"%s\"`", fieldName, key))
			}
			return
		}
		f.marshalArray(key, v, depth, buf, keyPath+"[0]", id)
	case string:
		f.getString(buf, fieldName, typePrefix, tagName, *id)
	case float64: //json默认解析的数字是float64类型
		// int
		if float64(int(v)) == v {
			f.getInt(buf, fieldName, typePrefix, tagName, *id)
			return
		}

		// float64
		f.getFloat64(buf, fieldName, typePrefix, tagName, *id)
	case int: //yaml解析成map[string]any，数值是int类型
		f.getInt(buf, fieldName, typePrefix, tagName, *id)
	case bool:
		f.getBool(buf, fieldName, typePrefix, tagName, *id)
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
