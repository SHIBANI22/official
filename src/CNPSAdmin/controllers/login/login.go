package login

import (
	"CNPSAdmin/model/db"
	"CNPSAdmin/session"
	"errors"
	"html/template"
	"runtime/debug"
	"strings"
	"time"

	"CNPSAdmin/model/utils"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"ominaya.com/database/sql"
	"ominaya.com/encoding/base64"
	"ominaya.com/util/log"
	"ominaya.com/util/pbkdf2"
)

type Login struct {
	beego.Controller
}

type LoginData struct {
	Uname    string `form:"username" valid:"Required"`
	Pass     string `form:"password" valid:"Required;MinSize(6);MaxSize(16)"`
	Language string `form:"language" valid:"Required"`
}

func (c *Login) Get() {
	defer func() {
		if l_exception := recover(); l_exception != nil {
			stack := debug.Stack()
			log.Println(beego.AppConfig.String("loglevel"), "Exception", string(stack))
			session.GlobalSessions.SessionDestroy(c.Ctx.ResponseWriter, c.Ctx.Request)
			c.Data["DisplayMessage"] = "Something went wrong.Please Contact CustomerCare."
			c.TplName = "error/error.html"
		}
		return
	}()
	utils.SetHTTPHeader(c.Ctx)

	session.GlobalSessions.SessionDestroy(c.Ctx.ResponseWriter, c.Ctx.Request)

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Client IP - ", c.Ctx.Input.IP())
	c.TplName = "login/login.html"

	c.Data["end"] = "pavan english"
	return
}

func (c *Login) Post() {

	log.Println(beego.AppConfig.String("loglevel"), "Info", "Login Start")
	pip := c.Ctx.Input.IP()
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Client IP - ", pip)

	var err error

	defer func() {

		if l_exception := recover(); l_exception != nil {
			stack := debug.Stack()
			log.Println(beego.AppConfig.String("loglevel"), "Exception", string(stack))
			session.GlobalSessions.SessionDestroy(c.Ctx.ResponseWriter, c.Ctx.Request)
			c.TplName = "error/error.html"
		}
		if err != nil {
			c.Data["DisplayMessage"] = err.Error()
			session.GlobalSessions.SessionDestroy(c.Ctx.ResponseWriter, c.Ctx.Request)
			c.TplName = "login/login.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Login Fail")
		} else {

			log.Println(beego.AppConfig.String("loglevel"), "Info", "Login Success")
		}
		return
	}()

	utils.SetHTTPHeader(c.Ctx)
	sess, err := session.GlobalSessions.SessionStart(c.Ctx.ResponseWriter, c.Ctx.Request)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("System is unable to process your request.Please contact customer care")
		return
	}

	defer func() {
		utils.EventLogs(c.Ctx, sess, c.Ctx.Input.Method(), c.Input(), c.Data, err)
		sess.SessionRelease(c.Ctx.ResponseWriter)
	}()

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Session ID - ", sess.SessionID())
	var l LoginData
	if err := c.ParseForm(&l); err != nil {
		err = errors.New("Invalid Request Received")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Form Data - ", l.Uname+l.Language)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Language - ", l.Language)
	c.Data["FormData"] = l
	valid := validation.Validation{}
	b, err := valid.Valid(&l)
	if err != nil {
		err = errors.New("Parameter validation failed")
		return
	}

	if utils.IsDisableCharacters(l.Uname) {
		beego.Error("suspicious special characters found")
		err = errors.New(beego.AppConfig.String(l.Language + "_USER_AUTHENTICATION_FAILED"))
		return
	}
	if !utils.IsLetter(l.Language) {
		beego.Error("suspicious special characters found")
		err = errors.New(beego.AppConfig.String(l.Language + "_USER_AUTHENTICATION_FAILED"))
		return
	}
	if !b {
		for _, err := range valid.Errors {
			log.Println(beego.AppConfig.String("loglevel"), "Debug", err.Key, ":", err.Message, ":", err.Field, ":", err.LimitValue, ":", err.Name, ":", err.Tmpl, ":", err.Value)
		}
		// if l.Language == "English" {
		// 	err = errors.New(beego.AppConfig.String("ENGLISH_INVALID_INPUT_ERROR_LABEL"))
		// } else {
		// 	err = errors.New(beego.AppConfig.String("FRENCH_INVALID_INPUT_ERROR_LABEL"))
		// }
		beego.Error("ENGLISH_INVALID_INPUT_ERROR_LABEL")
		err = errors.New(beego.AppConfig.String(l.Language + "_USER_AUTHENTICATION_FAILED"))
		return
	}

	err = session.CheckUserSession(l.Uname)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	}

	language, err := utils.UserLanguage(l.Uname)
	if l.Language == "English" {
		c.Data["OKLABEL"] = beego.AppConfig.String("ENGLISH_OK_BUTTON_LABEL")
	} else if l.Language == "French" {
		c.Data["OKLABEL"] = beego.AppConfig.String("FRENCH_OK_BUTTON_LABEL")
	}

	currentTime := time.Now()
	//YYY/MM/DD
	if currentTime.Format("2006-01-02") > "2020-12-31" && beego.AppConfig.String("VC") == "TRUE" {
		err = errors.New(beego.AppConfig.String("ENGLISH_EXPIRED"))
		return
	}

	status, err := utils.UserStatus(l.Uname)

	if status == "SUSPEND" {
		// if language == "English" {
		// 	err = errors.New(beego.AppConfig.String("ENGLISH_USER_SUSPEND_CONTACT_TECH_TEAM"))
		// 	return

		// } else if language == "French" {
		// 	err = errors.New(beego.AppConfig.String("FRENCH_USER_SUSPEND_CONTACT_TECH_TEAM"))
		// 	return
		// }
		beego.Error("ENGLISH_USER_SUSPEND_CONTACT_TECH_TEAM")
		err = errors.New(beego.AppConfig.String(l.Language + "_USER_AUTHENTICATION_FAILED"))

		//err = errors.New(")
		return
	} else if status == "INACTIVE" {
		// if language == "English" {
		// 	err = errors.New(beego.AppConfig.String("ENGLISH_USER_INACTIVE_ERROR_LABEL"))
		// } else {
		// 	err = errors.New(beego.AppConfig.String("FRENCH_USER_INACTIVE_ERROR_LABEL"))
		// }
		beego.Error("ENGLISH_USER_INACTIVE_ERROR_LABEL")
		err = errors.New(beego.AppConfig.String(l.Language + "_USER_AUTHENTICATION_FAILED"))

		return
	}
	count, err := utils.LoginCount(l.Uname)

	if err != nil {
		// if language == "English" {
		// 	err = errors.New("User Login Count Gets Failed")
		// 	return

		// } else if language == "French" {
		// 	err = errors.New("La tentative de login a échoué")
		// 	return
		// }
		beego.Error("User Login Count Gets Failed/User not found")
		err = errors.New(beego.AppConfig.String(l.Language + "_USER_AUTHENTICATION_FAILED"))
		return
	}

	loginCount, _ := beego.AppConfig.Int("LOGIN_COUNT")

	if count >= loginCount {

		// if language == "English" {
		// 	err = errors.New("User Login Attempt Exceeded.Please contact technical team")

		// } else if language == "French" {
		// 	err = errors.New("Les tentatives de connectiion max atteint. Veuillez contacter le support.")
		// }

		beego.Error("User Login Attempt Exceeded.Please contact technical team")
		err = errors.New(beego.AppConfig.String(l.Language + "_USER_AUTHENTICATION_FAILED"))

		result, err := db.Db.Exec(`UPDATE Users SET status=?,last_update=now() WHERE email= ?`,
			"SUSPEND", l.Uname)
		if err != nil {
			beego.Error(err)
			err = errors.New(beego.AppConfig.String(l.Language + "_USER_AUTHENTICATION_FAILED"))
			return
		}

		i, err := result.RowsAffected()
		if err != nil || i == 0 {
			beego.Error(err)
			err = errors.New(beego.AppConfig.String(l.Language + "_USER_AUTHENTICATION_FAILED"))
			return
		}

		return
	}

	name, id, role, department, language, firstlogin, fullname, rolename, menuJson, err := authinticate(l.Uname, l.Pass, l.Language)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		if language == "English" {
			err = errors.New("User Authentication Failed")

		} else if language == "French" {
			err = errors.New("L'authentification a échoué")
		}
		count++
		r_err := utils.PasswordMismatch(count, l.Uname)
		if r_err != nil {

			return
		}

		return
	}
	if count > 0 {
		err = utils.ResetLoginCount(l.Uname)
		if err != nil {
			err = errors.New("Admin User Login Count Reset Failed")
			return
		}
	}
	err = session.SetUserSession(sess.SessionID(), l.Uname)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New(beego.AppConfig.String(l.Language + "_USER_AUTHENTICATION_FAILED"))
		return
	}

	sess.Set("uname", l.Uname)
	sess.Set("name", name)
	sess.Set("uid", id)
	sess.Set("role", role)
	sess.Set("language", language)
	sess.Set("department", department)
	sess.Set("menujson", menuJson)
	sess.Set("passwordSet", firstlogin)
	sess.Set("photo", "")
	sess.Set("fullname", fullname)
	sess.Set("rolename", rolename)

	successAmount, servicecharge, totalTransCount, totalBanks, successCount, pendingCount, declainedCount, transErr := utils.SideBarTransactionData()
	if transErr != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get sidebar data")
		return
	}

	if firstlogin == "NO" {
		c.Redirect("/Common/ChangePassword", 302)
		return
	} else if role == "ADMIN" && language == "English" {

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "firstlogin - ", firstlogin)

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

		c.Redirect("/Dashboard", 302)
		return

	} else if role == "ADMIN" && language == "French" {
		sess.Set("uname", l.Uname)
		sess.Set("name", name)
		sess.Set("uid", id)
		sess.Set("role", role)
		sess.Set("language", language)
		sess.Set("department", department)
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "firstlogin - ", firstlogin)

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
		c.Redirect("/Dashboard", 302)
		return

	} else if role == "MERCHANT" && language == "English" {
		sess.Set("uname", l.Uname)
		sess.Set("name", name)
		sess.Set("uid", id)
		sess.Set("role", role)
		sess.Set("language", language)
		sess.Set("department", department)
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "firstlogin - ", firstlogin)

		c.Data["Menus"] = template.HTML(`` + beego.AppConfig.String("ENGLISH_USER_TEMPLATE") + ``)
		headerContent := strings.Replace(beego.AppConfig.String("MERCHANT_ENGLISH_HEADER_TEMPLATE"), "{{.Uname}}", sess.Get("uname").(string), -1)

		c.Data["Header"] = template.HTML(`` + headerContent + ``)

		c.Redirect("/MerchantDashboard", 302)
		return

	} else if role == "MERCHANT" && language == "French" {
		sess.Set("uname", l.Uname)
		sess.Set("name", name)
		sess.Set("uid", id)
		sess.Set("role", role)
		sess.Set("language", language)
		sess.Set("department", department)
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "firstlogin - ", firstlogin)

		c.Data["Menus1"] = template.HTML(`` + beego.AppConfig.String("FRENCH_USER_TEMPLATE") + ``)
		headerContent := strings.Replace(beego.AppConfig.String("MERCHANT_FRENCH_HEADER_TEMPLATE"), "{{.Uname}}", sess.Get("uname").(string), -1)

		c.Data["Header1"] = template.HTML(`` + headerContent + ``)
		c.Redirect("/MerchantDashboard", 302)
		return
	}

	return

}

func authinticate(uname, pass, lang string) (name, id, role, department, language, firstlogin, fullname, rolename, menuJson string, err error) {
	row, err := db.Db.Query("select password, email,Users.id ,role,department,language,password_set,role_id, CONCAT(first_name, ' ', last_name), Roles.name from Users LEFT JOIN Roles ON Roles.id=Users.role_id where email=? and status='ACTIVE' limit 1", &uname)
	authErrMsg := ""
	regErrMsg := ""
	if lang == "English" {
		authErrMsg = beego.AppConfig.String(lang + "_USER_AUTHENTICATION_FAILED")
		regErrMsg = beego.AppConfig.String(lang + "_USER_AUTHENTICATION_FAILED")
	} else {
		authErrMsg = beego.AppConfig.String(lang + "_USER_AUTHENTICATION_FAILED")
		regErrMsg = beego.AppConfig.String(lang + "_USER_AUTHENTICATION_FAILED")
	}
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New(authErrMsg)
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New(authErrMsg)
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	if len(data) <= 0 {
		log.Println(beego.AppConfig.String("loglevel"), "Error", "user not found")
		err = errors.New(regErrMsg)
		return
	}

	cp, err := base64.Decode([]byte(data[0][0]))
	if err != nil {
		err = errors.New(authErrMsg)
		return
	}

	var pbkdf pbkdf2.Pbkdf2
	pbkdf.Itr = 32
	pbkdf.KeyLen = 32
	pbkdf.Plain = []byte(pass)
	pbkdf.Salt = cp[:32]
	pbkdf.Cipher = cp[32:]
	result, err := pbkdf.Compare()

	name = data[0][1]
	id = data[0][2]
	role = data[0][3]
	department = data[0][4]
	language = data[0][5]
	firstlogin = data[0][6]
	fullname = data[0][8]
	rolename = data[0][9]

	errMsg := ""
	if language == "English" {
		errMsg = "User password incorrect"

	} else if language == "French" {
		errMsg = "User password incorrect FR"
	}
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New(errMsg)
		return
	}

	if !result {
		log.Println(beego.AppConfig.String("loglevel"), "Error", "incorrect password")
		err = errors.New(errMsg)
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "role - ", role)

	row2, err := db.Db.Query(`SELECT id,name,menus FROM Roles where id = ?`, data[0][7])

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get Role data")
		return
	}
	defer sql.Close(row2)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row2)
	_, data2, err := sql.Scan(row2)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get Role data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data2, "\nData len - ", len(data2))
	if len(data) <= 0 {
		err = errors.New("Unable to get Role data")
		return
	}

	menuJson = data2[0][2]

	return
}
