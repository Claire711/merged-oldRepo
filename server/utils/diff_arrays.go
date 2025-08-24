package utils

func DiffArrays(oldArray, newArray []string) (added, removed []string) {
	//added（新数组比旧数组多的元素）和removed（旧数组比新数组多的元素）
	oldMap := make(map[string]struct{})
	//将数组的 "线性查找" 转化为 map 的 "常数时间查找"
	for _, item := range oldArray {
		oldMap[item] = struct{}{}
	}
	for _, item := range newArray {
		if _, exists := oldMap[item]; !exists {
			added = append(added, item)
		}
	}
	newMap := make(map[string]struct{})
	for _, item := range newArray {
		newMap[item] = struct{}{}
	}
	for _, item := range oldArray {
		if _, exists := newMap[item]; !exists {
			removed = append(removed, item)
		}
	}
	return
}
