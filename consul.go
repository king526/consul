package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/resolver"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	gl sync.Mutex
)

type Consul struct {
	cfg *ConsulConfig
	reg *api.AgentServiceRegistration
	api *api.Client
}

func NewConsul(cc *ConsulConfig, reg *api.AgentServiceRegistration, chk *api.AgentServiceCheck) (*Consul, error) {
	setupServiceRegistration(reg)
	setupServiceCheck(chk, reg.ID)
	reg.Check = chk
	cfg := setupDefaultConfig(cc)
	cli, err := api.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	cs := &Consul{
		cfg: cc,
		reg: reg,
		api: cli,
	}
	gl.Lock()
	resolver.Register(&exampleResolverBuilder{
		ht:     cli.Health(),
		scheme: cc.DnsScheme,
	})
	gl.Unlock()
	return cs, nil
}

func (c *Consul) Regist() error {
	return c.api.Agent().ServiceRegister(c.reg)
}

func (c *Consul) DeRegist() error {
	return c.api.Agent().ServiceDeregister(c.reg.ID)
}

func (c *Consul) RegistWithTTL(f func() (status, note string)) error {
	ds, err := time.ParseDuration(c.reg.Check.TTL)
	if err != nil {
		return err
	}
	if f == nil {
		f = func() (status, note string) {
			return "passing", ""
		}
	}
	go c.doTTL(ds, f)
	return c.api.Agent().ServiceRegister(c.reg)
}

func (c *Consul) RegistWithGRPCHealth(s *grpc.Server, check GRPCHealthCheckFunc) error {
	gc := &grpcHealth{
		check: check,
	}
	grpc_health_v1.RegisterHealthServer(s, gc)
	if c.reg.Check.GRPC == "" {
		c.reg.Check.GRPC = c.reg.Address + ":" + strconv.Itoa(c.reg.Port)
	}
	if !strings.Contains(c.reg.Check.GRPC, "/") {
		c.reg.Check.GRPC += "/grpc.health.v1.Health"
	}
	return c.api.Agent().ServiceRegister(c.reg)
}

//func (c *Consul)ExitMaintenance()error{
//	return c.api.Agent().DisableServiceMaintenance(c.reg.ID)
//}
//
//func (c *Consul)Maintenance(reason string)error{
//	return c.api.Agent().EnableServiceMaintenance(c.reg.ID,reason)
//}

func (c *Consul) doTTL(ds time.Duration, f func() (string, string)) {
	ds -= time.Second
	if ds < time.Second {
		ds = time.Second * 10
	}
	var status, note string
	for range time.Tick(ds) {
		status, note = f()
		err := c.api.Agent().UpdateTTL(c.reg.Check.CheckID, note, status)
		if err != nil {
			fmt.Println(err)
		}
	}
}
