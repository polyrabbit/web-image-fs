package webpage

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
)

const ua = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.83 Safari/537.36"

type HTTPClient struct {
	urlHost *url.URL
	client  *http.Client
}

func MustNewHTTPClient(host string, timeoutSeconds time.Duration) *HTTPClient {
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

func (c *HTTPClient) Parse(ctx context.Context, rPath string) ([]DomNode, error) {
	resp, err := c.Fetch(ctx, http.MethodGet, rPath)
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
			var (
				contentLength int
				contentType   string
			)
			if header, err := c.MetaInfo(ctx, imgPath); err == nil {
				contentLength, _ = strconv.Atoi(header.Get("Content-Length"))
				contentType = header.Get("Content-Type")
			} else {
				logrus.WithError(err).WithField("path", imgPath).Debug("Failed to get image header")
			}
			doms = append(doms, &ImageNode{
				// TODO: fix duplicated names
				LinkNode: LinkNode{
					Name: s.AttrOr("alt", ""),
					// FIXME: parent path is missing
					SelfLink: imgPath,
				},
				Size:        uint64(contentLength),
				ContentType: contentType,
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

// GetInfo uses http HEAD method to get meta info in advance
func (c *HTTPClient) MetaInfo(ctx context.Context, rPath string) (http.Header, error) {
	resp, err := c.Fetch(ctx, http.MethodHead, rPath)
	if err != nil {
		return nil, fmt.Errorf("http HEAD: %w", err)
	}
	defer resp.Body.Close()
	return resp.Header, nil
}

func (c *HTTPClient) Download(ctx context.Context, rPath string) ([]byte, error) {
	resp, err := c.Fetch(ctx, http.MethodGet, rPath)
	if err != nil {
		return nil, fmt.Errorf("http GET: %w", err)
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func (c *HTTPClient) Fetch(ctx context.Context, method, rPath string) (*http.Response, error) {
	uRef, err := url.Parse(rPath)
	if err != nil {
		return nil, fmt.Errorf("parse ref url: %w", err)
	}
	absURL := c.urlHost.ResolveReference(uRef)
	req, err := http.NewRequestWithContext(ctx, method, absURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequestWithContext: %w", err)
	}
	req.Header.Set("User-Agent", ua)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http client Do: %w", err)
	}
	if resp.StatusCode != 200 {
		return nil, NewHTTPError(resp)
	}
	return resp, nil
}
