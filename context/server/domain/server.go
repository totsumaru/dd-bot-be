package domain

import (
	"encoding/json"

	"github.com/totsumaru/dd-bot-be/internal/errors"
)

// サーバーです
type Server struct {
	id          string
	dbChannelID string // DB関連の操作を実行するチャンネルです
}

// サーバーを作成します
func NewServer(id, dbChannelID string) (Server, error) {
	d := Server{
		id:          id,
		dbChannelID: dbChannelID,
	}

	if err := d.Validate(); err != nil {
		return Server{}, err
	}

	return d, nil
}

// IDを取得します
func (d Server) ID() string {
	return d.id
}

// DBチャンネルIDを取得します
func (d Server) DBChannelID() string {
	return d.dbChannelID
}

// バリデーションを行います
func (d Server) Validate() error {
	if d.id == "" {
		return errors.NewError("IDが空です")
	}

	if d.dbChannelID == "" {
		return errors.NewError("DBチャンネルIDが空です")
	}

	return nil
}

// JSONに変換します
func (d Server) MarshalJSON() ([]byte, error) {
	data := struct {
		ID          string `json:"id"`
		DBChannelID string `json:"db_channel_id"`
	}{
		ID:          d.id,
		DBChannelID: d.dbChannelID,
	}

	return json.Marshal(data)
}

// JSONからDiscordのIDを復元します
func (d *Server) UnmarshalJSON(b []byte) error {
	data := struct {
		ID          string `json:"id"`
		DBChannelID string `json:"db_channel_id"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return errors.NewError("JSONからDiscordIDの復元に失敗しました", err)
	}

	d.id = data.ID
	d.dbChannelID = data.DBChannelID

	return nil
}
