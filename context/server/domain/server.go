package domain

import (
	"encoding/json"
	"time"

	"github.com/totsumaru/dd-bot-be/internal/errors"
)

// サーバーです
type Server struct {
	id          ServerID
	dbChannelID ChannelID // DB関連の操作を実行するチャンネルです
	apiKey      APIKey
	createdAt   time.Time
	updatedAt   time.Time
}

// サーバーを作成します
func NewServer(
	id ServerID,
	dbChannelID ChannelID,
	createdAt, updatedAt time.Time,
) (Server, error) {
	d := Server{
		id:          id,
		dbChannelID: dbChannelID,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}

	if err := d.Validate(); err != nil {
		return Server{}, err
	}

	return d, nil
}

// APIキーを変更します
func (d *Server) UpdateAPIKey(apiKey APIKey) error {
	d.apiKey = apiKey

	if err := d.Validate(); err != nil {
		return errors.NewError("バリデーションに失敗しました", err)
	}

	return nil
}

// IDを取得します
func (d Server) ID() ServerID {
	return d.id
}

// DBチャンネルIDを取得します
func (d Server) DBChannelID() ChannelID {
	return d.dbChannelID
}

// APIキーを取得します
func (d Server) APIKey() APIKey {
	return d.apiKey
}

// 作成日時を取得します
func (d Server) CreatedAt() time.Time {
	return d.createdAt
}

// 更新日時を取得します
func (d Server) UpdatedAt() time.Time {
	return d.updatedAt
}

// バリデーションを行います
func (d Server) Validate() error {
	return nil
}

// JSONに変換します
func (d Server) MarshalJSON() ([]byte, error) {
	data := struct {
		ID          ServerID  `json:"id"`
		DBChannelID ChannelID `json:"db_channel_id"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
	}{
		ID:          d.id,
		DBChannelID: d.dbChannelID,
		CreatedAt:   d.createdAt,
		UpdatedAt:   d.updatedAt,
	}

	return json.Marshal(data)
}

// JSONからDiscordのIDを復元します
func (d *Server) UnmarshalJSON(b []byte) error {
	data := struct {
		ID          ServerID  `json:"id"`
		DBChannelID ChannelID `json:"db_channel_id"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return errors.NewError("JSONからDiscordIDの復元に失敗しました", err)
	}

	d.id = data.ID
	d.dbChannelID = data.DBChannelID
	d.createdAt = data.CreatedAt
	d.updatedAt = data.UpdatedAt

	return nil
}
