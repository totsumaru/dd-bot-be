package gateway

import (
	"encoding/json"
	defaultErrors "errors"

	"github.com/totsumaru/dd-bot-be/context/server/domain"
	"github.com/totsumaru/dd-bot-be/internal/db"
	"github.com/totsumaru/dd-bot-be/internal/errors"
	"gorm.io/gorm"
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

// サーバーを作成します
func (g Gateway) Create(server domain.Server) error {
	dbServer, err := castToDBStruct(server)
	if err != nil {
		return errors.NewError("ドメインモデルをDBの構造体に変換できません", err)
	}

	if err = g.tx.Create(&dbServer).Error; err != nil {
		return errors.NewError("サーバーを作成できません", err)
	}

	return nil
}

// サーバーを更新します
func (g Gateway) Update(server domain.Server) error {
	dbServer, err := castToDBStruct(server)
	if err != nil {
		return errors.NewError("ドメインモデルをDBの構造体に変換できません", err)
	}

	if err = g.tx.Save(&dbServer).Error; err != nil {
		return errors.NewError("サーバーを更新できません", err)
	}

	return nil
}

// IDでサーバーを取得します
//
// 取得できない場合はエラーを返します。
func (g Gateway) FindByID(id string) (domain.Server, error) {
	var dbServer db.Server
	if err := g.tx.First(&dbServer, "id = ?", id).Error; err != nil {
		if defaultErrors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Server{}, errors.NewError("サーバーが見つかりません", err)
		}

		return domain.Server{}, errors.NewError("サーバーを取得できません", err)
	}

	res, err := castToDomainModel(dbServer)
	if err != nil {
		return domain.Server{}, errors.NewError("DBの構造体をドメインモデルに変換できません", err)
	}

	return res, nil
}

// ドメインモデルからDBの構造体に変換します
func castToDBStruct(server domain.Server) (db.Server, error) {
	res := db.Server{}

	b, err := json.Marshal(&server)
	if err != nil {
		return res, errors.NewError("Marshalに失敗しました", err)
	}

	res.ID = server.ID()
	res.Data = b

	return res, nil
}

// DBの構造体からドメインモデルに変換します
func castToDomainModel(dbServer db.Server) (domain.Server, error) {
	res := domain.Server{}

	if err := json.Unmarshal(dbServer.Data, &res); err != nil {
		return res, errors.NewError("Unmarshalに失敗しました", err)
	}

	return res, nil
}
