package grafana

import (
	"fmt"
	"time"
)

type FoldersPermissionsService struct {
	client *Client
}

type FolderPermissions []FolderPermission

type FolderPermission struct {
	ID             int        `json:"id,omitempty"`
	UID            string     `json:"uid,omitempty"`
	FolderID       int        `json:"folderId,omitempty"`
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

func (s *FoldersPermissionsService) GetFolderPermissions(uid string) (*FolderPermissions, *Response, error) {
	u := fmt.Sprintf("folders/%s/permissions", uid)

	req, err := s.client.NewRequest("GET", u, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	fp := new(FolderPermissions)
	resp, err := s.client.Do(req, fp)
	if err != nil {
		return nil, resp, err
	}
	return fp, resp, err
}

type UpdateFolderPermissionsOptions struct {
	Items FolderPermissions `json:"items"`
}

func (s *FoldersPermissionsService) UpdateFolderPermissions(uid string, opt *UpdateFolderPermissionsOptions) (*Response, error) {
	u := fmt.Sprintf("folders/%s/permissions", uid)

	req, err := s.client.NewRequest("POST", u, opt, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
