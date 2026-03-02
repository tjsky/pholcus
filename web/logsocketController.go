package web

import (
	"regexp"
	"runtime/debug"
	"sync"
	"sync/atomic"

	ws "github.com/andeya/pholcus/common/websocket"
	"github.com/andeya/pholcus/logs"
)

// wsLogHandle handles WebSocket connections for streaming logs to the client.
func wsLogHandle(conn *ws.Conn) {
	defer func() {
		if p := recover(); p != nil {
			logs.Log.Error("panic recovered: %v\n%s", p, debug.Stack())
		}
	}()
	sess, _ := globalSessions.SessionStart(nil, conn.Request())
	sessID := sess.SessionID()
	connPool := Lsc.connPool.Load().(map[string]*ws.Conn)
	if connPool[sessID] == nil {
		Lsc.Add(sessID, conn)
	}
	defer func() {
		Lsc.Remove(sessID)
	}()
	for {
		if err := ws.JSON.Receive(conn, nil); err != nil {
			return
		}
	}
}

// LogSocketController manages WebSocket connections for log streaming.
type LogSocketController struct {
	connPool atomic.Value
	lock     sync.Mutex
}

var (
	// Lsc is the global LogSocketController for log streaming.
	Lsc = func() *LogSocketController {
		l := new(LogSocketController)
		l.connPool.Store(make(map[string]*ws.Conn))
		return l
	}()
	colorRegexp = regexp.MustCompile("\033\\[[0-9;]{1,4}m")
)

func (self *LogSocketController) Write(p []byte) (int, error) {
	defer func() {
		if r := recover(); r != nil {
			logs.Log.Error("panic recovered: %v\n%s", r, debug.Stack())
		}
	}()
	p = colorRegexp.ReplaceAll(p, []byte{})
	connPool := self.connPool.Load().(map[string]*ws.Conn)
	for sessID, conn := range connPool {
		_, err := ws.Message.Send(conn, (string(p) + "\r\n"))
		if err != nil {
			self.Remove(sessID)
		}
	}
	return len(p), nil
}

func (self *LogSocketController) Add(sessID string, conn *ws.Conn) {
	self.lock.Lock()
	defer self.lock.Unlock()

	connPool := self.connPool.Load().(map[string]*ws.Conn)
	newConnPool := make(map[string]*ws.Conn, len(connPool)+1)
	for k, v := range connPool {
		newConnPool[k] = v
	}
	newConnPool[sessID] = conn
	self.connPool.Store(newConnPool)
}

func (self *LogSocketController) Remove(sessID string) {
	self.lock.Lock()
	defer self.lock.Unlock()
	defer func() {
		if p := recover(); p != nil {
			logs.Log.Error("panic recovered: %v\n%s", p, debug.Stack())
		}
	}()
	connPool := self.connPool.Load().(map[string]*ws.Conn)
	conn := connPool[sessID]
	if conn == nil {
		return
	}
	conn.Close()
	newConnPool := make(map[string]*ws.Conn, len(connPool)+1)
	for k, v := range connPool {
		if k != sessID {
			newConnPool[k] = v
		}
	}
	self.connPool.Store(newConnPool)
}
