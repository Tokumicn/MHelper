package mhjl_query

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"golang.org/x/net/html"
	"io"
	"log/slog"
	"mhxyHelper/internal/data"
	"mhxyHelper/pkg/utils"
	"net/http"
	"strings"
)

type MHJLResponse struct {
	Answer    string `json:"answer"`
	Result    string `json:"result"`
	RawAnswer string `json:"raw_answer"`
}

func QueryMHJL(ctx context.Context, query string) (string, error) {

	// 查询梦幻精灵
	rawAnswer, err := query4MHJL(ctx, query)
	if err != nil {
		return "", err
	}

	// 清理无用属性
	cleanHtml, err := parseAndCleanAttributes(ctx, rawAnswer)
	if err != nil {
		return "", err
	}

	// 提取文本
	formatAnswer := extractFormattingText(ctx, cleanHtml)

	// 存储用户问题及精灵回答
	log := data.MHJLResponseLog{
		UserId:       0, // TODO 暂时为添加用户ID 命令行如何添加账户系统待定
		QueryMd5:     utils.MD5(query),
		Query:        query,
		RawAnswer:    rawAnswer,
		RawAnswerMd5: utils.MD5(rawAnswer),
		FormatAnswer: formatAnswer,
	}

	_, err = log.Save(ctx)
	if err != nil {
		return "", err
	}

	return formatAnswer, nil
}

func query4MHJL(ctx context.Context, query string) (string, error) {

	// TODO: 提取配置
	url := fmt.Sprintf("https://xyq.gm.163.com/cgi-bin/csa/csa_sprite.py?act=ask&question=%s&product_name=xyq", query)
	method := http.MethodGet

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return "", err
	}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		err := res.Body.Close()
		if err != nil {
			slog.ErrorContext(ctx, "close res body err: ", err.Error())
		}
	}()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	resp := MHJLResponse{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return "", err
	}

	return resp.RawAnswer, nil
}

// 对原始的 html 进行解析并清理无用属性
func parseAndCleanAttributes(ctx context.Context, htmlStr string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlStr))
	if err != nil {
		return "", err
	}

	// 清理属性
	doc.Find("*").Each(func(i int, s *goquery.Selection) {
		for _, node := range s.Nodes {
			cleanAttributes(node)
		}
	})

	cleanHtml, err := doc.Html()
	if err != nil {
		return "", err
	}

	slog.DebugContext(ctx, "clean html success. cleanHtml: ", cleanHtml)

	// 转换为纯文本
	return cleanHtml, nil
}

func cleanAttributes(n *html.Node) {
	if n.Type == html.ElementNode {
		n.Attr = nil // 清空所有属性
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		cleanAttributes(c)
	}
}

// 保留 html 格式提取文本
func extractFormattingText(ctx context.Context, html string) string {
	browser := rod.New().ControlURL(launcher.New().
		Headless(true).
		Set("default-encoding", "utf-8"). // 关键设置
		MustLaunch()).MustConnect()
	defer browser.MustClose()

	page := browser.MustPage("data:text/html;charset=UTF-8," + html).MustWaitLoad()
	// 获取可视区域文本（自动处理 CSS 样式）
	element := page.MustElement("body")
	text := element.MustText()
	slog.InfoContext(ctx, "build html success. html2text: ", text)
	return text
}
