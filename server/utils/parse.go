package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ParseDuration 解析时间间隔字符串（支持d/h/m/s单位）
// 参数d: 要解析的时间字符串（如"2d3h30m"）
// 返回值: time.Duration对象和可能的错误
func ParseDuration(d string) (time.Duration, error) {
	// 1. 去除字符串首尾空格
	d = strings.TrimSpace(d)
	// 2. 检查空字符串情况
	if len(d) == 0 {
		return 0, fmt.Errorf("empty duration string")
	}
	// 3. 定义单位映射关系
	unitPattern := map[string]time.Duration{
		"d": time.Hour * 24,
		"h": time.Hour,
		"m": time.Minute,
		"s": time.Second,
	}
	// 4. 初始化总时长
	var totalDuration time.Duration
	// 5. 遍历所有时间单位
	/*代码 强制按 d → h → m → s 的顺序 处理单位，与输入字符串中的顺序无关。
	无论原始字符串是 "1d12h" 还是 "12h1d"，解析时都会 先处理 d 再处理 h。*/
	for _, unit := range []string{"d", "h", "m", "s"} {
		// 6. 循环处理字符串中所有当前单位的出现
		for strings.Contains(d, unit) {
			// 7. 找到当前单位在字符串中的位置
			unitIndex := strings.Index(d, unit)
			// 8. 提取单位前的数字部分
			/*从字符串 d 中截取从开头到 unitIndex 位置（不包含该位置）的子字符串，并赋值给变量 part*/
			part := d[:unitIndex]
			//如果数字部分为空,默认为0
			if part == "" {
				part = "0"
			}
			// 10. 将字符串转换为整数
			val, err := strconv.Atoi(part)
			// 转换失败时返回错误
			if err != nil {
				return 0, fmt.Errorf("invalid duration part: %v", err)
			}
			totalDuration += time.Duration(val) * unitPattern[unit]
			d = d[unitIndex+len(unit):]
		}
	}
	if len(d) > 0 {
		return 0, fmt.Errorf("unrecognized duration format")
	}
	return totalDuration, nil
}
