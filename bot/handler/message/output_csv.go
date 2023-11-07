package message

import (
	"bytes"
	"encoding/csv"
	"log"
	"sort"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/totsumaru/dd-bot-be/context/record/app"
	"github.com/totsumaru/dd-bot-be/context/record/domain"
	serverApp "github.com/totsumaru/dd-bot-be/context/server/app"
	"github.com/totsumaru/dd-bot-be/internal/db"
	"github.com/totsumaru/dd-bot-be/internal/errors"
)

// CSVを出力します
func OutputCSV(s *discordgo.Session, m *discordgo.MessageCreate) {
	// チャンネルがDB専用のチャンネルかどうかを確認します
	{
		// サーバーを取得します
		server, err := serverApp.GetServer(db.DB, m.GuildID)
		if err != nil {
			errors.SendErrMsg(s, errors.NewError("サーバーを取得できません", err))
		}

		if !server.DBChannelID().Equal(m.ChannelID) {
			return
		}
	}

	// 全てを取得します
	records, err := app.GetAllRecords(db.DB, m.GuildID)
	if err != nil {
		errors.SendErrMsg(s, errors.NewError("レコードを取得できません", err))
		return
	}

	// CSVデータをバッファに書き込む
	var csvBuffer bytes.Buffer
	writer := csv.NewWriter(&csvBuffer)

	// すべてのレコードからユニークなvalueキーを取得してソートする
	valueKeys := getAllValueKeys(records)
	sort.Strings(valueKeys)

	// ヘッダーの作成
	headers := []string{"server_id", "namespace", "key", "updated_at"}
	headers = append(headers, valueKeys...)

	// ヘッダーの書き込み
	if err = writer.Write(headers); err != nil {
		panic(err)
	}

	// 各RecordをCSVに書き込み
	for _, record := range records {
		if err = writeRecord(writer, record, valueKeys); err != nil {
			errors.SendErrMsg(s, errors.NewError("CSVへの書き込みに失敗しました", err))
			return
		}
		writer.Flush() // バッファをフラッシュしてファイルに書き込む
		if err = writer.Error(); err != nil {
			errors.SendErrMsg(s, errors.NewError("CSV書き込みエラー", err))
			return
		}
	}

	file := &discordgo.File{
		Name:   "records.csv",
		Reader: bytes.NewReader(csvBuffer.Bytes()),
	}
	message := &discordgo.MessageSend{
		Files:     []*discordgo.File{file},
		Reference: m.Reference(),
	}

	// Discordにメッセージ送信
	_, err = s.ChannelMessageSendComplex(m.ChannelID, message)
	if err != nil {
		log.Fatalf("error sending message: %v", err)
	}
}

// 全てのRecordからユニークなvalueキーのスライスを返します
func getAllValueKeys(records []domain.Record) []string {
	keyMap := make(map[string]bool)
	for _, record := range records {
		for key := range record.Value() {
			keyMap[key] = true
		}
	}

	keys := make([]string, 0, len(keyMap))
	for key := range keyMap {
		keys = append(keys, key)
	}

	return keys
}

// writeRecord は1つのRecordをCSVに書き込みます。
func writeRecord(
	writer *csv.Writer,
	record domain.Record,
	valueKeys []string,
) error {
	// Recordの基本情報を書き込み(valueは除く)
	recordSlice := []string{
		record.ServerID().String(),
		record.Namespace().String(),
		record.Key().String(),
		record.UpdatedAt().Format(time.RFC3339),
	}

	// valueキーに対応する値を順に追加する
	for _, key := range valueKeys {
		val, ok := record.Value()[key]
		if !ok {
			val = "" // キーが存在しない場合は空文字を入れる
		}
		recordSlice = append(recordSlice, val)
	}

	// 書き込み
	return writer.Write(recordSlice)
}
