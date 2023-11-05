package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/totsumaru/dd-bot-be/handler"
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

	db.ConnectDB()

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

	// Deployedメッセージを送信
	if _, err = session.ChannelMessageSend(internal.ChannelID().LOG, "deployed!"); err != nil {
		log.Fatalln(err)
	}

	stopBot := make(chan os.Signal, 1)
	signal.Notify(stopBot, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-stopBot
}
