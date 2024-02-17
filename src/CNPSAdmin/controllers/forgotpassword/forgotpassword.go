package forgotPassword

import (
	"CNPSAdmin/model/db"
	"CNPSAdmin/session"

	"runtime/debug"

	"CNPSAdmin/model/utils"

	"github.com/astaxie/beego"

	"crypto/rand"

	"ominaya.com/database/sql"
	"ominaya.com/encoding/base64"
	"ominaya.com/util/pbkdf2"

	"errors"

	"io/ioutil"
	"net/mail"

	"net/smtp"

	"strings"

	"github.com/scorredoira/email"
	"ominaya.com/util/log"
	"ominaya.com/util/password"
)

type ForgotPassword struct {
	beego.Controller
}

type unencryptedAuth struct {
	smtp.Auth
}

func (c *ForgotPassword) Get() {
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Forgot Password Page Start")
	pip := c.Ctx.Input.IP()
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Client IP - ", pip)

	utils.SetHTTPHeader(c.Ctx)

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
			c.TplName = "forgotpassword/forgotpassword.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Forgot Password  Page Fail")
		} else {

			c.TplName = "forgotpassword/forgotpassword.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Forgot Password  Page Success")
		}
		return
	}()

	return
}

func (c *ForgotPassword) Post() {

	log.Println(beego.AppConfig.String("loglevel"), "Info", "Forgot Password Page Start")
	pip := c.Ctx.Input.IP()
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Client IP - ", pip)

	language := c.Input().Get("language")
	autmsg := ""
	pssSentSuccessMsg := ""
	sendMailError := ""
	if language == "English" {
		c.Data["OKLABEL"] = beego.AppConfig.String("ENGLISH_OK_BUTTON_LABEL")
		autmsg = "Admin Authentication Failed"
		pssSentSuccessMsg = "Password has been reset and sent on your registered email address successfully."
		sendMailError = beego.AppConfig.String("ENGLISH_USER_SENDMAIL__NOT_FOUND")
	} else {
		autmsg = "Authentification Admin a échoué"
		pssSentSuccessMsg = "Le mot de passe a été réinitialisé et envoyé a l'adresse mail enregistré."
		c.Data["OKLABEL"] = beego.AppConfig.String("FRENCH_OK_BUTTON_LABEL")
		sendMailError = beego.AppConfig.String("FRENCH_USER_SENDMAIL__NOT_FOUND")
	}
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
			c.TplName = "forgotpassword/forgotpassword.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Forgot Password Fail")
		} else {
			c.Data["DisplayMessage"] = pssSentSuccessMsg
			c.TplName = "forgotpassword/forgotpassword.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Forgot Password Success")
		}
		return
	}()

	utils.SetHTTPHeader(c.Ctx)
	uname := c.Input().Get("username")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Username - ", uname)

	if utils.IsDisableCharacters(uname) {
		beego.Error("suspicious special characters found")
		err = errors.New("Parameter validation failed")
		return
	}

	if uname == "" {
		log.Println(beego.AppConfig.String("loglevel"), "Error", "Blank User Name")
		err = errors.New("UserName can't be blank.")
		return
	}

	newPass, _ := password.AlphaNumericSpecial(6)

	UserFirstName, err := SearchUser(uname)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New(autmsg)
		return
	}

	err = UpdatePassword(uname, newPass)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Admin Update Password Failed")
		return
	}
	//log.Println(beego.AppConfig.String("loglevel"), "pass", newPass)

	err = SendEmail(uname, UserFirstName, newPass)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New(sendMailError)
		return
	}

	return
}

func SearchUser(uname string) (firstName string, err error) {

	row, err := db.Db.Query("SELECT id,email,first_name,status FROM Users where email=? limit 1", uname)
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

	if data[0][0] == beego.AppConfig.String("CNPS_DEFAULT_USER_ID") || data[0][0] == beego.AppConfig.String("CNPS_DEFAULT_MERCHANT_ID") {
		err = errors.New("Default User cannot be Update")
	}

	if data[0][3] != "ACTIVE" {
		err = errors.New("User Not Active")
		return
	}
	firstName = data[0][2]

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
		err = errors.New(" User Password Update Fail")
		return
	}
	var tmp []byte
	tmp = append(tmp, pbkdf.Salt...)
	tmp = append(tmp, pbkdf.Cipher...)

	out, err := base64.Encode(tmp)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New(" User2 Password Update Fail")
		return
	}

	result, err := db.Db.Exec("update Users set password=?, last_update=now(),password_set=? where email=? ", string(out), "NO", uname)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New(" User3 Password Update Fail")
		return
	}
	n, err := result.RowsAffected()
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New(" User4 Password Update Fail")
		return
	}

	if n != 1 {
		err = errors.New("System User Password Update Fail")
		return
	}
	//log.Println(beego.AppConfig.String("loglevel"), "Debug", "pass", password)
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

	// auth := unencryptedAuth{
	// 	smtp.PlainAuth(
	// 		"",
	// 		uname,
	// 		pass,
	// 		url,
	// 	),
	// }

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
