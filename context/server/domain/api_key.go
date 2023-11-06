package domain

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"

	"github.com/totsumaru/dd-bot-be/internal/errors"
)

// APIキーです
//
// 最初は平文で保存しますが、いずれHash化する予定です。
type APIKey struct {
	value string
}

// APIキーを生成します
func GenerateAPIKey() (APIKey, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return APIKey{}, errors.NewError("ランダムな値を作成できません")
	}
	apiKey := base64.URLEncoding.EncodeToString(b)

	res := APIKey{
		value: apiKey,
	}

	if err = res.validate(); err != nil {
		return res, errors.NewError("検証に失敗しました", err)
	}

	return res, nil
}

// APIキーを作成します
func NewAPIKey(value string) (APIKey, error) {
	a := APIKey{value: value}

	if err := a.validate(); err != nil {
		return a, errors.NewError("検証に失敗しました", err)
	}

	return a, nil
}

// APIキーを返します
func (a APIKey) String() string {
	return a.value
}

// APIキーが存在しているか確認します
func (a APIKey) IsEmpty() bool {
	return a.value == ""
}

// APIキーが一致しているか確認します
func (a APIKey) IsMatch(apiKey string) bool {
	return a.value == apiKey
}

// APIキーを検証します
func (a APIKey) validate() error {
	if a.IsEmpty() {
		return errors.NewError("APIキーが空です")
	}

	return nil
}

// APIキーをJSONに変換します
func (a APIKey) MarshalJSON() ([]byte, error) {
	data := struct {
		Value string `json:"value"`
	}{
		Value: a.value,
	}

	return json.Marshal(data)
}

// JSONからAPIキーを復元します
func (a *APIKey) UnmarshalJSON(b []byte) error {
	data := struct {
		Value string `json:"value"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return errors.NewError("JSONからAPIKeyの復元に失敗しました", err)
	}

	a.value = data.Value

	return nil
}
