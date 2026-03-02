package distribute

// Task is used for distributed task dispatch.
type Task struct {
	Id             int
	Spiders        []map[string]string // Spider rule name and keyin, format: map[string]string{"name":"baidu","keyin":"henry"}
	ThreadNum      int                 // Global max concurrency
	Pausetime      int64               // Pause duration in ms (random: Pausetime/2 ~ Pausetime*2)
	OutType        string              // Output method
	DockerCap      int                 // Segment dump container capacity
	DockerQueueCap int                 // Segment output pool capacity, >= 2
	SuccessInherit bool                // Inherit historical success records
	FailureInherit bool                // Inherit historical failure records
	Limit          int64               // Collection limit, 0=unlimited; if rule sets LIMIT then custom limit
	ProxyMinute    int64               // Proxy IP rotation interval in minutes
	Keyins         string              // Custom input, later split into Keyin config for multiple tasks
}
