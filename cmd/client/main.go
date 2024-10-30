package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/selector"
	"github.com/go-kratos/kratos/v2/selector/filter"
	"github.com/go-kratos/kratos/v2/selector/random"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/hashicorp/consul/api"
	"github.com/origadmin/toolkits/utils/replacer"

	"origadmin/basic-layout/api/v1/services/helloworld"
	"origadmin/basic-layout/internal/mods/helloworld/dto"
)

func main() {
	address := "${CONSUL_ADDRESS=127.0.0.1:8500}"
	m, err := replacer.NewMatchFile(
		"resources/env/env.toml",
		replacer.WithMatchFold(true),
		replacer.WithMatchSta("${"),
		replacer.WithMatchEnd("}"))
	if err != nil {
		panic(err)
	}
	address = m.Replace(address)
	cfg := api.DefaultConfig()
	cfg.Address = address
	fmt.Println("consul address:", address)
	consulCli, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	r := consul.New(consulCli)
	serviceName := "origadmin.service.v1.helloworld"
	// new grpc client
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///"+serviceName),
		grpc.WithDiscovery(r),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	gClient := helloworld.NewGreeterClient(conn)

	// 创建路由 Filter：筛选版本号为"2.0.0"的实例
	filter := filter.Version("v1.0.0")
	// 创建 P2C 负载均衡算法 Selector，并将路由 Filter 注入
	selector.SetGlobalSelector(random.NewBuilder())
	//selector.SetGlobalSelector(wrr.NewBuilder())
	// new http client
	hConn, err := http.NewClient(
		context.Background(),
		http.WithMiddleware(
			recovery.Recovery(),
		),
		http.WithEndpoint("discovery:///"+serviceName),
		//http.WithEndpoint("127.0.0.1:8000"),
		http.WithDiscovery(r),
		http.WithNodeFilter(filter),
		//http.WithBlock(),
		//http.WithTimeout(time.Second*5),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer hConn.Close()
	hClient := helloworld.NewGreeterHTTPClient(hConn)
	fmt.Println("start")
	for {
		time.Sleep(time.Second)
		callGRPC(gClient)
		_ = gClient
		callHTTP(hClient)
	}
}

func callGRPC(client helloworld.GreeterClient) {
	req := &helloworld.GreeterRequest{
		Id:   "kratos",
		Name: "kratos",
		Data: &dto.Greeter{
			Id:   "kratos",
			Name: "kratos",
		}}
	err := req.ValidateAll()
	if err != nil {
		log.Print("[grpc] SayHello validate ", err)
		return
	}
	reply, err := client.PostHello(context.Background(), req)
	if err != nil {
		log.Print("[grpc] SayHello ", err)
		return
	}
	log.Printf("[grpc] SayHello %+v\n", reply.Data)
}

func callHTTP(client helloworld.GreeterHTTPClient) {
	req := &helloworld.GreeterRequest{
		Id:   "kratos",
		Name: "kratos",
		Data: &dto.Greeter{
			Id:   "kratos",
			Name: "kratos",
		}}
	err := req.ValidateAll()
	if err != nil {
		log.Print("[http] SayHello validate ", err)
		return
	}
	reply, err := client.PostHello(context.Background(), req)
	if err != nil {
		log.Print("[http] SayHello ", err)
		return
	}
	log.Printf("[http] SayHello %v\n", reply.Data)
}
