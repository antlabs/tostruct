package json

import (
	"bytes"
	"fmt"
)

const (
	// protobuf message开始
	messageStart               = "message %s {\n"      // ok
	messageStartArrayStruct    = "message %s {\n"      //
	messageEndStruct           = "}"                   // TODO
	messageStartInlineMapAfter = "%s%s %s = %d;\n"     // repeated? int32 id = 4;
	messageStartInlineMap      = "message %s {\n"      // 内联message体开始
	messageStartMap            = "%s %s%s `%s:\"%s\"`" // 拆开结构体开始, TODO
	messageEndMap              = "}\n"                 // 拆开结构体结束, TODO
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
		"} `%s:\"%s\"`"
	defStructName = "AutoGenerated"
	nilFmt        = "%s interface{} `%s:\"%s\"`"
	stringFmt     = "%s %sstring `%s:\"%s\"`"
	stringPtrFmt  = "%s %s*string `%s:\"%s\"`"
	boolFmt       = "%s %sbool `%s:\"%s\"`"
	boolPtrFmt    = "%s %s*bool `%s:\"%s\"`"
	float64Fmt    = "%s %sfloat64 `%s:\"%s\"`"
	float64PtrFmt = "%s %s*float64 `%s:\"%s\"`"
	intFmt        = "%s %sint `%s:\"%s\"`"
	intPtrFmt     = "%s %s*int `%s:\"%s\"`"
	specifytFmt   = "%s %s `%s:\"%s\"`"
)

func (j *JSON) getSpecifytFmt(buf *bytes.Buffer, fieldName string, fieldType string, tagName string, id int, ptr bool) {

	if j.IsProtobuf {
		buf.WriteString(fmt.Sprintf(messageSpecifytFmt, fieldType, fieldName, id))
		return
	}

	// TODO 指针类型
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

func (j *JSON) getInt(buf *bytes.Buffer, fieldName string, typePrefix string, tagName string, id int, ptr bool) {
	if j.IsProtobuf {

		buf.WriteString(fmt.Sprintf(messageIntFmt, typePrefix, fieldName, id))
		return
	}

	format := intFmt
	if ptr {
		format = intPtrFmt
	}

	buf.WriteString(fmt.Sprintf(format, fieldName, typePrefix, j.Tag, tagName))
	return
}

func (j *JSON) getFloat64(buf *bytes.Buffer, fieldName string, typePrefix string, tagName string, id int, ptr bool) {
	if j.IsProtobuf {

		buf.WriteString(fmt.Sprintf(messageFloat64Fmt, typePrefix, fieldName, id))
		return
	}

	format := float64Fmt
	if ptr {
		format = float64PtrFmt
	}
	buf.WriteString(fmt.Sprintf(format, fieldName, typePrefix, j.Tag, tagName))
	return
}

func (j *JSON) getBool(buf *bytes.Buffer, fieldName string, typePrefix string, tagName string, id int, ptr bool) {
	if j.IsProtobuf {

		buf.WriteString(fmt.Sprintf(messageBoolFmt, typePrefix, fieldName, id))
		return
	}

	format := boolFmt
	if ptr {
		format = boolPtrFmt
	}
	buf.WriteString(fmt.Sprintf(format, fieldName, typePrefix, j.Tag, tagName))
	return
}

func (j *JSON) getString(buf *bytes.Buffer, fieldName string, typePrefix string, tagName string, id int, ptr bool) {
	if j.IsProtobuf {

		buf.WriteString(fmt.Sprintf(messageStringFmt, typePrefix, fieldName, id))
		return
	}

	format := stringFmt
	if ptr {
		format = stringPtrFmt
	}
	buf.WriteString(fmt.Sprintf(format, fieldName, typePrefix, j.Tag, tagName))
	return
}

func (j *JSON) getStartInlineMap(buf *bytes.Buffer, fieldName string, typePrefix string, depth int) {
	if j.IsProtobuf {

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
