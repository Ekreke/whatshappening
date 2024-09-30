package test

import (
	"bufio"
	"context"
	"os"
	"testing"

	// "whatshappening/biz"
	"whatshappening/internal/biz"
	"whatshappening/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
)

func TestGetRank(t *testing.T) {
	mockRankConf := &conf.Server{
		RankServer: &conf.Server_RankServer{
			Baseurl: "http://localhost:6688/",
		},
	}

	ru := biz.NewRankUsecase(log.GetLogger(), mockRankConf)
	res, err := ru.GetRankByPlatform(context.Background(), "baidu", false, 10)
	if err != nil {
		t.Error(err)
	}
	t.Log(res)
}

// StopWords 读取停用词表
func loadStopWords(filePath string) (map[string]struct{}, error) {
	stopWords := make(map[string]struct{})
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := scanner.Text()
		stopWords[word] = struct{}{}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return stopWords, nil
}
func TestLoadStopWords(t *testing.T) {
	t.Log(os.Getwd())
	res,_ := loadStopWords("../biz/stop_words.txt")
	t.Log(res,"aaa")
}
func TestWordCountTop50(t *testing.T) {
	// 模拟 RankServer 配置
	mockRankConf := &conf.Server{
		RankServer: &conf.Server_RankServer{
			Baseurl: "http://localhost:6688/",
		},
	}

	// 初始化 RankUsecase
	ru := biz.NewRankUsecase(log.GetLogger(), mockRankConf)
	t.Logf("stopWord: %v", ru.StopWords["”"])
	// 模拟请求数据，限制前50词频
	wcRequest := biz.WordCountRequest{
		Plantforms: []string{"baidu", "weibo"}, // 传入多个平台模拟数据
		IsExclude:  false,                      // 假设不排除某些数据
		Limit:      10,                         // 限制数量
	}

	// 调用 WordCount 方法
	wordCount, err := ru.WordCount(context.Background(), wcRequest)
	if err != nil {
		t.Errorf("Error in WordCount: %v", err)
	}

	// 输出结果，查看是否正常返回了词频前50
	for word, count := range wordCount {
		t.Logf("Word: %s, Count: %d", word, count)
	}

	// // 可以添加进一步的断言，例如检查某个特定的词是否存在
	// expectedWord := "exampleWord"
	// if count, ok := wordCount[expectedWord]; ok {
	// 	t.Logf("Expected word '%s' appears %d times.", expectedWord, count)
	// } else {
	// 	t.Errorf("Expected word '%s' not found in the top 50 words.", expectedWord)
	// }
}
