package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
)

type config struct {
	BotApi   string `env:"BOT_API"`
	BotToken string `env:"BOT_TOKEN"`
}

func main() {
	// Load .env config.
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err.Error())
	}

	botUrl := cfg.BotApi + cfg.BotToken

	// Main bot loop.
	offset := 0
	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		updates, err := getUpdates(botUrl, offset)
		if err != nil {
			log.Println("Smth went wrong: ", err.Error())
		}

		for _, update := range updates {
			err = respond(botUrl, update)
			if err != nil {
				log.Println("Smth went wrong: ", err.Error())
			}

			offset = update.UpdateId + 1
		}
		fmt.Println(updates)
	}
}

// Query bot api for updates.
func getUpdates(botUrl string, offset int) ([]Update, error) {
	resp, err := http.Get(botUrl + "/getUpdates" + "?offset=" + strconv.Itoa(offset))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var restResponse RestResponse
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		return nil, err
	}
	return restResponse.Result, nil
}

// Respond to incoming messages (updates).
func respond(botUrl string, update Update) error {
	botMessage := &BotMessage{
		ChatId: update.Message.Chat.ChatId,
	}

	if command, ok := commands[commandName(update.Message.Text)]; ok {
		command(botMessage)
	} else {
		botMessage.Text = "command not found"
	}

	buf, err := json.Marshal(botMessage)
	if err != nil {
		return err
	}
	_, err = http.Post(botUrl+"/sendMessage", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return err
	}

	return nil
}
