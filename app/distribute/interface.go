package distribute

// Distributer is the distributed interface.
type Distributer interface {
	// Send sends a task from the master to the jar.
	Send(clientNum int) Task
	// Receive receives a task into the jar on a slave node.
	Receive(task *Task)
	// CountNodes returns the number of connected nodes.
	CountNodes() int
}
