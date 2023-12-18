package httpgateway

import (
	"context"

	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	pb "github.com/golang/protobuf/proto"
	"github.com/oldjon/gutil/conv"
	"github.com/oldjon/gutil/env"
	gjwt "github.com/oldjon/gutil/jwt"
	gxgrpc "github.com/oldjon/gx/modules/grpc"
	gxhttp "github.com/oldjon/gx/modules/http"
	"github.com/oldjon/gx/service"
	"github.com/pkg/errors"
	com "gitlab.com/morbackend/mor_services/common"
	"gitlab.com/morbackend/mor_services/mpb"
	"gitlab.com/morbackend/mor_services/mpberr"
	"gitlab.com/morbackend/mor_services/util"
	etcd "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

var httpAESEncryptionKeyPairs = append(util.DefaultAESEncryptionKeyPairs, &util.AESEncryptionKeyPair{
	Index:   conv.Uint32ToString(1),
	Key:     []byte("ad#ASDoiwdasdUbn^asd!q1"),
	IV:      []byte("HUIA@asffs0908(&^@!cs"),
	Retired: false,
})

type HTTPGateway struct {
	name       string
	logger     *zap.Logger
	config     env.ModuleConfig
	mux        *http.ServeMux
	etcdClient *etcd.Client
	connMgr    *gxgrpc.ConnManager

	kvm              *service.KVMgr
	signingKey       []byte
	signingMethod    jwt.SigningMethod
	signingDuration  time.Duration
	protocolEncode   string
	isSandbox        bool
	enableEncryption bool
	// registrationLimiter *util.RedisLimiter

	// HTTPClient Client
	metrics *metrics

	// wfClient *MultiLangWordFilter // TODO

	// res          *gatewayRes //TODO
	centerRegion string
}

// NewHTTPGateway create a HTTPGateway entity
func NewHTTPGateway(driver service.ModuleDriver) (gxhttp.GXHttpHandler, error) {
	mux := http.NewServeMux()
	host := driver.Host()
	gateway := HTTPGateway{
		name:            driver.ModuleName(),
		logger:          driver.Logger(),
		config:          driver.ModuleConfig(),
		mux:             mux,
		etcdClient:      host.EtcdSession().Client(),
		kvm:             host.KVManager(),
		signingMethod:   jwt.SigningMethodHS256,
		signingDuration: 24 * time.Hour,
		metrics:         newGatewayMetrics(driver),
	}

	gateway.protocolEncode = gateway.config.GetString("protocol_code")
	gateway.centerRegion = gateway.config.GetString("center_region")
	gateway.isSandbox = gateway.config.GetBool("is_sandbox")
	gateway.enableEncryption = gateway.config.GetBool("enable_encryption")

	jm := gjwt.New(gjwt.Options{
		KeyGetter: func(token *jwt.Token) (interface{}, error) {
			return gateway.signingKey, nil
		},
		NewClaimsFunc: func() jwt.Claims {
			return &mpb.JWTClaims{}
		},
		SigningMethod: gateway.signingMethod,
	})
	eh := util.NewHTTPErrorHandler(driver.Logger())
	dialer := gxgrpc.Dialer{
		HostName:   driver.Host().Name(),
		EtcdClient: gateway.etcdClient,
		Logger:     gateway.logger,
		Tracer:     driver.Tracer(),
		EnableTLS:  gateway.config.GetBool("enable_tls"),
		CAFile:     gateway.config.GetString("ca_file"),
		CertFile:   gateway.config.GetString("cert_file"),
		KeyFile:    gateway.config.GetString("key_file"),
	}
	_ = jm
	gateway.connMgr = gxgrpc.NewConnManager(&dialer)

	// var err error
	// gateway.HTTPClient, err = wire.NewHTTPClient(host, wire.HTTPClientOptions{})
	// if err != nil {
	// 	return nil, err
	// }
	//
	// gateway.cache, err = wire.NewCacheClient(host, wire.CacheClientOptions{})
	// if err != nil {
	// 	return nil, err
	// }
	//
	// gateway.wfClient, err = util.GetMultiLangWordFilter(gateway.logger)
	// if err != nil {
	// 	return nil, err
	// }
	// gateway.res, err = newGatewayRes(gateway.logger, gateway.metrics)
	// if err != nil {
	// 	return nil, err
	// }
	//
	// if gateway.config.GetBool("registration_limit.open") {
	// 	bot, err := wire.NewRedisClient(host, wire.RedisClientOptions{ConfigKey: "registration_limit.gredis"})
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	gateway.registrationLimiter, err = util.NewRedisLimiter(bot, gateway.logger,
	// 		gateway.config.GetInt64("registration_limit.duration"),
	// 		gateway.config.GetInt64("registration_limit.cnt_per_dur"))
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	mux.Handle("/HelloWorld", eh.Handler(gateway.helloWorld))
	mux.Handle("/LoginTest", eh.Handler(gateway.loginTest))
	mux.Handle("/LoginByToken", eh.Handler(gateway.loginByToken))
	mux.Handle("/TelegramLogin", eh.Handler(gateway.telegramLogin))
	mux.Handle("/GetItems", jm.Handler(eh.Handler(gateway.getItems)))
	mux.Handle("/UpgradeBaseEquipStar", jm.Handler(eh.Handler(gateway.upgradeBaseEquipStar)))
	mux.Handle("/UpgradeBaseEquipLevel", jm.Handler(eh.Handler(gateway.upgradeBaseEquipLevel)))
	mux.Handle("/Fight", jm.Handler(eh.Handler(gateway.fight)))
	mux.Handle("/GetHiddenBoss", jm.Handler(eh.Handler(gateway.getHiddenBoss)))
	mux.Handle("/GetEnergy", jm.Handler(eh.Handler(gateway.getEnergy)))
	return &gateway, nil
}

func (hg *HTTPGateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger := hg.logger.With(zap.String("path", r.URL.Path))
	logger.Info("handling http")
	defer logger.Info("finish http")
	hg.mux.ServeHTTP(w, r)
}

func (hg *HTTPGateway) Serve(ctx context.Context) error {
	signingKey, err := hg.kvm.GetOrGenerate(ctx, com.JWTGatewayTokenKey, 32)
	if err != nil {
		return errors.WithStack(err)
	}
	hg.signingKey = signingKey

	<-ctx.Done()
	return ctx.Err()
}

func (hg *HTTPGateway) Logger() *zap.Logger {
	return hg.logger
}

func (hg *HTTPGateway) ConnMgr() *gxgrpc.ConnManager {
	return hg.connMgr
}

func (hg *HTTPGateway) Name() string {
	return hg.name
}

func (hg *HTTPGateway) readHTTPReq(w http.ResponseWriter, r *http.Request, msg interface{}) error {
	var err error

	//_, isLogin := msg.(*mpb.CReqWebLoginByWallet)
	isLogin := false
	//other login way
	options := util.HTTPEncryptionOptions{
		EnableEncryption:          hg.enableEncryption,
		AESEncryptionKeyPairs:     httpAESEncryptionKeyPairs,
		IsPlatformLoginMethodCall: isLogin,
	}

	if hg.protocolEncode == "json" ||
		r.Header.Get("Content-Type") == "application/json" {
		err = util.ReadHTTPJSONReq(w, r, msg, options)
	} else {
		pbMsg, ok := msg.(pb.Message)
		if !ok {
			return mpberr.ErrUnsupportedHTTPContentType
		}
		err = util.ReadHTTPReq(w, r, pbMsg, options)
	}
	if err != nil {
		hg.metrics.incReadHTTPFail(r.URL.Path, err)
	}
	return err
}

func (hg *HTTPGateway) writeHTTPRes(w http.ResponseWriter, msg pb.Message) error {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "content-type")
	if hg.protocolEncode == "json" {
		return util.WriteHTTPJSONRes(w, msg)
	}
	return util.WriteHTTPRes(w, msg)
}
