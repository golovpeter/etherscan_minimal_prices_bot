package prices

type GetPricesIn struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  Result `json:"result"`
}

type Result struct {
	SafeGasPrice    string `json:"SafeGasPrice"`
	ProposeGasPrice string `json:"ProposeGasPrice"`
	FastGasPrice    string `json:"FastGasPrice"`
}

type CurPrices struct {
	SafeGasPrice    int
	ProposeGasPrice int
	FastGasPrice    int
}
