package app

import (
	"github.com/totsumaru/dd-bot-be/context/server/domain"
	"github.com/totsumaru/dd-bot-be/context/server/gateway"
	"github.com/totsumaru/dd-bot-be/internal/errors"
	"github.com/totsumaru/dd-bot-be/internal/now"
	"gorm.io/gorm"
)

// サーバーを作成します
//
// bot導入時にコールされます。
func CreateServer(tx *gorm.DB, id, channelID string) error {
	serverID, err := domain.NewServerID(id)
	if err != nil {
		return errors.NewError("サーバーIDを作成できません", err)
	}

	chID, err := domain.NewChannelID(channelID)
	if err != nil {
		return errors.NewError("チャンネルIDを作成できません", err)
	}

	s, err := domain.NewServer(serverID, chID, now.NowJST(), now.NowJST())
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

// APIキーを生成します
//
// 初回の生成も更新もこの関数を使います。
func GenerateAPIKey(tx *gorm.DB, id string) (domain.APIKey, error) {
	serverID, err := domain.NewServerID(id)
	if err != nil {
		return domain.APIKey{}, errors.NewError("サーバーIDを作成できません", err)
	}

	gw, err := gateway.NewGateway(tx)
	if err != nil {
		return domain.APIKey{}, errors.NewError("ゲートウェイを作成できません", err)
	}

	// IDでサーバーを取得します
	server, err := gw.FindByIDForUpdate(serverID)
	if err != nil {
		return domain.APIKey{}, errors.NewError("サーバーを取得できません", err)
	}

	newApiKey, err := domain.GenerateAPIKey()
	if err != nil {
		return domain.APIKey{}, errors.NewError("APIキーを生成できません", err)
	}

	if err = server.UpdateAPIKey(newApiKey); err != nil {
		return domain.APIKey{}, errors.NewError("APIキーを更新できません", err)
	}

	if err = gw.Update(server); err != nil {
		return domain.APIKey{}, errors.NewError("サーバーを更新できません", err)
	}

	return newApiKey, nil
}

// サーバーを取得します
//
// 取得できない場合はエラーを返します。
func GetServer(tx *gorm.DB, id string) (domain.Server, error) {
	serverID, err := domain.NewServerID(id)
	if err != nil {
		return domain.Server{}, errors.NewError("サーバーIDを作成できません", err)
	}

	gw, err := gateway.NewGateway(tx)
	if err != nil {
		return domain.Server{}, errors.NewError("ゲートウェイを作成できません", err)
	}

	s, err := gw.FindByID(serverID)
	if err != nil {
		return domain.Server{}, errors.NewError("サーバーを取得できません", err)
	}

	return s, nil
}
