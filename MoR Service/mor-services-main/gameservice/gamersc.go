package gameservice

import (
	"fmt"
	com "gitlab.com/morbackend/mor_services/common"

	"github.com/oldjon/gutil"
	gcsv "github.com/oldjon/gutil/csv"
	gdm "github.com/oldjon/gutil/dirmonitor"
	"gitlab.com/morbackend/mor_services/mpb"
	"gitlab.com/morbackend/mor_services/mpberr"
	"gitlab.com/morbackend/mor_services/util"
	"go.uber.org/zap"
)

const (
	//csvSuffix   = ".csv"
	baseCSVPath = "./resources/game/"

	gameConfigCSV        = "GameConfig.csv"
	playerInitAttrsCSV   = "PlayerInitAttrs.csv"
	bossCSV              = "Boss.csv"
	hiddenBossTriggerCSV = "HiddenBossTrigger.csv"
)

type gameResourceMgr struct {
	logger *zap.Logger
	dm     *gdm.DirMonitor
	mtr    *util.ServiceMetrics

	gameConfig        *mpb.GameConfigRcs
	playerInitAttrs   *mpb.PlayerInitAttrs
	bosses            map[uint32]*mpb.BossRsc
	hiddenBossTrigger map[uint32]*mpb.HiddenBossTriggerRcs
}

func newGameResourceMgr(logger *zap.Logger, mtr *util.ServiceMetrics) (*gameResourceMgr, error) {
	rMgr := &gameResourceMgr{
		logger: logger,
		mtr:    mtr,
	}

	var err error
	rMgr.dm, err = gdm.NewDirMonitor(baseCSVPath)
	if err != nil {
		return nil, err
	}

	err = rMgr.load()
	if err != nil {
		return nil, err
	}

	err = rMgr.watch()
	if err != nil {
		return nil, err
	}

	return rMgr, nil
}

func (rm *gameResourceMgr) watch() error {
	return rm.dm.StartWatch()
}

func (rm *gameResourceMgr) load() error {
	err := rm.dm.BindAndExec(gameConfigCSV, rm.loadGameConfig)
	if err != nil {
		return err
	}
	err = rm.dm.BindAndExec(playerInitAttrsCSV, rm.loadPlayerInitAttrs)
	if err != nil {
		return err
	}
	err = rm.dm.BindAndExec(bossCSV, rm.loadBoss)
	if err != nil {
		return err
	}
	err = rm.dm.BindAndExec(hiddenBossTriggerCSV, rm.loadHiddenBossTrigger)
	if err != nil {
		return err
	}
	return nil
}

func (rm *gameResourceMgr) loadGameConfig(csvPath string) error {
	datas, err := gcsv.ReadCSV2Array(csvPath)
	if err != nil {
		rm.logger.Error(fmt.Sprintf("load %s failed: %s", csvPath, err.Error()))
		return err
	}
	if len(datas) != 1 {
		rm.logger.Error(fmt.Sprintf("load %s failed: config row num %d", csvPath, len(datas)))
		return mpberr.ErrConfig
	}
	node := &mpb.GameConfigRcs{
		EnergyLimit:       gutil.Str2Uint32(datas[0]["energylimit"]),
		EnergyRecoverTime: gutil.Str2Int64(datas[0]["energyrecovertime"]),
		FightHiddenBossCd: gutil.Str2Int64(datas[0]["fighthiddenbosscd"]),
	}

	rm.gameConfig = node
	rm.logger.Debug("loadGameConfig read finish:", zap.Any("row", node))
	return nil
}

func (rm *gameResourceMgr) getEnergyLimit() uint32 {
	return rm.gameConfig.GetEnergyLimit()
}

func (rm *gameResourceMgr) getEnergyRecoverTime() int64 {
	return rm.gameConfig.GetEnergyRecoverTime()
}

func (rm *gameResourceMgr) getFightHiddenBossCd() int64 {
	return rm.gameConfig.GetFightHiddenBossCd()
}

func (rm *gameResourceMgr) loadPlayerInitAttrs(csvPath string) error {
	datas, err := gcsv.ReadCSV2Array(csvPath)
	if err != nil {
		rm.logger.Error(fmt.Sprintf("load %s failed: %s", csvPath, err.Error()))
		return err
	}
	if len(datas) != 1 {
		rm.logger.Error(fmt.Sprintf("load %s failed: config row num %d", csvPath, len(datas)))
		return mpberr.ErrConfig
	}
	node := &mpb.PlayerInitAttrs{
		Attrs: com.ReadAttrs(datas[0]),
	}

	rm.logger.Debug("loadPlayerInitAttrs read finish:", zap.Any("row", node))
	rm.playerInitAttrs = node
	return nil
}

func (rm *gameResourceMgr) getPlayerInitAttrs() *mpb.Attrs {
	return com.CloneAttrs(rm.playerInitAttrs.Attrs)
}

func (rm *gameResourceMgr) loadBoss(csvPath string) error {
	datas, err := gcsv.ReadCSV2Array(csvPath)
	if err != nil {
		rm.logger.Error(fmt.Sprintf("load %s failed: %s", csvPath, err.Error()))
		return err
	}
	m := make(map[uint32]*mpb.BossRsc)
	for _, data := range datas {
		node := &mpb.BossRsc{
			BossId:         gutil.Str2Uint32(data["bossid"]),
			BossType:       mpb.ERole_BossType(gutil.Str2Uint32(data["bosstype"])),
			Class:          gutil.Str2Uint32(data["class"]),
			Level:          gutil.Str2Uint32(data["level"]),
			LiveTime:       gutil.Str2Int64(data["livetime"]),
			PreBoss:        gutil.Str2Uint32(data["preboss"]),
			NftEquips:      util.ReadUint32Slice(data["nftequips"], ";"),
			NftEquipsLevel: gutil.Str2Uint32(data["nftequipslevel"]),
			EnergyCost:     gutil.Str2Uint32(data["energycost"]),
			Attrs:          com.ReadAttrs(data),
			WinDmgRate:     gutil.Str2Uint64(data["windmgrate"]),
			DmgAwardsCoe1:  gutil.Str2Uint64(data["dmgawardscoe1"]),
			DmgAwardsCoe2:  gutil.Str2Uint64(data["dmgawardscoe2"]),
		}
		node.FirstWinAwards, err = com.ReadAwardsRsc(data["firstwinawards"])
		if err != nil {
			rm.logger.Error("loadBoss parse firstwinawards failed",
				zap.String("firstwinawards", data["firstwinawards"]))
			return err
		}
		node.Awards, err = com.ReadAwardsRsc(data["awards"])
		if err != nil {
			rm.logger.Error("loadBoss parse awards failed",
				zap.String("awards", data["awards"]))
			return err
		}
		node.FinderAwards, err = com.ReadAwardsRsc(data["finderawards"])
		if err != nil {
			rm.logger.Error("loadBoss parse finderawards failed",
				zap.String("finderawards", data["finderawards"]))
			return err
		}
		node.KillerAwards, err = com.ReadAwardsRsc(data["killerawards"])
		if err != nil {
			rm.logger.Error("loadBoss parse killerawards failed",
				zap.String("killerawards", data["killerawards"]))
			return err
		}
		node.DmgAwards, err = com.ReadAwardsRsc(data["dmgawards"])
		if err != nil {
			rm.logger.Error("loadBoss parse dmgawards failed",
				zap.String("dmgawards", data["dmgawards"]))
			return err
		}
		m[node.BossId] = node
		rm.logger.Debug("loadBoss read:", zap.Any("row", node))
	}

	rm.bosses = m
	rm.logger.Debug("loadBoss read finish:", zap.Any("rows", m))

	return nil
}

func (rm *gameResourceMgr) getBossRsc(bossId uint32) *mpb.BossRsc {
	return rm.bosses[bossId]
}

func (rm *gameResourceMgr) loadHiddenBossTrigger(csvPath string) error {
	datas, err := gcsv.ReadCSV2Array(csvPath)
	if err != nil {
		rm.logger.Error(fmt.Sprintf("load %s failed: %s", csvPath, err.Error()))
		return err
	}
	m := make(map[uint32]*mpb.HiddenBossTriggerRcs)
	for _, data := range datas {
		node := &mpb.HiddenBossTriggerRcs{
			BossId:       gutil.Str2Uint32(data["bossid"]),
			TriggerRate:  gutil.Str2Uint32(data["triggerrate"]),
			HiddenBossId: gutil.Str2Uint32(data["hiddenbossid"]),
		}

		m[node.BossId] = node
		rm.logger.Debug("loadHiddenBossTrigger read:", zap.Any("row", node))
	}

	rm.hiddenBossTrigger = m
	rm.logger.Debug("loadHiddenBossTrigger read finish:", zap.Any("rows", m))
	return nil
}

func (rm *gameResourceMgr) getHiddenBossTriggerRsc(bossId uint32) *mpb.HiddenBossTriggerRcs {
	return rm.hiddenBossTrigger[bossId]
}
