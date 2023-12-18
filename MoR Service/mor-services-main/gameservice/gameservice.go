package gameservice

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/oldjon/gutil/env"
	"github.com/oldjon/gutil/gdb"
	gprotocol "github.com/oldjon/gutil/protocol"
	grmux "github.com/oldjon/gutil/redismutex"
	gxgrpc "github.com/oldjon/gx/modules/grpc"
	"github.com/oldjon/gx/service"
	com "gitlab.com/morbackend/mor_services/common"
	"gitlab.com/morbackend/mor_services/mpb"
	"gitlab.com/morbackend/mor_services/mpberr"
	"gitlab.com/morbackend/mor_services/util"
	etcd "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type GameService struct {
	mpb.UnimplementedGameServiceServer
	name            string
	logger          *zap.Logger
	config          env.ModuleConfig
	etcdClient      *etcd.Client
	host            service.Host
	connMgr         *gxgrpc.ConnManager
	tcpMsgCoder     gprotocol.FrameCoder
	signingMethod   jwt.SigningMethod
	signingDuration time.Duration
	rm              *gameResourceMgr
	kvm             *service.KVMgr
	serverEnv       uint32
	sm              *util.ServiceMetrics
	dao             *gameDAO
	gl              *gameLevel
	bossUUIDSF      service.Snowflake
}

// NewGameService create a GameService entity
func NewGameService(driver service.ModuleDriver) (gxgrpc.GRPCServer, error) {
	svc := &GameService{
		name:            driver.ModuleName(),
		logger:          driver.Logger(),
		config:          driver.ModuleConfig(),
		etcdClient:      driver.Host().EtcdSession().Client(),
		host:            driver.Host(),
		kvm:             driver.Host().KVManager(),
		sm:              util.NewServiceMetrics(driver),
		signingMethod:   jwt.SigningMethodHS256,
		signingDuration: 24 * 30 * time.Hour,
	}

	dialer := gxgrpc.Dialer{
		HostName:   driver.Host().Name(),
		Tracer:     driver.Tracer(),
		EtcdClient: svc.etcdClient,
		Logger:     svc.logger,
		EnableTLS:  svc.config.GetBool("enable_tls"),
		CAFile:     svc.config.GetString("ca_file"),
		CertFile:   svc.config.GetString("cert_file"),
		KeyFile:    svc.config.GetString("key_file"),
	}
	svc.connMgr = gxgrpc.NewConnManager(&dialer)
	svc.gl = newGameLevel(svc)

	var err error
	svc.rm, err = newGameResourceMgr(svc.logger, svc.sm)
	if err != nil {
		return nil, err
	}

	redisMux, err := grmux.NewRedisMux(svc.config.SubConfig("redis_mutex"), nil, svc.logger, driver.Tracer())
	if err != nil {
		return nil, err
	}

	gameRedis, err := gdb.NewRedisClientByConfig(svc.config.SubConfig("game_redis"),
		svc.config.GetString("db_marshaller"), driver.Tracer())
	if err != nil {
		return nil, err
	}

	tmpRedis, err := gdb.NewRedisClientByConfig(svc.config.SubConfig("tmp_redis"),
		svc.config.GetString("db_marshaller"), driver.Tracer())
	if err != nil {
		return nil, err
	}

	svc.dao = newGameDAO(svc, redisMux, gameRedis, tmpRedis)

	svc.serverEnv = uint32(svc.config.GetInt64("server_env"))
	svc.tcpMsgCoder = gprotocol.NewFrameCoder(svc.config.GetString("protocol_code"))

	return svc, err
}

func (svc *GameService) Register(grpcServer *grpc.Server) {
	mpb.RegisterGameServiceServer(grpcServer, svc)
}

func (svc *GameService) Serve(ctx context.Context) error {
	var err error
	svc.bossUUIDSF, err = svc.host.Snowflake(ctx, com.SnowflakeBossUUID, service.SnowflakeType_Default)
	if err != nil {
		svc.logger.Error("failed to create snowflake", zap.Error(err))
		return err
	}

	<-ctx.Done()
	return ctx.Err()
}

func (svc *GameService) Logger() *zap.Logger {
	return svc.logger
}

func (svc *GameService) ConnMgr() *gxgrpc.ConnManager {
	return svc.connMgr
}

func (svc *GameService) Name() string {
	return svc.name
}

func (svc *GameService) Fight(ctx context.Context, req *mpb.ReqFight) (*mpb.ResFight, error) {
	// check energy
	nowUnix := time.Now().Unix()

	bossRsc := svc.rm.getBossRsc(req.BossId)
	if bossRsc == nil {
		return nil, mpberr.ErrParam
	}

	if bossRsc.PreBoss > 0 { // check pre boss
		fightHis, err := svc.dao.getFightHistory(ctx, req.UserId)
		if err != nil {
			return nil, err
		}
		if fightHis.WinTimes[bossRsc.PreBoss] == 0 {
			return nil, mpberr.ErrParam
		}
	}

	var res *mpb.ResFight
	var err error

	if bossRsc.BossType != mpb.ERole_BossType_Hidden {
		res, err = svc.fightBoss(ctx, req, bossRsc, nowUnix)
		if err != nil {
			return nil, err
		}
	} else { // normal boss or nft boss
		if req.BossUuid == 0 {
			return nil, mpberr.ErrParam
		}
		res, err = svc.fightHiddenBoss(ctx, req, bossRsc, nowUnix)
		if err != nil {
			return nil, err
		}
	}

	if !res.Win {
		return res, nil
	}

	// trigger hidden boss
	res.HiddenBoss, err = svc.triggerHiddenBoss(ctx, req.UserId, req.BossId, nowUnix)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (svc *GameService) fightBoss(ctx context.Context, req *mpb.ReqFight, bossRsc *mpb.BossRsc, nowUnix int64) (
	*mpb.ResFight, error) {
	dbEnergy, err := svc.dao.getEnergy(ctx, req.UserId, nowUnix)
	if err != nil {
		return nil, err
	}
	if dbEnergy.Energy < bossRsc.EnergyCost {
		return nil, mpberr.ErrEnergyNotEnough
	}

	equips, err := svc.getUserEquips(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	if len(bossRsc.NftEquips) > 0 {
		var totalLevel uint32
		for _, equipId := range bossRsc.NftEquips {
			equip := equips.NftEquips[equipId]
			if equip == nil {
				return nil, mpberr.ErrNFTWeaponNotEquipped
			}
			totalLevel += equip.Level
		}
		if totalLevel < bossRsc.NftEquipsLevel {
			return nil, mpberr.ErrBaseEquipMaxLevel
		}
	}

	boss := newGameBoss(bossRsc)
	player := newPlayer(svc, req.UserId)
	player.updateEquips(equips.BaseEquips, equips.NftEquips)

	res := &mpb.ResFight{}
	res.Win, res.Details, res.Dmg, res.DmgRate, res.BossDie = svc.gl.fight(player, boss)
	if !res.Win {
		return res, nil
	}
	// win
	// cost energy
	_, err = svc.dao.consumeEnergy(ctx, req.UserId, bossRsc.EnergyCost, nowUnix)
	if err != nil {
		return nil, err
	}
	res.EnergyCost = bossRsc.EnergyCost
	// update win times
	dbFightHis, err := svc.dao.updateFightHistory(ctx, req.UserId, req.BossId)
	if err != nil {
		return nil, err
	}
	// give awards
	var awards []*mpb.AwardRsc
	if dbFightHis.WinTimes[req.BossId] == 1 { // first win
		awards = bossRsc.FirstWinAwards
	}
	for _, v := range bossRsc.Awards {
		award := &mpb.AwardRsc{
			ItemId: v.ItemId,
		}
		award.Num = uint32(uint64(v.NumRange[1]-v.NumRange[0]) * (boss.totalHP - boss.hp - boss.winLoseHp) /
			(boss.totalHP - boss.winLoseHp))
		award.Num += v.NumRange[0]
		awards = append(awards, award)
	}

	res.Awards, err = com.AddItemsFromAwardRsc(ctx, svc, req.UserId, [][]*mpb.AwardRsc{awards},
		mpb.EItem_TransReason_GameFight, uint64(boss.bossId))
	if err != nil {
		svc.logger.Error("fight boss give awards failed", zap.Uint64("user_id", req.UserId),
			zap.Any("awards", awards), zap.Error(err))
		return nil, err
	}
	return res, nil
}

func (svc *GameService) fightHiddenBoss(ctx context.Context, req *mpb.ReqFight, bossRsc *mpb.BossRsc, nowUnix int64) (
	*mpb.ResFight, error) {

	// check fight cd
	if svc.inHiddenBossFightCD(ctx, req.UserId, nowUnix) {
		return nil, mpberr.ErrParam
	}

	equips, err := svc.getUserEquips(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	if len(bossRsc.NftEquips) > 0 {
		var totalLevel uint32
		for _, equipId := range bossRsc.NftEquips {
			equip := equips.NftEquips[equipId]
			if equip == nil {
				return nil, mpberr.ErrNFTWeaponNotEquipped
			}
			totalLevel += equip.Level
		}
		if totalLevel < bossRsc.NftEquipsLevel {
			return nil, mpberr.ErrBaseEquipMaxLevel
		}
	}

	boss := newGameBoss(bossRsc)
	player := newPlayer(svc, req.UserId)
	player.updateEquips(equips.BaseEquips, equips.NftEquips)
	res := &mpb.ResFight{}
	var dbBoss *mpb.DBHiddenBoss
	err = svc.dao.fightHiddenBoss(ctx, req.BossUuid, func() error {
		dbBoss, err = svc.dao.getHiddenBoss(ctx, req.BossUuid)
		if err != nil {
			return err
		}
		if dbBoss.ExpiredAt < nowUnix {
			return mpberr.ErrHiddenBossExpired
		}
		if dbBoss.Hp == 0 {
			return mpberr.ErrHiddenBossDied
		}
		if dbBoss.LastFightTime[req.UserId] > 0 {
			return mpberr.ErrHiddenBossFought
		}
		boss.hp = dbBoss.Hp
		// fight
		res.Win, res.Details, res.Dmg, res.DmgRate, res.BossDie = svc.gl.fight(player, boss)
		// update dbBoss
		fmt.Println(boss.hp)
		dbBoss.Hp = boss.hp
		if dbBoss.Hp == 0 {
			dbBoss.Killer = boss.killedBy
		}
		dbBoss.LastFightTime[req.UserId] = nowUnix
		if res.Dmg >= boss.totalHP/100 {
			dbBoss.Dmgs[req.UserId] = res.Dmg
		}
		err = svc.dao.updateHiddenBoss(ctx, dbBoss, com.KeepTTL)
		if err != nil {
			svc.logger.Error("fightHiddenBoss update boss failed", zap.Uint64("user_id", req.UserId),
				zap.Any("boss", dbBoss), zap.Error(err))
		}
		return nil
	})
	if err != nil {
		svc.logger.Error("fightHiddenBoss failed", zap.Uint64("user_id", req.UserId),
			zap.Uint32("boss_id", bossRsc.BossId), zap.Uint64("boss_uuid", req.BossUuid),
			zap.Error(err))
		return nil, err
	}
	// cost energy
	_, err = svc.dao.consumeEnergy(ctx, req.UserId, bossRsc.EnergyCost, nowUnix)
	if err != nil {
		return nil, err
	}

	err = svc.dao.updateHiddenBossFightCD(ctx, req.UserId, svc.rm.getFightHiddenBossCd()+nowUnix)
	if err != nil {
		return nil, err
	}

	if boss.hp > 0 { // boss still alive
		return res, nil
	}

	// boss die
	awards := make(map[uint64][][]*mpb.AwardRsc)
	// 1.finder awards
	awards[dbBoss.Finder] = [][]*mpb.AwardRsc{bossRsc.FinderAwards}
	// 2.killer awards
	awards[dbBoss.Killer] = [][]*mpb.AwardRsc{bossRsc.KillerAwards}
	// 3.dmg awards
	for uid, dmg := range dbBoss.Dmgs {
		dmgAwards := make([]*mpb.AwardRsc, 0, 1)
		for _, v := range bossRsc.DmgAwards {
			dmgAwards = append(dmgAwards, &mpb.AwardRsc{
				ItemId: v.ItemId,
				Num: uint32((bossRsc.DmgAwardsCoe1*uint64(v.Num)*dmg +
					bossRsc.DmgAwardsCoe2*uint64(v.Num)*boss.totalHP +
					boss.totalHP*com.RateBase/2) / // for round
					(boss.totalHP * com.RateBase)),
			})
		}
		uAwards := awards[uid]
		if len(uAwards) == 0 {
			uAwards = append(uAwards, [][]*mpb.AwardRsc{dmgAwards}...)
		} else {
			uAwards[0] = append(uAwards[0], dmgAwards...)
		}
		awards[uid] = uAwards
	}
	cAwards, err := com.BatchAddItemsFromAwardRsc(ctx, svc, awards, mpb.EItem_TransReason_GameFight, req.BossUuid)
	if err != nil {
		svc.logger.Error("fightHiddenBoss give awards failed", zap.Uint64("user_id", req.UserId),
			zap.Any("awards", awards), zap.Error(err))
		return nil, err
	}
	res.Awards = cAwards[req.UserId]
	return res, nil
}

func (svc *GameService) getUserEquips(ctx context.Context, userId uint64) (*mpb.ResGetEquips, error) {
	itemClient, err := com.GetItemServiceClient(ctx, svc)
	if err != nil {
		return nil, err
	}
	res, err := itemClient.GetEquips(ctx, &mpb.ReqUserId{
		UserId: userId,
	})
	if err != nil {
		svc.logger.Error("getUserEquips get user equips failed", zap.Uint64("user_id", userId),
			zap.Error(err))
		return nil, err
	}
	return res, nil
}

func (svc *GameService) triggerHiddenBoss(ctx context.Context, userId uint64, bossId uint32, nowUnix int64) (
	hiddenBoss *mpb.HiddenBoss, err error) {
	triggerRsc := svc.rm.getHiddenBossTriggerRsc(bossId)
	if triggerRsc == nil {
		return nil, nil
	}

	if !util.IsPick(triggerRsc.TriggerRate, com.RateBase) {
		return nil, nil
	}

	bossRsc := svc.rm.getBossRsc(triggerRsc.HiddenBossId)
	if bossRsc == nil {
		return nil, nil
	}

	// check whether user can trigger new boss or not
	dbFindHis, err := svc.dao.getHiddenBossFindHistory(ctx, userId)
	if err != nil {
		return nil, err
	}

	for bossUUID, expireAt := range dbFindHis.BossExpireAt {
		if expireAt < nowUnix {
			delete(dbFindHis.BossExpireAt, bossUUID)
			continue
		}
		dbBoss, err := svc.dao.getHiddenBoss(ctx, bossUUID)
		if err != nil {
			return nil, err
		}
		curBossRsc := svc.rm.getBossRsc(dbBoss.BossId)
		if curBossRsc == nil {
			return nil, mpberr.ErrConfig
		}
		boss := newGameBoss(curBossRsc)
		if dbBoss.Hp > (boss.totalHP - boss.winLoseHp) {
			return nil, nil
		}
		delete(dbFindHis.BossExpireAt, bossUUID)
	}

	// new a hidden boss
	bossUUID, err := svc.bossUUIDSF.Next()
	if err != nil {
		return nil, err
	}

	boss := newGameBoss(bossRsc)
	dbBoss := &mpb.DBHiddenBoss{
		BossUuid:  bossUUID,
		BossId:    bossRsc.BossId,
		Finder:    userId,
		Hp:        boss.hp,
		ExpiredAt: bossRsc.LiveTime + nowUnix,
	}
	err = svc.dao.updateHiddenBoss(ctx, dbBoss, time.Duration(bossRsc.LiveTime+com.Secs1Week)*time.Second)
	if err != nil {
		return nil, err
	}

	// update find history
	dbFindHis.BossExpireAt[bossUUID] = dbBoss.ExpiredAt
	err = svc.dao.updateHiddenBossFindHistory(ctx, userId, dbFindHis)
	if err != nil {
		return nil, err
	}
	return svc.dbHiddenBoss2HiddenBoss(dbBoss), nil
}

func (svc *GameService) inHiddenBossFightCD(ctx context.Context, userId uint64, nowUnix int64) bool {
	cd, err := svc.dao.getHiddenBossFightCD(ctx, userId)
	if err != nil {
		return true
	}
	return cd > nowUnix
}

func (svc *GameService) AddEnergy(ctx context.Context, req *mpb.ReqAddEnergy) (*mpb.ResAddEnergy, error) {
	dbEnergy, err := svc.dao.addEnergy(ctx, req.UserId, req.Energy, time.Now().Unix())
	if err != nil {
		return nil, err
	}
	return &mpb.ResAddEnergy{
		Energy:   dbEnergy.Energy,
		UpdateAt: dbEnergy.RecoverAt,
	}, nil
}

func (svc *GameService) GetHiddenBoss(ctx context.Context, req *mpb.ReqGetHiddenBoss) (*mpb.ResGetHiddenBoss, error) {
	res := &mpb.ResGetHiddenBoss{}
	dbHiddenBoss, err := svc.dao.getHiddenBoss(ctx, req.BossUuid)
	if err != nil {
		return nil, err
	}
	res.HiddenBoss = svc.dbHiddenBoss2HiddenBoss(dbHiddenBoss)
	res.Fought = dbHiddenBoss.LastFightTime[req.UserId] > 0
	dbFightCD, err := svc.dao.getHiddenBossFightCD(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	res.FightCd = dbFightCD
	return res, nil
}

func (svc *GameService) GetEnergy(ctx context.Context, req *mpb.ReqUserId) (*mpb.ResGetEnergy, error) {
	dbEnergy, err := svc.dao.getEnergy(ctx, req.UserId, time.Now().Unix())
	if err != nil {
		return nil, err
	}

	return &mpb.ResGetEnergy{
		Energy:   dbEnergy.Energy,
		UpdateAt: dbEnergy.RecoverAt,
	}, nil
}
