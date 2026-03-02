// Package queue provides a bounded channel-based queue.
package queue

// Queue is a bounded channel-based queue.
type Queue struct {
	PoolSize int
	PoolChan chan interface{}
}

// NewQueue creates a new Queue with the given capacity.
func NewQueue(size int) *Queue {
	return &Queue{
		PoolSize: size,
		PoolChan: make(chan interface{}, size),
	}
}

// Init reinitializes the Queue with a new capacity.
func (this *Queue) Init(size int) *Queue {
	this.PoolSize = size
	this.PoolChan = make(chan interface{}, size)
	return this
}

// Push adds an item to the queue. Returns false if the queue is full.
func (this *Queue) Push(i interface{}) bool {
	if len(this.PoolChan) == this.PoolSize {
		return false
	}
	this.PoolChan <- i
	return true
}

// PushSlice adds all items from the slice to the queue.
func (this *Queue) PushSlice(s []interface{}) {
	for _, i := range s {
		this.Push(i)
	}
}

// Pull removes and returns an item from the queue.
func (this *Queue) Pull() interface{} {
	return <-this.PoolChan
}

// Exchange resizes the queue for reuse. Returns the number of items that can be added.
func (this *Queue) Exchange(num int) (add int) {
	last := len(this.PoolChan)

	if last >= num {
		add = int(0)
		return
	}

	if this.PoolSize < num {
		pool := []interface{}{}
		for i := 0; i < last; i++ {
			pool = append(pool, <-this.PoolChan)
		}
		this.Init(num).PushSlice(pool)
	}

	add = num - last
	return
}
