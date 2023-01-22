// antlabs, guonaihong 2023
// apache 2.0
package name

import "github.com/gobeam/stringy"

func GetFieldAndTagName(key string) (string, string) {
	str := stringy.New(key)
	fieldName := str.CamelCase("?", "")
	tagName := str.SnakeCase("?", "").ToLower()
	return fieldName, tagName
}
