package telegoat

import (
	"log"
	"time"
)

func (t *TelegramClient) Poll(pollTime time.Duration, updatesHandler func(Update)) {
	offset := 0
	for {
		updates, err := t.GetUpdates(offset)
		if err != nil {
			log.Println(err)
			continue
		}

		for _, update := range updates {
			go updatesHandler(update)
		}

		if len(updates) != 0 {
			offset = updates[len(updates)-1].UpdateId + 1
		}

		time.Sleep(pollTime)
	}
}
