package gameservice

import "gitlab.com/morbackend/mor_services/mpb"

func (svc *GameService) dbHiddenBoss2HiddenBoss(dbBoss *mpb.DBHiddenBoss) *mpb.HiddenBoss {
	if dbBoss == nil {
		return nil
	}
	return &mpb.HiddenBoss{
		BossId:   dbBoss.BossId,
		BossUuid: dbBoss.BossUuid,
		Hp:       dbBoss.Hp,
		Finder:   dbBoss.Finder,
		ExpireAt: dbBoss.ExpiredAt,
	}
}
