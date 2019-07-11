package fetcher

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/0xb10c/memo/memod/config"
	"github.com/0xb10c/memo/memod/database"
	"github.com/0xb10c/memo/memod/logger"
	"github.com/0xb10c/memo/memod/types"
	"github.com/jasonlvhit/gocron"
	"github.com/tidwall/gjson"
)

// SetupFeerateAPIFetcher sets up a periodic feerate API fetch job
func SetupFeerateAPIFetcher() {
	feerateFetchInterval := config.GetInt("feeratefetcher.fetchInterval")
	s := gocron.NewScheduler()
	s.Every(uint64(feerateFetchInterval)).Seconds().Do(getFeerates)
	logger.Info.Println("Setup feerate API fetcher. First fetch in", feerateFetchInterval, "seconds")
	<-s.Start()
	defer s.Clear()
}

func getFeerates() {
	cBTCCom := make(chan types.FeeAPIResponse1)
	cBlockchairCom := make(chan types.FeeAPIResponse1)

	cBlockchainInfo := make(chan types.FeeAPIResponse2)

	cEarnCom := make(chan types.FeeAPIResponse3)
	cBitgoCom := make(chan types.FeeAPIResponse3)
	cBlockcypherCom := make(chan types.FeeAPIResponse3)
	cBitpayCom := make(chan types.FeeAPIResponse3)
	cWasabiWalletIoEcon := make(chan types.FeeAPIResponse3)
	cWasabiWalletIoCons := make(chan types.FeeAPIResponse3)
	cTrezorIo := make(chan types.FeeAPIResponse3)
	cLedgerCom := make(chan types.FeeAPIResponse3)
	cMyceliumIo := make(chan types.FeeAPIResponse3)
	cBitcoinerLive := make(chan types.FeeAPIResponse3)
	cBlockstreamInfo := make(chan types.FeeAPIResponse3)

	go getBTCCom(cBTCCom)
	go getBlockchairCom(cBlockchairCom)
	go getBlockchainInfo(cBlockchainInfo)
	go getEarnCom(cEarnCom)
	go getBitgoCom(cBitgoCom)
	go getBlockcypherCom(cBlockcypherCom)
	go getBitpayCom(cBitpayCom)
	go getWasabiWalletIo(cWasabiWalletIoEcon, cWasabiWalletIoCons)
	go getTrezorIo(cTrezorIo)
	go getLedgerCom(cLedgerCom)
	go getMyceliumIo(cMyceliumIo)
	go getBitcoinerLive(cBitcoinerLive)
	go getBlockstreamInfo(cBlockstreamInfo)

	respBTCCom := <-cBTCCom
	respBlockchairCom := <-cBlockchairCom
	respBlockchainInfo := <-cBlockchainInfo
	respEarnCom := <-cEarnCom
	respBitgoCom := <-cBitgoCom
	respBlockcypherCom := <-cBlockcypherCom
	respBitpayCom := <-cBitpayCom
	respWasabiWalletIoEcon := <-cWasabiWalletIoEcon
	respWasabiWalletIoCons := <-cWasabiWalletIoCons
	respTrezorIo := <-cTrezorIo
	respLedgerCom := <-cLedgerCom
	respMyceliumIo := <-cMyceliumIo
	respBitcoinerLive := <-cBitcoinerLive
	respBlockstreamInfo := <-cBlockstreamInfo

	entry := types.FeeRateAPIEntry{
		Timestamp:          time.Now().Unix(),
		BTCCom:             respBTCCom,
		BlockchairCom:      respBlockchairCom,
		BlockchainInfo:     respBlockchainInfo,
		EarnCom:            respEarnCom,
		BitgoCom:           respBitgoCom,
		BlockcypherCom:     respBlockcypherCom,
		BitpayCom:          respBitpayCom,
		WasabiWalletIoEcon: respWasabiWalletIoEcon,
		WasabiWalletIoCons: respWasabiWalletIoCons,
		TrezorIo:           respTrezorIo,
		LedgerCom:          respLedgerCom,
		MyceliumIo:         respMyceliumIo,
		BitcoinerLive:      respBitcoinerLive,
		BlockstreamInfo:    respBlockstreamInfo,
	}

	database.WriteFeerateAPIEntry(entry)
}

func getEarnCom(cEarnCom chan types.FeeAPIResponse3) {
	url := "https://bitcoinfees.earn.com/api/v1/fees/recommended" // response:  {"fastestFee":50,"halfHourFee":50,"hourFee":42}
	body, err := makeHTTPGETReq(url, 5)
	if err != nil {
		logger.Error.Printf("Could not fetch bitcoinfees.earn.com: %v", err.Error())
	}
	result := gjson.GetMany(string(body), "fastestFee", "halfHourFee", "hourFee")
	highFee, medFee, lowFee := result[0].Float(), result[1].Float(), result[2].Float()

	cEarnCom <- types.FeeAPIResponse3{highFee, medFee, lowFee}
}

func getBitgoCom(cBitgoCom chan types.FeeAPIResponse3) {
	url := "https://www.bitgo.com/api/v2/btc/tx/fee"
	body, err := makeHTTPGETReq(url, 5)
	if err != nil {
		logger.Error.Printf("Could not fetch bitgo.com: %v", err.Error())
		cBitgoCom <- types.FeeAPIResponse3{0, 0, 0}
		return
	}
	result := gjson.GetMany(string(body), "feeByBlockTarget.1", "feeByBlockTarget.3", "feeByBlockTarget.6")
	highFee, medFee, lowFee := result[0].Float()/1000, result[1].Float()/1000, result[2].Float()/1000

	cBitgoCom <- types.FeeAPIResponse3{highFee, medFee, lowFee}
}

func getBTCCom(cBTCCom chan types.FeeAPIResponse1) {
	url := "https://btc.com/service/fees/distribution"
	body, err := makeHTTPGETReq(url, 5)
	if err != nil {
		logger.Error.Printf("Could not fetch btc.com: %v", err.Error())
		cBTCCom <- types.FeeAPIResponse1{0}
		return
	}
	result := gjson.Get(string(body), "fees_recommended.one_block_fee")
	highFee := result.Float()

	cBTCCom <- types.FeeAPIResponse1{highFee}
}

func getBlockcypherCom(cBlockcypherCom chan types.FeeAPIResponse3) {
	url := "https://api.blockcypher.com/v1/btc/main"
	body, err := makeHTTPGETReq(url, 5)
	if err != nil {
		logger.Error.Printf("Could not fetch api.blockcypher.com: %v", err.Error())
		cBlockcypherCom <- types.FeeAPIResponse3{0, 0, 0}
		return
	}
	result := gjson.GetMany(string(body), "high_fee_per_kb", "medium_fee_per_kb", "low_fee_per_kb")
	highFee, medFee, lowFee := result[0].Float()/1000, result[1].Float()/1000, result[2].Float()/1000

	cBlockcypherCom <- types.FeeAPIResponse3{highFee, medFee, lowFee}
}

func getBlockchainInfo(cBlockchainInfo chan types.FeeAPIResponse2) {
	url := "https://api.blockchain.info/mempool/fees"
	body, err := makeHTTPGETReq(url, 5)
	if err != nil {
		logger.Error.Printf("Could not fetch api.blockchain.info: %v", err.Error())
		cBlockchainInfo <- types.FeeAPIResponse2{0, 0}
		return
	}
	result := gjson.GetMany(string(body), "priority", "regular")
	highFee, medFee := result[0].Float(), result[1].Float()

	cBlockchainInfo <- types.FeeAPIResponse2{highFee, medFee}
}

func getBlockchairCom(cBlockchairCom chan types.FeeAPIResponse1) {
	url := "https://api.blockchair.com/bitcoin/stats"
	body, err := makeHTTPGETReq(url, 5)
	if err != nil {
		logger.Error.Printf("Could not fetch api.blockchair.com: %v", err.Error())
		cBlockchairCom <- types.FeeAPIResponse1{0}
	}
	result := gjson.Get(string(body), "data.suggested_transaction_fee_per_byte_sat")
	highFee := result.Float()

	cBlockchairCom <- types.FeeAPIResponse1{highFee}
}

func getBitpayCom(cBitpayCom chan types.FeeAPIResponse3) {
	url := "https://insight.bitpay.com/api/utils/estimatefee?nbBlocks=2,3,6"
	body, err := makeHTTPGETReq(url, 5)
	if err != nil {
		logger.Error.Printf("Could not fetch insight.bitpay.com: %v", err.Error())
		cBitpayCom <- types.FeeAPIResponse3{0, 0, 0}
		return
	}
	result := gjson.GetMany(string(body), "2", "3", "6")
	highFee, medFee, lowFee := result[0].Float()*100000, result[1].Float()*100000, result[2].Float()*100000
	cBitpayCom <- types.FeeAPIResponse3{highFee, medFee, lowFee}
}

func getWasabiWalletIo(cWasabiWalletIoEcon chan types.FeeAPIResponse3, cWasabiWalletIoCons chan types.FeeAPIResponse3) {
	url := "https://wasabiwallet.io/api/v3/btc/Blockchain/fees/2,4,6"
	body, err := makeHTTPGETReq(url, 5)
	if err != nil {
		logger.Error.Printf("Could not fetch wasabiwallet.io: %v", err.Error())
		cWasabiWalletIoEcon <- types.FeeAPIResponse3{0, 0, 0}
		cWasabiWalletIoCons <- types.FeeAPIResponse3{0, 0, 0}
		return
	}
	resultEconomical := gjson.GetMany(string(body), "2.economical", "4.economical", "6.economical")
	resultConservative := gjson.GetMany(string(body), "2.conservative", "4.conservative", "6.conservative")
	highFeeEconomical, medFeeEconomical, lowFeeEconomical := resultEconomical[0].Float(), resultEconomical[1].Float(), resultEconomical[2].Float()
	highFeeConservative, medFeeConservative, lowFeeConservative := resultConservative[0].Float(), resultConservative[1].Float(), resultConservative[2].Float()

	cWasabiWalletIoEcon <- types.FeeAPIResponse3{highFeeEconomical, medFeeEconomical, lowFeeEconomical}
	cWasabiWalletIoCons <- types.FeeAPIResponse3{highFeeConservative, medFeeConservative, lowFeeConservative}
}

func getTrezorIo(cTrezorIo chan types.FeeAPIResponse3) {

	// Since the Trezor API is publicly accessible, but not publicly advertised I don't
	// want to have the plain text url crawable on GitHub. It's encoded as base64.
	urlPrefixBytes, _ := base64.StdEncoding.DecodeString("aHR0cHM6Ly9idGMxLnRyZXpvci5pby9hcGkvdjEvZXN0aW1hdGVmZWUv")
	urlPrefix := string(urlPrefixBytes)

	bodyBlocks2, err := makeHTTPGETReq(urlPrefix+"2", 5)
	if err != nil {
		logger.Error.Printf("Could not fetch trezor.io: %v", err.Error())
		cTrezorIo <- types.FeeAPIResponse3{0, 0, 0}
		return
	}

	bodyBlocks4, err := makeHTTPGETReq(urlPrefix+"4", 5)
	if err != nil {
		logger.Error.Printf("Could not fetch trezor.io: %v", err.Error())
		cTrezorIo <- types.FeeAPIResponse3{0, 0, 0}
		return
	}

	bodyBlocks6, err := makeHTTPGETReq(urlPrefix+"6", 5)
	if err != nil {
		logger.Error.Printf("Could not fetch trezor.io: %v", err.Error())
		cTrezorIo <- types.FeeAPIResponse3{0, 0, 0}
		return
	}

	resultBlocks2 := gjson.Get(string(bodyBlocks2), "result")
	resultBlocks4 := gjson.Get(string(bodyBlocks4), "result")
	resultBlocks6 := gjson.Get(string(bodyBlocks6), "result")

	highFee := resultBlocks2.Float() * 100000
	medFee := resultBlocks4.Float() * 100000
	lowFee := resultBlocks6.Float() * 100000

	cTrezorIo <- types.FeeAPIResponse3{highFee, medFee, lowFee}
}

func getMyceliumIo(cMyceliumIo chan types.FeeAPIResponse3) {

	// Since the Mycelium API is publicly accessible, but not publicly advertised I don't
	// want to have the plain text url crawable on GitHub. It's encoded as base64.
	urlBytes, _ := base64.StdEncoding.DecodeString("aHR0cHM6Ly9td3MyMC5teWNlbGl1bS5jb20vd2FwaS93YXBpL2dldE1pbmVyRmVlRXN0aW1hdGlvbnM=")
	url := string(urlBytes)

	// Mycelium uses a self signed cert. It's appended it to a copy of the cert store and trust it to query their API.
	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}

	// Append our cert to a copy of our system cert pool / trust store
	if ok := rootCAs.AppendCertsFromPEM(getMyceliumFeeAPICert()); !ok {
		logger.Error.Printf("Could not append Mycelium self-singed cert to trust store.")
		cMyceliumIo <- types.FeeAPIResponse3{0, 0, 0}
		return
	}

	// Trust the augmented cert pool in our client
	config := &tls.Config{
		RootCAs:            rootCAs,
		InsecureSkipVerify: true, // The https certificate mycelium uses expired on 13.08.2019.
	}
	tr := &http.Transport{TLSClientConfig: config}
	client := &http.Client{Transport: tr}

	resp, err := client.Post(url, "application/json", nil)
	if err != nil {
		logger.Error.Printf("Could not POST to mycelium.com: %v", err.Error())
		cMyceliumIo <- types.FeeAPIResponse3{0, 0, 0}
		return
	}
	defer resp.Body.Close()

	body, err := readResponseBody(resp)
	if err != nil {
		logger.Error.Printf("Could not read response body from mycelium.com: %v", err.Error())
		cMyceliumIo <- types.FeeAPIResponse3{0, 0, 0}
		return
	}

	result := gjson.GetMany(string(body), "r.feeEstimation.feeForNBlocks.2", "r.feeEstimation.feeForNBlocks.4", "r.feeEstimation.feeForNBlocks.10")
	highFee, medFee, lowFee := result[0].Float()/1000, result[1].Float()/1000, result[2].Float()/1000

	cMyceliumIo <- types.FeeAPIResponse3{highFee, medFee, lowFee}
}

func getBitcoinerLive(cBitcoinerLive chan types.FeeAPIResponse3) {
	url := "https://bitcoiner.live/api/fees/estimates/latest"
	body, err := makeHTTPGETReq(url, 5)
	if err != nil {
		logger.Error.Printf("Could not fetch bitcoiner.live: %v", err.Error())
		cBitcoinerLive <- types.FeeAPIResponse3{0, 0, 0}
		return
	}
	result := gjson.GetMany(string(body), "estimates.30.sat_per_vbyte", "estimates.60.sat_per_vbyte", "estimates.120.sat_per_vbyte")
	highFee, medFee, lowFee := result[0].Float(), result[1].Float(), result[2].Float()
	cBitcoinerLive <- types.FeeAPIResponse3{highFee, medFee, lowFee}
}

func getBlockstreamInfo(cBlockstreamInfo chan types.FeeAPIResponse3) {
	url := "https://blockstream.info/api/fee-estimates"
	body, err := makeHTTPGETReq(url, 5)
	if err != nil {
		logger.Error.Printf("Could not fetch blockstream.info: %v", err.Error())
		cBlockstreamInfo <- types.FeeAPIResponse3{0, 0, 0}
		return
	}
	result := gjson.GetMany(string(body), "2", "3", "6")
	highFee, medFee, lowFee := result[0].Float(), result[1].Float(), result[2].Float()
	cBlockstreamInfo <- types.FeeAPIResponse3{highFee, medFee, lowFee}
}

func getLedgerCom(cLedgerCom chan types.FeeAPIResponse3) {
	// Since the Ledger Live API is publicly accessible, but not publicly advertised I don't
	// want to have the plain text url crawable on GitHub. It's encoded as base64.
	urlBytes, _ := base64.StdEncoding.DecodeString("aHR0cHM6Ly9leHBsb3JlcnMuYXBpLmxpdmUubGVkZ2VyLmNvbS9ibG9ja2NoYWluL3YyL2J0Yy9mZWVz")
	url := string(urlBytes)
	body, err := makeHTTPGETReq(url, 5)
	if err != nil {
		logger.Error.Printf("Could not fetch ledger.com: %v", err.Error())
		cLedgerCom <- types.FeeAPIResponse3{0, 0, 0}
		return
	}
	result := gjson.GetMany(string(body), "1", "3", "6")
	highFee, medFee, lowFee := result[0].Float()/1000, result[1].Float()/1000, result[2].Float()/1000
	cLedgerCom <- types.FeeAPIResponse3{highFee, medFee, lowFee}
}
