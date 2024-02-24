package prices

import (
	"encoding/json"
	"etherscan_gastracker/internal/common"
	"etherscan_gastracker/internal/repository/prices"
	"fmt"
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
	location   *time.Location
	repository prices.Repository
}

func NewService(
	repository prices.Repository,
	location *time.Location,
) GetPricesService {
	return &service{
		repository: repository,
		location:   location,
	}
}

func (s *service) GetAllPrices() (map[int]int, error) {
	return s.repository.GetAllPrices()
}

func (s *service) GetCurrentPrices() (*CurPrices, error) {
	curPrices, err := s.repository.GetCurrentPrices()
	if err != nil {
		return nil, err
	}

	return &CurPrices{
		SafeGasPrice:    curPrices.SafeGasPrice,
		ProposeGasPrice: curPrices.ProposeGasPrice,
		FastGasPrice:    curPrices.FastGasPrice,
	}, nil
}

func (s *service) GetNewPrices(apiKey string) chan error {
	cleared := false
	errCh := make(chan error)

	curMinimalPrice := math.MaxInt32
	ticker := time.NewTicker(tickTime)
	prevHour := time.Now().In(s.location).Hour()

	go func() {
		for {
			hour, _, _ := time.Now().In(s.location).Clock()

			if hour == 1 && cleared == false {
				err := s.repository.ClearData()
				if err != nil {
					errCh <- err
					return
				}

				fmt.Println("cleared")
				cleared = true
			}

			if hour != prevHour {
				err := s.repository.UpdatePrices(hour, &prices.UpdatePricesIn{
					MinimalPriceInHour: curMinimalPrice,
				})
				if err != nil {
					errCh <- err
					return
				}

				prevHour = hour
				curMinimalPrice = math.MaxInt32

				if cleared {
					cleared = false
				}
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

				err = s.repository.UpdateCurrentPrices(&prices.CurPrices{
					SafeGasPrice:    curPrices.SafeGasPrice,
					ProposeGasPrice: curPrices.ProposeGasPrice,
					FastGasPrice:    curPrices.FastGasPrice,
				})
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
		SafeGasPrice:    result.SafeGasPrice,
		ProposeGasPrice: result.ProposeGasPrice,
		FastGasPrice:    result.FastGasPrice,
	}, nil
}
