# tostruct
json/http header/query string 转成struct

## 一、json字符串生成结构体
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
	all, _ := json.Marshal([]byte(obj), option.WithStructName("reqName"))
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
	all, _ := json.Marshal([]byte(obj), option.WithStructName("reqName"), option.WithNotInline())
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

## 二、http header生成结构体
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
## 三、查询字符串生成结构体
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