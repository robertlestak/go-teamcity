package build

import (
	"os"
	"testing"

	"github.com/robertlestak/go-teamcity/pkg/teamcity"
)

// TestBuildTriggers tests the response for BuildTriggers
func TestBuildTriggers(t *testing.T) {
	tc := teamcity.New(os.Getenv("TEAMCITY_HOST"), os.Getenv("TEAMCITY_USER"), os.Getenv("TEAMCITY_PASS"))
	c := &Config{Client: tc}
	bt, err := c.BuildTriggers(os.Getenv("BUILD_ID"))
	if err != nil {
		t.Error(err)
	}
	t.Logf("Build Triggers: %+v", bt)
}

// TestProjectTriggers tests the response for ProjectTriggers
func TestProjectTriggers(t *testing.T) {
	tc := teamcity.New(os.Getenv("TEAMCITY_HOST"), os.Getenv("TEAMCITY_USER"), os.Getenv("TEAMCITY_PASS"))
	c := &Config{Client: tc}
	bt, err := c.ProjectTriggers(os.Getenv("PROJECT_NAME"))
	if err != nil {
		t.Error(err)
	}
	t.Logf("Project Triggers: %+v", bt)
}

// TestSaveBuildTriggerState tests SaveBuildTriggerState
func TestSaveBuildTriggerState(t *testing.T) {
	tc := teamcity.New(os.Getenv("TEAMCITY_HOST"), os.Getenv("TEAMCITY_USER"), os.Getenv("TEAMCITY_PASS"))
	c := &Config{Client: tc}
	f := "/tmp/go-teamcity-test-build-trigger-state.json"
	os.Remove(f)
	err := c.SaveBuildTriggerState(os.Getenv("BUILD_ID"), f)
	if err != nil {
		t.Error(err)
	}
	if _, err := os.Stat(f); os.IsNotExist(err) {
		t.Error(err)
	}
	os.Remove(f)
}

// TestParseBuildTriggerState tests ParseBuildTriggerState
func TestParseBuildTriggerState(t *testing.T) {
	tc := teamcity.New(os.Getenv("TEAMCITY_HOST"), os.Getenv("TEAMCITY_USER"), os.Getenv("TEAMCITY_PASS"))
	c := &Config{Client: tc}
	f := "/tmp/go-teamcity-test-build-trigger-state.json"
	os.Remove(f)
	ts, err := c.BuildTriggers(os.Getenv("BUILD_ID"))
	if err != nil {
		t.Error(err.Error())
	}
	terr := c.SaveTriggerState(ts, f)
	if terr != nil {
		t.Error(err.Error())
	}
	if _, err := os.Stat(f); os.IsNotExist(err) {
		t.Error(err)
	}
	pts, perr := ParseTriggerState(f)
	if perr != nil {
		t.Error(perr)
	}
	t.Logf("Number of Build Triggers Parsed: %d", len(ts))
	if len(ts) != len(pts) {
		t.Errorf("Number of saved and parsed records do not match")
	}
	os.Remove(f)
}

// TestSaveProjectTriggerState tests SaveProjectTriggerState
func TestSaveProjectTriggerState(t *testing.T) {
	tc := teamcity.New(os.Getenv("TEAMCITY_HOST"), os.Getenv("TEAMCITY_USER"), os.Getenv("TEAMCITY_PASS"))
	c := &Config{Client: tc}
	f := "/tmp/go-teamcity-test-project-trigger-state.json"
	os.Remove(f)
	err := c.SaveProjectTriggerState(os.Getenv("PROJECT_NAME"), f)
	if err != nil {
		t.Error(err)
	}
	if _, err := os.Stat(f); os.IsNotExist(err) {
		t.Error(err)
	}
	os.Remove(f)
}

// TestParseProjectTriggerState tests ParseProjectTriggerState
func TestParseProjectTriggerState(t *testing.T) {
	tc := teamcity.New(os.Getenv("TEAMCITY_HOST"), os.Getenv("TEAMCITY_USER"), os.Getenv("TEAMCITY_PASS"))
	c := &Config{Client: tc}
	f := "/tmp/go-teamcity-test-project-trigger-state.json"
	os.Remove(f)
	ts, err := c.ProjectTriggers(os.Getenv("PROJECT_NAME"))
	if err != nil {
		t.Error(err.Error())
	}
	terr := c.SaveTriggerState(ts, f)
	if terr != nil {
		t.Error(err.Error())
	}
	if _, err := os.Stat(f); os.IsNotExist(err) {
		t.Error(err)
	}
	pts, perr := ParseTriggerState(f)
	if perr != nil {
		t.Error(perr)
	}
	t.Logf("Number of Project Triggers Parsed: %d", len(ts))
	if len(ts) != len(pts) {
		t.Errorf("Number of saved and parsed records do not match")
	}
	os.Remove(f)
}
