package webpage

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const ua = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.83 Safari/537.36"

type HTTPClient struct {
	urlHost *url.URL
	client  *http.Client
}

func MustWebBrowser(host string, timeoutSeconds time.Duration) *HTTPClient {
	u, err := url.Parse(host)
	if err != nil {
		panic(err)
	}
	if u.Scheme == "" {
		u.Scheme = "http"
	}
	return &HTTPClient{
		urlHost: u,
		client: &http.Client{
			Timeout: timeoutSeconds,
		},
	}
}

func (b *HTTPClient) Parse(ctx context.Context, rPath string) ([]DomNode, error) {
	resp, err := b.Fetch(ctx, http.MethodGet, rPath)
	if err != nil {
		return nil, fmt.Errorf("http GET: %w", err)
	}
	defer resp.Body.Close()
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("goquery.NewDocumentFromReader: %w", err)
	}
	var doms []DomNode
	// First, find all images
	doc.Find("img[src]").Each(func(i int, s *goquery.Selection) {
		if imgPath, exists := s.Attr("src"); exists {
			doms = append(doms, &ImageNode{
				// TODO: fix duplicated names
				LinkNode: LinkNode{
					Name:     s.AttrOr("alt", ""),
					SelfLink: imgPath,
				},
			})
		}
	})

	// Then, find all links/directories
	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		if linkPath, exists := s.Attr("href"); exists {
			doms = append(doms, &LinkNode{
				// TODO: fix duplicated names
				Name:     s.Text(),
				SelfLink: linkPath,
			})
		}
	})
	return doms, nil
}

func (b *HTTPClient) Fetch(ctx context.Context, method, rPath string) (*http.Response, error) {
	urlObj := *b.urlHost // A shadow copy
	urlObj.Path = path.Join(urlObj.Path, rPath)
	req, err := http.NewRequestWithContext(ctx, method, urlObj.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequestWithContext: %w", err)
	}
	req.Header.Set("User-Agent", ua)
	resp, err := b.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http client Do: %w", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected http status: %d(%s)", resp.StatusCode, resp.Status)
	}
	return resp, nil
}
