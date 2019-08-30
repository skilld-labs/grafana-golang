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
}

type DashboardResult struct {
	Meta      Meta      `json:"meta"`
	Dashboard Dashboard `json:"dashboard"`
}

type DashboardSearchResults []DashboardSearchResult

type DashboardSearchResult struct {
	ID        int      `json:"id,omitempty"`
	UID       string   `json:"uid,omitempty"`
	Title     string   `json:"title,omitempty"`
	URI       string   `json:"uri,omitempty"`
	URL       string   `json:"url,omitempty"`
	Slug      string   `json:"slug,omitempty"`
	Type      string   `json:"type,omitempty"`
	Tags      []string `json:"tags,omitempty"`
	IsStarred bool     `json:"isStarred,omitempty"`
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

//If opt.Dashboard.ID is null, dashboard will be created, else updated if opt.Overwrite is set to true.
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

//https://grafana.com/docs/http_api/folder_dashboard_search/
//can search by id, uid, title...
func (s *DashboardsService) SearchDashboardByQuery(query string) (*DashboardSearchResults, *Response, error) {
	u := fmt.Sprintf("search?query=%s&type=dash-db", query)

	req, err := s.client.NewRequest("GET", u, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	d := new(DashboardSearchResults)
	resp, err := s.client.Do(req, d)
	if err != nil {
		return nil, resp, err
	}
	return d, resp, err
}
