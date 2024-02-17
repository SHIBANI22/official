package error

import (
	"CNPSAdmin/model/utils"
	"CNPSAdmin/session"
	"log"

	"github.com/astaxie/beego"
)

type Error struct {
	beego.Controller
}

func (c *Error) Error404() {
	utils.SetHTTPHeader(c.Ctx)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "404 Path", c.Ctx.Request.URL.Path)

	c.TplName = "error/error.html"
}

func (c *Error) Error501() {
	utils.SetHTTPHeader(c.Ctx)
	pip := c.Ctx.Input.IP()
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Client IP - ", pip)
	sess, _ := session.GlobalSessions.SessionStart(c.Ctx.ResponseWriter, c.Ctx.Request)

	log.Println(beego.AppConfig.String("loglevel"), "Info", "UserName Nil Found")
	uname := sess.Get("uname")
	if uname != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Info", "UserName Nil Found")
		session.SeTLogoutSession(uname.(string))
	}
	sess.SessionRelease(c.Ctx.ResponseWriter)
	session.GlobalSessions.SessionDestroy(c.Ctx.ResponseWriter, c.Ctx.Request)
	c.Data["language"] = sess.Get("language").(string)
	c.Data["DisplayMessage"] = "501, server error"
	c.TplName = "error/error.html"
}
