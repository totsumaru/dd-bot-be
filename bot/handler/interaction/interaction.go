package interaction

import (
	"github.com/bwmarrin/discordgo"
)

// コマンドが実行された時のハンドラーです
func InteractionCreateHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	// ボタン
	case discordgo.InteractionMessageComponent:
		switch i.MessageComponentData().CustomID {
		}
	// コマンド
	case discordgo.InteractionApplicationCommand:
		switch i.ApplicationCommandData().Name {
		case "create-api-key":
			CreateAPIKeyHandler(s, i)
		}
	// Modal
	case discordgo.InteractionModalSubmit:
	}
}
