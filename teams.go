package grafana

import (
	"fmt"
)

type TeamsService struct {
	client *Client
}

type Teams []Team

type Team struct {
	ID          int    `json:"id,omitempty"`
	OrgID       int    `json:"orgId,omitempty"`
	Name        string `json:"name,omitempty"`
	Email       string `json:"email,omitempty"`
	AvatarUrl   string `json:"avatarUrl,omitempty"`
	MemberCount int    `json:"memberCount,omitempty"`
	Permission  int    `json:"permission,omitempty"`
}

type TeamSearchResults struct {
	Teams      *Teams `json:"teams,omitempty"`
	TotalCount int    `json:"totalCount,omitempty"`
	Page       int    `json:"page,omitempty"`
	PerPage    int    `json:"perPage,omitempty"`
}

type TeamMembers []TeamMember

type TeamMember struct {
	OrgID     int    `json:"orgId,omitempty"`
	TeamID    int    `json:"teamId,omitempty"`
	UserID    int    `json:"userId,omitempty"`
	Email     string `json:"email,omitempty"`
	Login     string `json:"login,omitempty"`
	AvatarUrl string `json:"avatarUrl,omitempty"`
}

type TeamPreferences struct {
	Theme      string `json:"theme,omitempty"`
	HomeTeamID int    `json:"homeTeamId,omitempty"`
	Timezone   string `json:"timezone,omitempty"`
}

func (s *TeamsService) GetTeamByID(id int) (*Team, *Response, error) {
	u := fmt.Sprintf("teams/%d", id)

	req, err := s.client.NewRequest("GET", u, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	t := new(Team)
	resp, err := s.client.Do(req, t)
	if err != nil {
		return nil, resp, err
	}
	return t, resp, err
}

type CreateTeamOptions struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

func (s *TeamsService) CreateTeam(opt *CreateTeamOptions) (*Response, error) {
	req, err := s.client.NewRequest("POST", "teams", opt, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

type UpdateTeamByIDOptions CreateTeamOptions //set all fields or else existing fields will be removed

func (s *TeamsService) UpdateTeamByID(id int, opt *UpdateTeamByIDOptions) (*Response, error) {
	u := fmt.Sprintf("teams/%d", id)

	req, err := s.client.NewRequest("PUT", u, opt, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

func (s *TeamsService) DeleteTeamByID(id int) (*Response, error) {
	u := fmt.Sprintf("teams/%d", id)

	req, err := s.client.NewRequest("DELETE", u, nil, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

func (s *TeamsService) GetTeamMembersByTeamID(id int) (*TeamMembers, *Response, error) {
	u := fmt.Sprintf("teams/%d/members", id)

	req, err := s.client.NewRequest("GET", u, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	tm := new(TeamMembers)
	resp, err := s.client.Do(req, tm)
	if err != nil {
		return nil, resp, err
	}
	return tm, resp, err
}

type AddTeamMemberOptions struct {
	UserID int `json:"userId,omitempty"`
}

func (s *TeamsService) AddTeamMember(id int, opt *AddTeamMemberOptions) (*Response, error) {
	u := fmt.Sprintf("teams/%d/members", id)

	req, err := s.client.NewRequest("POST", u, opt, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

func (s *TeamsService) RemoveMemberFromTeam(teamID, userID int) (*Response, error) {
	u := fmt.Sprintf("teams/%d/members/%d", teamID, userID)

	req, err := s.client.NewRequest("DELETE", u, nil, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

func (s *TeamsService) GetTeamPreferences(id int) (*TeamPreferences, *Response, error) {
	u := fmt.Sprintf("teams/%d/preferences", id)

	req, err := s.client.NewRequest("GET", u, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	tp := new(TeamPreferences)
	resp, err := s.client.Do(req, tp)
	if err != nil {
		return nil, resp, err
	}

	return tp, resp, err
}

type UpdateTeamPreferencesOptions TeamPreferences

func (s *TeamsService) UpdateTeamPreferences(id int, opt *UpdateTeamPreferencesOptions) (*Response, error) {
	u := fmt.Sprintf("teams/%d/preferences", id)

	req, err := s.client.NewRequest("PUT", u, opt, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

//https://grafana.com/docs/http_api/team/#team-search-with-paging
//can search by id, uid, name...
func (s *TeamsService) SearchTeamByQuery(query string) (*TeamSearchResults, *Response, error) {
	u := fmt.Sprintf("teams/search?query=%s", query)

	req, err := s.client.NewRequest("GET", u, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	tsr := new(TeamSearchResults)
	resp, err := s.client.Do(req, tsr)
	if err != nil {
		return nil, resp, err
	}
	return tsr, resp, err
}
