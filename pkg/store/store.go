package store

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

type AssignResp struct {
	Fid       string `json:"fid"`
	URL       string `json:"url"`
	PublicURL string `json:"publicUrl"`
	Count     int    `json:"count"`
	Error     string `json:"error"`
}

// Assign gets the designated url to store the media
func Assign(weedAddr string) (*AssignResp, error) {
	dirurl := fmt.Sprintf("%s/dir/assign", weedAddr)
	resp, err := http.Get(dirurl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	assignment := &AssignResp{}
	if err := json.NewDecoder(resp.Body).Decode(assignment); err != nil {
		return nil, err
	}
	if assignment.Error != "" {
		return nil, errors.New(assignment.Error)
	}
	return assignment, nil
}

type UploadResp struct {
	Name  int    `json:"name"`
	Size  int    `json:"size"`
	ETag  string `json:"eTag"`
	Error string `json:"error"`
}

// Upload sends the form to the URL
// courtesy of https://stackoverflow.com/a/20397167/6056991
func Upload(fidURL string, form map[string]io.Reader) (*UploadResp, error) {
	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range form {
		var fw io.Writer
		var err error
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		// Add an image file
		if x, ok := r.(*os.File); ok {
			if fw, err = w.CreateFormFile(key, x.Name()); err != nil {
				return nil, err
			}
		} else {
			// Add other fields
			if fw, err = w.CreateFormField(key); err != nil {
				return nil, err
			}
		}
		if _, err := io.Copy(fw, r); err != nil {
			return nil, err
		}
	}
	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", fidURL, &b)
	if err != nil {
		return nil, err
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Submit the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	upload := &UploadResp{}
	if err = json.NewDecoder(resp.Body).Decode(upload); err != nil {
		return nil, err
	}

	if upload.Error != "" {
		return nil, errors.New(upload.Error)
	}
	return upload, nil
}
