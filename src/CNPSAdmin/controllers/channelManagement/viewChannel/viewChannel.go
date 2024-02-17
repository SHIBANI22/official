package viewChannel

import (
	"CNPSAdmin/model/db"
	"CNPSAdmin/model/utils"
	"CNPSAdmin/session"
	"errors"
	"html/template"
	"runtime/debug"
	"strings"

	"github.com/astaxie/beego"
	"ominaya.com/database/sql"
	"ominaya.com/util/log"
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
type ViewChannel struct {
	beego.Controller
}

func (c *ViewChannel) Get() {
	AdminId := c.Ctx.Input.Param(":AdminID")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "AdminId", AdminId)
	log.Println(beego.AppConfig.String("loglevel"), "Info", "View Channel Start")
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
			c.TplName = "channelmanagement/viewChannel/viewChannel.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "View Channel Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "channelmanagement/viewChannel/viewChannel.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "View Channel  Page Success")
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
	auth, err := utils.IsAuthorized(sess.Get("menujson").(string), "ViewChannels")
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
		c.Data["ViewChannel"] = beego.AppConfig.String("ENGLISH_VIEW_CHANNELS")
		c.Data["ViewChannel1"] = beego.AppConfig.String("ENGLISH_VIEW_CHANNELS")
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
		c.Data["Statuss"] = beego.AppConfig.String("FRENCH_STATUS")
		c.Data["SearchFilters"] = beego.AppConfig.String("ENGLISH_SEARCH_FILTERS")
		//c.Data["ListOfChannels"] = beego.AppConfig.String("ENGLISH_LIST_OF_CHANNELS")
		c.Data["UpdateChannel"] = beego.AppConfig.String("ENGLISH_UPDATE_CHANNEL")
		c.Data["SysConfig"] = beego.AppConfig.String("ENGLISH_SYSTEM_CONFIGURATION")
		c.Data["Back"] = beego.AppConfig.String("ENGLISH_BACK")
		c.Data["View"] = beego.AppConfig.String("ENGLISH_VIEW")
		c.Data["Channels"] = beego.AppConfig.String("ENGLISH_CHANNELS")

		c.TplName = "channelmanagement/viewChannel/viewChannel.html"
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
		c.Data["ViewChannel"] = beego.AppConfig.String("FRENCH_VIEW_CHANNELS")
		c.Data["ViewChannel1"] = beego.AppConfig.String("FRENCH_VIEW_CHANNELS")
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
		c.Data["UpdateChannel"] = beego.AppConfig.String("FRENCH_UPDATE_CHANNEL")
		c.Data["Back"] = beego.AppConfig.String("FRENCH_BACK")
		c.Data["SysConfig"] = beego.AppConfig.String("FRENCH_SYSTEM_CONFIGURATION")
		c.Data["View"] = beego.AppConfig.String("FRENCH_VIEW")
		c.Data["Channels"] = beego.AppConfig.String("FRENCH_CHANNELS")
		c.TplName = "channelmanagement/viewChannel/viewChannel.html"
	}

	// select pg.id,pg.gateway_name,pc.id,pc.channel_name,pc.channel_status,pc.channel_desc,pc.create_date from Payment_Channel as pc inner join Payment_Gateway as pg where pc.payment_gateway_id = pg.id;

	row, err := db.Db.Query(`select pg.uuid,pg.gateway_name,pg.gateway_status,pg.gateway_image,pc.uuid,pc.channel_name,pc.channel_status,pc.channel_image,pc.channel_desc,pc.create_date from Payment_Channel as pc inner join Payment_Gateway as pg where (pc.uuid= ?) AND (pc.payment_gateway_id = pg.id)`, AdminId)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get user data")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get user data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("User data not found")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", data)

	for i := range data {

		c.Data["AdminID"] = data[i][4]
		c.Data["GatewayName"] = data[i][1]
		c.Data["GatewayStatus"] = data[i][2]
		c.Data["GatewayImage"] = data[i][3]
		c.Data["ChannelName"] = data[i][5]
		c.Data["ChnnelStatus"] = data[i][6]
		c.Data["ChannelImage"] = data[i][7]
		c.Data["ChannelDesc"] = data[i][8]
		c.Data["CreateDate"] = data[i][9]

	}

	return

}
