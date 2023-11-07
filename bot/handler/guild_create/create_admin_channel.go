package guild_create

import (
	"github.com/bwmarrin/discordgo"
	"github.com/totsumaru/dd-bot-be/context/server/app"
	"github.com/totsumaru/dd-bot-be/internal/db"
	"github.com/totsumaru/dd-bot-be/internal/errors"
	"gorm.io/gorm"
)

// DBの専用チャンネルを作成し、DBに情報を保存します
func CreateAdminChannel(s *discordgo.Session, g *discordgo.GuildCreate) {
	// 既に作成されている場合は終了します
	server, err := app.GetServer(db.DB, g.Guild.ID)
	if err == nil && server.ID().String() != "" {
		return
	}

	// DBの専用チャンネルを作成します
	channel, err := s.GuildChannelCreateComplex(g.Guild.ID, discordgo.GuildChannelCreateData{
		Name: "dd-bot-admin",
		Type: discordgo.ChannelTypeGuildText,
		PermissionOverwrites: []*discordgo.PermissionOverwrite{
			{
				ID:    g.Guild.ID, // @everyone ロール
				Type:  discordgo.PermissionOverwriteTypeRole,
				Allow: 0,
				Deny:  discordgo.PermissionViewChannel, // チャンネルの閲覧をOFFに
			}, {
				ID:    s.State.User.ID, // bot自身のユーザーID
				Type:  discordgo.PermissionOverwriteTypeMember,
				Allow: discordgo.PermissionViewChannel,
				Deny:  0,
			},
		},
	})
	if err != nil {
		errors.SendErrMsg(s, errors.NewError("専用チャンネルを作成できません", err))
	}

	// Tx
	err = db.DB.Transaction(func(tx *gorm.DB) error {
		if err = app.CreateServer(tx, g.Guild.ID, channel.ID); err != nil {
			return errors.NewError("サーバーを作成できません", err)
		}

		return nil
	})
	if err != nil {
		errors.SendErrMsg(s, errors.NewError("サーバーを作成できません", err))
	}
}
