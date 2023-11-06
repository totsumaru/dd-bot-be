package domain

import (
	"encoding/json"

	"github.com/totsumaru/dd-bot-be/internal/errors"
)

// サーバーIDです
type ServerID struct {
	value string
}

// サーバーIDを生成します
func NewServerID(value string) (ServerID, error) {
	s := ServerID{value: value}

	if err := s.validate(); err != nil {
		return s, errors.NewError("検証に失敗しました", err)
	}

	return s, nil
}

// サーバーIDを返します
func (s ServerID) String() string {
	return s.value
}

// サーバーIDが存在しているか確認します
func (s ServerID) IsEmpty() bool {
	return s.value == ""
}

// サーバーIDを検証します
func (s ServerID) validate() error {
	if s.IsEmpty() {
		return errors.NewError("サーバーIDが空です")
	}

	return nil
}

// サーバーIDをJSONに変換します
func (s ServerID) MarshalJSON() ([]byte, error) {
	data := struct {
		Value string `json:"value"`
	}{
		Value: s.value,
	}

	return json.Marshal(data)
}

// JSONからサーバーIDを復元します
func (s *ServerID) UnmarshalJSON(b []byte) error {
	data := struct {
		Value string `json:"value"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return errors.NewError("JSONからServerIDの復元に失敗しました", err)
	}

	s.value = data.Value

	return nil
}
