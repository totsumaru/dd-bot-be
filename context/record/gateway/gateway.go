package gateway

import (
	"encoding/json"
	"strings"

	"github.com/totsumaru/dd-bot-be/context/record/domain"
	"github.com/totsumaru/dd-bot-be/internal/db"
	"github.com/totsumaru/dd-bot-be/internal/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// gatewayです
type Gateway struct {
	tx *gorm.DB
}

// gatewayを作成します
func NewGateway(tx *gorm.DB) (Gateway, error) {
	if tx == nil {
		return Gateway{}, errors.NewError("引数が空です")
	}

	res := Gateway{
		tx: tx,
	}

	return res, nil
}

// Upsertします
//
// idはCreateの時のみ使用します。
// updateの時は、idは無視されます。
func (g Gateway) Upsert(record domain.Record) error {
	dbRecord, err := castToDBStruct(record)
	if err != nil {
		return errors.NewError("ドメインモデルをDBの構造体に変換できません", err)
	}

	// レコードが見つからない場合は新しいレコードを作成
	if err = g.tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "server_id_namespace_key"}},
		UpdateAll: true,
	}).Create(&dbRecord).Error; err != nil {
		return errors.NewError("レコードを作成または更新できません", err)
	}

	return nil
}

// 削除します
func (g Gateway) Remove(serverID, namespace, key string) error {
	joined := joinServerIDNamespaceKey(serverID, namespace, key)

	if err := g.tx.Delete(
		&db.Record{},
		"server_id_namespace_key = ?", joined,
	).Error; err != nil {
		return errors.NewError("レコードを削除できません", err)
	}

	return nil
}

// 条件に一致するものを取得します
//
// 取得できない場合はエラーを返します。
func (g Gateway) FindByCondition(serverID, namespace, key string) (domain.Record, error) {
	joined := joinServerIDNamespaceKey(serverID, namespace, key)

	var dbRecord db.Record
	if err := g.tx.First(
		&dbRecord,
		"server_id_namespace_key = ?", joined,
	).Error; err != nil {
		return domain.Record{}, errors.NewError("IDからレコードを取得できません", err)
	}

	res, err := castToDomainModel(dbRecord)
	if err != nil {
		return domain.Record{}, errors.NewError("DBの構造体をドメインモデルに変換できません", err)
	}

	return res, nil
}

// サーバーIDに一致する全ての情報を取得します
func (g Gateway) FindAllByServerID(serverID string) ([]domain.Record, error) {
	res := make([]domain.Record, 0)

	var dbRecord []db.Record
	if err := g.tx.Where(
		"server_id = ?", serverID,
	).Find(&dbRecord).Error; err != nil {
		return res, errors.NewError("サーバーIDからレコードを取得できません", err)
	}

	for _, v := range dbRecord {
		var record domain.Record
		if err := json.Unmarshal(v.Data, &record); err != nil {
			return res, errors.NewError("レコードをドメインモデルに変換できません", err)
		}

		res = append(res, record)
	}

	return res, nil
}

// ドメインモデルからDBの構造体に変換します
func castToDBStruct(record domain.Record) (db.Record, error) {
	res := db.Record{}

	b, err := json.Marshal(&record)
	if err != nil {
		return res, errors.NewError("Marshalに失敗しました", err)
	}

	joined := joinServerIDNamespaceKey(
		record.ServerID(),
		record.Namespace(),
		record.Key(),
	)

	res.ServerID = record.ServerID()
	res.Data = b
	res.ServerIDNamespaceKey = joined

	return res, nil
}

// DBの構造体からドメインモデルに変換します
func castToDomainModel(dbRecord db.Record) (domain.Record, error) {
	res := domain.Record{}

	if err := json.Unmarshal(dbRecord.Data, &res); err != nil {
		return res, errors.NewError("Unmarshalに失敗しました", err)
	}

	return res, nil
}

// ServerIDNamespaceKeyを結合します
func joinServerIDNamespaceKey(serverID, namespace, key string) string {
	return strings.Join([]string{serverID, namespace, key}, "-")
}
