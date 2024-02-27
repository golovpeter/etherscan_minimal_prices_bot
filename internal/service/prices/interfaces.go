package prices

type GetPricesService interface {
	GetNewPrices(apiKey string) chan string
	GetAllPrices() (map[int]int, error)
	GetCurrentPrices() (*CurPrices, error)
}
