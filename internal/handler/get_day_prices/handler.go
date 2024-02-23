package get_day_prices

import (
	"etherscan_gastracker/internal/service/prices"
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const hoursInDay = 24

type handler struct {
	service prices.GetPricesService
}

func NewHandler(service prices.GetPricesService) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) GetDayPrices(bot *tgbotapi.BotAPI, data *GetDayPricesIn) {
	currentPrices, err := h.service.GetCurrentPrices()
	if err != nil {
		log.Println(err)
		return
	}

	currentTime := time.Now().Format("01-02-2006")
	message := fmt.Sprintf("Minimal Price per hours on %s\n", currentTime)

	for i := 0; i <= hoursInDay; i++ {
		message += fmt.Sprintf("%d: %d\n", i, currentPrices[i])
	}

	_, err = bot.Send(tgbotapi.NewMessage(data.UserID, message))
	if err != nil {
		log.Println(err)
		return
	}
}
