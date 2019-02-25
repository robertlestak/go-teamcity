package teamcity

import (
	"errors"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	c := New(os.Getenv("TEAMCITY_HOST"), os.Getenv("TEAMCITY_USER"), os.Getenv("TEAMCITY_PASS"))
	if c.Host == "" || c.User == "" || c.Pass == "" {
		t.Error(errors.New("Host, User, Pass required"))
	}
}

func TestHTTPRequest(t *testing.T) {
	c := New(os.Getenv("TEAMCITY_HOST"), os.Getenv("TEAMCITY_USER"), os.Getenv("TEAMCITY_PASS"))
	rd, err := c.HTTPRequest("GET", "/httpAuth/app/rest/builds?locator=running:true", nil)
	if err != nil {
		t.Error(err)
	}
	if len(rd) == 0 {
		t.Error(errors.New("No HTTP Response"))
	}
}
