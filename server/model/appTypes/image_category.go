package appTypes

import "encoding/json"

type Category int

const (
	Null Category = iota
	System
	Carousel     //轮播图
	Cover        //封面
	Illustration //插图
	AdImage      //广告
	Friend       //友链
)

func (c Category) MarshalJSON() ([]byte, error) {
	// // 返回2个值：[]byte和error
	return json.Marshal(c.String())
}
func (c *Category) UnmarshalJSON(data []byte) error {
	// 返回1个值：error
	var str string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}
	*c = ToCategory(str)
	return nil
}
func (c Category) String() string {
	switch c {
	case Null:
		return "未使用"
	case System:
		return "系统"
	case Carousel:
		return "背景"
	case Cover:
		return "封面"
	case Illustration:
		return "插图"
	case AdImage:
		return "广告"
	case Friend:
		return "友链"
	//要加个默认返回，不然会缺少返回类型，会报错
	default:
		return "未知类型"
	}
}
func ToCategory(str string) Category {
	switch str {
	case "未使用":
		return Null
	case "系统":
		return System
	case "背景":
		return Carousel
	case "封面":
		return Cover
	case "插图":
		return Illustration
	case "广告":
		return AdImage
	case "友链":
		return Friend
	default:
		return -1
	}
}

/*func (接收者) 函数名(参数列表) (返回值列表) {// 函数体}*/
