package queue

import (
	"encoding/json"
	"strconv"

	"github.com/robertlestak/go-teamcity/pkg/build"
	"github.com/robertlestak/go-teamcity/pkg/teamcity"
)

// Config contains config data
type Config struct {
	Client       *teamcity.Client
	CancelReason string
}

// ActiveQueue returns list of all queued builds
func (c *Config) ActiveQueue() ([]build.Build, error) {
	type queuedBuilds struct {
		Build []build.Build `json:"build"`
	}
	qb := &queuedBuilds{}
	rd, err := c.Client.HTTPRequest("GET", "/httpAuth/app/rest/buildQueue", nil)
	if err != nil {
		return nil, err
	}
	jerr := json.Unmarshal(rd, &qb)
	if jerr != nil {
		return nil, jerr
	}
	return qb.Build, nil
}

// ActiveIDs returns just IDs for queued builds
func (c *Config) ActiveIDs() ([]int, error) {
	var i []int
	bs, err := c.ActiveQueue()
	if err != nil {
		return nil, err
	}
	for _, d := range bs {
		i = append(i, d.ID)
	}
	return i, nil
}

// CancelBuild cancels a given build ID
func (c *Config) CancelBuild(i int, cr string) ([]byte, error) {
	crs := "<buildCancelRequest comment='" + cr + "' readdIntoQueue='false'/>"
	c.Client.ContentType = "application/xml"
	rd, err := c.Client.HTTPRequest("GET", "/httpAuth/app/rest/builds/id:"+strconv.Itoa(i), []byte(crs))
	if err != nil {
		return nil, err
	}
	return rd, nil
}

// DeleteBuild deletes a given build ID
func (c *Config) DeleteBuild(i int) ([]byte, error) {
	rd, err := c.Client.HTTPRequest("DELETE", "/httpAuth/app/rest/builds/id:"+strconv.Itoa(i), nil)
	if err != nil {
		return nil, err
	}
	return rd, nil
}

// CancelAndDeleteBuild cancels and deletes a given build ID and cancelation reason
func (c *Config) CancelAndDeleteBuild(i int, cr string) ([]byte, error) {
	var od []byte
	cd, cerr := c.CancelBuild(i, cr)
	if cerr != nil {
		return od, cerr
	}
	od = append(od, cd...)
	d, derr := c.DeleteBuild(i)
	if derr != nil {
		return od, derr
	}
	od = append(od, d...)
	return od, nil
}

// cancelAndDeleteWorker concurrent worker for mass cancel and delete requests
func cancelAndDeleteWorker(c *Config, req chan int, res chan int) {
	for r := range req {
		c.CancelAndDeleteBuild(r, c.CancelReason)
		res <- r
	}
}

// ClearQueue clears all builds in queue
func (c *Config) ClearQueue() error {
	ids, err := c.ActiveIDs()
	if err != nil {
		return err
	}
	if len(ids) == 0 {
		return nil
	}
	req := make(chan int, len(ids))
	res := make(chan int, len(ids))
	for i := 0; i <= 100; i++ {
		go cancelAndDeleteWorker(c, req, res)
	}
	for j := 0; j < len(ids); j++ {
		req <- ids[j]
	}
	close(req)
	for a := 0; a < len(ids); a++ {
		<-res
	}
	return nil
}
