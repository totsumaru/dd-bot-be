package guild_create

import (
	"github.com/bwmarrin/discordgo"
)

// botが追加された時のハンドラーです
func GuildCreateHandler(s *discordgo.Session, g *discordgo.GuildCreate) {
	CreateAdminChannel(s, g)
	AddCommand(s, g)
}
