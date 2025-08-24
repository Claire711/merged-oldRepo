package flag

import (
	"os"
	"server/global"
	"strings"
)

func SQLImport(sqlPath string) (errs []error) {
	// 读取指定路径的 SQL 文件，按分号分割成多条 SQL 语句并逐条执行，返回执行过程中遇到的所有错误。
	byteData, err := os.ReadFile(sqlPath)
	if err != nil {
		return append(errs, err)
	}
	sqlList := strings.Split(string(byteData), ";")
	for _, sql := range sqlList {
		//strings.TrimSpace 是 Go 语言中一个常用的字符串处理函数，用于移除字符串开头和结尾的所有空白字符。
		sql = strings.TrimSpace(sql)
		if sql == "" {
			continue
		}
		//通过链式调用先执行 SQL，再获取可能的错误
		err = global.DB.Exec(sql).Error
		if err != nil {
			errs = append(errs, err)
			continue
		}
	}
	return nil
}
