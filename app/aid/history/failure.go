package history

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/andeya/pholcus/app/downloader/request"
	"github.com/andeya/pholcus/common/mgo"
	"github.com/andeya/pholcus/common/mysql"
	"github.com/andeya/pholcus/common/pool"
	"github.com/andeya/pholcus/config"
)

// Failure tracks failed requests for retry.
type Failure struct {
	tabName     string
	fileName    string
	list        map[string]*request.Request
	inheritable bool
	sync.RWMutex
}

func (self *Failure) PullFailure() map[string]*request.Request {
	list := self.list
	self.list = make(map[string]*request.Request)
	return list
}

// UpsertFailure updates or adds a failure record. Returns true if an insert occurred.
func (self *Failure) UpsertFailure(req *request.Request) bool {
	self.RWMutex.Lock()
	defer self.RWMutex.Unlock()
	if self.list[req.Unique()] != nil {
		return false
	}
	self.list[req.Unique()] = req
	return true
}

// DeleteFailure removes a failure record.
func (self *Failure) DeleteFailure(req *request.Request) {
	self.RWMutex.Lock()
	delete(self.list, req.Unique())
	self.RWMutex.Unlock()
}

// flush clears historical failure records first, then updates.
func (self *Failure) flush(provider string) (fLen int, err error) {
	self.RWMutex.Lock()
	defer self.RWMutex.Unlock()
	fLen = len(self.list)

	switch provider {
	case "mgo":
		if mgo.Error() != nil {
			err = fmt.Errorf(" *     Fail  [add failure record][mgo]: %v [ERROR]  %v\n", fLen, mgo.Error())
			return
		}
		mgo.Call(func(src pool.Src) error {
			c := src.(*mgo.MgoSrc).DB(config.DB_NAME).C(self.tabName)
			c.DropCollection()
			if fLen == 0 {
				return nil
			}

			var docs = []interface{}{}
			for key, req := range self.list {
				docs = append(docs, map[string]interface{}{"_id": key, "failure": req.Serialize()})
			}
			c.Insert(docs...)
			return nil
		})

	case "mysql":
		_, err := mysql.DB()
		if err != nil {
			return fLen, fmt.Errorf(" *     Fail  [add failure record][mysql]: %v [PING]  %v\n", fLen, err)
		}
		table, ok := getWriteMysqlTable(self.tabName)
		if !ok {
			table = mysql.New()
			table.SetTableName(self.tabName).CustomPrimaryKey(`id VARCHAR(255) NOT NULL PRIMARY KEY`).AddColumn(`failure MEDIUMTEXT`)
			setWriteMysqlTable(self.tabName, table)
			err = table.Create()
			if err != nil {
				return fLen, fmt.Errorf(" *     Fail  [add failure record][mysql]: %v [CREATE]  %v\n", fLen, err)
			}
		} else {
			err = table.Truncate()
			if err != nil {
				return fLen, fmt.Errorf(" *     Fail  [add failure record][mysql]: %v [TRUNCATE]  %v\n", fLen, err)
			}
		}

		for key, req := range self.list {
			table.AutoInsert([]string{key, req.Serialize()})
			err = table.FlushInsert()
			if err != nil {
				fLen--
			}
		}

	default:
		os.Remove(self.fileName)
		if fLen == 0 {
			return
		}

		f, _ := os.OpenFile(self.fileName, os.O_CREATE|os.O_WRONLY, 0777)

		docs := make(map[string]string, len(self.list))
		for key, req := range self.list {
			docs[key] = req.Serialize()
		}
		b, _ := json.Marshal(docs)
		b = bytes.Replace(b, []byte(`\u0026`), []byte(`&`), -1)
		f.Write(b)
		f.Close()
	}
	return
}
