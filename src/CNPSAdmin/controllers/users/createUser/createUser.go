package createUser

import (
	"CNPSAdmin/model/db"
	"CNPSAdmin/model/utils"
	"CNPSAdmin/session"
	"crypto/rand"

	"regexp"

	"ominaya.com/database/sql"
	"ominaya.com/encoding/base64"
	"ominaya.com/util/pbkdf2"

	"errors"
	"html/template"
	"io/ioutil"
	"net/mail"

	"net/smtp"
	"runtime/debug"

	"strings"

	"github.com/scorredoira/email"

	//	"crypto/tls"

	"github.com/astaxie/beego"
	"ominaya.com/util/log"
	"ominaya.com/util/password"
	"ominaya.com/util/txnno"
)

type CreateUser struct {
	beego.Controller
}

type unencryptedAuth struct {
	smtp.Auth
}

func (c *CreateUser) Get() {
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
			c.TplName = "user/createUser/createUser.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Creae Admin User Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "user/createUser/createUser.html"
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
		c.Data["ProfileManagement"] = beego.AppConfig.String("ENGLISH_PROFILEMANAGEMENT")

		c.TplName = "user/createUser/createUser.html"
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

		c.TplName = "user/createUser/createUser.html"
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

func (c *CreateUser) Post() {

	var IsEmail = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`).MatchString

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
			c.TplName = "user/createUser/createUser.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Creae Admin User Page Fail")
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
				c.Data["DisplayMessage"] = beego.AppConfig.String("ENGLISH_USER_ADDED_SUCCESSFULLY")
				c.TplName = "user/createUser/createUser.html"
				log.Println(beego.AppConfig.String("loglevel"), "Info", "Creae Admin User  Page Success")
			} else if sess.Get("language") == "French" {
				c.Data["DisplayMessage"] = beego.AppConfig.String("FRENCH_USER_ADDED_SUCCESSFULLY")
				c.TplName = "user/createUser/createUser.html"
				log.Println(beego.AppConfig.String("loglevel"), "Info", "Creae Admin User  Page Success")
			}

			// c.Data["DisplayMessage"] = "User Created Successfully"
			// c.TplName = "user/createUser/createUser.html"
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

		c.TplName = "user/createUser/createUser.html"
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

		c.TplName = "user/createUser/createUser.html"
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

	user_mobile := c.Input().Get("input_user_mobile")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "user_mobile - ", user_mobile)

	user_email_id := c.Input().Get("input_user_email_id")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "user_email_id - ", user_email_id)

	user_first_name := c.Input().Get("input_user_first_name")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "user_first_name - ", user_first_name)

	user_middele_name := c.Input().Get("input_user_middele_name")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "user_middele_name - ", user_middele_name)

	user_last_name := c.Input().Get("input_user_last_name")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "user_last_name - ", user_last_name)

	user_role := c.Input().Get("input_user_role")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "user_role - ", user_role)

	user_contact_number := c.Input().Get("input_user_contact_number")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "user_contact_number - ", user_contact_number)

	user_department := c.Input().Get("input_user_department")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "user_department - ", user_department)

	user_employee_id := c.Input().Get("input_user_employee_id")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "user_employee_id - ", user_employee_id)

	user_status := c.Input().Get("input_user_status")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "user_status - ", user_status)

	user_language := c.Input().Get("input_user_language")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "user_language - ", user_language)

	if !utils.IsLetter(user_first_name) || !utils.IsLetter(user_last_name) || !utils.IsLetter(user_middele_name) || !utils.IsLetter(user_department) {
		beego.Error("suspicious special characters found")
		err = errors.New(beego.AppConfig.String(sessionLanguage + "_PARAMETER_VALIDATION"))
		return
	}

	if !utils.IsNumber(user_contact_number) || !utils.IsNumber(user_mobile) || !utils.IsNumber(user_role) {
		beego.Error("suspicious special characters found")
		err = errors.New(beego.AppConfig.String(sessionLanguage + "_PARAMETER_VALIDATION"))
		return
	}
	if utils.IsDisableCharacters(user_status) || utils.IsDisableCharacters(user_language) {
		beego.Error("suspicious special characters found")
		err = errors.New(beego.AppConfig.String(sessionLanguage + "_PARAMETER_VALIDATION"))
		return
	}
	if utils.IsDisableCharacters(user_email_id) || utils.IsDisableCharacters(user_employee_id) {
		beego.Error("suspicious special characters found")
		err = errors.New(beego.AppConfig.String(sessionLanguage + "_PARAMETER_VALIDATION"))
		return
	}

	if !IsEmail(user_email_id) {
		beego.Error("invalid email found")
		err = errors.New(beego.AppConfig.String(sessionLanguage + "_PARAMETER_VALIDATION"))
		return
	}

	if user_status != "ACTIVE" && user_status != "INACTIVE" {
		beego.Error("suspicious special characters found")
		err = errors.New(beego.AppConfig.String(sessionLanguage + "_PARAMETER_VALIDATION"))
		return
	}

	if user_language != "English" && user_language != "French" {
		beego.Error("suspicious special characters found")
		err = errors.New(beego.AppConfig.String(sessionLanguage + "_PARAMETER_VALIDATION"))
		return
	}

	loginPass, _ := password.AlphaNumericSpecial(6)
	UUID := beego.AppConfig.String("INSTANCE_ID_PREFIX") + txnno.Generate()

	//password display on logs
	//log.Println(beego.AppConfig.String("loglevel"), "Debug", "loginPass - ", loginPass)

	creater_email := sess.Get("uname")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "creater_email - ", creater_email)

	name, err := utils.Usercheck(user_email_id)

	if name == user_email_id {
		if sess.Get("language") == "English" {
			err = errors.New(beego.AppConfig.String("ENGLISH_USER_EXIST"))
			return
		}
		if sess.Get("language") == "French" {
			err = errors.New(beego.AppConfig.String("FRENCH_USER_EXIST"))
			return
		}

	}

	empid, err := utils.UsercheckEmpid(user_employee_id)

	if empid == user_employee_id {
		if sess.Get("language") == "English" {
			err = errors.New(beego.AppConfig.String("ENGLISH_USER_EMPLOYEEID_EXIST"))
			return
		}
		if sess.Get("language") == "French" {
			err = errors.New(beego.AppConfig.String("FRENCH_USER_EMPLOYEEID_EXIST"))
			return
		}
	}
	err = SendEmail(user_email_id, user_first_name, loginPass)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		if sess.Get("language") == "English" {
			err = errors.New(beego.AppConfig.String("ENGLISH_USER_SENDMAIL__NOT_FOUND"))
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			return
		}
		if sess.Get("language") == "French" {
			err = errors.New(beego.AppConfig.String("FRENCH_USER_SENDMAIL__NOT_FOUND"))
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			return
		}
	}

	pass, err := EncryptPassword(loginPass)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
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
	// 	err = errors.New("User creation failed")
	// 	return
	// }

	result, err := db.Db.Exec(`INSERT INTO Users(mobile,email,password,first_name,middle_name,last_name,role,status,
	contact_number,department,employee_id,language,login_count,password_set,created_date,role_id,created_by,uuid) 
	VALUES (?,?,?,?,?,?,?,?,?,?,?,?,"0","NO",now(),?,?,?)`,
		user_mobile, user_email_id, pass, user_first_name,
		user_middele_name, user_last_name,
		"ADMIN", user_status, user_contact_number,
		user_department, user_employee_id,
		user_language, user_role, creater_email, UUID)
	if err != nil {
		err = errors.New("User creation failed")
		return
	}

	i, err := result.RowsAffected()
	if err != nil || i == 0 {
		err = errors.New("User creation failed")
		return
	}

	return
}

func EncryptPassword(pass string) (out []byte, err error) {

	//commenting display of password in logs
	//log.Println(beego.AppConfig.String("loglevel"), "Debug", "inside encryption password rec: ", pass)
	b := make([]byte, 32)
	_, err = rand.Read(b)
	var pbkdf pbkdf2.Pbkdf2
	pbkdf.Itr = 32
	pbkdf.KeyLen = 32
	pbkdf.Plain = []byte(pass)
	pbkdf.Salt = b
	err = pbkdf.Encrypt()
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to encrypt password")
		return
	}
	var tmp []byte
	tmp = append(tmp, pbkdf.Salt...)
	tmp = append(tmp, pbkdf.Cipher...)

	out, err = base64.Encode(tmp)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to encrypt password")
		return
	}
	//	log.Printf("%s %s \n", "inside encrypt pass after encryption:", out)
	return
}

func SendEmail(emilid string, name string, password string) (err error) {
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "called - ")

	uname := beego.AppConfig.String("EMAIL_NOTIFY_USERNAME")
	pass := beego.AppConfig.String("EMAIL_NOTIFY_PASSWORD")
	url := beego.AppConfig.String("EMAIL_NOTIFY_URL")
	to := beego.AppConfig.String("EMAIL_NOTIFY_TIMEOUT")
	loginurl := beego.AppConfig.String("EMAIL_APPLICATION_LOGIN_URL")
	recipients := strings.Split(emilid, "||")

	tmpFile := beego.AppConfig.String("EMAIL_TEMPLATE")

	buff, err := ioutil.ReadFile(tmpFile)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", "read file -", err)
		return
	}

	msg := string(buff)
	msg = strings.Replace(string(msg), "{{.Name}}", name, -1)
	msg = strings.Replace(string(msg), "{{.Email}}", emilid, -1)
	msg = strings.Replace(string(msg), "{{.Password}}", password, -1)
	msg = strings.Replace(string(msg), "{{.LoginURL}}", loginurl, -1)

	m := email.NewHTMLMessage("Email", msg)
	m.From = mail.Address{Name: "Supernet", Address: uname}
	m.To = recipients

	// send it
	//auth := smtp.PlainAuth("", uname, pass, url)

	config := beego.AppConfig.String("EMAIL_AUTH_CONFIG_MODE")

	if config == "1" {
		auth := smtp.PlainAuth("", uname, pass, url)
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "auth")
		if err = email.Send(url+":"+to, auth, m); err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			return
		}

	} else if config == "2" {
		auth := unencryptedAuth{
			smtp.PlainAuth(
				"",
				uname,
				pass,
				url,
			),
		}
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "no tls auth")
		if err = email.Send(url+":"+to, auth, m); err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			return
		}
	} else {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "no auth")
		if err = email.Send(url+":"+to, nil, m); err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			return
		}
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Email sent successfully")
	return
}

func (a unencryptedAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	s := *server
	s.TLS = true
	return a.Auth.Start(&s)
}
