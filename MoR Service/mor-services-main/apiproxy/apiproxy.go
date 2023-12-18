package apiproxy

import (
	"context"
	"net/http"

	pb "github.com/golang/protobuf/proto"
	"github.com/oldjon/gutil/conv"
	"github.com/oldjon/gutil/env"
	gxgrpc "github.com/oldjon/gx/modules/grpc"
	gxhttp "github.com/oldjon/gx/modules/http"
	"github.com/oldjon/gx/service"
	com "gitlab.com/morbackend/mor_services/common"
	"gitlab.com/morbackend/mor_services/mpb"
	"gitlab.com/morbackend/mor_services/mpberr"
	"gitlab.com/morbackend/mor_services/util"
	etcd "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

var httpAESEncryptionKeyPairs = append(util.DefaultAESEncryptionKeyPairs, &util.AESEncryptionKeyPair{ // nolint:unused
	Index:   conv.Uint32ToString(1),
	Key:     []byte("ad#oiwUbn^asd!q1"),
	IV:      []byte("HUI@as0908(&^@!cs"),
	Retired: false,
})

var apiproxy *APIProxy

type APIProxy struct {
	name       string
	logger     *zap.Logger
	config     env.ModuleConfig
	mux        *http.ServeMux
	etcdClient *etcd.Client
	connMgr    *gxgrpc.ConnManager

	kvm              *service.KVMgr
	protocolEncode   string
	isSandbox        bool
	enableEncryption bool
	centerRegion     string
	// HTTPClient Client
	metrics *metrics
}

// NewAPIProxy create an apiproxy entity
func NewAPIProxy(driver service.ModuleDriver) (gxhttp.GXHttpHandler, error) {
	mux := http.NewServeMux()
	host := driver.Host()
	ap := APIProxy{
		name:       driver.ModuleName(),
		logger:     driver.Logger(),
		config:     driver.ModuleConfig(),
		mux:        mux,
		etcdClient: host.EtcdSession().Client(),
		kvm:        host.KVManager(),
		metrics:    newMetrics(driver),
	}

	ap.protocolEncode = ap.config.GetString("protocol_code")
	ap.centerRegion = ap.config.GetString("center_region")
	ap.isSandbox = ap.config.GetBool("is_sandbox")
	ap.enableEncryption = ap.config.GetBool("enable_encryption")

	eh := util.NewHTTPErrorHandler(driver.Logger())
	dialer := gxgrpc.Dialer{
		HostName:   driver.Host().Name(),
		EtcdClient: ap.etcdClient,
		Logger:     ap.logger,
		Tracer:     driver.Tracer(),
		EnableTLS:  ap.config.GetBool("enable_tls"),
		CAFile:     ap.config.GetString("ca_file"),
		CertFile:   ap.config.GetString("cert_file"),
		KeyFile:    ap.config.GetString("key_file"),
	}
	ap.connMgr = gxgrpc.NewConnManager(&dialer)

	mux.Handle("/HelloWorld", eh.APIHandler(ap.helloWorld))
	mux.Handle("/TelegramLogin", eh.APIHandler(ap.telegramLogin))
	mux.Handle("/MORAPIForTelegram", eh.APIHandler(ap.morAPIForTelegram))

	apiproxy = &ap
	return &ap, nil
}

func (ap *APIProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger := ap.logger.With(zap.String("path", r.URL.Path))
	logger.Info("handling http")
	defer logger.Info("finish http")
	ap.mux.ServeHTTP(w, r)
}

func (ap *APIProxy) Serve(ctx context.Context) error {
	<-ctx.Done()
	return ctx.Err()
}

func (ap *APIProxy) Logger() *zap.Logger {
	return ap.logger
}

func (ap *APIProxy) ConnMgr() *gxgrpc.ConnManager {
	return ap.connMgr
}

func (ap *APIProxy) Name() string {
	return "apiproxygrpc"
}

func (ap *APIProxy) readHTTPReq(w http.ResponseWriter, r *http.Request, msg interface{}) error { // nolint:unused
	var err error

	//_, isLogin := msg.(*mpb.CReqLogin)// TODO

	options := util.HTTPEncryptionOptions{
		EnableEncryption:      ap.enableEncryption,
		AESEncryptionKeyPairs: httpAESEncryptionKeyPairs,
		//IsPlatformLoginMethodCall: isLogin,
	}

	if ap.protocolEncode == "json" ||
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
		ap.metrics.incReadHTTPFail(r.URL.Path, err)
	}
	return err
}

func (ap *APIProxy) writeHTTPRes(w http.ResponseWriter, msg pb.Message) error { // nolint:unused
	if ap.protocolEncode == "json" {
		return util.WriteHTTPJSONRes(w, msg)
	}
	return util.WriteHTTPRes(w, msg)
}

func (ap *APIProxy) helloWorld(w http.ResponseWriter, r *http.Request) error {
	//_, err := w.Write([]byte("hello world"))

	return APIProxyGRPCGetMe().sendEmailBindCode("lrunwow@gmail.com", "123456")
}

func (ap *APIProxy) telegramLogin(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	req := &mpb.TGMsgRecv{}
	err := ap.readHTTPReq(w, r, req)
	if err != nil {
		return err
	}
	client, err := com.GetAccountServiceClient(ctx, ap)
	if err != nil {
		return err
	}

	if req.Message == nil || req.Message.From == nil {
		return mpberr.ErrParam
	}

	rpcRes, err := client.TelegramLogin(ctx, &mpb.ReqTelegramLogin{
		Account: req.Message.From.FirstName + req.Message.From.LastName,
	})
	if err != nil {
		return err
	}

	return ap.writeHTTPRes(w, &mpb.CResTelegramLogin{Account: rpcRes.Account})
}

func (ap *APIProxy) morAPIForTelegram(w http.ResponseWriter, r *http.Request) error {
	ip := util.GetRemoteIPAddress(r)
	ctx := r.Context()
	//bys, _ := io.ReadAll(r.Body)
	//ap.logger.Debug("helloWorld payload", zap.String("", string(bys)))
	req := &mpb.TGMsgRecv{}
	err := ap.readHTTPReq(w, r, req)
	if err != nil {
		ap.logger.Error("morAPIForTelegram read req failed", zap.Error(err))
		return err
	}
	ap.logger.Debug("morAPIForTelegram payload", zap.Any("", req), zap.String("remote_ip", ip))

	if req.Message != nil { // for cmd or normal message
		var text = req.Message.Text
		if len(text) == 0 {
			return nil
		}
		if text[0] != '/' {
			fn, _ := TelegramManagerGetMe().getCmdHandler(com.TGCmd_Echo)
			return fn(ctx, w, req)
		}

		cmd, cmdParams, at := util.ReadTGCmdFromMsgText(text)
		_, _ = cmdParams, at
		_ = cmd
		fn, ok := TelegramManagerGetMe().getCmdHandler(cmd)
		if !ok {
			ap.logger.Error("morAPIForTelegram api not found", zap.String("cmd", cmd))
			return nil
		}
		return fn(ctx, w, req)
	}

	if req.CallbackQuery != nil { // for callback query
		if req.CallbackQuery.Data != "" { // for callback data
			var text = req.CallbackQuery.Data
			if text[0] != '/' {
				fn, _ := TelegramManagerGetMe().getCmdHandler(com.TGCBQ_Echo)
				return fn(ctx, w, req)
			}

			cmd, cmdParams, at := util.ReadTGCmdFromMsgText(text)
			_, _ = cmdParams, at
			_ = cmd
			fn, ok := TelegramManagerGetMe().getCBQHandler(cmd)
			if !ok {
				ap.logger.Error("morAPIForTelegram api not found", zap.String("cmd", cmd))
				return nil
			}
			return fn(ctx, w, req)
		}
		if req.CallbackQuery.GameShortName != "" { // for play game 			// TODO try launch the game https://core.telegram.org/bots/games
			fn, ok := TelegramManagerGetMe().getCBQHandler(com.TGCBQ_LaunchGame)
			if !ok {
				return mpberr.ErrParam
			}
			return fn(ctx, w, req)
		}
	}
	return nil
}
