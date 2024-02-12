package bot

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

// Msg action interface.
type replyCntrlrInterface interface {
	// Messages that controller reacts to.
	MatchMsg(update Update) bool
	// Creates bot reply message.
	CreateReply(update Update) *BotMessage
}

// Message pusher interface.
type pusherCntrlrInterface interface {
	// Pushes message to the user.
	Push(bot botInterface, t time.Time)
}

// Logger interface.
type LoggerInterface interface {
	Info(args ...interface{})
	Warn(args ...interface{})
}

// Bot application.
type bot struct {
	botUrl              string
	log                 LoggerInterface
	replyCntrlrs        []replyCntrlrInterface
	subscriptionCntrlrs []pusherCntrlrInterface
	wg                  sync.WaitGroup
	quitSignal          chan struct{}
	store               StoreInterface
}

func NewBot(botUrl string, log LoggerInterface, store StoreInterface) *bot {
	b := &bot{
		botUrl:     botUrl,
		log:        log,
		quitSignal: make(chan struct{}),
		store:      store,
	}

	return b
}

// Main bot loop.
func (bot *bot) Run() {

	// Listen to interrupts.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	bot.wg.Add(1)
	go bot.reactToMsgs()

	bot.wg.Add(1)
	go bot.pushMsgs()

	<-c

	bot.log.Info("Quitting...")

	// Send quit signal to goroutines.
	close(bot.quitSignal)

	bot.wg.Wait()
}

// Goroutine that reacts to incoming user messages.
func (bot *bot) reactToMsgs() {
	defer bot.wg.Done()

	offset := 0
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			{
				updates, err := bot.getUpdates(offset)
				if err != nil {
					bot.log.Warn("Smth went wrong: ", err.Error())
				}

				for _, update := range updates {
					// Save update to storage.
					err = bot.store.Save(update)
					if err != nil {
						bot.log.Warn("Smth went wrong: ", err.Error())
					}

					msg := bot.router(update)
					msg.ChatId = update.Message.Chat.ID

					msgJson, err := json.Marshal(msg)
					if err != nil {
						bot.log.Warn("Smth went wrong: ", err.Error())
					}

					err = bot.respond(msgJson)
					if err != nil {
						bot.log.Warn("Smth went wrong: ", err.Error())
					}

					offset = update.UpdateID + 1
				}

				bot.log.Info(updates)
			}
		case <-bot.quitSignal:
			{
				bot.log.Info("Cleaning up reactToMsgs()")
				return
			}
		}
	}
}

// Gorutine that pushes messages to registered users.
func (bot *bot) pushMsgs() {
	defer bot.wg.Done()

	ticker := time.NewTicker(time.Second * 1)
	for {
		select {
		case <-ticker.C:
			{
				for _, pusher := range bot.subscriptionCntrlrs {
					go func(pusher pusherCntrlrInterface, t time.Time) {
						pusher.Push(bot, t)
					}(pusher, time.Now())
				}
			}
		case <-bot.quitSignal:
			{
				bot.log.Info("Cleaning up pushMsgs()")
				return
			}
		}
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

	var restResponse TelegramResponse
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		return nil, err
	}
	return restResponse.Result, nil
}

// Respond to incoming messages (updates).
func (bot *bot) respond(msgJson []byte) error {
	resp, err := http.Post(bot.botUrl+"/sendMessage", "application/json", bytes.NewBuffer(msgJson))

	if err != nil {
		bot.log.Warn(err)
		return err
	}

	if resp.StatusCode != 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		bot.log.Warn(string(body))

		return errors.New("Non 200 status code returned from telegram api.")
	}

	return nil
}
func (bot *bot) AddReplyController(c replyCntrlrInterface) {
	bot.replyCntrlrs = append(bot.replyCntrlrs, c)
}

func (bot *bot) AddSubscription(p pusherCntrlrInterface) {
	bot.subscriptionCntrlrs = append(bot.subscriptionCntrlrs, p)
}

func (bot *bot) router(update Update) *BotMessage {
	// Loop over all registered controllers to find first match that reacts to incoming message.
	for _, controller := range bot.replyCntrlrs {
		if controller.MatchMsg(update) {
			return controller.CreateReply(update)
		}
	}
	return &BotMessage{Text: "Command not found..."}
}
