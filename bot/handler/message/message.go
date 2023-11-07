package message

import (
	"github.com/bwmarrin/discordgo"
	"github.com/totsumaru/dd-bot-be/bot/handler/message/test"
)

// メッセージが作成された時のハンドラーです
func MessageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	switch m.Content {
	case "!dd-test-store":
		test.SendStoreReq(s, m)
	case "!dd-test-delete":
		test.SendDeleteReq(s, m)
	case "!dd-csv":
		OutputCSV(s, m)
	}

	Store(s, m)
}
