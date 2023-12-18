package gameservice

import (
	"context"
	"github.com/oldjon/gutil"
	"github.com/oldjon/gutil/gdb"
	grmux "github.com/oldjon/gutil/redismutex"
	com "gitlab.com/morbackend/mor_services/common"
	"gitlab.com/morbackend/mor_services/mpb"
	"gitlab.com/morbackend/mor_services/mpberr"
	"go.uber.org/zap"
	"time"
)

type gameDAO struct {
	logger *zap.Logger
	rm     *gameResourceMgr
	rMux   *grmux.RedisMutex
	gameDB *gdb.DB
	tmpDB  *gdb.DB
}

func newGameDAO(svc *GameService, rMux *grmux.RedisMutex, gameRedis, tmpRedis gdb.RedisClient) *gameDAO {
	return &gameDAO{
		logger: svc.logger,
		rm:     svc.rm,
		rMux:   rMux,
		gameDB: gdb.NewDB(gameRedis),
		tmpDB:  gdb.NewDB(tmpRedis),
	}
}

func (dao *gameDAO) getEnergy(ctx context.Context, userId uint64, nowUnix int64) (*mpb.DBEnergy, error) {
	dbEnergy := &mpb.DBEnergy{}
	err := dao.gameDB.GetObject(ctx, com.EnergyKey(userId), dbEnergy)
	if dao.gameDB.IsErrNil(err) {
		dbEnergy.Energy = dao.rm.getEnergyLimit()
		dbEnergy.RecoverAt = nowUnix
		return dbEnergy, nil
	}
	if err != nil {
		dao.logger.Error("get energy failed", zap.Uint64("user_id", userId), zap.Error(err))
		return nil, mpberr.ErrDB
	}
	if dbEnergy.Energy >= dao.rm.getEnergyLimit() {
		dbEnergy.RecoverAt = time.Now().Unix()
		return dbEnergy, nil
	}
	er := uint32((nowUnix - dbEnergy.RecoverAt) / dao.rm.getEnergyRecoverTime())
	dbEnergy.Energy = gutil.Min(dbEnergy.Energy+er, dao.rm.getEnergyLimit())
	dbEnergy.RecoverAt = gutil.If(dbEnergy.Energy == dao.rm.getEnergyLimit(),
		nowUnix,
		dbEnergy.RecoverAt+int64(er)*dao.rm.getEnergyRecoverTime())
	return dbEnergy, err
}

func (dao *gameDAO) consumeEnergy(ctx context.Context, userId uint64, consumeEnergy uint32, nowUnix int64) (
	*mpb.DBEnergy, error) {
	var dbEnergy *mpb.DBEnergy
	key := com.EnergyKey(userId)
	err := dao.rMux.Safely(ctx, key, func() error {
		var err error
		dbEnergy, err = dao.getEnergy(ctx, userId, nowUnix)
		if err != nil {
			return err
		}
		if dbEnergy.Energy < consumeEnergy {
			return mpberr.ErrEnergyNotEnough
		}
		dbEnergy.Energy -= consumeEnergy
		err = dao.gameDB.SetObject(ctx, key, dbEnergy)
		if err != nil {
			dao.logger.Error("consumeEnergy update db failed", zap.Uint64("user_id", userId),
				zap.Any("energy", dbEnergy), zap.Error(err))
			return mpberr.ErrDB
		}
		return nil
	})
	if err != nil {
		dao.logger.Error("consumeEnergy failed", zap.Uint64("user_id", userId), zap.Error(err))
		return nil, err
	}
	return dbEnergy, nil
}

func (dao *gameDAO) addEnergy(ctx context.Context, userId uint64, addEnergy uint32, nowUnix int64) (*mpb.DBEnergy,
	error) {
	var dbEnergy *mpb.DBEnergy
	key := com.EnergyKey(userId)
	err := dao.rMux.Safely(ctx, key, func() error {
		var err error
		dbEnergy, err = dao.getEnergy(ctx, userId, nowUnix)
		if err != nil {
			return err
		}
		dbEnergy.Energy += addEnergy
		err = dao.gameDB.SetObject(ctx, key, dbEnergy)
		if err != nil {
			dao.logger.Error("addEnergy update db failed", zap.Uint64("user_id", userId),
				zap.Any("energy", dbEnergy), zap.Error(err))
			return mpberr.ErrDB
		}
		return nil
	})
	if err != nil {
		dao.logger.Error("addEnergy failed", zap.Uint64("user_id", userId), zap.Error(err))
		return nil, err
	}
	return dbEnergy, nil
}

func (dao *gameDAO) getFightHistory(ctx context.Context, userId uint64) (*mpb.DBFightHistory, error) {
	dbHis := &mpb.DBFightHistory{}
	err := dao.gameDB.GetObject(ctx, com.FightHistoryKey(userId), dbHis)
	if err != nil && !dao.gameDB.IsErrNil(err) {
		dao.logger.Error("getFightHistory get failed", zap.Uint64("user_id", userId), zap.Error(err))
		return nil, mpberr.ErrDB
	}
	if dbHis.WinTimes == nil {
		dbHis.WinTimes = make(map[uint32]uint32)
	}
	return dbHis, nil
}

func (dao *gameDAO) updateFightHistory(ctx context.Context, userId uint64, bossId uint32) (*mpb.DBFightHistory,
	error) {
	dbHis := &mpb.DBFightHistory{}
	key := com.FightHistoryKey(userId)
	err := dao.rMux.Safely(ctx, key, func() error {
		err := dao.gameDB.GetObject(ctx, key, dbHis)
		if err != nil && !dao.gameDB.IsErrNil(err) {
			dao.logger.Error("updateFightHistory get obj from db failed",
				zap.Uint64("user_id", userId), zap.Error(err))
			return mpberr.ErrDB
		}
		if dbHis.WinTimes == nil {
			dbHis.WinTimes = make(map[uint32]uint32)
		}
		dbHis.WinTimes[bossId] += 1
		err = dao.gameDB.SetObject(ctx, key, dbHis)
		if err != nil {
			dao.logger.Error("updateFightHistory set obj failed",
				zap.Uint64("user_id", userId), zap.Error(err))
			return mpberr.ErrDB
		}
		return nil
	})
	if err != nil {
		dao.logger.Error("updateFightHistory failed", zap.Uint64("user_id", userId), zap.Error(err))
		return nil, err
	}
	return dbHis, nil
}

func (dao *gameDAO) getHiddenBoss(ctx context.Context, bossUUID uint64) (*mpb.DBHiddenBoss, error) {
	dbHiddenBoss := &mpb.DBHiddenBoss{}
	err := dao.gameDB.GetObject(ctx, com.HiddenBossKey(bossUUID), dbHiddenBoss)
	if err != nil && !dao.gameDB.IsErrNil(err) {
		dao.logger.Error("getHiddenBoss get obj failed", zap.Uint64("boss_uuid", bossUUID), zap.Error(err))
		return nil, mpberr.ErrDB
	}
	if dao.gameDB.IsErrNil(err) {
		return nil, mpberr.ErrHiddenBossNotFound
	}
	if dbHiddenBoss.LastFightTime == nil {
		dbHiddenBoss.LastFightTime = make(map[uint64]int64)
	}
	if dbHiddenBoss.Dmgs == nil {
		dbHiddenBoss.Dmgs = make(map[uint64]uint64)
	}
	return dbHiddenBoss, nil
}

func (dao *gameDAO) updateHiddenBoss(ctx context.Context, dbHiddenBoss *mpb.DBHiddenBoss, expiration time.Duration,
) error {
	err := dao.gameDB.SetObjectEX(ctx, com.HiddenBossKey(dbHiddenBoss.BossUuid), dbHiddenBoss, expiration)
	if err != nil && !dao.gameDB.IsErrNil(err) {
		dao.logger.Error("updateHiddenBoss set obj failed",
			zap.Uint64("boss_uuid", dbHiddenBoss.BossUuid),
			zap.Error(err))
		return mpberr.ErrDB
	}
	return nil
}

func (dao *gameDAO) fightHiddenBoss(ctx context.Context, bossUUId uint64, fightFunc func() error) error {
	err := dao.rMux.Safely(ctx, com.HiddenBossKey(bossUUId), fightFunc)
	if err != nil {
		dao.logger.Error("fightHiddenBoss failed", zap.Uint64("boss_uuid", bossUUId),
			zap.Error(err))
		return err
	}
	return nil
}

func (dao *gameDAO) getHiddenBossFightCD(ctx context.Context, userId uint64) (int64, error) {
	t, err := gdb.ToUint64(dao.gameDB.Get(ctx, com.UserHiddenBossCDKey(userId)))
	if err != nil && !dao.gameDB.IsErrNil(err) {
		dao.logger.Error("getHiddenBossFightCD get cd failed",
			zap.Uint64("user_id", userId), zap.Error(err))
		return 0, mpberr.ErrDB
	}
	return int64(t), nil
}

func (dao *gameDAO) updateHiddenBossFightCD(ctx context.Context, userId uint64, cd int64) error {
	err := dao.gameDB.Set(ctx, com.UserHiddenBossCDKey(userId), cd)
	if err != nil {
		dao.logger.Error("updateHiddenBossFightCD failed", zap.Uint64("user_id", userId),
			zap.Int64("cd", cd), zap.Error(err))
		return err
	}
	return nil
}

func (dao *gameDAO) getHiddenBossFindHistory(ctx context.Context, userId uint64) (*mpb.DBHiddenBossFindHistory, error) {
	dbHis := &mpb.DBHiddenBossFindHistory{}
	err := dao.gameDB.GetObject(ctx, com.HiddenBossFindHistoryKey(userId), dbHis)
	if err != nil && !dao.gameDB.IsErrNil(err) {
		dao.logger.Error("getHiddenBossFindHistory get cd failed",
			zap.Uint64("user_id", userId), zap.Error(err))
		return nil, mpberr.ErrDB
	}
	if dbHis.BossExpireAt == nil {
		dbHis.BossExpireAt = make(map[uint64]int64)
	}
	return dbHis, nil
}

func (dao *gameDAO) updateHiddenBossFindHistory(ctx context.Context, userId uint64, dbHis *mpb.DBHiddenBossFindHistory,
) error {
	err := dao.gameDB.SetObject(ctx, com.HiddenBossFindHistoryKey(userId), dbHis)
	if err != nil {
		dao.logger.Error("updateHiddenBossFindHistory failed",
			zap.Uint64("user_id", userId), zap.Any("db_his", dbHis),
			zap.Error(err))
		return err
	}
	return nil
}
