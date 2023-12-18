package apiproxy

import (
	"context"
	"encoding/json"
	com "gitlab.com/morbackend/mor_services/common"
	"gitlab.com/morbackend/mor_services/mpb"
	"gitlab.com/morbackend/mor_services/mpberr"
	"gitlab.com/morbackend/mor_services/util"
	"go.uber.org/zap"
	"net/http"
)

type telegramAPIHandler func(ctx context.Context, w http.ResponseWriter, msg *mpb.TGMsgRecv) error

type TelegramManager struct {
	svc         *APIProxyGRPCService
	logger      *zap.Logger
	sendUrl     string
	botTokens   map[string]string
	cmdHandlers map[string]telegramAPIHandler
	cbqHandlers map[string]telegramAPIHandler
	rm          *apiProxyResourceMgr
}

func TelegramManagerGetMe() *TelegramManager {
	tm := APIProxyGRPCGetMe().tgMgr
	if tm.cmdHandlers == nil {
		tm.initCmdHandlers()
	}
	if tm.cbqHandlers == nil {
		tm.initCBQHandlers()
	}
	return tm
}

func newTelegramManager(as *APIProxyGRPCService, sendUrl string, botTokens map[string]string,
) *TelegramManager {
	return &TelegramManager{
		svc:       as,
		logger:    as.logger,
		sendUrl:   sendUrl,
		botTokens: botTokens,
		rm:        as.rm,
	}
}

func (tm *TelegramManager) initCmdHandlers() {
	tm.cmdHandlers = make(map[string]telegramAPIHandler)
	tm.cmdHandlers[com.TGCmd_Echo] = tm.tgAPICmdEcho
	tm.cmdHandlers[com.TGCmd_Menu] = tm.tgAPICmdMenu
	tm.cmdHandlers[com.TGCmd_Game] = tm.tgAPICmdGame
	return
}

func (tm *TelegramManager) initCBQHandlers() {
	tm.cbqHandlers = make(map[string]telegramAPIHandler)
	tm.cbqHandlers[com.TGCBQ_SendGame] = tm.tgAPICBQSendGame
	tm.cbqHandlers[com.TGCBQ_LaunchGame] = tm.tgAPICBQLaunchGame
	return
}

func (tm *TelegramManager) getCmdHandler(cmd string) (telegramAPIHandler, bool) {
	fn, ok := tm.cmdHandlers[cmd]
	return fn, ok
}

func (tm *TelegramManager) getCBQHandler(cmd string) (telegramAPIHandler, bool) {
	fn, ok := tm.cbqHandlers[cmd]
	return fn, ok
}

func (tm *TelegramManager) sendMsgToTelegram(ctx context.Context, bot string, msg []byte) error {
	botToken, ok := tm.botTokens[bot]
	if !ok {
		return mpberr.ErrNoTelegramBot
	}
	urlStr := tm.sendUrl + botToken + "/"
	headers := map[string]string{"Content-Type": "application/json"}
	tm.logger.Debug("sendMsgToTelegram url", zap.String("", urlStr), zap.String("data", string(msg)))
	resp, err := util.HttpsPost(ctx, urlStr, headers, msg)
	if err != nil {
		tm.logger.Info("sendMsgToTelegram failed", zap.Error(err))
		return err
	}
	tm.logger.Info("sendMsgToTelegram result", zap.String("", string(resp)))
	return nil
}

func (tm *TelegramManager) sendCmdReply(ctx context.Context, msg interface{}) error {
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return tm.sendMsgToTelegram(ctx, com.TGBot, msgBytes)
}

func (tm *TelegramManager) tgAPICmdEcho(ctx context.Context, _ http.ResponseWriter, msg *mpb.TGMsgRecv) error {
	if msg == nil || msg.Message == nil || msg.Message.Chat == nil {
		return mpberr.ErrParam
	}
	replyRsc := tm.rm.getTGReplyRsc(com.TGCmd_Echo)
	if replyRsc == nil {
		tm.logger.Error("tgAPICmdMenu cmd reply rsc not found")
		return mpberr.ErrConfig
	}
	replyMsg := &mpb.TGReply{
		Method: replyRsc.Method,
		ChatId: msg.Message.Chat.Id,
		Text:   msg.Message.Text,
	}
	err := tm.sendCmdReply(ctx, replyMsg)
	if err != nil {
		tm.logger.Error("tgAPICmdEcho failed", zap.Error(err))
		return err
	}
	return nil
}

func (tm *TelegramManager) tgAPICmdMenu(ctx context.Context, w http.ResponseWriter, msg *mpb.TGMsgRecv) error {
	if msg == nil || msg.Message == nil || msg.Message.Chat == nil {
		return mpberr.ErrParam
	}
	replyRsc := tm.rm.getTGReplyRsc(com.TGCmd_Menu)
	if replyRsc == nil {
		tm.logger.Error("tgAPICmdMenu cmd reply rsc not found")
		return mpberr.ErrConfig
	}
	replyMsg := &mpb.TGReply{
		Method:      replyRsc.Method,
		ChatId:      msg.Message.Chat.Id,
		Text:        replyRsc.Text,
		Photo:       replyRsc.Photo,
		ReplyMarkup: tm.rm.getTGInlineKeyBoard(com.TGCmd_Menu),
	}
	err := tm.sendCmdReply(ctx, replyMsg)
	if err != nil {
		tm.logger.Error("tgAPICmdMenu failed", zap.Error(err))
		return err
	}
	return nil
}

func (tm *TelegramManager) tgAPICmdGame(ctx context.Context, _ http.ResponseWriter, msg *mpb.TGMsgRecv) error {
	if msg == nil || msg.Message == nil || msg.Message.Chat == nil {
		return mpberr.ErrParam
	}
	replyRsc := tm.rm.getTGReplyRsc(com.TGCmd_Game)
	if replyRsc == nil {
		tm.logger.Error("tgAPICmdGame cmd reply rsc not found")
		return mpberr.ErrConfig
	}
	replyMsg := &mpb.TGReply{
		Method:        replyRsc.Method,
		ChatId:        msg.Message.Chat.Id,
		GameShortName: replyRsc.GameShortName,
	}
	err := tm.sendCmdReply(ctx, replyMsg)
	if err != nil {
		tm.logger.Error("tgAPICmdGame failed", zap.Error(err))
		return err
	}
	return nil
}

func (tm *TelegramManager) tgAPICBQSendGame(ctx context.Context, _ http.ResponseWriter, msg *mpb.TGMsgRecv) error {
	if msg == nil || msg.CallbackQuery == nil ||
		msg.CallbackQuery.Message == nil || msg.CallbackQuery.Message.Chat == nil {
		return mpberr.ErrParam
	}

	answerRsc := tm.rm.getTGCBQAnswerRsc(com.TGCBQ_SendGame)
	if answerRsc != nil {
		answerMsg := &mpb.TGAnswerCallbackQuery{
			Method:          "answerCallbackQuery",
			CallbackQueryId: msg.CallbackQuery.Id,
		}
		err := tm.sendCmdReply(ctx, answerMsg)
		if err != nil {
			tm.logger.Error("tgAPICBQSendGame failed", zap.Error(err))
			return err
		}
	}

	replyRsc := tm.rm.getTGReplyRsc(com.TGCBQ_SendGame)
	if replyRsc == nil {
		tm.logger.Error("tgAPICBQSendGame cmd reply rsc not found")
		return mpberr.ErrConfig
	}

	replyMsg := &mpb.TGReply{
		Method:        replyRsc.Method,
		ChatId:        msg.CallbackQuery.Message.Chat.Id,
		GameShortName: replyRsc.GameShortName,
		ReplyMarkup:   tm.rm.getTGInlineKeyBoard(com.TGCBQ_SendGame),
	}
	err := tm.sendCmdReply(ctx, replyMsg)
	if err != nil {
		tm.logger.Error("tgAPICBQSendGame failed", zap.Error(err))
		return err
	}
	return nil
}

func (tm *TelegramManager) tgAPICBQLaunchGame(ctx context.Context, _ http.ResponseWriter, msg *mpb.TGMsgRecv) error {
	if msg == nil || msg.CallbackQuery == nil || msg.CallbackQuery.From == nil {
		return nil
	}
	answerRsc := tm.rm.getTGCBQAnswerRsc(com.TGCBQ_LaunchGame)
	if answerRsc == nil {
		tm.logger.Error("tgAPICBQLaunchGame cmd reply rsc not found")
		return mpberr.ErrConfig
	}
	token, err := tm.generateLoginToken(ctx, msg.CallbackQuery.From)
	if err != nil {
		return err
	}
	answerMsg := &mpb.TGAnswerCallbackQuery{
		Method:          "answerCallbackQuery",
		CallbackQueryId: msg.CallbackQuery.Id,
		Url:             answerRsc.Url + "?token=" + token,
	}
	err = tm.sendCmdReply(ctx, answerMsg)
	if err != nil {
		tm.logger.Error("tgAPICBQLaunchGame failed", zap.Error(err))
		return err
	}
	return nil
}

func (tm *TelegramManager) generateLoginToken(ctx context.Context, tgUser *mpb.TGUser) (string, error) {
	client, err := com.GetAccountServiceClient(ctx, tm.svc)
	if err != nil {
		return "", err
	}
	res, err := client.GenerateLoginToken(ctx, &mpb.ReqGenerateLoginToken{
		TgId:         tgUser.Id,
		FirstName:    tgUser.FirstName,
		LastName:     tgUser.LastName,
		LanguageCode: tgUser.LanguageCode,
	})
	if err != nil {
		return "", err
	}

	return res.Token, nil
}
