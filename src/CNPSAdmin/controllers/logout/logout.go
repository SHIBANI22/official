package logout

import (
	"CNPSAdmin/model/utils"
	"CNPSAdmin/session"
	"log"
	"runtime/debug"

	"github.com/astaxie/beego"
)

type Logout struct {
	beego.Controller
}

func (c *Logout) Get() {
	pip := c.Ctx.Input.IP()
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Client IP - ", pip)
	var err error
	defer func() {
		if l_exception := recover(); l_exception != nil {
			stack := debug.Stack()
			log.Println(beego.AppConfig.String("loglevel"), "Exception", string(stack))
			session.GlobalSessions.SessionDestroy(c.Ctx.ResponseWriter, c.Ctx.Request)
			c.Data["DisplayMessage"] = "Something went wrong.Please Contact CustomerCare."
			c.TplName = "error/error.html"
		}
		return

	}()
	utils.SetHTTPHeader(c.Ctx)

	sess, _ := session.GlobalSessions.SessionStart(c.Ctx.ResponseWriter, c.Ctx.Request)

	log.Println(beego.AppConfig.String("loglevel"), "Info", "UserName Nil Found")
	uname := sess.Get("uname")
	if uname != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Info", "UserName Nil Found")
		session.SeTLogoutSession(uname.(string))
	}
	sess.SessionRelease(c.Ctx.ResponseWriter)
	session.GlobalSessions.SessionDestroy(c.Ctx.ResponseWriter, c.Ctx.Request)
	utils.EventLogs(c.Ctx, sess, c.Ctx.Input.Method(), c.Input(), c.Data, err)

	c.Redirect("/", 302)
	return
}
