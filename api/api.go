package api

import (
	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
	"github.com/totsumaru/dd-bot-be/api/record"
	"gorm.io/gorm"
)

// ルートを設定します
func RegisterRouter(e *gin.Engine, db *gorm.DB, s *discordgo.Session) {
	Route(e)
	record.GetRecord(e, db, s)
}

// ルートです
//
// Note: この関数は削除しても問題ありません
func Route(e *gin.Engine) {
	e.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "hello world",
		})
	})
}
