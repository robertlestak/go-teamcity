package queue

import (
	"os"
	"testing"

	"github.com/robertlestak/go-teamcity/pkg/teamcity"
)

func TestActiveQueue(t *testing.T) {
	tc := teamcity.New(os.Getenv("TEAMCITY_HOST"), os.Getenv("TEAMCITY_USER"), os.Getenv("TEAMCITY_PASS"))
	c := &Config{Client: tc}
	_, err := c.ActiveQueue()
	if err != nil {
		t.Error(err)
	}
}

func TestActiveIDs(t *testing.T) {
	tc := teamcity.New(os.Getenv("TEAMCITY_HOST"), os.Getenv("TEAMCITY_USER"), os.Getenv("TEAMCITY_PASS"))
	c := &Config{Client: tc}
	ids, err := c.ActiveIDs()
	if err != nil {
		t.Error(err)
	}
	t.Logf("Current Queue IDs: %d", len(ids))
}
