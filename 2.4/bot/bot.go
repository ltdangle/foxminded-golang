package bot

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

// Msg action interface.
type replyControllerInterface interface {
	// Messages that controller reacts to.
	MatchMsg(msg Message) bool
	// Creates bot reply message.
	CreateReply(msg Message) *BotMessage
}

// Logger interface.
type loggerInterface interface {
	Info(args ...interface{})
	Warn(args ...interface{})
}

// Bot application.
type bot struct {
	botUrl           string
	log              loggerInterface
	replyControllers []replyControllerInterface
}

func NewBot(botUrl string, log loggerInterface) *bot {
	return &bot{
		botUrl: botUrl,
		log:    log,
	}
}

// Main bot loop.
func (bot *bot) Run() {
	offset := 0
	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		updates, err := bot.getUpdates(offset)
		if err != nil {
			log.Warn("Smth went wrong: ", err.Error())
		}

		for _, update := range updates {
			msg := bot.router(update.Message)
			msg.ChatId = update.Message.Chat.ChatId

			msgJson, err := json.Marshal(msg)
			if err != nil {
				log.Warn("Smth went wrong: ", err.Error())
			}

			err = bot.respond(msgJson)
			if err != nil {
				log.Warn("Smth went wrong: ", err.Error())
			}

			offset = update.UpdateId + 1
		}

		log.Info(updates)
	}
}

// Query bot api for updates.
func (bot *bot) getUpdates(offset int) ([]Update, error) {
	resp, err := http.Get(bot.botUrl + "/getUpdates" + "?offset=" + strconv.Itoa(offset))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	bot.log.Info("Bot update: " + string(body))

	var restResponse RestResponse
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		return nil, err
	}
	return restResponse.Result, nil
}

// Respond to incoming messages (updates).
func (bot *bot) respond(msgJson []byte) error {
	_, err := http.Post(bot.botUrl+"/sendMessage", "application/json", bytes.NewBuffer(msgJson))
	if err != nil {
		return err
	}

	return nil
}
func (bot *bot) AddReplyController(c replyControllerInterface) {
	bot.replyControllers = append(bot.replyControllers, c)
}

func (bot *bot) router(msgIn Message) *BotMessage {
	// Loop over all registered controllers to find first match that reacts to incoming message.
	for _, controller := range bot.replyControllers {
		if controller.MatchMsg(msgIn) {
			return controller.CreateReply(msgIn)
		}
	}
	return &BotMessage{Text: "Command not found..."}
}
