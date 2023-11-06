package domain

import (
	"encoding/json"
	"time"

	"github.com/totsumaru/dd-bot-be/internal/errors"
)

// レコードです
type Record struct {
	serverID  ServerID
	namespace Namespace
	key       Key
	value     map[string]string
	updatedAt time.Time
}

// レコードを作成します
func NewRecord(
	serverID ServerID,
	namespace Namespace,
	key Key,
	value map[string]string,
	updatedAt time.Time,
) (Record, error) {
	d := Record{
		serverID:  serverID,
		namespace: namespace,
		key:       key,
		value:     value,
		updatedAt: updatedAt,
	}

	if err := d.Validate(); err != nil {
		return Record{}, errors.NewError("レコードを作成できません", err)
	}

	return d, nil
}

// サーバーIDを取得します
func (d Record) ServerID() ServerID {
	return d.serverID
}

// ネームスペースを取得します
func (d Record) Namespace() Namespace {
	return d.namespace
}

// キーを取得します
func (d Record) Key() Key {
	return d.key
}

// 値を取得します
func (d Record) Value() map[string]string {
	return d.value
}

// 更新日時を取得します
func (d Record) UpdatedAt() time.Time {
	return d.updatedAt
}

// 検証します
func (d Record) Validate() error {
	if len(d.value) == 0 {
		return errors.NewError("valueの値が空です")
	}

	if d.updatedAt.IsZero() {
		return errors.NewError("updatedの値が空です")
	}

	return nil
}

// JSONに変換します
func (d Record) MarshalJSON() ([]byte, error) {
	data := struct {
		ServerID  ServerID          `json:"server_id"`
		Namespace Namespace         `json:"namespace"`
		Key       Key               `json:"key"`
		Value     map[string]string `json:"value"`
		UpdatedAt time.Time         `json:"updated_at"`
	}{
		ServerID:  d.serverID,
		Namespace: d.namespace,
		Key:       d.key,
		Value:     d.value,
		UpdatedAt: d.updatedAt,
	}

	return json.Marshal(data)
}

// JSONからレコードを復元します
func (d *Record) UnmarshalJSON(b []byte) error {
	data := struct {
		ServerID  ServerID          `json:"server_id"`
		Namespace Namespace         `json:"namespace"`
		Key       Key               `json:"key"`
		Value     map[string]string `json:"value"`
		UpdatedAt time.Time         `json:"updated_at"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return errors.NewError("JSONからレコードの復元に失敗しました", err)
	}

	d.serverID = data.ServerID
	d.namespace = data.Namespace
	d.key = data.Key
	d.value = data.Value
	d.updatedAt = data.UpdatedAt

	return nil
}
