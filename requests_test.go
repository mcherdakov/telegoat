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

	var unmarshalledBody map[string]string

	err = json.Unmarshal(requestBuffer.Bytes(), &unmarshalledBody)
	if err != nil {
		t.Fatal(err)
	}

	if fmt.Sprint(requestBody) != fmt.Sprint(unmarshalledBody) {
		t.Fatalf(
			"Json changed after unmarshalling. Before: %v, after: %v",
			requestBody,
			unmarshalledBody,
		)
	}

}

func TestPostRequest(t *testing.T) {
	expectedResponseBody := map[string]string{
		"hello": "world",
	}
	marshalledBody, err := json.Marshal(expectedResponseBody)
	if err != nil {
		t.Fatal(err)
	}
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			_, err := w.Write(marshalledBody)
			if err != nil {
				t.Fatal(err)
			}
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
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		t.Fatal(err)
	}

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

	if err == nil {
		t.Fatal("Expected error")
	}

	if err.Error() != "500 status from external server" {
		t.Fatalf("Invalid error: %v", err.Error())
	}
}

func TestPostRequestWithError(t *testing.T) {
	requestBody := make(map[string]interface{})
	_, err := postRequest("fake_url", requestBody)

	if err == nil {
		t.Fatal("Expected error")
	}
}
