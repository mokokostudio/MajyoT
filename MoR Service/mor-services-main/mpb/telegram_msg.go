package mpb

type TGUser struct {
	Id           uint64 `json:"id,omitempty"`
	IsBot        bool   `json:"is_bot,omitempty"`
	FirstName    string `json:"first_name,omitempty"`
	LastName     string `json:"last_name,omitempty"`
	LanguageCode string `json:"language_code,omitempty"`
	Username     string `json:"username,omitempty"`
}

type TGMsgChat struct {
	Id        uint64 `json:"id,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Type      string `json:"type,omitempty"`
}

type TGMsgEntity struct {
	Type   string `json:"type"`
	Offset uint32 `json:"offset"`
	Length uint32 `json:"length"`
	Url    string `json:"url,omitempty"`
}

type TGMsg struct {
	MessageId uint32         `json:"message_id,omitempty"`
	From      *TGUser        `json:"from,omitempty"`
	Chat      *TGMsgChat     `json:"chat,omitempty"`
	Date      int64          `json:"date,omitempty"`
	Text      string         `json:"text,omitempty"`
	Entities  []*TGMsgEntity `json:"entities,omitempty"`
}

type TGMsgRecv struct {
	UpdateId      uint64           `json:"update_id,omitempty"`
	Message       *TGMsg           `json:"message,omitempty"`
	InlineQuery   *TGInlineQuery   `json:"inline_query,omitempty"`
	CallbackQuery *TGCallbackQuery `json:"callback_query,omitempty"`
}

type TGReply struct {
	Method                string                  `json:"method,omitempty"`
	ChatId                uint64                  `json:"chat_id,omitempty"`
	Text                  string                  `json:"text,omitempty"`
	Photo                 string                  `json:"photo,omitempty"`
	ParseMode             string                  `json:"parse_mode,omitempty"`
	Entities              []*TGMsgEntity          `json:"entities,omitempty"`
	GameShortName         string                  `json:"game_short_name,omitempty"`
	DisableWebPagePreview bool                    `json:"disable_web_page_preview"`
	ReplyMarkup           *TGInlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

type TGInlineKeyboardMarkup struct {
	InlineKeyboard [][]*TGInlineKeyboardButton `json:"inline_keyboard,omitempty"`
}

type TGInlineKeyboardButton struct {
	Text         string          `json:"text,omitempty"`
	Url          string          `json:"url,omitempty"`
	CallbackData string          `json:"callback_data,omitempty"`
	CallbackGame *TGCallbackGame `json:"callback_game,omitempty"`
}

type TGInputFile struct {
	FileId string `json:"file_id,omitempty"`
}

type TGInlineQuery struct {
	Id       string      `json:"id"`
	From     *TGUser     `json:"from"`
	Query    string      `json:"query"`
	Offset   string      `json:"offset,omitempty"`
	ChatType string      `json:"chat_type,omitempty"`
	Location *TGLocation `json:"location,omitempty"`
}

type TGLocation struct {
}

type TGCallbackGame struct {
}

type TGCallbackQuery struct {
	Id              string  `json:"id"`
	From            *TGUser `json:"from"`
	Message         *TGMsg  `json:"message,omitempty"`
	InlineMessageId string  `json:"inline_message_id,omitempty"`
	ChatInstance    string  `json:"chat_instance,omitempty"`
	Data            string  `json:"data,omitempty"`
	GameShortName   string  `json:"game_short_name,omitempty"`
}

type TGAnswerCallbackQuery struct {
	Method          string `json:"method,omitempty"`
	CallbackQueryId string `json:"callback_query_id"`
	Text            string `json:"text,omitempty"`
	Url             string `json:"url,omitempty"`
}
