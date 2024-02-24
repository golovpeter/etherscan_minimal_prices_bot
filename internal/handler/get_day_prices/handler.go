package get_day_prices

import (
	"etherscan_gastracker/internal/service/prices"
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const hoursInDay = 23

type handler struct {
	service  prices.GetPricesService
	location *time.Location
}

func NewHandler(
	service prices.GetPricesService,
	location *time.Location,
) *handler {
	return &handler{
		service:  service,
		location: location,
	}
}

func (h *handler) GetDayPrices(bot *tgbotapi.BotAPI, data *GetDayPricesIn) {
	allPrices, err := h.service.GetAllPrices()
	if err != nil {
		log.Println(err)
		return
	}

	currentTime := time.Now().In(h.location).Format("02-01-2006")
	message := fmt.Sprintf(
		"Minimal Prices per hours on %s\n\n",
		currentTime,
	)

	for i := 0; i <= hoursInDay; i++ {
		if i < 10 {
			message += "0"
		}

		message += fmt.Sprintf("%d: %d\n", i, allPrices[i])
	}

	if _, err = bot.Send(tgbotapi.NewMessage(data.UserID, message)); err != nil {
		log.Println(err)
		return
	}
}
