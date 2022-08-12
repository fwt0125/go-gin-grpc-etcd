package discovery

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
	"strings"
	"time"
)

type Register struct {
	EtcdAddress []string
	DialTimeout int

	closeCh     chan struct{}
	leasesId    clientv3.LeaseID
	keepAliveCh <-chan *clientv3.LeaseKeepAliveResponse

	serverInfo Server
	serverTTL  int64
	cli        *clientv3.Client
	logger     *logrus.Logger
}

// 基于ETCD创建一个register
func NewRegister(etcdAddress []string, logger *logrus.Logger) *Register {
	return &Register{
		EtcdAddress: etcdAddress,
		DialTimeout: 3,
		logger:      logger,
	}
}

func (r *Register) Register(srvInfo Server, ttl int64) (chan<- struct{}, error) {
	var err error
	if strings.Split(srvInfo.Addr, ":")[0] == "" {
		return nil, errors.New("invalid ip address")
	}

	r.cli, err = clientv3.New(clientv3.Config{
		Endpoints:            r.EtcdAddress,
		AutoSyncInterval:     0,
		DialTimeout:          time.Duration(r.DialTimeout) * time.Second,
		DialKeepAliveTime:    0,
		DialKeepAliveTimeout: 0,
		MaxCallSendMsgSize:   0,
		MaxCallRecvMsgSize:   0,
		TLS:                  nil,
		Username:             "",
		Password:             "",
		RejectOldCluster:     false,
		DialOptions:          nil,
		Context:              nil,
		Logger:               nil,
		LogConfig:            nil,
		PermitWithoutStream:  false,
	})
	if err != nil {
		return nil, err
	}

	r.serverInfo = srvInfo
	r.serverTTL = ttl
	if err = r.register(); err != nil {
		return nil, err
	}

	r.closeCh = make(chan struct{})
	go r.keepAlive()
	return r.closeCh, nil
}

func (r *Register) register() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.DialTimeout)*time.Second)
	defer cancel()

	leaseResp, err := r.cli.Grant(ctx, r.serverTTL)
	if err != nil {
		return err
	}

	r.leasesId = leaseResp.ID
	r.keepAliveCh, err = r.cli.KeepAlive(context.Background(), r.leasesId)
	if err != nil {
		return err
	}
	marshal, err := json.Marshal(r.serverInfo)
	if err != nil {
		return err
	}
	_, err = r.cli.Put(context.Background(), BuildRegisterPath(r.serverInfo), string(marshal), clientv3.WithLease(r.leasesId))
	return err
}

func (r *Register) keepAlive() error {
	ticker := time.NewTicker(time.Duration(r.serverTTL) * time.Second)
	for {
		select {
		case <-r.closeCh:
			if err := r.unRegister(); err != nil {
				fmt.Println("unregister failed error", err)
			}
			_, err := r.cli.Revoke(context.Background(), r.leasesId)
			if err != nil {
				fmt.Println("revoke fail")
			}
		case res := <-r.keepAliveCh:
			if res != nil {
				if err := r.register(); err != nil {
					fmt.Println("register err")
				}
			}
		case <-ticker.C:
			if r.keepAliveCh != nil {
				if err := r.register(); err != nil {
					fmt.Println("register ticker err")
				}
			}
		}
	}
}

func (r *Register) unRegister() error {
	_, err := r.cli.Delete(context.Background(), BuildRegisterPath(r.serverInfo))
	return err
}
