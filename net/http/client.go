package http

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

var client http.Client

const (
	DEFAULT_TIMEOUT = 2
)

func init() {
	client.Timeout = DEFAULT_TIMEOUT * time.Second
}

//发送http的GET或POST请求
//data如果为nil，则发GET请求，否则发POST请求
func HttpSend(host string, path string, params, header, cookies map[string]string, data []byte, timeout ...int) (body []byte, e error) {
	to := DEFAULT_TIMEOUT
	if len(timeout) > 0 {
		to = timeout[0]
	}
	return send("http", host, path, params, header, cookies, data, to)
}

//post body json
func HttpSendJson(host, path string, body map[string]interface{}, timeout ...int) (result []byte, e error) {
	b, e := json.Marshal(body)
	if e != nil {
		return
	}

	to := DEFAULT_TIMEOUT
	if len(timeout) > 0 {
		to = timeout[0]
	}

	return send("http", host, path, nil, nil, nil, b, to)
}

func HttpGet(host string, path string, params map[string]string, timeout int) (body []byte, e error) {
	return send("http", host, path, params, nil, nil, nil, timeout)
}

func HttpsGet(host string, path string, params map[string]string, timeout int) (body []byte, e error) {
	return send("https", host, path, params, nil, nil, nil, timeout)
}

func send(protocal string, host string, path string, params map[string]string, header map[string]string, cookies map[string]string, data []byte, timeout int) (body []byte, e error) {
	m := "GET"
	if data != nil {
		m = "POST"
	}

	// get param
	v := url.Values{}
	for key, value := range params {
		v.Set(key, value)
	}

	req_url := &url.URL{
		Host:     host,
		Scheme:   protocal,
		Path:     path,
		RawQuery: v.Encode(),
	}

	req, e := http.NewRequest(m, req_url.String(), bytes.NewBuffer(data))
	if e != nil {
		return nil, e
	}

	// head
	for k, v := range header {
		req.Header.Add(k, v)
	}

	// cookies
	for k, v := range cookies {
		var cookie http.Cookie
		cookie.Name = k
		cookie.Value = v
		req.AddCookie(&cookie)
	}

	c := &client
	if timeout != DEFAULT_TIMEOUT {
		c = &http.Client{}
		c.Timeout = time.Duration(timeout) * time.Second
	}

	// do req
	resp, e := c.Do(req)
	if e != nil {
		return nil, e
	}

	body, e = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return
}

//post 文件的请求
func PostFile(protocal string, host string, path string, params map[string]string, header map[string]string, cookies map[string]string, fileparam string, filename string, timeout ...int) (body []byte, e error) {
	to := DEFAULT_TIMEOUT
	if len(timeout) > 0 {
		to = timeout[0]
	}
	req_url := &url.URL{
		Host:   host,
		Scheme: protocal,
		Path:   path,
	}
	req, e := newfileUploadRequest(req_url.String(), params, fileparam, filename)
	if e != nil {
		return nil, e
	}

	for k, v := range header {
		req.Header.Add(k, v)
	}
	for k, v := range cookies {
		var cookie http.Cookie
		cookie.Name = k
		cookie.Value = v
		req.AddCookie(&cookie)
	}
	c := &client
	if to != DEFAULT_TIMEOUT {
		c = &http.Client{}
		c.Timeout = time.Duration(to) * time.Second
	}
	//fmt.Println(req.URL)
	resp, e := c.Do(req)
	if e != nil {
		return nil, e
	}
	body, e = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return
}

func newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	return http.NewRequest("POST", uri, body)
}
