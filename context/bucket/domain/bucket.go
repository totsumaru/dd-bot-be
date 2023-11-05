package domain

import (
	"unicode"
	"unicode/utf8"

	"github.com/totsumaru/dd-bot-be/internal/errors"

	"github.com/google/uuid"
)

// バケットです
type Bucket struct {
	id        string
	serverID  string
	namespace string
	key       string
	value     map[string]string
}

// バケットを作成します
func NewBucket(
	id, serverID, namespace, key string, value map[string]string,
) (Bucket, error) {
	b := Bucket{
		id:        id,
		serverID:  serverID,
		namespace: namespace,
		key:       key,
		value:     value,
	}

	if err := b.Validate(); err != nil {
		return Bucket{}, errors.NewError("バケットを作成できません", err)
	}

	return b, nil
}

// 検証します
func (b Bucket) Validate() error {
	// idがUUIDでない場合はエラー
	_, err := uuid.Parse(b.id)
	if err != nil {
		return errors.NewError("IDがUUIDではありません", err)
	}

	// サーバーIDの最大文字数を検証
	if len([]rune(b.serverID)) > 50 {
		return errors.NewError("サーバーIDの文字数を超えています")
	}

	// ネームスペースの最大文字数を検証
	if len([]rune(b.namespace)) > 50 {
		return errors.NewError("ネームスペースの文字数を超えています")
	}
	// ネームスペースで使用できるのは、半角英数字とアンダースコアのみ
	if !isAlphanumeric(b.namespace) {
		return errors.NewError("ネームスペースに使用できない文字が含まれています")
	}

	// キーの最大文字数を検証
	if len([]rune(b.key)) > 100 {
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
