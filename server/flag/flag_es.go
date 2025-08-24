package flag

import (
	"bufio"
	"fmt"
	"os"
	"server/model/elasticsearch"
	"server/service"
)

func Elasticsearch() error {
	esService := service.ServiceGroupApp.EsService
	indexExists, err := esService.IndexExists(elasticsearch.ArticleIndex())
	if err != nil {
		return err
	}
	if indexExists {
		fmt.Printf("The index already exists.Do you want to delete the data and recreate theindex? (y/n)")

		//读取用户输入
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input := scanner.Text()

		switch input {
		case "y":
			fmt.Println("Proceeding to delete the data and recreate the index...")
			//删除 Elasticsearch 中的文章索引，并且会对删除操作可能出现的错误进行检查和处理
			if err := esService.IndexDelete(elasticsearch.ArticleIndex()); err != nil {
				return err
			}
		case "n":
			fmt.Println("Exiting the program.")
			os.Exit(0)
		default:
			//如果用户输入无效，提示用户重新输入
			fmt.Println("Invalid input. Please enter 'y' to delete or 'n' to exit")
			return Elasticsearch()
		}
	}
	return esService.IndexCreate(elasticsearch.ArticleIndex(), elasticsearch.ArticleMapping())
}
