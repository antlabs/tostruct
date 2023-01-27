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
### 2.1 option.WithSpecifyType 指定生成类型
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
// 默认data成员会转成结构体。有时需要它转成map[string]string类型，这
all, err := Marshal([]byte(obj), option.WithStructName("reqName"), option.WithTagName("json"), option.WithSpecifyType(map[string]string{
	".data": "map[string]string",
}))

// output all
type reqName struct {
	Action    string            `json:"action"`
	Count     int               `json:"count"`
	Data      map[string]string `json:"data"`
	Duration  int               `json:"duration"`
	entities  interface{}       `json:"json:entities"`
	Timestamp int               `json:"timestamp"`
	URI       string            `json:"uri"`
}

```