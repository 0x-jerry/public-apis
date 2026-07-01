package search

import (
	"fmt"
	"strings"
)

type Result struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Snippet string `json:"snippet"`
}

type Response struct {
	Results []Result `json:"results"`
	Prev    string   `json:"prev,omitempty"`
	Next    string   `json:"next,omitempty"`
}

func ToMarkdown(resp *Response) string {
	if resp == nil {
		return "No results found."
	}

	var lines []string

	if len(resp.Results) == 0 {
		lines = append(lines, "No results found.")
	} else {
		for _, r := range resp.Results {
			lines = append(lines, "- ["+r.Title+"]("+r.URL+")")
			if r.Snippet != "" {
				lines = append(lines, "  "+r.Snippet)
			}
		}
	}

	var nav []string
	if resp.Prev != "" {
		nav = append(nav, "[Previous]("+resp.Prev+")")
	}
	if resp.Next != "" {
		nav = append(nav, "[Next]("+resp.Next+")")
	}
	if len(nav) > 0 {
		if len(nav) == 2 {
			lines = append(lines, "")
			lines = append(lines, fmt.Sprintf("%s | %s", nav[0], nav[1]))
		} else if len(nav) == 1 {
			lines = append(lines, "")
			lines = append(lines, nav[0])
		}
	}

	return strings.Join(lines, "\n")
}
