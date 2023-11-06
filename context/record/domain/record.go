package domain

import (
	"encoding/json"
	"unicode"
	"unicode/utf8"

	"github.com/totsumaru/dd-bot-be/internal/errors"
)

// レコードです
type Record struct {
	serverID  string
	namespace string
	key       string
	value     map[string]string
}

// レコードを作成します
func NewRecord(
	serverID, namespace, key string, value map[string]string,
) (Record, error) {
	d := Record{
		serverID:  serverID,
		namespace: namespace,
		key:       key,
		value:     value,
	}

	if err := d.Validate(); err != nil {
		return Record{}, errors.NewError("レコードを作成できません", err)
	}

	return d, nil
}

// サーバーIDを取得します
func (d Record) ServerID() string {
	return d.serverID
}

// ネームスペースを取得します
func (d Record) Namespace() string {
	return d.namespace
}

// キーを取得します
func (d Record) Key() string {
	return d.key
}

// 値を取得します
func (d Record) Value() map[string]string {
	return d.value
}

// 検証します
func (d Record) Validate() error {
	// サーバーIDを検証
	{
		if d.serverID == "" {
			return errors.NewError("サーバーIDがありません")
		}
		if len([]rune(d.serverID)) > 50 {
			return errors.NewError("サーバーIDの文字数を超えています")
		}
	}

	// ネームスペースを検証
	{
		if d.namespace == "" {
			return errors.NewError("ネームスペースがありません")
		}
		if len([]rune(d.namespace)) > 50 {
			return errors.NewError("ネームスペースの文字数を超えています")
		}
		// ネームスペースで使用できるのは、半角英数字とアンダースコアのみ
		if !isAlphanumeric(d.namespace) {
			return errors.NewError("ネームスペースに使用できない文字が含まれています")
		}
	}

	// キーを検証
	{
		if d.key == "" {
			return errors.NewError("キーがありません")
		}
		if len([]rune(d.key)) > 100 {
			return errors.NewError("キーの文字数を超えています")
		}
	}

	return nil
}

// JSONに変換します
func (d Record) MarshalJSON() ([]byte, error) {
	data := struct {
		ServerID  string            `json:"server_id"`
		Namespace string            `json:"namespace"`
		Key       string            `json:"key"`
		Value     map[string]string `json:"value"`
	}{
		ServerID:  d.serverID,
		Namespace: d.namespace,
		Key:       d.key,
		Value:     d.value,
	}

	return json.Marshal(data)
}

// JSONからレコードを復元します
func (d *Record) UnmarshalJSON(b []byte) error {
	data := struct {
		ServerID  string            `json:"server_id"`
		Namespace string            `json:"namespace"`
		Key       string            `json:"key"`
		Value     map[string]string `json:"value"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return errors.NewError("JSONからレコードの復元に失敗しました", err)
	}

	d.serverID = data.ServerID
	d.namespace = data.Namespace
	d.key = data.Key
	d.value = data.Value

	return nil
}

// 英数字かアンダースコアのみで構成されているか確認します
func isAlphanumeric(s string) bool {
	for len(s) > 0 {
		r, size := utf8.DecodeRuneInString(s)
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return false
		}
		s = s[size:]
	}
	return true
}
