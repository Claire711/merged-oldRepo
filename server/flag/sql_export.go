package flag

import (
	"fmt"
	"os"
	"os/exec"
	"server/global"
	"time"
)

func SQLExport() error {
	//从 Docker 容器中导出 MySQL 数据库
	mysql := global.Config.Mysql
	timer := time.Now().Format("20060102")
	sqlPath := fmt.Sprintf("mysql_%s.sql", timer)
	//mysqldump：MySQL 导出工具,在cmd命令行执行docker命令
	cmd := exec.Command("docker", "exec", "mysql", "mysqldump", "-u"+mysql.Username, "-p"+mysql.Password, mysql.DBName)
	outFile, err := os.Create(sqlPath)
	if err != nil {
		return err
	}
	defer outFile.Close()
	cmd.Stdout = outFile
	return cmd.Run()
}

/*
"20250717" 是硬编码的日期，不会随时间变化。
应使用动态格式（如 20060102 是 Go 的日期模板）
*/
/*
改进方向：
1.密码安全问题
2.容器名称硬编码
3.错误处理不足
*/
