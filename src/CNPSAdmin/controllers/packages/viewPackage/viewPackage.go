package viewPackage

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
type ViewPackage struct {
	beego.Controller
}

func (c *ViewPackage) Get() {
	AdminId := c.Ctx.Input.Param(":AdminID")
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Customer Start")
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
			c.TplName = "packages/viewPackage/viewPackage.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Package Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "packages/viewPackage/viewPackage.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Package  Page Success")
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
	auth, err := utils.IsAuthorized(sess.Get("menujson").(string), "Users")
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
		c.Data["ViewUser"] = beego.AppConfig.String("ENGLISH_VIEW_USERS")
		c.Data["ViewUser1"] = beego.AppConfig.String("ENGLISH_VIEW_USERS")
		c.Data["UserDetails"] = beego.AppConfig.String("ENGLISH_USER_DETAILS")
		c.Data["FirstName"] = beego.AppConfig.String("ENGLISH_PACKAGE_NAME")
		c.Data["MiddleName"] = beego.AppConfig.String("ENGLISH_PACKAGE_VOLUME")
		c.Data["LastName"] = beego.AppConfig.String("ENGLISH_PACKAGE_TXN_FEES")
		c.Data["Mobile"] = beego.AppConfig.String("ENGLISH_MOBILE")
		c.Data["Email"] = beego.AppConfig.String("ENGLISH_EMAIL")
		c.Data["Role"] = beego.AppConfig.String("ENGLISH_ROLE")
		c.Data["ContectNumber"] = beego.AppConfig.String("ENGLISH_CONTACT_NUMBER")
		c.Data["Department"] = beego.AppConfig.String("ENGLISH_DEPT_NAME")
		c.Data["EmployeeId"] = beego.AppConfig.String("ENGLISH_EMPLOYEE_ID")
		c.Data["Status"] = beego.AppConfig.String("ENGLISH_STATUS")
		c.Data["Language"] = beego.AppConfig.String("ENGLISH_LANGUAGE")
		c.Data["CreatedDate"] = beego.AppConfig.String("ENGLISH_CREATE_DATE")
		c.Data["Updateuser"] = beego.AppConfig.String("ENGLISH_UPDATE_PACKAGE")
		c.Data["Back"] = beego.AppConfig.String("ENGLISH_BACK")
		c.Data["ProfileManagement"] = beego.AppConfig.String("ENGLISH_PROFILEMANAGEMENT")
		c.Data["Users"] = beego.AppConfig.String("ENGLISH_SEARCH_Packages")
		c.Data["View"] = beego.AppConfig.String("ENGLISH_VIEW")
		c.Data["Package"] = beego.AppConfig.String("ENGLISH_PACKAGE")
		c.Data["AddPackage"] = beego.AppConfig.String("ENGLISH_ADDPACKAGE")
		c.Data["packageNamePlaceholder"] = beego.AppConfig.String("ENGLISH_PACKAGENAME")
		c.Data["volumePlaceholder"] = beego.AppConfig.String("ENGLISH_VOLUME")
		c.Data["taxPlaceholder"] = beego.AppConfig.String("ENGLISH_TAX")
		c.Data["SysConfig"] = beego.AppConfig.String("ENGLISH_SYSTEM_CONFIGURATION")
		c.TplName = "packages/viewPackage/viewPackage.html"
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
		c.Data["ViewUser"] = beego.AppConfig.String("FRENCH_VIEW_USERS")
		c.Data["ViewUser1"] = beego.AppConfig.String("FRENCH_VIEW_USERS")
		c.Data["UserDetails"] = beego.AppConfig.String("FRENCH_USER_DETAILS")
		c.Data["FirstName"] = beego.AppConfig.String("FRENCH_FIRSTNAME")
		c.Data["MiddleName"] = beego.AppConfig.String("FRENCH_MIDDLENAME")
		c.Data["LastName"] = beego.AppConfig.String("FRENCH_LASTNAME")
		c.Data["Mobile"] = beego.AppConfig.String("FRENCH_MOBILE")
		c.Data["Email"] = beego.AppConfig.String("FRENCH_EMAIL")
		c.Data["Role"] = beego.AppConfig.String("FRENCH_ROLE")
		c.Data["ContectNumber"] = beego.AppConfig.String("FRENCH_CONTACT_NUMBER")
		c.Data["Department"] = beego.AppConfig.String("FRENCH_DEPT_NAME")
		c.Data["EmployeeId"] = beego.AppConfig.String("FRENCH_EMPLOYEE_ID")
		c.Data["Status"] = beego.AppConfig.String("FRENCH_STATUS")
		c.Data["Language"] = beego.AppConfig.String("FRENCH_LANGUAGE")
		c.Data["CreatedDate"] = beego.AppConfig.String("FRENCH_CREATE_DATE")
		c.Data["Updateuser"] = beego.AppConfig.String("FRENCH_UPDATE_USER")
		c.Data["Back"] = beego.AppConfig.String("FRENCH_BACK")
		c.Data["ProfileManagement"] = beego.AppConfig.String("FRENCH_PROFILEMANAGEMENT")
		c.Data["Users"] = beego.AppConfig.String("FRENCH_SEARCH_USERS")
		c.Data["View"] = beego.AppConfig.String("FRENCH_VIEW")

		c.TplName = "packages/viewPackage/viewPackage.html"
	}

	row, err := db.Db.Query(`select id,name,volume,txn_fees,status,created_at from Txn_Package  where id= ?`, AdminId)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get package data")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get package data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("Package data not found")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", data)

	for i := range data {

		c.Data["Id"] = data[i][0]
		c.Data["UserFirstName"] = data[i][1]
		c.Data["UserMiddleName"] = data[i][2]
		c.Data["UserLastName"] = data[i][3]
		c.Data["UserStatus"] = data[i][4]
		c.Data["UserCreateDate"] = data[i][5]
		//c.Data["UserRole"] = data[i][6]
		/*c.Data["UserStatus"] = data[i][7]
		c.Data["UserContactNumber"] = data[i][8]
		c.Data["UserDepartment"] = data[i][9]
		c.Data["UserEmployeeID"] = data[i][10]
		c.Data["UserLanguage"] = data[i][11]
		c.Data["UserCreateDate"] = data[i][12]
		c.Data["UserRole"] = data[i][13]*/

	}

	return

}
