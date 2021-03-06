package utils

import (
	"bytes"
	"fmt"

	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type GPP struct {
	Uri     string
	Timeout time.Duration
	Headers map[string]string
	Params  interface{}
}

func Get(gpp *GPP) (body []byte, err error) {
	uri, timeout, headers, params := gpp.Uri, gpp.Timeout, gpp.Headers, gpp.Params
	switch v := params.(type) {
	case map[string]string:
		u, _ := url.Parse(uri)
		values := u.Query()
		for key, value := range v {
			values.Set(key, value)
		}
		u.RawQuery = values.Encode()
		uri = u.String()
	}
	return sendHttpRequest("GET", uri, timeout, headers, nil)
}

func Post(gpp *GPP) (body []byte, err error) {
	uri, timeout, headers, params := gpp.Uri, gpp.Timeout, gpp.Headers, gpp.Params
	var reader io.Reader
	switch v := params.(type) {
	case map[string]string:
		values := url.Values{}
		for key, value := range v {
			values.Set(key, value)
		}
		reader = strings.NewReader(values.Encode())
	case string:
		reader = strings.NewReader(v)
	case []byte:
		reader = bytes.NewReader(v)
	}
	return sendHttpRequest("POST", uri, timeout, headers, reader)
}

func sendHttpRequest(
	method string,
	uri string,
	timeout time.Duration,
	headers map[string]string,
	bodyReader io.Reader,
) (body []byte, err error) {
	req, err := http.NewRequest(method, uri, bodyReader)
	if err != nil {
		return
	}
	if host, ok := headers["Host"]; ok {
		req.Host = host
		delete(headers, "Host")
	}
	if connection, ok := headers["Connection"]; ok {
		if strings.ToLower(connection) == "close" {
			req.Close = true
			delete(headers, "Connection")
		}
	}
	for name, value := range headers {
		req.Header.Set(name, value)
	}
	if bodyReader != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	client := &http.Client{
		Timeout: timeout,
	}
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if g, e := resp.StatusCode, http.StatusOK; g != e {
		err = fmt.Errorf("http resp code: %d", g)
		return
	}

	return
}
