package infra

type LoggerInterface interface {
	Warn(args ...interface{})
	Info(args ...interface{})
}

type stubLogger struct{}

func NewStubLogger() *stubLogger {
	return &stubLogger{}
}
func (l *stubLogger) Warn(_...interface{}) {
	// Implement warning log logic
}

func (l *stubLogger) Info(_...interface{}) {
	// Implement info log logic
}
