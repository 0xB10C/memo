package types

// The external fee APIs return different estimate counts for which different structs are used
// FeeAPIResponse1 for one estimate, FeeAPIResponse2 for two and FeeAPIResponse3 for three estimates.

type FeeAPIResponse1 struct {
	HighFee float64 `json:"high"`
}

type FeeAPIResponse2 struct {
	HighFee float64 `json:"high"`
	MedFee  float64 `json:"med"`
}

type FeeAPIResponse3 struct {
	HighFee float64 `json:"high"`
	MedFee  float64 `json:"med"`
	LowFee  float64 `json:"low"`
}

type FeeRateAPIEntry struct {
	Timestamp          int64           `json:"timestamp"`
	BTCCom             FeeAPIResponse1 `json:"btccom"`
	BlockchairCom      FeeAPIResponse1 `json:"blockchaircom"`
	BlockchainInfo     FeeAPIResponse2 `json:"blockchaininfo"`
	EarnCom            FeeAPIResponse3 `json:"earncom"`
	BitgoCom           FeeAPIResponse3 `json:"bitgocom"`
	BlockcypherCom     FeeAPIResponse3 `json:"blockcyphercom"`
	BitpayCom          FeeAPIResponse3 `json:"bitpaycom"`
	WasabiWalletIoEcon FeeAPIResponse3 `json:"wasabiwalletioEcon"`
	WasabiWalletIoCons FeeAPIResponse3 `json:"wasabiwalletioCons"`
	TrezorIo           FeeAPIResponse3 `json:"trezorio"`
	LedgerCom          FeeAPIResponse3 `json:"ledgercom"`
	MyceliumIo         FeeAPIResponse3 `json:"myceliumio"`
	BitcoinerLive      FeeAPIResponse3 `json:"bitcoinerlive"`
	BlockstreamInfo    FeeAPIResponse3 `json:"blockstreaminfo"`
}
