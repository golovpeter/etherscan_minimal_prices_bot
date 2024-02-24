package prices

type UpdatePricesIn struct {
	MinimalPriceInHour int
}

type CurPrices struct {
	SafeGasPrice    int
	ProposeGasPrice int
	FastGasPrice    int
}
