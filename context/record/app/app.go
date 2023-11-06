package app

import (
	"github.com/totsumaru/dd-bot-be/context/record/domain"
	"github.com/totsumaru/dd-bot-be/context/record/gateway"
	"github.com/totsumaru/dd-bot-be/internal/errors"
	"github.com/totsumaru/dd-bot-be/internal/now"
	"gorm.io/gorm"
)

// Upsertのリクエストです
type UpsertRequest struct {
	ServerID  string
	Namespace string
	Key       string
	Value     map[string]string
}

// レコードをUpsertします
func UpsertRecord(tx *gorm.DB, req UpsertRequest) error {
	serverID, err := domain.NewServerID(req.ServerID)
	if err != nil {
		return errors.NewError("サーバーIDを作成できません", err)
	}

	namespace, err := domain.NewNamespace(req.Namespace)
	if err != nil {
		return errors.NewError("ネームスペースを作成できません", err)
	}

	key, err := domain.NewKey(req.Key)
	if err != nil {
		return errors.NewError("キーを作成できません", err)
	}

	record, err := domain.NewRecord(serverID, namespace, key, req.Value, now.NowJST())
	if err != nil {
		return errors.NewError("レコードを作成できません", err)
	}

	gw, err := gateway.NewGateway(tx)
	if err != nil {
		return errors.NewError("ゲートウェイを作成できません", err)
	}

	if err = gw.Upsert(record); err != nil {
		return errors.NewError("レコードを作成または更新できません", err)
	}

	return nil
}

// レコードを削除します
func RemoveRecord(tx *gorm.DB, serverID, namespace, key string) error {
	gw, err := gateway.NewGateway(tx)
	if err != nil {
		return errors.NewError("ゲートウェイを作成できません", err)
	}

	if err = gw.Remove(serverID, namespace, key); err != nil {
		return errors.NewError("レコードを削除できません", err)
	}

	return nil
}

// 条件に一致するレコードを取得します
//
// 取得できない場合はエラーを返します。
func GetRecord(tx *gorm.DB, serverID, namespace, key string) (domain.Record, error) {
	res := domain.Record{}

	gw, err := gateway.NewGateway(tx)
	if err != nil {
		return res, errors.NewError("ゲートウェイを作成できません", err)
	}

	res, err = gw.FindByCondition(serverID, namespace, key)
	if err != nil {
		return res, errors.NewError("レコードを取得できません", err)
	}

	return res, nil
}
