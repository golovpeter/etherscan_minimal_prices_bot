package main

import (
	"etherscan_gastracker/internal/handler/get_day_prices"
	pricesRepo "etherscan_gastracker/internal/repository/prices"
	"etherscan_gastracker/internal/service/prices"
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var numericKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("/get_prices"),
	),
)

var (
	BotToken   = os.Getenv("BOT_TOKEN")
	WebhookUrl = os.Getenv("WEBHOOK_URL")
	ApiKey     = os.Getenv("API_KEY")
	Port       = os.Getenv("PORT")
)

func main() {
	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		log.Panicln(err)
	}

	wh, err := tgbotapi.NewWebhook(WebhookUrl)
	if err != nil {
		log.Fatalf("NewWebhook failed: %s", err)

	}

	_, err = bot.Request(wh)
	if err != nil {
		log.Fatalf("SetWebhook failed: %s", err)
	}

	updates := bot.ListenForWebhook("/")

	port := Port
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/health", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte("OK"))
	})

	go func() {
		err = http.ListenAndServe(":"+port, nil)
		if err != nil {
			log.Fatalln("http err:", err)
		}
	}()

	pricesRepository := pricesRepo.NewRepository()

	getPricesService := prices.NewService(pricesRepository)

	errCh := getPricesService.GetNewPrices(ApiKey)
	go func() {
		<-errCh
		log.Println(err)
		return
	}()

	getPricesHandler := get_day_prices.NewHandler(getPricesService)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		command := update.Message.Command()

		switch command {
		case "start":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Команда /get_prices позволит тебе получить минимальные цена на газ!")
			msg.ReplyMarkup = numericKeyboard

			if _, err = bot.Send(msg); err != nil {
				log.Fatalln(err)
			}
		case "get_prices":
			getPricesHandler.GetDayPrices(bot, &get_day_prices.GetDayPricesIn{
				UserID: update.Message.Chat.ID,
				ApiKey: ApiKey,
			})
		default:
			if _, err = bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Такого я не понимаю!")); err != nil {
				log.Fatalln(err)
			}
		}
	}
}
