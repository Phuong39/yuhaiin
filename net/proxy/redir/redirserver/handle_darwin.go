package redirserver

import (
	"net"

	"github.com/Asutorufa/yuhaiin/net/common"
	"github.com/Asutorufa/yuhaiin/net/proxy/redir/pfutil"
)

func handle(req net.Conn, dst func(string) (net.Conn, error)) error {
	defer req.Close()
	_ = req.(*net.TCPConn).SetKeepAlive(true)
	target, err := pfutil.NatLookup(req.(*net.TCPConn))
	if err != nil {
		return err
	}

	rsp, err := dst(target.String())
	if err != nil {
		return err
	}
	switch rsp.(type) {
	case *net.TCPConn:
		_ = rsp.(*net.TCPConn).SetKeepAlive(true)
	}
	defer rsp.Close()
	common.Forward(req, rsp)
	return nil
}
