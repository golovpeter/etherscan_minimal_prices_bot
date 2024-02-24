package prices

type GetPricesService interface {
	GetNewPrices(apiKey string) chan error
	GetAllPrices() (map[int]int, error)
	GetCurrentPrices() (*CurPrices, error)
}
