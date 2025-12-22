package main

type InvertedIndex map[string][]int

// 构建倒排索引
func BuildIndex(docx []string) InvertedIndex {
	index := make(InvertedIndex)

	// 遍历所有的 docx，i 是行号，d 是文档内容
	for i, d := range docx {
		// 文本切词
		for _, word := range tokenize(d) {
			if _, ok := index[word]; !ok {
				// 如果词第一次出现，创建一个新列表并存入当前行号
				index[word] = []int{i}
			} else {
				// 如果已存在，直接把行号追加到列表末尾
				index[word] = append(index[word], i)
			}
		}
	}

	return index
}
