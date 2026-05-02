package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"

	"axell.me/rblxutils/common"
	"axell.me/rblxutils/resources"
)

func StartProxy() {
	ModifyHostsFile(false)
	ips, err := net.LookupIP("fts.rbxcdn.com")
	if err != nil {
		common.FatalError(err)
	}

	var rbxcdnIp net.IP
	for _, ip := range ips {
		fmt.Println("found ip: ", ip)
		if ip.To4() != nil {
			rbxcdnIp = ip
			break
		}
	}

	ips, err = net.LookupIP("assetdelivery.roblox.com")
	if err != nil {
		common.FatalError(err)
	}

	var assetdeliveryIp net.IP
	for _, ip := range ips {
		if ip.To4() != nil {
			assetdeliveryIp = ip
			break
		}
	}

	ModifyHostsFile(true)
	fmt.Println("rbxcdn remote is", rbxcdnIp.String(), "assetdelivery is", assetdeliveryIp.String())

	rbxcdnHostUrl, _ := url.Parse("https://fts.rbxcdn.com")
	assetdeliveryHostUrl, _ := url.Parse("https://assetdelivery.roblox.com")

	rbxcdnCert, err := tls.X509KeyPair(resources.RbxcdnCert, resources.RbxcdnKey)
	if err != nil {
		common.FatalError(err)
	}

	assetdeliveryCert, err := tls.X509KeyPair(resources.AssetdeliveryCert, resources.AssetdeliveryKey)
	if err != nil {
		common.FatalError(err)
	}

	killFunc := func() error { return nil }
	proxy := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			if r.In.Header.Get("x-rblxutils-kill-server") == "1" {
				killFunc()
				return
			}

			if r.In.Host == "assetdelivery.roblox.com" {
				fmt.Println(r.In.URL.Path)
				if r.In.URL.Path == "/v1/assets/batch" {
					bodyBa, err := io.ReadAll(r.In.Body)
					if err != nil {
						fmt.Println("failed reading body for assetdelivery request")
						return
					}

					fmt.Println(string(bodyBa))
					r.Out.Body = io.NopCloser(bytes.NewReader(bodyBa))
				}
			} else if r.In.Host == "fts.rbxcdn.com" {
				r.SetURL(rbxcdnHostUrl)

			}

			r.Out.Host = r.In.Host
			fmt.Println(r.In.URL.String())
		},
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				host, port, err := net.SplitHostPort(addr)
				if err != nil {
					return nil, err
				}

				switch host {
				case "fts.rbxcdn.com":
					host = rbxcdnIp.String()
				case "assetdelivery.roblox.com":
					host = assetdeliveryIp.String()
				}

				var dialer net.Dialer
				return dialer.DialContext(ctx, network, net.JoinHostPort(host, port))
			},
		},
	}

	server := http.Server{
		Addr:    ":443",
		Handler: proxy,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{
				rbxcdnCert,
				assetdeliveryCert,
			},
			InsecureSkipVerify: true,
		},
	}
	killFunc = server.Close

	common.LoadState()
	common.State.HelperAction = ""
	err = common.WriteState()
	if err != nil {
		common.FatalError(err)
	}

	err = server.ListenAndServeTLS("", "")
	if err != nil && err != http.ErrServerClosed {
		common.FatalError(err)
	}

	ModifyHostsFile(false)
}

func ModifyHostsFile(add bool) {
	ba, err := os.ReadFile("C:\\Windows\\System32\\drivers\\etc\\hosts")
	if err != nil {
		common.FatalError(err)
	}

	hosts := strings.ReplaceAll(string(ba), "\r", "")
	hosts = strings.ReplaceAll(hosts, "\n", "\r\n") //make sure the host file is clean

	lines := []string{}
	for line := range strings.SplitSeq(hosts, "\r\n") {
		if !strings.Contains(line, "fts.rbxcdn.com") && !strings.Contains(line, "assetdelivery.roblox.com") && !strings.Contains(line, "rblxutils") {
			lines = append(lines, line)
		}
	}

	if add {
		lines = append(lines, "# The following two lines were inserted by rblxutils. They should be automatically removed when rblxutils exits.")
		lines = append(lines, "  127.0.0.1     fts.rbxcdn.com")
		lines = append(lines, "  127.0.0.1     assetdelivery.roblox.com")
	}

	finalBa := strings.Join(lines, "\r\n")
	err = os.WriteFile("C:\\Windows\\System32\\drivers\\etc\\hosts", []byte(finalBa), 0666)
	if err != nil {
		common.FatalError(err)
	}

	time.Sleep(400 * time.Millisecond)
	err = exec.Command("ipconfig", "/flushdns").Run()
	if err != nil {
		common.FatalError(err)
	}
	time.Sleep(500 * time.Millisecond)
}
