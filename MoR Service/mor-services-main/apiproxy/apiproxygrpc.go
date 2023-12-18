package apiproxy

import (
	"context"
	"time"

	"github.com/oldjon/gutil/env"
	gxgrpc "github.com/oldjon/gx/modules/grpc"
	"github.com/oldjon/gx/service"
	"gitlab.com/morbackend/mor_services/mpb"
	"gitlab.com/morbackend/mor_services/util"
	etcd "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var apiproxygrpc *APIProxyGRPCService

func APIProxyGRPCGetMe() *APIProxyGRPCService { // nolint:unused
	if apiproxygrpc == nil {
		for i := 0; i < 60; i++ { //spinning for 60s
			time.Sleep(time.Second)
			if apiproxygrpc != nil {
				return apiproxygrpc
			}
		}
		panic("apiproxygrpc not initialize")
	}
	return apiproxygrpc
}

type APIProxyGRPCService struct {
	mpb.UnimplementedAPIProxyGRPCServer
	name           string
	logger         *zap.Logger
	config         env.ModuleConfig
	etcdClient     *etcd.Client
	connMgr        *gxgrpc.ConnManager // 转发至其他gateway的消息，需要通过grpc
	host           service.Host
	sm             *util.ServiceMetrics
	rm             *apiProxyResourceMgr
	tgMgr          *TelegramManager
	emailSendIndex uint32
}

// NewAPIProxyGRPCService create a APIProxyGRPCService entity
func NewAPIProxyGRPCService(driver service.ModuleDriver) (gxgrpc.GRPCServer, error) {
	as := &APIProxyGRPCService{
		name:       driver.ModuleName(),
		logger:     driver.Logger(),
		config:     driver.ModuleConfig(),
		etcdClient: driver.Host().EtcdSession().Client(),
		host:       driver.Host(),
		sm:         util.NewServiceMetrics(driver),
	}

	dialer := gxgrpc.Dialer{
		HostName:   driver.Host().Name(),
		EtcdClient: as.etcdClient,
		Logger:     as.logger,
		Tracer:     driver.Tracer(),
		EnableTLS:  as.config.GetBool("enable_tls"),
		CAFile:     as.config.GetString("ca_file"),
		CertFile:   as.config.GetString("cert_file"),
		KeyFile:    as.config.GetString("key_file"),
	}
	as.connMgr = gxgrpc.NewConnManager(&dialer)
	var err error
	as.rm, err = newAPIProxyResourceMgr(as.logger, nil)
	if err != nil {
		return nil, err
	}
	as.tgMgr = newTelegramManager(as, as.config.GetString("telegram_send_url"),
		as.config.GetStringMapString("bot_tokens"))
	as.logger.Info("apiproxy grpc service start success")
	apiproxygrpc = as
	return as, nil
}

func (svc *APIProxyGRPCService) Register(grpcServer *grpc.Server) {
	mpb.RegisterAPIProxyGRPCServer(grpcServer, svc)
}

func (svc *APIProxyGRPCService) Serve(ctx context.Context) error {
	<-ctx.Done()
	return ctx.Err()
}

func (svc *APIProxyGRPCService) SendEmailBindCode(ctx context.Context, req *mpb.ReqSendEmailBindCode) (*mpb.Empty, error) {
	err := svc.sendEmailBindCode(req.Email, req.Code)
	if err != nil {
		return nil, err
	}
	return &mpb.Empty{}, nil
}

func (svc *APIProxyGRPCService) SendEmailResetPasswordValidationCode(ctx context.Context, req *mpb.ReqSendEmailResetPasswordValidationCode) (*mpb.Empty, error) {
	err := svc.sendEmailResetPasswordValidationCode(req.Email, req.Code)
	if err != nil {
		return nil, err
	}
	return &mpb.Empty{}, nil
}

func (svc *APIProxyGRPCService) SendMsgToTelegram(ctx context.Context, req *mpb.ReqSendMsgToTelegram) (*mpb.Empty, error) {
	err := svc.tgMgr.sendMsgToTelegram(ctx, req.Bot, req.Msg)
	if err != nil {
		return nil, err
	}
	return &mpb.Empty{}, nil
}

func (svc *APIProxyGRPCService) Logger() *zap.Logger {
	return svc.logger
}

func (svc *APIProxyGRPCService) ConnMgr() *gxgrpc.ConnManager {
	return svc.connMgr
}

func (svc *APIProxyGRPCService) Name() string {
	return svc.name
}
