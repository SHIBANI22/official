package adminViewProfile

import (
	"CNPSAdmin/model/db"
	"CNPSAdmin/model/utils"
	"CNPSAdmin/session"

	"ominaya.com/database/sql"

	"errors"
	"html/template"
	"runtime/debug"
	"strings"

	"github.com/astaxie/beego"
	"ominaya.com/util/log"

	//	"ominaya.com/util/password"
	"os"
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
type AdminViewProfile struct {
	beego.Controller
}

func (c *AdminViewProfile) Get() {
	beego.ReadFromRequest(&c.Controller)
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Admin View Profile Page Start")
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
			c.TplName = "adminviewprofile/adminViewProfile.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Admin View Profile Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "adminviewprofile/adminViewProfile.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Admin View Profile Page Success")
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
	passSet := sess.Get("passwordSet").(string)
	if passSet != "YES" {
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
		c.Data["AdminViewProfile"] = beego.AppConfig.String("ENGLISH_ADMIN_VIEW_PROFILE")
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
		c.Data["UserInformation"] = beego.AppConfig.String("ENGLISH_USER_INFORMATION")

		c.Data["ProfileTitleLabel"] = beego.AppConfig.String("ENGLISH_PROFILE_TITLE_LABEL")
		c.Data["ChoosePlaceholderLabel"] = beego.AppConfig.String("ENGLISH_CHOOSE_FILE_LABEL")
		c.Data["ProfileChooseLabel"] = beego.AppConfig.String("ENGLISH_PROFILE_LABEL")
		c.Data["ProfileUpdateTooltipLabel"] = beego.AppConfig.String("ENGLISH_UPDATE_PROFILE_LABEL")
		c.Data["ProfileBrowseLabel"] = beego.AppConfig.String("ENGLISH_BROWSE_LABEL")
		c.Data["BackTooltipLabel"] = beego.AppConfig.String("ENGLISH_BACK_TO_PROFILE_LABEL")

		c.TplName = "adminviewprofile/adminViewProfile.html"
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
		c.Data["AdminViewProfile"] = beego.AppConfig.String("FRENCH_ADMIN_VIEW_PROFILE")
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
		c.Data["UserInformation"] = beego.AppConfig.String("FRENCH_USER_INFORMATION")

		c.Data["ProfileTitleLabel"] = beego.AppConfig.String("FRENCH_PROFILE_TITLE_LABEL")
		c.Data["ChoosePlaceholderLabel"] = beego.AppConfig.String("FRENCH_CHOOSE_FILE_LABEL")
		c.Data["ProfileChooseLabel"] = beego.AppConfig.String("FRENCH_PROFILE_LABEL")
		c.Data["ProfileUpdateTooltipLabel"] = beego.AppConfig.String("FRENCH_UPDATE_PROFILE_LABEL")
		c.Data["ProfileBrowseLabel"] = beego.AppConfig.String("FRENCH_BROWSE_LABEL")
		c.Data["BackTooltipLabel"] = beego.AppConfig.String("FRENCH_BACK_TO_PROFILE_LABEL")

		c.TplName = "adminviewprofile/adminViewProfile.html"
	}

	uname := sess.Get("uname")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "User Name - ", uname)

	row, err := db.Db.Query(`select Users.id,mobile,email,first_name,middle_name,last_name,Roles.name,status,contact_number,department,employee_id,language,Users.created_date from Users LEFT JOIN Roles ON Roles.id=Users.role_id where email= ?`, uname)
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
		//c.Data["Photo"] = data[i][13]

	}

	return

}

func (c *AdminViewProfile) Post() {
	flash := beego.NewFlash()
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Admin View Profile Page Start")
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
		log.Println(beego.AppConfig.String("loglevel"), "Info", "Redirecting", err)
		if err != nil {
			if sessErr == true {
				log.Println(beego.AppConfig.String("loglevel"), "Info", "Redirecting to login")
				c.Redirect(beego.AppConfig.String("LOGIN_PATH"), 302)

			} else {
				// c.Data["DisplayMessage"] = err.Error()
				flash.Error(err.Error())
				flash.Store(&c.Controller)
				c.Redirect(c.URLFor("AdminViewProfile.Get"), 302)
			}
			c.TplName = "adminviewprofile/adminViewProfile.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Admin View Profile Page Fail")
		} else {
			// c.Data["DisplayMessage"] = "User information updated successfully!"
			utils.SetHTTPHeader(c.Ctx)
			sess, err := session.GlobalSessions.SessionStart(c.Ctx.ResponseWriter, c.Ctx.Request)
			if err != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				err = errors.New("System is unable to process your request.Please contact customer care")
				sessErr = true
				return
			}

			if sess.Get("language") == "English" {
				flash.Success("User information updated successfully!")
			} else {
				flash.Success("Utilisateur mis à jour avec succès")
			}
			flash.Store(&c.Controller)
			c.Redirect(c.URLFor("AdminViewProfile.Get"), 302)
			c.TplName = "adminviewprofile/adminViewProfile.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Admin View Profile Page Success")
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
	passSet := sess.Get("passwordSet").(string)
	if passSet != "YES" {
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
		c.Data["AdminViewProfile"] = beego.AppConfig.String("ENGLISH_ADMIN_VIEW_PROFILE")
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
		c.Data["UserInformation"] = beego.AppConfig.String("ENGLISH_USER_INFORMATION")

		c.Data["ProfileTitleLabel"] = beego.AppConfig.String("ENGLISH_PROFILE_TITLE_LABEL")
		c.Data["ChoosePlaceholderLabel"] = beego.AppConfig.String("ENGLISH_CHOOSE_FILE_LABEL")
		c.Data["ProfileChooseLabel"] = beego.AppConfig.String("ENGLISH_PROFILE_LABEL")
		c.Data["ProfileUpdateTooltipLabel"] = beego.AppConfig.String("ENGLISH_UPDATE_PROFILE_LABEL")
		c.Data["ProfileBrowseLabel"] = beego.AppConfig.String("ENGLISH_BROWSE_LABEL")
		c.Data["BackTooltipLabel"] = beego.AppConfig.String("ENGLISH_BACK_TO_PROFILE_LABEL")

		c.TplName = "adminviewprofile/adminViewProfile.html"
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
		c.Data["AdminViewProfile"] = beego.AppConfig.String("FRENCH_ADMIN_VIEW_PROFILE")
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
		c.Data["UserInformation"] = beego.AppConfig.String("FRENCH_USER_INFORMATION")

		c.Data["ProfileTitleLabel"] = beego.AppConfig.String("FRENCH_PROFILE_TITLE_LABEL")
		c.Data["ChoosePlaceholderLabel"] = beego.AppConfig.String("FRENCH_CHOOSE_FILE_LABEL")
		c.Data["ProfileChooseLabel"] = beego.AppConfig.String("FRENCH_PROFILE_LABEL")
		c.Data["ProfileUpdateTooltipLabel"] = beego.AppConfig.String("FRENCH_UPDATE_PROFILE_LABEL")
		c.Data["ProfileBrowseLabel"] = beego.AppConfig.String("FRENCH_BROWSE_LABEL")
		c.Data["BackTooltipLabel"] = beego.AppConfig.String("FRENCH_BACK_TO_PROFILE_LABEL")

		c.TplName = "adminviewprofile/adminViewProfile.html"
	}

	uname := sess.Get("uname")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "User Name - ", uname)

	user_language := c.Input().Get("input_user_language")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "user_language - ", user_language)

	if utils.IsDisableCharacters(user_language) {
		beego.Error("suspicious special characters found")
		err = errors.New(beego.AppConfig.String(sessionLanguage + "_PARAMETER_VALIDATION"))
		return
	}

	if user_language != "English" && user_language != "French" {
		beego.Error("suspicious special characters found")
		err = errors.New(beego.AppConfig.String(sessionLanguage + "_PARAMETER_VALIDATION"))
		return
	}

	/*files, _ := c.GetFiles("profile_photo")

	if len(files) > 0 {
		f, h, err := c.GetFile("profile_photo")

		if err != nil {
			err = errors.New("Please choose file to upload")
			log.Println(beego.AppConfig.String("loglevel"), "Error", "Please choose file to upload")
			return
		}
		defer f.Close()
		uploadDir := beego.AppConfig.String("UPLOAD_DIR")
		_, errFol := os.Stat(uploadDir)
		if os.IsNotExist(errFol) {
			if errDir := os.MkdirAll(uploadDir, 0755); errDir != nil {
				err = errors.New("Unable to upload photo")
				flash.Error(err.Error())
				flash.Store(&c.Controller)
				c.Redirect(c.URLFor("AdminViewProfile.Get"), 302)
				return
			}
		}
		path := uploadDir + h.Filename

		if h.Header.Get("Content-Type") != "image/jpeg" && h.Header.Get("Content-Type") != "image/png" && h.Header.Get("Content-Type") != "image/gif" && h.Header.Get("Content-Type") != "image/jpg" {
			if sessionLanguage == "English" {
				err = errors.New("Please upload only jpg, png and gif format")
			} else {
				err = errors.New("Veuillez charger les formats Png, Jpeg and Gif format")
			}
			flash.Error(err.Error())
			flash.Store(&c.Controller)
			c.Redirect(c.URLFor("AdminViewProfile.Get"), 302)
			return
		} else if h.Size > 2000000 {
			if sessionLanguage == "English" {
				err = errors.New("Please upload a file less than 2 Mb")
			} else {
				err = errors.New("la taille de la photo doit être inférieure à 2Mb")
			}
			log.Println(beego.AppConfig.String("loglevel"), "Error", "Please upload a file less than 2 Mb")
			flash.Error(err.Error())
			flash.Store(&c.Controller)
			c.Redirect(c.URLFor("AdminViewProfile.Get"), 302)
			return
		}
		if err := c.SaveToFile("profile_photo", path); err != nil {
			if sessionLanguage == "English" {
				err = errors.New("Profile Photo not uploaded.Please try again!")
			} else {
				err = errors.New("Utilisateur non mis à jour. Veuillez réessayer svp!")
			}
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			return
		}
		res, err := db.Db.Exec(`UPDATE Users SET last_update=now(), photo=? WHERE email= ?`, h.Filename, uname)
		if err != nil {
			if sessionLanguage == "English" {
				err = errors.New("Profile Photo not uploaded.Please try again!")
			} else {
				err = errors.New("Utilisateur non mis à jour. Veuillez réessayer svp!")
			}
			return
		}

		i, err := res.RowsAffected()
		if err != nil || i == 0 {
			if sessionLanguage == "English" {
				err = errors.New("Profile Photo not uploaded.Please try again!")
			} else {
				err = errors.New("Utilisateur non mis à jour. Veuillez réessayer svp!")
			}
			return
		}
		sess.Set("photo", h.Filename)
	}*/
	defer func() {
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
	}()
	// c.Redirect("/Dashboard", 302)
	return
}

func (c *AdminViewProfile) RemoveProfilePhoto() {
	flash := beego.NewFlash()
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Admin View Profile Page Start")
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
				// c.Data["DisplayMessage"] = err.Error()
				flash.Error(err.Error())
				flash.Store(&c.Controller)
				c.Redirect(c.URLFor("AdminViewProfile.Get"), 302)
			}
			c.TplName = "adminviewprofile/adminViewProfile.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Admin View Profile Page Fail")
		} else {
			// c.Data["DisplayMessage"] = "Profile photo removed successfully!"
			utils.SetHTTPHeader(c.Ctx)
			sess, err := session.GlobalSessions.SessionStart(c.Ctx.ResponseWriter, c.Ctx.Request)
			if err != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				err = errors.New("System is unable to process your request.Please contact customer care")
				sessErr = true
				return
			}

			if sess.Get("language") == "English" {
				flash.Success("Profile photo removed successfully!")
			} else {
				flash.Success("Photo de profil retirée avec succès")
			}
			flash.Store(&c.Controller)
			c.Redirect(c.URLFor("AdminViewProfile.Get"), 302)
			c.TplName = "adminviewprofile/adminViewProfile.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Admin View Profile Page Success")
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
	passSet := sess.Get("passwordSet").(string)
	if passSet != "YES" {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	defer func() {
		utils.EventLogs(c.Ctx, sess, c.Ctx.Input.Method(), c.Input(), c.Data, err)
		sess.SessionRelease(c.Ctx.ResponseWriter)
	}()
	c.Data["language"] = sess.Get("language").(string)
	// sessionLanguage := sess.Get("language").(string)

	uname := sess.Get("uname")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "User Name - ", uname)

	res, err := db.Db.Exec(`UPDATE Users SET last_update=now(), photo=NULL WHERE email= ?`, uname)
	if err != nil {
		if sess.Get("language") == "English" {
			err = errors.New("Profile Photo not removed.Please try again!")
		} else {
			err = errors.New("Profile Photo not removed.Please try again!")
		}
		return
	}

	i, err := res.RowsAffected()
	if err != nil || i == 0 {
		err = errors.New("Profile Photo not removed.Please try again!")
		return
	}
	fileName := sess.Get("photo").(string)
	fileFullPath := beego.AppConfig.String("UPLOAD_DIR") + fileName

	if err := os.Remove(fileFullPath); err != nil {
		err = errors.New("Profile Photo not removed.Please try again!")
		return
	}
	sess.Set("photo", "")
	return
}
