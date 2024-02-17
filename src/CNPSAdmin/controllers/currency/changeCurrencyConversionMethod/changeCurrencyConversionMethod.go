package changeCurrencyConversionMethod

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

type ChangeCurrencyConversionMethod struct {
	beego.Controller
}

func (c *ChangeCurrencyConversionMethod) Get() {
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Change Currency Conversion Method Start")
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
			c.TplName = "currency/changeCurrencyConversionMethod/changeCurrencyConversionMethod.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Change Currency Conversion Method Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "currency/changeCurrencyConversionMethod/changeCurrencyConversionMethod.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Change Currency Conversion Method  Page Success")
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
	auth, err := utils.IsAuthorized(sess.Get("menujson").(string), "ChangeCurrencyConversionMethod")
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
		c.Data["Currency"] = beego.AppConfig.String("ENGLISH_CURRENCY")
		c.Data["AddCurrency"] = beego.AppConfig.String("ENGLISH_ADD_CURRENCY")
		c.Data["Code"] = beego.AppConfig.String("ENGLISH_CODE")
		c.Data["Country"] = beego.AppConfig.String("ENGLISH_COUNTRY")
		c.Data["Symbol"] = beego.AppConfig.String("ENGLISH_SYMBOL")
		c.Data["Back"] = beego.AppConfig.String("ENGLISH_BACK")
		c.Data["Submit"] = beego.AppConfig.String("ENGLISH_SUBMIT")
		c.Data["Cancel"] = beego.AppConfig.String("ENGLISH_CANCEL")
		c.Data["CurrencyDetails"] = beego.AppConfig.String("ENGLISH_ROLE_DETAILS")
		c.Data["EnterCurrencyName"] = beego.AppConfig.String("ENGLISH_ENTER_ROLE_NAME")
		c.Data["InputCurrencyName"] = beego.AppConfig.String("ENGLISH_INPUT_ROLE_NAME")

		c.Data["ProfileManagment"] = beego.AppConfig.String("ENGLISH_PROFILEMANAGEMENT")
		c.Data["Users"] = beego.AppConfig.String("ENGLISH_SEARCH_USERS")
		c.Data["CreateUser"] = beego.AppConfig.String("ENGLISH_CREATE_USER")
		c.Data["ViewUser"] = beego.AppConfig.String("ENGLISH_VIEW_USERS")
		c.Data["UpdateUser"] = beego.AppConfig.String("ENGLISH_UPDATE_USERS")
		c.Data["SystemConfiguration"] = beego.AppConfig.String("ENGLISH_SYSTEM_CONFIGURATION")
		c.Data["ViewSwitchStatus"] = beego.AppConfig.String("ENGLISH_VIEW_SWTICH_STATUS")
		c.Data["UpdateSwitchStatus"] = beego.AppConfig.String("ENGLISH_UPDATE_SWTICH_STATUS")
		c.Data["Viewcurrency"] = beego.AppConfig.String("ENGLISH_VIEW_ROLE")
		c.Data["ChangeCurrencyConversionMethod"] = beego.AppConfig.String("ENGLISH_ADD_ROLE")
		c.Data["Updatecurrency"] = beego.AppConfig.String("ENGLISH_UPDATE_ROLE")
		c.Data["ViewChannels"] = beego.AppConfig.String("ENGLISH_VIEW_CHANNELS")
		c.Data["UpdateChannels"] = beego.AppConfig.String("ENGLISH_UPDATE_CHANNELS")
		c.Data["Reports"] = beego.AppConfig.String("ENGLISH_REPORTS")
		c.Data["AuditReport"] = beego.AppConfig.String("ENGLISH_ADUIT_REPORT")
		c.Data["ChannelReport"] = beego.AppConfig.String("ENGLISH_CHANNEL_REPORT")
		c.Data["TransactionReport"] = beego.AppConfig.String("ENGLISH_TRANSACTION_REPORT")
		c.Data["SelectAll"] = beego.AppConfig.String("ENGLISH_SELECT_ALL_LABEL")
		c.Data["SystemMonitoring"] = beego.AppConfig.String("ENGLISH_SYSTEM_MONITORING")
		c.Data["SysConfig"] = beego.AppConfig.String("ENGLISH_SYSTEM_CONFIGURATION")
		c.Data["Viewprivileges"] = beego.AppConfig.String("ENGLISH_VIEW_PRIVILEGES")
		c.Data["Updateprivileges"] = beego.AppConfig.String("ENGLISH_UPDATE_PRIVILEGES")

		c.Data["ListLabel"] = beego.AppConfig.String("ENGLISH_LIST_LABEL")
		c.Data["MenuLabel"] = beego.AppConfig.String("ENGLISH_MENU_LABEL")
		c.Data["AddLabel"] = beego.AppConfig.String("ENGLISH_ROLE_ADD_LABEL")
		c.Data["UpdateLabel"] = beego.AppConfig.String("ENGLISH_ROLE_UPDATE_LABEL")
		c.Data["ViewLabel"] = beego.AppConfig.String("ENGLISH_ROLE_VIEW_LABEL")
		c.Data["PGWLabel"] = beego.AppConfig.String("ENGLISH_ROLE_PGW_SWITCH_LABEL")
		c.Data["CurrencyLabel"] = beego.AppConfig.String("ENGLISH_ROLE_MENU_LABEL")
		c.Data["PaymentChannelLabel"] = beego.AppConfig.String("ENGLISH_ROLE_PAYMENT_CHANNEL_LABEL")
		c.Data["ViewCurrency"] = beego.AppConfig.String("ENGLISH_VIEW_ROLES_MENU")
		c.Data["ProfileMgmt"] = beego.AppConfig.String("ENGLISH_PROFILEMANAGEMENT")
		c.Data["Language"] = sess.Get("language")
		c.TplName = "currency/changeCurrencyConversionMethod/changeCurrencyConversionMethod.html"
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
		c.Data["Currency"] = beego.AppConfig.String("FRENCH_CURRENCY")
		c.Data["AddCurrency"] = beego.AppConfig.String("FRENCH_ADD_CURRENCY")
		c.Data["Code"] = beego.AppConfig.String("FRENCH_CODE")
		c.Data["Symbol"] = beego.AppConfig.String("FRENCH_SYMBOL")
		c.Data["Country"] = beego.AppConfig.String("FRENCH_COUNTRY")
		c.Data["Back"] = beego.AppConfig.String("FRENCH_BACK")
		c.Data["Submit"] = beego.AppConfig.String("FRENCH_SUBMIT")
		c.Data["Cancel"] = beego.AppConfig.String("FRENCH_CANCEL")
		c.Data["CurrencyDetails"] = beego.AppConfig.String("FRENCH_ROLE_DETAILS")
		c.Data["EnterCurrencyName"] = beego.AppConfig.String("FRENCH_ENTER_ROLE_NAME")
		c.Data["InputCurrencyName"] = beego.AppConfig.String("FRENCH_INPUT_ROLE_NAME")

		c.Data["ProfileManagment"] = beego.AppConfig.String("FRENCH_PROFILEMANAGEMENT")
		c.Data["Users"] = beego.AppConfig.String("FRENCH_SEARCH_USERS")
		c.Data["CreateUser"] = beego.AppConfig.String("FRENCH_CREATE_USER")
		c.Data["ViewUser"] = beego.AppConfig.String("FRENCH_VIEW_USERS")
		c.Data["UpdateUser"] = beego.AppConfig.String("FRENCH_UPDATE_USERS")
		c.Data["SystemConfiguration"] = beego.AppConfig.String("FRENCH_SYSTEM_CONFIGURATION")
		c.Data["ViewSwitchStatus"] = beego.AppConfig.String("FRENCH_VIEW_SWTICH_STATUS")
		c.Data["UpdateSwitchStatus"] = beego.AppConfig.String("FRENCH_UPDATE_SWTICH_STATUS")
		c.Data["Viewcurrency"] = beego.AppConfig.String("FRENCH_VIEW_ROLE")
		c.Data["ChangeCurrencyConversionMethod"] = beego.AppConfig.String("FRENCH_ADD_ROLE")
		c.Data["Updatecurrency"] = beego.AppConfig.String("FRENCH_UPDATE_ROLE")
		c.Data["ViewChannels"] = beego.AppConfig.String("FRENCH_VIEW_CHANNELS")
		c.Data["UpdateChannels"] = beego.AppConfig.String("FRENCH_UPDATE_CHANNELS")
		c.Data["Reports"] = beego.AppConfig.String("FRENCH_REPORTS")
		c.Data["AuditReport"] = beego.AppConfig.String("FRENCH_ADUIT_REPORT")
		c.Data["ChannelReport"] = beego.AppConfig.String("FRENCH_CHANNEL_REPORT")
		c.Data["TransactionReport"] = beego.AppConfig.String("FRENCH_TRANSACTION_REPORT")
		c.Data["SystemMonitoring"] = beego.AppConfig.String("FRENCH_SYSTEM_MONITORING")
		c.Data["SysConfig"] = beego.AppConfig.String("FRENCH_SYSTEM_CONFIGURATION")
		c.Data["Viewprivileges"] = beego.AppConfig.String("FRENCH_VIEW_PRIVILEGES")
		c.Data["Updateprivileges"] = beego.AppConfig.String("FRENCH_UPDATE_PRIVILEGES")

		c.Data["ListLabel"] = beego.AppConfig.String("FRENCH_LIST_LABEL")
		c.Data["MenuLabel"] = beego.AppConfig.String("FRENCH_MENU_LABEL")
		c.Data["AddLabel"] = beego.AppConfig.String("FRENCH_ROLE_ADD_LABEL")
		c.Data["UpdateLabel"] = beego.AppConfig.String("FRENCH_ROLE_UPDATE_LABEL")
		c.Data["ViewLabel"] = beego.AppConfig.String("FRENCH_ROLE_VIEW_LABEL")
		c.Data["PGWLabel"] = beego.AppConfig.String("FRENCH_ROLE_PGW_SWITCH_LABEL")
		c.Data["CurrencyLabel"] = beego.AppConfig.String("FRENCH_ROLE_MENU_LABEL")
		c.Data["PaymentChannelLabel"] = beego.AppConfig.String("FRENCH_ROLE_PAYMENT_CHANNEL_LABEL")
		c.Data["SelectAll"] = beego.AppConfig.String("FRENCH_SELECT_ALL_LABEL")
		c.Data["ViewCurrency"] = beego.AppConfig.String("FRENCH_VIEW_ROLES_MENU")
		c.Data["ProfileMgmt"] = beego.AppConfig.String("FRENCH_PROFILEMANAGEMENT")
		c.Data["Language"] = sess.Get("language")
		c.TplName = "currency/changeCurrencyConversionMethod/changeCurrencyConversionMethod.html"
	}
	c.Data["CurrencyConversionMethod"], err = getCurrencyConversionMethod()

}

func getCurrencyConversionMethod() (currencyConversionMethod string, err error) {
	row, err := db.Db.Query(`SELECT currency_conversion_enables FROM Payment_Swtich`)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Merchants Not Exists.Please check admin.")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Merchants Detail Scan Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("Payment_Swtich Record Not Exists. Please Contact System Admin")
		return
	}
	if len(data) > 1 {
		err = errors.New("More Then One Payment_Switch Record Exists. Please Contact System Admin")
		return
	}

	currencyConversionMethod = data[0][0]
	return
}

func (c *ChangeCurrencyConversionMethod) Post() {

	log.Println(beego.AppConfig.String("loglevel"), "Info", "Change Currency Conversion Method Page Start")
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
		if err != nil {
			if sessErr == true {
				log.Println(beego.AppConfig.String("loglevel"), "Info", "Redirecting to login")
				c.Redirect(beego.AppConfig.String("LOGIN_PATH"), 302)
			} else {
				c.Data["DisplayMessage"] = err.Error()
			}
			c.TplName = "currency/changeCurrencyConversionMethod/changeCurrencyConversionMethod.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Change Currency Conversion Method Page Fail")
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
				c.Data["DisplayMessage"] = "Updated Currency Conversion Method Successfully"
				c.TplName = "currency/changeCurrencyConversionMethod/changeCurrencyConversionMethod.html"
				log.Println(beego.AppConfig.String("loglevel"), "Info", "Change Currency Conversion Method Page Success")
			} else if sess.Get("language") == "French" {
				c.Data["DisplayMessage"] = "Updated Currency Conversion Method Successfully"
				c.TplName = "currency/changeCurrencyConversionMethod/changeCurrencyConversionMethod.html"
				log.Println(beego.AppConfig.String("loglevel"), "Info", "Change Currency Conversion Method  Page Success")
			}
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
	auth, err := utils.IsAuthorized(sess.Get("menujson").(string), "ChangeCurrencyConversionMethod")
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
	// sessionLanguage := sess.Get("language").(string)
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
		c.Data["Currency"] = beego.AppConfig.String("ENGLISH_CURRENCY")
		c.Data["AddCurrency"] = beego.AppConfig.String("ENGLISH_ADD_CURRENCY")
		c.Data["Code"] = beego.AppConfig.String("ENGLISH_CODE")
		c.Data["Country"] = beego.AppConfig.String("ENGLISH_COUNTRY")
		c.Data["Symbol"] = beego.AppConfig.String("ENGLISH_SYMBOL")
		c.Data["Back"] = beego.AppConfig.String("ENGLISH_BACK")
		c.Data["Submit"] = beego.AppConfig.String("ENGLISH_SUBMIT")
		c.Data["Cancel"] = beego.AppConfig.String("ENGLISH_CANCEL")
		c.Data["CurrencyDetails"] = beego.AppConfig.String("ENGLISH_ROLE_DETAILS")
		c.Data["EnterCurrencyName"] = beego.AppConfig.String("ENGLISH_ENTER_ROLE_NAME")
		c.Data["InputCurrencyName"] = beego.AppConfig.String("ENGLISH_INPUT_ROLE_NAME")

		c.Data["ProfileManagment"] = beego.AppConfig.String("ENGLISH_PROFILEMANAGEMENT")
		c.Data["Users"] = beego.AppConfig.String("ENGLISH_SEARCH_USERS")
		c.Data["CreateUser"] = beego.AppConfig.String("ENGLISH_CREATE_USER")
		c.Data["ViewUser"] = beego.AppConfig.String("ENGLISH_VIEW_USERS")
		c.Data["UpdateUser"] = beego.AppConfig.String("ENGLISH_UPDATE_USERS")
		c.Data["SystemConfiguration"] = beego.AppConfig.String("ENGLISH_SYSTEM_CONFIGURATION")
		c.Data["ViewSwitchStatus"] = beego.AppConfig.String("ENGLISH_VIEW_SWTICH_STATUS")
		c.Data["UpdateSwitchStatus"] = beego.AppConfig.String("ENGLISH_UPDATE_SWTICH_STATUS")
		c.Data["ChangeCurrencyConversionMethod"] = beego.AppConfig.String("ENGLISH_ADD_ROLE")
		c.Data["Viewcurrency"] = beego.AppConfig.String("ENGLISH_VIEW_ROLE")
		c.Data["Updatecurrency"] = beego.AppConfig.String("ENGLISH_UPDATE_ROLE")
		c.Data["ViewChannels"] = beego.AppConfig.String("ENGLISH_VIEW_CHANNELS")
		c.Data["UpdateChannels"] = beego.AppConfig.String("ENGLISH_UPDATE_CHANNELS")
		c.Data["Reports"] = beego.AppConfig.String("ENGLISH_REPORTS")
		c.Data["AuditReport"] = beego.AppConfig.String("ENGLISH_ADUIT_REPORT")
		c.Data["ChannelReport"] = beego.AppConfig.String("ENGLISH_CHANNEL_REPORT")
		c.Data["TransactionReport"] = beego.AppConfig.String("ENGLISH_TRANSACTION_REPORT")
		c.Data["SystemMonitoring"] = beego.AppConfig.String("ENGLISH_SYSTEM_MONITORING")
		c.Data["SysConfig"] = beego.AppConfig.String("ENGLISH_SYSTEM_CONFIGURATION")
		c.Data["Viewprivileges"] = beego.AppConfig.String("ENGLISH_VIEW_PRIVILEGES")
		c.Data["Updateprivileges"] = beego.AppConfig.String("ENGLISH_UPDATE_PRIVILEGES")

		c.Data["ListLabel"] = beego.AppConfig.String("ENGLISH_LIST_LABEL")
		c.Data["MenuLabel"] = beego.AppConfig.String("ENGLISH_MENU_LABEL")
		c.Data["AddLabel"] = beego.AppConfig.String("ENGLISH_ROLE_ADD_LABEL")
		c.Data["UpdateLabel"] = beego.AppConfig.String("ENGLISH_ROLE_UPDATE_LABEL")
		c.Data["ViewLabel"] = beego.AppConfig.String("ENGLISH_ROLE_VIEW_LABEL")
		c.Data["PGWLabel"] = beego.AppConfig.String("ENGLISH_ROLE_PGW_SWITCH_LABEL")
		c.Data["CurrencyLabel"] = beego.AppConfig.String("ENGLISH_ROLE_MENU_LABEL")
		c.Data["PaymentChannelLabel"] = beego.AppConfig.String("ENGLISH_ROLE_PAYMENT_CHANNEL_LABEL")
		c.Data["SelectAll"] = beego.AppConfig.String("ENGLISH_SELECT_ALL_LABEL")
		c.Data["ViewCurrency"] = beego.AppConfig.String("ENGLISH_VIEW_ROLES_MENU")
		c.Data["ProfileMgmt"] = beego.AppConfig.String("ENGLISH_PROFILEMANAGEMENT")

		c.Data["Language"] = sess.Get("language")
		c.TplName = "currency/changeCurrencyConversionMethod/changeCurrencyConversionMethod.html"
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
		c.Data["Currency"] = beego.AppConfig.String("FRENCH_CURRENCY")
		c.Data["AddCurrency"] = beego.AppConfig.String("FRENCH_ADD_CURRENCY")
		c.Data["Code"] = beego.AppConfig.String("FRENCH_CODE")
		c.Data["Symbol"] = beego.AppConfig.String("FRENCH_SYMBOL")
		c.Data["Country"] = beego.AppConfig.String("FRENCH_COUNTRY")
		c.Data["Back"] = beego.AppConfig.String("FRENCH_BACK")
		c.Data["Submit"] = beego.AppConfig.String("FRENCH_SUBMIT")
		c.Data["Cancel"] = beego.AppConfig.String("FRENCH_CANCEL")
		c.Data["CurrencyDetails"] = beego.AppConfig.String("FRENCH_ROLE_DETAILS")
		c.Data["EnterCurrencyName"] = beego.AppConfig.String("FRENCH_ENTER_ROLE_NAME")
		c.Data["InputCurrencyName"] = beego.AppConfig.String("FRENCH_INPUT_ROLE_NAME")

		c.Data["ProfileManagment"] = beego.AppConfig.String("FRENCH_PROFILEMANAGEMENT")
		c.Data["Users"] = beego.AppConfig.String("FRENCH_SEARCH_USERS")
		c.Data["CreateUser"] = beego.AppConfig.String("FRENCH_CREATE_USER")
		c.Data["ViewUser"] = beego.AppConfig.String("FRENCH_VIEW_USERS")
		c.Data["UpdateUser"] = beego.AppConfig.String("FRENCH_UPDATE_USERS")
		c.Data["SystemConfiguration"] = beego.AppConfig.String("FRENCH_SYSTEM_CONFIGURATION")
		c.Data["ViewSwitchStatus"] = beego.AppConfig.String("FRENCH_VIEW_SWTICH_STATUS")
		c.Data["UpdateSwitchStatus"] = beego.AppConfig.String("FRENCH_UPDATE_SWTICH_STATUS")
		c.Data["Viewcurrency"] = beego.AppConfig.String("FRENCH_VIEW_ROLE")
		c.Data["ChangeCurrencyConversionMethod"] = beego.AppConfig.String("FRENCH_ADD_ROLE")
		c.Data["Updatecurrency"] = beego.AppConfig.String("FRENCH_UPDATE_ROLE")
		c.Data["ViewChannels"] = beego.AppConfig.String("FRENCH_VIEW_CHANNELS")
		c.Data["UpdateChannels"] = beego.AppConfig.String("FRENCH_UPDATE_CHANNELS")
		c.Data["Reports"] = beego.AppConfig.String("FRENCH_REPORTS")
		c.Data["AuditReport"] = beego.AppConfig.String("FRENCH_ADUIT_REPORT")
		c.Data["ChannelReport"] = beego.AppConfig.String("FRENCH_CHANNEL_REPORT")
		c.Data["TransactionReport"] = beego.AppConfig.String("FRENCH_TRANSACTION_REPORT")
		c.Data["SystemMonitoring"] = beego.AppConfig.String("FRENCH_SYSTEM_MONITORING")
		c.Data["SysConfig"] = beego.AppConfig.String("FRENCH_SYSTEM_CONFIGURATION")
		c.Data["Viewprivileges"] = beego.AppConfig.String("FRENCH_VIEW_PRIVILEGES")
		c.Data["Updateprivileges"] = beego.AppConfig.String("FRENCH_UPDATE_PRIVILEGES")

		c.Data["ListLabel"] = beego.AppConfig.String("FRENCH_LIST_LABEL")
		c.Data["MenuLabel"] = beego.AppConfig.String("FRENCH_MENU_LABEL")
		c.Data["AddLabel"] = beego.AppConfig.String("FRENCH_ROLE_ADD_LABEL")
		c.Data["UpdateLabel"] = beego.AppConfig.String("FRENCH_ROLE_UPDATE_LABEL")
		c.Data["ViewLabel"] = beego.AppConfig.String("FRENCH_ROLE_VIEW_LABEL")
		c.Data["PGWLabel"] = beego.AppConfig.String("FRENCH_ROLE_PGW_SWITCH_LABEL")
		c.Data["CurrencyLabel"] = beego.AppConfig.String("FRENCH_ROLE_MENU_LABEL")
		c.Data["PaymentChannelLabel"] = beego.AppConfig.String("FRENCH_ROLE_PAYMENT_CHANNEL_LABEL")
		c.Data["SelectAll"] = beego.AppConfig.String("FRENCH_SELECT_ALL_LABEL")
		c.Data["ViewCurrency"] = beego.AppConfig.String("FRENCH_VIEW_ROLES_MENU")
		c.Data["ProfileMgmt"] = beego.AppConfig.String("FRENCH_PROFILEMANAGEMENT")
		c.Data["Language"] = sess.Get("language")
		c.TplName = "currency/changeCurrencyConversionMethod/changeCurrencyConversionMethod.html"
	}

	input_currency_conversion_method := c.Input().Get("input_currency_conversion_method")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Input Currency Conversion Method - ", input_currency_conversion_method)

	result, err := db.Db.Exec(`UPDATE Payment_Swtich SET currency_conversion_enables=?, last_update=now(), updated_by=?`,
		input_currency_conversion_method, sess.Get("uname").(string))
	if err != nil {
		err = errors.New("Update currency_conversion_method failed")
		return
	}

	i, err := result.RowsAffected()
	if err != nil || i == 0 {
		err = errors.New("Update currency_conversion_method failed")
		return
	}
	c.Data["CurrencyConversionMethod"], err = getCurrencyConversionMethod()
}
