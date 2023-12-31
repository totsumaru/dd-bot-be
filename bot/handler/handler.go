package handler

import (
	"github.com/bwmarrin/discordgo"
	"github.com/totsumaru/dd-bot-be/bot/handler/guild_create"
	"github.com/totsumaru/dd-bot-be/bot/handler/interaction"
	"github.com/totsumaru/dd-bot-be/bot/handler/message"
)

// ハンドラを追加します
func AddHandler(s *discordgo.Session) {
	s.AddHandler(message.MessageCreateHandler)
	s.AddHandler(interaction.InteractionCreateHandler)
	s.AddHandler(guild_create.GuildCreateHandler)
}
