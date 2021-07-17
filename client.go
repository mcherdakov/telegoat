package telegoat

import (
	"encoding/json"
	"fmt"
)

type TelegramClient struct {
	telegramURL string
	token       string
}

func NewTelegramClient(token string) TelegramClient {
	telegramClient := TelegramClient{
		telegramURL: fmt.Sprintf("%s%s", cfg.telegramHost, token),
		token:       token,
	}

	return telegramClient
}

func (t *TelegramClient) SendMessage(message string, chat_id int) error {
	sendMessageContext := map[string]interface{}{
		"chat_id": chat_id,
		"text":    message,
	}

	_, err := postRequest(
		fmt.Sprintf("%s/sendMessage", t.telegramURL),
		sendMessageContext,
	)

	if err != nil {
		return fmt.Errorf("unable to send message: %v", err.Error())
	}

	return nil
}

func (t *TelegramClient) GetUpdates(offset int) ([]Update, error) {
	updateContext := map[string]interface{}{
		"offset": offset,
	}
	responseData, err := postRequest(
		fmt.Sprintf("%s/getUpdates", t.telegramURL),
		updateContext,
	)
	if err != nil {
		return nil, err
	}

	var updatesResponse GetUpdatesResponse
	err = json.Unmarshal([]byte(responseData), &updatesResponse)
	if err != nil {
		return nil, err
	}

	return updatesResponse.Result, nil
}
