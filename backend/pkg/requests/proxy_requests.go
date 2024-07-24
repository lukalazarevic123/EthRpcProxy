package requests

type SendTransactionArgs struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Gas      string `json:"gas"`
	GasPrice string `json:"gasPrice"`
	Value    string `json:"value"`
	Data     string `json:"data"`
	Nonce    string `json:"nonce"`
}

type SendTransactionRequest struct {
	JsonRpc string                `json:"jsonrpc"`
	Method  string                `json:"method"`
	Params  []SendTransactionArgs `json:"params"`
	ID      int                   `json:"id"`
}

type SendTransactionResponse struct {
	Id      int    `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Result  string `json:"result"`
}
