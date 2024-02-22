package prices

type GetPricesService interface {
	GetNewPrices(apiKey string) chan error
	GetCurrentPrices() (map[int]int, error)
}
