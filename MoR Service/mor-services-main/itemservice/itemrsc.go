package itemservice

import (
	"fmt"
	"github.com/oldjon/gutil"
	gcsv "github.com/oldjon/gutil/csv"
	gdm "github.com/oldjon/gutil/dirmonitor"
	com "gitlab.com/morbackend/mor_services/common"
	"gitlab.com/morbackend/mor_services/mpb"
	"gitlab.com/morbackend/mor_services/util"
	"go.uber.org/zap"
)

const (
	//csvSuffix   = ".csv"
	baseCSVPath = "./resources/item/"

	itemCSV            = "Items.csv"
	baseEquipLevelsCSV = "BaseEquipLevels.csv"
	baseEquipStarsCSV  = "BaseEquipStars.csv"
)

type itemResourceMgr struct {
	logger *zap.Logger
	dm     *gdm.DirMonitor
	mtr    *util.ServiceMetrics

	items             map[uint32]*mpb.ItemRsc
	baseEquipStars    map[mpb.EItem_BaseEquipType]map[uint32]*mpb.BaseEquipStarRsc
	baseEquipMaxStar  map[mpb.EItem_BaseEquipType]uint32
	baseEquipLevels   map[mpb.EItem_BaseEquipType]map[uint32]*mpb.BaseEquipLevelRsc
	baseEquipMaxLevel map[mpb.EItem_BaseEquipType]uint32
}

func newItemResourceMgr(logger *zap.Logger, mtr *util.ServiceMetrics) (*itemResourceMgr, error) {
	rMgr := &itemResourceMgr{
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

func (rm *itemResourceMgr) watch() error {
	return rm.dm.StartWatch()
}

func (rm *itemResourceMgr) load() error {
	err := rm.dm.BindAndExec(itemCSV, rm.loadItems)
	if err != nil {
		return err
	}
	err = rm.dm.BindAndExec(baseEquipStarsCSV, rm.loadBaseEquipStars)
	if err != nil {
		return err
	}
	err = rm.dm.BindAndExec(baseEquipLevelsCSV, rm.loadBaseEquipLevels)
	if err != nil {
		return err
	}

	return nil
}

func (rm *itemResourceMgr) loadItems(csvPath string) error {
	datas, err := gcsv.ReadCSV2Array(csvPath)
	if err != nil {
		rm.logger.Error(fmt.Sprintf("load %s failed: %s", csvPath, err.Error()))
		return err
	}
	m := make(map[uint32]*mpb.ItemRsc)
	for _, data := range datas {
		node := &mpb.ItemRsc{
			ItemId:     gcsv.Str2Uint32(data["itemid"]),
			ItemType:   mpb.EItem_Type(gcsv.Str2Uint32(data["itemtype"])),
			IsUnique:   gcsv.Str2Bool(data["isunique"]),
			OriginId:   gcsv.Str2Uint32(data["originid"]),
			ExpireTime: gcsv.Str2Int64(data["expiretime"]),
		}

		m[node.ItemId] = node
		rm.logger.Debug("loadItems read:", zap.Any("row", node))
	}

	rm.logger.Debug("loadItems read finish:", zap.Any("rows", m))
	rm.items = m

	return nil
}

func (rm itemResourceMgr) getItemRsc(itemId uint32) *mpb.ItemRsc {
	return rm.items[itemId]
}

func (rm *itemResourceMgr) loadBaseEquipStars(csvPath string) error {
	datas, err := gcsv.ReadCSV2Array(csvPath)
	if err != nil {
		rm.logger.Error(fmt.Sprintf("load %s failed: %s", csvPath, err.Error()))
		return err
	}
	m := make(map[mpb.EItem_BaseEquipType]map[uint32]*mpb.BaseEquipStarRsc)
	maxM := make(map[mpb.EItem_BaseEquipType]uint32)
	for _, data := range datas {
		node := &mpb.BaseEquipStarRsc{
			EquipType:          mpb.EItem_BaseEquipType(gcsv.Str2Uint32(data["equiptype"])),
			Star:               gcsv.Str2Uint32(data["star"]),
			Attrs:              com.ReadAttrs(data),
			UpgradeSuccessRate: gcsv.Str2Uint32(data["upgradesuccessrate"]),
			ProtectSuccessNum:  gcsv.Str2Uint32(data["protectsuccessnum"]),
		}
		node.UpgradeConsumeItems, err = com.ReadItems(data["upgradeconsumeitems"])
		if err != nil {
			rm.logger.Error("loadBaseEquipStars parse upgradeconsumeitems failed",
				zap.String("upgradeconsumeitems", data["upgradeconsumeitems"]))
			return err
		}

		subm := m[node.EquipType]
		if subm == nil {
			subm = map[uint32]*mpb.BaseEquipStarRsc{}
			m[node.EquipType] = subm
		}
		subm[node.Star] = node
		maxM[node.EquipType] = gutil.Max(maxM[node.EquipType], node.Star)
		rm.logger.Debug("loadBaseEquipStars read:", zap.Any("row", node))
	}

	rm.logger.Debug("loadBaseEquipStars read finish:", zap.Any("rows", m))
	rm.baseEquipStars = m
	rm.baseEquipMaxStar = maxM
	return nil
}

func (rm *itemResourceMgr) getBaseEquipStarRsc(equipType mpb.EItem_BaseEquipType, star uint32) *mpb.BaseEquipStarRsc {
	return rm.baseEquipStars[equipType][star]
}

func (rm *itemResourceMgr) getBaseEquipMaxStar(equipType mpb.EItem_BaseEquipType) uint32 {
	return rm.baseEquipMaxStar[equipType]
}

func (rm *itemResourceMgr) loadBaseEquipLevels(csvPath string) error {
	datas, err := gcsv.ReadCSV2Array(csvPath)
	if err != nil {
		rm.logger.Error(fmt.Sprintf("load %s failed: %s", csvPath, err.Error()))
		return err
	}
	m := make(map[mpb.EItem_BaseEquipType]map[uint32]*mpb.BaseEquipLevelRsc)
	maxM := make(map[mpb.EItem_BaseEquipType]uint32)
	for _, data := range datas {
		node := &mpb.BaseEquipLevelRsc{
			EquipType:          mpb.EItem_BaseEquipType(gcsv.Str2Uint32(data["equiptype"])),
			Level:              gcsv.Str2Uint32(data["level"]),
			Attrs:              com.ReadAttrs(data),
			UpgradeSuccessRate: gcsv.Str2Uint32(data["upgradesuccessrate"]),
			ProtectSuccessNum:  gcsv.Str2Uint32(data["protectsuccessnum"]),
		}

		node.UpgradeConsumeItems, err = com.ReadItems(data["upgradeconsumeitems"])
		if err != nil {
			rm.logger.Error("loadBaseEquipStars parse upgradeconsumeitems failed",
				zap.String("upgradeconsumeitems", data["upgradeconsumeitems"]))
			return err
		}

		subm := m[node.EquipType]
		if subm == nil {
			subm = map[uint32]*mpb.BaseEquipLevelRsc{}
			m[node.EquipType] = subm
		}
		subm[node.Level] = node
		maxM[node.EquipType] = gutil.Max(maxM[node.EquipType], node.Level)
		rm.logger.Debug("loadBaseEquipLevels read:", zap.Any("row", node))
	}

	rm.logger.Debug("loadBaseEquipLevels read finish:", zap.Any("rows", m))
	rm.baseEquipLevels = m
	rm.baseEquipMaxLevel = maxM
	return nil
}

func (rm *itemResourceMgr) getBaseEquipLevelRsc(equipType mpb.EItem_BaseEquipType, level uint32) *mpb.BaseEquipLevelRsc {
	return rm.baseEquipLevels[equipType][level]
}

func (rm *itemResourceMgr) getBaseEquipMaxLevel(equipType mpb.EItem_BaseEquipType) uint32 {
	return rm.baseEquipMaxLevel[equipType]
}
