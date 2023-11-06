package interaction

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/totsumaru/dd-bot-be/context/server/app"
	"github.com/totsumaru/dd-bot-be/internal/db"
	"github.com/totsumaru/dd-bot-be/internal/errors"
	"gorm.io/gorm"
)

// APIキーを作成します
func CreateAPIKeyHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	tmpl := `
APIキーを作成しました。
他の人に共有せず、大切に保管してください。

` + "```\n%s\n```" + `
`

	editFunc, err := SendInteractionWaitingMessage(s, i, false, true)
	if err != nil {
		errors.SendErrMsg(s, errors.NewError("Waitingを作成できません", err))
		return
	}

	description := "管理者権限を持っていないため、実行できません"

	// 管理者かどうかを確認します
	if hasAdminPermission(i.Member) {
		err = db.DB.Transaction(func(tx *gorm.DB) error {
			// APIキーを作成します
			apiKey, err := app.GenerateAPIKey(tx, i.GuildID)
			if err != nil {
				return errors.NewError("APIキーを作成できません", err)
			}

			description = fmt.Sprintf(tmpl, apiKey.String())
			return nil
		})
		if err != nil {
			errors.SendErrMsg(s, errors.NewError("APIキーを作成できません", err))
			return
		}
	}

	embed := &discordgo.MessageEmbed{
		Description: description,
	}

	webhook := &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{embed},
	}
	if _, err = editFunc(i.Interaction, webhook); err != nil {
		errors.SendErrMsg(s, errors.NewError("レスポンスを更新できません", err))
		return
	}
}

// 管理者権限を持っているかを確認します
func hasAdminPermission(member *discordgo.Member) bool {
	return member.Permissions&discordgo.PermissionAdministrator == discordgo.PermissionAdministrator
}
