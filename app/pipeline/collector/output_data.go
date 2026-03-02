package collector

import (
	"github.com/andeya/pholcus/logs"
)

var (
	// DataOutput maps output type names to their implementation functions.
	DataOutput = make(map[string]func(self *Collector) error)

	// DataOutputLib lists the names of supported text data output backends.
	DataOutputLib []string
)

// outputData writes collected text data to the configured output backend.
func (self *Collector) outputData() {
	defer func() {
		self.resetDataDocker()
	}()

	dataLen := uint64(len(self.dataDocker))
	if dataLen == 0 {
		return
	}

	defer func() {
		if p := recover(); p != nil {
			logs.Log.Informational(" * ")
			logs.Log.App(" *     Panic  [Data output: %v | KEYIN: %v | Batch: %v]  %v records! [ERROR]  %v\n",
				self.Spider.GetName(), self.Spider.GetKeyin(), self.dataBatch, dataLen, p)
		}
	}()

	self.addDataSum(dataLen)

	err := DataOutput[self.outType](self)

	logs.Log.Informational(" * ")
	if err != nil {
		logs.Log.App(" *     Fail  [Data output: %v | KEYIN: %v | Batch: %v]  %v records! [ERROR]  %v\n",
			self.Spider.GetName(), self.Spider.GetKeyin(), self.dataBatch, dataLen, err)
	} else {
		logs.Log.App(" *     [Data output: %v | KEYIN: %v | Batch: %v]  %v records!\n",
			self.Spider.GetName(), self.Spider.GetKeyin(), self.dataBatch, dataLen)
		self.Spider.TryFlushSuccess()
	}
}

// Register adds an output backend for the given type name.
func Register(outType string, outFunc func(self *Collector) (err error)) {
	DataOutput[outType] = outFunc
}
