# tostruct
[![Go](https://github.com/antlabs/tostruct/workflows/Go/badge.svg)](https://github.com/antlabs/tostruct/actions)
[![codecov](https://codecov.io/gh/antlabs/tostruct/branch/master/graph/badge.svg)](https://codecov.io/gh/antlabs/tostruct)

json/yaml/http header/query string 转成struct定义，免去手写struct的烦恼.    


## 一、json/yaml/http header/query string 生成结构体
### 1.1 json字符串生成结构体
```go
import (
    "github.com/antlabs/tostruct/json"
    "github.com/antlabs/tostruct/option"
)

func main() {
	var str = `{
	"action": "Deactivate user",
	"entities": [
		{
		"uuid": "4759aa70-XXXX-XXXX-925f-6fa0510823ba",
		"type": "user",
		"created": 1542595573399,
		"modified": 1542597578147,
		"username": "user1",
		"activated": false,
		"nickname": "user"
		}],
		"timestamp": 1542602157258,
		"duration": 12
	}`

	// 父子结构合在一起
	all, _ := json.Marshal([]byte(str), option.WithStructName("reqName"))
	fmt.Println(string(all))
	/*
	type reqName struct {
		Action   string `json:"action"`
		Duration int    `json:"duration"`
		Entities []struct {
			Activated bool   `json:"activated"`
			Created   int    `json:"created"`
			Modified  int    `json:"modified"`
			Nickname  string `json:"nickname"`
			Type      string `json:"type"`
			Username  string `json:"username"`
			UUID      string `json:"uuid"`
		} `json:"entities"`
		Timestamp int `json:"timestamp"`
	}
	*/

	// 子结构拆分
	all, _ := json.Marshal([]byte(str), option.WithStructName("reqName"), option.WithNotInline())
	fmt.Println(string(all))
	/*
	type reqName struct {
		Action    string     `json:"action"`
		Duration  int        `json:"duration"`
		Entities  []Entities `json:"entities"`
		Timestamp int        `json:"timestamp"`
	}

	type Entities struct {
		Activated bool   `json:"activated"`
		Created   int    `json:"created"`
		Modified  int    `json:"modified"`
		Nickname  string `json:"nickname"`
		Type      string `json:"type"`
		Username  string `json:"username"`
		UUID      string `json:"uuid"`
	}
	*/
}
```

### 1.2 http header生成结构体
```go
import (
    "github.com/antlabs/tostruct/header"
    "github.com/antlabs/tostruct/option"
	"net/http"
)

func main() {
	h := http.header{
		"bool": []string{"true"},
		"int" : []string{"1"},
		"string" : []string{"hello"},
		"float64" : []string{"3.14"},
	}
	res, err := header.Marshal(h, option.WithStructName("test"), option.WithTagName("header"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(string(res))
	/*
type test struct {
	Bool    bool    `header:"bool"`
	Float64 float64 `header:"float64"`
	Int     int     `header:"int"`
	String  string  `header:"string"`
}
	*/
}
```
### 1.3、查询字符串生成结构体
```go
import (
    "github.com/antlabs/tostruct/url"
    "github.com/antlabs/tostruct/option"
)

func main() {
	url := "http://127.0.0.1:8080?int=1&float64=1.1&bool=true&string=hello"
	res, err := header.Marshal(h, option.WithStructName("test"), option.WithTagName("form"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(res))
	/*
type test struct {
	Bool    bool    `form:"bool"`
	Float64 float64 `form:"float64"`
	Int     int     `form:"int"`
	String  string  `form:"string"`
}
	*/
}

```

### 1.4 yaml生成结构体
```go
import (
    "github.com/antlabs/tostruct/url"
    "github.com/antlabs/tostruct/option"
	"github.com/antlabs/tostruct/yaml"
)

func main() {
	str := `
a: 1
b: 3.14
c: hello
d:
  - aa
  - bb
e:
  - a: 1
    b: 3.14
    c: hello
f:
  first: 1
  second: 3.14
`

	all, err := Marshal([]byte(str), option.WithStructName("reqName"))
	if err != nil {
		return
	}
	fmt.Println(all)
	/*
type reqName struct {
	A int      `yaml:"a"`
	B float64  `yaml:"b"`
	C string   `yaml:"c"`
	D []string `yaml:"d"`
	E []struct {
		A int     `yaml:"a"`
		B float64 `yaml:"b"`
		C string  `yaml:"c"`
	} `yaml:"e"`
	F struct {
		First  int     `yaml:"first"`
		Second float64 `yaml:"second"`
	} `yaml:"f"`
}

	*/

	all, err := Marshal([]byte(str), option.WithStructName("reqName"), option.WithNotInline())
	if err != nil {
		return
	}
	fmt.Println(all)
/*
type reqName struct {
	A int      `yaml:"a"`
	B float64  `yaml:"b"`
	C string   `yaml:"c"`
	D []string `yaml:"d"`
	E []E      `yaml:"e"`
	F F        `yaml:"f"`
}

type E struct {
	A int     `yaml:"a"`
	B float64 `yaml:"b"`
	C string  `yaml:"c"`
}

type F struct {
	First  int     `yaml:"first"`
	Second float64 `yaml:"second"`
}
*/
}
```

## 二、各种配置项函数用法
### 2.1 option.WithSpecifyType 指定生成类型(支持json/yaml)
```go
obj := `
{
  "action": "get",
  "count": 0,
  "data": {
    "a123": "delivered"
  },
  "duration": 5,
  "entities": [],
  "timestamp": 1542601830084,
  "uri": "http://XXXX/XXXX/XXXX/users/user1/offline_msg_status/123"
}
  `
// 默认对象是转成结构体的，比如这里的data成员，有些接口返回的key是变动的，比如key是用户名value是消息是否投递成功(假装在im系统中)
// 此类业务就需要转成map[string]string类型，这里就可以用上option.WithSpecifyType
// 传参的类型是map[string]string， key为json路径，value值是要指定的类型
all, err := Marshal([]byte(obj), option.WithStructName("reqName"), option.WithTagName("json"), option.WithSpecifyType(map[string]string{
	".data": "map[string]string",
}))

// output all
type reqName struct {
	Action    string            `json:"action"`
	Count     int               `json:"count"`
	Data      map[string]string `json:"data"`
	Duration  int               `json:"duration"`
	Entities  interface{}       `json:"entities"`
	Timestamp int               `json:"timestamp"`
	URI       string            `json:"uri"`
}

```
### 2.2 option.WithTagNameFromKey tag直接使用key的名字(支持http header)
http header比较特殊，传递的标准header头会有-符。可以使用option.WithTagNameFromKey()指定直接使用key的名字当作tag名(header:"xxx")这里的xxx
```go
h := http.Header{
	"Content-Type":  []string{"application/json"},
	"Accept": []string{"application/json"},
}

res, err := http.Marshal(h, option.WithStructName("test"), option.WithTagNameFromKey())

// output res
type test struct {
	Accept      string `header:"Accept"`
	ContentType string `header:"Content-Type"`
}
```

### 2.3 option.WithOutputFmtBefore 支持(json/yaml/http header/query string)
获取格式代码之前的buf输出，方便debug
```go

obj := `{"a":"b"}`
var out bytes.Buffer
all, err := json.Marshal([]byte(obj), option.WithStructName("reqName"), option.WithTagName("json"), option.WithOutputFmtBefore(&out))
_ = all
_ = err
```

### 2.4 option.WithGetValue 支持(json/yaml)
把json串序列化成struct字符串的时候，顺带提取json字符串里面指定的key值
```go
obj := `{"a":"b"}`
getValue := map[string]string{
  ".a": "",
}

_, err := json.Marshal([]byte(obj), option.WithStructName("reqName"), option.WithTagName("json"), option.WithGetValue(getValue))
fmt.Println(getValue[".a"], "b"))
```
