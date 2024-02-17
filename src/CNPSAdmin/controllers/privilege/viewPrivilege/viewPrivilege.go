package viewPrivilege

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
	Id                string
	UserMobile        string
	UserEmail         string
	UserFirstName     string
	UserMiddleName    string
	UserLastName      string
	UserRole          string
	UserStatus        string
	UserContactNumber string
	UserDepartment    string
	UserEmployeeID    string
	UserLanguage      string
	UserCreateDate    string
}
type ViewPrivilege struct {
	beego.Controller
}

func (c *ViewPrivilege) Get() {
	AdminId := c.Ctx.Input.Param(":AdminID")
	log.Println(beego.AppConfig.String("loglevel"), "Info", "View Privilege Page Start")
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
			c.TplName = "privilege/viewPrivilege/viewPrivilege.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "View Privilege Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "privilege/viewPrivilege/viewPrivilege.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "View Privilege Page Success")
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
	auth, err := utils.IsAuthorized(sess.Get("menujson").(string), "Viewprivileges")
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
		c.Data["Roles"] = beego.AppConfig.String("ENGLISH_ROLE")
		c.Data["Back"] = beego.AppConfig.String("ENGLISH_BACK")
		c.Data["RoleDetails"] = beego.AppConfig.String("ENGLISH_ROLE_DETAILS")
		c.Data["RoleName"] = beego.AppConfig.String("ENGLISH_ROLE_NAME")
		c.Data["ProfileManagment"] = beego.AppConfig.String("ENGLISH_PROFILEMANAGEMENT")
		c.Data["Users"] = beego.AppConfig.String("ENGLISH_SEARCH_USERS")
		c.Data["CreateUser"] = beego.AppConfig.String("ENGLISH_CREATE_USER")
		c.Data["ViewUser"] = beego.AppConfig.String("ENGLISH_VIEW_USERS")
		c.Data["UpdateUser"] = beego.AppConfig.String("ENGLISH_UPDATE_USERS")
		c.Data["SystemConfiguration"] = beego.AppConfig.String("ENGLISH_SYSTEM_CONFIGURATION")
		c.Data["ViewSwitchStatus"] = beego.AppConfig.String("ENGLISH_SWTICH_STATUS")
		c.Data["UpdateSwitchStatus"] = beego.AppConfig.String("ENGLISH_UPDATE_SWTICH_STATUS")
		c.Data["Createrole"] = beego.AppConfig.String("ENGLISH_ADD_ROLE")
		c.Data["Viewrole"] = beego.AppConfig.String("ENGLISH_ROLES")
		c.Data["Updaterole"] = beego.AppConfig.String("ENGLISH_UPDATE_ROLE")
		c.Data["ViewChannels"] = beego.AppConfig.String("ENGLISH_CHANNELS")
		c.Data["UpdateChannels"] = beego.AppConfig.String("ENGLISH_UPDATE_CHANNELS")
		c.Data["Reports"] = beego.AppConfig.String("ENGLISH_REPORTS")
		c.Data["AuditReport"] = beego.AppConfig.String("ENGLISH_ADUIT_REPORT")
		c.Data["ChannelReport"] = beego.AppConfig.String("ENGLISH_CHANNEL_REPORT")
		c.Data["TransactionReport"] = beego.AppConfig.String("ENGLISH_TRANSACTION_REPORTS_MENU")
		c.Data["SystemMonitoring"] = beego.AppConfig.String("ENGLISH_SYSTEM_MONITORING")
		c.Data["ListLabel"] = beego.AppConfig.String("ENGLISH_LIST_LABEL")
		c.Data["MenuLabel"] = beego.AppConfig.String("ENGLISH_MENU_LABEL")
		c.Data["AddLabel"] = beego.AppConfig.String("ENGLISH_ROLE_ADD_LABEL")
		c.Data["UpdateLabel"] = beego.AppConfig.String("ENGLISH_ROLE_UPDATE_LABEL")
		c.Data["ViewLabel"] = beego.AppConfig.String("ENGLISH_ROLE_VIEW_LABEL")
		c.Data["PGWLabel"] = beego.AppConfig.String("ENGLISH_ROLE_PGW_SWITCH_LABEL")
		c.Data["RoleLabel"] = beego.AppConfig.String("ENGLISH_ROLE_MENU_LABEL")
		c.Data["PaymentChannelLabel"] = beego.AppConfig.String("ENGLISH_ROLE_PAYMENT_CHANNEL_LABEL")
		c.Data["SelectAll"] = beego.AppConfig.String("ENGLISH_SELECT_ALL_LABEL")

		c.Data["Privilege"] = beego.AppConfig.String("ENGLISH_PRIVILEGE")
		c.Data["ViewPrivileges"] = beego.AppConfig.String("ENGLISH_VIEW_PRIVILEGES")
		c.Data["UpdatePrivileges"] = beego.AppConfig.String("ENGLISH_UPDATE_PRIVILEGES")
		c.Data["PrivilegeDetails"] = beego.AppConfig.String("ENGLISH_PRIVILEGE_DETAILS")
		c.Data["PrivilegesNew"] = beego.AppConfig.String("ENGLISH_PRIVILEGES_NEW")
		c.Data["View"] = beego.AppConfig.String("ENGLISH_VIEW")
		c.Data["Reporting"] = beego.AppConfig.String("ENGLISH_REPORTING_MENU")
		c.Data["SyestemReports"] = beego.AppConfig.String("ENGLISH_SYSTEM_REPORTS")
		c.Data["MerchantRequests"] = beego.AppConfig.String("ENGLISH_MERCHANT_REQUESTS")
		c.Data["ProfileMgmt"] = beego.AppConfig.String("ENGLISH_PROFILEMANAGEMENT")
		c.Data["ANALYTICS"] = beego.AppConfig.String("ENGLISH_ANALYTICS")
		c.Data["TransactionProcessing"] = beego.AppConfig.String("ENGLISH_TRANSACTIONPROCESSING")
		c.Data["ManualTransactionLabel"] = beego.AppConfig.String("ENGLISH_MANUAL_TRASACTION_LABEL")
		c.Data["DebugLogReport"] = beego.AppConfig.String("ENGLISH_DEBUG_LOG_REPORT_LABEL")

		c.TplName = "privilege/viewPrivilege/viewPrivilege.html"
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
		c.Data["Roles"] = beego.AppConfig.String("FRENCH_ROLE")
		c.Data["Back"] = beego.AppConfig.String("FRENCH_BACK")
		c.Data["RoleDetails"] = beego.AppConfig.String("FRENCH_ROLE_DETAILS")
		c.Data["RoleName"] = beego.AppConfig.String("FRENCH_ROLE_NAME")
		c.Data["ProfileManagment"] = beego.AppConfig.String("FRENCH_PROFILEMANAGEMENT")
		c.Data["Users"] = beego.AppConfig.String("FRENCH_SEARCH_USERS")
		c.Data["CreateUser"] = beego.AppConfig.String("FRENCH_CREATE_USER")
		c.Data["ViewUser"] = beego.AppConfig.String("FRENCH_VIEW_USERS")
		c.Data["UpdateUser"] = beego.AppConfig.String("FRENCH_UPDATE_USERS")
		c.Data["SystemConfiguration"] = beego.AppConfig.String("FRENCH_SYSTEM_CONFIGURATION")
		c.Data["ViewSwitchStatus"] = beego.AppConfig.String("FRENCH_SWTICH_STATUS")
		c.Data["UpdateSwitchStatus"] = beego.AppConfig.String("FRENCH_UPDATE_SWTICH_STATUS")
		c.Data["Createrole"] = beego.AppConfig.String("FRENCH_ADD_ROLE")
		c.Data["Viewrole"] = beego.AppConfig.String("FRENCH_ROLES")
		c.Data["Updaterole"] = beego.AppConfig.String("FRENCH_UPDATE_ROLE")
		c.Data["ViewChannels"] = beego.AppConfig.String("FRENCH_CHANNELS")
		c.Data["UpdateChannels"] = beego.AppConfig.String("FRENCH_UPDATE_CHANNELS")
		c.Data["Reports"] = beego.AppConfig.String("FRENCH_REPORTS")
		c.Data["AuditReport"] = beego.AppConfig.String("FRENCH_ADUIT_REPORT")
		c.Data["ChannelReport"] = beego.AppConfig.String("FRENCH_CHANNEL_REPORT")
		c.Data["TransactionReport"] = beego.AppConfig.String("FRENCH_TRANSACTION_REPORTS_MENU")
		c.Data["SystemMonitoring"] = beego.AppConfig.String("FRENCH_SYSTEM_MONITORING")
		c.Data["ListLabel"] = beego.AppConfig.String("FRENCH_LIST_LABEL")
		c.Data["MenuLabel"] = beego.AppConfig.String("FRENCH_MENU_LABEL")
		c.Data["AddLabel"] = beego.AppConfig.String("FRENCH_ROLE_ADD_LABEL")
		c.Data["UpdateLabel"] = beego.AppConfig.String("FRENCH_ROLE_UPDATE_LABEL")
		c.Data["ViewLabel"] = beego.AppConfig.String("FRENCH_ROLE_VIEW_LABEL")
		c.Data["PGWLabel"] = beego.AppConfig.String("FRENCH_ROLE_PGW_SWITCH_LABEL")
		c.Data["RoleLabel"] = beego.AppConfig.String("FRENCH_ROLE_MENU_LABEL")
		c.Data["PaymentChannelLabel"] = beego.AppConfig.String("FRENCH_ROLE_PAYMENT_CHANNEL_LABEL")
		c.Data["SelectAll"] = beego.AppConfig.String("FRENCH_SELECT_ALL_LABEL")

		c.Data["Privilege"] = beego.AppConfig.String("FRENCH_PRIVILEGE")
		c.Data["ViewPrivileges"] = beego.AppConfig.String("FRENCH_VIEW_PRIVILEGES")
		c.Data["UpdatePrivileges"] = beego.AppConfig.String("FRENCH_UPDATE_PRIVILEGES")
		c.Data["PrivilegeDetails"] = beego.AppConfig.String("FRENCH_PRIVILEGE_DETAILS")
		c.Data["PrivilegesNew"] = beego.AppConfig.String("FRENCH_PRIVILEGES_NEW")
		c.Data["View"] = beego.AppConfig.String("FRENCH_VIEW")
		c.Data["Reporting"] = beego.AppConfig.String("FRENCH_REPORTING_MENU")
		c.Data["SyestemReports"] = beego.AppConfig.String("FRENCH_SYSTEM_REPORTS")
		c.Data["MerchantRequests"] = beego.AppConfig.String("FRENCH_MERCHANT_REQUESTS")
		c.Data["ProfileMgmt"] = beego.AppConfig.String("FRENCH_PROFILEMANAGEMENT")
		c.Data["ANALYTICS"] = beego.AppConfig.String("FRENCH_ANALYTICS")
		c.Data["TransactionProcessing"] = beego.AppConfig.String("FRENCH_TRANSACTIONPROCESSING")
		c.Data["ManualTransactionLabel"] = beego.AppConfig.String("FRENCH_MANUAL_TRASACTION_LABEL")
		c.Data["DebugLogReport"] = beego.AppConfig.String("FRENCH_DEBUG_LOG_REPORT_LABEL")

		c.TplName = "privilege/viewPrivilege/viewPrivilege.html"
	}

	row, err := db.Db.Query(`select uuid,name,menus,created_date,last_update from Roles where uuid= ?`, AdminId)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get  data")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get  data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	if len(data) <= 0 {
		err = errors.New(" data not found")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", data)

	for i := range data {

		c.Data["Id"] = data[i][0]
		c.Data["DepartmentName"] = data[i][1]
		c.Data["DepartmentMenus"] = data[i][2]
		c.Data["DepartmentCreatedate"] = data[i][3]
		c.Data["DepartmentUpdatedate"] = data[i][4]
	}

	return

}
