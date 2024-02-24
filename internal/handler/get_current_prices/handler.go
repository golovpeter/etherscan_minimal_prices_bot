package get_current_prices

import (
	"etherscan_gastracker/internal/service/prices"
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

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

func (h *handler) GetCurrentPrices(bot *tgbotapi.BotAPI, data *GetCurrentPricesIn) {
	curPrices, err := h.service.GetCurrentPrices()
	if err != nil {
		log.Println(err)
		return
	}

	nowTime := time.Now().In(h.location).Format("02-01-2006 15:04:05")
	message := fmt.Sprintf("Current gwei on %s.\nLow: %d, Average: %d, High: %d",
		nowTime,
		curPrices.SafeGasPrice,
		curPrices.ProposeGasPrice,
		curPrices.FastGasPrice)

	if _, err = bot.Send(tgbotapi.NewMessage(data.UserID, message)); err != nil {
		log.Println(err)
		return
	}
}
