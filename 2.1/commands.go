package main

// Command name.
type commandName string

// Command action.
type commandFn func(*BotMessage)

// Command map.
var commands = map[commandName]commandFn{
	"/about": func(m *BotMessage) {
		m.Text = "Some short information about me."
	},
	"/links": func(m *BotMessage) {
		m.Text = "List of my social links."
	},
	"/start": func(m *BotMessage) {
		m.Text = "List of commands with reply markup."
		m.ReplyMarkup = &ReplyKeyboardMarkup{
			Keyboard: [][]string{
				{"/start", "/help"},
				{"/about", "/links"},
			},
			ResizeKeyboard:  true,
			OneTimeKeyboard: true,
		}
	},
	"/help": func(m *BotMessage) {
		m.Text = "List of commands with reply markup."
		m.ReplyMarkup = &ReplyKeyboardMarkup{
			Keyboard: [][]string{
				{"/start", "/help"},
				{"/about", "/links"},
			},
			ResizeKeyboard:  true,
			OneTimeKeyboard: true,
		}
	},
}
