package updateChannel

import (
	"CNPSAdmin/model/db"
	"CNPSAdmin/model/utils"
	"CNPSAdmin/session"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime/debug"

	"html/template"
	"strings"

	"encoding/json"

	"ominaya.com/database/sql"
	"ominaya.com/util/log"

	"github.com/astaxie/beego"
)

type Row struct {
	ID            string
	GatewayName   string
	GatewayStatus string
	GatewayImage  string
	ChannelName   string
	ChnnelStatus  string
	ChannelImage  string
	ChannelDesc   string
	CreateDate    string
}

type UBA struct {
	Tok string `json:"token1"`
}

type UBAResponse struct {
	Secretkey   string `json:"secretkey"`
	Expirytime  string `json:"expirytime"`
	Returncode  string `json:"returncode"`
	Description string `json:"description"`
}
type UpdateChannel struct {
	beego.Controller
}

func (c *UpdateChannel) Get() {
	AdminId := c.Ctx.Input.Param(":AdminID")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "AdminId", AdminId)
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Channel User Start")
	pip := c.Ctx.Input.IP()
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Client IP - ", pip)
	var err error
	sessErr := false
	var Autherr error
	defer func() {

		if l_exception := recover(); l_exception != nil {
			stack := debug.Stack()
			log.Println(beego.AppConfig.String("loglevel"), "Exception", string(stack))
			session.GlobalSessions.SessionDestroy(c.Ctx.ResponseWriter, c.Ctx.Request)
			c.Data["DisplayMessage"] = "Something went wrong.Please Contact CustomerCare."
			c.TplName = "error/error.html"
		}
		if Autherr != nil {
			c.Data["DisplayMessage"] = Autherr.Error()
			c.TplName = "error/error.html"
			return
		}
		if err != nil {
			if sessErr == true {
				log.Println(beego.AppConfig.String("loglevel"), "Info", "Redirecting to login")
				c.Redirect(beego.AppConfig.String("LOGIN_PATH"), 302)

			} else {
				c.Data["DisplayMessage"] = err.Error()
			}
			c.TplName = "channelmanagement/updateChannel/updateChannel.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Channel Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "channelmanagement/updateChannel/updateChannel.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Channel Page Success")
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
	auth, err := utils.IsAuthorized(sess.Get("menujson").(string), "UpdateChannels")
	passSet := sess.Get("passwordSet").(string)
	if err != nil {
		beego.Error(err)
		Autherr = errors.New("Unable to get Menus")
		return
	}
	if !auth || passSet != "YES" {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")
	defer func() {
		utils.EventLogs(c.Ctx, sess, c.Ctx.Input.Method(), c.Input(), c.Data, err)
		sess.SessionRelease(c.Ctx.ResponseWriter)
	}()
	c.Data["MenuJson"] = sess.Get("menujson")

	c.Data["language"] = sess.Get("language").(string)
	c.Data["Photo"] = sess.Get("photo").(string)
	successAmount, servicecharge, totalTransCount, totalBanks, successCount, pendingCount, declainedCount, transErr := utils.SideBarTransactionData()
	if transErr != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get sidebar data")
		return
	}
	if sess.Get("role") == "ADMIN" && sess.Get("language") == "English" {
		menuContent := strings.Replace(beego.AppConfig.String("MENU_TEMPLATE"), "{{.SuccessAmount}}", beego.AppConfig.String("BaseCurrency")+" "+successAmount, -1)
		menuContent = strings.Replace(menuContent, "{{.TotalTransCount}}", totalTransCount, -1)
		menuContent = strings.Replace(menuContent, "{{.servicecharge}}", beego.AppConfig.String("BaseCurrency")+" "+servicecharge, -1)
		menuContent = strings.Replace(menuContent, "{{.BanksCount}}", totalBanks, -1)
		menuContent = strings.Replace(menuContent, "{{.SuccessCount}}", successCount, -1)
		menuContent = strings.Replace(menuContent, "{{.PendingCount}}", pendingCount, -1)
		menuContent = strings.Replace(menuContent, "{{.DeclainedCount}}", declainedCount, -1)
		c.Data["Menus"] = template.HTML(`` + menuContent + ``)
		headerContent := strings.Replace(beego.AppConfig.String("ENGLISH_HEADER_TEMPLATE"), "{{.Fullname}}", sess.Get("fullname").(string), -1)
		headerContent = strings.Replace(headerContent, "{{.Rolename}}", sess.Get("rolename").(string), -1)
		c.Data["Header"] = template.HTML(`` + headerContent + ``)
		c.Data["Dashboard"] = beego.AppConfig.String("ENGLISH_DASHBOARD")
		c.Data["UpdateChannel"] = beego.AppConfig.String("ENGLISH_UPDATE_CHANNELS")
		c.Data["UpdateChannel1"] = beego.AppConfig.String("ENGLISH_UPDATE_CHANNELS")
		c.Data["ChannelNamee"] = beego.AppConfig.String("ENGLISH_CHANNEL_NAME")
		c.Data["ChannelImagee"] = beego.AppConfig.String("ENGLISH_CHANNEL_IMAGE")
		c.Data["GatewayImagee"] = beego.AppConfig.String("ENGLISH_GATEWAY_IMAGE")
		c.Data["GatewayNamee"] = beego.AppConfig.String("ENGLISH_GATEWAY_NAME")
		c.Data["SelectDateRangee"] = beego.AppConfig.String("ENGLISH_USER_SELECT_DATARANGE")
		c.Data["GatewayStatuss"] = beego.AppConfig.String("ENGLISH_GATEWAY_STATUS")
		c.Data["ChannelStatuss"] = beego.AppConfig.String("ENGLISH_CHANNEL_STATUS")
		c.Data["ID"] = beego.AppConfig.String("ENGLISH_ID")
		c.Data["TimeStamp"] = beego.AppConfig.String("ENGLISH_TIMESTAMP")
		c.Data["Descc"] = beego.AppConfig.String("ENGLISH_DESC")
		c.Data["Statuss"] = beego.AppConfig.String("ENGLISH_STATUS")
		c.Data["SearchFilters"] = beego.AppConfig.String("ENGLISH_SEARCH_FILTERS")
		c.Data["Back"] = beego.AppConfig.String("ENGLISH_BACK")
		c.Data["Submit"] = beego.AppConfig.String("ENGLISH_SUBMIT")
		c.Data["Update"] = beego.AppConfig.String("ENGLISH_UPDATE")
		c.Data["pleaseselectPlaceholder"] = beego.AppConfig.String("ENGLISH_ENTER_PLEASESELECT")
		c.Data["Active"] = beego.AppConfig.String("ENGLISH_ACTIVE")
		c.Data["InActive"] = beego.AppConfig.String("ENGLISH_INACTIVE")
		c.Data["SysConfig"] = beego.AppConfig.String("ENGLISH_SYSTEM_CONFIGURATION")
		//c.Data["ListOfChannels"] = beego.AppConfig.String("ENGLISH_LIST_OF_CHANNELS")
		c.Data["Update"] = beego.AppConfig.String("ENGLISH_UPDATE")
		c.Data["Channels"] = beego.AppConfig.String("ENGLISH_CHANNELS")

		c.TplName = "channelmanagement/updateChannel/updateChannel.html"
	} else if sess.Get("role") == "ADMIN" && sess.Get("language") == "French" {
		menuContent := strings.Replace(beego.AppConfig.String("FRENCH_MENU_TEMPLATE"), "{{.SuccessAmount}}", beego.AppConfig.String("BaseCurrency")+" "+successAmount, -1)
		menuContent = strings.Replace(menuContent, "{{.TotalTransCount}}", totalTransCount, -1)
		menuContent = strings.Replace(menuContent, "{{.servicecharge}}", beego.AppConfig.String("BaseCurrency")+" "+servicecharge, -1)
		menuContent = strings.Replace(menuContent, "{{.BanksCount}}", totalBanks, -1)
		menuContent = strings.Replace(menuContent, "{{.SuccessCount}}", successCount, -1)
		menuContent = strings.Replace(menuContent, "{{.PendingCount}}", pendingCount, -1)
		menuContent = strings.Replace(menuContent, "{{.DeclainedCount}}", declainedCount, -1)
		c.Data["Menus"] = template.HTML(`` + menuContent + ``)
		headerContent := strings.Replace(beego.AppConfig.String("FRENCH_HEADER_TEMPLATE"), "{{.Fullname}}", sess.Get("fullname").(string), -1)
		headerContent = strings.Replace(headerContent, "{{.Rolename}}", sess.Get("rolename").(string), -1)
		c.Data["Header"] = template.HTML(`` + headerContent + ``)

		c.Data["Dashboard"] = beego.AppConfig.String("FRENCH_DASHBOARD")
		c.Data["UpdateChannel"] = beego.AppConfig.String("FRENCH_UPDATE_CHANNELS")
		c.Data["UpdateChannel1"] = beego.AppConfig.String("FRENCH_UPDATE_CHANNELS")
		c.Data["ChannelNamee"] = beego.AppConfig.String("FRENCH_CHANNEL_NAME")
		c.Data["GatewayNamee"] = beego.AppConfig.String("FRENCH_GATEWAY_NAME")
		c.Data["ChannelImagee"] = beego.AppConfig.String("FRENCH_CHANNEL_IMAGE")
		c.Data["GatewayImagee"] = beego.AppConfig.String("FRENCH_GATEWAY_IMAGE")
		c.Data["SelectDateRangee"] = beego.AppConfig.String("FRENCH_USER_SELECT_DATARANGE")

		c.Data["GatewayStatuss"] = beego.AppConfig.String("FRENCH_GATEWAY_STATUS")
		c.Data["ChannelStatuss"] = beego.AppConfig.String("FRENCH_CHANNEL_STATUS")
		c.Data["ID"] = beego.AppConfig.String("FRENCH_ID")
		c.Data["TimeStamp"] = beego.AppConfig.String("FRENCH_TIMESTAMP")
		c.Data["Statuss"] = beego.AppConfig.String("FRENCH_STATUS")
		c.Data["Descc"] = beego.AppConfig.String("FRENCH_DESC")
		c.Data["SearchFilters"] = beego.AppConfig.String("FRENCH_SEARCH_FILTERS")
		//c.Data["ListOfChannels"] = beego.AppConfig.String("FRENCH_LIST_OF_CHANNELS")
		c.Data["pleaseselectPlaceholder"] = beego.AppConfig.String("FRENCH_ENTER_PLEASESELECT")
		c.Data["Back"] = beego.AppConfig.String("FRENCH_BACK")
		c.Data["Submit"] = beego.AppConfig.String("FRENCH_SUBMIT")
		c.Data["Update"] = beego.AppConfig.String("FRENCH_UPDATE")
		c.Data["Active"] = beego.AppConfig.String("FRENCH_ACTIVE")
		c.Data["InActive"] = beego.AppConfig.String("FRENCH_INACTIVE")
		c.Data["SysConfig"] = beego.AppConfig.String("FRENCH_SYSTEM_CONFIGURATION")
		c.Data["Update"] = beego.AppConfig.String("FRENCH_UPDATE")
		c.Data["Channels"] = beego.AppConfig.String("FRENCH_CHANNELS")

		c.TplName = "channelmanagement/updateChannel/updateChannel.html"
	}

	// uname = sess.Get("uname").(string)
	// c.Data["Uname"] = uname

	row, err := db.Db.Query(`select pg.uuid,pg.gateway_name,pg.gateway_status,pg.gateway_image,pc.uuid,pc.channel_name,pc.channel_status,pc.channel_image,pc.channel_desc,pc.create_date from Payment_Channel as pc inner join Payment_Gateway as pg where (pc.uuid= ?) AND (pc.payment_gateway_id = pg.id)`, AdminId)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get data")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("Channel data  found")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", data)

	for i := range data {

		c.Data["GatewayID"] = data[i][0]
		c.Data["GatewayName"] = data[i][1]
		c.Data["GatewayStatus"] = data[i][2]
		c.Data["GatewayImage"] = data[i][3]
		c.Data["AdminID"] = data[i][4]
		c.Data["ChannelName"] = data[i][5]
		c.Data["ChannelStatus"] = data[i][6]
		c.Data["ChannelImage"] = data[i][7]
		c.Data["ChannelDesc"] = data[i][8]
		c.Data["CreateDate"] = data[i][9]

		if data[i][5] == "DIRECT DEBIT" {
			log.Println(beego.AppConfig.String("loglevel"), "Debug", "pavan kulkarni")
			c.Data["UBA"] = "TRUE"

		}

	}

	return

}

func (c *UpdateChannel) Post() {

	AdminId := c.Ctx.Input.Param(":AdminID")

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "AdminId", AdminId)

	log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Channel Page Start")
	pip := c.Ctx.Input.IP()
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Client IP - ", pip)
	var err error
	sessErr := false
	var Autherr error
	defer func() {
		if l_exception := recover(); l_exception != nil {
			stack := debug.Stack()
			log.Println(beego.AppConfig.String("loglevel"), "Exception", string(stack))
			session.GlobalSessions.SessionDestroy(c.Ctx.ResponseWriter, c.Ctx.Request)
			c.TplName = "error/error.html"
		}
		if Autherr != nil {
			c.Data["DisplayMessage"] = Autherr.Error()
			c.TplName = "error/error.html"
			return
		}
		row, rowerr := db.Db.Query(`select pg.uuid,pg.gateway_name,pg.gateway_status,pg.gateway_image,pc.uuid,pc.channel_name,pc.channel_status,pc.channel_image,pc.channel_desc,pc.create_date from Payment_Channel as pc inner join Payment_Gateway as pg where (pc.uuid= ?) AND (pc.payment_gateway_id = pg.id)`, AdminId)
		if rowerr != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", rowerr)
			err = errors.New("Unable to get channel data")
			return
		}
		defer sql.Close(row)
		_, data, dataerr := sql.Scan(row)
		if dataerr != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", dataerr)
			err = errors.New("Unable to get channel data")
			return
		}
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
		if len(data) <= 0 {
			err = errors.New("channel data not found")
			return
		}

		log.Println(beego.AppConfig.String("loglevel"), "Debug", data)

		for i := range data {

			c.Data["GatewayID"] = data[i][0]
			c.Data["GatewayName"] = data[i][1]
			c.Data["GatewayStatus"] = data[i][2]
			c.Data["GatewayImage"] = data[i][3]
			c.Data["AdminID"] = data[i][4]
			c.Data["ChannelName"] = data[i][5]
			c.Data["ChannelStatus"] = data[i][6]
			c.Data["ChannelImage"] = data[i][7]
			c.Data["ChannelDesc"] = data[i][8]
			c.Data["CreateDate"] = data[i][9]
		}
		if err != nil {
			if sessErr == true {
				log.Println(beego.AppConfig.String("loglevel"), "Info", "Redirecting to login")
				c.Redirect(beego.AppConfig.String("LOGIN_PATH"), 302)
			} else {
				c.Data["DisplayMessage"] = err.Error()
			}
			c.TplName = "channelmanagement/updateChannel/updateChannel.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Channel Fail")
		} else {
			utils.SetHTTPHeader(c.Ctx)

			sess, err := session.GlobalSessions.SessionStart(c.Ctx.ResponseWriter, c.Ctx.Request)
			if err != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				err = errors.New("System is unable to process your request.Please contact customer care")
				sessErr = true
				return
			}
			if sess.Get("language") == "English" {
				c.Data["DisplayMessage"] = beego.AppConfig.String("ENGLISH_CHANNEL_UPDATED_SUCCESSFULLY")
				c.TplName = "channelmanagement/updateChannel/updateChannel.html"
				log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Channel Page Success")
			} else if sess.Get("language") == "French" {
				c.Data["DisplayMessage"] = beego.AppConfig.String("FRENCH_CHANNEL_UPDATED_SUCCESSFULLY")
				c.TplName = "channelmanagement/updateChannel/updateChannel.html"
				log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Channel Page Success")
			}
			// c.Data["DisplayMessage"] = "Channel Updated Successfully"
			// c.TplName = "channelmanagement/updateChannel/updateChannel.html"
			// log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Channel Page Success")
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
	auth, err := utils.IsAuthorized(sess.Get("menujson").(string), "UpdateChannels")
	passSet := sess.Get("passwordSet").(string)
	if err != nil {
		beego.Error(err)
		Autherr = errors.New("Unable to get Menus")
		return
	}
	if !auth || passSet != "YES" {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")
	defer func() {
		utils.EventLogs(c.Ctx, sess, c.Ctx.Input.Method(), c.Input(), c.Data, err)
		sess.SessionRelease(c.Ctx.ResponseWriter)
	}()
	c.Data["MenuJson"] = sess.Get("menujson")

	c.Data["language"] = sess.Get("language").(string)
	sessionLanguage := sess.Get("language").(string)
	c.Data["Photo"] = sess.Get("photo").(string)
	successAmount, servicecharge, totalTransCount, totalBanks, successCount, pendingCount, declainedCount, transErr := utils.SideBarTransactionData()
	if transErr != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get sidebar data")
		return
	}
	if sess.Get("role") == "ADMIN" && sess.Get("language") == "English" {
		menuContent := strings.Replace(beego.AppConfig.String("MENU_TEMPLATE"), "{{.SuccessAmount}}", beego.AppConfig.String("BaseCurrency")+" "+successAmount, -1)
		menuContent = strings.Replace(menuContent, "{{.TotalTransCount}}", totalTransCount, -1)
		menuContent = strings.Replace(menuContent, "{{.servicecharge}}", beego.AppConfig.String("BaseCurrency")+" "+servicecharge, -1)
		menuContent = strings.Replace(menuContent, "{{.BanksCount}}", totalBanks, -1)
		menuContent = strings.Replace(menuContent, "{{.SuccessCount}}", successCount, -1)
		menuContent = strings.Replace(menuContent, "{{.PendingCount}}", pendingCount, -1)
		menuContent = strings.Replace(menuContent, "{{.DeclainedCount}}", declainedCount, -1)
		c.Data["Menus"] = template.HTML(`` + menuContent + ``)
		headerContent := strings.Replace(beego.AppConfig.String("ENGLISH_HEADER_TEMPLATE"), "{{.Fullname}}", sess.Get("fullname").(string), -1)
		headerContent = strings.Replace(headerContent, "{{.Rolename}}", sess.Get("rolename").(string), -1)
		c.Data["Header"] = template.HTML(`` + headerContent + ``)
		c.Data["Dashboard"] = beego.AppConfig.String("ENGLISH_DASHBOARD")
		c.Data["UpdateChannel"] = beego.AppConfig.String("ENGLISH_UPDATE_CHANNELS")
		c.Data["UpdateChannel1"] = beego.AppConfig.String("ENGLISH_UPDATE_CHANNELS")
		c.Data["ChannelNamee"] = beego.AppConfig.String("ENGLISH_CHANNEL_NAME")
		c.Data["ChannelImagee"] = beego.AppConfig.String("ENGLISH_CHANNEL_IMAGE")
		c.Data["GatewayImagee"] = beego.AppConfig.String("ENGLISH_GATEWAY_IMAGE")
		c.Data["GatewayNamee"] = beego.AppConfig.String("ENGLISH_GATEWAY_NAME")
		c.Data["SelectDateRangee"] = beego.AppConfig.String("ENGLISH_USER_SELECT_DATARANGE")
		c.Data["GatewayStatuss"] = beego.AppConfig.String("ENGLISH_GATEWAY_STATUS")
		c.Data["ChannelStatuss"] = beego.AppConfig.String("ENGLISH_CHANNEL_STATUS")
		c.Data["ID"] = beego.AppConfig.String("ENGLISH_ID")
		c.Data["TimeStamp"] = beego.AppConfig.String("ENGLISH_TIMESTAMP")
		c.Data["Descc"] = beego.AppConfig.String("ENGLISH_DESC")
		c.Data["Statuss"] = beego.AppConfig.String("ENGLISH_STATUS")
		c.Data["SearchFilters"] = beego.AppConfig.String("ENGLISH_SEARCH_FILTERS")
		c.Data["Back"] = beego.AppConfig.String("ENGLISH_BACK")
		c.Data["Submit"] = beego.AppConfig.String("ENGLISH_SUBMIT")
		c.Data["Update"] = beego.AppConfig.String("ENGLISH_UPDATE")
		c.Data["pleaseselectPlaceholder"] = beego.AppConfig.String("ENGLISH_ENTER_PLEASESELECT")
		c.Data["Active"] = beego.AppConfig.String("ENGLISH_ACTIVE")
		c.Data["InActive"] = beego.AppConfig.String("ENGLISH_INACTIVE")
		c.Data["SysConfig"] = beego.AppConfig.String("ENGLISH_SYSTEM_CONFIGURATION")
		//c.Data["ListOfChannels"] = beego.AppConfig.String("ENGLISH_LIST_OF_CHANNELS")
		c.Data["Update"] = beego.AppConfig.String("ENGLISH_UPDATE")
		c.Data["Channels"] = beego.AppConfig.String("ENGLISH_CHANNELS")

		c.TplName = "channelmanagement/updateChannel/updateChannel.html"
	} else if sess.Get("role") == "ADMIN" && sess.Get("language") == "French" {
		menuContent := strings.Replace(beego.AppConfig.String("FRENCH_MENU_TEMPLATE"), "{{.SuccessAmount}}", beego.AppConfig.String("BaseCurrency")+" "+successAmount, -1)
		menuContent = strings.Replace(menuContent, "{{.TotalTransCount}}", totalTransCount, -1)
		menuContent = strings.Replace(menuContent, "{{.servicecharge}}", beego.AppConfig.String("BaseCurrency")+" "+servicecharge, -1)
		menuContent = strings.Replace(menuContent, "{{.BanksCount}}", totalBanks, -1)
		menuContent = strings.Replace(menuContent, "{{.SuccessCount}}", successCount, -1)
		menuContent = strings.Replace(menuContent, "{{.PendingCount}}", pendingCount, -1)
		menuContent = strings.Replace(menuContent, "{{.DeclainedCount}}", declainedCount, -1)
		c.Data["Menus"] = template.HTML(`` + menuContent + ``)
		headerContent := strings.Replace(beego.AppConfig.String("FRENCH_HEADER_TEMPLATE"), "{{.Fullname}}", sess.Get("fullname").(string), -1)
		headerContent = strings.Replace(headerContent, "{{.Rolename}}", sess.Get("rolename").(string), -1)
		c.Data["Header"] = template.HTML(`` + headerContent + ``)
		c.Data["Dashboard"] = beego.AppConfig.String("FRENCH_DASHBOARD")
		c.Data["UpdateChannel"] = beego.AppConfig.String("FRENCH_UPDATE_CHANNELS")
		c.Data["UpdateChannel1"] = beego.AppConfig.String("FRENCH_UPDATE_CHANNELS")
		c.Data["ChannelNamee"] = beego.AppConfig.String("FRENCH_CHANNEL_NAME")
		c.Data["GatewayNamee"] = beego.AppConfig.String("FRENCH_GATEWAY_NAME")
		c.Data["ChannelImagee"] = beego.AppConfig.String("FRENCH_CHANNEL_IMAGE")
		c.Data["GatewayImagee"] = beego.AppConfig.String("FRENCH_GATEWAY_IMAGE")
		c.Data["SelectDateRangee"] = beego.AppConfig.String("FRENCH_USER_SELECT_DATARANGE")

		c.Data["GatewayStatuss"] = beego.AppConfig.String("FRENCH_GATEWAY_STATUS")
		c.Data["ChannelStatuss"] = beego.AppConfig.String("FRENCH_CHANNEL_STATUS")
		c.Data["ID"] = beego.AppConfig.String("FRENCH_ID")
		c.Data["TimeStamp"] = beego.AppConfig.String("FRENCH_TIMESTAMP")
		c.Data["Statuss"] = beego.AppConfig.String("FRENCH_STATUS")
		c.Data["Descc"] = beego.AppConfig.String("FRENCH_DESC")
		c.Data["SearchFilters"] = beego.AppConfig.String("FRENCH_SEARCH_FILTERS")
		//c.Data["ListOfChannels"] = beego.AppConfig.String("FRENCH_LIST_OF_CHANNELS")
		c.Data["pleaseselectPlaceholder"] = beego.AppConfig.String("FRENCH_ENTER_PLEASESELECT")
		c.Data["Back"] = beego.AppConfig.String("FRENCH_BACK")
		c.Data["Submit"] = beego.AppConfig.String("FRENCH_SUBMIT")
		c.Data["Update"] = beego.AppConfig.String("FRENCH_UPDATE")
		c.Data["Active"] = beego.AppConfig.String("FRENCH_ACTIVE")
		c.Data["InActive"] = beego.AppConfig.String("FRENCH_INACTIVE")
		c.Data["SysConfig"] = beego.AppConfig.String("FRENCH_SYSTEM_CONFIGURATION")
		c.Data["Update"] = beego.AppConfig.String("FRENCH_UPDATE")
		c.Data["Channels"] = beego.AppConfig.String("FRENCH_CHANNELS")

		c.TplName = "channelmanagement/updateChannel/updateChannel.html"
	}

	gateway_id := c.Input().Get("gateway_id")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "gateway_id - ", gateway_id)

	channel_id := c.Input().Get("channel_id")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "channel_id - ", channel_id)

	gateway_status := c.Input().Get("input_gateway_status")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Gateway_status - ", gateway_status)

	channel_status := c.Input().Get("input_channel_status")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "channel_status - ", channel_status)

	if utils.IsDisableCharacters(channel_status) || utils.IsDisableCharacters(gateway_status) {
		beego.Error("suspicious special characters found")
		err = errors.New(beego.AppConfig.String(sessionLanguage + "_PARAMETER_VALIDATION"))
		return
	}
	if channel_status != "ACTIVE" && channel_status != "INACTIVE" {
		beego.Error("suspicious special characters found")
		err = errors.New(beego.AppConfig.String(sessionLanguage + "_PARAMETER_VALIDATION"))
		return
	}

	if gateway_status != "ACTIVE" && gateway_status != "INACTIVE" {
		beego.Error("suspicious special characters found")
		err = errors.New(beego.AppConfig.String(sessionLanguage + "_PARAMETER_VALIDATION"))
		return
	}

	// gateway_image := c.Input().Get("gateway_image")
	// log.Println(beego.AppConfig.String("loglevel"), "Debug", "gateway_image - ", gateway_image)

	// channel_image := c.Input().Get("channel_image")
	// log.Println(beego.AppConfig.String("loglevel"), "Debug", "channel_image - ", channel_image)

	// channel_desc := c.Input().Get("channel_desc")
	// log.Println(beego.AppConfig.String("loglevel"), "Debug", "channel_desc - ", channel_desc)

	user_email := sess.Get("uname")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "user_email - ", user_email)

	input_tokennumber := c.Input().Get("input_tokennumber")

	var ubares string
	if input_tokennumber != "" {
		bytes, err := json.Marshal(UBA{
			Tok: input_tokennumber,
		})
		if err != nil {
			panic(err)
		}

		fmt.Println(string(bytes))
		UBAGenerateCredentailURL := beego.AppConfig.String("UBAGenerateCredentailURL")

		ubares = post(UBAGenerateCredentailURL, string(bytes))

		if err != nil {
			beego.Error(err)
			err = errors.New("Unable to get UBA Status")
			return
		}

		var ubaresJson UBAResponse
		err = json.Unmarshal([]byte(ubares), &ubaresJson)
		if err != nil {
			panic(err)
		}
		c.Data["ubares"] = ubares
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "ubares - ", ubares)

		if ubaresJson.Returncode != "1000" {
			err = errors.New("Channel updation failed")
			log.Println(beego.AppConfig.String("loglevel"), "Debug", ubaresJson.Returncode+" : "+ubaresJson.Description)
			return
		}
	}
	result, err := db.Db.Exec(`UPDATE Payment_Gateway as pg,Payment_Channel as pc SET pg.gateway_status=?,pc.channel_status=?,pc.last_update=now(),pg.updated_by=?,pc.updated_by=?, pc.channel_details=? WHERE (pg.uuid = ?) AND ( pc.uuid = ? )`,
		gateway_status, channel_status, user_email, user_email, ubares, gateway_id, channel_id)
	if err != nil {
		err = errors.New("Channel updation failed")
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Channel error")
		return
	}

	i, rowerr := result.RowsAffected()
	if rowerr != nil && i == 0 {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Channel error", err, i)
		err = errors.New("Channel updation failed")
		return
	}

	return
}

func post(url string, jsonData string) string {
	var jsonStr = []byte(jsonData)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}
