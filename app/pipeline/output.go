package pipeline

import (
	"sort"

	"github.com/andeya/pholcus/app/pipeline/collector"
	"github.com/andeya/pholcus/common/kafka"
	"github.com/andeya/pholcus/common/mgo"
	"github.com/andeya/pholcus/common/mysql"
	"github.com/andeya/pholcus/runtime/cache"
)

// init populates collector.DataOutputLib with available output backend names.
func init() {
	for out, _ := range collector.DataOutput {
		collector.DataOutputLib = append(collector.DataOutputLib, out)
	}
	sort.Strings(collector.DataOutputLib)
}

// RefreshOutput refreshes the state of the configured output backend.
func RefreshOutput() {
	switch cache.Task.OutType {
	case "mgo":
		mgo.Refresh()
	case "mysql":
		mysql.Refresh()
	case "kafka":
		kafka.Refresh()
	}
}
