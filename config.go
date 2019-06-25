package consul

import "github.com/hashicorp/consul/api"

type ConsulConfig struct {
	ConsulAddr string //consul 地址
	DnsScheme  string //默认为 consul
}

func setupDefaultConfig(c *ConsulConfig) *api.Config {
	if c.DnsScheme == "" {
		c.DnsScheme = "consul"
	}
	cfg := api.DefaultConfig()
	if c.ConsulAddr != "" {
		cfg.Address = c.ConsulAddr
	}
	return cfg
}

func setupServiceCheck(c *api.AgentServiceCheck, id string) {
	if c == nil {
		return
	}
	if c.CheckID == "" {
		c.CheckID = id
	}
	if c.TTL == "" && c.Interval == "" {
		c.Interval = "30s"
	}
	if c.Timeout == "" {
		c.Timeout = "3s"
	}

	//if c.DeregisterCriticalServiceAfter==""{
	//	c.DeregisterCriticalServiceAfter="1m"
	//}

	//if c.Status==""{
	//	c.Status="passing"
	//}
}

func setupServiceRegistration(reg *api.AgentServiceRegistration) {
	if reg.ID == "" {
		reg.ID = reg.Name
	}

}
