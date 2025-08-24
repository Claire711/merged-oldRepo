package other

import "encoding/json"

type Data struct {
	ID  *string         `json:"id"`  //类型加*的意思是允许该字段的值为 nil
	Doc json.RawMessage `json:"doc"` //存储原始 JSON 数据，延迟解析
}
type ESIndexResponse struct {
	Data []Data `json:"data"`
}
