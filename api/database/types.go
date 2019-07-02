package database

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
