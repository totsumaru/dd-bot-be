package message

import (
	"github.com/bwmarrin/discordgo"
	"github.com/totsumaru/dd-bot-be/bot"

	recordApp "github.com/totsumaru/dd-bot-be/context/record/app"

	"github.com/totsumaru/dd-bot-be/internal/db"
	"github.com/totsumaru/dd-bot-be/internal/errors"
	"gorm.io/gorm"
)

// 送信された情報を保存します
func Store(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Bot自身からのメッセージは無視します
	// ※タイトルが`テスト`となっている場合は登録します。
	if m.Author.ID == s.State.User.ID {
		if len(m.Embeds) == 0 {
			return
		}
		if m.Embeds[0].Title != "テスト" {
			return
		}
	}

	// チャンネルがDB専用のチャンネルかどうかを確認します
	ok, err := bot.IsDBChannel(m.GuildID, m.ChannelID)
	if err != nil {
		errors.SendErrMsg(s, errors.NewError("DBチャンネルを確認できません", err))
		return
	}
	if !ok {
		return
	}

	// 埋め込みでは無い場合は無視します
	if len(m.Embeds) == 0 {
		return
	}

	// メッセージを取得します
	kv := map[string]string{}
	for _, embed := range m.Embeds {
		for _, field := range embed.Fields {
			kv[field.Name] = field.Value
		}
	}

	isDelete := false

	// バリデーションを行います
	{
		if kv["namespace"] == "" {
			_, err = s.ChannelMessageSendReply(m.ChannelID, "ERROR: namespaceがありません", m.Reference())
			if err != nil {
				errors.SendErrMsg(s, errors.NewError("返信を送信できません", err))
			}
			return
		}
		if kv["key"] == "" {
			_, err = s.ChannelMessageSendReply(m.ChannelID, "ERROR: keyがありません", m.Reference())
			if err != nil {
				errors.SendErrMsg(s, errors.NewError("返信を送信できません", err))
			}
			return
		}

		switch kv["method"] {
		case "POST":
		case "DELETE":
			isDelete = true
		default:
			_, err = s.ChannelMessageSendReply(m.ChannelID, "ERROR: methodが不正です", m.Reference())
			if err != nil {
				errors.SendErrMsg(s, errors.NewError("返信を送信できません", err))
			}
			return
		}
	}

	// メソッドがDELETEの場合、レコードを削除します
	if isDelete {
		err = db.DB.Transaction(func(tx *gorm.DB) error {
			if err = recordApp.RemoveRecord(tx, m.GuildID, kv["namespace"], kv["key"]); err != nil {
				return errors.NewError("レコードを削除できません", err)
			}
			return nil
		})
		if err != nil {
			_, err = s.ChannelMessageSendReply(m.ChannelID, "ERROR: エラーが発生しました", m.Reference())
			if err != nil {
				errors.SendErrMsg(s, errors.NewError("返信を送信できません", err))
			}
			errors.SendErrMsg(s, errors.NewError("レコードを削除できません", err))
			return
		}

		// リアクションを追加します
		if err = s.MessageReactionAdd(m.ChannelID, m.ID, "💥"); err != nil {
			errors.SendErrMsg(s, errors.NewError("リアクションを追加できません", err))
			return
		}

		return
	}

	// valueにはnamespace,key,methodは含めません
	value := map[string]string{}
	for k, v := range kv {
		if k == "namespace" || k == "key" || k == "method" {
			continue
		}
		value[k] = v
	}

	// メッセージを保存します
	if err = recordApp.UpsertRecord(db.DB, recordApp.UpsertRequest{
		ServerID:  m.GuildID,
		Namespace: kv["namespace"],
		Key:       kv["key"],
		Value:     value,
	}); err != nil {
		_, err = s.ChannelMessageSendReply(m.ChannelID, "ERROR: エラーが発生しました", m.Reference())
		if err != nil {
			errors.SendErrMsg(s, errors.NewError("返信を送信できません", err))
		}
		errors.SendErrMsg(s, errors.NewError("メッセージを保存できません", err))
		return
	}

	// リアクションを追加します
	if err = s.MessageReactionAdd(m.ChannelID, m.ID, "👍"); err != nil {
		errors.SendErrMsg(s, errors.NewError("リアクションを追加できません", err))
		return
	}
}
