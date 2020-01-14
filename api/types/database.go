package types

type RecentBlock struct {
	Height    int   `json:"height"`
	Size      int   `json:"size"`
	Timestamp int64 `json:"timestamp"`
	TxCount   int   `json:"txCount"`
	Weight    int   `json:"weight"`
}

type HistoricalMempoolData struct {
	DataInBuckets interface{} `json:"dataInBuckets"`
	Timestamp     int64       `json:"timestamp"`
}

type TransactionStat struct {
	SegwitCount int   `json:"segwitCount"`
	RbfCount    int   `json:"rbfCount"`
	TxCount     int   `json:"txCount"`
	Timestamp   int64 `json:"timestamp"`
}

type MempoolEntry struct {
	EntryTime      int64          `json:"entryTime"`
	TxID           string         `json:"txid"`
	Fee            int64          `json:"fee"`
	Size           int64          `json:"size"`
	Version        int32          `json:"version"`
	InputCount     int            `json:"inputCount"`
	OutputCount    int            `json:"outputCount"`
	Locktime       uint32         `json:"locktime"`
	OutputSum      int64          `json:"outputValue"`
	SpendsSegWit   bool           `json:"spendsSegWit"`
	SpendsMultisig bool           `json:"spendsMultisig"`
	IsBIP69        bool           `json:"isBIP69"`
	SignalsRBF     bool           `json:"signalsRBF"`
	OPReturnData   string         `json:"opreturnData"`
	OPReturnLength int            `json:"opreturnLength"`
	Multisig       map[string]int `json:"multisigsSpend"`
	Spends         map[string]int `json:"spends"`
	PaysTo         map[string]int `json:"paysTo"`
}

// BlockEntry holds the height, the first-seen timestamp and 
// shortended TXIDs. It's used in the Bitcoin Transaction Monitor
// to mark transactions by block they were included in.
type BlockEntry struct {
	Height     int      `json:"height"`
	Timestamp  int64    `json:"timestamp"`
	ShortTXIDs []string `json:"shortTXIDs"`
}
