package initialize

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/leminhthai/train-ticket/booking-service/global"
	"go.uber.org/zap"
)

func checkErrorPanic(err error, errString string) {
	if err != nil {
		global.Logger.Error(errString, zap.Error(err))

		panic(err)
	}
}

func InitMySql() {
	m := global.Config.MySql

	dsn := "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	s := fmt.Sprintf(dsn, m.Username, m.Password, m.Host, m.Port, m.DbName)

	db, err := sql.Open("mysql", s)

	checkErrorPanic(err, "Init MySql Failed!!!")
	global.Logger.Info("Init MySql Successfully!!!")

	global.Mdb = db

	SetPool()
}

func SetPool() {
	if global.Mdb == nil {
		global.Logger.Error("SetPool: DB not initialized!!!")
		return
	}

	m := global.Config.MySql

	global.Mdb.SetMaxIdleConns(m.MaxIdleConns)
	global.Mdb.SetMaxOpenConns(m.MaxOpenConns)
	global.Mdb.SetConnMaxLifetime(time.Duration(m.ConnMaxLifetime))
}
