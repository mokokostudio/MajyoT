package mpberr

import (
	"errors"
	"net/http"

	gprotocol "github.com/oldjon/gutil/protocol"
	"gitlab.com/morbackend/mor_services/mpb"
)

var (
	// common error
	ErrOk                         = errors.New(mpb.ErrCode_ERR_OK.String())
	ErrCmd                        = errors.New(mpb.ErrCode_ERR_CMD.String())
	ErrParam                      = errors.New(mpb.ErrCode_ERR_PARAM.String())
	ErrDB                         = errors.New(mpb.ErrCode_ERR_DB.String())
	ErrConfig                     = errors.New(mpb.ErrCode_ERR_CONFIG.String())
	ErrRepeatedRequest            = errors.New(mpb.ErrCode_ERR_REPEATED_REQUEST.String())
	ErrUnsupportedHTTPContentType = errors.New(mpb.ErrCode_ERR_UNSUPPORTED_HTTP_CONTENT_TYPE.String())

	// tcpgateway error
	ErrMsgDecode    = errors.New(mpb.ErrCode_ERR_MSG_DECODE.String())
	ErrTokenVerify  = errors.New(mpb.ErrCode_ERR_TOKEN_VERIFY.String())
	ErrWrongGateway = errors.New(mpb.ErrCode_ERR_WRONG_GATEWAY.String())

	// account
	ErrAccountExist          = errors.New(mpb.ErrCode_ERR_ACCOUNT_EXIST.String())
	ErrAccountNotExist       = errors.New(mpb.ErrCode_ERR_ACCOUNT_NOT_EXIST.String())
	ErrPassword              = errors.New(mpb.ErrCode_ERR_PASSWORD.String())
	ErrEmailAddress          = errors.New(mpb.ErrCode_ERR_EMAIL_ADDRESS.String())
	ErrEmailSendMax          = errors.New(mpb.ErrCode_ERR_EMAIL_SEND_MAX.String())
	ErrEmailBindCode         = errors.New(mpb.ErrCode_ERR_EMAIL_BIND_CODE.String())
	ErrEmailBound            = errors.New(mpb.ErrCode_ERR_EMAIL_BOUND.String())
	ErrEmailNotExist         = errors.New(mpb.ErrCode_ERR_EMAIL_NOT_EXIST.String())
	ErrAccBoundEmail         = errors.New(mpb.ErrCode_ERR_ACC_BOUND_EMAIL.String())
	ErrAptosPublicKey        = errors.New(mpb.ErrCode_ERR_APTOS_PUBLIC_KEY.String())
	ErrAptosVerifySignature  = errors.New(mpb.ErrCode_ERR_APTOS_VERIFY_SIGNATURE.String())
	ErrNewPWSameWithOldPW    = errors.New(mpb.ErrCode_ERR_NEW_PASSWD_SAME_WITH_OLD_PASSWD.String())
	ErrEmailVerificationCode = errors.New(mpb.ErrCode_ERR_EMAIL_VERIFICATION_CODE.String())
	ErrLoginTokenExpired     = errors.New(mpb.ErrCode_ERR_LOGIN_TOKEN_EXPIRED.String())

	// item
	ErrNotEnoughMana     = errors.New(mpb.ErrCode_ERR_NOT_ENOUGH_MANA.String())
	ErrNotEnoughItem     = errors.New(mpb.ErrCode_ERR_NOT_ENOUGH_ITEM.String())
	ErrDuplicatedItem    = errors.New(mpb.ErrCode_ERR_DUPLICATED_ITEM.String())
	ErrItemInvalid       = errors.New(mpb.ErrCode_ERR_ITEM_INVALID.String())
	ErrBaseEquipMaxStar  = errors.New(mpb.ErrCode_ERR_BASE_EQUIP_MAX_STAR.String())
	ErrBaseEquipMaxLevel = errors.New(mpb.ErrCode_ERR_BASE_EQUIP_MAX_LEVEL.String())

	// nft
	ErrParseNFTId = errors.New(mpb.ErrCode_ERR_PARSE_NFT_ID.String())
	ErrNFTTokenId = errors.New(mpb.ErrCode_ERR_NFT_TOKEN_ID.String())
	ErrNFTNoOwner = errors.New(mpb.ErrCode_ERR_NFT_NO_OWNER.String())
	ErrNFTHashId  = errors.New(mpb.ErrCode_ERR_NFT_HASH_ID.String())

	// gm
	ErrAdminAccountOrPasswd = errors.New(mpb.ErrCode_ERR_ADMIN_ACCOUNT_OR_PASSWD.String())

	// apiproxy
	ErrNoTelegramBot      = errors.New(mpb.ErrCode_ERR_NO_TELEGRAM_BOT.String())
	ErrInvalidTelegramCmd = errors.New(mpb.ErrCode_ERR_INVALID_TELEGRAM_CMD.String())

	// game
	ErrEnergyNotEnough      = errors.New(mpb.ErrCode_ERR_ENERGY_NOT_ENOUGH.String())
	ErrHiddenBossFought     = errors.New(mpb.ErrCode_ERR_HIDDEN_BOSS_FOUGHT.String())
	ErrHiddenBossFightCD    = errors.New(mpb.ErrCode_ERR_HIDDEN_BOSS_FIGHT_CD.String())
	ErrHiddenBossExpired    = errors.New(mpb.ErrCode_ERR_HIDDEN_BOSS_EXPIRED.String())
	ErrNFTWeaponNotEquipped = errors.New(mpb.ErrCode_ERR_NFT_WEAPON_NOT_EQUIPPED.String())
	ErrNFTWeaponLevel       = errors.New(mpb.ErrCode_ERR_NFT_WEAPON_LEVEL.String())
	ErrHiddenBossNotFound   = errors.New(mpb.ErrCode_ERR_HIDDEN_BOSS_NOT_FOUND.String())
	ErrHiddenBossDied       = errors.New(mpb.ErrCode_ERR_HIDDEN_BOSS_DIED.String())
)

var HTTPErrMap = map[string]int{
	mpb.ErrCode_ERR_OK.String():                              http.StatusOK,
	mpb.ErrCode_ERR_CMD.String():                             http.StatusBadRequest,
	mpb.ErrCode_ERR_PARAM.String():                           http.StatusBadRequest,
	mpb.ErrCode_ERR_DB.String():                              http.StatusBadRequest,
	mpb.ErrCode_ERR_MSG_DECODE.String():                      http.StatusBadRequest,
	mpb.ErrCode_ERR_TOKEN_VERIFY.String():                    http.StatusBadRequest,
	mpb.ErrCode_ERR_WRONG_GATEWAY.String():                   http.StatusBadRequest,
	mpb.ErrCode_ERR_ACCOUNT_EXIST.String():                   http.StatusBadRequest,
	mpb.ErrCode_ERR_ACCOUNT_NOT_EXIST.String():               http.StatusBadRequest,
	mpb.ErrCode_ERR_PASSWORD.String():                        http.StatusBadRequest,
	mpb.ErrCode_ERR_CONFIG.String():                          http.StatusBadRequest,
	mpb.ErrCode_ERR_EMAIL_ADDRESS.String():                   http.StatusBadRequest,
	mpb.ErrCode_ERR_EMAIL_SEND_MAX.String():                  http.StatusBadRequest,
	mpb.ErrCode_ERR_REPEATED_REQUEST.String():                http.StatusBadRequest,
	mpb.ErrCode_ERR_EMAIL_BIND_CODE.String():                 http.StatusBadRequest,
	mpb.ErrCode_ERR_EMAIL_BOUND.String():                     http.StatusBadRequest,
	mpb.ErrCode_ERR_EMAIL_NOT_EXIST.String():                 http.StatusBadRequest,
	mpb.ErrCode_ERR_ACC_BOUND_EMAIL.String():                 http.StatusBadRequest,
	mpb.ErrCode_ERR_APTOS_PUBLIC_KEY.String():                http.StatusBadRequest,
	mpb.ErrCode_ERR_APTOS_VERIFY_SIGNATURE.String():          http.StatusBadRequest,
	mpb.ErrCode_ERR_NEW_PASSWD_SAME_WITH_OLD_PASSWD.String(): http.StatusBadRequest,
	mpb.ErrCode_ERR_EMAIL_VERIFICATION_CODE.String():         http.StatusBadRequest,
	mpb.ErrCode_ERR_PARSE_NFT_ID.String():                    http.StatusBadRequest,
	mpb.ErrCode_ERR_ADMIN_ACCOUNT_OR_PASSWD.String():         http.StatusBadRequest,
	mpb.ErrCode_ERR_NFT_TOKEN_ID.String():                    http.StatusBadRequest,
	mpb.ErrCode_ERR_NFT_NO_OWNER.String():                    http.StatusBadRequest,
	mpb.ErrCode_ERR_NFT_HASH_ID.String():                     http.StatusBadRequest,
	mpb.ErrCode_ERR_NO_TELEGRAM_BOT.String():                 http.StatusBadRequest,
	mpb.ErrCode_ERR_ENERGY_NOT_ENOUGH.String():               http.StatusBadRequest,
	mpb.ErrCode_ERR_HIDDEN_BOSS_FOUGHT.String():              http.StatusBadRequest,
	mpb.ErrCode_ERR_HIDDEN_BOSS_FIGHT_CD.String():            http.StatusBadRequest,
	mpb.ErrCode_ERR_HIDDEN_BOSS_EXPIRED.String():             http.StatusBadRequest,
	mpb.ErrCode_ERR_NOT_ENOUGH_MANA.String():                 http.StatusBadRequest,
	mpb.ErrCode_ERR_NOT_ENOUGH_ITEM.String():                 http.StatusBadRequest,
	mpb.ErrCode_ERR_DUPLICATED_ITEM.String():                 http.StatusBadRequest,
	mpb.ErrCode_ERR_ITEM_INVALID.String():                    http.StatusBadRequest,
	mpb.ErrCode_ERR_BASE_EQUIP_MAX_STAR.String():             http.StatusBadRequest,
	mpb.ErrCode_ERR_BASE_EQUIP_MAX_LEVEL.String():            http.StatusBadRequest,
	mpb.ErrCode_ERR_NFT_WEAPON_NOT_EQUIPPED.String():         http.StatusBadRequest,
	mpb.ErrCode_ERR_NFT_WEAPON_LEVEL.String():                http.StatusBadRequest,
	mpb.ErrCode_ERR_HIDDEN_BOSS_NOT_FOUND.String():           http.StatusBadRequest,
	mpb.ErrCode_ERR_HIDDEN_BOSS_DIED.String():                http.StatusBadRequest,
	mpb.ErrCode_ERR_INVALID_TELEGRAM_CMD.String():            http.StatusBadRequest,
	mpb.ErrCode_ERR_UNSUPPORTED_HTTP_CONTENT_TYPE.String():   http.StatusBadRequest,
	mpb.ErrCode_ERR_LOGIN_TOKEN_EXPIRED.String():             http.StatusBadRequest,
}

func ErrMsg(fc gprotocol.FrameCoder, err error) []byte {
	errMsg := &mpb.ErrorMsg{}
	errCode, ok := mpb.ErrCode_value[err.Error()]
	if ok {
		errMsg.Error = mpb.ErrCode(errCode)
	} else {
		errMsg.Error = mpb.ErrCode_ERR_UNKNOWN
	}
	data, _ := fc.EncodeMsg(uint8(mpb.MainCmd_Error), uint32(mpb.SubCmd_Error_None), errMsg)
	return data
}
