package webutility

// TODO(marko): If DecodeJSON() returns io.EOF treat it as if there is no response body, since response content length can sometimes be -1.

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

func GetContent(url string, params url.Values, headers http.Header) (content []byte, status int, err error) {
	if params != nil {
		p := params.Encode()
		if p != "" {
			url += "?" + p
		}
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, err
	}

	if headers != nil {
		for k, head := range headers {
			for i, h := range head {
				if i == 0 {
					req.Header.Set(k, h)
				} else {
					req.Header.Add(k, h)
				}
			}
		}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	status = resp.StatusCode

	if status != http.StatusOK {
		return nil, status, err
	}

	content, err = ioutil.ReadAll(resp.Body)

	return content, status, err
}

// DecodeJSON decodes JSON data from r to v.
// Returns an error if it fails.
func DecodeJSON(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

func GetJSON(url string, v interface{}, params url.Values, headers http.Header) (status int, err error) {
	if params != nil {
		p := params.Encode()
		if p != "" {
			url += "?" + p
		}
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}

	if headers != nil {
		for k, head := range headers {
			for i, h := range head {
				if i == 0 {
					req.Header.Set(k, h)
				} else {
					req.Header.Add(k, h)
				}
			}
		}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	status = resp.StatusCode

	if status != http.StatusOK {
		return status, err
	}

	if err = DecodeJSON(resp.Body, v); err == io.EOF {
		err = nil
	}

	return status, err
}

func PostJSON(url string, v, r interface{}, params url.Values, headers http.Header) (status int, err error) {
	buffer := bytes.NewBuffer(make([]byte, 0))
	json.NewEncoder(buffer).Encode(v)

	if params != nil {
		p := params.Encode()
		if p != "" {
			url += "?" + p
		}
	}

	req, err := http.NewRequest(http.MethodPost, url, buffer)
	if err != nil {
		return 0, err
	}

	if headers != nil {
		for k, head := range headers {
			for i, h := range head {
				if i == 0 {
					req.Header.Set(k, h)
				} else {
					req.Header.Add(k, h)
				}
			}
		}
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	status = resp.StatusCode
	defer resp.Body.Close()

	if status != http.StatusOK && status != http.StatusCreated {
		return status, err
	}

	if err = DecodeJSON(resp.Body, v); err == io.EOF {
		err = nil
	}

	return status, err
}
