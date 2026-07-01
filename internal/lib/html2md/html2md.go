package html2md

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"

	readability "codeberg.org/readeck/go-readability/v2"
	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/base"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/commonmark"
	"github.com/goccy/go-yaml"
	"golang.org/x/net/html"
)

var conv = converter.NewConverter(
	converter.WithPlugins(
		base.NewBasePlugin(),
		commonmark.NewCommonmarkPlugin(),
	),
)

func resolveRelativeLinks(doc *html.Node, rootUrl string) {
	baseURL, err := url.Parse(rootUrl)
	if err != nil {
		return
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for i, a := range n.Attr {
				if a.Key == "href" {
					if rel := relativize(baseURL, a.Val); rel != "" {
						n.Attr[i].Val = rel
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
}

func relativize(baseURL *url.URL, href string) string {
	if href == "" {
		return ""
	}
	lower := strings.ToLower(href)
	if strings.HasPrefix(lower, "javascript:") || strings.HasPrefix(lower, "mailto:") || strings.HasPrefix(lower, "tel:") {
		return ""
	}
	u, err := url.Parse(href)
	if err != nil || u.Host == "" {
		return ""
	}
	if u.Host == baseURL.Host {
		return u.RequestURI()
	}
	return ""
}

type Options struct {
	Mode string
	URL  string
}

type Metadata struct {
	Title       string `yaml:"title,omitempty"`
	Description string `yaml:"description,omitempty"`
	URL         string `yaml:"url,omitempty"`
}

func Convert(htmlStr string, opts *Options) (string, error) {
	if opts == nil {
		opts = &Options{Mode: "readable"}
	}

	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return "", fmt.Errorf("parse html: %w", err)
	}

	metadata := extractMetadata(doc)
	metadata.URL = opts.URL
	removeInvisibleTags(doc)

	if opts.URL != "" {
		resolveRelativeLinks(doc, opts.URL)
	}

	var content string

	if opts.Mode == "readable" {
		removeEmptyImageTags(doc)
		removeEmptyLinkTags(doc)
		removeEmbedTags(doc)

		if readability.CheckDocument(doc) {
			article, err := readability.FromDocument(doc, nil)
			if err == nil {
				var buf bytes.Buffer
				if err := article.RenderHTML(&buf); err == nil {
					content = buf.String()
				}
			}
		}
	}

	if content == "" {
		content = renderHTML(doc)
	}

	markdown, err := conv.ConvertString(content)
	if err != nil {
		return "", fmt.Errorf("convert to markdown: %w", err)
	}

	metaYAML, _ := yaml.Marshal(metadata)
	metaStr := strings.TrimRight(string(metaYAML), "\n")

	return "---\n" + metaStr + "\n---\n" + markdown, nil
}

func extractMetadata(doc *html.Node) Metadata {
	var m Metadata
	m.Title = extractText(doc, "title")
	if m.Title == "" {
		m.Title = extractMetaContent(doc, "property", "og:title")
	}
	m.Description = extractMetaContent(doc, "name", "description")
	if m.Description == "" {
		m.Description = extractMetaContent(doc, "property", "og:description")
	}
	return m
}

func extractText(doc *html.Node, tag string) string {
	var text string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == tag {
			if n.FirstChild != nil {
				text = strings.TrimSpace(n.FirstChild.Data)
			}
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return text
}

func extractMetaContent(doc *html.Node, attr, value string) string {
	var content string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "meta" {
			found := false
			for _, a := range n.Attr {
				if a.Key == attr && a.Val == value {
					found = true
				}
			}
			if found {
				for _, a := range n.Attr {
					if a.Key == "content" {
						content = a.Val
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return content
}

func removeInvisibleTags(doc *html.Node) {
	tags := []string{"script", "head", "style", "noscript", "template"}
	removeTags(doc, tags)
}

func removeEmptyImageTags(doc *html.Node) {
	removeIf(doc, "img", func(n *html.Node) bool {
		return getAttr(n, "alt") == ""
	})
}

func removeEmptyLinkTags(doc *html.Node) {
	removeIf(doc, "a", func(n *html.Node) bool {
		return strings.TrimSpace(textContent(n)) == ""
	})
}

func removeEmbedTags(doc *html.Node) {
	tags := []string{"object", "embed", "iframe"}
	removeTags(doc, tags)
}

func removeTags(doc *html.Node, tags []string) {
	tagSet := make(map[string]bool)
	for _, t := range tags {
		tagSet[t] = true
	}

	var removeList []*html.Node
	var find func(*html.Node)
	find = func(n *html.Node) {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode && tagSet[c.Data] {
				removeList = append(removeList, c)
			}
			find(c)
		}
	}
	find(doc)

	for _, n := range removeList {
		if n.Parent != nil {
			n.Parent.RemoveChild(n)
		}
	}
}

func removeIf(doc *html.Node, tag string, cond func(*html.Node) bool) {
	var removeList []*html.Node
	var find func(*html.Node)
	find = func(n *html.Node) {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode && c.Data == tag && cond(c) {
				removeList = append(removeList, c)
			}
			find(c)
		}
	}
	find(doc)

	for _, n := range removeList {
		if n.Parent != nil {
			n.Parent.RemoveChild(n)
		}
	}
}

func getAttr(n *html.Node, key string) string {
	for _, a := range n.Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}

func textContent(n *html.Node) string {
	var buf strings.Builder
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			buf.WriteString(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(n)
	return buf.String()
}

func renderHTML(doc *html.Node) string {
	var buf bytes.Buffer
	html.Render(&buf, doc)
	return buf.String()
}
