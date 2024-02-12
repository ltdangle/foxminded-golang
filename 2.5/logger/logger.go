package logger

import log "github.com/sirupsen/logrus"

type BotLogger struct{}

func NewBotLogger() *BotLogger {
	return &BotLogger{}
}
func (l *BotLogger) Info(args ...interface{}) {
	log.Info(args...)
}

func (l *BotLogger) Warn(args ...interface{}) {
	log.Warn(args...)
}
