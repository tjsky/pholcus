package mgo

import (
	"fmt"

	"gopkg.in/mgo.v2/bson"

	"github.com/andeya/pholcus/common/pool"
)

// Update updates the first document matching the selector.
type Update struct {
	Database   string                 // database name
	Collection string                 // collection name
	Selector   map[string]interface{} // document selector
	Change     map[string]interface{} // update document
}

func (self *Update) Exec(_ interface{}) error {
	return Call(func(src pool.Src) error {
		c := src.(*MgoSrc).DB(self.Database).C(self.Collection)

		if id, ok := self.Selector["_id"]; ok {
			if idStr, ok2 := id.(string); !ok2 {
				return fmt.Errorf("%v", "parameter _id must be of string type")
			} else {
				self.Selector["_id"] = bson.ObjectIdHex(idStr)
			}
		}

		return c.Update(self.Selector, self.Change)
	})
}
