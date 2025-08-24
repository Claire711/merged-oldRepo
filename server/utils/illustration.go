package utils

import "regexp"

func FindIllustrations(text string) ([]string, error) {
	// 定义正则表达式，匹配 Markdown 图片语法
	// 格式: ![alt text](image_url)
	regex := `!\[([^\]]*)\]\(([^)]+)\)`
	re, err := regexp.Compile(regex)
	if err != nil {
		return nil, err
	}
	matches := re.FindAllStringSubmatch(text, -1)
	var illustrations []string
	for _, match := range matches {
		if len(match) > 2 {
			illustrations = append(illustrations, match[2])
		}
	}
	return illustrations, nil
}

//从文本中批量提取 Markdown 格式的图片链接
