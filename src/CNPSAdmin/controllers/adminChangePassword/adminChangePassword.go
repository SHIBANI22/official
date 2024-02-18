package adminChangePassword

import (
	"CNPSAdmin/model/db"
	"CNPSAdmin/model/utils"
	"CNPSAdmin/session"
	"crypto/rand"
	"unicode"

	"ominaya.com/database/sql"

	"ominaya.com/encoding/base64"
	"ominaya.com/util/pbkdf2"

	"errors"
	"html/template"
	"runtime/debug"
	"strings"

	"github.com/astaxie/beego"
	"ominaya.com/util/log"
	//	"ominaya.com/util/password"
)

type AdminChangePassword struct {
	beego.Controller
}

func (c *AdminChangePassword) Get() {
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Admin Change Password Start")
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
			c.TplName = "adminchangePassword/adminchangePassword.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Admin Change Password Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "adminchangePassword/adminchangePassword.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Admin Change Password Page Success")
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
		c.Data["ChangePassword"] = beego.AppConfig.String("ENGLISH_CHANGE_PASSWORD")
		c.Data["OldPassword"] = beego.AppConfig.String("ENGLISH_OLD_PASSWORD")
		c.Data["NewPassword"] = beego.AppConfig.String("ENGLISH_NEW_PASSWORD")
		c.Data["ConfirmPassword"] = beego.AppConfig.String("ENGLISH_CONFIRM_PASSWORD")
		c.Data["input_old_password"] = beego.AppConfig.String("ENGLISH_ENTER_OLD_PASSWORD")
		c.Data["input_new_password"] = beego.AppConfig.String("ENGLISH_ENTER_NEW_PASSWORD")
		c.Data["input_confirm_password"] = beego.AppConfig.String("ENGLISH_ENTER_CONFIRM_PASSWORD")
		c.Data["msg"] = beego.AppConfig.String("ENGLISH_PASSWORD_MUST_CONTAIN_THE_FOLLOWING")
		c.Data["letter"] = beego.AppConfig.String("ENGLISH_A_LOWERCASE_LETTER")
		c.Data["capital"] = beego.AppConfig.String("ENGLISH_A_CAPITAL_LETTER")
		c.Data["number"] = beego.AppConfig.String("ENGLISH_A_NUMBER")
		c.Data["length"] = beego.AppConfig.String("ENGLISH_MINIMUM_6_CHARACTERS")
		c.Data["specialchar"] = beego.AppConfig.String("ENGLISH_A_SPECIAL_CHARACTER_LETTER")
		c.Data["Submit"] = beego.AppConfig.String("ENGLISH_VALIDATE")
		c.Data["filterDataReset"] = beego.AppConfig.String("ENGLISH_BACK_BUTTON")

		// headerContent := sess.Get("Header").(string)
		// c.Data["Header"] = template.HTML(`` + headerContent + ``)

		c.TplName = "adminchangePassword/adminchangePassword.html"
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
		c.Data["ChangePassword"] = beego.AppConfig.String("FRENCH_CHANGE_PASSWORD")

		c.Data["OldPassword"] = beego.AppConfig.String("FRENCH_OLD_PASSWORD")
		c.Data["NewPassword"] = beego.AppConfig.String("FRENCH_NEW_PASSWORD")
		c.Data["ConfirmPassword"] = beego.AppConfig.String("FRENCH_CONFIRM_PASSWORD")
		c.Data["input_old_password"] = beego.AppConfig.String("FRENCH_ENTER_OLD_PASSWORD")
		c.Data["input_new_password"] = beego.AppConfig.String("FRENCH_ENTER_NEW_PASSWORD")
		c.Data["input_confirm_password"] = beego.AppConfig.String("FRENCH_ENTER_CONFIRM_PASSWORD")
		c.Data["msg"] = beego.AppConfig.String("FRENCH_PASSWORD_MUST_CONTAIN_THE_FOLLOWING")
		c.Data["letter"] = beego.AppConfig.String("FRENCH_A_LOWERCASE_LETTER")
		c.Data["capital"] = beego.AppConfig.String("FRENCH_A_CAPITAL_LETTER")
		c.Data["number"] = beego.AppConfig.String("FRENCH_A_NUMBER")
		c.Data["length"] = beego.AppConfig.String("FRENCH_MINIMUM_6_CHARACTERS")
		c.Data["specialchar"] = beego.AppConfig.String("FRENCH_A_SPECIAL_CHARACTER_LETTER")
		c.Data["Submit"] = beego.AppConfig.String("FRENCH_VALIDATE")
		c.Data["filterDataReset"] = beego.AppConfig.String("FRENCH_BACK_BUTTON")

		// headerContent := sess.Get("Header").(string)
		// c.Data["Header"] = template.HTML(`` + headerContent + ``)

		c.TplName = "adminchangePassword/adminchangePassword.html"
	}

	return
}

func (c *AdminChangePassword) Post() {

	log.Println(beego.AppConfig.String("loglevel"), "Info", "Admin Change Password Start")
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
			c.TplName = "adminchangePassword/adminchangePassword.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Admin Chnage Password Page Fail")
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
				c.Data["DisplayMessage"] = beego.AppConfig.String("ENGLISH_CHANGE_PASSWORD_SUCESSFULLY")
				c.TplName = "adminchangePassword/adminchangePassword.html"
				log.Println(beego.AppConfig.String("loglevel"), "Info", "Admin Chnage Password  Page Success")
			} else if sess.Get("language") == "French" {
				c.Data["DisplayMessage"] = beego.AppConfig.String("FRENCH_CHANGE_PASSWORD_SUCESSFULLY")
				c.TplName = "adminchangePassword/adminchangePassword.html"
				log.Println(beego.AppConfig.String("loglevel"), "Info", "Admin Chnage Password  Page Success")

			}
			c.Data["LogoutCurrentUser"] = true
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
		c.Data["ChangePassword"] = beego.AppConfig.String("ENGLISH_CHANGE_PASSWORD")
		c.Data["OldPassword"] = beego.AppConfig.String("ENGLISH_OLD_PASSWORD")
		c.Data["NewPassword"] = beego.AppConfig.String("ENGLISH_NEW_PASSWORD")
		c.Data["ConfirmPassword"] = beego.AppConfig.String("ENGLISH_CONFIRM_PASSWORD")
		c.Data["input_old_password"] = beego.AppConfig.String("ENGLISH_ENTER_OLD_PASSWORD")
		c.Data["input_new_password"] = beego.AppConfig.String("ENGLISH_ENTER_NEW_PASSWORD")
		c.Data["input_confirm_password"] = beego.AppConfig.String("ENGLISH_ENTER_CONFIRM_PASSWORD")
		c.Data["msg"] = beego.AppConfig.String("ENGLISH_PASSWORD_MUST_CONTAIN_THE_FOLLOWING")
		c.Data["letter"] = beego.AppConfig.String("ENGLISH_A_LOWERCASE_LETTER")
		c.Data["capital"] = beego.AppConfig.String("ENGLISH_A_CAPITAL_LETTER")
		c.Data["number"] = beego.AppConfig.String("ENGLISH_A_NUMBER")
		c.Data["length"] = beego.AppConfig.String("ENGLISH_MINIMUM_6_CHARACTERS")
		c.Data["specialchar"] = beego.AppConfig.String("ENGLISH_A_SPECIAL_CHARACTER_LETTER")
		c.Data["Submit"] = beego.AppConfig.String("ENGLISH_VALIDATE")
		c.Data["filterDataReset"] = beego.AppConfig.String("ENGLISH_BACK_BUTTON")

		// headerContent := sess.Get("Header").(string)
		// c.Data["Header"] = template.HTML(`` + headerContent + ``)

		c.TplName = "adminchangePassword/adminchangePassword.html"
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
		c.Data["ChangePassword"] = beego.AppConfig.String("FRENCH_CHANGE_PASSWORD")

		c.Data["OldPassword"] = beego.AppConfig.String("FRENCH_OLD_PASSWORD")
		c.Data["NewPassword"] = beego.AppConfig.String("FRENCH_NEW_PASSWORD")
		c.Data["ConfirmPassword"] = beego.AppConfig.String("FRENCH_CONFIRM_PASSWORD")
		c.Data["input_old_password"] = beego.AppConfig.String("FRENCH_ENTER_OLD_PASSWORD")
		c.Data["input_new_password"] = beego.AppConfig.String("FRENCH_ENTER_NEW_PASSWORD")
		c.Data["input_confirm_password"] = beego.AppConfig.String("FRENCH_ENTER_CONFIRM_PASSWORD")
		c.Data["msg"] = beego.AppConfig.String("FRENCH_PASSWORD_MUST_CONTAIN_THE_FOLLOWING")
		c.Data["letter"] = beego.AppConfig.String("FRENCH_A_LOWERCASE_LETTER")
		c.Data["capital"] = beego.AppConfig.String("FRENCH_A_CAPITAL_LETTER")
		c.Data["number"] = beego.AppConfig.String("FRENCH_A_NUMBER")
		c.Data["length"] = beego.AppConfig.String("FRENCH_MINIMUM_6_CHARACTERS")
		c.Data["specialchar"] = beego.AppConfig.String("FRENCH_A_SPECIAL_CHARACTER_LETTER")
		c.Data["Submit"] = beego.AppConfig.String("FRENCH_VALIDATE")
		c.Data["filterDataReset"] = beego.AppConfig.String("FRENCH_BACK_BUTTON")

		// headerContent := sess.Get("Header").(string)
		// c.Data["Header"] = template.HTML(`` + headerContent + ``)

		c.TplName = "adminchangePassword/adminchangePassword.html"
	}

	old_password := c.Input().Get("input_old_password")
	//log.Println(beego.AppConfig.String("loglevel"), "Debug", "old_password - ", old_password)

	new_password := c.Input().Get("input_new_password")
	//log.Println(beego.AppConfig.String("loglevel"), "Debug", "new_password - ", new_password)

	confirm_password := c.Input().Get("input_confirm_password")
	//log.Println(beego.AppConfig.String("loglevel"), "Debug", "confirm_password - ", confirm_password)

	uname := sess.Get("uname")
	language := sess.Get("language").(string)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "User Name - ", uname)

	if old_password == "" || new_password == "" || confirm_password == "" {
		log.Println(beego.AppConfig.String("loglevel"), "Error", "Blank Passwords")
		err = errors.New("Passwords can't be blank")
		return
	}

	if len(old_password) > 16 || len(new_password) > 16 || len(confirm_password) > 16 {
		beego.Error("suspicious special characters found")
		err = errors.New(beego.AppConfig.String(language + "_PARAMETER_VALIDATION"))
		return
	}

	if !isValid(new_password) {
		beego.Error("password patter not match")
		err = errors.New(beego.AppConfig.String(language + "_PARAMETER_VALIDATION"))
		return
	}
	if !isValid(confirm_password) {
		beego.Error("password patter not match")
		err = errors.New(beego.AppConfig.String(language + "_PARAMETER_VALIDATION"))
		return
	}

	if new_password != confirm_password {
		log.Println(beego.AppConfig.String("loglevel"), "Error", "New Password Mismatch")
		err = errors.New("New passwords can't be different")
		return
	}

	err = Authenticate(uname.(string), old_password, "ADMIN", language)

	if err != nil {

		if sess.Get("language") == "English" {
			err = errors.New(beego.AppConfig.String("ENGLISH_ADMIN_AUTHENTICATION_FAILED"))
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			return
		}
		if sess.Get("language") == "French" {
			err = errors.New(beego.AppConfig.String("FRENCH_ADMIN_AUTHENTICATION_FAILED"))
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			return
		}
		// log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		// err = errors.New("Admin Authentication Failed")
		// return
	}

	err = UpdatePassword(uname.(string), new_password)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Admin Update Password Failed")
		return
	}

	return
}

func Authenticate(uname, pass, userType, lang string) (err error) {

	row, err := db.Db.Query("SELECT id,password,status FROM Users where email=? limit 1", uname)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("System User Not Found")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("System User Detail Scan Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("System User Not Found")
		return
	}

	cp, err := base64.Decode([]byte(data[0][1]))
	if err != nil {
		err = errors.New("System User Password Decoding Fail")
		return
	}

	var pbkdf pbkdf2.Pbkdf2
	pbkdf.Itr = 32
	pbkdf.KeyLen = 32
	pbkdf.Plain = []byte(pass)
	pbkdf.Salt = cp[:32]
	pbkdf.Cipher = cp[32:]
	result, err := pbkdf.Compare()
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("System User Password MisMatch")
		return
	}

	if !result {
		log.Println(beego.AppConfig.String("loglevel"), "Error", result)
		err = errors.New("System User Password MisMatch")
		return
	}

	if data[0][2] != "ACTIVE" {
		if lang == "English" {
			err = errors.New("System User Not Active")
		} else {
			err = errors.New("Utilisateur systeme non trouvé")
		}
		return
	}

	return
}

func UpdatePassword(uname, password string) (err error) {

	b := make([]byte, 32)
	_, err = rand.Read(b)
	var pbkdf pbkdf2.Pbkdf2
	pbkdf.Itr = 32
	pbkdf.KeyLen = 32
	pbkdf.Plain = []byte(password)
	pbkdf.Salt = b
	err = pbkdf.Encrypt()
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("System User Password Update Fail")
		return
	}
	var tmp []byte
	tmp = append(tmp, pbkdf.Salt...)
	tmp = append(tmp, pbkdf.Cipher...)

	out, err := base64.Encode(tmp)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("System User Password Update Fail")
		return
	}

	result, err := db.Db.Exec("update Users set password=?, last_update=now() where email=? ", string(out), uname)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("System User Password Update Fail")
		return
	}
	n, err := result.RowsAffected()
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("System User Password Update Fail")
		return
	}

	if n != 1 {
		err = errors.New("System User Password Update Fail")
		return
	}
	//log.Println(beego.AppConfig.String("loglevel"), "Debug", "pass", password)
	return
}

func isValid(s string) bool {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	if len(s) >= 6 {
		hasMinLen = true
	}
	for _, char := range s {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}
