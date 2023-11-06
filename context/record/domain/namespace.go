package domain

import (
	"encoding/json"
	"regexp"

	"github.com/totsumaru/dd-bot-be/internal/errors"
)

// ネームスペースです
type Namespace struct {
	value string
}

// ネームスペースを生成します
func NewNamespace(value string) (Namespace, error) {
	n := Namespace{value: value}

	if err := n.validate(); err != nil {
		return n, errors.NewError("検証に失敗しました", err)
	}

	return n, nil
}

// ネームスペースを返します
func (n Namespace) String() string {
	return n.value
}

// ネームスペースが存在しているか確認します
func (n Namespace) IsEmpty() bool {
	return n.value == ""
}

// ネームスペースを検証します
func (n Namespace) validate() error {
	if n.IsEmpty() {
		return errors.NewError("ネームスペースが空です")
	}

	if len([]rune(n.value)) > 50 {
		return errors.NewError("ネームスペースの最大文字数を超えています")
	}

	if !isAlphanumeric(n.value) {
		return errors.NewError("ネームスペースは英数字とアンダースコアのみで構成されている必要があります")
	}

	return nil
}

// ネームスペースをJSONに変換します
func (n Namespace) MarshalJSON() ([]byte, error) {
	data := struct {
		Value string `json:"value"`
	}{
		Value: n.value,
	}

	return json.Marshal(data)
}

// JSONからネームスペースを復元します
func (n *Namespace) UnmarshalJSON(b []byte) error {
	data := struct {
		Value string `json:"value"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return errors.NewError("JSONからNamespaceの復元に失敗しました", err)
	}

	n.value = data.Value

	return nil
}

// 英数字かアンダースコアのみで構成されているか確認します
func isAlphanumeric(s string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	return re.MatchString(s)
}
