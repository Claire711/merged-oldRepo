package appTypes

import "encoding/json"

type Register int

const (
	Email Register = iota
	QQ
)

// MarshalJSON 实现了 json.Marshaler 接口
func (r Register) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

// UnmarshalJSON 实现了 json.Unmarshaler 接口
func (r *Register) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	*r = ToRegister(str)
	return nil
}

// (r Register) 表示该方法属于 Register 类型，其中 r 是方法内部使用的实例对象
// String() string 表示该方法不接受参数，返回一个 string 类型的值。
func (r Register) String() string {
	var str string
	switch r {
	case Email:
		str = "邮箱"
	case QQ:
		str = "QQ"
	default:
		str = "未知"
	}
	return str
}

func ToRegister(str string) Register {
	switch str {
	case "邮箱":

		return Email
	case "QQ":

		return QQ
	default:
		return -1
	}
}

/*三种光标状态：(按 Insert 键 切换模式)
光标类型	显示效果	触发条件
插入模式	细白竖条（默认）	正常输入状态
覆盖模式	白色方块（覆盖字符）	按 Insert 键激活
选择模式	反色背景（选中文本）	拖动选择文本时*/
