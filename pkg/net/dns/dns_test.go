package dns

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net"
	"testing"
)

func TestDNS3(t *testing.T) {
	t.Log(0b10000001)
	t.Log(fmt.Sprintf("%08b", 0b00000011)[4:])
}

func TestDNS6(t *testing.T) {
	t.Log(base64.URLEncoding.EncodeToString(creatRequest("www.example.com", A, false)))
	t.Log(base64.URLEncoding.EncodeToString(creatRequest("www.google.com", A, false)))
	t.Log(base64.URLEncoding.EncodeToString(creatRequest("a.62characterlabel-makes-base64url-distinct-from-standard-base64.example.com", A, false)))
}

func TestDNSResolver(t *testing.T) {
	b := bytes.NewBuffer(nil)

	b.WriteByte(0x01)

	z := b.Bytes()

	t.Log(z)

	b.Reset()
	b.WriteByte(0x02)

	t.Log(z, b.Bytes())

}

func TestResolve(t *testing.T) {
	req := []byte{46, 230, 1, 0, 0, 1, 0, 0, 0, 0, 0, 1, 7, 98, 114, 111, 119, 115, 101, 114, 4, 112, 105, 112, 101, 4, 97, 114, 105, 97, 9, 109, 105, 99, 114, 111, 115, 111, 102, 116, 3, 99, 111, 109, 0, 0, 1, 0, 1, 0, 0, 41, 16, 0, 0, 0, 0, 0, 0, 12, 0, 8, 0, 8, 0, 1, 22, 0, 223, 5, 4, 0}
	ans := []byte{46, 230, 129, 128, 0, 1, 0, 3, 0, 0, 0, 1, 7, 98, 114, 111, 119, 115, 101, 114, 4, 112, 105, 112, 101, 4, 97, 114, 105, 97, 9, 109, 105, 99, 114, 111, 115, 111, 102, 116, 3, 99, 111, 109, 0, 0, 1, 0, 1, 192, 12, 0, 5, 0, 1, 0, 0, 13, 58, 0, 40, 7, 98, 114, 111, 119, 115, 101, 114, 6, 101, 118, 101, 110, 116, 115, 4, 100, 97, 116, 97, 14, 116, 114, 97, 102, 102, 105, 99, 109, 97, 110, 97, 103, 101, 114, 3, 110, 101, 116, 0, 192, 61, 0, 5, 0, 1, 0, 0, 0, 43, 0, 32, 20, 115, 107, 121, 112, 101, 100, 97, 116, 97, 112, 114, 100, 99, 111, 108, 99, 117, 115, 49, 50, 8, 99, 108, 111, 117, 100, 97, 112, 112, 192, 96, 192, 113, 0, 1, 0, 1, 0, 0, 0, 6, 0, 4, 13, 89, 202, 241, 0, 0, 41, 16, 0, 0, 0, 0, 0, 0, 0}

	t.Log(Resolve(req, ans))
}

func TestXxxx(t *testing.T) {
	i := net.ParseIP("10.0.2.1")
	t.Log([]byte(i.To16()))
	t.Log([]byte(net.ParseIP("ff::ff")))
}
