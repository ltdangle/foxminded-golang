package bot

type testController struct {
	matchedMsg string
}

func NewTestController() *testController {
	return &testController{}
}

func (c *testController) CreateReply(_ Message) *BotMessage {
	return &BotMessage{Text: "test ok"}
}

func (c *testController) MatchMsg(msg Message) bool {
	if msg.Text == nil {
		return false
	}

	return c.matchedMsg == *msg.Text
}

func (c *testController) SetMatchMsg(msg string) {
	c.matchedMsg = msg
}
