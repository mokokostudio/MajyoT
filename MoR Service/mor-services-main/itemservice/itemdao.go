package itemservice

import (
	"context"
	"gitlab.com/morbackend/mor_services/util"
	"time"

	gutil "github.com/oldjon/gutil"
	"github.com/oldjon/gutil/gdb"
	grmux "github.com/oldjon/gutil/redismutex"
	com "gitlab.com/morbackend/mor_services/common"
	"gitlab.com/morbackend/mor_services/mpb"
	"gitlab.com/morbackend/mor_services/mpberr"
	"go.uber.org/zap"
)

const (
	maxDeep     = 10
	weight10000 = 10000

	itemsShardNum  = 10
	uItemsShardNum = 10
)

type itemDAO struct {
	logger *zap.Logger
	rMux   *grmux.RedisMutex
	itemDB *gdb.DB
	tmpDB  *gdb.DB
	svc    *ItemService
	rm     *itemResourceMgr
}

func newItemDAO(svc *ItemService, rMux *grmux.RedisMutex, itemRedis, tmpRedis gdb.RedisClient) *itemDAO {
	return &itemDAO{
		svc:    svc,
		logger: svc.logger,
		rm:     svc.rm,
		rMux:   rMux,
		itemDB: gdb.NewDB(itemRedis),
		tmpDB:  gdb.NewDB(tmpRedis),
	}
}

func (dao *itemDAO) getItems(ctx context.Context, userId uint64) ([]*mpb.DBItem, error) {
	var dstItems = make([]interface{}, 0, itemsShardNum+uItemsShardNum)
	shards := make([]*mpb.DBItemListShard, itemsShardNum)
	shardKeys := com.ItemsShardKeysByShardCnt(itemsShardNum)
	uShards := make([]*mpb.DBUItemListShard, uItemsShardNum)
	uShardKeys := com.UItemsShardKeysByShardCnt(uItemsShardNum)
	for i := range shards {
		shards[i] = &mpb.DBItemListShard{}
		dstItems = append(dstItems, shards[i])
	}
	for i := range uShards {
		uShards[i] = &mpb.DBUItemListShard{}
		dstItems = append(dstItems, uShards[i])
	}
	err := dao.itemDB.HMGetObjects(ctx, com.ItemsKey(userId), append(shardKeys, uShardKeys...), dstItems)
	if err != nil && !dao.itemDB.IsErrNil(err) {
		dao.logger.Error("getItems failed", zap.Uint64("user_id", userId), zap.Error(err))
		return nil, mpberr.ErrDB
	}
	ret := make([]*mpb.DBItem, 0, 20)
	nowUnix := time.Now().Unix()
	for _, v := range shards {
		if v == nil || v.Items == nil {
			continue
		}
		for _, vv := range v.Items {
			if vv.ExpireAt != 0 && vv.ExpireAt < nowUnix { // don't show expired items
				continue
			}
			ret = append(ret, vv)
		}
	}
	for _, v := range uShards {
		if v == nil || v.Items == nil {
			continue
		}
		for _, vv := range v.Items {
			if vv.ExpireAt != 0 && vv.ExpireAt < nowUnix { // don't show expired items
				continue
			}
			ret = append(ret, vv)
		}
	}

	return ret, nil
}

func (dao *itemDAO) getWallet(ctx context.Context, userId uint64) (wallet *mpb.DBWallet, err error) {
	wallet = &mpb.DBWallet{}
	err = dao.itemDB.GetObject(ctx, com.WalletKey(userId), wallet)
	if err != nil && !dao.itemDB.IsErrNil(err) {
		dao.logger.Error("getWallet get wallet failed",
			zap.Uint64("user_ud", userId),
			zap.Error(err))
		return nil, mpberr.ErrDB
	}
	return wallet, nil
}

// AllItemsDB use to store items shard data temporary
type AllItemsDB struct {
	Shards    []*mpb.DBItemListShard
	UShards   []*mpb.DBUItemListShard
	DeltaMana int32
}

func (dao *itemDAO) exchangeMoney(ctx context.Context, userId uint64, deltaMana int32) (updateList []*mpb.Item, err error) {
	if deltaMana == 0 {
		return
	}

	err = dao.rMux.Safely(ctx, com.WalletKey(userId), func() error {
		wallet := &mpb.DBWallet{}
		err := dao.itemDB.GetObject(ctx, com.WalletKey(userId), wallet)
		if err != nil && !dao.itemDB.IsErrNil(err) {
			dao.logger.Error("exchangeMoney get wallet failed",
				zap.Uint64("user_ud", userId),
				zap.Int32("delta_mana", deltaMana),
				zap.Error(err))
			return mpberr.ErrDB
		}

		// mana
		if deltaMana != 0 {
			mana := int32(wallet.Mana)
			if mana+deltaMana < 0 {
				return mpberr.ErrNotEnoughMana
			}
			mana += deltaMana
			wallet.Mana = uint32(mana)
			updateList = append(updateList, &mpb.Item{ItemId: uint32(mpb.EItem_ItemId_Mana), Num: wallet.Mana})
		}

		err = dao.itemDB.SetObject(ctx, com.WalletKey(userId), wallet)
		if err != nil {
			dao.logger.Error("exchangeMoney write db failed",
				zap.Uint64("user_ud", userId),
				zap.Int32("delta_mana", deltaMana),
				zap.Error(err))
			return mpberr.ErrDB
		}
		return nil
	})
	if err != nil {
		dao.logger.Error("exchangeMoney failed",
			zap.Uint64("user_ud", userId),
			zap.Int32("delta_mana", deltaMana),
			zap.Error(err))
		return nil, err
	}
	return
}

// exchangeItems handle add and delete items.
// param addItems includes the origin items need to add.
// param delItems includes the origin items need to delete.
// showList is used to return to caller, usually return to client for show.
// updateList is used to notify client to update backpack.
// vAddList and vDelList contains virtual items those need to be processed separately.
func (dao *itemDAO) exchangeItems(ctx context.Context, userId uint64, addItems, delItems []*mpb.Item, nowUnix int64,
) (showList, updateList, vAddList, vDelList []*mpb.Item, err error) {
	if len(addItems)+len(delItems) == 0 {
		return
	}

	var deltaMana int32
	addAllMoney, addMana := dao.svc.isAllMoney(addItems)
	delAllMoney, delMana := dao.svc.isAllMoney(delItems)
	if addAllMoney && delAllMoney { // if just only need to handle money
		deltaMana = int32(addMana) - int32(delMana)
		updateList, err = dao.exchangeMoney(ctx, userId, deltaMana)
		if err != nil {
			return
		}
		showList = addItems
		return
	}

	if delMana > 0 { // check money is enough
		wallet, iErr := dao.getWallet(ctx, userId)
		if iErr != nil {
			err = iErr
			return
		}
		if delMana > wallet.Mana {
			err = mpberr.ErrNotEnoughMana
			return
		}
	}

	updateItemsMap := make(map[uint32]*mpb.DBItem)
	updateUItemsList := make([]*mpb.DBItem, 0)
	var dstItems = make([]interface{}, 0, itemsShardNum+uItemsShardNum)
	allItemsDB := &AllItemsDB{
		Shards:  make([]*mpb.DBItemListShard, itemsShardNum),
		UShards: make([]*mpb.DBUItemListShard, uItemsShardNum),
	}
	for i := range allItemsDB.Shards {
		allItemsDB.Shards[i] = &mpb.DBItemListShard{}
		dstItems = append(dstItems, allItemsDB.Shards[i])
	}
	for i := range allItemsDB.UShards {
		allItemsDB.UShards[i] = &mpb.DBUItemListShard{}
		dstItems = append(dstItems, allItemsDB.UShards[i])
	}
	shardKeys := com.ItemsShardKeysByShardCnt(itemsShardNum)
	uShardKeys := com.UItemsShardKeysByShardCnt(uItemsShardNum)
	err = dao.rMux.Safely(ctx, com.ItemsKey(userId), func() error {
		err := dao.itemDB.HMGetObjects(ctx, com.ItemsKey(userId), append(shardKeys, uShardKeys...), dstItems)
		if err != nil {
			dao.logger.Error("exchangeItems get items failed", zap.Uint64("user_id", userId), zap.Error(err))
			return mpberr.ErrDB
		}

		// remove expired items
		for i, shard := range allItemsDB.Shards {
			if shard == nil {
				shard = &mpb.DBItemListShard{Items: map[uint32]*mpb.DBItem{}}
				allItemsDB.Shards[i] = shard
			}
			for _, v := range shard.Items {
				if v.ExpireAt == 0 || v.ExpireAt >= nowUnix {
					continue
				} else {
					delete(shard.Items, v.ItemId)
				}
			}
		}
		// remove expired uitems
		for i, shard := range allItemsDB.UShards {
			if shard == nil {
				shard = &mpb.DBUItemListShard{Items: map[uint64]*mpb.DBItem{}, ItemNums: map[uint32]uint32{}}
				allItemsDB.UShards[i] = shard
			}
			shard.ItemNums = make(map[uint32]uint32)
			for _, v := range shard.Items {
				if v.ExpireAt == 0 || v.ExpireAt >= nowUnix {
					// count valid uItems
					shard.ItemNums[v.ItemId] += 1
				} else {
					delete(shard.Items, v.Uuid)
				}
			}
		}

		// handle del items, just updateList, no showList
		um, uul, vl, err := dao.delItemListFromItemsDB(delItems, nowUnix, allItemsDB)
		if err != nil {
			return err
		}
		for k, v := range um {
			updateItemsMap[k] = v
		}
		updateUItemsList = append(updateUItemsList, uul...)
		vDelList = append(vDelList, vl...)

		// handle add items
		sl, um, uul, vl, err := dao.addItemListInToItemsDB(addItems, nowUnix, 0, allItemsDB) // showList updateMap updateUItemList
		if err != nil {
			return err
		}
		showList = append(showList, sl...)
		for k, v := range um {
			updateItemsMap[k] = v
		}
		updateUItemsList = append(updateUItemsList, uul...)
		vAddList = append(vAddList, vl...)

		if len(updateItemsMap)+len(updateUItemsList) == 0 { // no need to update db
			return nil
		}

		// write back to db
		wValues := make([]interface{}, 0, (itemsShardNum+uItemsShardNum)*2)
		updateShardsKeyMap := make(map[uint32]struct{})
		updateUShardsKeyMap := make(map[uint32]struct{})
		for id := range updateItemsMap {
			updateShardsKeyMap[id%itemsShardNum] = struct{}{}
		}
		for _, item := range updateUItemsList {
			updateUShardsKeyMap[item.ItemId%uItemsShardNum] = struct{}{}
		}
		for k := range updateShardsKeyMap {
			wValues = append(wValues, shardKeys[k], allItemsDB.Shards[k])
		}

		for k := range updateUShardsKeyMap {
			wValues = append(wValues, uShardKeys[k], allItemsDB.UShards[k])
		}
		err = dao.itemDB.HSetObjects(ctx, com.ItemsKey(userId), wValues...)
		if err != nil {
			dao.logger.Error("exchangeItems write db failed",
				zap.Uint64("user_id", userId),
				zap.Any("write_db_data", wValues), // TO DO 为空？？？
				zap.Any("add_items", addItems),
				zap.Any("del_items", delItems),
				zap.Error(err))
			return mpberr.ErrDB
		}
		return nil
	})
	if err != nil {
		dao.logger.Error("exchangeItems write db failed",
			zap.Uint64("user_id", userId),
			zap.Any("add_items", addItems),
			zap.Any("del_items", delItems),
			zap.Error(err))
		return
	}

	for _, v := range updateItemsMap {
		updateList = append(updateList, dao.svc.DBItem2Item(v))
	}

	for _, v := range updateUItemsList {
		updateList = append(updateList, dao.svc.DBItem2Item(v))
	}

	if allItemsDB.DeltaMana != 0 {
		moneyUpdateList, iErr := dao.exchangeMoney(ctx, userId, allItemsDB.DeltaMana)
		if iErr != nil {
			err = iErr
			return
		}
		updateList = append(updateList, moneyUpdateList...)
	}

	return
}

func (dao *itemDAO) addItemInToItemDB(item *mpb.Item, rsc *mpb.ItemRsc, nowUnix int64, deep uint32, allItemsDB *AllItemsDB,
) (showItems []*mpb.Item, updateItems map[uint32]*mpb.DBItem, updateUItems []*mpb.DBItem, err error) {
	if deep > maxDeep {
		dao.logger.Error("addItemToItemDB too deep recursive call:", zap.Uint32("item_id", item.ItemId))
		err = mpberr.ErrConfig
		return
	}

	var oRes *mpb.ItemRsc
	if rsc.OriginId > 0 {
		oRes = dao.rm.getItemRsc(rsc.OriginId)
		if oRes == nil {
			return
		}
	}

	updateItems = make(map[uint32]*mpb.DBItem)
	if rsc.ExpireTime > 0 { // when add a time limited item
		itemId := gutil.If(rsc.OriginId == 0, rsc.ItemId, oRes.ItemId)
		expireTime := int64(item.Num) * rsc.ExpireTime
		shard := allItemsDB.Shards[itemId%itemsShardNum]
		if shard == nil || shard.Items == nil {
			shard = &mpb.DBItemListShard{Items: map[uint32]*mpb.DBItem{}}
			allItemsDB.Shards[itemId%itemsShardNum] = shard
		}
		itemDB, ok := shard.Items[itemId]
		if !ok { // not own
			itemDB = &mpb.DBItem{ItemId: itemId, Num: 1, ExpireAt: expireTime + nowUnix}
			shard.Items[itemId] = itemDB
			showItems = append(showItems, dao.svc.CloneItem(item))
			updateItems[itemDB.ItemId] = itemDB
			return
		}

		// already own
		if itemDB.ExpireAt == 0 { // own permanent item
			showItems = append(showItems, dao.svc.CloneItem(item))
			return
		}

		// own time limited item
		itemDB.ExpireAt = gutil.Max(itemDB.ExpireAt, nowUnix) + expireTime
		showItems = append(showItems, dao.svc.CloneItem(item))
		updateItems[itemDB.ItemId] = itemDB
		return
	}

	// when add a permanent item
	itemId := item.ItemId
	shard := allItemsDB.Shards[itemId%itemsShardNum]
	if shard == nil || shard.Items == nil {
		shard = &mpb.DBItemListShard{Items: map[uint32]*mpb.DBItem{}}
		allItemsDB.Shards[itemId%itemsShardNum] = shard
	}
	itemDB, ok := shard.Items[itemId]
	if !ok { // don't own
		itemDB = &mpb.DBItem{ItemId: itemId}
		shard.Items[itemId] = itemDB
	}
	if itemDB.ExpireAt > 0 {
		itemDB.Num = 0
		itemDB.ExpireAt = 0
	}

	itemDB.Num = item.Num + itemDB.Num
	showItems = append(showItems, &mpb.Item{
		ItemId:  itemId,
		Num:     item.Num,
		BatchId: item.BatchId,
	})
	updateItems[itemId] = itemDB
	return
}

func (dao *itemDAO) addUItemInToItemDB(item *mpb.Item, rsc *mpb.ItemRsc, nowUnix int64, deep uint32, allItemsDB *AllItemsDB,
) (showItems []*mpb.Item, updateItems map[uint32]*mpb.DBItem, updateUItems []*mpb.DBItem, err error) {
	if deep > maxDeep {
		dao.logger.Error("addUItemToItemDB too deep recursive call:", zap.Uint32("item_id", item.ItemId))
		err = mpberr.ErrConfig
		return
	}
	updateItems = make(map[uint32]*mpb.DBItem)
	itemId := gutil.If(rsc.OriginId == 0, rsc.ItemId, rsc.OriginId)
	uShard := allItemsDB.UShards[itemId%uItemsShardNum]
	if uShard == nil || uShard.Items == nil {
		uShard = &mpb.DBUItemListShard{
			Items:    map[uint64]*mpb.DBItem{},
			ItemNums: map[uint32]uint32{},
		}
		allItemsDB.UShards[itemId%uItemsShardNum] = uShard
	}

	// overflow
	addNum := item.Num
	if item.Uuid > 0 {
		_, ok := uShard.Items[item.Uuid]
		if ok {
			err = mpberr.ErrDuplicatedItem
			return
		}
		item.Num = 1 // when uuid exists, num must be one
		itemDB := dao.svc.Item2DBItem(item)
		uShard.Items[itemDB.Uuid] = itemDB
		uShard.ItemNums[itemDB.ItemId] += 1
		showItems = append(showItems, dao.svc.CloneItem(item))
		updateUItems = append(updateUItems, itemDB)
		return
	}

	for i := uint32(0); i < addNum; i++ {
		cItem := dao.svc.CloneItem(item)
		cItem.Num = 1
		cItem.Uuid, err = dao.svc.itemUUIDSF.Next()
		if err != nil {
			return
		}
		cItem.ExpireAt = gutil.If(rsc.ExpireTime == 0, 0, nowUnix+rsc.ExpireTime)
		itemDB := dao.svc.Item2DBItem(cItem)
		itemDB.ItemId = itemId
		uShard.Items[itemDB.Uuid] = itemDB
		uShard.ItemNums[itemDB.ItemId] += 1
		showItems = append(showItems, cItem)
		updateUItems = append(updateUItems, itemDB)
	}
	return
}

func (dao *itemDAO) addItemListInToItemsDB(items []*mpb.Item, nowUnix int64, deep uint32, allItemsDB *AllItemsDB,
) (showItems []*mpb.Item, updateItems map[uint32]*mpb.DBItem, updateUItems []*mpb.DBItem, virtualItems []*mpb.Item,
	err error) {
	if deep > maxDeep {
		dao.logger.Error("addItemListInToItemsDB too deep recursive call:", zap.Any("items", items))
		err = mpberr.ErrConfig
		return
	}
	updateItems = make(map[uint32]*mpb.DBItem)
	for _, item := range items {
		// handle money first
		if isMoney, mana := dao.svc.isMoney(item); isMoney {
			allItemsDB.DeltaMana += int32(mana)
			showItems = append(showItems, dao.svc.CloneItem(item))
			continue
		}

		rsc := dao.rm.getItemRsc(item.ItemId)
		if rsc == nil {
			err = mpberr.ErrConfig
			return
		}

		if dao.svc.isVirtualItem(rsc.ItemType) {
			showItems = append(showItems, item)
			virtualItems = append(virtualItems, item)
			continue
		}

		var sl []*mpb.Item
		var um map[uint32]*mpb.DBItem
		var uul []*mpb.DBItem
		if !rsc.IsUnique {
			sl, um, uul, err = dao.addItemInToItemDB(item, rsc, nowUnix, deep+1, allItemsDB)
		} else {
			sl, um, uul, err = dao.addUItemInToItemDB(item, rsc, nowUnix, deep+1, allItemsDB)
		}
		if err != nil {
			return
		}
		showItems = append(showItems, sl...)
		for k, v := range um {
			updateItems[k] = v
		}
		updateUItems = append(updateUItems, uul...)
	}
	return
}

func (dao *itemDAO) delItemFromItemsDB(item *mpb.Item, res *mpb.ItemRsc, nowUnix int64, allItemsDB *AllItemsDB,
) (updateItem *mpb.DBItem, err error) {
	if res.ExpireTime > 0 { // time limit item not support to delete
		err = mpberr.ErrItemInvalid
		return
	}

	shard := allItemsDB.Shards[item.ItemId%itemsShardNum]
	if shard == nil || shard.Items == nil {
		err = mpberr.ErrNotEnoughItem
		return
	}
	itemDB, ok := shard.Items[item.ItemId]
	if !ok || itemDB.ExpireAt > 0 && itemDB.ExpireAt < nowUnix || itemDB.Num < item.Num {
		err = mpberr.ErrNotEnoughItem
		return
	}
	itemDB.Num -= item.Num
	if itemDB.Num == 0 {
		delete(shard.Items, item.ItemId)
	}
	updateItem = itemDB
	return
}

func (dao *itemDAO) delUItemFromItemsDB(item *mpb.Item, nowUnix int64, allItemsDB *AllItemsDB,
) (updateUItem *mpb.DBItem, err error) {
	shard := allItemsDB.UShards[item.ItemId%uItemsShardNum]
	if shard == nil || shard.Items == nil {
		err = mpberr.ErrNotEnoughItem
		return
	}
	itemDB, ok := shard.Items[item.Uuid]
	if !ok || itemDB.ExpireAt > 0 && itemDB.ExpireAt < nowUnix {
		err = mpberr.ErrNotEnoughItem
		return
	}
	if itemDB.ItemId != item.ItemId {
		err = mpberr.ErrParam
		return
	}
	delete(shard.Items, item.Uuid)
	shard.ItemNums[item.ItemId] -= 1
	itemDB.Num = 0
	updateUItem = itemDB
	return
}

func (dao *itemDAO) delItemListFromItemsDB(items []*mpb.Item, nowUnix int64, allItemsDB *AllItemsDB,
) (updateItems map[uint32]*mpb.DBItem, updateUItems []*mpb.DBItem, virtualItems []*mpb.Item, err error) {
	updateItems = make(map[uint32]*mpb.DBItem)
	for _, item := range items {
		if isMoney, mana := dao.svc.isMoney(item); isMoney {
			allItemsDB.DeltaMana -= int32(mana)
			continue
		}

		rsc := dao.rm.getItemRsc(item.ItemId)
		if rsc != nil {
			err = mpberr.ErrConfig
			return
		}

		if dao.svc.isVirtualItem(rsc.ItemType) {
			virtualItems = append(virtualItems, item)
			continue
		}

		var ui *mpb.DBItem
		var uui *mpb.DBItem
		if !rsc.IsUnique {
			ui, err = dao.delItemFromItemsDB(item, rsc, nowUnix, allItemsDB)
		} else {
			uui, err = dao.delUItemFromItemsDB(item, nowUnix, allItemsDB)
		}
		if err != nil {
			return
		}
		if ui != nil {
			updateItems[ui.ItemId] = ui
		}
		if uui != nil {
			updateUItems = append(updateUItems, uui)
		}
	}
	return
}

func (dao *itemDAO) getBaseEquips(ctx context.Context, userId uint64) (*mpb.DBBaseEquips, error) {
	dbBaseEquips := &mpb.DBBaseEquips{}
	err := dao.itemDB.GetObject(ctx, com.BaseEquipsKey(userId), dbBaseEquips)
	if err != nil && !dao.itemDB.IsErrNil(err) {
		dao.logger.Error("getBaseEquips get failed", zap.Uint64("user_id", userId),
			zap.Error(err))
		return nil, mpberr.ErrDB
	}
	if dbBaseEquips.Equips == nil {
		dbBaseEquips.Equips = dao.svc.initBaseEquips()
		dbBaseEquips.UpgradeStarFailedTimes = make(map[uint32]uint32)
		dbBaseEquips.UpgradeLevelFailedTimes = make(map[uint32]uint32)
	}
	return dbBaseEquips, nil
}

func (dao *itemDAO) upgradeBaseEquips(ctx context.Context, userId uint64, equipType mpb.EItem_BaseEquipType, upStar,
	upLevel bool, curStar, curLevel uint32) (*mpb.DBBaseEquip, bool, error) {
	var err error
	var dbBaseEquips *mpb.DBBaseEquips
	var success bool
	err = dao.rMux.Safely(ctx, com.BaseEquipsKey(userId), func() error {
		dbBaseEquips, err = dao.getBaseEquips(ctx, userId)
		if err != nil {
			return err
		}
		dbBaseEquip := dbBaseEquips.Equips[uint32(equipType)]
		if dbBaseEquip == nil {
			return mpberr.ErrDB
		}
		if upStar {
			if dbBaseEquip.Star != curStar {
				return mpberr.ErrParam
			}
			if dbBaseEquip.Star == dao.rm.getBaseEquipMaxStar(equipType) {
				return mpberr.ErrBaseEquipMaxStar
			}
			starRsc := dao.rm.getBaseEquipStarRsc(equipType, dbBaseEquip.Star)
			if starRsc == nil {
				return mpberr.ErrConfig
			}
			success = dbBaseEquips.UpgradeStarFailedTimes[uint32(equipType)] >= starRsc.ProtectSuccessNum ||
				util.IsPick(starRsc.UpgradeSuccessRate, weight10000)
			if success {
				dbBaseEquip.Star++
				dbBaseEquips.UpgradeStarFailedTimes[uint32(equipType)] = 0
			} else {
				dbBaseEquips.UpgradeStarFailedTimes[uint32(equipType)] += 1
			}
		}
		if upLevel {
			if dbBaseEquip.Level != curLevel {
				return mpberr.ErrParam
			}
			if dbBaseEquip.Level == dao.rm.getBaseEquipMaxLevel(equipType) {
				return mpberr.ErrBaseEquipMaxLevel
			}
			levelRsc := dao.rm.getBaseEquipLevelRsc(equipType, dbBaseEquip.Level)
			if levelRsc == nil {
				return mpberr.ErrConfig
			}
			success = dbBaseEquips.UpgradeLevelFailedTimes[uint32(equipType)] >= levelRsc.ProtectSuccessNum ||
				util.IsPick(levelRsc.UpgradeSuccessRate, weight10000)
			if success {
				dbBaseEquip.Level++
				dbBaseEquips.UpgradeLevelFailedTimes[uint32(equipType)] = 0
			} else {
				dbBaseEquips.UpgradeLevelFailedTimes[uint32(equipType)] += 1
			}
		}

		err = dao.itemDB.SetObject(ctx, com.BaseEquipsKey(userId), dbBaseEquips)
		if err != nil {
			dao.logger.Error("upgradeBaseEquips update db failed", zap.Uint64("user_id", userId),
				zap.Any("db_base_equips", dbBaseEquips), zap.Error(err))
			return err
		}
		return nil
	})
	if err != nil {
		dao.logger.Error("upgradeBaseEquips failed", zap.Uint64("user_id", userId),
			zap.Uint32("equip_type", uint32(equipType)), zap.Bool("up_star", upStar),
			zap.Bool("up_level", upLevel), zap.Error(err))
		return nil, false, err
	}
	return dbBaseEquips.Equips[uint32(equipType)], success, nil
}
