package accountservice

import "gitlab.com/morbackend/mor_services/mpb"

func (svc *AccountService) DBAccountInfo2AccountInfo(in *mpb.DBAccountInfo) *mpb.AccountInfo {
	if in == nil {
		return nil
	}
	return &mpb.AccountInfo{
		Account:    in.Account,
		UserId:     in.UserId,
		Email:      in.Email,
		Nickname:   in.Nickname,
		Icon:       in.Icon,
		WalletAddr: in.WalletAddr,
	}
}
