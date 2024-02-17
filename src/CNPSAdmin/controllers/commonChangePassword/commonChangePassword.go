package commonChangePassword

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
	//	"html/template"

	"runtime/debug"

	"github.com/astaxie/beego"
	"ominaya.com/util/log"
	//	"ominaya.com/util/password"
)

type CommonChangePassword struct {
	beego.Controller
}

func (c *CommonChangePassword) Get() {
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Common Change Password Page Start")
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
			c.TplName = "commonchangePassword/commonchangePassword.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Common Change Password Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "commonchangePassword/commonchangePassword.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Common Change Password Page Success")
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
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Client IP - ", sess.Get("language").(string))

	return
}

func (c *CommonChangePassword) Post() {

	log.Println(beego.AppConfig.String("loglevel"), "Info", "Common Change password Page Start")
	pip := c.Ctx.Input.IP()
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Client IP - ", pip)

	language := c.Input().Get("language")

	notMatchErr := ""
	blankErr := ""
	authErr := ""
	updateErr := ""
	passSuccessMsg := ""
	if language == "English" {
		c.Data["OKLABEL"] = beego.AppConfig.String("ENGLISH_OK_BUTTON_LABEL")
		blankErr = "Passwords can't be blank"
		notMatchErr = "Confirm Password can't be different"
		authErr = "Admin Authentication Failed"
		updateErr = "Admin Update Password Failed"
		passSuccessMsg = beego.AppConfig.String("ENGLISH_PASSWORD_CHANGE_SUCCCESS_LABEL")
	} else {
		c.Data["OKLABEL"] = beego.AppConfig.String("FRENCH_OK_BUTTON_LABEL")
		blankErr = "Les mots de passe ne peuvent pas être vides"
		notMatchErr = "Saisir le mot de Passe à confirmer, il doit être identique au nouveau mot de Passe"
		authErr = "Authentification Admin a échoué"
		updateErr = "Échec du mot de passe de mise à jour administrateur"
		passSuccessMsg = beego.AppConfig.String("FRENCH_PASSWORD_CHANGE_SUCCCESS_LABEL")
	}

	var err error
	sessErr := false
	defer func() {
		if l_exception := recover(); l_exception != nil {
			stack := debug.Stack()
			log.Println(beego.AppConfig.String("loglevel"), "Exception", string(stack))
			session.GlobalSessions.SessionDestroy(c.Ctx.ResponseWriter, c.Ctx.Request)
			c.TplName = "error/error.html"
		}
		if err != nil {
			if sessErr == true {
				log.Println(beego.AppConfig.String("loglevel"), "Info", "Redirecting to login")
				c.Redirect(beego.AppConfig.String("LOGIN_PATH"), 302)
			} else {
				c.Data["DisplayMessage"] = err.Error()
			}
			c.TplName = "commonchangePassword/commonchangePassword.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "change password Page Fail")
		} else {
			c.Data["DisplayMessage"] = passSuccessMsg
			//c.TplName = "commonchangePassword/commonchangePassword.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "change password Page Success")
			c.Redirect("/Login", 302)
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

	old_password := c.Input().Get("input_old_password")
	//log.Println(beego.AppConfig.String("loglevel"), "Debug", "old_password - ", old_password)

	new_password := c.Input().Get("input_new_password")
	//log.Println(beego.AppConfig.String("loglevel"), "Debug", "new_password - ", new_password)

	confirm_password := c.Input().Get("input_confirm_password")
	//log.Println(beego.AppConfig.String("loglevel"), "Debug", "confirm_password - ", confirm_password)

	uname := sess.Get("uname")
	sessLang := sess.Get("language").(string)

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "User Name - ", uname)

	if old_password == "" || new_password == "" || confirm_password == "" {
		log.Println(beego.AppConfig.String("loglevel"), "Error", "Blank Passwords")
		err = errors.New(blankErr)
		return
	}
	if len(old_password) > 16 || len(new_password) > 16 || len(confirm_password) > 16 {
		beego.Error("suspicious special characters found")
		err = errors.New(beego.AppConfig.String(sessLang + "_PARAMETER_VALIDATION"))
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
		err = errors.New(notMatchErr)
		return
	}

	err = Authenticate(uname.(string), old_password, "ADMIN", sess.Get("language").(string))

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New(authErr)
		return
	}

	err = UpdatePassword(uname.(string), new_password)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New(updateErr)
		return
	}

	c.Data["language"] = sess.Get("language").(string)

	//c.Redirect("/Login", 302)
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

	result, err := db.Db.Exec("update Users set password=?, last_update=now(),password_set=? where email=? ", string(out), "YES", uname)
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
