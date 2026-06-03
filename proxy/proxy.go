package proxy

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/axellse/rblxutils/common"
	"github.com/axellse/rblxutils/resources"
)

const DelaySamples = 30

var RewriteDelaysNs = []int{}
var ModifyResponseAssetDeliveryDelaysNs = []int{}
var ModifyResponseCdnDelaysNs = []int{}

func IndentBytes(input []byte, prefix, indent string) []byte {
	var out bytes.Buffer
	json.Indent(&out, input, prefix, indent)
	return out.Bytes()
}

func StartProxy() {
	fmt.Println("Launched rblxutils proxy!")
	fmt.Println("modifying hosts file...")
	ModifyHostsFile(false)

	rbxcdnIp, assetdeliveryIp := LookupIps()

	ModifyHostsFile(true)
	fmt.Println("rbxcdn remote is", rbxcdnIp.String(), "assetdelivery is", assetdeliveryIp.String())

	rbxcdnHostUrl, _ := url.Parse("https://fts.rbxcdn.com")
	assetdeliveryHostUrl, _ := url.Parse("https://assetdelivery.roblox.com")
	fmt.Println("parsing certs...")
	rbxcdnCert, err := tls.X509KeyPair(resources.RbxcdnCert, resources.RbxcdnKey)
	assetdeliveryCert, err2 := tls.X509KeyPair(resources.AssetdeliveryCert, resources.AssetdeliveryKey)
	if err != nil || err2 != nil {
		common.FatalError(err)
	}

	rules := ConsolidateMods()
	urlBlobMap := map[string][]byte{}
	urlBlobMapMutex := sync.Mutex{}
	var killFunc func() error

	proxy := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			st := time.Now()

			switch r.In.Host {
			case "assetdelivery.roblox.com":
				r.SetURL(assetdeliveryHostUrl)
				if r.In.URL.Path == "/v1/assets/batch" {
					rawBody := bytes.Buffer{}

					bodyRd := io.TeeReader(r.In.Body, &rawBody)
					uncBodyRd, err := NewBodyDecodingReader(r.In.Header.Get("content-encoding"), bodyRd)
					if err != nil {
						fmt.Println("failed creating body decompressor")
					}

					bodyBa, err := io.ReadAll(uncBodyRd)
					if err != nil {
						fmt.Println("failed reading body for assetdelivery request")
						return
					}

					var reqs []V1BatchRequest
					err = json.Unmarshal(bodyBa, &reqs)
					if err != nil {
						fmt.Println("failed Unmarshal body for assetdelivery request")
						return
					}

					ctx := context.WithValue(r.Out.Context(), RequestIds, reqs)
					r.Out = r.Out.WithContext(ctx)

					r.Out.Body = io.NopCloser(&rawBody)
					if len(RewriteDelaysNs) >= DelaySamples {
						RewriteDelaysNs = RewriteDelaysNs[1:]
					}
					RewriteDelaysNs = append(RewriteDelaysNs, int(time.Since(st).Nanoseconds()))

				}
			case "fts.rbxcdn.com":
				r.SetURL(rbxcdnHostUrl)
			}

			r.Out.Host = r.In.Host
		},
		ModifyResponse: func(r *http.Response) error {
			t1 := time.Now()
			if r.Request.Host == "assetdelivery.roblox.com" && r.Request.URL.Path == "/v1/assets/batch" {
				rawBody := bytes.Buffer{}

				bodyRd := io.TeeReader(r.Body, &rawBody)
				uncBodyRd, err := NewBodyDecodingReader(r.Header.Get("content-encoding"), bodyRd)
				if err != nil {
					fmt.Println("failed creating body decompressor")
				}

				bodyBa, err := io.ReadAll(uncBodyRd)
				if err != nil {
					fmt.Println("failed reading body for assetdelivery response")
					return err
				}

				responses := []V1BatchResponse{}
				err = json.Unmarshal(bodyBa, &responses)
				if err != nil {
					fmt.Println("failed unmarshal body for assetdelivery response")
					return err
				}

				requests := r.Request.Context().Value(RequestIds).([]V1BatchRequest)
				for i, req := range requests {
					if responses[i].ContentRepresentationSpecifier.Format != "" && responses[i].AssetTypeId == 1 {
						fmt.Println("non-png Image found:", req.AssetId, "of format", responses[i].ContentRepresentationSpecifier.Format) //roblox doesn't seem to care if we serve a png even though it expects a ktx.
					}

					for _, rule := range rules {
						if slices.Contains(rule.Sources.Ids, req.AssetId) || slices.Contains(rule.Sources.Types, responses[i].AssetTypeId) {
							urlBlobMapMutex.Lock()
							urlBlobMap[responses[i].Location] = rule.Data.Blob
							urlBlobMapMutex.Unlock()
						}
					}
				}

				r.Body = io.NopCloser(&rawBody)
				if len(ModifyResponseAssetDeliveryDelaysNs) >= DelaySamples {
					ModifyResponseAssetDeliveryDelaysNs = ModifyResponseAssetDeliveryDelaysNs[1:]
				}
				ModifyResponseAssetDeliveryDelaysNs = append(ModifyResponseAssetDeliveryDelaysNs, int(time.Since(t1).Nanoseconds()))
			} else if r.Request.Host == "assetdelivery.roblox.com" && (r.Request.URL.Path == "/v1/asset" || r.Request.URL.Path == "/v1/asset/") {
				// we need data like the assetTypeId to apply any mods which we cant get of this endpoint,
				// so we just run the asset id through /v1/assets/batch which will give us a new fts.rbxcdn.com
				// url that has mods applied to it via urlBlobMapMutex. This endpoint normally redirects to a
				// fts.rbxcdn.com url so we just need to change the Location header to our new url and we're
				// good to go.

				id, err := strconv.Atoi(r.Request.URL.Query().Get("id"))
				if err != nil {
					return nil
				}

				reqBody := []V1BatchRequest{
					{
						AssetId: id,
						RequestId: "0",
					},
				}
				ba, err := json.Marshal(reqBody)
				if err != nil {
					fmt.Println("failed marshal body for /v1/assets/batch in /v1/asset: ", err)
					return err
				}
				
				rd := bytes.NewReader(ba)

				client := &http.Client{
					Transport: &http.Transport{
						TLSClientConfig: &tls.Config{
							InsecureSkipVerify: true,
						},
					},
				}
				
				req, err := http.NewRequest("POST", "https://assetdelivery.roblox.com/v1/assets/batch", rd)
				if err != nil {
					fmt.Println(err)
					return err
				}

				token, err := r.Request.Cookie(".ROBLOSECURITY")
				if err != nil {
					fmt.Println("no roblox token found, cant make request for /v1/asset")
					return err
				}

				req.AddCookie(token)
				req.Header.Set("Content-Type", "application/json")

				resp, err := client.Do(req)
				if err != nil {
					fmt.Println(err)
					return err
				}

				ba, err = io.ReadAll(resp.Body)
				if err != nil {
					fmt.Println("failed read body for /v1/assets/batch in /v1/asset: ", err)
					return err
				}

				responses := []V1BatchResponse{}
				err = json.Unmarshal(ba, &responses)
				if err != nil {
					fmt.Println("failed unmarshal body for /v1/assets/batch in /v1/asset:", err)
					return err
				}

				if len(responses) == 0 {
					fmt.Println("/v1/assets/batch in /v1/asset responses is empty.")
					return errors.New("empty responses")
				}

				r.Header.Set("Location", responses[0].Location)
			} else if r.Request.Host == "fts.rbxcdn.com" {
				urlBlobMapMutex.Lock()
				blob, ok := urlBlobMap[r.Request.URL.String()]
				urlBlobMapMutex.Unlock()
				if !ok {
					return nil
				}
				urlBlobMapMutex.Lock()
				delete(urlBlobMap, r.Request.URL.String())
				urlBlobMapMutex.Unlock()

				/*				f, err := os.Create(common.LPath("./junk/random_image-" + strconv.Itoa(rand.IntN(10e3)) + ".png"))
								if err != nil {
									fmt.Println("could not open file")
									return errors.New("filerr")
								}

								defer f.Close()
								_, err = io.Copy(f, r.Body)
								if err != nil {
									fmt.Println("could not copy file")
									return errors.New("copyerr")
								}*/

				for _, v := range []string{"Transfer-Encoding", "Content-Encoding"} {
					r.Header.Del(v)
				}

				r.Header.Set("Content-Length", strconv.Itoa(len(blob)))
				r.Body = io.NopCloser(bytes.NewReader(blob))

				if len(ModifyResponseCdnDelaysNs) >= DelaySamples {
					ModifyResponseCdnDelaysNs = ModifyResponseCdnDelaysNs[1:]
				}
				ModifyResponseCdnDelaysNs = append(ModifyResponseCdnDelaysNs, int(time.Since(t1).Nanoseconds()))
			}

			return nil
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
		Handler: GetProxyServemux(proxy, &killFunc),
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{
				rbxcdnCert,
				assetdeliveryCert,
			},
			InsecureSkipVerify: true,
		},
	}
	killFunc = server.Close

	fmt.Println("everything ready, now starting server and awaiting lock...")

	err = server.ListenAndServeTLS("", "")
	if err != nil && err != http.ErrServerClosed {
		ModifyHostsFile(false)
		common.FatalError(err)
	}

	ModifyHostsFile(false)
}

func ModifyHostsFile(add bool) {
	fmt.Println("modifying hosts file", add)
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

// not up to spec but who cares
func NewBodyDecodingReader(encoding string, body io.Reader) (io.Reader, error) {
	switch encoding {
	case "gzip":
		return gzip.NewReader(body)
	case "deflate":
		return zlib.NewReader(body)
	}
	return body, nil
}

func LookupIps() (cdnIp net.IP, assetDeliveryIp net.IP) {
	fmt.Println("looking up fts.rbxcdn.com")
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

	fmt.Println("looking up assetdelivery.roblox.com")
	ips, err = net.LookupIP("assetdelivery.roblox.com")
	if err != nil {
		common.FatalError(err)
	}

	var assetdeliveryIp net.IP
	for _, ip := range ips {
		fmt.Println("found ip: ", ip)
		if ip.To4() != nil {
			assetdeliveryIp = ip
			break
		}
	}

	return rbxcdnIp, assetdeliveryIp
}

type ProxyContextKey int

const (
	RequestIds ProxyContextKey = 0
)

type V1BatchRequest struct {
	AssetId int `json:"assetId"`
	RequestId string `json:"requestId"`
}

type ContentRepresentationSpecifier struct {
	Format string `json:"format"`
}
type V1BatchResponse struct {
	Location                       string                         `json:"location"`
	AssetTypeId                    int                            `json:"assetTypeId"`
	ContentRepresentationSpecifier ContentRepresentationSpecifier `json:"contentRepresentationSpecifier"`
}