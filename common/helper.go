package common

import (
	"crypto/tls"
	"net/http"
	"strings"
)

func KillHelper() {
	req, err := http.NewRequest("GET", "https://127.0.0.1", nil)
	if err != nil {
		FatalError(err)
	}
	
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	req.Header.Add("x-rblxutils-kill-server", "1")
	_, err = client.Do(req)
	if err != nil && !strings.Contains(err.Error(), "EOF") {
		FatalError(err)
	}
}