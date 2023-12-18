package httpgateway

import (
	"net/http"

	com "gitlab.com/morbackend/mor_services/common"
	"gitlab.com/morbackend/mor_services/mpb"
	"gitlab.com/morbackend/mor_services/util"
)

func (hg *HTTPGateway) getItems(w http.ResponseWriter, r *http.Request) error {
	claim, ctx, err := util.ClaimFromContext(r.Context())
	if err != nil {
		return err
	}
	rpcReq := &mpb.ReqUserId{
		UserId: claim.UserId,
	}
	client, err := com.GetItemServiceClient(ctx, hg)
	if err != nil {
		return err
	}
	rpcRes, err := client.GetItems(ctx, rpcReq)
	if err != nil {
		return err
	}
	res := &mpb.CResGetItems{
		Items:      rpcRes.Items,
		BaseEquips: rpcRes.BaseEquips,
		NftEquips:  rpcRes.NftEquips,
		Mana:       rpcRes.Mana,
	}
	return hg.writeHTTPRes(w, res)
}

func (hg *HTTPGateway) upgradeBaseEquipStar(w http.ResponseWriter, r *http.Request) error {
	claim, ctx, err := util.ClaimFromContext(r.Context())
	if err != nil {
		return err
	}

	req := &mpb.CReqUpgradeBaseEquipStar{}
	err = hg.readHTTPReq(w, r, req)
	if err != nil {
		return err
	}

	rpcReq := &mpb.ReqUpgradeBaseEquipStar{
		UserId:     claim.UserId,
		EquiptType: req.EquiptType,
		CurStar:    req.CurStar,
	}
	client, err := com.GetItemServiceClient(ctx, hg)
	if err != nil {
		return err
	}
	rpcRes, err := client.UpgradeBaseEquipStar(ctx, rpcReq)
	if err != nil {
		return err
	}
	res := &mpb.CResUpgradeBaseEquipStar{
		Success: rpcRes.Success,
		NewStar: rpcRes.NewStar,
	}
	return hg.writeHTTPRes(w, res)
}

func (hg *HTTPGateway) upgradeBaseEquipLevel(w http.ResponseWriter, r *http.Request) error {
	claim, ctx, err := util.ClaimFromContext(r.Context())
	if err != nil {
		return err
	}

	req := &mpb.CReqUpgradeBaseEquipLevel{}
	err = hg.readHTTPReq(w, r, req)
	if err != nil {
		return err
	}

	rpcReq := &mpb.ReqUpgradeBaseEquipLevel{
		UserId:     claim.UserId,
		EquiptType: req.EquiptType,
		CurLevel:   req.CurLevel,
	}
	client, err := com.GetItemServiceClient(ctx, hg)
	if err != nil {
		return err
	}
	rpcRes, err := client.UpgradeBaseEquipLevel(ctx, rpcReq)
	if err != nil {
		return err
	}
	res := &mpb.CResUpgradeBaseEquipLevel{
		Success:  rpcRes.Success,
		NewLevel: rpcRes.NewLevel,
	}
	return hg.writeHTTPRes(w, res)
}
