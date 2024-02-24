package main

import (
	customconfig "etherscan_gastracker/internal/config"
	"etherscan_gastracker/internal/handler/get_current_prices"
	"etherscan_gastracker/internal/handler/get_day_prices"
	pricesRepo "etherscan_gastracker/internal/repository/prices"
	"etherscan_gastracker/internal/service/prices"
	"log"
	"net/http"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var numericKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Minimal hours gwei"),
		tgbotapi.NewKeyboardButton("Current gwei"),
	),
)

func main() {
	config := customconfig.NewConfig()

	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		log.Fatalln(err)
	}

	bot, err := tgbotapi.NewBotAPI(config.BotToken)
	if err != nil {
		log.Fatalln(err)
	}

	wh, err := tgbotapi.NewWebhook(config.WebHookURL)
	if err != nil {
		log.Fatalf("NewWebhook failed: %s", err)

	}

	_, err = bot.Request(wh)
	if err != nil {
		log.Fatalf("SetWebhook failed: %s", err)
	}

	updates := bot.ListenForWebhook("/")

	http.HandleFunc("/health", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte("OK"))
	})

	pricesRepository := pricesRepo.NewRepository()

	getPricesService := prices.NewService(pricesRepository, location)

	errCh := getPricesService.GetNewPrices(config.ApiToken)
	go func() {
		<-errCh
		log.Println(err)
		return
	}()

	getPricesHandler := get_day_prices.NewHandler(getPricesService, location)
	currentPricesHandler := get_current_prices.NewHandler(getPricesService, location)

	go func() {
		err = http.ListenAndServe(":"+config.Port, nil)
		if err != nil {
			log.Fatalln("http err:", err)
		}
	}()

	for update := range updates {
		if update.Message == nil {
			continue
		}

		command := update.Message.Text
		switch command {
		case "/start":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Команда /get_prices позволит тебе получить минимальные цена на газ!")
			msg.ReplyMarkup = numericKeyboard

			if _, err = bot.Send(msg); err != nil {
				log.Fatalln(err)
			}
		case "Minimal hours gwei":
			getPricesHandler.GetDayPrices(bot, &get_day_prices.GetDayPricesIn{
				UserID: update.Message.Chat.ID,
				ApiKey: config.ApiToken,
			})
		case "Current gwei":
			currentPricesHandler.GetCurrentPrices(bot, &get_current_prices.GetCurrentPricesIn{
				UserID: update.Message.Chat.ID,
			})
		default:
			if _, err = bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Такого я не понимаю!")); err != nil {
				log.Fatalln(err)
			}
		}
	}
}
