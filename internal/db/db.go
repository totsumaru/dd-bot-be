package db

import (
	"os"

	"github.com/totsumaru/dd-bot-be/internal/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ユーザーデータのテーブルです
type UserData struct {
	ID       string `gorm:"primaryKey;"`
	ServerID string `gorm:"index;"`
	Data     []byte `gorm:"type:jsonb"`
	// ここはドメインには定義されておらず、UNIQUE制約を付けるために使用しています
	// serverID_namespace_keyの形式でユニークになります
	// 例: 1234567890-namespace1-key1
	ServerIDNamespaceKey string `gorm:"uniqueIndex;"`
}

// DBに接続します
func ConnectDB() {
	dialector := postgres.Open(os.Getenv("DB_URL"))
	db, err := gorm.Open(dialector, &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic(errors.NewError("DBに接続できません", err))
	}

	// テーブルが存在していない場合のみテーブルを作成します
	// 存在している場合はスキーマを同期します
	if err = db.AutoMigrate(&UserData{}); err != nil {
		panic(errors.NewError("テーブルのスキーマが一致しません", err))
	}
}
