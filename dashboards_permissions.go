package grafana

import (
	"fmt"
	"time"
)

type DashboardsPermissionsService struct {
	client *Client
}

type DashboardPermissions []DashboardPermission

type DashboardPermission struct {
	ID             int        `json:"id,omitempty"`
	UID            string     `json:"uid,omitempty"`
	DashboardID    int        `json:"folderId,omitempty"`
	UserID         int        `json:"userId,omitempty"`
	UserLogin      string     `json:"userLogin,omitempty"`
	UserEmail      string     `json:"userEmail,omitempty"`
	TeamID         int        `json:"teamId,omitempty"`
	Team           string     `json:"team,omitempty"`
	Role           string     `json:"role,omitempty"`
	Permission     int        `json:"permission,omitempty"`
	PermissionName string     `json:"permissionName,omitempty"`
	Title          string     `json:"title,omitempty"`
	Slug           string     `json:"slug,omitempty"`
	IsFolder       bool       `json:"isFolder,omitempty"`
	URL            string     `json:"url,omitempty"`
	Created        *time.Time `json:"created,omitempty"`
	Updated        *time.Time `json:"updated,omitempty"`
}

func (s *DashboardsPermissionsService) GetDashboardPermissions(id int) (*DashboardPermissions, *Response, error) {
	u := fmt.Sprintf("dashboards/id/%d/permissions", id)

	req, err := s.client.NewRequest("GET", u, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	dp := new(DashboardPermissions)
	resp, err := s.client.Do(req, dp)
	if err != nil {
		return nil, resp, err
	}
	return dp, resp, err
}

type UpdateDashboardPermissionsOptions struct {
	Items DashboardPermissions `json:"items"`
}

func (s *DashboardsPermissionsService) UpdateDashboardPermissions(id int, opt *UpdateDashboardPermissionsOptions) (*Response, error) {
	u := fmt.Sprintf("dashboards/id/%d/permissions", id)

	req, err := s.client.NewRequest("POST", u, opt, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
