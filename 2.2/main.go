package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
)

type config struct {
	BotApi    string `env:"BOT_API"`
	BotToken  string `env:"BOT_TOKEN"`
	LogOutput string `env:"LOG_OUTPUT"`
	LogFile   string `env:"LOG_FILE"`
	LogLevel  string `env:"LOG_LEVEL"`
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

	// Configure log output.
	switch cfg.LogOutput {
	case "file":
		file, err := os.Create(cfg.LogFile)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		log.SetOutput(file)
	case "stdout":
		log.SetOutput(os.Stdout)
	default:
		log.Fatal("Log output not configured.")
	}

	// Configure log level.
	switch cfg.LogLevel {
	case "panic":
		log.SetLevel(log.PanicLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "trace":
		log.SetLevel(log.TraceLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}

	botUrl := cfg.BotApi + cfg.BotToken

	// Main bot loop.
	offset := 0
	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		updates, err := getUpdates(botUrl, offset)
		if err != nil {
			log.Warn("Smth went wrong: ", err.Error())
		}

		for _, update := range updates {
			err = respond(botUrl, update)
			if err != nil {
				log.Warn("Smth went wrong: ", err.Error())
			}

			offset = update.UpdateId + 1
		}

		log.Info(updates)
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
