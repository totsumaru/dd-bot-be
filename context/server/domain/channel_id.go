package domain

import (
	"encoding/json"

	"github.com/totsumaru/dd-bot-be/internal/errors"
)

// チャンネルIDです
type ChannelID struct {
	value string
}

// チャンネルIDを生成します
func NewChannelID(value string) (ChannelID, error) {
	c := ChannelID{value: value}

	if err := c.validate(); err != nil {
		return c, errors.NewError("検証に失敗しました", err)
	}

	return c, nil
}

// チャンネルIDを返します
func (c ChannelID) String() string {
	return c.value
}

// チャンネルIDと一致しているか確認します
func (c ChannelID) Equal(value string) bool {
	return c.value == value
}

// チャンネルIDが存在しているか確認します
func (c ChannelID) IsEmpty() bool {
	return c.value == ""
}

// チャンネルIDを検証します
func (c ChannelID) validate() error {
	if c.IsEmpty() {
		return errors.NewError("チャンネルIDが空です")
	}

	return nil
}

// チャンネルIDをJSONに変換します
func (c ChannelID) MarshalJSON() ([]byte, error) {
	data := struct {
		Value string `json:"value"`
	}{
		Value: c.value,
	}

	return json.Marshal(data)
}

// JSONからチャンネルIDを復元します
func (c *ChannelID) UnmarshalJSON(b []byte) error {
	data := struct {
		Value string `json:"value"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return errors.NewError("JSONからChannelIDの復元に失敗しました", err)
	}

	c.value = data.Value

	return nil
}
