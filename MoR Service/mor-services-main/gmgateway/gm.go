package gmgateway

import (
	com "gitlab.com/morbackend/mor_services/common"
	"gitlab.com/morbackend/mor_services/mpb"
	"gitlab.com/morbackend/mor_services/mpberr"
	"net/http"
	"strings"
)

func (gg *GMGateway) adminLoginByPassword(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	req := &mpb.CReqAdminLoginByPassword{}
	err := gg.readHTTPReq(w, r, req)
	if err != nil {
		return err
	}
	if len(req.Account) < com.MinAdminAccountLength || len(req.Account) > com.MaxAdminAccountLength {
		return mpberr.ErrAdminAccountOrPasswd
	}

	req.Password = strings.ToLower(req.Password)

	if len(req.Password) != com.PasswordLen {
		return mpberr.ErrAdminAccountOrPasswd
	}

	client, err := com.GetGMServiceClient(ctx, gg)
	if err != nil {
		return err
	}
	rpcReq := mpb.ReqAdminLoginByPassword{
		Account:  req.Account,
		Password: req.Password,
	}
	res, err := client.AdminLoginByPassword(ctx, &rpcReq)
	if err != nil {
		return err
	}
	cres := &mpb.CResAdminLoginByPassword{
		Token: res.Token,
	}
	return gg.writeHTTPRes(w, cres)
}
