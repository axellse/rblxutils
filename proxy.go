package main

import (
	"context"
	"crypto/tls"
	"fmt"
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

			fmt.Println(r.In.Host)
			switch r.In.Host {
			case "fts.rbxcdn.com":
				fmt.Println("fts: " + rbxcdnIp.String())
				r.SetURL(rbxcdnHostUrl)
			case "assetdelivery.roblox.com":
				fmt.Println("assetdelivery: " + assetdeliveryIp.String())
				r.SetURL(assetdeliveryHostUrl)
			}
			r.Out.Host = r.In.Host

			fmt.Println("what kind of asset request is this? this is a ", r.In.URL.String(), "asset request.")
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
