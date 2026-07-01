package search

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/enetx/g"
	"github.com/enetx/surf"
)

func DuckDuckGo(query string) (*Response, error) {
	searchURL := "https://html.duckduckgo.com/html/?q=" + url.QueryEscape(query)

	client := surf.NewClient()
	result := client.Get(g.NewString(searchURL)).Do()
	if result.IsErr() {
		return nil, fmt.Errorf("fetch duckduckgo: %w", result.Err())
	}
	resp := result.Ok()
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body.Reader)
	if err != nil {
		return nil, fmt.Errorf("parse duckduckgo: %w", err)
	}

	var results []Result
	doc.Find(".result").Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find(".result__a").Text())
		resultURL := strings.TrimSpace(s.Find(".result__url").Text())
		snippet := strings.TrimSpace(s.Find(".result__snippet").Text())
		if title != "" && resultURL != "" {
			results = append(results, Result{
				Title:   title,
				URL:     resultURL,
				Snippet: snippet,
			})
		}
	})

	var prev, next string

	doc.Find(".nav-link a").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		href, exists := s.Attr("href")
		if exists {
			if strings.Contains(strings.ToLower(text), "prev") {
				prev = href
			}
			if strings.Contains(strings.ToLower(text), "next") {
				next = href
			}
		}
	})

	doc.Find(".nav-link form").Each(func(i int, s *goquery.Selection) {
		value := strings.TrimSpace(s.Find("input[type=\"submit\"]").AttrOr("value", ""))
		action := s.AttrOr("action", "/html/")

		params := url.Values{}
		s.Find("input[type=\"hidden\"]").Each(func(j int, input *goquery.Selection) {
			name := input.AttrOr("name", "")
			val := input.AttrOr("value", "")
			if name != "" {
				params.Set(name, val)
			}
		})

		formURL := fmt.Sprintf("https://html.duckduckgo.com%s?%s", action, params.Encode())
		if strings.Contains(strings.ToLower(value), "prev") {
			prev = formURL
		}
		if strings.Contains(strings.ToLower(value), "next") {
			next = formURL
		}
	})

	return &Response{
		Results: results,
		Prev:    prev,
		Next:    next,
	}, nil
}
