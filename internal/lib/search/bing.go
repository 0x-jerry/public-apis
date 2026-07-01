package search

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"public-apis/internal/lib/browser"
	"github.com/go-rod/rod/lib/proto"
)

func Bing(ctx context.Context, bm *browser.Manager, query string) (*Response, error) {
	b, err := bm.Get()
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, fmt.Errorf("browser not available")
	}

	page, err := b.Page(proto.TargetCreateTarget{})
	if err != nil {
		return nil, err
	}
	defer page.Close()

	page = page.Context(ctx)
	searchURL := "https://www.bing.com/search?q=" + url.QueryEscape(query)
	if err := page.Timeout(15 * time.Second).Navigate(searchURL); err != nil {
		return nil, err
	}
	page.WaitLoad()

	if _, err := page.Element("[aria-label=\"Search Results\"]"); err != nil {
		return nil, fmt.Errorf("no search results found")
	}

	result, err := page.Eval(`() => {
		const results = [];
		const items = document.querySelectorAll('li.b_algo');
		for (const item of items) {
			const titleEl = item.querySelector('h2 a');
			const title = titleEl?.textContent?.trim() ?? '';
			const url = titleEl?.href ?? '';
			const snippetEl = item.querySelector('.b_caption p, .b_lineclamp2');
			const snippet = snippetEl?.textContent?.trim() ?? '';
			if (title && url) results.push({ title, url, snippet });
		}

		let prev, next;
		const nav = document.querySelector('nav[role="navigation"]');
		if (nav) {
			const prevLink = nav.querySelector('a.sb_pagP');
			const nextLink = nav.querySelector('a.sb_pagN');
			if (prevLink) prev = prevLink.href;
			if (nextLink) next = nextLink.href;
		}

		return { results, prev, next };
	}`)
	if err != nil {
		return nil, err
	}

	var resp Response
	if err := result.Value.Unmarshal(&resp); err != nil {
		return nil, err
	}

	for i := range resp.Results {
		resp.Results[i].Title = strings.TrimSpace(resp.Results[i].Title)
		resp.Results[i].URL = strings.TrimSpace(resp.Results[i].URL)
		resp.Results[i].Snippet = strings.TrimSpace(resp.Results[i].Snippet)
	}

	return &resp, nil
}
