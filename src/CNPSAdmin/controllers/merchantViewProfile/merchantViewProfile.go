package merchantViewProfile

import (
	"CNPSAdmin/model/db"
	"CNPSAdmin/model/utils"
	"CNPSAdmin/session"
	"strings"

	"ominaya.com/database/sql"

	"errors"
	"html/template"

	"runtime/debug"

	"github.com/astaxie/beego"
	"ominaya.com/util/log"
	//	"ominaya.com/util/password"
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
type MerchantViewProfile struct {
	beego.Controller
}

func (c *MerchantViewProfile) Get() {

	log.Println(beego.AppConfig.String("loglevel"), "Info", "Merchnat View Proflie  Start")
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
			c.TplName = "merchantviewprofile/merchantViewProfile.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Merchnat View Proflie Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "merchantviewprofile/merchantViewProfile.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Merchnat View Proflie  Page Success")
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
	c.Data["language"] = sess.Get("language").(string)
	if sess.Get("role") == "MERCHANT" && sess.Get("language") == "English" {
		c.Data["Menus1"] = template.HTML(`` + beego.AppConfig.String("ENGLISH_USER_TEMPLATE") + ``)
		headerContent := strings.Replace(beego.AppConfig.String("MERCHANT_ENGLISH_HEADER_TEMPLATE"), "{{.Uname}}", sess.Get("uname").(string), -1)

		c.Data["Header1"] = template.HTML(`` + headerContent + ``)

		c.Data["Dashboard"] = beego.AppConfig.String("ENGLISH_DASHBOARD")
		c.Data["MerchantViewProfile"] = beego.AppConfig.String("ENGLISH_MERCHANT_VIEW_PROFILE")
		c.Data["MerchantViewProfile1"] = beego.AppConfig.String("ENGLISH_MERCHANT_VIEW_PROFILE")
		c.Data["GeneralInformation"] = beego.AppConfig.String("ENGLISH_GENERAL_INFORMATION")
		c.Data["ID"] = beego.AppConfig.String("ENGLISH_ID")
		c.Data["Mobile"] = beego.AppConfig.String("ENGLISH_MOBILE")
		c.Data["Email"] = beego.AppConfig.String("ENGLISH_EMAIL")
		c.Data["FirstName"] = beego.AppConfig.String("ENGLISH_FIRSTNAME")
		c.Data["MiddleName"] = beego.AppConfig.String("ENGLISH_MIDDLENAME")
		c.Data["LastName"] = beego.AppConfig.String("ENGLISH_LASTNAME")
		c.Data["Role"] = beego.AppConfig.String("ENGLISH_ROLE")

		c.Data["ContactNumber"] = beego.AppConfig.String("ENGLISH_CONTATC_NUMBER")
		c.Data["DepartmentName"] = beego.AppConfig.String("ENGLISH_DEPT_NAME")
		c.Data["EmployeeId"] = beego.AppConfig.String("ENGLISH_EMPLOYEE_ID")

		c.Data["Status"] = beego.AppConfig.String("ENGLISH_STATUS")
		c.Data["Language"] = beego.AppConfig.String("ENGLISH_LANGUAGE")
		c.Data["updateUser"] = beego.AppConfig.String("ENGLISH_UPDATEUSER")
		c.Data["viewuser"] = beego.AppConfig.String("ENGLISH_VIEW_USER")
		c.Data["cancel"] = beego.AppConfig.String("ENGLISH_CANCEL")
		c.Data["pleaseselect"] = beego.AppConfig.String("ENGLISH_ENTER_PLEASESELECT")
		c.Data["ViewUser"] = beego.AppConfig.String("ENGLISH_VIEW_USERS")

		c.TplName = "merchantviewprofile/merchantViewProfile.html"
	} else if sess.Get("role") == "MERCHANT" && sess.Get("language") == "French" {
		c.Data["Menus1"] = template.HTML(`` + beego.AppConfig.String("FRENCH_USER_TEMPLATE") + ``)
		headerContent := strings.Replace(beego.AppConfig.String("MERCHANT_FRENCH_HEADER_TEMPLATE"), "{{.Uname}}", sess.Get("uname").(string), -1)

		c.Data["Header1"] = template.HTML(`` + headerContent + ``)

		c.Data["Dashboard"] = beego.AppConfig.String("FRENCH_DASHBOARD")
		c.Data["MerchantViewProfile"] = beego.AppConfig.String("FRENCH_MERCHANT_VIEW_PROFILE")
		c.Data["MerchantViewProfile1"] = beego.AppConfig.String("FRENCH_MERCHANT_VIEW_PROFILE")
		c.Data["GeneralInformation"] = beego.AppConfig.String("FRENCH_GENERAL_INFORMATION")
		c.Data["ID"] = beego.AppConfig.String("FRENCH_ID")
		c.Data["Mobile"] = beego.AppConfig.String("FRENCH_MOBILE")
		c.Data["Email"] = beego.AppConfig.String("FRENCH_EMAIL")
		c.Data["FirstName"] = beego.AppConfig.String("FRENCH_FIRSTNAME")
		c.Data["MiddleName"] = beego.AppConfig.String("FRENCH_MIDDLENAME")
		c.Data["LastName"] = beego.AppConfig.String("FRENCH_LASTNAME")
		c.Data["Role"] = beego.AppConfig.String("FRENCH_ROLE")

		c.Data["ContactNumber"] = beego.AppConfig.String("FRENCH_CONTATC_NUMBER")
		c.Data["DepartmentName"] = beego.AppConfig.String("FRENCH_DEPT_NAME")
		c.Data["EmployeeId"] = beego.AppConfig.String("FRENCH_EMPLOYEE_ID")

		c.Data["Status"] = beego.AppConfig.String("FRENCH_STATUS")
		c.Data["Language"] = beego.AppConfig.String("FRENCH_LANGUAGE")
		c.Data["updateUser"] = beego.AppConfig.String("FRENCH_UPDATEUSER")
		c.Data["viewuser"] = beego.AppConfig.String("FRENCH_VIEW_USER")
		c.Data["cancel"] = beego.AppConfig.String("FRENCH_CANCEL")
		c.Data["pleaseselect"] = beego.AppConfig.String("FRENCH_ENTER_PLEASESELECT")
		c.Data["ViewUser"] = beego.AppConfig.String("FRENCH_VIEW_USERS")

		c.TplName = "merchantviewprofile/merchantViewProfile.html"
	}

	uname := sess.Get("uname")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "User Name - ", uname)

	row, err := db.Db.Query(`select id,mobile,email,first_name,middle_name,last_name,role,status,contact_number,department,employee_id,language,created_date from Users where email= ?`, uname)
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

	}

	return

}

func (c *MerchantViewProfile) Post() {

	log.Println(beego.AppConfig.String("loglevel"), "Info", "Admin View Profile Page Start")
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
			c.TplName = "adminviewprofile/adminViewProfile.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Customer Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "adminviewprofile/adminViewProfile.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Customer  Page Success")
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
	c.Data["language"] = sess.Get("language").(string)
	if sess.Get("role") == "MERCHANT" && sess.Get("language") == "English" {
		c.Data["Menus1"] = template.HTML(`` + beego.AppConfig.String("ENGLISH_USER_TEMPLATE") + ``)
		headerContent := strings.Replace(beego.AppConfig.String("MERCHANT_ENGLISH_HEADER_TEMPLATE"), "{{.Uname}}", sess.Get("uname").(string), -1)

		c.Data["Header1"] = template.HTML(`` + headerContent + ``)

		c.Data["Dashboard"] = beego.AppConfig.String("ENGLISH_DASHBOARD")
		c.Data["MerchantViewProfile"] = beego.AppConfig.String("ENGLISH_MERCHANT_VIEW_PROFILE")
		c.Data["MerchantViewProfile1"] = beego.AppConfig.String("ENGLISH_MERCHANT_VIEW_PROFILE")
		c.Data["GeneralInformation"] = beego.AppConfig.String("ENGLISH_GENERAL_INFORMATION")
		c.Data["ID"] = beego.AppConfig.String("ENGLISH_ID")
		c.Data["Mobile"] = beego.AppConfig.String("ENGLISH_MOBILE")
		c.Data["Email"] = beego.AppConfig.String("ENGLISH_EMAIL")
		c.Data["FirstName"] = beego.AppConfig.String("ENGLISH_FIRSTNAME")
		c.Data["MiddleName"] = beego.AppConfig.String("ENGLISH_MIDDLENAME")
		c.Data["LastName"] = beego.AppConfig.String("ENGLISH_LASTNAME")
		c.Data["Role"] = beego.AppConfig.String("ENGLISH_ROLE")

		c.Data["ContactNumber"] = beego.AppConfig.String("ENGLISH_CONTATC_NUMBER")
		c.Data["DepartmentName"] = beego.AppConfig.String("ENGLISH_DEPT_NAME")
		c.Data["EmployeeId"] = beego.AppConfig.String("ENGLISH_EMPLOYEE_ID")

		c.Data["Status"] = beego.AppConfig.String("ENGLISH_STATUS")
		c.Data["Language"] = beego.AppConfig.String("ENGLISH_LANGUAGE")
		c.Data["updateUser"] = beego.AppConfig.String("ENGLISH_UPDATEUSER")
		c.Data["viewuser"] = beego.AppConfig.String("ENGLISH_VIEW_USER")
		c.Data["cancel"] = beego.AppConfig.String("ENGLISH_CANCEL")
		c.Data["pleaseselect"] = beego.AppConfig.String("ENGLISH_ENTER_PLEASESELECT")

		c.TplName = "merchantviewprofile/merchantViewProfile.html"
	} else if sess.Get("role") == "MERCHANT" && sess.Get("language") == "French" {
		c.Data["Menus1"] = template.HTML(`` + beego.AppConfig.String("FRENCH_USER_TEMPLATE") + ``)
		headerContent := strings.Replace(beego.AppConfig.String("MERCHANT_FRENCH_HEADER_TEMPLATE"), "{{.Uname}}", sess.Get("uname").(string), -1)

		c.Data["Header1"] = template.HTML(`` + headerContent + ``)

		c.Data["Dashboard"] = beego.AppConfig.String("FRENCH_DASHBOARD")
		c.Data["MerchantViewProfile"] = beego.AppConfig.String("FRENCH_MERCHANT_VIEW_PROFILE")
		c.Data["MerchantViewProfile1"] = beego.AppConfig.String("FRENCH_MERCHANT_VIEW_PROFILE")
		c.Data["GeneralInformation"] = beego.AppConfig.String("FRENCH_GENERAL_INFORMATION")
		c.Data["ID"] = beego.AppConfig.String("FRENCH_ID")
		c.Data["Mobile"] = beego.AppConfig.String("FRENCH_MOBILE")
		c.Data["Email"] = beego.AppConfig.String("FRENCH_EMAIL")
		c.Data["FirstName"] = beego.AppConfig.String("FRENCH_FIRSTNAME")
		c.Data["MiddleName"] = beego.AppConfig.String("FRENCH_MIDDLENAME")
		c.Data["LastName"] = beego.AppConfig.String("FRENCH_LASTNAME")
		c.Data["Role"] = beego.AppConfig.String("FRENCH_ROLE")

		c.Data["ContactNumber"] = beego.AppConfig.String("FRENCH_CONTATC_NUMBER")
		c.Data["DepartmentName"] = beego.AppConfig.String("FRENCH_DEPT_NAME")
		c.Data["EmployeeId"] = beego.AppConfig.String("FRENCH_EMPLOYEE_ID")

		c.Data["Status"] = beego.AppConfig.String("FRENCH_STATUS")
		c.Data["Language"] = beego.AppConfig.String("FRENCH_LANGUAGE")
		c.Data["updateUser"] = beego.AppConfig.String("FRENCH_UPDATEUSER")
		c.Data["viewuser"] = beego.AppConfig.String("FRENCH_VIEW_USER")
		c.Data["cancel"] = beego.AppConfig.String("FRENCH_CANCEL")
		c.Data["pleaseselect"] = beego.AppConfig.String("FRENCH_ENTER_PLEASESELECT")

		c.TplName = "merchantviewprofile/merchantViewProfile.html"
	}

	uname := sess.Get("uname")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "User Name - ", uname)

	user_language := c.Input().Get("input_user_language")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "user_language - ", user_language)

	result, err := db.Db.Exec(`UPDATE Users SET language=?,last_update=now() WHERE email= ?`,
		user_language, uname)
	if err != nil {
		err = errors.New("Customer updation failed")
		return
	}

	i, err := result.RowsAffected()
	if err != nil || i == 0 {
		err = errors.New("Customer updation failed")
		return
	}

	sess.Set("language", user_language)
	c.Redirect("/MerchantDashboard", 302)

	return

}
