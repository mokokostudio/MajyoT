package nftservice

import (
	"context"
	"github.com/oldjon/gutil/env"
	"github.com/oldjon/gutil/gdb"
	gprotocol "github.com/oldjon/gutil/protocol"
	grmux "github.com/oldjon/gutil/redismutex"
	gxgrpc "github.com/oldjon/gx/modules/grpc"
	"github.com/oldjon/gx/service"
	"gitlab.com/morbackend/mor_services/mpb"
	"gitlab.com/morbackend/mor_services/util"
	etcd "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type NFTService struct {
	mpb.UnimplementedNFTServiceServer
	name        string
	logger      *zap.Logger
	config      env.ModuleConfig
	etcdClient  *etcd.Client
	host        service.Host
	connMgr     *gxgrpc.ConnManager
	rm          *NFTResourceMgr
	kvm         *service.KVMgr
	serverEnv   uint32
	sm          *util.ServiceMetrics
	dao         *nftDAO
	tcpMsgCoder gprotocol.FrameCoder
}

// NewNFTService create an NFTService entity
func NewNFTService(driver service.ModuleDriver) (gxgrpc.GRPCServer, error) {
	svc := &NFTService{
		name:       driver.ModuleName(),
		logger:     driver.Logger(),
		config:     driver.ModuleConfig(),
		etcdClient: driver.Host().EtcdSession().Client(),
		host:       driver.Host(),
		kvm:        driver.Host().KVManager(),
		sm:         util.NewServiceMetrics(driver),
	}

	dialer := gxgrpc.Dialer{
		HostName:   driver.Host().Name(),
		Tracer:     driver.Tracer(),
		EtcdClient: svc.etcdClient,
		Logger:     svc.logger,
		EnableTLS:  svc.config.GetBool("enable_tls"),
		CAFile:     svc.config.GetString("ca_file"),
		CertFile:   svc.config.GetString("cert_file"),
		KeyFile:    svc.config.GetString("key_file"),
	}
	svc.connMgr = gxgrpc.NewConnManager(&dialer)

	var err error
	svc.rm, err = newNFTResourceMgr(svc.logger, svc.sm)
	if err != nil {
		return nil, err
	}

	redisMux, err := grmux.NewRedisMux(svc.config.SubConfig("redis_mutex"), nil, svc.logger, driver.Tracer())
	if err != nil {
		return nil, err
	}

	nftRedis, err := gdb.NewRedisClientByConfig(svc.config.SubConfig("nft_redis"),
		svc.config.GetString("db_marshaller"), driver.Tracer())
	if err != nil {
		return nil, err
	}

	tranRedis, err := gdb.NewRedisClientByConfig(svc.config.SubConfig("tran_redis"),
		svc.config.GetString("db_marshaller"), driver.Tracer())
	if err != nil {
		return nil, err
	}

	tmpRedis, err := gdb.NewRedisClientByConfig(svc.config.SubConfig("tmp_redis"),
		svc.config.GetString("db_marshaller"), driver.Tracer())
	if err != nil {
		return nil, err
	}

	svc.dao = newNftDAO(svc, redisMux, nftRedis, tmpRedis, tranRedis)

	svc.serverEnv = uint32(svc.config.GetInt64("server_env"))
	svc.tcpMsgCoder = gprotocol.NewFrameCoder(svc.config.GetString("protocol_code"))

	return svc, err
}

func (svc *NFTService) Register(grpcServer *grpc.Server) {
	mpb.RegisterNFTServiceServer(grpcServer, svc)
}

func (svc *NFTService) Serve(ctx context.Context) error {
	<-ctx.Done()
	return ctx.Err()
}

func (svc *NFTService) Logger() *zap.Logger {
	return svc.logger
}

func (svc *NFTService) ConnMgr() *gxgrpc.ConnManager {
	return svc.connMgr
}

func (svc *NFTService) Name() string {
	return svc.name
}
