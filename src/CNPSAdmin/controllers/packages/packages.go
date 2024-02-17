package packages

import (
	"CNPSAdmin/model/db"
	"CNPSAdmin/model/utils"
	"CNPSAdmin/session"

	"ominaya.com/database/sql"

	"errors"
	"html/template"

	"net/smtp"
	"runtime/debug"

	"strings"

	//	"crypto/tls"

	"github.com/astaxie/beego"
	"ominaya.com/util/log"
)

type PackagesController struct {
	beego.Controller
}

type unencryptedAuth struct {
	smtp.Auth
}

func (c *PackagesController) Get() {
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Creae Admin User Start")
	pip := c.Ctx.Input.IP()
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Client IP - ", pip)
	var err error
	var Autherr error
	sessErr := false
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
			c.TplName = "packages/createPackage/package.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Creae Admin User Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "packages/createPackage/package.html"
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
	auth, err := utils.IsAuthorized(sess.Get("menujson").(string), "CreateUser")
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
		c.Data["AddUser"] = beego.AppConfig.String("ENGLISH_ADD_USERS")
		c.Data["AddUser1"] = beego.AppConfig.String("ENGLISH_ADD_USERS")
		c.Data["UserDetails"] = beego.AppConfig.String("ENGLISH_USER_DETAILS")
		c.Data["FirstName"] = beego.AppConfig.String("ENGLISH_FIRSTNAME")
		c.Data["MiddleName"] = beego.AppConfig.String("ENGLISH_MIDDLENAME")
		c.Data["LastName"] = beego.AppConfig.String("ENGLISH_LASTNAME")
		c.Data["Mobile"] = beego.AppConfig.String("ENGLISH_MOBILE")
		c.Data["Email"] = beego.AppConfig.String("ENGLISH_EMAIL")
		c.Data["Role"] = beego.AppConfig.String("ENGLISH_ROLE")
		c.Data["ContectNumber"] = beego.AppConfig.String("ENGLISH_CONTACT_NUMBER")
		c.Data["Department"] = beego.AppConfig.String("ENGLISH_DEPT_NAME")
		c.Data["EmployeeId"] = beego.AppConfig.String("ENGLISH_EMPLOYEE_ID")
		c.Data["Status"] = beego.AppConfig.String("ENGLISH_STATUS")
		c.Data["Language"] = beego.AppConfig.String("ENGLISH_LANGUAGE")
		c.Data["CreatedDate"] = beego.AppConfig.String("ENGLISH_CREATE_DATE")
		c.Data["Back"] = beego.AppConfig.String("ENGLISH_BACK")
		c.Data["Submit"] = beego.AppConfig.String("ENGLISH_SUBMIT")
		c.Data["Cancel"] = beego.AppConfig.String("ENGLISH_CANCEL")
		c.Data["usernamePlaceholder"] = beego.AppConfig.String("ENGLISH_ENTER_USERNAME")
		c.Data["mobilePlaceholde"] = beego.AppConfig.String("ENGLISH_ENTER_MOBILE")
		c.Data["middleNamePlaceholde"] = beego.AppConfig.String("ENGLISH_ENTER_MIDDLENAME")
		c.Data["lastnamePlaceholde"] = beego.AppConfig.String("ENGLISH_ENTER_LASTNAME")
		c.Data["emailPlaceholder"] = beego.AppConfig.String("ENGLISH_ENTER_EMAIL")
		c.Data["pleaseselectPlaceholder"] = beego.AppConfig.String("ENGLISH_ENTER_PLEASESELECT")
		c.Data["contactnumberPlaceholder"] = beego.AppConfig.String("ENGLISH_ENTER_CONTACTNUMBER")
		c.Data["deptPlaceholder"] = beego.AppConfig.String("ENGLISH_ENTER_DEPARTMENT")
		c.Data["empidPlaceholder"] = beego.AppConfig.String("ENGLISH_ENTER_EMPLOYEEID")
		c.Data["Active"] = beego.AppConfig.String("ENGLISH_ACTIVE")
		c.Data["InActive"] = beego.AppConfig.String("ENGLISH_INACTIVE")
		c.Data["English"] = beego.AppConfig.String("ENGLISH_ENGLISH")
		c.Data["French"] = beego.AppConfig.String("ENGLISH_FRENCH")
		c.Data["Dashboard"] = beego.AppConfig.String("ENGLISH_DASHBOARD")
		c.Data["Users"] = beego.AppConfig.String("ENGLISH_SEARCH_USERS")

		//adding package fields
		c.Data["Package"] = beego.AppConfig.String("ENGLISH_PACKAGE")
		c.Data["AddPackage"] = beego.AppConfig.String("ENGLISH_ADDPACKAGE")
		c.Data["packageNamePlaceholder"] = beego.AppConfig.String("ENGLISH_PACKAGENAME")
		c.Data["volumePlaceholder"] = beego.AppConfig.String("ENGLISH_VOLUME")
		c.Data["taxPlaceholder"] = beego.AppConfig.String("ENGLISH_TAX")
		c.Data["SysConfig"] = beego.AppConfig.String("ENGLISH_SYSTEM_CONFIGURATION")
		c.TplName = "packages/createPackage/package.html"
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
		c.Data["AddUser"] = beego.AppConfig.String("FRENCH_ADD_USERS")
		c.Data["AddUser1"] = beego.AppConfig.String("FRENCH_ADD_USERS")
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
		c.Data["Back"] = beego.AppConfig.String("FRENCH_BACK")
		c.Data["Submit"] = beego.AppConfig.String("FRENCH_SUBMIT")
		c.Data["Cancel"] = beego.AppConfig.String("FRENCH_CANCEL")
		c.Data["usernamePlaceholder"] = beego.AppConfig.String("FRENCH_ENTER_USERNAME")
		c.Data["mobilePlaceholde"] = beego.AppConfig.String("FRENCH_ENTER_MOBILE")
		c.Data["middleNamePlaceholde"] = beego.AppConfig.String("FRENCH_ENTER_MIDDLENAME")
		c.Data["lastnamePlaceholde"] = beego.AppConfig.String("FRENCH_ENTER_LASTNAME")
		c.Data["emailPlaceholder"] = beego.AppConfig.String("FRENCH_ENTER_EMAIL")
		c.Data["pleaseselectPlaceholder"] = beego.AppConfig.String("FRENCH_ENTER_PLEASESELECT")
		c.Data["contactnumberPlaceholder"] = beego.AppConfig.String("FRENCH_ENTER_CONTACTNUMBER")
		c.Data["deptPlaceholder"] = beego.AppConfig.String("FRENCH_ENTER_DEPARTMENT")
		c.Data["empidPlaceholder"] = beego.AppConfig.String("FRENCH_ENTER_EMPLOYEEID")
		c.Data["Active"] = beego.AppConfig.String("FRENCH_ACTIVE")
		c.Data["InActive"] = beego.AppConfig.String("FRENCH_INACTIVE")
		c.Data["English"] = beego.AppConfig.String("FRENCH_ENGLISH")
		c.Data["French"] = beego.AppConfig.String("FRENCH_FRENCH")
		c.Data["Dashboard"] = beego.AppConfig.String("FRENCH_DASHBOARD")
		c.Data["Users"] = beego.AppConfig.String("FRENCH_SEARCH_USERS")
		c.Data["ProfileManagement"] = beego.AppConfig.String("FRENCH_PROFILEMANAGEMENT")

		c.TplName = "packages/createPackage/package.html"
	}

	// row, err := db.Db.Query("select id,name from Roles")

	row, err := db.Db.Query("select id,name from Roles where JSON_SEARCH(menus, 'one', 'true') IS NOT NULL")

	//SELECT id,name FROM test.Roles WHERE JSON_SEARCH(menus, 'one', 'true') IS NOT NULL;

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
	c.Data["DepartArray"] = data

	return
}

func (c *PackagesController) Post() {

	//var IsEmail = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`).MatchString

	log.Println(beego.AppConfig.String("loglevel"), "Info", "Creae Admin User Page Start")
	pip := c.Ctx.Input.IP()
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Client IP - ", pip)
	var err error
	var Autherr error
	sessErr := false
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
			c.TplName = "packages/createPackage/package.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Creae Admin User Page Fail")
		} else {
			utils.SetHTTPHeader(c.Ctx)
			_, err := session.GlobalSessions.SessionStart(c.Ctx.ResponseWriter, c.Ctx.Request)
			if err != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				err = errors.New("System is unable to process your request.Please contact customer care")
				sessErr = true
				return
			}

			// c.Data["DisplayMessage"] = "Package Created Successfully"
			// c.TplName = "packages/createPackage/package.html"
			// log.Println(beego.AppConfig.String("loglevel"), "Info", "Creae Admin User  Page Success")
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
	auth, err := utils.IsAuthorized(sess.Get("menujson").(string), "CreateUser")
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
		c.Data["AddUser"] = beego.AppConfig.String("ENGLISH_ADD_USERS")
		c.Data["AddUser1"] = beego.AppConfig.String("ENGLISH_ADD_USERS")
		c.Data["UserDetails"] = beego.AppConfig.String("ENGLISH_USER_DETAILS")
		c.Data["FirstName"] = beego.AppConfig.String("ENGLISH_FIRSTNAME")
		c.Data["MiddleName"] = beego.AppConfig.String("ENGLISH_MIDDLENAME")
		c.Data["LastName"] = beego.AppConfig.String("ENGLISH_LASTNAME")
		c.Data["Mobile"] = beego.AppConfig.String("ENGLISH_MOBILE")
		c.Data["Email"] = beego.AppConfig.String("ENGLISH_EMAIL")
		c.Data["Role"] = beego.AppConfig.String("ENGLISH_ROLE")
		c.Data["ContectNumber"] = beego.AppConfig.String("ENGLISH_CONTACT_NUMBER")
		c.Data["Department"] = beego.AppConfig.String("ENGLISH_DEPT_NAME")
		c.Data["EmployeeId"] = beego.AppConfig.String("ENGLISH_EMPLOYEE_ID")
		c.Data["Status"] = beego.AppConfig.String("ENGLISH_STATUS")
		c.Data["Language"] = beego.AppConfig.String("ENGLISH_LANGUAGE")
		c.Data["CreatedDate"] = beego.AppConfig.String("ENGLISH_CREATE_DATE")
		c.Data["Back"] = beego.AppConfig.String("ENGLISH_BACK")
		c.Data["Submit"] = beego.AppConfig.String("ENGLISH_SUBMIT")
		c.Data["Cancel"] = beego.AppConfig.String("ENGLISH_CANCEL")
		c.Data["usernamePlaceholder"] = beego.AppConfig.String("ENGLISH_ENTER_USERNAME")
		c.Data["mobilePlaceholde"] = beego.AppConfig.String("ENGLISH_ENTER_MOBILE")
		c.Data["middleNamePlaceholde"] = beego.AppConfig.String("ENGLISH_ENTER_MIDDLENAME")
		c.Data["lastnamePlaceholde"] = beego.AppConfig.String("ENGLISH_ENTER_LASTNAME")
		c.Data["emailPlaceholder"] = beego.AppConfig.String("ENGLISH_ENTER_EMAIL")
		c.Data["pleaseselectPlaceholder"] = beego.AppConfig.String("ENGLISH_ENTER_PLEASESELECT")
		c.Data["contactnumberPlaceholder"] = beego.AppConfig.String("ENGLISH_ENTER_CONTACTNUMBER")
		c.Data["deptPlaceholder"] = beego.AppConfig.String("ENGLISH_ENTER_DEPARTMENT")
		c.Data["empidPlaceholder"] = beego.AppConfig.String("ENGLISH_ENTER_EMPLOYEEID")
		c.Data["Active"] = beego.AppConfig.String("ENGLISH_ACTIVE")
		c.Data["InActive"] = beego.AppConfig.String("ENGLISH_INACTIVE")
		c.Data["English"] = beego.AppConfig.String("ENGLISH_ENGLISH")
		c.Data["French"] = beego.AppConfig.String("ENGLISH_FRENCH")
		c.Data["Dashboard"] = beego.AppConfig.String("ENGLISH_DASHBOARD")
		c.Data["Users"] = beego.AppConfig.String("ENGLISH_SEARCH_USERS")
		c.Data["ProfileManagement"] = beego.AppConfig.String("ENGLISH_PROFILEMANAGEMENT")
		c.Data["Package"] = beego.AppConfig.String("ENGLISH_PACKAGE")
		c.Data["AddPackage"] = beego.AppConfig.String("ENGLISH_ADDPACKAGE")
		c.Data["packageNamePlaceholder"] = beego.AppConfig.String("ENGLISH_PACKAGENAME")
		c.Data["volumePlaceholder"] = beego.AppConfig.String("ENGLISH_VOLUME")
		c.Data["taxPlaceholder"] = beego.AppConfig.String("ENGLISH_TAX")
		c.Data["SysConfig"] = beego.AppConfig.String("ENGLISH_SYSTEM_CONFIGURATION")
		c.Data["SysConfig"] = beego.AppConfig.String("ENGLISH_SYSTEM_CONFIGURATION")
		c.TplName = "packages/createPackage/package.html"
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
		c.Data["AddUser"] = beego.AppConfig.String("FRENCH_ADD_USERS")
		c.Data["AddUser1"] = beego.AppConfig.String("FRENCH_ADD_USERS")
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
		c.Data["Back"] = beego.AppConfig.String("FRENCH_BACK")
		c.Data["Submit"] = beego.AppConfig.String("FRENCH_SUBMIT")
		c.Data["Cancel"] = beego.AppConfig.String("FRENCH_CANCEL")
		c.Data["usernamePlaceholder"] = beego.AppConfig.String("FRENCH_ENTER_USERNAME")
		c.Data["mobilePlaceholde"] = beego.AppConfig.String("FRENCH_ENTER_MOBILE")
		c.Data["middleNamePlaceholde"] = beego.AppConfig.String("FRENCH_ENTER_MIDDLENAME")
		c.Data["lastnamePlaceholde"] = beego.AppConfig.String("FRENCH_ENTER_LASTNAME")
		c.Data["emailPlaceholder"] = beego.AppConfig.String("FRENCH_ENTER_EMAIL")
		c.Data["pleaseselectPlaceholder"] = beego.AppConfig.String("FRENCH_ENTER_PLEASESELECT")
		c.Data["contactnumberPlaceholder"] = beego.AppConfig.String("FRENCH_ENTER_CONTACTNUMBER")
		c.Data["deptPlaceholder"] = beego.AppConfig.String("FRENCH_ENTER_DEPARTMENT")
		c.Data["empidPlaceholder"] = beego.AppConfig.String("FRENCH_ENTER_EMPLOYEEID")
		c.Data["Active"] = beego.AppConfig.String("FRENCH_ACTIVE")
		c.Data["InActive"] = beego.AppConfig.String("FRENCH_INACTIVE")
		c.Data["English"] = beego.AppConfig.String("FRENCH_ENGLISH")
		c.Data["French"] = beego.AppConfig.String("FRENCH_FRENCH")
		c.Data["Dashboard"] = beego.AppConfig.String("FRENCH_DASHBOARD")
		c.Data["Users"] = beego.AppConfig.String("FRENCH_SEARCH_USERS")
		c.Data["ProfileManagement"] = beego.AppConfig.String("FRENCH_PROFILEMANAGEMENT")

		c.TplName = "packages/createPackage/package.html"
	}

	row, err := db.Db.Query("select id,name from Roles where JSON_SEARCH(menus, 'one', 'true') IS NOT NULL")

	//SELECT id,name FROM test.Roles WHERE JSON_SEARCH(menus, 'one', 'true') IS NOT NULL;

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

	// log.Println(beego.AppConfig.String("loglevel"), "Debug", data)

	c.Data["DepartArray"] = data

	//user_mobile := c.Input().Get("input_user_mobile")
	//log.Println(beego.AppConfig.String("loglevel"), "Debug", "user_mobile - ", user_mobile)

	//user_email_id := c.Input().Get("input_user_email_id")
	//log.Println(beego.AppConfig.String("loglevel"), "Debug", "user_email_id - ", user_email_id)

	name := c.Input().Get("input_user_first_name1")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "package_name - ", name)

	volume := c.Input().Get("input_user_middele_name")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "volume - ", volume)

	transaction_fee := c.Input().Get("input_user_last_name1")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "transaction_fee - ", transaction_fee)

	//user_status := c.Input().Get("input_user_status")
	//log.Println(beego.AppConfig.String("loglevel"), "Debug", "user_status - ", user_status)

	if !utils.IsLetter(name) {
		beego.Error("suspicious special characters found")
		err = errors.New(beego.AppConfig.String(sessionLanguage + "_PARAMETER_VALIDATION"))
		return
	}

	if !utils.IsNumber(volume) || !utils.IsNumber(transaction_fee) {
		beego.Error("suspicious special characters found")
		err = errors.New(beego.AppConfig.String(sessionLanguage + "_PARAMETER_VALIDATION"))
		return
	}

	// result, err := db.Db.Exec(`INSERT INTO Users(mobile,email,password,first_name,middle_name,last_name,role,status,
	// contact_number,department,employee_id,language,login_count,password_set,created_date,role_id,created_by)
	// VALUES (?,?,?,?,?,?,?,?,?,?,?,?,"0","NO",now(),?,?)`,
	// 	user_mobile, user_email_id, pass, user_first_name,
	// 	user_middele_name, user_last_name,
	// 	"ADMIN", user_status, user_contact_number,
	// 	user_department, user_employee_id,
	// 	user_language, user_role, creater_email)
	// if err != nil {
	// 	err = errors.New("Package creation failed")
	// 	return
	// }

	nameExists, err := checkNameExists(name)
	volumeExists, err := checkVolumeExists(volume)

	if err != nil {
		beego.Error("Error checking merchant existence:", err)
	}
	if nameExists && volumeExists {
		c.Data["DisplayMessage"] = "Name and Volume already exists"
		// c.Data["nameErrorMsg"] = "name already exists"
		// c.Data["volumeErrorMsg"] = "volume already exists"
		c.TplName = "packages/createPackage/package.html"
		return
	}
	if nameExists {
		c.Data["DisplayMessage"] = "Name already exists"
		// c.Data["nameErrorMsg"] = "name already exists"
		c.TplName = "packages/createPackage/package.html"
		return
	}
	if volumeExists {
		c.Data["DisplayMessage"] = "Volume already exists"
		// c.Data["volumeErrorMsg"] = "volume already exists"
		c.TplName = "packages/createPackage/package.html"
		return
	}

	result, err := db.Db.Exec(`INSERT INTO PGS.Txn_Package (name, volume, txn_fees,created_by) VALUES (?, ?, ?,?)`, name, volume, transaction_fee, "Admin")
	if err != nil {
		err = errors.New("Package creation failed")
		return
	}

	i, err := result.RowsAffected()
	if err != nil || i == 0 {
		err = errors.New("Package creation failed")
		return
	}

	if sess.Get("language") == "English" {
		c.Data["DisplayMessage"] = beego.AppConfig.String("ENGLISH_PACKAGE_ADDED_SUCCESSFULLY")
		c.TplName = "packages/createPackage/package.html"
		log.Println(beego.AppConfig.String("loglevel"), "Info", "Creae Admin User  Page Success")
	} else if sess.Get("language") == "French" {
		c.Data["DisplayMessage"] = beego.AppConfig.String("FRENCH_USER_ADDED_SUCCESSFULLY")
		c.TplName = "packages/createPackage/package.html"
		log.Println(beego.AppConfig.String("loglevel"), "Info", "Creae Admin User  Page Success")
	}

	return
}
func checkNameExists(name string) (bool, error) {

	row3, err := db.Db.Query(`SELECT name FROM Txn_Package where name=?`, name)
	var nameExists bool
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		//err = errors.New("Unable to get data")
		return false, err
	}
	defer sql.Close(row3)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row3)
	_, data, err := sql.Scan(row3)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		//err = errors.New("Unable to get  data")
		return false, err
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	if len(data) <= 0 {
		//err = errors.New("Unable to get  data")
		return false, err
	}

	temp := data[0][0]

	if temp == name {
		nameExists = true
	}

	return nameExists, nil
}
func checkVolumeExists(volume string) (bool, error) {

	row3, err := db.Db.Query(`SELECT volume FROM Txn_Package where volume=?`, volume)
	var volumeExists bool
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		//err = errors.New("Unable to get data")
		return false, err
	}
	defer sql.Close(row3)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row3)
	_, data1, err := sql.Scan(row3)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		//err = errors.New("Unable to get  data")
		return false, err
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data1, "\nData len - ", len(data1))
	if len(data1) <= 0 {
		//err = errors.New("Unable to get  data")
		return false, err
	}
	temp1 := data1[0][0]
	if temp1 == volume {
		volumeExists = true
	}
	return volumeExists, nil
}
