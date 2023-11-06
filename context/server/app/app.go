package app

import (
	"github.com/totsumaru/dd-bot-be/context/server/domain"
	"github.com/totsumaru/dd-bot-be/context/server/gateway"
	"github.com/totsumaru/dd-bot-be/internal/errors"
	"gorm.io/gorm"
)

// サーバーを作成します
func CreateServer(tx *gorm.DB, id, channelID string) error {
	s, err := domain.NewServer(id, channelID)
	if err != nil {
		return errors.NewError("サーバーを作成できません", err)
	}

	gw, err := gateway.NewGateway(tx)
	if err != nil {
		return errors.NewError("ゲートウェイを作成できません", err)
	}

	if err = gw.Create(s); err != nil {
		return errors.NewError("サーバーを作成できません", err)
	}

	return nil
}

// サーバーを取得します
//
// 取得できない場合はエラーを返します。
func GetServer(tx *gorm.DB, id string) (domain.Server, error) {
	gw, err := gateway.NewGateway(tx)
	if err != nil {
		return domain.Server{}, errors.NewError("ゲートウェイを作成できません", err)
	}

	s, err := gw.FindByID(id)
	if err != nil {
		return domain.Server{}, errors.NewError("サーバーを取得できません", err)
	}

	return s, nil
}
