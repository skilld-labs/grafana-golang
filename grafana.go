package grafana

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
)

const (
	defaultBaseURL = "https://grafana.com/api/"
	userAgent      = "go-grafana"
)

type authType int

const (
	privateToken authType = iota
	basicAuth
)

type ISOTime time.Time

// ISO 8601 date format
const iso8601 = "2006-01-02"

// MarshalJSON implements the json.Marshaler interface
func (t ISOTime) MarshalJSON() ([]byte, error) {
	if y := time.Time(t).Year(); y < 0 || y >= 10000 {
		// ISO 8901 uses 4 digits for the years
		return nil, errors.New("ISOTime.MarshalJSON: year outside of range [0,9999]")
	}

	b := make([]byte, 0, len(iso8601)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, iso8601)
	b = append(b, '"')

	return b, nil
}

func (t *ISOTime) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package
	if string(data) == "null" {
		return nil
	}

	isotime, err := time.Parse(`"`+iso8601+`"`, string(data))
	*t = ISOTime(isotime)

	return err
}

func (t ISOTime) String() string {
	return time.Time(t).Format(iso8601)
}

type Client struct {
	client   *http.Client
	Username string
	Password string

	baseURL *url.URL

	authType authType

	token string

	UserAgent string

	Dashboards            *DashboardsService
	DashboardsPermissions *DashboardsPermissionsService
	Folders               *FoldersService
	FoldersPermissions    *FoldersPermissionsService
	Teams                 *TeamsService
	Users                 *UsersService
}

func NewTokenClient(httpClient *http.Client, endpoint, token string) (*Client, error) {
	client := newClient(httpClient)
	if err := client.SetBaseURL(endpoint + "/api"); err != nil {
		return nil, err
	}
	client.authType = privateToken
	client.token = token
	return client, nil
}

func NewBasicAuthClient(httpClient *http.Client, endpoint, username, password string) (*Client, error) {
	client := newClient(httpClient)
	if err := client.SetBaseURL(endpoint + "/api"); err != nil {
		return nil, err
	}
	client.Username = username
	client.Password = password
	client.authType = basicAuth
	return client, nil
}

func newClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	c := &Client{client: httpClient, UserAgent: userAgent}
	if err := c.SetBaseURL(defaultBaseURL); err != nil {
		// Should never happen since defaultBaseURL is our constant.
		panic(err)
	}

	c.Dashboards = &DashboardsService{client: c}
	c.DashboardsPermissions = &DashboardsPermissionsService{client: c}
	c.Folders = &FoldersService{client: c}
	c.FoldersPermissions = &FoldersPermissionsService{client: c}
	c.Teams = &TeamsService{client: c}
	c.Users = &UsersService{client: c}

	return c
}

func (c *Client) BaseURL() *url.URL {
	u := *c.baseURL
	return &u
}

func (c *Client) SetBaseURL(urlStr string) error {
	// Make sure the given URL end with a slash
	if !strings.HasSuffix(urlStr, "/") {
		urlStr += "/"
	}

	var err error
	c.baseURL, err = url.Parse(urlStr)
	return err
}

func (c *Client) NewRequest(method, path string, opt interface{}, options []OptionFunc) (*http.Request, error) {
	u := *c.baseURL
	// Set the encoded opaque data
	u.Opaque = c.baseURL.Path + path

	if opt != nil {
		q, err := query.Values(opt)
		if err != nil {
			return nil, err
		}
		u.RawQuery = q.Encode()
	}

	req := &http.Request{
		Method:     method,
		URL:        &u,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(http.Header),
		Host:       u.Host,
	}

	for _, fn := range options {
		if fn == nil {
			continue
		}

		if err := fn(req); err != nil {
			return nil, err
		}
	}

	if method == "POST" || method == "PUT" {
		bodyBytes, err := json.Marshal(opt)
		if err != nil {
			return nil, err
		}
		bodyReader := bytes.NewReader(bodyBytes)

		u.RawQuery = ""
		req.Body = ioutil.NopCloser(bodyReader)
		req.ContentLength = int64(bodyReader.Len())
		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("Accept", "application/json")

	switch c.authType {
	case privateToken:
		req.Header.Set("Authorization", "Bearer "+c.token)
	case basicAuth:
		req.SetBasicAuth(c.Username, c.Password)
	}

	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}

	return req, nil
}

type Response struct {
	*http.Response
}

func newResponse(r *http.Response) *Response {
	return &Response{Response: r}
}

func (c *Client) Do(req *http.Request, v interface{}) (*Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	response := newResponse(resp)

	err = CheckResponse(resp)
	if err != nil {
		return response, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
		}
	}

	return response, err
}

type ErrorResponse struct {
	Response *http.Response
	Message  string
}

func (e *ErrorResponse) Error() string {
	path, _ := url.QueryUnescape(e.Response.Request.URL.Opaque)
	u := fmt.Sprintf("%s://%s%s", e.Response.Request.URL.Scheme, e.Response.Request.URL.Host, path)
	return fmt.Sprintf("%s %s: %d %s", e.Response.Request.Method, u, e.Response.StatusCode, e.Message)
}

func CheckResponse(r *http.Response) error {
	switch r.StatusCode {
	case 200, 201, 202, 204, 304:
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		var raw interface{}
		if err := json.Unmarshal(data, &raw); err != nil {
			errorResponse.Message = "failed to parse unknown error format"
		}

		errorResponse.Message = parseError(raw)
	}

	return errorResponse
}

// Format:
// {
//     "message": {
//         "<property-name>": [
//             "<error-message>",
//             "<error-message>",
//             ...
//         ],
//         "<embed-entity>": {
//             "<property-name>": [
//                 "<error-message>",
//                 "<error-message>",
//                 ...
//             ],
//         }
//     },
//     "error": "<error-message>"
// }
func parseError(raw interface{}) string {
	switch raw := raw.(type) {
	case string:
		return raw

	case []interface{}:
		var errs []string
		for _, v := range raw {
			errs = append(errs, parseError(v))
		}
		return fmt.Sprintf("[%s]", strings.Join(errs, ", "))

	case map[string]interface{}:
		var errs []string
		for k, v := range raw {
			errs = append(errs, fmt.Sprintf("{%s: %s}", k, parseError(v)))
		}
		sort.Strings(errs)
		return strings.Join(errs, ", ")

	default:
		return fmt.Sprintf("failed to parse unexpected error type: %T", raw)
	}
}

type OptionFunc func(*http.Request) error
