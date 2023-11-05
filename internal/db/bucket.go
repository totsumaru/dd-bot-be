package db

import (
	defaultErrors "errors"
	"strings"

	"github.com/totsumaru/dd-bot-be/internal/errors"
	"gorm.io/gorm"
)

// バケットのテーブルです
type Bucket struct {
	ID       string `gorm:"primaryKey;"`
	ServerID string `gorm:"index;"`
	Data     []byte `gorm:"type:jsonb"`
	// ここはドメインには定義されておらず、UNIQUE制約を付けるために使用しています
	// serverID_namespace_keyの形式でユニークになります
	// 例: 1234567890-namespace1-key1
	ServerIDNamespaceKey string `gorm:"uniqueIndex;"`
}

// Upsertします
//
// idはCreateの時のみ使用します。
// updateの時は、idは無視されます。
func UpsertBucket(tx *gorm.DB, id, serverID, namespace, key string, data []byte) error {
	var bucket Bucket

	joined := strings.Join([]string{serverID, namespace, key}, "-")

	if err := tx.Where(
		"server_id_namespace_key = ?", joined,
	).First(&bucket).Error; err != nil {
		// レコードが見つからない場合は新しいレコードを作成
		if defaultErrors.Is(err, gorm.ErrRecordNotFound) {
			bucket = Bucket{
				ID:                   id,
				ServerID:             serverID,
				Data:                 data,
				ServerIDNamespaceKey: joined,
			}
			if err = tx.Create(&bucket).Error; err != nil {
				return errors.NewError("レコードを作成できません", err)
			}
		} else {
			return errors.NewError("IDからレコードを取得できません", err)
		}
	} else {
		// レコードが見つかった場合はQuantityを更新
		bucket.Data = data
		if err = tx.Save(&bucket).Error; err != nil {
			return errors.NewError("レコードを更新できません", err)
		}
	}

	return nil
}

// 削除します
func RemoveBucket(tx *gorm.DB, serverID, namespace, key string) error {
	var bucket Bucket

	joined := strings.Join([]string{serverID, namespace, key}, "-")

	if err := tx.Where(
		"server_id_namespace_key = ?", joined,
	).First(&bucket).Error; err != nil {
		// レコードが見つからない場合
		if defaultErrors.Is(err, gorm.ErrRecordNotFound) {
			return errors.NewError("削除するレコードが見つかりません", err)
		}
		// その他のエラーの場合
		return errors.NewError("レコードの検索に失敗しました", err)
	}

	// レコードを削除
	if err := tx.Delete(&bucket).Error; err != nil {
		return errors.NewError("レコードの削除に失敗しました", err)
	}

	return nil
}

// 複数条件で取得します
func FindBucketByCondition(serverID, namespace, key string) (Bucket, error) {
	var bucket Bucket

	joined := strings.Join([]string{serverID, namespace, key}, "-")

	// IDをキーとして検索
	if err := DB.Where(
		"server_id_namespace_key = ?", joined,
	).First(&bucket).Error; err != nil {
		// レコードが見つからない場合
		if defaultErrors.Is(err, gorm.ErrRecordNotFound) {
			return Bucket{}, nil
		}
		// その他のエラーの場合
		return Bucket{}, errors.NewError("レコードの取得に失敗しました", err)
	}

	return bucket, nil
}

// サーバーIDに一致する全ての情報を取得します
func FindAllBucketsByServerID(serverID string) ([]Bucket, error) {
	var buckets []Bucket

	if err := DB.Where(
		"server_id = ?", serverID,
	).First(&buckets).Error; err != nil {
		// レコードが見つからない場合
		if defaultErrors.Is(err, gorm.ErrRecordNotFound) {
			return buckets, nil
		}
		// その他のエラーの場合
		return buckets, errors.NewError("レコードの取得に失敗しました", err)
	}

	return buckets, nil
}
