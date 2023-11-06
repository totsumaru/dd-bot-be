package record

import (
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	recordApp "github.com/totsumaru/dd-bot-be/context/record/app"
	serverApp "github.com/totsumaru/dd-bot-be/context/server/app"
	"gorm.io/gorm"
)

// レスポンスです
//
// jsonへの変換はMarshalJSONを実装しています。
type Res struct {
	Namespace string
	Key       string
	UpdatedAt time.Time         // 例: 2023-09-01T12:00:00+09:00
	Value     map[string]string // ここは`value`は無視され、展開して返されます
}

// レコードを取得します
func GetRecord(e *gin.Engine, db *gorm.DB) {
	e.GET("/record", func(c *gin.Context) {
		guild := c.Query("guild")
		namespace := c.Query("namespace")
		key := c.Query("key")

		apiKey := c.GetHeader("X-API-KEY")

		if guild == "" || namespace == "" || key == "" {
			c.JSON(400, "guild, namespace, keyを指定してください")
			return
		}

		// サーバーを取得します
		serverRes, err := serverApp.GetServer(db, guild)
		if err != nil {
			c.JSON(500, "サーバーを取得できません")
			return
		}

		if !serverRes.APIKey().IsMatch(apiKey) {
			c.JSON(401, "APIキーが不正です")
			return
		}

		// レコードを取得します
		recordRes, err := recordApp.GetRecord(db, guild, namespace, key)
		if err != nil {
			c.JSON(500, "レコードを取得できません")
			return
		}

		value := map[string]string{}
		for k, v := range recordRes.Value() {
			value[k] = v
		}

		res := Res{
			Namespace: recordRes.Namespace().String(),
			Key:       recordRes.Key().String(),
			Value:     value,
			UpdatedAt: recordRes.UpdatedAt(),
		}

		c.JSON(200, res)
	})
}

// MarshalJSON は Res のカスタム JSON マーシャリングを実装します。
func (r Res) MarshalJSON() ([]byte, error) {
	// Value マップを含む一時的なマップを作成
	tmp := map[string]interface{}{
		"namespace":  r.Namespace,
		"key":        r.Key,
		"updated_at": r.UpdatedAt,
	}
	// Value マップのキーと値を一時的なマップに追加
	for k, v := range r.Value {
		tmp[k] = v
	}

	return json.Marshal(tmp)
}
