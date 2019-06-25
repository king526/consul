package consul

import (
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc/resolver"

	"strconv"
)

type exampleResolverBuilder struct {
	scheme string
	ht     *api.Health
}

func (b *exampleResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOption) (resolver.Resolver, error) {
	r := &exampleResolver{
		ht:     b.ht,
		target: target,
		cc:     cc,
	}
	r.resolve()
	go r.watch()
	return r, nil
}

func (b *exampleResolverBuilder) Scheme() string { return b.scheme }

// exampleResolver is a
// Resolver(https://godoc.org/google.golang.org/grpc/resolver#Resolver).
type exampleResolver struct {
	ht     *api.Health
	target resolver.Target
	cc     resolver.ClientConn
	q      api.QueryOptions
}

func (r *exampleResolver) watch() {
	for {
		r.resolve()
	}
}

func (r *exampleResolver) resolve() error {
	ce := r.parseEndpoint(r.target.Endpoint)
	entrys, meta, err := r.ht.ServiceMultipleTags(ce.service, ce.tags, true, &r.q)
	if err != nil {
		return err
	}
	r.q.WaitIndex = meta.LastIndex
	var addrs []resolver.Address
	for _, item := range entrys {
		addrs = append(addrs, r.parseAddr(item, ce.addrMeta))
	}
	r.cc.NewAddress(addrs)
	return nil
}

type consulEndpoint struct {
	service  string
	addrMeta string
	tags     []string
}

// parseEndpoint  Endpoint now just service Name
func (r *exampleResolver) parseEndpoint(ep string) (ce consulEndpoint) {
	ce.service = ep
	return
}

func (r *exampleResolver) parseAddr(item *api.ServiceEntry, addrMeta string) resolver.Address {
	addr := item.Service.Address + ":" + strconv.Itoa(item.Service.Port)
	if addrMeta != "" && item.Service.Meta != nil {
		if addr2, ok := item.Service.Meta[addrMeta]; ok {
			addr = addr2
		}
	}
	return resolver.Address{
		ServerName: r.target.Endpoint,
		Addr:       addr,
		Metadata:   item.Service,
	}
}

func (*exampleResolver) ResolveNow(o resolver.ResolveNowOption) {}
func (*exampleResolver) Close()                                 {}
