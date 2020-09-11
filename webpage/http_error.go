package webpage

import (
	"fmt"
	"net/http"
)

type HTTPError struct {
	StatusCode int
	Status     string
	URL        string
}

func NewHTTPError(resp *http.Response) *HTTPError {
	return &HTTPError{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		URL:        resp.Request.URL.String(),
	}
}

func (e HTTPError) Error() string {
	return fmt.Sprintf("http status: %d(%s), url: %s", e.StatusCode, e.Status, e.URL)
}
