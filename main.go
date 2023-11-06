package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/totsumaru/dd-bot-be/api"
	"github.com/totsumaru/dd-bot-be/bot/handler"
	"github.com/totsumaru/dd-bot-be/internal"
	"github.com/totsumaru/dd-bot-be/internal/db"
)

func init() {
	godotenv.Load(".env")

	location := "Asia/Tokyo"
	loc, err := time.LoadLocation(location)
	if err != nil {
		loc = time.FixedZone(location, 9*60*60)
	}
	time.Local = loc
}

func main() {
	var Token = "Bot " + os.Getenv("APP_BOT_TOKEN")

	session, err := discordgo.New(Token)
	session.Token = Token
	if err != nil {
		log.Fatalln(err)
	}

	// DBのセットアップ
	db.ConnectDB()

	// ハンドラを追加
	handler.AddHandler(session)

	if err = session.Open(); err != nil {
		log.Fatalln(err)
	}
	defer func() {
		if err = session.Close(); err != nil {
			log.Fatalln(err)
		}
		return
	}()

	// Ginの設定
	{
		engine := gin.Default()

		// CORSの設定
		// ここからCorsの設定
		engine.Use(cors.New(cors.Config{
			// アクセスを許可したいアクセス元
			AllowOrigins: []string{
				"*",
			},
			// アクセスを許可したいHTTPメソッド
			AllowMethods: []string{
				"GET",
				"POST",
				"OPTIONS",
			},
			// 許可したいHTTPリクエストヘッダ
			AllowHeaders: []string{
				"Origin",
				"Content-Length",
				"Content-Type",
				"Authorization",
				"Accept",
				"X-Requested-With",
			},
			ExposeHeaders: []string{"Content-Length"},
			// cookieなどの情報を必要とするかどうか
			AllowCredentials: false,
			// preflightリクエストの結果をキャッシュする時間
			//MaxAge: 24 * time.Hour,
		}))

		// ルートを設定する
		api.RegisterRouter(engine, db.DB, session)

		if err = engine.Run(":8080"); err != nil {
			log.Fatal("起動に失敗しました", err)
		}
	}

	// Deployedメッセージを送信
	if _, err = session.ChannelMessageSend(internal.ChannelID().LOG, "deployed!"); err != nil {
		log.Fatalln(err)
	}

	stopBot := make(chan os.Signal, 1)
	signal.Notify(stopBot, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-stopBot
}
