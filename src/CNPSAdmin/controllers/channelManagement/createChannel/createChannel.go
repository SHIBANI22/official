package createChannel

import (
	"CNPSAdmin/model/db"
	"CNPSAdmin/model/utils"
	"CNPSAdmin/session"
	"errors"
	"html/template"
	"runtime/debug"

	"github.com/astaxie/beego"
	"ominaya.com/util/log"
)

type CreateChannel struct {
	beego.Controller
}

func (c *CreateChannel) Get() {
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Create Channel Page Start")
	pip := c.Ctx.Input.IP()
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Client IP - ", pip)
	var err error
	sessErr := false
	defer func() {

		if l_exception := recover(); l_exception != nil {
			stack := debug.Stack()
			log.Println(beego.AppConfig.String("loglevel"), "Exception", string(stack))
			session.GlobalSessions.SessionDestroy(c.Ctx.ResponseWriter, c.Ctx.Request)
			c.Data["DisplayMessage"] = "Something went wrong.Please Contact CustomerCare."
			c.TplName = "error/error.html"
		}
		if err != nil {
			if sessErr == true {
				log.Println(beego.AppConfig.String("loglevel"), "Info", "Redirecting to login")
				c.Redirect(beego.AppConfig.String("LOGIN_PATH"), 302)

			} else {
				c.Data["DisplayMessage"] = err.Error()
			}
			c.TplName = "channelmanagement/createChannel/createChannel.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Creae Admin User Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "channelmanagement/createChannel/createChannel.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Creae Admin User  Page Success")
		}
		return
	}()

	utils.SetHTTPHeader(c.Ctx)

	sess, err := session.GlobalSessions.SessionStart(c.Ctx.ResponseWriter, c.Ctx.Request)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("System is unable to process your request.Please contact customer care")
		sessErr = true
		return
	}

	if err = session.ValidateSession(sess); err != nil {
		sess.SessionRelease(c.Ctx.ResponseWriter)
		session.GlobalSessions.SessionDestroy(c.Ctx.ResponseWriter, c.Ctx.Request)
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		sessErr = true
		return
	}
	defer func() {
		utils.EventLogs(c.Ctx, sess, c.Ctx.Input.Method(), c.Input(), c.Data, err)
		sess.SessionRelease(c.Ctx.ResponseWriter)
	}()

	content := sess.Get("Menus").(string)
	c.Data["Menus"] = template.HTML(`` + content + ``)
	headerContent := sess.Get("Header").(string)
	c.Data["Header"] = template.HTML(`` + headerContent + ``)
	c.Data["language"] = sess.Get("language").(string)
	return
}

func (c *CreateChannel) Post() {

	log.Println(beego.AppConfig.String("loglevel"), "Info", "Create Channel Page Start")
	pip := c.Ctx.Input.IP()
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Client IP - ", pip)
	var err error
	sessErr := false
	defer func() {
		if l_exception := recover(); l_exception != nil {
			stack := debug.Stack()
			log.Println(beego.AppConfig.String("loglevel"), "Exception", string(stack))
			session.GlobalSessions.SessionDestroy(c.Ctx.ResponseWriter, c.Ctx.Request)
			c.TplName = "error/error.html"
		}
		if err != nil {
			if sessErr == true {
				log.Println(beego.AppConfig.String("loglevel"), "Info", "Redirecting to login")
				c.Redirect(beego.AppConfig.String("LOGIN_PATH"), 302)
			} else {
				c.Data["DisplayMessage"] = err.Error()
			}
			c.TplName = "channelmanagement/createChannel/createChannel.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Creae Admin User Page Fail")
		} else {
			c.Data["DisplayMessage"] = " "
			c.TplName = "channelmanagement/createChannel/createChannel.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Creae Admin User  Page Success")
		}
		return
	}()
	utils.SetHTTPHeader(c.Ctx)
	sess, err := session.GlobalSessions.SessionStart(c.Ctx.ResponseWriter, c.Ctx.Request)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("System is unable to process your request.Please contact customer care")
		sessErr = true
		return
	}

	if err = session.ValidateSession(sess); err != nil {
		sess.SessionRelease(c.Ctx.ResponseWriter)
		session.GlobalSessions.SessionDestroy(c.Ctx.ResponseWriter, c.Ctx.Request)
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		sessErr = true
		return
	}
	defer func() {
		utils.EventLogs(c.Ctx, sess, c.Ctx.Input.Method(), c.Input(), c.Data, err)
		sess.SessionRelease(c.Ctx.ResponseWriter)
	}()
	content := sess.Get("Menus").(string)
	c.Data["Menus"] = template.HTML(`` + content + ``)
	headerContent := sess.Get("Header").(string)
	c.Data["Header"] = template.HTML(`` + headerContent + ``)
	c.Data["language"] = sess.Get("language").(string)
	gateway_name := c.Input().Get("input_gatewayname")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Gateway Name - ", gateway_name)

	gateway_status := c.Input().Get("input_gatewaystatus")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Gateway Status - ", gateway_status)

	gateway_image := c.Input().Get("input_gatewayimagename")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Gateway Image - ", gateway_image)

	channel_name := c.Input().Get("input_channelname")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Channel Name - ", channel_name)

	channel_status := c.Input().Get("input_channelstatus")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Channel Status - ", channel_status)

	channel_desc := c.Input().Get("input_channeldesc")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Channel Desc - ", channel_desc)

	channel_image := c.Input().Get("input_channelimagename")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Channel Desc - ", channel_image)

	result, err := db.Db.Exec(`INSERT INTO Payment_Gateway(gateway_name,gateway_status,gateway_image,channel_name,channel_status,channel_image ,channel_desc,create_date) 
	VALUES (?,?,?,now())`,
		channel_name, channel_desc, channel_status)
	if err != nil {
		err = errors.New("Channel creation failed")
		return
	}

	i, err := result.RowsAffected()
	if err != nil || i == 0 {
		err = errors.New("Channel creation failed")
		return
	}

	return
}
