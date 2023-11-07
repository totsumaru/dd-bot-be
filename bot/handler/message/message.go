package message

import (
	"github.com/bwmarrin/discordgo"
)

// メッセージが作成された時のハンドラーです
func MessageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	switch m.Content {
	case "!dd-kv":
		SendKeyValue(s, m)
	case "!dd-csv":
		OutputCSV(s, m)
	}

	Store(s, m)
}
