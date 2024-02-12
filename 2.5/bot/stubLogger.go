package bot

type stubLogger struct{}

func NewStubLogger() *stubLogger {
	return &stubLogger{}
}

func (s *stubLogger) Info(_...interface{}) {
}

func (s *stubLogger) Warn(_...interface{}) {
}
