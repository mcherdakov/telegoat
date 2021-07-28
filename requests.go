package telegoat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var httpClient = http.Client{
	Timeout: time.Minute,
}

func formRequestBuffer(requestBody map[string]interface{}) (*bytes.Buffer, error) {
	bytesRepresentation, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(bytesRepresentation), nil
}

func postRequest(url string, requestBody map[string]interface{}) ([]byte, error) {
	requestBuffer, err := formRequestBuffer(requestBody)
	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Post(
		url,
		"application/json",
		requestBuffer,
	)

	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%v status from external server", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
