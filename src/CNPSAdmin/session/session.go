package session

import (
	"errors"
	"strings"

	"ominaya.com/net/redis"
	"ominaya.com/util/datetime"
	"ominaya.com/util/log"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/session"
	_ "github.com/astaxie/beego/session/redis"
)

var GlobalSessions *session.Manager
var ru redis.Redis
var rc redis.Redis

func Init() (err error) {

	var cf session.ManagerConfig
	cf.CookieName = beego.AppConfig.String("appname")
	cf.EnableSetCookie = true
	cf.Gclifetime, _ = beego.AppConfig.Int64("gclifetime")
	cf.Maxlifetime, _ = beego.AppConfig.Int64("maxLifetime")
	cf.Secure, _ = beego.AppConfig.Bool("EnableHTTPS")
	cf.CookieLifeTime, _ = beego.AppConfig.Int("cookieLifeTime")
	cf.ProviderConfig = beego.AppConfig.String("SessionProviderConfig")

	GlobalSessions, err = session.NewManager(beego.AppConfig.String("SessionProviderName"), &cf)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		return
	}

	go GlobalSessions.GC()

	ru.SavePath = beego.AppConfig.String("RedisConnectionString")
	err = ru.Connect()
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		return
	}
	rc.SavePath = beego.AppConfig.String("SessionProviderConfig")
	err = rc.Connect()
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		return
	}
	return
}

func ValidateSession(sess session.Store) (err error) {
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Session ID - ", sess.SessionID())
	if !rc.Exists(sess.SessionID()) {
		log.Println(beego.AppConfig.String("loglevel"), "Error", "Session Expire")
		err = errors.New("Session Time Out.Please Logout and Login Again.")
		return
	}

	uname := sess.Get("uname").(string)
	log.Println(beego.AppConfig.String("loglevel"), "Info", "UserName - ", uname)
	if ru.Exists(uname) {
		log.Println(beego.AppConfig.String("loglevel"), "Info", "User Session Found")
		val, err := ru.Get(uname)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			return err
		}
		data := strings.Split(val, "*")
		log.Println(beego.AppConfig.String("loglevel"), "Info", "User Session Status - ", data[3])
		if data[0] == sess.SessionID() && data[3] == "LOGGED_IN" {
			log.Println(beego.AppConfig.String("loglevel"), "Info", "User Session Good")
			return nil
		} else {
			log.Println(beego.AppConfig.String("loglevel"), "Info", "User Session Bad")
			err = errors.New("Invalid session found in cache.Please Logout and Login Again.")
			return err
		}
	} else {
		log.Println(beego.AppConfig.String("loglevel"), "Info", "User Session Not Found")
		err = errors.New("User session not found.Please Logout and Login Again.")
	}
	return
}

func SetUserSession(sid, uname string) (err error) {
	dt, _ := datetime.Get("", "", "Africa/Harare")
	log.Println(beego.AppConfig.String("loglevel"), "Info", "UserName - ", uname)
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Session ID - ", sid)
	if ru.Exists(uname) {
		log.Println(beego.AppConfig.String("loglevel"), "Info", "User Session Found")
		val, err := ru.Get(uname)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			return err
		}
		log.Println(beego.AppConfig.String("loglevel"), "Info", "User Session Value - ", val)
		data := strings.Split(val, "*")

		newSess := sid + "*" + data[2] + "*" + dt + "*" + "LOGGED_IN" + "*"
		log.Println(beego.AppConfig.String("loglevel"), "Info", "New User Session Value - ", newSess)
		err = ru.Set(uname, newSess)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			return err
		}
	} else {
		log.Println(beego.AppConfig.String("loglevel"), "Info", "User Session Not Found")
		newSess := sid + "**" + dt + "*" + "LOGGED_IN" + "*"
		log.Println(beego.AppConfig.String("loglevel"), "Info", "New User Session Value - ", newSess)
		err = ru.Set(uname, newSess)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			return err
		}
	}
	return
}

func SeTLogoutSession(uname string) (err error) {
	dt, _ := datetime.Get("", "", "Africa/Harare")
	log.Println(beego.AppConfig.String("loglevel"), "Info", "UserName - ", uname)
	if ru.Exists(uname) {
		val, err := ru.Get(uname)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			return err
		}
		data := strings.Split(val, "*")
		log.Println(beego.AppConfig.String("loglevel"), "Info", "User Session Value - ", val)
		newSess := data[0] + "*" + data[1] + "*" + data[2] + "*" + "LOGGED_OUT" + "*" + dt
		log.Println(beego.AppConfig.String("loglevel"), "Info", "LogedOut User Session Value - ", newSess)
		err = ru.Set(uname, newSess)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			return err
		}
	} else {
		log.Println(beego.AppConfig.String("loglevel"), "Info", "User Session Not Found")
		return
	}
	return
}

func CheckUserSession(uname string) (err error) {
	dt, _ := datetime.Get("", "", "Africa/Harare")
	if ru.Exists(uname) {
		log.Println(beego.AppConfig.String("loglevel"), "Info", "User Session Found")
		val, err := ru.Get(uname)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			return err
		}
		log.Println(beego.AppConfig.String("loglevel"), "Info", "User Session Value - ", val)
		data := strings.Split(val, "*")
		log.Println(beego.AppConfig.String("loglevel"), "Info", "Status - ", data[0])
		log.Println(beego.AppConfig.String("loglevel"), "Info", "Status - ", data[1])
		log.Println(beego.AppConfig.String("loglevel"), "Info", "Status - ", data[2])
		log.Println(beego.AppConfig.String("loglevel"), "Info", "Status - ", data[3])
		if data[3] == "LOGGED_OUT" {
			log.Println(beego.AppConfig.String("loglevel"), "Info", "User Logout Session found")
			return nil
		}
		err = ru.Set(uname, data[0]+"*"+data[1]+"*"+data[2]+"*"+"LOGGED_OUT"+"*"+dt)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			return err
		}
		log.Println(beego.AppConfig.String("loglevel"), "Info", "User logout session set")
		err = errors.New("Multiple user login not allowed.Now you have been logged out from all browsers.Please login again")
		return err
	} else {
		log.Println(beego.AppConfig.String("loglevel"), "Info", "User Session Not Found")
		return
	}
	return
}
