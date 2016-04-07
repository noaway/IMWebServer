package maps

import (
	"github.com/gorilla/websocket"
	"sync"
)

var (
	conns Connmap
)

type Connmap struct {
	_connmap map[string]*websocket.Conn
}

func init() {
	conns = Connmap{_connmap: make(map[string]*websocket.Conn, 100000)}
}

func Conns(id string) *websocket.Conn {
	return conns._connmap[id]
}

func Register(id string, c *websocket.Conn) {
	var m sync.RWMutex
	m.RLock()
	defer m.RUnlock()
	conns._connmap[id] = c
}
