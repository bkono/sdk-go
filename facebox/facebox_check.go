package facebox

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/url"

	"github.com/pkg/errors"
)

// Check checks the image in the io.Reader for faces.
func (c *Client) Check(image io.Reader) ([]Face, error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	fw, err := w.CreateFormFile("file", "image.dat")
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(fw, image)
	if err != nil {
		return nil, err
	}
	if err = w.Close(); err != nil {
		return nil, err
	}
	u, err := url.Parse(c.addr + "/facebox/check")
	if err != nil {
		return nil, err
	}
	if !u.IsAbs() {
		return nil, errors.New("box address must be absolute")
	}
	resp, err := c.HTTPClient.Post(u.String(), w.FormDataContentType(), &buf)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return nil, errors.New(string(b))
	}

	return c.parseCheckResponse(resp.Body)
}

// CheckURL checks the image at the specified URL for faces.
func (c *Client) CheckURL(imageURL *url.URL) ([]Face, error) {
	u, err := url.Parse(c.addr + "/facebox/check")
	if err != nil {
		return nil, err
	}
	if !u.IsAbs() {
		return nil, errors.New("box address must be absolute")
	}
	if !imageURL.IsAbs() {
		return nil, errors.New("url must be absolute")
	}
	form := url.Values{}
	form.Set("url", imageURL.String())
	resp, err := c.HTTPClient.PostForm(u.String(), form)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return c.parseCheckResponse(resp.Body)
}

type checkResponse struct {
	Success    bool   `json:"success"`
	FacesCount int    `json:"facesCount"`
	Error      string `json:"error,omitempty"`
	Faces      []Face `json:"faces"`
}

func (c *Client) parseCheckResponse(r io.Reader) ([]Face, error) {
	var resp checkResponse
	if err := json.NewDecoder(r).Decode(&resp); err != nil {
		return nil, errors.Wrap(err, "decoding response")
	}
	if !resp.Success {
		return nil, ErrFacebox(resp.Error)
	}
	return resp.Faces, nil
}
