package updateUser

import (
	"CNPSAdmin/model/db"
	"CNPSAdmin/model/utils"
	"CNPSAdmin/session"
	"errors"
	"runtime/debug"
	"strings"

	"ominaya.com/database/sql"
	"ominaya.com/util/log"

	"html/template"

	"github.com/astaxie/beego"
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

type UpdateUser struct {
	beego.Controller
}

func (c *UpdateUser) Get() {
	AdminId := c.Ctx.Input.Param(":AdminID")
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Admin User Start")
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
			c.TplName = "user/updateUser/updateUser.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Admin User Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "user/updateUser/updateUser.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Admin User Page Success")
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
	auth, err := utils.IsAuthorized(sess.Get("menujson").(string), "UpdateUser")
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
		c.Data["UpdateUser"] = beego.AppConfig.String("ENGLISH_UPDATE_USERS")
		c.Data["UpdateUser1"] = beego.AppConfig.String("ENGLISH_UPDATE_USERS")
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
		c.Data["ProfileManagement"] = beego.AppConfig.String("ENGLISH_PROFILEMANAGEMENT")
		c.Data["input_user_email_id"] = beego.AppConfig.String("ENGLISH_USEREMAIL")
		c.Data["please_select"] = beego.AppConfig.String("ENGLISH_PLEASESELECTCOMMON")
		c.Data["Active"] = beego.AppConfig.String("ENGLISH_ACTIVE")
		c.Data["InActive"] = beego.AppConfig.String("ENGLISH_INACTIVE")
		c.Data["English"] = beego.AppConfig.String("ENGLISH_ENGLISH")
		c.Data["French"] = beego.AppConfig.String("ENGLISH_FRENCH")
		c.Data["Update"] = beego.AppConfig.String("ENGLISH_UPDATE")
		c.Data["Users"] = beego.AppConfig.String("ENGLISH_SEARCH_USERS")

		c.TplName = "user/updateUser/updateUser.html"
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
		c.Data["UpdateUser"] = beego.AppConfig.String("FRENCH_UPDATE_USERS")
		c.Data["UpdateUser1"] = beego.AppConfig.String("FRENCH_UPDATE_USERS")
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
		c.Data["ProfileManagement"] = beego.AppConfig.String("FRENCH_PROFILEMANAGEMENT")
		c.Data["input_user_email_id"] = beego.AppConfig.String("FRENCH_USEREMAIL")
		c.Data["please_select"] = beego.AppConfig.String("FRENCH_PLEASESELECTCOMMON")
		c.Data["Active"] = beego.AppConfig.String("FRENCH_ACTIVE")
		c.Data["InActive"] = beego.AppConfig.String("FRENCH_INACTIVE")
		c.Data["English"] = beego.AppConfig.String("FRENCH_ENGLISH")
		c.Data["French"] = beego.AppConfig.String("FRENCH_FRENCH")
		c.Data["Update"] = beego.AppConfig.String("FRENCH_UPDATE")
		c.Data["Users"] = beego.AppConfig.String("FRENCH_SEARCH_USERS")

		c.TplName = "user/updateUser/updateUser.html"
	}

	// uname = sess.Get("uname").(string)
	// c.Data["Uname"] = uname

	row, err := db.Db.Query(`select uuid,mobile,email,first_name,middle_name,last_name,role,status,contact_number,department,employee_id,language,created_date,role_id from Users where uuid= ?`, AdminId)
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

		c.Data["Id"] = data[i][0]
		c.Data["UserMobile"] = data[i][1]
		c.Data["UserEmail"] = data[i][2]
		c.Data["UserFirstName"] = data[i][3]
		c.Data["UserMiddleName"] = data[i][4]
		c.Data["UserLastName"] = data[i][5]
		c.Data["UserRole"] = data[i][6]
		c.Data["UserStatus"] = data[i][7]
		c.Data["UserContactNumber"] = data[i][8]
		c.Data["UserDepartment"] = data[i][9]
		c.Data["UserEmployeeID"] = data[i][10]
		c.Data["UserLanguage"] = data[i][11]
		c.Data["UserCreateDate"] = data[i][12]
		c.Data["UserRoleId"] = data[i][13]

	}

	row1, err := db.Db.Query("select id,name from Roles where JSON_SEARCH(menus, 'one', 'true') IS NOT NULL")

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get user data")
		return
	}
	defer sql.Close(row1)
	_, data1, err := sql.Scan(row1)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get user data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data1, "\nData len - ", len(data1))
	if len(data1) <= 0 {
		err = errors.New("User data not found")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", data1)

	c.Data["DepartArray"] = data1

	return

}

func (c *UpdateUser) Post() {

	AdminId := c.Ctx.Input.Param(":AdminID")

	log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Admin User Page Start")
	pip := c.Ctx.Input.IP()
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Client IP - ", pip)
	var err error
	var derr error
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
			c.TplName = "user/updateUser/updateUser.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Admin UserPage Fail")
		} else if derr != nil {
			if sessErr == true {
				log.Println(beego.AppConfig.String("loglevel"), "Info", "Redirecting to login")
				c.Redirect(beego.AppConfig.String("LOGIN_PATH"), 302)
			} else {
				c.Data["DisplayMessage"] = derr.Error()
			}
			c.TplName = "user/updateUser/updateUser.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Admin UserPage Fail")

		} else {
			sess, err := session.GlobalSessions.SessionStart(c.Ctx.ResponseWriter, c.Ctx.Request)
			if err != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				err = errors.New("System is unable to process your request.Please contact customer care")
				sessErr = true
				return
			}
			if sess.Get("language") == "English" {
				c.Data["DisplayMessage"] = beego.AppConfig.String("ENGLISH_USER_UPDATED_SUCCESSFULLY")
				c.TplName = "user/updateUser/updateUser.html"
				log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Admin User  Page Success")
			} else if sess.Get("language") == "French" {
				c.Data["DisplayMessage"] = beego.AppConfig.String("FRENCH_USER_UPDATED_SUCCESSFULLY")
				c.TplName = "user/updateUser/updateUser.html"
				log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Admin User  Page Success")
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
	auth, err := utils.IsAuthorized(sess.Get("menujson").(string), "UpdateUser")
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
		c.Data["UpdateUser"] = beego.AppConfig.String("ENGLISH_UPDATE_USERS")
		c.Data["UpdateUser1"] = beego.AppConfig.String("ENGLISH_UPDATE_USERS")
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
		c.Data["ProfileManagement"] = beego.AppConfig.String("ENGLISH_PROFILEMANAGEMENT")
		c.Data["input_user_email_id"] = beego.AppConfig.String("ENGLISH_USEREMAIL")
		c.Data["please_select"] = beego.AppConfig.String("ENGLISH_PLEASESELECTCOMMON")
		c.Data["Active"] = beego.AppConfig.String("ENGLISH_ACTIVE")
		c.Data["InActive"] = beego.AppConfig.String("ENGLISH_INACTIVE")
		c.Data["English"] = beego.AppConfig.String("ENGLISH_ENGLISH")
		c.Data["French"] = beego.AppConfig.String("ENGLISH_FRENCH")
		c.Data["Update"] = beego.AppConfig.String("ENGLISH_UPDATE")
		c.Data["Users"] = beego.AppConfig.String("ENGLISH_SEARCH_USERS")

		c.TplName = "user/updateUser/updateUser.html"
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
		c.Data["UpdateUser"] = beego.AppConfig.String("FRENCH_UPDATE_USERS")
		c.Data["UpdateUser1"] = beego.AppConfig.String("FRENCH_UPDATE_USERS")
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
		c.Data["ProfileManagement"] = beego.AppConfig.String("FRENCH_PROFILEMANAGEMENT")
		c.Data["input_user_email_id"] = beego.AppConfig.String("FRENCH_USEREMAIL")
		c.Data["please_select"] = beego.AppConfig.String("FRENCH_PLEASESELECTCOMMON")
		c.Data["Active"] = beego.AppConfig.String("FRENCH_ACTIVE")
		c.Data["InActive"] = beego.AppConfig.String("FRENCH_INACTIVE")
		c.Data["English"] = beego.AppConfig.String("FRENCH_ENGLISH")
		c.Data["French"] = beego.AppConfig.String("FRENCH_FRENCH")
		c.Data["Update"] = beego.AppConfig.String("FRENCH_UPDATE")
		c.Data["Users"] = beego.AppConfig.String("FRENCH_SEARCH_USERS")

		c.TplName = "user/updateUser/updateUser.html"
	}
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

	user_contact_number := c.Input().Get("input_user_contact_number")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "user_contact_number - ", user_contact_number)

	user_role := c.Input().Get("input_user_role")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "user_role - ", user_role)

	user_department := c.Input().Get("input_user_department")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "user_department - ", user_department)

	user_employee_id := c.Input().Get("input_user_employee_id")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "user_employee_id - ", user_employee_id)

	user_status := c.Input().Get("input_user_status")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "user_status - ", user_status)

	user_language := c.Input().Get("input_user_language")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "user_language - ", user_language)

	creater_email := sess.Get("uname")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "creater_email - ", creater_email)

	if !utils.IsLetter(user_first_name) || !utils.IsLetter(user_last_name) || !utils.IsLetter(user_middele_name) || !utils.IsLetter(user_department) {
		beego.Error("suspicious special characters found")
		err = errors.New(beego.AppConfig.String(sessionLanguage + "_PARAMETER_VALIDATION"))
		Getuserdata(c, AdminId)
		return
	}

	if !utils.IsNumber(user_contact_number) || !utils.IsNumber(user_mobile) || !utils.IsNumber(user_role) {
		beego.Error("suspicious special characters found")
		err = errors.New(beego.AppConfig.String(sessionLanguage + "_PARAMETER_VALIDATION"))
		Getuserdata(c, AdminId)
		return
	}
	if utils.IsDisableCharacters(user_status) || utils.IsDisableCharacters(user_language) {
		beego.Error("suspicious special characters found")
		err = errors.New(beego.AppConfig.String(sessionLanguage + "_PARAMETER_VALIDATION"))
		Getuserdata(c, AdminId)
		return
	}
	if utils.IsDisableCharacters(user_email_id) || utils.IsDisableCharacters(user_employee_id) {
		beego.Error("suspicious special characters found")
		err = errors.New(beego.AppConfig.String(sessionLanguage + "_PARAMETER_VALIDATION"))
		Getuserdata(c, AdminId)
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

	if AdminId == beego.AppConfig.String("CNPS_DEFAULT_USER_ID") || AdminId == beego.AppConfig.String("CNPS_DEFAULT_MERCHANT_ID") {
		derr = errors.New(beego.AppConfig.String(sessionLanguage + "_DEFAULT_USER_CANNOT_UPDATE"))
	} else {

		if user_status == "ACTIVE" {
			result, err := db.Db.Exec(`UPDATE Users SET mobile=?, first_name =?,middle_name=?,last_name=?,status=?,contact_number=?,department=?,employee_id=?,language=?,last_update=now(),login_count=?,role_id=?,updated_by=? WHERE email= ?`,
				user_mobile, user_first_name, user_middele_name, user_last_name,
				user_status, user_contact_number, user_department, user_employee_id, user_language, "0", user_role, creater_email, user_email_id)
			if err != nil {
				err = errors.New("Customer updation failed")
				return
			}

			i, err := result.RowsAffected()
			if err != nil || i == 0 {
				err = errors.New("Customer updation failed")
				return
			}

		} else {
			result, err := db.Db.Exec(`UPDATE Users SET mobile=?,first_name =?,middle_name=?,last_name=?,status=?,contact_number=?,department=?,employee_id=?,language=?,last_update=now(),updated_by=? WHERE email= ?`,
				user_mobile, user_first_name, user_middele_name, user_last_name,
				user_status, user_contact_number, user_department, user_employee_id, user_language, creater_email, user_email_id)
			if err != nil {
				err = errors.New("Customer updation failed")
				return
			}

			i, err := result.RowsAffected()
			if err != nil || i == 0 {
				err = errors.New("Customer updation failed")
				return
			}
		}
	}
	Getuserdata(c, AdminId)

	// row, err := db.Db.Query(`select id,mobile,email,first_name,middle_name,last_name,role,status,contact_number,department,employee_id,language,created_date,role_id from Users where id= ?`, AdminId)
	// if err != nil {
	// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	// 	err = errors.New("Unable to get user data")
	// 	return
	// }
	// defer sql.Close(row)
	// _, data, err := sql.Scan(row)
	// if err != nil {
	// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	// 	err = errors.New("Unable to get user data")
	// 	return
	// }
	// log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	// if len(data) <= 0 {
	// 	err = errors.New("User data not found")
	// 	return
	// }

	// log.Println(beego.AppConfig.String("loglevel"), "Debug", data)

	// for i := range data {

	// 	c.Data["Id"] = data[i][0]
	// 	c.Data["UserMobile"] = data[i][1]
	// 	c.Data["UserEmail"] = data[i][2]
	// 	c.Data["UserFirstName"] = data[i][3]
	// 	c.Data["UserMiddleName"] = data[i][4]
	// 	c.Data["UserLastName"] = data[i][5]
	// 	c.Data["UserRole"] = data[i][6]
	// 	c.Data["UserStatus"] = data[i][7]
	// 	c.Data["UserContactNumber"] = data[i][8]
	// 	c.Data["UserDepartment"] = data[i][9]
	// 	c.Data["UserEmployeeID"] = data[i][10]
	// 	c.Data["UserLanguage"] = data[i][11]
	// 	c.Data["UserCreateDate"] = data[i][12]
	// 	c.Data["UserRoleId"] = data[i][13]

	// }

	// row1, err := db.Db.Query("select id,name from Roles")

	// if err != nil {
	// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	// 	err = errors.New("Unable to get user data")
	// 	return
	// }
	// defer sql.Close(row1)
	// _, data1, err := sql.Scan(row1)
	// if err != nil {
	// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	// 	err = errors.New("Unable to get user data")
	// 	return
	// }
	// log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data1, "\nData len - ", len(data1))
	// if len(data1) <= 0 {
	// 	err = errors.New("User data not found")
	// 	return
	// }

	// log.Println(beego.AppConfig.String("loglevel"), "Debug", data1)

	// c.Data["DepartArray"] = data1

	// return
}

func Getuserdata(c *UpdateUser, AdminId string) {

	row, err := db.Db.Query(`select uuid,mobile,email,first_name,middle_name,last_name,role,status,contact_number,department,employee_id,language,created_date,role_id from Users where uuid= ?`, AdminId)
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

		c.Data["Id"] = data[i][0]
		c.Data["UserMobile"] = data[i][1]
		c.Data["UserEmail"] = data[i][2]
		c.Data["UserFirstName"] = data[i][3]
		c.Data["UserMiddleName"] = data[i][4]
		c.Data["UserLastName"] = data[i][5]
		c.Data["UserRole"] = data[i][6]
		c.Data["UserStatus"] = data[i][7]
		c.Data["UserContactNumber"] = data[i][8]
		c.Data["UserDepartment"] = data[i][9]
		c.Data["UserEmployeeID"] = data[i][10]
		c.Data["UserLanguage"] = data[i][11]
		c.Data["UserCreateDate"] = data[i][12]
		c.Data["UserRoleId"] = data[i][13]

	}

	row1, err := db.Db.Query("select id,name from Roles")

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get user data")
		return
	}
	defer sql.Close(row1)
	_, data1, err := sql.Scan(row1)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get user data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data1, "\nData len - ", len(data1))
	if len(data1) <= 0 {
		err = errors.New("User data not found")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", data1)

	c.Data["DepartArray"] = data1

	return
}
