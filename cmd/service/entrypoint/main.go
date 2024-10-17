package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/hashicorp/consul/api"

	"origadmin/basic-layout/api/v1/services/helloworld"
)

func main() {
	cfg := api.DefaultConfig()
	cfg.Address = "192.168.28.42:8500"
	consulCli, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	r := consul.New(consulCli)

	// new grpc client
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///helloworld"),
		grpc.WithDiscovery(r),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	gClient := helloworld.NewGreeterServiceClient(conn)

	// new http client
	hConn, err := http.NewClient(
		context.Background(),
		http.WithMiddleware(
			recovery.Recovery(),
		),
		http.WithEndpoint("discovery:///helloworld"),
		//http.WithEndpoint("127.0.0.1:8000"),
		http.WithDiscovery(r),
		//http.WithBlock(),
		//http.WithTimeout(time.Second*5),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer hConn.Close()
	hClient := helloworld.NewGreeterServiceHTTPClient(hConn)
	fmt.Println("start")
	for {
		time.Sleep(time.Second)
		callGRPC(gClient)
		callHTTP(hClient)
	}
}

func callGRPC(client helloworld.GreeterServiceClient) {
	req := &helloworld.GreeterRequest{Id: "kratos", Name: "kratos", Data: &helloworld.Greeter{
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

func callHTTP(client helloworld.GreeterServiceHTTPClient) {
	req := &helloworld.GreeterRequest{Id: "kratos", Name: "kratos", Data: &helloworld.Greeter{
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