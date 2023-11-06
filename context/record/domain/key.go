package domain

import (
	"encoding/json"

	"github.com/totsumaru/dd-bot-be/internal/errors"
)

// キーです
type Key struct {
	value string
}

// キーを生成します
func NewKey(value string) (Key, error) {
	k := Key{value: value}

	if err := k.validate(); err != nil {
		return k, errors.NewError("検証に失敗しました", err)
	}

	return k, nil
}

// キーを返します
func (k Key) String() string {
	return k.value
}

// キーが存在しているか確認します
func (k Key) IsEmpty() bool {
	return k.value == ""
}

// キーを検証します
func (k Key) validate() error {
	if k.IsEmpty() {
		return errors.NewError("キーが空です")
	}

	if len([]rune(k.value)) > 100 {
		return errors.NewError("キーの最大文字数を超えています")
	}

	if !isAlphanumeric(k.value) {
		return errors.NewError("キーは英数字とアンダースコアのみで構成されている必要があります")
	}

	return nil
}

// キーをJSONに変換します
func (k Key) MarshalJSON() ([]byte, error) {
	data := struct {
		Value string `json:"value"`
	}{
		Value: k.value,
	}

	return json.Marshal(data)
}

// JSONからキーを復元します
func (k *Key) UnmarshalJSON(b []byte) error {
	data := struct {
		Value string `json:"value"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return errors.NewError("JSONからKeyの復元に失敗しました", err)
	}

	k.value = data.Value

	return nil
}
