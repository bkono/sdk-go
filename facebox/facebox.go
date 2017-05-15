// Package facebox provides a client for accessing facebox services.
package facebox

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/machinebox/sdk-go/x/boxutil"
)

// Face represents a face in an image.
type Face struct {
	Rect    Rect   `json:"rect"`
	ID      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Matched bool   `json:"matched"`
}

// Rect represents the coordinates of a face within an image.
type Rect struct {
	Top    int `json:"top"`
	Left   int `json:"left"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

// Similar represents a similar face.
type Similar struct {
	ID   string
	Name string
}

// Client is an HTTP client that can make requests to the box.
type Client struct {
	addr string

	// HTTPClient is the http.Client that will be used to
	// make requests.
	HTTPClient *http.Client
}

// make sure the Client implements boxutil.Box
var _ boxutil.Box = (*Client)(nil)

// New creates a new Client.
func New(addr string) *Client {
	return &Client{
		addr: addr,
		HTTPClient: &http.Client{
			Timeout: 1 * time.Minute,
		},
	}
}

// Info gets the details about the box.
func (c *Client) Info() (*boxutil.Info, error) {
	var info boxutil.Info
	u, err := url.Parse(c.addr + "/info")
	if err != nil {
		return nil, err
	}
	if !u.IsAbs() {
		return nil, errors.New("box address must be absolute")
	}
	resp, err := c.HTTPClient.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}
	return &info, nil
}

// ErrFacebox represents an error from nudebox.
type ErrFacebox string

func (e ErrFacebox) Error() string {
	return "facebox: " + string(e)
}
