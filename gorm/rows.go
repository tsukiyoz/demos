package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Interactive struct {
	Id    int64  `gorm:"primaryKey,autoIncrement"`
	BizId int64  `gorm:"uniqueIndex:biz_type_id"`
	Biz   string `gorm:"type:varchar(128);uniqueIndex:biz_type_id"`

	ReadCnt     int64
	FavoriteCnt int64
	LikeCnt     int64

	Ctime int64
	Utime int64
}

func main() {
	dsn := "root:for.nothing@tcp(localhost:3306)/mercury"
	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		err = sqlDB.PingContext(ctx)
		cancel()
		if err == nil {
			break
		}
		log.Println("waiting for connect MySQL", err)
	}
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic(err)
	}
	rows, err := db.Model(&Interactive{}).Limit(0).Rows()
	if err != nil {
		log.Printf("rows errors %v\n", err)
	}
	columns, err := rows.Columns()
	if err != nil {
		log.Printf("rows errors %v\n", err)
	}
	log.Printf("rows %v\n", columns)
}
