package option

type OptionFunc func(c *Option)

type Option struct {
	Inline     bool
	Tag        string
	StructName string
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
