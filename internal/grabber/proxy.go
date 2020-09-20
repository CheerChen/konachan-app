package grabber

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"net/url"

	"github.com/CheerChen/konachan-app/internal/log"

	"golang.org/x/net/proxy"
)

var proxyClient *http.Client
var hostname string
var ErrRecordNotFound = errors.New("record not found")

func SetHost(host string) {
	hostname = host
}

func SetProxy(enable bool, socketUrl string) {
	proxyClient = &http.Client{}
	if enable {
		url, err := url.Parse(socketUrl)
		if err != nil {
			log.Warnf("Error parsing proxy url, %s", err)
			return
		}
		dialer, err := proxy.FromURL(url, proxy.Direct)
		if err != nil {
			log.Warnf("Error Dialer, %s", err)
			return
		}
		proxyClient.Transport = &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
				c, e := dialer.Dial(network, addr)
				return c, e
			},
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
}
