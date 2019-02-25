package teamcity

import (
	"io/ioutil"
	"net/http"
)

// Client is a TeamCity client
type Client struct {
	Host        string
	User        string
	Pass        string
	Accept      string
	ContentType string
}

// New returns a new TeamCity client
func New(h string, u string, p string) *Client {
	return &Client{
		Host: h,
		User: u,
		Pass: p,
	}
}

// HTTPRequest is a generic HTTP Request to TeamCity
func (c *Client) HTTPRequest(m string, u string, b []byte) ([]byte, error) {
	hc := &http.Client{}
	req, rerr := http.NewRequest("GET", c.Host+u, nil)
	if rerr != nil {
		return nil, rerr
	}
	if c.Accept == "" {
		c.Accept = "application/json"
	}
	req.Header.Set("Accept", c.Accept)
	if c.ContentType != "" {
		req.Header.Set("Content-Type", c.ContentType)
	}
	req.SetBasicAuth(c.User, c.Pass)
	res, err := hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	bd, ierr := ioutil.ReadAll(res.Body)
	if ierr != nil {
		return nil, ierr
	}
	return bd, nil
}
