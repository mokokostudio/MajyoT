// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v4.24.3
// source: grpc_game.proto

package mpb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ReqFight struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId   uint64 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	BossId   uint32 `protobuf:"varint,2,opt,name=boss_id,json=bossId,proto3" json:"boss_id,omitempty"`
	BossUuid uint64 `protobuf:"varint,3,opt,name=boss_uuid,json=bossUuid,proto3" json:"boss_uuid,omitempty"`
}

func (x *ReqFight) Reset() {
	*x = ReqFight{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_game_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReqFight) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReqFight) ProtoMessage() {}

func (x *ReqFight) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_game_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReqFight.ProtoReflect.Descriptor instead.
func (*ReqFight) Descriptor() ([]byte, []int) {
	return file_grpc_game_proto_rawDescGZIP(), []int{0}
}

func (x *ReqFight) GetUserId() uint64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *ReqFight) GetBossId() uint32 {
	if x != nil {
		return x.BossId
	}
	return 0
}

func (x *ReqFight) GetBossUuid() uint64 {
	if x != nil {
		return x.BossUuid
	}
	return 0
}

type ResFight struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Win        bool           `protobuf:"varint,1,opt,name=win,proto3" json:"win,omitempty"`
	BossDie    bool           `protobuf:"varint,2,opt,name=boss_die,json=bossDie,proto3" json:"boss_die,omitempty"`
	Details    []*FightDetail `protobuf:"bytes,3,rep,name=details,proto3" json:"details,omitempty"`
	Awards     *CAwards       `protobuf:"bytes,4,opt,name=awards,proto3" json:"awards,omitempty"`
	EnergyCost uint32         `protobuf:"varint,5,opt,name=energy_cost,json=energyCost,proto3" json:"energy_cost,omitempty"`
	Dmg        uint64         `protobuf:"varint,6,opt,name=dmg,proto3" json:"dmg,omitempty"`
	DmgRate    uint64         `protobuf:"varint,7,opt,name=dmg_rate,json=dmgRate,proto3" json:"dmg_rate,omitempty"`
	HiddenBoss *HiddenBoss    `protobuf:"bytes,8,opt,name=hidden_boss,json=hiddenBoss,proto3" json:"hidden_boss,omitempty"`
}

func (x *ResFight) Reset() {
	*x = ResFight{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_game_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResFight) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResFight) ProtoMessage() {}

func (x *ResFight) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_game_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResFight.ProtoReflect.Descriptor instead.
func (*ResFight) Descriptor() ([]byte, []int) {
	return file_grpc_game_proto_rawDescGZIP(), []int{1}
}

func (x *ResFight) GetWin() bool {
	if x != nil {
		return x.Win
	}
	return false
}

func (x *ResFight) GetBossDie() bool {
	if x != nil {
		return x.BossDie
	}
	return false
}

func (x *ResFight) GetDetails() []*FightDetail {
	if x != nil {
		return x.Details
	}
	return nil
}

func (x *ResFight) GetAwards() *CAwards {
	if x != nil {
		return x.Awards
	}
	return nil
}

func (x *ResFight) GetEnergyCost() uint32 {
	if x != nil {
		return x.EnergyCost
	}
	return 0
}

func (x *ResFight) GetDmg() uint64 {
	if x != nil {
		return x.Dmg
	}
	return 0
}

func (x *ResFight) GetDmgRate() uint64 {
	if x != nil {
		return x.DmgRate
	}
	return 0
}

func (x *ResFight) GetHiddenBoss() *HiddenBoss {
	if x != nil {
		return x.HiddenBoss
	}
	return nil
}

type ResGetEnergy struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Energy   uint32 `protobuf:"varint,1,opt,name=energy,proto3" json:"energy,omitempty"`
	UpdateAt int64  `protobuf:"varint,2,opt,name=update_at,json=updateAt,proto3" json:"update_at,omitempty"`
}

func (x *ResGetEnergy) Reset() {
	*x = ResGetEnergy{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_game_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResGetEnergy) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResGetEnergy) ProtoMessage() {}

func (x *ResGetEnergy) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_game_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResGetEnergy.ProtoReflect.Descriptor instead.
func (*ResGetEnergy) Descriptor() ([]byte, []int) {
	return file_grpc_game_proto_rawDescGZIP(), []int{2}
}

func (x *ResGetEnergy) GetEnergy() uint32 {
	if x != nil {
		return x.Energy
	}
	return 0
}

func (x *ResGetEnergy) GetUpdateAt() int64 {
	if x != nil {
		return x.UpdateAt
	}
	return 0
}

type ReqAddEnergy struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId uint64 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Energy uint32 `protobuf:"varint,2,opt,name=energy,proto3" json:"energy,omitempty"`
}

func (x *ReqAddEnergy) Reset() {
	*x = ReqAddEnergy{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_game_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReqAddEnergy) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReqAddEnergy) ProtoMessage() {}

func (x *ReqAddEnergy) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_game_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReqAddEnergy.ProtoReflect.Descriptor instead.
func (*ReqAddEnergy) Descriptor() ([]byte, []int) {
	return file_grpc_game_proto_rawDescGZIP(), []int{3}
}

func (x *ReqAddEnergy) GetUserId() uint64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *ReqAddEnergy) GetEnergy() uint32 {
	if x != nil {
		return x.Energy
	}
	return 0
}

type ResAddEnergy struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Energy   uint32 `protobuf:"varint,1,opt,name=energy,proto3" json:"energy,omitempty"`
	UpdateAt int64  `protobuf:"varint,2,opt,name=update_at,json=updateAt,proto3" json:"update_at,omitempty"`
}

func (x *ResAddEnergy) Reset() {
	*x = ResAddEnergy{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_game_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResAddEnergy) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResAddEnergy) ProtoMessage() {}

func (x *ResAddEnergy) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_game_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResAddEnergy.ProtoReflect.Descriptor instead.
func (*ResAddEnergy) Descriptor() ([]byte, []int) {
	return file_grpc_game_proto_rawDescGZIP(), []int{4}
}

func (x *ResAddEnergy) GetEnergy() uint32 {
	if x != nil {
		return x.Energy
	}
	return 0
}

func (x *ResAddEnergy) GetUpdateAt() int64 {
	if x != nil {
		return x.UpdateAt
	}
	return 0
}

type ReqGetHiddenBoss struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId   uint64 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	BossUuid uint64 `protobuf:"varint,2,opt,name=boss_uuid,json=bossUuid,proto3" json:"boss_uuid,omitempty"`
	FightCd  int64  `protobuf:"varint,3,opt,name=fight_cd,json=fightCd,proto3" json:"fight_cd,omitempty"`
}

func (x *ReqGetHiddenBoss) Reset() {
	*x = ReqGetHiddenBoss{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_game_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReqGetHiddenBoss) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReqGetHiddenBoss) ProtoMessage() {}

func (x *ReqGetHiddenBoss) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_game_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReqGetHiddenBoss.ProtoReflect.Descriptor instead.
func (*ReqGetHiddenBoss) Descriptor() ([]byte, []int) {
	return file_grpc_game_proto_rawDescGZIP(), []int{5}
}

func (x *ReqGetHiddenBoss) GetUserId() uint64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *ReqGetHiddenBoss) GetBossUuid() uint64 {
	if x != nil {
		return x.BossUuid
	}
	return 0
}

func (x *ReqGetHiddenBoss) GetFightCd() int64 {
	if x != nil {
		return x.FightCd
	}
	return 0
}

type ResGetHiddenBoss struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	HiddenBoss *HiddenBoss `protobuf:"bytes,1,opt,name=hidden_boss,json=hiddenBoss,proto3" json:"hidden_boss,omitempty"`
	Fought     bool        `protobuf:"varint,2,opt,name=fought,proto3" json:"fought,omitempty"`
	FightCd    int64       `protobuf:"varint,3,opt,name=fight_cd,json=fightCd,proto3" json:"fight_cd,omitempty"`
}

func (x *ResGetHiddenBoss) Reset() {
	*x = ResGetHiddenBoss{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_game_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResGetHiddenBoss) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResGetHiddenBoss) ProtoMessage() {}

func (x *ResGetHiddenBoss) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_game_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResGetHiddenBoss.ProtoReflect.Descriptor instead.
func (*ResGetHiddenBoss) Descriptor() ([]byte, []int) {
	return file_grpc_game_proto_rawDescGZIP(), []int{6}
}

func (x *ResGetHiddenBoss) GetHiddenBoss() *HiddenBoss {
	if x != nil {
		return x.HiddenBoss
	}
	return nil
}

func (x *ResGetHiddenBoss) GetFought() bool {
	if x != nil {
		return x.Fought
	}
	return false
}

func (x *ResGetHiddenBoss) GetFightCd() int64 {
	if x != nil {
		return x.FightCd
	}
	return 0
}

var File_grpc_game_proto protoreflect.FileDescriptor

var file_grpc_game_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x67, 0x72, 0x70, 0x63, 0x5f, 0x67, 0x61, 0x6d, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x03, 0x6d, 0x70, 0x62, 0x1a, 0x0c, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x59, 0x0a, 0x08, 0x52, 0x65, 0x71, 0x46, 0x69, 0x67, 0x68, 0x74,
	0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x62, 0x6f, 0x73,
	0x73, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x62, 0x6f, 0x73, 0x73,
	0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x62, 0x6f, 0x73, 0x73, 0x5f, 0x75, 0x75, 0x69, 0x64, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x62, 0x6f, 0x73, 0x73, 0x55, 0x75, 0x69, 0x64, 0x22,
	0x89, 0x02, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x46, 0x69, 0x67, 0x68, 0x74, 0x12, 0x10, 0x0a, 0x03,
	0x77, 0x69, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x03, 0x77, 0x69, 0x6e, 0x12, 0x19,
	0x0a, 0x08, 0x62, 0x6f, 0x73, 0x73, 0x5f, 0x64, 0x69, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x07, 0x62, 0x6f, 0x73, 0x73, 0x44, 0x69, 0x65, 0x12, 0x2a, 0x0a, 0x07, 0x64, 0x65, 0x74,
	0x61, 0x69, 0x6c, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x6d, 0x70, 0x62,
	0x2e, 0x46, 0x69, 0x67, 0x68, 0x74, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x52, 0x07, 0x64, 0x65,
	0x74, 0x61, 0x69, 0x6c, 0x73, 0x12, 0x24, 0x0a, 0x06, 0x61, 0x77, 0x61, 0x72, 0x64, 0x73, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x6d, 0x70, 0x62, 0x2e, 0x43, 0x41, 0x77, 0x61,
	0x72, 0x64, 0x73, 0x52, 0x06, 0x61, 0x77, 0x61, 0x72, 0x64, 0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x65,
	0x6e, 0x65, 0x72, 0x67, 0x79, 0x5f, 0x63, 0x6f, 0x73, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x0a, 0x65, 0x6e, 0x65, 0x72, 0x67, 0x79, 0x43, 0x6f, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03,
	0x64, 0x6d, 0x67, 0x18, 0x06, 0x20, 0x01, 0x28, 0x04, 0x52, 0x03, 0x64, 0x6d, 0x67, 0x12, 0x19,
	0x0a, 0x08, 0x64, 0x6d, 0x67, 0x5f, 0x72, 0x61, 0x74, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x07, 0x64, 0x6d, 0x67, 0x52, 0x61, 0x74, 0x65, 0x12, 0x30, 0x0a, 0x0b, 0x68, 0x69, 0x64,
	0x64, 0x65, 0x6e, 0x5f, 0x62, 0x6f, 0x73, 0x73, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f,
	0x2e, 0x6d, 0x70, 0x62, 0x2e, 0x48, 0x69, 0x64, 0x64, 0x65, 0x6e, 0x42, 0x6f, 0x73, 0x73, 0x52,
	0x0a, 0x68, 0x69, 0x64, 0x64, 0x65, 0x6e, 0x42, 0x6f, 0x73, 0x73, 0x22, 0x43, 0x0a, 0x0c, 0x52,
	0x65, 0x73, 0x47, 0x65, 0x74, 0x45, 0x6e, 0x65, 0x72, 0x67, 0x79, 0x12, 0x16, 0x0a, 0x06, 0x65,
	0x6e, 0x65, 0x72, 0x67, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x65, 0x6e, 0x65,
	0x72, 0x67, 0x79, 0x12, 0x1b, 0x0a, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x5f, 0x61, 0x74,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x41, 0x74,
	0x22, 0x3f, 0x0a, 0x0c, 0x52, 0x65, 0x71, 0x41, 0x64, 0x64, 0x45, 0x6e, 0x65, 0x72, 0x67, 0x79,
	0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x65, 0x6e, 0x65,
	0x72, 0x67, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x65, 0x6e, 0x65, 0x72, 0x67,
	0x79, 0x22, 0x43, 0x0a, 0x0c, 0x52, 0x65, 0x73, 0x41, 0x64, 0x64, 0x45, 0x6e, 0x65, 0x72, 0x67,
	0x79, 0x12, 0x16, 0x0a, 0x06, 0x65, 0x6e, 0x65, 0x72, 0x67, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x06, 0x65, 0x6e, 0x65, 0x72, 0x67, 0x79, 0x12, 0x1b, 0x0a, 0x09, 0x75, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x5f, 0x61, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x75, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x41, 0x74, 0x22, 0x63, 0x0a, 0x10, 0x52, 0x65, 0x71, 0x47, 0x65, 0x74,
	0x48, 0x69, 0x64, 0x64, 0x65, 0x6e, 0x42, 0x6f, 0x73, 0x73, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73,
	0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x75, 0x73, 0x65,
	0x72, 0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x62, 0x6f, 0x73, 0x73, 0x5f, 0x75, 0x75, 0x69, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x62, 0x6f, 0x73, 0x73, 0x55, 0x75, 0x69, 0x64,
	0x12, 0x19, 0x0a, 0x08, 0x66, 0x69, 0x67, 0x68, 0x74, 0x5f, 0x63, 0x64, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x07, 0x66, 0x69, 0x67, 0x68, 0x74, 0x43, 0x64, 0x22, 0x77, 0x0a, 0x10, 0x52,
	0x65, 0x73, 0x47, 0x65, 0x74, 0x48, 0x69, 0x64, 0x64, 0x65, 0x6e, 0x42, 0x6f, 0x73, 0x73, 0x12,
	0x30, 0x0a, 0x0b, 0x68, 0x69, 0x64, 0x64, 0x65, 0x6e, 0x5f, 0x62, 0x6f, 0x73, 0x73, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x6d, 0x70, 0x62, 0x2e, 0x48, 0x69, 0x64, 0x64, 0x65,
	0x6e, 0x42, 0x6f, 0x73, 0x73, 0x52, 0x0a, 0x68, 0x69, 0x64, 0x64, 0x65, 0x6e, 0x42, 0x6f, 0x73,
	0x73, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x6f, 0x75, 0x67, 0x68, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x06, 0x66, 0x6f, 0x75, 0x67, 0x68, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x66, 0x69, 0x67,
	0x68, 0x74, 0x5f, 0x63, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x66, 0x69, 0x67,
	0x68, 0x74, 0x43, 0x64, 0x32, 0xd6, 0x01, 0x0a, 0x0b, 0x47, 0x61, 0x6d, 0x65, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x12, 0x3d, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x48, 0x69, 0x64, 0x64, 0x65,
	0x6e, 0x42, 0x6f, 0x73, 0x73, 0x12, 0x15, 0x2e, 0x6d, 0x70, 0x62, 0x2e, 0x52, 0x65, 0x71, 0x47,
	0x65, 0x74, 0x48, 0x69, 0x64, 0x64, 0x65, 0x6e, 0x42, 0x6f, 0x73, 0x73, 0x1a, 0x15, 0x2e, 0x6d,
	0x70, 0x62, 0x2e, 0x52, 0x65, 0x73, 0x47, 0x65, 0x74, 0x48, 0x69, 0x64, 0x64, 0x65, 0x6e, 0x42,
	0x6f, 0x73, 0x73, 0x12, 0x25, 0x0a, 0x05, 0x46, 0x69, 0x67, 0x68, 0x74, 0x12, 0x0d, 0x2e, 0x6d,
	0x70, 0x62, 0x2e, 0x52, 0x65, 0x71, 0x46, 0x69, 0x67, 0x68, 0x74, 0x1a, 0x0d, 0x2e, 0x6d, 0x70,
	0x62, 0x2e, 0x52, 0x65, 0x73, 0x46, 0x69, 0x67, 0x68, 0x74, 0x12, 0x2e, 0x0a, 0x09, 0x47, 0x65,
	0x74, 0x45, 0x6e, 0x65, 0x72, 0x67, 0x79, 0x12, 0x0e, 0x2e, 0x6d, 0x70, 0x62, 0x2e, 0x52, 0x65,
	0x71, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x1a, 0x11, 0x2e, 0x6d, 0x70, 0x62, 0x2e, 0x52, 0x65,
	0x73, 0x47, 0x65, 0x74, 0x45, 0x6e, 0x65, 0x72, 0x67, 0x79, 0x12, 0x31, 0x0a, 0x09, 0x41, 0x64,
	0x64, 0x45, 0x6e, 0x65, 0x72, 0x67, 0x79, 0x12, 0x11, 0x2e, 0x6d, 0x70, 0x62, 0x2e, 0x52, 0x65,
	0x71, 0x41, 0x64, 0x64, 0x45, 0x6e, 0x65, 0x72, 0x67, 0x79, 0x1a, 0x11, 0x2e, 0x6d, 0x70, 0x62,
	0x2e, 0x52, 0x65, 0x73, 0x41, 0x64, 0x64, 0x45, 0x6e, 0x65, 0x72, 0x67, 0x79, 0x42, 0x07, 0x5a,
	0x05, 0x2e, 0x2f, 0x6d, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_grpc_game_proto_rawDescOnce sync.Once
	file_grpc_game_proto_rawDescData = file_grpc_game_proto_rawDesc
)

func file_grpc_game_proto_rawDescGZIP() []byte {
	file_grpc_game_proto_rawDescOnce.Do(func() {
		file_grpc_game_proto_rawDescData = protoimpl.X.CompressGZIP(file_grpc_game_proto_rawDescData)
	})
	return file_grpc_game_proto_rawDescData
}

var file_grpc_game_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_grpc_game_proto_goTypes = []interface{}{
	(*ReqFight)(nil),         // 0: mpb.ReqFight
	(*ResFight)(nil),         // 1: mpb.ResFight
	(*ResGetEnergy)(nil),     // 2: mpb.ResGetEnergy
	(*ReqAddEnergy)(nil),     // 3: mpb.ReqAddEnergy
	(*ResAddEnergy)(nil),     // 4: mpb.ResAddEnergy
	(*ReqGetHiddenBoss)(nil), // 5: mpb.ReqGetHiddenBoss
	(*ResGetHiddenBoss)(nil), // 6: mpb.ResGetHiddenBoss
	(*FightDetail)(nil),      // 7: mpb.FightDetail
	(*CAwards)(nil),          // 8: mpb.CAwards
	(*HiddenBoss)(nil),       // 9: mpb.HiddenBoss
	(*ReqUserId)(nil),        // 10: mpb.ReqUserId
}
var file_grpc_game_proto_depIdxs = []int32{
	7,  // 0: mpb.ResFight.details:type_name -> mpb.FightDetail
	8,  // 1: mpb.ResFight.awards:type_name -> mpb.CAwards
	9,  // 2: mpb.ResFight.hidden_boss:type_name -> mpb.HiddenBoss
	9,  // 3: mpb.ResGetHiddenBoss.hidden_boss:type_name -> mpb.HiddenBoss
	5,  // 4: mpb.GameService.GetHiddenBoss:input_type -> mpb.ReqGetHiddenBoss
	0,  // 5: mpb.GameService.Fight:input_type -> mpb.ReqFight
	10, // 6: mpb.GameService.GetEnergy:input_type -> mpb.ReqUserId
	3,  // 7: mpb.GameService.AddEnergy:input_type -> mpb.ReqAddEnergy
	6,  // 8: mpb.GameService.GetHiddenBoss:output_type -> mpb.ResGetHiddenBoss
	1,  // 9: mpb.GameService.Fight:output_type -> mpb.ResFight
	2,  // 10: mpb.GameService.GetEnergy:output_type -> mpb.ResGetEnergy
	4,  // 11: mpb.GameService.AddEnergy:output_type -> mpb.ResAddEnergy
	8,  // [8:12] is the sub-list for method output_type
	4,  // [4:8] is the sub-list for method input_type
	4,  // [4:4] is the sub-list for extension type_name
	4,  // [4:4] is the sub-list for extension extendee
	0,  // [0:4] is the sub-list for field type_name
}

func init() { file_grpc_game_proto_init() }
func file_grpc_game_proto_init() {
	if File_grpc_game_proto != nil {
		return
	}
	file_common_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_grpc_game_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReqFight); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_grpc_game_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResFight); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_grpc_game_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResGetEnergy); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_grpc_game_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReqAddEnergy); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_grpc_game_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResAddEnergy); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_grpc_game_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReqGetHiddenBoss); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_grpc_game_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResGetHiddenBoss); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_grpc_game_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_grpc_game_proto_goTypes,
		DependencyIndexes: file_grpc_game_proto_depIdxs,
		MessageInfos:      file_grpc_game_proto_msgTypes,
	}.Build()
	File_grpc_game_proto = out.File
	file_grpc_game_proto_rawDesc = nil
	file_grpc_game_proto_goTypes = nil
	file_grpc_game_proto_depIdxs = nil
}
