package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"google.golang.org/grpc"

	"github.com/yakaa/grpcx"
	"github.com/yakaa/grpcx/config"
	"github.com/yakaa/grpcx/example/proto"
)

var Endpoints = []string{"192.168.1.200:2379"}

func main() {

	for i := 1; i < 5; i++ {
		go func(j int) {
			Server(j)
		}(i)
	}

	Client()

}

func Server(count int) {
	conf := &config.ServiceConf{
		EtcdAuth:      config.EtcdAuth{},
		Schema:        "www.vector.com",
		ServerName:    "knowing",
		Endpoints:     Endpoints,
		ServerAddress: ":2000" + strconv.Itoa(count),
	}
	demo := &RegionHandlerServer{ServerAddress: conf.ServerAddress}
	rpcServer, err := grpcx.MustNewGrpcxServer(conf, func(server *grpc.Server) {
		proto.RegisterRegionHandlerServer(server, demo)
	})
	if err != nil {
		panic(err)
	}
	log.Fatal(rpcServer.Run())
}

type RegionHandlerServer struct {
	ServerAddress string
}

func (s *RegionHandlerServer) GetListenAudio(ctx context.Context, r *proto.FindRequest) (*proto.HasOptionResponse, error) {

	has := []*proto.HasOption(nil)
	for _, t := range r.Tokens {

		has = append(has, &proto.HasOption{Token: t + " FROM " + s.ServerAddress, Listen: 1})
	}
	res := &proto.HasOptionResponse{
		Items: has,
	}
	return res, nil
}

func Client() {
	conf := &config.ClientConf{
		EtcdAuth:  config.EtcdAuth{},
		Target:    "www.vector.com:///knowing",
		Endpoints: Endpoints,
		WithBlock: false,
	}

	r, err := grpcx.MustNewGrpcxClient(conf)
	if err != nil {
		panic(err)
	}
	conn, err := r.GetConnection()
	if err != nil {
		panic(err)
	}
	regionHandlerClient := proto.NewRegionHandlerClient(conn)
	for {
		res, err := regionHandlerClient.GetListenAudio(
			context.Background(),
			&proto.FindRequest{Tokens: []string{"a_"}},
		)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(res)
		time.Sleep(1 * time.Second)
	}
}
