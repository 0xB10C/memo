package types

type Fees struct {
	Base       float64 `json:"base"`
	Modified   float64 `json:"modified"`
	Ancestor   float64 `json:"ancestor"`
	Descendant float64 `json:"descendant"`
}

type Transaction struct {
	Fees              Fees          `son:"fees"`
	Size              int           `json:"size"`
	Fee               float64       `json:"fee"`
	Modifiedfee       float64       `json:"modifiedfee"`
	Time              int           `json:"time"`
	Height            int           `json:"height"`
	Descendantcount   int           `json:"descendantcount"`
	Descendantsize    int           `json:"descendantsize"`
	Descendantfees    int           `json:"descendantfees"`
	Ancestorcount     int           `json:"ancestorcount"`
	Ancestorsize      int           `json:"ancestorsize"`
	Ancestorfees      int           `json:"ancestorfees"`
	Wtxid             string        `json:"wtxid"`
	Depends           []interface{} `json:"depends"`
	Spentby           []interface{} `json:"spentby"`
	Bip125Replaceable bool          `json:"bip125-replaceable"`
}

// PartialTransaction is a part-struct of `Transaction` which contains
// only the for-now-used values to be more memory efficient
type PartialTransaction struct {
	Size int     `json:"size"`
	Fee  float64 `json:"fee"`
}
