package statistic

import (
	"errors"

	"github.com/Asutorufa/yuhaiin/pkg/net/interfaces/proxy"
	"github.com/Asutorufa/yuhaiin/pkg/net/proxy/direct"
	"github.com/Asutorufa/yuhaiin/pkg/net/utils/resolver"
	protoconfig "github.com/Asutorufa/yuhaiin/pkg/protos/config"
	"github.com/Asutorufa/yuhaiin/pkg/protos/statistic"
)

type router struct {
	remotedns *remotedns
	localdns  *localdns
	bootstrap *bootstrap
	statistic *counter
	shunt     *shunt
}

func NewRouter(dialer proxy.Proxy) *router {
	c := &router{statistic: NewStatistic()}

	c.localdns = newLocaldns(c.statistic)
	c.bootstrap = newBootstrap(c.statistic)
	resolver.Bootstrap = c.bootstrap
	c.remotedns = newRemotedns(direct.Default, dialer, c.statistic)

	c.shunt = newShunt(c.remotedns, c.statistic)
	c.shunt.AddDialer(PROXY, dialer, c.remotedns)
	c.shunt.AddDialer(DIRECT, direct.Default, c.localdns)
	c.shunt.AddDialer(BLOCK, proxy.NewErrProxy(errors.New("block")), c.localdns)

	return c
}

func (a *router) Update(s *protoconfig.Setting) {
	a.shunt.Update(s)
	a.localdns.Update(s)
	a.bootstrap.Update(s)
	a.remotedns.Update(s)
}

func (a *router) Proxy() proxy.Proxy { return a.shunt }

func (a *router) Statistic() statistic.ConnectionsServer { return a.statistic }