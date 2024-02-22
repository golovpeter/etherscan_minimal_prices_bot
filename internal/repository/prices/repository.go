package prices

import "sync"

type repository struct {
	DayPrices map[int]int
	mu        *sync.RWMutex
}

func NewRepository() Repository {
	return &repository{
		DayPrices: make(map[int]int),
		mu:        &sync.RWMutex{},
	}
}

func (r *repository) UpdatePrices(hour int, data *UpdatePricesIn) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.DayPrices[hour] = data.MinimalPriceInHour

	return nil
}

func (r *repository) GetAllPrices() (map[int]int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.DayPrices, nil
}
