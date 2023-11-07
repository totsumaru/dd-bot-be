package test

import (
	"github.com/bwmarrin/discordgo"
	"github.com/totsumaru/dd-bot-be/bot"
	"github.com/totsumaru/dd-bot-be/internal/errors"
)

// 検証用に、key-valueを送信する関数を作成します
//
// DBチャンネルでのみ起動します
func SendDeleteReq(s *discordgo.Session, m *discordgo.MessageCreate) {
	ok, err := bot.IsDBChannel(m.GuildID, m.ChannelID)
	if err != nil {
		errors.SendErrMsg(s, errors.NewError("DBチャンネルかを判定できません", err))
		return
	}
	if !ok {
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:  "テスト",
		Fields: []*discordgo.MessageEmbedField{},
	}

	// Add fields to the embed
	for k, v := range TestFieldsDelete {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  k,
			Value: v,
		})
	}

	// Send the embed
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if err != nil {
		errors.SendErrMsg(s, errors.NewError("key-valueを送信できません", err))
	}
}
