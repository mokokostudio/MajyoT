package itemservice

import (
	"context"
	"gitlab.com/morbackend/mor_services/mpberr"
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
	"gitlab.com/morbackend/mor_services/util"
	etcd "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const (
	baseEquipInitStar  = 1
	baseEquipInitLevel = 1
)

type ItemService struct {
	mpb.UnimplementedItemServiceServer
	name            string
	logger          *zap.Logger
	config          env.ModuleConfig
	etcdClient      *etcd.Client
	host            service.Host
	connMgr         *gxgrpc.ConnManager
	signingMethod   jwt.SigningMethod
	signingDuration time.Duration
	rm              *itemResourceMgr
	kvm             *service.KVMgr
	serverEnv       uint32
	sm              *util.ServiceMetrics
	dao             *itemDAO
	tcpMsgCoder     gprotocol.FrameCoder
	itemUUIDSF      service.Snowflake
	transUUIDSF     service.Snowflake
}

// NewItemService create an ItemService entity
func NewItemService(driver service.ModuleDriver) (gxgrpc.GRPCServer, error) {
	svc := &ItemService{
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

	var err error
	svc.rm, err = newItemResourceMgr(svc.logger, svc.sm)
	if err != nil {
		return nil, err
	}

	redisMux, err := grmux.NewRedisMux(svc.config.SubConfig("redis_mutex"), nil, svc.logger, driver.Tracer())
	if err != nil {
		return nil, err
	}

	itemRedis, err := gdb.NewRedisClientByConfig(svc.config.SubConfig("item_redis"),
		svc.config.GetString("db_marshaller"), driver.Tracer())
	if err != nil {
		return nil, err
	}

	tmpRedis, err := gdb.NewRedisClientByConfig(svc.config.SubConfig("tmp_redis"),
		svc.config.GetString("db_marshaller"), driver.Tracer())
	if err != nil {
		return nil, err
	}

	svc.dao = newItemDAO(svc, redisMux, itemRedis, tmpRedis)

	svc.serverEnv = uint32(svc.config.GetInt64("server_env"))
	svc.tcpMsgCoder = gprotocol.NewFrameCoder(svc.config.GetString("protocol_code"))

	return svc, err
}

func (svc *ItemService) Register(grpcServer *grpc.Server) {
	mpb.RegisterItemServiceServer(grpcServer, svc)
}

func (svc *ItemService) Serve(ctx context.Context) error {
	var err error
	svc.itemUUIDSF, err = svc.host.Snowflake(ctx, com.SnowflakeItemUUID, service.SnowflakeType_Default)
	if err != nil {
		svc.logger.Error("failed to create snowflake", zap.Error(err))
		return err
	}
	svc.transUUIDSF, err = svc.host.Snowflake(ctx, com.SnowflakeTransactionUUID, service.SnowflakeType_Default)
	if err != nil {
		svc.logger.Error("failed to create snowflake", zap.Error(err))
		return err
	}
	<-ctx.Done()
	return ctx.Err()
}

func (svc *ItemService) Logger() *zap.Logger {
	return svc.logger
}

func (svc *ItemService) ConnMgr() *gxgrpc.ConnManager {
	return svc.connMgr
}

func (svc *ItemService) Name() string {
	return svc.name
}

func (svc *ItemService) GetEquips(ctx context.Context, req *mpb.ReqUserId) (*mpb.ResGetEquips, error) {
	dbBaseEquips, err := svc.dao.getBaseEquips(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	// TODO(fishillte):  nft equips
	res := &mpb.ResGetEquips{}
	res.BaseEquips, err = svc.DBBaseEquips2BaseEquips(dbBaseEquips, true)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *ItemService) BatchAddItems(ctx context.Context, req *mpb.ReqBatchAddItems) (*mpb.ResBatchAddItems, error) {
	var userExchangeItems = make(map[uint64]*mpb.ReqExchangeItems)
	for userId, items := range req.AddItems {
		userExchangeItems[userId] = &mpb.ReqExchangeItems{
			UserId:         userId,
			AddItems:       items.Items,
			TransReason:    req.TransReason,
			TransSubReason: req.TransSubReason,
		}
	}
	for userId, mana := range req.AddManas {
		node := userExchangeItems[userId]
		if node == nil {
			node = &mpb.ReqExchangeItems{
				UserId:         userId,
				TransReason:    req.TransReason,
				TransSubReason: req.TransSubReason,
			}
			userExchangeItems[userId] = node
		}
		node.DeltaMana = int32(mana)
	}

	res := &mpb.ResBatchAddItems{
		AddItems:    make(map[uint64]*mpb.Items),
		UpdateItems: make(map[uint64]*mpb.Items),
	}

	for userId, v := range userExchangeItems {
		addItems, _, updateItems, err := svc.exchangeItems(ctx, v)
		if err != nil {
			svc.logger.Error("BatchAddItems exchangeItems failed", zap.Uint64("user_id", userId),
				zap.Any("add_items", v), zap.Error(err))
			continue
		}
		res.AddItems[userId] = &mpb.Items{Items: addItems}
		res.UpdateItems[userId] = &mpb.Items{Items: updateItems}
	}

	return res, nil
}

func (svc *ItemService) GetItems(ctx context.Context, req *mpb.ReqUserId) (*mpb.ResGetItems, error) {
	items, err := svc.dao.getItems(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	wallet, err := svc.dao.getWallet(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	dbBaseEquips, err := svc.dao.getBaseEquips(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	// TODO(fishillte):  nft equips

	res := &mpb.ResGetItems{Items: make([]*mpb.Item, 0, len(items)), Mana: wallet.Mana}
	for _, v := range items {
		res.Items = append(res.Items, svc.DBItem2Item(v))
	}
	res.BaseEquips, err = svc.DBBaseEquips2BaseEquips(dbBaseEquips, true)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (svc *ItemService) ExchangeItems(ctx context.Context, req *mpb.ReqExchangeItems) (*mpb.ResExchangeItems, error) {
	addItems, delItems, updateItems, err := svc.exchangeItems(ctx, req)
	if err != nil {
		return nil, err
	}
	return &mpb.ResExchangeItems{
		AddItems:    addItems,
		DelItems:    delItems,
		UpdateItems: updateItems,
	}, nil
}

func (svc *ItemService) exchangeItems(ctx context.Context, req *mpb.ReqExchangeItems) (addItems, delItems,
	updateItems []*mpb.Item, err error) {
	// pre check
	nowUnix := time.Now().Unix()
	addItems = make([]*mpb.Item, 0, len(req.AddItems))
	for _, v := range req.AddItems {
		if v.Num == 0 {
			continue
		}
		rsc := svc.rm.getItemRsc(v.ItemId)
		if rsc == nil {
			return nil, nil, nil, mpberr.ErrConfig
		}
		addItems = append(addItems, v)
	}

	for _, v := range req.DelItems {
		rsc := svc.rm.getItemRsc(v.ItemId)
		if rsc == nil {
			return nil, nil, nil, mpberr.ErrConfig
		}
	}

	if req.DeltaMana > 0 {
		addItems = append(addItems, &mpb.Item{ItemId: uint32(mpb.EItem_ItemId_Mana), Num: uint32(req.DeltaMana)})
	} else if req.DeltaMana < 0 {
		req.DelItems = append(req.DelItems,
			&mpb.Item{ItemId: uint32(mpb.EItem_ItemId_Mana), Num: uint32(-req.DeltaMana)})
	}

	var vAddItems, vDelItems []*mpb.Item
	addItems, updateItems, vAddItems, vDelItems, err = svc.dao.exchangeItems(ctx, req.UserId, addItems, req.DelItems,
		nowUnix)
	if err != nil {
		return nil, nil, nil, err
	}

	vUpdateItems, err := svc.handleVirtualItems(ctx, req.UserId, vAddItems, vDelItems)
	if err != nil {
		return nil, nil, nil, err
	}
	updateItems = append(updateItems, vUpdateItems...)
	// TODO log transaction
	return addItems, delItems, updateItems, nil
}

func (svc *ItemService) isAllMoney(items []*mpb.Item) (allMoney bool, mana uint32) {
	allMoney = true
	for _, v := range items {
		money, m := svc.isMoney(v)
		allMoney = allMoney && money
		mana += m
	}
	return
}

func (svc *ItemService) isMoney(item *mpb.Item) (money bool, mana uint32) {
	money = true
	switch mpb.EItem_ItemId(item.ItemId) {
	case mpb.EItem_ItemId_Mana:
		mana = item.Num
		return
	default:
		money = false
		return
	}
}

func (svc *ItemService) initBaseEquips() map[uint32]*mpb.DBBaseEquip {
	m := map[uint32]*mpb.DBBaseEquip{}
	m[uint32(mpb.EItem_BaseEquipType_Weapon)] = &mpb.DBBaseEquip{
		EquipType: mpb.EItem_BaseEquipType_Weapon, Star: baseEquipInitStar, Level: baseEquipInitLevel}
	m[uint32(mpb.EItem_BaseEquipType_Armor)] = &mpb.DBBaseEquip{
		EquipType: mpb.EItem_BaseEquipType_Armor, Star: baseEquipInitStar, Level: baseEquipInitLevel}
	m[uint32(mpb.EItem_BaseEquipType_Helmet)] = &mpb.DBBaseEquip{
		EquipType: mpb.EItem_BaseEquipType_Helmet, Star: baseEquipInitStar, Level: baseEquipInitLevel}
	m[uint32(mpb.EItem_BaseEquipType_Glove)] = &mpb.DBBaseEquip{
		EquipType: mpb.EItem_BaseEquipType_Glove, Star: baseEquipInitStar, Level: baseEquipInitLevel}
	m[uint32(mpb.EItem_BaseEquipType_Shoes)] = &mpb.DBBaseEquip{
		EquipType: mpb.EItem_BaseEquipType_Shoes, Star: baseEquipInitStar, Level: baseEquipInitLevel}
	return m
}

func (svc *ItemService) isVirtualItem(itemType mpb.EItem_Type) bool {
	return itemType == mpb.EItem_Type_Energy
}

func (svc *ItemService) handleVirtualItems(ctx context.Context, userId uint64, vAddItems, vDelItems []*mpb.Item,
) (vUpdateItems []*mpb.Item, err error) {
	if len(vAddItems)+len(vDelItems) == 0 {
		return nil, nil
	}
	// merge items by item id
	var vAddNumMap = make(map[uint32]uint32)
	for _, item := range vAddItems {
		vAddNumMap[item.ItemId] += item.Num
	}

	for itemId, num := range vAddNumMap {
		rsc := svc.rm.getItemRsc(itemId)
		if rsc == nil {
			return nil, mpberr.ErrConfig
		}
		switch rsc.ItemType {
		case mpb.EItem_Type_Energy:
			uItem, err := svc.handAddGameEnergy(ctx, userId, itemId, num)
			if err != nil {
				return nil, err
			}
			vUpdateItems = append(vUpdateItems, uItem)
		default:
		}
	}
	return
}

func (svc *ItemService) handAddGameEnergy(ctx context.Context, userId uint64, itemId, energy uint32) (*mpb.Item, error) {
	client, err := com.GetGameServiceClient(ctx, svc)
	if err != nil {
		return nil, err
	}
	res, err := client.AddEnergy(ctx, &mpb.ReqAddEnergy{UserId: userId, Energy: energy})
	if err != nil {
		svc.logger.Error("handAddGameEnergy failed", zap.Uint64("user_id", userId),
			zap.Uint32("energy", energy), zap.Error(err))
		return nil, err
	}
	return &mpb.Item{ItemId: itemId, Num: res.Energy}, nil
}

func (svc *ItemService) UpgradeBaseEquipStar(ctx context.Context, req *mpb.ReqUpgradeBaseEquipStar) (
	*mpb.ResUpgradeBaseEquipStar, error) {
	// check and sub consume
	starRsc := svc.rm.getBaseEquipStarRsc(req.EquiptType, req.CurStar)
	if starRsc == nil {
		return nil, mpberr.ErrParam
	}
	_, _, _, err := svc.exchangeItems(ctx, &mpb.ReqExchangeItems{
		DelItems:       starRsc.UpgradeConsumeItems,
		TransReason:    mpb.EItem_TransReason_UpgradeBaseEquipStar,
		TransSubReason: uint64(req.EquiptType),
	})
	if err != nil {
		return nil, err
	}
	// check upgrade
	dbBaseEquip, success, err := svc.dao.upgradeBaseEquips(ctx, req.UserId, req.EquiptType, true, false,
		req.CurStar, 0)
	if err != nil {
		svc.logger.Error("UpgradeBaseEquipStar handel upgrade failed", zap.Uint64("user_id", req.UserId),
			zap.Any("base_equip_type", req.EquiptType), zap.Uint32("cur_star_in_params", req.CurStar),
			zap.Error(err))
		if _, _, _, err := svc.exchangeItems(ctx, &mpb.ReqExchangeItems{
			AddItems:       starRsc.UpgradeConsumeItems,
			TransReason:    mpb.EItem_TransReason_UpgradeBaseEquipStarRollback,
			TransSubReason: uint64(req.EquiptType),
		}); err != nil {
			svc.logger.Error("UpgradeBaseEquipStar rollback consume failed",
				zap.Uint64("user_id", req.UserId),
				zap.Any("base_equip_type", req.EquiptType),
				zap.Uint32("cur_star_in_params", req.CurStar),
				zap.Error(err))
		}
		return nil, err
	}

	return &mpb.ResUpgradeBaseEquipStar{
		Success: success,
		NewStar: dbBaseEquip.Star,
	}, nil
}

func (svc *ItemService) UpgradeBaseEquipLevel(ctx context.Context, req *mpb.ReqUpgradeBaseEquipLevel) (*mpb.ResUpgradeBaseEquipLevel, error) {
	// check and sub consume
	levelRsc := svc.rm.getBaseEquipLevelRsc(req.EquiptType, req.CurLevel)
	if levelRsc == nil {
		return nil, mpberr.ErrParam
	}
	_, _, _, err := svc.exchangeItems(ctx, &mpb.ReqExchangeItems{
		DelItems:       levelRsc.UpgradeConsumeItems,
		TransReason:    mpb.EItem_TransReason_UpgradeBaseEquipLevel,
		TransSubReason: uint64(req.EquiptType),
	})
	if err != nil {
		return nil, err
	}
	// check upgrade
	dbBaseEquip, success, err := svc.dao.upgradeBaseEquips(ctx, req.UserId, req.EquiptType, false, true,
		0, req.CurLevel)
	if err != nil {
		svc.logger.Error("UpgradeBaseEquipLevel handel upgrade failed", zap.Uint64("user_id", req.UserId),
			zap.Any("base_equip_type", req.EquiptType), zap.Uint32("cur_level_in_params", req.CurLevel),
			zap.Error(err))
		if _, _, _, err := svc.exchangeItems(ctx, &mpb.ReqExchangeItems{
			AddItems:       levelRsc.UpgradeConsumeItems,
			TransReason:    mpb.EItem_TransReason_UpgradeBaseEquipLevelRollback,
			TransSubReason: uint64(req.EquiptType),
		}); err != nil {
			svc.logger.Error("UpgradeBaseEquipLevel rollback consume failed",
				zap.Uint64("user_id", req.UserId),
				zap.Any("base_equip_type", req.EquiptType),
				zap.Uint32("cur_level_in_params", req.CurLevel),
				zap.Error(err))
		}
		return nil, err
	}

	return &mpb.ResUpgradeBaseEquipLevel{
		Success:  success,
		NewLevel: dbBaseEquip.Level,
	}, nil
}
