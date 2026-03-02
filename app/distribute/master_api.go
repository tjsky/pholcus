package distribute

import (
	"encoding/json"

	"github.com/andeya/pholcus/app/distribute/teleport"
	"github.com/andeya/pholcus/logs"
)

// MasterApi creates the master node API.
func MasterApi(n Distributer) teleport.API {
	return teleport.API{
		"task": &masterTaskHandle{n},
		"log":  &masterLogHandle{},
	}
}

// masterTaskHandle assigns tasks to clients.
type masterTaskHandle struct {
	Distributer
}

func (self *masterTaskHandle) Process(receive *teleport.NetData) *teleport.NetData {
	b, _ := json.Marshal(self.Send(self.CountNodes()))
	return teleport.ReturnData(string(b))
}

// masterLogHandle receives and prints log messages from slave nodes.
type masterLogHandle struct{}

func (*masterLogHandle) Process(receive *teleport.NetData) *teleport.NetData {
	logs.Log.Informational(" * ")
	logs.Log.Informational(" *     [ %s ]    %s", receive.From, receive.Body)
	logs.Log.Informational(" * ")
	return nil
}
