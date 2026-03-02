package mgo

import (
	"fmt"

	"gopkg.in/mgo.v2/bson"

	"github.com/andeya/pholcus/common/pool"
)

// Count returns the number of documents matching the query.
type Count struct {
	Database   string                 // database name
	Collection string                 // collection name
	Query      map[string]interface{} // query filter
}

func (self *Count) Exec(resultPtr interface{}) (err error) {
	defer func() {
		if re := recover(); re != nil {
			err = fmt.Errorf("%v", re)
		}
	}()
	resultPtr2 := resultPtr.(*int)
	*resultPtr2 = 0

	err = Call(func(src pool.Src) error {
		c := src.(*MgoSrc).DB(self.Database).C(self.Collection)

		if id, ok := self.Query["_id"]; ok {
			if idStr, ok2 := id.(string); !ok2 {
				return fmt.Errorf("%v", "parameter _id must be of string type")
			} else {
				self.Query["_id"] = bson.ObjectIdHex(idStr)
			}
		}

		*resultPtr2, err = c.Find(self.Query).Count()
		return err
	})
	return
}
