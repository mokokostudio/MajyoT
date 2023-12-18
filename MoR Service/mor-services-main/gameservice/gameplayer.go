package gameservice

import (
	com "gitlab.com/morbackend/mor_services/common"
	"gitlab.com/morbackend/mor_services/mpb"
)

type gamePlayer struct {
	userId          uint64
	baseEquips      []*mpb.BaseEquip
	nftEquips       []*mpb.NFTEquip
	attrs           *mpb.Attrs
	totalHP         uint64
	hp              uint64
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
}

func newPlayer(svc *GameService, userId uint64) *gamePlayer {
	player := &gamePlayer{
		userId: userId,
		attrs:  svc.rm.getPlayerInitAttrs(),
	}
	return player
}

func (p *gamePlayer) updateEquips(baseEquips []*mpb.BaseEquip, nftEquips []*mpb.NFTEquip) {
	p.baseEquips = baseEquips
	p.nftEquips = nftEquips
	p.calcPlayerAttrs()
}

func (p *gamePlayer) calcPlayerAttrs() {
	for _, v := range p.baseEquips {
		p.attrs = com.AddAttrs(p.attrs, v.Attrs)
	}
	for _, v := range p.nftEquips {
		p.attrs = com.AddAttrs(p.attrs, v.Attrs)
	}
	p.hp = p.attrs.Hp
	p.hp = p.hp * (com.RateBase + p.attrs.HpAddRate) / com.RateBase
	p.totalHP = p.hp
	p.atk = p.attrs.Atk
	p.atk = p.atk * (com.RateBase + p.attrs.AtkAddRate) / com.RateBase
	p.atkGap = p.attrs.AtkGap
	p.atkGap = p.atkGap * com.RateBase / int64(com.RateBase+p.attrs.AtkSpeedAddRate)
	p.criRate = p.attrs.CriRate
	p.criDmgRate = p.attrs.CriDmgRate
	p.hitRate = p.attrs.HitRate
	p.dodgeRate = p.attrs.DodgeRate
	p.dmgAddRate = p.attrs.DmgAddRate
	p.dmgReduceRate = p.attrs.DmgReduceRate
	p.atkBuffRate = p.attrs.AtkBuffRate
	p.defenseBuffRate = p.attrs.DefenseBuffRate
}

func (p *gamePlayer) RoleType() mpb.ERole_Type { return mpb.ERole_RoleType_Player }
func (p *gamePlayer) HP() uint64               { return p.hp }
func (p *gamePlayer) SetHP(hp uint64)          { p.hp = hp }
func (p *gamePlayer) ATK() uint64              { return p.atk }
func (p *gamePlayer) ATKGap() int64            { return p.atkGap }
func (p *gamePlayer) CriRate() uint64          { return p.criRate }
func (p *gamePlayer) CriDmgRate() uint64       { return p.criDmgRate }
func (p *gamePlayer) HitRate() uint64          { return p.hitRate }
func (p *gamePlayer) DodgeRate() uint64        { return p.dodgeRate }
func (p *gamePlayer) DmgAddRate() uint64       { return p.dmgAddRate }
func (p *gamePlayer) DmgReduceRate() uint64    { return p.dmgReduceRate }
func (p *gamePlayer) ATKBuffRate() uint64      { return p.atkBuffRate }
func (p *gamePlayer) DefenseBuffRate() uint64  { return p.defenseBuffRate }
