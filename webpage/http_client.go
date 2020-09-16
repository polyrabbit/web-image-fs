package webpage

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
)

const ua = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.83 Safari/537.36"

type HTTPClient struct {
	client *http.Client
}

func NewHTTPClient(timeoutSeconds time.Duration) *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: timeoutSeconds,
		},
	}
}

func (c *HTTPClient) Parse(ctx context.Context, absURL string) ([]DomNode, error) {
	resp, err := c.Fetch(ctx, http.MethodGet, absURL)
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
			if strings.HasPrefix(imgPath, "data:") { // TODO: should support base64-encoded image
				return
			}
			imgURL, err := c.URLJoin(absURL, imgPath)
			if err != nil {
				logrus.WithError(err).WithField("path", imgPath).Debug("Failed to join image url")
				return
			}
			var (
				contentLength int
				contentType   string
			)
			if header, err := c.MetaInfo(ctx, imgURL); err == nil {
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
					SelfLink: imgURL,
				},
				Size:        uint64(contentLength),
				ContentType: contentType,
			})
		}
	})

	// Then, find all links/directories
	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		if linkPath, exists := s.Attr("href"); exists {
			if strings.HasPrefix(linkPath, "#") || strings.HasPrefix(linkPath, "javascript:") {
				return
			}
			linkURL, err := c.URLJoin(absURL, linkPath)
			if err != nil {
				logrus.WithError(err).WithField("path", linkPath).Debug("Failed to join link url")
				return
			}
			doms = append(doms, &LinkNode{
				// TODO: fix duplicated names
				Name:     s.Text(),
				SelfLink: linkURL,
			})
		}
	})
	return doms, nil
}

// GetInfo uses http HEAD method to get meta info in advance
func (c *HTTPClient) MetaInfo(ctx context.Context, absURL string) (http.Header, error) {
	resp, err := c.Fetch(ctx, http.MethodHead, absURL)
	if err != nil {
		return nil, fmt.Errorf("http HEAD: %w", err)
	}
	defer resp.Body.Close()
	return resp.Header, nil
}

func (c *HTTPClient) Download(ctx context.Context, absURL string) ([]byte, error) {
	resp, err := c.Fetch(ctx, http.MethodGet, absURL)
	if err != nil {
		return nil, fmt.Errorf("http GET: %w", err)
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func (c *HTTPClient) Fetch(ctx context.Context, method, absURL string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, absURL, nil)
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

func (c *HTTPClient) URLJoin(base, relativePath string) (string, error) {
	uBase, err := url.Parse(base)
	if err != nil {
		return "", fmt.Errorf("parse base url: %w", err)
	}
	uRef, err := url.Parse(relativePath)
	if err != nil {
		return "", fmt.Errorf("parse relative path: %w", err)
	}
	return uBase.ResolveReference(uRef).String(), nil
}
