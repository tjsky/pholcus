package proxy

import (
	"sync"
	"time"
)

// ProxyForHost manages proxy IPs for a host, sorted by response time.
type ProxyForHost struct {
	curIndex  int // Index of current proxy IP
	proxys    []string
	timedelay []time.Duration
	isEcho    bool // Whether to print proxy switch info
	sync.Mutex
}

// Len implements sort.Interface.
func (self *ProxyForHost) Len() int {
	return len(self.proxys)
}

func (self *ProxyForHost) Less(i, j int) bool {
	return self.timedelay[i] < self.timedelay[j]
}

func (self *ProxyForHost) Swap(i, j int) {
	self.proxys[i], self.proxys[j] = self.proxys[j], self.proxys[i]
	self.timedelay[i], self.timedelay[j] = self.timedelay[j], self.timedelay[i]
}
