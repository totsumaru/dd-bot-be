package record

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
	recordApp "github.com/totsumaru/dd-bot-be/context/record/app"
	serverApp "github.com/totsumaru/dd-bot-be/context/server/app"
	"github.com/totsumaru/dd-bot-be/internal/errors"
	"gorm.io/gorm"
)

// レコードを取得します
func GetRecord(e *gin.Engine, db *gorm.DB, s *discordgo.Session) {
	e.GET("/record", func(c *gin.Context) {
		guild := c.Query("guild")
		namespace := c.Query("namespace")
		key := c.Query("key")

		if guild == "" || namespace == "" || key == "" {
			c.JSON(400, "guild, namespace, keyを指定してください")
			return
		}

		// サーバーを取得します
		serverRes, err := serverApp.GetServer(db, guild)
		if err != nil {
			c.JSON(500, "サーバーを取得できません")
			return
		}

		// レコードを取得します
		recordRes, err := recordApp.GetRecord(db, guild, namespace, key)
		if err != nil {
			c.JSON(500, "レコードを取得できません")
			return
		}

		fmt.Printf("%+v\n", recordRes)

		// Discordに送信します
		kv := map[string]string{}
		kv["namespace"] = recordRes.Namespace()
		kv["key"] = recordRes.Key()
		for k, v := range recordRes.Value() {
			kv[k] = v
		}

		embed := &discordgo.MessageEmbed{
			Title:  "Response",
			Fields: []*discordgo.MessageEmbedField{},
		}

		// Add fields to the embed
		for k, v := range kv {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:  k,
				Value: v,
			})
		}

		// Send the embed
		msg, err := s.ChannelMessageSendEmbed(serverRes.DBChannelID(), embed)
		if err != nil {
			errors.SendErrMsg(s, errors.NewError("key-valueを送信できません", err))
		}

		messageURL := fmt.Sprintf(
			"https://discord.com/channels/%s/%s/%s",
			guild,
			serverRes.DBChannelID(),
			msg.ID,
		)

		c.JSON(200, messageURL)
	})
}
