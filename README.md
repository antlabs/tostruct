# tostruct
json/http header/query string 转成struct

## json字序列化成结构体
### 
```go
import (
    "github.com/antlabs/tostruct/json"
)
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
all, _ := json.Marshal([]byte(obj), json.WithStructName("reqName"), json.WithTagName("json"))
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
all, _ := json.Marshal([]byte(obj), json.WithStructName("reqName"), json.WithTagName("json"), json.WithNotInline())
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
/*
```