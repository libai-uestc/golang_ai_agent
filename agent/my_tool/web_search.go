package mytool

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/karust/openserp/core"
	"github.com/karust/openserp/duckduckgo" // duckduckgo 官方并没有维护任何search API
)

var (
	browser    *core.Browser
	browseOnce sync.Once
)

func GetBrowser() *core.Browser {
	browseOnce.Do(func() {
		var err error
		browser, err = core.NewBrowser(core.BrowserOpts{})
		if err != nil {
			panic(err)
		}
	})
	return browser
}

func CloseBrowser() {
	if browser != nil {
		browser.Close()
	}
}

type SearchRequest struct {
	Query string `json:"query" jsonschema:"required,description=搜索词"`
	Count int    `json:"count" jsonschema:"description=最多返回搜索结果的条数，默认值为3"`
}

func DuckDuckSearch(ctx context.Context, request *SearchRequest) (string, error) {
	if request.Query == "" {
		return "", errors.New("搜索query为空")
	}
	if request.Count == 0 {
		request.Count = 3
	}
	if request.Count > 3 {
		request.Count = 3
	}
	query := core.Query{
		Text:     request.Query,
		LangCode: "zh-cn",
		Limit:    request.Count,
		Site:     "",
		/*
			# SOCKS5 proxy
			./openserp serve --proxy socks5://127.0.0.1:1080
			# HTTP proxy with authentication
			./openserp search bing "query" --proxy http://user:pass@127.0.0.1:8080
		*/
		ProxyURL: "",
		Insecure: false,
	}

	ddg := duckduckgo.New(*GetBrowser(), core.SearchEngineOptions{})
	results, err := ddg.Search(ctx, query)
	if err != nil {
		return "", err
	}
	s, err := sonic.MarshalString(results)
	if err != nil {
		return "", err
	}

	return s, nil

}

func CreateDuckDuckSearchTool() tool.InvokableTool {
	tool, err := utils.InferTool("web_search_tool", "调用搜索引擎，获取搜索结果", DuckDuckSearch)
	if err != nil {
		log.Fatal(err)
	}
	return tool
}

func UseDuckDuckSearch(query string) {
	ctx := context.Background()
	timeTool := CreateDuckDuckSearchTool()

	request := SearchRequest{
		Query: query,
		Count: 2,
	}
	jsonReq, err := sonic.MarshalString(request)
	if err != nil {
		log.Fatalf("Marshal of search request failed: %s", err)
	}

	resp, err := timeTool.InvokableRun(ctx, jsonReq)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}

type Url struct {
	Url string `json:"url" jsonschema:"required,description=被爬取的网页的url"`
}

// func Crawl(ctx context.Context, url *Url) (string, error) {
// 	if url.Url == "" {
// 		return "", errors.New("没有指定爬取的URL")
// 	}

//		page, err := GetBrowser().Navigate(ctx, url.Url)
//		if err != nil {
//			return "", err
//		}
//		element, _ := page.Element("title")
//		title, _ := element.Text()
//		element, _ = page.Element("body")
//		body, _ := element.Text()
//		result := "网页标题: " + title + "\n网页内容：" + body
//		return result, nil
//	}
func Crawl(ctx context.Context, url *Url) (string, error) {
	if url.Url == "" {
		return "工具调用失败: 没有指定爬取的URL", nil // 返回普通字符串提示LLM，且 error 传 nil
	}

	page, err := GetBrowser().Navigate(ctx, url.Url)
	if err != nil {
		// 核心修改：将错误信息作为字符串反馈给模型，让模型知道被拦截了，并让其换个方法
		errMsg := fmt.Sprintf("爬取该网页失败，可能遭遇了反爬拦截 (错误: %v)。请使用现有的搜索结果，或尝试爬取其他搜索结果中的链接。", err)
		return errMsg, nil
	}

	element, err := page.Element("title")
	title := "未知标题"
	if err == nil {
		title, _ = element.Text()
	}

	element, err = page.Element("body")
	body := ""
	if err == nil {
		body, _ = element.Text()
	} else {
		return "工具调用失败: 网页加载成功但无法获取页面正文", nil
	}

	// 【强烈建议添加】：限制网页正文长度，防止网页内容过长导致大模型的 Context Token 爆满
	maxLength := 4000
	runeBody := []rune(body)
	if len(runeBody) > maxLength {
		body = string(runeBody[:maxLength]) + "\n...(内容过长已截断)"
	}

	result := "网页标题: " + title + "\n网页内容：" + body
	return result, nil
}

func CreateCrawlTool() tool.InvokableTool {
	tool, err := utils.InferTool("web_crawler_tool", "网页爬虫", Crawl)
	if err != nil {
		log.Fatal(err)
	}
	return tool
}

func UseCrawlTool(url string) {
	ctx := context.Background()
	timeTool := CreateCrawlTool()

	request := Url{
		Url: url,
	}
	jsonReq, err := sonic.MarshalString(request)
	if err != nil {
		log.Fatalf("Marshal of crawl request failed: %s", err)
	}

	resp, err := timeTool.InvokableRun(ctx, jsonReq)
	if err != nil {
		log.Printf("crawl failed: %s", err)
	}

	fmt.Println(resp)
}
