package prices

import "sync"

type repository struct {
	DayPrices map[int]int
	CurPrices *CurPrices
	mu        *sync.RWMutex
}

func NewRepository() Repository {
	return &repository{
		DayPrices: make(map[int]int),
		CurPrices: &CurPrices{},
		mu:        &sync.RWMutex{},
	}
}

func (r *repository) UpdateCurrentPrices(data *CurPrices) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.CurPrices.SafeGasPrice = data.SafeGasPrice
	r.CurPrices.ProposeGasPrice = data.ProposeGasPrice
	r.CurPrices.FastGasPrice = data.FastGasPrice

	return nil
}

func (r *repository) GetCurrentPrices() (*CurPrices, error) {
	return r.CurPrices, nil
}

func (r *repository) UpdatePrices(hour int, data *UpdatePricesIn) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if hour == 0 {
		r.DayPrices[23] = data.MinimalPriceInHour
	} else {
		r.DayPrices[hour-1] = data.MinimalPriceInHour
	}

	return nil
}

func (r *repository) GetAllPrices() (map[int]int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.DayPrices, nil
}

func (r *repository) ClearData() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for k := range r.DayPrices {
		delete(r.DayPrices, k)
	}

	return nil
}
