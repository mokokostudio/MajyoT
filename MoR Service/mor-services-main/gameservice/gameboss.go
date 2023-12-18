package gameservice

import (
	com "gitlab.com/morbackend/mor_services/common"
	"gitlab.com/morbackend/mor_services/mpb"
)

type gameBoss struct {
	bossId          uint32
	rsc             *mpb.BossRsc
	totalHP         uint64
	hp              uint64
	winLoseHp       uint64
	atk             uint64
	atkGap          int64
	criRate         uint64
	criDmgRate      uint64
	hitRate         uint64
	dodgeRate       uint64
	dmgAddRate      uint64
	dmgReduceRate   uint64
	atkBuffRate     uint64
	defenseBuffRate uint64
	killedBy        uint64
}

func newGameBoss(bossRsc *mpb.BossRsc) *gameBoss {
	boss := &gameBoss{
		bossId: bossRsc.BossId,
		rsc:    bossRsc,
	}
	boss.calcBossAttrs()
	return boss
}

func (b *gameBoss) calcBossAttrs() {
	b.hp = b.rsc.Attrs.Hp
	b.hp = b.hp * (com.RateBase + b.rsc.Attrs.HpAddRate) / com.RateBase
	b.totalHP = b.hp
	b.winLoseHp = b.totalHP * b.rsc.WinDmgRate / com.RateBase
	b.atk = b.rsc.Attrs.Atk
	b.atk = b.atk * (com.RateBase + b.rsc.Attrs.AtkAddRate) / com.RateBase
	b.atkGap = b.rsc.Attrs.AtkGap
	b.atkGap = b.atkGap * com.RateBase / int64(com.RateBase+b.rsc.Attrs.AtkSpeedAddRate)
	b.criRate = b.rsc.Attrs.CriRate
	b.criDmgRate = b.rsc.Attrs.CriDmgRate
	b.hitRate = b.rsc.Attrs.HitRate
	b.dodgeRate = b.rsc.Attrs.DodgeRate
	b.dmgAddRate = b.rsc.Attrs.DmgAddRate
	b.dmgReduceRate = b.rsc.Attrs.DmgReduceRate
	b.atkBuffRate = b.rsc.Attrs.AtkBuffRate
	b.defenseBuffRate = b.rsc.Attrs.DefenseBuffRate
}

func (b *gameBoss) RoleType() mpb.ERole_Type { return mpb.ERole_RoleType_Boss }
func (b *gameBoss) HP() uint64               { return b.hp }
func (b *gameBoss) SetHP(hp uint64)          { b.hp = hp }
func (b *gameBoss) ATK() uint64              { return b.atk }
func (b *gameBoss) ATKGap() int64            { return b.atkGap }
func (b *gameBoss) CriRate() uint64          { return b.criRate }
func (b *gameBoss) CriDmgRate() uint64       { return b.criDmgRate }
func (b *gameBoss) HitRate() uint64          { return b.hitRate }
func (b *gameBoss) DodgeRate() uint64        { return b.dodgeRate }
func (b *gameBoss) DmgAddRate() uint64       { return b.dmgAddRate }
func (b *gameBoss) DmgReduceRate() uint64    { return b.dmgReduceRate }
func (b *gameBoss) ATKBuffRate() uint64      { return b.atkBuffRate }
func (b *gameBoss) DefenseBuffRate() uint64  { return b.defenseBuffRate }
