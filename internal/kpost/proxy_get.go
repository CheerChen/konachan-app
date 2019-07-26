package kpost

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func proxyGet(url string) (b []byte, err error) {
	client := &http.Client{}
	//dialer, _ := proxy.SOCKS5("tcp", "127.0.0.1:1080",
	//	nil,
	//	&net.Dialer{
	//		Timeout:   5 * time.Second,
	//		KeepAlive: 5 * time.Second,
	//	},
	//)
	//client.Transport = &http.Transport{
	//	Dial: dialer.Dial,
	//	//Proxy:           http.ProxyURL(proxy),
	//	//TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	//}
	fmt.Println("start fetching url", url)
	resp, err := client.Get(url)
	if err != nil {
		return b, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
