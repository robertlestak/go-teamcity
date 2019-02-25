package build

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/robertlestak/go-teamcity/pkg/teamcity"
)

// Config contains config data
type Config struct {
	Client *teamcity.Client
}

// Build contains build data
type Build struct {
	ID                 int    `json:"id"`
	BuildTypeID        string `json:"buildTypeId"`
	Number             string `json:"number"`
	Status             string `json:"status"`
	State              string `json:"string"`
	PercentageComplete int    `json:"percentageComplete"`
	HREF               string `json:"href"`
	WebURL             string `json:"webUrl"`
}

// Type contains buildType data
type Type struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Paused      bool    `json:"paused"`
	ProjectName string  `json:"projectName"`
	ProjectID   string  `json:"projectId"`
	HREF        string  `json:"href"`
	WebURL      string  `json:"webUrl"`
	Project     Project `json:"project"`
}

// Project contains project data
type Project struct {
	ID              string      `json:"id"`
	Name            string      `json:"name"`
	ParentProjectID string      `json:"parentProjectId"`
	Description     string      `json:"description"`
	Archived        bool        `json:"archived"`
	HREF            string      `json:"href"`
	WebURL          string      `json:"webUrl"`
	Parameters      Parameters  `json:"parameters"`
	Projects        ProjectList `json:"projects"`
}

// ProjectList contains project list data
type ProjectList struct {
	Count   int       `json:"count"`
	Project []Project `json:"project"`
}

// Parameters contains project parameter data
type Parameters struct {
	Count               int                `json:"count"`
	HREF                string             `json:"href"`
	ParameterProperties []ParameterPropery `json:"propery"`
}

// ParameterPropery contains project parameter propery data
type ParameterPropery struct {
	Name      string `json:"name"`
	Value     string `json:"value"`
	Inherited bool   `json:"inherited"`
}

var typeCache []*Type

// RunningBuilds returns list of all running builds
func (c *Config) RunningBuilds() ([]Build, error) {
	type runningBuilds struct {
		Count int     `json:"count"`
		Build []Build `json:"build"`
	}
	rb := &runningBuilds{}
	rd, err := c.Client.HTTPRequest("GET", "/httpAuth/app/rest/builds?locator=running:true", nil)
	if err != nil {
		return nil, err
	}
	jerr := json.Unmarshal(rd, &rb)
	if jerr != nil {
		return nil, jerr
	}
	return rb.Build, nil
}

// BuildsContainsProject checks if an array of builds contains a build for project p
func (c *Config) BuildsContainsProject(bs []Build, p string) (bool, error) {
	for _, b := range bs {
		var pid string
		var perr error
		pid, perr = c.ProjectID(b.BuildTypeID)
		if perr != nil {
			return false, perr
		}
		if strings.ToLower(p) == strings.ToLower(pid) {
			return true, nil
		}
		pid, perr = c.ParentProjectID(b.BuildTypeID)
		if perr != nil {
			return false, perr
		}
		if strings.ToLower(p) == strings.ToLower(pid) {
			return true, nil
		}
	}
	return false, nil
}

// RunningBuildsPercentages returns the percentages for all running builds
func RunningBuildsPercentages(rbs []Build, p string) string {
	var bd string
	if p != "" {
		bd = "Running Builds For " + p + ": "
	} else {
		bd = "Running Builds: "
	}
	bd += strconv.Itoa(len(rbs)) + "\n"
	for _, b := range rbs {
		bd += b.BuildTypeID + " (" + strconv.Itoa(b.ID) + "): " + strconv.Itoa(b.PercentageComplete) + "%\n"
	}
	return bd
}

// WaitForRunningBuilds waits for all running builds to complete with timeout t
// if Project/Parent Project ID p provided, only wait for these builds
func (c *Config) WaitForRunningBuilds(p string, t time.Duration) error {
	var rbs []Build
	var e error
	rbs, e = c.RunningBuilds()
	if e != nil {
		return e
	}
	if len(rbs) == 0 {
		return nil
	}
	w := make(chan error, 1)
	go func() {
		for len(rbs) > 0 {
			fmt.Print(RunningBuildsPercentages(rbs, p))
			rbs, e = c.RunningBuilds()
			if e != nil {
				w <- e
				return
			}
			if p != "" {
				cp, cerr := c.BuildsContainsProject(rbs, p)
				if cerr != nil {
					w <- e
					return
				}
				if !cp {
					w <- nil
					return
				}
			}
			time.Sleep(time.Second * 10)
		}
		w <- nil
	}()
	if t > 0 {
		select {
		case <-time.After(t):
			return errors.New("Timeout reached")
		case we := <-w:
			return we
		}
	} else {
		return <-w
	}
}

// Types returns list of all running builds
func (c *Config) Types() ([]Type, error) {
	type runningBuilds struct {
		Count int    `json:"count"`
		Type  []Type `json:"buildType"`
	}
	rb := &runningBuilds{}
	rd, err := c.Client.HTTPRequest("GET", "/httpAuth/app/rest/buildTypes", nil)
	if err != nil {
		return nil, err
	}
	jerr := json.Unmarshal(rd, &rb)
	if jerr != nil {
		return nil, jerr
	}
	return rb.Type, nil
}

// GetProject returns list of all project data
func (c *Config) GetProject(p string) (*Project, error) {
	pr := &Project{}
	rd, err := c.Client.HTTPRequest("GET", "/httpAuth/app/rest/projects/"+p, nil)
	if err != nil {
		return pr, err
	}
	jerr := json.Unmarshal(rd, &pr)
	if jerr != nil {
		return pr, jerr
	}
	return pr, nil
}

// TypesForProject returns all buildTypes for project p
func (c *Config) TypesForProject(p string) ([]Type, error) {
	ts, err := c.Types()
	if err != nil {
		return nil, err
	}
	pr, perr := c.GetProject(p)
	if perr != nil {
		return nil, perr
	}
	var nts []Type
	for _, t := range ts {
		for _, pd := range pr.Projects.Project {
			if t.ProjectID == pd.ID {
				nts = append(nts, t)
			}
		}
	}
	return nts, nil
}

// GetType gets type data for buildType ID
func (c *Config) GetType(id string) (*Type, error) {
	for _, tc := range typeCache {
		if tc.ID == id {
			return tc, nil
		}
	}
	rd, err := c.Client.HTTPRequest("GET", "/httpAuth/app/rest/buildTypes/id:"+id, nil)
	if err != nil {
		return nil, err
	}
	t := &Type{}
	jerr := json.Unmarshal(rd, &t)
	if jerr != nil {
		return t, jerr
	}
	typeCache = append(typeCache, t)
	return t, nil
}

// ProjectID returns projectId for buildType id
func (c *Config) ProjectID(id string) (string, error) {
	t, err := c.GetType(id)
	if err != nil {
		return "", err
	}
	return t.ProjectID, nil
}

// ParentProjectID returns projectId for buildType id
func (c *Config) ParentProjectID(id string) (string, error) {
	t, err := c.GetType(id)
	if err != nil {
		return "", err
	}
	return t.Project.ParentProjectID, nil
}
