package bot

import (
	"sync"
)

// サーバーの情報を保持するためのストアです
var ServerMemoryStore = NewServerStore()

// サーバーの情報を保持する構造体です。
type ServerData struct {
	ServerID    string
	DBChannelID string
}

// ServerStore は ServerData を管理するためのスレッドセーフなストアです。
type ServerStore struct {
	sync.Map
}

// キーと値のペアを保存します。
//
// もしキーが存在しない場合は新しい値を保存します。
func (s *ServerStore) Insert(key string, value ServerData) {
	s.Store(key, value)
}

// 取得します
func (s *ServerStore) Get(key string) (ServerData, bool) {
	result, ok := s.Load(key)
	if !ok {
		return ServerData{}, false
	}
	return result.(ServerData), true
}

// 削除します
func (s *ServerStore) Remove(key string) {
	s.Delete(key)
}

// ServerStore を初期化します。
func NewServerStore() *ServerStore {
	return &ServerStore{}
}
