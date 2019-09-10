package model

import (
	"fmt"
	"github.com/hunterhug/fafacms/core/util/rdb"
)

var FafaRdb *rdb.MyDb

func CreateTable(tables []interface{}) {
	for _, table := range tables {
		ok, err := FafaRdb.IsTableExist(table)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		if !ok {
			//  change the Charset
			sess := FafaRdb.Client.NewSession()
			sess.Charset("utf8mb4")
			err = sess.CreateTable(table)
			if err != nil {
				sess.Close()
				fmt.Println(err.Error())
				continue
			}

			sess.Close()
		}

		err = FafaRdb.Client.CreateIndexes(table)
		if err != nil {
			fmt.Println(err.Error())
			//continue
		}
		err = FafaRdb.Client.CreateUniques(table)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
	}

	u := new(User)
	u.Name = "admin"
	u.Email = "admin@admin"
	u.NickName = "admin"
	u.Password = "admin"
	u.Status = 1
	u.InsertOne()
}
