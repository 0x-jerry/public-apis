package fetch

import (
	"github.com/enetx/g"
	"github.com/enetx/surf"
)

type FetchResult struct {
	Body        string
	ContentType string
	StatusCode  int
}

func Fetch(rawURL string, accept string) (*FetchResult, error) {
	client := surf.NewClient()
	defer client.Close()

	req := client.Get(g.NewString(rawURL))
	if accept != "" {
		req.SetHeaders(g.NewString("Accept"), g.NewString(accept))
	}

	result := req.Do()
	if result.IsErr() {
		return nil, result.Err()
	}
	resp := result.Ok()
	defer resp.Body.Close()

	contentType := string(resp.Headers.Get(g.NewString("Content-Type")))
	statusCode := int(resp.StatusCode)

	bodyResult := resp.Body.String()
	if bodyResult.IsErr() {
		return nil, bodyResult.Err()
	}

	return &FetchResult{
		Body:        string(bodyResult.Ok()),
		ContentType: contentType,
		StatusCode:  statusCode,
	}, nil
}
