package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"regexp"

	"github.com/go-ego/gse"
)

// 读取数据源，构造索引数据
func fileOpen() []string {
	file, err := os.Open("data_source/movies.csv")
	if err != nil {
		fmt.Println("err,", err)
	}

	defer file.Close()

	// 创建一个 Scanner 用来读取文件内容
	docx := make([]string, 0)
	scanner := bufio.NewScanner(file)

	// 逐行读取文件内容并打印
	for scanner.Scan() {
		re := make([]string, 0)
		line := scanner.Text()
		re = strings.Split(line, ",")
		docx = append(docx, re[4])
	}

	docx = docx[1:]

	return docx
}

var StopWord = []string{",", ".", "。", "*", "(", ")", "'", "\""}

// 数据清洗
func removeStopWord(word string) string {
	for i := range StopWord {
		word = strings.Replace(word, StopWord[i], "", -1)
	}

	return word
}

var enStopWords = map[string]struct{}{
	"a": {}, "is": {}, "the": {}, "of": {}, "and": {},
}

func filterStopWords(tokens []string) []string {
	res := make([]string, 0, len(tokens))
	for _, t := range tokens {
		if _, ok := enStopWords[t]; !ok {
			res = append(res, t)
		}
	}
	return res
}

var gobalGse gse.Segmenter

func InitConfig() {
	var err error
	gobalGse, err = gse.New()
	if err != nil {
		fmt.Println("init config err,", err)
	}
}

var englishWordRe = regexp.MustCompile(`[a-zA-Z]+`)

// token化
func tokenize(text string) []string {
	if isASCII(text) {
		tokens := englishWordRe.FindAllString(text, -1)
		return filterStopWords(tokens)
	}

	text = removeStopWord(text)
	return gobalGse.CutSearch(text)
}

func isASCII(s string) bool {
	for _, r := range s {
		if r > 127 {
			return false
		}
	}
	return true
}
