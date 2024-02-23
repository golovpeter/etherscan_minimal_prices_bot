package prices

import (
	"encoding/json"
	"etherscan_gastracker/internal/common"
	"etherscan_gastracker/internal/repository/prices"
	"io"
	"math"
	"net/http"
	"time"
)

const (
	getPricesURL = "https://api.etherscan.io/api?module=gastracker&action=gasoracle&apikey="
	tickTime     = time.Second * 15
)

type service struct {
	repository prices.Repository
}

func NewService(repository prices.Repository) GetPricesService {
	return &service{
		repository: repository,
	}
}

func (s *service) GetCurrentPrices() (map[int]int, error) {
	return s.repository.GetAllPrices()
}

func (s *service) GetNewPrices(apiKey string) chan error {
	errCh := make(chan error)

	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		errCh <- err
		return errCh
	}

	curMinimalPrice := math.MaxInt32
	ticker := time.NewTicker(tickTime)
	prevHour := time.Now().Hour()

	go func() {
		for {
			hour, _, _ := time.Now().In(loc).Clock()

			if hour != prevHour && curMinimalPrice != math.MaxInt32 {
				err = s.repository.UpdatePrices(hour, &prices.UpdatePricesIn{
					MinimalPriceInHour: curMinimalPrice,
				})
				if err != nil {
					errCh <- err
					return
				}

				prevHour = hour
				curMinimalPrice = math.MaxInt32
			}

			select {
			case <-ticker.C:
				reqPrices, err := getRequest(apiKey)
				if err != nil {
					errCh <- err
					return
				}

				curPrices, err := common.ConvertPrices(
					reqPrices.SafeGasPrice,
					reqPrices.ProposeGasPrice,
					reqPrices.FastGasPrice,
				)
				if err != nil {
					errCh <- err
					return
				}

				if curPrices.SafeGasPrice < curMinimalPrice {
					curMinimalPrice = curPrices.ProposeGasPrice
				}
			}
		}
	}()

	return errCh
}

func getRequest(apiKey string) (*Result, error) {
	var in GetPricesIn

	req, err := http.Get(getPricesURL + apiKey)
	if err != nil {
		return nil, err
	}

	defer req.Body.Close()

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &in)
	if err != nil {
		return nil, err
	}

	result := in.Result

	return &Result{
		LastBlock:       result.LastBlock,
		SafeGasPrice:    result.SafeGasPrice,
		ProposeGasPrice: result.ProposeGasPrice,
		FastGasPrice:    result.FastGasPrice,
		SuggestBaseFee:  result.SuggestBaseFee,
		GasUsedRatio:    result.GasUsedRatio,
	}, nil
}
