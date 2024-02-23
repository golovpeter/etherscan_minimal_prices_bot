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

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Panicln(err)
	}

	wh, err := tgbotapi.NewWebhook(os.Getenv("WEBHOOK_URL"))
	if err != nil {
		log.Fatalf("NewWebhook failed: %s", err)

	}

	_, err = bot.Request(wh)
	if err != nil {
		log.Fatalf("SetWebhook failed: %s", err)
	}

	updates := bot.ListenForWebhook("/")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	go func() {
		err = http.ListenAndServe(":"+port, nil)
		if err != nil {
			log.Fatalln("http err:", err)
		}
	}()

	pricesRepository := pricesRepo.NewRepository()

	getPricesService := prices.NewService(pricesRepository)

	errCh := getPricesService.GetNewPrices(os.Getenv("API_KEY"))
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
		case "get_prices":
			getPricesHandler.GetDayPrices(bot, &get_day_prices.GetDayPricesIn{
				UserID: update.Message.Chat.ID,
				ApiKey: os.Getenv("API_KEY"),
			})
		}
	}
}
