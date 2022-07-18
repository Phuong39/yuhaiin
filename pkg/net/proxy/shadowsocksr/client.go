package shadowsocksr

import (
	"fmt"
	"net"
	"strconv"

	"github.com/Asutorufa/yuhaiin/pkg/net/interfaces/proxy"
	"github.com/Asutorufa/yuhaiin/pkg/net/proxy/shadowsocks"
	"github.com/Asutorufa/yuhaiin/pkg/net/proxy/shadowsocksr/cipher"
	"github.com/Asutorufa/yuhaiin/pkg/net/proxy/shadowsocksr/obfs"
	"github.com/Asutorufa/yuhaiin/pkg/net/proxy/shadowsocksr/protocol"
	ssr "github.com/Asutorufa/yuhaiin/pkg/net/proxy/shadowsocksr/utils"
	s5c "github.com/Asutorufa/yuhaiin/pkg/net/proxy/socks5/client"
	"github.com/Asutorufa/yuhaiin/pkg/net/utils/resolver"
	"github.com/Asutorufa/yuhaiin/pkg/protos/node"
)

var _ proxy.Proxy = (*Shadowsocksr)(nil)

type Shadowsocksr struct {
	protocol *protocol.Protocol
	obfs     *obfs.Obfs
	cipher   *cipher.Cipher
	dial     proxy.Proxy

	addr string
}

func NewShadowsocksr(config *node.PointProtocol_Shadowsocksr) node.WrapProxy {
	c := config.Shadowsocksr

	return func(p proxy.Proxy) (proxy.Proxy, error) {
		cipher, err := cipher.NewCipher(c.Method, c.Password)
		if err != nil {
			return nil, fmt.Errorf("new cipher failed: %w", err)
		}

		port, err := strconv.ParseUint(c.Port, 10, 16)
		if err != nil {
			return nil, fmt.Errorf("parse port failed: %w", err)
		}

		info := ssr.Info{
			IVSize:  cipher.IVSize(),
			Key:     cipher.Key(),
			KeySize: cipher.KeySize(),
		}

		obfs, err := obfs.NewObfs(c.Obfs, ssr.ObfsInfo{
			Host:  c.Server,
			Port:  uint16(port),
			Param: c.Obfsparam,
			Info:  info,
		})
		if err != nil {
			return nil, fmt.Errorf("new obfs failed: %w", err)
		}

		protocol, err := protocol.NewProtocol(
			c.Protocol,
			protocol.ProtocolInfo{
				Auth:         &protocol.AuthData{},
				Param:        c.Protoparam,
				TcpMss:       1460,
				Info:         info,
				ObfsOverhead: obfs.Overhead(),
			})
		if err != nil {
			return nil, fmt.Errorf("new protocol failed: %w", err)
		}

		return &Shadowsocksr{protocol, obfs, cipher, p, net.JoinHostPort(c.Server, c.Port)}, nil
	}
}

func (s *Shadowsocksr) Conn(addr proxy.Address) (net.Conn, error) {
	c, err := s.dial.Conn(addr)
	if err != nil {
		return nil, fmt.Errorf("get conn failed: %w", err)
	}
	// obfsServerInfo.SetHeadLen(b, 30)
	// protocolServerInfo.SetHeadLen(b, 30)
	obfs := s.obfs.Stream(c)
	cipher := s.cipher.StreamConn(obfs)

	var iv []byte
	if z, ok := cipher.(interface{ WriteIV() []byte }); ok {
		iv = z.WriteIV()
	}

	conn := s.protocol.Stream(cipher, iv)
	if _, err := conn.Write(s5c.ParseAddr(addr)); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("write target failed: %w", err)
	}

	return conn, nil
}

func (s *Shadowsocksr) PacketConn(addr proxy.Address) (net.PacketConn, error) {
	c, err := s.dial.PacketConn(addr)
	if err != nil {
		return nil, fmt.Errorf("get packet conn failed: %w", err)
	}
	cipher := s.cipher.PacketConn(c)
	proto := s.protocol.Packet(cipher)
	uaddr, err := resolver.ResolveUDPAddr(s.addr)
	if err != nil {
		c.Close()
		return nil, fmt.Errorf("resolve udp addr failed: %w", err)
	}
	return shadowsocks.NewSsPacketConn(proto, uaddr, addr), nil
}
