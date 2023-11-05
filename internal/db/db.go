package db

import (
	"os"

	"github.com/totsumaru/dd-bot-be/internal/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

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
	if err = db.AutoMigrate(&Bucket{}); err != nil {
		panic(errors.NewError("テーブルのスキーマが一致しません", err))
	}

	DB = db
}
