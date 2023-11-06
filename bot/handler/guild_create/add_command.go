package guild_create

import (
	"github.com/bwmarrin/discordgo"
	"github.com/totsumaru/dd-bot-be/bot"
	"github.com/totsumaru/dd-bot-be/internal/errors"
)

// コマンドを追加します
func AddCommand(s *discordgo.Session, g *discordgo.GuildCreate) {
	var adminPermission int64 = discordgo.PermissionAdministrator

	commands := []discordgo.ApplicationCommand{
		// `/create-api-key`コマンドを追加します
		{
			Name:                     bot.SlashCommand_CreateAPIKey,
			Description:              "[管理者のみ]新しいAPIキーを作成します(過去のキーは使えなくなります)",
			Options:                  []*discordgo.ApplicationCommandOption{},
			DefaultMemberPermissions: &adminPermission,
		},
	}

	for _, command := range commands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, g.Guild.ID, &command)
		if err != nil {
			errors.SendErrMsg(s, errors.NewError("コマンドを登録できません", err))
			return
		}
	}
}
