package main

import (
	"encoding/json"
	"os"

	log "github.com/sirupsen/logrus"

	"2.1/bot"
	"2.1/holidayApi"
	"2.1/logger"
	weatherapi "2.1/weatherApi"
	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
)

type config struct {
	BotApi            string `env:"BOT_API"`
	BotToken          string `env:"BOT_TOKEN"`
	LogOutput         string `env:"LOG_OUTPUT"`
	LogFile           string `env:"LOG_FILE"`
	LogLevel          string `env:"LOG_LEVEL"`
	AbstractApiKey    string `env:"ABSTRACT_API_KEY"`
	AbstractApiUrl    string `env:"ABSTRACT_API_URL"`
	CountriesDataFile string `env:"COUNTRIES_DATA_FILE"`
	WeathertApiKey    string `env:"WEATHER_API_KEY"`
	WeatherApiUrl     string `env:"WEATHER_API_URL"`
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

	// Init bot with dependencies.
	botUrl := cfg.BotApi + cfg.BotToken
	hApi := holidayApi.NewAbstractApi(cfg.AbstractApiUrl, cfg.AbstractApiKey)
	botLogger := logger.NewBotLogger()

	b := bot.NewBot(botUrl, botLogger)

	// /testCntrl controller
	testCntrl := bot.NewTestController()
	testCntrl.SetMatchMsg("/test")

	// /startCntrl controller
	startCntrl := bot.NewStartController()
	startCntrl.SetMatchMsg("/start")

	// country controller
	countryCntrl := bot.NewCountryController(hApi, botLogger)

	// read countries data
	countries := readCountriesData(cfg.CountriesDataFile)
	for _, country := range countries {
		startCntrl.AddCountry(country)
		countryCntrl.AddCountry(country)
	}

	// weather controller
	wApi := weatherapi.NewWeatherApi(cfg.WeatherApiUrl, cfg.WeathertApiKey)
	weatherCntrl := bot.NewWeatherController(wApi, botLogger)

	b.AddReplyController(testCntrl)
	b.AddReplyController(startCntrl)
	b.AddReplyController(countryCntrl)
	b.AddReplyController(weatherCntrl)

	b.Run()
}

func readCountriesData(path string) []bot.CountryData {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err.Error())
	}

	countries := []bot.CountryData{}
	err = json.Unmarshal(file, &countries)
	if err != nil {
		log.Fatal(err.Error())
	}

	return countries
}
