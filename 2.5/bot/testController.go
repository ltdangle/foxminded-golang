package bot

type testController struct {
	matchedMsg string
}

func NewTestController() *testController {
	return &testController{}
}

func (c *testController) CreateReply(_ Update) *BotMessage {
	return NewBotMessage("test ok")
}

func (c *testController) MatchMsg(update Update) bool {
	if update.Message.Text == nil {
		return false
	}

	return c.matchedMsg == *update.Message.Text
}

func (c *testController) SetMatchMsg(msg string) {
	c.matchedMsg = msg
}
