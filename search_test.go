package main

import (
	"fmt"
	"testing"
)

func TestSearch(t *testing.T) {
	query := "城市刑警普通案件"

	// 初始化配置
	InitConfig()
	docx := fileOpen()

	// 构建索引
	index := BuildIndex(docx)

	for w, i := range index {
		fmt.Printf("w:%v i:%v", w, i)
	}

	// 召回
	res, qy := retrieval(index, query, docx)
	fmt.Printf("一共%d记录，query分词结果%v\n", len(res), qy)

	// 排序
	resList := sortRes(qy, res)
	for i := range resList {
		fmt.Println(resList[i].Score, resList[i].Docx)
	}
}

func TestTokenize(t *testing.T) {
	query := "must"

	qy := tokenize(query)
	fmt.Printf("分词结果:%v", qy)
}
