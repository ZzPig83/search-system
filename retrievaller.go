package main

func retrieval(index InvertedIndex, query string, docs []string) ([]string, []string) {
	// 对 query 分词
	qy := tokenize(query)

	result := make(map[int]bool)
	for _, word := range qy {
		if doc, ok := index[word]; ok {
			// 搜索倒排索引中，term对应的doc数组，doc数组就是存在该term词条的所有的doc id
			for _, d := range doc {
				// 对doc数组进行遍历，获取所有的doc id，并且进行标志。
				result[d] = true
			}
		}
	}

	output := []string{}
	for d := range result {
		output = append(output, docs[d])
		// 利用正排索引，找到id对应的存储内容并返回
	}
	return output, qy
}
