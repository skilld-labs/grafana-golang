package grafana

import (
	"fmt"
	"time"
)

type FoldersService struct {
	client *Client
}

type Folders []Folder

type Folder struct {
	ID        int        `json:"id,omitempty"`
	UID       string     `json:"uid,omitempty"`
	Title     string     `json:"title,omitempty"`
	URL       string     `json:"url,omitempty"`
	HasACL    bool       `json:"hasAcl,omitempty"`
	CanSave   bool       `json:"canSave,omitempty"`
	CanEdit   bool       `json:"canEdit,omitempty"`
	CanAdmin  bool       `json:"canAdmin,omitempty"`
	CreatedBy string     `json:"createdBy,omitempty"`
	Created   *time.Time `json:"created,omitempty"`
	UpdatedBy string     `json:"updatedBy,omitempty"`
	Updated   *time.Time `json:"updated,omitempty"`
	Version   int        `json:"version,omitempty"`
}

func (s *FoldersService) GetAllFolders() (*Folders, *Response, error) {
	req, err := s.client.NewRequest("GET", "folders", nil, nil)
	if err != nil {
		return nil, nil, err
	}

	f := new(Folders)
	resp, err := s.client.Do(req, f)
	if err != nil {
		return nil, resp, err
	}
	return f, resp, err
}

func (s *FoldersService) GetFolderByID(id int) (*Folder, *Response, error) {
	u := fmt.Sprintf("folders/id/%d", id)

	req, err := s.client.NewRequest("GET", u, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	f := new(Folder)
	resp, err := s.client.Do(req, f)
	if err != nil {
		return nil, resp, err
	}

	return f, resp, err
}

func (s *FoldersService) GetFolderByUID(uid string) (*Folder, *Response, error) {
	u := fmt.Sprintf("folders/%s", uid)

	req, err := s.client.NewRequest("GET", u, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	f := new(Folder)
	resp, err := s.client.Do(req, f)
	if err != nil {
		return nil, resp, err
	}

	return f, resp, err
}

type CreateFolderOptions struct {
	UID   string `json:"uid,omitempty"` //optional, will be generated if not set
	Title string `json:"title,omitempty"`
}

func (s *FoldersService) CreateFolder(opt *CreateFolderOptions) (*Folder, *Response, error) {
	req, err := s.client.NewRequest("POST", "folders", opt, nil)
	if err != nil {
		return nil, nil, err
	}

	f := new(Folder)
	resp, err := s.client.Do(req, f)
	if err != nil {
		return nil, resp, err
	}

	return f, resp, err
}

type UpdateFolderByUIDOptions struct {
	UID       string `json:"uid,omitempty"` //updates the uid
	Title     string `json:"title,omitempty"`
	Version   int    `json:"version,omitempty"`
	Overwrite bool   `json:"overwrite,omitempty"`
}

func (s *FoldersService) UpdateFolderByUID(uid string, opt *UpdateFolderByUIDOptions) (*Folder, *Response, error) {
	u := fmt.Sprintf("folders/%s", uid)

	req, err := s.client.NewRequest("PUT", u, opt, nil)
	if err != nil {
		return nil, nil, err
	}

	f := new(Folder)
	resp, err := s.client.Do(req, f)
	if err != nil {
		return nil, resp, err
	}

	return f, resp, err
}

func (s *FoldersService) DeleteFolderByUID(uid string) (*Response, error) {
	u := fmt.Sprintf("folders/%s", uid)

	req, err := s.client.NewRequest("DELETE", u, nil, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
