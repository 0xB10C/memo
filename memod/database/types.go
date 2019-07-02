package database

type recentBlock struct {
	Height    int   `json:"height"`
	Size      int   `json:"size"`
	Timestamp int64 `json:"timestamp"`
	TxCount   int   `json:"txCount"`
	Weight    int   `json:"weight"`
}

type historicalMempoolData struct {
	DataInBuckets interface{} `json:"dataInBuckets"`
	Timestamp     int64       `json:"timestamp"`
}

type transactionStat struct {
	SegwitCount int   `json:"segwitCount"`
	RbfCount    int   `json:"rbfCount"`
	TxCount     int   `json:"txCount"`
	Timestamp   int64 `json:"timestamp"`
}
