package utils

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"net/smtp"
	"strings"
	"text/template"

	"encoding/json"
	"strconv"

	"regexp"

	"github.com/astaxie/beego"
	"github.com/SHIBANI22/official/CNPSAdmin/model/db"
	"github.com/SHIBANI22/official/CNPSAdmin/encoding/base64"
	"github.com/SHIBANI22/official/CNPSAdmin/util/log"
	"github.com/SHIBANI22/official/CNPSAdmin/util/password"
	"github.com/SHIBANI22/official/CNPSAdmin//util/pbkdf2"

	"net/url"

	"CNPSAdmin/model/db"

	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/session"
)

type MenuArrayJson struct {
	Identity string `json:"identity"`
	Status   string `json:"status"`
}

type MenuArrayJsons []MenuArrayJson

func GeneratePassword() (login_pass, encrypted_pwd string, err error) {
	login_pass, _ = p.Numeric(6)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Login Password", login_pass)

	b := make([]byte, 32)
	_, err = rand.Read(b)
	var pbkdf pbkdf2.Pbkdf2
	pbkdf.Itr = 32
	pbkdf.KeyLen = 32
	pbkdf.Plain = []byte(login_pass)
	pbkdf.Salt = b
	err = pbkdf.Encrypt()
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to create password")
		return
	}
	var tmp []byte
	tmp = append(tmp, pbkdf.Salt...)
	tmp = append(tmp, pbkdf.Cipher...)

	out, err := base64.Encode(tmp)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to create password")
		return
	}
	encrypted_pwd = string(out)
	return
}

func SendEmail(host string, port int, userName string, password string, to []string, subject string, message string) (err error) {
	parameters := struct {
		From    string
		To      string
		Subject string
		Message string
	}{
		userName,
		strings.Join([]string(to), ","),
		subject,
		message,
	}

	buffer := new(bytes.Buffer)

	template := template.Must(template.New("emailTemplate").Parse(emailScript()))
	template.Execute(buffer, &parameters)

	auth := smtp.PlainAuth("", userName, password, host)

	err = smtp.SendMail(
		fmt.Sprintf("%s:%d", host, port),
		auth,
		userName,
		to,
		buffer.Bytes())

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	}
	return err
}
func emailScript() (script string) {
	return `From: {{.From}}
To: {{.To}}
Subject: {{.Subject}}
MIME-version: 1.0
Content-Type: text/html; charset="UTF-8"

{{.Message}}`
}

func SetHTTPHeader(Ctx *context.Context) {
	Ctx.Output.Header("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate")
	Ctx.Output.Header("Pragma", "no-cache")
	Ctx.Output.Header("Expires", "0")
	Ctx.Output.Header("X-Content-Type-Options", "nosniff")
	Ctx.Output.Header("Strict-Transport-Security", "max-age=31536000 ; includeSubDomains")
	Ctx.Output.Header("X-Frame-Options", "deny")
	Ctx.Output.Header("X-XSS-Protection", "1; mode=block")
	Ctx.Output.Header("X-Content-Security-Policy", "default-src 'self'")
	//Ctx.Output.Header("X-WebKit-CSP", "default-src 'self'")
}

func EventLogs(c *context.Context, sess session.Store, method string, input url.Values, output map[interface{}]interface{}, err_r error) (res string, err error) {
	m2 := make(map[string]interface{})

	for key, value := range output {
		switch key := key.(type) {
		case string:
			m2[key] = value
		}
	}
	event := make(map[string]interface{})

	event["PIP"] = c.Input.IP()
	event["URL"] = c.Input.URL()
	event["SessionID"] = sess.SessionID()
	if sess.Get("uname") != nil {
		event["UserName"] = sess.Get("uname")
	}
	event["Host"] = c.Input.Host()
	event["Method"] = method
	if err_r != nil {
		event["Status"] = err_r

	} else {
		event["Status"] = "Success"
	}

	event["Input"] = input
	event["Output"] = m2

	jsonString, err := json.Marshal(event)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Host Data - ", c.Input.Host())
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "URL Data - ", c.Input.URL())
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "User Data - ", sess.Get("uname"))
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "method Data - ", method)

	event_json := string(jsonString)

	_, err = db.Db.Exec("INSERT INTO web_event(event, created_on,user_id,url,ip,host,status) VALUES ( ?,now(),?,?,?,?,?)", event_json, sess.Get("uname"), c.Input.URL(), c.Input.IP(), c.Input.Host(), event["Status"])

	return
}

func UserStatus(uname string) (status string, err error) {

	row, err := db.Db.Query("select status from Users where email=?", uname)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("User Status Not Found")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("User Login Count Scan Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("User Login Count Not Found")
		return
	}
	status = data[0][0]
	//count, _ = strconv.Atoi(count_str)
	return
}

func LoginCount(uname string) (count int, err error) {

	row, err := db.Db.Query("select login_count from Users where email=?", uname)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("User Login Count Not Found")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("User Login Count Scan Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("User Login Count Not Found")
		return
	}
	count_str := data[0][0]
	count, _ = strconv.Atoi(count_str)
	return
}

func PasswordMismatch(count int, uname string) (err error) {
	result, err := db.Db.Exec("update Users set login_count=? where email=? ", count, uname)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New(" User Count Update Fail")
		return
	}
	n, err := result.RowsAffected()
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New(" User Count Update Fail")
		return
	}

	if n != 1 {
		err = errors.New(" User Count Update Fail")
		return
	}

	return
}

func Usercheck(uname string) (eamil string, err error) {

	row, err := db.Db.Query("select email from Users where email=?", uname)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("User Status Not Found")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("User Login Count Scan Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("User Login Count Not Found")
		return
	}
	eamil = data[0][0]
	//count, _ = strconv.Atoi(count_str)
	return
}

func UsercheckEmpid(uname string) (empid string, err error) {

	row, err := db.Db.Query("select employee_id from Users where employee_id=?", uname)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("User query data not Found")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("User query data not Found")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("User query data not Found")
		return
	}
	empid = data[0][0]
	//count, _ = strconv.Atoi(count_str)
	return
}

func UserLanguage(uname string) (status string, err error) {

	row, err := db.Db.Query("select language from Users where email=?", uname)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("User language Not Found")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("User Login Count Scan Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("User Login Count Not Found")
		return
	}
	status = data[0][0]
	//count, _ = strconv.Atoi(count_str)
	return
}

var IsNumber = regexp.MustCompile(`^[[:digit:][:space:].]+$|^$`).MatchString
var IsLetter = regexp.MustCompile(`^[a-zA-z[:space:]']+$|^$`).MatchString
var IsDisableCharacters = regexp.MustCompile(`[<|>|(|)|'|"|(|)|%|!|/|;]`).MatchString
var IsDateRange = regexp.MustCompile(`^[[:digit:][:space:]-/]+$`).MatchString

func IsAuthorized(menusjson, page string) (result bool, err error) {

	result = false
	var msg MenuArrayJsons
	err = json.Unmarshal([]byte(menusjson), &msg)
	if err != nil {
		return
	}

	for i := 0; i < len(msg); i++ {
		if msg[i].Identity == page {
			if msg[i].Status == "true" {
				result = true
			}
		}
	}
	return
}
func SideBarTransactionData() (successAmount, servicecharge, totalTransCount, totalBanks, successCount, pendingCount, declainedCount string, err error) {

	rows, errTrans := db.Db.Query("SELECT ROUND(SUM(CASE WHEN status='APPROVED' THEN amount ELSE 0 END),2), ROUND(SUM(CASE WHEN status='APPROVED' THEN service_charge ELSE 0 END),2), COUNT(CASE WHEN status='APPROVED' THEN 1 WHEN status='PENDING' THEN 1 WHEN status='DECLINED' THEN 1 ELSE NULL END),COUNT(CASE WHEN status='APPROVED' THEN 1 ELSE NULL END),COUNT(CASE WHEN status='PENDING' THEN 1 ELSE NULL END),COUNT(CASE WHEN status='DECLINED' THEN 1 ELSE NULL END) FROM Transactions")
	if errTrans != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", errTrans)
		err = errors.New("Transaction Data not found")
		return
	}

	_, data, err := sql.Scan(rows)
	defer sql.Close(rows)

	if err != nil {
		err = errors.New("Columns not found")
		return
	}
	dataArr := data[0]
	successAmount = dataArr[0]
	if successAmount == "" {
		successAmount = "0"
	}
	servicecharge = dataArr[1]
	if servicecharge == "" {
		servicecharge = "0"
	}
	totalTransCount = dataArr[2]
	successCount = dataArr[3]
	pendingCount = dataArr[4]
	declainedCount = dataArr[5]
	totalBanks = strconv.Itoa(len(strings.Split(beego.AppConfig.String("CNPS_OPERATORS"), "|")))
	return
}
func ResetLoginCount(uname string) (err error) {
	count := 0
	result, err := db.Db.Exec("UPDATE Users set login_count=? where email=? ", count, uname)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("System User Count Update Fail")
		return
	}
	n, err := result.RowsAffected()
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("System User Count Update Fail")
		return
	}

	if n != 1 {
		err = errors.New("System User Count Update Fail")
		return
	}

	return
}
func CurrencyFormat(m string) string {
	n, _ := strconv.ParseInt(m, 10, 64)
	in := strconv.FormatInt(n, 10)

	numOfDigits := len(in)
	// fmt.Printf("%d", numOfDigits)
	if n < 0 {
		numOfDigits-- // First character is the - sign (not a digit)
	}
	numOfCommas := (numOfDigits - 1) / 3

	out := make([]byte, len(in)+numOfCommas)
	if n < 0 {
		in, out[0] = in[1:], '-'
	}

	for i, j, k := len(in)-1, len(out)-1, 0; ; i, j = i-1, j-1 {
		out[j] = in[i]
		if i == 0 {
			return string(out)
		}
		if k++; k == 3 {
			j, k = j-1, 0
			out[j] = ' '
		}
	}
}
func CustomDateFormat(df, tf string) (string, error) {
	switch df + tf {
	case "YY-MM-DD" + "SS:MM:HH":
		return "2006-01-02 03:04:05PM", nil
	case "DD-MM-YY" + "HH:MM:SS":
		return "02-01-2006 03:04:05PM", nil
	case "DD/MM/YY" + "MM:HH:SS":
		return "02/01/2006 04:03:05PM", nil
	case "YY/MM/DD" + "HH:MM:SS":
		return "2006/01/02 15:04:05", nil
	case "YY/DD/MM" + "HH:SS:MM":
		return "2006/02/01 15:04:05", nil
	case "DD/MM/YYYY" + "":
		return "02/01/2006", nil
	case "DD-MM-YYYY" + "":
		return "02-01-2006", nil
	case "MM/DD/YYYY" + "":
		return "01/02/2006", nil
	case "YY/MM/DD" + "":
		return "2006/01/02", nil
	case "MM/DD/YY" + "":
		return "01/02/2006", nil
	case "YYYY-MM-DD" + "":
		return "2006-01-02", nil
	case "YYYY/MM/DD" + "":
		return "2006/01/02", nil
	case "YYYYMMDD" + "HHMMSS":
		return "20060102150405", nil
	case "YYYYMMDD" + "":
		return "20060102", nil
	case "DDMMYYYY" + "":
		return "02012006", nil
	case "MMDDYY" + "HHMMSS":
		return "010206150405", nil

	case "YYYY-MM-DD" + "HH:MM:SS":
		return "2006-01-02 15:04:05", nil

	case "DD.MM.YYYY" + "HH:MM:SS":
		return "02.01.2006 15:04:05", nil
	default:
		return "2006-01-02T15:04:05", nil
	}
}
