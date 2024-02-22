package common

import "strconv"

func ConvertPrices(
	safePrice string,
	proposePrice string,
	fastPrice string,
) (*Prices, error) {
	safePriceInt, err := strconv.Atoi(safePrice)
	if err != nil {
		return nil, err
	}

	proposePriceInt, err := strconv.Atoi(proposePrice)
	if err != nil {
		return nil, err
	}

	fastPriceInt, err := strconv.Atoi(fastPrice)
	if err != nil {
		return nil, err
	}

	return &Prices{
		SafeGasPrice:    safePriceInt,
		ProposeGasPrice: proposePriceInt,
		FastGasPrice:    fastPriceInt,
	}, nil
}
