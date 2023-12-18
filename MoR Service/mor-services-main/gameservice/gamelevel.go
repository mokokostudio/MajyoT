package gameservice

import (
	"github.com/oldjon/gutil"
	com "gitlab.com/morbackend/mor_services/common"
	"gitlab.com/morbackend/mor_services/mpb"
	"go.uber.org/zap"
	"math/rand"
)

const (
	FightMaxRound = 100000
)

type gameLevel struct {
	svc    *GameService
	logger *zap.Logger
}

func newGameLevel(svc *GameService) *gameLevel {
	return &gameLevel{
		svc:    svc,
		logger: svc.logger,
	}
}

func (gl *gameLevel) fight(player *gamePlayer, boss *gameBoss) (win bool, details []*mpb.FightDetail,
	dmg uint64, dmgRate uint64, bossDie bool) {
	var ms int64
	var round int64
	details = make([]*mpb.FightDetail, 0, 20)

	for player.hp > 0 && boss.hp > 0 && round < FightMaxRound {
		ms++
		if ms%player.atkGap == 0 {
			// do player atk
			detail := gl.atk(ms, player, boss)
			if detail != nil {
				details = append(details, detail)
			}
			round++
		}
		if boss.hp == 0 {
			boss.killedBy = player.userId
			break
		}
		if ms%boss.atkGap == 0 {
			// do boss atk
			detail := gl.atk(ms, boss, player)
			if detail != nil {
				details = append(details, detail)
			}
			round++
		}
	}

	dmg = boss.totalHP - boss.hp
	dmgRate = dmg * com.RateBase / boss.totalHP
	win = boss.hp <= (boss.totalHP - boss.winLoseHp)
	return win, details, dmg, dmgRate, boss.hp == 0
}

type roleI interface {
	RoleType() mpb.ERole_Type
	HP() uint64
	SetHP(uint642 uint64)
	ATK() uint64
	ATKGap() int64
	CriRate() uint64
	CriDmgRate() uint64
	HitRate() uint64
	DodgeRate() uint64
	DmgAddRate() uint64
	DmgReduceRate() uint64
	ATKBuffRate() uint64
	DefenseBuffRate() uint64
}

func (gl *gameLevel) atk(ms int64, from, to roleI) *mpb.FightDetail {
	hp := to.HP()
	if hp == 0 {
		return nil
	}
	detail := &mpb.FightDetail{
		AttackerType: from.RoleType(),
		BeAttacker:   to.RoleType(),
		HpBefore:     hp,
		HpAfter:      hp,
		AttackTime:   ms,
	}
	hitRate := gutil.Bound(
		gutil.If(from.HitRate() > to.DodgeRate(), from.HitRate()-to.DodgeRate(), 0),
		30, com.RateBase)
	if hitRate < uint64(rand.Int31n(com.RateBase)+1) { // miss
		detail.IsMiss = true
		return detail
	}

	dmg := from.ATK() *
		gutil.If(com.RateBase+from.DmgAddRate() > to.DmgReduceRate(),
			com.RateBase+from.DmgAddRate()-to.DmgReduceRate(), 0) *
		gutil.If(com.RateBase+from.ATKBuffRate() > to.DefenseBuffRate(),
			com.RateBase+from.ATKBuffRate()-to.DefenseBuffRate(), 0) /
		(com.RateBase * com.RateBase)

	detail.IsCri = from.CriRate() >= uint64(rand.Int31n(com.RateBase)+1)
	if detail.IsCri {
		dmg = dmg * (com.RateBase + from.CriDmgRate()) / com.RateBase
	}

	dmg = gutil.Min(dmg, hp)
	hp = hp - dmg
	to.SetHP(hp)
	detail.HpAfter = hp
	detail.Dmg = dmg
	return detail
}
