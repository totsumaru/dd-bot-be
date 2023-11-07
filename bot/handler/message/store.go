package message

import (
	"github.com/bwmarrin/discordgo"
	"github.com/totsumaru/dd-bot-be/bot"
	recordApp "github.com/totsumaru/dd-bot-be/context/record/app"
	serverApp "github.com/totsumaru/dd-bot-be/context/server/app"
	"github.com/totsumaru/dd-bot-be/internal/db"
	"github.com/totsumaru/dd-bot-be/internal/errors"
)

// 送信された情報を保存します
func Store(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Bot自身からのメッセージは無視します
	// ※タイトルが`テスト登録用`となっている場合は登録します。
	if m.Author.ID == s.State.User.ID {
		if len(m.Embeds) == 0 {
			return
		}
		if m.Embeds[0].Title != "テスト登録用" {
			return
		}
	}

	// チャンネルがDB専用のチャンネルかどうかを確認します
	ok, err := IsDBChannel(m.GuildID, m.ChannelID)
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
	}

	// valueにはnamespaceとkeyは含めません
	value := map[string]string{}
	for k, v := range kv {
		if k == "namespace" || k == "key" {
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

// チャンネルが指定したDBチャンネルかを判定します
func IsDBChannel(serverID, channelID string) (bool, error) {
	var dbChannelID string

	// メモリストアからサーバーの情報を取得します
	serverStore, ok := bot.ServerMemoryStore.Get(serverID)
	if ok {
		dbChannelID = serverStore.DBChannelID
	} else {
		// サーバーを取得します
		server, err := serverApp.GetServer(db.DB, serverID)
		if err != nil {
			return false, errors.NewError("サーバーを取得できません", err)
		}

		// Storeに登録します
		bot.ServerMemoryStore.Insert(server.ID().String(), bot.ServerData{
			ServerID:    server.ID().String(),
			DBChannelID: server.DBChannelID().String(),
		})

		dbChannelID = server.DBChannelID().String()
	}

	return dbChannelID == channelID, nil
}
