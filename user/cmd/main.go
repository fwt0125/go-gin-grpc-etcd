package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"net"
	"user/config"
	"user/discovery"
	"user/internal/handler"
	"user/internal/repository"
	"user/internal/service"
)

func main() {
	config.InitConfig()
	repository.InitDb()

	//etcd的地址
	etcdAddress := []string{viper.GetString("etcd.address")}
	//服务的注册
	EtcdRegister := discovery.NewRegister(etcdAddress, logrus.New())
	grpcAddress := viper.GetString("server.grpcAddress")
	serverNode := discovery.Server{
		Name: viper.GetString("server.domain"),
		Addr: grpcAddress,
	}
	server := grpc.NewServer()
	defer server.Stop()

	//绑定服务
	service.RegisterUserServiceServer(server, handler.NewUserService())

	listen, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		panic(err)
	}

	if _, err := EtcdRegister.Register(serverNode, 10); err != nil {
		panic(err)
	}

	if err := server.Serve(listen); err != nil {
		panic(err)
	}

}
