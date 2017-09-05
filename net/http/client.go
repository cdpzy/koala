package http

import (
	"bytes"
	"io/ioutil"
	nhttp "net/http"
)

func Get(url string) ([]byte, error) {
	resp, err := nhttp.Get(url)
	if err != nil {
		return make([]byte, 0), nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}

func PostJSON(url string, data []byte) ([]byte, error) {
	req, err := nhttp.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return make([]byte, 0), nil
	}

	defer req.Body.Close()

	req.Header.Set("X-Custom-Header", "val")
	req.Header.Set("Content-Type", "application/json")

	client := &nhttp.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return make([]byte, 0), nil
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}
