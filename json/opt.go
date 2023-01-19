package json

type JSONConfig func(f *JSON)

// 控制生成的结构体是否内联
func WithNotInline() JSONConfig {
	return func(f *JSON) {
		f.inline = false
	}
}

// 设置tag名, 修改third里面的json字段
/*
type Third struct {
	B1 string `json:"b1"`
	B2 string `json:"b2"`
}
*/
func WithTagName(name string) JSONConfig {
	return func(f *JSON) {
		f.tag = name
	}
}

// 设置最外层结构体的名字，WithStructName("Third")
/*
type Third struct {
	B1 string `json:"b1"`
	B2 string `json:"b2"`
}
*/
func WithStructName(name string) JSONConfig {
	return func(f *JSON) {
		f.structName = name
	}
}
