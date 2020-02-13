package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"google.golang.org/grpc"

	"github.com/holenat/smartdns/proto"

	"go.etcd.io/etcd/clientv3"
)

type record struct {
	Host string `json:"host"`
}

type ResolverConfig struct {
	Endpoints  []string `toml:"endpoints"`
	ListenAddr string   `toml:"listen"`
}

type Resolver struct {
	cfg *ResolverConfig
	cli *clientv3.Client
}

func NewResolve(cfg *ResolverConfig) (*Resolver, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   cfg.Endpoints,
		DialTimeout: time.Minute * 1,
	})
	if err != nil {
		return nil, err
	}

	return &Resolver{
		cfg: cfg,
		cli: cli,
	}, nil
}

func (r *Resolver) Run() error {
	lis, err := net.Listen("tcp", r.cfg.ListenAddr)
	if err != nil {
		return err
	}

	defer lis.Close()

	srv := grpc.NewServer()
	proto.RegisterSmartDNSServer(srv, r)

	return srv.Serve(lis)
}

func (r *Resolver) UpdateDomain(ctx context.Context, req *proto.UpdateDomainReq) (*proto.UpdateDomainReply, error) {
	domain := req.Domain
	ip := req.Ip

	sp := strings.Split(domain, ".")
	if len(sp) == 0 {
		return &proto.UpdateDomainReply{}, fmt.Errorf("invalid domain: %s", domain)
	}

	key := "/skydns"
	for i := len(sp) - 1; i >= 0; i-- {
		key = fmt.Sprintf("%s/%s", key, sp[i])
	}

	value := &record{Host: ip}
	b, err := json.Marshal(value)
	if err != nil {
		return &proto.UpdateDomainReply{}, err
	}

	log.Printf("put etcd, key %s value %s\n", key, string(b))

	_, err = r.cli.Put(context.Background(), key, string(b))
	return &proto.UpdateDomainReply{}, err
}
