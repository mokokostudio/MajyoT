package itemservice

import (
	com "gitlab.com/morbackend/mor_services/common"
	"gitlab.com/morbackend/mor_services/mpb"
	"gitlab.com/morbackend/mor_services/mpberr"
)

func (svc *ItemService) DBItem2Item(in *mpb.DBItem) *mpb.Item {
	if in == nil {
		return nil
	}
	return &mpb.Item{
		ItemId:   in.ItemId,
		Num:      in.Num,
		Uuid:     in.Uuid,
		ExpireAt: in.ExpireAt,
	}
}

func (svc *ItemService) Item2DBItem(in *mpb.Item) *mpb.DBItem {
	if in == nil {
		return nil
	}
	return &mpb.DBItem{
		ItemId:   in.ItemId,
		Num:      in.Num,
		Uuid:     in.Uuid,
		ExpireAt: in.ExpireAt,
	}
}

func (svc *ItemService) CloneItem(in *mpb.Item) *mpb.Item {
	return &mpb.Item{
		ItemId:   in.ItemId,
		Num:      in.Num,
		Uuid:     in.Uuid,
		BatchId:  in.BatchId,
		ExpireAt: in.ExpireAt,
	}
}

func (svc *ItemService) DBBaseEquips2BaseEquips(dbEquips *mpb.DBBaseEquips, withAttrs bool) ([]*mpb.BaseEquip, error) {
	out := make([]*mpb.BaseEquip, 0, len(dbEquips.Equips))
	assemble := func(equipType mpb.EItem_BaseEquipType) error {
		dbEquip := dbEquips.Equips[uint32(equipType)]
		equip := &mpb.BaseEquip{
			EquipType: equipType,
			Star:      dbEquip.Star,
			Level:     dbEquip.Level,
		}
		if withAttrs {
			starRsc := svc.rm.getBaseEquipStarRsc(equipType, dbEquip.Star)
			if starRsc == nil {
				return mpberr.ErrConfig
			}
			levelRsc := svc.rm.getBaseEquipLevelRsc(equipType, dbEquip.Level)
			if levelRsc == nil {
				return mpberr.ErrConfig
			}
			equip.Attrs = com.AddAttrs(starRsc.Attrs, levelRsc.Attrs)
		}
		out = append(out, equip)
		return nil
	}
	err := assemble(mpb.EItem_BaseEquipType_Weapon)
	if err != nil {
		return nil, err
	}
	err = assemble(mpb.EItem_BaseEquipType_Armor)
	if err != nil {
		return nil, err
	}
	err = assemble(mpb.EItem_BaseEquipType_Helmet)
	if err != nil {
		return nil, err
	}
	err = assemble(mpb.EItem_BaseEquipType_Glove)
	if err != nil {
		return nil, err
	}
	err = assemble(mpb.EItem_BaseEquipType_Shoes)
	if err != nil {
		return nil, err
	}

	return out, nil
}
