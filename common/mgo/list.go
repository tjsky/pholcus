package mgo

import (
	"fmt"

	"github.com/andeya/pholcus/common/pool"
)

// List returns a map of database names to their collection names.
type List struct {
	Dbs []string // list of database names to query (empty = all)
}

func (self *List) Exec(resultPtr interface{}) (err error) {
	defer func() {
		if re := recover(); re != nil {
			err = fmt.Errorf("%v", re)
		}
	}()
	resultPtr2 := resultPtr.(*map[string][]string)
	*resultPtr2 = map[string][]string{}

	err = Call(func(src pool.Src) error {
		var (
			s   = src.(*MgoSrc)
			dbs []string
		)

		if dbs, err = s.DatabaseNames(); err != nil {
			return err
		}

		if len(self.Dbs) == 0 {
			for _, dbname := range dbs {
				(*resultPtr2)[dbname], err = s.DB(dbname).CollectionNames()
				if err != nil {
					return err
				}
			}
			return err
		}

		for _, dbname := range self.Dbs {
			(*resultPtr2)[dbname], err = s.DB(dbname).CollectionNames()
			if err != nil {
				return err
			}
		}
		return err
	})

	return
}
