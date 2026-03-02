package distribute

// TaskJar is the task storage.
type TaskJar struct {
	Tasks chan *Task
}

func NewTaskJar() *TaskJar {
	return &TaskJar{
		Tasks: make(chan *Task, 1024),
	}
}

// Push adds a task to the jar (server side).
func (self *TaskJar) Push(task *Task) {
	id := len(self.Tasks)
	task.Id = id
	self.Tasks <- task
}

// Pull gets a task from the local jar (client side).
func (self *TaskJar) Pull() *Task {
	return <-self.Tasks
}

// Len returns number of tasks in the jar.
func (self *TaskJar) Len() int {
	return len(self.Tasks)
}

// Send sends a task from the jar (master side).
func (self *TaskJar) Send(clientNum int) Task {
	return *<-self.Tasks
}

// Receive receives a task into the jar (slave side).
func (self *TaskJar) Receive(task *Task) {
	self.Tasks <- task
}
