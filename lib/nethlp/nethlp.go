package nethlp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func GetJSON(ret interface{}, url string, headers map[string]string) (*http.Response, error) {
	return Request("GET", ret, url, headers, nil)
}

func Request(method string, ret interface{}, url string, headers map[string]string, payload interface{}) (*http.Response, error) {
	var payloadr io.Reader
	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		payloadr = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, url, payloadr)

	if headers != nil {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return rsp, err
	}

	err = parseResponse(ret, rsp)
	return rsp, err
}

func parseResponse(v interface{}, rsp *http.Response) error {
	body := rsp.Body
	defer body.Close()

	// body := bytes.NewBuffer([]byte(`{"error":"check your params","code":1001}`))
	// b, err := ioutil.ReadAll(rsp.Body)
	// fmt.Println(string(b), err)

	if err := json.NewDecoder(body).Decode(v); err != nil {
		return errors.New(fmt.Sprintf("[%v]%v", rsp.StatusCode, err))
	}

	if rsp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("%v", rsp.StatusCode))
	}

	return nil
}
