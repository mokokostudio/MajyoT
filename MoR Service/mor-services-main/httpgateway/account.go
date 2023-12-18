package httpgateway

import (
	"encoding/json"
	"net/http"
	"strconv"

	com "gitlab.com/morbackend/mor_services/common"
	"gitlab.com/morbackend/mor_services/mpb"
	"gitlab.com/morbackend/mor_services/mpberr"
	"gitlab.com/morbackend/mor_services/util"
	"go.uber.org/zap"
)

func (hg *HTTPGateway) helloWorld(w http.ResponseWriter, r *http.Request) error {
	ip := util.GetRemoteIPAddress(r)
	ctx := r.Context()

	//bys, _ := io.ReadAll(r.Body)
	//hg.logger.Debug("helloWorld payload", zap.String("", string(bys)))
	//return nil
	req := &mpb.TGMsgRecv{}
	err := hg.readHTTPReq(w, r, req)
	if err != nil {
		return err
	}
	hg.logger.Debug("helloWorld payload", zap.Any("", req), zap.String("remote_ip", ip))

	accClient, err := com.GetAccountServiceClient(ctx, hg)
	if err != nil {
		return err
	}
	if req.Message == nil || req.Message.From == nil {
		return mpberr.ErrParam
	}
	rpcRes, err := accClient.TelegramLogin(ctx, &mpb.ReqTelegramLogin{
		Account: req.Message.From.FirstName + req.Message.From.LastName,
	})
	if err != nil {
		return err
	}

	apiClient, err := com.GetAPIProxyGRPCClient(ctx, hg)
	if err != nil {
		return err
	}
	msg := &mpb.TGReply{
		Method: "sendMessage",
		ChatId: req.Message.Chat.Id,
		Text:   "Hello " + rpcRes.Account.Account + "! your uid is " + strconv.Itoa(int(rpcRes.Account.UserId)),
	}
	msgBytes, _ := json.Marshal(msg)
	_, err = apiClient.SendMsgToTelegram(ctx, &mpb.ReqSendMsgToTelegram{
		Bot: com.TGBot,
		Msg: msgBytes,
	})

	return err
}

func (hg *HTTPGateway) telegramLogin(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	req := &mpb.TGMsgRecv{}
	err := hg.readHTTPReq(w, r, req)
	if err != nil {
		return err
	}
	client, err := com.GetAccountServiceClient(ctx, hg)
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

	return hg.writeHTTPRes(w, &mpb.CResTelegramLogin{Account: rpcRes.Account})
}

func (hg *HTTPGateway) loginTest(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	req := &mpb.CReqLoginTest{}
	err := hg.readHTTPReq(w, r, req)
	if err != nil {
		return err
	}
	client, err := com.GetAccountServiceClient(ctx, hg)
	if err != nil {
		return err
	}

	rpcRes, err := client.LoginTest(ctx, &mpb.ReqUserId{
		UserId: req.UserId,
	})
	if err != nil {
		return err
	}
	res := &mpb.CResLoginTest{
		Account:        rpcRes.Account,
		Token:          rpcRes.Token,
		Energy:         rpcRes.Energy,
		EnergyUpdateAt: rpcRes.EnergyUpdateAt,
	}
	return hg.writeHTTPRes(w, res)
}

func (hg *HTTPGateway) loginByToken(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	req := &mpb.CReqLoginByToken{}
	err := hg.readHTTPReq(w, r, req)
	if err != nil {
		return err
	}
	client, err := com.GetAccountServiceClient(ctx, hg)
	if err != nil {
		return err
	}

	rpcRes, err := client.LoginByToken(ctx, &mpb.ReqLoginByToken{
		Token: req.Token,
	})
	if err != nil {
		return err
	}
	res := &mpb.CResLoginByToken{
		Account:        rpcRes.Account,
		Token:          rpcRes.Token,
		Energy:         rpcRes.Energy,
		EnergyUpdateAt: rpcRes.EnergyUpdateAt,
	}
	return hg.writeHTTPRes(w, res)
}
