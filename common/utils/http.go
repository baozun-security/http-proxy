package utils

import (
	"golang.org/x/net/idna"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"unicode/utf8"
)

func canonicalAddr(url *url.URL) string {
	addr := url.Hostname()
	if v, err := toASCII(addr); err == nil {
		addr = v
	}
	port := url.Port()
	if port == "" {
		switch url.Scheme {
		case "http":
			port = "80"
		case "https":
			port = "443"
		case "socks5":
			port = "1080"
		}
	}
	return net.JoinHostPort(addr, port)
}

func toASCII(v string) (string, error) {
	if isASCII(v) {
		return v, nil
	}
	return idna.Lookup.ToASCII(v)
}

func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] >= utf8.RuneSelf {
			return false
		}
	}
	return true
}

// 根据请求获取服务器地址信息
func GetServerAddress(r *http.Request) string {
	var hasPort = regexp.MustCompile(`:\d+$`)
	if r.Method == "CONNECT" { // https
		host := r.URL.Host
		if !hasPort.MatchString(host) { // http 是没有端口的
			host += ":80"
		}
		return host
	} else {
		return canonicalAddr(r.URL)
	}
}
