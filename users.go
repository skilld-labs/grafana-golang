package grafana

import (
	"fmt"
)

type UsersService struct {
	client *Client
}

type Users []User

type User struct {
	ID             int    `json:"id,omitempty"`
	Email          string `json:"email,omitempty"`
	Name           string `json:"name,omitempty"`
	Login          string `json:"login,omitempty"`
	Theme          string `json:"theme,omitempty"`
	OrgID          int    `json:"orgId,omitempty"`
	IsGrafanaAdmin bool   `json:"isGrafanaAdmin,omitempty"`
}

func (s *UsersService) GetUserByID(id int) (*User, *Response, error) {
	u := fmt.Sprintf("users/%d", id)

	req, err := s.client.NewRequest("GET", u, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	usr := new(User)
	resp, err := s.client.Do(req, usr)
	if err != nil {
		return nil, resp, err
	}
	return usr, resp, err
}

func (s *UsersService) GetUserByLoginOrEmail(id string) (*User, *Response, error) {
	u := fmt.Sprintf("users/lookup?loginOrEmail=%s", id)

	req, err := s.client.NewRequest("GET", u, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	usr := new(User)
	resp, err := s.client.Do(req, usr)
	if err != nil {
		return nil, resp, err
	}
	return usr, resp, err
}

type CreateUserOptions struct {
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}

//must be authenticated as admin to do this
func (s *UsersService) CreateUser(opt *CreateUserOptions) (*Response, error) {
	req, err := s.client.NewRequest("POST", "admin/users", opt, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

type UpdateUserOptions struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	Login string `json:"login,omitempty"`
	Theme string `json:"theme,omitempty"`
}

func (s *UsersService) UpdateUserByID(id int, opt *UpdateUserOptions) (*Response, error) {
	u := fmt.Sprintf("users/%d", id)

	req, err := s.client.NewRequest("PUT", u, opt, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

//must be authenticated as admin to do this
func (s *UsersService) DeleteUserByID(id int) (*Response, error) {
	u := fmt.Sprintf("admin/users/%d", id)

	req, err := s.client.NewRequest("DELETE", u, nil, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

func (s *UsersService) GetUserTeams(id int) (*Teams, *Response, error) {
	u := fmt.Sprintf("users/%d/teams", id)

	req, err := s.client.NewRequest("GET", u, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	t := new(Teams)
	resp, err := s.client.Do(req, t)
	if err != nil {
		return nil, resp, err
	}

	return t, resp, err
}
