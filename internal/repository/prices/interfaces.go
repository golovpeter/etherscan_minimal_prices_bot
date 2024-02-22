package prices

type Repository interface {
	UpdatePrices(hour int, data *UpdatePricesIn) error
	GetAllPrices() (map[int]int, error)
}
