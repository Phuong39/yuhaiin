package parser

import (
	"context"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"testing"

	"github.com/Asutorufa/yuhaiin/pkg/net/dns"
	"github.com/Asutorufa/yuhaiin/pkg/net/interfaces/proxy"
	ss "github.com/Asutorufa/yuhaiin/pkg/net/proxy/shadowsocks"
	"github.com/Asutorufa/yuhaiin/pkg/net/proxy/simple"
	"github.com/Asutorufa/yuhaiin/pkg/node/register"
	"github.com/Asutorufa/yuhaiin/pkg/protos/config"
	"github.com/Asutorufa/yuhaiin/pkg/protos/node"
	"github.com/Asutorufa/yuhaiin/pkg/utils/assert"
	"google.golang.org/protobuf/encoding/protojson"
)

func TestSsrParse2(t *testing.T) {
	ssr := []string{"ssr://MS4xLjEuMTo1MzphdXRoX2NoYWluX2E6bm9uZTpodHRwX3NpbXBsZTo2YUtkNW9HcDZMcXIvP29iZnNwYXJhbT02YUtkNW9HcDZMcXImcHJvdG9wYXJhbT02YUtkNW9HcDZMcXImcmVtYXJrcz02YUtkNW9HcDZMcXImZ3JvdXA9NmFLZDVvR3A2THFy",
		"ssr://MS4xLjEuMTo1MzphdXRoX2NoYWluX2E6bm9uZTpodHRwX3NpbXBsZTo2YUtkNW9HcDZMcXIvP29iZnNwYXJhbT02YUtkNW9HcDZMcXImcHJvdG9wYXJhbT02YUtkNW9HcDZMcXImcmVtYXJrcz1jMlZqYjI1ayZncm91cD02YUtkNW9HcDZMcXIK",
		"ssr://MS4xLjEuMTo1MzphdXRoX2NoYWluX2E6bm9uZTpodHRwX3NpbXBsZTo2YUtkNW9HcDZMcXIvP29iZnNwYXJhbT02YUtkNW9HcDZMcXImcHJvdG9wYXJhbT02YUtkNW9HcDZMcXImcmVtYXJrcz1jM056YzNOeiZncm91cD1jM056YzNOego",
		"ssr://MjIyLjIyMi4yMjIuMjIyOjQ0MzphdXRoX2FlczEyOF9tZDU6Y2hhY2hhMjAtaWV0ZjpodHRwX3Bvc3Q6ZEdWemRBby8/b2Jmc3BhcmFtPWRHVnpkQW8mcHJvdG9wYXJhbT1kR1Z6ZEFvJnJlbWFya3M9ZEdWemRBbyZncm91cD1kR1Z6ZEFvCg"}

	for x := range ssr {
		log.Println(Parse(node.NodeLink_shadowsocksr, []byte(ssr[x])))
	}
}

func TestConnections(t *testing.T) {
	p := simple.NewSimple(proxy.ParseAddressSplit("", "127.0.0.1", proxy.ParsePort(1090)), nil)

	z, err := ss.NewHTTPOBFS(
		&node.PointProtocol_ObfsHttp{
			ObfsHttp: &node.ObfsHttp{
				Host: "example.com",
				Port: "80",
			},
		})(p)
	assert.NoError(t, err)
	z, err = ss.NewShadowsocks(
		&node.PointProtocol_Shadowsocks{
			Shadowsocks: &node.Shadowsocks{
				Method:   "AEAD_AES_128_GCM",
				Password: "test",
				Server:   "127.0.0.1",
				Port:     "1090",
			},
		})(z)
	assert.NoError(t, err)
	tt := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				ad, err := proxy.ParseAddress(network, addr)
				assert.NoError(t, err)
				return z.Conn(ad)
			},
		},
	}

	req := http.Request{
		Method: "GET",
		URL: &url.URL{
			Scheme: "http",
			Host:   "ip.sb",
		},
		Header: make(http.Header),
	}
	req.Header.Set("User-Agent", "curl/v2.4.1")
	resp, err := tt.Do(&req)
	assert.NoError(t, err)
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	t.Log(string(data))
}

func TestConnectionSsr(t *testing.T) {
	p := &node.Point{
		Protocols: []*node.PointProtocol{},
	}

	err := protojson.Unmarshal([]byte(``), p)
	assert.NoError(t, err)
	z, err := register.Dialer(p)
	assert.NoError(t, err)

	tt := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				ad, err := proxy.ParseAddress(network, addr)
				assert.NoError(t, err)
				return z.Conn(ad)
			},
		},
	}

	dns := dns.New(dns.Config{
		Type: config.Dns_udp,
		Host: "1.1.1.1:53", Dialer: z})
	t.Log(dns.LookupIP("www.google.com"))

	req := http.Request{
		Method: "GET",
		URL: &url.URL{
			Scheme: "http",
			Host:   "ip.sb",
		},
		Header: make(http.Header),
	}
	req.Header.Set("User-Agent", "curl/v2.4.1")
	resp, err := tt.Do(&req)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	assert.NoError(t, err)
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	t.Log(string(data))
}

func TestSSr(t *testing.T) {
	p := &node.Point{
		Protocols: []*node.PointProtocol{},
	}
	z, err := register.Dialer(p)
	assert.NoError(t, err)

	tt := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				ad, err := proxy.ParseAddress(network, addr)
				assert.NoError(t, err)
				return z.Conn(ad)
			},
		},
	}

	req := http.Request{
		Method: "GET",
		URL: &url.URL{
			Scheme: "http",
			Host:   "ip.sb",
		},
		Header: make(http.Header),
	}
	req.Header.Set("User-Agent", "curl/v2.4.1")
	resp, err := tt.Do(&req)
	t.Error(err)
	assert.NoError(t, err)
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	t.Log(string(data))
}
