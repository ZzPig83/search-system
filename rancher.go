package main

import (
	"sort"
	"strings"
)

type SortRes struct {
	Docx  string
	Score float64
	Id    int
}

func sortRes(qy []string, res []string) []*SortRes {
	exist := make(map[int]*SortRes)
	// 遍历每一个 query 的分词后的 token
	for _, v := range qy {
		// 遍历每一个召回的文档
		for i, v2 := range res {
			score := calculateTFIDF(v, v2, res)

			// 记录分数构成，计算每个词条对每个文档结构的score
			if _, ok := exist[i]; !ok {
				// 如果exist中还没存在这个词条，则进行进行初始化
				tmp := &SortRes{
					Docx:  v2,
					Score: score,
					Id:    i,
				}
				exist[i] = tmp
			} else {
				// 如果已经存在了，则进行分数的相加
				// 意思就是每个res中的doc对于每个token的权重之和的结果。权重的对象始终都是res中doc
				exist[i].Score += score
			}
		}
	}
	resList := make([]*SortRes, 0)
	for _, v := range exist { // 构建结构体
		resList = append(resList, &SortRes{
			Docx:  v.Docx,
			Score: v.Score,
			Id:    v.Id,
		})
	}
	sort.Slice(resList, func(i, j int) bool { // 按照score进行排序
		return resList[i].Score > resList[j].Score
	})
	return resList
}

// 计算 TFDF
func calculateTFIDF(term string, document string, documents []string) float64 {
	tf := calculateTF(term, document)
	idf := calculateIDF(term, documents)
	return tf * idf * 100.0
}

// 计算 TF
func calculateTF(term string, document string) float64 {
	termCount := strings.Count(document, term)
	totalWords := len(tokenize(document))
	return float64(termCount) / float64(totalWords)
}

// 计算 IDF
func calculateIDF(term string, documents []string) float64 {
	docWithTerm := 0
	for _, doc := range documents {
		if strings.Contains(doc, term) {
			docWithTerm++
		}
	}
	return float64(len(documents)) / float64(docWithTerm)
}
