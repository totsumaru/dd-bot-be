package internal

import "os"

type Channel struct {
	LOG string
}

func ChannelID() Channel {
	if os.Getenv("ENV") == "dev" {
		// テスト環境
		return Channel{
			LOG: "1170723690644766852",
		}
	} else {
		// 本番環境
		return Channel{
			LOG: "1170723690644766852",
		}
	}
}
