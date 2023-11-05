package gateway

import (
	"encoding/json"
	defaultErrors "errors"
	"strings"

	"github.com/totsumaru/dd-bot-be/context/user_data/domain"
	"github.com/totsumaru/dd-bot-be/internal/db"
	"github.com/totsumaru/dd-bot-be/internal/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// TODO: FOR UPDATEを使って排他制御を行う

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
func (g Gateway) Upsert(userData domain.UserData) error {
	dbUserData, err := castToDBStruct(userData)
	if err != nil {
		return errors.NewError("ドメインモデルをDBの構造体に変換できません", err)
	}

	// レコードが見つからない場合は新しいレコードを作成
	if err = g.tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "server_id_namespace_key"}},
		UpdateAll: true,
	}).Create(&dbUserData).Error; err != nil {
		return errors.NewError("レコードを作成または更新できません", err)
	}

	return nil
}

// 削除します
func (g Gateway) Remove(serverID, namespace, key string) error {
	joined := strings.Join([]string{serverID, namespace, key}, "-")

	if err := g.tx.Delete(
		&db.UserData{},
		"server_id_namespace_key = ?", joined,
	).Error; err != nil {
		return errors.NewError("レコードを削除できません", err)
	}

	return nil
}

// 条件に一致するものを取得します
//
// 取得できない場合はnilを返します。
func (g Gateway) FindByCondition(serverID, namespace, key string) (*domain.UserData, error) {
	res := &domain.UserData{}

	joined := strings.Join([]string{serverID, namespace, key}, "-")

	var dbUserData db.UserData
	if err := g.tx.First(
		&dbUserData,
		"server_id_namespace_key = ?", joined,
	).Error; err != nil {
		// レコードが見つからない場合はnilを返す
		if defaultErrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, errors.NewError("IDからレコードを取得できません", err)
		}
	}

	return res, nil
}

// サーバーIDに一致する全ての情報を取得します
func (g Gateway) FindAllByServerID(serverID string) ([]domain.UserData, error) {
	res := make([]domain.UserData, 0)

	var dbUserData []db.UserData
	if err := g.tx.Where(
		"server_id = ?", serverID,
	).Find(&dbUserData).Error; err != nil {
		return res, errors.NewError("サーバーIDからレコードを取得できません", err)
	}

	for _, v := range dbUserData {
		var userData domain.UserData
		if err := json.Unmarshal(v.Data, &userData); err != nil {
			return res, errors.NewError("レコードをドメインモデルに変換できません", err)
		}

		res = append(res, userData)
	}

	return res, nil
}

// ドメインモデルからDBの構造体に変換します
func castToDBStruct(userData domain.UserData) (db.UserData, error) {
	res := db.UserData{}

	b, err := json.Marshal(&userData)
	if err != nil {
		return res, errors.NewError("Marshalに失敗しました", err)
	}

	res.ID = userData.ID()
	res.ServerID = userData.ServerID()
	res.Data = b
	res.ServerIDNamespaceKey = strings.Join(
		[]string{userData.ServerID(), userData.Namespace(), userData.Key()},
		"-",
	)

	return res, nil
}
