package models

type SignRequest struct {
	CallData           string `json:"call_data"`
	AccountGasLimits   string `json:"account_gas_limits"`
	PreVerificationGas string `json:"pre_verification_gas"`
	GasFees            string `json:"gas_fees"`
}

type SignResponse struct {
	TxHash string `json:"tx_hash"`
}
