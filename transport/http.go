package transport

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// HTTPGet makes a request to a remote endpoint and returns the response as a string
func HTTPGet(endpoint string, headers map[string]string) []byte {
	req, _ := http.NewRequest("GET", endpoint, nil)

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := HTTPClient().Do(req)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return body
}

// HTTPPost makes a POST request to a remote endpoint and returns the response as a string
func HTTPPost(endpoint string, headers map[string]string, body interface{}) []byte {
	str, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer(str))

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// Always make sure we send the Content-Length
	req.Header.Set("Content-Length", strconv.Itoa(len(str)))

	resp, err := HTTPClient().Do(req)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Handle HTTP status codes here
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return data
}

// HTTPClient returns a client with a 10 second timeout
func HTTPClient() *http.Client {
	return &http.Client{
		Timeout: time.Second * 10,
	}
}
