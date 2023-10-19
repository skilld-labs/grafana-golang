package grafana

import (
	"fmt"
	"time"
)

type DashboardsService struct {
	client *Client
}

type Dashboard struct {
	ID            int           `json:"id,omitempty"`
	UID           string        `json:"uid,omitempty"`
	Title         string        `json:"title,omitempty"`
	Editable      bool          `json:"editable,omitempty"`
	SchemaVersion int           `json:"schemaVersion,omitempty"`
	Style         string        `json:"style,omitempty"`
	Tags          []string      `json:"tags,omitempty"`
	Annotations   *Annotations  `json:"annotations,omitempty"`
	Panels        []interface{} `json:"panels,omitempty"`
	Time          *Time         `json:"time,omitempty"`
	Templating    *Templating   `json:"templating,omitempty"`
}

type DashboardResult struct {
	Meta      Meta      `json:"meta"`
	Dashboard Dashboard `json:"dashboard"`
}

type Meta struct {
	URL     string     `json:"url"`
	Version int        `json:"version"`
	Created *time.Time `json:"created"`
	Updated *time.Time `json:"updated"`
}

type Annotations struct {
	List []interface{} `json:"list"`
}

type Time struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type Templating struct {
	List []interface{} `json:"list"`
}

func (s *DashboardsService) GetDashboardByUID(uid string) (*DashboardResult, *Response, error) {
	u := fmt.Sprintf("dashboards/uid/%s", uid)

	req, err := s.client.NewRequest("GET", u, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	d := new(DashboardResult)
	resp, err := s.client.Do(req, d)
	if err != nil {
		return nil, resp, err
	}
	return d, resp, err
}

type CreateOrUpdateDashboardOptions struct {
	Dashboard Dashboard `json:"dashboard,omitempty"`
	FolderID  int       `json:"folderId,omitempty"`
	Overwrite bool      `json:"overwrite,omitempty"`
}

// If opt.Dashboard.ID is null, dashboard will be created, else updated if opt.Overwrite is set to true.
func (s *DashboardsService) CreateOrUpdateDashboard(opt *CreateOrUpdateDashboardOptions) (*Dashboard, *Response, error) {
	u := fmt.Sprintf("dashboards/db")

	req, err := s.client.NewRequest("POST", u, opt, nil)
	if err != nil {
		return nil, nil, err
	}

	d := new(Dashboard)
	resp, err := s.client.Do(req, d)
	if err != nil {
		return nil, resp, err
	}

	return d, resp, err
}

func (s *DashboardsService) DeleteDashboard(uid string) (*Response, error) {
	u := fmt.Sprintf("dashboards/uid/%s", uid)

	req, err := s.client.NewRequest("DELETE", u, nil, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// https://grafana.com/docs/http_api/folder_dashboard_search/
// can search by id, uid, title...
// Deprecated: Use client.SearchByOptions instead with specifying SearchType.
func (s *DashboardsService) SearchDashboardByQuery(query string) (*SearchResults, *Response, error) {
	return s.client.SearchByOptions(SearchOptions{
		SearchType: SearchForDashboards,
		Query:      query,
	})
}
