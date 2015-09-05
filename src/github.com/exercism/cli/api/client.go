package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/exercism/cli/config"
)

const (
	apiIssueTracker  = "https://github.com/exercism/exercism.io/issues"
	xapiIssueTracker = "https://github.com/exercism/x-api/issues"
)

var (
	// UserAgent lets the API know where the call is being made from.
	// It's set from main() so that we have access to the version.
	UserAgent string
)

// Client contains the necessary information to contact the Exercism APIs
type Client struct {
	client   *http.Client
	APIHost  string
	XAPIHost string
	APIKey   string
}

// NewClient returns an Exercism API Client
func NewClient(c *config.Config) *Client {
	return &Client{
		client:   http.DefaultClient,
		APIHost:  c.API,
		XAPIHost: c.XAPI,
		APIKey:   c.APIKey,
	}
}

// NewRequest returns an http.Request with information for the Exercism API
func (c *Client) NewRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

// Do performs an http.Request and optionally parses the response body into the given interface
func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	switch res.StatusCode {
	case http.StatusNoContent:
		return res, nil
	case http.StatusInternalServerError:
		issueTracker := apiIssueTracker
		if strings.Contains(req.URL.Host, "x.exercism.io") {
			issueTracker = xapiIssueTracker
		}
		return nil, fmt.Errorf("an internal server error was received.\nPlease file a bug report with the contents of 'exercism debug' at: %s ", issueTracker)
	default:
		if v != nil {
			defer res.Body.Close()
			if err := json.NewDecoder(res.Body).Decode(v); err != nil {
				return nil, fmt.Errorf("error parsing API response - %s", err)
			}
		}
	}

	return res, nil
}
