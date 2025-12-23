package service

import (
	"bufio"
	"fmt"
	"github.com/go-ego/gse"
	"os"
	"regexp"
	"sort"
	"strings"
)

func Search(query string) []*SortRes {
	query = "城市刑警普通案件"

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
	return resList
}

func GetDataSourceLenth() int {
	return 888
}

// 读取数据源，构造索引数据
func fileOpen() []string {
	file, err := os.Open("./data_source/movies.csv")
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
