package utils

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// GenerateVerificationCode 生成一个指定长度的随机验证码
func GenerateVerificationCode(len int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%0*d", len, r.Intn(int(math.Pow10(len))))
}

/*
UnixNano()：是 time.Time 类型的方法，
会把当前时间转换成从 UTC 时间 1970 年 1 月 1 日 00:00:00（ Unix 纪元起点 ）到当前时间所经过的纳秒数，
返回结果是 int64 类型的整数。
*/
