package build

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strings"
)

// Trigger contains trigger data
type Trigger struct {
	ID         string            `json:"id"`
	Type       string            `json:"type"`
	Disabled   bool              `json:"disabled"`
	Properties TriggerProperties `json:"properties"`
}

// TriggerProperties contains trigger property data
type TriggerProperties struct {
	Count    int               `json:"count"`
	Property []TriggerProperty `json:"property"`
}

// TriggerProperty contains trigger property data
type TriggerProperty struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// BuildTriggers returns triggers for a build ID
func (c *Config) BuildTriggers(id string) ([]Trigger, error) {
	type triggers struct {
		Count   int       `json:"count"`
		Trigger []Trigger `json:"trigger"`
	}
	rb := &triggers{}
	rd, err := c.Client.HTTPRequest("GET", "/httpAuth/app/rest/buildTypes/id:"+id+"/triggers", nil)
	if err != nil {
		return nil, err
	}
	jerr := json.Unmarshal(rd, &rb)
	if jerr != nil {
		return nil, jerr
	}
	return rb.Trigger, nil
}

// ProjectTriggers returns triggers for a project
func (c *Config) ProjectTriggers(p string) ([]Trigger, error) {
	ts, err := c.TypesForProject(p)
	if err != nil {
		return nil, err
	}
	var nts []Trigger
	for _, t := range ts {
		tts, err := c.BuildTriggers(t.ID)
		if err != nil {
			return nil, err
		}
		nts = append(nts, tts...)
	}
	return nts, nil
}

// SetBuildTriggerDisable sets disabled status a build trigger for a build
func (c *Config) SetBuildTriggerDisable(id string, t string, d bool) error {
	u := "/httpAuth/app/rest/buildTypes/id:" + id + "/triggers/" + t + "/disabled"
	c.Client.Accept = "text/plain"
	c.Client.ContentType = "text/plain"
	var sb string
	if d {
		sb = "true"
	} else {
		sb = "false"
	}
	_, err := c.Client.HTTPRequest("PUT", u, []byte(sb))
	if err != nil {
		return err
	}
	return nil
}

// DisableBuildTrigger disables a build trigger
func (c *Config) DisableBuildTrigger(id string, t string) error {
	return c.SetBuildTriggerDisable(id, t, true)
}

// EnableBuildTrigger enables a build trigger
func (c *Config) EnableBuildTrigger(id string, t string) error {
	return c.SetBuildTriggerDisable(id, t, false)
}

// SaveTriggerState saves the trigger state to file f
func (c *Config) SaveTriggerState(ts []Trigger, f string) error {
	of, ferr := os.Create(f)
	if ferr != nil {
		return ferr
	}
	defer of.Close()
	jd, jerr := json.Marshal(ts)
	if jerr != nil {
		return jerr
	}
	_, werr := of.Write(jd)
	if werr != nil {
		return werr
	}
	return nil
}

// SaveBuildTriggerState saves the build trigger state to file f
func (c *Config) SaveBuildTriggerState(id string, f string) error {
	ts, err := c.BuildTriggers(id)
	if err != nil {
		return err
	}
	terr := c.SaveTriggerState(ts, f)
	if terr != nil {
		return terr
	}
	return nil
}

// SaveProjectTriggerState saves the project trigger state to file f
func (c *Config) SaveProjectTriggerState(p string, f string) error {
	ts, err := c.ProjectTriggers(p)
	if err != nil {
		return err
	}
	terr := c.SaveTriggerState(ts, f)
	if terr != nil {
		return terr
	}
	return nil
}

// ParseTriggerState parses a trigger state array from a file
func ParseTriggerState(f string) ([]Trigger, error) {
	if _, err := os.Stat(f); os.IsNotExist(err) {
		return nil, err
	}
	var ts []Trigger
	bd, rerr := ioutil.ReadFile(f)
	if rerr != nil {
		return nil, rerr
	}
	jerr := json.Unmarshal(bd, &ts)
	if jerr != nil {
		return nil, jerr
	}
	return ts, nil
}

// TriggerStateFromFile sets the build trigger state for all Triggers
// in build id from file f
func (c *Config) TriggerStateFromFile(id string, f string) error {
	ts, err := ParseTriggerState(f)
	if err != nil {
		return err
	}
	var erstrs []string
	for _, t := range ts {
		derr := c.SetBuildTriggerDisable(id, t.ID, t.Disabled)
		if derr != nil {
			erstrs = append(erstrs, derr.Error())
		}
	}
	if len(erstrs) > 0 {
		return errors.New(strings.Join(erstrs, "; "))
	}
	return nil
}

// SaveBuildStateAndDisableAll saves current state to file f and disables all triggers
func (c *Config) SaveBuildStateAndDisableAll(id string, f string) error {
	ts, berr := c.BuildTriggers(id)
	if berr != nil {
		return berr
	}
	err := c.SaveTriggerState(ts, f)
	if err != nil {
		return err
	}
	var erstrs []string
	for _, b := range ts {
		derr := c.DisableBuildTrigger(id, b.ID)
		if derr != nil {
			erstrs = append(erstrs, derr.Error())
		}
	}
	if len(erstrs) > 0 {
		return errors.New(strings.Join(erstrs, "; "))
	}
	return nil
}
