package telegoat

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFormRequestBuffer(t *testing.T) {
	requestBody := map[string]interface{}{
		"test": "string",
	}

	requestBuffer, err := formRequestBuffer(requestBody)

	if err != nil {
		t.Fatal(err)
	}

	var unmarshaledBody map[string]string

	json.Unmarshal(requestBuffer.Bytes(), &unmarshaledBody)

	if fmt.Sprint(requestBody) != fmt.Sprint(unmarshaledBody) {
		t.Fatalf(
			"Json changed after unmarshalling. Before: %v, after: %v",
			requestBody,
			unmarshaledBody,
		)
	}

}

func TestPostRequest(t *testing.T) {
	expectedResponseBody := map[string]string{
		"hello": "world",
	}
	marshalledBody, err := json.Marshal(expectedResponseBody)
	if err != nil {
		t.Error(err)
	}
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Write(marshalledBody)
		}),
	)
	defer ts.Close()

	requestBody := map[string]interface{}{
		"test": "body",
	}

	body, err := postRequest(ts.URL, requestBody)
	if err != nil {
		t.Fatal(err)
	}

	var responseBody map[string]string
	json.Unmarshal(body, &responseBody)

	if fmt.Sprint(expectedResponseBody) != fmt.Sprint(responseBody) {
		t.Fatalf("Unexpected response: %v != %v", expectedResponseBody, responseBody)
	}

}

func TestPostRequestNotOkResponse(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}),
	)
	defer ts.Close()

	requestBody := make(map[string]interface{})
	_, err := postRequest(ts.URL, requestBody)

	if err.Error() != "500 status from external server" {
		t.Fatalf("Invalid error: %v", err.Error())
	}
}
