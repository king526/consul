package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/king526/consul"
	"google.golang.org/grpc"
	"google.golang.org/grpc/examples/helloworld/helloworld"
)

func main1(port int) {
	id := strconv.Itoa(port)
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := &gs{
		idx: id,
	}
	svr := grpc.NewServer()
	cc, _ := consul.NewConsul(
		&consul.ConsulConfig{
			ConsulAddr: "192.168.82.2:8500",
		},
		&api.AgentServiceRegistration{
			ID:      id,
			Name:    "ct",
			Address: "192.168.85.10",
			Port:    port,
		},
		&api.AgentServiceCheck{
			TTL: "10s",
		},
	)
	helloworld.RegisterGreeterServer(svr, s)
	//fmt.Println(id,cc.RegistWithGRPCHealth(svr,nil))
	fmt.Println(id, cc.RegistWithTTL(nil))
	svr.Serve(lis)
}

func main() {
	go main1(9810)
	go main1(9811)
	time.Sleep(time.Second)
	request()

}

type gs struct {
	idx string
}

func (s *gs) SayHello(ctx context.Context, r *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	fmt.Println(s.idx, r.String())
	return &helloworld.HelloReply{Message: fmt.Sprintf(s.idx)}, nil
}

func request() {
	t, _ := context.WithTimeout(context.Background(), time.Second)
	conn, _ := grpc.DialContext(t, "consul:///ct", grpc.WithInsecure(), grpc.WithBalancerName("round_robin"))
	for range time.Tick(time.Second * 3) {
		rep, err := helloworld.NewGreeterClient(conn).SayHello(context.Background(), &helloworld.HelloRequest{
			Name: fmt.Sprint("req"),
		})
		fmt.Println(rep, err)
	}
}
