package accountservice

import (
	"gitlab.com/morbackend/mor_services/util"
	"go.uber.org/zap"
)

type AccountResourceMgr struct {
}

func newAccountResourceMgr(logger *zap.Logger, sm *util.ServiceMetrics) (*AccountResourceMgr, error) {
	return &AccountResourceMgr{}, nil
}
