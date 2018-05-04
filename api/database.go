package api

import (
	"github.com/coreos/bbolt"
	"fmt"
	"zuijinbuzai/fundtop/api/types"
	"time"
	"os"
	"encoding/json"
)

const (
	DB_PATH = "fundtop.db"
	TABLE_FUND = "fund"
)

var (
	db			*bolt.DB
	tableName 	string
)

func DBOpen() {
	var err error
	db, err = bolt.Open(DB_PATH, 0600, nil)
	if err != nil {
		return
	}

	t := time.Now()
	tableName = fmt.Sprintf("%d%02d%02d", t.Year(), t.Month(), t.Day())
	//fmt.Println(tableName)

	exist := true
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(tableName))
		if bucket == nil {
			exist = false
		}
		return nil
	})

	if !exist {
		db.Close()
		os.Remove(DB_PATH)
		db, err = bolt.Open(DB_PATH, 0600, nil)
		if err != nil {
			return
		}
		err = db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte(tableName))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
			return nil
		})
	}
}

func DBPutArray(arr []*types.Fund) {
	tx, _ := db.Begin(true)
	b := tx.Bucket([]byte(tableName))
	for _, v := range arr {
		if len(v.FArray) > 1 {
			data, _ := json.Marshal(v.FArray[1:])
			b.Put([]byte(v.Code), data)
		}
	}
	tx.Commit()
}

func DBGet(code string) *[]*types.FundItem {
	result := &[]*types.FundItem{}
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(tableName))
		data := b.Get([]byte(code))

		err := json.Unmarshal(data, result)
		return err
	})
	if err != nil {
		return nil
	}
	return result
}

func DBClose()  {
	if db != nil {
		db.Close()
		db = nil
	}
}