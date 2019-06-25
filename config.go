package consul

import "github.com/hashicorp/consul/api"

type ConsulConfig struct {
	ConsulAddr string //consul 地址
	DnsScheme string //默认为 consul

	ServerAddr string //本地地址(重要)



}

func setupDefaultConfig(c *ConsulConfig)*api.Config{
	if c.DnsScheme==""{
		c.DnsScheme="consul"
	}
	cfg:=api.DefaultConfig()
	if c.ConsulAddr!="" {
		cfg.Address = c.ConsulAddr
	}
	return cfg
}

func setupServiceCheck(c *api.AgentServiceCheck){
	if c.Interval==""{
		c.Interval="30s"
	}
}

func setupServiceRegistration(reg *api.AgentServiceRegistration){
	if reg.ID==""{
		reg.ID=reg.Name
	}
}