package scheduler

import (
	"runtime/debug"
	"sync"

	"github.com/andeya/pholcus/app/aid/proxy"
	"github.com/andeya/pholcus/logs"
	"github.com/andeya/pholcus/runtime/cache"
	"github.com/andeya/pholcus/runtime/status"
)

// scheduler coordinates crawl tasks and resource allocation.
type scheduler struct {
	status       int          // running status
	count        chan bool    // total concurrency count
	useProxy     bool         // whether proxy IP is used
	proxy        *proxy.Proxy // global proxy IP
	matrices     []*Matrix    // request matrices per Spider instance
	sync.RWMutex              // global read-write lock
}

// sdl is the global scheduler instance.
var sdl = &scheduler{
	status: status.RUN,
	count:  make(chan bool, cache.Task.ThreadNum),
	proxy:  proxy.New(),
}

// Init initialize scheduler.
func Init() {
	sdl.matrices = []*Matrix{}
	sdl.count = make(chan bool, cache.Task.ThreadNum)

	if cache.Task.ProxyMinute > 0 {
		if sdl.proxy.Count() > 0 {
			sdl.useProxy = true
			sdl.proxy.UpdateTicker(cache.Task.ProxyMinute)
			logs.Log.Informational(" *     Using proxy IP, rotation interval: %v minutes\n", cache.Task.ProxyMinute)
		} else {
			sdl.useProxy = false
			logs.Log.Informational(" *     Proxy IP list is empty, cannot use proxy\n")
		}
	} else {
		sdl.useProxy = false
		logs.Log.Informational(" *     Not using proxy IP\n")
	}

	sdl.status = status.RUN
}

// ReloadProxyLib reload proxy ip list from config file.
func ReloadProxyLib() {
	sdl.proxy.Update()
}

// AddMatrix registers a resource queue for the given spider and returns its Matrix.
func AddMatrix(spiderName, spiderSubName string, maxPage int64) *Matrix {
	matrix := newMatrix(spiderName, spiderSubName, maxPage)
	sdl.RLock()
	defer sdl.RUnlock()
	sdl.matrices = append(sdl.matrices, matrix)
	return matrix
}

// PauseRecover toggles pause/resume for all crawl tasks.
func PauseRecover() {
	sdl.Lock()
	defer sdl.Unlock()
	switch sdl.status {
	case status.PAUSE:
		sdl.status = status.RUN
	case status.RUN:
		sdl.status = status.PAUSE
	}
}

// Stop terminates all crawl tasks.
func Stop() {
	sdl.Lock()
	defer sdl.Unlock()
	sdl.status = status.STOP
	defer func() {
		if p := recover(); p != nil {
			logs.Log.Error("panic recovered: %v\n%s", p, debug.Stack())
		}
	}()
	// for _, matrix := range sdl.matrices {
	// 	matrix.windup()
	// }
	close(sdl.count)
	sdl.matrices = []*Matrix{}
}

// avgRes returns the average resources allocated per spider instance.
func (self *scheduler) avgRes() int32 {
	avg := int32(cap(sdl.count) / len(sdl.matrices))
	if avg == 0 {
		avg = 1
	}
	return avg
}

func (self *scheduler) checkStatus(s int) bool {
	self.RLock()
	b := self.status == s
	self.RUnlock()
	return b
}
