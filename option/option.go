// antlabs, guonaihong 2023
// apache 2.0
package option

type OptionFunc func(c *Option)

type Option struct {
	Inline         bool
	Tag            string
	StructName     string
	TypeMap        map[string]string
	TagNameFromKey bool
}

// 控制生成的结构体是否内联
func WithNotInline() OptionFunc {
	return func(c *Option) {
		c.Inline = false
	}
}

// 设置tag名, 修改third里面的json字段
/*
type Third struct {
	B1 string `json:"b1"`
	B2 string `json:"b2"`
}
*/
func WithTagName(name string) OptionFunc {
	return func(c *Option) {
		c.Tag = name
	}
}

// 设置最外层结构体的名字，WithStructName("Third")
/*
type Third struct {
	B1 string `json:"b1"`
	B2 string `json:"b2"`
}
*/
func WithStructName(name string) OptionFunc {
	return func(c *Option) {
		c.StructName = name
	}
}

// 指定类型, datal默认转成struct， 这里直接指定生成map[string]string类型
// {
//
//		"data" : {
//		  "user1": "111"
//		}
//	}
//
//	WithSpecifyType(map[string]string{
//	   ".data": "map[string]string"
//	})
//
// 目前只支持json/yaml
func WithSpecifyType(typeMap map[string]string) OptionFunc {
	return func(c *Option) {
		c.TypeMap = typeMap
	}
}

// tag使用变量名, http header特殊一点
// 目前仅仅支持http header marshal
func WithTagNameFromKey() OptionFunc {
	return func(c *Option) {
		c.TagNameFromKey = true
	}
}
