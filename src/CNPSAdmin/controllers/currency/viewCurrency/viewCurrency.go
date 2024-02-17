package viewCurrency

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

type ViewCurrency struct {
	beego.Controller
}

func (c *ViewCurrency) Get() {
	Currency := c.Ctx.Input.Param(":currency")
	log.Println(beego.AppConfig.String("loglevel"), "Info", "View Currency Page Start")
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
			c.TplName = "currency/viewCurrency/viewCurrency.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "View Currency Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "currency/viewCurrency/viewCurrency.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "View Currency Page Success")
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
	auth, err := utils.IsAuthorized(sess.Get("menujson").(string), "ViewCurrency")
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
		c.Data["CreateDate"] = beego.AppConfig.String("ENGLISH_CREATE_DATE")
		c.Data["UpdateDate"] = beego.AppConfig.String("ENGLISH_UPDATED_DATE")
		c.Data["Back"] = beego.AppConfig.String("ENGLISH_BACK")
		c.Data["Currency"] = beego.AppConfig.String("ENGLISH_CURRENCY")
		c.Data["AddCurrency"] = beego.AppConfig.String("ENGLISH_ADD_CURRENCY")
		c.Data["Codes"] = beego.AppConfig.String("ENGLISH_CODE")
		c.Data["Symbols"] = beego.AppConfig.String("ENGLISH_SYMBOL")
		c.Data["Countrys"] = beego.AppConfig.String("ENGLISH_COUNTRY")
		c.Data["RoleDetails"] = beego.AppConfig.String("ENGLISH_ROLE_DETAILS")
		c.Data["RoleName"] = beego.AppConfig.String("ENGLISH_ROLE_NAME")
		c.Data["Statuss"] = beego.AppConfig.String("ENGLISH_STATUS")
		c.Data["ProfileManagment"] = beego.AppConfig.String("ENGLISH_PROFILEMANAGEMENT")
		c.Data["Users"] = beego.AppConfig.String("ENGLISH_SEARCH_USERS")
		c.Data["CreateUser"] = beego.AppConfig.String("ENGLISH_CREATE_USER")
		c.Data["ViewUser"] = beego.AppConfig.String("ENGLISH_VIEW_USERS")
		c.Data["UpdateUser"] = beego.AppConfig.String("ENGLISH_UPDATE_USERS")
		c.Data["SystemConfiguration"] = beego.AppConfig.String("ENGLISH_SYSTEM_CONFIGURATION")
		c.Data["ViewSwitchStatus"] = beego.AppConfig.String("ENGLISH_VIEW_SWTICH_STATUS")
		c.Data["UpdateSwitchStatus"] = beego.AppConfig.String("ENGLISH_UPDATE_SWTICH_STATUS")
		c.Data["Createrole"] = beego.AppConfig.String("ENGLISH_ADD_ROLE")
		c.Data["Viewrole"] = beego.AppConfig.String("ENGLISH_VIEW_ROLE")
		c.Data["Updaterole"] = beego.AppConfig.String("ENGLISH_UPDATE_ROLE")
		c.Data["ViewChannels"] = beego.AppConfig.String("ENGLISH_VIEW_CHANNELS")
		c.Data["UpdateChannels"] = beego.AppConfig.String("ENGLISH_UPDATE_CHANNELS")
		c.Data["Reports"] = beego.AppConfig.String("ENGLISH_REPORTS")
		c.Data["AuditReport"] = beego.AppConfig.String("ENGLISH_ADUIT_REPORT")
		c.Data["ChannelReport"] = beego.AppConfig.String("ENGLISH_CHANNEL_REPORT")
		c.Data["TransactionReport"] = beego.AppConfig.String("ENGLISH_TRANSACTION_REPORT")

		c.Data["SystemMonitoring"] = beego.AppConfig.String("ENGLISH_SYSTEM_MONITORING")
		c.Data["SysConfig"] = beego.AppConfig.String("ENGLISH_SYSTEM_CONFIGURATION")

		c.Data["ListLabel"] = beego.AppConfig.String("ENGLISH_LIST_LABEL")
		c.Data["MenuLabel"] = beego.AppConfig.String("ENGLISH_MENU_LABEL")
		c.Data["AddLabel"] = beego.AppConfig.String("ENGLISH_ROLE_ADD_LABEL")
		c.Data["UpdateLabel"] = beego.AppConfig.String("ENGLISH_ROLE_UPDATE_LABEL")
		c.Data["ViewLabel"] = beego.AppConfig.String("ENGLISH_ROLE_VIEW_LABEL")
		c.Data["PGWLabel"] = beego.AppConfig.String("ENGLISH_ROLE_PGW_SWITCH_LABEL")
		c.Data["RoleLabel"] = beego.AppConfig.String("ENGLISH_ROLE_MENU_LABEL")
		c.Data["PaymentChannelLabel"] = beego.AppConfig.String("ENGLISH_ROLE_PAYMENT_CHANNEL_LABEL")
		c.Data["SelectAll"] = beego.AppConfig.String("ENGLISH_SELECT_ALL_LABEL")
		c.Data["View"] = beego.AppConfig.String("ENGLISH_VIEW")
		c.Data["ProfileMgmt"] = beego.AppConfig.String("ENGLISH_PROFILEMANAGEMENT")

		c.TplName = "currency/viewCurrency/viewCurrency.html"
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
		c.Data["CreateDate"] = beego.AppConfig.String("FRENCH_CREATE_DATE")
		c.Data["UpdateDate"] = beego.AppConfig.String("FRENCH_UPDATED_DATE")
		c.Data["Back"] = beego.AppConfig.String("FRENCH_BACK")
		c.Data["Currency"] = beego.AppConfig.String("FRENCH_CURRENCY")
		c.Data["AddCurrency"] = beego.AppConfig.String("FRENCH_ADD_CURRENCY")
		c.Data["Codes"] = beego.AppConfig.String("FRENCH_CODE")
		c.Data["Symbols"] = beego.AppConfig.String("FRENCH_SYMBOL")
		c.Data["Countrys"] = beego.AppConfig.String("FRENCH_COUNTRY")
		c.Data["Statuss"] = beego.AppConfig.String("FRENCH_STATUS")
		c.Data["RoleDetails"] = beego.AppConfig.String("FRENCH_ROLE_DETAILS")
		c.Data["RoleName"] = beego.AppConfig.String("FRENCH_ROLE_NAME")

		c.Data["ProfileManagment"] = beego.AppConfig.String("FRENCH_PROFILEMANAGEMENT")
		c.Data["Users"] = beego.AppConfig.String("FRENCH_SEARCH_USERS")
		c.Data["CreateUser"] = beego.AppConfig.String("FRENCH_CREATE_USER")
		c.Data["ViewUser"] = beego.AppConfig.String("FRENCH_VIEW_USERS")
		c.Data["UpdateUser"] = beego.AppConfig.String("FRENCH_UPDATE_USERS")
		c.Data["SystemConfiguration"] = beego.AppConfig.String("FRENCH_SYSTEM_CONFIGURATION")
		c.Data["ViewSwitchStatus"] = beego.AppConfig.String("FRENCH_VIEW_SWTICH_STATUS")
		c.Data["UpdateSwitchStatus"] = beego.AppConfig.String("FRENCH_UPDATE_SWTICH_STATUS")
		c.Data["Createrole"] = beego.AppConfig.String("FRENCH_ADD_ROLE")
		c.Data["Viewrole"] = beego.AppConfig.String("FRENCH_VIEW_ROLE")
		c.Data["Updaterole"] = beego.AppConfig.String("FRENCH_UPDATE_ROLE")
		c.Data["ViewChannels"] = beego.AppConfig.String("FRENCH_VIEW_CHANNELS")
		c.Data["UpdateChannels"] = beego.AppConfig.String("FRENCH_UPDATE_CHANNELS")
		c.Data["Reports"] = beego.AppConfig.String("FRENCH_REPORTS")
		c.Data["AuditReport"] = beego.AppConfig.String("FRENCH_ADUIT_REPORT")
		c.Data["ChannelReport"] = beego.AppConfig.String("FRENCH_CHANNEL_REPORT")
		c.Data["TransactionReport"] = beego.AppConfig.String("FRENCH_TRANSACTION_REPORT")

		c.Data["SystemMonitoring"] = beego.AppConfig.String("FRENCH_SYSTEM_MONITORING")
		c.Data["SysConfig"] = beego.AppConfig.String("FRENCH_SYSTEM_CONFIGURATION")

		c.Data["ListLabel"] = beego.AppConfig.String("FRENCH_LIST_LABEL")
		c.Data["MenuLabel"] = beego.AppConfig.String("FRENCH_MENU_LABEL")
		c.Data["AddLabel"] = beego.AppConfig.String("FRENCH_ROLE_ADD_LABEL")
		c.Data["UpdateLabel"] = beego.AppConfig.String("FRENCH_ROLE_UPDATE_LABEL")
		c.Data["ViewLabel"] = beego.AppConfig.String("FRENCH_ROLE_VIEW_LABEL")
		c.Data["PGWLabel"] = beego.AppConfig.String("FRENCH_ROLE_PGW_SWITCH_LABEL")
		c.Data["RoleLabel"] = beego.AppConfig.String("FRENCH_ROLE_MENU_LABEL")
		c.Data["PaymentChannelLabel"] = beego.AppConfig.String("FRENCH_ROLE_PAYMENT_CHANNEL_LABEL")
		c.Data["SelectAll"] = beego.AppConfig.String("FRENCH_SELECT_ALL_LABEL")
		c.Data["View"] = beego.AppConfig.String("FRENCH_VIEW")
		c.Data["ProfileMgmt"] = beego.AppConfig.String("FRENCH_PROFILEMANAGEMENT")
		c.TplName = "currency/viewCurrency/viewCurrency.html"
	}

	row, err := db.Db.Query(`select code,symbol,country,status from Currency where code= ?`, Currency)
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
		err = errors.New("Currency data not found")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", data)

	for i := range data {

		c.Data["Code"] = data[i][0]
		c.Data["Symbol"] = data[i][1]
		c.Data["Country"] = data[i][2]
		c.Data["Status"] = data[i][3]
	}

	return

}
