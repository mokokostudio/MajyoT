package nftservice

import (
	"github.com/oldjon/gutil/gdb"
	grmux "github.com/oldjon/gutil/redismutex"
	"go.uber.org/zap"
)

type nftDAO struct {
	svc    *NFTService
	logger *zap.Logger
	rMux   *grmux.RedisMutex
	nftDB  *gdb.DB
	tranDB *gdb.DB
	tmpDB  *gdb.DB
	rm     *NFTResourceMgr
}

func newNftDAO(svc *NFTService, rMux *grmux.RedisMutex, nftRedis, tmpRedis, tranRedis gdb.RedisClient) *nftDAO {
	return &nftDAO{
		logger: svc.logger,
		rMux:   rMux,
		nftDB:  gdb.NewDB(nftRedis),
		tranDB: gdb.NewDB(tranRedis),
		tmpDB:  gdb.NewDB(tmpRedis),
		rm:     svc.rm,
	}
}
