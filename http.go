package utils

import (
	"bytes"

	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func Get(uri string, timeout time.Duration, headers map[string]string, params map[string]string) (body []byte, err error) {
	if params != nil {
		u, _ := url.Parse(uri)
		values := u.Query()
		for key, value := range params {
			values.Set(key, value)
		}
		u.RawQuery = values.Encode()
		uri = u.String()
	}
	return sendHttpRequest("GET", uri, timeout, headers, nil)
}

func Post(uri string, timeout time.Duration, headers map[string]string, params map[string]string) (body []byte, err error) {
	values := url.Values{}
	for key, value := range params {
		values.Set(key, value)
	}
	return sendHttpRequest("POST", uri, timeout, headers, strings.NewReader(values.Encode()))
}

func PostRaw(uri string, timeout time.Duration, headers map[string]string, content []byte) (body []byte, err error) {
	return sendHttpRequest("POST", uri, timeout, headers, bytes.NewReader(content))
}

func sendHttpRequest(method string, uri string, timeout time.Duration, headers map[string]string, bodyReader io.Reader) (body []byte, err error) {
	req, err := http.NewRequest(method, uri, bodyReader)
	if err != nil {
		return
	}
	if host, ok := headers["Host"]; ok {
		req.Host = host
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

	return
}
