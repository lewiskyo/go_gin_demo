package model

import (
	"database/sql"
	"fmt"
	"go_gin_demo/bo"
	"go_gin_demo/db"

	_ "github.com/go-sql-driver/mysql"
)

func InsertMsg(msg bo.Msg) {
	dbconn := db.GetDB()

	_, err := dbconn.Exec("insert into msg(uid,type,region,device,ip,network,version,timestamp) values(?,?,?,?,?,?,?)", msg.Uid, msg.Type, msg.Region, msg.Device, msg.Ip, msg.Network, msg.Version, msg.Timestamp)
	if err != nil {
		fmt.Printf("insert msg fail, err: %s", err.Error())
		return
	}
}

func GetLastMsg(uid int64) (bo.Msg, error) {
	dbconn := db.GetDB()

	msg := bo.Msg{}
	row := dbconn.QueryRow(fmt.Sprintf("select uid,type,region,device,ip,network,version,timestamp from msg where uid = %d order by timestamp desc limit 1", uid))
	err := row.Scan(&msg.Uid, &msg.Type, &msg.Region, &msg.Device, &msg.Ip, &msg.Network, &msg.Version, &msg.Timestamp)
	if err != nil {
		if err == sql.ErrNoRows {
			return msg, err
		}
	}

	return msg, nil
}
