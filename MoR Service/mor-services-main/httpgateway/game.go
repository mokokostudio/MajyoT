package httpgateway

import (
	"net/http"

	com "gitlab.com/morbackend/mor_services/common"
	"gitlab.com/morbackend/mor_services/mpb"
	"gitlab.com/morbackend/mor_services/util"
)

func (hg *HTTPGateway) fight(w http.ResponseWriter, r *http.Request) error {
	claim, ctx, err := util.ClaimFromContext(r.Context())
	if err != nil {
		return err
	}
	req := &mpb.CReqFight{}
	err = hg.readHTTPReq(w, r, req)
	if err != nil {
		return err
	}
	client, err := com.GetGameServiceClient(ctx, hg)
	if err != nil {
		return err
	}

	rpcRes, err := client.Fight(ctx, &mpb.ReqFight{
		UserId:   claim.UserId,
		BossId:   req.BossId,
		BossUuid: req.BossUuid,
	})
	if err != nil {
		return err
	}

	res := &mpb.CResFight{
		Win:        rpcRes.Win,
		Details:    rpcRes.Details,
		Awards:     rpcRes.Awards,
		EnergyCost: rpcRes.EnergyCost,
		Dmg:        rpcRes.Dmg,
		DmgRate:    rpcRes.DmgRate,
		HiddenBoss: rpcRes.HiddenBoss,
	}
	return hg.writeHTTPRes(w, res)
}

func (hg *HTTPGateway) getHiddenBoss(w http.ResponseWriter, r *http.Request) error {
	claim, ctx, err := util.ClaimFromContext(r.Context())
	if err != nil {
		return err
	}
	req := &mpb.CReqGetHiddenBoss{}
	err = hg.readHTTPReq(w, r, req)
	if err != nil {
		return err
	}
	client, err := com.GetGameServiceClient(ctx, hg)
	if err != nil {
		return err
	}

	rpcRes, err := client.GetHiddenBoss(ctx, &mpb.ReqGetHiddenBoss{
		UserId:   claim.UserId,
		BossUuid: req.BossUuid,
	})
	if err != nil {
		return err
	}

	res := &mpb.CResGetHiddenBoss{
		HiddenBoss: rpcRes.HiddenBoss,
		Fought:     rpcRes.Fought,
		FightCd:    rpcRes.FightCd,
	}
	return hg.writeHTTPRes(w, res)
}

func (hg *HTTPGateway) getEnergy(w http.ResponseWriter, r *http.Request) error {
	claim, ctx, err := util.ClaimFromContext(r.Context())
	if err != nil {
		return err
	}
	client, err := com.GetGameServiceClient(ctx, hg)
	if err != nil {
		return err
	}

	rpcRes, err := client.GetEnergy(ctx, &mpb.ReqUserId{
		UserId: claim.UserId,
	})
	if err != nil {
		return err
	}

	res := &mpb.CResGetEnergy{
		Energy:   rpcRes.Energy,
		UpdateAt: rpcRes.UpdateAt,
	}
	return hg.writeHTTPRes(w, res)
}
