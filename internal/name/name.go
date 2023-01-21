package name

import "github.com/gobeam/stringy"

func GetFieldAndTagName(key string) (string, string) {
	str := stringy.New(key)
	fieldName := str.CamelCase("?", "")
	tagName := str.SnakeCase("?", "").ToLower()
	return fieldName, tagName
}
