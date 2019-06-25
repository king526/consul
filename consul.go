package consul

import (
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"sync"
	"google.golang.org/grpc/resolver"
	"strconv"
	"strings"
)

var (
	gl sync.Mutex
)
type Consul struct{
	cfg *ConsulConfig
	reg *api.AgentServiceRegistration
	api *api.Client
}

func NewConsul(cc *ConsulConfig,reg *api.AgentServiceRegistration,chk *api.AgentServiceCheck)(*Consul,error){
	if  chk==nil{
		chk=&api.AgentServiceCheck{}
	}
	setupServiceCheck(chk)
	setupServiceRegistration(reg)
	reg.Check=chk
	cfg:=setupDefaultConfig(cc)
	cli,err:= api.NewClient(cfg)
	if err!=nil{
		return nil,err
	}
	cs:=&Consul{
		cfg:cc,
		reg:reg,
		api:cli,
	}
	gl.Lock()
	resolver.Register(&exampleResolverBuilder{
		ht:cli.Health(),
		scheme:cc.DnsScheme,
	})
	gl.Unlock()
	return cs,nil
}

func (c *Consul)Regist()error{
	return c.api.Agent().ServiceRegister(c.reg)
}

func (c *Consul)DeRegist()error{
	return c.api.Agent().ServiceDeregister(c.reg.ID)
}

//func (c *Consul)ExitMaintenance()error{
//	return c.api.Agent().DisableServiceMaintenance(c.reg.ID)
//}
//
//func (c *Consul)Maintenance(reason string)error{
//	return c.api.Agent().EnableServiceMaintenance(c.reg.ID,reason)
//}

func (c *Consul)RegistWithGRPCHealth(s *grpc.Server, check GRPCHealthCheckFunc)error{
	gc:=&grpcHealth{
		check:check,
	}
	grpc_health_v1.RegisterHealthServer(s,gc )
	if c.reg.Check.GRPC=="" {
		c.reg.Check.GRPC=c.reg.Address+":"+strconv.Itoa(c.reg.Port)
	}
	if !strings.Contains( c.reg.Check.GRPC,"/"){
		c.reg.Check.GRPC+= "/grpc.health.v1.Health"
	}
	return c.api.Agent().ServiceRegister(c.reg)
}


