package consul

import (
	"google.golang.org/grpc/resolver"
	"github.com/hashicorp/consul/api"

	"strconv"
)



type exampleResolverBuilder struct{
	scheme string
	ht *api.Health
}

func (b *exampleResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOption) (resolver.Resolver, error) {
	r := &exampleResolver{
		ht:b.ht,
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
	ht        *api.Health
	target     resolver.Target
	cc         resolver.ClientConn
	q api.QueryOptions

}

func (r *exampleResolver) watch(){
	for {
		r.resolve()
	}
}

//TODO Endpoint now just service Name
func (r *exampleResolver) resolve()error {
	var addrs []resolver.Address
	entrys,meta,err:= r.ht.Service(r.target.Endpoint,"",true,&r.q)
	if err!=nil{
		return err
	}
	r.q.WaitIndex=meta.LastIndex
	for _,item:=range entrys{
		addrs=append(addrs,resolver.Address{
			ServerName:r.target.Endpoint,
			Addr:item.Service.Address+":"+strconv.Itoa(item.Service.Port),
			//Metadata:item.Service,
		})
	}
	r.cc.NewAddress(addrs)
	return nil
}

func (*exampleResolver) ResolveNow(o resolver.ResolveNowOption) {}
func (*exampleResolver) Close()                                 {}

