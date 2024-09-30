package biz

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"whatshappening/internal/conf"

	"os"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/yanyiwu/gojieba"
)

// Platforms 定义所有平台的常量
var platforms = []string{
	"bilibili",
	"acfun",
	"weibo",
	"zhihu",
	"zhihu-daily",
	"baidu",
	"douyin",
	"douban-movie",
	"douban-group",
	"tieba",
	"sspai",
	"ithome",
	"ithome-xijiayi",
	"jianshu",
	"thepaper",
	"toutiao",
	"36kr",
	"51cto",
	"csdn",
	"nodeseek",
	"juejin",
	"qq-news",
	"sina",
	"sina-news",
	"netease-news",
	"52pojie",
	"hostloc",
	"huxiu",
	"hupu",
	"ifanr",
	"lol",
	"genshin",
	"honkai",
	"starrail",
	"weread",
	"ngabbs",
	"v2ex",
	"hellogithub",
	"weatheralarm",
	"earthquake",
	"history",
}

type RankRepo interface {
}

// RankUsecase 定义排行榜业务逻辑
type RankUsecase struct {
	// repo RankRepo
	RankConf  *conf.Server_RankServer
	log       *log.Helper
	jieba     *gojieba.Jieba
	StopWords map[string]struct{}
}

func NewRankUsecase(logger log.Logger, s *conf.Server) *RankUsecase {
	stopWords, _ := loadStopWords("./stop_words.txt")

	return &RankUsecase{
		log:       log.NewHelper(logger),
		RankConf:  s.RankServer,
		jieba:     gojieba.NewJieba(),
		StopWords: stopWords,
	}
}

// RankResp 定义响应结构
type RankResp struct {
	Code        int         `json:"code"`
	Name        string      `json:"name"`
	Title       string      `json:"title"`
	Type        string      `json:"type"`
	Description string      `json:"description"` // 新增描述字段
	Params      interface{} `json:"params"`
	Link        string      `json:"link"`
	Total       int         `json:"total"`
	UpdateTime  string      `json:"updateTime"` // 修改为string类型
	FromCache   bool        `json:"fromCache"`
	Data        []RankEntry `json:"data"`
}

// Params 定义 params 字段的结构
type Params struct {
	Type ParamType `json:"type"`
}

// ParamType 定义 type 字段的详细内容
type ParamType struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// RankEntry 定义每个热搜榜条目的结构
type RankEntry struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}
type WordCountRequest struct {
	Plantforms []string `json:"plantforms"`
	IsExclude  bool     `json:"isExclude"`
	Limit      int      `json:"limit"`
}

// StopWords 读取停用词表
func loadStopWords(filePath string) (map[string]struct{}, error) {
	pwd, _ := os.Getwd()
	fmt.Println("pwd: ", pwd)
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

// 分词并统计词频
func (ru *RankUsecase) wordCount(rs []*RankResp) map[string]int32 {
	wordCount := make(map[string]int32)

	for _, r := range rs {
		for _, entry := range r.Data {
			words := ru.jieba.Cut(entry.Title, true) // 使用gojieba进行分词
			for _, word := range words {
				// 忽略停用词和单个字符的词
				if _, ok := ru.StopWords[word]; !ok && len(word) > 1 {
					wordCount[word]++
				}
			}
		}
	}
	return wordCount
}

// 计算并返回前50个词
func (ru *RankUsecase) WordCount(ctx context.Context, wc WordCountRequest) (map[string]int32, error) {
	if wc.Plantforms == nil {
		return nil, errors.New(400, "INVALID_PLANTFORM", "Invalid plantform")
	}
	cacheEnable := ru.RankConf.CacheEnable
	var (
		rs  []*RankResp
		err error
	)

	if wc.IsExclude {
		rs, err = ru.GetRankListExclude(ctx, wc.Plantforms, cacheEnable, wc.Limit)
	} else {
		rs, err = ru.GetRankList(ctx, wc.Plantforms, cacheEnable, wc.Limit)
	}
	if err != nil {
		return nil, err
	}

	// 统计词频
	wordCount := ru.wordCount(rs)

	// 返回词频前50的词
	return topNWords(wordCount, 50), nil
}

// 获取词频前N个词
func topNWords(wordCount map[string]int32, N int) map[string]int32 {
	type wordFreq struct {
		Word  string
		Count int32
	}
	words := make([]wordFreq, 0, len(wordCount))

	for word, count := range wordCount {
		words = append(words, wordFreq{Word: word, Count: count})
	}

	sort.Slice(words, func(i, j int) bool {
		return words[i].Count > words[j].Count
	})

	topWords := make(map[string]int32)
	for i := 0; i < N && i < len(words); i++ {
		topWords[words[i].Word] = words[i].Count
	}

	return topWords
}

func (ru *RankUsecase) GetRankByPlatform(ctx context.Context, platform string, cache bool, limit int) (*RankResp, error) {
	//后面优化要修改这个字符串是否符合上面plantforms的值
	if platform == "" {
		return nil, errors.New(400, "INVALID_PLATFORM", "Invalid platform")
	}
	url := fmt.Sprintf("%s%s?cache=%t", ru.RankConf.Baseurl, platform, cache)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.New(500, "REQUEST_CREATION_FAILED", "Failed to create request")
	}

	// 发起 HTTP 请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New(500, "REQUEST_EXECUTION_FAILED", "Failed to execute request")
	}
	defer resp.Body.Close()

	// 检查 HTTP 响应状态码
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, errors.New(500, "REQUEST_FAILED", fmt.Sprintf("Request failed with status %d: %s", resp.StatusCode, body))
	}

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New(500, "READ_BODY_FAILED", "Failed to read response body")
	}
	// fmt.Println("Response Body:", string(body))

	// 解析 JSON 响应数据
	var result RankResp
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, errors.New(500, "JSON_UNMARSHAL_FAILED", "Failed to unmarshal JSON data")
	}

	return &result, nil
}

// GetRankList 获取排行榜数据
func (ru *RankUsecase) GetRankList(ctx context.Context, rankInfo []string, cache bool, limit int) ([]*RankResp, error) {
	fmt.Println("GetRankList")

	rs := make([]*RankResp, 0)
	for _, planform := range rankInfo {
		result, err := ru.GetRankByPlatform(ctx, planform, cache, limit)
		if err != nil {
			continue
		}
		//添加result.Data数据到哪里？
		rs = append(rs, result)
	}
	return rs, nil
}
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
func (ru *RankUsecase) GetRankListExclude(ctx context.Context, rankInfo []string, cache bool, limit int) ([]*RankResp, error) {
	fmt.Println("GetRankList")

	rs := make([]*RankResp, 0)
	for _, planform := range platforms {
		if !contains(rankInfo, planform) {
			result, err := ru.GetRankByPlatform(ctx, planform, cache, limit)
			if err != nil {
				continue
			}
			//添加result.Data数据到哪里？
			rs = append(rs, result)
		}
	}
	return rs, nil
}

// func (ru *RankUsecase) WordCount(ctx context.Context, wc WordCountRequest) (map[string]int32, error) {
// 	// 检查输入是否有效
// 	if wc.Plantforms == nil {
// 		return nil, errors.New(400, "INVALID_PLANTFORM", "Invalid plantform")
// 	}
// 	cacheEnable := ru.RankConf.CacheEnable
// 	// 声明用于存储结果的切片和错误变量
// 	var (
// 		rs  []*RankResp
// 		err error
// 	)

// 	// 根据 wc.IsExclude 调用不同的方法
// 	if wc.IsExclude {
// 		rs, err = ru.GetRankListExclude(ctx, wc.Plantforms, cacheEnable, wc.Limit)
// 	} else {
// 		rs, err = ru.GetRankList(ctx, wc.Plantforms, cacheEnable, wc.Limit)
// 	}

// 	// 如果调用出错，返回错误
// 	if err != nil {
// 		return nil, err
// 	}
// 	wordCount := wordCount(rs)
// 	return wordCount, nil
// }
// func wordCount(rs []*RankResp) map[string]int32 {
// 	wordCount := make(map[string]int32)
// 	for _, r := range rs {
// 		for _, d := range r.Data {
// 			wordCount[d.Title] = wordCount[d.Title] + 1
// 		}
// 	}
// 	return wordCount
// }

//当前路径./stop_words.txt有一张停用词表

// func main() {
// 	ctx := context.Background()
// 	rankInfo := []string{} // 示例：你可以在这里传递一些参数
// 	result, err := GetRankList(ctx, rankInfo, false, 10)
// 	if err != nil {
// 		fmt.Printf("Error occurred: %v\n", err)
// 		return
// 	}

// 	fmt.Println("Received Rank Response:")
// 	fmt.Println(result)
// }
