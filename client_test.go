package telegoat

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendMessage(t *testing.T) {
	expectedCallBody := map[string]interface{}{
		"chat_id": 123,
		"text":    "message",
	}
	callBody := make(map[string]interface{})

	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, _ := ioutil.ReadAll(r.Body)
			err := json.Unmarshal(body, &callBody)
			if err != nil {
				t.Fatal(err)
			}
		}),
	)
	defer ts.Close()

	telegramClient := NewTelegramClient("fake_token")
	telegramClient.telegramURL = ts.URL

	err := telegramClient.SendMessage("message", 123)

	if err != nil {
		t.Error(err)
	}

	if fmt.Sprint(expectedCallBody) != fmt.Sprint(callBody) {
		t.Errorf("Unexpected call body: %v != %v", expectedCallBody, callBody)
	}

}

func TestGetUpdates(t *testing.T) {
	updatesResponse := GetUpdatesResponse{
		Result: []Update{
			{
				UpdateId: 3,
				Message: Message{
					MessageId: 2,
					From:      User{Id: 1, Username: "bromigo"},
					Text:      "hello",
				},
			},
		},
	}

	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			marshalledBody, _ := json.Marshal(updatesResponse)
			_, err := w.Write(marshalledBody)
			if err != nil {
				t.Fatal(err)
			}
		}),
	)
	defer ts.Close()

	telegramClient := NewTelegramClient("fake_token")
	telegramClient.telegramURL = ts.URL

	updates, err := telegramClient.GetUpdates(0)
	if err != nil {
		t.Fatal(err)
	}

	if updates[0] != updatesResponse.Result[0] {
		t.Fatalf("Unexpected response: %v != %v", updates[0], updatesResponse.Result[0])
	}

}
