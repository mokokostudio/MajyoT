package apiproxy

import (
	"fmt"
	"github.com/oldjon/gutil"

	gcsv "github.com/oldjon/gutil/csv"
	gdm "github.com/oldjon/gutil/dirmonitor"
	"gitlab.com/morbackend/mor_services/mpb"
	"go.uber.org/zap"
)

const (
	//csvSuffix   = ".csv"
	baseCSVPath = "./resources/apiproxy/"

	tgReplyCSV          = "TelegramReply.csv"
	tgInlineKeyBoardCSV = "TelegramInlineKeyboard.csv"
	tgGameCSV           = "TelegramGame.csv"
)

type apiProxyResourceMgr struct {
	logger *zap.Logger
	dm     *gdm.DirMonitor
	mtr    *metrics

	tgReply                map[string]*mpb.TGReplyRsc
	tgCBQAnswer            map[string]*mpb.TGReplyRsc
	tgCmdInlineKeyBoardMap map[string][][]*mpb.TGInlineKeyboardRsc
	tgGameShortNameMap     map[string]*mpb.TGGameRsc
}

func newAPIProxyResourceMgr(logger *zap.Logger, mtr *metrics) (*apiProxyResourceMgr, error) {
	rMgr := &apiProxyResourceMgr{
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

func (rm *apiProxyResourceMgr) load() error {
	var err error
	err = rm.dm.BindAndExec(tgReplyCSV, rm.loadTGReply)
	if err != nil {
		return err
	}
	err = rm.dm.BindAndExec(tgInlineKeyBoardCSV, rm.loadTGInlineKeyboard)
	if err != nil {
		return err
	}
	err = rm.dm.BindAndExec(tgGameCSV, rm.loadTGGame)
	if err != nil {
		return err
	}
	return nil
}

func (rm *apiProxyResourceMgr) watch() error {
	return rm.dm.StartWatch()
}

func (rm *apiProxyResourceMgr) loadTGInlineKeyboard(csvPath string) error {
	datas, err := gcsv.ReadCSV2Array(csvPath)
	if err != nil {
		rm.logger.Error(fmt.Sprintf("load %s failed: %s", csvPath, err.Error()))
		return err
	}
	m := make(map[string][][]*mpb.TGInlineKeyboardRsc)
	for _, data := range datas {
		node := &mpb.TGInlineKeyboardRsc{
			Cmd:          data["cmd"],
			Row:          gcsv.Str2Uint32(data["row"]),
			Text:         data["text"],
			Url:          data["url"],
			CallbackData: data["callbackdata"],
			CallbackGame: gcsv.Str2Bool(data["callbackgame"]),
		}

		l := m[node.Cmd]
		var ok bool
		for i, ll := range l {
			if ll[0].Row == node.Row {
				ok = true
				ll = append(ll, node)
				l[i] = ll
				break
			}
		}
		if !ok {
			l = append(l, []*mpb.TGInlineKeyboardRsc{node})
		}

		m[node.Cmd] = l
		rm.logger.Debug("loadTGInlineKeyboard read:", zap.Any("row", node))
	}

	rm.tgCmdInlineKeyBoardMap = m
	rm.logger.Debug("loadTGInlineKeyboard read finish:", zap.Any("rows", m))

	return nil
}

func (rm *apiProxyResourceMgr) getTGInlineKeyBoard(cmd string) *mpb.TGInlineKeyboardMarkup {
	ret := &mpb.TGInlineKeyboardMarkup{}
	l := rm.tgCmdInlineKeyBoardMap[cmd]
	for _, ll := range l {
		row := make([]*mpb.TGInlineKeyboardButton, 0, 1)
		for _, v := range ll {
			row = append(row, &mpb.TGInlineKeyboardButton{
				Text:         v.Text,
				Url:          v.Url,
				CallbackData: v.CallbackData,
				CallbackGame: gutil.If(v.CallbackGame, &mpb.TGCallbackGame{}, nil),
			})
		}
		ret.InlineKeyboard = append(ret.InlineKeyboard, row)
	}
	return ret
}

func (rm *apiProxyResourceMgr) loadTGReply(csvPath string) error {
	datas, err := gcsv.ReadCSV2Array(csvPath)
	if err != nil {
		rm.logger.Error(fmt.Sprintf("load %s failed: %s", csvPath, err.Error()))
		return err
	}
	m := make(map[string]*mpb.TGReplyRsc)
	cbqam := make(map[string]*mpb.TGReplyRsc)
	for _, data := range datas {
		node := &mpb.TGReplyRsc{
			Cmd:           data["cmd"],
			Method:        data["method"],
			Text:          data["text"],
			Photo:         data["photo"],
			GameShortName: data["gameshortname"],
			Url:           data["url"],
		}
		if node.Method == "answerCallbackQuery" {
			cbqam[node.Cmd] = node
		} else {
			m[node.Cmd] = node
		}
		rm.logger.Debug("loadTGReply read:", zap.Any("row", node))
	}

	rm.tgReply = m
	rm.tgCBQAnswer = cbqam
	rm.logger.Debug("loadTGReply read finish:", zap.Any("rows", m), zap.Any("cbqa_rows", cbqam))
	return nil
}

func (rm *apiProxyResourceMgr) getTGReplyRsc(cmd string) *mpb.TGReplyRsc {
	return rm.tgReply[cmd]
}

func (rm *apiProxyResourceMgr) getTGCBQAnswerRsc(cmd string) *mpb.TGReplyRsc {
	return rm.tgCBQAnswer[cmd]
}

func (rm *apiProxyResourceMgr) getEmailAddrs() []*mpb.EmailAddrRsc {
	return nil
}

func (rm *apiProxyResourceMgr) loadTGGame(csvPath string) error {
	datas, err := gcsv.ReadCSV2Array(csvPath)
	if err != nil {
		rm.logger.Error(fmt.Sprintf("load %s failed: %s", csvPath, err.Error()))
		return err
	}
	m := make(map[string]*mpb.TGGameRsc)
	for _, data := range datas {
		node := &mpb.TGGameRsc{
			GameName:      data["gamename"],
			GameShortName: data["gameshortname"],
			GameUrl:       data["gameurl"],
		}

		m[node.GameShortName] = node
		rm.logger.Debug("loadTGGame read:", zap.Any("row", node))
	}

	rm.tgGameShortNameMap = m
	rm.logger.Debug("loadTGGame read finish:", zap.Any("rows", m))
	return nil
}

func (rm *apiProxyResourceMgr) getTGGameRscByGameShortName(name string) *mpb.TGGameRsc {
	return rm.tgGameShortNameMap[name]
}
