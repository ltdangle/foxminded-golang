package bot

type startController struct {
	// Messages that this controller reacts to.
	matchedMsg string
	countries  []CountryData
}

func NewStartController() *startController {
	return &startController{}
}

func (c *startController) CreateReply(_ Message) *BotMessage {
	botMessage := &BotMessage{}
	var kb [][]string
	botMessage.Text = "Pick a country..."
	kb = append(kb, c.countryNames())
	botMessage.ReplyMarkup = &ReplyKeyboardMarkup{
		Keyboard:        kb,
		ResizeKeyboard:  true,
		OneTimeKeyboard: true,
	}

	return botMessage
}

func (c *startController) countryNames() []string {
	var names []string
	for _, country := range c.countries {
		names = append(names, country.Emoji)
	}
	return names
}

func (c *startController) AddCountry(data CountryData) {
	c.countries = append(c.countries, data)
}

func (c *startController) MatchMsg(msg Message) bool {
	if msg.Text == nil {
		return false
	}

	return c.matchedMsg == *msg.Text
}

func (c *startController) SetMatchMsg(msg string) {
	c.matchedMsg = msg
}
