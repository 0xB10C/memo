package types

// PartialMempoolEntry is slimmed down version of Bitcoin Core's mempool entry
// returned by the getrawmempool RPC. This allows us to be more memory efficient.
// See https://github.com/bitcoin/bitcoin/blob/6d5790956f45e3de5c6c4ee6fda21878b0d1287b/src/rpc/mempool.cpp#L253-L279
type PartialMempoolEntry struct {
	Size int     `json:"vsize"`
	Fee  float64 `json:"fees.base"`
	Fees struct {
		Base float64 `json:"base"`
	} `json:"fees"`
	Time              int    `json:"time"`
	Wtxid             string `json:"wtxid"`
	Bip125Replaceable bool   `json:"bip125-replaceable"`
}
