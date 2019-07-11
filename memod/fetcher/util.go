package fetcher

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/0xb10c/memo/memod/logger"
)

// read the response body
func readResponseBody(resp *http.Response) ([]byte, error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error.Println(err.Error())
	}
	defer resp.Body.Close()

	return body, nil
}

func makeHTTPGETReq(url string, timeout int) (body []byte, err error) {
	client := http.Client{Timeout: time.Duration(timeout) * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return
	}

	body, err = readResponseBody(resp)
	if err != nil {
		return
	}

	return
}

func getMyceliumFeeAPICert() (cert []byte) {

	// Mycelium uses a self signed certificate for their feerate API.
	// This cert is stored here.

	certAsPem := `-----BEGIN CERTIFICATE-----
MIIDTTCCAjWgAwIBAgIJAOudLU3vQBxpMA0GCSqGSIb3DQEBCwUAMD0xGzAZBgNV
BAMMEm13czIwLm15Y2VsaXVtLmNvbTERMA8GA1UECgwITXljZWxpdW0xCzAJBgNV
BAYTAlhYMB4XDTE2MTExNjE1MTExMFoXDTE5MDgxMzE1MTExMFowPTEbMBkGA1UE
AwwSbXdzMjAubXljZWxpdW0uY29tMREwDwYDVQQKDAhNeWNlbGl1bTELMAkGA1UE
BhMCWFgwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC7jTD9MLV/k4y2
qJH+aLfJU3pdMmdlciAz1L3T03rT6QDOaT4eanGqpeYNqO/udp8GOFzvnJJ1qXU9
XI3zYOVwK+m/3JBG5B4olQjibkahqi4yterxxXjSzGP5apG//9kfKPx8Q2P47EG+
wt2cdVF2WPicK5+42D7QPwM3cEcohPzzHHaVB0tFt9bcFgooGCIlCz/mo7rD3PsK
nxmdx/0T3Tyh8iLCuLh/PqcLPWKPtpgy3vo3W9gxVnfZj0KR8qMMbQV8KjjCR74s
l3Mfs3bfSKuJKLxK/qElu3BZJRZ6CZzodSb6n0+s9qHrwfJ2FjNvj7jzACv+lTU6
DHGgRd3rAgMBAAGjUDBOMB0GA1UdDgQWBBT2Di9Vrc35R3LSmUfQbPzJrFo6PTAf
BgNVHSMEGDAWgBT2Di9Vrc35R3LSmUfQbPzJrFo6PTAMBgNVHRMEBTADAQH/MA0G
CSqGSIb3DQEBCwUAA4IBAQBFVqonJzLPL/5OY+yy/AJnMqscgGNMiKn9lMi9xU1H
2mx1Tk4ziJgfT7OtoPIwTPMkqGLRfh6gGQnuePmvuG9MrjNNYBEPB0/esvVOws8V
wgYESOC6b4uGyLCv79gyUQFQgwo7CgMxs1vltZEk3DUx1y6eiHGyLiSEE5fmxLQY
xUGiv1w/ZWSlqDiYRl7BdjVtVDxqXyPgcBIW9+k+iRhae1M/nCB8ZpU2IR8x2IjC
fDm1IEMFno98J0hAFyHPVBXsXWLKuDRNbRAWJpe2TmK8G4/MxwgOS41RASTG6icK
mFjasRwenevsfGl0e5Nsth0ToynsQzuO3Tv2NzQbYgOG
-----END CERTIFICATE-----`

	return []byte(certAsPem)
}
