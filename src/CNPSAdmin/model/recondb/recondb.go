package recondb

import (
	"log"

	"time"

	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
	"ominaya.com/database/sql"
)

var Db sql.Database

func Init() (err error) {
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Trying to connect Recon DB")
	Db.Ip = beego.AppConfig.String("ReconDBIP")
	Db.Port = beego.AppConfig.String("ReconDBPort")
	Db.Type = beego.AppConfig.String("ReconDBType")
	Db.Schema = beego.AppConfig.String("ReconDBName")
	Db.Username = beego.AppConfig.String("ReconDBUsername")
	Db.Password = beego.AppConfig.String("ReconDBPassword")
	Db.LogLevel = beego.AppConfig.String("loglevel")

	err = Db.Connect()
	Db.ConnPtr.SetConnMaxLifetime(time.Second * 300)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", "Recon DB Connect fail")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Info", "Recon DB Connected successfully")
	//////////////////////////////////////////////////////////////////////////////////////////////////

	return
}
