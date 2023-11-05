package domain

import (
	"unicode"
	"unicode/utf8"

	"github.com/totsumaru/dd-bot-be/internal/errors"

	"github.com/google/uuid"
)

// ユーザーのデータです
type UserData struct {
	id        string
	serverID  string
	namespace string
	key       string
	value     map[string]string
}

// ユーザーデータを作成します
func NewUserData(
	id, serverID, namespace, key string, value map[string]string,
) (UserData, error) {
	d := UserData{
		id:        id,
		serverID:  serverID,
		namespace: namespace,
		key:       key,
		value:     value,
	}

	if err := d.Validate(); err != nil {
		return UserData{}, errors.NewError("ユーザーデータを作成できません", err)
	}

	return d, nil
}

// IDを取得します
func (d UserData) ID() string {
	return d.id
}

// サーバーIDを取得します
func (d UserData) ServerID() string {
	return d.serverID
}

// ネームスペースを取得します
func (d UserData) Namespace() string {
	return d.namespace
}

// キーを取得します
func (d UserData) Key() string {
	return d.key
}

// 値を取得します
func (d UserData) Value() map[string]string {
	return d.value
}

// 検証します
func (d UserData) Validate() error {
	// idがUUIDでない場合はエラー
	_, err := uuid.Parse(d.id)
	if err != nil {
		return errors.NewError("IDがUUIDではありません", err)
	}

	// サーバーIDの最大文字数を検証
	if len([]rune(d.serverID)) > 50 {
		return errors.NewError("サーバーIDの文字数を超えています")
	}

	// ネームスペースの最大文字数を検証
	if len([]rune(d.namespace)) > 50 {
		return errors.NewError("ネームスペースの文字数を超えています")
	}
	// ネームスペースで使用できるのは、半角英数字とアンダースコアのみ
	if !isAlphanumeric(d.namespace) {
		return errors.NewError("ネームスペースに使用できない文字が含まれています")
	}

	// キーの最大文字数を検証
	if len([]rune(d.key)) > 100 {
		return errors.NewError("キーの文字数を超えています")
	}

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
