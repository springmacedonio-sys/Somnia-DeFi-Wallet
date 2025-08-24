package types

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type TxRecord struct {
	DexName    string `json:"dexName"`
	FromToken  string `json:"fromToken"`
	FromAmount string `json:"fromAmount"`
	ToToken    string `json:"toToken"`
	ToAmount   string `json:"toAmount"`
	Timestamp  string `json:"timestamp"`
}
