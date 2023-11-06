package message

import (
	"github.com/bwmarrin/discordgo"
	"github.com/totsumaru/dd-bot-be/internal"
	"github.com/totsumaru/dd-bot-be/internal/errors"
)

// 検証用に、key-valueを送信する関数を作成します
//
// 管理者のみ使用できます。
func SendKeyValue(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID != internal.UserID().TOTSUMARU {
		return
	}

	fields := map[string]string{
		"namespace": "ns1",
		"key":       "12345",
		"key1":      "value1",
	}

	embed := &discordgo.MessageEmbed{
		Title:  "テスト送信用",
		Fields: []*discordgo.MessageEmbedField{},
	}

	// Add fields to the embed
	for k, v := range fields {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  k,
			Value: v,
		})
	}

	// Send the embed
	_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if err != nil {
		errors.SendErrMsg(s, errors.NewError("key-valueを送信できません", err))
	}
}
