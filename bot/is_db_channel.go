package bot

import (
	serverApp "github.com/totsumaru/dd-bot-be/context/server/app"
	"github.com/totsumaru/dd-bot-be/internal/db"
	"github.com/totsumaru/dd-bot-be/internal/errors"
)

// チャンネルが指定したDBチャンネルかを判定します
func IsDBChannel(serverID, channelID string) (bool, error) {
	var dbChannelID string

	// メモリストアからサーバーの情報を取得します
	serverStore, ok := ServerMemoryStore.Get(serverID)
	if ok {
		dbChannelID = serverStore.DBChannelID
	} else {
		// サーバーを取得します
		server, err := serverApp.GetServer(db.DB, serverID)
		if err != nil {
			return false, errors.NewError("サーバーを取得できません", err)
		}

		// Storeに登録します
		ServerMemoryStore.Insert(server.ID().String(), ServerData{
			ServerID:    server.ID().String(),
			DBChannelID: server.DBChannelID().String(),
		})

		dbChannelID = server.DBChannelID().String()
	}

	return dbChannelID == channelID, nil
}
