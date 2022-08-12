package main

import (
	"api-getway/config"
	"api-getway/discovery"
	"api-getway/internal/service"
	"api-getway/routers"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	config.InitConfig()
	etcdAddress := []string{viper.GetString("etcd.address")}
	etcdRegister := discovery.NewResolver(etcdAddress, logrus.New())
	resolver.Register(etcdRegister)
	go startListen()
	{
		osSignal := make(chan os.Signal, 1)
		signal.Notify(osSignal, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
		s := <-osSignal
		fmt.Println("exit!", s)
	}
	fmt.Println("gateway listen on : ", viper.GetString("server.port"))
}

func startListen() {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	fmt.Println(viper.GetString("user.address"))
	userConn, err := grpc.Dial(viper.GetString("user.address"), opts...)
	if err != nil {
		panic("grpc " + viper.GetString("user.address") + " error")
	}

	userService := service.NewUserServiceClient(userConn)
	ginRouter := routers.NewRouter(userService)
	server := &http.Server{
		Addr:           viper.GetString("server.port"),
		Handler:        ginRouter,
		ReadTimeout:    time.Second * 10,
		WriteTimeout:   time.Second * 10,
		MaxHeaderBytes: 1 << 20,
	}

	err = server.ListenAndServe()
	if err != nil {
		fmt.Println("绑定失败，端口可能被占用:" + viper.GetString("server.port"))
		fmt.Println(err)
	}
}
