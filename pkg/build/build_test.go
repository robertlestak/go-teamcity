package build

import (
	"os"
	"testing"

	"github.com/robertlestak/go-teamcity/pkg/teamcity"
)

// TestRunningBuilds tests the response for RunningBuilds
func TestRunningBuilds(t *testing.T) {
	tc := teamcity.New(os.Getenv("TEAMCITY_HOST"), os.Getenv("TEAMCITY_USER"), os.Getenv("TEAMCITY_PASS"))
	c := &Config{Client: tc}
	bs, err := c.RunningBuilds()
	if err != nil {
		t.Error(err)
	}
	t.Logf("Running Builds: %+v", bs)
}

// TestTypes tests the response for Types
func TestTypes(t *testing.T) {
	tc := teamcity.New(os.Getenv("TEAMCITY_HOST"), os.Getenv("TEAMCITY_USER"), os.Getenv("TEAMCITY_PASS"))
	c := &Config{Client: tc}
	ts, err := c.Types()
	if err != nil {
		t.Error(err)
	}
	t.Logf("Number of Types: %+v", len(ts))
}

// TestWaitForRunningBuilds tests the response for WaitForRunningBuilds
/*
func TestWaitForRunningBuilds(t *testing.T) {
	tc := teamcity.New(os.Getenv("TEAMCITY_HOST"), os.Getenv("TEAMCITY_USER"), os.Getenv("TEAMCITY_PASS"))
	c := &Config{Client: tc}
	t.Log("Wait for all builds")
	err := c.WaitForRunningBuilds("", time.Second*20)
	if err != nil {
		t.Log(err.Error())
	}
	t.Logf("Wait for %s builds", "swift")
	perr := c.WaitForRunningBuilds("swift", time.Second*20)
	if perr != nil {
		t.Log(perr.Error())
	}
}
*/

// TestProjectID tests the response for ProjectID
func TestProjectID(t *testing.T) {
	tc := teamcity.New(os.Getenv("TEAMCITY_HOST"), os.Getenv("TEAMCITY_USER"), os.Getenv("TEAMCITY_PASS"))
	c := &Config{Client: tc}
	id, err := c.ProjectID(os.Getenv("BUILD_ID"))
	if err != nil {
		t.Error(err)
	}
	t.Logf("Project ID: %s", id)
}

// TestParentProjectID tests the response for ParentProjectID
func TestParentProjectID(t *testing.T) {
	tc := teamcity.New(os.Getenv("TEAMCITY_HOST"), os.Getenv("TEAMCITY_USER"), os.Getenv("TEAMCITY_PASS"))
	c := &Config{Client: tc}
	id, err := c.ParentProjectID(os.Getenv("BUILD_ID"))
	if err != nil {
		t.Error(err)
	}
	t.Logf("Parent Project ID: %s", id)
}
