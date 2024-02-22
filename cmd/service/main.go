package main

import (
	"etherscan_gastracker/internal/handler/get_day_prices"
	pricesRepo "etherscan_gastracker/internal/repository/prices"
	"etherscan_gastracker/internal/service/prices"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Panicln(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Println(err)
		return
	}

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
				ApiKey: APIKEY,
			})
		}
	}
}
